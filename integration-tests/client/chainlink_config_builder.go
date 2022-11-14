package client

import (
	"fmt"
	"net"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/core/assets"
	v2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	chainlinkConfig "github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ChainlinkConfigBuilder struct {
	config *chainlinkConfig.Config `toml:"-"`
}

// NewDefaultConfig gets the base config for
func NewDefaultConfig() *ChainlinkConfigBuilder {
	conf := &chainlinkConfig.Config{}
	conf.RootDir = Pointer("./clroot")
	conf.DevMode = false

	conf.Log.JSONConsole = Pointer(true)

	conf.WebServer.AllowOrigins = Pointer("*")
	conf.WebServer.SecureCookies = Pointer(false)
	conf.WebServer.SessionTimeout = models.MustNewDuration(time.Hour * 999)
	conf.WebServer.TLS.HTTPSPort = Pointer(uint16(0))
	return &ChainlinkConfigBuilder{conf}
}

// AddNetworks adds networks to the config, enabling forwarders for each network if specified
func (c *ChainlinkConfigBuilder) AddNetworks(
	enableForwarders bool,
	networks ...*blockchain.EVMNetwork,
) *ChainlinkConfigBuilder {
	for _, network := range networks {
		nodes := []*v2.Node{}
		for nodeId := range network.URLs {
			wsURL, err := models.ParseURL(network.URLs[nodeId])
			if err != nil {
				log.Fatal().Str("URL", network.URLs[nodeId]).Err(err).Msg("Error in URL formatting")
			}
			httpURL, err := models.ParseURL(network.HTTPURLs[nodeId])
			if err != nil {
				log.Fatal().Str("URL", network.HTTPURLs[nodeId]).Err(err).Msg("Error in URL formatting")
			}
			nodes = append(nodes, &v2.Node{
				Name:    Pointer(fmt.Sprintf("node-%d", nodeId)),
				WSURL:   wsURL,
				HTTPURL: httpURL,
			})
		}

		c.config.EVM = append(c.config.EVM, &v2.EVMConfig{
			ChainID: utils.NewBigI(network.ChainID),
			Enabled: Pointer(true),
			Nodes:   nodes,
			Chain: v2.Chain{
				MinContractPayment: assets.NewLinkFromJuels(0),
				Transactions: v2.Transactions{
					ForwardersEnabled: Pointer(enableForwarders),
				},
			},
		})
	}
	return c
}

// AddP2PNetworkingV1 adds defaults for V1 P2P networking
func (c *ChainlinkConfigBuilder) AddP2PNetworkingV1() *ChainlinkConfigBuilder {
	c.config.P2P.V1.Enabled = Pointer(true)
	c.config.P2P.V1.ListenIP = &net.IPv4zero
	c.config.P2P.V1.ListenPort = Pointer(uint16(6690))
	return c
}

// AddP2PNetworkingV2 adds defaults for V2 P2P networking (also enables LogPoller)
func (c *ChainlinkConfigBuilder) AddP2PNetworkingV2() *ChainlinkConfigBuilder {
	c.config.P2P.V2.Enabled = Pointer(true)
	c.config.P2P.V2.ListenAddresses = &[]string{"0.0.0.0:6690"}
	c.config.P2P.V2.AnnounceAddresses = &[]string{"0.0.0.0:6690"}

	c.config.Feature.LogPoller = Pointer(true)
	return c
}

// AddKeeperDefaults enables default testing behavior for Keepers
func (c *ChainlinkConfigBuilder) AddKeeperDefaults() *ChainlinkConfigBuilder {
	c.config.Keeper.TurnFlagEnabled = Pointer(true)
	c.config.Keeper.TurnLookBack = Pointer(int64(0))
	c.config.Keeper.Registry.SyncInterval = models.MustNewDuration(time.Second * 5)
	c.config.Keeper.Registry.PerformGasOverhead = Pointer(uint32(150_000))
	return c
}

// AddOCRDefaults enables OCR functionality
func (c *ChainlinkConfigBuilder) AddOCRDefaults() *ChainlinkConfigBuilder {
	c.config.OCR.Enabled = Pointer(true)
	return c
}

// AddOCR2Defaults enables OCR2 functionality
func (c *ChainlinkConfigBuilder) AddOCR2Defaults() *ChainlinkConfigBuilder {
	c.config.OCR2.Enabled = Pointer(true)
	return c
}

// MustTOML marshals the config to a TOML string. Will fail if there is an error.
func (c *ChainlinkConfigBuilder) MustTOML() string {
	rawTOML, err := toml.Marshal(c.config)
	if err != nil {
		log.Fatal().Err(err).Msg("Error marshalling config TOML")
	}
	return fmt.Sprintf("%s", rawTOML)
}

// RawConfig returns the raw config object so you can modify it directly
func (c *ChainlinkConfigBuilder) RawConfig() *chainlinkConfig.Config {
	return c.config
}

// Pointer converts normal types to a pointer type, for use in building config values
func Pointer[T any](v T) *T { return &v }
