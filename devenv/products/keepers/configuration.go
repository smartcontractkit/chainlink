package keepers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "keepers"}).Logger()

type Configurator struct {
	Config []*Keepers `toml:"keepers"`
}

type Keepers struct {
	RegistryVersion  string           `toml:"registry_version"`
	RegistrySettings RegistrySettings `toml:"registry_settings"`

	CLNodesFundingETH float64              `toml:"cl_nodes_funding_eth"`
	GasSettings       products.GasSettings `toml:"gas_settings"`

	DeployedContracts DeployedContracts `toml:"deployed_contracts"`
}

type DeployedContracts struct {
	LinkToken   string   `toml:"link_token"`
	LinkEthFeed string   `toml:"link_eth_feed"`
	EthGasFeed  string   `toml:"eth_gas_feed"`
	Transcoder  string   `toml:"transcoder"`
	Registry    string   `toml:"registry"`
	Registrar   string   `toml:"registrar"`
	Upkeeps     []string `toml:"upkeeps"`
}

type RegistrySettings struct {
	PaymentPremiumPPB    *uint32  `toml:"payment_premium_ppb"`
	FlatFeeMicroLINK     *uint32  `toml:"flat_fee_micro_link"`
	CheckGasLimit        *uint32  `toml:"check_gas_limit"`
	StalenessSeconds     *big.Int `toml:"staleness_seconds"`
	GasCeilingMultiplier *uint16  `toml:"gas_ceiling_multiplier"`
	MaxPerformGas        *uint32  `toml:"max_perform_gas"`
	MinUpkeepSpend       *big.Int `toml:"min_upkeep_spend"`
	FallbackGasPrice     *big.Int `toml:"fallback_gas_price"`
	FallbackLinkPrice    *big.Int `toml:"fallback_link_price"`
	FallbackNativePrice  *big.Int `toml:"fallback_native_price"`
	MaxCheckDataSize     *uint32  `toml:"max_check_data_size"`
	MaxPerformDataSize   *uint32  `toml:"max_perform_data_size"`
	MaxRevertDataSize    *uint32  `toml:"max_revert_data_size"`
	BlockCountPerTurn    *big.Int `toml:"block_count_per_turn"`
}

func (m *Keepers) GetRegistryConfig() contracts.KeeperRegistrySettings {
	registrySettings := m.RegistrySettings
	return contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    *registrySettings.PaymentPremiumPPB,
		FlatFeeMicroLINK:     *registrySettings.FlatFeeMicroLINK,
		CheckGasLimit:        *registrySettings.CheckGasLimit,
		StalenessSeconds:     registrySettings.StalenessSeconds,
		GasCeilingMultiplier: *registrySettings.GasCeilingMultiplier,
		MinUpkeepSpend:       registrySettings.MinUpkeepSpend,
		MaxPerformGas:        *registrySettings.MaxPerformGas,
		FallbackGasPrice:     registrySettings.FallbackGasPrice,
		FallbackLinkPrice:    registrySettings.FallbackLinkPrice,
		FallbackNativePrice:  registrySettings.FallbackNativePrice,
		MaxCheckDataSize:     *registrySettings.MaxCheckDataSize,
		MaxPerformDataSize:   *registrySettings.MaxPerformDataSize,
		MaxRevertDataSize:    *registrySettings.MaxRevertDataSize,
		BlockCountPerTurn:    registrySettings.BlockCountPerTurn,
		RegistryVersion:      m.MustGetRegistryVersion(),
	}
}

func NewConfigurator() *Configurator {
	return &Configurator{}
}

func (m *Configurator) Load() error {
	cfg, err := products.Load[Configurator]()
	if err != nil {
		return fmt.Errorf("failed to load product config: %w", err)
	}
	m.Config = cfg.Config
	return nil
}

func (m *Configurator) Store(path string, instanceIdx int) error {
	if err := products.Store(".", &Configurator{Config: []*Keepers{m.Config[instanceIdx]}}); err != nil {
		return fmt.Errorf("failed to store product config: %w", err)
	}
	return nil
}

