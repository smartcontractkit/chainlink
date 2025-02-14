package deployment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type MyChangeSet struct{}

func (m MyChangeSet) Apply(e Environment, config uint64) (ChangesetOutput, error) {
	return ChangesetOutput{AddressBook: NewMemoryAddressBook()}, nil
}
func (m MyChangeSet) VerifyPreconditions(e Environment, config uint64) error {
	return nil
}

func TestChangeSetNormalType(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)

	var cs ChangeSetV2[uint64] = MyChangeSet{}
	verify := cs.VerifyPreconditions(e, 5)
	require.NoError(t, verify)
	out, _ := cs.Apply(e, 5)
	require.Equal(t, NewMemoryAddressBook(), out.AddressBook)
}

func TestChangeSetConstructionComposedFromLambdas(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)

	var cs = CreateChangeSet(
		// Don't do this in real life, this is for a test. Make nice tested functions.
		func(e Environment, config string) (ChangesetOutput, error) {
			return ChangesetOutput{AddressBook: NewMemoryAddressBook()}, nil
		},
		func(e Environment, config string) error {
			return nil
		},
	)
	verify := cs.VerifyPreconditions(e, "foo")
	require.NoError(t, verify)
	out, _ := cs.Apply(e, "foo")
	require.Equal(t, NewMemoryAddressBook(), out.AddressBook)
}

var fakeChangeSet = CreateChangeSet(oldSchool, oldSchoolVerify)

func oldSchool(e Environment, config uint32) (ChangesetOutput, error) {
	return ChangesetOutput{AddressBook: NewMemoryAddressBook()}, nil
}
func oldSchoolVerify(e Environment, _ uint32) error {
	return nil
}

// This is likely the best example of how to use this API.
func TestChangeSetComposedType(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)

	verify := fakeChangeSet.VerifyPreconditions(e, 5)
	require.NoError(t, verify)
	out, _ := fakeChangeSet.Apply(e, 5)
	require.Equal(t, NewMemoryAddressBook(), out.AddressBook)
}

// TestChangeSetLegacyFunction tests using legacy ChangeSet functions (but just naturally conforming to the type,
// via duck-typing, in the new wrapper.
func TestChangeSetLegacyFunctionWithStandardChangeSetFunction(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)
	var cs = CreateLegacyChangeSet(oldSchool)
	verify := cs.VerifyPreconditions(e, 5)
	require.NoError(t, verify)
	out, _ := cs.Apply(e, 5)
	require.Equal(t, NewMemoryAddressBook(), out.AddressBook)
}

// TestChangeSetLegacyFunction tests using legacy ChangeSet (strongly declared as a ChangeSet[C]) in the wrapper.
func TestChangeSetLegacyFunction(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)
	var csFunc ChangeSet[uint32] = oldSchool // Cast to a ChangeSet and use in CreateLegacyChangeSet
	var cs = CreateLegacyChangeSet(csFunc)
	verify := cs.VerifyPreconditions(e, 5)
	require.NoError(t, verify)
	out, _ := cs.Apply(e, 5)
	require.Equal(t, NewMemoryAddressBook(), out.AddressBook)
}

func NewNoopEnvironment(t *testing.T) Environment {
	return *NewEnvironment(
		"noop",
		logger.TestLogger(t),
		NewMemoryAddressBook(),
		map[uint64]Chain{},
		map[uint64]SolChain{},
		[]string{},
		nil,
		func() context.Context { return tests.Context(t) },
		XXXGenerateTestOCRSecrets(),
	)
}
