package models

import (
	"time"
)

// TODO - RYAN - constants for error messages

// JobSpecError represents an asynchronous error caused by a JobSpec
type JobSpecError struct {
	JobSpec     JobSpec
	ID          uint
	JobSpecID   *ID
	Description string
	CreatedAt   time.Time
}

// NewJobSpecError creates a new JobSpecError struct
func NewJobSpecError(jobSpecID *ID, description string) JobSpecError {
	return JobSpecError{
		JobSpecID:   jobSpecID,
		Description: description,
	}
}
