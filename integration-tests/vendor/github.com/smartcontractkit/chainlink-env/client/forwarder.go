package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type Forwarder struct {
	Client         *K8sClient
	mu             *sync.Mutex
	KeepConnection bool
	Info           map[string]interface{}
}

type ConnectionInfo struct {
	Ports portforward.ForwardedPort
	Host  string
}

func NewForwarder(client *K8sClient, keepConnection bool) *Forwarder {
	return &Forwarder{
		Client:         client,
		mu:             &sync.Mutex{},
		KeepConnection: keepConnection,
		Info:           make(map[string]interface{}),
	}
}

func (m *Forwarder) forwardPodPorts(pod v1.Pod, namespaceName string) error {
	if pod.Status.Phase != v1.PodRunning {
		log.Debug().Str("Pod", pod.Name).Interface("Phase", pod.Status.Phase).Msg("Skipping pod for port forwarding")
		return nil
	}
	roundTripper, upgrader, err := spdy.RoundTripperFor(m.Client.RESTConfig)
	if err != nil {
		return err
	}
	httpPath := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespaceName, pod.Name)
	hostIP := strings.TrimLeft(m.Client.RESTConfig.Host, "htps:/")
	serverURL := url.URL{Scheme: "https", Path: httpPath, Host: hostIP}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)

	portRules := m.portRulesForPod(pod)
	if len(portRules) == 0 {
		return nil
	}

	// porforward is not thread safe for using multiple rules in the same forwarder,
	// at least not until this pr is merged: https://github.com/kubernetes/kubernetes/pull/114342
	forwardedPorts := []portforward.ForwardedPort{}
	for _, portRule := range portRules {
		stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
		out, errOut := new(bytes.Buffer), new(bytes.Buffer)

		log.Debug().
			Str("Pod", pod.Name).
			Msg("Attempting to forward ports")

		forwarder, err := portforward.New(dialer, []string{portRule}, stopChan, readyChan, out, errOut)
		if err != nil {
			return err
		}
		go func() {
			if err := forwarder.ForwardPorts(); err != nil {
				log.Error().Str("Pod", pod.Name).Err(err)
			}
		}()

		<-readyChan
		if len(errOut.String()) > 0 {
			return fmt.Errorf("error on forwarding k8s port: %v", errOut.String())
		}
		fP, err := forwarder.GetPorts()
		if err != nil {
			return err
		}
		forwardedPorts = append(forwardedPorts, fP...)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	namedPorts := m.podPortsByName(pod, forwardedPorts)
	if pod.Labels["app"] != "" {
		m.Info[fmt.Sprintf("%s:%s", pod.Labels["app"], pod.Labels["instance"])] = namedPorts
	}
	return nil
}

func (m *Forwarder) collectPodPorts(pod v1.Pod) error {
	namedPorts := make(map[string]interface{})
	for _, c := range pod.Spec.Containers {
		for _, cp := range c.Ports {
			if namedPorts[c.Name] == nil {
				namedPorts[c.Name] = make(map[string]interface{})
			}
			namedPorts[c.Name].(map[string]interface{})[cp.Name] = ConnectionInfo{
				Host:  pod.Status.PodIP,
				Ports: portforward.ForwardedPort{Remote: uint16(cp.ContainerPort)},
			}
		}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if pod.Labels["app"] != "" {
		m.Info[fmt.Sprintf("%s:%s", pod.Labels["app"], pod.Labels["instance"])] = namedPorts
	}
	return nil
}

func (m *Forwarder) podPortsByName(pod v1.Pod, fp []portforward.ForwardedPort) map[string]interface{} {
	ports := make(map[string]interface{})
	for _, forwardedPort := range fp {
		for _, c := range pod.Spec.Containers {
			for _, cp := range c.Ports {
				if cp.ContainerPort == int32(forwardedPort.Remote) {
					if ports[c.Name] == nil {
						ports[c.Name] = make(map[string]interface{})
					}
					ports[c.Name].(map[string]interface{})[cp.Name] = ConnectionInfo{
						Host:  pod.Status.PodIP,
						Ports: forwardedPort,
					}
				}
			}
		}
	}
	return ports
}

func (m *Forwarder) portRulesForPod(pod v1.Pod) []string {
	rules := make([]string, 0)
	for _, c := range pod.Spec.Containers {
		for _, port := range c.Ports {
			rules = append(rules, fmt.Sprintf(":%d", port.ContainerPort))
		}
	}
	return rules
}

func (m *Forwarder) Connect(namespaceName string, selector string, insideK8s bool) error {
	m.Info = make(map[string]interface{})
	pods, err := m.Client.ListPods(namespaceName, selector)
	if err != nil {
		return err
	}
	eg := &errgroup.Group{}
	for _, p := range pods.Items {
		p := p
		if insideK8s {
			eg.Go(func() error {
				return m.collectPodPorts(p)
			})
		} else {
			eg.Go(func() error {
				return m.forwardPodPorts(p, namespaceName)
			})
		}
	}
	return eg.Wait()
}

// PrintLocalPorts prints all local forwarded ports
func (m *Forwarder) PrintLocalPorts() {
	for labeledAppPodName, labeledAppPod := range m.Info {
		for containerName, container := range labeledAppPod.(map[string]interface{}) {
			for fpName, portsData := range container.(map[string]interface{}) {
				log.Info().
					Str("Label", labeledAppPodName).
					Str("Container", containerName).
					Str("PortNames", fpName).
					Uint16("Port", portsData.(ConnectionInfo).Ports.Local).
					Msg("Local ports")
			}
		}
	}
}

func (m *Forwarder) FindPort(ks ...string) *URLConverter {
	d, err := lookupMap(m.Info, ks...)
	return NewURLConverter(d.(ConnectionInfo), err)
}

func lookupMap(m map[string]interface{}, ks ...string) (rval interface{}, err error) {
	var ok bool
	if len(ks) == 0 {
		return nil, fmt.Errorf("select port path like $app_name:$instance $container_name $port_name")
	}
	if rval, ok = m[ks[0]]; !ok {
		return ConnectionInfo{}, fmt.Errorf("key not found: '%s' remaining keys: %s, provided map: %s", ks[0], ks, m)
	} else if len(ks) == 1 {
		return rval, nil
	} else {
		return lookupMap(m[ks[0]].(map[string]interface{}), ks[1:]...)
	}
}
