// Licensed Materials - Property of IBM
// (c) Copyright IBM Corporation 2018. All Rights Reserved.
// Note to U.S. Government Users Restricted Rights:
// Use, duplication or disclosure restricted by GSA ADP Schedule
// Contract with IBM Corp.

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/mcm/v1alpha1"
	scheme "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/clientset_generated/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ClusterStatusesGetter has a method to return a ClusterStatusInterface.
// A group's client should implement this interface.
type ClusterStatusesGetter interface {
	ClusterStatuses(namespace string) ClusterStatusInterface
}

// ClusterStatusInterface has methods to work with ClusterStatus resources.
type ClusterStatusInterface interface {
	Create(*v1alpha1.ClusterStatus) (*v1alpha1.ClusterStatus, error)
	Update(*v1alpha1.ClusterStatus) (*v1alpha1.ClusterStatus, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.ClusterStatus, error)
	List(opts v1.ListOptions) (*v1alpha1.ClusterStatusList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ClusterStatus, err error)
	ClusterStatusExpansion
}

// clusterStatuses implements ClusterStatusInterface
type clusterStatuses struct {
	client rest.Interface
	ns     string
}

// newClusterStatuses returns a ClusterStatuses
func newClusterStatuses(c *McmV1alpha1Client, namespace string) *clusterStatuses {
	return &clusterStatuses{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the clusterStatus, and returns the corresponding clusterStatus object, and an error if there is any.
func (c *clusterStatuses) Get(name string, options v1.GetOptions) (result *v1alpha1.ClusterStatus, err error) {
	result = &v1alpha1.ClusterStatus{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("clusterstatuses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterStatuses that match those selectors.
func (c *clusterStatuses) List(opts v1.ListOptions) (result *v1alpha1.ClusterStatusList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ClusterStatusList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("clusterstatuses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterStatuses.
func (c *clusterStatuses) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("clusterstatuses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a clusterStatus and creates it.  Returns the server's representation of the clusterStatus, and an error, if there is any.
func (c *clusterStatuses) Create(clusterStatus *v1alpha1.ClusterStatus) (result *v1alpha1.ClusterStatus, err error) {
	result = &v1alpha1.ClusterStatus{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("clusterstatuses").
		Body(clusterStatus).
		Do().
		Into(result)
	return
}

// Update takes the representation of a clusterStatus and updates it. Returns the server's representation of the clusterStatus, and an error, if there is any.
func (c *clusterStatuses) Update(clusterStatus *v1alpha1.ClusterStatus) (result *v1alpha1.ClusterStatus, err error) {
	result = &v1alpha1.ClusterStatus{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("clusterstatuses").
		Name(clusterStatus.Name).
		Body(clusterStatus).
		Do().
		Into(result)
	return
}

// Delete takes name of the clusterStatus and deletes it. Returns an error if one occurs.
func (c *clusterStatuses) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("clusterstatuses").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterStatuses) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("clusterstatuses").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched clusterStatus.
func (c *clusterStatuses) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ClusterStatus, err error) {
	result = &v1alpha1.ClusterStatus{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("clusterstatuses").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
