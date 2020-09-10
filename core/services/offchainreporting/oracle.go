package offchainreporting

import (
	"github.com/golangci/golangci-lint/pkg/result"
	// "math/big"

	"github.com/pkg/errors"
	// "github.com/pkg/errors"

	// "github.com/smartcontractkit/chainlink/core/services/eth"

	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"

	// "github.com/smartcontractkit/chainlink/core/utils"
	// "github.com/smartcontractkit/offchain-reporting-design/prototype/gethwrappers/ulairi"
	// ocr "github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting"
	// ocrimpls "github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/to_be_internal/testimplementations"
	ocrtypes "github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/types"
)

func RegisterJobTypes(jobSpawner pipeline.Spawner, orm ormInterface, ethClient eth.Client, logBroadcaster eth.LogBroadcaster) {
	jobSpawner.RegisterJobType(
		JobType,
		func(jobSpec pipeline.JobSpec) ([]pipeline.JobService, error) {
			concreteSpec, ok := jobSpec.(*OracleSpec)
			if !ok {
				return nil, errors.Errorf("expected an offchainreporting.OracleSpec, got %T", jobSpec)
			}

			db := database{JobSpecID: concreteSpec.ID, orm: orm}

			config, err := db.ReadConfig()
			if err != nil {
				return nil, err
			}

			// aggregator := ocrcontracts.NewOffchainReportingAggregator(
			// 	concreteSpec.ContractAddress,
			// 	ethClient,
			// 	logBroadcaster,
			// 	concreteSpec.ID,
			// )

			privateKeys := ocrimpls.NewPrivateKeys(nil)            // @@TODO
			netEndpoint := ocrtypes.BinaryNetworkEndpoint(nil)     // @@TODO
			monitoringEndpoint := ocrtypes.MonitoringEndpoint(nil) // @@TODO
			localConfig := ocrtypes.LocalConfig{
				DatasourceTimeout: concreteSpec.ObservationTimeout,
				LogLevel:          concreteSpec.LogLevel,
			}

			service, err := ocr.Run(ocr.Params{
				LocalConfig: localConfig,
				PrivateKeys: privateKeys,
				NetEndPoint: netEndpoint,
				Datasource:  dataSource(concreteSpec.ObservationSource),
				// ContractTransmitter:   aggregator,
				// ContractConfigTracker: aggregator,
				MonitoringEndpoint: monitoringEndpoint,
			}), nil
			if err != nil {
				return nil, err
			}

			return []pipeline.JobService{service}, nil
		},
	)
}

// dataSource is an abstraction over the process of initiating a pipeline run
// and capturing the result.  Additionally, it converts the result to a *big.Int,
// as expected by the offchain reporting library.
type dataSource struct {
	pipelineRunner pipeline.Runner
	pipelineSpecID int64
}

var _ ocr.DataSource = (*dataSource)(nil)

func (ds dataSource) Fetch() (*big.Int, error) {
	runID, err := ds.pipelineRunner.CreatePipelineRun(ds.pipelineSpecID)
	if err != nil {
		return nil, err
	}

	<-ds.pipelineRunner.AwaitRun(runID)

	results, err := ds.pipelineRunner.ResultsForRun(runID)
	if err != nil {
		return nil, err
	} else if len(results) != 1 {
		return nil, errors.Errorf("offchain reporting pipeline should have a single output (pipeline spec ID: %v, pipeline run ID: %v)", ds.pipelineSpecID, runID)
	}
	result := results[0]
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error, "pipeline error")
	}

	asDecimal, err := utils.ToDecimal(result.Value)
	if err != nil {
		return nil, err
	}
	return asDecimal.BigInt(), nil
}

// // database is an abstraction that conforms to the Database interface in the
// // offchain reporting prototype, which is unaware of job IDs.
// type database struct {
// 	orm       ormInterface
// 	JobSpecID models.ID
// }

// var _ ocr.Database = database{}

// type ormInterface interface {
// 	FindOffchainReportingPersistentState(jobID models.ID, groupID ocrtypes.GroupID) (PersistentState, error)
// 	SaveOffchainReportingPersistentState(state PersistentState) error
// 	FindOffchainReportingConfig(jobID models.ID) (Config, error)
// 	SaveOffchainReportingConfig(config Config) error
// }

// type PersistentState struct {
// 	JobSpecID models.ID
// 	GroupID   ocrtypes.GroupID
// 	ocr.PersistentState
// }

// type Config struct {
// 	JobSpecID models.ID
// 	ocrtypes.Config
// }

// func (db database) ReadState(groupID ocrtypes.GroupID) (*ocr.PersistentState, error) {
// 	state, err := db.orm.FindOffchainReportingPersistentState(db.JobSpecID, groupID)
// 	if err != nil {
// 		return &ocr.PersistentState{}, err
// 	}
// 	return state.PersistentState, nil
// }

// func (db database) WriteState(groupID ocrtypes.GroupID, state ocr.PersistentState) error {
// 	return db.orm.SaveOffchainReportingPersistentState(PersistentState{
// 		ID:              db.JobSpecID,
// 		PersistentState: state,
// 	})
// }

// func (db database) ReadConfig() (ocrtypes.Config, error) {
// 	config, err := db.orm.FindOffchainReportingConfig(db.JobSpecID)
// 	if err != nil {
// 		return ocr.Config{}, err
// 	}
// 	return config.Config, nil
// }

// func (db database) WriteConfig(config ocrtypes.Config) error {
// 	return db.orm.SaveOffchainReportingConfig(Config{
// 		ID:     db.JobSpecID,
// 		Config: config,
// 	})
// }
