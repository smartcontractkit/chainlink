package ccip_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_5_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	integrationtesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/integration"
)

func TestIntegration_CCIP(t *testing.T) {
	// Run the batches of tests for both pipeline and dynamic price getter setups.
	// We will remove the pipeline batch once the feature is deleted from the code.
	tests := []struct {
		name                     string
		withPipeline             bool
		allowOutOfOrderExecution bool
	}{
		{
			name:                     "with pipeline allowOutOfOrderExecution true",
			withPipeline:             true,
			allowOutOfOrderExecution: true,
		},
		{
			name:                     "with dynamic price getter allowOutOfOrderExecution false",
			withPipeline:             false,
			allowOutOfOrderExecution: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ccipTH := integrationtesthelpers.SetupCCIPIntegrationTH(
				t,
				testhelpers.SourceChainID,
				testhelpers.SourceChainSelector,
				testhelpers.DestChainID,
				testhelpers.DestChainSelector,
				ccip.DefaultSourceFinalityDepth,
				ccip.DefaultDestFinalityDepth,
			)

			tokenPricesUSDPipeline := ""
			priceGetterConfigJson := ""

			if test.withPipeline {
				// Set up a test pipeline.
				testPricePipeline, linkUSD, ethUSD := ccipTH.CreatePricesPipeline(t)
				defer linkUSD.Close()
				defer ethUSD.Close()
				tokenPricesUSDPipeline = testPricePipeline
			} else {
				// Set up a test price getter.
				// Set up the aggregators here to avoid modifying ccipTH.
				aggSrcNatAddr, _, aggSrcNat, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(ccipTH.Source.User, ccipTH.Source.Chain, 18, big.NewInt(2e18))
				require.NoError(t, err)
				_, err = aggSrcNat.UpdateRoundData(ccipTH.Source.User, big.NewInt(50), big.NewInt(17000000), big.NewInt(1000), big.NewInt(1000))
				require.NoError(t, err)
				ccipTH.Source.Chain.Commit()

				aggDstLnkAddr, _, aggDstLnk, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(ccipTH.Dest.User, ccipTH.Dest.Chain, 18, big.NewInt(3e18))
				require.NoError(t, err)
				ccipTH.Dest.Chain.Commit()
				_, err = aggDstLnk.UpdateRoundData(ccipTH.Dest.User, big.NewInt(50), big.NewInt(8000000), big.NewInt(1000), big.NewInt(1000))
				require.NoError(t, err)
				ccipTH.Dest.Chain.Commit()

				priceGetterConfig := config.DynamicPriceGetterConfig{
					AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
						ccipTH.Source.WrappedNative.Address(): {
							ChainID:                   ccipTH.Source.ChainID,
							AggregatorContractAddress: aggSrcNatAddr,
						},
						ccipTH.Dest.LinkToken.Address(): {
							ChainID:                   ccipTH.Dest.ChainID,
							AggregatorContractAddress: aggDstLnkAddr,
						},
						ccipTH.Dest.WrappedNative.Address(): {
							ChainID:                   ccipTH.Dest.ChainID,
							AggregatorContractAddress: aggDstLnkAddr,
						},
					},
					StaticPrices: map[common.Address]config.StaticPriceConfig{},
				}
				priceGetterConfigBytes, err := json.MarshalIndent(priceGetterConfig, "", " ")
				require.NoError(t, err)
				priceGetterConfigJson = string(priceGetterConfigBytes)
			}
			jobParams := ccipTH.SetUpNodesAndJobs(t, tokenPricesUSDPipeline, priceGetterConfigJson, "")

			// track sequence number and nonce separately since nonce doesn't bump for messages with allowOutOfOrderExecution == true,
			// but sequence number always bumps.
			// for this test, when test.outOfOrder == false, sequence number and nonce are equal.
			// when test.outOfOrder == true, nonce is not bumped at all, so sequence number and nonce are NOT equal.
			currentSeqNum := 1
			currentNonce := uint64(1)

			t.Run("single", func(t *testing.T) {
				tokenAmount := big.NewInt(500000003) // prime number
				gasLimit := big.NewInt(200_003)      // prime number

				extraArgs, err2 := testhelpers.GetEVMExtraArgsV2(gasLimit, test.allowOutOfOrderExecution)
				require.NoError(t, err2)

				sourceBalances, err2 := testhelpers.GetBalances(t, []testhelpers.BalanceReq{
					{Name: testhelpers.SourcePool, Addr: ccipTH.Source.LinkTokenPool.Address(), Getter: ccipTH.GetSourceLinkBalance},
					{Name: testhelpers.OnRamp, Addr: ccipTH.Source.OnRamp.Address(), Getter: ccipTH.GetSourceLinkBalance},
					{Name: testhelpers.SourceRouter, Addr: ccipTH.Source.Router.Address(), Getter: ccipTH.GetSourceLinkBalance},
					{Name: testhelpers.SourcePriceRegistry, Addr: ccipTH.Source.PriceRegistry.Address(), Getter: ccipTH.GetSourceLinkBalance},
				})
				require.NoError(t, err2)
				destBalances, err2 := testhelpers.GetBalances(t, []testhelpers.BalanceReq{
					{Name: testhelpers.Receiver, Addr: ccipTH.Dest.Receivers[0].Receiver.Address(), Getter: ccipTH.GetDestLinkBalance},
					{Name: testhelpers.DestPool, Addr: ccipTH.Dest.LinkTokenPool.Address(), Getter: ccipTH.GetDestLinkBalance},
					{Name: testhelpers.OffRamp, Addr: ccipTH.Dest.OffRamp.Address(), Getter: ccipTH.GetDestLinkBalance},
				})
				require.NoError(t, err2)

				ccipTH.Source.User.Value = tokenAmount
				_, err2 = ccipTH.Source.WrappedNative.Deposit(ccipTH.Source.User)
				require.NoError(t, err2)
				ccipTH.Source.Chain.Commit()
				ccipTH.Source.User.Value = nil

				msg := router.ClientEVM2AnyMessage{
					Receiver: testhelpers.MustEncodeAddress(t, ccipTH.Dest.Receivers[0].Receiver.Address()),
					Data:     []byte("hello"),
					TokenAmounts: []router.ClientEVMTokenAmount{
						{
							Token:  ccipTH.Source.LinkToken.Address(),
							Amount: tokenAmount,
						},
						{
							Token:  ccipTH.Source.WrappedNative.Address(),
							Amount: tokenAmount,
						},
					},
					FeeToken:  ccipTH.Source.LinkToken.Address(),
					ExtraArgs: extraArgs,
				}
				fee, err2 := ccipTH.Source.Router.GetFee(nil, testhelpers.DestChainSelector, msg)
				require.NoError(t, err2)
				// Currently no overhead and 10gwei dest gas price. So fee is simply (gasLimit * gasPrice)* link/native
				// require.Equal(t, new(big.Int).Mul(gasLimit, gasPrice).String(), fee.String())
				// Approve the fee amount + the token amount
				_, err2 = ccipTH.Source.LinkToken.Approve(ccipTH.Source.User, ccipTH.Source.Router.Address(), new(big.Int).Add(fee, tokenAmount))
				require.NoError(t, err2)
				ccipTH.Source.Chain.Commit()
				_, err2 = ccipTH.Source.WrappedNative.Approve(ccipTH.Source.User, ccipTH.Source.Router.Address(), tokenAmount)
				require.NoError(t, err2)
				ccipTH.Source.Chain.Commit()

				beforeNonce, err := ccipTH.Source.OnRamp.GetSenderNonce(nil, ccipTH.Source.User.From)
				require.NoError(t, err)
				ccipTH.SendRequest(t, msg)
				// TODO: can this be moved into SendRequest?
				if test.allowOutOfOrderExecution {
					// the nonce for that sender must not be bumped for allowOutOfOrderExecution == true messages.
					nonce, err2 := ccipTH.Source.OnRamp.GetSenderNonce(nil, ccipTH.Source.User.From)
					require.NoError(t, err2)
					require.Equal(t, beforeNonce, nonce, "nonce must not be bumped for allowOutOfOrderExecution == true requests")
				} else {
					// the nonce for that sender must be bumped for allowOutOfOrderExecution == false messages.
					nonce, err2 := ccipTH.Source.OnRamp.GetSenderNonce(nil, ccipTH.Source.User.From)
					require.NoError(t, err2)
					require.Equal(t, beforeNonce+1, nonce, "nonce must be bumped for allowOutOfOrderExecution == false requests")
				}

				// Should eventually see this executed.
				ccipTH.AllNodesHaveReqSeqNum(t, currentSeqNum)
				ccipTH.EventuallyReportCommitted(t, currentSeqNum)

				executionLogs := ccipTH.AllNodesHaveExecutedSeqNums(t, currentSeqNum, currentSeqNum)
				assert.Len(t, executionLogs, 1)
				ccipTH.AssertExecState(t, executionLogs[0], testhelpers.ExecutionStateSuccess)

				// Asserts
				// 1) The total pool input == total pool output
				// 2) Pool flow equals tokens sent
				// 3) Sent tokens arrive at the receiver
				ccipTH.AssertBalances(t, []testhelpers.BalanceAssertion{
					{
						Name:     testhelpers.SourcePool,
						Address:  ccipTH.Source.LinkTokenPool.Address(),
						Expected: testhelpers.MustAddBigInt(sourceBalances[testhelpers.SourcePool], tokenAmount.String()).String(),
						Getter:   ccipTH.GetSourceLinkBalance,
					},
					{
						Name:     testhelpers.SourcePriceRegistry,
						Address:  ccipTH.Source.PriceRegistry.Address(),
						Expected: sourceBalances[testhelpers.SourcePriceRegistry].String(),
						Getter:   ccipTH.GetSourceLinkBalance,
					},
					{
						// Fees end up in the onramp.
						Name:     testhelpers.OnRamp,
						Address:  ccipTH.Source.OnRamp.Address(),
						Expected: testhelpers.MustAddBigInt(sourceBalances[testhelpers.SourcePriceRegistry], fee.String()).String(),
						Getter:   ccipTH.GetSourceLinkBalance,
					},
					{
						Name:     testhelpers.SourceRouter,
						Address:  ccipTH.Source.Router.Address(),
						Expected: sourceBalances[testhelpers.SourceRouter].String(),
						Getter:   ccipTH.GetSourceLinkBalance,
					},
					{
						Name:     testhelpers.Receiver,
						Address:  ccipTH.Dest.Receivers[0].Receiver.Address(),
						Expected: testhelpers.MustAddBigInt(destBalances[testhelpers.Receiver], tokenAmount.String()).String(),
						Getter:   ccipTH.GetDestLinkBalance,
					},
					{
						Name:     testhelpers.DestPool,
						Address:  ccipTH.Dest.LinkTokenPool.Address(),
						Expected: testhelpers.MustSubBigInt(destBalances[testhelpers.DestPool], tokenAmount.String()).String(),
						Getter:   ccipTH.GetDestLinkBalance,
					},
					{
						Name:     testhelpers.OffRamp,
						Address:  ccipTH.Dest.OffRamp.Address(),
						Expected: destBalances[testhelpers.OffRamp].String(),
						Getter:   ccipTH.GetDestLinkBalance,
					},
				})
				currentSeqNum++
				if !test.allowOutOfOrderExecution {
					currentNonce = uint64(currentSeqNum)
				}
			})

			t.Run("multiple batches", func(t *testing.T) {
				tokenAmount := big.NewInt(500000003)
				gasLimit := big.NewInt(250_000)

				var txs []*gethtypes.Transaction
				// Enough to require batched executions as gasLimit per tx is 250k -> 500k -> 750k ....
				// The actual gas usage of executing 15 messages is higher than the gas limit for
				// a single tx. This means that when batching is turned off, and we simply include
				// all txs without checking gas, this also fails.
				n := 15
				for i := 0; i < n; i++ {
					txGasLimit := new(big.Int).Mul(gasLimit, big.NewInt(int64(i+1)))

					// interleave ordered and non-ordered messages.
					allowOutOfOrderExecution := false
					if i%2 == 0 {
						allowOutOfOrderExecution = true
					}
					extraArgs, err2 := testhelpers.GetEVMExtraArgsV2(txGasLimit, allowOutOfOrderExecution)
					require.NoError(t, err2)
					msg := router.ClientEVM2AnyMessage{
						Receiver: testhelpers.MustEncodeAddress(t, ccipTH.Dest.Receivers[0].Receiver.Address()),
						Data:     []byte("hello"),
						TokenAmounts: []router.ClientEVMTokenAmount{
							{
								Token:  ccipTH.Source.LinkToken.Address(),
								Amount: tokenAmount,
							},
						},
						FeeToken:  ccipTH.Source.LinkToken.Address(),
						ExtraArgs: extraArgs,
					}
					fee, err2 := ccipTH.Source.Router.GetFee(nil, testhelpers.DestChainSelector, msg)
					require.NoError(t, err2)
					// Currently no overhead and 1gwei dest gas price. So fee is simply gasLimit * gasPrice.
					// require.Equal(t, new(big.Int).Mul(txGasLimit, gasPrice).String(), fee.String())
					// Approve the fee amount + the token amount
					_, err2 = ccipTH.Source.LinkToken.Approve(ccipTH.Source.User, ccipTH.Source.Router.Address(), new(big.Int).Add(fee, tokenAmount))
					require.NoError(t, err2)
					tx, err2 := ccipTH.Source.Router.CcipSend(ccipTH.Source.User, ccipTH.Dest.ChainSelector, msg)
					require.NoError(t, err2)
					txs = append(txs, tx)
					if !allowOutOfOrderExecution {
						currentNonce++
					}
				}

				// Send a batch of requests in a single block
				testhelpers.ConfirmTxs(t, txs, ccipTH.Source.Chain)
				for i := 0; i < n; i++ {
					ccipTH.AllNodesHaveReqSeqNum(t, currentSeqNum+i)
				}
				// Should see a report with the full range
				ccipTH.EventuallyReportCommitted(t, currentSeqNum+n-1)
				// Should all be executed
				executionLogs := ccipTH.AllNodesHaveExecutedSeqNums(t, currentSeqNum, currentSeqNum+n-1)
				for _, execLog := range executionLogs {
					ccipTH.AssertExecState(t, execLog, testhelpers.ExecutionStateSuccess)
				}

				currentSeqNum += n
			})

			// Deploy new on ramp,Commit store,off ramp
			// Delete v1 jobs
			// Send a number of requests
			// Upgrade the router with new contracts
			// create new jobs
			// Verify all pending requests are sent after the contracts are upgraded
			t.Run("upgrade contracts and verify requests can be sent with upgraded contract", func(t *testing.T) {
				gasLimit := big.NewInt(200_003) // prime number
				tokenAmount := big.NewInt(100)
				commitStoreV1 := ccipTH.Dest.CommitStore
				offRampV1 := ccipTH.Dest.OffRamp
				onRampV1 := ccipTH.Source.OnRamp
				// deploy v2 contracts
				ccipTH.DeployNewOnRamp(t)
				ccipTH.DeployNewCommitStore(t)
				ccipTH.DeployNewOffRamp(t)

				// send a request as the v2 contracts are not enabled in router it should route through the v1 contracts
				t.Logf("sending request for seqnum %d", currentSeqNum)
				ccipTH.SendMessage(t, gasLimit, tokenAmount, ccipTH.Dest.Receivers[0].Receiver.Address())
				ccipTH.Source.Chain.Commit()
				ccipTH.Dest.Chain.Commit()
				t.Logf("verifying seqnum %d on previous onRamp %s", currentSeqNum, onRampV1.Address().Hex())
				ccipTH.AllNodesHaveReqSeqNum(t, currentSeqNum, onRampV1.Address())
				ccipTH.EventuallyReportCommitted(t, currentSeqNum, commitStoreV1.Address())
				executionLog := ccipTH.AllNodesHaveExecutedSeqNums(t, currentSeqNum, currentSeqNum, offRampV1.Address())
				ccipTH.AssertExecState(t, executionLog[0], testhelpers.ExecutionStateSuccess, offRampV1.Address())

				nonceAtOnRampV1, err := onRampV1.GetSenderNonce(nil, ccipTH.Source.User.From)
				require.NoError(t, err, "getting nonce from onRamp")
				require.Equal(t, currentNonce, nonceAtOnRampV1, "nonce should be synced from v1 onRamp")
				nonceAtOffRampV1, err := offRampV1.GetSenderNonce(nil, ccipTH.Source.User.From)
				require.NoError(t, err, "getting nonce from offRamp")
				require.Equal(t, currentNonce, nonceAtOffRampV1, "nonce should be synced from v1 offRamp")

				// enable the newly deployed contracts
				newConfigBlock := ccipTH.Dest.Chain.Blockchain().CurrentBlock().Number.Int64()
				ccipTH.EnableOnRamp(t)
				ccipTH.EnableCommitStore(t)
				ccipTH.EnableOffRamp(t)
				srcStartBlock := ccipTH.Source.Chain.Blockchain().CurrentBlock().Number.Uint64()

				// send a number of requests, the requests should not be delivered yet as the previous contracts are not configured
				// with the router anymore
				startSeq := 1
				noOfRequests := 5
				endSeqNum := startSeq + noOfRequests
				for i := startSeq; i <= endSeqNum; i++ {
					t.Logf("sending request for seqnum %d", i)
					ccipTH.SendMessage(t, gasLimit, tokenAmount, ccipTH.Dest.Receivers[0].Receiver.Address())
					ccipTH.Source.Chain.Commit()
					ccipTH.Dest.Chain.Commit()
					ccipTH.EventuallySendRequested(t, uint64(i))
				}

				// delete v1 jobs
				for _, node := range ccipTH.Nodes {
					id := node.FindJobIDForContract(t, commitStoreV1.Address())
					require.Greater(t, id, int32(0))
					t.Logf("deleting job %d", id)
					err = node.App.DeleteJob(context.Background(), id)
					require.NoError(t, err)
					id = node.FindJobIDForContract(t, offRampV1.Address())
					require.Greater(t, id, int32(0))
					t.Logf("deleting job %d", id)
					err = node.App.DeleteJob(context.Background(), id)
					require.NoError(t, err)
				}

				// Commit on both chains to reach Finality
				ccipTH.Source.Chain.Commit()
				ccipTH.Dest.Chain.Commit()

				// create new jobs
				jobParams = ccipTH.NewCCIPJobSpecParams(tokenPricesUSDPipeline, priceGetterConfigJson, newConfigBlock, "")
				jobParams.Version = "v2"
				jobParams.SourceStartBlock = srcStartBlock
				ccipTH.AddAllJobs(t, jobParams)
				committedSeqNum := uint64(0)
				// Now the requests should be delivered
				for i := startSeq; i <= endSeqNum; i++ {
					t.Logf("verifying seqnum %d", i)
					ccipTH.AllNodesHaveReqSeqNum(t, i)
					if committedSeqNum < uint64(i+1) {
						committedSeqNum = ccipTH.EventuallyReportCommitted(t, i)
					}
					ccipTH.EventuallyExecutionStateChangedToSuccess(t, []uint64{uint64(i)}, uint64(newConfigBlock))
				}

				// nonces should be correctly synced from v1 contracts for the sender
				nonceAtOnRampV2, err := ccipTH.Source.OnRamp.GetSenderNonce(nil, ccipTH.Source.User.From)
				require.NoError(t, err, "getting nonce from onRamp")
				nonceAtOffRampV2, err := ccipTH.Dest.OffRamp.GetSenderNonce(nil, ccipTH.Source.User.From)
				require.NoError(t, err, "getting nonce from offRamp")
				require.Equal(t, nonceAtOnRampV1+uint64(noOfRequests)+1, nonceAtOnRampV2, "nonce should be synced from v1 onRamps")
				require.Equal(t, nonceAtOffRampV1+uint64(noOfRequests)+1, nonceAtOffRampV2, "nonce should be synced from v1 offRamps")
				currentSeqNum = endSeqNum + 1
				if !test.allowOutOfOrderExecution {
					currentNonce = uint64(currentSeqNum)
				}
			})

			t.Run("pay nops", func(t *testing.T) {
				linkToTransferToOnRamp := big.NewInt(1e18)

				// transfer some link to onramp to pay the nops
				_, err := ccipTH.Source.LinkToken.Transfer(ccipTH.Source.User, ccipTH.Source.OnRamp.Address(), linkToTransferToOnRamp)
				require.NoError(t, err)
				ccipTH.Source.Chain.Commit()

				srcBalReq := []testhelpers.BalanceReq{
					{
						Name:   testhelpers.Sender,
						Addr:   ccipTH.Source.User.From,
						Getter: ccipTH.GetSourceWrappedTokenBalance,
					},
					{
						Name:   testhelpers.OnRampNative,
						Addr:   ccipTH.Source.OnRamp.Address(),
						Getter: ccipTH.GetSourceWrappedTokenBalance,
					},
					{
						Name:   testhelpers.OnRamp,
						Addr:   ccipTH.Source.OnRamp.Address(),
						Getter: ccipTH.GetSourceLinkBalance,
					},
					{
						Name:   testhelpers.SourceRouter,
						Addr:   ccipTH.Source.Router.Address(),
						Getter: ccipTH.GetSourceWrappedTokenBalance,
					},
				}

				var nopsAndWeights []evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight
				var totalWeight uint16
				nodes := ccipTH.Nodes
				for i := range nodes {
					// For now set the transmitter addresses to be the same as the payee addresses
					nodes[i].PaymentReceiver = nodes[i].Transmitter
					nopsAndWeights = append(nopsAndWeights, evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{
						Nop:    nodes[i].PaymentReceiver,
						Weight: 5,
					})
					totalWeight += 5
					srcBalReq = append(srcBalReq, testhelpers.BalanceReq{
						Name:   fmt.Sprintf("node %d", i),
						Addr:   nodes[i].PaymentReceiver,
						Getter: ccipTH.GetSourceLinkBalance,
					})
				}
				srcBalances, err := testhelpers.GetBalances(t, srcBalReq)
				require.NoError(t, err)

				// set nops on the onramp
				ccipTH.SetNopsOnRamp(t, nopsAndWeights)

				// send a message
				extraArgs, err := testhelpers.GetEVMExtraArgsV2(big.NewInt(200_000), test.allowOutOfOrderExecution)
				require.NoError(t, err)

				// FeeToken is empty, indicating it should use native token
				msg := router.ClientEVM2AnyMessage{
					Receiver:     testhelpers.MustEncodeAddress(t, ccipTH.Dest.Receivers[1].Receiver.Address()),
					Data:         []byte("hello"),
					TokenAmounts: []router.ClientEVMTokenAmount{},
					ExtraArgs:    extraArgs,
					FeeToken:     common.Address{},
				}
				fee, err := ccipTH.Source.Router.GetFee(nil, testhelpers.DestChainSelector, msg)
				require.NoError(t, err)

				// verify message is sent
				ccipTH.Source.User.Value = fee
				ccipTH.SendRequest(t, msg)
				ccipTH.Source.User.Value = nil
				ccipTH.AllNodesHaveReqSeqNum(t, currentSeqNum)
				ccipTH.EventuallyReportCommitted(t, currentSeqNum)

				executionLogs := ccipTH.AllNodesHaveExecutedSeqNums(t, currentSeqNum, currentSeqNum)
				assert.Len(t, executionLogs, 1)
				ccipTH.AssertExecState(t, executionLogs[0], testhelpers.ExecutionStateSuccess)
				currentSeqNum++
				if test.allowOutOfOrderExecution {
					currentNonce = uint64(currentSeqNum)
				}

				// get the nop fee
				nopFee, err := ccipTH.Source.OnRamp.GetNopFeesJuels(nil)
				require.NoError(t, err)
				t.Log("nopFee", nopFee)

				// withdraw fees and verify there is still fund left for nop payment
				_, err = ccipTH.Source.OnRamp.WithdrawNonLinkFees(
					ccipTH.Source.User,
					ccipTH.Source.WrappedNative.Address(),
					ccipTH.Source.User.From,
				)
				require.NoError(t, err)
				ccipTH.Source.Chain.Commit()

				// pay nops
				_, err = ccipTH.Source.OnRamp.PayNops(ccipTH.Source.User)
				require.NoError(t, err)
				ccipTH.Source.Chain.Commit()

				srcBalanceAssertions := []testhelpers.BalanceAssertion{
					{
						// Onramp should not have any balance left in wrapped native
						Name:     testhelpers.OnRampNative,
						Address:  ccipTH.Source.OnRamp.Address(),
						Expected: big.NewInt(0).String(),
						Getter:   ccipTH.GetSourceWrappedTokenBalance,
					},
					{
						// Onramp should have the remaining link after paying nops
						Name:     testhelpers.OnRamp,
						Address:  ccipTH.Source.OnRamp.Address(),
						Expected: new(big.Int).Sub(srcBalances[testhelpers.OnRamp], nopFee).String(),
						Getter:   ccipTH.GetSourceLinkBalance,
					},
					{
						Name:     testhelpers.SourceRouter,
						Address:  ccipTH.Source.Router.Address(),
						Expected: srcBalances[testhelpers.SourceRouter].String(),
						Getter:   ccipTH.GetSourceWrappedTokenBalance,
					},
					// onRamp's balance (of previously sent fee during message sending) should have been transferred to
					// the owner as a result of WithdrawNonLinkFees
					{
						Name:     testhelpers.Sender,
						Address:  ccipTH.Source.User.From,
						Expected: fee.String(),
						Getter:   ccipTH.GetSourceWrappedTokenBalance,
					},
				}

				// the nodes should be paid according to the weights assigned
				for i, node := range nodes {
					paymentWeight := float64(nopsAndWeights[i].Weight) / float64(totalWeight)
					paidInFloat := paymentWeight * float64(nopFee.Int64())
					paid, _ := new(big.Float).SetFloat64(paidInFloat).Int64()
					bal := new(big.Int).Add(
						new(big.Int).SetInt64(paid),
						srcBalances[fmt.Sprintf("node %d", i)]).String()
					srcBalanceAssertions = append(srcBalanceAssertions, testhelpers.BalanceAssertion{
						Name:     fmt.Sprintf("node %d", i),
						Address:  node.PaymentReceiver,
						Expected: bal,
						Getter:   ccipTH.GetSourceLinkBalance,
					})
				}
				ccipTH.AssertBalances(t, srcBalanceAssertions)
			})

			// Keep on sending a bunch of messages
			// In the meantime update onchainConfig with new price registry address
			// Verify if the jobs can pick up updated config
			// Verify if all the messages are sent
			t.Run("config change or price registry update while requests are inflight", func(t *testing.T) {
				gasLimit := big.NewInt(200_003) // prime number
				tokenAmount := big.NewInt(100)
				msgWg := &sync.WaitGroup{}
				msgWg.Add(1)
				ticker := time.NewTicker(100 * time.Millisecond)
				defer ticker.Stop()
				startSeq := currentSeqNum
				endSeq := currentSeqNum + 20

				// send message with the old configs
				ccipTH.SendMessage(t, gasLimit, tokenAmount, ccipTH.Dest.Receivers[0].Receiver.Address())
				ccipTH.Source.Chain.Commit()

				go func(ccipContracts testhelpers.CCIPContracts, currentSeqNum int) {
					seqNumber := currentSeqNum + 1
					defer msgWg.Done()
					for {
						<-ticker.C // wait for ticker
						t.Logf("sending request for seqnum %d", seqNumber)
						ccipContracts.SendMessage(t, gasLimit, tokenAmount, ccipTH.Dest.Receivers[0].Receiver.Address())
						ccipContracts.Source.Chain.Commit()
						seqNumber++
						if seqNumber == endSeq {
							return
						}
					}
				}(ccipTH.CCIPContracts, currentSeqNum)

				ccipTH.DeployNewPriceRegistry(t)
				commitOnchainConfig := ccipTH.CreateDefaultCommitOnchainConfig(t)
				commitOffchainConfig := ccipTH.CreateDefaultCommitOffchainConfig(t)
				execOnchainConfig := ccipTH.CreateDefaultExecOnchainConfig(t)
				execOffchainConfig := ccipTH.CreateDefaultExecOffchainConfig(t)

				ccipTH.SetupOnchainConfig(t, commitOnchainConfig, commitOffchainConfig, execOnchainConfig, execOffchainConfig)

				// wait for all requests to be complete
				msgWg.Wait()
				for i := startSeq; i < endSeq; i++ {
					ccipTH.AllNodesHaveReqSeqNum(t, i)
					ccipTH.EventuallyReportCommitted(t, i)

					executionLogs := ccipTH.AllNodesHaveExecutedSeqNums(t, i, i)
					assert.Len(t, executionLogs, 1)
					ccipTH.AssertExecState(t, executionLogs[0], testhelpers.ExecutionStateSuccess)
				}

				for i, node := range ccipTH.Nodes {
					t.Logf("verifying node %d", i)
					node.EventuallyNodeUsesNewCommitConfig(t, ccipTH, ccipdata.CommitOnchainConfig{
						PriceRegistry: ccipTH.Dest.PriceRegistry.Address(),
					})
					node.EventuallyNodeUsesNewExecConfig(t, ccipTH, v1_5_0.ExecOnchainConfig{
						PermissionLessExecutionThresholdSeconds: testhelpers.PermissionLessExecutionThresholdSeconds,
						Router:                                  ccipTH.Dest.Router.Address(),
						PriceRegistry:                           ccipTH.Dest.PriceRegistry.Address(),
						MaxDataBytes:                            1e5,
						MaxNumberOfTokensPerMsg:                 5,
					})
					node.EventuallyNodeUsesUpdatedPriceRegistry(t, ccipTH)
				}
				currentSeqNum = endSeq
				if test.allowOutOfOrderExecution {
					currentNonce = uint64(currentSeqNum)
				}
			})
		})
	}
}

