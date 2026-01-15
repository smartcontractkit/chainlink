package automation

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "automation"}).Logger()

type Configurator struct {
	Automation *Automation `toml:"automation"`
}

type Automation struct {
	RegistryVersion  string           `toml:"registryVersion"`
	RegistrySettings RegistrySettings `toml:"RegistrySettings"`

	MercuryVersion *MercurySettings `toml:"MercuryVersion"`

	PluginConfig PluginConfig `toml:"PluginConfig"`
	PublicConfig PublicConfig `toml:"PublicConfig"`

	//TODO add fields from EVMConfigData
}

type MercurySettings struct {
	MercuryVersion *string `toml:"mercuryVersion"`
	FakeEndpoint   string  `toml:"fakeEndpoint"`
	FakePort       uint    `toml:"fakePort"`
}

type HeadTrackerData struct {
	HistoryDepth int
}

type GasEstimatorData struct {
	Mode         string
	LimitDefault int64
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

func NewAutomationConfigurator() *Configurator {
	return &Configurator{}
}

func (m *Configurator) Load() error {
	cfg, err := products.Load[Configurator]()
	if err != nil {
		return fmt.Errorf("failed to load product config: %w", err)
	}
	m.Automation = cfg.Automation
	return nil
}

func (m *Configurator) Store(path string, _ int) error {
	if err := products.Store(".", m); err != nil {
		return fmt.Errorf("failed to store product config: %w", err)
	}
	return nil
}

func (m *Configurator) GenerateCLNodesBlockchainConfig(ctx context.Context, bc *blockchain.Input) (string, error) {
	L.Info().Msg("Applying default CL nodes configuration")
	// configure node set and generate CL nodes configs
	config := `
	[Feature]
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
AnnounceAddresses = ['0.0.0.0:6690']`

	netConfigTemplate := `
[[EVM]]
AutoCreateKey = true
MinContractPayment = 0
BlockBackfillDepth = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.0000001 link'

ChainID = '{{.ChainID}}'
LogPollInterval = '{{.LogPollInterval}}'

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
		LogPollInterval     string
		ChainID             string
		WsUrl               string
		HttpUrl             string
		LinkContractAddress *string
		FinalityDepth       *int
		FinalityTag         *bool
		SafeTagSupported    *bool
		HeadTracker         *HeadTrackerData
		GasEstimator        *GasEstimatorData
	}

	d := data{
		LinkContractAddress: nil,
		ChainID:             bc.Out.ChainID,
		FinalityDepth:       nil, // TODO: set to &value if needed
		FinalityTag:         nil, // TODO: set to &value if needed
		WsUrl:               bc.Out.Nodes[0].InternalWSUrl,
		HttpUrl:             bc.Out.Nodes[0].InternalHTTPUrl,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, d); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	L.Info().Msg("Nodes network configuration is finished")
	return config + buf.String(), nil
}

func (m *Configurator) GenerateCLNodesSecrets(ctx context.Context) (string, error) {
	if m.Automation.MercuryVersion == nil {
		L.Info().Msg("Skipping CL nodes secrets configuration")
		return "", nil
	}

	L.Info().Msg("Applying default CL nodes secrets configuration")
	mercurySecretsTemplate := `
	[Mercury.Credentials.cred1]
	LegacyURL = '{{.URL}}'
	URL = '{{.URL}}'
	Username = 'node'
	Password = 'nodepass'`

	type data struct {
		URL string
	}

	u, err := url.JoinPath(framework.HostDockerInternal(), m.Automation.MercuryVersion.FakeEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to join URL path: %w", err)
	}

	d := data{
		URL: u,
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
	fake *fake.Input,
	bc *blockchain.Input,
	ns *nodeset.Input,
) error {
	L.Info().Msg("Connecting to CL nodes")
	cl, err := clclient.New(ns.Out.CLNodes)
	if err != nil {
		return err
	}
	pkey := getNetworkPrivateKey()
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
	c, auth, rootAddr, err := ETHClient(
		ctx,
		bcNode.ExternalWSUrl,
		m.Config[0].GasSettings.FeeCapMultiplier,
		m.Config[0].GasSettings.TipCapMultiplier,
	)
	if err != nil {
		return fmt.Errorf("could not create basic eth client: %w", err)
	}
	for _, addr := range ethKeyAddresses {
		if cErr := FundNodeEIP1559(ctx, c, pkey, addr, m.Config[0].CLNodesFundingETH); cErr != nil {
			return cErr
		}
	}
	ocrv2Config, ocr2Addr, err := m.configureContracts(
		ctx,
		c,
		auth,
		cl,
		rootAddr,
		transmitters,
		m.Config[0].CLNodesFundingLink,
	)
	if err != nil {
		return err
	}
	m.Config[0].OCR2SetConfigOut = ocrv2Config
	if cErr := m.configureJobs(ctx, fake, bc, ns, cl, ocr2Addr); cErr != nil {
		return cErr
	}
	r := resty.New().SetBaseURL(fake.Out.BaseURLHost)

	_, err = r.R().Post(`/trigger_deviation?result=200`)
	if err != nil {
		return fmt.Errorf("could not set ea fake values: %w", err)
	}
	L.Info().
		Msg("Setting fake external adapter (data feed) values")
	m.Config[0].DeployedContracts = &DeployedContracts{OCRv2AggregatorAddr: ocr2Addr}
	return nil
}
