package models

import (
	"fmt"
	"strings"
)

// DatabaseAccessError is an error that occurs during database access.
type DatabaseAccessError struct {
	msg string
}

func (e *DatabaseAccessError) Error() string { return e.msg }

// NewDatabaseAccessError returns a database access error.
func NewDatabaseAccessError(msg string) error {
	return &DatabaseAccessError{msg}
}

// ValidationError is an error that occurs during validation.
type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string { return e.msg }

// NewValidationError returns a validation error.
func NewValidationError(msg string, values ...interface{}) error {
	return &ValidationError{msg: fmt.Sprintf(msg, values...)}
}

type JSONAPIErrors struct {
	Errors []JSONAPIError `json:"errors"`
}

type JSONAPIError struct {
	Detail string `json:"detail"`
}

func NewJSONAPIErrors() *JSONAPIErrors {
	fe := JSONAPIErrors{
		Errors: []JSONAPIError{},
	}
	return &fe
}

func NewJSONAPIErrorsWith(detail string) *JSONAPIErrors {
	fe := NewJSONAPIErrors()
	fe.Errors = append(fe.Errors, JSONAPIError{Detail: detail})
	return fe
}

func (errors *JSONAPIErrors) Error() string {
	var messages []string
	for _, e := range errors.Errors {
		messages = append(messages, e.Detail)
	}
	return strings.Join(messages, ",")
}

func (e *JSONAPIErrors) Add(detail string) {
	e.Errors = append(e.Errors, JSONAPIError{Detail: detail})
}

func (errors *JSONAPIErrors) Merge(e error) {
	switch typed := e.(type) {
	case *JSONAPIErrors:
		errors.Errors = append(errors.Errors, typed.Errors...)
	default:
		errors.Add(e.Error())
	}
}

func (e *JSONAPIErrors) CoerceEmptyToNil() error {
	if len(e.Errors) == 0 {
		return nil
	}
	return e
}
