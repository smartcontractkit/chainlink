package src

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Taken from https://github.com/smartcontractkit/chainlink-testing-framework/blob/7ae47c88ecbb8a483ffe4d1e796704d1245f54d0/k8s/client/client.go

// K8sClient high level k8s client
type K8sClient struct {
	ClientSet  *kubernetes.Clientset
	RESTConfig *rest.Config
	KubeConfig clientcmd.ClientConfig
	namespace  string
}

// getLocalK8sDeps get local k8s context config
func getLocalK8sDeps() (*kubernetes.Clientset, *rest.Config, clientcmd.ClientConfig, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	k8sConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, nil, nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, nil, nil, err
	}
	return k8sClient, k8sConfig, kubeConfig, nil
}

// MustNewK8sClient creates a new k8s client with a REST config
func MustNewK8sClient() *K8sClient {
	cs, restConfig, kubeConfig, err := getLocalK8sDeps()
	if err != nil {
		log.Fatalf("Failed to create k8s client: %v", err)
	}

	namespace, overridden, err := kubeConfig.Namespace()
	if err != nil {
		log.Fatalf("Failed to get namespace: %v", err)
	}
	if overridden {
		fmt.Println("Namespace overridden to: ", namespace)
	}

	return &K8sClient{
		ClientSet:  cs,
		RESTConfig: restConfig,
		KubeConfig: kubeConfig,
		namespace:  namespace,
	}
}

type DeploymentWithConfigMap struct {
	apps.Deployment
	ServiceName string
	ConfigMap   v1.ConfigMap
	Host        string
}

func (m *K8sClient) GetDeploymentsWithConfigMap() ([]DeploymentWithConfigMap, error) {
	deployments, err := m.ListDeployments("app=app")
	if err != nil {
		return nil, err
	}
	if len(deployments.Items) == 0 {
		return nil, errors.New("no deployments found, is your nodeset deployed?")
	}

	deploymentsWithConfigMaps := []DeploymentWithConfigMap{}
	ingressList, err := m.ListIngresses()
	if err != nil {
		return nil, err
	}
	if len(ingressList.Items) == 0 {
		return nil, errors.New("no ingress found, is your nodeset deployed?")
	}

	for _, deployment := range deployments.Items {
		for _, v := range deployment.Spec.Template.Spec.Volumes {
			if v.ConfigMap == nil {
				continue
			}
			cm, err := m.GetConfigMap(v.ConfigMap.Name)
			if err != nil {
				return nil, err
			}
			instance := deployment.Labels["instance"]
			var host string
			var serviceName string
			for _, ingress := range ingressList.Items {
				for _, rule := range ingress.Spec.Rules {
					for _, path := range rule.HTTP.Paths {
						if strings.Contains(path.Backend.Service.Name, instance) {
							host = rule.Host
							serviceName = path.Backend.Service.Name
						}
					}
				}
			}

			if host == "" {
				return nil, fmt.Errorf("could not find host for deployment %s", deployment.Name)
			}

			deploymentWithConfigMap := DeploymentWithConfigMap{
				Host:        host,
				ServiceName: serviceName,
				Deployment:  deployment,
				ConfigMap:   *cm,
			}
			deploymentsWithConfigMaps = append(deploymentsWithConfigMaps, deploymentWithConfigMap)
		}
	}

	fmt.Printf("Found %d deployments with config maps\n", len(deploymentsWithConfigMaps))
	return deploymentsWithConfigMaps, nil
}

// ListDeployments lists deployments for a namespace
func (m *K8sClient) ListDeployments(selector string) (*apps.DeploymentList, error) {
	deployments, err := m.ClientSet.AppsV1().Deployments(m.namespace).List(context.Background(), metaV1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	sort.Slice(deployments.Items, func(i, j int) bool {
		return deployments.Items[i].CreationTimestamp.Before(deployments.Items[j].CreationTimestamp.DeepCopy())
	})
	return deployments.DeepCopy(), nil
}

// Get a config map
func (m *K8sClient) GetConfigMap(name string) (*v1.ConfigMap, error) {
	configMap, err := m.ClientSet.CoreV1().ConfigMaps(m.namespace).Get(context.Background(), name, metaV1.GetOptions{})

	return configMap.DeepCopy(), err
}

func (m *K8sClient) ListIngresses() (*networkingV1.IngressList, error) {
	ingressList, err := m.ClientSet.NetworkingV1().Ingresses(m.namespace).List(context.Background(), metaV1.ListOptions{})

	return ingressList.DeepCopy(), err
}
