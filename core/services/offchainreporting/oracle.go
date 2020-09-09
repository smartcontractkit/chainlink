package offchainreporting

import (

	// "math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

const PipelineType pipeline.Type = "offchainreporting"

// PipelineSpec conforms to the pipeline.Spec interface
var _ pipeline.Spec = JobSpec{}

func (spec JobSpec) JobID() *models.ID {
	return spec.UUID
}

func (spec JobSpec) Type() pipeline.Type {
	return PipelineType
}

func (spec JobSpec) Tasks() []Task {
	return spec.ObservationSource.Tasks()
}

// func (n *P2PBootstrapNode) Scan(value interface{}) error { return json.Unmarshal(value.([]byte), n) }
// func (n P2PBootstrapNode) Value() (driver.Value, error)  { return json.Marshal(n) }

func RegisterJobTypes(jobSpawner job.Spawner, orm ormInterface, ethClient eth.Client, logBroadcaster eth.LogBroadcaster) {
	jobSpawner.RegisterJobType(
		JobType,
		func(jobSpec job.JobSpec) ([]job.JobService, error) {
			concreteSpec, ok := jobSpec.(JobSpec)
			if !ok {
				return nil, errors.Errorf("expected an offchainreporting.JobSpec, got %T", jobSpec)
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

			return []job.JobService{service}, nil
		},
	)
}

// // dataSource is a simple wrapper around an existing job.Fetcher that converts
// // whatever value is fetched into a *big.Int, as the offchain reporting prototype
// // expects.
// type dataSource job.Fetcher

// var _ ocr.DataSource = dataSource(nil)

// func (ds dataSource) Fetch() (*big.Int, error) {
// 	val, err := job.Fetcher(ds).Fetch()
// 	if err != nil {
// 		return nil, err
// 	}
// 	asDecimal, err := utils.ToDecimal(val)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return asDecimal.BigInt(), nil
// }

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
