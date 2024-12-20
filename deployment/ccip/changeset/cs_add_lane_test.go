package changeset

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

func TestAddLanesWithTestRouter(t *testing.T) {
	t.Parallel()
	e := NewMemoryEnvironment(t)
	// Here we have CR + nodes set up, but no CCIP contracts deployed.
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	selectors := e.Env.AllChainSelectors()
	chain1, chain2 := selectors[0], selectors[1]

	stateChain1 := state.Chains[chain1]
	e.Env, err = commoncs.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commoncs.ChangesetApplication{
		{
			Changeset: commoncs.WrapChangeSet(UpdateOnRampsDests),
			Config: UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]OnRampDestinationUpdate{
					chain1: {
						chain2: {
							IsEnabled:        true,
							TestRouter:       true,
							AllowListEnabled: false,
						},
					},
				},
			},
		},
		{
			Changeset: commoncs.WrapChangeSet(UpdateFeeQuoterPricesCS),
			Config: UpdateFeeQuoterPricesConfig{
				PricesByChain: map[uint64]FeeQuoterPriceUpdatePerSource{
					chain1: {
						TokenPrices: map[common.Address]*big.Int{
							stateChain1.LinkToken.Address(): DefaultLinkPrice,
							stateChain1.Weth9.Address():     DefaultWethPrice,
						},
						GasPrices: map[uint64]*big.Int{
							chain2: DefaultGasPrice,
						},
					},
				},
			},
		},
		{
			Changeset: commoncs.WrapChangeSet(UpdateFeeQuoterDests),
			Config: UpdateFeeQuoterDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
					chain1: {
						chain2: DefaultFeeQuoterDestChainConfig(),
					},
				},
			},
		},
		{
			Changeset: commoncs.WrapChangeSet(UpdateOffRampSources),
			Config: UpdateOffRampSourcesConfig{
				UpdatesByChain: map[uint64]map[uint64]OffRampSourceUpdate{
					chain2: {
						chain1: {
							IsEnabled:  true,
							TestRouter: true,
						},
					},
				},
			},
		},
		{
			Changeset: commoncs.WrapChangeSet(UpdateRouterRamps),
			Config: UpdateRouterRampsConfig{
				TestRouter: true,
				UpdatesByChain: map[uint64]RouterUpdates{
					// onRamp update on source chain
					chain1: {
						OnRampUpdates: map[uint64]bool{
							chain2: true,
						},
					},
					// offramp update on dest chain
					chain2: {
						OffRampUpdates: map[uint64]bool{
							chain1: true,
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNumExec := make(map[SourceDestPair][]uint64)
	latesthdr, err := e.Env.Chains[chain2].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[chain2] = &block
	msgSentEvent := TestSendRequest(t, e.Env, state, chain1, chain2, true, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[chain2].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	expectedSeqNumExec[SourceDestPair{
		SourceChainSelector: chain1,
		DestChainSelector:   chain2,
	}] = []uint64{msgSentEvent.SequenceNumber}
	ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)
}
