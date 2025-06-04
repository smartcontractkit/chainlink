package llo_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocrconfigurationstoreevmsimple"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/wsrpc/credentials"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	lloevm "github.com/smartcontractkit/chainlink-data-streams/llo/reportcodecs/evm"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/destination_verifier"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/destination_verifier_proxy"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	reportcodecv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
	mercuryverifier "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/verifier"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

var (
	fNodes = uint8(1)
	nNodes = 4 // number of nodes (not including bootstrap)
)

// TODO(gg) see also:
// https://github.com/smartcontractkit/mercury-pipeline/blob/9f0bc5d457d57d5807122446cb936306ecf1b263/e2e_tests/mercuryhelpers/helpers.go#L308 for example of onchain config

func setupBlockchain(t *testing.T) (
	*bind.TransactOpts,
	evmtypes.Backend,
	*destination_verifier.DestinationVerifier,
	*channel_config_store.ChannelConfigStore,
	common.Address,
	*verifier.Verifier,
	common.Address,
) {
	steve := evmtestutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := gethtypes.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	backend.Commit()
	backend.Commit() // ensure starting block number at least 1

	// Configurator
	_, _, _, err := configurator.DeployConfigurator(steve, backend.Client())
	require.NoError(t, err)
	backend.Commit()

	// DestinationVerifierProxy
	destinationVerifierProxyAddr, _, verifierProxy, err := destination_verifier_proxy.DeployDestinationVerifierProxy(steve, backend.Client())
	require.NoError(t, err)
	backend.Commit()
	// DestinationVerifier
	destinationVerifierAddr, _, destinationVerifier, err := destination_verifier.DeployDestinationVerifier(steve, backend.Client(), destinationVerifierProxyAddr)
	require.NoError(t, err)
	backend.Commit()
	// AddVerifier
	_, err = verifierProxy.SetVerifier(steve, destinationVerifierAddr)
	require.NoError(t, err)
	backend.Commit()

	// Legacy mercury verifier
	legacyVerifier, legacyVerifierAddr, _, _ := setupLegacyMercuryVerifier(t, steve, backend)

	// ChannelConfigStore
	configStoreAddress, _, configStore, err := channel_config_store.DeployChannelConfigStore(steve, backend.Client())
	require.NoError(t, err)

	backend.Commit()

	return steve, backend, destinationVerifier, configStore, configStoreAddress, legacyVerifier, legacyVerifierAddr
}

func setupLegacyMercuryVerifier(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend) (*verifier.Verifier, common.Address, *verifier_proxy.VerifierProxy, common.Address) {
	linkTokenAddress, _, linkToken, err := link_token_interface.DeployLinkToken(steve, backend.Client())
	require.NoError(t, err)
	backend.Commit()
	_, err = linkToken.Transfer(steve, steve.From, big.NewInt(1000))
	require.NoError(t, err)
	backend.Commit()
	nativeTokenAddress, _, nativeToken, err := link_token_interface.DeployLinkToken(steve, backend.Client())
	require.NoError(t, err)
	backend.Commit()
	_, err = nativeToken.Transfer(steve, steve.From, big.NewInt(1000))
	require.NoError(t, err)
	backend.Commit()
	verifierProxyAddr, _, verifierProxy, err := verifier_proxy.DeployVerifierProxy(steve, backend.Client(), common.Address{}) // zero address for access controller disables access control
	require.NoError(t, err)
	backend.Commit()
	verifierAddress, _, verifier, err := verifier.DeployVerifier(steve, backend.Client(), verifierProxyAddr)
	require.NoError(t, err)
	backend.Commit()
	_, err = verifierProxy.InitializeVerifier(steve, verifierAddress)
	require.NoError(t, err)
	backend.Commit()
	rewardManagerAddr, _, rewardManager, err := reward_manager.DeployRewardManager(steve, backend.Client(), linkTokenAddress)
	require.NoError(t, err)
	backend.Commit()
	feeManagerAddr, _, _, err := fee_manager.DeployFeeManager(steve, backend.Client(), linkTokenAddress, nativeTokenAddress, verifierProxyAddr, rewardManagerAddr)
	require.NoError(t, err)
	backend.Commit()
	_, err = verifierProxy.SetFeeManager(steve, feeManagerAddr)
	require.NoError(t, err)
	backend.Commit()
	_, err = rewardManager.SetFeeManager(steve, feeManagerAddr)
	require.NoError(t, err)
	backend.Commit()
	return verifier, verifierAddress, verifierProxy, verifierProxyAddr
}

type Stream struct {
	id                 uint32
	baseBenchmarkPrice decimal.Decimal
	baseBid            decimal.Decimal
	baseAsk            decimal.Decimal
}

const (
	ethStreamID    = 52
	linkStreamID   = 53
	quoteStreamID1 = 55
	quoteStreamID2 = 56
)

