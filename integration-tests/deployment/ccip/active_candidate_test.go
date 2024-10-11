package ccipdeployment

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"go.uber.org/zap/zapcore"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
)

func Test_ActiveCandidateMigration(t *testing.T) {
	// [SETUP]
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	ab := deployment.NewMemoryAddressBook()
	fromCS, toCS := allocateCCIPChainSelectors(e.Chains)
	feeTokenContracts, _ := DeployTestContracts(t, lggr, ab, fromCS, toCS, e.Chains)
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)

	state, err := LoadOnchainState(e, ab)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[fromCS].CapabilityRegistry)
	require.NotNil(t, state.Chains[fromCS].CCIPHome)
	require.NotNil(t, state.Chains[toCS].USDFeeds)

	err = DeployCCIPContracts(e, ab, DeployCCIPContractConfig{
		HomeChainSel:       fromCS,
		FeedChainSel:       toCS,
		ChainsToDeploy:     e.AllChainSelectors(),
		TokenConfig:        NewTokenConfig(),
		CapabilityRegistry: state.Chains[fromCS].CapabilityRegistry.Address(),
		MCMSConfig:         NewTestMCMSConfig(t, e),
		FeeTokenContracts:  feeTokenContracts,
	})
	require.NoError(t, err)
	state, err = LoadOnchainState(e, ab)
	require.NoError(t, err)

	// Add a uni directional lane
	err = AddLane(e, state, fromCS, toCS)
	require.NoError(t, err)
	// [SETUP] done

	// [ACTIVE ONLY, NO CANDIDATE] send successful request on active
	//h, err := e.Chains[toCS].Client.HeaderByNumber(testcontext.Get(t), nil)
	//require.NoError(t, err)
	//
	//startBlock := h.Number.Uint64()
	//seqNum := SendRequest(t, e, state, fromCS, toCS, false)
	//require.Equal(t, uint64(1), seqNum)
	//require.NoError(t, ConfirmExecWithSeqNr(t, e.Chains[fromCS], e.Chains[toCS], state.Chains[toCS].OffRamp, &startBlock, seqNum))
	// [ACTIVE ONLY, NO CANDIDATE] done

	// [ACTIVE, CANDIDATE] setup by setting candidate on CCIPHome
	state, err = LoadOnchainState(e, ab)
	require.NoError(t, err)

	donArgMap, err := BuildAddDONArgs(lggr, state.Chains[toCS].OffRamp, e.Chains[toCS], toCS, map[ocrtypes.Account]pluginconfig.TokenInfo{}, nodes.NonBootstraps())
	allDons, err := state.Chains[fromCS].CapabilityRegistry.GetDONs(&bind.CallOpts{})
	require.NoError(t, err)
	require.Equal(t, len(allDons), 2)
	commitDonId := allDons[0].Id
	execDonId := allDons[1].Id

	activeCommitDigest, err := state.Chains[fromCS].CCIPHome.GetActiveDigest(&bind.CallOpts{}, commitDonId, uint8(cctypes.PluginTypeCCIPCommit))
	require.NoError(t, err)

	tx, err := state.Chains[fromCS].CCIPHome.SetCandidate(
		e.Chains[fromCS].DeployerKey,
		commitDonId,
		uint8(cctypes.PluginTypeCCIPCommit),
		donArgMap[cctypes.PluginTypeCCIPCommit],
		activeCommitDigest)
	require.NoError(t, err)
	_, err = e.Chains[fromCS].Confirm(tx)
	require.NoError(t, err)

	activeExecDigest, err := state.Chains[fromCS].CCIPHome.GetActiveDigest(&bind.CallOpts{}, execDonId, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err)
	tx, err = state.Chains[fromCS].CCIPHome.SetCandidate(
		e.Chains[fromCS].DeployerKey,
		execDonId,
		uint8(cctypes.PluginTypeCCIPExec),
		donArgMap[cctypes.PluginTypeCCIPExec],
		activeExecDigest)
	require.NoError(t, err)
	_, err = e.Chains[fromCS].Confirm(tx)
	require.NoError(t, err)
	// [ACTIVE, CANDIDATE] done setup

	// [ACTIVE, CANDIDATE] send unsuccessful request on candidate

	// [ACTIVE, CANDIDATE] done send unsuccessful request on candidate

	// [NEW ACTIVE, NO CANDIDATE] promote to active
	// promote candidate and confirm by getting all configs
	candidateCommitDigest, err := state.Chains[fromCS].CCIPHome.GetCandidateDigest(&bind.CallOpts{}, commitDonId, uint8(cctypes.PluginTypeCCIPCommit))
	require.NoError(t, err)
	tx, err = state.Chains[fromCS].CCIPHome.PromoteCandidateAndRevokeActive(
		e.Chains[fromCS].DeployerKey,
		commitDonId,
		uint8(cctypes.PluginTypeCCIPCommit),
		candidateCommitDigest,
		activeCommitDigest)
	require.NoError(t, err)
	_, err = e.Chains[fromCS].Confirm(tx)
	require.NoError(t, err)

	allCommitConfigs, err := state.Chains[fromCS].CCIPHome.GetAllConfigs(&bind.CallOpts{}, commitDonId, uint8(cctypes.PluginTypeCCIPCommit))
	require.NoError(t, err)
	require.Nil(t, allCommitConfigs.CandidateConfig)
	require.NotNil(t, allCommitConfigs.ActiveConfig)

	// repeat above for exec don
	candidateExecDigest, err := state.Chains[fromCS].CCIPHome.GetCandidateDigest(&bind.CallOpts{}, commitDonId, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err)
	tx, err = state.Chains[fromCS].CCIPHome.PromoteCandidateAndRevokeActive(
		e.Chains[fromCS].DeployerKey,
		commitDonId,
		uint8(cctypes.PluginTypeCCIPCommit),
		candidateExecDigest,
		activeExecDigest)
	require.NoError(t, err)
	_, err = e.Chains[fromCS].Confirm(tx)
	require.NoError(t, err)

	allExecConfigs, err := state.Chains[fromCS].CCIPHome.GetAllConfigs(&bind.CallOpts{}, execDonId, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err)
	require.Nil(t, allExecConfigs.CandidateConfig)
	require.NotNil(t, allExecConfigs.ActiveConfig)
	// [NEW ACTIVE, NO CANDIDATE] done promoting

	// [NEW ACTIVE, NO CANDIDATE] send successful request on new active

	//h, err := e.Chains[toCS].Client.HeaderByNumber(testcontext.Get(t), nil)
	//require.NoError(t, err)
	//
	//startBlock := h.Number.Uint64()
	//seqNum := SendRequest(t, e, state, fromCS, toCS, false)
	//require.Equal(t, uint64(1), seqNum)
	//require.NoError(t, ConfirmExecWithSeqNr(t, e.Chains[fromCS], e.Chains[toCS], state.Chains[toCS].OffRamp, &startBlock, seqNum))

	// [NEW ACTIVE, NO CANDIDATE] done sending successful request
}
