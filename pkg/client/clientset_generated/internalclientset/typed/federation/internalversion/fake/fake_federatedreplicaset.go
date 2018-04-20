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
package fake

import (
	federation "github.com/kubernetes-sigs/federation-v2/pkg/apis/federation"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFederatedReplicaSets implements FederatedReplicaSetInterface
type FakeFederatedReplicaSets struct {
	Fake *FakeFederation
	ns   string
}

var federatedreplicasetsResource = schema.GroupVersionResource{Group: "federation.k8s.io", Version: "", Resource: "federatedreplicasets"}

var federatedreplicasetsKind = schema.GroupVersionKind{Group: "federation.k8s.io", Version: "", Kind: "FederatedReplicaSet"}

// Get takes name of the federatedReplicaSet, and returns the corresponding federatedReplicaSet object, and an error if there is any.
func (c *FakeFederatedReplicaSets) Get(name string, options v1.GetOptions) (result *federation.FederatedReplicaSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(federatedreplicasetsResource, c.ns, name), &federation.FederatedReplicaSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*federation.FederatedReplicaSet), err
}

// List takes label and field selectors, and returns the list of FederatedReplicaSets that match those selectors.
func (c *FakeFederatedReplicaSets) List(opts v1.ListOptions) (result *federation.FederatedReplicaSetList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(federatedreplicasetsResource, federatedreplicasetsKind, c.ns, opts), &federation.FederatedReplicaSetList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &federation.FederatedReplicaSetList{}
	for _, item := range obj.(*federation.FederatedReplicaSetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested federatedReplicaSets.
func (c *FakeFederatedReplicaSets) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(federatedreplicasetsResource, c.ns, opts))

}

// Create takes the representation of a federatedReplicaSet and creates it.  Returns the server's representation of the federatedReplicaSet, and an error, if there is any.
func (c *FakeFederatedReplicaSets) Create(federatedReplicaSet *federation.FederatedReplicaSet) (result *federation.FederatedReplicaSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(federatedreplicasetsResource, c.ns, federatedReplicaSet), &federation.FederatedReplicaSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*federation.FederatedReplicaSet), err
}

// Update takes the representation of a federatedReplicaSet and updates it. Returns the server's representation of the federatedReplicaSet, and an error, if there is any.
func (c *FakeFederatedReplicaSets) Update(federatedReplicaSet *federation.FederatedReplicaSet) (result *federation.FederatedReplicaSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(federatedreplicasetsResource, c.ns, federatedReplicaSet), &federation.FederatedReplicaSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*federation.FederatedReplicaSet), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFederatedReplicaSets) UpdateStatus(federatedReplicaSet *federation.FederatedReplicaSet) (*federation.FederatedReplicaSet, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(federatedreplicasetsResource, "status", c.ns, federatedReplicaSet), &federation.FederatedReplicaSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*federation.FederatedReplicaSet), err
}

// Delete takes name of the federatedReplicaSet and deletes it. Returns an error if one occurs.
func (c *FakeFederatedReplicaSets) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(federatedreplicasetsResource, c.ns, name), &federation.FederatedReplicaSet{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFederatedReplicaSets) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(federatedreplicasetsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &federation.FederatedReplicaSetList{})
	return err
}

// Patch applies the patch and returns the patched federatedReplicaSet.
func (c *FakeFederatedReplicaSets) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *federation.FederatedReplicaSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(federatedreplicasetsResource, c.ns, name, data, subresources...), &federation.FederatedReplicaSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*federation.FederatedReplicaSet), err
}
