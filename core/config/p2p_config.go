package config

import (
	"time"

	"github.com/pkg/errors"

	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
)

// P2PNetworking is a subset of global config relevant to p2p networking.
type P2PNetworking interface {
	P2PNetworkingStack() (n ocrnetworking.NetworkingStack)
	P2PNetworkingStackRaw() string
	P2PPeerID() p2pkey.PeerID
	P2PPeerIDRaw() string
	P2PIncomingMessageBufferSize() int
	P2POutgoingMessageBufferSize() int
	P2PDeprecated
}

// P2PNetworkingStack returns the preferred networking stack for libocr
func (c *generalConfig) P2PNetworkingStack() (n ocrnetworking.NetworkingStack) {
	str := c.P2PNetworkingStackRaw()
	err := n.UnmarshalText([]byte(str))
	if err != nil {
		c.lggr.Panicf("P2PNetworkingStack failed to unmarshal '%s': %s", str, err)
	}
	return n
}

// P2PNetworkingStackRaw returns the raw string passed as networking stack
func (c *generalConfig) P2PNetworkingStackRaw() string {
	return c.viper.GetString(envvar.Name("P2PNetworkingStack"))
}

// P2PPeerID is the default peer ID that will be used, if not overridden
func (c *generalConfig) P2PPeerID() p2pkey.PeerID {
	pidStr := c.viper.GetString(envvar.Name("P2PPeerID"))
	if pidStr == "" {
		return ""
	}
	var pid p2pkey.PeerID
	if err := pid.UnmarshalText([]byte(pidStr)); err != nil {
		c.lggr.Critical(errors.Wrapf(ErrInvalid, "P2P_PEER_ID is invalid %v", err))
		return ""
	}
	return pid
}

// P2PPeerIDRaw returns the string value of whatever P2P_PEER_ID was set to with no parsing
func (c *generalConfig) P2PPeerIDRaw() string {
	return c.viper.GetString(envvar.Name("P2PPeerID"))
}

func (c *generalConfig) P2PIncomingMessageBufferSize() int {
	if c.OCRIncomingMessageBufferSize() != 0 {
		return c.OCRIncomingMessageBufferSize()
	}
	return int(getEnvWithFallback(c, envvar.NewUint16("P2PIncomingMessageBufferSize")))
}

func (c *generalConfig) P2POutgoingMessageBufferSize() int {
	if c.OCROutgoingMessageBufferSize() != 0 {
		return c.OCRIncomingMessageBufferSize()
	}
	return int(getEnvWithFallback(c, envvar.NewUint16("P2PIncomingMessageBufferSize")))
}

type P2PDeprecated interface {
	// DEPRECATED - HERE FOR BACKWARDS COMPATIBILITY
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
