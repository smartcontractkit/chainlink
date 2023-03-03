// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package wrappers

import (
	"strings"
)

var _ error = &aggregate{}

type Errs struct{ Err error }

func (errs *Errs) Errored() bool { return errs.Err != nil }

func (errs *Errs) Add(errors ...error) {
	if errs.Err == nil {
		for _, err := range errors {
			if err != nil {
				errs.Err = err
				break
			}
		}
	}
}

// NewAggregate returns an aggregate error from a list of errors
func NewAggregate(errs []error) error {
	err := &aggregate{errs}
	if len(err.Errors()) == 0 {
		return nil
	}
	return err
}

type aggregate struct{ errs []error }

// Error returns the slice of errors with comma separated messsages wrapped in brackets
// [ error string 0 ], [ error string 1 ] ...
func (a *aggregate) Error() string {
	errString := make([]string, len(a.errs))
	for i, err := range a.errs {
		errString[i] = "[" + err.Error() + "]"
	}
	return strings.Join(errString, ",")
}

func (a *aggregate) Errors() []error {
	return a.errs
}
