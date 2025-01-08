package changeset_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestUpdateNodes(t *testing.T) {
	t.Parallel()

	t.Run("no mcms", func(t *testing.T) {
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
		})

		updates := make(map[p2pkey.PeerID]changeset.NodeUpdate)
		i := uint8(0)
		for id := range te.WFNodes {
			k, err := p2pkey.MakePeerID(id)
			require.NoError(t, err)
			pubKey := [32]byte{31: i + 1}
			// don't set capabilities or nop b/c those must already exist in the contract
			// those ops must be a different proposal when using MCMS
			updates[k] = changeset.NodeUpdate{
				EncryptionPublicKey: hex.EncodeToString(pubKey[:]),
				Signer:              [32]byte{0: i + 1},
			}
			i++
		}

		cfg := changeset.UpdateNodesRequest{
			RegistryChainSel: te.RegistrySelector,
			P2pToUpdates:     updates,
		}

		csOut, err := changeset.UpdateNodes(te.Env, &cfg)
		require.NoError(t, err)
		require.Empty(t, csOut.Proposals)
		require.Nil(t, csOut.AddressBook)

		validateUpdate(t, te, updates)
	})

	t.Run("with mcms", func(t *testing.T) {
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		updates := make(map[p2pkey.PeerID]changeset.NodeUpdate)
		i := uint8(0)
		for id := range te.WFNodes {
			k, err := p2pkey.MakePeerID(id)
			require.NoError(t, err)
			pubKey := [32]byte{31: i + 1}
			// don't set capabilities or nop b/c those must already exist in the contract
			// those ops must be a different proposal when using MCMS
			updates[k] = changeset.NodeUpdate{
				EncryptionPublicKey: hex.EncodeToString(pubKey[:]),
				Signer:              [32]byte{0: i + 1},
			}
			i++
		}

		cfg := changeset.UpdateNodesRequest{
			RegistryChainSel: te.RegistrySelector,
			P2pToUpdates:     updates,
			MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
		}

		csOut, err := changeset.UpdateNodes(te.Env, &cfg)
		require.NoError(t, err)
		require.Len(t, csOut.Proposals, 1)
		require.Nil(t, csOut.AddressBook)

		// now apply the changeset such that the proposal is signed and execed
		contracts := te.ContractSets()[te.RegistrySelector]
		timelockContracts := map[uint64]*proposalutils.TimelockExecutionContracts{
			te.RegistrySelector: {
				Timelock:  contracts.Timelock,
				CallProxy: contracts.CallProxy,
			},
		}
		_, err = commonchangeset.ApplyChangesets(t, te.Env, timelockContracts, []commonchangeset.ChangesetApplication{
			{
				Changeset: commonchangeset.WrapChangeSet(changeset.UpdateNodes),
				Config: &changeset.UpdateNodesRequest{
					RegistryChainSel: te.RegistrySelector,
					P2pToUpdates:     updates,
					MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
				},
			},
		})
		require.NoError(t, err)

		validateUpdate(t, te, updates)
	})
}

// validateUpdate checks reads nodes from the registry and checks they have the expected updates
func validateUpdate(t *testing.T, te test.TestEnv, expected map[p2pkey.PeerID]changeset.NodeUpdate) {
	registry := te.ContractSets()[te.RegistrySelector].CapabilitiesRegistry
	wfP2PIDs := p2pIDs(t, maps.Keys(te.WFNodes))
	nodes, err := registry.GetNodesByP2PIds(nil, wfP2PIDs)
	require.NoError(t, err)
	require.Len(t, nodes, len(wfP2PIDs))
	for _, node := range nodes {
		// only check the fields that were updated
		assert.Equal(t, expected[node.P2pId].EncryptionPublicKey, hex.EncodeToString(node.EncryptionPublicKey[:]))
		assert.Equal(t, expected[node.P2pId].Signer, node.Signer)
	}
}

func p2pIDs(t *testing.T, vals []string) [][32]byte {
	var out [][32]byte
	for _, v := range vals {
		id, err := p2pkey.MakePeerID(v)
		require.NoError(t, err)
		out = append(out, id)
	}
	return out
}
