package integration_tests

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type triggerFactory = func(t *testing.T, reports []datastreams.FeedReport) commoncap.TriggerCapability
type targetFactory = func(t *testing.T, reportsSink chan commoncap.CapabilityResponse) commoncap.TargetCapability

type consensusFactory = func(t *testing.T, keys []ocr2key.KeyBundle) commoncap.ConsensusCapability

func Test_HardcodedWorkflow_DonTopologies(t *testing.T) {
	return
	ctx := testutils.Context(t)

	reportsSink := make(chan commoncap.CapabilityResponse, 1000)

	numTriggerDonPeers := 7
	triggerKeyBundles, triggerDonPeerIDs := getDonKeysAndPeers(t, numTriggerDonPeers)

	workflowKeyBundles, workflowDonPeerIDs := getDonKeysAndPeers(t, 5)

	reportCtx := ocrTypes.ReportContext{}
	rawCtx := RawReportContext(reportCtx)

	feedOneIDBytes, feedOneIDString := newFeedID(t)
	feedTwoIDBytes, feedTwoIDString := newFeedID(t)
	feedThreeIDBytes, feedThreeIDString := newFeedID(t)

	reportsRef := []*datastreams.FeedReport{
		{
			FeedID:               feedOneIDString,
			FullReport:           newReport(t, feedOneIDBytes, big.NewInt(1), 5),
			BenchmarkPrice:       big.NewInt(1).Bytes(),
			ObservationTimestamp: 5,
			Signatures:           [][]byte{},
			ReportContext:        rawCtx,
		},
		{
			FeedID:               feedThreeIDString,
			FullReport:           newReport(t, feedThreeIDBytes, big.NewInt(3), 7),
			BenchmarkPrice:       big.NewInt(2).Bytes(),
			ObservationTimestamp: 8,
			Signatures:           [][]byte{},
			ReportContext:        rawCtx,
		},
		{
			FeedID:               feedTwoIDString,
			FullReport:           newReport(t, feedTwoIDBytes, big.NewInt(2), 6),
			BenchmarkPrice:       big.NewInt(3).Bytes(),
			ObservationTimestamp: 10,
			Signatures:           [][]byte{},
			ReportContext:        rawCtx,
		},
	}

	for _, report := range reportsRef {
		var signatures [][]byte

		for _, key := range triggerKeyBundles {
			sig, err := key.Sign(reportCtx, report.FullReport)
			require.NoError(t, err)
			signatures = append(signatures, sig)
		}

		report.Signatures = signatures
	}

	workflowDonF := uint8(1)
	triggerDonF := uint8(1)
	targetDonF := uint8(1)

	var reportsCopy []datastreams.FeedReport
	for _, reportRef := range reportsRef {
		reportsCopy = append(reportsCopy, *reportRef)
	}

	triggerDonPeers := make([]p2ptypes.PeerID, len(triggerDonPeerIDs))
	for i := 0; i < len(triggerDonPeerIDs); i++ {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(triggerDonPeerIDs[i].PeerID)))
		triggerDonPeers[i] = capabilityPeerID
	}

	triggerDonInfo := commoncap.DON{
		ID:      "trigger-don",
		Members: triggerDonPeers,
		F:       triggerDonF,
	}

	targetDonPeers := make([]p2ptypes.PeerID, len(targetDonPeerIDs))
	for i := 0; i < len(targetDonPeerIDs); i++ {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(targetDonPeerIDs[i].PeerID)))
		targetDonPeers[i] = capabilityPeerID
	}

	targetDonInfo := commoncap.DON{
		ID:      "target-don",
		Members: targetDonPeers,
		F:       targetDonF,
	}

	workflowPeers := make([]p2ptypes.PeerID, len(workflowDonPeerIDs))
	for i := 0; i < len(workflowDonPeerIDs); i++ {
		workflowPeerID := p2ptypes.PeerID{}
		require.NoError(t, workflowPeerID.UnmarshalText([]byte(workflowDonPeerIDs[i].PeerID)))
		workflowPeers[i] = workflowPeerID
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      "workflow-don",
		F:       workflowDonF,
	}

	peerIdToSigner := make(map[string]string)
	for _, peer := range triggerDonPeerIDs {
		peerIdToSigner[peer.PeerID] = peer.Signer
	}
	for _, peer := range targetDonPeerIDs {
		peerIdToSigner[peer.PeerID] = peer.Signer
	}
	for _, peer := range workflowDonPeerIDs {
		peerIdToSigner[peer.PeerID] = peer.Signer
	}

	simulatedEthBlockchain, transactor := setupBlockchain(t, triggerDonPeerIDs)
	capabilitiesRegistryAddr := setupCapabilitiesRegistry(t, ctx, workflowDonPeerIDs, triggerDonPeerIDs, transactor, simulatedEthBlockchain)
	forwarderAddr := setupForwarder(t, ctx, triggerDonPeerIDs, transactor, simulatedEthBlockchain)

	workflowDonNodes, _, _ := createDons(ctx, t, []triggerFactory{mockMercuryTrigger},
		[]targetFactory{mockEthereumTestnetSepoliaTarget},
		[]consensusFactory{mockConsensus}, workflowDonInfo, triggerDonInfo, targetDonInfo,
		reportsSink, reportsCopy, simulatedEthBlockchain, capabilitiesRegistryAddr, forwarderAddr, peerIdToSigner,
		workflowKeyBundles, transactor)
	for _, node := range workflowDonNodes {
		AddWorkflowJob(t, node)
	}

	/*
		logs := make(chan ethtypes.Log, 10000)
		sub, err := simulatedEthBlockchain.SubscribeFilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(0),
			ToBlock:   big.NewInt(math.MaxInt64),
		}, logs)
		require.NoError(t, err)
		defer sub.Unsubscribe()

		for {
			select {
			case err := <-sub.Err():
				log.Fatalf("Received subscription error: %v", err)
			case vLog := <-logs:
				fmt.Println("Log: ", vLog)
			}
		} */

	reportCount := 0
	for r := range reportsSink {
		reportCount++
		fmt.Printf("received report %v\n", r)
		// Todo - -1?
		if reportCount == len(workflowDonPeerIDs)-1 {
			break
		}
	}

}

