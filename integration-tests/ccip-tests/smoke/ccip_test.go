package smoke

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/lock_release_token_pool"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
)

type testDefinition struct {
	testName string
	lane     *actions.CCIPLane
}

func TestSmokeCCIPForBidirectionalLane(t *testing.T) {
	t.Parallel()
	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	require.NotNil(t, TestCfg.TestGroupInput.MsgDetails.DestGasLimit)
	gasLimit := big.NewInt(*TestCfg.TestGroupInput.MsgDetails.DestGasLimit)
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, &log, "smoke-ccip", nil, TestCfg)
	if len(setUpOutput.Lanes) == 0 {
		log.Info().Msg("No lanes found")
		return
	}

	t.Cleanup(func() {
		// If we are running a test that is a token transfer, we need to verify the balance.
		// skip the balance check for existing deployment, there can be multiple external requests in progress for existing deployments
		// other than token transfer initiated by the test, which can affect the balance check
		// therefore we check the balance only for the ccip environment created by the test
		if TestCfg.TestGroupInput.MsgDetails.IsTokenTransfer() &&
			!pointer.GetBool(TestCfg.TestGroupInput.USDCMockDeployment) &&
			!pointer.GetBool(TestCfg.TestGroupInput.ExistingDeployment) {
			setUpOutput.Balance.Verify(t)
		}
		require.NoError(t, setUpOutput.TearDown())
	})

	// Create test definitions for each lane.
	var tests []testDefinition
	for _, lane := range setUpOutput.Lanes {
		tests = append(tests, testDefinition{
			testName: fmt.Sprintf("CCIP message transfer from network %s to network %s",
				lane.ForwardLane.SourceNetworkName, lane.ForwardLane.DestNetworkName),
			lane: lane.ForwardLane,
		})
		if lane.ReverseLane != nil {
			tests = append(tests, testDefinition{
				testName: fmt.Sprintf("CCIP message transfer from network %s to network %s",
					lane.ReverseLane.SourceNetworkName, lane.ReverseLane.DestNetworkName),
				lane: lane.ReverseLane,
			})
		}
	}

	// Execute tests.
	log.Info().Int("Total Lanes", len(tests)).Msg("Starting CCIP test")
	for _, test := range tests {
		tc := test
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			tc.lane.Test = t
			log.Info().
				Str("Source", tc.lane.SourceNetworkName).
				Str("Destination", tc.lane.DestNetworkName).
				Msgf("Starting lane %s -> %s", tc.lane.SourceNetworkName, tc.lane.DestNetworkName)

			tc.lane.RecordStateBeforeTransfer()
			err := tc.lane.SendRequests(1, gasLimit)
			require.NoError(t, err)
			tc.lane.ValidateRequests()
		})
	}
}

