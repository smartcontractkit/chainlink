package example

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

/**
TerminalErrorExampleChangeset demonstrates how to use Operations API to return a terminal error from an operation.
By returning an UnrecoverableError, the operation will not be retried by the framework and is returned immediately.
This is useful when an operation encounters an error that should not be retried.
*/

var _ cldf.ChangeSetV2[operations.EmptyInput] = TerminalErrorExampleChangeset{}

type TerminalErrorExampleChangeset struct{}

func (l TerminalErrorExampleChangeset) VerifyPreconditions(e cldf.Environment, config operations.EmptyInput) error {
	// perform any preconditions checks here
	return nil
}

func (l TerminalErrorExampleChangeset) Apply(e cldf.Environment, config operations.EmptyInput) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()

	_, err := operations.ExecuteOperation(e.OperationsBundle, TerminalErrorOperation, nil, operations.EmptyInput{})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

var TerminalErrorOperation = operations.NewOperation(
	"terminal-error-operation",
	semver.MustParse("1.0.0"),
	"Operation that returns a terminal error",
	func(b operations.Bundle, _ any, input operations.EmptyInput) (any, error) {
		// by returning an UnrecoverableError, the operation will not be retried by the framework
		return nil, operations.NewUnrecoverableError(errors.New("terminal error"))
	},
)
