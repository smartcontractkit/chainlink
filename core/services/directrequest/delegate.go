package directrequest

import (
	"context"
	"fmt"
	"reflect"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
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
}

func NewDelegate(logBroadcaster log.Broadcaster, pipelineRunner pipeline.Runner, db *gorm.DB) *Delegate {
	return &Delegate{
		logBroadcaster,
		pipelineRunner,
		db,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.DirectRequest
}

// ServicesForSpec returns the log listener service for a direct request job
// TODO: This will need heavy test coverage
func (d *Delegate) ServicesForSpec(spec job.SpecDB) (services []job.Service, err error) {
	if spec.DirectRequestSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *job.DirectRequestSpec to be present, got %v", spec)
	}
	concreteSpec := spec.DirectRequestSpec

	logListener := listener{
		d.logBroadcaster,
		concreteSpec.ContractAddress.Address(),
		d.pipelineRunner,
		d.db,
		spec.ID,
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
	contractAddress gethCommon.Address
	pipelineRunner  pipeline.Runner
	db              *gorm.DB
	jobID           int32
}

// Start complies with job.Service
func (d listener) Start() error {
	connected := d.logBroadcaster.Register(nil, d)
	if !connected {
		return errors.New("Failed to register listener with logBroadcaster")
	}
	return nil
}

// Close complies with job.Service
func (d listener) Close() error {
	d.logBroadcaster.Unregister(nil, d)
	return nil
}

// OnConnect complies with log.Listener
func (listener) OnConnect() {}

// OnDisconnect complies with log.Listener
func (listener) OnDisconnect() {}

// OnConnect complies with log.Listener
func (d listener) HandleLog(lb log.Broadcast, err error) {
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

	// TODO: Need to filter each job to the _jobId - can we filter upstream?
	// TODO: Will need to generate a jobID somehow... hash of DAG?
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

func (d *listener) handleOracleRequest(req contracts.OracleRequest) {
	meta := make(map[string]interface{})
	meta["oracleRequest"] = req.ToMap()
	panic("HERE")
	ctx := context.TODO()
	_, err := d.pipelineRunner.CreateRun(ctx, d.jobID, meta)
	if err != nil {
		logger.Errorw("DirectRequest failed to create run", "err", err)
	}
}

// Cancels runs that haven't been started yet, with the given request ID
// TODO: Boy does this ever need testing
func (d *listener) handleCancelOracleRequest(requestID [32]byte) {
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

// JobID complies with log.Listener
func (listener) JobID() models.JobID {
	return models.NilJobID
}

// SpecDB complies with log.Listener
func (d listener) JobIDV2() int32 {
	return d.jobID
}

// IsV2Job complies with log.Listener
func (listener) IsV2Job() bool {
	return true
}
