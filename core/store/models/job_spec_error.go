package models

import (
	"time"
)

// JobSpecError represents an asynchronous error caused by a JobSpec
type JobSpecError struct {
	ID          *ID       `json:"id"`
	JobSpecID   *ID       `json:"-"`
	Description string    `json:"description"`
	Occurences  uint      `json:"occurences"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewJobSpecError creates a new JobSpecError struct
func NewJobSpecError(jobSpecID *ID, description string) JobSpecError {
	return JobSpecError{
		ID:          NewID(),
		JobSpecID:   jobSpecID,
		Description: description,
		Occurences:  1,
	}
}