var (
	quoteStreamFeedID1 = common.HexToHash(`0x0003111111111111111111111111111111111111111111111111111111111111`)
	quoteStreamFeedID2 = common.HexToHash(`0x0003222222222222222222222222222222222222222222222222222222222222`)
	ethStream          = Stream{
		id:                 ethStreamID,
		baseBenchmarkPrice: decimal.NewFromFloat32(2_976.39),
	}
	linkStream = Stream{
		id:                 linkStreamID,
		baseBenchmarkPrice: decimal.NewFromFloat32(13.25),
	}
	quoteStream1 = Stream{
		id:                 quoteStreamID1,
		baseBenchmarkPrice: decimal.NewFromFloat32(1000.1212),
		baseBid:            decimal.NewFromFloat32(998.5431),
		baseAsk:            decimal.NewFromFloat32(1001.6999),
	}
	quoteStream2 = Stream{
		id:                 quoteStreamID2,
		baseBenchmarkPrice: decimal.NewFromFloat32(500.1212),
		baseBid:            decimal.NewFromFloat32(499.5431),
		baseAsk:            decimal.NewFromFloat32(502.6999),
	}
)

// see: https://github.com/smartcontractkit/offchain-reporting/blob/master/lib/offchainreporting2plus/internal/config/ocr3config/public_config.go
type OCRConfig struct {
	DeltaProgress                           time.Duration
	DeltaResend                             time.Duration
	DeltaInitial                            time.Duration
	DeltaRound                              time.Duration
	DeltaGrace                              time.Duration
	DeltaCertifiedCommitRequest             time.Duration
	DeltaStage                              time.Duration
	RMax                                    uint64
	S                                       []int
	Oracles                                 []confighelper.OracleIdentityExtra
	ReportingPluginConfig                   []byte
	MaxDurationInitialization               *time.Duration
	MaxDurationQuery                        time.Duration
	MaxDurationObservation                  time.Duration
	MaxDurationShouldAcceptAttestedReport   time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration
	F                                       int
	OnchainConfig                           []byte
}

func makeDefaultOCRConfig() *OCRConfig {
	defaultOnchainConfig, err := (&datastreamsllo.EVMOnchainConfigCodec{}).Encode(datastreamsllo.OnchainConfig{
		Version:                 1,
		PredecessorConfigDigest: nil,
	})
	if err != nil {
		panic(err)
	}
	return &OCRConfig{
		DeltaProgress:                           2 * time.Second,
		DeltaResend:                             20 * time.Second,
		DeltaInitial:                            400 * time.Millisecond,
		DeltaRound:                              500 * time.Millisecond,
		DeltaGrace:                              250 * time.Millisecond,
		DeltaCertifiedCommitRequest:             300 * time.Millisecond,
		DeltaStage:                              1 * time.Minute,
		RMax:                                    100,
		ReportingPluginConfig:                   []byte{},
		MaxDurationInitialization:               nil,
		MaxDurationQuery:                        0,
		MaxDurationObservation:                  250 * time.Millisecond,
		MaxDurationShouldAcceptAttestedReport:   0,
		MaxDurationShouldTransmitAcceptedReport: 0,
		F:                                       int(fNodes),
		OnchainConfig:                           defaultOnchainConfig,
	}
}

func withOffchainConfig(offchainConfig datastreamsllo.OffchainConfig) OCRConfigOption {
	return func(cfg *OCRConfig) {
		offchainConfigEncoded, err := offchainConfig.Encode()
		if err != nil {
			panic(err)
		}
		cfg.ReportingPluginConfig = offchainConfigEncoded
	}
}

func withOracles(oracles []confighelper.OracleIdentityExtra) OCRConfigOption {
	return func(cfg *OCRConfig) {
		cfg.Oracles = oracles
		cfg.S = []int{len(oracles)} // all oracles transmit by default
	}
}

type OCRConfigOption func(*OCRConfig)

func generateConfig(t *testing.T, opts ...OCRConfigOption) (signers []types.OnchainPublicKey, transmitters []types.Account, f uint8, outOnchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) {
	cfg := makeDefaultOCRConfig()

	for _, opt := range opts {
		opt(cfg)
	}
	var err error
	signers, transmitters, f, outOnchainConfig, offchainConfigVersion, offchainConfig, err = ocr3confighelper.ContractSetConfigArgsForTests(
		cfg.DeltaProgress,
		cfg.DeltaResend,
		cfg.DeltaInitial,
		cfg.DeltaRound,
		cfg.DeltaGrace,
		cfg.DeltaCertifiedCommitRequest,
		cfg.DeltaStage,
		cfg.RMax,
		cfg.S,
		cfg.Oracles,
		cfg.ReportingPluginConfig,
		cfg.MaxDurationInitialization,
		cfg.MaxDurationQuery,
		cfg.MaxDurationObservation,
		cfg.MaxDurationShouldAcceptAttestedReport,
		cfg.MaxDurationShouldTransmitAcceptedReport,
		cfg.F,
		cfg.OnchainConfig,
	)

	require.NoError(t, err)

	return
}

func setLegacyConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, legacyVerifier *verifier.Verifier, legacyVerifierAddr common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra, inOffchainConfig datastreamsllo.OffchainConfig) ocr2types.ConfigDigest {
	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig := generateConfig(t, withOracles(oracles), withOffchainConfig(inOffchainConfig))

	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	offchainTransmitters := make([][32]byte, nNodes)
	for i := 0; i < nNodes; i++ {
		offchainTransmitters[i] = nodes[i].ClientPubKey
	}
	donIDPadded := llo.DonIDToBytes32(donID)
	_, err = legacyVerifier.SetConfig(steve, donIDPadded, signerAddresses, offchainTransmitters, fNodes, onchainConfig, offchainConfigVersion, offchainConfig, nil)
	require.NoError(t, err)

	// libocr requires a few confirmations to accept the config
	backend.Commit()
	backend.Commit()
	backend.Commit()
	backend.Commit()

	l, err := legacyVerifier.LatestConfigDigestAndEpoch(&bind.CallOpts{}, donIDPadded)
	require.NoError(t, err)

	return l.ConfigDigest
}

