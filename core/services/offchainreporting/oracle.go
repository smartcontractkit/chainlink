package offchainreporting

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocr "github.com/smartcontractkit/offchain-reporting/lib/offchainreporting"
	ocrtypes "github.com/smartcontractkit/offchain-reporting/lib/offchainreporting/types"
	ocrcontracts "github.com/smartcontractkit/offchain-reporting/lib/prototype/contracts"
)

const JobType job.Type = "offchainreporting"

func RegisterJobType(
	db *gorm.DB,
	jobSpawner job.Spawner,
	pipelineRunner pipeline.Runner,
	ethClient eth.Client,
	logBroadcaster eth.LogBroadcaster,
) {
	jobSpawner.RegisterDelegate(
		NewJobSpawnerDelegate(db, pipelineRunner, ethClient, logBroadcaster),
	)
}

type jobSpawnerDelegate struct {
	db             *gorm.DB
	pipelineRunner pipeline.Runner
	ethClient      eth.Client
	logBroadcaster eth.LogBroadcaster
}

func NewJobSpawnerDelegate(
	db *gorm.DB,
	pipelineRunner pipeline.Runner,
	ethClient eth.Client,
	logBroadcaster eth.LogBroadcaster,
) *jobSpawnerDelegate {
	return &jobSpawnerDelegate{db, pipelineRunner, ethClient, logBroadcaster}
}

func (d jobSpawnerDelegate) JobType() job.Type {
	return JobType
}

func (d jobSpawnerDelegate) ToDBRow(spec job.Spec) models.JobSpecV2 {
	concreteSpec, ok := spec.(*OracleSpec)
	if !ok {
		panic(fmt.Sprintf("expected an *offchainreporting.OracleSpec, got %T", spec))
	}
	return models.JobSpecV2{OffchainreportingOracleSpec: &concreteSpec.OffchainReportingOracleSpec}
}

func (d jobSpawnerDelegate) FromDBRow(spec models.JobSpecV2) job.Spec {
	if spec.OffchainreportingOracleSpec == nil {
		return nil
	}
	return &OracleSpec{
		OffchainReportingOracleSpec: *spec.OffchainreportingOracleSpec,
		jobID:                       spec.ID,
	}
}

func (d jobSpawnerDelegate) ServicesForSpec(spec job.Spec) ([]job.Service, error) {
	concreteSpec := spec.(*OracleSpec)

	aggregator, err := ocrcontracts.NewOffchainReportingAggregator(
		concreteSpec.ContractAddress,
		d.ethClient,
		d.logBroadcaster,
		concreteSpec.JobID(),
		nil, // auth *bind.TransactOpts,
		"",  // rpcURL string,
	)
	if err != nil {
		return nil, err
	}

	service := ocr.NewOracle(ocr.OracleArgs{
		LocalConfig: ocrtypes.LocalConfig{
			DataSourceTimeout: time.Duration(concreteSpec.ObservationTimeout),
			BlockchainTimeout: time.Duration(concreteSpec.BlockchainTimeout),
			// ContractConfigTrackerSubscribeInterval: time.Duration(concreteSpec.ContractConfigTrackerSubscribeInterval),
			ContractConfigTrackerPollInterval: time.Duration(concreteSpec.ContractConfigTrackerPollInterval),
			ContractConfigConfirmations:       concreteSpec.ContractConfigConfirmations,
		},
		Database:              nil, //orm{jobID: concreteSpec.JobID(), db: d.db},
		Datasource:            dataSource{jobID: concreteSpec.JobID(), pipelineRunner: d.pipelineRunner},
		ContractTransmitter:   aggregator,
		ContractConfigTracker: aggregator,
		PrivateKeys:           ocrtypes.PrivateKeys(nil),
		NetEndpointFactory:    ocrtypes.BinaryNetworkEndpointFactory(nil),
		MonitoringEndpoint:    ocrtypes.MonitoringEndpoint(nil),
		Logger:                ocrtypes.Logger(nil),
		Bootstrappers:         []string{},
	})

	return []job.Service{service}, nil
}

// dataSource is an abstraction over the process of initiating a pipeline run
// and capturing the result.  Additionally, it converts the result to an
// ocrtypes.Observation (*big.Int), as expected by the offchain reporting library.
type dataSource struct {
	pipelineRunner pipeline.Runner
	jobID          int32
}

var _ ocrtypes.DataSource = (*dataSource)(nil)

func (ds dataSource) Observe(ctx context.Context) (ocrtypes.Observation, error) {
	runID, err := ds.pipelineRunner.CreateRun(ds.jobID, nil)
	if err != nil {
		return nil, err
	}

	err = ds.pipelineRunner.AwaitRun(ctx, runID)
	if err != nil {
		return nil, err
	}

	results, err := ds.pipelineRunner.ResultsForRun(runID)
	if err != nil {
		return nil, err
	} else if len(results) != 1 {
		return nil, errors.Errorf("offchain reporting pipeline should have a single output (job spec ID: %v, pipeline run ID: %v)", ds.jobID, runID)
	}
	result := results[0]
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "pipeline error")
	}

	asDecimal, err := utils.ToDecimal(result.Value)
	if err != nil {
		return nil, err
	}
	return ocrtypes.Observation(asDecimal.BigInt()), nil
}
