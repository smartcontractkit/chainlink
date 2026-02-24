package infra

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
