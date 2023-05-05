package mercury_test

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/credentials"
	"github.com/smartcontractkit/wsrpc/peer"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	mercury "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Feed struct {
	name               string
	id                 [32]byte
	baseBenchmarkPrice *big.Int
	baseBid            *big.Int
	baseAsk            *big.Int
}

func randomFeedID() [32]byte {
	return [32]byte(utils.NewHash())
}

func TestIntegration_Mercury(t *testing.T) {
	t.Parallel()

	// test constants
	const f = uint8(1)
	const n = 4         // number of nodes
	const fromBlock = 1 // cannot use zero, start from block 1
	const multiplier = 100000000
	testStartTimeStamp := uint32(time.Now().Unix())

	// test vars
	// pError is the probability that an EA will return an error instead of a result, as integer percentage
	// pError = 0 means it will never return error
	pError := atomic.Int64{}

	// feeds
	btcFeed := Feed{"BTC/USD", randomFeedID(), big.NewInt(20_000 * multiplier), big.NewInt(19_997 * multiplier), big.NewInt(20_004 * multiplier)}
	ethFeed := Feed{"ETH/USD", randomFeedID(), big.NewInt(1_568 * multiplier), big.NewInt(1_566 * multiplier), big.NewInt(1_569 * multiplier)}
	linkFeed := Feed{"LINK/USD", randomFeedID(), big.NewInt(7150 * multiplier / 1000), big.NewInt(7123 * multiplier / 1000), big.NewInt(7177 * multiplier / 1000)}
	feeds := []Feed{btcFeed, ethFeed, linkFeed}
	feedM := make(map[[32]byte]Feed, len(feeds))
	for i := range feeds {
		feedM[feeds[i].id] = feeds[i]
	}

	lggr := logger.TestLogger(t)

	reqs := make(chan request)
	serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	serverPubKey := serverKey.PublicKey
	srv := NewMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), reqs)
	min := big.NewInt(0)
	max := big.NewInt(math.MaxInt64)

	clientCSAKeys := make([]csakey.KeyV2, n+1)
	clientPubKeys := make([]ed25519.PublicKey, n+1)
	for i := 0; i < n+1; i++ {
		k := big.NewInt(int64(i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}
	serverURL := startMercuryServer(t, srv, clientPubKeys)
	chainID := testutils.SimulatedChainID

	// Setup blockchain
	steve := testutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := core.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	backend.Commit()                                  // ensure starting block number at least 1
	stopMining := cltest.Mine(backend, 1*time.Second) // Should be greater than deltaRound since we cannot access old blocks on simulated blockchain
	t.Cleanup(stopMining)

	// Deploy config contract
	verifierProxyAddr, _, _, err := mercury_verifier_proxy.DeployMercuryVerifierProxy(steve, backend, common.Address{}) // zero address for access controller disables access control
	require.NoError(t, err)
	verifierAddress, _, verifier, err := mercury_verifier.DeployMercuryVerifier(steve, backend, verifierProxyAddr)
	require.NoError(t, err)
	backend.Commit()

	// Setup bootstrap + oracle nodes
	bootstrapNodePort := int64(19700)
	appBootstrap, bootstrapPeerID, _, bootstrapKb := setupNode(t, bootstrapNodePort, "bootstrap_mercury", nil, backend, clientCSAKeys[n])
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}
	var (
		oracles []confighelper.OracleIdentityExtra
		nodes   []Node
	)
	// Set up n oracles

	for i := int64(0); i < int64(n); i++ {
		app, peerID, transmitter, kb := setupNode(t, bootstrapNodePort+i+1, fmt.Sprintf("oracle_mercury%d", i), []commontypes.BootstrapperLocator{
			// Supply the bootstrap IP and port as a V2 peer address
			{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
		}, backend, clientCSAKeys[i])

		nodes = append(nodes, Node{
			app, transmitter, kb,
		})
		offchainPublicKey, _ := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  offchainPublicKey,
				TransmitAccount:   ocr2types.Account(fmt.Sprintf("%x", transmitter[:])),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}

	for _, feed := range feeds {
		addBootstrapJob(t, bootstrapNode, chainID, verifierAddress, feed.name, feed.id)
	}

	createBridge := func(name string, i int, p *big.Int, borm bridges.ORM) (bridgeName string) {
		bridge := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			b, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			require.Equal(t, `{"data":{"from":"ETH","to":"USD"}}`, string(b))

			r := rand.Int63n(101)
			if r > pError.Load() {
				res.WriteHeader(http.StatusOK)
				val := decimal.NewFromBigInt(p, 0).Div(decimal.NewFromInt(multiplier)).Add(decimal.NewFromInt(int64(i)).Div(decimal.NewFromInt(100))).String()
				resp := fmt.Sprintf(`{"result": %s}`, val)
				res.Write([]byte(resp))
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				resp := fmt.Sprintf(`{"error": "pError test error"}`)
				res.Write([]byte(resp))
			}
		}))
		t.Cleanup(bridge.Close)
		u, _ := url.Parse(bridge.URL)
		bridgeName = fmt.Sprintf("bridge-%s-%d", name, i)
		require.NoError(t, borm.CreateBridgeType(&bridges.BridgeType{
			Name: bridges.BridgeName(bridgeName),
			URL:  models.WebURL(*u),
		}))

		return bridgeName
	}

	// Add OCR jobs - one per feed on each node
	for i, node := range nodes {
		for j, feed := range feeds {
			bmBridge := createBridge(fmt.Sprintf("benchmarkprice-%d", j), i, feed.baseBenchmarkPrice, node.App.BridgeORM())
			askBridge := createBridge(fmt.Sprintf("ask-%d", j), i, feed.baseAsk, node.App.BridgeORM())
			bidBridge := createBridge(fmt.Sprintf("bid-%d", j), i, feed.baseBid, node.App.BridgeORM())

			addMercuryJob(
				t,
				node,
				i,
				verifierAddress,
				bootstrapPeerID,
				bootstrapNodePort,
				bmBridge,
				bidBridge,
				askBridge,
				serverURL,
				serverPubKey,
				clientPubKeys[i],
				feed.name,
				feed.id,
				chainID,
				fromBlock,
			)
		}
	}

	// Setup config on contract
	c := relaymercury.OnchainConfig{Min: min, Max: max}
	onchainConfig, err := (relaymercury.StandardOnchainConfigCodec{}).Encode(c)
	require.NoError(t, err)

	require.NoError(t, err)
	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		2*time.Second,        // DeltaProgress
		20*time.Second,       // DeltaResend
		100*time.Millisecond, // DeltaRound
		0,                    // DeltaGrace
		1*time.Minute,        // DeltaStage
		100,                  // rMax
		[]int{len(nodes)},    // S
		oracles,
		[]byte{},             // reportingPluginConfig []byte,
		0,                    // Max duration query
		250*time.Millisecond, // Max duration observation
		250*time.Millisecond, // MaxDurationReport
		250*time.Millisecond, // MaxDurationShouldAcceptFinalizedReport
		250*time.Millisecond, // MaxDurationShouldTransmitAcceptedReport
		int(f),               // f
		onchainConfig,
	)
	require.NoError(t, err)
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)

	offchainTransmitters := make([][32]byte, n)
	for i := 0; i < n; i++ {
		offchainTransmitters[i] = nodes[i].ClientPubKey
	}

	for i, feed := range feeds {
		lggr.Infow("Setting Config on Oracle Contract",
			"i", i,
			"feedID", feed.id,
			"feedName", feed.name,
			"signerAddresses", signerAddresses,
			"offchainTransmitters", offchainTransmitters,
			"f", f,
			"onchainConfig", onchainConfig,
			"offchainConfigVersion", offchainConfigVersion,
			"offchainConfig", offchainConfig,
		)

		_, err = verifier.SetConfig(
			steve,
			feed.id,
			signerAddresses,
			offchainTransmitters,
			f,
			onchainConfig,
			offchainConfigVersion,
			offchainConfig,
		)
		require.NoError(t, err)
		backend.Commit()
	}

	// Bury it with finality depth
	ch, err := bootstrapNode.App.GetChains().EVM.Get(testutils.SimulatedChainID)
	require.NoError(t, err)
	finalityDepth := ch.Config().EvmFinalityDepth()
	for i := 0; i < int(finalityDepth); i++ {
		backend.Commit()
	}

	// Expect at least one report per feed from each oracle
	seen := make(map[[32]byte]map[credentials.StaticSizedPublicKey]struct{})
	for i := range feeds {
		// feedID will be deleted when all n oracles have reported
		seen[feeds[i].id] = make(map[credentials.StaticSizedPublicKey]struct{}, n)
	}

	for req := range reqs {
		v := make(map[string]interface{})
		err := mercury.PayloadTypes.UnpackIntoMap(v, req.req.Payload)
		require.NoError(t, err)
		report, exists := v["report"]
		if !exists {
			t.Fatalf("expected payload %#v to contain 'report'", v)
		}
		reportElems := make(map[string]interface{})
		err = reportcodec.ReportTypes.UnpackIntoMap(reportElems, report.([]byte))
		require.NoError(t, err)

		feedID := ([32]byte)(reportElems["feedId"].([32]uint8))
		feed, exists := feedM[feedID]
		require.True(t, exists)

		if _, exists := seen[feedID]; !exists {
			continue // already saw all oracles for this feed
		}

		num, err := (&reportcodec.EVMReportCodec{}).CurrentBlockNumFromReport(ocr2types.Report(report.([]byte)))
		require.NoError(t, err)
		currentBlock, err := backend.BlockByNumber(testutils.Context(t), nil)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, currentBlock.Number().Int64(), num)

		expectedBm := feed.baseBenchmarkPrice
		expectedBid := feed.baseBid
		expectedAsk := feed.baseAsk

		assert.GreaterOrEqual(t, int(reportElems["observationsTimestamp"].(uint32)), int(testStartTimeStamp))
		assert.InDelta(t, expectedBm.Int64(), reportElems["benchmarkPrice"].(*big.Int).Int64(), 5000000)
		assert.InDelta(t, expectedBid.Int64(), reportElems["bid"].(*big.Int).Int64(), 5000000)
		assert.InDelta(t, expectedAsk.Int64(), reportElems["ask"].(*big.Int).Int64(), 5000000)
		assert.GreaterOrEqual(t, int(currentBlock.Number().Int64()), int(reportElems["currentBlockNum"].(uint64)))
		assert.NotEqual(t, common.Hash{}, common.Hash(reportElems["currentBlockHash"].([32]uint8)))
		assert.LessOrEqual(t, int(reportElems["validFromBlockNum"].(uint64)), int(reportElems["currentBlockNum"].(uint64)))

		t.Logf("oracle %x reported for feed %s (0x%x)", req.pk, feed.name, feed.id)

		seen[feedID][req.pk] = struct{}{}
		if len(seen[feedID]) == n {
			t.Logf("all oracles reported for feed %x (0x%x)", feed.name, feed.id)
			delete(seen, feedID)
			if len(seen) == 0 {
				break // saw all oracles; success!
			}
		}
	}

	pError.Store(20) // 20% chance of EA error
	for i := range feeds {
		// feedID will be deleted when all n oracles have reported
		seen[feeds[i].id] = make(map[credentials.StaticSizedPublicKey]struct{}, n)
	}

	for req := range reqs {
		v := make(map[string]interface{})
		err := mercury.PayloadTypes.UnpackIntoMap(v, req.req.Payload)
		require.NoError(t, err)
		report, exists := v["report"]
		if !exists {
			t.Fatalf("expected payload %#v to contain 'report'", v)
		}
		reportElems := make(map[string]interface{})
		err = reportcodec.ReportTypes.UnpackIntoMap(reportElems, report.([]byte))
		require.NoError(t, err)

		feedID := ([32]byte)(reportElems["feedId"].([32]uint8))
		feed, exists := feedM[feedID]
		require.True(t, exists)

		if _, exists := seen[feedID]; !exists {
			continue // already saw all oracles for this feed
		}

		num, err := (&reportcodec.EVMReportCodec{}).CurrentBlockNumFromReport(ocr2types.Report(report.([]byte)))
		require.NoError(t, err)
		currentBlock, err := backend.BlockByNumber(testutils.Context(t), nil)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, currentBlock.Number().Int64(), num)

		expectedBm := feed.baseBenchmarkPrice
		expectedBid := feed.baseBid
		expectedAsk := feed.baseAsk

		assert.GreaterOrEqual(t, int(reportElems["observationsTimestamp"].(uint32)), int(testStartTimeStamp))
		assert.InDelta(t, expectedBm.Int64(), reportElems["benchmarkPrice"].(*big.Int).Int64(), 5000000)
		assert.InDelta(t, expectedBid.Int64(), reportElems["bid"].(*big.Int).Int64(), 5000000)
		assert.InDelta(t, expectedAsk.Int64(), reportElems["ask"].(*big.Int).Int64(), 5000000)
		assert.GreaterOrEqual(t, int(currentBlock.Number().Int64()), int(reportElems["currentBlockNum"].(uint64)))
		assert.NotEqual(t, common.Hash{}, common.Hash(reportElems["currentBlockHash"].([32]uint8)))
		assert.LessOrEqual(t, int(reportElems["validFromBlockNum"].(uint64)), int(reportElems["currentBlockNum"].(uint64)))

		t.Logf("oracle %x reported for feed %s (0x%x)", req.pk, feed.name, feed.id)

		seen[feedID][req.pk] = struct{}{}
		if len(seen[feedID]) == n {
			t.Logf("all oracles reported for feed %x (0x%x)", feed.name, feed.id)
			delete(seen, feedID)
			if len(seen) == 0 {
				break // saw all oracles; success!
			}
		}
	}
}

