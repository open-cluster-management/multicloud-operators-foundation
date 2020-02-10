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

// FakeWorkSets implements WorkSetInterface
type FakeWorkSets struct {
	Fake *FakeMcm
	ns   string
}

var worksetsResource = schema.GroupVersionResource{Group: "mcm.ibm.com", Version: "", Resource: "worksets"}

var worksetsKind = schema.GroupVersionKind{Group: "mcm.ibm.com", Version: "", Kind: "WorkSet"}

// Get takes name of the workSet, and returns the corresponding workSet object, and an error if there is any.
func (c *FakeWorkSets) Get(name string, options v1.GetOptions) (result *mcm.WorkSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(worksetsResource, c.ns, name), &mcm.WorkSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.WorkSet), err
}

// List takes label and field selectors, and returns the list of WorkSets that match those selectors.
func (c *FakeWorkSets) List(opts v1.ListOptions) (result *mcm.WorkSetList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(worksetsResource, worksetsKind, c.ns, opts), &mcm.WorkSetList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &mcm.WorkSetList{ListMeta: obj.(*mcm.WorkSetList).ListMeta}
	for _, item := range obj.(*mcm.WorkSetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested workSets.
func (c *FakeWorkSets) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(worksetsResource, c.ns, opts))

}

// Create takes the representation of a workSet and creates it.  Returns the server's representation of the workSet, and an error, if there is any.
func (c *FakeWorkSets) Create(workSet *mcm.WorkSet) (result *mcm.WorkSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(worksetsResource, c.ns, workSet), &mcm.WorkSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.WorkSet), err
}

// Update takes the representation of a workSet and updates it. Returns the server's representation of the workSet, and an error, if there is any.
func (c *FakeWorkSets) Update(workSet *mcm.WorkSet) (result *mcm.WorkSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(worksetsResource, c.ns, workSet), &mcm.WorkSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.WorkSet), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeWorkSets) UpdateStatus(workSet *mcm.WorkSet) (*mcm.WorkSet, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(worksetsResource, "status", c.ns, workSet), &mcm.WorkSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.WorkSet), err
}

// Delete takes name of the workSet and deletes it. Returns an error if one occurs.
func (c *FakeWorkSets) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(worksetsResource, c.ns, name), &mcm.WorkSet{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeWorkSets) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(worksetsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &mcm.WorkSetList{})
	return err
}

// Patch applies the patch and returns the patched workSet.
func (c *FakeWorkSets) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *mcm.WorkSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(worksetsResource, c.ns, name, pt, data, subresources...), &mcm.WorkSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.WorkSet), err
}