func getDonKeysAndPeers(t *testing.T, numTriggerDonPeers int) ([]ocr2key.KeyBundle, []peer) {
	var keyBundles []ocr2key.KeyBundle
	var donPeerIDs []peer
	for i := 0; i < numTriggerDonPeers; i++ {
		peerID := NewPeerID()

		keyBundle, err := ocr2key.New(chaintype.EVM)
		require.NoError(t, err)
		keyBundles = append(keyBundles, keyBundle)

		pk := keyBundle.PublicKey()

		p := peer{
			PeerID: peerID,
			Signer: fmt.Sprintf("0x%x", pk),
		}

		donPeerIDs = append(donPeerIDs, p)
	}
	return keyBundles, donPeerIDs
}

func newFeedID(t *testing.T) ([32]byte, string) {
	buf := [32]byte{}
	_, err := rand.Read(buf[:])
	require.NoError(t, err)
	return buf, "0x" + hex.EncodeToString(buf[:])
}

func newReport(t *testing.T, feedID [32]byte, price *big.Int, timestamp int64) []byte {
	v3Codec := reportcodec.NewReportCodec(feedID, logger.TestLogger(t))
	raw, err := v3Codec.BuildReport(v3.ReportFields{
		BenchmarkPrice: price,
		Timestamp:      uint32(timestamp),
		Bid:            big.NewInt(0),
		Ask:            big.NewInt(0),
		LinkFee:        big.NewInt(0),
		NativeFee:      big.NewInt(0),
	})
	require.NoError(t, err)
	return raw
}

