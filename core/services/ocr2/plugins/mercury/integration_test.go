package mercury_test

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
	"github.com/umbracle/ethgo/abi"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/credentials"
	"github.com/smartcontractkit/wsrpc/examples/simple/keys"
	"github.com/smartcontractkit/wsrpc/peer"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ pb.MercuryServer = &mercuryServer{}

type request struct {
	pk  credentials.StaticSizedPublicKey
	req *pb.TransmitRequest
}

type mercuryServer struct {
	reqs chan request
}

func NewMercuryServer(reqs chan request) *mercuryServer {
	return &mercuryServer{reqs}
}

func (s *mercuryServer) Transmit(ctx context.Context, req *pb.TransmitRequest) (*pb.TransmitResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("could not extract public key")
	}
	s.reqs <- request{p.PublicKey, req}

	return &pb.TransmitResponse{
		Code:  1,
		Error: "",
	}, nil
}

func (s *mercuryServer) LatestReport(context.Context, *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
	return nil, nil
}

func startMercuryServer(t *testing.T, srv *mercuryServer) (url string) {
	privKey := keys.FromHex(keys.ServerPrivKey)
	pubKeys := []ed25519.PublicKey{}

	// Set up the wsrpc server
	lis, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		t.Fatalf("[MAIN] failed to listen: %v", err)
	}
	url = lis.Addr().String()
	s := wsrpc.NewServer(wsrpc.Creds(privKey, pubKeys))

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

func TestIntegration_Mercury(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Sample feed
	feedID := [32]byte(utils.NewHash())
	reqs := make(chan request)
	srv := NewMercuryServer(reqs)
	reportURL := startMercuryServer(t, srv)
	chainID := testutils.SimulatedChainID
	f := 1
	n := 4

	// Setup blockchain
	steve := testutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := core.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	stopMining := cltest.Mine(backend, 3*time.Second) // Should be greater than deltaRound since we cannot access old blocks on simulated blockchain
	t.Cleanup(stopMining)

	// Deploy config contract
	// TODO: build/import contracts and wrappers
	configContractAddr := testutils.NewAddress()

	// Setup bootstrap + oracle nodes
	bootstrapNodePort := int64(19700)
	appBootstrap, bootstrapPeerID, _, bootstrapKb := setupNode(t, bootstrapNodePort, "bootstrap_mercury", nil, backend)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}
	var (
		oracles []confighelper.OracleIdentityExtra
		nodes   []Node
	)
	// Set up n oracles
	for i := int64(0); i < int64(n); i++ {
		app, peerID, transmitter, kb := setupNode(t, bootstrapNodePort+i+1, fmt.Sprintf("oracle_keeper%d", i), []commontypes.BootstrapperLocator{
			// Supply the bootstrap IP and port as a V2 peer address
			{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
		}, backend)

		nodes = append(nodes, Node{
			app, transmitter, kb,
		})
		offchainPublicKey, _ := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  offchainPublicKey,
				TransmitAccount:   ocr2types.Account(hexutil.Encode(transmitter[:])),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}

	// Add the bootstrap job
	bootstrapNode.AddBootstrapJob(t, fmt.Sprintf(`
		type                              = "bootstrap"
		relay                             = "evm"
		schemaVersion                     = 1
		name                              = "boot"
		contractID                        = "%s"
		contractConfigTrackerPollInterval = "1s"

		[relayConfig]
		chainID = 1337
	`, configContractAddr))

	// Add OCR jobs
	for i, node := range nodes {
		node.AddJob(t, fmt.Sprintf(`
		type = "offchainreporting2"
		schemaVersion = 1
		name = "mercury-%d"
		forwardingAllowed = false
		maxTaskDuration = "1s"
		contractID = "%s"
		ocrKeyBundleID = "%s"
		p2pv2Bootstrappers = [
		  "%s"
		]
		relay = "evm"
		pluginType = "mercury"
		transmitterID = ""
		observationSource = """
			// Block Num + Hash
			block           [type=ethgetblock];
			bn_lookup       [type=lookup key="number"];
			bh_lookup       [type=lookup key="hash"];

			b1 -> bn_lookup;
			b1 -> bh_lookup;
			
			// Benchmark Price
			price1          [type=bridge name="bridge-cfbenchmarks-test" timeout="50ms" requestData="{\\"data\\":{\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
			price1_parse    [type=jsonparse path="result"];
			price1_multiply [type=multiply times=100000000];

			price1 -> price1_parse -> price1_multiply;

			// Bid
			bid          [type=bridge name="bridge-cfbenchmarks-test" timeout="50ms" requestData="{\\"data\\":{\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
			bid_parse    [type=jsonparse path="result"];
			bid_multiply [type=multiply times=100000000];

			bid -> bid_parse -> bid_multiply;

			// Ask
			ask          [type=bridge name="bridge-cfbenchmarks-test" timeout="50ms" requestData="{\\"data\\":{\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
			ask_parse    [type=jsonparse path="result"];
			ask_multiply [type=multiply times=100000000];

			ask -> ask_parse -> ask_multiply;
		"""

		[pluginConfig]
		feedID = "0x%x"
		url = "%s"

		[relayConfig]
		chainID = %d
		fromBlock = %d
		`, i, configContractAddr, node.KeyBundle.ID(), fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort), feedID, reportURL, chainID, 0))
	}

	// Setup config on contract
	configType := abi.MustNewType("tuple()")
	onchainConfig, err := abi.Encode(map[string]interface{}{}, configType)
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
		f,                    // f
		onchainConfig,
	)
	require.NoError(t, err)
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	lggr.Infow("Setting Config on Oracle Contract",
		"feedID", feedID,
		"signerAddresses", signerAddresses,
		"f", f,
		"onchainConfig", onchainConfig,
		"offchainConfigVersion", offchainConfigVersion,
		"offchainConfig", offchainConfig,
	)
	// TODO: need the verifier
	// _, err = verifier.SetConfig(
	//     steve,
	//     feedID,
	//     signerAddresses,
	//     f,
	//     onchainConfig,
	//     offchainConfigVersion,
	//     offchainConfig,
	// )
	// require.NoError(t, err)
	// backend.Commit()

	<-reqs

	panic("END")
}

