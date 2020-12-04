package services

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gopkg.in/guregu/null.v4"
)

// EthRequestEvent is a wrapper for `models.EthRequestEvent`, the DB
// representation of the job spec. It fulfills the job.Spec interface
// and has facilities for unmarshaling the pipeline DAG from the job spec text.
type EthRequestEventSpec struct {
	Type            string          `toml:"type"`
	SchemaVersion   uint32          `toml:"schemaVersion"`
	Name            null.String     `toml:"name"`
	MaxTaskDuration models.Interval `toml:"maxTaskDuration"`

	models.EthRequestEventSpec

	// The `jobID` field exists to cache the ID from the jobs table that joins
	// to the eth_request_events table.
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

// EthRequestEventSpec conforms to the job.Spec interface
var _ job.Spec = EthRequestEventSpec{}

func (spec EthRequestEventSpec) JobID() int32 {
	return spec.jobID
}

func (spec EthRequestEventSpec) JobType() job.Type {
	return models.EthRequestEventJobType
}

func (spec EthRequestEventSpec) TaskDAG() pipeline.TaskDAG {
	return spec.Pipeline
}

type ethRequestEventSpecDelegate struct{}

func (d *ethRequestEventSpecDelegate) JobType() job.Type {
	return models.EthRequestEventJobType
}

func (d *ethRequestEventSpecDelegate) ToDBRow(spec job.Spec) models.JobSpecV2 {
	concreteSpec, ok := spec.(EthRequestEventSpec)
	if !ok {
		panic(fmt.Sprintf("expected a services.EthRequestEventSpec, got %T", spec))
	}
	return models.JobSpecV2{
		EthRequestEventSpec: &concreteSpec.EthRequestEventSpec,
		Type:                string(models.EthRequestEventJobType),
		SchemaVersion:       concreteSpec.SchemaVersion,
		MaxTaskDuration:     concreteSpec.MaxTaskDuration,
	}
}

func (d *ethRequestEventSpecDelegate) FromDBRow(spec models.JobSpecV2) job.Spec {
	if spec.EthRequestEventSpec == nil {
		return nil
	}
	return &EthRequestEventSpec{
		EthRequestEventSpec: *spec.EthRequestEventSpec,
		jobID:               spec.ID,
	}
}

func (d *ethRequestEventSpecDelegate) ServicesForSpec(job.Spec) (services []job.Service, err error) {
	return
}

func RegisterEthRequestEventDelegate(jobSpawner job.Spawner) {
	jobSpawner.RegisterDelegate(
		NewEthRequestEventDelegate(jobSpawner),
	)
}

func NewEthRequestEventDelegate(jobSpawner job.Spawner) *ethRequestEventSpecDelegate {
	return &ethRequestEventSpecDelegate{}
}
