package stellar

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stellar/go-stellar-sdk/keypair"
	protocolrpc "github.com/stellar/go-stellar-sdk/protocols/rpc"
	"github.com/stellar/go-stellar-sdk/xdr"
	"github.com/stretchr/testify/require"

	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	cldflogger "github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

const (
	testContractID = "CA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJUWDA"
	testAdmin      = "GAAZI4TCR3TY5OJHCTJC2A4QSY6CJWJH5IAJTGKIN2ER7LBNVKOCCWN7"
	// testProxyAddress is a second, distinct contract address used wherever a
	// test needs to tell the proxy and cache refs apart (e.g. asserting Apply
	// invoked the proxy, not the cache).
	testProxyAddress = "CBPROXY00000000000000000000000000000000000000000000000AA"
)

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

type invocation struct {
	ContractID string
	Function   string
	Args       []xdr.ScVal
}

type fakeInvoker struct {
	calls []invocation
}

func (f *fakeInvoker) InvokeContract(_ context.Context, contractID, fn string, args []xdr.ScVal) (*xdr.ScVal, error) {
	f.calls = append(f.calls, invocation{contractID, fn, args})
	v := xdr.ScVal{Type: xdr.ScValTypeScvVoid}
	return &v, nil
}

func (f *fakeInvoker) SimulateContract(_ context.Context, contractID, fn string, args []xdr.ScVal) (*xdr.ScVal, error) {
	f.calls = append(f.calls, invocation{contractID, fn, args})
	v := xdr.ScVal{Type: xdr.ScValTypeScvVoid}
	return &v, nil
}

func (f *fakeInvoker) GetEvents(_ context.Context, _ string, _ uint32, _ []string) ([]protocolrpc.EventInfo, error) {
	return nil, nil
}

type deployCall struct {
	WasmPath string
	Salt     [32]byte
	Args     []xdr.ScVal
}

type fakeDeployer struct {
	deploys    []deployCall
	contractID string
	wasmHash   xdr.Hash
	uploads    []string // wasm paths passed to UploadContractWASM, in call order
}

func (f *fakeDeployer) DeployContractWithArgs(_ context.Context, wasmPath string, salt [32]byte, args []xdr.ScVal) (string, error) {
	f.deploys = append(f.deploys, deployCall{wasmPath, salt, args})
	return f.contractID, nil
}

func (f *fakeDeployer) UploadContractWASM(_ context.Context, wasmPath string) (xdr.Hash, error) {
	f.uploads = append(f.uploads, wasmPath)
	return f.wasmHash, nil
}
