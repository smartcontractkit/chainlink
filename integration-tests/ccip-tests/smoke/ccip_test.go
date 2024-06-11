package smoke

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"

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
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, log, "smoke-ccip", nil, TestCfg)
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
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, log, "smoke-ccip", nil, TestCfg)
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
				addFund := func(ccipCommon *actions.CCIPCommon) {
					for i, btp := range ccipCommon.BridgeTokenPools {
						token := ccipCommon.BridgeTokens[i]
						err := btp.AddLiquidity(token.Approve, token.Address(), new(big.Int).Mul(AggregatedRateLimitCapacity, big.NewInt(20)))
						require.NoError(t, err)
					}
				}
				addFund(src.Common)
				addFund(tc.lane.Dest.Common)
			}
			log.Info().
				Str("Source", tc.lane.SourceNetworkName).
				Str("Destination", tc.lane.DestNetworkName).
				Msgf("Starting lane %s -> %s", tc.lane.SourceNetworkName, tc.lane.DestNetworkName)

			// capture the rate limit config before we change it
			prevRLOnRamp, err := src.OnRamp.Instance.CurrentRateLimiterState(nil)
			require.NoError(t, err)
			tc.lane.Logger.Info().Interface("rate limit", prevRLOnRamp).Msg("Initial OnRamp rate limiter state")

			prevOnRampRLTokenPool, err := src.Common.BridgeTokenPools[0].Instance.GetCurrentOutboundRateLimiterState(nil, tc.lane.Source.DestChainSelector) // TODO RENS maybe?
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

			prevOffRampRLTokenPool, err := tc.lane.Dest.Common.BridgeTokenPools[0].Instance.GetCurrentInboundRateLimiterState(nil, tc.lane.Dest.SourceChainSelector) // TODO RENS maybe?
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
			require.NoError(t, src.Common.BridgeTokens[0].Approve(src.Common.Router.Address(), src.TransferAmount[0]))
			require.NoError(t, tc.lane.Source.Common.ChainClient.WaitForEvents())
			failedTx, _, _, err := tc.lane.Source.SendRequest(
				tc.lane.Dest.ReceiverDapp.EthAddress,
				big.NewInt(600_000), // gas limit
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
			err = tc.lane.SendRequests(1, big.NewInt(600_000))
			require.NoError(t, err)

			// try to send again with amount more than the amount refilled by rate and
			// this should fail, as the refill rate is not enough to refill the capacity
			src.TransferAmount[0] = new(big.Int).Mul(AggregatedRateLimitRate, big.NewInt(10))
			failedTx, _, _, err = tc.lane.Source.SendRequest(
				tc.lane.Dest.ReceiverDapp.EthAddress,
				big.NewInt(600_000), // gas limit
			)
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

			failedTx, _, _, err = tc.lane.Source.SendRequest(
				tc.lane.Dest.ReceiverDapp.EthAddress,
				big.NewInt(600_000), // gas limit
			)
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
			err = tc.lane.SendRequests(1, big.NewInt(600_000))
			require.NoError(t, err)

			// try to send again with amount more than the amount refilled by token pool rate and
			// this should fail, as the refill rate is not enough to refill the capacity
			tokensToSend = new(big.Int).Mul(TokenPoolRateLimitRate, big.NewInt(20))
			tc.lane.Logger.Info().Str("tokensToSend", tokensToSend.String()).Msg("More than TokenPool Rate")
			src.TransferAmount[0] = tokensToSend
			// approve the tokens
			require.NoError(t, src.Common.BridgeTokens[0].Approve(src.Common.Router.Address(), src.TransferAmount[0]))
			require.NoError(t, tc.lane.Source.Common.ChainClient.WaitForEvents())
			failedTx, _, _, err = tc.lane.Source.SendRequest(
				tc.lane.Dest.ReceiverDapp.EthAddress,
				big.NewInt(600_000),
			)
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

func TestSmokeCCIPSelfServeRateLimitOnRamp(t *testing.T) {
	t.Parallel()

	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	if offRampVersion, exists := TestCfg.VersionInput[contracts.OffRampContract]; exists {
		require.NotEqual(t, offRampVersion, contracts.V1_2_0, "Provided OffRamp contract version '%s' is not supported for this test", offRampVersion)
	} else {
		require.FailNow(t, "OffRamp contract version not found in test config")
	}
	if onRampVersion, exists := TestCfg.VersionInput[contracts.OnRampContract]; exists {
		require.NotEqual(t, onRampVersion, contracts.V1_2_0, "Provided OnRamp contract version '%s' is not supported for this test", onRampVersion)
	} else {
		require.FailNow(t, "OnRamp contract version not found in test config")
	}

	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, log, "smoke-ccip", nil, TestCfg)
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

	aggregateRateLimit := big.NewInt(1e16)

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s - Self Serve Rate Limit", tc.testName), func(t *testing.T) {
			tc.lane.Test = t
			src := tc.lane.Source
			dest := tc.lane.Dest
			require.GreaterOrEqual(t, len(src.Common.BridgeTokens), 2, "At least two bridge tokens needed for test")
			require.GreaterOrEqual(t, len(src.Common.BridgeTokenPools), 2, "At least two bridge token pools needed for test")
			require.GreaterOrEqual(t, len(dest.Common.BridgeTokens), 2, "At least two bridge tokens needed for test")
			require.GreaterOrEqual(t, len(dest.Common.BridgeTokenPools), 2, "At least two bridge token pools needed for test")
			// add liquidity to pools on both networks
			if !pointer.GetBool(TestCfg.TestGroupInput.ExistingDeployment) {
				addFund := func(ccipCommon *actions.CCIPCommon) {
					for i, btp := range ccipCommon.BridgeTokenPools {
						token := ccipCommon.BridgeTokens[i]
						err := btp.AddLiquidity(token.Approve, token.Address(), new(big.Int).Mul(aggregateRateLimit, big.NewInt(20)))
						require.NoError(t, err)
					}
				}
				addFund(src.Common)
				addFund(dest.Common)
			}

			var (
				freeTokenIndex    = 0
				limitedTokenIndex = 1

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
			overLimitAmount := new(big.Int).Add(aggregateRateLimit, big.NewInt(1))
			src.TransferAmount[freeTokenIndex] = overLimitAmount
			src.TransferAmount[limitedTokenIndex] = overLimitAmount
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(600_000))
			require.NoError(t, err)
			tc.lane.ValidateRequests()

			// Enable aggregate rate limiting on the destination and source chains for the limited token
			err = dest.AddRateLimitTokens([]*contracts.ERC20Token{limitedSrcToken}, []*contracts.ERC20Token{limitedDestToken})
			require.NoError(t, err, "Error setting destination rate limits")
			err = dest.OffRamp.SetRateLimit(contracts.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  aggregateRateLimit,
				Rate:      aggregateRateLimit,
			})
			require.NoError(t, err, "Error setting destination rate limits")
			err = dest.Common.ChainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")
			tc.lane.Logger.Debug().Str("Token", limitedSrcToken.ContractAddress.Hex()).Msg("Enabled aggregate rate limit on destination chain")

			err = src.OnRamp.SetTokenTransferFeeConfig([]evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
				{
					Token:                     limitedSrcToken.ContractAddress,
					AggregateRateLimitEnabled: true,
				},
			})
			require.NoError(t, err, "Error setting OnRamp rate limits")
			err = src.OnRamp.SetRateLimit(evm_2_evm_onramp.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  aggregateRateLimit,
				Rate:      aggregateRateLimit,
			})
			require.NoError(t, err, "Error setting OnRamp rate limits")
			err = src.Common.ChainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

			// Send free token that should not have a rate limit and should succeed
			src.TransferAmount[freeTokenIndex] = overLimitAmount
			src.TransferAmount[limitedTokenIndex] = big.NewInt(0)
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(600_000))
			require.NoError(t, err, "Free token transfer failed")
			tc.lane.ValidateRequests()
			tc.lane.Logger.Info().Str("Token", freeSrcToken.ContractAddress.Hex()).Msg("Free token transfer succeeded")

			// Send limited token with rate limit that should fail and revert on the source chain
			src.TransferAmount[freeTokenIndex] = big.NewInt(0)
			src.TransferAmount[limitedTokenIndex] = overLimitAmount
			tc.lane.Logger.Info().Str("Token", limitedSrcToken.ContractAddress.Hex()).Msg("Enabled aggregate rate limit on OnRamp")
			failedTx, _, _, err := tc.lane.Source.SendRequest(tc.lane.Dest.ReceiverDapp.EthAddress, big.NewInt(600_000))
			require.Error(t, err, "Limited token transfer should immediately revert")
			errReason, _, err := src.Common.ChainClient.RevertReasonFromTx(failedTx, evm_2_evm_onramp.EVM2EVMOnRampABI)
			require.NoError(t, err)
			require.Equal(t, "AggregateValueMaxCapacityExceeded", errReason, "Expected rate limit reached error")
			tc.lane.Logger.
				Info().
				Str("Token", limitedSrcToken.ContractAddress.Hex()).
				Msg("Limited token transfer failed on source chain (a good thing in this context)")
		})
	}
}

