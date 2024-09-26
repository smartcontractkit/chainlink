package generic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	ocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type oracleFactoryConfig struct {
	Enabled        bool
	TraceLogging   bool
	BootstrapPeers []commontypes.BootstrapperLocator
}

func NewOracleFactoryConfig(config job.JSONConfig) (*oracleFactoryConfig, error) {
	var ofc struct {
		Enabled        bool     `json:"enabled"`
		TraceLogging   bool     `json:"traceLogging"`
		BootstrapPeers []string `json:"bootstrapPeers"`
	}
	err := json.Unmarshal(config.Bytes(), &ofc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal oracle factory config")
	}

	if !ofc.Enabled {
		return &oracleFactoryConfig{}, nil
	}

	// If Oracle Factory is enabled, it must have at least one bootstrap peer
	if len(ofc.BootstrapPeers) == 0 {
		return nil, errors.New("no bootstrap peers found")
	}

	bootstrapPeers, err := ocrcommon.ParseBootstrapPeers(ofc.BootstrapPeers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse bootstrap peers")
	}

	return &oracleFactoryConfig{
		Enabled:        true,
		TraceLogging:   ofc.TraceLogging,
		BootstrapPeers: bootstrapPeers,
	}, nil
}

type oracleFactory struct {
	database    ocr3types.Database
	jobID       int32
	jobName     string
	jobORM      job.ORM
	kb          ocr2key.KeyBundle
	lggr        logger.Logger
	config      *oracleFactoryConfig
	peerWrapper *ocrcommon.SingletonPeerWrapper
}

type OracleFactoryParams struct {
	Database    ocr3types.Database
	JobID       int32
	JobName     string
	JobORM      job.ORM
	Kb          ocr2key.KeyBundle
	Logger      logger.Logger
	Config      *oracleFactoryConfig
	PeerWrapper *ocrcommon.SingletonPeerWrapper
}

func NewOracleFactory(params OracleFactoryParams) (core.OracleFactory, error) {
	return &oracleFactory{
		database:    params.Database,
		jobID:       params.JobID,
		jobName:     params.JobName,
		jobORM:      params.JobORM,
		kb:          params.Kb,
		lggr:        params.Logger,
		config:      params.Config,
		peerWrapper: params.PeerWrapper,
	}, nil
}

func (of *oracleFactory) NewOracle(ctx context.Context, args core.OracleArgs) (core.Oracle, error) {
	if !of.peerWrapper.IsStarted() {
		return nil, errors.New("peer wrapper not started")
	}
	oracle, err := ocr.NewOracle(ocr.OCR3OracleArgs[[]byte]{
		LocalConfig:                  args.LocalConfig,
		ContractConfigTracker:        args.ContractConfigTracker,
		ContractTransmitter:          args.ContractTransmitter,
		OffchainConfigDigester:       args.OffchainConfigDigester,
		ReportingPluginFactory:       args.ReportingPluginFactoryService,
		BinaryNetworkEndpointFactory: of.peerWrapper.Peer2,
		V2Bootstrappers:              of.config.BootstrapPeers,
		Database:                     of.database,
		Logger: ocrcommon.NewOCRWrapper(of.lggr, of.config.TraceLogging, func(ctx context.Context, msg string) {
			logger.Sugared(of.lggr).ErrorIf(of.jobORM.RecordError(ctx, of.jobID, msg), "unable to record error")
		}),
		// TODO?
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
