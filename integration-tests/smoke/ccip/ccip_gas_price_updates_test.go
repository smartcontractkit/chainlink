package ccip

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// Test_CCIPGasPriceUpdates tests that chain fee price updates are propagated correctly when expiry time is reached.
func Test_CCIPGasPriceUpdatesWriteFrequency(t *testing.T) {
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

	// assert that the chain fees are updated by the commit plugin reports
	// Already because initial fee we set in the test is different from the fee that will be calculated
	// based on weth price, both chains will have their fees updated to new prices after gasPriceExpiry time.
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

		// First time the value changes anyways because the initial fee is different from the calculated fee
		// To determine if the fee was updated due to expiration and not deviation,
		// we need to check if the timestamp was updated after the first update and make sure value stays the same.
		if chain1FeeNow.Value.Cmp(initialChain1Fee.Value) != 0 {
			t.Logf("chainFee1 changed: %v prev:%v", chain1FeeNow, initialChain1Fee.Value)
			initialChain1Fee = chain1FeeNow
			return false
		}
		if chain2FeeNow.Value.Cmp(initialChain2Fee.Value) != 0 {
			t.Logf("chainFee2 changed: %v prev:%v", chain2FeeNow, initialChain2Fee.Value)
			initialChain2Fee = chain2FeeNow
			return false
		}

		chain1TimeChanged := chain1FeeNow.Value.Cmp(initialChain1Fee.Value) == 0 &&
			chain1FeeNow.Timestamp > initialChain1Fee.Timestamp
		chain2TimeChanged := chain2FeeNow.Value.Cmp(initialChain2Fee.Value) == 0 &&
			chain2FeeNow.Timestamp > initialChain2Fee.Timestamp
		if chain1TimeChanged {
			t.Logf("chainFee1 changed: %v prev:%v", chain1FeeNow, initialChain1Fee)
		}
		if chain2TimeChanged {
			t.Logf("chainFee2 changed: %v prev:%v", chain2FeeNow, initialChain2Fee)
		}

		if chain1TimeChanged && chain2TimeChanged {
			t.Logf("Calculated fees updated on both chains because expiration is reached")
			return true
		}
		return false
	}, tests.WaitTimeout(t), gasPriceExpiry)
}

// price reaches some deviation threshold
func Test_CCIPGasPriceUpdatesDeviation(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-576")

	ctx := testhelpers.Context(t)
	callOpts := &bind.CallOpts{Context: ctx}

	// Big expiry time to make sure that the update was due to price deviation
	var gasPriceExpiry = 30 * time.Minute
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

	// As weth price and ExecutionFee are constant during the test, when commit plugin calculates the fee it
	// will be different from the initial fees with a big enough deviation to trigger a price update.
	// We set up the test with big expiry time to make sure that the update was due to price deviation
	// and not expiration.
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

		chain1Changed := chain1FeeNow.Value.Cmp(initialChain1Fee.Value) != 0
		chain2Changed := chain2FeeNow.Value.Cmp(initialChain2Fee.Value) != 0
		if chain1Changed {
			t.Logf("chainFee1 changed: %v prev:%v", chain1FeeNow, initialChain1Fee.Value)
		}
		if chain2Changed {
			t.Logf("chainFee2 changed: %v prev:%v", chain2FeeNow, initialChain2Fee.Value)
		}

		if chain1Changed && chain2Changed {
			t.Logf("Calculated fees updated on both chains because deviation is reached")
			return true
		}
		return false
	}, tests.WaitTimeout(t), 5*time.Second)
}
