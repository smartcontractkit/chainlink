package mercury_test

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	grpc_creds "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/credentials"
	"github.com/smartcontractkit/wsrpc/peer"

	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

var _ pb.MercuryServer = &mercuryServer{}

type request struct {
	pk  string
	req *pb.TransmitRequest
}

func (r request) TransmitterID() ocr2types.Account {
	return ocr2types.Account(r.pk)
}

func (r request) PublicKeyBytes() [32]byte {
	var arr [32]byte
	copy(arr[:], r.pk)
	return arr
}

type mercuryServer struct {
	privKey     ed25519.PrivateKey
	reqsCh      chan request
	t           *testing.T
	buildReport func() []byte

	// tls fields are for grpc client testing
	tlsCertFile    *string
	tlsPrivKeyFile *string

	pb.UnimplementedMercuryServer
}

func NewMercuryServer(t *testing.T, privKey ed25519.PrivateKey, reqsCh chan request, buildReport func() []byte, tlsCertFile *string, tlsKeyFile *string) *mercuryServer {
	return &mercuryServer{privKey, reqsCh, t, buildReport, tlsCertFile, tlsKeyFile, pb.UnimplementedMercuryServer{}}
}

func (s *mercuryServer) Transmit(ctx context.Context, req *pb.TransmitRequest) (*pb.TransmitResponse, error) {
	var peerID string

	if s.tlsCertFile != nil {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("could not extract public key from grpc outgoing context")
		}
		peerID = md.Get("csa-key")[0]
		signature := md.Get("signature")[0]
		require.NoError(s.t, VerifySignature(peerID, req, signature), "signature verification failed")
	} else {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, errors.New("could not extract public key")
		}

		peerID = p.PublicKey.String()
	}

	r := request{peerID, req}
	s.reqsCh <- r

	return &pb.TransmitResponse{
		Code:  1,
		Error: "",
	}, nil
}

func (s *mercuryServer) LatestReport(ctx context.Context, lrr *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
	var peerID string

	if s.tlsCertFile != nil {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("could not extract public key from grpc outgoing context")
		}
		peerID = md.Get("csa-key")[0]
		signature := md.Get("signature")[0]
		require.NoError(s.t, VerifySignature(peerID, lrr, signature), "signature verification failed")
	} else {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, errors.New("could not extract public key")
		}
		peerID = p.PublicKey.String()
	}
	s.t.Logf("mercury server got latest report from %v for feed id 0x%x", peerID, lrr.FeedId)

	out := new(pb.LatestReportResponse)
	out.Report = new(pb.Report)
	out.Report.FeedId = lrr.FeedId

	report := s.buildReport()
	payload, err := mercury.PayloadTypes.Pack(evmutil.RawReportContext(ocrtypes.ReportContext{}), report, [][32]byte{}, [][32]byte{}, [32]byte{})
	if err != nil {
		panic(err)
	}
	out.Report.Payload = payload
	return out, nil
}

type PbRequest interface {
	String() string
}

// VerifySignature verifies the signature of the request
// TODO: Should live in a separate module
func VerifySignature(publicKeyStr string, request PbRequest, signature string) error {
	publicKeyBytes, err := hexutil.Decode("0x" + publicKeyStr)
	if err != nil {
		return err
	}

	canonicalRequestString := request.String()
	signedBytes, err := hexutil.Decode("0x" + signature)
	if err != nil {
		return err
	}
	if !ed25519.Verify(publicKeyBytes, []byte(canonicalRequestString), signedBytes) {
		return errors.New("signature verification failed")
	}
	return nil
}

// Routes start logic based on TLS attributes of srv
func startMercuryServer(t *testing.T, srv *mercuryServer, pubKeys []ed25519.PublicKey) (serverURL string) {
	if srv.tlsCertFile != nil {
		return startMercuryGrpcServer(t, srv)
	}
	return startMercuryWsrpcServer(t, srv, pubKeys)
}

// Set up the wsrpc server
func startMercuryWsrpcServer(t *testing.T, srv *mercuryServer, pubKeys []ed25519.PublicKey) (serverURL string) {
	port := freeport.GetOne(t)
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		t.Fatalf("[MAIN] failed to listen: %v", err)
	}
	serverURL = lis.Addr().String()
	s := wsrpc.NewServer(wsrpc.WithCreds(srv.privKey, pubKeys))

	// Register mercury implementation with the wsrpc server
	pb.RegisterMercuryServer(s, srv)

	// Start serving
	go s.Serve(lis)
	t.Cleanup(s.Stop)

	return
}

