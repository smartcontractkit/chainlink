package smoke

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

func TestMigrateFromV1_5ToV1_6(t *testing.T) {
	// Deploy CCIP 1.5 with 3 chains and 4 nodes + 1 bootstrap
	// Deploy 1.5 contracts (excluding pools and real RMN, use MockRMN to start, but including MCMS) .
	e := changeset.NewMemoryEnvironment(
		t,
		changeset.WithLegacyDeployment(),
		changeset.WithChains(3),
		changeset.WithChainIds([]uint64{chainselectors.GETH_TESTNET.EvmChainID}))
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChainsExcept1337 := e.Env.AllChainSelectorsExcluding([]uint64{chainselectors.GETH_TESTNET.Selector})
	require.Contains(t, e.Env.AllChainSelectors(), chainselectors.GETH_TESTNET.Selector)
	require.Len(t, allChainsExcept1337, 2)
	src1, src2, dest := allChainsExcept1337[0], allChainsExcept1337[1], chainselectors.GETH_TESTNET.Selector
	destChain := e.Env.Chains[dest]
	pairs := []changeset.SourceDestPair{
		{SourceChainSelector: src1, DestChainSelector: dest},
		{SourceChainSelector: src2, DestChainSelector: dest},
	}
	// wire up all lanes
	// deploy onRamp, commit store, offramp , set ocr2config and send corresponding jobs
	e.Env = v1_5.AddLanes(t, e.Env, state, pairs)
	// reload state after adding lanes
	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	// ensure that all lanes are functional
	for _, pair := range pairs {
		sentEvent, err := v1_5.SendRequest(t, e.Env, state,
			changeset.WithSourceChain(pair.SourceChainSelector),
			changeset.WithDestChain(pair.DestChainSelector),
			changeset.WithTestRouter(false),
			changeset.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[pair.DestChainSelector].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			}),
		)
		require.NoError(t, err)
		require.NotNil(t, sentEvent)
		destStartBlock, err := destChain.Client.HeaderByNumber(context.Background(), nil)
		require.NoError(t, err)
		v1_5.WaitForCommit(t, e.Env.Chains[pair.SourceChainSelector], destChain, state.Chains[dest].CommitStore[src1], sentEvent.Message.SequenceNumber)
		v1_5.WaitForExecute(t, e.Env.Chains[pair.SourceChainSelector], destChain, state.Chains[dest].EVM2EVMOffRamp[src1], []uint64{sentEvent.Message.SequenceNumber}, destStartBlock.Number.Uint64())
	}
	// now that all lanes work transfer ownership of the contracts to MCMS
	// add 1.6 contracts to the environment and send 1.6 jobs
}
