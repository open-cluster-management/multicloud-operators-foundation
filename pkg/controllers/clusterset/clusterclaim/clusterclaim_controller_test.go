package clusterclaim

import (
	"context"
	"os"
	"testing"

	clusterv1alapha1 "github.com/open-cluster-management/api/cluster/v1alpha1"
	"github.com/open-cluster-management/multicloud-operators-foundation/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
)

var (
	scheme = runtime.NewScheme()
)

func TestMain(m *testing.M) {
	// AddToSchemes may be used to add all resources defined in the project to a Scheme
	var AddToSchemes runtime.SchemeBuilder
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back

	if err := AddToSchemes.AddToScheme(scheme); err != nil {
		klog.Errorf("Failed adding apis to scheme, %v", err)
		os.Exit(1)
	}
	if err := clusterv1alapha1.Install(scheme); err != nil {
		klog.Errorf("Failed adding cluster v1alph1 to scheme, %v", err)
		os.Exit(1)
	}

	if err := hivev1.AddToScheme(scheme); err != nil {
		klog.Errorf("Failed adding hive to scheme, %v", err)
		os.Exit(1)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func newTestReconciler(clusterpoolObjs []runtime.Object, clusterclaimObjs []runtime.Object) *Reconciler {
	objs := append(clusterclaimObjs, clusterpoolObjs...)
	r := &Reconciler{
		client: fake.NewFakeClientWithScheme(scheme, objs...),
		scheme: scheme,
	}
	return r
}

func TestReconcile(t *testing.T) {

	tests := []struct {
		name             string
		clusterClaimObjs []runtime.Object
		clusterPools     []runtime.Object
		expectedlabel    map[string]string
		req              reconcile.Request
	}{
		{
			name: "add Clusterclaim",
			clusterClaimObjs: []runtime.Object{
				&hivev1.ClusterClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "clusterClaim1",
						Namespace: "ns1",
					},
					Spec: hivev1.ClusterClaimSpec{
						ClusterPoolName: "pool1",
					},
				},
			},
			clusterPools: []runtime.Object{
				&hivev1.ClusterPool{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pool1",
						Namespace: "ns1",
						Labels: map[string]string{
							utils.ClusterSetLabel: "clusterSet1",
						},
					},
					Spec: hivev1.ClusterPoolSpec{},
				},
			},
			req: reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "clusterClaim1",
					Namespace: "ns1",
				},
			},
			expectedlabel: map[string]string{
				utils.ClusterSetLabel: "clusterSet1",
			},
		},
		{
			name: "clusterclaim related clusterpool has no set label",
			clusterClaimObjs: []runtime.Object{
				&hivev1.ClusterClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "clusterClaim1",
						Namespace: "ns1",
						Labels: map[string]string{
							utils.ClusterSetLabel: "clusterSet1",
						},
					},
					Spec: hivev1.ClusterClaimSpec{
						ClusterPoolName: "pool1",
					},
				},
			},
			clusterPools: []runtime.Object{
				&hivev1.ClusterPool{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pool1",
						Namespace: "ns1",
					},
					Spec: hivev1.ClusterPoolSpec{},
				},
			},
			req: reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "clusterClaim1",
					Namespace: "ns1",
				},
			},
			expectedlabel: map[string]string{},
		},
	}

	for _, test := range tests {
		r := newTestReconciler(test.clusterPools, test.clusterClaimObjs)
		r.Reconcile(test.req)
		validateResult(t, r, test.name, test.req, test.expectedlabel)
	}
}

func validateResult(t *testing.T, r *Reconciler, caseName string, req reconcile.Request, expectedlabels map[string]string) {
	clusterclaim := &hivev1.ClusterClaim{}

	err := r.client.Get(context.TODO(), req.NamespacedName, clusterclaim)
	if err != nil {
		t.Errorf("case: %v, failed to get clusterclaim: %v", caseName, req.NamespacedName)
	}

	if len(clusterclaim.Labels) != len(expectedlabels) {
		t.Errorf("case: %v, expect:%v  actual:%v", caseName, expectedlabels, clusterclaim.Labels)
	}
}
