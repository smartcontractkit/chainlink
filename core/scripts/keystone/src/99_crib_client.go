// A k8s client specifically for the crib project. Contains methods that are hard coded
// to the crib project's needs.
package src

import (
	"fmt"
	"sort"
	"strings"
)

type CribClient struct {
	k8sClient *K8sClient
}

// SimpleURL lets us marshal a URL with only the fields we need.
type SimpleURL struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Path   string `json:"path"`
}

func (s SimpleURL) String() string {
	return fmt.Sprintf("%s://%s%s", s.Scheme, s.Host, s.Path)
}

type NodeWithCreds struct {
	URL              SimpleURL
	RemoteURL        SimpleURL
	ServiceName      string
	APILogin         string
	APIPassword      string
	KeystorePassword string
}

func NewCribClient() *CribClient {
	k8sClient := MustNewK8sClient()
	return &CribClient{
		k8sClient: k8sClient,
	}
}

func (m *CribClient) getCLNodes() ([]NodeWithCreds, error) {
	fmt.Println("Getting CL node deployments with config maps...")
	deployments, err := m.k8sClient.GetDeploymentsWithConfigMap()
	if err != nil {
		return nil, err
	}
	nodes := []NodeWithCreds{}

	for _, deployment := range deployments {
		apiCredentials := deployment.ConfigMap.Data["apicredentials"]
		splitCreds := strings.Split(strings.TrimSpace(apiCredentials), "\n")
		username := splitCreds[0]
		password := splitCreds[1]
		keystorePassword := deployment.ConfigMap.Data["node-password"]
		url := SimpleURL{
			Scheme: "https",
			Host:   deployment.Host,
			Path:   "",
		}

		node := NodeWithCreds{
			// We dont handle both in-cluster and out-of-cluster deployments
			// Hence why both URL and RemoteURL are the same
			URL:              url,
			RemoteURL:        url,
			ServiceName:      deployment.ServiceName,
			APILogin:         username,
			APIPassword:      password,
			KeystorePassword: keystorePassword,
		}

		nodes = append(nodes, node)
	}

	// Sort nodes by URL
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].URL.Host < nodes[j].URL.Host
	})

	return nodes, nil
}