func TestSmokeCCIPRateLimit(t *testing.T) {
	t.Parallel()
	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	require.True(t, TestCfg.TestGroupInput.MsgDetails.IsTokenTransfer(), "Test config should have token transfer message type")
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, &log, "smoke-ccip", nil, TestCfg)
	if len(setUpOutput.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		require.NoError(t, setUpOutput.TearDown())
	})

	var tests []testDefinition
	for _, lane := range setUpOutput.Lanes {
		tests = append(tests, testDefinition{
			testName: fmt.Sprintf("Network %s to network %s",
				lane.ForwardLane.SourceNetworkName, lane.ForwardLane.DestNetworkName),
			lane: lane.ForwardLane,
		})
	}

	// if we are running in simulated or in testnet mode, we can set the rate limit to test friendly values
	// For mainnet, we need to set this as false to avoid changing the deployed contract config
	setRateLimit := true
	AggregatedRateLimitCapacity := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(30))
	AggregatedRateLimitRate := big.NewInt(1e17)

	TokenPoolRateLimitCapacity := new(big.Int).Mul(big.NewInt(1e17), big.NewInt(1))
	TokenPoolRateLimitRate := big.NewInt(1e14)

	for _, test := range tests {
		tc := test
		t.Run(fmt.Sprintf("%s - Rate Limit", tc.testName), func(t *testing.T) {
			tc.lane.Test = t
			src := tc.lane.Source
			// add liquidity to pools on both networks
			if !pointer.GetBool(TestCfg.TestGroupInput.ExistingDeployment) {
				addLiquidity(t, src.Common, new(big.Int).Mul(AggregatedRateLimitCapacity, big.NewInt(20)))
				addLiquidity(t, tc.lane.Dest.Common, new(big.Int).Mul(AggregatedRateLimitCapacity, big.NewInt(20)))
			}
			log.Info().
				Str("Source", tc.lane.SourceNetworkName).
				Str("Destination", tc.lane.DestNetworkName).
				Msgf("Starting lane %s -> %s", tc.lane.SourceNetworkName, tc.lane.DestNetworkName)

			// capture the rate limit config before we change it
			prevRLOnRamp, err := src.OnRamp.Instance.CurrentRateLimiterState(nil)
			require.NoError(t, err)
			tc.lane.Logger.Info().Interface("rate limit", prevRLOnRamp).Msg("Initial OnRamp rate limiter state")

			prevOnRampRLTokenPool, err := src.Common.BridgeTokenPools[0].Instance.GetCurrentOutboundRateLimiterState(
				nil, tc.lane.Source.DestChainSelector,
			) // TODO RENS maybe?
			require.NoError(t, err)
			tc.lane.Logger.Info().
				Interface("rate limit", prevOnRampRLTokenPool).
				Str("pool", src.Common.BridgeTokenPools[0].Address()).
				Str("onRamp", src.OnRamp.Address()).
				Msg("Initial Token Pool rate limiter state")

			// some sanity checks
			rlOffRamp, err := tc.lane.Dest.OffRamp.Instance.CurrentRateLimiterState(nil)
			require.NoError(t, err)
			tc.lane.Logger.Info().Interface("rate limit", rlOffRamp).Msg("Initial OffRamp rate limiter state")
			if rlOffRamp.IsEnabled {
				require.GreaterOrEqual(t, rlOffRamp.Capacity.Cmp(prevRLOnRamp.Capacity), 0,
					"OffRamp Aggregated capacity should be greater than or equal to OnRamp Aggregated capacity",
				)
			}

			prevOffRampRLTokenPool, err := tc.lane.Dest.Common.BridgeTokenPools[0].Instance.GetCurrentInboundRateLimiterState(
				nil, tc.lane.Dest.SourceChainSelector,
			) // TODO RENS maybe?
			require.NoError(t, err)
			tc.lane.Logger.Info().
				Interface("rate limit", prevOffRampRLTokenPool).
				Str("pool", tc.lane.Dest.Common.BridgeTokenPools[0].Address()).
				Str("offRamp", tc.lane.Dest.OffRamp.Address()).
				Msg("Initial Token Pool rate limiter state")
			if prevOffRampRLTokenPool.IsEnabled {
				require.GreaterOrEqual(t, prevOffRampRLTokenPool.Capacity.Cmp(prevOnRampRLTokenPool.Capacity), 0,
					"OffRamp Token Pool capacity should be greater than or equal to OnRamp Token Pool capacity",
				)
			}

			AggregatedRateLimitChanged := false
			TokenPoolRateLimitChanged := false

			// reset the rate limit config to what it was before the tc
			t.Cleanup(func() {
				if AggregatedRateLimitChanged {
					require.NoError(t, src.OnRamp.SetRateLimit(evm_2_evm_onramp.RateLimiterConfig{
						IsEnabled: prevRLOnRamp.IsEnabled,
						Capacity:  prevRLOnRamp.Capacity,
						Rate:      prevRLOnRamp.Rate,
					}), "setting rate limit")
					require.NoError(t, src.Common.ChainClient.WaitForEvents(), "waiting for events")
				}
				if TokenPoolRateLimitChanged {
					require.NoError(t, src.Common.BridgeTokenPools[0].SetRemoteChainRateLimits(src.DestChainSelector,
						token_pool.RateLimiterConfig{
							Capacity:  prevOnRampRLTokenPool.Capacity,
							IsEnabled: prevOnRampRLTokenPool.IsEnabled,
							Rate:      prevOnRampRLTokenPool.Rate,
						}))
					require.NoError(t, src.Common.ChainClient.WaitForEvents(), "waiting for events")
				}
			})

			if setRateLimit {
				if prevRLOnRamp.Capacity.Cmp(AggregatedRateLimitCapacity) != 0 ||
					prevRLOnRamp.Rate.Cmp(AggregatedRateLimitRate) != 0 ||
					!prevRLOnRamp.IsEnabled {
					require.NoError(t, src.OnRamp.SetRateLimit(evm_2_evm_onramp.RateLimiterConfig{
						IsEnabled: true,
						Capacity:  AggregatedRateLimitCapacity,
						Rate:      AggregatedRateLimitRate,
					}), "setting rate limit on onramp")
					require.NoError(t, src.Common.ChainClient.WaitForEvents(), "waiting for events")
					AggregatedRateLimitChanged = true
				}
			} else {
				AggregatedRateLimitCapacity = prevRLOnRamp.Capacity
				AggregatedRateLimitRate = prevRLOnRamp.Rate
			}

			rlOnRamp, err := src.OnRamp.Instance.CurrentRateLimiterState(nil)
			require.NoError(t, err)
			tc.lane.Logger.Info().Interface("rate limit", rlOnRamp).Msg("OnRamp rate limiter state")
			require.True(t, rlOnRamp.IsEnabled, "OnRamp rate limiter should be enabled")

			tokenPrice, err := src.Common.PriceRegistry.Instance.GetTokenPrice(nil, src.Common.BridgeTokens[0].ContractAddress)
			require.NoError(t, err)
			tc.lane.Logger.Info().Str("tokenPrice.Value", tokenPrice.String()).Msg("Price Registry Token Price")

			totalTokensForOnRampCapacity := new(big.Int).Mul(
				big.NewInt(1e18),
				new(big.Int).Div(rlOnRamp.Capacity, tokenPrice),
			)

			tc.lane.Source.Common.ChainClient.ParallelTransactions(true)

			// current tokens are equal to the full capacity  - should fail
			src.TransferAmount[0] = rlOnRamp.Tokens
			tc.lane.Logger.Info().Str("tokensToSend", rlOnRamp.Tokens.String()).Msg("Aggregated Capacity")
			// approve the tokens
			require.NoError(t, src.Common.BridgeTokens[0].Approve(
				tc.lane.Source.Common.ChainClient.GetDefaultWallet(), src.Common.Router.Address(), src.TransferAmount[0]),
			)
			require.NoError(t, tc.lane.Source.Common.ChainClient.WaitForEvents())
			failedTx, _, _, _, err := tc.lane.Source.SendRequest(
				tc.lane.Dest.ReceiverDapp.EthAddress,
				big.NewInt(actions.DefaultDestinationGasLimit),
			)
			require.NoError(t, err)
			require.Error(t, tc.lane.Source.Common.ChainClient.WaitForEvents())
			errReason, v, err := tc.lane.Source.Common.ChainClient.RevertReasonFromTx(failedTx, evm_2_evm_onramp.EVM2EVMOnRampABI)
			require.NoError(t, err)
			tc.lane.Logger.Info().
				Str("Revert Reason", errReason).
				Interface("Args", v).
				Str("TokensSent", src.TransferAmount[0].String()).
				Str("Token", tc.lane.Source.Common.BridgeTokens[0].Address()).
				Str("FailedTx", failedTx.Hex()).
				Msg("Msg sent with tokens more than AggregateValueMaxCapacity")
			require.Equal(t, "AggregateValueMaxCapacityExceeded", errReason)

			// 99% of the aggregated capacity - should succeed
			tokensToSend := new(big.Int).Div(new(big.Int).Mul(totalTokensForOnRampCapacity, big.NewInt(99)), big.NewInt(100))
			tc.lane.Logger.Info().Str("tokensToSend", tokensToSend.String()).Msg("99% of Aggregated Capacity")
			tc.lane.RecordStateBeforeTransfer()
			src.TransferAmount[0] = tokensToSend
			err = tc.lane.SendRequests(1, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)

			// try to send again with amount more than the amount refilled by rate and
			// this should fail, as the refill rate is not enough to refill the capacity
			src.TransferAmount[0] = new(big.Int).Mul(AggregatedRateLimitRate, big.NewInt(10))
			failedTx, _, _, _, err = tc.lane.Source.SendRequest(tc.lane.Dest.ReceiverDapp.EthAddress, big.NewInt(actions.DefaultDestinationGasLimit))
			tc.lane.Logger.Info().Str("tokensToSend", src.TransferAmount[0].String()).Msg("More than Aggregated Rate")
			require.NoError(t, err)
			require.Error(t, tc.lane.Source.Common.ChainClient.WaitForEvents())
			errReason, v, err = tc.lane.Source.Common.ChainClient.RevertReasonFromTx(failedTx, evm_2_evm_onramp.EVM2EVMOnRampABI)
			require.NoError(t, err)
			tc.lane.Logger.Info().
				Str("Revert Reason", errReason).
				Interface("Args", v).
				Str("TokensSent", src.TransferAmount[0].String()).
				Str("Token", tc.lane.Source.Common.BridgeTokens[0].Address()).
				Str("FailedTx", failedTx.Hex()).
				Msg("Msg sent with tokens more than AggregateValueRate")
			require.Equal(t, "AggregateValueRateLimitReached", errReason)

			// validate the  successful request was delivered to the destination
			tc.lane.ValidateRequests()

			// now set the token pool rate limit
			if setRateLimit {
				if prevOnRampRLTokenPool.Capacity.Cmp(TokenPoolRateLimitCapacity) != 0 ||
					prevOnRampRLTokenPool.Rate.Cmp(TokenPoolRateLimitRate) != 0 ||
					!prevOnRampRLTokenPool.IsEnabled {
					require.NoError(t, src.Common.BridgeTokenPools[0].SetRemoteChainRateLimits(
						src.DestChainSelector,
						token_pool.RateLimiterConfig{
							IsEnabled: true,
							Capacity:  TokenPoolRateLimitCapacity,
							Rate:      TokenPoolRateLimitRate,
						}), "error setting rate limit on token pool")
					require.NoError(t, src.Common.ChainClient.WaitForEvents(), "waiting for events")
					TokenPoolRateLimitChanged = true
				}
			} else {
				TokenPoolRateLimitCapacity = prevOnRampRLTokenPool.Capacity
				TokenPoolRateLimitRate = prevOnRampRLTokenPool.Rate
			}

			rlOnPool, err := src.Common.BridgeTokenPools[0].Instance.GetCurrentOutboundRateLimiterState(nil, src.DestChainSelector)
			require.NoError(t, err)
			require.True(t, rlOnPool.IsEnabled, "Token Pool rate limiter should be enabled")

			// try to send more than token pool capacity - should fail
			tokensToSend = new(big.Int).Add(TokenPoolRateLimitCapacity, big.NewInt(2))

			// wait for the AggregateCapacity to be refilled
			onRampState, err := src.OnRamp.Instance.CurrentRateLimiterState(nil)
			if err != nil {
				return
			}
			if AggregatedRateLimitCapacity.Cmp(onRampState.Capacity) > 0 {
				capacityToBeFilled := new(big.Int).Sub(AggregatedRateLimitCapacity, onRampState.Capacity)
				durationToFill := time.Duration(new(big.Int).Div(capacityToBeFilled, AggregatedRateLimitRate).Int64())
				tc.lane.Logger.Info().
					Dur("wait duration", durationToFill).
					Str("current capacity", onRampState.Capacity.String()).
					Str("tokensToSend", tokensToSend.String()).
					Msg("Waiting for aggregated capacity to be available")
				time.Sleep(durationToFill * time.Second)
			}

			src.TransferAmount[0] = tokensToSend
			tc.lane.Logger.Info().Str("tokensToSend", tokensToSend.String()).Msg("More than Token Pool Capacity")

			failedTx, _, _, _, err = tc.lane.Source.SendRequest(tc.lane.Dest.ReceiverDapp.EthAddress, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)
			require.Error(t, tc.lane.Source.Common.ChainClient.WaitForEvents())
			errReason, v, err = tc.lane.Source.Common.ChainClient.RevertReasonFromTx(failedTx, lock_release_token_pool.LockReleaseTokenPoolABI)
			require.NoError(t, err)
			tc.lane.Logger.Info().
				Str("Revert Reason", errReason).
				Interface("Args", v).
				Str("TokensSent", src.TransferAmount[0].String()).
				Str("Token", tc.lane.Source.Common.BridgeTokens[0].Address()).
				Str("FailedTx", failedTx.Hex()).
				Msg("Msg sent with tokens more than token pool capacity")
			require.Equal(t, "TokenMaxCapacityExceeded", errReason)

			// try to send 99% of token pool capacity - should succeed
			tokensToSend = new(big.Int).Div(new(big.Int).Mul(TokenPoolRateLimitCapacity, big.NewInt(99)), big.NewInt(100))
			src.TransferAmount[0] = tokensToSend
			tc.lane.Logger.Info().Str("tokensToSend", tokensToSend.String()).Msg("99% of Token Pool Capacity")
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)

			// try to send again with amount more than the amount refilled by token pool rate and
			// this should fail, as the refill rate is not enough to refill the capacity
			tokensToSend = new(big.Int).Mul(TokenPoolRateLimitRate, big.NewInt(20))
			tc.lane.Logger.Info().Str("tokensToSend", tokensToSend.String()).Msg("More than TokenPool Rate")
			src.TransferAmount[0] = tokensToSend
			// approve the tokens
			require.NoError(t, src.Common.BridgeTokens[0].Approve(
				src.Common.ChainClient.GetDefaultWallet(), src.Common.Router.Address(), src.TransferAmount[0]),
			)
			require.NoError(t, tc.lane.Source.Common.ChainClient.WaitForEvents())
			failedTx, _, _, _, err = tc.lane.Source.SendRequest(tc.lane.Dest.ReceiverDapp.EthAddress, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)
			require.Error(t, tc.lane.Source.Common.ChainClient.WaitForEvents())
			errReason, v, err = tc.lane.Source.Common.ChainClient.RevertReasonFromTx(failedTx, lock_release_token_pool.LockReleaseTokenPoolABI)
			require.NoError(t, err)
			tc.lane.Logger.Info().
				Str("Revert Reason", errReason).
				Interface("Args", v).
				Str("TokensSent", src.TransferAmount[0].String()).
				Str("Token", tc.lane.Source.Common.BridgeTokens[0].Address()).
				Str("FailedTx", failedTx.Hex()).
				Msg("Msg sent with tokens more than TokenPool Rate")
			require.Equal(t, "TokenRateLimitReached", errReason)

			// validate that the successful transfers are reflected in destination
			tc.lane.ValidateRequests()
		})
	}
}

