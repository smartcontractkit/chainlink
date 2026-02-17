package cron

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "cron"}).Logger()

type Configurator struct {
	Config []*Cron `toml:"cron"`
}

type Cron struct {
	Schedule string `toml:"schedule" comment:"Cron schedule string in format: 'CRON_TZ=UTC * * * * * *'"`
	Out      *Out   `toml:"out"`
}

type Out struct {
	JobID string `toml:"job_id"`
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
	if err := products.Store(".", m); err != nil {
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
	node := bc[0].Out.Nodes[0]
	chainID := bc[0].Out.ChainID
	netConfig := fmt.Sprintf(`
    [[EVM]]
    LogPollInterval = '1s'
    BlockBackfillDepth = 100
    ChainID = '%s'
    MinIncomingConfirmations = 1
    MinContractPayment = '0.0000001 link'

    [[EVM.Nodes]]
    Name = 'default'
    WsUrl = '%s'
    HttpUrl = '%s'

    [Feature]
    FeedsManager = true
    LogPoller = true
    UICSAKeys = true
    [OCR2]
    Enabled = true
    SimulateTransactions = false
    DefaultTransactionQueueDepth = 1
    [P2P.V2]
    Enabled = true
    ListenAddresses = ['0.0.0.0:6690']
	[Log]
	JSONConsole = true
	Level = 'debug'
	[Pyroscope]
	ServerAddress = 'http://host.docker.internal:4040'
	Environment = 'local'
	[WebServer]
	SessionTimeout = '999h0m0s'
	HTTPWriteTimeout = '3m'
	SecureCookies = false
	HTTPPort = 6688
	[WebServer.TLS]
	HTTPSPort = 0
	[WebServer.RateLimit]
	Authenticated = 5000
	Unauthenticated = 5000
	[JobPipeline]
	[JobPipeline.HTTPRequest]
	DefaultTimeout = '1m'
	[Log.File]
	MaxSize = '0b'
`,
		chainID,
		node.InternalWSUrl,
		node.InternalHTTPUrl,
	)
	L.Info().Msg("Nodes network configuration is finished")
	return netConfig, nil
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
	cls, err := clclient.New(ns[0].Out.CLNodes)
	if err != nil {
		return err
	}
	L.Info().Msg("Creating bridge and cron schedule")
	bta := &clclient.BridgeTypeAttributes{
		Name:        "cron-" + uuid.NewString(),
		URL:         fmt.Sprintf("%s/cron_response", fs.Out.BaseURLDocker),
		RequestData: "{}",
	}
	if err := cls[0].MustCreateBridge(bta); err != nil {
		return fmt.Errorf("failed to create bridge: %w", err)
	}
	j, err := cls[0].MustCreateJob(&clclient.CronJobSpec{
		Schedule:          m.Config[0].Schedule,
		ObservationSource: clclient.ObservationSourceSpecBridge(bta),
	})
	if err != nil {
		return fmt.Errorf("failed to create cron job: %w", err)
	}
	m.Config[0].Out = &Out{
		JobID: j.Data.ID,
	}
	return nil
}
