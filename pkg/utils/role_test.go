package utils

import (
	"context"
	"reflect"
	"testing"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMergesubjects(t *testing.T) {
	type args struct {
		subjects    []rbacv1.Subject
		cursubjects []rbacv1.Subject
	}
	tests := []struct {
		name string
		args args
		want []rbacv1.Subject
	}{
		{"test1", args{subjects: []rbacv1.Subject{{Kind: "R1", APIGroup: "G1", Name: "N1"}}, cursubjects: []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}}}, []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}, {Kind: "R1", APIGroup: "G1", Name: "N1"}}},
		{"test2", args{cursubjects: []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}}}, []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}}},
		{"test3", args{subjects: []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}}}, []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Mergesubjects(tt.args.subjects, tt.args.cursubjects)
			if len(res) != len(tt.want) {
				t.Errorf("Mergesubjects() = %v, want %v", res, tt.want)
			}
		})
	}
}

func createPolicyRule(groups, verbs, res, resnames []string) *rbacv1.PolicyRule {
	return &rbacv1.PolicyRule{
		APIGroups:     groups,
		Verbs:         verbs,
		Resources:     res,
		ResourceNames: resnames,
	}
}

func TestEqualSubjects(t *testing.T) {
	type args struct {
		subjects1 []rbacv1.Subject
		subjects2 []rbacv1.Subject
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{subjects1: []rbacv1.Subject{{Kind: "R1", APIGroup: "G1", Name: "N1"}}, subjects2: []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}}}, false},
		{"test2", args{subjects1: []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}}}, false},
		{"test2", args{subjects2: []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}}}, false},
		{"test3", args{subjects1: []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}}, subjects2: []rbacv1.Subject{{Kind: "R2", APIGroup: "G2", Name: "N2"}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualSubjects(tt.args.subjects1, tt.args.subjects2); got != tt.want {
				t.Errorf("EqualSubjects() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createClusterrolebinding(name, roleName string, labels map[string]string, subjects []rbacv1.Subject) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     roleName,
		},
		Subjects: subjects,
	}
}

func TestApplyClusterRoleBinding(t *testing.T) {
	ctx := context.Background()
	var objs []runtime.Object
	var labels = make(map[string]string)
	rb1 := createClusterrolebinding("crb1", "r1", labels, []rbacv1.Subject{})
	rb2 := createClusterrolebinding("crb1", "r2", labels, []rbacv1.Subject{})

	objs = append(objs, rb1)
	client := fake.NewFakeClient(objs...)
	req := rb1
	err := ApplyClusterRoleBinding(ctx, client, req)
	if err != nil {
		t.Errorf("Error to apply clusterolebinding. Error:%v", err)
	}
	applied := verifyApply(ctx, client, req)
	if !applied {
		t.Errorf("Error to apply clusterolebinding.")
	}

	req = rb2
	err = ApplyClusterRoleBinding(ctx, client, req)
	if err != nil {
		t.Errorf("Error to apply clusterolebinding. Error:%v", err)
	}
	applied = verifyApply(ctx, client, req)
	if !applied {
		t.Errorf("Error to apply clusterolebinding.")
	}
}

func verifyApply(ctx context.Context, client client.Client, required *rbacv1.ClusterRoleBinding) bool {
	existing := &rbacv1.ClusterRoleBinding{}
	err := client.Get(ctx, types.NamespacedName{Name: required.Name}, existing)
	if err != nil {
		return false
	}
	if !reflect.DeepEqual(existing.RoleRef, required.RoleRef) {
		return false
	}
	if !EqualSubjects(existing.Subjects, required.Subjects) {
		return false
	}
	return true
}

func TestIsManagedClusterClusterrolebinding(t *testing.T) {
	type args struct {
		rolebindingName string
		role            string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{rolebindingName: "not:hanlde", role: "admin"}, false},
		{"test2", args{rolebindingName: "open-cluster-management:managedclusterset:admin:managedcluster:managedcluster1", role: "admin"}, true},
		{"test3", args{rolebindingName: "open-cluster-management:managedclusterset:view:managedcluster:managedcluster1", role: "false"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IsManagedClusterClusterrolebinding(tt.args.rolebindingName, tt.args.role)
			if res != tt.want {
				t.Errorf("Failed to test IsManagedClusterClusterrolebinding, rolebinding name: %v, role: %v, want: %v", tt.args.rolebindingName, tt.args.role, tt.want)
			}
		})
	}
}

func TestContainsSubject(t *testing.T) {
	type args struct {
		rolebindingName string
		role            string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{rolebindingName: "not:hanlde", role: "admin"}, false},
		{"test2", args{rolebindingName: "open-cluster-management:managedclusterset:admin:managedcluster:managedcluster1", role: "admin"}, true},
		{"test3", args{rolebindingName: "open-cluster-management:managedclusterset:view:managedcluster:managedcluster1", role: "false"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IsManagedClusterClusterrolebinding(tt.args.rolebindingName, tt.args.role)
			if res != tt.want {
				t.Errorf("Failed to test IsManagedClusterClusterrolebinding, rolebinding name: %v, role: %v, want: %v", tt.args.rolebindingName, tt.args.role, tt.want)
			}
		})
	}
}
