package ccip_integration_tests

import (
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/hashicorp/consul/sdk/freeport"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/stretchr/testify/require"
)

const STATE_SUCCESS = uint8(2)

/*
*   If you want to debug, set log level to info and use the following commands for easier logs filtering.
*
*   // Run the test and redirect logs to logs.txt
*   go test -v -run "^TestIntegration_OCR3Nodes" ./core/capabilities/ccip/ccip_integration_tests 2>&1 > logs.txt
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
		oracleLogLevel            = zapcore.InfoLevel      // Log level for the oracle / plugins.
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

	cfgs, err := homeChainUni.ccipConfig.GetAllChainConfigs(callCtx, big.NewInt(0), big.NewInt(100))
	require.NoError(t, err)
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

	t.Logf("Sending ccip requests from each chain to all other chains")
	for _, uni := range universes {
		requests := genRequestData(uni.chainID, universes)
		uni.SendCCIPRequests(t, requests)
	}

	// Wait for the oracles to come up.
	// TODO: We need some data driven way to do this e.g. wait until LP filters to be registered.
	time.Sleep(oraclesBootWaitTime)

	// Replay the log poller on all the chains so that the logs are in the db.
	// otherwise the plugins won't pick them up.
	for _, node := range nodes {
		for chainID := range universes {
			t.Logf("Replaying logs for chain %d from block %d", chainID, 1)
			require.NoError(t, node.app.ReplayFromBlock(big.NewInt(int64(chainID)), 1, false), "failed to replay logs")
		}
	}

	// with only one request sent from each chain to each other chain,
	// and with sequence numbers on incrementing by 1 on a per-dest chain
	// basis, we expect the min sequence number to be 1 on all chains.
	expectedSeqNrRange := ccipocr3.NewSeqNumRange(1, 1)
	var wg sync.WaitGroup
	for _, uni := range universes {
		for remoteSelector := range universes {
			if remoteSelector == uni.chainID {
				continue
			}
			wg.Add(1)
			go func(uni onchainUniverse, remoteSelector uint64) {
				defer wg.Done()
				waitForCommitWithInterval(t, uni, getSelector(remoteSelector), expectedSeqNrRange)
			}(uni, remoteSelector)
		}
	}

	start := time.Now()
	wg.Wait()
	t.Logf("All chains received the expected commit report in %s", time.Since(start))

	// with only one request sent from each chain to each other chain,
	// all ExecutionStateChanged events should have the sequence number 1.
	expectedSeqNr := uint64(1)
	for _, uni := range universes {
		for remoteSelector := range universes {
			if remoteSelector == uni.chainID {
				continue
			}
			wg.Add(1)
			go func(uni onchainUniverse, remoteSelector uint64) {
				defer wg.Done()
				waitForExecWithSeqNr(t, uni, getSelector(remoteSelector), expectedSeqNr)
			}(uni, remoteSelector)
		}
	}

	start = time.Now()
	wg.Wait()
	t.Logf("All chains received the expected ExecutionStateChanged event in %s", time.Since(start))
}

func genRequestData(chainID uint64, universes map[uint64]onchainUniverse) []requestData {
	var res []requestData
	for destChainID, destUni := range universes {
		if destChainID == chainID {
			continue
		}
		res = append(res, requestData{
			destChainSelector: getSelector(destChainID),
			receiverAddress:   destUni.receiver.Address(),
			data:              []byte(fmt.Sprintf("msg from chain %d to chain %d", chainID, destChainID)),
		})
	}
	return res
}

func waitForCommitWithInterval(
	t *testing.T,
	uni onchainUniverse,
	expectedSourceChainSelector uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
) {
	sink := make(chan *offramp.OffRampCommitReportAccepted)
	subscription, err := uni.offramp.WatchCommitReportAccepted(&bind.WatchOpts{
		Context: testutils.Context(t),
	}, sink)
	require.NoError(t, err)

	for {
		select {
		case <-time.After(10 * time.Second):
			t.Logf("Waiting for commit report on chain id %d (selector %d) from source selector %d expected seq nr range %s",
				uni.chainID, getSelector(uni.chainID), expectedSourceChainSelector, expectedSeqNumRange.String())
		case subErr := <-subscription.Err():
			t.Fatalf("Subscription error: %+v", subErr)
		case report := <-sink:
			if len(report.Report.MerkleRoots) > 0 {
				// Check the interval of sequence numbers and make sure it matches
				// the expected range.
				for _, mr := range report.Report.MerkleRoots {
					if mr.SourceChainSelector == expectedSourceChainSelector &&
						uint64(expectedSeqNumRange.Start()) == mr.Interval.Min &&
						uint64(expectedSeqNumRange.End()) == mr.Interval.Max {
						t.Logf("Received commit report on chain id %d (selector %d) from source selector %d expected seq nr range %s",
							uni.chainID, getSelector(uni.chainID), expectedSourceChainSelector, expectedSeqNumRange.String())
						return
					}
				}
			}
		}
	}
}

func waitForExecWithSeqNr(t *testing.T, uni onchainUniverse, expectedSourceChainSelector, expectedSeqNr uint64) {
	for {
		scc, err := uni.offramp.GetSourceChainConfig(nil, expectedSourceChainSelector)
		require.NoError(t, err)
		t.Logf("Waiting for ExecutionStateChanged on chain %d (selector %d) from chain %d with expected sequence number %d, current onchain minSeqNr: %d",
			uni.chainID, getSelector(uni.chainID), expectedSourceChainSelector, expectedSeqNr, scc.MinSeqNr)
		iter, err := uni.offramp.FilterExecutionStateChanged(nil, []uint64{expectedSourceChainSelector}, []uint64{expectedSeqNr}, nil)
		require.NoError(t, err)
		var count int
		for iter.Next() {
			if iter.Event.SequenceNumber == expectedSeqNr && iter.Event.SourceChainSelector == expectedSourceChainSelector {
				count++
			}
		}
		if count == 1 {
			t.Logf("Received ExecutionStateChanged on chain %d (selector %d) from chain %d with expected sequence number %d",
				uni.chainID, getSelector(uni.chainID), expectedSourceChainSelector, expectedSeqNr)
			return
		}
		time.Sleep(5 * time.Second)
	}
}
