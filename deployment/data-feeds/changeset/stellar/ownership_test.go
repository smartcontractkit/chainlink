package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

func TestOwnershipChangesets(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedContractRefs(t, &env,
		contractRefSpec{CacheContract, testContractID, "test-cache", "1.0.0"},
		contractRefSpec{ProxyContract, testProxyAddress, "test-proxy", "1.0.0"},
	)

	cacheReq := &OwnershipRequest{
		ChainSel:        testChainSel,
		Qualifier:       "test-cache",
		Version:         "1.0.0",
		Contract:        CacheContract,
		NewOwner:        testAdmin,
		LiveUntilLedger: 100,
	}

	require.NoError(t, TransferOwnership{}.VerifyPreconditions(env, cacheReq))
	_, err := TransferOwnership{}.Apply(env, cacheReq)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "transfer_ownership", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)
	require.Len(t, inv.calls[0].Args, 2)

	require.NoError(t, AcceptOwnership{}.VerifyPreconditions(env, cacheReq))
	_, err = AcceptOwnership{}.Apply(env, cacheReq)
	require.NoError(t, err)
	require.Len(t, inv.calls, 2)
	require.Equal(t, "accept_ownership", inv.calls[1].Function)
	require.Equal(t, testContractID, inv.calls[1].ContractID)
	require.Empty(t, inv.calls[1].Args)

	require.NoError(t, RenounceOwnership{}.VerifyPreconditions(env, cacheReq))
	_, err = RenounceOwnership{}.Apply(env, cacheReq)
	require.NoError(t, err)
	require.Len(t, inv.calls, 3)
	require.Equal(t, "renounce_ownership", inv.calls[2].Function)
	require.Equal(t, testContractID, inv.calls[2].ContractID)

	// same three ops against the proxy: IsProxy selects the proxy client and
	// the resolved contractID is the proxy's address, not the cache's.
	proxyReq := *cacheReq
	proxyReq.Qualifier = "test-proxy"
	proxyReq.Contract = ProxyContract

	require.NoError(t, TransferOwnership{}.VerifyPreconditions(env, &proxyReq))
	_, err = TransferOwnership{}.Apply(env, &proxyReq)
	require.NoError(t, err)
	require.Equal(t, "transfer_ownership", inv.calls[3].Function)
	require.Equal(t, testProxyAddress, inv.calls[3].ContractID)

	require.NoError(t, AcceptOwnership{}.VerifyPreconditions(env, &proxyReq))
	_, err = AcceptOwnership{}.Apply(env, &proxyReq)
	require.NoError(t, err)
	require.Equal(t, "accept_ownership", inv.calls[4].Function)
	require.Equal(t, testProxyAddress, inv.calls[4].ContractID)

	require.NoError(t, RenounceOwnership{}.VerifyPreconditions(env, &proxyReq))
	_, err = RenounceOwnership{}.Apply(env, &proxyReq)
	require.NoError(t, err)
	require.Equal(t, "renounce_ownership", inv.calls[5].Function)
	require.Equal(t, testProxyAddress, inv.calls[5].ContractID)

	// unknown Contract type must fail preconditions for all three changesets
	badContract := *cacheReq
	badContract.Contract = datastore.ContractType("SomethingElse")
	require.Error(t, TransferOwnership{}.VerifyPreconditions(env, &badContract))
	require.Error(t, AcceptOwnership{}.VerifyPreconditions(env, &badContract))
	require.Error(t, RenounceOwnership{}.VerifyPreconditions(env, &badContract))

	// TransferOwnership requires a valid NewOwner and nonzero LiveUntilLedger;
	// Accept/Renounce don't touch those fields so the same request stays valid.
	missingOwner := *cacheReq
	missingOwner.NewOwner = ""
	require.Error(t, TransferOwnership{}.VerifyPreconditions(env, &missingOwner))
	require.NoError(t, AcceptOwnership{}.VerifyPreconditions(env, &missingOwner))

	badOwner := *cacheReq
	badOwner.NewOwner = "not-a-key"
	require.Error(t, TransferOwnership{}.VerifyPreconditions(env, &badOwner))

	zeroLedger := *cacheReq
	zeroLedger.LiveUntilLedger = 0
	require.Error(t, TransferOwnership{}.VerifyPreconditions(env, &zeroLedger))

	// missing address ref must fail preconditions
	badQualifier := *cacheReq
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, TransferOwnership{}.VerifyPreconditions(env, &badQualifier))
}
