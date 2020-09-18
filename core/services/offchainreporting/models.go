package offchainreporting

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type OracleSpec struct {
	models.OffchainReportingOracleSpec
	ObservationSource pipeline.TaskDAG `toml:"observationSource"`
}

type KeyBundle models.OffchainReportingKeyBundle

func (OracleSpec) TableName() string { return "offchainreporting_oracle_specs" }
func (KeyBundle) TableName() string  { return "offchainreporting_key_bundles" }

const JobType job.Type = "offchainreporting"

// OracleSpec conforms to the job.Spec interface
var _ job.Spec = OracleSpec{}

func (spec OracleSpec) JobID() int32 {
	return spec.JID
}

func (spec OracleSpec) JobType() job.Type {
	return JobType
}

func (spec OracleSpec) TaskDAG() pipeline.TaskDAG {
	return spec.ObservationSource
}