func TestIntegration_LLO_evm_premium_legacy(t *testing.T) {
	offchainConfig := datastreamsllo.OffchainConfig{ProtocolVersion: 0, DefaultMinReportIntervalNanoseconds: 0}
	testStartTimeStamp := time.Now()
	multiplier := decimal.New(1, 18)
	expirationWindow := time.Hour / time.Second

	const salt = 100

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)
	for i := 0; i < nNodes; i++ {
		k := big.NewInt(int64(salt + i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	steve, backend, verifier, configStore, configStoreAddress, legacyVerifier, legacyVerifierAddr := setupBlockchain(t)
	t.Logf("configStoreAddress: %s", configStoreAddress.Hex())
	fromBlock, err := backend.Client().BlockNumber(testutils.Context(t))
	require.NoError(t, err)

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo", backend, bootstrapCSAKey, nil)
	t.Logf("bootstrapPeerID: %s", bootstrapPeerID)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	reqs := make(chan wsrpcRequest, 100000)
	serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
	serverPubKey := serverKey.PublicKey
	t.Logf("serverPubKey: %s", hex.EncodeToString(serverPubKey[:]))
	srv := NewWSRPCMercuryServer(t, serverKey, reqs)

	serverURL := startWSRPCMercuryServer(t, srv, clientPubKeys)
	t.Logf("serverURL: %s", serverURL)

	donID := uint32(995544)
	streams := []Stream{ethStream, linkStream, quoteStream1, quoteStream2}
	streamMap := make(map[uint32]Stream)
	for _, strm := range streams {
		streamMap[strm.id] = strm
	}

	// Setup oracle nodes
	oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
		c.Mercury.Transmitter.Protocol = ptr(config.MercuryTransmitterProtocolWSRPC)

		// TODO(gg): something like this + extra config
		// c.Feature.SecureMint.Enabled = true
	})

	chainID := testutils.SimulatedChainID
	relayType := "evm"
	relayConfig := fmt.Sprintf(`
chainID = "%s"
fromBlock = %d
lloDonID = %d
lloConfigMode = "mercury"
`, chainID, fromBlock, donID)
	addBootstrapJob(t, bootstrapNode, legacyVerifierAddr, "job-2", relayType, relayConfig)

	// Channel definitions
	channelDefinitions := llotypes.ChannelDefinitions{
		1: {
			ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
			Streams: []llotypes.Stream{
				{
					StreamID:   ethStreamID,
					Aggregator: llotypes.AggregatorMedian,
				},
				{
					StreamID:   linkStreamID,
					Aggregator: llotypes.AggregatorMedian,
				},
				{
					StreamID:   quoteStreamID1,
					Aggregator: llotypes.AggregatorQuote,
				},
			},
			Opts: llotypes.ChannelOpts([]byte(fmt.Sprintf(`{"baseUSDFee":"0.1","expirationWindow":%d,"feedId":"0x%x","multiplier":"%s"}`, expirationWindow, quoteStreamFeedID1, multiplier.String()))),
		},
		2: {
			ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
			Streams: []llotypes.Stream{
				{
					StreamID:   ethStreamID,
					Aggregator: llotypes.AggregatorMedian,
				},
				{
					StreamID:   linkStreamID,
					Aggregator: llotypes.AggregatorMedian,
				},
				{
					StreamID:   quoteStreamID2,
					Aggregator: llotypes.AggregatorQuote,
				},
			},
			Opts: llotypes.ChannelOpts([]byte(fmt.Sprintf(`{"baseUSDFee":"0.1","expirationWindow":%d,"feedId":"0x%x","multiplier":"%s"}`, expirationWindow, quoteStreamFeedID2, multiplier.String()))),
		},
	}

	url, sha := newChannelDefinitionsServer(t, channelDefinitions)

	// Set channel definitions
	_, err = configStore.SetChannelDefinitions(steve, donID, url, sha)
	require.NoError(t, err)
	backend.Commit()

	pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
	donID = %d
	channelDefinitionsContractAddress = "0x%x"
	channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock)
	addOCRJobsEVMPremiumLegacy(t, streams, serverPubKey, serverURL, legacyVerifierAddr, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

	allowedSenders := make([]common.Address, len(nodes))
	for i, node := range nodes {
		keys, err := node.App.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		allowedSenders[i] = keys[0].Address // assuming the first key is the transmitter
	}

	aggregatorAddress := setSecureMintOnchainConfigUsingAggregator(t, steve, backend, nodes, oracles)

	ocrConfigStoreAddress, ocrConfigStore := setSecureMintOnchainConfigUsingEvmSimpleConfig(t, steve, backend, nodes, oracles)
	t.Logf("Deployed and configured OCRConfigStore contract at: %s", ocrConfigStoreAddress.Hex())
	ds, err := ocrConfigStore.TypeAndVersion(&bind.CallOpts{})
	require.NoError(t, err)
	t.Logf("OCRConfigStore description: %s", ds)

	// TODO(gg): enable this for writing step
	// TODO(gg): deduplicate
	// feedIDBytes := [16]byte{}
	// copy(feedIDBytes[:], common.FromHex("0xA1B2C3D4E5F600010203040506070809"))

	// dfCacheAddress, dfCacheContract := setupDataFeedsCacheContract(t, steve, backend, allowedSenders, steve.From.Hex(), "securemint")
	// t.Logf("Deployed and configured DataFeedsCache contract at: %s", dfCacheAddress.Hex())
	// desc, err := dfCacheContract.GetDescription(&bind.CallOpts{}, feedIDBytes)
	// require.NoError(t, err)
	// t.Logf("DataFeedsCache description: %s", desc)

	// setSecureMintOnchainConfig(t, steve, backend, nodes, oracles, dfCacheAddress, dfCache)

	// configDetails, err := ocrContract.LatestConfigDetails(&bind.CallOpts{})
	// require.NoError(t, err)
	// t.Logf("configDetails: %+v", configDetails)

	// latestConfigDigestAndEpoch, err := ocrContract.LatestConfigDigestAndEpoch(&bind.CallOpts{})
	// require.NoError(t, err)
	// t.Logf("latestConfigDigestAndEpoch: %+v", latestConfigDigestAndEpoch)

	jobIDs := addSecureMintOCRJobs(t, nodes, aggregatorAddress)

	t.Logf("Configuring contract again")
	configureIt(t, ocrConfigStore, steve, backend, nodes, oracles)
	t.Logf("Configured contract again")

	t.Logf("jobIDs: %v", jobIDs)
	validateJobsRunningSuccessfully(t, nodes, jobIDs)

	// Set config on configurator
	setLegacyConfig(
		t, donID, steve, backend, legacyVerifier, legacyVerifierAddr, nodes, oracles, offchainConfig,
	)

	// Set config on the destination verifier
	signerAddresses := make([]common.Address, len(oracles))
	for i, oracle := range oracles {
		signerAddresses[i] = common.BytesToAddress(oracle.OracleIdentity.OnchainPublicKey)
	}
	{
		recipientAddressesAndWeights := []destination_verifier.CommonAddressAndWeight{}

		_, err := verifier.SetConfig(steve, signerAddresses, fNodes, recipientAddressesAndWeights)
		require.NoError(t, err)
		backend.Commit()
	}

	// Expect at least one report per feed from each oracle
	seen := make(map[[32]byte]map[credentials.StaticSizedPublicKey]struct{})
	for _, cd := range channelDefinitions {
		var opts lloevm.ReportFormatEVMPremiumLegacyOpts
		err := json.Unmarshal(cd.Opts, &opts)
		require.NoError(t, err)
		// feedID will be deleted when all n oracles have reported
		seen[opts.FeedID] = make(map[credentials.StaticSizedPublicKey]struct{}, nNodes)
	}
	for req := range reqs {
		assert.Equal(t, uint32(llotypes.ReportFormatEVMPremiumLegacy), req.req.ReportFormat)
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

		feedID := reportElems["feedId"].([32]uint8)

		if _, exists := seen[feedID]; !exists {
			continue // already saw all oracles for this feed
		}

		var expectedBm, expectedBid, expectedAsk *big.Int
		if feedID == quoteStreamFeedID1 {
			expectedBm = quoteStream1.baseBenchmarkPrice.Mul(multiplier).BigInt()
			expectedBid = quoteStream1.baseBid.Mul(multiplier).BigInt()
			expectedAsk = quoteStream1.baseAsk.Mul(multiplier).BigInt()
		} else if feedID == quoteStreamFeedID2 {
			expectedBm = quoteStream2.baseBenchmarkPrice.Mul(multiplier).BigInt()
			expectedBid = quoteStream2.baseBid.Mul(multiplier).BigInt()
			expectedAsk = quoteStream2.baseAsk.Mul(multiplier).BigInt()
		} else {
			t.Fatalf("unrecognized feedID: 0x%x", feedID)
		}

		assert.GreaterOrEqual(t, reportElems["validFromTimestamp"].(uint32), uint32(testStartTimeStamp.Unix()))
		assert.GreaterOrEqual(t, int(reportElems["observationsTimestamp"].(uint32)), int(testStartTimeStamp.Unix()))
		assert.Equal(t, "33597747607000", reportElems["nativeFee"].(*big.Int).String())
		assert.Equal(t, "7547169811320755", reportElems["linkFee"].(*big.Int).String())
		assert.Equal(t, reportElems["observationsTimestamp"].(uint32)+uint32(expirationWindow), reportElems["expiresAt"].(uint32))
		assert.Equal(t, expectedBm.String(), reportElems["benchmarkPrice"].(*big.Int).String())
		assert.Equal(t, expectedBid.String(), reportElems["bid"].(*big.Int).String())
		assert.Equal(t, expectedAsk.String(), reportElems["ask"].(*big.Int).String())

		// emulate mercury server verifying report (local verification)
		{
			rv := mercuryverifier.NewVerifier()

			reportSigners, err := rv.Verify(mercuryverifier.SignedReport{
				RawRs:         v["rawRs"].([][32]byte),
				RawSs:         v["rawSs"].([][32]byte),
				RawVs:         v["rawVs"].([32]byte),
				ReportContext: v["reportContext"].([3][32]byte),
				Report:        v["report"].([]byte),
			}, fNodes, signerAddresses)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(reportSigners), int(fNodes+1))
			assert.Subset(t, signerAddresses, reportSigners)
		}

		seen[feedID][req.pk] = struct{}{}
		if len(seen[feedID]) == nNodes {
			delete(seen, feedID)
			if len(seen) == 0 {
				break // saw all oracles; success!
			}
		}
	}
}

