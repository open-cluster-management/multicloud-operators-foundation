// Licensed Materials - Property of IBM
// (c) Copyright IBM Corporation 2018. All Rights Reserved.
// Note to U.S. Government Users Restricted Rights:
// Use, duplication or disclosure restricted by GSA ADP Schedule
// Contract with IBM Corp.

// Code generated by client-gen. DO NOT EDIT.

package internalversion

import (
	"time"

	mcm "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/mcm"
	scheme "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/clientset_generated/internalclientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ClusterJoinRequestsGetter has a method to return a ClusterJoinRequestInterface.
// A group's client should implement this interface.
type ClusterJoinRequestsGetter interface {
	ClusterJoinRequests() ClusterJoinRequestInterface
}

// ClusterJoinRequestInterface has methods to work with ClusterJoinRequest resources.
type ClusterJoinRequestInterface interface {
	Create(*mcm.ClusterJoinRequest) (*mcm.ClusterJoinRequest, error)
	Update(*mcm.ClusterJoinRequest) (*mcm.ClusterJoinRequest, error)
	UpdateStatus(*mcm.ClusterJoinRequest) (*mcm.ClusterJoinRequest, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*mcm.ClusterJoinRequest, error)
	List(opts v1.ListOptions) (*mcm.ClusterJoinRequestList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *mcm.ClusterJoinRequest, err error)
	ClusterJoinRequestExpansion
}

// clusterJoinRequests implements ClusterJoinRequestInterface
type clusterJoinRequests struct {
	client rest.Interface
}

// newClusterJoinRequests returns a ClusterJoinRequests
func newClusterJoinRequests(c *McmClient) *clusterJoinRequests {
	return &clusterJoinRequests{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterJoinRequest, and returns the corresponding clusterJoinRequest object, and an error if there is any.
func (c *clusterJoinRequests) Get(name string, options v1.GetOptions) (result *mcm.ClusterJoinRequest, err error) {
	result = &mcm.ClusterJoinRequest{}
	err = c.client.Get().
		Resource("clusterjoinrequests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterJoinRequests that match those selectors.
func (c *clusterJoinRequests) List(opts v1.ListOptions) (result *mcm.ClusterJoinRequestList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &mcm.ClusterJoinRequestList{}
	err = c.client.Get().
		Resource("clusterjoinrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterJoinRequests.
func (c *clusterJoinRequests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("clusterjoinrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a clusterJoinRequest and creates it.  Returns the server's representation of the clusterJoinRequest, and an error, if there is any.
func (c *clusterJoinRequests) Create(clusterJoinRequest *mcm.ClusterJoinRequest) (result *mcm.ClusterJoinRequest, err error) {
	result = &mcm.ClusterJoinRequest{}
	err = c.client.Post().
		Resource("clusterjoinrequests").
		Body(clusterJoinRequest).
		Do().
		Into(result)
	return
}

// Update takes the representation of a clusterJoinRequest and updates it. Returns the server's representation of the clusterJoinRequest, and an error, if there is any.
func (c *clusterJoinRequests) Update(clusterJoinRequest *mcm.ClusterJoinRequest) (result *mcm.ClusterJoinRequest, err error) {
	result = &mcm.ClusterJoinRequest{}
	err = c.client.Put().
		Resource("clusterjoinrequests").
		Name(clusterJoinRequest.Name).
		Body(clusterJoinRequest).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *clusterJoinRequests) UpdateStatus(clusterJoinRequest *mcm.ClusterJoinRequest) (result *mcm.ClusterJoinRequest, err error) {
	result = &mcm.ClusterJoinRequest{}
	err = c.client.Put().
		Resource("clusterjoinrequests").
		Name(clusterJoinRequest.Name).
		SubResource("status").
		Body(clusterJoinRequest).
		Do().
		Into(result)
	return
}

// Delete takes name of the clusterJoinRequest and deletes it. Returns an error if one occurs.
func (c *clusterJoinRequests) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("clusterjoinrequests").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterJoinRequests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("clusterjoinrequests").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched clusterJoinRequest.
func (c *clusterJoinRequests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *mcm.ClusterJoinRequest, err error) {
	result = &mcm.ClusterJoinRequest{}
	err = c.client.Patch(pt).
		Resource("clusterjoinrequests").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}