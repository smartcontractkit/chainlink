package config

import (
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/config/parse"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// OCR1Config is a subset of global config relevant to OCR v1.
type OCR1Config interface {
	// OCR1 config, can override in jobs, only ethereum.
	OCRBlockchainTimeout() time.Duration
	OCRContractPollInterval() time.Duration
	OCRContractSubscribeInterval() time.Duration
	OCRKeyBundleID() (string, error)
	OCRObservationTimeout() time.Duration
	OCRSimulateTransactions() bool
	OCRTransmitterAddress() (ethkey.EIP55Address, error) // OCR2 can support non-evm changes
	// OCR1 config, cannot override in jobs
	OCRTraceLogging() bool
	OCRDefaultTransactionQueueDepth() uint32
}

func (c *generalConfig) getDuration(field string) time.Duration {
	return c.getWithFallback(field, parse.Duration).(time.Duration)
}

func (c *generalConfig) GlobalOCRContractConfirmations() (uint16, bool) {
	return lookupEnv(c, envvar.Name("OCRContractConfirmations"), parse.Uint16)
}

func (c *generalConfig) GlobalOCRObservationGracePeriod() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("OCRObservationGracePeriod"), time.ParseDuration)
}

func (c *generalConfig) GlobalOCRContractTransmitterTransmitTimeout() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("OCRContractTransmitterTransmitTimeout"), time.ParseDuration)
}

func (c *generalConfig) GlobalOCRDatabaseTimeout() (time.Duration, bool) {
	return lookupEnv(c, envvar.Name("OCRDatabaseTimeout"), time.ParseDuration)
}

func (c *generalConfig) OCRContractPollInterval() time.Duration {
	return c.getDuration("OCRContractPollInterval")
}

func (c *generalConfig) OCRContractSubscribeInterval() time.Duration {
	return c.getDuration("OCRContractSubscribeInterval")
}

func (c *generalConfig) OCRBlockchainTimeout() time.Duration {
	return c.getDuration("OCRBlockchainTimeout")
}

func (c *generalConfig) OCRKeyBundleID() (string, error) {
	kbStr := c.viper.GetString(envvar.Name("OCRKeyBundleID"))
	if kbStr != "" {
		_, err := models.Sha256HashFromHex(kbStr)
		if err != nil {
			return "", errors.Wrapf(ErrEnvInvalid, "OCR_KEY_BUNDLE_ID is an invalid sha256 hash hex string %v", err)
		}
	}
	return kbStr, nil
}

// OCRDefaultTransactionQueueDepth controls the queue size for DropOldestStrategy in OCR
// Set to 0 to use SendEvery strategy instead
func (c *generalConfig) OCRDefaultTransactionQueueDepth() uint32 {
	return c.viper.GetUint32(envvar.Name("OCRDefaultTransactionQueueDepth"))
}

// OCRTraceLogging determines whether OCR logs at TRACE level are enabled. The
// option to turn them off is given because they can be very verbose
func (c *generalConfig) OCRTraceLogging() bool {
	return c.viper.GetBool(envvar.Name("OCRTraceLogging"))
}

func (c *generalConfig) OCRObservationTimeout() time.Duration {
	return c.getDuration("OCRObservationTimeout")
}

// OCRSimulateTransactions enables using eth_call transaction simulation before
// sending when set to true
func (c *generalConfig) OCRSimulateTransactions() bool {
	return c.viper.GetBool(envvar.Name("OCRSimulateTransactions"))
}

func (c *generalConfig) OCRTransmitterAddress() (ethkey.EIP55Address, error) {
	taStr := c.viper.GetString(envvar.Name("OCRTransmitterAddress"))
	if taStr != "" {
		ta, err := ethkey.NewEIP55Address(taStr)
		if err != nil {
			return "", errors.Wrapf(ErrEnvInvalid, "OCR_TRANSMITTER_ADDRESS is invalid EIP55 %v", err)
		}
		return ta, nil
	}
	return "", errors.Wrap(ErrEnvUnset, "OCR_TRANSMITTER_ADDRESS env var is not set")
}
