package models

import (
	"time"
)

// JobSpecError represents an asynchronous error caused by a JobSpec
type JobSpecError struct {
	ID          uint      `json:"id"`
	JobSpecID   *ID       `json:"-"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

// NewJobSpecError creates a new JobSpecError struct
func NewJobSpecError(jobSpecID *ID, description string) JobSpecError {
	return JobSpecError{
		JobSpecID:   jobSpecID,
		Description: description,
	}
}
