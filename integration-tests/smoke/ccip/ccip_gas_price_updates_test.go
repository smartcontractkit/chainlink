package smoke

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Test_CCIPGasPriceUpdates tests that chain fee price updates are propagated correctly when
// price reaches some deviation threshold or when the price has expired.
func Test_CCIPGasPriceUpdates(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := changeset.Context(t)

	var gasPriceExpiry = 5 * time.Second
	e, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, &changeset.TestConfigs{
		OCRConfigOverride: func(params changeset.CCIPOCRParams) changeset.CCIPOCRParams {
			params.CommitOffChainConfig.RemoteGasPriceBatchWriteFrequency = *config.MustNewDuration(gasPriceExpiry)
			return params
		},
	})
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	require.NoError(t, changeset.AddLanesForAll(e.Env, state))

	allChainSelectors := maps.Keys(e.Env.Chains)
	assert.GreaterOrEqual(t, len(allChainSelectors), 2, "test requires at least 2 chains")

	sourceChain1 := allChainSelectors[0]
	sourceChain2 := allChainSelectors[1]

	feeQuoter1 := state.Chains[sourceChain1].FeeQuoter
	feeQuoter2 := state.Chains[sourceChain2].FeeQuoter

	// get initial chain fees
	chainFee2, err := feeQuoter1.GetDestinationChainGasPrice(&bind.CallOpts{Context: ctx}, sourceChain2)
	require.NoError(t, err)
	chainFee1, err := feeQuoter2.GetDestinationChainGasPrice(&bind.CallOpts{Context: ctx}, sourceChain1)
	require.NoError(t, err)
	t.Logf("initial chain1 fee (stored in chain2): %v", chainFee1)
	t.Logf("initial chain2 fee (stored in chain1): %v", chainFee2)

	// get latest price updates sequence number from the offRamps
	offRampChain1 := state.Chains[sourceChain1].OffRamp
	offRampChain2 := state.Chains[sourceChain2].OffRamp
	priceUpdatesSeqNumChain1, err := offRampChain1.GetLatestPriceSequenceNumber(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	priceUpdatesSeqNumChain2, err := offRampChain2.GetLatestPriceSequenceNumber(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	t.Logf("priceUpdatesSeqNumChain1: %v", priceUpdatesSeqNumChain1)
	t.Logf("priceUpdatesSeqNumChain2: %v", priceUpdatesSeqNumChain2)

	// update the price of chain2
	tx, err := feeQuoter1.UpdatePrices(e.Env.Chains[sourceChain1].DeployerKey, fee_quoter.InternalPriceUpdates{
		TokenPriceUpdates: nil,
		GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{
			{DestChainSelector: sourceChain2, UsdPerUnitGas: big.NewInt(5123)},
		},
	})
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[sourceChain1], tx, err)
	require.NoError(t, err)

	// assert that the chain fees are updated by the commit plugin reports
	priceDeviationChecked := false // flag to check if price deviation condition was met
	assert.Eventually(t, func() bool {
		// offRamps should have updated the sequence number
		priceUpdatesSeqNumChain1Now, err := offRampChain1.GetLatestPriceSequenceNumber(&bind.CallOpts{Context: ctx})
		require.NoError(t, err)
		priceUpdatesSeqNumChain2Now, err := offRampChain2.GetLatestPriceSequenceNumber(&bind.CallOpts{Context: ctx})
		require.NoError(t, err)
		t.Logf("priceUpdatesSeqNumChain1: %v", priceUpdatesSeqNumChain1Now)
		t.Logf("priceUpdatesSeqNumChain2: %v", priceUpdatesSeqNumChain2Now)
		if priceUpdatesSeqNumChain1Now <= priceUpdatesSeqNumChain1 {
			return false
		}
		if priceUpdatesSeqNumChain2Now <= priceUpdatesSeqNumChain2 {
			return false
		}

		chainFee2Now, err := feeQuoter1.GetDestinationChainGasPrice(&bind.CallOpts{Context: ctx}, sourceChain2)
		require.NoError(t, err)
		chainFee1Now, err := feeQuoter2.GetDestinationChainGasPrice(&bind.CallOpts{Context: ctx}, sourceChain1)
		require.NoError(t, err)
		t.Logf("chainFee1 (stored in chain2): %v", chainFee1Now)
		t.Logf("chainFee2 (stored in chain1): %v", chainFee2Now)

		if !priceDeviationChecked {
			// make sure there was a price update for chain2 when price deviation was reached
			if chainFee2Now.Value.Cmp(chainFee2.Value) == 0 {
				t.Logf("chainFee2 not updated: %v original=%v", chainFee2Now, chainFee2.Value)
				return false
			}
			require.NotEqual(t, chainFee2Now.Timestamp, chainFee2.Timestamp)
			priceDeviationChecked = true
		}

		// make sure there was a price update for chain1 but with the same price - when expiration is reached
		if chainFee1Now.Timestamp == chainFee1.Timestamp {
			t.Logf("chainFee1 timestamp not updated: %v original=%v", chainFee1Now, chainFee1.Timestamp)
			chainFee1 = chainFee1Now
			return false
		}
		if chainFee1Now.Value.Cmp(chainFee1.Value) != 0 {
			t.Logf("chainFee1 changed: %v prev:%v", chainFee1Now, chainFee1.Value)
			chainFee1 = chainFee1Now
			return false
		}

		return priceDeviationChecked
	}, tests.WaitTimeout(t), 500*time.Millisecond)
}
