package integration_tests

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

const (
	// As a  default set the logging to info otherwise 10s/100s of MB of logs are created on each test run
	TestLogLevel = zapcore.InfoLevel
)

var (
	workflowName    = "abcdef0123"
	workflowOwnerID = "0100000000000000000000000000000000000001"
)

type donInfo struct {
	commoncap.DON
	keys       []ethkey.KeyV2
	keyBundles []ocr2key.KeyBundle
	peerIDs    []peer
}

func setupStreamDonsWithTransmissionSchedule(ctx context.Context, t *testing.T, workflowDonInfo donInfo, triggerDonInfo donInfo, targetDonInfo donInfo,
	feedCount int, deltaStage string, schedule string) (*feeds_consumer.KeystoneFeedsConsumer, []string, *reportsSink) {
	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(TestLogLevel)

	ethBlockchain, transactor := setupBlockchain(t, 1000, 1*time.Second)
	capabilitiesRegistryAddr := setupCapabilitiesRegistryContract(ctx, t, workflowDonInfo, triggerDonInfo, targetDonInfo, transactor, ethBlockchain)
	forwarderAddr, _ := setupForwarderContract(t, workflowDonInfo, transactor, ethBlockchain)
	consumerAddr, consumer := setupConsumerContract(t, transactor, ethBlockchain, forwarderAddr, workflowOwnerID, workflowName)

	var feedIDs []string
	for i := 0; i < feedCount; i++ {
		feedIDs = append(feedIDs, newFeedID(t))
	}

	sink := newReportsSink()

	libocr := newMockLibOCR(t, workflowDonInfo.F, 1*time.Second)
	workflowDonNodes, _, _ := createDons(ctx, t, lggr, sink,
		workflowDonInfo, triggerDonInfo, targetDonInfo,
		ethBlockchain, capabilitiesRegistryAddr, forwarderAddr,
		workflowDonInfo.keyBundles, transactor, libocr)
	for _, node := range workflowDonNodes {
		addWorkflowJob(t, node, workflowName, workflowOwnerID, feedIDs, consumerAddr, deltaStage, schedule)
	}

	servicetest.Run(t, ethBlockchain)
	servicetest.Run(t, libocr)
	servicetest.Run(t, sink)
	return consumer, feedIDs, sink
}

func createDons(ctx context.Context, t *testing.T, lggr logger.Logger, reportsSink *reportsSink,
	workflowDon donInfo,
	triggerDon donInfo,
	targetDon donInfo,
	simulatedEthBlockchain *ethBackend,
	capRegistryAddr common.Address,
	forwarderAddr common.Address,
	workflowNodeKeyBundles []ocr2key.KeyBundle,
	transactor *bind.TransactOpts,
	libocr *mockLibOCR,
) ([]*cltest.TestApplication, []*cltest.TestApplication, []*cltest.TestApplication) {
	broker := newTestAsyncMessageBroker(t, 1000)

	var triggerNodes []*cltest.TestApplication
	for i, triggerPeer := range triggerDon.Members {
		triggerPeerDispatcher := broker.NewDispatcherForNode(triggerPeer)
		nodeInfo := commoncap.Node{
			PeerID: &triggerPeer,
		}

		capabilityRegistry := capabilities.NewRegistry(lggr)
		trigger := reportsSink.getNewTrigger(t)
		err := capabilityRegistry.Add(ctx, trigger)
		require.NoError(t, err)

		triggerNode := startNewNode(ctx, t, lggr.Named("Trigger-"+strconv.Itoa(i)), nodeInfo, simulatedEthBlockchain, capRegistryAddr, triggerPeerDispatcher,
			testPeerWrapper{peer: testPeer{triggerPeer}}, capabilityRegistry, nil, transactor,
			triggerDon.keys[i])

		require.NoError(t, triggerNode.Start(testutils.Context(t)))
		triggerNodes = append(triggerNodes, triggerNode)
	}

	var targetNodes []*cltest.TestApplication
	for i, targetPeer := range targetDon.Members {
		targetPeerDispatcher := broker.NewDispatcherForNode(targetPeer)
		nodeInfo := commoncap.Node{
			PeerID: &targetPeer,
		}

		capabilityRegistry := capabilities.NewRegistry(lggr)

		targetNode := startNewNode(ctx, t, lggr.Named("Target-"+strconv.Itoa(i)), nodeInfo, simulatedEthBlockchain, capRegistryAddr, targetPeerDispatcher,
			testPeerWrapper{peer: testPeer{targetPeer}}, capabilityRegistry, &forwarderAddr, transactor,
			targetDon.keys[i])

		require.NoError(t, targetNode.Start(testutils.Context(t)))
		targetNodes = append(triggerNodes, targetNode)
	}

	var workflowNodes []*cltest.TestApplication
	for i, workflowPeer := range workflowDon.Members {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeer)
		capabilityRegistry := capabilities.NewRegistry(lggr)

		requestTimeout := 10 * time.Minute
		cfg := ocr3.Config{
			Logger:            lggr,
			EncoderFactory:    evm.NewEVMEncoder,
			AggregatorFactory: capabilities.NewAggregator,
			RequestTimeout:    &requestTimeout,
		}

		ocr3Capability := ocr3.NewOCR3(cfg)
		servicetest.Run(t, ocr3Capability)

		pluginCfg := coretypes.ReportingPluginServiceConfig{}
		pluginFactory, err := ocr3Capability.NewReportingPluginFactory(ctx, pluginCfg, nil,
			nil, nil, nil, capabilityRegistry, nil, nil)
		require.NoError(t, err)

		repConfig := ocr3types.ReportingPluginConfig{
			F: int(workflowDon.F),
		}
		plugin, _, err := pluginFactory.NewReportingPlugin(repConfig)
		require.NoError(t, err)

		transmitter := ocr3.NewContractTransmitter(lggr, capabilityRegistry, "")

		libocr.AddNode(plugin, transmitter, workflowNodeKeyBundles[i])

		nodeInfo := commoncap.Node{
			PeerID:         &workflowPeer,
			WorkflowDON:    workflowDon.DON,
			CapabilityDONs: []commoncap.DON{triggerDon.DON, targetDon.DON},
		}

		workflowNode := startNewNode(ctx, t, lggr.Named("Workflow-"+strconv.Itoa(i)), nodeInfo, simulatedEthBlockchain, capRegistryAddr, workflowPeerDispatcher,
			testPeerWrapper{peer: testPeer{workflowPeer}}, capabilityRegistry, nil, transactor,
			workflowDon.keys[i])

		require.NoError(t, workflowNode.Start(testutils.Context(t)))
		workflowNodes = append(workflowNodes, workflowNode)
	}

	servicetest.Run(t, broker)

	return workflowNodes, triggerNodes, targetNodes
}