var _ pb.MercuryServer = &mercuryServer{}

type request struct {
	pk  credentials.StaticSizedPublicKey
	req *pb.TransmitRequest
}

type mercuryServer struct {
	privKey ed25519.PrivateKey
	reqsCh  chan request
	t       *testing.T
}

func NewMercuryServer(t *testing.T, privKey ed25519.PrivateKey, reqsCh chan request) *mercuryServer {
	return &mercuryServer{privKey, reqsCh, t}
}

func (s *mercuryServer) Transmit(ctx context.Context, req *pb.TransmitRequest) (*pb.TransmitResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("could not extract public key")
	}
	r := request{p.PublicKey, req}
	s.reqsCh <- r

	return &pb.TransmitResponse{
		Code:  1,
		Error: "",
	}, nil
}

func (s *mercuryServer) LatestReport(ctx context.Context, lrr *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
	// not implemented in test
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("could not extract public key")
	}
	s.t.Logf("mercury server got latest report from %x for feed id 0x%x", p.PublicKey, lrr.FeedId)
	return nil, nil
}

func startMercuryServer(t *testing.T, srv *mercuryServer, pubKeys []ed25519.PublicKey) (serverURL string) {
	// Set up the wsrpc server
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("[MAIN] failed to listen: %v", err)
	}
	serverURL = fmt.Sprintf("%s", lis.Addr().String())
	s := wsrpc.NewServer(wsrpc.Creds(srv.privKey, pubKeys))

	// Register mercury implementation with the wsrpc server
	pb.RegisterMercuryServer(s, srv)

	// Start serving
	go s.Serve(lis)
	t.Cleanup(s.Stop)

	return
}

