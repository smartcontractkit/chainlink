package configtest

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/config"
)

var _ config.P2PV1Networking = &TestGeneralConfig{}

// P2PListenPort returns the overridden value, if one exists.
func (c *TestGeneralConfig) P2PListenPort() uint16 {
	if c.Overrides.P2PListenPort.Valid {
		return uint16(c.Overrides.P2PListenPort.Int64)
	}
	return 12345
}

// P2PBootstrapPeers returns the overridden value, if one exists.
func (c *TestGeneralConfig) P2PBootstrapPeers() ([]string, error) {
	if c.Overrides.P2PBootstrapPeers != nil {
		return c.Overrides.P2PBootstrapPeers, nil
	}
	return c.GeneralConfig.P2PBootstrapPeers()
}

// P2PBootstrapCheckInterval returns the overridden value, if one exists.
func (c *TestGeneralConfig) P2PBootstrapCheckInterval() time.Duration {
	if c.Overrides.P2PBootstrapCheckInterval != nil {
		return *c.Overrides.P2PBootstrapCheckInterval
	}
	return c.GeneralConfig.P2PBootstrapCheckInterval()
}
