package configtest

import (
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

var _ config.P2PV1Networking = &TestGeneralConfig{}

// P2PV2DeltaDial returns the overridden value, if one exists.
func (c *TestGeneralConfig) P2PV2DeltaDial() models.Duration {
	if c.Overrides.P2PV2DeltaDial != nil {
		return models.MustMakeDuration(*c.Overrides.P2PV2DeltaDial)
	}
	return c.GeneralConfig.P2PV2DeltaDial()
}

// P2PV2Bootstrappers returns the overridden value, if one exists.
func (c *TestGeneralConfig) P2PV2Bootstrappers() []ocrcommontypes.BootstrapperLocator {
	if len(c.Overrides.P2PV2Bootstrappers) != 0 {
		return c.Overrides.P2PV2Bootstrappers
	}
	return c.GeneralConfig.P2PV2Bootstrappers()
}

// P2PV2ListenAddresses returns the overridden value, if one exists.
func (c *TestGeneralConfig) P2PV2ListenAddresses() []string {
	if len(c.Overrides.P2PV2ListenAddresses) != 0 {
		return c.Overrides.P2PV2ListenAddresses
	}
	return c.GeneralConfig.P2PV2ListenAddresses()
}

// P2PV2AnnounceAddresses returns the overridden value, if one exists.
func (c *TestGeneralConfig) P2PV2AnnounceAddresses() []string {
	if len(c.Overrides.P2PV2AnnounceAddresses) != 0 {
		return c.Overrides.P2PV2AnnounceAddresses
	}
	return c.GeneralConfig.P2PV2AnnounceAddresses()
}

// P2PV2DeltaReconcile returns the overridden value, if one exists.
func (c *TestGeneralConfig) P2PV2DeltaReconcile() models.Duration {
	if c.Overrides.P2PV2DeltaReconcile != nil {
		return models.MustMakeDuration(*c.Overrides.P2PV2DeltaReconcile)
	}
	return c.GeneralConfig.P2PV2DeltaReconcile()
}
