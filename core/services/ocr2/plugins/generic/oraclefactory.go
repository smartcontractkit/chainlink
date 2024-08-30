package generic

import (
	"context"
	"fmt"

	ocr "github.com/smartcontractkit/libocr/offchainreporting2plus"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type oracleFactory struct{}

type adaptedOracle struct {
	oracle ocr.Oracle
}

func (a *adaptedOracle) Start(ctx context.Context) error {
	return a.oracle.Start()
}

func (a *adaptedOracle) Close(ctx context.Context) error {
	return a.oracle.Close()
}

func NewOracleFactory() (core.OracleFactory, error) {
	return &oracleFactory{}, nil
}

func (o *oracleFactory) NewOracle(ctx context.Context, args core.OracleArgs) (core.Oracle, error) {
	// create the oracle from config values
	oracle, err := ocr.NewOracle(ocr.OCR3OracleArgs[[]byte]{
		LocalConfig:            args.LocalConfig,
		ContractConfigTracker:  args.ContractConfigTracker,
		ContractTransmitter:    args.ContractTransmitter,
		OffchainConfigDigester: args.OffchainConfigDigester,
		ReportingPluginFactory: args.ReportingPluginFactoryService,
		// BinaryNetworkEndpointFactory: d.cfg.BinaryNetworkEndpointFactory,
		// V2Bootstrappers:              d.cfg.V2Bootstrappers,
		// Database:                     d.cfg.Database,
		// Logger:                       d.cfg.OCRLogger,
		// MonitoringEndpoint:           d.cfg.MonitoringEndpoint,
		// OffchainKeyring:              d.cfg.OffchainKeyring,
		// OnchainKeyring:               d.cfg.OnchainKeyring,
		// MetricsRegisterer: prometheus.WrapRegistererWith(map[string]string{"job_name": d.cfg.JobName.ValueOrZero()}, prometheus.DefaultRegisterer),
	})

	if err != nil {
		return nil, fmt.Errorf("%w: failed to create new OCR oracle", err)
	}

	return &adaptedOracle{oracle: oracle}, nil
}
