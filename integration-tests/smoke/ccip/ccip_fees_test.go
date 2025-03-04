package ccip

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/feestest"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

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
				Sender:         e.Chains[sourceChain].DeployerKey,
				IsTestRouter:   true,
				SourceChain:    sourceChain,
				DestChain:      destChain,
				Evm2AnyMessage: ccipMessage,
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
}