func createDons(ctx context.Context, t *testing.T, triggerFactories []triggerFactory,
	targetFactories []targetFactory,
	consensusFactories []consensusFactory,
	workflowDon commoncap.DON,
	triggerDon commoncap.DON,
	targetDon commoncap.DON,
	reportsSink chan commoncap.CapabilityResponse,
	reports []datastreams.FeedReport,
	simulatedEthBlockchain *backends.SimulatedBackend,
	capRegistryAddr common.Address,
	forwarderAddr common.Address,
	peerIDToSigner map[string]string,
	workflowNodeKeyBundles []ocr2key.KeyBundle,
	transactor *bind.TransactOpts,
) ([]*cltest.TestApplication, []*cltest.TestApplication, []*cltest.TestApplication) {

	lggr := logger.TestLogger(t)

	broker := newTestAsyncMessageBroker(t, 1000)

	var triggerNodes []*cltest.TestApplication
	for _, triggerPeer := range triggerDon.Members {
		triggerPeerDispatcher := broker.NewDispatcherForNode(triggerPeer)
		nodeInfo := commoncap.Node{
			PeerID: &triggerPeer,
		}

		capabilityRegistry := capabilities.NewRegistry(lggr)
		for _, factory := range triggerFactories {
			trig := factory(t, reports)
			err := capabilityRegistry.Add(ctx, trig)
			require.NoError(t, err)
		}

		triggerNode := StartNewNode(ctx, t, nodeInfo, simulatedEthBlockchain, capRegistryAddr, triggerPeerDispatcher,
			testPeerWrapper{peer: testPeer{triggerPeer}}, capabilityRegistry, nil, transactor)

		require.NoError(t, triggerNode.Start(testutils.Context(t)))
		triggerNodes = append(triggerNodes, triggerNode)
	}

	var targetNodes []*cltest.TestApplication
	for _, targetPeer := range targetDon.Members {
		targetPeerDispatcher := broker.NewDispatcherForNode(targetPeer)
		nodeInfo := commoncap.Node{
			PeerID: &targetPeer,
		}

		capabilityRegistry := capabilities.NewRegistry(lggr)
		for _, factory := range targetFactories {
			targ := factory(t, reportsSink)
			err := capabilityRegistry.Add(ctx, targ)
			require.NoError(t, err)
		}

		targetNode := StartNewNode(ctx, t, nodeInfo, simulatedEthBlockchain, capRegistryAddr, targetPeerDispatcher,
			testPeerWrapper{peer: testPeer{targetPeer}}, capabilityRegistry, &forwarderAddr, transactor)

		require.NoError(t, targetNode.Start(testutils.Context(t)))
		targetNodes = append(triggerNodes, targetNode)
	}

	var workflowNodes []*cltest.TestApplication
	for _, workflowPeer := range workflowDon.Members {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeer)
		capabilityRegistry := capabilities.NewRegistry(lggr)

		for _, consensusFactory := range consensusFactories {
			consensus := consensusFactory(t, workflowNodeKeyBundles)
			err := capabilityRegistry.Add(ctx, consensus)
			require.NoError(t, err)
		}

		nodeInfo := commoncap.Node{
			PeerID:         &workflowPeer,
			WorkflowDON:    workflowDon,
			CapabilityDONs: []commoncap.DON{triggerDon, targetDon},
		}

		workflowNode := StartNewNode(ctx, t, nodeInfo, simulatedEthBlockchain, capRegistryAddr, workflowPeerDispatcher,
			testPeerWrapper{peer: testPeer{workflowPeer}}, capabilityRegistry, nil, transactor)

		require.NoError(t, workflowNode.Start(testutils.Context(t)))
		workflowNodes = append(workflowNodes, workflowNode)
	}

	servicetest.Run(t, broker)

	return workflowNodes, triggerNodes, targetNodes
}

