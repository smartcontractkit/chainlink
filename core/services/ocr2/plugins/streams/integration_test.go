package streams_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink-data-streams/streams"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocrconfigurationstoreevmsimple"
	"github.com/smartcontractkit/wsrpc/credentials"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var (
	f                = uint8(1)
	n                = 4 // number of nodes
	multiplier int64 = 100000000
)

func setupBlockchain(t *testing.T) (*bind.TransactOpts, *backends.SimulatedBackend, *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimple, common.Address) {
	steve := testutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := core.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	backend.Commit()                                  // ensure starting block number at least 1
	stopMining := cltest.Mine(backend, 1*time.Second) // Should be greater than deltaRound since we cannot access old blocks on simulated blockchain
	t.Cleanup(stopMining)

	// Deploy contracts
	configAddress, _, configContract, err := ocrconfigurationstoreevmsimple.DeployOCRConfigurationStoreEVMSimple(steve, backend)
	require.NoError(t, err)

	// linkTokenAddress, _, linkToken, err := link_token_interface.DeployLinkToken(steve, backend)
	// require.NoError(t, err)
	// _, err = linkToken.Transfer(steve, steve.From, big.NewInt(1000))
	// require.NoError(t, err)
	// nativeTokenAddress, _, nativeToken, err := link_token_interface.DeployLinkToken(steve, backend)
	// require.NoError(t, err)
	// _, err = nativeToken.Transfer(steve, steve.From, big.NewInt(1000))
	// require.NoError(t, err)
	// verifierProxyAddr, _, verifierProxy, err := verifier_proxy.DeployVerifierProxy(steve, backend, common.Address{}) // zero address for access controller disables access control
	// require.NoError(t, err)
	// verifierAddress, _, verifier, err := verifier.DeployVerifier(steve, backend, verifierProxyAddr)
	// require.NoError(t, err)
	// _, err = verifierProxy.InitializeVerifier(steve, verifierAddress)
	// require.NoError(t, err)
	// rewardManagerAddr, _, rewardManager, err := reward_manager.DeployRewardManager(steve, backend, linkTokenAddress)
	// require.NoError(t, err)
	// feeManagerAddr, _, _, err := fee_manager.DeployFeeManager(steve, backend, linkTokenAddress, nativeTokenAddress, verifierProxyAddr, rewardManagerAddr)
	// require.NoError(t, err)
	// _, err = verifierProxy.SetFeeManager(steve, feeManagerAddr)
	// require.NoError(t, err)
	// _, err = rewardManager.SetFeeManager(steve, feeManagerAddr)
	// require.NoError(t, err)

	backend.Commit()

	return steve, backend, configContract, configAddress
}

func detectPanicLogs(t *testing.T, logObservers []*observer.ObservedLogs) {
	var panicLines []string
	for _, observedLogs := range logObservers {
		panicLogs := observedLogs.Filter(func(e observer.LoggedEntry) bool {
			return e.Level >= zapcore.DPanicLevel
		})
		for _, log := range panicLogs.All() {
			line := fmt.Sprintf("%v\t%s\t%s\t%s\t%s", log.Time.Format(time.RFC3339), log.Level.CapitalString(), log.LoggerName, log.Caller.TrimmedPath(), log.Message)
			panicLines = append(panicLines, line)
		}
	}
	if len(panicLines) > 0 {
		t.Errorf("Found logs with DPANIC or higher level:\n%s", strings.Join(panicLines, "\n"))
	}
}

// type mercuryServer struct {
//     privKey ed25519.PrivateKey
//     reqsCh  chan request
//     t       *testing.T
// }

// func NewMercuryServer(t *testing.T, privKey ed25519.PrivateKey, reqsCh chan request) *mercuryServer {
//     return &mercuryServer{privKey, reqsCh, t}
// }

// func (s *mercuryServer) Transmit(ctx context.Context, req *pb.TransmitRequest) (*pb.TransmitResponse, error) {
//     p, ok := peer.FromContext(ctx)
//     if !ok {
//         return nil, errors.New("could not extract public key")
//     }
//     r := request{p.PublicKey, req}
//     s.reqsCh <- r

//     return &pb.TransmitResponse{
//         Code:  1,
//         Error: "",
//     }, nil
// }

// func (s *mercuryServer) LatestReport(ctx context.Context, lrr *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
//     panic("not needed for llo")
// }

