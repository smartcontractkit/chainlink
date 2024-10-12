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
	URL            *url.URL
	DeploymentName string
	ServiceName    string
	Username       string
	Password       string
	NodePassword   string
}

func NewCribClient() *CribClient {
	k8sClient := MustNewK8sClient()
	return &CribClient{
		k8sClient: k8sClient,
	}
}

func (m *CribClient) GetCLNodeCredentials() ([]CLNodeCredentials, error) {
	fmt.Println("Getting CL node deployments with config maps...")
	deployments, err := m.k8sClient.GetDeploymentsWithConfigMap()
	if err != nil {
		return nil, err
	}
	clNodeCredentials := []CLNodeCredentials{}

	for _, deployment := range deployments {
		apiCredentials := deployment.ConfigMap.Data["apicredentials"]
		splitCreds := strings.Split(strings.TrimSpace(apiCredentials), "\n")
		username := splitCreds[0]
		password := splitCreds[1]
		nodePassword := deployment.ConfigMap.Data["node-password"]
		url, err := url.Parse("https://" + deployment.Host)
		if err != nil {
			return nil, err
		}

		clNodeCredential := CLNodeCredentials{
			URL:            url,
			DeploymentName: deployment.Name,
			Username:       username,
			Password:       password,
			NodePassword:   nodePassword,
		}

		clNodeCredentials = append(clNodeCredentials, clNodeCredential)
	}

	return clNodeCredentials, nil
}
