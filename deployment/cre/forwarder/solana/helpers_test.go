package solana

import (
	"math/big"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go/rpc"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/onchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

const (
	testQualifier = "test-deploy"
	testVersion   = "1.0.0"
)

var testSelector = chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector

// forwarderTestEnv is a Solana chain wired into a keystone environment, ready to run the forwarder
// changesets against.
type forwarderTestEnv struct {
	Runtime  *runtime.Runtime
	Chain    cldfsol.Chain
	DON      forwarder.DonConfiguration
	Selector uint64
}

// setupForwarderTestEnv starts a Solana container preloaded with the keystone programs, registers
// the forwarder program in the datastore and joins the chain to a keystone test environment that
// runs a workflow DON. When withMCMS is set, MCMS with timelock is deployed and its signer PDAs are
// funded, so changesets can be applied through proposals.
func setupForwarderTestEnv(t *testing.T, programsPath string, programIDs map[string]string, withMCMS bool) forwarderTestEnv {
	t.Helper()

	solChains, err := onchain.
		NewSolanaContainerLoader(programsPath, programIDs).
		Load(t, []uint64{testSelector})
	require.NoError(t, err)

	solChain, ok := cldfchain.NewBlockChainsFromSlice(solChains).SolanaChains()[testSelector]
	require.True(t, ok, "solana loader must return a solana chain for selector %d", testSelector)
	require.NotEmpty(t, solChain.ProgramsPath, "programs dir required for solana program deploy CLI")

	te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4, ChainSelectors: []uint64{testSelector}},
		AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
		WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
		NumChains:       1,
	})

	// Inject the solana chain into the environment (merge EVM first, then sol so sol is never overwritten).
	blockchains := make(map[uint64]cldfchain.BlockChain)
	for _, ch := range te.Env.BlockChains.All() {
		blockchains[ch.ChainSelector()] = ch
	}
	blockchains[testSelector] = solChain
	te.Env.BlockChains = cldfchain.NewBlockChains(blockchains)

	// Populate the datastore with the keystone forwarder contract
	ds := datastore.NewMemoryDataStore()
	if withMCMS {
		soltestutils.RegisterMCMSPrograms(t, testSelector, ds)
	}
	err = ds.AddressRefStore.Add(datastore.AddressRef{
		Address:       programIDs["keystone_forwarder"],
		ChainSelector: testSelector,
		Type:          ForwarderContract,
		Version:       semver.MustParse(testVersion),
		Qualifier:     testQualifier,
	})
	require.NoError(t, err)
	te.Env.DataStore = ds.Seal()

	// We set up a new runtime to execute the changesets based on the previously set up environment
	rt := runtime.NewFromEnvironment(te.Env)

	if withMCMS {
		mcmsState, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolanaV2(
			rt.Environment(),
			ds,
			solChain,
			cldfproposalutils.MCMSWithTimelockConfig{
				Canceller:        cldftesthelpers.SingleGroupMCMS(t),
				Proposer:         cldftesthelpers.SingleGroupMCMS(t),
				Bypasser:         cldftesthelpers.SingleGroupMCMS(t),
				TimelockMinDelay: big.NewInt(0),
			},
		)
		require.NoError(t, err)

		soltestutils.FundSignerPDAs(t, solChain, mcmsState)
	}

	var wfNodes []string
	for _, id := range te.GetP2PIDs("wfDon") {
		wfNodes = append(wfNodes, id.String())
	}

	return forwarderTestEnv{
		Runtime:  rt,
		Chain:    solChain,
		Selector: testSelector,
		DON: forwarder.DonConfiguration{
			Name:    "test-wf-don",
			ID:      1,
			F:       1,
			Version: 1,
			NodeIDs: wfNodes,
		},
	}
}

// requireOraclesConfig asserts whether the oracles config account of the given DON exists onchain.
func requireOraclesConfig(t *testing.T, env cldf.Environment, chainSel uint64, donID, configVersion uint32, wantExists bool) {
	t.Helper()

	chain, ok := env.BlockChains.SolanaChains()[chainSel]
	require.True(t, ok, "solana chain not found for chain selector %d", chainSel)

	target, err := resolveForwarderConfigTarget(env, chain, semver.MustParse(testVersion), testQualifier, donID, configVersion, nil)
	require.NoError(t, err)

	_, err = chain.Client.GetAccountInfoWithOpts(t.Context(), target.ConfigPDA, &rpc.GetAccountInfoOpts{
		Commitment: cldfsol.SolDefaultCommitment,
	})
	if wantExists {
		require.NoError(t, err, "expected oracles config %s of don %d to exist", target.ConfigPDA, donID)
		return
	}

	require.ErrorIs(t, err, rpc.ErrNotFound, "expected oracles config %s of don %d to be closed", target.ConfigPDA, donID)
}