func (m *Configurator) GenerateNodesConfig(
	ctx context.Context,
	fs *fake.Input,
	bc []*blockchain.Input,
	ns []*nodeset.Input,
) (string, error) {
	L.Info().Msg("Applying default CL nodes configuration")
	// configure node set and generate CL nodes configs
	config := `[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true

[Log]
Level = 'debug'
JSONConsole = true

[Log.File]
MaxSize = '0b'

[WebServer]
AllowOrigins = '*'
HTTPPort = 6688
SecureCookies = false
HTTPWriteTimeout = '3m'
SessionTimeout = '999h0m0s'

[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 1000

[WebServer.TLS]
HTTPSPort = 0

[Keeper]
TurnLookBack = 0
`

	netConfigTemplate := `
[[EVM]]
AutoCreateKey = true
MinContractPayment = 0
BlockBackfillDepth = 100
MinIncomingConfirmations = 1

ChainID = '{{.ChainID}}'

[[EVM.Nodes]]
Name = 'default'
WsUrl = '{{.WsURL}}'
HttpUrl = '{{.HTTPURL}}'
`

	tmpl, err := template.New("config").Parse(netConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	type data struct {
		ChainID string
		WsURL   string
		HTTPURL string
	}

	d := data{
		ChainID: bc[0].Out.ChainID,
		WsURL:   bc[0].Out.Nodes[0].InternalWSUrl,
		HTTPURL: bc[0].Out.Nodes[0].InternalHTTPUrl,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, d); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	L.Info().Msg("Nodes network configuration is finished")
	return config + buf.String(), nil
}

func (m *Configurator) GenerateNodesSecrets(
	_ context.Context,
	_ *fake.Input,
	_ []*blockchain.Input,
	_ []*nodeset.Input,
) (string, error) {
	return "", nil
}

func (m *Configurator) ConfigureJobsAndContracts(
	ctx context.Context,
	instanceIdx int,
	fs *fake.Input,
	bc []*blockchain.Input,
	ns []*nodeset.Input,
) error {
	L.Info().Msg("Connecting to CL nodes")
	cl, err := clclient.New(ns[0].Out.CLNodes)
	if err != nil {
		return err
	}
	pkey := products.NetworkPrivateKey()
	if pkey == "" {
		return errors.New("PRIVATE_KEY environment variable not set")
	}

	ethKeyAddresses := make([]string, 0)
	for i, nc := range cl {
		addr, cErr := nc.ReadPrimaryETHKey(bc[0].Out.ChainID)
		if cErr != nil {
			return cErr
		}
		ethKeyAddresses = append(ethKeyAddresses, addr.Attributes.Address)
		L.Info().
			Int("Idx", i).
			Str("ETH", addr.Attributes.Address).
			Msg("Node info")
	}

	bcNode := bc[0].Out.Nodes[0]
	c, _, _, err := products.ETHClient(
		ctx,
		bcNode.ExternalWSUrl,
		m.Config[instanceIdx].GasSettings.FeeCapMultiplier,
		m.Config[instanceIdx].GasSettings.TipCapMultiplier,
	)

	if err != nil {
		return fmt.Errorf("could not create basic eth client: %w", err)
	}
	for _, addr := range ethKeyAddresses {
		if cErr := products.FundAddressEIP1559(ctx, c, pkey, addr, m.Config[instanceIdx].CLNodesFundingETH); cErr != nil {
			return cErr
		}
	}

	chainID, err := strconv.ParseUint(bc[0].Out.ChainID, 10, 64)
	if err != nil {
		return err
	}

	chainClient, err := products.InitSeth(bcNode.ExternalWSUrl, []string{products.NetworkPrivateKey()}, &chainID)
	if err != nil {
		return err
	}

	if err := deployContracts(L, chainClient, m.Config[instanceIdx]); err != nil {
		return err
	}

	if err := createJobs(L, cl, m.Config[instanceIdx], bc[0].Out.ChainID); err != nil {
		return err
	}

	return setConfigOnRegistry(L, chainClient, cl, m.Config[instanceIdx])
}

func setConfigOnRegistry(l zerolog.Logger, chainClient *seth.Client, chainlinkNodes []*clclient.ChainlinkClient, config *Keepers) error {
	registry, err := contracts.LoadKeeperRegistry(L, chainClient, common.HexToAddress(config.DeployedContracts.Registry), config.MustGetRegistryVersion(), ZeroAddress)
	if err != nil {
		return err
	}

	primaryNode := chainlinkNodes[0]
	primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
	if err != nil {
		l.Error().Err(err).Msg("Reading ETH Keys from Chainlink Client shouldn't fail")
		return err
	}

	nodeAddresses := make([]string, 0)
	for _, clNode := range chainlinkNodes {
		clNodeAddress, err := clNode.PrimaryEthAddress()
		if err != nil {
			l.Error().Err(err).Msg("Error retrieving chainlink node address")
			return err
		}
		nodeAddresses = append(nodeAddresses, clNodeAddress)
	}

	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla)
		payees = append(payees, primaryNodeAddress)
	}

	err = registry.SetKeepers(nodeAddressesStr, payees, contracts.OCRv2Config{})
	if err != nil {
		l.Error().Err(err).Msg("Setting keepers in the registry shouldn't fail")
		return err
	}

	return nil
}

func (m *Keepers) MustGetRegistryVersion() contracts.KeeperRegistryVersion {
	version := semver.MustParse(m.RegistryVersion)
	switch {
	case version.Equal(semver.MustParse("1.0")):
		return contracts.RegistryVersion_1_0
	case version.Equal(semver.MustParse("1.1")):
		return contracts.RegistryVersion_1_1
	case version.Equal(semver.MustParse("1.2")):
		return contracts.RegistryVersion_1_2
	case version.Equal(semver.MustParse("1.3")):
		return contracts.RegistryVersion_1_3
	default:
		panic("unsupported registry version: " + m.RegistryVersion)
	}
}
