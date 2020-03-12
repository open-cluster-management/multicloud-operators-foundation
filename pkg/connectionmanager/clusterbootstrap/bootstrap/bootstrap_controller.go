// licensed Materials - Property of IBM
// 5737-E67
// (C) Copyright IBM Corporation 2016, 2019 All Rights Reserved
// US Government Users Restricted Rights - Use, duplication or disclosure restricted by GSA ADP Schedule Contract with IBM Corp.

package bootstrap

import (
	"fmt"
	"reflect"
	"time"

	"github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/mcm"
	hcmv1alpha1 "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/mcm/v1alpha1"
	clientset "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/clientset_generated/clientset"
	clusterinformers "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/cluster_informers_generated/externalversions"
	clusterlisters "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/cluster_listers_generated/clusterregistry/v1alpha1"
	informers "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/informers_generated/externalversions"
	listers "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/listers_generated/mcm/v1alpha1"
	"github.com/open-cluster-management/multicloud-operators-foundation/pkg/connectionmanager/clusterbootstrap/rbac"
	"github.com/open-cluster-management/multicloud-operators-foundation/pkg/connectionmanager/common"
	csrv1beta1 "k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	csrlister "k8s.io/client-go/listers/certificates/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	clusterv1alpha1 "k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
	"k8s.io/klog"
)

// Controller control the cluster bootstrap process
type Controller struct {
	// hcmclientset is a clientset for our own API group
	hcmclientset                  clientset.Interface
	kubeclientset                 kubernetes.Interface
	clusterLister                 clusterlisters.ClusterLister
	clusterSyced                  cache.InformerSynced
	csrLister                     csrlister.CertificateSigningRequestLister
	csrSynced                     cache.InformerSynced
	hcmjoinLister                 listers.ClusterJoinRequestLister
	hcmSynced                     cache.InformerSynced
	autoApproveClusterJoinRequest bool
	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	csrworkqueue     workqueue.RateLimitingInterface
	hcmjoinworkqueue workqueue.RateLimitingInterface

	stopCh <-chan struct{}
}

// controllerKind contains the schema.GroupVersionKind for this controller type.
var hcmjoinControllerKind = mcm.SchemeGroupVersion.WithKind("ClusterJoinRequest")

type queueHandlerFunc func(key string) error

var clientCertUsage = []csrv1beta1.KeyUsage{
	csrv1beta1.UsageDigitalSignature,
	csrv1beta1.UsageKeyEncipherment,
	csrv1beta1.UsageClientAuth,
}

// NewController create a bootstrapcontroller object
func NewController(
	kubeclientset kubernetes.Interface,
	hcmclientset clientset.Interface,
	kubeInformerFactory kubeinformers.SharedInformerFactory,
	informerFactory informers.SharedInformerFactory,
	clusterInformerFactory clusterinformers.SharedInformerFactory,
	autoApproveClusterJoinRequest bool,
	stopCh <-chan struct{}) *Controller {
	csrInformers := kubeInformerFactory.Certificates().V1beta1().CertificateSigningRequests()
	clusterInformer := clusterInformerFactory.Clusterregistry().V1alpha1().Clusters()
	bootstrapInformers := informerFactory.Mcm().V1alpha1().ClusterJoinRequests()

	controller := &Controller{
		hcmclientset:                  hcmclientset,
		kubeclientset:                 kubeclientset,
		clusterLister:                 clusterInformer.Lister(),
		clusterSyced:                  clusterInformer.Informer().HasSynced,
		csrLister:                     csrInformers.Lister(),
		csrSynced:                     csrInformers.Informer().HasSynced,
		hcmjoinLister:                 bootstrapInformers.Lister(),
		hcmSynced:                     bootstrapInformers.Informer().HasSynced,
		autoApproveClusterJoinRequest: autoApproveClusterJoinRequest,
		csrworkqueue:                  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "CSR"),
		hcmjoinworkqueue:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "HCMJoin"),
		stopCh:                        stopCh,
	}

	csrInformers.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			csr := new.(*csrv1beta1.CertificateSigningRequest)
			controllerRef := metav1.GetControllerOf(csr)
			if controllerRef == nil {
				return
			}
			if controllerRef.Kind != hcmjoinControllerKind.Kind {
				return
			}
			controller.enqueue(new, controller.csrworkqueue)
		},
		UpdateFunc: func(old, new interface{}) {
			oldcsr := old.(*csrv1beta1.CertificateSigningRequest)
			newcsr := new.(*csrv1beta1.CertificateSigningRequest)
			controllerRef := metav1.GetControllerOf(newcsr)
			if controllerRef == nil {
				return
			}
			if controllerRef.Kind != hcmjoinControllerKind.Kind {
				return
			}
			if !reflect.DeepEqual(&oldcsr.Status, &newcsr.Status) {
				controller.enqueue(new, controller.csrworkqueue)
			}
		},
	})

	bootstrapInformers.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			controller.enqueue(new, controller.hcmjoinworkqueue)
		},
		UpdateFunc: func(old, new interface{}) {
			controller.enqueue(new, controller.hcmjoinworkqueue)
		},
	})

	clusterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			cluster := obj.(*clusterv1alpha1.Cluster)
			controller.cleanCluster(cluster)
		},
	})

	return controller
}

