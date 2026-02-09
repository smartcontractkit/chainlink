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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "automation"}).Logger()

type Configurator struct {
	Config []*Automation `toml:"automation"`
}

type Automation struct {
	RegistryVersion  string           `toml:"registry_version"`
	RegistrySettings RegistrySettings `toml:"registry_settings"`
	MercurySettings  *MercurySettings `toml:"mercury_settings"`

	PluginConfig PluginConfig `toml:"plugin_config"`
	PublicConfig PublicConfig `toml:"public_config"`

	CLNodesFundingETH  float64              `toml:"cl_nodes_funding_eth"`
	GasSettings        products.GasSettings `toml:"gas_settings"`
	DeployedContracts  DeployedContracts    `toml:"deployed_contracts"`
	EVMNetworkSettings EVMNetworkSettings   `toml:"evm_network_settings"`
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
	CredentialsName string `toml:"credentials_name"`
}

type EVMNetworkSettings struct {
	FinalityTagEnabled *bool `toml:"finality_tag_enabled"`
	FinalityDepth      *uint `toml:"finality_depth"`
	SafeTagSupported   *bool `toml:"safe_tag_supported"`

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
	LogProviderConfig    *LogProviderConfig `toml:"log_provider_config"`
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

{{- if .FinalityDepth}}
FinalityDepth = {{.FinalityDepth}}
{{- end}}
{{- if .FinalityTagEnabled}}
FinalityTagEnabled = {{.FinalityTagEnabled}}
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
WsUrl = '{{.WsURL}}'
HttpUrl = '{{.HTTPURL}}'
`

	tmpl, err := template.New("config").Parse(netConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	type data struct {
		LogPollInterval           *string
		BackupLogPollerBlockDelay *uint
		ChainID                   string
		WsURL                     string
		HTTPURL                   string
		FinalityDepth             *uint
		FinalityTagEnabled        *bool
		SafeTagSupported          *bool
		HeadTracker               *HeadTrackerData
		GasEstimator              *GasEstimatorData
	}

	d := data{
		ChainID:                   bc[0].Out.ChainID,
		FinalityDepth:             m.Config[0].EVMNetworkSettings.FinalityDepth,
		FinalityTagEnabled:        m.Config[0].EVMNetworkSettings.FinalityTagEnabled,
		SafeTagSupported:          m.Config[0].EVMNetworkSettings.SafeTagSupported,
		LogPollInterval:           m.Config[0].EVMNetworkSettings.LogPollerInterval,
		BackupLogPollerBlockDelay: m.Config[0].EVMNetworkSettings.BackupLogPollerBlockDelay,
		HeadTracker:               m.Config[0].EVMNetworkSettings.HeadTrackerData,
		GasEstimator:              m.Config[0].EVMNetworkSettings.GasEstimatorData,
		WsURL:                     bc[0].Out.Nodes[0].InternalWSUrl,
		HTTPURL:                   bc[0].Out.Nodes[0].InternalHTTPUrl,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, d); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	L.Info().Msg("Nodes network configuration is finished")
	return config + buf.String(), nil
}

func (m *Configurator) GenerateNodesSecrets(
	ctx context.Context,
	fs *fake.Input,
	_ []*blockchain.Input,
	_ []*nodeset.Input,
) (string, error) {
	if m.Config[0].MercurySettings == nil {
		L.Info().Msg("Product doesn't use Mercury. Skipping CL nodes secrets configuration")
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
		URL:             fs.Out.BaseURLDocker,
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

	if err := deployContracts(chainClient, m.Config[instanceIdx]); err != nil {
		return err
	}

	nodeDetails, err := collectNodeDetails(chainClient.Cfg.Network.ChainID, cl, ns[0].Out.CLNodes)
	if err != nil {
		return fmt.Errorf("error collecting node details: %w", err)
	}

	// it is crucial to create jobs before setting the config on the registry, as otherwise event emitted will not be picked up by the log poller
	if err := createJobs(cl, nodeDetails, int(chainClient.Cfg.Network.ChainID), m.Config[instanceIdx].MustGetRegistryVersion(), m.Config[instanceIdx].DeployedContracts.Registry, m.Config[instanceIdx].GetMercuryCredentialsName()); err != nil { //nolint:gosec // G115: chainID will never be big enough to overflow int
		return err
	}

	if err := waitForConfigWatcherToBeHealthy(cl); err != nil {
		return fmt.Errorf("failed to wait for ConfigWatcher health check: %w", err)
	}

	return setConfigOnRegistry(nodeDetails, m.Config[instanceIdx], chainClient)
}

func waitForConfigWatcherToBeHealthy(nodes []*clclient.ChainlinkClient) error {
	eg := &errgroup.Group{}
	for _, node := range nodes {
		eg.Go(func() error {
			return node.WaitHealthy(".*ConfigWatcher", "passing", 100)
		})
	}
	if waitErr := eg.Wait(); waitErr != nil {
		return fmt.Errorf("failed to wait for ConfigWatcher health check: %w", waitErr)
	}

	return nil
}

func (m *Automation) MustGetRegistryVersion() contracts.KeeperRegistryVersion {
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
	case version.Equal(semver.MustParse("2.0")):
		return contracts.RegistryVersion_2_0
	case version.Equal(semver.MustParse("2.1")):
		return contracts.RegistryVersion_2_1
	case version.Equal(semver.MustParse("2.2")):
		return contracts.RegistryVersion_2_2
	case version.Equal(semver.MustParse("2.3")):
		return contracts.RegistryVersion_2_3
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

func (m *Automation) GetRegistryConfig() contracts.KeeperRegistrySettings {
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
		RegistryVersion:      m.MustGetRegistryVersion(),
	}
}

func (m *Automation) GetPluginConfig() ocr2keepers30config.OffchainConfig {
	plCfg := m.PluginConfig
	return ocr2keepers30config.OffchainConfig{
		TargetProbability:    *plCfg.TargetProbability,
		TargetInRounds:       *plCfg.TargetInRounds,
		PerformLockoutWindow: *plCfg.PerformLockoutWindow,
		GasLimitPerReport:    *plCfg.GasLimitPerReport,
		GasOverheadPerUpkeep: *plCfg.GasOverheadPerUpkeep,
		MinConfirmations:     *plCfg.MinConfirmations,
		MaxUpkeepBatchSize:   *plCfg.MaxUpkeepBatchSize,
		LogProviderConfig: ocr2keepers30config.LogProviderConfig{
			BlockRate: *plCfg.LogProviderConfig.BlockRate,
			LogLimit:  *plCfg.LogProviderConfig.LogLimit,
		},
	}
}

func (m *Automation) GetPublicConfig() ocr3.PublicConfig {
	pubCfg := m.PublicConfig
	return ocr3.PublicConfig{
		DeltaProgress:                           *pubCfg.DeltaProgress,
		DeltaResend:                             *pubCfg.DeltaResend,
		DeltaInitial:                            *pubCfg.DeltaInitial,
		DeltaRound:                              *pubCfg.DeltaRound,
		DeltaGrace:                              *pubCfg.DeltaGrace,
		DeltaCertifiedCommitRequest:             *pubCfg.DeltaCertifiedCommitRequest,
		DeltaStage:                              *pubCfg.DeltaStage,
		RMax:                                    *pubCfg.RMax,
		MaxDurationQuery:                        *pubCfg.MaxDurationQuery,
		MaxDurationObservation:                  *pubCfg.MaxDurationObservation,
		MaxDurationShouldAcceptAttestedReport:   *pubCfg.MaxDurationShouldAcceptAttestedReport,
		MaxDurationShouldTransmitAcceptedReport: *pubCfg.MaxDurationShouldTransmitAcceptedReport,
		F:                                       *pubCfg.F,
	}
}
