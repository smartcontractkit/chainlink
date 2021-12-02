package config

import (
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
)

type OCRConfig interface {
	// OCR2 config, can override in jobs, all chains
	OCR2ContractPollInterval() time.Duration
	OCR2ContractSubscribeInterval() time.Duration
	OCR2ContractTransmitterTransmitTimeout() time.Duration
	OCR2BlockchainTimeout() time.Duration
	OCR2DatabaseTimeout() time.Duration
	OCR2MonitoringEndpoint() string
	OCR2KeyBundleID() (string, error)
	// OCR2 config, cannot override in jobs
	OCR2TraceLogging() bool

	// OCR1 config, can override in jobs, only ethereum.
	OCRContractPollInterval() time.Duration
	OCRContractSubscribeInterval() time.Duration
	OCRContractTransmitterTransmitTimeout() time.Duration
	OCRBlockchainTimeout() time.Duration
	OCRDatabaseTimeout() time.Duration
	OCRMonitoringEndpoint() string
	OCRKeyBundleID() (string, error)
	OCRObservationGracePeriod() time.Duration
	OCRObservationTimeout() time.Duration
	OCRSimulateTransactions() bool
	OCRTransmitterAddress() (ethkey.EIP55Address, error) // OCR2 can support non-evm changes
	// OCR1 config, cannot override in jobs
	OCRTraceLogging() bool
	OCRDefaultTransactionQueueDepth() uint32
}

func (c *generalConfig) getDuration(field string) time.Duration {
	return c.getWithFallback(field, ParseDuration).(time.Duration)
}

func (c *generalConfig) OCR2ContractPollInterval() time.Duration {
	return c.getDuration("OCR2ContractPollInterval")
}

func (c *generalConfig) OCR2ContractSubscribeInterval() time.Duration {
	return c.getDuration("OCR2ContractSubscribeInterval")
}

func (c *generalConfig) OCR2ContractTransmitterTransmitTimeout() time.Duration {
	return c.getWithFallback("OCR2ContractTransmitterTransmitTimeout", ParseDuration).(time.Duration)
}

func (c *generalConfig) OCR2BlockchainTimeout() time.Duration {
	return c.getDuration("OCR2BlockchainTimeout")
}

func (c *generalConfig) OCR2DatabaseTimeout() time.Duration {
	return c.getWithFallback("OCR2DatabaseTimeout", ParseDuration).(time.Duration)
}

func (c *generalConfig) OCR2MonitoringEndpoint() string {
	return c.viper.GetString(EnvVarName("OCR2MonitoringEndpoint"))
}

func (c *generalConfig) OCR2KeyBundleID() (string, error) {
	kbStr := c.viper.GetString(EnvVarName("OCR2KeyBundleID"))
	if kbStr != "" {
		_, err := models.Sha256HashFromHex(kbStr)
		if err != nil {
			return "", errors.Wrapf(ErrInvalid, "OCR_KEY_BUNDLE_ID is an invalid sha256 hash hex string %v", err)
		}
	}
	return kbStr, nil
}

func (c *generalConfig) OCR2TraceLogging() bool {
	return c.viper.GetBool(EnvVarName("OCRTraceLogging"))
}

func (c *generalConfig) OCRContractPollInterval() time.Duration {
	return c.getDuration("OCRContractPollInterval")
}

func (c *generalConfig) OCRContractSubscribeInterval() time.Duration {
	return c.getDuration("OCRContractSubscribeInterval")
}

func (c *generalConfig) OCRContractTransmitterTransmitTimeout() time.Duration {
	return c.getWithFallback("OCRContractTransmitterTransmitTimeout", ParseDuration).(time.Duration)
}

func (c *generalConfig) OCRBlockchainTimeout() time.Duration {
	return c.getDuration("OCRBlockchainTimeout")
}

func (c *generalConfig) OCRDatabaseTimeout() time.Duration {
	return c.getWithFallback("OCRDatabaseTimeout", ParseDuration).(time.Duration)
}

func (c *generalConfig) OCRMonitoringEndpoint() string {
	return c.viper.GetString(EnvVarName("OCRMonitoringEndpoint"))
}

func (c *generalConfig) OCRKeyBundleID() (string, error) {
	kbStr := c.viper.GetString(EnvVarName("OCRKeyBundleID"))
	if kbStr != "" {
		_, err := models.Sha256HashFromHex(kbStr)
		if err != nil {
			return "", errors.Wrapf(ErrInvalid, "OCR_KEY_BUNDLE_ID is an invalid sha256 hash hex string %v", err)
		}
	}
	return kbStr, nil
}

// OCRDefaultTransactionQueueDepth controls the queue size for DropOldestStrategy in OCR
// Set to 0 to use SendEvery strategy instead
func (c *generalConfig) OCRDefaultTransactionQueueDepth() uint32 {
	return c.viper.GetUint32(EnvVarName("OCRDefaultTransactionQueueDepth"))
}

