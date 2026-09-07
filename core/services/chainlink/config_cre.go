package chainlink

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

// defaultLinkingRequestTimeout mirrors the [CRE.Linking] RequestTimeout default in docs/core.toml.
const defaultLinkingRequestTimeout = 2 * time.Second

type creConfig struct {
	s toml.CreSecrets
	c toml.CreConfig
}

func (c *creConfig) StreamsAPIKey() string {
	if c.s.Streams == nil || c.s.Streams.APIKey == nil {
		return ""
	}
	return string(*c.s.Streams.APIKey)
}

func (c *creConfig) StreamsAPISecret() string {
	if c.s.Streams == nil || c.s.Streams.APISecret == nil {
		return ""
	}
	return string(*c.s.Streams.APISecret)
}

func (c *creConfig) WsURL() string {
	if c.c.Streams == nil || c.c.Streams.WsURL == nil {
		return ""
	}
	return *c.c.Streams.WsURL
}

func (c *creConfig) RestURL() string {
	if c.c.Streams == nil || c.c.Streams.RestURL == nil {
		return ""
	}
	return *c.c.Streams.RestURL
}

func (c *creConfig) DebugMode() bool {
	if c.c.DebugMode == nil {
		return false // disabled by default
	}
	return *c.c.DebugMode
}

type workflowFetcherConfig struct {
	url string
}

func (w *workflowFetcherConfig) URL() string {
	return w.url
}

func (c *creConfig) WorkflowFetcher() config.WorkflowFetcher {
	if c.c.WorkflowFetcher == nil || c.c.WorkflowFetcher.URL == nil {
		return &workflowFetcherConfig{url: ""}
	}
	return &workflowFetcherConfig{url: *c.c.WorkflowFetcher.URL}
}

func (c *creConfig) UseLocalTimeProvider() bool {
	if c.c.UseLocalTimeProvider == nil {
		return true // default to local time provider since DON Time plugin may not be running
	}
	return *c.c.UseLocalTimeProvider
}

func (c *creConfig) EnableDKGRecipient() bool {
	if c.c.EnableDKGRecipient == nil {
		return false
	}
	return *c.c.EnableDKGRecipient
}

type linkingConfig struct {
	url            string
	tlsEnabled     bool
	requestTimeout time.Duration
}

func (l *linkingConfig) URL() string {
	return l.url
}

func (l *linkingConfig) TLSEnabled() bool {
	return l.tlsEnabled
}

func (l *linkingConfig) RequestTimeout() time.Duration {
	return l.requestTimeout
}

func (c *creConfig) Linking() config.CRELinking {
	if c.c.Linking == nil {
		return &linkingConfig{url: "", tlsEnabled: true, requestTimeout: defaultLinkingRequestTimeout}
	}

	url := ""
	if c.c.Linking.URL != nil {
		url = *c.c.Linking.URL
	}

	tlsEnabled := true // default
	if c.c.Linking.TLSEnabled != nil {
		tlsEnabled = *c.c.Linking.TLSEnabled
	}

	requestTimeout := defaultLinkingRequestTimeout
	if c.c.Linking.RequestTimeout != nil {
		requestTimeout = c.c.Linking.RequestTimeout.Duration()
	}

	return &linkingConfig{url: url, tlsEnabled: tlsEnabled, requestTimeout: requestTimeout}
}

type confidentialRelayConfig struct {
	enabled          bool
	trustEnclaves    bool
	requireBFTQuorum bool
}

func (cr *confidentialRelayConfig) Enabled() bool          { return cr.enabled }
func (cr *confidentialRelayConfig) TrustEnclaves() bool    { return cr.trustEnclaves }
func (cr *confidentialRelayConfig) RequireBFTQuorum() bool { return cr.requireBFTQuorum }

func (c *creConfig) ConfidentialRelay() config.CREConfidentialRelay {
	if c.c.ConfidentialRelay == nil {
		return &confidentialRelayConfig{}
	}
	enabled := false
	if c.c.ConfidentialRelay.Enabled != nil {
		enabled = *c.c.ConfidentialRelay.Enabled
	}
	trustEnclaves := false
	if c.c.ConfidentialRelay.TrustEnclaves != nil {
		trustEnclaves = *c.c.ConfidentialRelay.TrustEnclaves
	}
	requireBFTQuorum := false
	if c.c.ConfidentialRelay.RequireBFTQuorum != nil {
		requireBFTQuorum = *c.c.ConfidentialRelay.RequireBFTQuorum
	}
	return &confidentialRelayConfig{enabled: enabled, trustEnclaves: trustEnclaves, requireBFTQuorum: requireBFTQuorum}
}

func (c *creConfig) LocalSecretOverrides() map[string]map[string]string {
	return c.s.LocalSecretOverrides
}
