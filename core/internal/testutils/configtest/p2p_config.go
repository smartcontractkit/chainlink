package configtest

import (
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
)

var _ config.P2PNetworking = &TestGeneralConfig{}

// P2PNetworkingStack stack returns the overridden value, if one exists.
func (c *TestGeneralConfig) P2PNetworkingStack() ocrnetworking.NetworkingStack {
	if c.Overrides.P2PNetworkingStack != 0 {
		return c.Overrides.P2PNetworkingStack
	}
	return c.GeneralConfig.P2PNetworkingStack()
}

// P2PPeerID returns the overridden value or empty.
func (c *TestGeneralConfig) P2PPeerID() p2pkey.PeerID {
	if c.Overrides.P2PPeerID.String() != "" {
		return c.Overrides.P2PPeerID
	}
	return ""
}
