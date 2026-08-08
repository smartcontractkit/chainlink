package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetProxyCacheChangeset(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedContractRefs(t, &env,
		contractRefSpec{ProxyContract, testProxyAddress, "test-proxy", "1.0.0"},
		contractRefSpec{CacheContract, testContractID, "test-cache", "1.0.0"},
	)

	req := &SetProxyCacheRequest{
		ChainSel:       testChainSel,
		Qualifier:      "test-proxy",
		Version:        "1.0.0",
		CacheQualifier: "test-cache",
	}
	require.NoError(t, SetProxyCache{}.VerifyPreconditions(env, req))

	out, err := SetProxyCache{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "set_cache", inv.calls[0].Function)
	require.Equal(t, testProxyAddress, inv.calls[0].ContractID)
	require.Len(t, inv.calls[0].Args, 1)

	// the metadata mirror records the proxy's new cache target
	meta := outputMetadata(t, out, testProxyAddress)
	require.Equal(t, testContractID, meta.Cache)

	// missing cache ref must fail preconditions (proxy ref alone isn't enough)
	proxyOnly, _, _ := newTestEnv(t)
	seedProxyRef(t, &proxyOnly, testProxyAddress, "test-proxy", "1.0.0")
	require.Error(t, SetProxyCache{}.VerifyPreconditions(proxyOnly, req))

	// missing proxy ref must fail preconditions (cache ref alone isn't enough)
	cacheOnly, _, _ := newTestEnv(t)
	seedCacheRef(t, &cacheOnly, testContractID, "test-cache", "1.0.0")
	require.Error(t, SetProxyCache{}.VerifyPreconditions(cacheOnly, req))
}

// CacheVersion resolves a cache recorded under a different version than the
// proxy; empty CacheVersion falls back to Version.
func TestSetProxyCacheChangeset_CrossVersion(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedContractRefs(t, &env,
		contractRefSpec{ProxyContract, testProxyAddress, "test-proxy", "1.1.0"},
		contractRefSpec{CacheContract, testContractID, "test-cache", "1.0.0"},
	)

	req := &SetProxyCacheRequest{
		ChainSel:       testChainSel,
		Qualifier:      "test-proxy",
		Version:        "1.1.0",
		CacheQualifier: "test-cache",
		CacheVersion:   "1.0.0",
	}
	require.NoError(t, SetProxyCache{}.VerifyPreconditions(env, req))

	out, err := SetProxyCache{}.Apply(env, req)
	require.NoError(t, err)
	require.Equal(t, "set_cache", inv.calls[0].Function)

	meta := outputMetadata(t, out, testProxyAddress)
	require.Equal(t, testContractID, meta.Cache)

	// without CacheVersion the cache ref must not resolve at the proxy's version
	noCacheVersion := *req
	noCacheVersion.CacheVersion = ""
	require.Error(t, SetProxyCache{}.VerifyPreconditions(env, &noCacheVersion))
}
