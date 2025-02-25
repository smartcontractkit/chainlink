package v1_5_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	v1_5changeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/rmn_contract"
)

// TestV1_5_Message_RMNRemote this test verify that 1.5 lane can send message when using RMNRemote
func TestV1_5_Message_RMNRemote(t *testing.T) {
	// Deploy CCIP 1.5 with 3 chains and 4 nodes + 1 bootstrap
	// Deploy 1.5 contracts (excluding pools to start, but including MCMS) .
	e, tEnv := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithPrerequisiteDeploymentOnly(
			&changeset.V1_5DeploymentConfig{
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
	)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChains := e.Env.AllChainSelectors()
	src1, dest := allChains[0], allChains[1]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: src1, DestChainSelector: dest},
	}
	// wire up all lanes
	// deploy onRamp, commit store, offramp , set ocr2config and send corresponding jobs
	e.Env = v1_5.AddLanes(t, e.Env, state, pairs)

	// permabless the commit stores
	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_5changeset.PermaBlessCommitStoreChangeset),
			v1_5changeset.PermaBlessCommitStoreConfig{
				Configs: map[uint64]v1_5changeset.PermaBlessCommitStoreConfigPerDest{
					dest: {
						Sources: []v1_5changeset.PermaBlessConfigPerSourceChain{
							{
								SourceChainSelector: src1,
								PermaBless:          true,
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
	oldState, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	e = testhelpers.AddCCIPContractsToEnvironment(t, e.Env.AllChainSelectors(), tEnv, false)
	// reload state after adding lanes

	tEnv.UpdateDeployedEnvironment(e)

	_, err = deployment.CreateLegacyChangeSet(v1_6.SetRMNRemoteOnRMNProxyChangeset).Apply(e.Env,
		v1_6.SetRMNRemoteOnRMNProxyConfig{
			ChainSelectors: e.Env.AllChainSelectors(),
		})
	require.NoError(t, err)

	// send continuous messages in real router until done is closed
	// send a message from the other lane src1 -> dest
	sentEvent, err := v1_5.SendRequest(t, e.Env, oldState,
		testhelpers.WithSourceChain(src1),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(false),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(oldState.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, sentEvent)
	destChain := e.Env.Chains[dest]
	require.NoError(t, err)
	v1_5.WaitForCommit(t, e.Env.Chains[src1], destChain, oldState.Chains[dest].CommitStore[src1],
		sentEvent.Message.SequenceNumber)
}

// TestV1_5_Message_RMNRemote this test verify that 1.5 lane can be cursed when using RMNRemote
func TestV1_5_Message_RMNRemote_Curse(t *testing.T) {
	// Deploy CCIP 1.5 with 3 chains and 4 nodes + 1 bootstrap
	// Deploy 1.5 contracts (excluding pools to start, but including MCMS) .
	e, tEnv := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithPrerequisiteDeploymentOnly(
			&changeset.V1_5DeploymentConfig{
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
	)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChains := e.Env.AllChainSelectors()
	src1, dest := allChains[0], allChains[1]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: src1, DestChainSelector: dest},
	}
	// wire up all lanes
	// deploy onRamp, commit store, offramp , set ocr2config and send corresponding jobs
	e.Env = v1_5.AddLanes(t, e.Env, state, pairs)

	// permabless the commit stores
	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_5changeset.PermaBlessCommitStoreChangeset),
			v1_5changeset.PermaBlessCommitStoreConfig{
				Configs: map[uint64]v1_5changeset.PermaBlessCommitStoreConfigPerDest{
					dest: {
						Sources: []v1_5changeset.PermaBlessConfigPerSourceChain{
							{
								SourceChainSelector: src1,
								PermaBless:          true,
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
	oldState, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	e = testhelpers.AddCCIPContractsToEnvironment(t, e.Env.AllChainSelectors(), tEnv, false)

	// reload state after adding lanes
	tEnv.UpdateDeployedEnvironment(e)

	_, err = deployment.CreateLegacyChangeSet(v1_6.SetRMNRemoteOnRMNProxyChangeset).Apply(e.Env,
		v1_6.SetRMNRemoteOnRMNProxyConfig{
			ChainSelectors: e.Env.AllChainSelectors(),
		})
	require.NoError(t, err)

	// send continuous messages in real router until done is closed
	// send a message from the other lane src1 -> dest
	sentEvent, err := v1_5.SendRequest(t, e.Env, oldState,
		testhelpers.WithSourceChain(src1),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(false),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(oldState.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}),
	)
	require.NoError(t, err)

	_, err = deployment.CreateLegacyChangeSet(v1_6.RMNCurseChangeset).Apply(e.Env, v1_6.RMNCurseConfig{
		CurseActions: []v1_6.CurseAction{v1_6.CurseChain(e.Env.AllChainSelectors()[0])},
		Reason:       "Curse test",
	})
	require.NoError(t, err)

	require.NotNil(t, sentEvent)
	destChain := e.Env.Chains[dest]
	require.NoError(t, err)
	v1_5.WaitForNoCommit(t, e.Env.Chains[src1], destChain, oldState.Chains[dest].CommitStore[src1],
		sentEvent.Message.SequenceNumber)
}

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
	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_5changeset.PermaBlessCommitStoreChangeset),
			v1_5changeset.PermaBlessCommitStoreConfig{
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
		),
	)
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
