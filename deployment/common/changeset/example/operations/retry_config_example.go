package example

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
)

/**
DisableRetryExampleChangeset demonstrates how to use Operations API to disable retry for an operation.
UpdateInputExampleChangeset demonstrates how to use Operations API to update input for an operation (eg for changing gas limit)
*/

var _ deployment.ChangeSetV2[operations.EmptyInput] = DisableRetryExampleChangeset{}

type DisableRetryExampleChangeset struct{}

func (l DisableRetryExampleChangeset) VerifyPreconditions(e deployment.Environment, config operations.EmptyInput) error {
	// perform any preconditions checks here
	return nil
}

func (l DisableRetryExampleChangeset) Apply(e deployment.Environment, config operations.EmptyInput) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()

	operationInput := SuccessFailOperationInput{ShouldFail: true}

	// Disable retry for this operation
	// If retry was not disable, the operation would be retried for 10 times with exponential backoff.
	_, err := operations.ExecuteOperation(e.OperationsBundle, SuccessFailOperation, nil, operationInput,
		operations.WithRetryConfig[SuccessFailOperationInput, any](operations.RetryConfig[SuccessFailOperationInput, any]{
			DisableRetry: true,
		}))
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

var _ deployment.ChangeSetV2[operations.EmptyInput] = UpdateInputExampleChangeset{}

type UpdateInputExampleChangeset struct{}

func (l UpdateInputExampleChangeset) VerifyPreconditions(e deployment.Environment, config operations.EmptyInput) error {
	// perform any preconditions checks here
	return nil
}

func (l UpdateInputExampleChangeset) Apply(e deployment.Environment, config operations.EmptyInput) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()

	operationInput := SuccessFailOperationInput{ShouldFail: true}

	// Retry operation with updated input
	// This operation will fail once and then succeed because the input was updated
	_, err := operations.ExecuteOperation(e.OperationsBundle, SuccessFailOperation, nil, operationInput,
		operations.WithRetryConfig[SuccessFailOperationInput, any](operations.RetryConfig[SuccessFailOperationInput, any]{
			InputHook: func(input SuccessFailOperationInput, _ any) SuccessFailOperationInput {
				// Update input to false, so it stops failing
				input.ShouldFail = false
				return input
			},
		}))
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

type SuccessFailOperationInput struct {
	ShouldFail bool
}

// SuccessFailOperation is an operation that always fails if ShouldFail is true.
// Else it succeeds.
var SuccessFailOperation = operations.NewOperation(
	"success-fail-operation",
	semver.MustParse("1.0.0"),
	"Operation that always fails",
	func(b operations.Bundle, _ any, input SuccessFailOperationInput) (any, error) {
		if input.ShouldFail {
			return nil, errors.New("operation failed")
		}
		return nil, nil
	},
)
