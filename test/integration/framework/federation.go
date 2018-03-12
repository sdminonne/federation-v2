/*
Copyright 2018 The Kubernetes Authors.

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

package framework

import (
	"fmt"
	"time"

	fedv1a1 "github.com/marun/fnord/pkg/apis/federation/v1alpha1"
	"github.com/marun/fnord/pkg/controller/federatedcluster"
	"github.com/marun/fnord/pkg/controller/util"
	"github.com/marun/fnord/test/common"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	crv1a1 "k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
)

// TODO(marun) In fedv1 namespace cleanup required that a kube api
// fixture run a namespace controller to ensure cleanup on deletion.
// Will this be required?

const userAgent = "federation-framework"

// FederationFixture manages servers for kube, cluster registry and
// federation along with a set of member clusters.
type FederationFixture struct {
	stopChan chan struct{}

	KubeApi *KubernetesApiFixture
	CrApi   *ClusterRegistryApiFixture
	FedApi  *FederationApiFixture

	Clusters map[string]*KubernetesApiFixture
}

func SetUpFederationFixture(tl common.TestLogger, clusterCount int) *FederationFixture {
	if clusterCount < 1 {
		tl.Fatal("Cluster count must be greater than 0")
	}
	tl.Logf("Starting a federation of %d clusters...", clusterCount)
	f := &FederationFixture{}
	f.setUp(tl, clusterCount)
	return f
}

func (f *FederationFixture) setUp(tl common.TestLogger, clusterCount int) {
	defer TearDownOnPanic(tl, f)

	tl.Logf("In FederationFixture.setUp 1")

	f.CrApi = SetUpClusterRegistryApiFixture(tl)
	tl.Logf("In FederationFixture.setUp 2")
	f.FedApi = SetUpFederationApiFixture(tl)
	tl.Logf("In FederationFixture.setUp 3")

	f.Clusters = make(map[string]*KubernetesApiFixture)
	for i := 0; i < clusterCount; i++ {
		clusterName := f.AddMemberCluster(tl)
		tl.Logf("Added cluster %s to the federation", clusterName)
	}

	// TODO(marun) Consider running the cluster controller as soon as
	// the kube api is available to speed up setting cluster status.
	f.stopChan = make(chan struct{})
	monitorPeriod := 1 * time.Second
	tl.Logf("Starting cluster controller")
	federatedcluster.StartClusterController(f.FedApi.NewConfig(tl), f.KubeApi.NewConfig(tl), f.CrApi.NewConfig(tl), f.stopChan, monitorPeriod)

	tl.Log("Federation started.")
}

func (f *FederationFixture) TearDown(tl common.TestLogger) {
	// Stop the cluster controller first to avoid spurious connection
	// errors when the target urls become unavailable.
	if f.stopChan != nil {
		close(f.stopChan)
		f.stopChan = nil
	}
	fixtures := []TestFixture{
		// KubeApi will be torn down via f.Clusters
		f.CrApi,
		f.FedApi,
	}
	for _, cluster := range f.Clusters {
		fixtures = append(fixtures, cluster)
	}
	for _, fixture := range fixtures {
		fixture.TearDown(tl)
	}
	tl.Log("FederationFixture.TearDown")
}

// AddCluster adds a new member cluster to the federation.
func (f *FederationFixture) AddMemberCluster(tl common.TestLogger) string {
	tl.Logf("In FederationFixture.AddMemeberCluster!")

	kubeApi := SetUpKubernetesApiFixture(tl)

	tl.Logf("In FederationFixture.AddMemeberCluster 2")

	// Pick the first added cluster to be the primary
	if f.KubeApi == nil {
		f.KubeApi = kubeApi
	}

	clusterName := f.registerCluster(tl, kubeApi.Host)
	tl.Logf("In FederationFixture.AddMemeberCluster 30")
	secretName := f.createSecret(tl, kubeApi, clusterName)
	tl.Logf("In FederationFixture.AddMemeberCluster 31")
	f.createFederatedCluster(tl, clusterName, secretName)
	tl.Logf("In FederationFixture.AddMemeberCluster 32")

	// Track clusters by name
	f.Clusters[clusterName] = kubeApi

	tl.Logf("At end of FederationFixture.AddMemeberCluster!")

	return clusterName
}

// registerCluster registers a cluster with the cluster registry
func (f *FederationFixture) registerCluster(tl common.TestLogger, host string) string {
	// Registry the kube api with the cluster registry
	crClient := f.CrApi.NewClient(tl, userAgent)
	cluster, err := crClient.ClusterregistryV1alpha1().Clusters().Create(&crv1a1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-cluster-",
		},
		Spec: crv1a1.ClusterSpec{
			KubernetesAPIEndpoints: crv1a1.KubernetesAPIEndpoints{
				ServerEndpoints: []crv1a1.ServerAddressByClientCIDR{
					{
						ClientCIDR:    "0.0.0.0/0",
						ServerAddress: host,
					},
				},
			},
		},
	})
	if err != nil {
		tl.Fatal(err)
	}
	return cluster.Name
}

func (f *FederationFixture) ensureNamespace(tl common.TestLogger) error {
	kubeClient := f.KubeApi.NewClient(tl, userAgent)

	_, err := kubeClient.CoreV1().Namespaces().Create(&apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: util.FederationSystemNamespace,
		},
	})
	return err
}

// createSecret creates a secret resource containing the credentials
// necessary to access the fixture-managed cluster.
func (f *FederationFixture) createSecret(tl common.TestLogger, clusterFixture *KubernetesApiFixture, clusterName string) string {
	// Do not include the host - it will need to be sourced from the
	// Cluster resource.
	config := clusterFixture.SecureConfigFixture.NewClientConfig(tl, "")
	kubeConfig := CreateKubeConfig(config)
	tl.Log("FedarionFixture.createSecret 0")

	// Flatten the kubeconfig to ensure that all the referenced file
	// contents are inlined.
	err := clientcmdapi.FlattenConfig(kubeConfig)
	if err != nil {
		tl.Fatal(err)
	}
	tl.Log("FedarionFixture.createSecret 1")
	configBytes, err := clientcmd.Write(*kubeConfig)
	if err != nil {
		tl.Fatal(err)
	}
	tl.Log("FedarionFixture.createSecret 3")

	if err := f.ensureNamespace(tl); err != nil && !errors.IsAlreadyExists(err) {
		tl.Fatal(err)
	}

	// Build the secret object with the flattened kubeconfig content.
	// TODO(marun) enforce some kind of relationship between federated cluster and secret?
	kubeClient := f.KubeApi.NewClient(tl, userAgent)

	secret, err := kubeClient.CoreV1().Secrets(util.FederationSystemNamespace).Create(&apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-credentials", clusterName),
			Namespace:    util.FederationSystemNamespace,
		},
		Data: map[string][]byte{
			util.KubeconfigSecretDataKey: configBytes,
		},
	})
	if err != nil {
		tl.Fatal(err)
	}
	return secret.Name
}

// createFederatedCluster create a federated cluster resource that
// associates the cluster and secret.
func (f *FederationFixture) createFederatedCluster(tl common.TestLogger, clusterName, secretName string) {
	fedClient := f.FedApi.NewClient(tl, userAgent)
	_, err := fedClient.FederationV1alpha1().FederatedClusters().Create(&fedv1a1.FederatedCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterName,
		},
		Spec: fedv1a1.FederatedClusterSpec{
			ClusterRef: apiv1.LocalObjectReference{
				Name: clusterName,
			},
			SecretRef: &apiv1.LocalObjectReference{
				Name: secretName,
			},
		},
	})
	if err != nil {
		tl.Fatal(err)
	}
}

func (f *FederationFixture) ClusterClients(tl common.TestLogger, userAgent string) map[string]clientset.Interface {
	clientMap := make(map[string]clientset.Interface)
	for name, cluster := range f.Clusters {
		clientMap[name] = cluster.NewClient(tl, userAgent)
	}
	return clientMap
}

func (f *FederationFixture) ClusterNames() []string {
	clusterNames := []string{}
	for name, _ := range f.Clusters {
		clusterNames = append(clusterNames, name)
	}
	return clusterNames
}
