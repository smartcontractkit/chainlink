package smoke

import (
	"context"
	"math/big"
	"testing"

	"golang.org/x/exp/maps"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	sel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := tests.Context(t)
	config := &changeset.TestConfigs{}
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, config)
	inMemoryEnv := false

	// use this if you are testing locally in memory
	// tenv := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, lggr, 2, 4, config)
	// inMemoryEnv := true

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Chains)
	sourceChain, destChain := allChainSelectors[0], allChainSelectors[1]
	ownerSourceChain := e.Chains[sourceChain].DeployerKey
	ownerDestChain := e.Chains[destChain].DeployerKey

	oneE18 := new(big.Int).SetUint64(1e18)
	funds := new(big.Int).Mul(oneE18, new(big.Int).SetUint64(10))

	// Deploy and fund self-serve actors
	selfServeSrcTokenPoolDeployer := createAndFundSelfServeActor(ctx, t, ownerSourceChain, e.Chains[sourceChain], funds, inMemoryEnv)
	selfServeDestTokenPoolDeployer := createAndFundSelfServeActor(ctx, t, ownerDestChain, e.Chains[destChain], funds, inMemoryEnv)

	// Deploy tokens and pool by CCIP Owner
	srcToken, _, destToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain,
		destChain,
		ownerSourceChain,
		ownerDestChain,
		state,
		e.ExistingAddresses,
		"OWNER_TOKEN",
	)
	require.NoError(t, err)

	// Deploy Self Serve tokens and pool
	selfServeSrcToken, _, selfServeDestToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain,
		destChain,
		selfServeSrcTokenPoolDeployer,
		selfServeDestTokenPoolDeployer,
		state,
		e.ExistingAddresses,
		"SELF_SERVE_TOKEN",
	)
	require.NoError(t, err)
	require.NoError(t, changeset.AddLanesForAll(e, state))

	changeset.MintAndAllow(
		t,
		e,
		state,
		map[uint64][]changeset.MintTokenInfo{
			sourceChain: {
				changeset.NewMintTokenInfo(selfServeSrcTokenPoolDeployer, selfServeSrcToken),
				changeset.NewMintTokenInfo(ownerSourceChain, srcToken),
			},
			destChain: {
				changeset.NewMintTokenInfo(selfServeDestTokenPoolDeployer, selfServeDestToken),
				changeset.NewMintTokenInfo(ownerDestChain, destToken),
			},
		},
	)

	tcs := []struct {
		name                   string
		srcChain               uint64
		dstChain               uint64
		tokenAmounts           []router.ClientEVMTokenAmount
		receiver               common.Address
		data                   []byte
		extraData              []byte
		expectedTokenBalances  map[common.Address]*big.Int
		expectedExecutionState int
	}{
		{
			name:     "Send token to EOA",
			srcChain: sourceChain,
			dstChain: destChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: oneE18,
				},
			},
			receiver: utils.RandomAddress(),
			expectedTokenBalances: map[common.Address]*big.Int{
				destToken.Address(): oneE18,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:     "Send token to contract",
			srcChain: sourceChain,
			dstChain: destChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: oneE18,
				},
			},
			receiver: state.Chains[destChain].Receiver.Address(),
			expectedTokenBalances: map[common.Address]*big.Int{
				destToken.Address(): oneE18,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:     "Send N tokens to contract",
			srcChain: destChain,
			dstChain: sourceChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  selfServeDestToken.Address(),
					Amount: oneE18,
				},
				{
					Token:  destToken.Address(),
					Amount: oneE18,
				},
				{
					Token:  selfServeDestToken.Address(),
					Amount: oneE18,
				},
			},
			receiver:  state.Chains[sourceChain].Receiver.Address(),
			extraData: changeset.MakeEVMExtraArgsV2(300_000, false),
			expectedTokenBalances: map[common.Address]*big.Int{
				selfServeSrcToken.Address(): new(big.Int).Add(oneE18, oneE18),
				srcToken.Address():          oneE18,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:     "Sending token transfer with custom gasLimits to the EOA is successful",
			srcChain: destChain,
			dstChain: sourceChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  selfServeDestToken.Address(),
					Amount: oneE18,
				},
				{
					Token:  destToken.Address(),
					Amount: new(big.Int).Add(oneE18, oneE18),
				},
			},
			receiver:  utils.RandomAddress(),
			extraData: changeset.MakeEVMExtraArgsV2(1, false),
			expectedTokenBalances: map[common.Address]*big.Int{
				selfServeSrcToken.Address(): oneE18,
				srcToken.Address():          new(big.Int).Add(oneE18, oneE18),
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:     "Sending PTT with too low gas limit leads to the revert when receiver is a contract",
			srcChain: destChain,
			dstChain: sourceChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  selfServeDestToken.Address(),
					Amount: oneE18,
				},
				{
					Token:  destToken.Address(),
					Amount: oneE18,
				},
			},
			receiver:  state.Chains[sourceChain].Receiver.Address(),
			data:      []byte("this should be reverted because gasLimit is too low, no tokens are transferred as well"),
			extraData: changeset.MakeEVMExtraArgsV2(1, false),
			expectedTokenBalances: map[common.Address]*big.Int{
				selfServeSrcToken.Address(): big.NewInt(0),
				srcToken.Address():          big.NewInt(0),
			},
			expectedExecutionState: changeset.EXECUTION_STATE_FAILURE,
		},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			initialBalances := map[common.Address]*big.Int{}
			for token := range tt.expectedTokenBalances {
				initialBalance := changeset.GetTokenBalance(ctx, t, token, tt.receiver, e.Chains[tt.dstChain])
				initialBalances[token] = initialBalance
			}

			changeset.TransferAndWaitForSuccess(
				ctx,
				t,
				e,
				state,
				tt.srcChain,
				tt.dstChain,
				tt.tokenAmounts,
				tt.receiver,
				tt.data,
				tt.expectedExecutionState,
				tt.extraData,
			)

			for token, balance := range tt.expectedTokenBalances {
				expected := new(big.Int).Add(initialBalances[token], balance)
				changeset.WaitForTheTokenBalance(ctx, t, token, tt.receiver, e.Chains[tt.dstChain], expected)
			}
		})
	}
}

func createAndFundSelfServeActor(
	ctx context.Context,
	t *testing.T,
	deployer *bind.TransactOpts,
	chain deployment.Chain,
	amountToFund *big.Int,
	isInMemory bool,
) *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	// Simulated backend sets chainID to 1337 always
	chainID := big.NewInt(1337)
	if !isInMemory {
		// Docker environment runs real geth so chainID has to be set accordingly
		stringChainID, err1 := sel.GetChainIDFromSelector(chain.Selector)
		require.NoError(t, err1)
		chainID, _ = new(big.Int).SetString(stringChainID, 10)
	}

	actor, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	require.NoError(t, err)

	nonce, err := chain.Client.PendingNonceAt(ctx, deployer.From)
	require.NoError(t, err)

	gasPrice, err := chain.Client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &actor.From,
		Value:    amountToFund,
		Gas:      uint64(21000),
		GasPrice: gasPrice,
		Data:     nil,
	})

	signedTx, err := deployer.Signer(deployer.From, tx)
	require.NoError(t, err)

	err = chain.Client.SendTransaction(ctx, signedTx)
	require.NoError(t, err)

	_, err = chain.Confirm(signedTx)
	require.NoError(t, err)

	return actor
}
