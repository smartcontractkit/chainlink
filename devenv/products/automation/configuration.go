package automation

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/devenv/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "automation"}).Logger()

type Configurator struct {
	Config []*Automation `toml:"automation"`
}

type Automation struct {
	RegistryVersion  string           `toml:"registryVersion"`
	RegistrySettings RegistrySettings `toml:"RegistrySettings"`

	MercurySettings *MercurySettings `toml:"MercurySettings"`

	PluginConfig PluginConfig `toml:"PluginConfig"`
	PublicConfig PublicConfig `toml:"PublicConfig"`

	CLNodesFundingETH float64 `toml:"cl_nodes_funding_eth"`
	// CLNodesFundingLink float64 `toml:"cl_nodes_funding_link"`

	GasSettings *products.GasSettings `toml:"gas_settings"`

	DeployedContracts DeployedContracts `toml:"deployed_contracts"`

	EVMNetworkSettings EVMNetworkSettings `toml:"evm_network_settings"`

	TestKeysMinFundingEth float64 `toml:"test_keys_min_funding_eth"`
}

type DeployedContracts struct {
	LinkToken   string   `toml:"link_token"`
	Weth        string   `toml:"weth"`
	LinkEthFeed string   `toml:"link_eth_feed"`
	EthGasFeed  string   `toml:"eth_gas_feed"`
	EthUSDFeed  string   `toml:"eth_usd_feed"`
	LinkUSDFeed string   `toml:"link_usd_feed"`
	Transcoder  string   `toml:"transcoder"`
	ChainModule string   `toml:"chain_module"`
	Registry    string   `toml:"registry"`
	Registrar   string   `toml:"registrar"`
	MultiCall   string   `toml:"multi_call"`
	Upkeeps     []string `toml:"upkeeps"`
}

type MercurySettings struct {
	Version         string `toml:"version"`
	CredentialsName string `toml:"credentialsName"`
}

type EVMNetworkSettings struct {
	LinkTokenAddress   *string `toml:"link_token_address"`
	FinalityTagEnabled *bool   `toml:"finality_tag_enabled"`
	FinalityDepth      *uint   `toml:"finality_depth"`
	SafeTagSupported   *bool   `toml:"safe_tag_supported"`

	BackupLogPollerBlockDelay *uint   `toml:"backup_log_poller_block_delay"`
	LogPollerInterval         *string `toml:"log_poller_interval"`

	HeadTrackerData  *HeadTrackerData  `toml:"head_tracker"`
	GasEstimatorData *GasEstimatorData `toml:"gas_estimator"`
}

type HeadTrackerData struct {
	HistoryDepth int `toml:"history_depth"`
}

type GasEstimatorData struct {
	Mode         string `toml:"mode"`
	LimitDefault int64  `toml:"limit_default"`
}

type PluginConfig struct {
	PerformLockoutWindow *int64             `toml:"perform_lockout_window"`
	TargetProbability    *string            `toml:"target_probability"`
	TargetInRounds       *int               `toml:"target_in_rounds"`
	MinConfirmations     *int               `toml:"min_confirmations"`
	GasLimitPerReport    *uint32            `toml:"gas_limit_per_report"`
	GasOverheadPerUpkeep *uint32            `toml:"gas_overhead_per_upkeep"`
	MaxUpkeepBatchSize   *int               `toml:"max_upkeep_batch_size"`
	LogProviderConfig    *LogProviderConfig `toml:"LogProviderConfig"`
}

type LogProviderConfig struct {
	BlockRate *uint32 `toml:"block_rate"`
	LogLimit  *uint32 `toml:"log_limit"`
}

type PublicConfig struct {
	DeltaProgress                           *time.Duration `toml:"delta_progress"`
	DeltaResend                             *time.Duration `toml:"delta_resend"`
	DeltaInitial                            *time.Duration `toml:"delta_initial"`
	DeltaRound                              *time.Duration `toml:"delta_round"`
	DeltaGrace                              *time.Duration `toml:"delta_grace"`
	DeltaCertifiedCommitRequest             *time.Duration `toml:"delta_certified_commit_request"`
	DeltaStage                              *time.Duration `toml:"delta_stage"`
	RMax                                    *uint64        `toml:"r_max"`
	F                                       *int           `toml:"f"`
	MaxDurationQuery                        *time.Duration `toml:"max_duration_query"`
	MaxDurationObservation                  *time.Duration `toml:"max_duration_observation"`
	MaxDurationShouldAcceptAttestedReport   *time.Duration `toml:"max_duration_should_accept_attested_report"`
	MaxDurationShouldTransmitAcceptedReport *time.Duration `toml:"max_duration_should_transmit_accepted_report"`
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
	if err := products.Store(".", &Configurator{Config: []*Automation{m.Config[instanceIdx]}}); err != nil {
		return fmt.Errorf("failed to store product config: %w", err)
	}
	return nil
}

