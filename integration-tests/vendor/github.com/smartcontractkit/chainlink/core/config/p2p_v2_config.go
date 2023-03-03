package config

import (
	"time"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/config/parse"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// P2PV2Networking is a subset of global config relevant to p2p v2 networking.
type P2PV2Networking interface {
	P2PV2AnnounceAddresses() []string
	P2PV2Bootstrappers() (locators []ocrcommontypes.BootstrapperLocator)
	P2PV2BootstrappersRaw() []string
	P2PV2DeltaDial() models.Duration
	P2PV2DeltaReconcile() models.Duration
	P2PV2ListenAddresses() []string
}

// P2PV2ListenAddresses contains the addresses the peer will listen to on the network in <host>:<port> form as
// accepted by net.Listen, but host and port must be fully specified and cannot be empty.
func (c *generalConfig) P2PV2ListenAddresses() []string {
	return c.viper.GetStringSlice(envvar.Name("P2PV2ListenAddresses"))
}

// P2PV2AnnounceAddresses contains the addresses the peer will advertise on the network in <host>:<port> form as
// accepted by net.Dial. The addresses should be reachable by peers of interest.
func (c *generalConfig) P2PV2AnnounceAddresses() []string {
	if c.viper.IsSet(envvar.Name("P2PV2AnnounceAddresses")) {
		return c.viper.GetStringSlice(envvar.Name("P2PV2AnnounceAddresses"))
	}
	return c.P2PV2ListenAddresses()
}

// P2PV2Bootstrappers returns the default bootstrapper peers for libocr's v2
// networking stack
func (c *generalConfig) P2PV2Bootstrappers() (locators []ocrcommontypes.BootstrapperLocator) {
	bootstrappers := c.P2PV2BootstrappersRaw()
	for _, s := range bootstrappers {
		var locator ocrcommontypes.BootstrapperLocator
		err := locator.UnmarshalText([]byte(s))
		if err != nil {
			c.lggr.Panicf("invalid format for bootstrapper '%s', got error: %s", s, err)
		}
		locators = append(locators, locator)
	}
	return
}

// P2PV2BootstrappersRaw returns the raw strings for v2 bootstrap peers
func (c *generalConfig) P2PV2BootstrappersRaw() []string {
	return c.viper.GetStringSlice(envvar.Name("P2PV2Bootstrappers"))
}

// P2PV2DeltaDial controls how far apart Dial attempts are
func (c *generalConfig) P2PV2DeltaDial() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("P2PV2DeltaDial", parse.Duration).(time.Duration))
}

// P2PV2DeltaReconcile controls how often a Reconcile message is sent to every peer.
func (c *generalConfig) P2PV2DeltaReconcile() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("P2PV2DeltaReconcile", parse.Duration).(time.Duration))
}
