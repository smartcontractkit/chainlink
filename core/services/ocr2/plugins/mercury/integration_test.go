package mercury_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/wsrpc/credentials"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	relaycodecv1 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v1"
	relaycodecv2 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v2"
	relaycodecv3 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v3"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	token "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/fee_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/reward_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	reportcodecv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1/reportcodec"
	reportcodecv2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2/reportcodec"
	reportcodecv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var (
	f                      = uint8(1)
	n                      = 4 // number of nodes
	multiplier       int64 = 100000000
	rawOnchainConfig       = relaymercury.OnchainConfig{
		Min: big.NewInt(0),
		Max: big.NewInt(math.MaxInt64),
	}
	rawReportingPluginConfig = relaymercury.OffchainConfig{
		ExpirationWindow: 1,
		BaseUSDFee:       decimal.NewFromInt(100),
	}
)

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

func setupBlockchain(t *testing.T) (*bind.TransactOpts, *backends.SimulatedBackend, *verifier.Verifier, common.Address) {
	steve := testutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := core.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	backend.Commit()                                  // ensure starting block number at least 1
	stopMining := cltest.Mine(backend, 1*time.Second) // Should be greater than deltaRound since we cannot access old blocks on simulated blockchain
	t.Cleanup(stopMining)

	// Deploy contracts
	linkTokenAddress, _, linkToken, err := token.DeployLinkToken(steve, backend)
	require.NoError(t, err)
	_, err = linkToken.Transfer(steve, steve.From, big.NewInt(1000))
	require.NoError(t, err)
	nativeTokenAddress, _, nativeToken, err := token.DeployLinkToken(steve, backend)
	require.NoError(t, err)
	_, err = nativeToken.Transfer(steve, steve.From, big.NewInt(1000))
	require.NoError(t, err)
	verifierProxyAddr, _, verifierProxy, err := verifier_proxy.DeployVerifierProxy(steve, backend, common.Address{}) // zero address for access controller disables access control
	require.NoError(t, err)
	verifierAddress, _, verifier, err := verifier.DeployVerifier(steve, backend, verifierProxyAddr)
	require.NoError(t, err)
	_, err = verifierProxy.InitializeVerifier(steve, verifierAddress)
	require.NoError(t, err)
	rewardManagerAddr, _, rewardManager, err := reward_manager.DeployRewardManager(steve, backend, linkTokenAddress)
	require.NoError(t, err)
	feeManagerAddr, _, _, err := fee_manager.DeployFeeManager(steve, backend, linkTokenAddress, nativeTokenAddress, verifierProxyAddr, rewardManagerAddr)
	require.NoError(t, err)
	_, err = verifierProxy.SetFeeManager(steve, feeManagerAddr)
	require.NoError(t, err)
	_, err = rewardManager.SetFeeManager(steve, feeManagerAddr)
	require.NoError(t, err)
	backend.Commit()

	return steve, backend, verifier, verifierAddress
}

