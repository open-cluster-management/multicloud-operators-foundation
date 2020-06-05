package clusterrbac

import (
	"os"
	"testing"
	"time"

	clusterv1 "github.com/open-cluster-management/api/cluster/v1"
	"github.com/stretchr/testify/assert"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	scheme = runtime.NewScheme()
)

func TestMain(m *testing.M) {
	// AddToSchemes may be used to add all resources defined in the project to a Scheme
	var AddToSchemes runtime.SchemeBuilder
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, clusterv1.Install)

	if err := AddToSchemes.AddToScheme(scheme); err != nil {
		klog.Errorf("Failed adding apis to scheme, %v", err)
		os.Exit(1)
	}

	if err := clusterv1.Install(scheme); err != nil {
		klog.Errorf("Failed adding cluster to scheme, %v", err)
		os.Exit(1)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

const (
	managedClusterName = "foo"
)

func newRoleObjs() []runtime.Object {
	return []runtime.Object{
		&rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namePrefix + managedClusterName,
				Namespace: managedClusterName,
			},
			Rules: nil,
		},
		&rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namePrefix + managedClusterName,
				Namespace: managedClusterName,
			},
			Subjects: nil,
			RoleRef:  rbacv1.RoleRef{},
		},
	}
}

func validateError(t *testing.T, err, expectedErrorType error) {
	if expectedErrorType != nil {
		assert.EqualError(t, err, expectedErrorType.Error())
	} else {
		assert.NoError(t, err)
	}
}

func newTestReconciler(existingObjs, existingRoleOjb []runtime.Object) *Reconciler {
	return &Reconciler{
		client:     fake.NewFakeClientWithScheme(scheme, existingObjs...),
		scheme:     scheme,
		kubeClient: k8sfake.NewSimpleClientset(existingRoleOjb...),
	}
}

func TestReconcile(t *testing.T) {
	tests := []struct {
		name              string
		existingObjs      []runtime.Object
		existingRoleOjbs  []runtime.Object
		expectedErrorType error
		req               reconcile.Request
	}{
		{
			name:         "ManagedClusterNotFound",
			existingObjs: []runtime.Object{},
			req: reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: managedClusterName,
				},
			},
		},
		{
			name: "ManagedClusterConditionFalse",
			existingObjs: []runtime.Object{
				&clusterv1.ManagedCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name: managedClusterName,
					},
					Spec: clusterv1.ManagedClusterSpec{},
					Status: clusterv1.ManagedClusterStatus{
						Conditions: []clusterv1.StatusCondition{
							{
								Type:   clusterv1.ManagedClusterConditionJoined,
								Status: v1beta1.ConditionFalse,
							},
						},
					},
				},
			},
			req: reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: managedClusterName,
				},
			},
		},
		{
			name: "ManagedClusterNoFinalizer",
			existingObjs: []runtime.Object{
				&clusterv1.ManagedCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name: managedClusterName,
					},
					Spec: clusterv1.ManagedClusterSpec{},
					Status: clusterv1.ManagedClusterStatus{
						Conditions: []clusterv1.StatusCondition{
							{
								Type:   clusterv1.ManagedClusterConditionJoined,
								Status: v1beta1.ConditionTrue,
							},
						},
					},
				},
			},
			existingRoleOjbs: newRoleObjs(),
			req: reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: managedClusterName,
				},
			},
		},
		{
			name: "ManagedClusterHasFinalizer",
			existingObjs: []runtime.Object{
				&clusterv1.ManagedCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name: managedClusterName,
						DeletionTimestamp: &metav1.Time{
							Time: time.Now(),
						},
						Finalizers: []string{
							clusterRBACFinalizerName,
						},
					},
					Spec: clusterv1.ManagedClusterSpec{},
					Status: clusterv1.ManagedClusterStatus{
						Conditions: []clusterv1.StatusCondition{
							{
								Type:   clusterv1.ManagedClusterConditionJoined,
								Status: v1beta1.ConditionFalse,
							},
						},
					},
				},
			},
			req: reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: managedClusterName,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svrc := newTestReconciler(test.existingObjs, test.existingRoleOjbs)
			res, err := svrc.ReconcileByManagedCluster(test.req)
			validateError(t, err, test.expectedErrorType)
			assert.Equal(t, res.Requeue, false)
		})
	}
}