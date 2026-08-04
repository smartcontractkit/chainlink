package stellar

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stellar/go-stellar-sdk/keypair"
	"github.com/stellar/go-stellar-sdk/xdr"
	"github.com/stretchr/testify/require"

	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	cldflogger "github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/smartcontractkit/chainlink-stellar/bindings/scval"
	stellardeploy "github.com/smartcontractkit/chainlink-stellar/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

const testContractID = "CA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJUWDA"

// testChainSel is a real Stellar chain selector (from chain-selectors) so
// that ChainMetadata methods (Family, NetworkType, ...) resolve correctly.
var testChainSel = chainsel.STELLAR_TESTNET.Selector

// newTestEnv returns an Environment with one fake Stellar chain and swapped-in fake deps.
// It swaps the package-level newStellarDeps seam to hand back fakeInvoker/fakeDeployer
// instead of talking to a real Soroban RPC endpoint, and restores the seam on cleanup.
func newTestEnv(t *testing.T) (cldf.Environment, *fakeInvoker, *fakeDeployer) {
	t.Helper()

	kp, err := keypair.Random()
	require.NoError(t, err)

	chain := cldfstellar.Chain{
		ChainMetadata: cldfstellar.ChainMetadata{Selector: testChainSel},
		Signer:        cldfstellar.NewStellarKeypairSigner(kp),
	}

	invoker := &fakeInvoker{}
	deployer := &fakeDeployer{contractID: testContractID}

	origDeps := newStellarDeps
	t.Cleanup(func() { newStellarDeps = origDeps })
	newStellarDeps = func(_ cldfstellar.Chain) (operation.StellarDeps, error) {
		return operation.StellarDeps{Deploy: deployer, Invoker: invoker}, nil
	}

	blockChains := cldfchain.NewBlockChains(map[uint64]cldfchain.BlockChain{
		testChainSel: chain,
	})

	env := cldf.NewEnvironment(
		"test",
		cldflogger.Test(t),
		nil, // ExistingAddresses: unused by these changesets, DataStore is authoritative
		datastore.NewMemoryDataStore().Seal(),
		nil,
		nil,
		func() context.Context { return context.Background() },
		ocr.OCRSecrets{},
		blockChains,
	)

	return *env, invoker, deployer
}

// contractRefSpec describes one AddressRef to seed via seedContractRefs.
type contractRefSpec struct {
	contractType datastore.ContractType
	address      string
	qualifier    string
	version      string
}

// seedContractRefs pre-seeds env.DataStore with the given AddressRefs (all in
// one MemoryDataStore) so changesets can resolve contracts by
// (ChainSel, Type, Version, Qualifier). Replaces any refs seeded by a prior
// call on the same env — call once per test with every ref the test needs
// (e.g. both a proxy and a cache ref for SetProxyCache).
func seedContractRefs(t *testing.T, env *cldf.Environment, refs ...contractRefSpec) {
	t.Helper()
	ds := datastore.NewMemoryDataStore()
	for _, r := range refs {
		require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
			Address:       r.address,
			ChainSelector: testChainSel,
			Type:          r.contractType,
			Version:       semver.MustParse(r.version),
			Qualifier:     r.qualifier,
		}))
	}
	env.DataStore = ds.Seal()
}

// seedCacheRef pre-seeds env.DataStore with a single cache AddressRef so
// config changesets (SetFeedConfigs, RemoveFeedConfigs, Add/RemoveFeedAdmin)
// can resolve the cache contract by (ChainSel, CacheContract, Version, Qualifier).
func seedCacheRef(t *testing.T, env *cldf.Environment, address, qualifier, version string) {
	t.Helper()
	seedContractRefs(t, env, contractRefSpec{CacheContract, address, qualifier, version})
}

// seedProxyRef pre-seeds env.DataStore with a single proxy AddressRef so
// changesets can resolve the proxy contract by
// (ChainSel, ProxyContract, Version, Qualifier).
func seedProxyRef(t *testing.T, env *cldf.Environment, address, qualifier, version string) {
	t.Helper()
	seedContractRefs(t, env, contractRefSpec{ProxyContract, address, qualifier, version})
}

// writeDummyWasm writes a placeholder wasm file under t.TempDir() so
// VerifyPreconditions' filesystem check (os.Stat) has a real path to find.
// The contents are never read — deployment is faked via fakeDeployer.
func writeDummyWasm(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte("dummy wasm"), 0o600))
	return path
}

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
// seeds the deploy salt directly. deploy_cache.go's Apply only falls back to
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

const testAdmin = "GAAZI4TCR3TY5OJHCTJC2A4QSY6CJWJH5IAJTGKIN2ER7LBNVKOCCWN7"

func TestSetFeedConfigsChangeset(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test", "1.0.0")

	req := &SetFeedConfigsRequest{
		ChainSel:     testChainSel,
		Qualifier:    "test",
		Version:      "1.0.0",
		Admin:        testAdmin,
		DataIDs:      []string{"0x018e16c39e00032000000"},
		Descriptions: []string{"BTC/USD"},
		Permissions: []FeedPermission{{
			AllowedSender:        testContractID, // the forwarder
			AllowedWorkflowOwner: "0x0102030405060708090a0b0c0d0e0f1011121314",
			AllowedWorkflowName:  "abc",
		}},
	}
	require.NoError(t, SetFeedConfigs{}.VerifyPreconditions(env, req))

	_, err := SetFeedConfigs{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "set_feed_configs", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)

	// Permissions apply to every feed in the batch: a single permission set
	// with two DataIDs must still pass length validation.
	multi := *req
	multi.DataIDs = []string{"0x01", "0x02"}
	multi.Descriptions = []string{"BTC/USD", "ETH/USD"}
	require.NoError(t, SetFeedConfigs{}.VerifyPreconditions(env, &multi))

	// mismatched lengths must fail preconditions
	bad := *req
	bad.Descriptions = nil
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &bad))

	// empty DataIDs must fail preconditions
	empty := *req
	empty.DataIDs = nil
	empty.Descriptions = nil
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &empty))

	// invalid admin address must fail preconditions
	badAdmin := *req
	badAdmin.Admin = "not-a-key"
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &badAdmin))

	// missing cache ref must fail preconditions
	badQualifier := *req
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &badQualifier))

	// invalid permission fields must fail preconditions
	badPerm := *req
	badPerm.Permissions = []FeedPermission{{
		AllowedSender:        "not-a-key",
		AllowedWorkflowOwner: "0x0102030405060708090a0b0c0d0e0f1011121314",
		AllowedWorkflowName:  "abc",
	}}
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &badPerm))

	// empty Permissions must fail preconditions
	emptyPerm := *req
	emptyPerm.Permissions = nil
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &emptyPerm))
}