type Node struct {
	App          chainlink.Application
	ClientPubKey credentials.StaticSizedPublicKey
	KeyBundle    ocr2key.KeyBundle
}

func (node *Node) AddJob(t *testing.T, spec string) {
	job, err := validate.ValidatedOracleSpecToml(node.App.GetConfig(), spec)
	require.NoError(t, err)
	err = node.App.AddJobV2(context.Background(), &job)
	require.NoError(t, err)
}

func (node *Node) AddBootstrapJob(t *testing.T, spec string) {
	job, err := ocrbootstrap.ValidatedBootstrapSpecToml(spec)
	require.NoError(t, err)
	err = node.App.AddJobV2(context.Background(), &job)
	require.NoError(t, err)
}

func setupNode(
	t *testing.T,
	port int64,
	dbName string,
	p2pV2Bootstrappers []commontypes.BootstrapperLocator,
	backend *backends.SimulatedBackend,
	csaKey csakey.KeyV2,
) (app chainlink.Application, peerID string, clientPubKey credentials.StaticSizedPublicKey, ocr2kb ocr2key.KeyBundle) {
	k := big.NewInt(port) // keys unique to port
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(k)
	rdr := keystest.NewRandReaderFromSeed(port)
	ocr2kb = ocr2key.MustNewInsecure(rdr, chaintype.EVM)

	p2paddresses := []string{fmt.Sprintf("127.0.0.1:%d", port)}

	config, _ := heavyweight.FullTestDBV2(t, fmt.Sprintf("%s%d", dbName, port), func(c *chainlink.Config, s *chainlink.Secrets) {
		// [JobPipeline]
		// MaxSuccessfulRuns = 0
		c.JobPipeline.MaxSuccessfulRuns = ptr(uint64(0))

		// [Feature]
		// UICSAKeys=true
		// LogPoller = true
		// FeedsManager = false
		c.Feature.UICSAKeys = ptr(true)
		c.Feature.LogPoller = ptr(true)
		c.Feature.FeedsManager = ptr(false)

		// [OCR]
		// Enabled = false
		c.OCR.Enabled = ptr(false)

		// [OCR2]
		// Enabled = true
		c.OCR2.Enabled = ptr(true)

		// [P2P]
		// PeerID = '$PEERID'
		// TraceLogging = true
		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.TraceLogging = ptr(true)

		// [P2P.V1]
		// Enabled = false
		c.P2P.V1.Enabled = ptr(false)

		// [P2P.V2]
		// Enabled = true
		// AnnounceAddresses = ['$EXT_IP:17775']
		// ListenAddresses = ['0.0.0.0:17775']
		// DeltaDial = 500ms
		// DeltaReconcile = 5s
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.AnnounceAddresses = &p2paddresses
		c.P2P.V2.ListenAddresses = &p2paddresses
		c.P2P.V2.DeltaDial = models.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = models.MustNewDuration(5 * time.Second)
	})

	app = cltest.NewApplicationWithConfigV2OnSimulatedBlockchain(t, config, backend, p2pKey, ocr2kb, csaKey, logger.TestLogger(t).Named(dbName))
	err := app.Start(testutils.Context(t))
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, app.Stop())
	})

	return app, p2pKey.PeerID().Raw(), csaKey.StaticSizedPublicKey(), ocr2kb
}

