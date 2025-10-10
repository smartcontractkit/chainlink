package utils

import (
	"os"
	"path"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/fakes"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type standardCapConfig struct {
	Config string

	// Set enabled to true to run the loop plugin.  Requires the plugin be installed.
	// Config will be passed to Initialise method of plugin.
	Enabled bool
}

var (
	goBinPath            = os.Getenv("GOBIN")
	standardCapabilities = map[string]standardCapConfig{
		"cron": {
			Config:  `{"fastestScheduleIntervalSeconds": 1}`,
			Enabled: true,
		},
		"readcontract":  {},
		"kvstore":       {},
		"workflowevent": {},
	}
)

// standaloneLoopWrapper wraps a StandardCapabilities to implement services.Service
type standaloneLoopWrapper struct {
	*standardcapabilities.StandardCapabilities
}

func (l *standaloneLoopWrapper) Ready() error { return l.StandardCapabilities.Ready() }

func (l *standaloneLoopWrapper) HealthReport() map[string]error { return make(map[string]error) }

func (l *standaloneLoopWrapper) Name() string { return "wrapped" }

func newStandardCapabilities(
	standardCapabilities map[string]standardCapConfig,
	lggr logger.Logger,
	registry *capabilities.Registry,
) []services.Service {
	caps := make([]services.Service, 0)

	pluginRegistrar := plugins.NewRegistrarConfig(
		loop.GRPCOpts{},
		func(name string) (*plugins.RegisteredLoop, error) { return &plugins.RegisteredLoop{}, nil },
		func(loopId string) {})

	for name, config := range standardCapabilities {
		if !config.Enabled {
			continue
		}

		spec := &job.StandardCapabilitiesSpec{
			Command: path.Join(goBinPath, name),
			Config:  config.Config,
		}

		loop := standardcapabilities.NewStandardCapabilities(lggr, spec,
			pluginRegistrar, core.StandardCapabilitiesDependencies{
				Config:             spec.Config,
				TelemetryService:   &fakes.TelemetryServiceMock{},
				Store:              &fakes.KVStoreMock{},
				CapabilityRegistry: registry,
				ErrorLog:           &fakes.ErrorLogMock{},
				PipelineRunner:     &fakes.PipelineRunnerServiceMock{},
				RelayerSet:         &fakes.RelayerSetMock{},
				OracleFactory:      &fakes.OracleFactoryMock{},
				GatewayConnector:   &fakes.GatewayConnectorMock{},
				P2PKeystore:        &fakes.KeystoreMock{},
			})

		service := &standaloneLoopWrapper{
			StandardCapabilities: loop,
		}
		caps = append(caps, service)
	}

	return caps
}
