package jobs

import "github.com/google/uuid"

// JobSpecType is the type of job spec set in the TOML.
type JobSpecType string

var (
	JobSpecTypeLLO       JobSpecType = "offchainreporting2" // used for LLO specs
	JobSpecTypeBootstrap JobSpecType = "bootstrap"          // used for bootstrap specs in LLO dons
	JobSpecTypeStream    JobSpecType = "stream"             // used for stream specs in LLO dons
)

type Base struct {
	Name          string      `toml:"name"`
	Type          JobSpecType `toml:"type"`
	SchemaVersion uint        `toml:"schemaVersion"`
	ExternalJobID uuid.UUID   `toml:"externalJobID"`
}

// Job interface
type JobSpec interface {
	MarshalTOML() ([]byte, error)
}
