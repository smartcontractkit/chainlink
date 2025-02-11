package v1_5_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	v1_5changeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/rmn_contract"
)

// This test only works if the destination chain id is 1337
// Otherwise it shows error for offchain and onchain config digest mismatch
func TestE2ELegacy(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithPrerequisiteDeploymentOnly(&changeset.V1_5DeploymentConfig{
			PriceRegStalenessThreshold: 60 * 60 * 24 * 14, // two weeks
			RMNConfig: &rmn_contract.RMNConfig{
				BlessWeightThreshold: 2,
				CurseWeightThreshold: 2,
				// setting dummy voters, we will permabless this later
				Voters: []rmn_contract.RMNVoter{
					{
						BlessWeight:   2,
						CurseWeight:   2,
						BlessVoteAddr: utils.RandomAddress(),
						CurseVoteAddr: utils.RandomAddress(),
					},
				},
			},
		}),
		testhelpers.WithNumOfChains(3),
		testhelpers.WithChainIDs([]uint64{chainselectors.GETH_TESTNET.EvmChainID}))
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChains := e.Env.AllChainSelectorsExcluding([]uint64{chainselectors.GETH_TESTNET.Selector})
	require.Contains(t, e.Env.AllChainSelectors(), chainselectors.GETH_TESTNET.Selector)
	require.Len(t, allChains, 2)
	src, dest := allChains[1], chainselectors.GETH_TESTNET.Selector
	srcChain := e.Env.Chains[src]
	destChain := e.Env.Chains[dest]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: src, DestChainSelector: dest},
	}
	e.Env = v1_5.AddLanes(t, e.Env, state, pairs)
	// permabless the commit stores
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(v1_5changeset.PermaBlessCommitStoreChangeset),
			Config: v1_5changeset.PermaBlessCommitStoreConfig{
				Configs: map[uint64]v1_5changeset.PermaBlessCommitStoreConfigPerDest{
					dest: {
						Sources: []v1_5changeset.PermaBlessConfigPerSourceChain{
							{
								SourceChainSelector: src,
								PermaBless:          true,
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	// reload state after adding lanes
	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	sentEvent, err := v1_5.SendRequest(t, e.Env, state,
		testhelpers.WithSourceChain(src),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(false),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
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
	v1_5.WaitForCommit(t, srcChain, destChain, state.Chains[dest].CommitStore[src], sentEvent.Message.SequenceNumber)
	v1_5.WaitForExecute(t, srcChain, destChain, state.Chains[dest].EVM2EVMOffRamp[src], []uint64{sentEvent.Message.SequenceNumber}, destStartBlock.Number.Uint64())
}
