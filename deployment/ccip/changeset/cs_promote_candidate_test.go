package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func Test_PromoteCandidate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv := NewMemoryEnvironment(t,
				WithChains(2),
				WithNodes(4))
			state, err := LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			// Deploy to all chains.
			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]
			t.Logf("source: %d, dest: %d, home chain: %d, feed chain: %d", source, dest, tenv.HomeChainSel, tenv.FeedChainSel)
			newAddresses := deployment.NewMemoryAddressBook()
			err = deployPrerequisiteChainContracts(tenv.Env, newAddresses, allChains)
			require.NoError(t, err)
			require.NoError(t, tenv.Env.ExistingAddresses.Merge(newAddresses))

			nodes, err := deployment.NodeInfo(tenv.Env.NodeIDs, tenv.Env.Offchain)
			require.NoError(t, err)

			var nodeIDs []string
			for _, node := range nodes {
				nodeIDs = append(nodeIDs, node.NodeID)
			}

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				_, err = commonchangeset.ApplyChangesets(t, tenv.Env, map[uint64]*commonchangeset.TimelockExecutionContracts{
					source: {
						Timelock:  state.Chains[source].Timelock,
						CallProxy: state.Chains[source].CallProxy,
					},
					dest: {
						Timelock:  state.Chains[dest].Timelock,
						CallProxy: state.Chains[dest].CallProxy,
					},
					tenv.HomeChainSel: {
						Timelock:  state.Chains[tenv.HomeChainSel].Timelock,
						CallProxy: state.Chains[tenv.HomeChainSel].CallProxy,
					},
				}, []commonchangeset.ChangesetApplication{
					{
						Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
						Config:    genTestTransferOwnershipConfig(tenv, allChains, state),
					},
				})
				require.NoError(t, err)
				assertTimelockOwnership(t, tenv, allChains, state)
			}

			var (
				capReg   = state.Chains[tenv.HomeChainSel].CapabilityRegistry
				ccipHome = state.Chains[tenv.HomeChainSel].CCIPHome
			)
			donID, _, err := internal.DonIDForChain(capReg, ccipHome, dest)
			require.NoError(t, err)
			candidateDigestCommitBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestCommitBefore)
			candidateDigestExecBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestExecBefore)

			var mcmsConfig *MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &MCMSConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, map[uint64]*commonchangeset.TimelockExecutionContracts{
				tenv.HomeChainSel: {
					Timelock:  state.Chains[tenv.HomeChainSel].Timelock,
					CallProxy: state.Chains[tenv.HomeChainSel].CallProxy,
				},
			}, []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(PromoteAllCandidatesChangeset),
					Config: PromoteAllCandidatesChangesetConfig{
						HomeChainSelector: tenv.HomeChainSel,
						NewChainSelector:  tenv.FeedChainSel,
						NodeIDs:           nodeIDs,
						MCMS:              mcmsConfig,
					},
				},
			})
			require.NoError(t, err)

			// There seems to be some flakiness where the chain state isn't properly updated.
			for i := 0; i < 10; i++ {
				tenv.Env.Chains[tenv.HomeChainSel].Client.(*memory.Backend).Commit()
			}

			// after promoting the zero digest, active digest should also be zero
			activeDigestCommit, err := ccipHome.GetActiveDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, activeDigestCommit)

			activeDigestExec, err := ccipHome.GetActiveDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, activeDigestExec)
		})
	}
}