func TestRemoveFeedConfigsChangeset(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test", "1.0.0")

	req := &RemoveFeedConfigsRequest{
		ChainSel:  testChainSel,
		Qualifier: "test",
		Version:   "1.0.0",
		Admin:     testAdmin,
		DataIDs:   []string{"0x018e16c39e00032000000"},
	}
	require.NoError(t, RemoveFeedConfigs{}.VerifyPreconditions(env, req))

	_, err := RemoveFeedConfigs{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "remove_feed_configs", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)

	// empty DataIDs must fail preconditions
	empty := *req
	empty.DataIDs = nil
	require.Error(t, RemoveFeedConfigs{}.VerifyPreconditions(env, &empty))

	// invalid admin address must fail preconditions
	badAdmin := *req
	badAdmin.Admin = "not-a-key"
	require.Error(t, RemoveFeedConfigs{}.VerifyPreconditions(env, &badAdmin))

	// missing cache ref must fail preconditions
	badQualifier := *req
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, RemoveFeedConfigs{}.VerifyPreconditions(env, &badQualifier))
}

func TestFeedAdminChangesets(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test", "1.0.0")

	req := &FeedAdminRequest{
		ChainSel:  testChainSel,
		Qualifier: "test",
		Version:   "1.0.0",
		Admin:     testAdmin,
	}

	require.NoError(t, AddFeedAdmin{}.VerifyPreconditions(env, req))
	_, err := AddFeedAdmin{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "add_feed_admin", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)

	require.NoError(t, RemoveFeedAdmin{}.VerifyPreconditions(env, req))
	_, err = RemoveFeedAdmin{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 2)
	require.Equal(t, "remove_feed_admin", inv.calls[1].Function)
	require.Equal(t, testContractID, inv.calls[1].ContractID)

	// invalid admin address must fail preconditions for both changesets
	badAdmin := *req
	badAdmin.Admin = "not-a-key"
	require.Error(t, AddFeedAdmin{}.VerifyPreconditions(env, &badAdmin))
	require.Error(t, RemoveFeedAdmin{}.VerifyPreconditions(env, &badAdmin))

	// missing cache ref must fail preconditions for both changesets
	badQualifier := *req
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, AddFeedAdmin{}.VerifyPreconditions(env, &badQualifier))
	require.Error(t, RemoveFeedAdmin{}.VerifyPreconditions(env, &badQualifier))
}

// testProxyAddress is a second, distinct contract address used wherever a
// test needs to tell the proxy and cache refs apart (e.g. asserting Apply
// invoked the proxy, not the cache).
const testProxyAddress = "CBPROXY00000000000000000000000000000000000000000000000AA"

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

func TestUpgradeCacheChangeset(t *testing.T) {
	env, inv, dep := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test", "1.0.0")

	var wantHash xdr.Hash
	for i := range wantHash {
		wantHash[i] = byte(i + 1)
	}
	dep.wasmHash = wantHash

	req := &UpgradeCacheRequest{
		ChainSel:  testChainSel,
		Qualifier: "test",
		Version:   "1.0.0",
		WasmPath:  writeDummyWasm(t, "new_cache.wasm"),
	}
	require.NoError(t, UpgradeCache{}.VerifyPreconditions(env, req))

	_, err := UpgradeCache{}.Apply(env, req)
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

	// missing wasm path must fail preconditions
	badPath := *req
	badPath.WasmPath = "/does/not/exist.wasm"
	require.Error(t, UpgradeCache{}.VerifyPreconditions(env, &badPath))

	// missing cache ref must fail preconditions
	badQualifier := *req
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, UpgradeCache{}.VerifyPreconditions(env, &badQualifier))
}

func TestRecoverTokensChangeset(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test", "1.0.0")

	req := &RecoverTokensRequest{
		ChainSel:  testChainSel,
		Qualifier: "test",
		Version:   "1.0.0",
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

	_, err := SetProxyCache{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "set_cache", inv.calls[0].Function)
	require.Equal(t, testProxyAddress, inv.calls[0].ContractID)
	require.Len(t, inv.calls[0].Args, 1)

	// missing cache ref must fail preconditions (proxy ref alone isn't enough)
	proxyOnly, _, _ := newTestEnv(t)
	seedProxyRef(t, &proxyOnly, testProxyAddress, "test-proxy", "1.0.0")
	require.Error(t, SetProxyCache{}.VerifyPreconditions(proxyOnly, req))

	// missing proxy ref must fail preconditions (cache ref alone isn't enough)
	cacheOnly, _, _ := newTestEnv(t)
	seedCacheRef(t, &cacheOnly, testContractID, "test-cache", "1.0.0")
	require.Error(t, SetProxyCache{}.VerifyPreconditions(cacheOnly, req))
}
