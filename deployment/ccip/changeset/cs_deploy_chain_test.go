package changeset

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func init() {
	AssertAllContractsArePresent()
}

func AssertAllContractsArePresent() {
	// list all contracts we'll expect to have
	expectedContracts := []string{
		"ccip_router.so",
	}

	// check if all contracts are present in the correct path
	programsPath := memory.GetProgramsPath()

	// check if all contracts are present
	for _, contract := range expectedContracts {
		contractPath := fmt.Sprintf("%s/%s", programsPath, contract)

		_, err := os.Stat(contractPath)

		if err != nil {
			panic(fmt.Sprintf("Contract %s not found in %s. Please run script TODO to populate them", contract, contractPath))
		}
	}

	// read current version of contracts and compare with expected
	currentContractsRevision, err := os.ReadFile(programsPath + "/.solana_contracts_rev")
	currentContractsRevision = []byte(strings.TrimSpace(string(currentContractsRevision)))
	if err != nil {
		panic(fmt.Sprintf("Could not read .solana_contracts_rev file: %v", err))
	}

	gomodContent, err := os.ReadFile("../../../deployment/go.mod")
	if err != nil {
		panic(fmt.Sprintf("Could not read deployment/go.mod file: %v", err))
	}

	modfile, err := modfile.Parse("deployment/go.mod", gomodContent, nil)
	if err != nil {
		panic(fmt.Sprintf("Could not parse deployment/go.mod file: %v", err))
	}

	var solanaModVersion string = ""
	for _, r := range modfile.Require {
		if r.Mod.Path == "github.com/smartcontractkit/chainlink-ccip/chains/solana" {
			solanaModVersion = r.Mod.Version
			break
		}
	}
	if solanaModVersion == "" {
		panic("Could not find ccip solana module in deployment/go.mod")
	}

	solanaModVersionRev, err := module.PseudoVersionRev(solanaModVersion)

	if err != nil {
		panic(fmt.Sprintf("Could not parse ccip solana module version: %v", err))
	}

	if string(currentContractsRevision) != solanaModVersionRev {
		panic(fmt.Sprintf("solana_contracts_rev file is outdated. Please run script TODO to update it"))
	}
}

func TestDeployChainContractsChangeset(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.AllChainSelectorsSolana()
	selectors := make([]uint64, 0, len(evmSelectors)+len(solChainSelectors))
	selectors = append(selectors, evmSelectors...)
	selectors = append(selectors, solChainSelectors...)
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIds := nodes.NonBootstraps().PeerIDs()
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfig(t)
	}
	var prereqCfg []DeployPrerequisiteConfigPerChain
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	SavePreloadedSolAddresses(t, e, solChainSelectors[0])
	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(DeployHomeChain),
			Config: DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  NewTestRMNStaticConfig(),
				RMNDynamicConfig: NewTestRMNDynamicConfig(),
				NodeOperators:    NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    selectors,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    cfg,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployPrerequisites),
			Config: DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployChainContracts),
			Config: DeployChainContractsConfig{
				ChainSelectors:    selectors,
				HomeChainSelector: homeChainSel,
			},
		},
	})
	require.NoError(t, err)

	// load onchain state
	state, err := LoadOnchainState(e)
	require.NoError(t, err)

	// verify all contracts populated
	require.NotNil(t, state.Chains[homeChainSel].CapabilityRegistry)
	require.NotNil(t, state.Chains[homeChainSel].CCIPHome)
	require.NotNil(t, state.Chains[homeChainSel].RMNHome)
	for _, sel := range evmSelectors {
		require.NotNil(t, state.Chains[sel].LinkToken)
		require.NotNil(t, state.Chains[sel].Weth9)
		require.NotNil(t, state.Chains[sel].TokenAdminRegistry)
		require.NotNil(t, state.Chains[sel].RegistryModule)
		require.NotNil(t, state.Chains[sel].Router)
		require.NotNil(t, state.Chains[sel].RMNRemote)
		require.NotNil(t, state.Chains[sel].TestRouter)
		require.NotNil(t, state.Chains[sel].NonceManager)
		require.NotNil(t, state.Chains[sel].FeeQuoter)
		require.NotNil(t, state.Chains[sel].OffRamp)
		require.NotNil(t, state.Chains[sel].OnRamp)
	}

	solState, err := LoadOnchainStateSolana(e)
	require.NoError(t, err)
	for _, sel := range solChainSelectors {
		require.NotNil(t, solState.SolChains[sel].LinkToken)
		require.NotNil(t, solState.SolChains[sel].SolCcipRouter)
	}

}

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()
	e := NewMemoryEnvironment(t)
	// Deploy all the CCIP contracts.
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)
	snap, err := state.View(e.Env.AllChainSelectors())
	require.NoError(t, err)

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}
