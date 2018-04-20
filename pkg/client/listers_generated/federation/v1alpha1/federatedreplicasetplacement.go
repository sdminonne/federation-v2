/*
Copyright 2018 The Federation v2 Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file was automatically generated by lister-gen

package v1alpha1

import (
	v1alpha1 "github.com/kubernetes-sigs/federation-v2/pkg/apis/federation/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// FederatedReplicaSetPlacementLister helps list FederatedReplicaSetPlacements.
type FederatedReplicaSetPlacementLister interface {
	// List lists all FederatedReplicaSetPlacements in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.FederatedReplicaSetPlacement, err error)
	// FederatedReplicaSetPlacements returns an object that can list and get FederatedReplicaSetPlacements.
	FederatedReplicaSetPlacements(namespace string) FederatedReplicaSetPlacementNamespaceLister
	FederatedReplicaSetPlacementListerExpansion
}

// federatedReplicaSetPlacementLister implements the FederatedReplicaSetPlacementLister interface.
type federatedReplicaSetPlacementLister struct {
	indexer cache.Indexer
}

// NewFederatedReplicaSetPlacementLister returns a new FederatedReplicaSetPlacementLister.
func NewFederatedReplicaSetPlacementLister(indexer cache.Indexer) FederatedReplicaSetPlacementLister {
	return &federatedReplicaSetPlacementLister{indexer: indexer}
}

// List lists all FederatedReplicaSetPlacements in the indexer.
func (s *federatedReplicaSetPlacementLister) List(selector labels.Selector) (ret []*v1alpha1.FederatedReplicaSetPlacement, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FederatedReplicaSetPlacement))
	})
	return ret, err
}

// FederatedReplicaSetPlacements returns an object that can list and get FederatedReplicaSetPlacements.
func (s *federatedReplicaSetPlacementLister) FederatedReplicaSetPlacements(namespace string) FederatedReplicaSetPlacementNamespaceLister {
	return federatedReplicaSetPlacementNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FederatedReplicaSetPlacementNamespaceLister helps list and get FederatedReplicaSetPlacements.
type FederatedReplicaSetPlacementNamespaceLister interface {
	// List lists all FederatedReplicaSetPlacements in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.FederatedReplicaSetPlacement, err error)
	// Get retrieves the FederatedReplicaSetPlacement from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.FederatedReplicaSetPlacement, error)
	FederatedReplicaSetPlacementNamespaceListerExpansion
}

// federatedReplicaSetPlacementNamespaceLister implements the FederatedReplicaSetPlacementNamespaceLister
// interface.
type federatedReplicaSetPlacementNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FederatedReplicaSetPlacements in the indexer for a given namespace.
func (s federatedReplicaSetPlacementNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.FederatedReplicaSetPlacement, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FederatedReplicaSetPlacement))
	})
	return ret, err
}

// Get retrieves the FederatedReplicaSetPlacement from the indexer for a given namespace and name.
func (s federatedReplicaSetPlacementNamespaceLister) Get(name string) (*v1alpha1.FederatedReplicaSetPlacement, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("federatedreplicasetplacement"), name)
	}
	return obj.(*v1alpha1.FederatedReplicaSetPlacement), nil
}