func TestIntegration_MercuryV1(t *testing.T) {
	t.Parallel()

	var logObservers []*observer.ObservedLogs
	t.Cleanup(func() {
		detectPanicLogs(t, logObservers)
	})
	lggr := logger.TestLogger(t)
	const fromBlock = 1 // cannot use zero, start from block 1
	testStartTimeStamp := uint32(time.Now().Unix())

	// test vars
	// pError is the probability that an EA will return an error instead of a result, as integer percentage
	// pError = 0 means it will never return error
	pError := atomic.Int64{}

	// feeds
	btcFeed := Feed{"BTC/USD", randomFeedID(1), big.NewInt(20_000 * multiplier), big.NewInt(19_997 * multiplier), big.NewInt(20_004 * multiplier)}
	ethFeed := Feed{"ETH/USD", randomFeedID(1), big.NewInt(1_568 * multiplier), big.NewInt(1_566 * multiplier), big.NewInt(1_569 * multiplier)}
	linkFeed := Feed{"LINK/USD", randomFeedID(1), big.NewInt(7150 * multiplier / 1000), big.NewInt(7123 * multiplier / 1000), big.NewInt(7177 * multiplier / 1000)}
	feeds := []Feed{btcFeed, ethFeed, linkFeed}
	feedM := make(map[[32]byte]Feed, len(feeds))
	for i := range feeds {
		feedM[feeds[i].id] = feeds[i]
	}

	reqs := make(chan request)
	serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	serverPubKey := serverKey.PublicKey
	srv := NewMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), reqs, func() []byte {
		report, err := (&reportcodecv1.ReportCodec{}).BuildReport(relaycodecv1.ReportFields{BenchmarkPrice: big.NewInt(234567), Bid: big.NewInt(1), Ask: big.NewInt(1), CurrentBlockHash: make([]byte, 32)})
		if err != nil {
			panic(err)
		}
		return report
	})
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

	steve, backend, verifier, verifierAddress := setupBlockchain(t)

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
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_mercury%d", i), backend, clientCSAKeys[i])

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
				_, err := res.Write([]byte(resp))
				require.NoError(t, err)
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				resp := `{"error": "pError test error"}`
				_, err := res.Write([]byte(resp))
				require.NoError(t, err)
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

			addV1MercuryJob(
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
	onchainConfig, err := (relaymercury.StandardOnchainConfigCodec{}).Encode(rawOnchainConfig)
	require.NoError(t, err)

	reportingPluginConfig, err := json.Marshal(rawReportingPluginConfig)
	require.NoError(t, err)

	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTestsMercuryV02(
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
			nil,
		)
		require.NoError(t, err)
		backend.Commit()
	}

	// Bury it with finality depth
	ch, err := bootstrapNode.App.GetRelayers().LegacyEVMChains().Get(testutils.SimulatedChainID.String())
	require.NoError(t, err)
	finalityDepth := ch.Config().EVM().FinalityDepth()
	for i := 0; i < int(finalityDepth); i++ {
		backend.Commit()
	}

	t.Run("receives at least one report per feed from each oracle when EAs are at 100% reliability", func(t *testing.T) {
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
			err = reportcodecv1.ReportTypes.UnpackIntoMap(reportElems, report.([]byte))
			require.NoError(t, err)

			feedID := ([32]byte)(reportElems["feedId"].([32]uint8))
			feed, exists := feedM[feedID]
			require.True(t, exists)

			if _, exists := seen[feedID]; !exists {
				continue // already saw all oracles for this feed
			}

			num, err := (&reportcodecv1.ReportCodec{}).CurrentBlockNumFromReport(ocr2types.Report(report.([]byte)))
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
			assert.GreaterOrEqual(t, currentBlock.Time(), reportElems["currentBlockTimestamp"].(uint64))
			assert.NotEqual(t, common.Hash{}, common.Hash(reportElems["currentBlockHash"].([32]uint8)))
			assert.LessOrEqual(t, int(reportElems["validFromBlockNum"].(uint64)), int(reportElems["currentBlockNum"].(uint64)))
			assert.Less(t, int64(0), int64(reportElems["validFromBlockNum"].(uint64)))

			t.Logf("oracle %x reported for feed %s (0x%x)", req.pk, feed.name, feed.id)

			seen[feedID][req.pk] = struct{}{}
			if len(seen[feedID]) == n {
				t.Logf("all oracles reported for feed %s (0x%x)", feed.name, feed.id)
				delete(seen, feedID)
				if len(seen) == 0 {
					break // saw all oracles; success!
				}
			}
		}
	})

	t.Run("receives at least one report per feed from each oracle when EAs are at 80% reliability", func(t *testing.T) {
		pError.Store(20) // 20% chance of EA error

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
			err = reportcodecv1.ReportTypes.UnpackIntoMap(reportElems, report.([]byte))
			require.NoError(t, err)

			feedID := ([32]byte)(reportElems["feedId"].([32]uint8))
			feed, exists := feedM[feedID]
			require.True(t, exists)

			if _, exists := seen[feedID]; !exists {
				continue // already saw all oracles for this feed
			}

			num, err := (&reportcodecv1.ReportCodec{}).CurrentBlockNumFromReport(ocr2types.Report(report.([]byte)))
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
			assert.GreaterOrEqual(t, currentBlock.Time(), reportElems["currentBlockTimestamp"].(uint64))
			assert.NotEqual(t, common.Hash{}, common.Hash(reportElems["currentBlockHash"].([32]uint8)))
			assert.LessOrEqual(t, int(reportElems["validFromBlockNum"].(uint64)), int(reportElems["currentBlockNum"].(uint64)))

			t.Logf("oracle %x reported for feed %s (0x%x)", req.pk, feed.name, feed.id)

			seen[feedID][req.pk] = struct{}{}
			if len(seen[feedID]) == n {
				t.Logf("all oracles reported for feed %s (0x%x)", feed.name, feed.id)
				delete(seen, feedID)
				if len(seen) == 0 {
					break // saw all oracles; success!
				}
			}
		}
	})
}