// OCRTraceLogging determines whether OCR logs at TRACE level are enabled. The
// option to turn them off is given because they can be very verbose
func (c *generalConfig) OCRTraceLogging() bool {
	return c.viper.GetBool(EnvVarName("OCRTraceLogging"))
}

func (c *generalConfig) OCRObservationTimeout() time.Duration {
	return c.getDuration("OCRObservationTimeout")
}

func (c *generalConfig) OCRObservationGracePeriod() time.Duration {
	return c.getWithFallback("OCRObservationGracePeriod", ParseDuration).(time.Duration)
}

// OCRSimulateTransactions enables using eth_call transaction simulation before
// sending when set to true
func (c *generalConfig) OCRSimulateTransactions() bool {
	return c.viper.GetBool(EnvVarName("OCRSimulateTransactions"))
}

func (c *generalConfig) OCRTransmitterAddress() (ethkey.EIP55Address, error) {
	taStr := c.viper.GetString(EnvVarName("OCRTransmitterAddress"))
	if taStr != "" {
		ta, err := ethkey.NewEIP55Address(taStr)
		if err != nil {
			return "", errors.Wrapf(ErrInvalid, "OCR_TRANSMITTER_ADDRESS is invalid EIP55 %v", err)
		}
		return ta, nil
	}
	return "", errors.Wrap(ErrUnset, "OCR_TRANSMITTER_ADDRESS env var is not set")
}

type P2PNetworking interface {
	P2PNetworkingStack() (n ocrnetworking.NetworkingStack)
	P2PNetworkingStackRaw() string
	P2PPeerID() p2pkey.PeerID
	P2PPeerIDRaw() string
	P2PIncomingMessageBufferSize() int
	P2POutgoingMessageBufferSize() int

	P2PV1Networking
	P2PV2Networking

	P2PDeprecated
}

// P2PNetworkingStack returns the preferred networking stack for libocr
func (c *generalConfig) P2PNetworkingStack() (n ocrnetworking.NetworkingStack) {
	str := c.P2PNetworkingStackRaw()
	err := n.UnmarshalText([]byte(str))
	if err != nil {
		logger.Fatalf("P2PNetworkingStack failed to unmarshal '%s': %s", str, err)
	}
	return n
}

// P2PNetworkingStackRaw returns the raw string passed as networking stack
func (c *generalConfig) P2PNetworkingStackRaw() string {
	return c.viper.GetString(EnvVarName("P2PNetworkingStack"))
}

// P2PPeerID is the default peer ID that will be used, if not overridden
func (c *generalConfig) P2PPeerID() p2pkey.PeerID {
	pidStr := c.viper.GetString(EnvVarName("P2PPeerID"))
	if pidStr == "" {
		return ""
	}
	var pid p2pkey.PeerID
	if err := pid.UnmarshalText([]byte(pidStr)); err != nil {
		logger.Error(errors.Wrapf(ErrInvalid, "P2P_PEER_ID is invalid %v", err))
		return ""
	}
	return pid
}

// P2PPeerIDRaw returns the string value of whatever P2P_PEER_ID was set to with no parsing
func (c *generalConfig) P2PPeerIDRaw() string {
	return c.viper.GetString(EnvVarName("P2PPeerID"))
}

func (c *generalConfig) P2PIncomingMessageBufferSize() int {
	if c.OCRIncomingMessageBufferSize() != 0 {
		return c.OCRIncomingMessageBufferSize()
	}
	return int(c.getWithFallback("P2PIncomingMessageBufferSize", ParseUint16).(uint16))
}

func (c *generalConfig) P2POutgoingMessageBufferSize() int {
	if c.OCROutgoingMessageBufferSize() != 0 {
		return c.OCRIncomingMessageBufferSize()
	}
	return int(c.getWithFallback("P2PIncomingMessageBufferSize", ParseUint16).(uint16))
}

type P2PDeprecated interface {
	// DEPRECATED - HERE FOR BACKWARDS COMPATABILITY
	OCRNewStreamTimeout() time.Duration
	OCRBootstrapCheckInterval() time.Duration
	OCRDHTLookupInterval() int
	OCRIncomingMessageBufferSize() int
	OCROutgoingMessageBufferSize() int
}

// DEPRECATED, do not use defaults, use only if specified and the
// newer env vars is not
func (c *generalConfig) OCRBootstrapCheckInterval() time.Duration {
	return c.viper.GetDuration("OCRBootstrapCheckInterval")
}

// DEPRECATED
func (c *generalConfig) OCRDHTLookupInterval() int {
	return c.viper.GetInt("OCRDHTLookupInterval")
}

// DEPRECATED
func (c *generalConfig) OCRNewStreamTimeout() time.Duration {
	return c.viper.GetDuration("OCRNewStreamTimeout")
}

// DEPRECATED
func (c *generalConfig) OCRIncomingMessageBufferSize() int {
	return c.viper.GetInt("OCRIncomingMessageBufferSize")
}

// DEPRECATED
func (c *generalConfig) OCROutgoingMessageBufferSize() int {
	return c.viper.GetInt("OCRIncomingMessageBufferSize")
}
