package mcmsnew_test

import (
	"encoding/json"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	evminternal "github.com/smartcontractkit/chainlink/deployment/common/changeset/evm/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployMCMSWithConfig(t *testing.T) {
	lggr := logger.TestLogger(t)

	chains := cldf_chain.NewBlockChainsFromSlice(
		memory.NewMemoryChainsEVMWithChainIDs(t, []uint64{chainsel.TEST_90000001.EvmChainID}, 1),
	).EVMChains()
	ab := cldf.NewMemoryAddressBook()

	// 1) Test WITHOUT a label
	mcmNoLabel, err := evminternal.DeployMCMSWithConfigEVM(
		types.ProposerManyChainMultisig,
		lggr,
		chains[chainsel.TEST_90000001.Selector],
		ab,
		proposalutils.SingleGroupMCMSV2(t),
	)
	require.NoError(t, err)
	require.Empty(t, mcmNoLabel.Tv.Labels, "expected no label to be set")

	// 2) Test WITH a label
	label := "SA"
	mcmWithLabel, err := evminternal.DeployMCMSWithConfigEVM(
		types.ProposerManyChainMultisig,
		lggr,
		chains[chainsel.TEST_90000001.Selector],
		ab,
		proposalutils.SingleGroupMCMSV2(t),
		evminternal.WithLabel(label),
	)
	require.NoError(t, err)
	require.NotNil(t, mcmWithLabel.Tv.Labels, "expected labels to be set")
	require.Contains(t, mcmWithLabel.Tv.Labels, label, "label mismatch")
}

func TestDeployMCMSWithTimelockContracts(t *testing.T) {
	lggr := logger.TestLogger(t)

	chains := cldf_chain.NewBlockChainsFromSlice(
		memory.NewMemoryChainsEVMWithChainIDs(t, []uint64{chainsel.TEST_90000001.EvmChainID}, 1),
	).EVMChains()

	ab := cldf.NewMemoryAddressBook()
	tenv := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	_, err := evminternal.DeployMCMSWithTimelockContractsEVM(tenv,
		chains[chainsel.TEST_90000001.Selector],
		ab, proposalutils.SingleGroupTimelockConfigV2(t), nil)
	require.NoError(t, err)
	addresses, err := ab.AddressesForChain(chainsel.TEST_90000001.Selector)
	require.NoError(t, err)
	require.Len(t, addresses, 5)
	mcmsState, err := state.MaybeLoadMCMSWithTimelockChainState(chains[chainsel.TEST_90000001.Selector], addresses)
	require.NoError(t, err)
	v, err := mcmsState.GenerateMCMSWithTimelockView()
	require.NoError(t, err)
	b, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	t.Log(string(b))
}