func TestIntegration_MercuryV2(t *testing.T) {
	t.Parallel()

	var logObservers []*observer.ObservedLogs
	t.Cleanup(func() {
		detectPanicLogs(t, logObservers)
	})

	testStartTimeStamp := uint32(time.Now().Unix())

	// test vars
	// pError is the probability that an EA will return an error instead of a result, as integer percentage
	// pError = 0 means it will never return error
	pError := atomic.Int64{}

	// feeds
	btcFeed := Feed{
		name:               "BTC/USD",
		id:                 randomFeedID(2),
		baseBenchmarkPrice: big.NewInt(20_000 * multiplier),
	}
	ethFeed := Feed{
		name:               "ETH/USD",
		id:                 randomFeedID(2),
		baseBenchmarkPrice: big.NewInt(1_568 * multiplier),
	}
	linkFeed := Feed{
		name:               "LINK/USD",
		id:                 randomFeedID(2),
		baseBenchmarkPrice: big.NewInt(7150 * multiplier / 1000),
	}
	feeds := []Feed{btcFeed, ethFeed, linkFeed}
	feedM := make(map[[32]byte]Feed, len(feeds))
	for i := range feeds {
		feedM[feeds[i].id] = feeds[i]
	}

	reqs := make(chan request)
	serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	serverPubKey := serverKey.PublicKey
	srv := NewMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), reqs, func() []byte {
		report, err := (&reportcodecv2.ReportCodec{}).BuildReport(relaycodecv2.ReportFields{BenchmarkPrice: big.NewInt(234567), LinkFee: big.NewInt(1), NativeFee: big.NewInt(1)})
		if err != nil {
			panic(err)
		}
		return report
	})
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

	steve, backend, verifier, verifierAddress := setupBlockchain(t)

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
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_mercury%d", i), backend, clientCSAKeys[i])

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
				_, err := res.Write([]byte(resp))
				require.NoError(t, err)
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				resp := `{"error": "pError test error"}`
				_, err := res.Write([]byte(resp))
				require.NoError(t, err)
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

			addV2MercuryJob(
				t,
				node,
				i,
				verifierAddress,
				bootstrapPeerID,
				bootstrapNodePort,
				bmBridge,
				serverURL,
				serverPubKey,
				clientPubKeys[i],
				feed.name,
				feed.id,
				randomFeedID(2),
				randomFeedID(2),
			)
		}
	}

	// Setup config on contract
	onchainConfig, err := (relaymercury.StandardOnchainConfigCodec{}).Encode(rawOnchainConfig)
	require.NoError(t, err)

	reportingPluginConfig, err := json.Marshal(rawReportingPluginConfig)
	require.NoError(t, err)

	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTestsMercuryV02(
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

	for _, feed := range feeds {
		_, err = verifier.SetConfig(
			steve,
			feed.id,
			signerAddresses,
			offchainTransmitters,
			f,
			onchainConfig,
			offchainConfigVersion,
			offchainConfig,
			nil,
		)
		require.NoError(t, err)
		backend.Commit()
	}

	// Bury it with finality depth
	ch, err := bootstrapNode.App.GetRelayers().LegacyEVMChains().Get(testutils.SimulatedChainID.String())
	require.NoError(t, err)
	finalityDepth := ch.Config().EVM().FinalityDepth()
	for i := 0; i < int(finalityDepth); i++ {
		backend.Commit()
	}

	runTestSetup := func() {
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
			err = reportcodecv2.ReportTypes.UnpackIntoMap(reportElems, report.([]byte))
			require.NoError(t, err)

			feedID := ([32]byte)(reportElems["feedId"].([32]uint8))
			feed, exists := feedM[feedID]
			require.True(t, exists)

			if _, exists := seen[feedID]; !exists {
				continue // already saw all oracles for this feed
			}

			expectedFee := relaymercury.CalculateFee(big.NewInt(234567), rawReportingPluginConfig.BaseUSDFee)
			expectedExpiresAt := reportElems["observationsTimestamp"].(uint32) + rawReportingPluginConfig.ExpirationWindow

			assert.GreaterOrEqual(t, int(reportElems["observationsTimestamp"].(uint32)), int(testStartTimeStamp))
			assert.InDelta(t, feed.baseBenchmarkPrice.Int64(), reportElems["benchmarkPrice"].(*big.Int).Int64(), 5000000)
			assert.NotZero(t, reportElems["validFromTimestamp"].(uint32))
			assert.GreaterOrEqual(t, reportElems["observationsTimestamp"].(uint32), reportElems["validFromTimestamp"].(uint32))
			assert.Equal(t, expectedExpiresAt, reportElems["expiresAt"].(uint32))
			assert.Equal(t, expectedFee, reportElems["linkFee"].(*big.Int))
			assert.Equal(t, expectedFee, reportElems["nativeFee"].(*big.Int))

			t.Logf("oracle %x reported for feed %s (0x%x)", req.pk, feed.name, feed.id)

			seen[feedID][req.pk] = struct{}{}
			if len(seen[feedID]) == n {
				t.Logf("all oracles reported for feed %s (0x%x)", feed.name, feed.id)
				delete(seen, feedID)
				if len(seen) == 0 {
					break // saw all oracles; success!
				}
			}
		}
	}

	t.Run("receives at least one report per feed from each oracle when EAs are at 100% reliability", func(t *testing.T) {
		runTestSetup()
	})

	t.Run("receives at least one report per feed from each oracle when EAs are at 80% reliability", func(t *testing.T) {
		pError.Store(20)
		runTestSetup()
	})
}

