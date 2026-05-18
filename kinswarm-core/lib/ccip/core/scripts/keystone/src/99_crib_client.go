// A k8s client specifically for the crib project. Contains methods that are hard coded
// to the crib project's needs.
package src

import (
	"fmt"
	"net/url"
	"strings"
)

type CribClient struct {
	k8sClient *K8sClient
}

type CLNodeCredentials struct {
	URL          *url.URL
	PodName      string
	Username     string
	Password     string
	NodePassword string
}

func NewCribClient() *CribClient {
	k8sClient := MustNewK8sClient()
	return &CribClient{
		k8sClient: k8sClient,
	}
}

func (m *CribClient) GetCLNodeCredentials() ([]CLNodeCredentials, error) {
	fmt.Println("Getting CL node pods with config maps...")
	pods, err := m.k8sClient.GetPodsWithConfigMap()
	if err != nil {
		return nil, err
	}
	clNodeCredentials := []CLNodeCredentials{}

	for _, pod := range pods {
		apiCredentials := pod.ConfigMap.Data["apicredentials"]
		splitCreds := strings.Split(strings.TrimSpace(apiCredentials), "\n")
		username := splitCreds[0]
		password := splitCreds[1]
		nodePassword := pod.ConfigMap.Data["node-password"]
		url, err := url.Parse("https://" + pod.Host)
		if err != nil {
			return nil, err
		}

		clNodeCredential := CLNodeCredentials{
			URL:          url,
			PodName:      pod.Name,
			Username:     username,
			Password:     password,
			NodePassword: nodePassword,
		}

		clNodeCredentials = append(clNodeCredentials, clNodeCredential)
	}

	return clNodeCredentials, nil
}