// Set up the grpc server
func startMercuryGrpcServer(t *testing.T, srv *mercuryServer) (serverURL string) {
	port := freeport.GetOne(t)
	serverUrl := fmt.Sprintf("127.0.0.1:%v", port)
	lis, err := net.Listen("tcp", serverUrl)
	t.Cleanup(func() { lis.Close() })
	if err != nil {
		t.Fatalf("FAIL: startMercuryGrpcServer failed to listen: %v", err)
	}
	t.Logf("serverUrl: %s, tlsCertFile: %s, tlsKeyFile: %s", serverUrl, *srv.tlsCertFile, *srv.tlsPrivKeyFile)
	serverCreds, err := grpc_creds.NewServerTLSFromFile(*srv.tlsCertFile, *srv.tlsPrivKeyFile)
	if err != nil {
		t.Fatalf("FAIL: startMercuryGrpcServer failed to create server creds: %v", err)
	}
	s := grpc.NewServer(
		grpc.Creds(
			serverCreds,
		),
	)

	// Register mercury implementation with the grpc server
	pb.RegisterGrpcMercuryServer(s, srv)

	// Start serving
	go s.Serve(lis)
	t.Cleanup(s.Stop)

	return serverUrl
}

type Feed struct {
	name               string
	id                 [32]byte
	baseBenchmarkPrice *big.Int
	baseBid            *big.Int
	baseAsk            *big.Int
	baseMarketStatus   uint32
}

func randomFeedID(version uint16) [32]byte {
	id := [32]byte(utils.NewHash())
	binary.BigEndian.PutUint16(id[:2], version)
	return id
}

type Node struct {
	App          chainlink.Application
	ClientPubKey credentials.StaticSizedPublicKey
	KeyBundle    ocr2key.KeyBundle
}

func (node *Node) AddJob(t *testing.T, spec string) {
	c := node.App.GetConfig()
	job, err := validate.ValidatedOracleSpecToml(testutils.Context(t), c.OCR2(), c.Insecure(), spec, nil)
	require.NoError(t, err)
	err = node.App.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
}

func (node *Node) AddBootstrapJob(t *testing.T, spec string) {
	job, err := ocrbootstrap.ValidatedBootstrapSpecToml(spec)
	require.NoError(t, err)
	err = node.App.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
}

// setupNode sets up a chainlink node with the given port, database name, simulated backend, and key
func setupNode(
	t *testing.T,
	port int,
	dbName string,
	backend *backends.SimulatedBackend,
	csaKey csakey.KeyV2,
	tlsCertFile string, // nil if using WSRPC, uses GRPC otherwise
) (app cltest.TestApplication, peerID string, clientPubKey credentials.StaticSizedPublicKey, ocr2kb ocr2key.KeyBundle, observedLogs *observer.ObservedLogs) {
	k := big.NewInt(int64(port)) // keys unique to port
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(k)
	rdr := keystest.NewRandReaderFromSeed(int64(port))
	ocr2kb = ocr2key.MustNewInsecure(rdr, chaintype.EVM)

	p2paddresses := []string{fmt.Sprintf("127.0.0.1:%d", port)}

	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		// [JobPipeline]
		// MaxSuccessfulRuns = 0
		c.JobPipeline.MaxSuccessfulRuns = ptr(uint64(0))
		c.JobPipeline.VerboseLogging = ptr(true)

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

		// [P2P.V2]
		// Enabled = true
		// AnnounceAddresses = ['$EXT_IP:17775']
		// ListenAddresses = ['127.0.0.1:17775']
		// DeltaDial = 500ms
		// DeltaReconcile = 5s
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.AnnounceAddresses = &p2paddresses
		c.P2P.V2.ListenAddresses = &p2paddresses
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)

		// [Mercury.TLS]
		// CertFile = '$TLS_CERT_FILE'
		c.Mercury.TLS.CertFile = &tlsCertFile
	})

	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	app = *cltest.NewApplicationWithConfigV2OnSimulatedBlockchain(t, config, backend, p2pKey, ocr2kb, csaKey, lggr.Named(dbName))
	os.Setenv("CHAINLINK_MERCURY_USE_GRPC", "true")
	err := app.Start(testutils.Context(t))
	require.NoError(t, err)

	return app, p2pKey.PeerID().Raw(), csaKey.StaticSizedPublicKey(), ocr2kb, observedLogs
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