func setupNodes(t *testing.T, nNodes int, backend evmtypes.Backend, clientCSAKeys []csakey.KeyV2, f func(*chainlink.Config)) (oracles []confighelper.OracleIdentityExtra, nodes []Node) {
	ports := freeport.GetN(t, nNodes)
	for i := 0; i < nNodes; i++ {
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_streams_%d", i), backend, clientCSAKeys[i], f)

		nodes = append(nodes, Node{
			App:          app,
			ClientPubKey: transmitter,
			KeyBundle:    kb,
			ObservedLogs: observedLogs,
		})
		offchainPublicKey, err := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
		require.NoError(t, err)
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  offchainPublicKey,
				TransmitAccount:   ocr2types.Account(hex.EncodeToString(transmitter[:])),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}
	return
}

func newChannelDefinitionsServer(t *testing.T, channelDefinitions llotypes.ChannelDefinitions) (url string, sha [32]byte) {
	channelDefinitionsJSON, err := json.MarshalIndent(channelDefinitions, "", "  ")
	require.NoError(t, err)
	channelDefinitionsSHA := sha3.Sum256(channelDefinitionsJSON)

	// Set up channel definitions server
	channelDefinitionsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(channelDefinitionsJSON)
		require.NoError(t, err)
	}))
	t.Cleanup(channelDefinitionsServer.Close)
	return channelDefinitionsServer.URL, channelDefinitionsSHA
}

