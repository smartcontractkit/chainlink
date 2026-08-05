package stellar

import (
	"testing"

	"github.com/stellar/go-stellar-sdk/xdr"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-stellar/bindings/scval"
)

func TestUpgradeChangeset(t *testing.T) {
	env, inv, dep := newTestEnv(t)
	seedContractRefs(t, &env,
		contractRefSpec{CacheContract, testContractID, "test", "1.0.0"},
		contractRefSpec{ProxyContract, testProxyAddress, "test-proxy", "1.0.0"},
	)

	var wantHash xdr.Hash
	for i := range wantHash {
		wantHash[i] = byte(i + 1)
	}
	dep.wasmHash = wantHash

	req := &UpgradeRequest{
		ChainSel:  testChainSel,
		Qualifier: "test",
		Version:   "1.0.0",
		Contract:  CacheContract,
		WasmPath:  writeDummyWasm(t, "new_cache.wasm"),
	}
	require.NoError(t, Upgrade{}.VerifyPreconditions(env, req))

	_, err := Upgrade{}.Apply(env, req)
	require.NoError(t, err)

	// the upload happened, against the requested wasm path
	require.Len(t, dep.uploads, 1)
	require.Equal(t, req.WasmPath, dep.uploads[0])

	// the upgrade invocation carried the fake's uploaded hash through
	require.Len(t, inv.calls, 1)
	require.Equal(t, "upgrade", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)
	require.Len(t, inv.calls[0].Args, 1)

	gotHash, err := scval.Bytes32FromScVal(inv.calls[0].Args[0])
	require.NoError(t, err)
	require.Equal(t, [32]byte(wantHash), gotHash)

	// proxy path targets the proxy ref
	proxyReq := *req
	proxyReq.Qualifier = "test-proxy"
	proxyReq.Contract = ProxyContract
	require.NoError(t, Upgrade{}.VerifyPreconditions(env, &proxyReq))
	_, err = Upgrade{}.Apply(env, &proxyReq)
	require.NoError(t, err)
	require.Len(t, inv.calls, 2)
	require.Equal(t, "upgrade", inv.calls[1].Function)
	require.Equal(t, testProxyAddress, inv.calls[1].ContractID)

	// unsupported contract type must fail preconditions
	badContract := *req
	badContract.Contract = "NotAContract"
	require.Error(t, Upgrade{}.VerifyPreconditions(env, &badContract))

	// missing wasm path must fail preconditions
	badPath := *req
	badPath.WasmPath = "/does/not/exist.wasm"
	require.Error(t, Upgrade{}.VerifyPreconditions(env, &badPath))

	// missing cache ref must fail preconditions
	badQualifier := *req
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, Upgrade{}.VerifyPreconditions(env, &badQualifier))
}
