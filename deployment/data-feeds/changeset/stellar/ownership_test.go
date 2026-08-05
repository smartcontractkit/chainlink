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

	base := OwnershipRequest{
		ChainSel:  testChainSel,
		Qualifier: "test-cache",
		Version:   "1.0.0",
		Contract:  CacheContract,
	}
	transferReq := &TransferOwnershipRequest{
		OwnershipRequest: base,
		NewOwner:         testAdmin,
		LiveUntilLedger:  100,
	}

	require.NoError(t, TransferOwnership{}.VerifyPreconditions(env, transferReq))
	_, err := TransferOwnership{}.Apply(env, transferReq)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "transfer_ownership", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)
	require.Len(t, inv.calls[0].Args, 2)

	require.NoError(t, AcceptOwnership{}.VerifyPreconditions(env, &base))
	_, err = AcceptOwnership{}.Apply(env, &base)
	require.NoError(t, err)
	require.Len(t, inv.calls, 2)
	require.Equal(t, "accept_ownership", inv.calls[1].Function)
	require.Equal(t, testContractID, inv.calls[1].ContractID)
	require.Empty(t, inv.calls[1].Args)

	// same ops against the proxy: IsProxy selects the proxy client and the
	// resolved contractID is the proxy's address, not the cache's.
	proxyBase := base
	proxyBase.Qualifier = "test-proxy"
	proxyBase.Contract = ProxyContract
	proxyTransfer := &TransferOwnershipRequest{
		OwnershipRequest: proxyBase,
		NewOwner:         testAdmin,
		LiveUntilLedger:  100,
	}

	require.NoError(t, TransferOwnership{}.VerifyPreconditions(env, proxyTransfer))
	_, err = TransferOwnership{}.Apply(env, proxyTransfer)
	require.NoError(t, err)
	require.Equal(t, "transfer_ownership", inv.calls[2].Function)
	require.Equal(t, testProxyAddress, inv.calls[2].ContractID)

	require.NoError(t, AcceptOwnership{}.VerifyPreconditions(env, &proxyBase))
	_, err = AcceptOwnership{}.Apply(env, &proxyBase)
	require.NoError(t, err)
	require.Equal(t, "accept_ownership", inv.calls[3].Function)
	require.Equal(t, testProxyAddress, inv.calls[3].ContractID)

	// unknown Contract type must fail preconditions for both changesets
	badContract := base
	badContract.Contract = datastore.ContractType("SomethingElse")
	require.Error(t, AcceptOwnership{}.VerifyPreconditions(env, &badContract))
	badTransfer := *transferReq
	badTransfer.Contract = badContract.Contract
	require.Error(t, TransferOwnership{}.VerifyPreconditions(env, &badTransfer))

	// TransferOwnership requires a valid NewOwner and nonzero LiveUntilLedger;
	// AcceptOwnership's request doesn't carry those fields at all.
	missingOwner := *transferReq
	missingOwner.NewOwner = ""
	require.Error(t, TransferOwnership{}.VerifyPreconditions(env, &missingOwner))
	require.NoError(t, AcceptOwnership{}.VerifyPreconditions(env, &missingOwner.OwnershipRequest))

	badOwner := *transferReq
	badOwner.NewOwner = "not-a-key"
	require.Error(t, TransferOwnership{}.VerifyPreconditions(env, &badOwner))

	zeroLedger := *transferReq
	zeroLedger.LiveUntilLedger = 0
	require.Error(t, TransferOwnership{}.VerifyPreconditions(env, &zeroLedger))

	// missing address ref must fail preconditions
	badQualifier := *transferReq
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, TransferOwnership{}.VerifyPreconditions(env, &badQualifier))
}
