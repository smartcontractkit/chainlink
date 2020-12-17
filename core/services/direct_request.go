package services

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"gopkg.in/guregu/null.v4"
)

var (
	DirectRequestLogTopic = getDirectRequestLogTopic()
)

func getDirectRequestLogTopic() gethCommon.Hash {
	abi, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
	if err != nil {
		panic("could not parse OffchainAggregator ABI: " + err.Error())
	}
	return abi.Events["ConfigSet"].ID
}

// DirectRequest is a wrapper for `models.DirectRequest`, the DB
// representation of the job spec. It fulfills the job.Spec interface
// and has facilities for unmarshaling the pipeline DAG from the job spec text.
type DirectRequestSpec struct {
	Type            string          `toml:"type"`
	SchemaVersion   uint32          `toml:"schemaVersion"`
	Name            null.String     `toml:"name"`
	MaxTaskDuration models.Interval `toml:"maxTaskDuration"`

	models.DirectRequestSpec

	// The `jobID` field exists to cache the ID from the jobs table that joins
	// to the direct_requests table.
	jobID int32

	// The `Pipeline` field is only used during unmarshaling.  A pipeline.TaskDAG
	// is a type that implements gonum.org/v1/gonum/graph#Graph, which means that
	// you can dot.Unmarshal(...) raw DOT source directly into it, and it will
	// be a fully-instantiated DAG containing information about all of the nodes
	// and edges described by the DOT.  Our pipeline.TaskDAG type has a method
	// called `.TasksInDependencyOrder()` which converts this node/edge data
	// structure into task specs which can then be saved to the database.
	Pipeline pipeline.TaskDAG `toml:"observationSource"`
}

// DirectRequestSpec conforms to the job.Spec interface
var _ job.Spec = DirectRequestSpec{}

func (spec DirectRequestSpec) JobID() int32 {
	return spec.jobID
}

func (spec DirectRequestSpec) JobType() job.Type {
	return models.DirectRequestJobType
}

func (spec DirectRequestSpec) TaskDAG() pipeline.TaskDAG {
	return spec.Pipeline
}

type DirectRequestSpecDelegate struct {
	logBroadcaster eth.LogBroadcaster
	pipelineRunner pipeline.Runner
	db             *gorm.DB
}

func (d *DirectRequestSpecDelegate) JobType() job.Type {
	return models.DirectRequestJobType
}

func (d *DirectRequestSpecDelegate) ToDBRow(spec job.Spec) models.JobSpecV2 {
	concreteSpec, ok := spec.(DirectRequestSpec)
	if !ok {
		panic(fmt.Sprintf("expected a services.DirectRequestSpec, got %T", spec))
	}
	return models.JobSpecV2{
		DirectRequestSpec: &concreteSpec.DirectRequestSpec,
		Type:              string(models.DirectRequestJobType),
		SchemaVersion:     concreteSpec.SchemaVersion,
		MaxTaskDuration:   concreteSpec.MaxTaskDuration,
	}
}

func (d *DirectRequestSpecDelegate) FromDBRow(spec models.JobSpecV2) job.Spec {
	if spec.DirectRequestSpec == nil {
		return nil
	}
	return &DirectRequestSpec{
		DirectRequestSpec: *spec.DirectRequestSpec,
		jobID:             spec.ID,
	}
}

// ServicesForSpec TODO
func (d *DirectRequestSpecDelegate) ServicesForSpec(spec job.Spec) (services []job.Service, err error) {
	concreteSpec, is := spec.(*DirectRequestSpec)
	if !is {
		return nil, errors.Errorf("services.DirectRequestSpecDelegate expects a *services.DirectRequestSpec, got %T", spec)
	}

	logListener := directRequestListener{
		d.logBroadcaster,
		concreteSpec.ContractAddress.Address(),
		d.pipelineRunner,
		d.db,
		spec.JobID(),
	}
	services = append(services, logListener)

	return
}

func RegisterDirectRequestDelegate(jobSpawner job.Spawner, logBroadcaster eth.LogBroadcaster, pipelineRunner pipeline.Runner, db *gorm.DB) {
	jobSpawner.RegisterDelegate(
		NewDirectRequestDelegate(jobSpawner, logBroadcaster, pipelineRunner, db),
	)
}

