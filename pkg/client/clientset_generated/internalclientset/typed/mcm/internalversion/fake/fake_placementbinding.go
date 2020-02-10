// Licensed Materials - Property of IBM
// (c) Copyright IBM Corporation 2018. All Rights Reserved.
// Note to U.S. Government Users Restricted Rights:
// Use, duplication or disclosure restricted by GSA ADP Schedule
// Contract with IBM Corp.

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	mcm "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/mcm"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePlacementBindings implements PlacementBindingInterface
type FakePlacementBindings struct {
	Fake *FakeMcm
	ns   string
}

var placementbindingsResource = schema.GroupVersionResource{Group: "mcm.ibm.com", Version: "", Resource: "placementbindings"}

var placementbindingsKind = schema.GroupVersionKind{Group: "mcm.ibm.com", Version: "", Kind: "PlacementBinding"}

// Get takes name of the placementBinding, and returns the corresponding placementBinding object, and an error if there is any.
func (c *FakePlacementBindings) Get(name string, options v1.GetOptions) (result *mcm.PlacementBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(placementbindingsResource, c.ns, name), &mcm.PlacementBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.PlacementBinding), err
}

// List takes label and field selectors, and returns the list of PlacementBindings that match those selectors.
func (c *FakePlacementBindings) List(opts v1.ListOptions) (result *mcm.PlacementBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(placementbindingsResource, placementbindingsKind, c.ns, opts), &mcm.PlacementBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &mcm.PlacementBindingList{ListMeta: obj.(*mcm.PlacementBindingList).ListMeta}
	for _, item := range obj.(*mcm.PlacementBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested placementBindings.
func (c *FakePlacementBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(placementbindingsResource, c.ns, opts))

}

// Create takes the representation of a placementBinding and creates it.  Returns the server's representation of the placementBinding, and an error, if there is any.
func (c *FakePlacementBindings) Create(placementBinding *mcm.PlacementBinding) (result *mcm.PlacementBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(placementbindingsResource, c.ns, placementBinding), &mcm.PlacementBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.PlacementBinding), err
}

// Update takes the representation of a placementBinding and updates it. Returns the server's representation of the placementBinding, and an error, if there is any.
func (c *FakePlacementBindings) Update(placementBinding *mcm.PlacementBinding) (result *mcm.PlacementBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(placementbindingsResource, c.ns, placementBinding), &mcm.PlacementBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.PlacementBinding), err
}

// Delete takes name of the placementBinding and deletes it. Returns an error if one occurs.
func (c *FakePlacementBindings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(placementbindingsResource, c.ns, name), &mcm.PlacementBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePlacementBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(placementbindingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &mcm.PlacementBindingList{})
	return err
}

// Patch applies the patch and returns the patched placementBinding.
func (c *FakePlacementBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *mcm.PlacementBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(placementbindingsResource, c.ns, name, pt, data, subresources...), &mcm.PlacementBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.PlacementBinding), err
}
