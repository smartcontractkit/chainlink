package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// BulkTaskStatus indicates what a bulk task is doing.
type BulkTaskStatus string

const (
	// BulkTaskStatusInProgress is the default state of any run status.
	BulkTaskStatusInProgress = BulkTaskStatus("")
	// BulkTaskStatusErrored means a bulk task stopped because it encountered an error.
	BulkTaskStatusErrored = BulkTaskStatus("errored")
	// BulkTaskStatusCompleted means a bulk task finished.
	BulkTaskStatusCompleted = BulkTaskStatus("completed")
)

// Value returns this instance serialized for database storage.
func (t BulkTaskStatus) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan reads the database value and returns an instance.
func (t *BulkTaskStatus) Scan(value interface{}) error {
	temp, ok := value.(string)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to BulkTaskStatus", value, value)
	}

	*t = BulkTaskStatus(temp)
	return nil
}

// BulkDeleteRunRequest describes the query for deletion of runs
type BulkDeleteRunRequest struct {
	ID            uint                `gorm:"primary_key"`
	Status        RunStatusCollection `json:"status" gorm:"type:text"`
	UpdatedBefore time.Time           `json:"updatedBefore"`
}

// ValidateBulkDeleteRunRequest returns a task from a request to make a task
func ValidateBulkDeleteRunRequest(request *BulkDeleteRunRequest) error {
	for _, status := range request.Status {
		if status != RunStatusCompleted && status != RunStatusErrored {
			return fmt.Errorf("cannot delete Runs with status %s", status)
		}
	}

	return nil
}

// RunStatusCollection is an array of RunStatus.
type RunStatusCollection []RunStatus

// ToStrings returns a copy of RunStatusCollection as an array of strings.
func (r RunStatusCollection) ToStrings() []string {
	// Unable to convert copy-free without unsafe:
	// https://stackoverflow.com/a/48554123/639773
	converted := make([]string, len(r))
	for i, e := range r {
		converted[i] = string(e)
	}
	return converted
}

// Value returns this instance serialized for database storage.
func (r RunStatusCollection) Value() (driver.Value, error) {
	return strings.Join(r.ToStrings(), ","), nil
}

// Scan reads the database value and returns an instance.
func (r *RunStatusCollection) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to RunStatusCollection", value, value)
	}

	if len(str) == 0 {
		return nil
	}

	arr := strings.Split(str, ",")
	collection := make(RunStatusCollection, len(arr))
	for i, r := range arr {
		collection[i] = RunStatus(r)
	}
	*r = collection
	return nil
}
