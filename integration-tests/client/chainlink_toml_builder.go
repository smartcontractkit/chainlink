package client

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

// TOMLBuilder helps facilitate cleanly building Chainlink settings TOMLs for different tests
type TOMLBuilder struct {
	tomlBuilder strings.Builder
}

// NewTOMLBuilder creates raw TOMLBuilder
func NewTOMLBuilder() *TOMLBuilder {
	return &TOMLBuilder{}
}

// NewDefaultTOMLBuilder shortcut to create TOML with defaults
func NewDefaultTOMLBuilder() *TOMLBuilder {
	return NewTOMLBuilder().AddGeneralDefaults()
}

// NewDefaultNetworksTOMLBuilder shortcut to create TOML with defaults and network settings
func NewDefaultNetworksTOMLBuilder(forwardingEnabled bool, networks ...*blockchain.EVMNetwork) *TOMLBuilder {
	return NewTOMLBuilder().AddGeneralDefaults().AddNetworks(forwardingEnabled, networks...)
}

// String builds the string value of the TOML to pass to config
func (t *TOMLBuilder) String() string {
	return t.tomlBuilder.String()
}

var defaultTOML = `RootDir = './clroot'

[Log]
Level = 'debug'
JSONConsole = true

[WebServer]
AllowOrigins = '*'
SecureCookies = false
SessionTimeout = '999h'

[WebServer.TLS]
HTTPSPort = 0`

// AddGeneralDefaults adds general testing defaults that are recommended for most tests
func (t *TOMLBuilder) AddGeneralDefaults() *TOMLBuilder {
	t.tomlBuilder.WriteString(fmt.Sprintf("\n%s\n", defaultTOML))
	return t
}

// AddNetworks adds TOML entries to connect the Chainlink node to provided networks
func (t *TOMLBuilder) AddNetworks(forwardingEnabled bool, networks ...*blockchain.EVMNetwork) *TOMLBuilder {
	for _, network := range networks {
		clNetwork, err := network.ChainlinkTOML(forwardingEnabled)
		if err != nil {
			log.Fatal().Err(err).Str("Network", network.Name).Msg("Error building network config for Chainlink TOML")
		}
		t.tomlBuilder.WriteString(fmt.Sprintf("\n%s\n", clNetwork))
	}
	return t
}

var p2pV1 = `[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690`

func (t *TOMLBuilder) AddP2PNetworkingV1() *TOMLBuilder {
	t.tomlBuilder.WriteString(fmt.Sprintf("\n%s\n", p2pV1))
	return t
}

var p2pV2 = `[P2P]
[P2P.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:6690']
AnnounceAddresses = ['0.0.0.0:6690']

[Feature]
LogPoller = true`

func (t *TOMLBuilder) AddP2PNetworkingV2() *TOMLBuilder {
	t.tomlBuilder.WriteString(fmt.Sprintf("\n%s\n", p2pV2))
	return t
}

var keeperDefaultTOML = `[Keeper]
TurnLookBack = 0
TurnFlagEnabled = true

[Keeper.Registry]
SyncInterval = '5s'
PerformGasOverhead = 150_000`

// AddKeeperDefaults adds default Keeper test settings
func (t *TOMLBuilder) AddKeeperDefaults() *TOMLBuilder {
	t.tomlBuilder.WriteString(fmt.Sprintf("\n%s\n", keeperDefaultTOML))
	return t
}

var ocrDefaultTOML = `[OCR]
Enabled = true`

// AddOCRDefaults adds default OCRv1 test settings
func (t *TOMLBuilder) AddOCRDefaults() *TOMLBuilder {
	t.tomlBuilder.WriteString(fmt.Sprintf("\n%s\n", ocrDefaultTOML))
	return t
}

var ocr2DefaultTOML = `[OCR2]
Enabled = true`

// AddOCR2Defaults adds default OCRv2 test settings
func (t *TOMLBuilder) AddOCR2Defaults() *TOMLBuilder {
	t.tomlBuilder.WriteString(fmt.Sprintf("\n%s\n", ocr2DefaultTOML))
	return t
}

// AddRaw adds a raw string to the TOML. Make sure it's properly formatted, or you'll see errors on the Chainlink node
func (t *TOMLBuilder) AddRaw(rawTOML string) *TOMLBuilder {
	t.tomlBuilder.WriteString(fmt.Sprintf("\n%s\n", rawTOML))
	return t
}
