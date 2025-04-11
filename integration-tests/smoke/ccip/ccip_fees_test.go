package ccip

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/feestest"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func setupNewFeeToken(
	t *testing.T,
	tenv testhelpers.DeployedEnv,
	deployer *bind.TransactOpts,
	chainSelector uint64,
	state changeset.CCIPOnChainState,
	tokenSymbol string,
	tokenDecimals uint8,
) (feeToken *burn_mint_erc677.BurnMintERC677) {
	lggr := logger.TestLogger(t)
	chain := tenv.Env.Chains[chainSelector]
	tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
		deployer,
		chain.Client,
		tokenSymbol,
		tokenSymbol,
		tokenDecimals,
		big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
	)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)
	lggr.Infow("Deployed new fee token", "tokenAddress", tokenAddress)

	// grant mint role
	tx, err = token.GrantMintRole(deployer, deployer.From)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	// mint token and approve to router
	tx, err = token.Mint(deployer, deployer.From, deployment.E18Mult(10_000))
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = token.Approve(deployer, state.Chains[chainSelector].Router.Address(), math.MaxBig256)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	// add this new fee token to fee quoter
	tx, err = state.Chains[chain.Selector].FeeQuoter.ApplyFeeTokensUpdates(deployer, []common.Address{}, []common.Address{tokenAddress})
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)
	lggr.Infow("Added new fee token to fee quoter", "tokenAddress", tokenAddress)

	// set price for this new token
	tx, err = state.Chains[chain.Selector].FeeQuoter.UpdatePrices(
		deployer,
		fee_quoter.InternalPriceUpdates{
			TokenPriceUpdates: []fee_quoter.InternalTokenPriceUpdate{
				{
					SourceToken: tokenAddress,
					UsdPerToken: big.NewInt(1e18),
				},
			},
			GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{},
		},
	)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)
	lggr.Infow("Set price for new fee token", "tokenAddress", tokenAddress)

	return token
}

// setupTokens deploys transferable tokens on the source and dest, mints tokens for the source and dest, and
// approves the router to spend the tokens
func setupTokens(
	t *testing.T,
	state changeset.CCIPOnChainState,
	tenv testhelpers.DeployedEnv,
	src, dest uint64,
	transferTokenMintAmount,
	feeTokenMintAmount *big.Int,
) (
	srcToken *burn_mint_erc677.BurnMintERC677,
	dstToken *burn_mint_erc677.BurnMintERC677,
) {
	lggr := logger.TestLogger(t)
	e := tenv.Env

	// Deploy the token to test transferring
	srcToken, _, dstToken, _, err := testhelpers.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		src,
		dest,
		tenv.Env.Chains[src].DeployerKey,
		tenv.Env.Chains[dest].DeployerKey,
		state,
		tenv.Env.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	linkToken := state.Chains[src].LinkToken

	tx, err := srcToken.Mint(
		e.Chains[src].DeployerKey,
		e.Chains[src].DeployerKey.From,
		transferTokenMintAmount,
	)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	// Mint a destination token
	tx, err = dstToken.Mint(
		e.Chains[dest].DeployerKey,
		e.Chains[dest].DeployerKey.From,
		transferTokenMintAmount,
	)
	_, err = deployment.ConfirmIfNoError(e.Chains[dest], tx, err)
	require.NoError(t, err)

	// Approve the router to spend the tokens and confirm the tx's
	// To prevent having to approve the router for every transfer, we approve a sufficiently large amount
	tx, err = srcToken.Approve(e.Chains[src].DeployerKey, state.Chains[src].Router.Address(), math.MaxBig256)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	tx, err = dstToken.Approve(e.Chains[dest].DeployerKey, state.Chains[dest].Router.Address(), math.MaxBig256)
	_, err = deployment.ConfirmIfNoError(e.Chains[dest], tx, err)
	require.NoError(t, err)

	// Grant mint and burn roles to the deployer key for the newly deployed linkToken
	// Since those roles are not granted automatically
	tx, err = linkToken.GrantMintAndBurnRoles(e.Chains[src].DeployerKey, e.Chains[src].DeployerKey.From)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	// Mint link token and confirm the tx
	tx, err = linkToken.Mint(
		e.Chains[src].DeployerKey,
		e.Chains[src].DeployerKey.From,
		feeTokenMintAmount,
	)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	return srcToken, dstToken
}