func TestSmokeCCIPOnRampLimits(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-124")
	t.Parallel()

	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke, testsetups.WithNoTokensPerMessage(4), testsetups.WithTokensPerChain(4))
	require.False(t, pointer.GetBool(TestCfg.TestGroupInput.ExistingDeployment),
		"This test modifies contract state. Before running it, ensure you are willing and able to do so.",
	)
	err := contracts.MatchContractVersionsOrAbove(map[contracts.Name]contracts.Version{
		contracts.OffRampContract: contracts.V1_5_0,
		contracts.OnRampContract:  contracts.V1_5_0,
	})
	require.NoError(t, err, "Required contract versions not met")

	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, &log, "smoke-ccip", nil, TestCfg)
	if len(setUpOutput.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		require.NoError(t, setUpOutput.TearDown())
	})

	var tests []testDefinition
	for _, lane := range setUpOutput.Lanes {
		tests = append(tests, testDefinition{
			testName: fmt.Sprintf("Network %s to network %s",
				lane.ForwardLane.SourceNetworkName, lane.ForwardLane.DestNetworkName),
			lane: lane.ForwardLane,
		})
	}

	var (
		capacityLimit      = big.NewInt(1e16)
		overCapacityAmount = new(big.Int).Add(capacityLimit, big.NewInt(1))

		// token without any transfer config
		freeTokenIndex = 0
		// token with bps non-zero, no agg rate limit
		bpsTokenIndex = 1
		// token with bps zero, with agg rate limit on
		aggRateTokenIndex = 2
		// token with both bps and agg rate limit
		bpsAndAggTokenIndex = 3
	)

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s - OnRamp Limits", tc.testName), func(t *testing.T) {
			tc.lane.Test = t
			src := tc.lane.Source
			dest := tc.lane.Dest
			require.GreaterOrEqual(t, len(src.Common.BridgeTokens), 2, "At least two bridge tokens needed for test")
			require.GreaterOrEqual(t, len(src.Common.BridgeTokenPools), 2, "At least two bridge token pools needed for test")
			require.GreaterOrEqual(t, len(dest.Common.BridgeTokens), 2, "At least two bridge tokens needed for test")
			require.GreaterOrEqual(t, len(dest.Common.BridgeTokenPools), 2, "At least two bridge token pools needed for test")
			addLiquidity(t, src.Common, new(big.Int).Mul(capacityLimit, big.NewInt(20)))
			addLiquidity(t, dest.Common, new(big.Int).Mul(capacityLimit, big.NewInt(20)))

			var (
				freeToken      = src.Common.BridgeTokens[freeTokenIndex]
				bpsToken       = src.Common.BridgeTokens[bpsTokenIndex]
				aggRateToken   = src.Common.BridgeTokens[aggRateTokenIndex]
				bpsAndAggToken = src.Common.BridgeTokens[bpsAndAggTokenIndex]
			)
			tc.lane.Logger.Info().
				Str("Free Token", freeToken.ContractAddress.Hex()).
				Str("BPS Token", bpsToken.ContractAddress.Hex()).
				Str("Agg Rate Token", aggRateToken.ContractAddress.Hex()).
				Str("BPS and Agg Rate Token", bpsAndAggToken.ContractAddress.Hex()).
				Msg("Tokens for rate limit testing")
			err := tc.lane.DisableAllRateLimiting()
			require.NoError(t, err, "Error disabling rate limits")

			// Set reasonable rate limits for the tokens
			err = src.OnRamp.SetTokenTransferFeeConfig([]evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
				{
					Token:                     bpsToken.ContractAddress,
					AggregateRateLimitEnabled: false,
					DeciBps:                   10,
				},
				{
					Token:                     aggRateToken.ContractAddress,
					AggregateRateLimitEnabled: true,
				},
				{
					Token:                     bpsAndAggToken.ContractAddress,
					AggregateRateLimitEnabled: true,
					DeciBps:                   10,
				},
			})
			require.NoError(t, err, "Error setting OnRamp transfer fee config")
			err = src.OnRamp.SetRateLimit(evm_2_evm_onramp.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  capacityLimit,
				Rate:      new(big.Int).Mul(capacityLimit, big.NewInt(500)), // Set a high rate to avoid it getting in the way
			})
			require.NoError(t, err, "Error setting OnRamp rate limits")
			err = src.Common.ChainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

			// Send all tokens under their limits and ensure they succeed
			src.TransferAmount[freeTokenIndex] = overCapacityAmount
			src.TransferAmount[bpsTokenIndex] = overCapacityAmount
			src.TransferAmount[aggRateTokenIndex] = big.NewInt(1)
			src.TransferAmount[bpsAndAggTokenIndex] = big.NewInt(1)
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)
			tc.lane.ValidateRequests()

			// Check that capacity limits are enforced
			src.TransferAmount[freeTokenIndex] = big.NewInt(0)
			src.TransferAmount[bpsTokenIndex] = big.NewInt(0)
			src.TransferAmount[aggRateTokenIndex] = overCapacityAmount
			src.TransferAmount[bpsAndAggTokenIndex] = big.NewInt(0)
			failedTx, _, _, _, err := tc.lane.Source.SendRequest(tc.lane.Dest.ReceiverDapp.EthAddress, big.NewInt(actions.DefaultDestinationGasLimit))
			require.Error(t, err, "Limited token transfer should immediately revert")
			errReason, _, err := src.Common.ChainClient.RevertReasonFromTx(failedTx, evm_2_evm_onramp.EVM2EVMOnRampABI)
			require.NoError(t, err)
			require.Equal(t, "AggregateValueMaxCapacityExceeded", errReason, "Expected capacity limit reached error")
			tc.lane.Logger.
				Info().
				Str("Token", aggRateToken.ContractAddress.Hex()).
				Msg("Limited token transfer failed on source chain (a good thing in this context)")

			src.TransferAmount[aggRateTokenIndex] = big.NewInt(0)
			src.TransferAmount[bpsAndAggTokenIndex] = overCapacityAmount
			failedTx, _, _, _, err = tc.lane.Source.SendRequest(tc.lane.Dest.ReceiverDapp.EthAddress, big.NewInt(actions.DefaultDestinationGasLimit))
			require.Error(t, err, "Limited token transfer should immediately revert")
			errReason, _, err = src.Common.ChainClient.RevertReasonFromTx(failedTx, evm_2_evm_onramp.EVM2EVMOnRampABI)
			require.NoError(t, err)
			require.Equal(t, "AggregateValueMaxCapacityExceeded", errReason, "Expected capacity limit reached error")
			tc.lane.Logger.
				Info().
				Str("Token", aggRateToken.ContractAddress.Hex()).
				Msg("Limited token transfer failed on source chain (a good thing in this context)")

			// Set a high price for the tokens to more easily trigger aggregate rate limits
			// Aggregate rate limits are based on USD price of the tokens
			err = src.Common.PriceRegistry.UpdatePrices([]contracts.InternalTokenPriceUpdate{
				{
					SourceToken: aggRateToken.ContractAddress,
					UsdPerToken: big.NewInt(100),
				},
				{
					SourceToken: bpsAndAggToken.ContractAddress,
					UsdPerToken: big.NewInt(100),
				},
			}, []contracts.InternalGasPriceUpdate{})
			require.NoError(t, err, "Error updating prices")
			// Enable aggregate rate limiting for the limited tokens
			err = src.OnRamp.SetRateLimit(evm_2_evm_onramp.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  new(big.Int).Mul(capacityLimit, big.NewInt(5000)), // Set a high capacity to avoid it getting in the way
				Rate:      big.NewInt(1),
			})
			require.NoError(t, err, "Error setting OnRamp rate limits")
			err = src.Common.ChainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

			// Send aggregate unlimited tokens and ensure they succeed
			src.TransferAmount[freeTokenIndex] = overCapacityAmount
			src.TransferAmount[bpsTokenIndex] = overCapacityAmount
			src.TransferAmount[aggRateTokenIndex] = big.NewInt(0)
			src.TransferAmount[bpsAndAggTokenIndex] = big.NewInt(0)
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)
			tc.lane.ValidateRequests()

			// Check that aggregate rate limits are enforced on limited tokens
			src.TransferAmount[freeTokenIndex] = big.NewInt(0)
			src.TransferAmount[bpsTokenIndex] = big.NewInt(0)
			src.TransferAmount[aggRateTokenIndex] = capacityLimit
			src.TransferAmount[bpsAndAggTokenIndex] = big.NewInt(0)
			failedTx, _, _, _, err = tc.lane.Source.SendRequest(tc.lane.Dest.ReceiverDapp.EthAddress, big.NewInt(actions.DefaultDestinationGasLimit))
			require.Error(t, err, "Aggregate rate limited token transfer should immediately revert")
			errReason, _, err = src.Common.ChainClient.RevertReasonFromTx(failedTx, evm_2_evm_onramp.EVM2EVMOnRampABI)
			require.NoError(t, err)
			require.Equal(t, "AggregateValueRateLimitReached", errReason, "Expected aggregate rate limit reached error")
			tc.lane.Logger.
				Info().
				Str("Token", aggRateToken.ContractAddress.Hex()).
				Msg("Limited token transfer failed on source chain (a good thing in this context)")

			src.TransferAmount[aggRateTokenIndex] = big.NewInt(0)
			src.TransferAmount[bpsAndAggTokenIndex] = capacityLimit
			failedTx, _, _, _, err = tc.lane.Source.SendRequest(tc.lane.Dest.ReceiverDapp.EthAddress, big.NewInt(actions.DefaultDestinationGasLimit))
			require.Error(t, err, "Aggregate rate limited token transfer should immediately revert")
			errReason, _, err = src.Common.ChainClient.RevertReasonFromTx(failedTx, evm_2_evm_onramp.EVM2EVMOnRampABI)
			require.NoError(t, err)
			require.Equal(t, "AggregateValueRateLimitReached", errReason, "Expected aggregate rate limit reached error")
			tc.lane.Logger.
				Info().
				Str("Token", aggRateToken.ContractAddress.Hex()).
				Msg("Limited token transfer failed on source chain (a good thing in this context)")
		})
	}
}