func TestIntegration_MercuryV3(t *testing.T) {
	t.Parallel()

	var logObservers []*observer.ObservedLogs
	t.Cleanup(func() {
		detectPanicLogs(t, logObservers)
	})

	testStartTimeStamp := uint32(time.Now().Unix())

	// test vars
	// pError is the probability that an EA will return an error instead of a result, as integer percentage
	// pError = 0 means it will never return error
	pError := atomic.Int64{}

	// feeds
	btcFeed := Feed{
		name:               "BTC/USD",
		id:                 randomFeedID(3),
		baseBenchmarkPrice: big.NewInt(20_000 * multiplier),
		baseBid:            big.NewInt(19_997 * multiplier),
		baseAsk:            big.NewInt(20_004 * multiplier),
	}
	ethFeed := Feed{
		name:               "ETH/USD",
		id:                 randomFeedID(3),
		baseBenchmarkPrice: big.NewInt(1_568 * multiplier),
		baseBid:            big.NewInt(1_566 * multiplier),
		baseAsk:            big.NewInt(1_569 * multiplier),
	}
	linkFeed := Feed{
		name:               "LINK/USD",
		id:                 randomFeedID(3),
		baseBenchmarkPrice: big.NewInt(7150 * multiplier / 1000),
		baseBid:            big.NewInt(7123 * multiplier / 1000),
		baseAsk:            big.NewInt(7177 * multiplier / 1000),
	}
	feeds := []Feed{btcFeed, ethFeed, linkFeed}
	feedM := make(map[[32]byte]Feed, len(feeds))
	for i := range feeds {
		feedM[feeds[i].id] = feeds[i]
	}

	reqs := make(chan request)
	serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	serverPubKey := serverKey.PublicKey
	srv := NewMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), reqs, func() []byte {
		report, err := (&reportcodecv3.ReportCodec{}).BuildReport(relaycodecv3.ReportFields{BenchmarkPrice: big.NewInt(234567), Bid: big.NewInt(1), Ask: big.NewInt(1), LinkFee: big.NewInt(1), NativeFee: big.NewInt(1)})
		if err != nil {
			panic(err)
		}
		return report
	})
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

	steve, backend, verifier, verifierAddress := setupBlockchain(t)

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
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_mercury%d", i), backend, clientCSAKeys[i])

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
				_, err := res.Write([]byte(resp))
				require.NoError(t, err)
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				resp := `{"error": "pError test error"}`
				_, err := res.Write([]byte(resp))
				require.NoError(t, err)
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
			bidBridge := createBridge(fmt.Sprintf("bid-%d", j), i, feed.baseBid, node.App.BridgeORM())
			askBridge := createBridge(fmt.Sprintf("ask-%d", j), i, feed.baseAsk, node.App.BridgeORM())

			addV3MercuryJob(
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
				randomFeedID(2),
				randomFeedID(2),
			)
		}
	}

	// Setup config on contract
	onchainConfig, err := (relaymercury.StandardOnchainConfigCodec{}).Encode(rawOnchainConfig)
	require.NoError(t, err)

	reportingPluginConfig, err := json.Marshal(rawReportingPluginConfig)
	require.NoError(t, err)

	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTestsMercuryV02(
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

	for _, feed := range feeds {
		_, err = verifier.SetConfig(
			steve,
			feed.id,
			signerAddresses,
			offchainTransmitters,
			f,
			onchainConfig,
			offchainConfigVersion,
			offchainConfig,
			nil,
		)
		require.NoError(t, err)
		backend.Commit()
	}

	// Bury it with finality depth
	ch, err := bootstrapNode.App.GetRelayers().LegacyEVMChains().Get(testutils.SimulatedChainID.String())
	require.NoError(t, err)
	finalityDepth := ch.Config().EVM().FinalityDepth()
	for i := 0; i < int(finalityDepth); i++ {
		backend.Commit()
	}

	runTestSetup := func() {
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
			err = reportcodecv3.ReportTypes.UnpackIntoMap(reportElems, report.([]byte))
			require.NoError(t, err)

			feedID := ([32]byte)(reportElems["feedId"].([32]uint8))
			feed, exists := feedM[feedID]
			require.True(t, exists)

			if _, exists := seen[feedID]; !exists {
				continue // already saw all oracles for this feed
			}

			expectedFee := relaymercury.CalculateFee(big.NewInt(234567), rawReportingPluginConfig.BaseUSDFee)
			expectedExpiresAt := reportElems["observationsTimestamp"].(uint32) + rawReportingPluginConfig.ExpirationWindow

			assert.GreaterOrEqual(t, int(reportElems["observationsTimestamp"].(uint32)), int(testStartTimeStamp))
			assert.InDelta(t, feed.baseBenchmarkPrice.Int64(), reportElems["benchmarkPrice"].(*big.Int).Int64(), 5000000)
			assert.InDelta(t, feed.baseBid.Int64(), reportElems["bid"].(*big.Int).Int64(), 5000000)
			assert.InDelta(t, feed.baseAsk.Int64(), reportElems["ask"].(*big.Int).Int64(), 5000000)
			assert.NotZero(t, reportElems["validFromTimestamp"].(uint32))
			assert.GreaterOrEqual(t, reportElems["observationsTimestamp"].(uint32), reportElems["validFromTimestamp"].(uint32))
			assert.Equal(t, expectedExpiresAt, reportElems["expiresAt"].(uint32))
			assert.Equal(t, expectedFee, reportElems["linkFee"].(*big.Int))
			assert.Equal(t, expectedFee, reportElems["nativeFee"].(*big.Int))

			t.Logf("oracle %x reported for feed %s (0x%x)", req.pk, feed.name, feed.id)

			seen[feedID][req.pk] = struct{}{}
			if len(seen[feedID]) == n {
				t.Logf("all oracles reported for feed %s (0x%x)", feed.name, feed.id)
				delete(seen, feedID)
				if len(seen) == 0 {
					break // saw all oracles; success!
				}
			}
		}
	}

	t.Run("receives at least one report per feed from each oracle when EAs are at 100% reliability", func(t *testing.T) {
		runTestSetup()
	})

	t.Run("receives at least one report per feed from each oracle when EAs are at 80% reliability", func(t *testing.T) {
		pError.Store(20)
		runTestSetup()
	})
}
