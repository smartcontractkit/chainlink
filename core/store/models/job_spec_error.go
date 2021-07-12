package models

import (
	"time"
)

// JobSpecError represents an asynchronous error caused by a JobSpec
type JobSpecError struct {
	ID          int64     `json:"id"`
	JobSpecID   JobID     `json:"-"`
	Description string    `json:"description"`
	Occurrences uint      `json:"occurrences"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewJobSpecError creates a new JobSpecError struct
func NewJobSpecError(jobSpecID JobID, description string) JobSpecError {
	return JobSpecError{
		JobSpecID:   jobSpecID,
		Description: description,
		Occurrences: 1,
	}
}
