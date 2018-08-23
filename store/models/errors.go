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

// JSONAPIErrors holds errors conforming to the JSONAPI spec.
type JSONAPIErrors struct {
	Errors []JSONAPIError `json:"errors"`
}

// JSONAPIError is an individual JSONAPI Error.
type JSONAPIError struct {
	Detail string `json:"detail"`
}

// NewJSONAPIErrors creates an instance of JSONAPIErrors, with the intention
// of managing a collection of them.
func NewJSONAPIErrors() *JSONAPIErrors {
	fe := JSONAPIErrors{
		Errors: []JSONAPIError{},
	}
	return &fe
}

// NewJSONAPIErrorsWith creates an instance of JSONAPIErrors populated with this
// single detail.
func NewJSONAPIErrorsWith(detail string) *JSONAPIErrors {
	fe := NewJSONAPIErrors()
	fe.Errors = append(fe.Errors, JSONAPIError{Detail: detail})
	return fe
}

// Error collapses the collection of errors into a collection of comma separated
// strings.
func (jae *JSONAPIErrors) Error() string {
	var messages []string
	for _, e := range jae.Errors {
		messages = append(messages, e.Detail)
	}
	return strings.Join(messages, ",")
}

// Add adds a new error to JSONAPIErrors with the passed detail.
func (jae *JSONAPIErrors) Add(detail string) {
	jae.Errors = append(jae.Errors, JSONAPIError{Detail: detail})
}

// Merge combines the arrays of the passed error if it is of type JSONAPIErrors,
// otherwise simply adds a single error with the error string as detail.
func (jae *JSONAPIErrors) Merge(e error) {
	switch typed := e.(type) {
	case *JSONAPIErrors:
		jae.Errors = append(jae.Errors, typed.Errors...)
	default:
		jae.Add(e.Error())
	}
}

// CoerceEmptyToNil will return nil if JSONAPIErrors has no errors.
func (jae *JSONAPIErrors) CoerceEmptyToNil() error {
	if len(jae.Errors) == 0 {
		return nil
	}
	return jae
}
