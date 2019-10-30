package models

import (
	null "gopkg.in/guregu/null.v3"
)

// RunResult keeps track of the outcome of a TaskRun or JobRun. It stores the
// Data and ErrorMessage.
type RunResult struct {
	ID           uint        `json:"-" gorm:"primary_key;auto_increment"`
	Data         JSON        `json:"data" gorm:"type:text"`
	ErrorMessage null.String `json:"error"`
}
