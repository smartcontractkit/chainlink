package configtest

import (
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/store/models"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
)

var _ config.P2PV1Networking = &TestGeneralConfig{}

func (c *TestGeneralConfig) P2PV2DeltaDial() models.Duration {
	if c.Overrides.P2PV2DeltaDial != nil {
		return models.MustMakeDuration(*c.Overrides.P2PV2DeltaDial)
	}
	return c.GeneralConfig.P2PV2DeltaDial()
}

func (c *TestGeneralConfig) P2PV2Bootstrappers() []ocrcommontypes.BootstrapperLocator {
	if len(c.Overrides.P2PV2Bootstrappers) != 0 {
		return c.Overrides.P2PV2Bootstrappers
	}
	return c.GeneralConfig.P2PV2Bootstrappers()
}

func (c *TestGeneralConfig) P2PV2ListenAddresses() []string {
	if len(c.Overrides.P2PV2ListenAddresses) != 0 {
		return c.Overrides.P2PV2ListenAddresses
	}
	return c.GeneralConfig.P2PV2ListenAddresses()
}

func (c *TestGeneralConfig) P2PV2AnnounceAddresses() []string {
	if len(c.Overrides.P2PV2AnnounceAddresses) != 0 {
		return c.Overrides.P2PV2AnnounceAddresses
	}
	return c.GeneralConfig.P2PV2AnnounceAddresses()
}

func (c *TestGeneralConfig) P2PV2DeltaReconcile() models.Duration {
	if c.Overrides.P2PV2DeltaReconcile != nil {
		return models.MustMakeDuration(*c.Overrides.P2PV2DeltaReconcile)
	}
	return c.GeneralConfig.P2PV2DeltaReconcile()
}