func TestSmokeCCIPOffRampCapacityLimit(t *testing.T) {
	t.Parallel()

	capacityLimited := contracts.RateLimiterConfig{
		IsEnabled: true,
		Capacity:  big.NewInt(1e16),
		Rate:      new(big.Int).Mul(big.NewInt(1e16), big.NewInt(10)), // Set a high rate limit to avoid it getting in the way
	}
	testOffRampRateLimits(t, capacityLimited)
}

func TestSmokeCCIPOffRampAggRateLimit(t *testing.T) {
	t.Parallel()

	aggRateLimited := contracts.RateLimiterConfig{
		IsEnabled: true,
		Capacity:  new(big.Int).Mul(big.NewInt(1e16), big.NewInt(10)), // Set a high capacity limit to avoid it getting in the way
		Rate:      big.NewInt(1),
	}
	testOffRampRateLimits(t, aggRateLimited)
}

func TestSmokeCCIPTokenPoolRateLimits(t *testing.T) {
	t.Parallel()

	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke, testsetups.WithNoTokensPerMessage(4), testsetups.WithTokensPerChain(4))
	require.False(t, pointer.GetBool(TestCfg.TestGroupInput.ExistingDeployment),
		"This test modifies contract state. Before running it, ensure you are willing and able to do so.",
	)
	err := contracts.MatchContractVersionsOrAbove(map[contracts.Name]contracts.Version{
		contracts.OffRampContract: contracts.V1_5_0,
		contracts.OnRampContract:  contracts.V1_5_0,
	})
	require.NoError(t, err, "Required contract versions not met")

	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, &log, "smoke-ccip", nil, TestCfg)
	if len(setUpOutput.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		require.NoError(t, setUpOutput.TearDown())
	})

	var tests []testDefinition
	for _, lane := range setUpOutput.Lanes {
		tests = append(tests, testDefinition{
			testName: fmt.Sprintf("Network %s to network %s",
				lane.ForwardLane.SourceNetworkName, lane.ForwardLane.DestNetworkName),
			lane: lane.ForwardLane,
		})
	}

	var (
		capacityLimit      = big.NewInt(1e16)
		overCapacityAmount = new(big.Int).Add(capacityLimit, big.NewInt(1))

		// token without any limits
		freeTokenIndex = 0
		// token with rate limits
		limitedTokenIndex = 1
	)

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s - Token Pool Rate Limits", tc.testName), func(t *testing.T) {
			tc.lane.Test = t
			src := tc.lane.Source
			dest := tc.lane.Dest
			require.GreaterOrEqual(t, len(src.Common.BridgeTokens), 2, "At least two bridge tokens needed for test")
			require.GreaterOrEqual(t, len(src.Common.BridgeTokenPools), 2, "At least two bridge token pools needed for test")
			require.GreaterOrEqual(t, len(dest.Common.BridgeTokens), 2, "At least two bridge tokens needed for test")
			require.GreaterOrEqual(t, len(dest.Common.BridgeTokenPools), 2, "At least two bridge token pools needed for test")
			addLiquidity(t, src.Common, new(big.Int).Mul(capacityLimit, big.NewInt(20)))
			addLiquidity(t, dest.Common, new(big.Int).Mul(capacityLimit, big.NewInt(20)))

			var (
				freeToken        = src.Common.BridgeTokens[freeTokenIndex]
				limitedToken     = src.Common.BridgeTokens[limitedTokenIndex]
				limitedTokenPool = src.Common.BridgeTokenPools[limitedTokenIndex]
			)
			tc.lane.Logger.Info().
				Str("Free Token", freeToken.ContractAddress.Hex()).
				Str("Limited Token", limitedToken.ContractAddress.Hex()).
				Msg("Tokens for rate limit testing")
			err := tc.lane.DisableAllRateLimiting() // Make sure this is pure
			require.NoError(t, err, "Error disabling rate limits")

			// Check capacity limits
			err = limitedTokenPool.SetRemoteChainRateLimits(src.DestChainSelector, token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  capacityLimit,
				Rate:      new(big.Int).Sub(capacityLimit, big.NewInt(1)), // Set as high rate as possible to avoid it getting in the way
			})
			require.NoError(t, err, "Error setting token pool rate limit")
			err = src.Common.ChainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

			// Send all tokens under their limits and ensure they succeed
			src.TransferAmount[freeTokenIndex] = overCapacityAmount
			src.TransferAmount[limitedTokenIndex] = big.NewInt(1)
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)
			tc.lane.ValidateRequests()

			// Send limited token over capacity and ensure it fails
			src.TransferAmount[freeTokenIndex] = big.NewInt(0)
			src.TransferAmount[limitedTokenIndex] = overCapacityAmount
			failedTx, _, _, _, err := tc.lane.Source.SendRequest(tc.lane.Dest.ReceiverDapp.EthAddress, big.NewInt(actions.DefaultDestinationGasLimit))
			require.Error(t, err, "Limited token transfer should immediately revert")
			errReason, _, err := src.Common.ChainClient.RevertReasonFromTx(failedTx, lock_release_token_pool.LockReleaseTokenPoolABI)
			require.NoError(t, err)
			require.Equal(t, "TokenMaxCapacityExceeded", errReason, "Expected token capacity error")
			tc.lane.Logger.
				Info().
				Str("Token", limitedToken.ContractAddress.Hex()).
				Msg("Limited token transfer failed on source chain (a good thing in this context)")

			// Check rate limit
			err = limitedTokenPool.SetRemoteChainRateLimits(src.DestChainSelector, token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  new(big.Int).Mul(capacityLimit, big.NewInt(2)), // Set a high capacity to avoid it getting in the way
				Rate:      big.NewInt(1),
			})
			require.NoError(t, err, "Error setting token pool rate limit")
			err = src.Common.ChainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

			// Send all tokens under their limits and ensure they succeed
			src.TransferAmount[freeTokenIndex] = overCapacityAmount
			src.TransferAmount[limitedTokenIndex] = capacityLimit
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)
			tc.lane.ValidateRequests()

			// Send limited token over rate limit and ensure it fails
			src.TransferAmount[freeTokenIndex] = big.NewInt(0)
			src.TransferAmount[limitedTokenIndex] = capacityLimit
			failedTx, _, _, _, err = tc.lane.Source.SendRequest(tc.lane.Dest.ReceiverDapp.EthAddress, big.NewInt(actions.DefaultDestinationGasLimit))
			require.Error(t, err, "Limited token transfer should immediately revert")
			errReason, _, err = src.Common.ChainClient.RevertReasonFromTx(failedTx, lock_release_token_pool.LockReleaseTokenPoolABI)
			require.NoError(t, err)
			require.Equal(t, "TokenRateLimitReached", errReason, "Expected rate limit reached error")
			tc.lane.Logger.
				Info().
				Str("Token", limitedToken.ContractAddress.Hex()).
				Msg("Limited token transfer failed on source chain (a good thing in this context)")
		})
	}
}