func startNewNode(ctx context.Context,
	t *testing.T, lggr logger.Logger, nodeInfo commoncap.Node,
	backend *ethBackend, capRegistryAddr common.Address,
	dispatcher remotetypes.Dispatcher,
	peerWrapper p2ptypes.PeerWrapper,
	localCapabilities *capabilities.Registry,
	forwarderAddress *common.Address,
	transactor *bind.TransactOpts,
	keyV2 ethkey.KeyV2,
) *cltest.TestApplication {
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Capabilities.ExternalRegistry.ChainID = ptr(fmt.Sprintf("%d", testutils.SimulatedChainID))
		c.Capabilities.ExternalRegistry.Address = ptr(capRegistryAddr.String())
		c.Capabilities.Peering.V2.Enabled = ptr(true)

		if forwarderAddress != nil {
			eip55Address := types.EIP55AddressFromAddress(*forwarderAddress)
			c.EVM[0].Chain.Workflow.ForwarderAddress = &eip55Address
			c.EVM[0].Chain.Workflow.FromAddress = &keyV2.EIP55Address
		}

		c.Feature.FeedsManager = ptr(false)
	})

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

	return cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, backend.SimulatedBackend, nodeInfo,
		dispatcher, peerWrapper, localCapabilities, keyV2, lggr)
}

type don struct {
	id       uint32
	numNodes int
	f        uint8
}

func createDonInfo(t *testing.T, don don) donInfo {
	keyBundles, peerIDs := getKeyBundlesAndPeerIDs(t, don.numNodes)

	donPeers := make([]p2ptypes.PeerID, len(peerIDs))
	var donKeys []ethkey.KeyV2
	for i := 0; i < len(peerIDs); i++ {
		peerID := p2ptypes.PeerID{}
		require.NoError(t, peerID.UnmarshalText([]byte(peerIDs[i].PeerID)))
		donPeers[i] = peerID
		newKey, err := ethkey.NewV2()
		require.NoError(t, err)
		donKeys = append(donKeys, newKey)
	}

	triggerDonInfo := donInfo{
		DON: commoncap.DON{
			ID:            don.id,
			Members:       donPeers,
			F:             don.f,
			ConfigVersion: 1,
		},
		peerIDs:    peerIDs,
		keys:       donKeys,
		keyBundles: keyBundles,
	}
	return triggerDonInfo
}

func createFeedReport(t *testing.T, price *big.Int, observationTimestamp int64,
	feedIDString string,
	keyBundles []ocr2key.KeyBundle) *datastreams.FeedReport {
	reportCtx := ocrTypes.ReportContext{}
	rawCtx := RawReportContext(reportCtx)

	bytes, err := hex.DecodeString(feedIDString[2:])
	require.NoError(t, err)
	var feedIDBytes [32]byte
	copy(feedIDBytes[:], bytes)

	report := &datastreams.FeedReport{
		FeedID:               feedIDString,
		FullReport:           newReport(t, feedIDBytes, price, observationTimestamp),
		BenchmarkPrice:       price.Bytes(),
		ObservationTimestamp: observationTimestamp,
		Signatures:           [][]byte{},
		ReportContext:        rawCtx,
	}

	for _, key := range keyBundles {
		sig, err := key.Sign(reportCtx, report.FullReport)
		require.NoError(t, err)
		report.Signatures = append(report.Signatures, sig)
	}

	return report
}

func getKeyBundlesAndPeerIDs(t *testing.T, numNodes int) ([]ocr2key.KeyBundle, []peer) {
	var keyBundles []ocr2key.KeyBundle
	var donPeerIDs []peer
	for i := 0; i < numNodes; i++ {
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

func newFeedID(t *testing.T) string {
	buf := [32]byte{}
	_, err := rand.Read(buf[:])
	require.NoError(t, err)
	return "0x" + hex.EncodeToString(buf[:])
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

func RawReportContext(reportCtx ocrTypes.ReportContext) []byte {
	rc := evmutil.RawReportContext(reportCtx)
	flat := []byte{}
	for _, r := range rc {
		flat = append(flat, r[:]...)
	}
	return flat
}
