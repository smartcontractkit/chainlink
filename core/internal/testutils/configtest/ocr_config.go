package configtest

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
)

var _ config.OCRConfig = &TestGeneralConfig{}

func (c *TestGeneralConfig) OCR2DatabaseTimeout() time.Duration {
	if c.Overrides.OCR2DatabaseTimeout != nil {
		return *c.Overrides.OCR2DatabaseTimeout
	}
	return c.GeneralConfig.OCR2DatabaseTimeout()
}

func (c *TestGeneralConfig) OCRKeyBundleID() (string, error) {
	if c.Overrides.OCRKeyBundleID.Valid {
		return c.Overrides.OCRKeyBundleID.String, nil
	}
	return c.GeneralConfig.OCRKeyBundleID()
}

func (c *TestGeneralConfig) OCRDatabaseTimeout() time.Duration {
	if c.Overrides.OCRDatabaseTimeout != nil {
		return *c.Overrides.OCRDatabaseTimeout
	}
	return c.GeneralConfig.OCRDatabaseTimeout()
}

func (c *TestGeneralConfig) OCRObservationGracePeriod() time.Duration {
	if c.Overrides.OCRObservationGracePeriod != nil {
		return *c.Overrides.OCRObservationGracePeriod
	}
	return c.GeneralConfig.OCRObservationGracePeriod()
}

func (c *TestGeneralConfig) OCRObservationTimeout() time.Duration {
	if c.Overrides.OCRObservationTimeout != nil {
		return *c.Overrides.OCRObservationTimeout
	}
	return c.GeneralConfig.OCRObservationTimeout()
}

func (c *TestGeneralConfig) OCRTransmitterAddress() (ethkey.EIP55Address, error) {
	if c.Overrides.OCRTransmitterAddress != nil {
		return *c.Overrides.OCRTransmitterAddress, nil
	}
	return c.GeneralConfig.OCRTransmitterAddress()
}

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
