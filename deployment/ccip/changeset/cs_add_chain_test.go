package changeset

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	changesetcommon "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func Test_AddChain(t *testing.T) {
	t.Parallel()

	const (
		numChains     = 4
		usersPerChain = 2
	)

	// Set up an env with 4 chains but initially
	// only deploy and configure 3 of them.
	e, tEnv := NewMemoryEnvironment(
		t,
		WithChains(numChains),
		WithNodes(4),
		WithPrerequisiteDeployment(),
		WithUsersPerChain(usersPerChain),
		WithNoJobsAndContracts(),
	)

	allChains := maps.Keys(e.Env.Chains)
	toDeploy := e.Env.AllChainSelectorsExcluding([]uint64{allChains[0]})
	require.Len(t, toDeploy, numChains-1)
	e = AddCCIPContractsToEnvironment(t, toDeploy, tEnv)

	// Need to update what the RMNProxy is pointing to, otherwise plugin will not work.
	var err error
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(SetRMNRemoteOnRMNProxy),
			Config: SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: toDeploy,
			},
		},
	})
	require.NoError(t, err)

	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Setup densely connected lanes between all chains.
	for _, source := range toDeploy {
		for _, dest := range toDeploy {
			if source == dest {
				continue
			}
			AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, source, dest, false)
		}
	}

	// Transfer ownership of all contracts to the MCMS and renounce the timelock deployer.
	transferToMCMSAndRenounceTimelockDeployer(t, e, toDeploy, state)

	// At this stage we can send some requests and confirm the setup is working.
	sendMsgs := func(chains []uint64) {
		var (
			startBlocks        = make(map[uint64]*uint64)
			expectedSeqNum     = make(map[SourceDestPair]uint64)
			expectedSeqNumExec = make(map[SourceDestPair][]uint64)
		)
		for _, source := range chains {
			for _, dest := range chains {
				if source == dest {
					continue
				}
				latesthdr, err := e.Env.Chains[dest].Client.HeaderByNumber(testcontext.Get(t), nil)
				require.NoError(t, err)
				block := latesthdr.Number.Uint64()
				msgSentEvent := TestSendRequest(t, e.Env, state, source, dest, false, router.ClientEVM2AnyMessage{
					Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
					Data:         []byte("hello world"),
					TokenAmounts: nil,
					FeeToken:     common.HexToAddress("0x0"),
					ExtraArgs:    nil,
				})

				startBlocks[dest] = &block
				expectedSeqNum[SourceDestPair{
					SourceChainSelector: source,
					DestChainSelector:   dest,
				}] = msgSentEvent.SequenceNumber
				expectedSeqNumExec[SourceDestPair{
					SourceChainSelector: source,
					DestChainSelector:   dest,
				}] = append(expectedSeqNumExec[SourceDestPair{
					SourceChainSelector: source,
					DestChainSelector:   dest,
				}], msgSentEvent.SequenceNumber)
			}
		}

		// Confirm execution of the message
		ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNum, startBlocks)
		ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)
	}

	// check RMNRemote is up.
	for _, chain := range toDeploy {
		require.NotEqual(t, common.Address{}, state.Chains[chain].RMNRemote.Address())
		_, err := state.Chains[chain].RMNRemote.GetCursedSubjects(&bind.CallOpts{
			Context: tests.Context(t),
		})
		require.NoError(t, err)

		// check which address RMNProxy is pointing to
		rmnAddress, err := state.Chains[chain].RMNProxy.GetARM(&bind.CallOpts{
			Context: tests.Context(t),
		})
		require.NoError(t, err)
		require.Equal(t, state.Chains[chain].RMNRemote.Address(), rmnAddress)

		t.Log("RMNRemote address for chain", chain, "is:", state.Chains[chain].RMNRemote.Address().Hex())
		t.Log("RMNProxy address for chain", chain, "is:", state.Chains[chain].RMNProxy.Address().Hex())
	}

	// wait for plugins to come up.
	time.Sleep(30 * time.Second)
	sendMsgs(toDeploy)
}

func transferToMCMSAndRenounceTimelockDeployer(
	t *testing.T,
	e DeployedEnv,
	chains []uint64,
	state CCIPOnChainState,
) {
	var apps []changesetcommon.ChangesetApplication
	apps = append(apps, changesetcommon.ChangesetApplication{
		Changeset: changesetcommon.WrapChangeSet(changesetcommon.TransferToMCMSWithTimelock),
		Config:    genTestTransferOwnershipConfig(e, chains, state),
	})
	for _, chain := range chains {
		apps = append(apps, changesetcommon.ChangesetApplication{
			Changeset: changesetcommon.WrapChangeSet(changesetcommon.RenounceTimelockDeployer),
			Config: changesetcommon.RenounceTimelockDeployerConfig{
				ChainSel: chain,
			},
		})
	}
	var err error
	e.Env, err = changesetcommon.ApplyChangesets(t, e.Env, e.TimelockContracts(t), apps)
	require.NoError(t, err)
}