func setupNode(
	t *testing.T,
	port int64,
	dbName string,
	p2pV2Bootstrappers []commontypes.BootstrapperLocator,
	backend *backends.SimulatedBackend,
) (app chainlink.Application, peerID string, clientPubKey credentials.StaticSizedPublicKey, ocr2kb ocr2key.KeyBundle) {
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(port))
	p2paddresses := []string{fmt.Sprintf("127.0.0.1:%d", port)}
	config, _ := heavyweight.FullTestDBV2(t, fmt.Sprintf("%s%d", dbName, port), func(c *chainlink.Config, s *chainlink.Secrets) {
		// [JobPipeline]
		// MaxSuccessfulRuns = 0
		c.JobPipeline.MaxSuccessfulRuns = ptr(uint64(0))

		// [Feature]
		// UICSAKeys=true
		// LogPoller = true
		c.Feature.UICSAKeys = ptr(true)
		c.Feature.LogPoller = ptr(true)

		// [OCR]
		// Enabled = false
		c.OCR.Enabled = ptr(false)

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

		// [OCR2]
		// Enabled = true
		c.OCR2.Enabled = ptr(true)
	})

	app = cltest.NewApplicationWithConfigV2OnSimulatedBlockchain(t, config, backend, p2pKey)
	ks := app.GetKeyStore()
	csaKey, err := ks.CSA().Create()
	require.NoError(t, err)
	ocr2kb, err = ks.OCR2().Create(chaintype.EVM)
	require.NoError(t, err)

	err = app.Start(testutils.Context(t))
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, app.Stop())
	})

	return app, p2pKey.PeerID().Raw(), csaKey.StaticSizedPublicKey(), ocr2kb
}

func ptr[T any](t T) *T { return &t }
