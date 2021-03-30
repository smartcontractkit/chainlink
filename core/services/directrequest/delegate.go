package directrequest

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/oracle_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
)

type Delegate struct {
	logBroadcaster log.Broadcaster
	pipelineRunner pipeline.Runner
	db             *gorm.DB
	ethClient      eth.Client
}

func NewDelegate(logBroadcaster log.Broadcaster, pipelineRunner pipeline.Runner, ethClient eth.Client, db *gorm.DB) *Delegate {
	return &Delegate{
		logBroadcaster,
		pipelineRunner,
		db,
		ethClient,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.DirectRequest
}

// ServicesForSpec returns the log listener service for a direct request job
// TODO: This will need heavy test coverage
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	if spec.DirectRequestSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *job.DirectRequestSpec to be present, got %v", spec)
	}
	concreteSpec := spec.DirectRequestSpec

	oracle, err := oracle_wrapper.NewOracle(concreteSpec.ContractAddress.Address(), d.ethClient)
	if err != nil {
		return
	}

	logListener := listener{
		logBroadcaster: d.logBroadcaster,
		oracle:         oracle,
		pipelineRunner: d.pipelineRunner,
		db:             d.db,
		jobID:          spec.ID,
	}
	services = append(services, logListener)

	return
}

var (
	_ log.Listener = &listener{}
	_ job.Service  = &listener{}
)

type listener struct {
	logBroadcaster  log.Broadcaster
	unsubscribeLogs func()
	oracle          oracle_wrapper.OracleInterface
	pipelineRunner  pipeline.Runner
	db              *gorm.DB
	jobID           int32
}

// Start complies with job.Service
func (d listener) Start() error {
	connected, unsubscribe := d.logBroadcaster.Register(d, log.ListenerOpts{
		Contract: d.oracle,
		Logs:     []generated.AbigenLog{},
	})
	if !connected {
		return errors.New("Failed to register listener with logBroadcaster")
	}
	d.unsubscribeLogs = unsubscribe
	return nil
}

// Close complies with job.Service
func (d listener) Close() error {
	if d.unsubscribeLogs != nil {
		d.unsubscribeLogs()
	}
	return nil
}

// OnConnect complies with log.Listener
func (listener) OnConnect() {}

// OnDisconnect complies with log.Listener
func (listener) OnDisconnect() {}

// OnConnect complies with log.Listener
func (d listener) HandleLog(lb log.Broadcast) {
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

	// TODO: Need to filter each job to the _jobId - can we filter upstream?
	// TODO: Will need to generate a jobID somehow... hash of DAG?
	switch log := log.(type) {
	case *oracle_wrapper.OracleOracleRequest:
		d.handleOracleRequest(log)
		err = lb.MarkConsumed()
		if err != nil {
			logger.Errorf("Error marking log as consumed: %v", err)
		}
	case *oracle_wrapper.OracleCancelOracleRequest:
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

func oracleRequestToMap(req *oracle_wrapper.OracleOracleRequest) map[string]interface{} {
	result := make(map[string]interface{})
	result["specId"] = fmt.Sprintf("0x%x", req.SpecId)
	result["requester"] = req.Requester.Hex()
	result["requestId"] = fmt.Sprintf("0x%x", req.RequestId)
	result["payment"] = fmt.Sprintf("%v", req.Payment)
	result["callbackAddr"] = req.CallbackAddr.Hex()
	result["callbackFunctionId"] = fmt.Sprintf("0x%x", req.CallbackFunctionId)
	result["cancelExpiration"] = fmt.Sprintf("%v", req.CancelExpiration)
	result["dataVersion"] = fmt.Sprintf("%v", req.DataVersion)
	result["data"] = fmt.Sprintf("0x%x", req.Data)
	return result
}

func (d *listener) handleOracleRequest(req *oracle_wrapper.OracleOracleRequest) {
	meta := make(map[string]interface{})
	meta["oracleRequest"] = oracleRequestToMap(req)
	ctx := context.TODO()
	_, err := d.pipelineRunner.CreateRun(ctx, d.jobID, meta)
	if err != nil {
		logger.Errorw("DirectRequest failed to create run", "err", err)
	}
}

// Cancels runs that haven't been started yet, with the given request ID
// TODO: Boy does this ever need testing
func (d *listener) handleCancelOracleRequest(requestID [32]byte) {
	err := d.db.Exec(`
DELETE FROM pipeline_runs
WHERE id IN (
	SELECT pipeline_runs.id
	FROM pipeline_runs
	INNER JOIN pipeline_task_runs
	ON pipeline_task_runs.pipeline_run_id = pipeline_runs.id
	WHERE pipeline_spec_id = ?
		AND pipeline_runs.meta->'oracleRequest'->>'requestId' = ?
		AND pipeline_task_runs.finished_at IS NULL
	FOR UPDATE OF pipeline_task_runs
	SKIP LOCKED
)`, d.jobID, fmt.Sprintf("0x%x", requestID)).Error
	if err != nil {
		logger.Errorw("Error while deleting pipeline_runs", "error", err)
	}
}

// JobID complies with log.Listener
func (listener) JobID() models.JobID {
	return models.NilJobID
}

// Job complies with log.Listener
func (d listener) JobIDV2() int32 {
	return d.jobID
}

// IsV2Job complies with log.Listener
func (listener) IsV2Job() bool {
	return true
}
