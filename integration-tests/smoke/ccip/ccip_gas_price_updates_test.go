package ccip

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
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"
)

// Test_CCIPGasPriceUpdates tests that chain fee price updates are propagated correctly when
// price reaches some deviation threshold or when the price has expired.
func Test_CCIPGasPriceUpdates(t *testing.T) {
	ctx := testhelpers.Context(t)
	callOpts := &bind.CallOpts{Context: ctx}

	var gasPriceExpiry = 5 * time.Second
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithOCRConfigOverride(func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			if params.CommitOffChainConfig != nil {
				params.CommitOffChainConfig.RemoteGasPriceBatchWriteFrequency = *config.MustNewDuration(gasPriceExpiry)
			}
			return params
		}),
	)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	testhelpers.AddLanesForAll(t, &e, state)

	allChainSelectors := maps.Keys(e.Env.Chains)
	assert.GreaterOrEqual(t, len(allChainSelectors), 2, "test requires at least 2 chains")

	sourceChain1 := allChainSelectors[0]
	sourceChain2 := allChainSelectors[1]

	feeQuoter1 := state.Chains[sourceChain1].FeeQuoter
	feeQuoter2 := state.Chains[sourceChain2].FeeQuoter

	// get initial chain fees
	initialChain2Fee, err := feeQuoter1.GetDestinationChainGasPrice(callOpts, sourceChain2)
	require.NoError(t, err)
	initialChain1Fee, err := feeQuoter2.GetDestinationChainGasPrice(callOpts, sourceChain1)
	require.NoError(t, err)
	t.Logf("initial chain1 fee (stored in chain2): %v", initialChain1Fee)
	t.Logf("initial chain2 fee (stored in chain1): %v", initialChain2Fee)

	// get latest price updates sequence number from the offRamps
	offRampChain1 := state.Chains[sourceChain1].OffRamp
	offRampChain2 := state.Chains[sourceChain2].OffRamp
	priceUpdatesSeqNumChain1, err := offRampChain1.GetLatestPriceSequenceNumber(callOpts)
	require.NoError(t, err)
	priceUpdatesSeqNumChain2, err := offRampChain2.GetLatestPriceSequenceNumber(callOpts)
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
		priceUpdatesSeqNumChain1Now, err := offRampChain1.GetLatestPriceSequenceNumber(callOpts)
		require.NoError(t, err)
		priceUpdatesSeqNumChain2Now, err := offRampChain2.GetLatestPriceSequenceNumber(callOpts)
		require.NoError(t, err)
		t.Logf("priceUpdatesSeqNumChain1: %v", priceUpdatesSeqNumChain1Now)
		t.Logf("priceUpdatesSeqNumChain2: %v", priceUpdatesSeqNumChain2Now)
		if priceUpdatesSeqNumChain1Now <= priceUpdatesSeqNumChain1 {
			return false
		}
		if priceUpdatesSeqNumChain2Now <= priceUpdatesSeqNumChain2 {
			return false
		}

		chain2FeeNow, err := feeQuoter1.GetDestinationChainGasPrice(callOpts, sourceChain2)
		require.NoError(t, err)
		chain1FeeNow, err := feeQuoter2.GetDestinationChainGasPrice(callOpts, sourceChain1)
		require.NoError(t, err)
		t.Logf("chainFee1 (stored in chain2): %v", chain1FeeNow)
		t.Logf("chainFee2 (stored in chain1): %v", chain2FeeNow)

		if !priceDeviationChecked {
			// make sure there was a price update for chain2 when price deviation was reached
			if chain2FeeNow.Value.Cmp(initialChain2Fee.Value) == 0 {
				t.Logf("chainFee2 not updated: %v original=%v", chain2FeeNow, initialChain2Fee.Value)
				return false
			}
			require.NotEqual(t, chain2FeeNow.Timestamp, initialChain2Fee.Timestamp)
			priceDeviationChecked = true
		}

		// make sure there was a price update for chain1 but with the same price - when expiration is reached
		if chain1FeeNow.Timestamp == initialChain1Fee.Timestamp {
			t.Logf("chainFee1 timestamp not updated: %v original=%v", chain1FeeNow, initialChain1Fee.Timestamp)
			initialChain1Fee = chain1FeeNow
			return false
		}
		if chain1FeeNow.Value.Cmp(initialChain1Fee.Value) != 0 {
			t.Logf("chainFee1 changed: %v prev:%v", chain1FeeNow, initialChain1Fee.Value)
			initialChain1Fee = chain1FeeNow
			return false
		}

		return priceDeviationChecked
	}, tests.WaitTimeout(t), 500*time.Millisecond)
}