func TestSmokeCCIPMulticall(t *testing.T) {
	t.Parallel()
	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	// enable multicall in one tx for this test
	TestCfg.TestGroupInput.MulticallInOneTx = ptr.Ptr(true)
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, &log, "smoke-ccip", nil, TestCfg)
	if len(setUpOutput.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		require.NoError(t, setUpOutput.TearDown())
	})

	var tests []testDefinition
	for _, lane := range setUpOutput.Lanes {
		tests = append(tests, testDefinition{
			testName: fmt.Sprintf("CCIP message transfer from network %s to network %s",
				lane.ForwardLane.SourceNetworkName, lane.ForwardLane.DestNetworkName),
			lane: lane.ForwardLane,
		})
		if lane.ReverseLane != nil {
			tests = append(tests, testDefinition{
				testName: fmt.Sprintf("CCIP message transfer from network %s to network %s",
					lane.ReverseLane.SourceNetworkName, lane.ReverseLane.DestNetworkName),
				lane: lane.ReverseLane,
			})
		}
	}

	log.Info().Int("Total Lanes", len(tests)).Msg("Starting CCIP test")
	for _, test := range tests {
		tc := test
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			tc.lane.Test = t
			log.Info().
				Str("Source", tc.lane.SourceNetworkName).
				Str("Destination", tc.lane.DestNetworkName).
				Msgf("Starting lane %s -> %s", tc.lane.SourceNetworkName, tc.lane.DestNetworkName)

			tc.lane.RecordStateBeforeTransfer()
			err := tc.lane.Multicall(TestCfg.TestGroupInput.NoOfSendsInMulticall, tc.lane.Source.Common.MulticallContract)
			require.NoError(t, err)
			tc.lane.ValidateRequests()
		})
	}
}

