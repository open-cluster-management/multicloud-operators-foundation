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

// FakeWorks implements WorkInterface
type FakeWorks struct {
	Fake *FakeMcm
	ns   string
}

var worksResource = schema.GroupVersionResource{Group: "mcm.ibm.com", Version: "", Resource: "works"}

var worksKind = schema.GroupVersionKind{Group: "mcm.ibm.com", Version: "", Kind: "Work"}

// Get takes name of the work, and returns the corresponding work object, and an error if there is any.
func (c *FakeWorks) Get(name string, options v1.GetOptions) (result *mcm.Work, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(worksResource, c.ns, name), &mcm.Work{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.Work), err
}

// List takes label and field selectors, and returns the list of Works that match those selectors.
func (c *FakeWorks) List(opts v1.ListOptions) (result *mcm.WorkList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(worksResource, worksKind, c.ns, opts), &mcm.WorkList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &mcm.WorkList{ListMeta: obj.(*mcm.WorkList).ListMeta}
	for _, item := range obj.(*mcm.WorkList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested works.
func (c *FakeWorks) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(worksResource, c.ns, opts))

}

// Create takes the representation of a work and creates it.  Returns the server's representation of the work, and an error, if there is any.
func (c *FakeWorks) Create(work *mcm.Work) (result *mcm.Work, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(worksResource, c.ns, work), &mcm.Work{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.Work), err
}

// Update takes the representation of a work and updates it. Returns the server's representation of the work, and an error, if there is any.
func (c *FakeWorks) Update(work *mcm.Work) (result *mcm.Work, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(worksResource, c.ns, work), &mcm.Work{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.Work), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeWorks) UpdateStatus(work *mcm.Work) (*mcm.Work, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(worksResource, "status", c.ns, work), &mcm.Work{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.Work), err
}

// Delete takes name of the work and deletes it. Returns an error if one occurs.
func (c *FakeWorks) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(worksResource, c.ns, name), &mcm.Work{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeWorks) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(worksResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &mcm.WorkList{})
	return err
}

// Patch applies the patch and returns the patched work.
func (c *FakeWorks) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *mcm.Work, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(worksResource, c.ns, name, pt, data, subresources...), &mcm.Work{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mcm.Work), err
}
