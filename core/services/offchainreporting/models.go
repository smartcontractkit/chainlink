package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// OracleSpec is a wrapper for `models.OffchainReportingOracleSpec`, the DB
// representation of the OCR job spec.  It fulfills the job.Spec interface
// and has facilities for unmarshaling the pipeline DAG from the job spec text.
type OracleSpec struct {
	models.OffchainReportingOracleSpec
	jobID    int32
	Pipeline pipeline.TaskDAG `toml:"observationSource"` // This field is only used during unmarshaling
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