func validateJobsRunningSuccessfully(t *testing.T, nodes []Node, jobIDs map[int]int32) {

	// 1. Assert no job spec errors
	for i, node := range nodes {
		jobs, _, err := node.App.JobORM().FindJobs(testutils.Context(t), 0, 1000)
		require.NoErrorf(t, err, "assert error finding jobs for node %d", i)
		t.Logf("%d jobs found for node %d", len(jobs), i)
		for _, j := range jobs {
			t.Logf("job %d on node %d oracle spec: %#v", j.ID, i, j.OCR2OracleSpec)
			t.Logf("job %d on node %d pipeline spec: %#v", j.ID, i, j.PipelineSpec)
		}
		// No spec errors
		for _, j := range jobs {
			ignore := 0
			for _, jse := range j.JobSpecErrors {
				// Non-fatal timing related error, ignore for testing.
				if strings.Contains(jse.Description, "leader's phase conflicts tGrace timeout") {
					ignore++
				} else {
					t.Errorf("assert error: job spec error on node %d: %v", i, jse)
				}
			}
			require.Lenf(t, j.JobSpecErrors, ignore, "assert error: job spec errors on node %d", i)
		}
	}

	t.Logf("No job spec errors identified for any node")

	// 2. Assert that all the Secure Mint jobs get a run with valid values eventually
	// var wg sync.WaitGroup
	// for i, node := range nodes {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		// t.Logf("finding pipeline runs for job %d on node %d", jobIDs[i], i)
	// 		// completedRuns, err := node.App.JobORM().FindPipelineRunIDsByJobID(testutils.Context(t), jobIDs[i], 0, 10)
	// 		// if !assert.NoError(t, err) {
	// 		// 	t.Logf("assert error finding pipeline runs for job %d: %v", jobIDs[i], err)
	// 		// 	return
	// 		// }
	// 		// t.Logf("found pipeline runs for job %d on node %d: %v", jobIDs[i], i, completedRuns)

	// 		// Want at least 2 runs so we see all the metadata.

	// 		pr := cltest.WaitForPipelineComplete(t, i, jobIDs[i], 1, 4, node.App.JobORM(), 30*time.Second, 1*time.Second)
	// 		jb, err := pr[0].Outputs.MarshalJSON()
	// 		if !assert.NoError(t, err) {
	// 			t.Logf("assert error marshalling outputs for job %d: %v", jobIDs[i], err)
	// 			return
	// 		}
	// 		assert.Equalf(t, []byte(fmt.Sprintf("[\"%d\"]", 1000*i)), jb, "pr[0] %+v pr[1] %+v", pr[0], pr[1], "assert error: something unexpected happened")
	// 	}()
	// }
	// t.Logf("waiting for pipeline runs to complete")
	// wg.Wait()
}

// func setSecureMintOnchainConfig(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, nodes []Node, oracles []confighelper.OracleIdentityExtra, dfCacheAddress common.Address, dfCacheContract *data_feeds_cache.DataFeedsCache) [32]byte {

// 	minAnswer, maxAnswer := new(big.Int), new(big.Int)
// 	minAnswer.Exp(big.NewInt(-2), big.NewInt(191), nil)
// 	maxAnswer.Exp(big.NewInt(2), big.NewInt(191), nil)
// 	maxAnswer.Sub(maxAnswer, big.NewInt(1))

// 	// TODO(gg): this uses the median codec, not sure if this is correct
// 	// onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(minAnswer, maxAnswer)
// 	// require.NoError(t, err)