func NewDirectRequestDelegate(jobSpawner job.Spawner, logBroadcaster eth.LogBroadcaster, pipelineRunner pipeline.Runner, db *gorm.DB) *DirectRequestSpecDelegate {
	return &DirectRequestSpecDelegate{
		logBroadcaster,
		pipelineRunner,
		db,
	}
}

var (
	_ eth.LogListener = &directRequestListener{}
	_ job.Service     = &directRequestListener{}
)

type directRequestListener struct {
	logBroadcaster  eth.LogBroadcaster
	contractAddress gethCommon.Address
	pipelineRunner  pipeline.Runner
	db              *gorm.DB
	jobID           int32
}

// Start complies with job.Service
func (d directRequestListener) Start() error {
	connected := d.logBroadcaster.Register(d.contractAddress, d)
	if !connected {
		return errors.New("Failed to register directRequestListener with logBroadcaster")
	}
	return nil
}

// Close complies with job.Service
func (d directRequestListener) Close() error {
	d.logBroadcaster.Unregister(d.contractAddress, d)
	return nil
}

// OnConnect complies with eth.LogListener
func (directRequestListener) OnConnect() {}

// OnDisconnect complies with eth.LogListener
func (directRequestListener) OnDisconnect() {}

// OnConnect complies with eth.LogListener
func (d directRequestListener) HandleLog(lb eth.LogBroadcast, err error) {
	if err != nil {
		logger.Errorw("DirectRequestListener: error in previous LogListener", "err", err)
		return
	}

	was, err := lb.WasAlreadyConsumed()
	if err != nil {
		logger.Errorw("DirectRequestListener: could not determine if log was already consumed", "error", err)
		return
	} else if was {
		return
	}

	log := lb.DecodedLog()
	if log == nil || reflect.ValueOf(log).IsNil() {
		logger.Error("HandleLog: ignoring nil value")
		return
	}

	switch log := log.(type) {
	case *contracts.LogOracleRequest:
		d.handleOracleRequest(log.ToOracleRequest())
		err = lb.MarkConsumed()
		if err != nil {
			logger.Errorf("Error marking log as consumed: %v", err)
		}
	case *contracts.LogCancelOracleRequest:
		d.handleCancelOracleRequest(log.RequestId)
		// TODO: Transactional/atomic log consumption would be nice
		err = lb.MarkConsumed()
		if err != nil {
			logger.Errorf("Error marking log as consumed: %v", err)
		}

	default:
		logger.Warnf("unexpected log type %T", log)
	}
}

func (d *directRequestListener) handleOracleRequest(req contracts.OracleRequest) {
	meta := make(map[string]interface{})
	meta["oracleRequest"] = req.ToMap()
	ctx := context.TODO()
	_, err := d.pipelineRunner.CreateRun(ctx, d.jobID, meta)
	if err != nil {
		logger.Errorw("DirectRequest failed to create run", "err", err)
	}
}

// Cancels runs that haven't been started yet, with the given request ID
// TODO: Boy does this ever need testing
func (d *directRequestListener) handleCancelOracleRequest(requestID [32]byte) {
	d.db.Exec(`
	DELETE FROM pipeline_runs 
	WHERE id IN (
		SELECT id FROM pipeline_runs FOR UPDATE OF pipeline_task_runs SKIP LOCKED
		INNER JOIN pipeline_task_runs WHERE pipeline_task_runs.pipeline_run_id = pipeline_runs.id
		WHERE pipeline_spec_id = ?
		AND pipeline_runs.meta->'oracleRequest'->'requestId' = ?
		HAVING bool_and(pipeline_task_runs.finished_at IS NULL)
	)
	`)
}

// JobID complies with eth.LogListener
func (directRequestListener) JobID() *models.ID {
	return nil
}

// JobSpecV2 complies with eth.LogListener
func (d directRequestListener) JobIDV2() int32 {
	return d.jobID
}

// IsV2Job complies with eth.LogListener
func (directRequestListener) IsV2Job() bool {
	return true
}
