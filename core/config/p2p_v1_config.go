package config

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
)

type P2PV1Networking interface {
	P2PAnnounceIP() net.IP
	P2PAnnouncePort() uint16
	P2PBootstrapPeers() ([]string, error)
	P2PDHTAnnouncementCounterUserPrefix() uint32
	P2PListenIP() net.IP
	P2PListenPort() uint16
	P2PListenPortRaw() string
	P2PNewStreamTimeout() time.Duration
	P2PBootstrapCheckInterval() time.Duration
	P2PDHTLookupInterval() int
	P2PPeerstoreWriteInterval() time.Duration
}

func (c *generalConfig) P2PPeerstoreWriteInterval() time.Duration {
	return c.getWithFallback("P2PPeerstoreWriteInterval", ParseDuration).(time.Duration)
}

func (c *generalConfig) P2PBootstrapPeers() ([]string, error) {
	if c.viper.IsSet(EnvVarName("P2PBootstrapPeers")) {
		bps := c.viper.GetStringSlice(EnvVarName("P2PBootstrapPeers"))
		if bps != nil {
			return bps, nil
		}
		return nil, errors.Wrap(ErrUnset, "P2P_BOOTSTRAP_PEERS env var is not set")
	}
	return []string{}, nil
}

// P2PListenIP is the ip that libp2p willl bind to and listen on
func (c *generalConfig) P2PListenIP() net.IP {
	return c.getWithFallback("P2PListenIP", ParseIP).(net.IP)
}

// P2PListenPort is the port that libp2p will bind to and listen on
func (c *generalConfig) P2PListenPort() uint16 {
	if c.viper.IsSet(EnvVarName("P2PListenPort")) {
		return uint16(c.viper.GetUint32(EnvVarName("P2PListenPort")))
	}
	// Fast path in case it was already set
	c.randomP2PPortMtx.RLock()
	if c.randomP2PPort > 0 {
		c.randomP2PPortMtx.RUnlock()
		return c.randomP2PPort
	}
	c.randomP2PPortMtx.RUnlock()
	// Path for initial set
	c.randomP2PPortMtx.Lock()
	defer c.randomP2PPortMtx.Unlock()
	if c.randomP2PPort > 0 {
		return c.randomP2PPort
	}
	r, err := rand.Int(rand.Reader, big.NewInt(65535-1023))
	if err != nil {
		panic(fmt.Errorf("unexpected error generating random port: %w", err))
	}
	randPort := uint16(r.Int64() + 1024)
	logger.Warnw(fmt.Sprintf("P2P_LISTEN_PORT was not set, listening on random port %d. A new random port will be generated on every boot, for stability it is recommended to set P2P_LISTEN_PORT to a fixed value in your environment", randPort), "p2pPort", randPort)
	c.randomP2PPort = randPort
	return c.randomP2PPort
}

// P2PListenPortRaw returns the raw string value of P2P_LISTEN_PORT
func (c *generalConfig) P2PListenPortRaw() string {
	return c.viper.GetString(EnvVarName("P2PListenPort"))
}

// P2PAnnounceIP is an optional override. If specified it will force the p2p
// layer to announce this IP as the externally reachable one to the DHT
// If this is set, P2PAnnouncePort MUST also be set.
func (c *generalConfig) P2PAnnounceIP() net.IP {
	str := c.viper.GetString(EnvVarName("P2PAnnounceIP"))
	return net.ParseIP(str)
}

// P2PAnnouncePort is an optional override. If specified it will force the p2p
// layer to announce this port as the externally reachable one to the DHT.
// If this is set, P2PAnnounceIP MUST also be set.
func (c *generalConfig) P2PAnnouncePort() uint16 {
	return uint16(c.viper.GetUint32(EnvVarName("P2PAnnouncePort")))
}

// P2PDHTAnnouncementCounterUserPrefix can be used to restore the node's
// ability to announce its IP/port on the P2P network after a database
// rollback. Make sure to only increase this value, and *never* decrease it.
// Don't use this variable unless you really know what you're doing, since you
// could semi-permanently exclude your node from the P2P network by
// misconfiguring it.
func (c *generalConfig) P2PDHTAnnouncementCounterUserPrefix() uint32 {
	return c.viper.GetUint32(EnvVarName("P2PDHTAnnouncementCounterUserPrefix"))
}

// FIXME: Add comments to all of these
func (c *generalConfig) P2PBootstrapCheckInterval() time.Duration {
	if c.OCRBootstrapCheckInterval() != 0 {
		return c.OCRBootstrapCheckInterval()
	}
	return c.getWithFallback("P2PBootstrapCheckInterval", ParseDuration).(time.Duration)
}

func (c *generalConfig) P2PDHTLookupInterval() int {
	if c.OCRDHTLookupInterval() != 0 {
		return c.OCRDHTLookupInterval()
	}
	return int(c.getWithFallback("P2PDHTLookupInterval", ParseUint16).(uint16))
}

func (c *generalConfig) P2PNewStreamTimeout() time.Duration {
	if c.OCRNewStreamTimeout() != 0 {
		return c.OCRNewStreamTimeout()
	}
	return c.getWithFallback("P2PNewStreamTimeout", ParseDuration).(time.Duration)
}