func TestSmokeCCIPManuallyExecuteAfterExecutionFailingDueToInsufficientGas(t *testing.T) {
	t.Parallel()
	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, &log, "smoke-ccip", nil, TestCfg)
	if len(setUpOutput.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		if TestCfg.TestGroupInput.MsgDetails.IsTokenTransfer() {
			setUpOutput.Balance.Verify(t)
		}
		require.NoError(t, setUpOutput.TearDown())
	})

	var tests []testDefinition
	for _, lane := range setUpOutput.Lanes {
		tests = append(tests, testDefinition{
			testName: fmt.Sprintf("CCIP message transfer from network %s to network %s",
				lane.ForwardLane.SourceNetworkName, lane.ForwardLane.DestNetworkName),
			lane: lane.ForwardLane,
		})
		if lane.ReverseLane != nil {
			tests = append(tests, testDefinition{
				testName: fmt.Sprintf("CCIP message transfer from network %s to network %s",
					lane.ReverseLane.SourceNetworkName, lane.ReverseLane.DestNetworkName),
				lane: lane.ReverseLane,
			})
		}
	}

	log.Info().Int("Total Lanes", len(tests)).Msg("Starting CCIP test")
	for _, test := range tests {
		tc := test
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			tc.lane.Test = t
			log.Info().
				Str("Source", tc.lane.SourceNetworkName).
				Str("Destination", tc.lane.DestNetworkName).
				Msgf("Starting lane %s -> %s", tc.lane.SourceNetworkName, tc.lane.DestNetworkName)

			tc.lane.RecordStateBeforeTransfer()
			// send with insufficient gas for ccip-receive to fail
			err := tc.lane.SendRequests(1, big.NewInt(0))
			require.NoError(t, err)
			tc.lane.ValidateRequests(actions.ExpectPhaseToFail(testreporters.ExecStateChanged))
			// wait for events
			err = tc.lane.Dest.Common.ChainClient.WaitForEvents()
			require.NoError(t, err)
			// execute all failed ccip requests manually
			err = tc.lane.ExecuteManually()
			require.NoError(t, err)
			if len(tc.lane.Source.TransferAmount) > 0 {
				tc.lane.Source.UpdateBalance(int64(tc.lane.NumberOfReq), tc.lane.TotalFee, tc.lane.Balance)
				tc.lane.Dest.UpdateBalance(tc.lane.Source.TransferAmount, int64(tc.lane.NumberOfReq), tc.lane.Balance)
			}
		})
	}
}

// Test expects to generate below finality reorg in both source and destination and
// expect CCIP transactions to go through successful.
func TestSmokeCCIPReorgBelowFinality(t *testing.T) {
	t.Parallel()
	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	gasLimit := big.NewInt(*TestCfg.TestGroupInput.MsgDetails.DestGasLimit)
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, &log, "smoke-ccip", nil, TestCfg)
	require.False(t, len(setUpOutput.Lanes) == 0, "No lanes found.")
	t.Cleanup(func() {
		require.NoError(t, setUpOutput.TearDown())
	})

	lane := setUpOutput.Lanes[0].ForwardLane
	log.Info().
		Str("Source", lane.SourceNetworkName).
		Str("Destination", lane.DestNetworkName).
		Msg("Starting CCIP reorg test")
	t.Run(fmt.Sprintf("CCIP reorg below finality test from network %s to network %s",
		lane.SourceNetworkName, lane.DestNetworkName), func(t *testing.T) {
		t.Parallel()
		lane.Test = t
		lane.RecordStateBeforeTransfer()
		// sending multiple request and expect all should go through though there is below finality reorg
		err := lane.SendRequests(5, gasLimit)
		require.NoError(t, err, "Send requests failed")
		rs := SetupReorgSuite(t, &log, setUpOutput)
		// run below finality reorg in both source and destination chain
		blocksBackSrc := rs.Cfg.SrcFinalityDepth - rs.Cfg.FinalityDelta
		blocksBackDst := rs.Cfg.DstFinalityDepth - rs.Cfg.FinalityDelta
		rs.RunReorg(rs.DstClient, blocksBackSrc, "Source", 2*time.Second)
		rs.RunReorg(rs.DstClient, blocksBackDst, "Destination", 2*time.Second)
		time.Sleep(1 * time.Minute)
		lane.ValidateRequests()
	})
}

// Test creates above finality reorg at destination and
// expects ccip transactions in-flight and the one initiated after reorg
// doesn't go through and verifies f+1 nodes are able to detect reorg.
// Note: LogPollInterval interval is set as 1s to detect the reorg immediately
func TestSmokeCCIPReorgAboveFinalityAtDestination(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-4401")
	t.Parallel()
	t.Run("Above finality reorg in destination chain", func(t *testing.T) {
		performAboveFinalityReorgAndValidate(t, "Destination")
	})
}

// Test creates above finality reorg at destination and
// expects ccip transactions in-flight doesn't go through, the transaction initiated after reorg
// shouldn't even get initiated and verifies f+1 nodes are able to detect reorg.
// Note: LogPollInterval interval is set as 1s to detect the reorg immediately
func TestSmokeCCIPReorgAboveFinalityAtSource(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-4401")
	t.Parallel()
	t.Run("Above finality reorg in source chain", func(t *testing.T) {
		performAboveFinalityReorgAndValidate(t, "Source")
	})
}

// TestSmokeCCIPForGivenNetworkPairs is designed specifically for scheduled mainnet testing. This test checks for recent
// transaction and skip the lanes accordingly. This test also has capability to take override input on network pairs and phase timeout.
func TestSmokeCCIPForGivenNetworkPairs(t *testing.T) {
	t.Parallel()
	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	// override network pairs
	var temp []testsetups.NetworkPair
	overrideNetworkPairs, err := osutil.GetEnv("OVERRIDE_NETWORK_PAIRS")
	require.NoError(t, err, "Error getting OVERRIDE_NETWORK_PAIRS environment variable")
	if overrideNetworkPairs != "" {
		networkPairs := strings.Split(overrideNetworkPairs, ";")
		for _, networkPair := range networkPairs {
			// check for any malformed inputs
			if !strings.Contains(networkPair, ",") || len(strings.Split(networkPair, ",")) != 2 {
				log.Error().Msgf("malformed OVERRIDE_NETWORK_PAIRS environment variable for network pair: %s ", networkPair)
				return
			}
			networkPair = strings.ToUpper(strings.ReplaceAll(networkPair, "_", " "))
			for _, network := range TestCfg.NetworkPairs {
				if strings.Contains(networkPair, strings.ToUpper(network.NetworkA.Name)) && strings.Contains(networkPair, strings.ToUpper(network.NetworkB.Name)) {
					temp = append(temp, network)
					break
				}
			}
		}
		log.Info().Int("Pairs", len(temp)).Msg("Number of lanes overridden in the test")
		log.Info().Interface("Lanes", networkPairs).Msg("Lanes under test")
		TestCfg.NetworkPairs = temp
	}

	// phase timeout override
	phaseTimeout, err := osutil.GetEnv("OVERRIDE_PHASE_TIMEOUT")
	require.NoError(t, err, "Error getting OVERRIDE_PHASE_TIMEOUT environment variable")
	if phaseTimeout != "" {
		configDuration, err := config.ParseDuration(phaseTimeout)
		require.NoError(t, err, "Error parsing phase timeout value")
		TestCfg.TestGroupInput.PhaseTimeout = &configDuration
		log.Info().Float64("Timeout in minutes", configDuration.Duration().Minutes()).Msg("Phase timeout is overridden")
	}

	gasLimit := big.NewInt(*TestCfg.TestGroupInput.MsgDetails.DestGasLimit)
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, &log, "smoke-ccip", nil, TestCfg)
	if len(setUpOutput.Lanes) == 0 {
		log.Error().Msg("No lanes found")
		return
	}

	t.Cleanup(func() {
		// If we are running a test that is a token transfer, we need to verify the balance.
		// skip the balance check for existing deployment, there can be multiple external requests in progress for existing deployments
		// other than token transfer initiated by the test, which can affect the balance check
		// therefore we check the balance only for the ccip environment created by the test
		if TestCfg.TestGroupInput.MsgDetails.IsTokenTransfer() &&
			!pointer.GetBool(TestCfg.TestGroupInput.USDCMockDeployment) &&
			!pointer.GetBool(TestCfg.TestGroupInput.ExistingDeployment) {
			setUpOutput.Balance.Verify(t)
		}
		require.NoError(t, setUpOutput.TearDown(), "error in tear down step")
	})

	var tests []testDefinition
	lookBackDuration := TestCfg.TestGroupInput.SkipRequestIfAnotherRequestTriggeredWithin
	var recentTxFound *types.Log

	addLanesToTest := func(lane *actions.CCIPLane) {
		// Create test definitions for given lane if no previous request has been triggered within the specified timeframe.
		// By default, the timeframe is set to nil. To define a timeframe, assign a duration to the variable
		// SkipRequestIfAnotherRequestTriggeredWithin.
		if lookBackDuration != nil {
			recentTxFound, err = lane.Source.IsPastRequestTriggeredWithinTimeframe(lane.Context, lookBackDuration)
			require.NoError(t, err, "error while finding recent request for lane network %s to network %s",
				lane.SourceNetworkName, lane.DestNetworkName)
		}
		if recentTxFound == nil {
			tests = append(tests, testDefinition{
				testName: fmt.Sprintf("CCIP message transfer from network %s to network %s",
					lane.SourceNetworkName, lane.DestNetworkName),
				lane: lane,
			})
		} else {
			log.Info().
				Str("TX", recentTxFound.TxHash.Hex()).
				Uint64("Block Number", recentTxFound.BlockNumber).
				Str("Source", lane.SourceNetworkName).
				Str("Dest", lane.DestNetworkName).
				Msgf("Lane Skipped. Recent request found within %v minutes.", lookBackDuration.Duration().Minutes())
		}
	}
	for _, lane := range setUpOutput.Lanes {
		addLanesToTest(lane.ForwardLane)
		if lane.ReverseLane != nil {
			recentTxFound = nil
			addLanesToTest(lane.ReverseLane)
		}
	}

	// Execute tests.
	log.Info().Int("Total Lanes", len(tests)).Msg("Starting CCIP test")
	for _, test := range tests {
		tc := test
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			tc.lane.Test = t
			log.Info().
				Str("Source", tc.lane.SourceNetworkName).
				Str("Destination", tc.lane.DestNetworkName).
				Msgf("Starting lane %s -> %s", tc.lane.SourceNetworkName, tc.lane.DestNetworkName)

			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, gasLimit)
			require.NoError(t, err, "error sending requests")
			tc.lane.ValidateRequests()
		})
	}
}