// TestReorg ensures that CCIP works even when a below finality depth reorg happens
func TestReorg(t *testing.T) {
	// We need higher finality depth on the destination to perform reorg deep enough to revert commit and execution reports
	destinationFinalityDepth := uint32(50)
	ccipTH := integrationtesthelpers.SetupCCIPIntegrationTH(
		t,
		testhelpers.SourceChainID,
		testhelpers.SourceChainSelector,
		testhelpers.DestChainID,
		testhelpers.DestChainSelector,
		ccip.DefaultSourceFinalityDepth,
		destinationFinalityDepth,
	)
	testPricePipeline, linkUSD, ethUSD := ccipTH.CreatePricesPipeline(t)
	defer linkUSD.Close()
	defer ethUSD.Close()
	ccipTH.SetUpNodesAndJobs(t, testPricePipeline, "", "")

	gasLimit := big.NewInt(200_00)
	tokenAmount := big.NewInt(1)

	forkBlock, err := ccipTH.Dest.Chain.BlockByNumber(context.Background(), nil)
	require.NoError(t, err, "Error while fetching the destination chain current block number")

	// Adjust time to start next blocks with timestamps two hours after the fork block.
	// This is critical to have two forks with different block_timestamps.
	err = ccipTH.Dest.Chain.AdjustTime(2 * time.Hour)
	require.NoError(t, err, "Error while adjusting the destination chain time")
	ccipTH.Dest.Chain.Commit()

	// Send request for the first time and make sure it's executed on the destination
	ccipTH.SendMessage(t, gasLimit, tokenAmount, ccipTH.Dest.Receivers[0].Receiver.Address())
	ccipTH.Dest.User.GasLimit = 100000
	ccipTH.EventuallySendRequested(t, uint64(1))
	ccipTH.EventuallyReportCommitted(t, 1)
	executionLog := ccipTH.AllNodesHaveExecutedSeqNums(t, 1, 1)
	ccipTH.AssertExecState(t, executionLog[0], testhelpers.ExecutionStateSuccess)

	currentBlock, err := ccipTH.Dest.Chain.BlockByNumber(context.Background(), nil)
	require.NoError(t, err, "Error while fetching the current block number of destination chain")

	// Reorg back to the `forkBlock`. Next blocks in the fork will have block_timestamps right after the fork,
	// but before the 2 hours interval defined above for the canonical chain
	require.NoError(t, ccipTH.Dest.Chain.Fork(testutils.Context(t), forkBlock.Hash()),
		"Error while forking the chain")
	// Make sure that fork is longer than the canonical chain to enforce switch
	noOfBlocks := int(currentBlock.NumberU64() - forkBlock.NumberU64())
	for i := 0; i < noOfBlocks+1; i++ {
		ccipTH.Dest.Chain.Commit()
	}

	// State of the chain (block_timestamps) after reorg:
	//            / --> block1 (02:01) --> block2 (02:02) --> commit report (02:03) --> ...
	// forkBlock (00:00) --> block1' (00:01) --> block2' (00:02) --> commit report' (00:03) --> ...

	// CCIP should commit and executed messages that was reorged away
	ccipTH.EventuallyReportCommitted(t, 1)
	executionLog = ccipTH.AllNodesHaveExecutedSeqNums(t, 1, 1)
	ccipTH.AssertExecState(t, executionLog[0], testhelpers.ExecutionStateSuccess)

	// Sending another message and make sure it's executed on the destination
	ccipTH.SendMessage(t, gasLimit, tokenAmount, ccipTH.Dest.Receivers[0].Receiver.Address())
	ccipTH.EventuallySendRequested(t, uint64(2))
	ccipTH.EventuallyReportCommitted(t, 2)
	executionLog = ccipTH.AllNodesHaveExecutedSeqNums(t, 1, 2)
	ccipTH.AssertExecState(t, executionLog[0], testhelpers.ExecutionStateSuccess)
}
