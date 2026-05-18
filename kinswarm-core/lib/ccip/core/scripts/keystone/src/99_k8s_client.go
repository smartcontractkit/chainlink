package src

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

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

type PodWithConfigMap struct {
	v1.Pod
	ConfigMap v1.ConfigMap
	Host      string
}

func (m *K8sClient) GetPodsWithConfigMap() ([]PodWithConfigMap, error) {
	pods, err := m.ListPods("app=app")
	if err != nil {
		return nil, err
	}
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no chainlink node crib pods found, is your crib cluster deployed?")
	}

	podsWithConfigMaps := []PodWithConfigMap{}
	ingressList, err := m.ListIngresses()
	if err != nil {
		return nil, err
	}
	if len(ingressList.Items) == 0 {
		return nil, fmt.Errorf("no ingress found, is your crib cluster deployed?")
	}

	for _, pod := range pods.Items {
		for _, v := range pod.Spec.Volumes {
			if v.ConfigMap == nil {
				continue
			}
			cm, err := m.GetConfigMap(v.ConfigMap.Name)
			if err != nil {
				return nil, err
			}
			// - host: crib-henry-keystone-node2.main.stage.cldev.sh
			// http:
			//   paths:
			//   - backend:
			//       service:
			//         name: app-node-2
			//         port:
			//           number: 6688
			//     path: /*
			//     pathType: ImplementationSpecific
			instance := pod.Labels["instance"]
			var host string
			for _, ingress := range ingressList.Items {
				for _, rule := range ingress.Spec.Rules {
					for _, path := range rule.HTTP.Paths {
						if strings.Contains(path.Backend.Service.Name, instance) {
							host = rule.Host
						}
					}
				}
			}

			if host == "" {
				return nil, fmt.Errorf("could not find host for pod %s", pod.Name)
			}

			podWithConfigMap := PodWithConfigMap{
				Host:      host,
				Pod:       pod,
				ConfigMap: *cm,
			}
			podsWithConfigMaps = append(podsWithConfigMaps, podWithConfigMap)
		}
	}

	fmt.Printf("Found %d chainlink node crib pods\n", len(podsWithConfigMaps))
	return podsWithConfigMaps, nil
}

// ListPods lists pods for a namespace and selector
func (m *K8sClient) ListPods(selector string) (*v1.PodList, error) {
	pods, err := m.ClientSet.CoreV1().Pods(m.namespace).List(context.Background(), metaV1.ListOptions{LabelSelector: selector})
	sort.Slice(pods.Items, func(i, j int) bool {
		return pods.Items[i].CreationTimestamp.Before(pods.Items[j].CreationTimestamp.DeepCopy())
	})

	return pods.DeepCopy(), err
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
