package example

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
)

/**
TerminalErrorExampleChangeset demonstrates how to use Operations API to return a terminal error from an operation.
By returning an UnrecoverableError, the operation will not be retried by the framework and is returned immediately.
This is useful when an operation encounters an error that should not be retried.
*/

var _ deployment.ChangeSetV2[operations.EmptyInput] = TerminalErrorExampleChangeset{}

type TerminalErrorExampleChangeset struct{}

func (l TerminalErrorExampleChangeset) VerifyPreconditions(e deployment.Environment, config operations.EmptyInput) error {
	// perform any preconditions checks here
	return nil
}

func (l TerminalErrorExampleChangeset) Apply(e deployment.Environment, config operations.EmptyInput) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()

	_, err := operations.ExecuteOperation(e.OperationsBundle, TerminalErrorOperation, nil, operations.EmptyInput{})
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
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
