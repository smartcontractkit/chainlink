package shared

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

const testChainSelector = 5009297550715157269

// testEnvironment builds a minimal environment for preflight tests: a memory address book, a
// sealed memory datastore, and a test logger. It mirrors the framework's NewNoopEnvironment
// without pulling in the full engine test runtime.
func testEnvironment(t *testing.T) deployment.Environment {
	t.Helper()

	return deployment.Environment{
		Logger:            logger.Test(t),
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		DataStore:         datastore.NewMemoryDataStore().Seal(),
	}
}

func poolRef(addr, qualifier string) datastore.AddressRef {
	return datastore.AddressRef{
		ChainSelector: testChainSelector,
		Address:       addr,
		Type:          datastore.ContractType("BurnMintTokenPool"),
		Version:       semver.MustParse("1.5.0"),
		Qualifier:     qualifier,
	}
}

// TestValidateAddressRefsRejectsSharedKey verifies the case the check exists for: a changeset
// about to deploy two contracts of one type and version that have not been given distinct
// qualifiers. Called from VerifyPreconditions this fails before anything is deployed.
func TestValidateAddressRefsRejectsSharedKey(t *testing.T) {
	t.Parallel()
	e := testEnvironment(t)

	err := ValidateAddressRefs(e, []datastore.AddressRef{
		poolRef("0xaaa", ""),
		poolRef("0xbbb", ""),
	})

	require.ErrorIs(t, err, deployment.ErrInvalidConfig)

	detail, ok := AsAddressRefValidationError(err)
	require.True(t, ok)
	require.Len(t, detail.Conflicts, 1)
	assert.Equal(t, "0xaaa", detail.Conflicts[0].First.Address)
	assert.Equal(t, "0xbbb", detail.Conflicts[0].Second.Address)
}

// TestValidateAddressRefsAcceptsDistinctQualifiers verifies the same two contracts pass once
// the caller has said which is which.
func TestValidateAddressRefsAcceptsDistinctQualifiers(t *testing.T) {
	t.Parallel()
	e := testEnvironment(t)

	require.NoError(t, ValidateAddressRefs(e, []datastore.AddressRef{
		poolRef("0xaaa", "LINK"),
		poolRef("0xbbb", "USDC"),
	}))
}

// TestValidateAddressRefsRejectsVersionless verifies a ref with no version is reported: the
// datastore key is built from the version, so such a ref has no key at all.
func TestValidateAddressRefsRejectsVersionless(t *testing.T) {
	t.Parallel()
	e := testEnvironment(t)

	err := ValidateAddressRefs(e, []datastore.AddressRef{{
		ChainSelector: testChainSelector,
		Address:       "0xaaa",
		Type:          datastore.ContractType("Router"),
	}})

	detail, ok := AsAddressRefValidationError(err)
	require.True(t, ok)
	require.Len(t, detail.Versionless, 1)
}

// TestValidateAddressRefsRepeatedAddress verifies that the same address listed twice under one
// key is not a conflict. It is one contract, declared twice.
func TestValidateAddressRefsRepeatedAddress(t *testing.T) {
	t.Parallel()
	e := testEnvironment(t)

	require.NoError(t, ValidateAddressRefs(e, []datastore.AddressRef{
		poolRef("0xaaa", "LINK"),
		poolRef("0xaaa", "LINK"),
	}))
}

// TestValidateAddressRefsSupersedingEnv verifies that replacing a contract the environment
// already holds is allowed by default, since that is what a redeploy is, and rejected by the
// strict variant for changesets that only ever deploy new contracts.
func TestValidateAddressRefsSupersedingEnv(t *testing.T) {
	t.Parallel()
	newEnv := func(t *testing.T) deployment.Environment {
		t.Helper()

		e := testEnvironment(t)
		mutable, ok := e.DataStore.Addresses().(datastore.MutableAddressRefStore)
		require.True(t, ok)
		require.NoError(t, mutable.Upsert(poolRef("0xold", "LINK")))

		return e
	}

	t.Run("allowed by default", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, ValidateAddressRefs(newEnv(t), []datastore.AddressRef{poolRef("0xnew", "LINK")}))
	})

	t.Run("rejected when strict", func(t *testing.T) {
		t.Parallel()
		err := ValidateAddressRefsStrict(newEnv(t), []datastore.AddressRef{poolRef("0xnew", "LINK")})
		require.ErrorIs(t, err, deployment.ErrInvalidEnvironment)
		assert.ErrorContains(t, err, "0xold")
	})

	t.Run("same address is not a replacement", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, ValidateAddressRefsStrict(newEnv(t), []datastore.AddressRef{poolRef("0xold", "LINK")}))
	})
}

// TestValidateAddressRefsRejectsSharedKeyBeforeDeploy verifies the check works in the situation
// it is documented for: called from VerifyPreconditions, where nothing has been deployed and the
// planned refs therefore carry no address. Two instances of one type with no distinguishing
// qualifier must be rejected here, not discovered once the contracts are already on chain.
func TestValidateAddressRefsRejectsSharedKeyBeforeDeploy(t *testing.T) {
	t.Parallel()
	e := testEnvironment(t)

	err := ValidateAddressRefs(e, []datastore.AddressRef{
		poolRef("", ""),
		poolRef("", ""),
	})

	require.ErrorIs(t, err, deployment.ErrInvalidConfig)

	detail, ok := AsAddressRefValidationError(err)
	require.True(t, ok)
	require.Len(t, detail.Conflicts, 1)
}

// TestValidateAddressRefsAcceptsDistinctQualifiersBeforeDeploy verifies the same undeployed refs
// pass once the caller has given them distinct qualifiers, so the rule above does not simply
// reject every ref that has no address yet.
func TestValidateAddressRefsAcceptsDistinctQualifiersBeforeDeploy(t *testing.T) {
	t.Parallel()
	e := testEnvironment(t)

	require.NoError(t, ValidateAddressRefs(e, []datastore.AddressRef{
		poolRef("", "LINK"),
		poolRef("", "USDC"),
	}))
}

// TestValidateAddressRefsStrictUndeployedOverEnv verifies that an undeployed ref claiming a key
// the environment already holds is still caught by the strict variant. The address is unknown,
// so it cannot be shown to be the contract already recorded there.
func TestValidateAddressRefsStrictUndeployedOverEnv(t *testing.T) {
	t.Parallel()
	e := testEnvironment(t)
	mutable, ok := e.DataStore.Addresses().(datastore.MutableAddressRefStore)
	require.True(t, ok)
	require.NoError(t, mutable.Upsert(poolRef("0xold", "LINK")))

	err := ValidateAddressRefsStrict(e, []datastore.AddressRef{poolRef("", "LINK")})
	require.ErrorIs(t, err, deployment.ErrInvalidEnvironment)
	assert.ErrorContains(t, err, "0xold")
}
