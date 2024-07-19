package ccip_integration_tests

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/consul/sdk/freeport"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/libocr/commontypes"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/stretchr/testify/require"
)

func TestIntegration_OCR3Nodes(t *testing.T) {
	t.Skip("Currently failing, will fix in follow-ups")

	numChains := 3
	homeChainUni, universes := createUniverses(t, numChains)
	numNodes := 4
	t.Log("creating ocr3 nodes")
	var (
		oracles = make(map[uint64][]confighelper2.OracleIdentityExtra)
		apps    []chainlink.Application
		nodes   []*ocr3Node
		p2pIDs  [][32]byte

		// The bootstrap node will be the first node (index 0)
		bootstrapPort  int
		bootstrapP2PID p2pkey.PeerID
		bootstrappers  []commontypes.BootstrapperLocator
	)

	ports := freeport.GetN(t, numNodes)
	for i := 0; i < numNodes; i++ {
		node := setupNodeOCR3(t, ports[i], bootstrappers, universes, homeChainUni)

		apps = append(apps, node.app)
		for chainID, transmitter := range node.transmitters {
			identity := confighelper2.OracleIdentityExtra{
				OracleIdentity: confighelper2.OracleIdentity{
					OnchainPublicKey:  node.keybundle.PublicKey(), // Different for each chain
					TransmitAccount:   ocrtypes.Account(transmitter.Hex()),
					OffchainPublicKey: node.keybundle.OffchainPublicKey(), // Same for each family
					PeerID:            node.peerID,
				},
				ConfigEncryptionPublicKey: node.keybundle.ConfigEncryptionPublicKey(), // Different for each chain
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
	tick := time.NewTicker(900 * time.Millisecond)
	defer tick.Stop()
	commitBlocksBackground(t, universes, tick)

	ctx := testutils.Context(t)

	ccipCapabilityID, err := homeChainUni.capabilityRegistry.GetHashedCapabilityId(&bind.CallOpts{
		Context: ctx,
	}, CapabilityLabelledName, CapabilityVersion)
	require.NoError(t, err, "failed to get hashed capability id for ccip")
	require.NotEqual(t, [32]byte{}, ccipCapabilityID, "ccip capability id is empty")

	// Need to Add nodes and assign capabilities to them before creating DONS
	homeChainUni.AddNodes(t, p2pIDs, [][32]byte{ccipCapabilityID})

	// Add homechain configs
	for _, uni := range universes {
		AddChainConfig(t, homeChainUni, getSelector(uni.chainID), p2pIDs, 1)
	}

	cfgs, err3 := homeChainUni.ccipConfig.GetAllChainConfigs(&bind.CallOpts{})
	require.NoError(t, err3)
	t.Logf("homechain_configs %+v", cfgs)
	require.Len(t, cfgs, numChains)

	// Create a DON for each chain
	for _, uni := range universes {
		// Add nodes and give them the capability
		t.Log("AddingDON for universe: ", uni.chainID)
		chainSelector := getSelector(uni.chainID)
		homeChainUni.AddDON(
			t,
			ccipCapabilityID,
			chainSelector,
			uni,
			1, // f
			bootstrapP2PID,
			p2pIDs,
			oracles[uni.chainID],
		)
	}

	t.Log("creating ocr3 jobs")
	for i := 0; i < len(nodes); i++ {
		err1 := nodes[i].app.Start(ctx)
		require.NoError(t, err1)
		tApp := apps[i]
		t.Cleanup(func() {
			require.NoError(t, tApp.Stop())
		})

		jb := mustGetJobSpec(t, bootstrapP2PID, bootstrapPort, nodes[i].peerID, nodes[i].keybundle.ID())
		require.NoErrorf(t, tApp.AddJobV2(ctx, &jb), "Wasn't able to create ccip job for node %d", i)
	}

	// sourceChain map[uint64],  destChain [32]byte
	var messageIDs = make(map[uint64]map[uint64][32]byte)
	// map[uint64] chainID, blocks
	var replayBlocks = make(map[uint64]uint64)
	pingPongs := initializePingPongContracts(t, universes)
	for chainID, uni := range universes {
		var replayBlock uint64
		for otherChain, pingPong := range pingPongs[chainID] {
			t.Log("PingPong From: ", chainID, " To: ", otherChain)

			expNextSeqNr, err1 := uni.onramp.GetExpectedNextSequenceNumber(&bind.CallOpts{}, getSelector(otherChain))
			require.NoError(t, err1)
			require.Equal(t, uint64(1), expNextSeqNr, "expected next sequence number should be 1")

			uni.backend.Commit()

			_, err2 := pingPong.StartPingPong(uni.owner)
			require.NoError(t, err2)
			uni.backend.Commit()

			endBlock := uni.backend.Blockchain().CurrentBlock().Number.Uint64()
			logIter, err3 := uni.onramp.FilterCCIPSendRequested(&bind.FilterOpts{
				Start: endBlock - 1,
				End:   &endBlock,
			}, []uint64{getSelector(otherChain)})
			require.NoError(t, err3)
			// Iterate until latest event
			var count int
			for logIter.Next() {
				count++
			}
			require.Equal(t, 1, count, "expected 1 CCIPSendRequested log only")

			log := logIter.Event
			require.Equal(t, getSelector(otherChain), log.DestChainSelector)
			require.Equal(t, pingPong.Address(), log.Message.Sender)
			chainPingPongAddr := pingPongs[otherChain][chainID].Address().Bytes()

			// Receiver address is abi-encoded if destination is EVM.
			paddedAddr := common.LeftPadBytes(chainPingPongAddr, len(log.Message.Receiver))
			require.Equal(t, paddedAddr, log.Message.Receiver)

			// check that sequence number is equal to the expected next sequence number.
			// and that the sequence number is bumped in the onramp.
			require.Equalf(t, log.Message.Header.SequenceNumber, expNextSeqNr, "incorrect sequence number in CCIPSendRequested event on chain %d", log.DestChainSelector)
			newExpNextSeqNr, err := uni.onramp.GetExpectedNextSequenceNumber(&bind.CallOpts{}, getSelector(otherChain))
			require.NoError(t, err)
			require.Equal(t, expNextSeqNr+1, newExpNextSeqNr, "expected next sequence number should be bumped by 1")

			_, ok := messageIDs[chainID]
			if !ok {
				messageIDs[chainID] = make(map[uint64][32]byte)
			}
			messageIDs[chainID][otherChain] = log.Message.Header.MessageId

			// replay block should be the earliest block that has a ccip message.
			if replayBlock == 0 {
				replayBlock = endBlock
			}
		}
		replayBlocks[chainID] = replayBlock
	}

	// HACK: wait for the oracles to come up.
	// Need some data driven way to do this.
	time.Sleep(30 * time.Second)

	// replay the log poller on all the chains so that the logs are in the db.
	// otherwise the plugins won't pick them up.
	// TODO: this is happening too early, we need to wait for the chain readers to get their config
	// and register the LP filters before this has any effect.
	for _, node := range nodes {
		for chainID, replayBlock := range replayBlocks {
			t.Logf("Replaying logs for chain %d from block %d", chainID, replayBlock)
			require.NoError(t, node.app.ReplayFromBlock(big.NewInt(int64(chainID)), replayBlock, false), "failed to replay logs")
		}
	}

	for _, uni := range universes {
		waitForCommit(t, uni)
	}
}

func waitForCommit(t *testing.T, uni onchainUniverse) {
	sink := make(chan *evm_2_evm_multi_offramp.EVM2EVMMultiOffRampCommitReportAccepted)
	subscipriton, err := uni.offramp.WatchCommitReportAccepted(&bind.WatchOpts{}, sink)
	require.NoError(t, err)

	for {
		select {
		case <-time.After(5 * time.Second):
			t.Logf("Waiting for commit report on chain id %d (selector %d)", uni.chainID, getSelector(uni.chainID))
		case subErr := <-subscipriton.Err():
			t.Fatalf("Subscription error: %+v", subErr)
		case report := <-sink:
			if len(report.Report.MerkleRoots) > 0 {
				t.Logf("Received commit report with merkle roots: %+v", report)
			} else {
				t.Logf("Received commit report without merkle roots: %+v", report)
			}
			return
		}
	}
}