// 	// TODO(gg): use DF Cache onchain conifg
// 	onchainConfig := por.PorOffchainConfig{} // TODO(gg): set config values
// 	onchainConfigBytes, err := onchainConfig.Serialize()
// 	require.NoError(t, err)

// 	signers, transmitters, f, outOnchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
// 		2*time.Second,        // deltaProgress,
// 		20*time.Second,       // deltaResend,
// 		400*time.Millisecond, // deltaInitial,
// 		500*time.Millisecond, // deltaRound,
// 		250*time.Millisecond, // deltaGrace,
// 		300*time.Millisecond, // deltaCertifiedCommitRequest,
// 		1*time.Minute,        // deltaStage,
// 		100,                  // rMax,
// 		[]int{len(oracles)},  // s,
// 		oracles,              // oracles,
// 		[]byte{},             // reportingPluginConfig, // TODO(gg): put something here?
// 		nil,                  // maxDurationInitialization,
// 		0,                    // maxDurationQuery,
// 		250*time.Millisecond, // maxDurationObservation,
// 		0,                    // maxDurationShouldAcceptAttestedReport,
// 		0,                    // maxDurationShouldTransmitAcceptedReport,
// 		int(fNodes),          // f,
// 		onchainConfigBytes,   // onchainConfig (binary blob containing configuration passed through to the ReportingPlugin and also available to the contract. Unlike ReportingPluginConfig which is only available offchain.)
// 	)
// 	require.NoError(t, err)

// 	t.Logf("offchainConfig: %s", hex.EncodeToString(offchainConfig))

// 	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
// 	require.NoError(t, err)

// 	transmitterAddresses := make([]common.Address, len(transmitters))
// 	for i := range transmitters {
// 		keys, err := nodes[i].App.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
// 		require.NoError(t, err)
// 		transmitterAddresses[i] = keys[0].Address // assuming the first key is the transmitter
// 	}

// 	_, err = dfCacheContract.SetConfig(steve, signerAddresses, transmitterAddresses, f, outOnchainConfig, offchainConfigVersion, offchainConfig)
// 	if err != nil {
// 		errString, err := rPCErrorFromError(err)
// 		require.NoError(t, err)

// 		t.Fatalf("Failed to configure contract: %s", errString)
// 	}

// 	// donIDPadded := llo.DonIDToBytes32(donID)
// 	// _, err = legacyVerifier.SetConfig(steve, donIDPadded, signerAddresses, offchainTransmitters, fNodes, onchainConfig, offchainConfigVersion, offchainConfig, nil)
// 	// require.NoError(t, err)

// 	// libocr requires a few confirmations to accept the config
// 	backend.Commit()
// 	backend.Commit()
// 	backend.Commit()
// 	backend.Commit()

// 	// l, err := legacyVerifier.LatestConfigDigestAndEpoch(&bind.CallOpts{}, donIDPadded)
// 	// require.NoError(t, err)

// 	l, err := dfCacheContract.LatestConfigDigestAndEpoch(&bind.CallOpts{})
// 	require.NoError(t, err)

// 	return l.ConfigDigest
// }

// setSecureMintOnchainConfigUsingEvmSimpleConfig deploys the OCRConfigurationStoreEVMSimple contract and sets the configuration for Secure Mint using it.
// Normal data feeds use the Aggregator contract to set onchain configuration for startup, but for Secure Mint we want to write to the DF Cache, so it would be weird/confusing to deploy an Aggregator
// contract just to set the configuration. Instead, we use the OCRConfigurationStoreEVMSimple contract for this purpose.
func setSecureMintOnchainConfigUsingEvmSimpleConfig(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, nodes []Node, oracles []confighelper.OracleIdentityExtra) (common.Address, *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimple) {

	ocrConfigStoreAddress, _, ocrConfigStore, err := ocrconfigurationstoreevmsimple.DeployOCRConfigurationStoreEVMSimple(steve, backend.Client())
	if err != nil {
		rPCError, err := rPCErrorFromError(err)
		require.NoError(t, err)
		t.Fatalf("Failed to deploy OCRConfigurationStoreEVMSimple contract: %s", rPCError)
	}
	backend.Commit()

	configCh := make(chan *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleNewConfiguration)
	ocrConfigStore.WatchNewConfiguration(&bind.WatchOpts{}, configCh, nil)
	go func() {
		for config := range configCh {
			t.Logf("TRACE New configuration added to OCRConfigurationStoreEVMSimple: %s", fmt.Sprintf("0x%x", config.ConfigDigest))
		}
	}()

	configureIt(t, ocrConfigStore, steve, backend, nodes, oracles)

	return ocrConfigStoreAddress, ocrConfigStore
}