// performAboveFinalityReorgAndValidate is to perform the above finality reorg test
func performAboveFinalityReorgAndValidate(t *testing.T, network string) {
	t.Helper()

	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	gasLimit := big.NewInt(*TestCfg.TestGroupInput.MsgDetails.DestGasLimit)
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, &log, "smoke-ccip", nil, TestCfg)
	require.False(t, len(setUpOutput.Lanes) == 0, "No lanes found.")
	t.Cleanup(func() {
		require.NoError(t, setUpOutput.TearDown())
	})
	rs := SetupReorgSuite(t, &log, setUpOutput)
	lane := setUpOutput.Lanes[0].ForwardLane
	log.Info().
		Str("Source", lane.SourceNetworkName).
		Str("Destination", lane.DestNetworkName).
		Msg("Starting ccip reorg test")
	lane.Test = t
	lane.RecordStateBeforeTransfer()
	err := lane.SendRequests(1, gasLimit)
	require.NoError(t, err, "Send requests failed")
	logPollerName := ""
	if network == "Destination" {
		logPollerName = fmt.Sprintf("EVM.%d.LogPoller", lane.DestChain.GetChainID())
		rs.RunReorg(rs.DstClient, rs.Cfg.DstFinalityDepth+rs.Cfg.FinalityDelta, network, 2*time.Second)
	} else {
		logPollerName = fmt.Sprintf("EVM.%d.LogPoller", lane.SourceChain.GetChainID())
		rs.RunReorg(rs.SrcClient, rs.Cfg.SrcFinalityDepth+rs.Cfg.FinalityDelta, network, 2*time.Second)
	}
	// DON is 3F+1, finding f+1 from the given number of nodes in the environment
	fPlus1Nodes := int(math.Ceil(float64(len(setUpOutput.Env.CLNodes)-1)/3)) + 1
	// assert at least f+1 nodes is detecting the reorg (LogPollInterval is set as 1s for faster detection)
	// additional context: Commit requires 2f+1 observations, so f+1 nodes need to detect it in order to force the entire DON to stop processing messages.
	nodesDetectedViolation := make(map[string]bool)
	assert.Eventually(t, func() bool {
		for _, node := range setUpOutput.Env.CLNodes {
			if _, ok := nodesDetectedViolation[node.ChainlinkClient.URL()]; ok {
				continue
			}
			resp, _, err := node.Health()
			require.NoError(t, err)
			for _, d := range resp.Data {
				if d.Attributes.Name == logPollerName && d.Attributes.Output == "finality violated" && d.Attributes.Status == "failing" {
					log.Debug().Str("Node", node.ChainlinkClient.URL()).Msg("Finality violated is detected by node")
					nodesDetectedViolation[node.ChainlinkClient.URL()] = true
				}
			}
		}
		return len(nodesDetectedViolation) >= fPlus1Nodes
	}, 3*time.Minute, 20*time.Second, "Reorg above finality depth is not detected by f+1 nodes")
	log.Debug().Interface("Nodes", nodesDetectedViolation).Msg("Violation detection details")
	// send another request and verify it fails
	err = lane.SendRequests(1, gasLimit)
	if network == "Source" {
		// if it is source chain reorg, the transaction will not even be initiated
		require.Error(t, err,
			"CCIP send transaction shouldn't be initiated as there is above finality depth reorg in source chain")
	} else {
		// if it is destination chain reorg, the transaction will be initiated and will fail in the process
		require.NoError(t, err,
			"CCIP send transaction should be initiated even when there above finality reorg in dest chain")
	}

	lane.ValidateRequests(actions.ExpectAnyPhaseToFail(actions.WithTimeout(time.Minute)))
}

// add liquidity to pools on both networks
func addLiquidity(t *testing.T, ccipCommon *actions.CCIPCommon, amount *big.Int) {
	t.Helper()

	for i, btp := range ccipCommon.BridgeTokenPools {
		token := ccipCommon.BridgeTokens[i]
		err := btp.AddLiquidity(
			token, token.OwnerWallet, amount,
		)
		require.NoError(t, err)
	}
}

