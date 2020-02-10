// Licensed Materials - Property of IBM
// (c) Copyright IBM Corporation 2018. All Rights Reserved.
// Note to U.S. Government Users Restricted Rights:
// Use, duplication or disclosure restricted by GSA ADP Schedule
// Contract with IBM Corp.

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/clientset_generated/clientset/typed/mcm/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeMcmV1alpha1 struct {
	*testing.Fake
}

func (c *FakeMcmV1alpha1) ClusterJoinRequests() v1alpha1.ClusterJoinRequestInterface {
	return &FakeClusterJoinRequests{c}
}

func (c *FakeMcmV1alpha1) ClusterStatuses(namespace string) v1alpha1.ClusterStatusInterface {
	return &FakeClusterStatuses{c, namespace}
}

func (c *FakeMcmV1alpha1) PlacementBindings(namespace string) v1alpha1.PlacementBindingInterface {
	return &FakePlacementBindings{c, namespace}
}

func (c *FakeMcmV1alpha1) PlacementPolicies(namespace string) v1alpha1.PlacementPolicyInterface {
	return &FakePlacementPolicies{c, namespace}
}

func (c *FakeMcmV1alpha1) ResourceViews(namespace string) v1alpha1.ResourceViewInterface {
	return &FakeResourceViews{c, namespace}
}

func (c *FakeMcmV1alpha1) Works(namespace string) v1alpha1.WorkInterface {
	return &FakeWorks{c, namespace}
}

func (c *FakeMcmV1alpha1) WorkSets(namespace string) v1alpha1.WorkSetInterface {
	return &FakeWorkSets{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeMcmV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
