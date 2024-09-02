package generic

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	ocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type oracleFactory struct {
	database ocr3types.Database
	jobID    int32
	jobName  string
	jobORM   job.ORM
	kb       ocr2key.KeyBundle
	lggr     logger.Logger
}

type OracleFactoryParams struct {
	Database ocr3types.Database
	JobID    int32
	JobName  string
	JobORM   job.ORM
	Kb       ocr2key.KeyBundle
	Logger   logger.Logger
}

func NewOracleFactory(params OracleFactoryParams) (core.OracleFactory, error) {
	return &oracleFactory{
		database: params.Database,
		jobID:    params.JobID,
		jobName:  params.JobName,
		jobORM:   params.JobORM,
		kb:       params.Kb,
		lggr:     params.Logger,
	}, nil
}

func (of *oracleFactory) NewOracle(ctx context.Context, args core.OracleArgs) (core.Oracle, error) {
	// Could come from the capability spec config. Unsure about this as it feels wrong to expose implementation details of OCR config to the capability spec.
	traceLogging := false
	ocrLogger := ocrcommon.NewOCRWrapper(of.lggr, traceLogging, func(ctx context.Context, msg string) {
		logger.Sugared(of.lggr).ErrorIf(of.jobORM.RecordError(ctx, of.jobID, msg), "unable to record error")
	})

	// create the oracle from config values
	oracle, err := ocr.NewOracle(ocr.OCR3OracleArgs[[]byte]{
		LocalConfig:            args.LocalConfig,
		ContractConfigTracker:  args.ContractConfigTracker,
		ContractTransmitter:    args.ContractTransmitter,
		OffchainConfigDigester: args.OffchainConfigDigester,
		ReportingPluginFactory: args.ReportingPluginFactoryService,
		// BinaryNetworkEndpointFactory: d.cfg.BinaryNetworkEndpointFactory,
		// V2Bootstrappers:              d.cfg.V2Bootstrappers,
		Database: of.database,
		Logger:   ocrLogger,
		// TODO: This is Relayer Dependent - can we make it Relayer agnostic?
		MonitoringEndpoint: &telemetry.NoopAgent{},
		OffchainKeyring:    of.kb,
		OnchainKeyring:     ocrcommon.NewOCR3OnchainKeyringAdapter(of.kb),
		MetricsRegisterer:  prometheus.WrapRegistererWith(map[string]string{"job_name": of.jobName}, prometheus.DefaultRegisterer),
	})

	if err != nil {
		return nil, fmt.Errorf("%w: failed to create new OCR oracle", err)
	}

	return &adaptedOracle{oracle: oracle}, nil
}

type adaptedOracle struct {
	oracle ocr.Oracle
}

func (a *adaptedOracle) Start(ctx context.Context) error {
	return a.oracle.Start()
}

func (a *adaptedOracle) Close(ctx context.Context) error {
	return a.oracle.Close()
}
