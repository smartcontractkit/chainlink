package models

import (
	"encoding/json"
	"time"
)

// JobSpecError represents an asynchronous error caused by a JobSpec
type JobSpecError struct {
	JobSpec     JobSpec
	ID          uint      `json:"id,omitempty"`
	JobSpecID   *ID       `json:"job_spec_id,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

// NewJobSpecError creates a new JobSpecError struct
func NewJobSpecError(jobSpecID *ID, description string) JobSpecError {
	return JobSpecError{
		JobSpecID:   jobSpecID,
		Description: description,
	}
}

func (jse JobSpecError) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			ID          uint
			Description string
			CreatedAt   time.Time
		}{
			ID:          jse.ID,
			Description: jse.Description,
			CreatedAt:   jse.CreatedAt,
		},
	)
}
