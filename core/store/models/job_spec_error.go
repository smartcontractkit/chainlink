package models

import (
	"time"
)

// JobSpecError represents an asynchronous error caused by a JobSpec
type JobSpecError struct {
	ID          uint      `json:"id"`
	JobSpecID   *ID       `json:"-"`
	Description string    `json:"description"`
	Occurances  uint      `json:"occurances"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewJobSpecError creates a new JobSpecError struct
func NewJobSpecError(jobSpecID *ID, description string) JobSpecError {
	return JobSpecError{
		JobSpecID:   jobSpecID,
		Description: description,
		Occurances:  1,
	}
}
