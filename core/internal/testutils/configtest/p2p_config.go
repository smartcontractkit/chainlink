package configtest

import (
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
)

var _ config.P2PNetworking = &TestGeneralConfig{}

func (c *TestGeneralConfig) P2PNetworkingStack() ocrnetworking.NetworkingStack {
	if c.Overrides.P2PNetworkingStack != 0 {
		return c.Overrides.P2PNetworkingStack
	}
	return c.GeneralConfig.P2PNetworkingStack()
}

func (c *TestGeneralConfig) P2PPeerID() p2pkey.PeerID {
	if c.Overrides.P2PPeerID.String() != "" {
		return c.Overrides.P2PPeerID
	}
	return ""
}
