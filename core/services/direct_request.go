package services

import (
	"fmt"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gopkg.in/guregu/null.v4"
)

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
	logBroadcaster log.Broadcaster
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

// ServicesForSpec returns the log listener service for a direct request job
// TODO: This will need heavy test coverage
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

// JobSpecV2 complies with log.Listener
func (d directRequestListener) JobIDV2() int32 {
	return d.jobID
}

// IsV2Job complies with log.Listener
func (directRequestListener) IsV2Job() bool {
	return true
}