// testOffRampRateLimits tests the rate limiting functionality of the OffRamp contract
// it's broken into a helper to help parallelize and keep the tests DRY
func testOffRampRateLimits(t *testing.T, rateLimiterConfig contracts.RateLimiterConfig) {
	t.Helper()

	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	require.False(t, pointer.GetBool(TestCfg.TestGroupInput.ExistingDeployment),
		"This test modifies contract state. Before running it, ensure you are willing and able to do so.",
	)
	err := contracts.MatchContractVersionsOrAbove(map[contracts.Name]contracts.Version{
		contracts.OffRampContract: contracts.V1_5_0,
	})
	require.NoError(t, err, "Required contract versions not met")
	require.False(t, pointer.GetBool(TestCfg.TestGroupInput.ExistingDeployment), "This test modifies contract state and cannot be run on existing deployments")

	// Set the default permissionless exec threshold lower so that we can manually execute the transactions faster
	// Tuning this too low stops any transactions from being realistically executed
	actions.DefaultPermissionlessExecThreshold = 1 * time.Minute

	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, &log, "smoke-ccip", nil, TestCfg)
	if len(setUpOutput.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		require.NoError(t, setUpOutput.TearDown())
	})

	var tests []testDefinition
	for _, lane := range setUpOutput.Lanes {
		tests = append(tests, testDefinition{
			testName: fmt.Sprintf("Network %s to network %s",
				lane.ForwardLane.SourceNetworkName, lane.ForwardLane.DestNetworkName),
			lane: lane.ForwardLane,
		})
	}

	var (
		freeTokenIndex    = 0
		limitedTokenIndex = 1
	)

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s - OffRamp Limits", tc.testName), func(t *testing.T) {
			tc.lane.Test = t
			src := tc.lane.Source
			dest := tc.lane.Dest
			var (
				capacityLimit   = rateLimiterConfig.Capacity
				overLimitAmount = new(big.Int).Add(capacityLimit, big.NewInt(1))
			)
			require.GreaterOrEqual(t, len(src.Common.BridgeTokens), 2, "At least two bridge tokens needed for test")
			require.GreaterOrEqual(t, len(src.Common.BridgeTokenPools), 2, "At least two bridge token pools needed for test")
			require.GreaterOrEqual(t, len(dest.Common.BridgeTokens), 2, "At least two bridge tokens needed for test")
			require.GreaterOrEqual(t, len(dest.Common.BridgeTokenPools), 2, "At least two bridge token pools needed for test")
			addLiquidity(t, src.Common, new(big.Int).Mul(capacityLimit, big.NewInt(20)))
			addLiquidity(t, dest.Common, new(big.Int).Mul(capacityLimit, big.NewInt(20)))

			var (
				freeSrcToken     = src.Common.BridgeTokens[freeTokenIndex]
				freeDestToken    = dest.Common.BridgeTokens[freeTokenIndex]
				limitedSrcToken  = src.Common.BridgeTokens[limitedTokenIndex]
				limitedDestToken = dest.Common.BridgeTokens[limitedTokenIndex]
			)
			tc.lane.Logger.Info().
				Str("Free Source Token", freeSrcToken.Address()).
				Str("Free Dest Token", freeDestToken.Address()).
				Str("Limited Source Token", limitedSrcToken.Address()).
				Str("Limited Dest Token", limitedDestToken.Address()).
				Msg("Tokens for rate limit testing")

			err := tc.lane.DisableAllRateLimiting()
			require.NoError(t, err, "Error disabling rate limits")

			// Send both tokens with no rate limits and ensure they succeed
			src.TransferAmount[freeTokenIndex] = overLimitAmount
			src.TransferAmount[limitedTokenIndex] = overLimitAmount
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)
			tc.lane.ValidateRequests()

			// Enable capacity limiting on the destination chain for the limited token
			err = dest.AddRateLimitTokens([]*contracts.ERC20Token{limitedSrcToken}, []*contracts.ERC20Token{limitedDestToken})
			require.NoError(t, err, "Error setting destination rate limits")
			err = dest.OffRamp.SetRateLimit(rateLimiterConfig)
			require.NoError(t, err, "Error setting destination rate limits")
			err = dest.Common.ChainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")
			tc.lane.Logger.Debug().Str("Token", limitedSrcToken.ContractAddress.Hex()).Msg("Enabled capacity limit on destination chain")

			// Send free token that should not have a rate limit and should succeed
			src.TransferAmount[freeTokenIndex] = overLimitAmount
			src.TransferAmount[limitedTokenIndex] = big.NewInt(0)
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err, "Free token transfer failed")
			tc.lane.ValidateRequests()
			tc.lane.Logger.Info().Str("Token", freeSrcToken.ContractAddress.Hex()).Msg("Free token transfer succeeded")

			// Send limited token with rate limit that should fail on the destination chain
			src.TransferAmount[freeTokenIndex] = big.NewInt(0)
			src.TransferAmount[limitedTokenIndex] = overLimitAmount
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err, "Failed to send rate limited token transfer")

			// We should see the ExecStateChanged phase fail on the OffRamp
			tc.lane.ValidateRequests(actions.ExpectPhaseToFail(testreporters.ExecStateChanged))
			tc.lane.Logger.Info().
				Str("Token", limitedSrcToken.ContractAddress.Hex()).
				Msg("Limited token transfer failed on destination chain (a good thing in this context)")

			// Manually execute the rate limited token transfer and expect a similar error
			tc.lane.Logger.Info().Str("Wait Time", actions.DefaultPermissionlessExecThreshold.String()).Msg("Waiting for Exec Threshold to Expire")
			time.Sleep(actions.DefaultPermissionlessExecThreshold) // Give time to exit the window
			// See above comment on timeout
			err = tc.lane.ExecuteManually(actions.WithConfirmationTimeout(time.Minute))
			require.Error(t, err, "There should be errors executing manually at this point")
			tc.lane.Logger.Debug().Str("Error", err.Error()).Msg("Manually executed rate limited token transfer failed as expected")

			// Change limits to make it viable
			err = dest.OffRamp.SetRateLimit(contracts.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  new(big.Int).Mul(capacityLimit, big.NewInt(100)),
				Rate:      new(big.Int).Mul(capacityLimit, big.NewInt(100)),
			})
			require.NoError(t, err, "Error setting destination rate limits")
			err = dest.Common.ChainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

			// Execute again manually and expect a pass
			err = tc.lane.ExecuteManually()
			require.NoError(t, err, "Error manually executing transaction after rate limit is lifted")
		})
	}
}

// SetupReorgSuite defines the setup required to perform re-org step
func SetupReorgSuite(t *testing.T, lggr *zerolog.Logger, setupOutput *testsetups.CCIPTestSetUpOutputs) *actions.ReorgSuite {
	var finalitySrc int
	var finalityDst int
	if setupOutput.Cfg.SelectedNetworks[0].FinalityTag {
		finalitySrc = 10
	} else {
		finalityDepth := setupOutput.Cfg.SelectedNetworks[0].FinalityDepth
		if finalityDepth > math.MaxInt {
			t.Fatalf("source finality depth overflows int: %d", finalityDepth)
		}
		finalitySrc = int(finalityDepth)
	}
	if setupOutput.Cfg.SelectedNetworks[1].FinalityTag {
		finalityDst = 10
	} else {
		finalityDepth := setupOutput.Cfg.SelectedNetworks[1].FinalityDepth
		if finalityDepth > math.MaxInt {
			t.Fatalf("destination finality depth overflows int: %d", finalityDepth)
		}
		finalityDst = int(finalityDepth)
	}
	var srcGethHTTPURL, dstGethHTTPURL string
	if setupOutput.Env.LocalCluster != nil {
		srcGethHTTPURL = setupOutput.Env.LocalCluster.EVMNetworks[0].HTTPURLs[0]
		dstGethHTTPURL = setupOutput.Env.LocalCluster.EVMNetworks[1].HTTPURLs[0]
	} else {
		srcGethHTTPURL = setupOutput.Env.K8Env.URLs["source-chain_http"][0]
		dstGethHTTPURL = setupOutput.Env.K8Env.URLs["dest-chain_http"][0]
	}
	rs, err := actions.NewReorgSuite(t, lggr, &actions.ReorgConfig{
		SrcGethHTTPURL:   srcGethHTTPURL,
		DstGethHTTPURL:   dstGethHTTPURL,
		SrcFinalityDepth: finalitySrc,
		DstFinalityDepth: finalityDst,
		FinalityDelta:    setupOutput.Cfg.TestGroupInput.ReorgProfile.FinalityDelta,
	})
	require.NoError(t, err)
	return rs
}
