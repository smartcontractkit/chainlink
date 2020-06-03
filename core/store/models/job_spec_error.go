package models

type JobSpecError struct {
	JobSpecID   *ID    `json:"jobId"`
	Description string `json:"description"`
}

func NewJobSpecError(jobSpecID *ID, description string) JobSpecError {
	return JobSpecError{
		JobSpecID:   jobSpecID,
		Description: description,
	}
}
