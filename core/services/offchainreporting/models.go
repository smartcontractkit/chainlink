package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gopkg.in/guregu/null.v4"
)

// OracleSpec is a wrapper for `models.OffchainReportingOracleSpec`, the DB
// representation of the OCR job spec.  It fulfills the job.Spec interface
// and has facilities for unmarshaling the pipeline DAG from the job spec text.
type OracleSpec struct {
	Type          string      `toml:"type"`
	SchemaVersion uint32      `toml:"schemaVersion"`
	Name          null.String `toml:"name"`

	models.OffchainReportingOracleSpec

	// The `jobID` field exists to cache the ID from the jobs table that joins
	// to the offchainreporting_oracle_specs table.
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

// OracleSpec conforms to the job.Spec interface
var _ job.Spec = OracleSpec{}

func (spec OracleSpec) JobID() int32 {
	return spec.jobID
}

func (spec OracleSpec) JobType() job.Type {
	return JobType
}

func (spec OracleSpec) TaskDAG() pipeline.TaskDAG {
	return spec.Pipeline
}
