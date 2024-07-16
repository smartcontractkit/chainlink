package ccip_integration_tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/consul/sdk/freeport"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/libocr/commontypes"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/stretchr/testify/require"
)

func TestIntegration_OCR3Nodes(t *testing.T) {
	numChains := 3
	homeChainUni, universes := createUniverses(t, numChains)
	numNodes := 4
	t.Log("creating ocr3 nodes")
	var (
		oracles      = make(map[uint64][]confighelper2.OracleIdentityExtra)
		transmitters = make(map[uint64][]common.Address)
		apps         []chainlink.Application
		nodes        []*ocr3Node
		p2pIDs       [][32]byte

		// The bootstrap node will be the first node (index 0)
		bootstrapPort  int
		bootstrapP2PID p2pkey.PeerID
		bootstrappers  []commontypes.BootstrapperLocator
	)

	ports := freeport.GetN(t, numNodes)
	capabilitiesPorts := freeport.GetN(t, numNodes)
	for i := 0; i < numNodes; i++ {
		node := setupNodeOCR3(t, ports[i], capabilitiesPorts[i], bootstrappers, universes, homeChainUni)

		apps = append(apps, node.app)
		for chainID, transmitter := range node.transmitters {
			transmitters[chainID] = append(transmitters[chainID], transmitter)
			identity := confighelper2.OracleIdentityExtra{
				OracleIdentity: confighelper2.OracleIdentity{
					OnchainPublicKey:  node.keybundle.PublicKey(),
					TransmitAccount:   ocrtypes.Account(transmitter.Hex()),
					OffchainPublicKey: node.keybundle.OffchainPublicKey(),
					PeerID:            node.peerID,
				},
				ConfigEncryptionPublicKey: node.keybundle.ConfigEncryptionPublicKey(),
			}
			oracles[chainID] = append(oracles[chainID], identity)
		}
		nodes = append(nodes, node)
		peerID, err := p2pkey.MakePeerID(node.peerID)
		require.NoError(t, err)
		p2pIDs = append(p2pIDs, peerID)

		// First Node is the bootstrap node
		if i == 0 {
			bootstrapPort = ports[i]
			bootstrapP2PID = peerID
			bootstrappers = []commontypes.BootstrapperLocator{
				{PeerID: node.peerID, Addrs: []string{
					fmt.Sprintf("127.0.0.1:%d", bootstrapPort),
				}},
			}
		}
	}

	// Start committing periodically in the background for all the chains
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	commitBlocksBackground(t, universes, tick)

	ctx := testutils.Context(t)
	t.Log("creating ocr3 jobs")
	for i := 0; i < len(nodes); i++ {
		err := nodes[i].app.Start(ctx)
		require.NoError(t, err)
		tApp := apps[i]
		t.Cleanup(func() {
			require.NoError(t, tApp.Stop())
		})
		//TODO: Create Jobs and add them to the app
	}

	ccipCapabilityID, err := homeChainUni.capabilityRegistry.GetHashedCapabilityId(&bind.CallOpts{
		Context: ctx,
	}, CapabilityLabelledName, CapabilityVersion)
	require.NoError(t, err, "failed to get hashed capability id for ccip")
	require.NotEqual(t, [32]byte{}, ccipCapabilityID, "ccip capability id is empty")

	// Need to Add nodes and assign capabilities to them before creating DONS
	homeChainUni.AddNodes(t, p2pIDs, [][32]byte{ccipCapabilityID})
	// Create a DON for each chain
	for _, uni := range universes {
		// Add nodes and give them the capability
		t.Log("AddingDON for universe: ", uni.chainID)
		homeChainUni.AddDON(t,
			ccipCapabilityID,
			uni.chainID,
			uni.offramp.Address().Bytes(),
			1, // f
			bootstrapP2PID,
			p2pIDs,
			oracles[uni.chainID],
		)
	}
}