func TestSmokeCCIPSelfServeRateLimitOffRamp(t *testing.T) {
	t.Parallel()

	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	if offRampVersion, exists := TestCfg.VersionInput[contracts.OffRampContract]; exists {
		require.NotEqual(t, offRampVersion, contracts.V1_2_0, "Provided OffRamp contract version '%s' is not supported for this test", offRampVersion)
	} else {
		require.FailNow(t, "OffRamp contract version not found in test config")
	}
	require.True(t, TestCfg.SelectedNetworks[0].Simulated, "This test relies on timing assumptions and should only be run on simulated networks")

	// Set the default permissionless exec threshold lower so that we can manually execute the transactions faster
	// Tuning this too low stops any transactions from being realistically executed
	actions.DefaultPermissionlessExecThreshold = 1 * time.Minute
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, log, "smoke-ccip", nil, TestCfg)
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

	aggregateRateLimit := big.NewInt(1e16)

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s - Self Serve Rate Limit", tc.testName), func(t *testing.T) {
			tc.lane.Test = t
			src := tc.lane.Source
			dest := tc.lane.Dest
			require.GreaterOrEqual(t, len(src.Common.BridgeTokens), 2, "At least two bridge tokens needed for test")
			require.GreaterOrEqual(t, len(src.Common.BridgeTokenPools), 2, "At least two bridge token pools needed for test")
			require.GreaterOrEqual(t, len(dest.Common.BridgeTokens), 2, "At least two bridge tokens needed for test")
			require.GreaterOrEqual(t, len(dest.Common.BridgeTokenPools), 2, "At least two bridge token pools needed for test")
			// add liquidity to pools on both networks
			if !pointer.GetBool(TestCfg.TestGroupInput.ExistingDeployment) {
				addFund := func(ccipCommon *actions.CCIPCommon) {
					for i, btp := range ccipCommon.BridgeTokenPools {
						token := ccipCommon.BridgeTokens[i]
						err := btp.AddLiquidity(token.Approve, token.Address(), new(big.Int).Mul(aggregateRateLimit, big.NewInt(20)))
						require.NoError(t, err)
					}
				}
				addFund(src.Common)
				addFund(dest.Common)
			}

			var (
				freeTokenIndex    = 0
				limitedTokenIndex = 1

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
			overLimitAmount := new(big.Int).Add(aggregateRateLimit, big.NewInt(1))
			src.TransferAmount[freeTokenIndex] = overLimitAmount
			src.TransferAmount[limitedTokenIndex] = overLimitAmount
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(600_000))
			require.NoError(t, err)
			tc.lane.ValidateRequests()

			// Enable aggregate rate limiting on the destination chain for the limited token
			err = dest.AddRateLimitTokens([]*contracts.ERC20Token{limitedSrcToken}, []*contracts.ERC20Token{limitedDestToken})
			require.NoError(t, err, "Error setting destination rate limits")
			err = dest.OffRamp.SetRateLimit(contracts.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  aggregateRateLimit,
				Rate:      aggregateRateLimit,
			})
			require.NoError(t, err, "Error setting destination rate limits")
			err = dest.Common.ChainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")
			tc.lane.Logger.Debug().Str("Token", limitedSrcToken.ContractAddress.Hex()).Msg("Enabled aggregate rate limit on destination chain")

			// Send free token that should not have a rate limit and should succeed
			src.TransferAmount[freeTokenIndex] = overLimitAmount
			src.TransferAmount[limitedTokenIndex] = big.NewInt(0)
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(600_000))
			require.NoError(t, err, "Free token transfer failed")
			tc.lane.ValidateRequests()
			tc.lane.Logger.Info().Str("Token", freeSrcToken.ContractAddress.Hex()).Msg("Free token transfer succeeded")

			// Send limited token with rate limit that should fail on the destination chain
			src.TransferAmount[freeTokenIndex] = big.NewInt(0)
			src.TransferAmount[limitedTokenIndex] = overLimitAmount
			tc.lane.RecordStateBeforeTransfer()
			err = tc.lane.SendRequests(1, big.NewInt(600_000))
			require.NoError(t, err, "Failed to send rate limited token transfer")
			// Expect the ExecutionStateChanged event to never show up
			// Since we're looking to confirm that an event has NOT occurred, this can lead to some imperfect assumptions and results
			// We set the timeout to stop waiting for the event after a minute
			// 99% of transactions occur in under a minute in ideal simulated conditions, so this is an okay assumption there
			// but on real chains this risks false negatives
			// If we don't set this timeout, this test can take a long time and hold up CI
			tc.lane.ValidateRequests(actions.ExpectPhaseToFail(testreporters.ExecStateChanged, actions.WithTimeout(time.Minute)))
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

			// Change rate limit to make it viable
			err = dest.OffRamp.SetRateLimit(contracts.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  big.NewInt(0).Mul(aggregateRateLimit, big.NewInt(100)),
				Rate:      big.NewInt(0).Mul(aggregateRateLimit, big.NewInt(100)),
			})
			require.NoError(t, err, "Error setting destination rate limits")
			err = dest.Common.ChainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")
			tc.lane.Logger.Debug().Str("Token", limitedSrcToken.ContractAddress.Hex()).Msg("Enabled aggregate rate limit on destination chain")

			// Execute again manually and expect a pass
			err = tc.lane.ExecuteManually()
			require.NoError(t, err, "Error manually executing transaction after rate limit is lifted")
		})
	}
}

func TestSmokeCCIPMulticall(t *testing.T) {
	t.Parallel()
	log := logging.GetTestLogger(t)
	TestCfg := testsetups.NewCCIPTestConfig(t, log, testconfig.Smoke)
	// enable multicall in one tx for this test
	TestCfg.TestGroupInput.MulticallInOneTx = ptr.Ptr(true)
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, log, "smoke-ccip", nil, TestCfg)
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
	setUpOutput := testsetups.CCIPDefaultTestSetUp(t, log, "smoke-ccip", nil, TestCfg)
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
			tc.lane.ValidateRequests(actions.ExpectPhaseToFail(testreporters.ExecStateChanged, actions.ShouldExist()))
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