func (m *Configurator) GenerateCLNodesBlockchainConfig(ctx context.Context, bc *blockchain.Input) (string, error) {
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

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:6690']
AnnounceAddresses = ['0.0.0.0:6690']
`

	netConfigTemplate := `
[[EVM]]
AutoCreateKey = true
MinContractPayment = 0
BlockBackfillDepth = 100
MinIncomingConfirmations = 1

ChainID = '{{.ChainID}}'
{{- if .LogPollInterval}}
LogPollInterval = '{{.LogPollInterval}}'
{{- end}}

{{- if .BackupLogPollerBlockDelay}}
BackupLogPollerBlockDelay = {{.BackupLogPollerBlockDelay}}
{{- end}}
{{- if .LinkContractAddress}}
LinkContractAddress = '{{.LinkContractAddress}}'
{{- end}}

{{- if .FinalityDepth}}
FinalityDepth = {{.FinalityDepth}}
{{- end}}
{{- if .FinalityTag}}
FinalityTag = {{.FinalityTag}}
{{- end}}
{{- if .SafeTagSupported}}
SafeTagSupported = {{.SafeTagSupported}}
{{- end}}

{{- if .HeadTracker}}
[HeadTracker]
HistoryDepth = {{.HeadTracker.HistoryDepth}}
{{- end}}
{{- if .GasEstimator}}
[GasEstimator]
Mode = '{{.GasEstimator.Mode}}'
LimitDefault = {{.GasEstimator.LimitDefault}}
{{- end}}

[[EVM.Nodes]]
Name = 'default'
WsUrl = '{{.WsUrl}}'
HttpUrl = '{{.HttpUrl}}'
`

	tmpl, err := template.New("config").Parse(netConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	type data struct {
		LogPollInterval           *string
		BackupLogPollerBlockDelay *uint
		ChainID                   string
		WsUrl                     string
		HttpUrl                   string
		LinkContractAddress       *string
		FinalityDepth             *uint
		FinalityTag               *bool
		SafeTagSupported          *bool
		HeadTracker               *HeadTrackerData
		GasEstimator              *GasEstimatorData
	}

	d := data{
		LinkContractAddress:       m.Config[0].EVMNetworkSettings.LinkTokenAddress, // TODO think whether we need and how to set if it is deployed later. Is the sequence deterministic enough?
		ChainID:                   bc.Out.ChainID,
		FinalityDepth:             m.Config[0].EVMNetworkSettings.FinalityDepth,
		FinalityTag:               m.Config[0].EVMNetworkSettings.FinalityTagEnabled,
		SafeTagSupported:          m.Config[0].EVMNetworkSettings.SafeTagSupported,
		LogPollInterval:           m.Config[0].EVMNetworkSettings.LogPollerInterval,
		BackupLogPollerBlockDelay: m.Config[0].EVMNetworkSettings.BackupLogPollerBlockDelay,
		HeadTracker:               m.Config[0].EVMNetworkSettings.HeadTrackerData,
		GasEstimator:              m.Config[0].EVMNetworkSettings.GasEstimatorData,
		WsUrl:                     bc.Out.Nodes[0].InternalWSUrl,
		HttpUrl:                   bc.Out.Nodes[0].InternalHTTPUrl,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, d); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	L.Info().Msg("Nodes network configuration is finished")
	return config + buf.String(), nil
}

func (m *Configurator) GenerateCLNodesSecrets(ctx context.Context, fake *fake.Input) (string, error) {
	if m.Config[0].MercurySettings == nil {
		L.Info().Msg("Skipping CL nodes secrets configuration")
		return "", nil
	}

	L.Info().Msg("Applying default CL nodes secrets configuration")
	mercurySecretsTemplate := `
	[Mercury.Credentials.{{.CredentialsName}}]
	LegacyURL = '{{.URL}}'
	URL = '{{.URL}}'
	Username = 'node'
	Password = 'nodepass'`

	type data struct {
		CredentialsName string
		URL             string
	}

	d := data{
		URL:             fake.Out.BaseURLDocker,
		CredentialsName: m.Config[0].MercurySettings.CredentialsName,
	}

	tmpl, err := template.New("secrets").Parse(mercurySecretsTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, d); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	L.Info().Msg("Nodes secrets configuration is finished")
	return buf.String(), nil
}

func (m *Configurator) ConfigureJobsAndContracts(
	ctx context.Context,
	instanceIdx int,
	fake *fake.Input,
	bc *blockchain.Input,
	ns *nodeset.Input,
) error {
	L.Info().Msg("Connecting to CL nodes")
	cl, err := clclient.New(ns.Out.CLNodes)
	if err != nil {
		return err
	}
	pkey := products.NetworkPrivateKey()
	if pkey == "" {
		return errors.New("PRIVATE_KEY environment variable not set")
	}

	transmitters := make([]common.Address, 0)
	ethKeyAddresses := make([]string, 0)
	for i, nc := range cl {
		addr, cErr := nc.ReadPrimaryETHKey(bc.Out.ChainID)
		if cErr != nil {
			return cErr
		}
		ethKeyAddresses = append(ethKeyAddresses, addr.Attributes.Address)
		transmitters = append(transmitters, common.HexToAddress(addr.Attributes.Address))
		L.Info().
			Int("Idx", i).
			Str("ETH", addr.Attributes.Address).
			Msg("Node info")
	}

	bcNode := bc.Out.Nodes[0]
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
		if cErr := products.FundNodeEIP1559(ctx, c, pkey, addr, m.Config[instanceIdx].CLNodesFundingETH); cErr != nil {
			return cErr
		}
	}

	chainID, err := strconv.ParseUint(bc.ChainID, 10, 64)
	if err != nil {
		return err
	}

	var chainClient *seth.Client
	if os.Getenv(seth.CONFIG_FILE_ENV_VAR) != "" {
		sethCfg, err := seth.ReadConfig()
		if err != nil {
			return err
		}

		chainClient, err = seth.NewClientBuilderWithConfig(sethCfg).
			UseNetworkWithChainId(chainID).
			WithRpcUrl(bc.Out.Nodes[0].ExternalWSUrl).
			Build()
	} else {
		chainClient, err = seth.NewClientBuilder().
			WithPrivateKeys([]string{products.NetworkPrivateKey()}).
			WithRpcUrl(bc.Out.Nodes[0].ExternalWSUrl).
			Build()
	}
	if err != nil {
		return err
	}

	if err := deployContracts(chainClient, m.Config[instanceIdx]); err != nil {
		return err
	}

	nodeDetails, err := CollectNodeDetails(chainClient.Cfg.Network.ChainID, cl, ns.Out.CLNodes)
	if err != nil {
		return fmt.Errorf("error collecting node details: %w", err)
	}

	if err := SetConfigOnRegistry(nodeDetails, m.Config[instanceIdx], chainClient); err != nil {
		return err
	}

	return createJobs(cl, nodeDetails, int(chainClient.Cfg.Network.ChainID), m.Config[instanceIdx].MustGetRegistryVersion(), m.Config[instanceIdx].DeployedContracts.Registry, m.Config[instanceIdx].GetMercuryCredentialsName())
}

func (m *Automation) MustGetRegistryVersion() ethereum.KeeperRegistryVersion {
	version := semver.MustParse(m.RegistryVersion)
	switch {
	case version.Equal(semver.MustParse("1.0")):
		return ethereum.RegistryVersion_1_0
	case version.Equal(semver.MustParse("1.1")):
		return ethereum.RegistryVersion_1_1
	case version.Equal(semver.MustParse("1.2")):
		return ethereum.RegistryVersion_1_2
	case version.Equal(semver.MustParse("1.3")):
		return ethereum.RegistryVersion_1_3
	case version.Equal(semver.MustParse("2.0")):
		return ethereum.RegistryVersion_2_0
	case version.Equal(semver.MustParse("2.1")):
		return ethereum.RegistryVersion_2_1
	case version.Equal(semver.MustParse("2.2")):
		return ethereum.RegistryVersion_2_2
	case version.Equal(semver.MustParse("2.3")):
		return ethereum.RegistryVersion_2_3
	default:
		panic("unsupported registry version: " + m.RegistryVersion)
	}
}

func (m *Automation) GetMercuryCredentialsName() string {
	if m.MercurySettings != nil {
		return m.MercurySettings.CredentialsName
	}

	return ""
}