func ptr[T any](t T) *T { return &t }

func addBootstrapJob(t *testing.T, bootstrapNode Node, chainID *big.Int, verifierAddress common.Address, feedName string, feedID [32]byte) {
	bootstrapNode.AddBootstrapJob(t, fmt.Sprintf(`
type                              = "bootstrap"
relay                             = "evm"
schemaVersion                     = 1
name                              = "boot-%s"
contractID                        = "%s"
feedID 							  = "0x%x"
contractConfigTrackerPollInterval = "1s"

[relayConfig]
chainID = %d
	`, feedName, verifierAddress, feedID, chainID))
}

func addMercuryJob(
	t *testing.T,
	node Node,
	i int,
	verifierAddress common.Address,
	bootstrapPeerID string,
	bootstrapNodePort int64,
	bmBridge,
	bidBridge,
	askBridge,
	serverURL string,
	serverPubKey,
	clientPubKey ed25519.PublicKey,
	feedName string,
	feedID [32]byte,
	chainID *big.Int,
	fromBlock int,
) {
	node.AddJob(t, fmt.Sprintf(`
type = "offchainreporting2"
schemaVersion = 1
name = "mercury-%[1]d-%[14]s"
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "%[2]s"
feedID = "0x%[11]x"
contractConfigTrackerPollInterval = "1s"
ocrKeyBundleID = "%[3]s"
p2pv2Bootstrappers = [
  "%[4]s"
]
relay = "evm"
pluginType = "mercury"
transmitterID = "%[10]x"
observationSource = """
	// Benchmark Price
	price1          [type=bridge name="%[5]s" timeout="50ms" requestData="{\\"data\\":{\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
	price1_parse    [type=jsonparse path="result"];
	price1_multiply [type=multiply times=100000000 index=0];

	price1 -> price1_parse -> price1_multiply;

	// Bid
	bid          [type=bridge name="%[6]s" timeout="50ms" requestData="{\\"data\\":{\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
	bid_parse    [type=jsonparse path="result"];
	bid_multiply [type=multiply times=100000000 index=1];

	bid -> bid_parse -> bid_multiply;

	// Ask
	ask          [type=bridge name="%[7]s" timeout="50ms" requestData="{\\"data\\":{\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
	ask_parse    [type=jsonparse path="result"];
	ask_multiply [type=multiply times=100000000 index=2];

	ask -> ask_parse -> ask_multiply;

	// Block Num + Hash
	b1                 [type=ethgetblock];
	bnum_lookup        [type=lookup key="number" index=3];
	bhash_lookup       [type=lookup key="hash" index=4];

	b1 -> bnum_lookup;
	b1 -> bhash_lookup;
"""

[pluginConfig]
serverURL = "%[8]s"
serverPubKey = "%[9]x"

[relayConfig]
chainID = %[12]d
fromBlock = %[13]d
		`,
		i,
		verifierAddress,
		node.KeyBundle.ID(),
		fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort),
		bmBridge,
		bidBridge,
		askBridge,
		serverURL,
		serverPubKey,
		clientPubKey,
		feedID,
		chainID,
		fromBlock,
		feedName,
	))
}