type testPeerWrapper struct {
	peer testPeer
}

func (t testPeerWrapper) Start(ctx context.Context) error {
	return nil
}

func (t testPeerWrapper) Close() error {
	return nil
}

func (t testPeerWrapper) Ready() error {
	return nil
}

func (t testPeerWrapper) HealthReport() map[string]error {
	return nil
}

func (t testPeerWrapper) Name() string {
	return "testPeerWrapper"
}

func (t testPeerWrapper) GetPeer() p2ptypes.Peer {
	return t.peer
}

type testPeer struct {
	id p2ptypes.PeerID
}

func (t testPeer) Start(ctx context.Context) error {
	return nil
}

func (t testPeer) Close() error {
	return nil
}

func (t testPeer) Ready() error {
	return nil
}

func (t testPeer) HealthReport() map[string]error {
	return nil
}

func (t testPeer) Name() string {
	return "testPeer"
}

func (t testPeer) ID() p2ptypes.PeerID {
	return t.id
}

func (t testPeer) UpdateConnections(peers map[p2ptypes.PeerID]p2ptypes.StreamConfig) error {
	return nil
}

func (t testPeer) Send(peerID p2ptypes.PeerID, msg []byte) error {
	return nil
}

func (t testPeer) Receive() <-chan p2ptypes.Message {
	return nil
}

func StartNewNode(ctx context.Context,
	t *testing.T, nodeInfo commoncap.Node,
	backend *backends.SimulatedBackend, capRegistryAddr common.Address,
	dispatcher remotetypes.Dispatcher,
	peerWrapper p2ptypes.PeerWrapper,
	localCapabilities coretypes.CapabilitiesRegistry,
	forwarderAddress *common.Address,
	transactor *bind.TransactOpts,
) *cltest.TestApplication {

	var keyV2 *ethkey.KeyV2
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {

		c.Capabilities.ExternalRegistry.ChainID = ptr(fmt.Sprintf("%d", testutils.SimulatedChainID))
		c.Capabilities.ExternalRegistry.Address = ptr(capRegistryAddr.String())
		c.Capabilities.Peering.V2.Enabled = ptr(true)

		//here - need to get the key out to sign ? or can pass in key?

		if forwarderAddress != nil {
			newKey, err := ethkey.NewV2()
			keyV2 = &newKey
			require.NoError(t, err)
			eip55Address := types.EIP55AddressFromAddress(*forwarderAddress)
			c.EVM[0].Chain.Workflow.ForwarderAddress = &eip55Address
			c.EVM[0].Chain.Workflow.FromAddress = &keyV2.EIP55Address
		}

		c.Feature.FeedsManager = ptr(false)
	})

	if keyV2 != nil {
		n, err := backend.NonceAt(ctx, transactor.From, nil)
		require.NoError(t, err)

		tx := cltest.NewLegacyTransaction(
			n, keyV2.Address,
			assets.Ether(1).ToInt(),
			21000,
			assets.GWei(1).ToInt(),
			nil)
		signedTx, err := transactor.Signer(transactor.From, tx)
		require.NoError(t, err)
		err = backend.SendTransaction(ctx, signedTx)
		require.NoError(t, err)
		backend.Commit()

		fmt.Printf("keyV2: %v\n", keyV2)

		return cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, backend, nodeInfo,
			dispatcher, peerWrapper, localCapabilities, *keyV2)
	} else {
		return cltest.NewApplicationWithConfigV2OnSimulatedBlockchain(t, config, backend, nodeInfo,
			dispatcher, peerWrapper, localCapabilities)
	}

}

func NewPeerID() string {
	var privKey [32]byte
	_, err := rand.Read(privKey[:])
	if err != nil {
		panic(err)
	}

	peerID := append(libp2pMagic(), privKey[:]...)

	return base58.Encode(peerID[:])
}

func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}

func ptr[T any](t T) *T { return &t }
