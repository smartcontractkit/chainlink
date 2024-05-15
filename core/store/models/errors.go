package models

import (
	"errors"
	"strings"
)

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
	var jsonErr *JSONAPIErrors
	if errors.As(e, &jsonErr) {
		jae.Errors = append(jae.Errors, jsonErr.Errors...)
		return
	}
	jae.Add(e.Error())
}

// CoerceEmptyToNil will return nil if JSONAPIErrors has no errors.
func (jae *JSONAPIErrors) CoerceEmptyToNil() error {
	if len(jae.Errors) == 0 {
		return nil
	}
	return jae
}
