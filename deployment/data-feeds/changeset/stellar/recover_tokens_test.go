package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecoverTokensChangeset(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedContractRefs(t, &env,
		contractRefSpec{CacheContract, testContractID, "test", "1.0.0"},
		contractRefSpec{ProxyContract, testProxyAddress, "test-proxy", "1.0.0"},
	)

	req := &RecoverTokensRequest{
		ChainSel:  testChainSel,
		Qualifier: "test",
		Version:   "1.0.0",
		Contract:  CacheContract,
		Token:     testContractID,
		To:        testAdmin,
		Amount:    1000,
	}
	require.NoError(t, RecoverTokens{}.VerifyPreconditions(env, req))

	_, err := RecoverTokens{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "recover_tokens", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)
	require.Len(t, inv.calls[0].Args, 3)

	// proxy path targets the proxy ref
	proxyReq := *req
	proxyReq.Qualifier = "test-proxy"
	proxyReq.Contract = ProxyContract
	require.NoError(t, RecoverTokens{}.VerifyPreconditions(env, &proxyReq))
	_, err = RecoverTokens{}.Apply(env, &proxyReq)
	require.NoError(t, err)
	require.Len(t, inv.calls, 2)
	require.Equal(t, "recover_tokens", inv.calls[1].Function)
	require.Equal(t, testProxyAddress, inv.calls[1].ContractID)

	// unsupported contract type must fail preconditions
	badContract := *req
	badContract.Contract = "NotAContract"
	require.Error(t, RecoverTokens{}.VerifyPreconditions(env, &badContract))

	// bad token address must fail preconditions
	badToken := *req
	badToken.Token = "not-a-key"
	require.Error(t, RecoverTokens{}.VerifyPreconditions(env, &badToken))

	// bad recipient address must fail preconditions
	badTo := *req
	badTo.To = "not-a-key"
	require.Error(t, RecoverTokens{}.VerifyPreconditions(env, &badTo))

	// non-positive amount must fail preconditions
	badAmount := *req
	badAmount.Amount = 0
	require.Error(t, RecoverTokens{}.VerifyPreconditions(env, &badAmount))

	// missing cache ref must fail preconditions
	badQualifier := *req
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, RecoverTokens{}.VerifyPreconditions(env, &badQualifier))
}