func Test_CCIPFees(t *testing.T) {
	t.Parallel()
	tenv, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithMultiCall3(),
	)
	e := tenv.Env

	allChains := e.AllChainSelectors()
	require.Len(t, allChains, 2, "need two chains for this test")
	sourceChain := allChains[0]
	destChain := allChains[1]

	// Get new state after migration.
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	srcToken, dstToken := setupTokens(
		t,
		state,
		tenv,
		sourceChain,
		destChain,
		deployment.E18Mult(10_000),
		deployment.E18Mult(10_000),
	)

	// Ensure capreg logs are up to date.
	testhelpers.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Add all lanes
	testhelpers.AddLanesForAll(t, &tenv, state)

	t.Run("Send programmable token transfer pay with Link token", func(t *testing.T) {
		feestest.RunFeeTokenTestCase(feestest.NewFeeTokenTestCase(
			t,
			e,
			sourceChain,
			destChain,
			state.Chains[sourceChain].LinkToken.Address(), // feeToken
			[]router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: deployment.E18Mult(2),
				},
			},
			srcToken,
			dstToken,
			common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32), // receiver
			[]byte("hello ptt world"), // data
			true,                      // assertTokenBalance
			true,                      // assertExecution
		))
	})

	t.Run("Send programmable token transfer pay with native", func(t *testing.T) {
		feestest.RunFeeTokenTestCase(feestest.NewFeeTokenTestCase(
			t,
			e,
			// note the order of src and dest is reversed here
			destChain,
			sourceChain,
			common.HexToAddress("0x0"), // feeToken
			[]router.ClientEVMTokenAmount{
				{
					Token:  dstToken.Address(),
					Amount: deployment.E18Mult(2),
				},
			},
			// note the order of src and dest is reversed here
			dstToken,
			srcToken,
			common.LeftPadBytes(state.Chains[sourceChain].Receiver.Address().Bytes(), 32), // receiver
			[]byte("hello ptt world"), // data
			true,                      // assertTokenBalance
			true,                      // assertExecution
		))
	})

	t.Run("Send programmable toke transfer pay with custom fee token", func(t *testing.T) {
		feeToken := setupNewFeeToken(t, tenv, e.Chains[sourceChain].DeployerKey, sourceChain, state, "FEE", 18)
		feestest.RunFeeTokenTestCase(feestest.NewFeeTokenTestCase(
			t,
			e,
			sourceChain,
			destChain,
			feeToken.Address(),
			[]router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: deployment.E18Mult(2),
				},
			},
			srcToken,
			dstToken,
			common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32), // receiver
			[]byte("hello ptt world"), // data
			true,                      // assertTokenBalance
			true,                      // assertExecution
		))
	})

	t.Run("Send programmable token transfer pay with wrapped native", func(t *testing.T) {
		feestest.RunFeeTokenTestCase(feestest.NewFeeTokenTestCase(
			t,
			e,
			sourceChain,
			destChain,
			state.Chains[sourceChain].Weth9.Address(), // feeToken
			[]router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: deployment.E18Mult(2),
				},
			},
			srcToken,
			dstToken,
			common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32), // receiver
			[]byte("hello ptt world"), // data
			true,                      // assertTokenBalance
			true,                      // assertExecution
		))
	})

	t.Run("Send programmable token transfer but revert not enough tokens", func(t *testing.T) {
		// Send to the receiver on the destination chain paying with LINK token
		var (
			receiver = common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32)
			data     = []byte("")
			feeToken = state.Chains[sourceChain].LinkToken.Address()
		)

		// Increase the token send amount to more than available to intentionally cause a revert
		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver: receiver,
			Data:     data,
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: deployment.E18Mult(100_000_000),
				},
			},
			FeeToken:  feeToken,
			ExtraArgs: nil,
		}

		_, _, err = testhelpers.CCIPSendRequest(
			e,
			state,
			&testhelpers.CCIPSendReqConfig{
				Sender:       e.Chains[sourceChain].DeployerKey,
				IsTestRouter: true,
				SourceChain:  sourceChain,
				DestChain:    destChain,
				Message:      ccipMessage,
			},
		)
		require.Error(t, err)
	})

	t.Run("Send data-only message pay with link token", func(t *testing.T) {
		feestest.RunFeeTokenTestCase(feestest.NewFeeTokenTestCase(
			t,
			e,
			sourceChain,
			destChain,
			// no tokens, only data
			state.Chains[sourceChain].LinkToken.Address(), // feeToken
			nil, // tokenAmounts
			srcToken,
			dstToken,
			common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32), // receiver
			[]byte("hello link world"), // data
			false,                      // assertTokenBalance
			true,                       // assertExecution
		))
	})

	t.Run("Send message pay with native", func(t *testing.T) {
		feestest.RunFeeTokenTestCase(feestest.NewFeeTokenTestCase(
			t,
			e,
			sourceChain,
			destChain,
			common.HexToAddress("0x0"), // feeToken
			// no tokens, only data
			nil, // tokenAmounts
			srcToken,
			dstToken,
			common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32), // receiver
			[]byte("hello native world"), // data
			false,                        // assertTokenBalance
			true,                         // assertExecution
		))
	})

	t.Run("Send message pay with wrapped native", func(t *testing.T) {
		feestest.RunFeeTokenTestCase(feestest.NewFeeTokenTestCase(
			t,
			e,
			sourceChain,
			destChain,
			state.Chains[sourceChain].Weth9.Address(), // feeToken
			// no tokens, only data
			nil, // tokenAmounts
			srcToken,
			dstToken,
			common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32), // receiver
			[]byte("hello wrapped native world"),                                        // data
			false,                                                                       // assertTokenBalance
			true,                                                                        // assertExecution
		))
	})

	t.Run("Send message pay with custom fee token", func(t *testing.T) {
		feeToken := setupNewFeeToken(t, tenv, e.Chains[sourceChain].DeployerKey, sourceChain, state, "FEE", 18)
		feestest.RunFeeTokenTestCase(feestest.NewFeeTokenTestCase(
			t,
			e,
			sourceChain,
			destChain,
			feeToken.Address(),
			// no tokens, only data
			nil, // tokenAmounts
			srcToken,
			dstToken,
			common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32), // receiver
			[]byte("hello custom fee token world"),                                      // data
			false,                                                                       // assertTokenBalance
			true,                                                                        // assertExecution
		))
	})
}
