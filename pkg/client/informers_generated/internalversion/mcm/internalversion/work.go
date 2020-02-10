// Licensed Materials - Property of IBM
// (c) Copyright IBM Corporation 2018. All Rights Reserved.
// Note to U.S. Government Users Restricted Rights:
// Use, duplication or disclosure restricted by GSA ADP Schedule
// Contract with IBM Corp.

// Code generated by informer-gen. DO NOT EDIT.

package internalversion

import (
	time "time"

	mcm "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/mcm"
	internalclientset "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/clientset_generated/internalclientset"
	internalinterfaces "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/informers_generated/internalversion/internalinterfaces"
	internalversion "github.com/open-cluster-management/multicloud-operators-foundation/pkg/client/listers_generated/mcm/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// WorkInformer provides access to a shared informer and lister for
// Works.
type WorkInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.WorkLister
}

type workInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewWorkInformer constructs a new informer for Work type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewWorkInformer(client internalclientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredWorkInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredWorkInformer constructs a new informer for Work type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredWorkInformer(client internalclientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Mcm().Works(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Mcm().Works(namespace).Watch(options)
			},
		},
		&mcm.Work{},
		resyncPeriod,
		indexers,
	)
}

func (f *workInformer) defaultInformer(client internalclientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredWorkInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *workInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&mcm.Work{}, f.defaultInformer)
}

func (f *workInformer) Lister() internalversion.WorkLister {
	return internalversion.NewWorkLister(f.Informer().GetIndexer())
}
