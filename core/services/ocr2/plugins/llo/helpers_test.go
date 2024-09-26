package llo_test

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/credentials"
	"github.com/smartcontractkit/wsrpc/peer"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

var _ pb.MercuryServer = &mercuryServer{}

type request struct {
	pk  credentials.StaticSizedPublicKey
	req *pb.TransmitRequest
}

func (r request) TransmitterID() ocr2types.Account {
	return ocr2types.Account(fmt.Sprintf("%x", r.pk))
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
	panic("should not be called")
}

func startMercuryServer(t *testing.T, srv *mercuryServer, pubKeys []ed25519.PublicKey) (serverURL string) {
	// Set up the wsrpc server
	lis, err := net.Listen("tcp", "127.0.0.1:0")
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

type Node struct {
	App          chainlink.Application
	ClientPubKey credentials.StaticSizedPublicKey
	KeyBundle    ocr2key.KeyBundle
	ObservedLogs *observer.ObservedLogs
}

func (node *Node) AddStreamJob(t *testing.T, spec string) (id int32) {
	job, err := streams.ValidatedStreamSpec(spec)
	require.NoError(t, err)
	err = node.App.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
	return job.ID
}

func (node *Node) DeleteJob(t *testing.T, id int32) {
	err := node.App.DeleteJob(testutils.Context(t), id)
	require.NoError(t, err)
}

func (node *Node) AddLLOJob(t *testing.T, spec string) {
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

func setupNode(
	t *testing.T,
	port int,
	dbName string,
	backend *backends.SimulatedBackend,
	csaKey csakey.KeyV2,
) (app chainlink.Application, peerID string, clientPubKey credentials.StaticSizedPublicKey, ocr2kb ocr2key.KeyBundle, observedLogs *observer.ObservedLogs) {
	k := big.NewInt(int64(port)) // keys unique to port
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(k)
	rdr := keystest.NewRandReaderFromSeed(int64(port))
	ocr2kb = ocr2key.MustNewInsecure(rdr, chaintype.EVM)

	p2paddresses := []string{fmt.Sprintf("127.0.0.1:%d", port)}

	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		// [JobPipeline]
		c.JobPipeline.MaxSuccessfulRuns = ptr(uint64(0))
		c.JobPipeline.VerboseLogging = ptr(true)

		// [Feature]
		c.Feature.UICSAKeys = ptr(true)
		c.Feature.LogPoller = ptr(true)
		c.Feature.FeedsManager = ptr(false)

		// [OCR]
		c.OCR.Enabled = ptr(false)

		// [OCR2]
		c.OCR2.Enabled = ptr(true)
		c.OCR2.ContractPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)

		// [P2P]
		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.TraceLogging = ptr(true)

		// [P2P.V2]
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.AnnounceAddresses = &p2paddresses
		c.P2P.V2.ListenAddresses = &p2paddresses
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)

		// [Mercury]
		c.Mercury.VerboseLogging = ptr(true)
	})

	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	if backend != nil {
		app = cltest.NewApplicationWithConfigV2OnSimulatedBlockchain(t, config, backend, p2pKey, ocr2kb, csaKey, lggr.Named(dbName))
	} else {
		app = cltest.NewApplicationWithConfig(t, config, p2pKey, ocr2kb, csaKey, lggr.Named(dbName))
	}
	err := app.Start(testutils.Context(t))
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, app.Stop())
	})

	return app, p2pKey.PeerID().Raw(), csaKey.StaticSizedPublicKey(), ocr2kb, observedLogs
}

func ptr[T any](t T) *T { return &t }

func addSingleDecimalStreamJob(
	t *testing.T,
	node Node,
	streamID uint32,
	bridgeName string,
) (id int32) {
	return node.AddStreamJob(t, fmt.Sprintf(`
type = "stream"
schemaVersion = 1
name = "strm-spec-%d"
streamID = %d
observationSource = """
	// Benchmark Price
	price1          [type=bridge name="%s" requestData="{\\"data\\":{\\"data\\":\\"foo\\"}}"];
	price1_parse    [type=jsonparse path="result"];

	price1 -> price1_parse;
"""

		`,
		streamID,
		streamID,
		bridgeName,
	))
}

