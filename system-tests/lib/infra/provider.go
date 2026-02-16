package infra

import (
	"fmt"
	"strings"
)

type Type = string

const (
	Kubernetes Type = "kubernetes"
	Docker     Type = "docker"
)

type Provider struct {
	Type       string           `toml:"type" validate:"oneof=kubernetes docker"`
	Kubernetes *KubernetesInput `toml:"kubernetes"`
}

func (i *Provider) IsKubernetes() bool {
	return strings.EqualFold(i.Type, Kubernetes)
}

func (i *Provider) IsDocker() bool {
	return strings.EqualFold(i.Type, Docker)
}

// Unfortunately, we need to construct some of these URLs before any environment is created, because they are used
// in CL node configs. This introduces a coupling between Helm charts used by Kubernetes and Docker container names used by CTFv2.
func (i *Provider) InternalHost(nodeIndex int, isBootstrap bool, donName string) string {
	if i.IsKubernetes() {
		if isBootstrap {
			return fmt.Sprintf("%s-bt-%d", donName, nodeIndex)
		}
		return fmt.Sprintf("%s-%d", donName, nodeIndex)
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}

func (i *Provider) InternalGatewayHost(nodeIndex int, isBootstrap bool, donName string) string {
	if i.IsKubernetes() {
		host := fmt.Sprintf("%s-%d", donName, nodeIndex)
		if isBootstrap {
			host = fmt.Sprintf("%s-bt-%d", donName, nodeIndex)
		}
		host += "-gtwnode"

		return host
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}

func (i *Provider) ExternalGatewayHost() string {
	if i.IsKubernetes() {
		return fmt.Sprintf("%s-gateway.%s", i.Kubernetes.Namespace, i.Kubernetes.ExternalDomain)
	}

	return "localhost"
}

func (i *Provider) ExternalGatewayPort(dockerPort int) int {
	if i.IsKubernetes() {
		return 80
	}

	return dockerPort
}
