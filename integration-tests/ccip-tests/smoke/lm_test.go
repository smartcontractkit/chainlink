package smoke

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
)

func TestLmBasic(t *testing.T) {
	t.Parallel()
	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	//gasLimit := big.NewInt(*TestCfg.TestGroupInput.MsgDetails.DestGasLimit)
	lmTestSetup := testsetups.LMDefaultTestSetup(t, &log, "smoke-lm", TestCfg)

	l1ChainId := lmTestSetup.Cfg.SelectedNetworks[0].ChainID
	l2ChainId := lmTestSetup.Cfg.SelectedNetworks[1].ChainID

	l1liquidityStart, err := lmTestSetup.LMModules[l1ChainId].LM.GetLiquidity()
	require.NoError(t, err, "Failed to get liquidity from L1")
	l2liquidityStart, err := lmTestSetup.LMModules[l2ChainId].LM.GetLiquidity()
	require.NoError(t, err, "Failed to get liquidity from L2")
	log.Info().Str("L1 Liquidity", l1liquidityStart.String()).Str("L2 Liquidity", l2liquidityStart.String()).Msg("Liquidity at start")

	// TODO: Improve this wait
	log.Info().Msg("Waiting 3 minutes for liquidity to change")
	time.Sleep(3 * time.Minute)

	l1liquidityEnd, err := lmTestSetup.LMModules[l1ChainId].LM.GetLiquidity()
	require.NoError(t, err, "Failed to get liquidity from L1")
	l2liquidityEnd, err := lmTestSetup.LMModules[l2ChainId].LM.GetLiquidity()
	require.NoError(t, err, "Failed to get liquidity from L2")
	log.Info().Str("L1 Liquidity", l1liquidityEnd.String()).Str("L2 Liquidity", l2liquidityEnd.String()).Msg("Liquidity at end")

	// Check if liquidity changed.
	// Ideally liquidity should change on both chains, but sometimes I have seen it not change on L2
	// TODO: Investigate why liquidity is not changing on L2
	if (l1liquidityEnd.Cmp(l1liquidityStart) == 0) && (l2liquidityEnd.Cmp(l2liquidityStart) == 0) {
		t.Errorf("Liquidity did not change during the test")
	}
	// TODO: The test should also check for LiquidityTransferred event on both chains
}