func addQuoteStreamJob(
	t *testing.T,
	node Node,
	streamID uint32,
	benchmarkBridgeName string,
	bidBridgeName string,
	askBridgeName string,
) (id int32) {
	return node.AddStreamJob(t, fmt.Sprintf(`
type = "stream"
schemaVersion = 1
name = "strm-spec-%d"
streamID = %d
observationSource = """
	// Benchmark Price
	price1          [type=bridge name="%s" requestData="{\\"data\\":{\\"data\\":\\"foo\\"}}"];
	price1_parse    [type=jsonparse path="result" index=0];

	price1 -> price1_parse;

	// Bid
	price2          [type=bridge name="%s" requestData="{\\"data\\":{\\"data\\":\\"foo\\"}}"];
	price2_parse    [type=jsonparse path="result" index=1];

	price2 -> price2_parse;

	// Ask
	price3          [type=bridge name="%s" requestData="{\\"data\\":{\\"data\\":\\"foo\\"}}"];
	price3_parse    [type=jsonparse path="result" index=2];

	price3 -> price3_parse;
"""

		`,
		streamID,
		streamID,
		benchmarkBridgeName,
		bidBridgeName,
		askBridgeName,
	))
}
func addBootstrapJob(t *testing.T, bootstrapNode Node, configuratorAddress common.Address, name string, relayType, relayConfig string) {
	bootstrapNode.AddBootstrapJob(t, fmt.Sprintf(`
type                              = "bootstrap"
relay                             = "%s"
schemaVersion                     = 1
name                              = "boot-%s"
contractID                        = "%s"
contractConfigTrackerPollInterval = "1s"

[relayConfig]
%s
providerType = "llo"`, relayType, name, configuratorAddress.Hex(), relayConfig))
}

func addLLOJob(
	t *testing.T,
	node Node,
	configuratorAddr common.Address,
	bootstrapPeerID string,
	bootstrapNodePort int,
	clientPubKey ed25519.PublicKey,
	jobName string,
	pluginConfig,
	relayType,
	relayConfig string,
) {
	node.AddLLOJob(t, fmt.Sprintf(`
type = "offchainreporting2"
schemaVersion = 1
name = "%s"
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "%s"
contractConfigTrackerPollInterval = "1s"
ocrKeyBundleID = "%s"
p2pv2Bootstrappers = [
  "%s"
]
relay = "%s"
pluginType = "llo"
transmitterID = "%x"

[pluginConfig]
%s

[relayConfig]
%s`,
		jobName,
		configuratorAddr.Hex(),
		node.KeyBundle.ID(),
		fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort),
		relayType,
		clientPubKey,
		pluginConfig,
		relayConfig,
	))
}

func createBridge(t *testing.T, name string, i int, p decimal.Decimal, borm bridges.ORM) (bridgeName string) {
	ctx := testutils.Context(t)
	bridge := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		require.Equal(t, `{"data":{"data":"foo"}}`, string(b))

		res.WriteHeader(http.StatusOK)
		val := p.String()
		resp := fmt.Sprintf(`{"result": %s}`, val)
		_, err = res.Write([]byte(resp))
		require.NoError(t, err)
	}))
	t.Cleanup(bridge.Close)
	u, _ := url.Parse(bridge.URL)
	bridgeName = fmt.Sprintf("bridge-%s-%d", name, i)
	require.NoError(t, borm.CreateBridgeType(ctx, &bridges.BridgeType{
		Name: bridges.BridgeName(bridgeName),
		URL:  models.WebURL(*u),
	}))

	return bridgeName
}

func addOCRJobsEVMPremiumLegacy(
	t *testing.T,
	streams []Stream,
	serverPubKey ed25519.PublicKey,
	serverURL string,
	configuratorAddress common.Address,
	bootstrapPeerID string,
	bootstrapNodePort int,
	nodes []Node,
	configStoreAddress common.Address,
	clientPubKeys []ed25519.PublicKey,
	pluginConfig,
	relayType,
	relayConfig string) (jobIDs map[int]map[uint32]int32) {
	// node idx => stream id => job id
	jobIDs = make(map[int]map[uint32]int32)
	// Add OCR jobs - one per feed on each node
	for i, node := range nodes {
		if jobIDs[i] == nil {
			jobIDs[i] = make(map[uint32]int32)
		}
		for j, strm := range streams {
			// assume that streams are native, link and additionals are quote
			if j < 2 {
				var name string
				if j == 0 {
					name = "nativeprice"
				} else {
					name = "linkprice"
				}
				name = fmt.Sprintf("%s-%d-%d", name, strm.id, j)
				bmBridge := createBridge(t, name, i, strm.baseBenchmarkPrice, node.App.BridgeORM())
				jobID := addSingleDecimalStreamJob(
					t,
					node,
					strm.id,
					bmBridge,
				)
				jobIDs[i][strm.id] = jobID
			} else {
				bmBridge := createBridge(t, fmt.Sprintf("benchmarkprice-%d-%d", strm.id, j), i, strm.baseBenchmarkPrice, node.App.BridgeORM())
				bidBridge := createBridge(t, fmt.Sprintf("bid-%d-%d", strm.id, j), i, strm.baseBid, node.App.BridgeORM())
				askBridge := createBridge(t, fmt.Sprintf("ask-%d-%d", strm.id, j), i, strm.baseAsk, node.App.BridgeORM())
				jobID := addQuoteStreamJob(
					t,
					node,
					strm.id,
					bmBridge,
					bidBridge,
					askBridge,
				)
				jobIDs[i][strm.id] = jobID
			}
		}
		addLLOJob(
			t,
			node,
			configuratorAddress,
			bootstrapPeerID,
			bootstrapNodePort,
			clientPubKeys[i],
			"feed-1",
			pluginConfig,
			relayType,
			relayConfig,
		)
	}
	return jobIDs
}
