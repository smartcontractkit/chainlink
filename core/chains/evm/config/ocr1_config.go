package config

import "time"

func (c *chainScopedConfig) OCRContractConfirmations() uint16 {
	val, ok := c.GeneralConfig.GlobalOCRContractConfirmations()
	if ok {
		c.logEnvOverrideOnce("OCRContractConfirmations", val)
		return val
	}
	return c.defaultSet.ocrContractConfirmations
}

func (c *chainScopedConfig) OCRContractTransmitterTransmitTimeout() time.Duration {
	val, ok := c.GeneralConfig.GlobalOCRContractTransmitterTransmitTimeout()
	if ok {
		c.logEnvOverrideOnce("OCRContractTransmitterTransmitTimeout", val)
		return val
	}
	return c.defaultSet.ocrContractTransmitterTransmitTimeout
}

func (c *chainScopedConfig) OCRDatabaseTimeout() time.Duration {
	val, ok := c.GeneralConfig.GlobalOCRDatabaseTimeout()
	if ok {
		c.logEnvOverrideOnce("OCRDatabaseTimeout", val)
		return val
	}
	return c.defaultSet.ocrDatabaseTimeout
}

func (c *chainScopedConfig) OCRObservationGracePeriod() time.Duration {
	val, ok := c.GeneralConfig.GlobalOCRObservationGracePeriod()
	if ok {
		c.logEnvOverrideOnce("OCRObservationGracePeriod", val)
		return val
	}
	return c.defaultSet.ocrObservationGracePeriod
}
