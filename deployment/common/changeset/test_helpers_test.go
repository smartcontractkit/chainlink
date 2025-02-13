package changeset

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestChangeSetLegacyFunction_PassingCase(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)

	executedCs := false
	executedValidator := false

	csv2 := deployment.CreateChangeSet(
		func(e deployment.Environment, config uint32) (deployment.ChangesetOutput, error) {
			executedCs = true
			return deployment.ChangesetOutput{AddressBook: deployment.NewMemoryAddressBook()}, nil
		},
		func(e deployment.Environment, config uint32) error {
			executedValidator = true
			return nil
		},
	)
	assert.False(t, executedCs, "Not expected to have executed the changeset yet")
	assert.False(t, executedValidator, "Not expected to have executed the validator yet")
	_, err := Apply(t, e, nil, Configure(csv2, 1))
	assert.True(t, executedCs, "Validator should have returned nil, allowing changeset execution")
	assert.True(t, executedValidator, "Not expected to have executed the validator yet")
	assert.NoError(t, err)
}

func TestChangeSetLegacyFunction_ErrorCase(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)

	executedCs := false
	executedValidator := false

	csv2 := deployment.CreateChangeSet(
		func(e deployment.Environment, config uint32) (deployment.ChangesetOutput, error) {
			executedCs = true
			return deployment.ChangesetOutput{AddressBook: deployment.NewMemoryAddressBook()}, nil
		},
		func(e deployment.Environment, config uint32) error {
			executedValidator = true
			return errors.New("you shall not pass")
		},
	)
	assert.False(t, executedCs, "Not expected to have executed the changeset yet")
	assert.False(t, executedValidator, "Not expected to have executed the validator yet")
	_, err := Apply(t, e, nil, Configure(csv2, 1))
	assert.False(t, executedCs, "Validator should have fired, preventing changeset execution")
	assert.True(t, executedValidator, "Not expected to have executed the validator yet")
	assert.Equal(t, "failed to apply changeset at index 0: you shall not pass", err.Error())
}

func NewNoopEnvironment(t *testing.T) deployment.Environment {
	return *deployment.NewEnvironment(
		"noop",
		logger.TestLogger(t),
		deployment.NewMemoryAddressBook(),
		map[uint64]deployment.Chain{},
		map[uint64]deployment.SolChain{},
		[]string{},
		nil,
		func() context.Context { return tests.Context(t) },
		deployment.XXXGenerateTestOCRSecrets(),
	)
}
