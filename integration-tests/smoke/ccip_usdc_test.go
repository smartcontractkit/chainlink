package smoke

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestUSDCTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)

	e := tenv.Env
	state, err := ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Chains)
	sourceChain := allChainSelectors[0]
	destChain := allChainSelectors[1]

	srcUSDC, dstUSDC, err := ccdeploy.ConfigureUSDCTokenPools(lggr, e.Chains, sourceChain, destChain, state)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))

	mintAndAllow(t, e, state, map[uint64]*burn_mint_erc677.BurnMintERC677{
		sourceChain: srcUSDC,
		destChain:   dstUSDC,
	})

	err = ccdeploy.UpdateFeeQuoterForUSDC(lggr, e.Chains[sourceChain], state.Chains[sourceChain], destChain, srcUSDC)
	require.NoError(t, err)

	err = ccdeploy.UpdateFeeQuoterForUSDC(lggr, e.Chains[destChain], state.Chains[destChain], sourceChain, dstUSDC)
	require.NoError(t, err)

	// MockE2EUSDCTransmitter always mint 1, see MockE2EUSDCTransmitter.sol for more details
	tinyOneCoin := new(big.Int).SetUint64(1)

	srcDstTokenMapping := map[common.Address]*burn_mint_erc677.BurnMintERC677{
		srcUSDC.Address(): dstUSDC,
		dstUSDC.Address(): srcUSDC,
	}

	tcs := []struct {
		name        string
		receiver    common.Address
		sourceChain uint64
		destChain   uint64
		tokens      []router.ClientEVMTokenAmount
		data        []byte
	}{
		{
			name:        "single USDC token transfer to EOA",
			receiver:    utils.RandomAddress(),
			sourceChain: destChain,
			destChain:   sourceChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  dstUSDC.Address(),
					Amount: tinyOneCoin,
				}},
		},
		{
			name:        "programmable token transfer to valid contract receiver",
			receiver:    state.Chains[destChain].Receiver.Address(),
			sourceChain: sourceChain,
			destChain:   destChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  srcUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			data: []byte("hello world"),
		},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			initialBalances := map[common.Address]*big.Int{}
			for _, token := range tt.tokens {
				destToken := srcDstTokenMapping[token.Token]

				initialBalance, err := destToken.BalanceOf(&bind.CallOpts{Context: tests.Context(t)}, tt.receiver)
				require.NoError(t, err)
				initialBalances[token.Token] = initialBalance
			}

			transferAndWaitForSuccess(
				t,
				e,
				state,
				tt.sourceChain,
				tt.destChain,
				tt.tokens,
				tt.receiver,
				tt.data,
			)

			for _, token := range tt.tokens {
				destToken := srcDstTokenMapping[token.Token]

				balance, err := destToken.BalanceOf(&bind.CallOpts{Context: tests.Context(t)}, tt.receiver)
				require.NoError(t, err)
				require.Equal(t, new(big.Int).Add(initialBalances[token.Token], tinyOneCoin), balance)
			}
		})
	}
}

// mintAndAllow mints tokens for deployers and allow router to spend them
func mintAndAllow(
	t *testing.T,
	e deployment.Environment,
	state ccdeploy.CCIPOnChainState,
	tokens map[uint64]*burn_mint_erc677.BurnMintERC677,
) {
	for chain, token := range tokens {
		twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))

		tx, err := token.Mint(
			e.Chains[chain].DeployerKey,
			e.Chains[chain].DeployerKey.From,
			new(big.Int).Mul(twoCoins, big.NewInt(10)),
		)
		require.NoError(t, err)
		_, err = e.Chains[chain].Confirm(tx)
		require.NoError(t, err)

		tx, err = token.Approve(e.Chains[chain].DeployerKey, state.Chains[chain].Router.Address(), twoCoins)
		require.NoError(t, err)
		_, err = e.Chains[chain].Confirm(tx)
		require.NoError(t, err)
	}
}

// transferAndWaitForSuccess sends a message from sourceChain to destChain and waits for it to be executed
func transferAndWaitForSuccess(
	t *testing.T,
	env deployment.Environment,
	state ccdeploy.CCIPOnChainState,
	sourceChain, destChain uint64,
	tokens []router.ClientEVMTokenAmount,
	receiver common.Address,
	data []byte,
) {
	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[uint64]uint64)

	latesthdr, err := env.Chains[destChain].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[destChain] = &block

	msgSentEvent := ccdeploy.TestSendRequest(t, env, state, sourceChain, destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
		Data:         data,
		TokenAmounts: tokens,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	expectedSeqNum[destChain] = msgSentEvent.SequenceNumber

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, env, state, expectedSeqNum, startBlocks)

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, env, state, expectedSeqNum, startBlocks)
}