func addV1MercuryJob(
	t *testing.T,
	node Node,
	i int,
	verifierAddress common.Address,
	bootstrapPeerID string,
	bootstrapNodePort int,
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
"""

[pluginConfig]
serverURL = "%[8]s"
serverPubKey = "%[9]x"
initialBlockNumber = %[13]d

[relayConfig]
chainID = %[12]d

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

func addV2MercuryJob(
	t *testing.T,
	node Node,
	i int,
	verifierAddress common.Address,
	bootstrapPeerID string,
	bootstrapNodePort int,
	bmBridge,
	serverURL string,
	serverPubKey,
	clientPubKey ed25519.PublicKey,
	feedName string,
	feedID [32]byte,
	linkFeedID [32]byte,
	nativeFeedID [32]byte,
) {
	node.AddJob(t, fmt.Sprintf(`
type = "offchainreporting2"
schemaVersion = 1
name = "mercury-%[1]d-%[10]s"
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "%[2]s"
feedID = "0x%[9]x"
contractConfigTrackerPollInterval = "1s"
ocrKeyBundleID = "%[3]s"
p2pv2Bootstrappers = [
  "%[4]s"
]
relay = "evm"
pluginType = "mercury"
transmitterID = "%[8]x"
observationSource = """
	// Benchmark Price
	price1          [type=bridge name="%[5]s" timeout="50ms" requestData="{\\"data\\":{\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
	price1_parse    [type=jsonparse path="result"];
	price1_multiply [type=multiply times=100000000 index=0];

	price1 -> price1_parse -> price1_multiply;
"""

[pluginConfig]
serverURL = "%[6]s"
serverPubKey = "%[7]x"
linkFeedID = "0x%[11]x"
nativeFeedID = "0x%[12]x"

[relayConfig]
chainID = 1337
		`,
		i,
		verifierAddress,
		node.KeyBundle.ID(),
		fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort),
		bmBridge,
		serverURL,
		serverPubKey,
		clientPubKey,
		feedID,
		feedName,
		linkFeedID,
		nativeFeedID,
	))
}

func addV3MercuryJob(
	t *testing.T,
	node Node,
	i int,
	verifierAddress common.Address,
	bootstrapPeerID string,
	bootstrapNodePort int,
	bmBridge,
	bidBridge,
	askBridge string,
	servers map[string]string,
	clientPubKey ed25519.PublicKey,
	feedName string,
	feedID [32]byte,
	linkFeedID [32]byte,
	nativeFeedID [32]byte,
) {
	srvs := make([]string, 0, len(servers))
	for u, k := range servers {
		srvs = append(srvs, fmt.Sprintf("%q = %q", u, k))
	}
	serversStr := fmt.Sprintf("{ %s }", strings.Join(srvs, ", "))

	node.AddJob(t, fmt.Sprintf(`
type = "offchainreporting2"
schemaVersion = 1
name = "mercury-%[1]d-%[11]s"
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "%[2]s"
feedID = "0x%[10]x"
contractConfigTrackerPollInterval = "1s"
ocrKeyBundleID = "%[3]s"
p2pv2Bootstrappers = [
  "%[4]s"
]
relay = "evm"
pluginType = "mercury"
transmitterID = "%[9]x"
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
"""

[pluginConfig]
servers = %[8]s
linkFeedID = "0x%[12]x"
nativeFeedID = "0x%[13]x"

[relayConfig]
chainID = 1337
		`,
		i,
		verifierAddress,
		node.KeyBundle.ID(),
		fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort),
		bmBridge,
		bidBridge,
		askBridge,
		serversStr,
		clientPubKey,
		feedID,
		feedName,
		linkFeedID,
		nativeFeedID,
	))
}

func addV4MercuryJob(
	t *testing.T,
	node Node,
	i int,
	verifierAddress common.Address,
	bootstrapPeerID string,
	bootstrapNodePort int,
	bmBridge,
	marketStatusBridge string,
	servers map[string]string,
	clientPubKey ed25519.PublicKey,
	feedName string,
	feedID [32]byte,
	linkFeedID [32]byte,
	nativeFeedID [32]byte,
) {
	srvs := make([]string, 0, len(servers))
	for u, k := range servers {
		srvs = append(srvs, fmt.Sprintf("%q = %q", u, k))
	}
	serversStr := fmt.Sprintf("{ %s }", strings.Join(srvs, ", "))

	node.AddJob(t, fmt.Sprintf(`
type = "offchainreporting2"
schemaVersion = 1
name = "mercury-%[1]d-%[9]s"
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "%[2]s"
feedID = "0x%[8]x"
contractConfigTrackerPollInterval = "1s"
ocrKeyBundleID = "%[3]s"
p2pv2Bootstrappers = [
  "%[4]s"
]
relay = "evm"
pluginType = "mercury"
transmitterID = "%[7]x"
observationSource = """
	// Benchmark Price
	price1          [type=bridge name="%[5]s" timeout="50ms" requestData="{\\"data\\":{\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
	price1_parse    [type=jsonparse path="result"];
	price1_multiply [type=multiply times=100000000 index=0];

	price1 -> price1_parse -> price1_multiply;

	// Market Status
	marketstatus       [type=bridge name="%[12]s" timeout="50ms" requestData="{\\"data\\":{\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
	marketstatus_parse [type=jsonparse path="result" index=1];

	marketstatus -> marketstatus_parse;
"""

[pluginConfig]
servers = %[6]s
linkFeedID = "0x%[10]x"
nativeFeedID = "0x%[11]x"

[relayConfig]
chainID = 1337
		`,
		i,
		verifierAddress,
		node.KeyBundle.ID(),
		fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort),
		bmBridge,
		serversStr,
		clientPubKey,
		feedID,
		feedName,
		linkFeedID,
		nativeFeedID,
		marketStatusBridge,
	))
}

func sign(t *testing.T, key csakey.KeyV2, request PbRequest) string {
	canonicalRequestString := request.String()
	signableKey := ed25519.PrivateKey(key.Raw())
	signedBytes, err := signableKey.Sign(nil, []byte(canonicalRequestString), crypto.Hash(0))
	require.NoError(t, err)
	return fmt.Sprintf("%x", signedBytes)
}
