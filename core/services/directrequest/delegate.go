package directrequest

import (
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
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
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	if spec.DirectRequestSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *job.DirectRequestSpec to be present, got %v", spec)
	}
	concreteSpec := spec.DirectRequestSpec

	logListener := listener{
		d.logBroadcaster,
		nil,
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
	unsubscribeLogs func()
	contractAddress gethCommon.Address
	pipelineRunner  pipeline.Runner
	db              *gorm.DB
	jobID           int32
}

// Start complies with job.Service
func (d listener) Start() error {
	connected, unsubscribe := d.logBroadcaster.Register(d, log.ListenerOpts{
		Contract: nil,
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
	d.unsubscribeLogs()
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

	// TODO: Logic to handle log will go here

	err = lb.MarkConsumed()
	if err != nil {
		logger.Errorf("Error marking log as consumed: %v", err)
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
