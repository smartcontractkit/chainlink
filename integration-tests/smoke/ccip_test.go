package smoke

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInitialDeployOnLocal(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)
	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello world"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			expectedSeqNum[changeset.SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = msgSentEvent.SequenceNumber
		}
	}

	// Wait for all commit reports to land.
	changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	changeset.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	// TODO: Apply the proposal.
}

func TestTokenTransfer(t *testing.T) {
	t.Parallel()
	tenv, state := setupEnvironment(t)
	e := tenv.Env
<<<<<<< HEAD
	output, err := changeset.DeployPrerequisites(tenv.Env, changeset.DeployPrerequisiteConfig{
		ChainSelectors: e.AllChainSelectors(),
	})

	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))

	// Deploy CCIP contracts
	output, err = changeset.InitialDeploy(e, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: e.AllChainSelectors(),
		TokenConfig:    ccdeploy.NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds),
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))
	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			_, err := e.Offchain.ProposeJob(ccdeploy.Context(t),
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes.
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))

	// Deploy and approve tokens.
	srcToken1, dstToken1 := deployAndApproveTokens(t, tenv, state, "Token1")
	srcToken2, dstToken2 := deployAndApproveTokens(t, tenv, state, "Token2")

	// Test scenarios are defined here
	scenarios := []struct {
		name         string
		srcChain     uint64
		dstChain     uint64
		tokenAmounts []router.ClientEVMTokenAmount
		receiver     common.Address
		data         []byte
	}{
		{
			name:     "Send token to EOA",
			srcChain: tenv.HomeChainSel,
			dstChain: tenv.FeedChainSel,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken1.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			receiver: e.Chains[tenv.FeedChainSel].DeployerKey.From,
			data:     []byte(""),
		},
		{
			name:     "Send token to contract",
			srcChain: tenv.HomeChainSel,
			dstChain: tenv.FeedChainSel,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken1.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			receiver: state.Chains[tenv.FeedChainSel].Receiver.Address(),
			data:     []byte(""),
		},
		{
			name:     "Send 2 tokens to receiver",
			srcChain: tenv.HomeChainSel,
			dstChain: tenv.FeedChainSel,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken1.Address(),
					Amount: big.NewInt(1e18),
				},
				{
					Token:  srcToken2.Address(),
					Amount: big.NewInt(2e18),
				},
			},
			receiver: e.Chains[tenv.FeedChainSel].DeployerKey.From,
			data:     []byte(""),
		},
		{
			name:     "Send N tokens to contract",
			srcChain: tenv.HomeChainSel,
			dstChain: tenv.FeedChainSel,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken1.Address(),
					Amount: big.NewInt(1e18),
				},
				{
					Token:  srcToken2.Address(),
					Amount: big.NewInt(2e18),
				},
				{
					Token:  srcToken1.Address(),
					Amount: big.NewInt(3e18),
				},
			},
			receiver: state.Chains[tenv.FeedChainSel].Receiver.Address(),
			data:     []byte(""),
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
=======
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	srcToken, _, dstToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		tenv.HomeChainSel,
		tenv.FeedChainSel,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)

	twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))
	tx, err := srcToken.Mint(
		e.Chains[tenv.HomeChainSel].DeployerKey,
		e.Chains[tenv.HomeChainSel].DeployerKey.From,
		new(big.Int).Mul(twoCoins, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = dstToken.Mint(
		e.Chains[tenv.FeedChainSel].DeployerKey,
		e.Chains[tenv.FeedChainSel].DeployerKey.From,
		new(big.Int).Mul(twoCoins, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = srcToken.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), twoCoins)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)
	tx, err = dstToken.Approve(e.Chains[tenv.FeedChainSel].DeployerKey, state.Chains[tenv.FeedChainSel].Router.Address(), twoCoins)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	tokens := map[uint64][]router.ClientEVMTokenAmount{
		tenv.HomeChainSel: {{
			Token:  srcToken.Address(),
			Amount: twoCoins,
		}},
		tenv.FeedChainSel: {{
			Token:  dstToken.Address(),
			Amount: twoCoins,
		}},
	}

	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block

			var (
				receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
				data     = []byte("hello world")
				feeToken = common.HexToAddress("0x0")
			)
			if src == tenv.HomeChainSel && dest == tenv.FeedChainSel {
				msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
					Receiver:     receiver,
					Data:         data,
					TokenAmounts: tokens[src],
					FeeToken:     feeToken,
					ExtraArgs:    nil,
				})
				expectedSeqNum[changeset.SourceDestPair{
					SourceChainSelector: src,
					DestChainSelector:   dest,
				}] = msgSentEvent.SequenceNumber
			} else {
				msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
					Receiver:     receiver,
					Data:         data,
					TokenAmounts: nil,
					FeeToken:     feeToken,
					ExtraArgs:    nil,
				})
				expectedSeqNum[changeset.SourceDestPair{
					SourceChainSelector: src,
					DestChainSelector:   dest,
				}] = msgSentEvent.SequenceNumber
			}
		}
	}

	// Wait for all commit reports to land.
	changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	changeset.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)
