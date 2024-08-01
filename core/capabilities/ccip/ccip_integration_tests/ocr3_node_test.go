package ccip_integration_tests

import (
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/consul/sdk/freeport"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ping_pong_demo"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/stretchr/testify/require"
)

/*
*   If you want to debug, set log level to info and use the following commands for easier logs filtering.
*
*   // Run the test and redirect logs to logs.txt
*   go test -v -run "^TestIntegration_OCR3Nodes" ./core/services/ocr3/plugins/ccip_integration_tests 2>&1 > logs.txt
*
*   // Reads logs.txt as a stream and apply filters using grep
*   tail -fn0 logs.txt | grep "CCIPExecPlugin"
 */
func TestIntegration_OCR3Nodes(t *testing.T) {
	const (
		numChains = 3 // number of chains that this test will run on
		numNodes  = 4 // number of OCR3 nodes, test assumes that every node supports every chain

		simulatedBackendBlockTime = 900 * time.Millisecond // Simulated backend blocks committing interval
		oraclesBootWaitTime       = 30 * time.Second       // Time to wait for oracles to come up (HACK)
		fChain                    = 1                      // fChain value for all the chains
		oracleLogLevel            = zapcore.ErrorLevel     // Log level for the oracle / plugins.
	)

	t.Logf("creating %d universes", numChains)
	homeChainUni, universes := createUniverses(t, numChains)

	var (
		oracles = make(map[uint64][]confighelper2.OracleIdentityExtra)
		apps    []chainlink.Application
		nodes   []*ocr3Node
		p2pIDs  [][32]byte

		// The bootstrap node will be: nodes[0]
		bootstrapPort  int
		bootstrapP2PID p2pkey.PeerID
	)

	ports := freeport.GetN(t, numNodes)
	ctx := testutils.Context(t)
	callCtx := &bind.CallOpts{Context: ctx}

	for i := 0; i < numNodes; i++ {
		t.Logf("Setting up ocr3 node:%d at port:%d", i, ports[i])
		node := setupNodeOCR3(t, ports[i], universes, homeChainUni, oracleLogLevel)

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

		apps = append(apps, node.app)
		nodes = append(nodes, node)

		peerID, err := p2pkey.MakePeerID(node.peerID)
		require.NoError(t, err)
		p2pIDs = append(p2pIDs, peerID)
	}

	bootstrapPort = ports[0]
	bootstrapP2PID = p2pIDs[0]
	bootstrapAddr := fmt.Sprintf("127.0.0.1:%d", bootstrapPort)
	t.Logf("[bootstrap node] peerID:%s p2pID:%d address:%s", nodes[0].peerID, bootstrapP2PID, bootstrapAddr)

	// Start committing periodically in the background for all the chains
	tick := time.NewTicker(simulatedBackendBlockTime)
	defer tick.Stop()
	commitBlocksBackground(t, universes, tick)

	ccipCapabilityID, err := homeChainUni.capabilityRegistry.GetHashedCapabilityId(
		callCtx, CapabilityLabelledName, CapabilityVersion)
	require.NoError(t, err, "failed to get hashed capability id for ccip")
	require.NotEqual(t, [32]byte{}, ccipCapabilityID, "ccip capability id is empty")

	// Need to Add nodes and assign capabilities to them before creating DONS
	homeChainUni.AddNodes(t, p2pIDs, [][32]byte{ccipCapabilityID})

	for _, uni := range universes {
		t.Logf("Adding chainconfig for chain %d", uni.chainID)
		AddChainConfig(t, homeChainUni, getSelector(uni.chainID), p2pIDs, fChain)
	}

	cfgs, err := homeChainUni.ccipConfig.GetAllChainConfigs(callCtx)
	require.NoError(t, err)
	t.Logf("Got all homechain configs %#v", cfgs)
	require.Len(t, cfgs, numChains)

	// Create a DON for each chain
	for _, uni := range universes {
		// Add nodes and give them the capability
		t.Log("Adding DON for universe: ", uni.chainID)
		chainSelector := getSelector(uni.chainID)
		homeChainUni.AddDON(
			t,
			ccipCapabilityID,
			chainSelector,
			uni,
			fChain,
			bootstrapP2PID,
			p2pIDs,
			oracles[uni.chainID],
		)
	}

	t.Log("Creating ocr3 jobs, starting oracles")
	for i := 0; i < len(nodes); i++ {
		err1 := nodes[i].app.Start(ctx)
		require.NoError(t, err1)
		tApp := apps[i]
		t.Cleanup(func() { require.NoError(t, tApp.Stop()) })

		jb := mustGetJobSpec(t, bootstrapP2PID, bootstrapPort, nodes[i].peerID, nodes[i].keybundle.ID())
		require.NoErrorf(t, tApp.AddJobV2(ctx, &jb), "Wasn't able to create ccip job for node %d", i)
	}

	t.Logf("Initializing PingPong contracts")
	pingPongs := initializePingPongContracts(t, universes)

	// NOTE: messageIDs are populated in the sendPingPong function
	var messageIDs = make(map[uint64]map[uint64][32]byte) // sourceChain->destChain->messageID
	var replayBlocks = make(map[uint64]uint64)            // chainID -> blocksToReplay

	t.Logf("Sending ping pong from each chain to each other")
	sendPingPong(t, universes, pingPongs, messageIDs, replayBlocks, 1)

	// Wait for the oracles to come up.
	// TODO: We need some data driven way to do this e.g. wait until LP filters to be registered.
	time.Sleep(oraclesBootWaitTime)

	// Replay the log poller on all the chains so that the logs are in the db.
	// otherwise the plugins won't pick them up.
	for _, node := range nodes {
		for chainID, replayBlock := range replayBlocks {
			t.Logf("Replaying logs for chain %d from block %d", chainID, replayBlock)
			require.NoError(t, node.app.ReplayFromBlock(big.NewInt(int64(chainID)), replayBlock, false), "failed to replay logs")
		}
	}

	// Wait for the commit reports to be generated and reported on all chains.
	numUnis := len(universes)
	var wg sync.WaitGroup
	for _, uni := range universes {
		wg.Add(1)
		go func(uni onchainUniverse) {
			defer wg.Done()
			waitForCommit(t, uni, numUnis, nil)
		}(uni)
	}

	tStart := time.Now()
	t.Log("Waiting for commit reports")
	wg.Wait()
	t.Logf("Commit reports received after %s", time.Since(tStart))

	var preRequestBlocks = make(map[uint64]uint64)
	for _, uni := range universes {
		preRequestBlocks[uni.chainID] = uni.backend.Blockchain().CurrentBlock().Number.Uint64()
	}

	t.Log("Sending ping pong from each chain to each other again for a second time")
	sendPingPong(t, universes, pingPongs, messageIDs, replayBlocks, 2)

	for _, uni := range universes {
		startBlock := preRequestBlocks[uni.chainID]
		wg.Add(1)
		go func(uni onchainUniverse, startBlock *uint64) {
			defer wg.Done()
			waitForCommit(t, uni, numUnis, startBlock)
		}(uni, &startBlock)
	}

	tStart = time.Now()
	t.Log("Waiting for second batch of commit reports")
	wg.Wait()
	t.Logf("Second batch of commit reports received after %s", time.Since(tStart))
}

