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
	return products.DefaultLegacyCLNodeConfig(bc)
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
		URL:         fs.Out.BaseURLDocker + "/cron_response",
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
