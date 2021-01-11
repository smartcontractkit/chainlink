package services

import (
	"context"
	"fmt"
	"reflect"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type DirectRequestSpecDelegate struct {
	logBroadcaster log.Broadcaster
	pipelineRunner pipeline.Runner
	db             *gorm.DB
}

func (d *DirectRequestSpecDelegate) JobType() job.Type {
	return job.DirectRequest
}

// ServicesForSpec returns the log listener service for a direct request job
// TODO: This will need heavy test coverage
func (d *DirectRequestSpecDelegate) ServicesForSpec(spec job.SpecDB) (services []job.Service, err error) {
	if spec.DirectRequestSpec == nil {
		return nil, errors.Errorf("services.DirectRequestSpecDelegate expects a *job.DirectRequestSpec to be present, got %v", spec)
	}
	concreteSpec := spec.DirectRequestSpec

	logListener := directRequestListener{
		d.logBroadcaster,
		concreteSpec.ContractAddress.Address(),
		d.pipelineRunner,
		d.db,
		spec.ID,
	}
	services = append(services, logListener)

	return
}

func RegisterDirectRequestDelegate(jobSpawner job.Spawner, logBroadcaster log.Broadcaster, pipelineRunner pipeline.Runner, db *gorm.DB) {
	jobSpawner.RegisterDelegate(
		NewDirectRequestDelegate(jobSpawner, logBroadcaster, pipelineRunner, db),
	)
}

func NewDirectRequestDelegate(jobSpawner job.Spawner, logBroadcaster log.Broadcaster, pipelineRunner pipeline.Runner, db *gorm.DB) *DirectRequestSpecDelegate {
	return &DirectRequestSpecDelegate{
		logBroadcaster,
		pipelineRunner,
		db,
	}
}

var (
	_ log.Listener = &directRequestListener{}
	_ job.Service  = &directRequestListener{}
)

type directRequestListener struct {
	logBroadcaster  log.Broadcaster
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

// OnConnect complies with log.Listener
func (directRequestListener) OnConnect() {}

// OnDisconnect complies with log.Listener
func (directRequestListener) OnDisconnect() {}

// OnConnect complies with log.Listener
func (d directRequestListener) HandleLog(lb log.Broadcast, err error) {
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

func (d *directRequestListener) handleOracleRequest(req contracts.OracleRequest) {
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

// SpecDB complies with log.Listener
func (d directRequestListener) JobIDV2() int32 {
	return d.jobID
}

// IsV2Job complies with log.Listener
func (directRequestListener) IsV2Job() bool {
	return true
}
