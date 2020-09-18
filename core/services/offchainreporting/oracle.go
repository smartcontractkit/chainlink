package offchainreporting

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrcontracts "github.com/smartcontractkit/offchain-reporting-design/prototype/contracts"
	ocr "github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting"
	ocrtypes "github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/types"
)

func RegisterJobType(
	db *gorm.DB,
	jobSpawner job.Spawner,
	pipelineRunner pipeline.Runner,
	ethClient eth.Client,
	logBroadcaster eth.LogBroadcaster,
) {
	jobSpawner.RegisterJobType(job.Registration{
		JobType: JobType,
		Spec:    &OracleSpec{},
		ServicesFactory: func(jobSpec job.Spec) ([]job.Service, error) {
			concreteSpec, ok := jobSpec.(*OracleSpec)
			if !ok {
				return nil, errors.Errorf("expected an offchainreporting.OracleSpec, got %T", jobSpec)
			} else if concreteSpec.JobID() == nil {
				return nil, errors.New("offchainreporting: got nil job ID")
			}

			aggregator, err := ocrcontracts.NewOffchainReportingAggregator(
				concreteSpec.ContractAddress,
				ethClient,
				logBroadcaster,
				*concreteSpec.JobID(),
				nil, // auth *bind.TransactOpts,
				"",  // rpcURL string,
			)
			if err != nil {
				return nil, err
			}

			service := ocr.NewOracle(ocr.OracleArgs{
				LocalConfig: ocrtypes.LocalConfig{
					DataSourceTimeout:                 concreteSpec.ObservationTimeout,
					BlockchainTimeout:                 concreteSpec.BlockchainTimeout,
					ContractConfigTrackerPollInterval: concreteSpec.ContractConfigTrackerPollInterval,
					ContractConfigConfirmations:       concreteSpec.ContractConfigConfirmations,
				},
				Database:              orm{jobSpecID: *concreteSpec.JobID(), db: db},
				Datasource:            dataSource{jobSpecID: *concreteSpec.JobID(), pipelineRunner: pipelineRunner},
				ContractTransmitter:   aggregator,
				ContractConfigTracker: aggregator,
				PrivateKeys:           ocrtypes.PrivateKeys(nil),
				NetEndpointFactory:    ocrtypes.BinaryNetworkEndpointFactory(nil),
				MonitoringEndpoint:    ocrtypes.MonitoringEndpoint(nil),
				Logger:                ocrtypes.Logger(nil),
				Bootstrappers:         []string{},
			})

			return []job.Service{service}, nil
		},
	})
}

// dataSource is an abstraction over the process of initiating a pipeline run
// and capturing the result.  Additionally, it converts the result to an
// ocrtypes.Observation (*big.Int), as expected by the offchain reporting library.
type dataSource struct {
	pipelineRunner pipeline.Runner
	jobSpecID      models.ID
}

var _ ocrtypes.DataSource = (*dataSource)(nil)

func (ds dataSource) Observe(ctx context.Context) (ocrtypes.Observation, error) {
	runID, err := ds.pipelineRunner.CreateRun(ds.jobSpecID)
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
		return nil, errors.Errorf("offchain reporting pipeline should have a single output (job spec ID: %v, pipeline run ID: %v)", ds.jobSpecID, runID)
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