// Run is the main run loop of kluster server
func (bc *Controller) Run() {
	defer runtime.HandleCrash()
	defer bc.csrworkqueue.ShutDown()
	defer bc.hcmjoinworkqueue.ShutDown()

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for hcm informer caches to sync")
	if ok := cache.WaitForCacheSync(bc.stopCh, bc.csrSynced, bc.hcmSynced, bc.clusterSyced); !ok {
		klog.Errorf("failed to wait for kubernetes caches to sync")
		return
	}

	go wait.Until(bc.runCSRWorker, time.Second, bc.stopCh)
	go wait.Until(bc.runHCMJoinWorker, time.Second, bc.stopCh)

	<-bc.stopCh
	klog.Info("Shutting controller")
}

// runCSRWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (bc *Controller) runCSRWorker() {
	for bc.processNextWorkItem(bc.csrworkqueue, bc.csrHandler) {
	}
}

// runHCMJoinWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (bc *Controller) runHCMJoinWorker() {
	for bc.processNextWorkItem(bc.hcmjoinworkqueue, bc.hcmJoinHandler) {
	}
}

func (bc *Controller) processNextWorkItem(queue workqueue.RateLimitingInterface, fn queueHandlerFunc) bool {
	obj, shutdown := queue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer queue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			queue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := fn(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		queue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (bc *Controller) hcmJoinHandler(key string) error {
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	hcmjoin, err := bc.hcmjoinLister.Get(name)
	if err != nil {
		// The csr resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("work '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	if len(hcmjoin.ObjectMeta.Finalizers) > 0 {
		policy := metav1.DeletePropagationBackground
		err = bc.hcmclientset.McmV1alpha1().ClusterJoinRequests().Delete(hcmjoin.Name, &metav1.DeleteOptions{PropagationPolicy: &policy})
		if err != nil {
			return err
		}
	}

	csr, err := bc.csrLister.Get(hcmjoin.Name)
	var createErr error
	if err != nil {
		if errors.IsNotFound(err) {
			// Create csr here
			csr = &csrv1beta1.CertificateSigningRequest{
				ObjectMeta: metav1.ObjectMeta{
					Name:            hcmjoin.Name,
					Labels:          hcmjoin.Labels,
					OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(hcmjoin, hcmjoinControllerKind)},
				},
				Spec: csrv1beta1.CertificateSigningRequestSpec{
					Request: hcmjoin.Spec.Request,
					Usages:  clientCertUsage,
				},
			}
			csr, createErr = bc.kubeclientset.CertificatesV1beta1().CertificateSigningRequests().Create(csr)

			if createErr != nil {
				return createErr
			}
		} else {
			return err
		}
	}

	return bc.updateHCMJoinStatus(hcmjoin, csr)
}

func (bc *Controller) updateHCMJoinStatus(hcmjoin *hcmv1alpha1.ClusterJoinRequest, csr *csrv1beta1.CertificateSigningRequest) error {
	var alreadyApproved bool
	var alreadyDenied bool
	for _, c := range csr.Status.Conditions {
		if c.Type == csrv1beta1.CertificateApproved {
			alreadyApproved = true
		} else if c.Type == csrv1beta1.CertificateDenied {
			alreadyDenied = true
		}
	}

	// If hcmjoin is denied, deny csr and return
	if hcmjoin.Status.Phase == hcmv1alpha1.JoinPhaseDenied {
		if !alreadyDenied {
			csr.Status.Conditions = append(csr.Status.Conditions, csrv1beta1.CertificateSigningRequestCondition{
				Type:           csrv1beta1.CertificateDenied,
				Reason:         "CertificateDenied",
				Message:        "This csr " + csr.Name + " was denied by Hub.",
				LastUpdateTime: metav1.Now(),
			})
			_, err := bc.kubeclientset.CertificatesV1beta1().CertificateSigningRequests().UpdateApproval(csr)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if hcmjoin.Status.Phase != hcmv1alpha1.JoinPhaseApproved {
		// Check if cluster is unique
		clusters, err := bc.clusterLister.List(labels.Everything())
		if err != nil {
			return err
		}
		return bc.approveOrDenyClusterJoinRequest(hcmjoin, csr, alreadyApproved, clusters)
	}

	return nil
}

func (bc *Controller) approveOrDenyClusterJoinRequest(
	hcmjoin *hcmv1alpha1.ClusterJoinRequest,
	csr *csrv1beta1.CertificateSigningRequest, alreadyApproved bool,
	clusters []*clusterv1alpha1.Cluster) error {
	var denied bool

	if isRenewalRequest(hcmjoin) {
		// deny the renewal request if the cluster does not exist
		found := false
		for _, cl := range clusters {
			if cl.Name == hcmjoin.Spec.ClusterName && cl.Namespace == hcmjoin.Spec.ClusterNamespace {
				// If cluster is pending or offline, do not approve or deny the request
				// and let admin to handle it manually
				if len(cl.Status.Conditions) == 0 {
					return nil
				}
				found = true
				break
			}
		}
		denied = !found
	} else {
		for _, cl := range clusters {
			// If cluster is pending or offline, always approve the csr
			if len(cl.Status.Conditions) == 0 {
				continue
			}

			//If the joinrequest's cluster name and namespace are same as the existing one. it's maybe a reinstall,
			//so admin need to approve it, do not deny automatically
			if cl.Name == hcmjoin.Spec.ClusterName && cl.Namespace == hcmjoin.Spec.ClusterNamespace {
				return nil
			}
			//Do not allow one namespace has more than one cluster, and do not allow diffirent cluster has same name.
			if cl.Name == hcmjoin.Spec.ClusterName {
				denied = true
				break
			}
			if cl.Namespace == hcmjoin.Spec.ClusterNamespace {
				denied = true
				break
			}
		}
	}

	if denied {
		// update hcmjoinrequest
		hcmjoin.Status.Phase = hcmv1alpha1.JoinPhaseDenied
		if _, err := bc.hcmclientset.McmV1alpha1().ClusterJoinRequests().UpdateStatus(hcmjoin); err != nil {
			return err
		}
	} else if bc.autoApproveClusterJoinRequest && !alreadyApproved {
		csr.Status.Conditions = append(csr.Status.Conditions, csrv1beta1.CertificateSigningRequestCondition{
			Type:           csrv1beta1.CertificateApproved,
			Reason:         "CertificateApproved",
			Message:        "This CSR " + csr.Name + " was approved by Hub.",
			LastUpdateTime: metav1.Now(),
		})
		if _, err := bc.kubeclientset.CertificatesV1beta1().CertificateSigningRequests().UpdateApproval(csr); err != nil {
			return err
		}
	}
	return nil
}

func (bc *Controller) createRoles(hcmjoin *hcmv1alpha1.ClusterJoinRequest) error {
	_, err := bc.kubeclientset.CoreV1().Namespaces().Get(hcmjoin.Spec.ClusterNamespace, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if errors.IsNotFound(err) {
		var createErr error
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: hcmjoin.Spec.ClusterNamespace,
			},
		}
		_, createErr = bc.kubeclientset.CoreV1().Namespaces().Create(ns)
		if createErr != nil {
			return createErr
		}
	}

	// Create Role/Rolebinding
	err = rbac.CreateOrUpdateRole(
		bc.kubeclientset,
		hcmjoin.Spec.ClusterName,
		hcmjoin.Spec.ClusterNamespace,
		*metav1.NewControllerRef(hcmjoin, hcmjoinControllerKind),
	)
	if err != nil {
		return err
	}

	err = rbac.CreateOrUpdateRoleBinding(
		bc.kubeclientset,
		hcmjoin.Spec.ClusterName,
		hcmjoin.Spec.ClusterNamespace,
		*metav1.NewControllerRef(hcmjoin, hcmjoinControllerKind),
	)

	return err
}

func (bc *Controller) csrHandler(key string) error {
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	csr, err := bc.csrLister.Get(name)
	if err != nil {
		// The csr resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("work '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	// Update hcmjoin if csr is approved
	var alreadyApproved bool
	var alreadyDenied bool
	for _, c := range csr.Status.Conditions {
		if c.Type == csrv1beta1.CertificateApproved {
			alreadyApproved = true
		} else if c.Type == csrv1beta1.CertificateDenied {
			alreadyDenied = true
		}
	}

	controllerRef := metav1.GetControllerOf(csr)
	if controllerRef == nil {
		return nil
	}
	hcmjoin := bc.resolveControllerRef(controllerRef)
	if hcmjoin == nil {
		return nil
	}

	if alreadyDenied && hcmjoin.Status.Phase != hcmv1alpha1.JoinPhaseDenied {
		hcmjoin.Status.Phase = hcmv1alpha1.JoinPhaseDenied
		condition := hcmv1alpha1.CLusterJoinRequestConditions{
			Type:           hcmv1alpha1.JoinTypeDenied,
			Reason:         "ClusterJoinRequestDenied",
			Message:        "This csr " + csr.Name + " was denied by Hub.",
			LastUpdateTime: metav1.Now(),
		}
		hcmjoin.Status.Conditions = append(hcmjoin.Status.Conditions, condition)
		hcmjoin.Status.Certificate = []byte{}
	} else if alreadyApproved && len(csr.Status.Certificate) != 0 {
		hcmjoin.Status.Phase = hcmv1alpha1.JoinPhaseApproved
		condition := hcmv1alpha1.CLusterJoinRequestConditions{
			Type:           hcmv1alpha1.JoinTypeApproved,
			Reason:         "ClusterJoinRequestApproved",
			Message:        "This csr " + csr.Name + " was approved by Hub.",
			LastUpdateTime: metav1.Now(),
		}
		hcmjoin.Status.Conditions = append(hcmjoin.Status.Conditions, condition)
		hcmjoin.Status.Certificate = csr.Status.Certificate

		err = bc.createRoles(hcmjoin)
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	_, err = bc.hcmclientset.McmV1alpha1().ClusterJoinRequests().UpdateStatus(hcmjoin)

	return err
}

// enqueue takes a resource and converts it into a name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Work.
func (bc *Controller) enqueue(obj interface{}, queue workqueue.RateLimitingInterface) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	queue.AddRateLimited(key)
}

// resolveControllerRef returns the controller referenced by a ControllerRef,
// or nil if the ControllerRef could not be resolved to a matching controller
// of the correct Kind.
func (bc *Controller) resolveControllerRef(controllerRef *metav1.OwnerReference) *hcmv1alpha1.ClusterJoinRequest {
	// We can't look up by UID, so look up by Name and then verify UID.
	// Don't even try to look up by Name if it's the wrong Kind.
	if controllerRef.Kind != hcmjoinControllerKind.Kind {
		return nil
	}
	hcmjoin, err := bc.hcmjoinLister.Get(controllerRef.Name)
	if err != nil {
		return nil
	}
	if hcmjoin.UID != controllerRef.UID {
		// The controller we found with this Name is not the same one that the
		// ControllerRef points to.
		return nil
	}
	return hcmjoin
}

func (bc *Controller) cleanCluster(cluster *clusterv1alpha1.Cluster) {
	if err := bc.hcmclientset.McmV1alpha1().ClusterJoinRequests().Delete(
		cluster.Namespace+"-"+cluster.Name, &metav1.DeleteOptions{}); err != nil {
		klog.Errorf("failed to delete cluster, %v", err)
	}
}

func isRenewalRequest(hcmjoin *hcmv1alpha1.ClusterJoinRequest) bool {
	if hcmjoin.Annotations == nil {
		return false
	}

	value, ok := hcmjoin.Annotations[common.RenewalAnnotation]
	if !ok {
		return false
	}

	if value != "true" {
		return false
	}

	return true
}
