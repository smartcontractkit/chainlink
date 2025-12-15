package infra

import (
	"fmt"
	"strings"
)

// NodeCredentials holds API credentials for Chainlink nodes
type NodeCredentials struct {
	User     string
	Password string
}

// GetNodeCredentials returns the appropriate API credentials based on infrastructure configuration
// Priority: TOML config > Defaults
func GetNodeCredentials(p *Provider) NodeCredentials {
	creds := NodeCredentials{
		User:     "admin@chain.link", // default
		Password: "password",         // default
	}

	// Kubernetes can override via TOML config
	if p.IsKubernetes() && p.Kubernetes != nil {
		if p.Kubernetes.NodeAPIUser != "" {
			creds.User = p.Kubernetes.NodeAPIUser
		}
		if p.Kubernetes.NodeAPIPassword != "" {
			creds.Password = p.Kubernetes.NodeAPIPassword
		}
	}

	return creds
}

type Type = string
type CribProvider = string

const (
	CRIB       Type         = "crib"
	Kubernetes Type         = "kubernetes"
	Docker     Type         = "docker"
	AWS        CribProvider = "aws"
	Kind       CribProvider = "kind"

	CribConfigsDir = "crib-configs"
)

type Provider struct {
	Type       string           `toml:"type" validate:"oneof=crib kubernetes docker"`
	CRIB       *CRIBInput       `toml:"crib"`
	Kubernetes *KubernetesInput `toml:"kubernetes"`
}

func (i *Provider) IsCRIB() bool {
	return strings.EqualFold(i.Type, CRIB)
}

func (i *Provider) IsKubernetes() bool {
	return strings.EqualFold(i.Type, Kubernetes)
}

func (i *Provider) IsDocker() bool {
	return strings.EqualFold(i.Type, Docker)
}

// Unfortunately, we need to construct some of these URLs before any environment is created, because they are used
// in CL node configs. This introduces a coupling between Helm charts used by CRIB/Kubernetes and Docker container names used by CTFv2.
func (i *Provider) InternalHost(nodeIndex int, isBootstrap bool, donName string) string {
	if i.IsCRIB() || i.IsKubernetes() {
		if isBootstrap {
			return fmt.Sprintf("%s-bt-%d", donName, nodeIndex)
		}
		return fmt.Sprintf("%s-%d", donName, nodeIndex)
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}

func (i *Provider) InternalGatewayHost(nodeIndex int, isBootstrap bool, donName string) string {
	if i.IsCRIB() || i.IsKubernetes() {
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
	if i.IsCRIB() {
		domain := i.CRIB.ExternalDomain
		if domain == "" {
			domain = "main.stage.cldev.sh" // fallback for backward compatibility
		}
		return fmt.Sprintf("%s-gateway.%s", i.CRIB.Namespace, domain)
	}
	if i.IsKubernetes() {
		return fmt.Sprintf("%s-gateway.%s", i.Kubernetes.Namespace, i.Kubernetes.ExternalDomain)
	}

	return "localhost"
}

func (i *Provider) ExternalGatewayPort(dockerPort int) int {
	if i.IsCRIB() || i.IsKubernetes() {
		return 80
	}

	return dockerPort
}

type CRIBInput struct {
	Namespace string `toml:"namespace" validate:"required"`
	// absolute path to the folder with CRIB CRE
	FolderLocation string `toml:"folder_location" validate:"required"`
	ExternalDomain string `toml:"external_domain"`
	Provider       string `toml:"provider" validate:"oneof=aws kind"`
	// required for cost attribution in AWS
	TeamInput *Team `toml:"team_input" validate:"required_if=Provider aws"`
}

// k8s cost attribution
type Team struct {
	Team       string `toml:"team" validate:"required"`
	Product    string `toml:"product" validate:"required"`
	CostCenter string `toml:"cost_center" validate:"required"`
	Component  string `toml:"component" validate:"required"`
}