func TestIntegration_Streams(t *testing.T) {
	// TODO:

	t.Parallel()

	var logObservers []*observer.ObservedLogs
	t.Cleanup(func() {
		detectPanicLogs(t, logObservers)
	})
	const fromBlock = 1 // cannot use zero, start from block 1
	// testStartTimeStamp := uint32(time.Now().Unix())

	reqs := make(chan request)
	serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	serverPubKey := serverKey.PublicKey
	srv := NewMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), reqs, nil)
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

	steve, backend, configContract, configAddress := setupBlockchain(t)
	// TODO

	// Setup bootstrap + oracle nodes
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, observedLogs := setupNode(t, bootstrapNodePort, "bootstrap_mercury", backend, clientCSAKeys[n])
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}
	logObservers = append(logObservers, observedLogs)

	// Set up n oracles
	var (
		oracles []confighelper.OracleIdentityExtra
		nodes   []Node
	)
	ports := freeport.GetN(t, n)
	for i := 0; i < n; i++ {
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_streams_%d", i), backend, clientCSAKeys[i])

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
		logObservers = append(logObservers, observedLogs)
	}

	addBootstrapJob(t, bootstrapNode, chainID, configAddress, "job-1")

	// createBridge := func(name string, i int, p *big.Int, borm bridges.ORM) (bridgeName string) {
	//     bridge := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	//         b, err := io.ReadAll(req.Body)
	//         require.NoError(t, err)
	//         require.Equal(t, `{"data":{"from":"ETH","to":"USD"}}`, string(b))

	//         res.WriteHeader(http.StatusOK)
	//         val := decimal.NewFromBigInt(p, 0).Div(decimal.NewFromInt(multiplier)).Add(decimal.NewFromInt(int64(i)).Div(decimal.NewFromInt(100))).String()
	//         resp := fmt.Sprintf(`{"result": %s}`, val)
	//         _, err = res.Write([]byte(resp))
	//         require.NoError(t, err)
	//     }))
	//     t.Cleanup(bridge.Close)
	//     u, _ := url.Parse(bridge.URL)
	//     bridgeName = fmt.Sprintf("bridge-%s-%d", name, i)
	//     require.NoError(t, borm.CreateBridgeType(&bridges.BridgeType{
	//         Name: bridges.BridgeName(bridgeName),
	//         URL:  models.WebURL(*u),
	//     }))

	//     return bridgeName
	// }

	// Add OCR jobs - one per feed on each node
	for i, node := range nodes {
		addStreamsJob(
			t,
			node,
			configAddress,
			bootstrapPeerID,
			bootstrapNodePort,
			serverURL,
			serverPubKey,
			clientPubKeys[i],
			"feed-1",
			chainID,
			fromBlock,
		)
	}

	// Setup config on contract
	rawOnchainConfig := streams.OnchainConfig{}
	// TODO: Move away from JSON
	onchainConfig, err := (&streams.JSONOnchainConfigCodec{}).Encode(rawOnchainConfig)
	require.NoError(t, err)

	rawReportingPluginConfig := streams.OffchainConfig{}
	reportingPluginConfig, err := rawReportingPluginConfig.Encode()
	require.NoError(t, err)

	signers, _, _, onchainConfig, _, _, err := ocr3confighelper.ContractSetConfigArgsForTestsMercuryV02(
		2*time.Second,        // DeltaProgress
		20*time.Second,       // DeltaResend
		400*time.Millisecond, // DeltaInitial
		100*time.Millisecond, // DeltaRound
		0,                    // DeltaGrace
		300*time.Millisecond, // DeltaCertifiedCommitRequest
		1*time.Minute,        // DeltaStage
		100,                  // rMax
		[]int{len(nodes)},    // S
		oracles,
		reportingPluginConfig, // reportingPluginConfig []byte,
		250*time.Millisecond,  // Max duration observation
		int(f),                // f
		onchainConfig,
	)

	require.NoError(t, err)
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)

	offchainTransmitters := make([][32]byte, n)
	for i := 0; i < n; i++ {
		offchainTransmitters[i] = nodes[i].ClientPubKey
	}

	cfg := ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleConfigurationEVMSimple{
		signerAddresses,
		offchainTransmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
		nil,
	}
	_, err = configContract.AddConfig(
		steve,
		cfg,
	)
	require.NoError(t, err)
	backend.Commit()

	// Bury it with finality depth
	ch, err := bootstrapNode.App.GetRelayers().LegacyEVMChains().Get(testutils.SimulatedChainID.String())
	require.NoError(t, err)
	finalityDepth := ch.Config().EVM().FinalityDepth()
	for i := 0; i < int(finalityDepth); i++ {
		backend.Commit()
	}

	t.Run("receives at least one report per feed from each oracle when EAs are at 100% reliability", func(t *testing.T) {
		// Expect at least one report per feed from each oracle
		seen := make(map[credentials.StaticSizedPublicKey]struct{})

		for req := range reqs {
			v := make(map[string]interface{})
			err := mercury.PayloadTypes.UnpackIntoMap(v, req.req.Payload)
			require.NoError(t, err)
			report, exists := v["report"]
			if !exists {
				t.Fatalf("expected payload %#v to contain 'report'", v)
			}

			assert.Equal(t, "foo", report)

			seen[req.pk] = struct{}{}
			if len(seen) == n {
				t.Logf("all oracles reported")
				break // saw all oracles; success!
			}
		}
	})
}

func addStreamsJob(
	t *testing.T,
	node Node,
	verifierAddress common.Address,
	bootstrapPeerID string,
	bootstrapNodePort int,
	serverURL string,
	serverPubKey,
	clientPubKey ed25519.PublicKey,
	jobName string,
	chainID *big.Int,
	fromBlock int,
) {
	node.AddJob(t, fmt.Sprintf(`
type = "offchainreporting2"
schemaVersion = 1
name = "%[1]s"
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "%[2]s"
contractConfigTrackerPollInterval = "1s"
ocrKeyBundleID = "%[3]s"
p2pv2Bootstrappers = [
  "%[4]s"
]
relay = "evm"
pluginType = "streams"
transmitterID = "%[5]x"

[pluginConfig]
serverURL = "%[6]s"
serverPubKey = "%[7]x"
fromBlock = %[8]d

[relayConfig]
chainID = %[9]d

		`,
		jobName,
		verifierAddress,
		node.KeyBundle.ID(),
		fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort),
		clientPubKey,
		serverURL,
		serverPubKey,
		fromBlock,
		chainID,
	))
}
