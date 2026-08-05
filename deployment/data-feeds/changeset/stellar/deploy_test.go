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
}

func TestDeployCacheChangeset_MissingChain(t *testing.T) {
	env, _, _ := newTestEnv(t)

	req := &DeployCacheRequest{
		ChainSel:  999999,
		WasmPath:  writeDummyWasm(t, "data_feeds_cache.wasm"),
		Qualifier: "test-cache",
		Version:   "1.0.0",
	}

	require.Error(t, DeployCache{}.VerifyPreconditions(env, req))
}

// TestDeployCacheChangeset_ExplicitOwner pins that an explicit req.Owner
// seeds the deploy salt directly. deploy.go's DeployCache.Apply only falls back to
// ch.Signer.Address() when req.Owner is empty; when it's set, the salt is
// GenerateDeterministicSalt(req.Owner, ...) -- not the chain signer's address.
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

	// Seed the datastore with a cache AddressRef so DeployProxy can resolve it.
	seedDS := datastore.NewMemoryDataStore()
	require.NoError(t, seedDS.Addresses().Add(datastore.AddressRef{
		Address:       "CCACHEFAKE0000000000000000000000000000000000000000000000",
		ChainSelector: testChainSel,
		Type:          CacheContract,
		Version:       semver.MustParse("1.0.0"),
		Qualifier:     "test-cache",
	}))
	env.DataStore = seedDS.Seal()

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
}

func TestDeployProxyChangeset_MissingCache(t *testing.T) {
	env, _, _ := newTestEnv(t)

	req := &DeployProxyRequest{
		ChainSel:       testChainSel,
		WasmPath:       writeDummyWasm(t, "data_feeds_proxy.wasm"),
		CacheQualifier: "does-not-exist",
		Qualifier:      "test-proxy",
		Version:        "1.0.0",
	}

	require.Error(t, DeployProxy{}.VerifyPreconditions(env, req))
}
