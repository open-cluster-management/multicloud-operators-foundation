package main

import (
	"context"
	"os"

	"github.com/open-cluster-management/multicloud-operators-foundation/cmd/acm-agent/app"
	"github.com/open-cluster-management/multicloud-operators-foundation/cmd/acm-agent/app/options"
	actionv1beta1 "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/action/v1beta1"
	clusterv1beta1 "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/cluster/v1beta1"
	viewv1beta1 "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/view/v1beta1"
	actionctrl "github.com/open-cluster-management/multicloud-operators-foundation/pkg/klusterlet/action"
	clusterinfoctl "github.com/open-cluster-management/multicloud-operators-foundation/pkg/klusterlet/clusterinfo"
	viewctrl "github.com/open-cluster-management/multicloud-operators-foundation/pkg/klusterlet/view"
	restutils "github.com/open-cluster-management/multicloud-operators-foundation/pkg/utils/rest"
	routev1 "github.com/openshift/client-go/route/clientset/versioned"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Needed for misc auth.
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = actionv1beta1.AddToScheme(scheme)
	_ = viewv1beta1.AddToScheme(scheme)
	_ = clusterv1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	o := options.NewAgentOptions()
	o.AddFlags(pflag.CommandLine)

	stopCh := signals.SetupSignalHandler()
	startManager(o, stopCh)
}

func startManager(o *options.AgentOptions, stopCh <-chan struct{}) {
	hubConfig, err := clientcmd.BuildConfigFromFlags("", o.HubKubeConfig)
	if err != nil {
		setupLog.Error(err, "Unable to get hub kube config.")
		os.Exit(1)
	}
	managedClusterConfig, err := clientcmd.BuildConfigFromFlags("", o.KubeConfig)
	if err != nil {
		setupLog.Error(err, "Unable to get managed cluster kube config.")
		os.Exit(1)
	}
	managedClusterDynamicClient, err := dynamic.NewForConfig(managedClusterConfig)
	if err != nil {
		setupLog.Error(err, "Unable to create managed cluster dynamic client.")
		os.Exit(1)
	}
	managedClusterKubeClient, err := kubernetes.NewForConfig(managedClusterConfig)
	if err != nil {
		setupLog.Error(err, "Unable to create managed cluster kube client.")
		os.Exit(1)
	}
	routeV1Client, err := routev1.NewForConfig(managedClusterConfig)
	if err != nil {
		setupLog.Error(err, "New route client config error:")
	}

	managedClusterClient, err := kubernetes.NewForConfig(managedClusterConfig)
	if err != nil {
		setupLog.Error(err, "Unable to create managed cluster clientset.")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(hubConfig, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: o.MetricsAddr,
		Namespace:          o.ClusterName,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	go app.ServeHealthProbes(stopCh)

	run := func(ctx context.Context) {
		// run agent server
		agent, err := app.AgentServerRun(o, managedClusterClient)
		if err != nil {
			setupLog.Error(err, "unable to run agent server")
			os.Exit(1)
		}

		// run mapper
		discoveryClient := cacheddiscovery.NewMemCacheClient(managedClusterClient.Discovery())
		mapper := restutils.NewMapper(discoveryClient, stopCh)
		mapper.Run()

		// Add controller into manager
		actionReconciler := actionctrl.NewActionReconciler(
			mgr.GetClient(),
			ctrl.Log.WithName("controllers").WithName("ManagedClusterAction"),
			mgr.GetScheme(),
			managedClusterDynamicClient,
			restutils.NewKubeControl(mapper, managedClusterConfig),
			o.EnableImpersonation,
		)
		viewReconciler := &viewctrl.ViewReconciler{
			Client:                      mgr.GetClient(),
			Log:                         ctrl.Log.WithName("controllers").WithName("ManagedClusterView"),
			Scheme:                      mgr.GetScheme(),
			ManagedClusterDynamicClient: managedClusterDynamicClient,
			Mapper:                      mapper,
		}

		clusterInfoReconciler := clusterinfoctl.ClusterInfoReconciler{
			Client:                      mgr.GetClient(),
			Log:                         ctrl.Log.WithName("controllers").WithName("ManagedClusterInfo"),
			Scheme:                      mgr.GetScheme(),
			KubeClient:                  managedClusterKubeClient,
			ManagedClusterDynamicClient: managedClusterDynamicClient,
			AgentRoute:                  o.AgentRoute,
			AgentAddress:                o.AgentAddress,
			AgentIngress:                o.AgentIngress,
			AgentPort:                   int32(o.AgentPort),
			RouteV1Client:               routeV1Client,
			Agent:                       agent,
		}

		if err = actionReconciler.SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ManagedClusterAction")
			os.Exit(1)
		}

		if err = viewReconciler.SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ManagedClusterView")
			os.Exit(1)
		}

		if err = clusterInfoReconciler.SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ManagedClusterInfo")
			os.Exit(1)
		}

		setupLog.Info("starting manager")
		if err := mgr.Start(stopCh); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}
	}

	if !o.EnableLeaderElection {
		run(context.TODO())
		panic("unreachable")
	}

	lec, err := app.NewLeaderElection(scheme, managedClusterClient, run)
	if err != nil {
		setupLog.Error(err, "cannot create leader election")
		os.Exit(1)
	}

	leaderelection.RunOrDie(context.TODO(), *lec)
	panic("unreachable")
}
