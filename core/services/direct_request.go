package services

import (
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
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

func NewDirectRequestDelegate(logBroadcaster log.Broadcaster, pipelineRunner pipeline.Runner, db *gorm.DB) *DirectRequestSpecDelegate {
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

	// TODO: Logic to handle log will go here

	err = lb.MarkConsumed()
	if err != nil {
		logger.Errorf("Error marking log as consumed: %v", err)
	}
}

// JobID complies with log.Listener
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