func configureIt(t *testing.T, ocrConfigStore *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimple, steve *bind.TransactOpts, backend evmtypes.Backend, nodes []Node, oracles []confighelper.OracleIdentityExtra) {

	onchainConfig := por.PorOffchainConfig{} // TODO(gg): set config values
	onchainConfigBytes, err := onchainConfig.Serialize()
	require.NoError(t, err)

	signers, transmitters, f, outOnchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress,
		20*time.Second,       // deltaResend,
		400*time.Millisecond, // deltaInitial,
		500*time.Millisecond, // deltaRound,
		250*time.Millisecond, // deltaGrace,
		300*time.Millisecond, // deltaCertifiedCommitRequest,
		1*time.Minute,        // deltaStage,
		100,                  // rMax,
		[]int{len(oracles)},  // s,
		oracles,              // oracles,
		[]byte{},             // reportingPluginConfig, // TODO(gg): put something here?
		nil,                  // maxDurationInitialization,
		0,                    // maxDurationQuery,
		250*time.Millisecond, // maxDurationObservation,
		0,                    // maxDurationShouldAcceptAttestedReport,
		0,                    // maxDurationShouldTransmitAcceptedReport,
		int(fNodes),          // f,
		onchainConfigBytes,   // onchainConfig (binary blob containing configuration passed through to the ReportingPlugin and also available to the contract. Unlike ReportingPluginConfig which is only available offchain.)
	)
	require.NoError(t, err)

	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)

	transmitterAddresses := make([]common.Address, len(transmitters))
	for i := range transmitters {
		keys, err := nodes[i].App.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		transmitterAddresses[i] = keys[0].Address // assuming the first key is the transmitter
	}

	ocrConfig := ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleConfigurationEVMSimple{
		ContractAddress:       common.Address{},
		ConfigCount:           1,
		Signers:               signerAddresses,
		Transmitters:          transmitterAddresses,
		F:                     f,
		OnchainConfig:         outOnchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
	_, err = ocrConfigStore.AddConfig(steve, ocrConfig)
	if err != nil {
		errString, err := rPCErrorFromError(err)
		require.NoError(t, err)

		t.Fatalf("Failed to configure contract: %s", errString)
	}

	// donIDPadded := llo.DonIDToBytes32(donID)
	// _, err = legacyVerifier.SetConfig(steve, donIDPadded, signerAddresses, offchainTransmitters, fNodes, onchainConfig, offchainConfigVersion, offchainConfig, nil)
	// require.NoError(t, err)

	// libocr requires a few confirmations to accept the config
	backend.Commit()
	backend.Commit()
	backend.Commit()
	backend.Commit()

	// l, err := legacyVerifier.LatestConfigDigestAndEpoch(&bind.CallOpts{}, donIDPadded)
	// require.NoError(t, err)

	// l, err := dfCacheContract.LatestConfigDigestAndEpoch(&bind.CallOpts{})
	// require.NoError(t, err)
}

func setSecureMintOnchainConfigUsingAggregator(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, nodes []Node, oracles []confighelper.OracleIdentityExtra) common.Address {

	// 1. Deploy aggregator contract

	// these min and max answers are not used by the secure mint oracle but they're needed for validation in aggregator.setConfig()
	// TODO(gg): maybe these could be 0 and max int?
	minAnswer, maxAnswer := new(big.Int), new(big.Int)
	minAnswer.Exp(big.NewInt(-2), big.NewInt(191), nil)
	maxAnswer.Exp(big.NewInt(2), big.NewInt(191), nil)
	maxAnswer.Sub(maxAnswer, big.NewInt(1))

	aggregatorAddress, _, aggregatorContract, err := ocr2aggregator.DeployOCR2Aggregator(
		steve,
		backend.Client(),
		common.Address{},   // _link common.Address,
		minAnswer,          // -2**191
		maxAnswer,          // 2**191 - 1
		common.Address{},   // accessAddress
		common.Address{},   // accessAddress
		9,                  // decimals
		"secure mint test", // description
	)
	if err != nil {
		rPCError, err := rPCErrorFromError(err)
		require.NoError(t, err)
		t.Fatalf("Failed to deploy OCR2Aggregator contract: %s", rPCError)
	}
	// Ensure we have finality depth worth of blocks to start.
	for i := 0; i < 20; i++ {
		backend.Commit()
	}
	t.Logf("Deployed OCR2Aggregator contract at: %s", aggregatorAddress.Hex())

	// 2. Create config
	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(minAnswer, maxAnswer) // TODO(gg): this uses the median codec, not sure if this is correct
	require.NoError(t, err)

	smPluginConfig := por.PorOffchainConfig{MaxChains: 5} // TODO(gg): set config values
	smPluginConfigBytes, err := smPluginConfig.Serialize()
	require.NoError(t, err)

	signers, _, f, outOnchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress,
		20*time.Second,       // deltaResend,
		400*time.Millisecond, // deltaInitial,
		500*time.Millisecond, // deltaRound,
		250*time.Millisecond, // deltaGrace,
		300*time.Millisecond, // deltaCertifiedCommitRequest,
		1*time.Minute,        // deltaStage,
		100,                  // rMax,
		[]int{len(oracles)},  // s,
		oracles,              // oracles,
		smPluginConfigBytes,  // reportingPluginConfig,
		nil,                  // maxDurationInitialization,
		0,                    // maxDurationQuery,
		250*time.Millisecond, // maxDurationObservation,
		0,                    // maxDurationShouldAcceptAttestedReport,
		0,                    // maxDurationShouldTransmitAcceptedReport,
		int(fNodes),          // f,
		onchainConfig,        // onchainConfig (binary blob containing configuration passed through to the ReportingPlugin and also available to the contract. Unlike ReportingPluginConfig which is only available offchain.)
	)
	require.NoError(t, err)

	// 3. Set config on the contract
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)

	transmitterAddresses := make([]common.Address, len(nodes))
	for i := range nodes {
		keys, err := nodes[i].App.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		transmitterAddresses[i] = keys[0].Address // assuming the first key is the transmitter
	}

	_, err = aggregatorContract.SetConfig(steve, signerAddresses, transmitterAddresses, f, outOnchainConfig, offchainConfigVersion, offchainConfig)
	if err != nil {
		errString, err := rPCErrorFromError(err)
		require.NoError(t, err)
		t.Fatalf("Failed to configure contract: %s", errString)
	}

	// libocr requires a few confirmations to accept the config
	backend.Commit()
	backend.Commit()
	backend.Commit()
	backend.Commit()

	aggregatorConfigDigest, err := aggregatorContract.LatestConfigDigestAndEpoch(&bind.CallOpts{})
	if err != nil {
		rPCError, err := rPCErrorFromError(err)
		require.NoError(t, err)
		t.Fatalf("Failed to get latest config digest: %s", rPCError)
	}
	t.Logf("Aggregator config digest: 0x%x", aggregatorConfigDigest.ConfigDigest)

	return aggregatorAddress
}

