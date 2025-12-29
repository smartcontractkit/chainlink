package infra

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

const (
	K8sConfigsDir = "k8s-configs"

	// Kubernetes service naming patterns
	NodeBootstrapNameFormat = "%s-bt-%d" // e.g., workflow-bt-0
	NodePluginNameFormat    = "%s-%d"    // e.g., workflow-0

	// Job Distributor service names
	JDGRPCServiceName    = "job-distributor:80"
	JDNodeRPCServiceName = "job-distributor-noderpc-lb:80"
)

// KubernetesInput contains configuration for Kubernetes-based infrastructure
type KubernetesInput struct {
	Namespace string `toml:"namespace" validate:"required"`
	// absolute path to the folder with Kubernetes CRE configs
	FolderLocation  string `toml:"folder_location"`
	ExternalDomain  string `toml:"external_domain"`
	ExternalPort    int    `toml:"external_port"`
	LabelSelector   string `toml:"label_selector"`
	NodeAPIUser     string `toml:"node_api_user"`
	NodeAPIPassword string `toml:"node_api_password"`
	// Use TLS for external JD GRPC connections (default: true for port 443, false otherwise)
	UseTLSForJD *bool `toml:"use_tls_for_jd"`
	// Kubernetes provider type (aws, kind, etc.) - optional
	Provider string `toml:"provider" validate:"oneof=aws kind ''"`
	// Required for cost attribution in AWS
	TeamInput *Team `toml:"team_input" validate:"required_if=Provider aws"`
}

// GetKubernetesClient creates a Kubernetes client from the default kubeconfig or in-cluster config
func GetKubernetesClient() (*kubernetes.Clientset, error) {
	// Try to use in-cluster config first (when running inside a pod)
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		config, err = kubeConfig.ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return clientset, nil
}

// GenerateKubernetesJDOutput generates JD service URLs for Kubernetes-deployed Job Distributor
func GenerateKubernetesJDOutput(infraInput *Provider, lggr zerolog.Logger) (*jd.Output, error) {
	if infraInput.Kubernetes == nil {
		return nil, errors.New("kubernetes configuration is required for GenerateKubernetesJDOutput")
	}

	externalDomain := infraInput.Kubernetes.ExternalDomain

	// Kubernetes service naming for JD
	// Internal URLs use short service names (no namespace suffix) - they're accessed from within same namespace
	internalGRPCUrl := JDGRPCServiceName
	internalWSRPCUrl := JDNodeRPCServiceName // Load balancer for node registration

	// External URLs
	externalGRPCUrl := internalGRPCUrl
	if externalDomain != "" {
		// Use external domain - just hostname:port (no protocol prefix)
		externalGRPCUrl = fmt.Sprintf("%s-job-distributor-grpc.%s:443", infraInput.Kubernetes.Namespace, externalDomain)
	}

	lggr.Info().Msgf("Generated JD URLs - External GRPC: %s, Internal GRPC: %s, Internal WSRPC: %s",
		externalGRPCUrl, internalGRPCUrl, internalWSRPCUrl)

	return &jd.Output{
		UseCache:         true,
		ExternalGRPCUrl:  externalGRPCUrl,
		InternalGRPCUrl:  internalGRPCUrl,
		InternalWSRPCUrl: internalWSRPCUrl,
	}, nil
}

// GenerateNodeInstanceNames creates Kubernetes-compatible instance names for nodes
// Bootstrap nodes get names like "workflow-bt-0", plugin nodes get "workflow-0", "workflow-1", etc.
// The nodeMetadataRoles slice indicates whether each node is a bootstrap node (true) or plugin node (false)
func GenerateNodeInstanceNames(nodeSetName string, nodeMetadataRoles []bool) []string {
	instanceNames := make([]string, len(nodeMetadataRoles))
	pluginNodeCounter := 0
	bootstrapNodeCounter := 0

	for i, isBootstrap := range nodeMetadataRoles {
		if isBootstrap {
			instanceNames[i] = fmt.Sprintf(NodeBootstrapNameFormat, nodeSetName, bootstrapNodeCounter)
			bootstrapNodeCounter++
		} else {
			instanceNames[i] = fmt.Sprintf(NodePluginNameFormat, nodeSetName, pluginNodeCounter)
			pluginNodeCounter++
		}
	}

	return instanceNames
}

// GenerateKubernetesNodeSetOutput generates node URLs for Kubernetes-deployed Chainlink nodes
// It creates service URLs matching the Helm chart naming convention:
// - Bootstrap nodes: {donName}-bt-{index}
// - Plugin nodes: {donName}-{counter} (counter starts from 0, separate from bootstrap)
func GenerateKubernetesNodeSetOutput(infraInput *Provider, nodeSetName string, nodeCount int, nodeMetadataRoles []bool, creds NodeCredentials, lggr zerolog.Logger) *ns.Output {
	externalDomain := ""
	if infraInput.Kubernetes != nil {
		externalDomain = infraInput.Kubernetes.ExternalDomain
	}

	clNodes := make([]*clnode.Output, nodeCount)

	// Generate instance names for all nodes
	instanceNames := GenerateNodeInstanceNames(nodeSetName, nodeMetadataRoles)

	for i := 0; i < nodeCount; i++ {
		serviceName := instanceNames[i]
		namespace := ""
		if infraInput.Kubernetes != nil {
			namespace = infraInput.Kubernetes.Namespace
		}
		externalServiceName := namespace + "-" + serviceName

		// Internal URLs - use short service names (no namespace suffix when in same namespace)
		internalURL := fmt.Sprintf("http://%s:6688", serviceName)
		internalP2PUrl := fmt.Sprintf("http://%s:5001", serviceName)

		// External URLs
		externalURL := internalURL
		if externalDomain != "" {
			externalURL = fmt.Sprintf("https://%s.%s", externalServiceName, externalDomain)
		}

		lggr.Info().Msgf("Generated URLs for node %d of DON %s (bootstrap=%v) - Internal: %s, External: %s",
			i, nodeSetName, nodeMetadataRoles[i], internalURL, externalURL)

		clNodes[i] = &clnode.Output{
			UseCache: true,
			Node: &clnode.NodeOut{
				APIAuthUser:     creds.User,
				APIAuthPassword: creds.Password,
				ExternalURL:     externalURL,
				InternalURL:     internalURL,
				InternalP2PUrl:  internalP2PUrl,
				InternalIP:      serviceName,
			},
		}
	}

	return &ns.Output{
		UseCache: true,
		CLNodes:  clNodes,
	}
}
