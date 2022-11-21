package config

import (
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/config/parse"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// OCR2Config is a subset of global config relevant to OCR v2.
type OCR2Config interface {
	// OCR2 config, can override in jobs, all chains
	OCR2ContractConfirmations() uint16
	OCR2ContractTransmitterTransmitTimeout() time.Duration
	OCR2BlockchainTimeout() time.Duration
	OCR2DatabaseTimeout() time.Duration
	OCR2ContractPollInterval() time.Duration
	OCR2ContractSubscribeInterval() time.Duration
	OCR2KeyBundleID() (string, error)
	// OCR2 config, cannot override in jobs
	OCR2TraceLogging() bool
}

func (c *generalConfig) OCR2ContractConfirmations() uint16 {
	return getEnvWithFallback(c, envvar.NewUint16("OCR2ContractConfirmations"))
}

func (c *generalConfig) OCR2ContractPollInterval() time.Duration {
	return c.getDuration("OCR2ContractPollInterval")
}

func (c *generalConfig) OCR2ContractSubscribeInterval() time.Duration {
	return c.getDuration("OCR2ContractSubscribeInterval")
}

func (c *generalConfig) OCR2ContractTransmitterTransmitTimeout() time.Duration {
	return c.getWithFallback("OCR2ContractTransmitterTransmitTimeout", parse.Duration).(time.Duration)
}

func (c *generalConfig) OCR2BlockchainTimeout() time.Duration {
	return c.getDuration("OCR2BlockchainTimeout")
}

func (c *generalConfig) OCR2DatabaseTimeout() time.Duration {
	return c.getWithFallback("OCR2DatabaseTimeout", parse.Duration).(time.Duration)
}

func (c *generalConfig) OCR2KeyBundleID() (string, error) {
	kbStr := c.viper.GetString(envvar.Name("OCR2KeyBundleID"))
	if kbStr != "" {
		_, err := models.Sha256HashFromHex(kbStr)
		if err != nil {
			return "", errors.Wrapf(ErrEnvInvalid, "OCR_KEY_BUNDLE_ID is an invalid sha256 hash hex string %v", err)
		}
	}
	return kbStr, nil
}

func (c *generalConfig) OCR2TraceLogging() bool {
	return c.viper.GetBool(envvar.Name("OCRTraceLogging"))
}