// func generateSmConfig(t *testing.T, opts ...OCRConfigOption) (signers []types.OnchainPublicKey, transmitters []types.Account, f uint8, outOnchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) {

// 	return
// }

// func setSmConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, legacyVerifier *verifier.Verifier, legacyVerifierAddr common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra, inOffchainConfig datastreamsllo.OffchainConfig) ocr2types.ConfigDigest {

// 	return l.ConfigDigest
// }

func rPCErrorFromError(txError error) (string, error) {
	errBytes, err := json.Marshal(txError)
	if err != nil {
		return "", err
	}
	var callErr struct {
		Code    int
		Data    string `json:"data"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(errBytes, &callErr)
	if err != nil {
		return "", err
	}
	// If the error data is blank
	if len(callErr.Data) == 0 {
		return callErr.Data, nil
	}
	// Some nodes prepend "Reverted " and we also remove the 0x
	trimmed := strings.TrimPrefix(callErr.Data, "Reverted ")[2:]
	data, err := hex.DecodeString(trimmed)
	if err != nil {
		return "", err
	}
	revert, err := abi.UnpackRevert(data)
	// If we can't decode the revert reason, return the raw data
	if err != nil {
		return callErr.Data, nil
	}
	return revert, nil
}

/**
blockBeforeConfig, err = b.Client().BlockByNumber(testutils.Context(t), nil)
require.NoError(t, err)
signers, effectiveTransmitters, threshold, _, encodedConfigVersion, encodedConfig, err := confighelper2.ContractSetConfigArgsForEthereumIntegrationTest(
	oracles,
	1,
	1000000000/100, // threshold PPB
)
require.NoError(t, err)

minAnswer, maxAnswer := new(big.Int), new(big.Int)
minAnswer.Exp(big.NewInt(-2), big.NewInt(191), nil)
maxAnswer.Exp(big.NewInt(2), big.NewInt(191), nil)
maxAnswer.Sub(maxAnswer, big.NewInt(1))

onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(minAnswer, maxAnswer)
require.NoError(t, err)

lggr.Debugw("Setting Config on Oracle Contract",
	"signers", signers,
	"transmitters", transmitters,
	"effectiveTransmitters", effectiveTransmitters,
	"threshold", threshold,
	"onchainConfig", onchainConfig,
	"encodedConfigVersion", encodedConfigVersion,
)
_, err = ocrContract.SetConfig(
	owner,
	signers,
	effectiveTransmitters,
	threshold,
	onchainConfig,
	encodedConfigVersion,
	encodedConfig,
)
require.NoError(t, err)
b.Commit()
*/

func setupDataFeedsCacheContract(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, allowedSenders []common.Address, workflowOwner, workflowName string) (
	common.Address, *data_feeds_cache.DataFeedsCache) {

	addr, _, dataFeedsCache, err := data_feeds_cache.DeployDataFeedsCache(steve, backend.Client())
	require.NoError(t, err)
	backend.Commit()

	var nameBytes [10]byte
	copy(nameBytes[:], workflowName)

	ownerAddr := common.HexToAddress(workflowOwner)

	_, err = dataFeedsCache.SetFeedAdmin(steve, ownerAddr, true)
	require.NoError(t, err)

	backend.Commit()

	metadatas := make([]data_feeds_cache.DataFeedsCacheWorkflowMetadata, len(allowedSenders))
	for i, sender := range allowedSenders {
		metadatas[i] =
			data_feeds_cache.DataFeedsCacheWorkflowMetadata{
				AllowedSender:        sender,
				AllowedWorkflowOwner: ownerAddr,
				AllowedWorkflowName:  nameBytes,
			}
	}

	feedIDBytes := [16]byte{}
	copy(feedIDBytes[:], common.FromHex("0xA1B2C3D4E5F600010203040506070809"))

	_, err = dataFeedsCache.SetDecimalFeedConfigs(steve, [][16]byte{feedIDBytes}, []string{"securemint"}, metadatas)
	if err != nil {
		errString, err := rPCErrorFromError(err)
		require.NoError(t, err)

		t.Fatalf("Failed to configure contract: %s", errString)
	}

	backend.Commit()

	return addr, dataFeedsCache
}
