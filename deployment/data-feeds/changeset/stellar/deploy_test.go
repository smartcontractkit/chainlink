package stellar

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	stellardeploy "github.com/smartcontractkit/chainlink-stellar/deployment"
)

func TestDeployCacheChangeset(t *testing.T) {
	env, _, dep := newTestEnv(t)

	req := &DeployCacheRequest{
		ChainSel:  testChainSel,
		WasmPath:  writeDummyWasm(t, "data_feeds_cache.wasm"),
		Qualifier: "test-cache",
		Version:   "1.0.0",
	}

	require.NoError(t, DeployCache{}.VerifyPreconditions(env, req))

	out, err := DeployCache{}.Apply(env, req)
	require.NoError(t, err)
	require.NotNil(t, out.DataStore)

	require.Len(t, dep.deploys, 1)
	require.Len(t, dep.deploys[0].Args, 1) // owner ctor arg

	owner := env.BlockChains.StellarChains()[testChainSel].Signer.Address()
	wantSalt := stellardeploy.GenerateDeterministicSalt(owner, "data_feeds_cache-"+req.Qualifier)
	require.Equal(t, wantSalt, dep.deploys[0].Salt)

	key := datastore.NewAddressRefKey(testChainSel, CacheContract, semver.MustParse("1.0.0"), "test-cache")
	ref, err := out.DataStore.Addresses().Get(key)
	require.NoError(t, err)
	require.Equal(t, testContractID, ref.Address)

	// invalid explicit Owner must fail preconditions
	badOwner := *req
	badOwner.Owner = "not-a-key"
	require.Error(t, DeployCache{}.VerifyPreconditions(env, &badOwner))

	// unknown chain must fail preconditions
	badChain := *req
	badChain.ChainSel = 999999
	require.Error(t, DeployCache{}.VerifyPreconditions(env, &badChain))

	// invalid version must fail preconditions
	badVersion := *req
	badVersion.Version = "not-semver"
	require.Error(t, DeployCache{}.VerifyPreconditions(env, &badVersion))

	// missing wasm file must fail preconditions
	badWasm := *req
	badWasm.WasmPath = "/does/not/exist.wasm"
	require.Error(t, DeployCache{}.VerifyPreconditions(env, &badWasm))
}

// An explicit Owner seeds the deploy salt; the chain signer is only the
// fallback when Owner is empty.
func TestDeployCacheChangeset_ExplicitOwner(t *testing.T) {
	env, _, dep := newTestEnv(t)

	req := &DeployCacheRequest{
		ChainSel:  testChainSel,
		WasmPath:  writeDummyWasm(t, "data_feeds_cache.wasm"),
		Owner:     testAdmin,
		Qualifier: "test-cache-owner",
		Version:   "1.0.0",
	}

	require.NoError(t, DeployCache{}.VerifyPreconditions(env, req))

	out, err := DeployCache{}.Apply(env, req)
	require.NoError(t, err)
	require.NotNil(t, out.DataStore)

	require.Len(t, dep.deploys, 1)
	require.Len(t, dep.deploys[0].Args, 1) // owner ctor arg

	wantSalt := stellardeploy.GenerateDeterministicSalt(testAdmin, "data_feeds_cache-"+req.Qualifier)
	require.Equal(t, wantSalt, dep.deploys[0].Salt)
}

func TestDeployProxyChangeset(t *testing.T) {
	env, _, dep := newTestEnv(t)

	cacheAddr := "CCACHEFAKE0000000000000000000000000000000000000000000000"
	seedCacheRef(t, &env, cacheAddr, "test-cache", "1.0.0")

	req := &DeployProxyRequest{
		ChainSel:       testChainSel,
		WasmPath:       writeDummyWasm(t, "data_feeds_proxy.wasm"),
		CacheQualifier: "test-cache",
		Qualifier:      "test-proxy",
		Version:        "1.0.0",
	}

	require.NoError(t, DeployProxy{}.VerifyPreconditions(env, req))

	out, err := DeployProxy{}.Apply(env, req)
	require.NoError(t, err)
	require.NotNil(t, out.DataStore)

	require.Len(t, dep.deploys, 1)
	require.Len(t, dep.deploys[0].Args, 2) // owner + cache ctor args

	owner := env.BlockChains.StellarChains()[testChainSel].Signer.Address()
	wantSalt := stellardeploy.GenerateDeterministicSalt(owner, "data_feeds_proxy-"+req.Qualifier)
	require.Equal(t, wantSalt, dep.deploys[0].Salt)

	key := datastore.NewAddressRefKey(testChainSel, ProxyContract, semver.MustParse("1.0.0"), "test-proxy")
	ref, err := out.DataStore.Addresses().Get(key)
	require.NoError(t, err)
	require.Equal(t, testContractID, ref.Address)

	// the metadata mirror records which cache the proxy points at
	meta := outputMetadata(t, out, testContractID)
	require.Equal(t, cacheAddr, meta.Cache)

	// missing cache ref must fail preconditions
	missing := *req
	missing.CacheQualifier = "does-not-exist"
	require.Error(t, DeployProxy{}.VerifyPreconditions(env, &missing))
}