>>>>>>> develop

			startBlocks := make(map[uint64]*uint64)
			expectedSeqNum := make(map[uint64]uint64)

			destChain := e.Chains[scenario.dstChain]
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[scenario.dstChain] = &block

			// Fetch initial balances and aggregate total amounts
			initialBalances := make(map[common.Address]*big.Int)
			totalAmountsTransferred := make(map[common.Address]*big.Int)

			for _, tokenAmount := range scenario.tokenAmounts {
				dstToken := getDestinationToken(scenario, tenv, srcToken1, dstToken1, srcToken2, dstToken2, tokenAmount.Token)
				require.NotNil(t, dstToken, "Destination token not found")

				// Fetch initial balance
				if _, exists := initialBalances[dstToken.Address()]; !exists {
					balance, err := dstToken.BalanceOf(nil, scenario.receiver)
					require.NoError(t, err)
					initialBalances[dstToken.Address()] = balance
				}

				// Initialize and sum up total amounts
				if _, exists := totalAmountsTransferred[dstToken.Address()]; !exists {
					totalAmountsTransferred[dstToken.Address()] = big.NewInt(0)
				}
				totalAmountsTransferred[dstToken.Address()].Add(totalAmountsTransferred[dstToken.Address()], tokenAmount.Amount)
			}

			msg := router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(scenario.receiver.Bytes(), 32),
				Data:         scenario.data,
				TokenAmounts: scenario.tokenAmounts,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			}

			// Send the message
			msgSentEvent := ccdeploy.TestSendRequest(t, e, state, scenario.srcChain, scenario.dstChain, false, msg)
			expectedSeqNum[scenario.dstChain] = msgSentEvent.SequenceNumber

			// Wait for commit and execution
			ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)
			ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

			// Fetch final balances and assert
			for tokenAddress, totalAmount := range totalAmountsTransferred {
				dstToken := getTokenByAddress(dstToken1, dstToken2, tokenAddress)
				require.NotNil(t, dstToken, "Destination token not found for address %s", tokenAddress.Hex())

				finalBalance, err := dstToken.BalanceOf(nil, scenario.receiver)
				require.NoError(t, err)

				initialBalance := initialBalances[dstToken.Address()]
				expectedBalance := new(big.Int).Add(initialBalance, totalAmount)

				require.Equal(t, expectedBalance, finalBalance, "Incorrect balance for token %s", dstToken.Address().Hex())
			}
		})
	}
}

func getTokenByAddress(dstToken1, dstToken2 *burn_mint_erc677.BurnMintERC677, tokenAddress common.Address) *burn_mint_erc677.BurnMintERC677 {
	if dstToken1.Address() == tokenAddress {
		return dstToken1
	} else if dstToken2.Address() == tokenAddress {
		return dstToken2
	}
	return nil
}

// Helper function to determine the destination token
func getDestinationToken(scenario struct {
	name         string
	srcChain     uint64
	dstChain     uint64
	tokenAmounts []router.ClientEVMTokenAmount
	receiver     common.Address
	data         []byte
}, tenv ccdeploy.DeployedEnv, srcToken1, dstToken1, srcToken2, dstToken2 *burn_mint_erc677.BurnMintERC677, tokenAddress common.Address) *burn_mint_erc677.BurnMintERC677 {
	if scenario.srcChain == tenv.HomeChainSel && scenario.dstChain == tenv.FeedChainSel {
		if tokenAddress == srcToken1.Address() {
			return dstToken1
		} else if tokenAddress == srcToken2.Address() {
			return dstToken2
		}
	} else if scenario.srcChain == tenv.FeedChainSel && scenario.dstChain == tenv.HomeChainSel {
		if tokenAddress == dstToken1.Address() {
			return srcToken1
		} else if tokenAddress == dstToken2.Address() {
			return srcToken2
		}
	}
	return nil
}

func setupEnvironment(t *testing.T) (ccdeploy.DeployedEnv, ccdeploy.CCIPOnChainState) {
	lggr := logger.TestLogger(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)
	state, err := ccdeploy.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	return tenv, state
}

func deployAndApproveTokens(t *testing.T, e ccdeploy.DeployedEnv, state ccdeploy.CCIPOnChainState, tokenName string) (*burn_mint_erc677.BurnMintERC677, *burn_mint_erc677.BurnMintERC677) {
	srcToken, _, dstToken, _, err := ccdeploy.DeployTransferableToken(
		logger.TestLogger(t),
		e.Env.Chains,
		e.HomeChainSel,
		e.FeedChainSel,
		state,
		e.Env.ExistingAddresses,
		tokenName,
	)
	require.NoError(t, err)

	tenTokens := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(10))
	mintAndApprove(t, e, state, srcToken, e.HomeChainSel, tenTokens)
	mintAndApprove(t, e, state, dstToken, e.FeedChainSel, tenTokens)

	return srcToken, dstToken
}

func mintAndApprove(t *testing.T, e ccdeploy.DeployedEnv, state ccdeploy.CCIPOnChainState, token *burn_mint_erc677.BurnMintERC677, chainSel uint64, amount *big.Int) {
	tx, err := token.Mint(
		e.Env.Chains[chainSel].DeployerKey,
		e.Env.Chains[chainSel].DeployerKey.From,
		amount,
	)
	require.NoError(t, err)
	_, err = e.Env.Chains[chainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = token.Approve(
		e.Env.Chains[chainSel].DeployerKey,
		state.Chains[chainSel].Router.Address(),
		amount,
	)
	require.NoError(t, err)
	_, err = e.Env.Chains[chainSel].Confirm(tx)
	require.NoError(t, err)
}