func sendPingPong(t *testing.T, universes map[uint64]onchainUniverse, pingPongs map[uint64]map[uint64]*ping_pong_demo.PingPongDemo, messageIDs map[uint64]map[uint64][32]byte, replayBlocks map[uint64]uint64, expectedSeqNum uint64) {
	for chainID, uni := range universes {
		var replayBlock uint64
		for otherChain, pingPong := range pingPongs[chainID] {
			t.Log("PingPong From: ", chainID, " To: ", otherChain)

			expNextSeqNr, err1 := uni.onramp.GetExpectedNextSequenceNumber(&bind.CallOpts{}, getSelector(otherChain))
			require.NoError(t, err1)
			require.Equal(t, expectedSeqNum, expNextSeqNr, "expected next sequence number should be 1")

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
}

func waitForCommit(t *testing.T, uni onchainUniverse, numUnis int, startBlock *uint64) {
	sink := make(chan *evm_2_evm_multi_offramp.EVM2EVMMultiOffRampCommitReportAccepted)
	subscription, err := uni.offramp.WatchCommitReportAccepted(&bind.WatchOpts{
		Start:   startBlock,
		Context: testutils.Context(t),
	}, sink)
	require.NoError(t, err)

	for {
		select {
		case <-time.After(5 * time.Second):
			t.Logf("Waiting for commit report on chain id %d (selector %d)", uni.chainID, getSelector(uni.chainID))
		case subErr := <-subscription.Err():
			t.Fatalf("Subscription error: %+v", subErr)
		case report := <-sink:
			if len(report.Report.MerkleRoots) > 0 {
				if len(report.Report.MerkleRoots) == numUnis-1 {
					t.Logf("Received commit report with %d merkle roots on chain id %d (selector %d): %+v",
						len(report.Report.MerkleRoots), uni.chainID, getSelector(uni.chainID), report)
					return
				}
				t.Fatalf("Received commit report with %d merkle roots, expected %d", len(report.Report.MerkleRoots), numUnis)
			} else {
				t.Logf("Received commit report without merkle roots on chain id %d (selector %d): %+v", uni.chainID, getSelector(uni.chainID), report)
			}
		}
	}
}
