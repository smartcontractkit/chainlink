package offchainreporting

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrcontract "github.com/smartcontractkit/offchain-reporting-design/contract"
	"github.com/smartcontractkit/offchain-reporting-design/prototype/gethwrappers/ulairi"
	ocr "github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting"
	ocrconfig "github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/config"
	ocrimpls "github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/to_be_internal/testimplementations"
	ocrtypes "github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/types"
)

const JobType = "offchainreporting"

type JobSpec struct {
	ID                 models.ID          `json:"id"`
	ContractAddress    common.Address     `json:"contractAddress"`
	P2PNodeID          string             `json:"p2pNodeID"`
	P2PBootstrapNodes  []P2PBootstrapNode `json:"p2pBootstrapNodes"`
	KeyBundle          string             `json:"keyBundle"`
	MonitoringEndpoint string             `json:"monitoringEndpoint"`
	NodeAddress        common.Address     `json:"nodeAddress"`
	ObservationTimeout time.Duration      `json:"observationTimeout"`
	ObservationSource  job.Fetcher        `json:"observationSource"`
	LogLevel           ocrtypes.LogLevel  `json:"logLevel,omitempty"`
}

// JobSpec conforms to the job.JobSpec interface
var _ job.JobSpec = JobSpec{}

func (spec JobSpec) JobID() *models.ID {
	return &spec.ID
}

func (spec JobSpec) JobType() string {
	return JobType
}

func (spec *JobSpec) UnmarshalJSON(bs []byte) error {
	if spec == nil {
		*spec = JobSpec{}
	}

	err := json.Unmarshal(bs, spec)
	if err != nil {
		return err
	}

	var obsSrc struct {
		ObservationSource json.RawMessage `json:"observationSource"`
	}
	err = json.Unmarshal(bs, &obsSrc)
	if err != nil {
		return err
	}
	fetcher, err := job.UnmarshalFetcherJSON([]byte(obsSrc.ObservationSource))
	if err != nil {
		return err
	}
	spec.ObservationSource = fetcher
	return nil
}

type P2PBootstrapNode struct {
	PeerID    string `json:"peerID"`
	Multiaddr string `json:"multiAddr"`
}

func RegisterJobTypes(jobSpawner job.Spawner, orm ormInterface, ethClient eth.Client, logBroadcaster eth.LogBroadcaster) {
	jobSpawner.RegisterJobType(
		JobType,
		func(jobSpec job.JobSpec) (job.JobService, error) {
			concreteSpec, ok := jobSpec.(JobSpec)
			if !ok {
				return nil, errors.Errorf("expected an offchainreporting.JobSpec, got %T", jobSpec)
			}

			db := database{JobSpecID: concreteSpec.ID, orm: orm}

			config, err := db.ReadConfig()
			if err != nil {
				return nil, err
			}

			aggregator := ocrcontracts.NewOffchainReportingAggregator(
				concreteSpec.ContractAddress,
				ethClient,
				logBroadcaster,
				concreteSpec.ID,
			)

			privateKeys := ocrimpls.NewPrivateKeys(nil)            // @@TODO
			netEndpoint := ocrtypes.BinaryNetworkEndpoint(nil)     // @@TODO
			monitoringEndpoint := ocrtypes.MonitoringEndpoint(nil) // @@TODO
			localConfig := ocrtypes.LocalConfig{
				DatasourceTimeout: concreteSpec.ObservationTimeout,
				LogLevel:          concreteSpec.LogLevel,
			}

			return ocr.Run(ocr.Params{
				LocalConfig:           localConfig,
				PrivateKeys:           privateKeys,
				NetEndPoint:           netEndpoint,
				Datasource:            dataSource(concreteSpec.ObservationSource),
				ContractTransmitter:   aggregator,
				ContractConfigTracker: aggregator,
				MonitoringEndpoint:    monitoringEndpoint,
			}), nil
		},
	)
}

// dataSource is a simple wrapper around an existing job.Fetcher that converts
// whatever value is fetched into a *big.Int, as the offchain reporting prototype
// expects.
type dataSource job.Fetcher

var _ ocr.DataSource = dataSource(nil)

func (ds dataSource) Fetch() (*big.Int, error) {
	val, err := job.Fetcher(ds).Fetch()
	if err != nil {
		return nil, err
	}
	asDecimal, err := utils.ToDecimal(val)
	if err != nil {
		return nil, err
	}
	return asDecimal.BigInt(), nil
}

// database is an abstraction that conforms to the Database interface in the
// offchain reporting prototype, which is unaware of job IDs.
type database struct {
	orm       ormInterface
	JobSpecID models.ID
}

var _ ocr.Database = database{}

type ormInterface interface {
	FindOffchainReportingPersistentState(jobID models.ID, groupID ocrtypes.GroupID) (PersistentState, error)
	SaveOffchainReportingPersistentState(state PersistentState) error
	FindOffchainReportingConfig(jobID models.ID) (Config, error)
	SaveOffchainReportingConfig(config Config) error
}

type PersistentState struct {
	JobSpecID models.ID
	GroupID   ocrtypes.GroupID
	ocr.PersistentState
}

type Config struct {
	JobSpecID models.ID
	ocrconfig.Config
}

func (db database) ReadState(groupID ocrtypes.GroupID) (*ocr.PersistentState, error) {
	state, err := db.orm.FindOffchainReportingPersistentState(db.JobSpecID, groupID)
	if err != nil {
		return &ocr.PersistentState{}, err
	}
	return state.PersistentState, nil
}

func (db database) WriteState(groupID ocrtypes.GroupID, state ocr.PersistentState) error {
	return db.orm.SaveOffchainReportingPersistentState(PersistentState{
		ID:              db.JobSpecID,
		PersistentState: state,
	})
}

func (db database) ReadConfig() (ocrconfig.Config, error) {
	config, err := db.orm.FindOffchainReportingConfig(db.JobSpecID)
	if err != nil {
		return ocr.Config{}, err
	}
	return config.Config, nil
}

func (db database) WriteConfig(config ocrconfig.Config) error {
	return db.orm.SaveOffchainReportingConfig(Config{
		ID:     db.JobSpecID,
		Config: config,
	})
}
