package changeset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

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

const (
	repoURL   = "https://github.com/smartcontractkit/chainlink-ccip.git"
	revision  = "main" // TODO get this from go.mod
	cloneDir  = "./.temp-repo"
	solanaDir = cloneDir + "/chains/solana"
)

func runCommand(command string, args []string, workDir string) (string, error) {
	fmt.Printf("Running command %s %v in %s\n", command, args, workDir)
	cmd := exec.Command(command, args...)
	cmd.Dir = workDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("stdout: %s\nstderr: %s", stdoutBuf.String(), stderrBuf.String()), err
	}
	return fmt.Sprintf("stdout: %s\nstderr: %s", stdoutBuf.String(), stderrBuf.String()), nil
}

func AssertAllContractsArePresent() {
	// list all contracts we'll expect to have
	expectedContracts := []string{
		"ccip_router.so",
	}

	// check if all contracts are present in the correct path
	programsPath := memory.GetProgramsPath()

	// rebuildContracts := false

	// check if all contracts are present
	for _, contract := range expectedContracts {
		contractPath := fmt.Sprintf("%s/%s", programsPath, contract)

		_, err := os.Stat(contractPath)

		if err != nil {
			panic(fmt.Sprintf("Contract %s not found in %s. Please run script TODO to populate them", contract, contractPath))
			// fmt.Sprintf("Contract %s not found in %s", contract, contractPath)
			// rebuildContracts = true
			// break
		}
	}

	// if !rebuildContracts {
	// 	// if all contracts are present, we can skip rebuilding them
	// 	return
	// }

	// fmt.Println("Cleaning up local repo...")
	// _, err := runCommand("rm", []string{"-rf", cloneDir}, ".")
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to clear folder: %v", err))
	// }

	// fmt.Println("Cloning repository...")
	// _, err = runCommand("git", []string{"clone", repoURL, cloneDir}, ".")
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to clone repository: %v", err))
	// }

	// fmt.Println("Checking out specific revision...")
	// _, err = runCommand("git", []string{"checkout", revision}, cloneDir)
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to checkout revision %s: %v", revision, err))
	// }

	// fmt.Println("Building contracts...")
	// _, err = runCommand("make", []string{"build-contracts"}, solanaDir)
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to build contracts: %s", err))
	// }

	// // move the contracts to the correct path
	// runCommand("pwd", []string{}, solanaDir)
	// runCommand("ls", []string{}, solanaDir)

	// // Check if the target directory exists
	// targetDir := solanaDir + "/contracts/target/deploy"

	// // List the .so files in the target directory
	// files, err := os.ReadDir(targetDir)
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to read target directory: %v", err))
	// }

	// // Copy each .so file individually, I don't know why the wildcard doesn't work
	// for _, file := range files {
	// 	if !file.IsDir() && strings.HasSuffix(file.Name(), ".so") {
	// 		sourcePath := fmt.Sprintf("%s/%s", targetDir, file.Name())
	// 		destPath := fmt.Sprintf("%s/%s", programsPath, file.Name())

	// 		_, err = runCommand("cp", []string{sourcePath, destPath}, ".")
	// 		if err != nil {
	// 			panic(fmt.Sprintf("Failed to copy contract %s: %v", sourcePath, err))
	// 		}
	// 	}
	// }
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
