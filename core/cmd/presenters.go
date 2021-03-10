package cmd

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// JAID represents a JSON API ID.
// It implements the api2go MarshalIdentifier and UnmarshalIdentitier interface.
type JAID struct {
	ID string `json:"-"`
}

// GetID implements the api2go MarshalIdentifier interface.
func (jaid JAID) GetID() string {
	return jaid.ID
}

// SetID implements the api2go UnmarshalIdentitier interface.
func (jaid *JAID) SetID(value string) error {
	jaid.ID = value

	return nil
}

// JobType defines the the job type
type JobType string

func (t JobType) String() string {
	return string(t)
}

const (
	// DirectRequestJob defines a Direct Request Job
	DirectRequestJob JobType = "directrequest"
	// FluxMonitorJob defines a Flux Monitor Job
	FluxMonitorJob JobType = "fluxmonitor"
	// OffChainReportingJob defines an OCR Job
	OffChainReportingJob JobType = "offchainreporting"
)

// DirectRequestSpec defines the spec details of a DirectRequest Job
type DirectRequestSpec struct {
	ContractAddress  string    `json:"contractAddress"`
	OnChainJobSpecID string    `json:"OnChainJobSpecID"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// FluxMonitorSpec defines the spec details of a FluxMonitor Job
type FluxMonitorSpec struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// OffChainReportingSpec defines the spec details of a OffChainReporting Job
type OffChainReportingSpec struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PipelineSpec defines the spec details of the pipeline
type PipelineSpec struct {
	ID           int32  `json:"ID"`
	DotDAGSource string `json:"dotDagSource"`
}

// Job represents a V2 Job
type Job struct {
	JAID
	Name                  string                 `json:"name"`
	Type                  JobType                `json:"type"`
	DirectRequestSpec     *DirectRequestSpec     `json:"DirectRequestSpec"`
	FluxMonitorSpec       *FluxMonitorSpec       `json:"fluxMonitorSpec"`
	OffChainReportingSpec *OffChainReportingSpec `json:"offChainReportingOracleSpec"`
	PipelineSpec          PipelineSpec           `json:"pipelineSpec"`
}

// GetName implements the api2go EntityNamer interface
func (j Job) GetName() string {
	return "jobs"
}

// GetTasks extracts the tasks from the dependency graph
//
// TODO - Remove dependency on the pipeline package
func (j Job) GetTasks() ([]string, error) {
	types := []string{}
	dag := pipeline.NewTaskDAG()
	err := dag.UnmarshalText([]byte(j.PipelineSpec.DotDAGSource))
	if err != nil {
		return nil, err
	}

	tasks, err := dag.TasksInDependencyOrder()
	if err != nil {
		return nil, err
	}

	// Reverse the order as dependency tasks start from output and end at the
	// inputs.
	for i := len(tasks) - 1; i >= 0; i-- {
		t := tasks[i]
		types = append(types, fmt.Sprintf("%s %s", t.GetDotID(), t.Type()))
	}

	return types, nil
}

// FriendlyTasks returns the tasks
func (j Job) FriendlyTasks() []string {
	taskTypes, err := j.GetTasks()
	if err != nil {
		return []string{"error parsing DAG"}
	}

	return taskTypes
}

// FriendlyCreatedAt returns the created at timestamp of the spec which matches the
// type in RFC3339 format.
func (j Job) FriendlyCreatedAt() string {
	switch j.Type {
	case DirectRequestJob:
		if j.DirectRequestSpec != nil {
			return j.DirectRequestSpec.CreatedAt.Format(time.RFC3339)
		}
	case FluxMonitorJob:
		if j.FluxMonitorSpec != nil {
			return j.FluxMonitorSpec.CreatedAt.Format(time.RFC3339)
		}
	case OffChainReportingJob:
		if j.OffChainReportingSpec != nil {
			return j.OffChainReportingSpec.CreatedAt.Format(time.RFC3339)
		}
	default:
		return "unknown"
	}

	// This should never occur since the job should always have a spec matching
	// the type
	return "N/A"
}

// ToRow returns the job as a multiple rows per task
func (j Job) ToRow() [][]string {
	row := [][]string{}

	for _, t := range j.FriendlyTasks() {
		row = append(row, []string{
			j.ID,
			j.Name,
			j.Type.String(),
			t,
			j.FriendlyCreatedAt(),
		})
	}
	return row
}
