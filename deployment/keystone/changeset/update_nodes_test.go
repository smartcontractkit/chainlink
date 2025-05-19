package changeset_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestUpdateNodes(t *testing.T) {
	t.Parallel()

	t.Run("no mcms", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		updates := make(map[p2pkey.PeerID]changeset.NodeUpdate)
		i := uint8(0)
		for _, id := range te.GetP2PIDs("wfDon") {
			pubKey := [32]byte{31: i + 1}
			// don't set capabilities or nop b/c those must already exist in the contract
			// those ops must be a different proposal when using MCMS
			updates[id] = changeset.NodeUpdate{
				EncryptionPublicKey: hex.EncodeToString(pubKey[:]),
				Signer:              [32]byte{0: i + 1},
			}
			i++
		}

		cfg := changeset.UpdateNodesRequest{
			RegistryChainSel: te.RegistrySelector,
			P2pToUpdates:     updates,
			RegistryRef:      te.CapabilityRegistryAddressRef(),
		}

		csOut, err := changeset.UpdateNodes(te.Env, &cfg)
		require.NoError(t, err)
		require.Empty(t, csOut.Proposals)
		require.Nil(t, csOut.AddressBook)

		validateUpdate(t, te, updates)
	})

	t.Run("with mcms", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		updates := make(map[p2pkey.PeerID]changeset.NodeUpdate)
		i := uint8(0)
		for _, id := range te.GetP2PIDs("wfDon") {
			pubKey := [32]byte{31: i + 1}
			// don't set capabilities or nop b/c those must already exist in the contract
			// those ops must be a different proposal when using MCMS
			updates[id] = changeset.NodeUpdate{
				EncryptionPublicKey: hex.EncodeToString(pubKey[:]),
				Signer:              [32]byte{0: i + 1},
			}
			i++
		}

		cfg := changeset.UpdateNodesRequest{
			RegistryChainSel: te.RegistrySelector,
			P2pToUpdates:     updates,
			MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
			RegistryRef:      te.CapabilityRegistryAddressRef(),
		}

		csOut, err := changeset.UpdateNodes(te.Env, &cfg)
		require.NoError(t, err)
		require.Len(t, csOut.MCMSTimelockProposals, 1)
		require.Nil(t, csOut.AddressBook)

		err = applyProposal(t, te, commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.UpdateNodes),
			&changeset.UpdateNodesRequest{
				RegistryChainSel: te.RegistrySelector,
				P2pToUpdates:     updates,
				MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
				RegistryRef:      te.CapabilityRegistryAddressRef(),
			},
		))
		require.NoError(t, err)

		validateUpdate(t, te, updates)
	})
}

// validateUpdate checks reads nodes from the registry and checks they have the expected updates
func validateUpdate(t *testing.T, te test.EnvWrapper, expected map[p2pkey.PeerID]changeset.NodeUpdate) {
	registry := te.CapabilitiesRegistry()
	wfP2PIDs := te.GetP2PIDs("wfDon").Bytes32()
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
