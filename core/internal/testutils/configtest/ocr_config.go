package configtest

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
)

var _ config.OCR1Config = &TestGeneralConfig{}

// OCRKeyBundleID returns the overridden value, if one exists.
func (c *TestGeneralConfig) OCRKeyBundleID() (string, error) {
	if c.Overrides.OCRKeyBundleID.Valid {
		return c.Overrides.OCRKeyBundleID.String, nil
	}
	return c.GeneralConfig.OCRKeyBundleID()
}

// OCRDatabaseTimeout returns the overridden value, if one exists.
func (c *TestGeneralConfig) OCRDatabaseTimeout() time.Duration {
	if c.Overrides.OCRDatabaseTimeout != nil {
		return *c.Overrides.OCRDatabaseTimeout
	}
	v, ok := c.GeneralConfig.GlobalOCRDatabaseTimeout()
	if !ok {
		return 1 * time.Second
	}
	return v
}

// OCRObservationGracePeriod returns the overridden value, if one exists.
func (c *TestGeneralConfig) OCRObservationGracePeriod() time.Duration {
	if c.Overrides.OCRObservationGracePeriod != nil {
		return *c.Overrides.OCRObservationGracePeriod
	}
	v, ok := c.GeneralConfig.GlobalOCRObservationGracePeriod()
	if !ok {
		return 100 * time.Millisecond
	}
	return v
}

// OCRObservationTimeout returns the overridden value, if one exists.
func (c *TestGeneralConfig) OCRObservationTimeout() time.Duration {
	if c.Overrides.OCRObservationTimeout != nil {
		return *c.Overrides.OCRObservationTimeout
	}
	return c.GeneralConfig.OCRObservationTimeout()
}

// OCRTransmitterAddress returns the overridden value, if one exists.
func (c *TestGeneralConfig) OCRTransmitterAddress() (ethkey.EIP55Address, error) {
	if c.Overrides.OCRTransmitterAddress != nil {
		return *c.Overrides.OCRTransmitterAddress, nil
	}
	return c.GeneralConfig.OCRTransmitterAddress()
}
