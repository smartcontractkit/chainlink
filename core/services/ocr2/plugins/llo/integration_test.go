package llo_test

import (
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/wsrpc/credentials"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	lloevm "github.com/smartcontractkit/chainlink-data-streams/llo/reportcodecs/evm"
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
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	reportcodecv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
	mercuryverifier "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/verifier"
)

var (
	fNodes = uint8(1)
	nNodes = 4 // number of nodes (not including bootstrap)
)

func setupBlockchain(t *testing.T) (
	*bind.TransactOpts,
	evmtypes.Backend,
	*configurator.Configurator,
	common.Address,
	*destination_verifier.DestinationVerifier,
	common.Address,
	*destination_verifier_proxy.DestinationVerifierProxy,
	common.Address,
	*channel_config_store.ChannelConfigStore,
	common.Address,
	*verifier.Verifier,
	common.Address,
	*verifier_proxy.VerifierProxy,
	common.Address,
) {
	steve := evmtestutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := gethtypes.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	backend.Commit()
	backend.Commit() // ensure starting block number at least 1

	// Configurator
	configuratorAddress, _, configurator, err := configurator.DeployConfigurator(steve, backend.Client())
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
	legacyVerifier, legacyVerifierAddr, legacyVerifierProxy, legacyVerifierProxyAddr := setupLegacyMercuryVerifier(t, steve, backend)

	// ChannelConfigStore
	configStoreAddress, _, configStore, err := channel_config_store.DeployChannelConfigStore(steve, backend.Client())
	require.NoError(t, err)

	backend.Commit()

	return steve, backend, configurator, configuratorAddress, destinationVerifier, destinationVerifierAddr, verifierProxy, destinationVerifierProxyAddr, configStore, configStoreAddress, legacyVerifier, legacyVerifierAddr, legacyVerifierProxy, legacyVerifierProxyAddr
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

func WithPredecessorConfigDigest(predecessorConfigDigest ocr2types.ConfigDigest) OCRConfigOption {
	return func(cfg *OCRConfig) {
		onchainConfig, err := (&datastreamsllo.EVMOnchainConfigCodec{}).Encode(datastreamsllo.OnchainConfig{
			Version:                 1,
			PredecessorConfigDigest: &predecessorConfigDigest,
		})
		if err != nil {
			panic(err)
		}
		cfg.OnchainConfig = onchainConfig
	}
}

func WithOffchainConfig(offchainConfig datastreamsllo.OffchainConfig) OCRConfigOption {
	return func(cfg *OCRConfig) {
		offchainConfigEncoded, err := offchainConfig.Encode()
		if err != nil {
			panic(err)
		}
		cfg.ReportingPluginConfig = offchainConfigEncoded
	}
}

func WithOracles(oracles []confighelper.OracleIdentityExtra) OCRConfigOption {
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
	t.Logf("Using OCR config: %+v\n", cfg)
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
	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig := generateConfig(t, WithOracles(oracles), WithOffchainConfig(inOffchainConfig))

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

func setStagingConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, opts ...OCRConfigOption) ocr2types.ConfigDigest {
	return setBlueGreenConfig(t, donID, steve, backend, configurator, configuratorAddress, nodes, opts...)
}

func setProductionConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, opts ...OCRConfigOption) ocr2types.ConfigDigest {
	return setBlueGreenConfig(t, donID, steve, backend, configurator, configuratorAddress, nodes, opts...)
}

func setBlueGreenConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, opts ...OCRConfigOption) ocr2types.ConfigDigest {
	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig := generateConfig(t, opts...)

	var onchainPubKeys [][]byte
	for _, signer := range signers {
		onchainPubKeys = append(onchainPubKeys, signer)
	}
	offchainTransmitters := make([][32]byte, nNodes)
	for i := 0; i < nNodes; i++ {
		offchainTransmitters[i] = nodes[i].ClientPubKey
	}
	donIDPadded := llo.DonIDToBytes32(donID)
	var isProduction bool
	{
		cfg, err := (&datastreamsllo.EVMOnchainConfigCodec{}).Decode(onchainConfig)
		require.NoError(t, err)
		isProduction = cfg.PredecessorConfigDigest == nil
	}
	var err error
	if isProduction {
		_, err = configurator.SetProductionConfig(steve, donIDPadded, onchainPubKeys, offchainTransmitters, fNodes, onchainConfig, offchainConfigVersion, offchainConfig)
	} else {
		_, err = configurator.SetStagingConfig(steve, donIDPadded, onchainPubKeys, offchainTransmitters, fNodes, onchainConfig, offchainConfigVersion, offchainConfig)
	}
	require.NoError(t, err)

	// libocr requires a few confirmations to accept the config
	backend.Commit()
	backend.Commit()
	backend.Commit()
	backend.Commit()

	var topic common.Hash
	if isProduction {
		topic = llo.ProductionConfigSet
	} else {
		topic = llo.StagingConfigSet
	}
	logs, err := backend.Client().FilterLogs(testutils.Context(t), ethereum.FilterQuery{Addresses: []common.Address{configuratorAddress}, Topics: [][]common.Hash{[]common.Hash{topic, donIDPadded}}})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(logs), 1)

	cfg, err := mercury.ConfigFromLog(logs[len(logs)-1].Data)
	require.NoError(t, err)

	return cfg.ConfigDigest
}

func promoteStagingConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, configurator *configurator.Configurator, configuratorAddress common.Address, isGreenProduction bool) {
	donIDPadded := llo.DonIDToBytes32(donID)
	_, err := configurator.PromoteStagingConfig(steve, donIDPadded, isGreenProduction)
	require.NoError(t, err)

	// libocr requires a few confirmations to accept the config
	backend.Commit()
	backend.Commit()
	backend.Commit()
	backend.Commit()
}

func TestIntegration_LLO_evm_premium_legacy(t *testing.T) {
	t.Parallel()
	offchainConfigs := []datastreamsllo.OffchainConfig{
		{
			ProtocolVersion:                     0,
			DefaultMinReportIntervalNanoseconds: 0,
		},
		{
			ProtocolVersion:                     1,
			DefaultMinReportIntervalNanoseconds: 1,
		},
	}
	for _, offchainConfig := range offchainConfigs {
		t.Run(fmt.Sprintf("offchainConfig=%+v", offchainConfig), func(t *testing.T) {
			t.Parallel()

			testIntegrationLLOEVMPremiumLegacy(t, offchainConfig)
		})
	}
}

func testIntegrationLLOEVMPremiumLegacy(t *testing.T, offchainConfig datastreamsllo.OffchainConfig) {
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

	steve, backend, _, _, verifier, _, verifierProxy, _, configStore, configStoreAddress, legacyVerifier, legacyVerifierAddr, _, _ := setupBlockchain(t)
	fromBlock := 1

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo", backend, bootstrapCSAKey, nil)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	t.Run("using legacy verifier configuration contract, produces reports in v0.3 format", func(t *testing.T) {
		reqs := make(chan wsrpcRequest, 100000)
		serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
		serverPubKey := serverKey.PublicKey
		srv := NewWSRPCMercuryServer(t, serverKey, reqs)

		serverURL := startWSRPCMercuryServer(t, srv, clientPubKeys)

		donID := uint32(995544)
		streams := []Stream{ethStream, linkStream, quoteStream1, quoteStream2}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
			c.Mercury.Transmitter.Protocol = ptr(config.MercuryTransmitterProtocolWSRPC)
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
		_, err := configStore.SetChannelDefinitions(steve, donID, url, sha)
		require.NoError(t, err)
		backend.Commit()

		pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
donID = %d
channelDefinitionsContractAddress = "0x%x"
channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock)
		addOCRJobsEVMPremiumLegacy(t, streams, serverPubKey, serverURL, legacyVerifierAddr, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

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

		t.Run("receives at least one report per channel from each oracle when EAs are at 100% reliability", func(t *testing.T) {
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

				// test on-chain verification
				t.Run("on-chain verification", func(t *testing.T) {
					t.Skip("SKIP - MERC-6637")
					// Disabled because it flakes, sometimes returns "execution reverted"
					// No idea why
					// https://smartcontract-it.atlassian.net/browse/MERC-6637
					_, err = verifierProxy.Verify(steve, req.req.Payload, []byte{})
					require.NoError(t, err)
				})

				t.Logf("oracle %x reported for 0x%x", req.pk[:], feedID[:])

				seen[feedID][req.pk] = struct{}{}
				if len(seen[feedID]) == nNodes {
					t.Logf("all oracles reported for 0x%x", feedID[:])
					delete(seen, feedID)
					if len(seen) == 0 {
						break // saw all oracles; success!
					}
				}
			}
		})
	})
}

func TestIntegration_LLO_multi_formats(t *testing.T) {
	t.Parallel()
	offchainConfigs := []datastreamsllo.OffchainConfig{
		{
			ProtocolVersion:                     0,
			DefaultMinReportIntervalNanoseconds: 0,
		},
		{
			ProtocolVersion:                     1,
			DefaultMinReportIntervalNanoseconds: 1,
		},
	}
	for _, offchainConfig := range offchainConfigs {
		t.Run(fmt.Sprintf("offchainConfig=%+v", offchainConfig), func(t *testing.T) {
			t.Parallel()

			testIntegrationLLOMultiFormats(t, offchainConfig)
		})
	}
}

func testIntegrationLLOMultiFormats(t *testing.T, offchainConfig datastreamsllo.OffchainConfig) {
	testStartTimeStamp := time.Now()
	expirationWindow := uint32(3600)

	const salt = 200

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)
	for i := 0; i < nNodes; i++ {
		k := big.NewInt(int64(salt + i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	steve, backend, configurator, configuratorAddress, _, _, _, _, configStore, configStoreAddress, _, _, _, _ := setupBlockchain(t)
	fromBlock := 1

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo", backend, bootstrapCSAKey, nil)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	t.Run("generates reports using multiple formats", func(t *testing.T) {
		packetCh := make(chan *packet, 100000)
		serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
		serverPubKey := serverKey.PublicKey
		srv := NewMercuryServer(t, serverKey, packetCh)

		serverURL := startMercuryServer(t, srv, clientPubKeys)

		donID := uint32(888333)
		streams := []Stream{ethStream, linkStream}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
			c.Mercury.Transmitter.Protocol = ptr(config.MercuryTransmitterProtocolGRPC)
		})

		chainID := testutils.SimulatedChainID
		relayType := "evm"
		relayConfig := fmt.Sprintf(`
chainID = "%s"
fromBlock = %d
lloDonID = %d
lloConfigMode = "bluegreen"
`, chainID, fromBlock, donID)
		addBootstrapJob(t, bootstrapNode, configuratorAddress, "job-4", relayType, relayConfig)

		dexBasedAssetPriceStreamID := uint32(1)
		marketStatusStreamID := uint32(2)
		baseMarketDepthStreamID := uint32(3)
		quoteMarketDepthStreamID := uint32(4)
		benchmarkPriceStreamID := uint32(5)
		binanceFundingRateStreamID := uint32(6)
		binanceFundingTimeStreamID := uint32(7)
		binanceFundingIntervalHoursStreamID := uint32(8)
		deribitFundingRateStreamID := uint32(9)
		deribitFundingTimeStreamID := uint32(10)
		deribitFundingIntervalHoursStreamID := uint32(11)
		timestampedStonkPriceStreamID := uint32(12)
		nullTimestampPriceStreamID := uint32(13)
		missingTimestampPriceStreamID := uint32(14)

		mustEncodeOpts := func(opts any) []byte {
			encoded, err := json.Marshal(opts)
			require.NoError(t, err)
			return encoded
		}

		standardMultiplier := ubig.NewI(1e18)

		const simpleStreamlinedChannelID = 5
		const complexStreamlinedChannelID = 6
		const sampleTimestampsStockPriceChannelID = 7

		dexBasedAssetFeedID := utils.NewHash()
		rwaFeedID := utils.NewHash()
		benchmarkPriceFeedID := utils.NewHash()
		fundingRateFeedID := utils.NewHash()
		simpleStreamlinedFeedID := pad32bytes(simpleStreamlinedChannelID)
		complexStreamlinedFeedID := pad32bytes(complexStreamlinedChannelID)
		sampleTimestampsStockPriceFeedID := pad32bytes(sampleTimestampsStockPriceChannelID)

		// Channel definitions
		channelDefinitions := llotypes.ChannelDefinitions{
			// Sample DEX-based asset schema
			1: {
				ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
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
						StreamID:   dexBasedAssetPriceStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   baseMarketDepthStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   quoteMarketDepthStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
				},
				Opts: mustEncodeOpts(&lloevm.ReportFormatEVMABIEncodeOpts{
					BaseUSDFee:       decimal.NewFromFloat32(0.1),
					ExpirationWindow: expirationWindow,
					FeedID:           dexBasedAssetFeedID,
					ABI: []lloevm.ABIEncoder{
						newSingleABIEncoder("int192", standardMultiplier),
						newSingleABIEncoder("int192", nil),
						newSingleABIEncoder("int192", nil),
					},
				}),
			},
			// Sample RWA schema
			2: {
				ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
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
						StreamID:   marketStatusStreamID,
						Aggregator: llotypes.AggregatorMode,
					},
				},
				Opts: mustEncodeOpts(&lloevm.ReportFormatEVMABIEncodeOpts{
					BaseUSDFee:       decimal.NewFromFloat32(0.1),
					ExpirationWindow: expirationWindow,
					FeedID:           rwaFeedID,
					ABI: []lloevm.ABIEncoder{
						newSingleABIEncoder("uint32", nil),
					},
				}),
			},
			// Sample Benchmark price schema
			3: {
				ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
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
						StreamID:   benchmarkPriceStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
				},
				Opts: mustEncodeOpts(&lloevm.ReportFormatEVMABIEncodeOpts{
					BaseUSDFee:       decimal.NewFromFloat32(0.1),
					ExpirationWindow: expirationWindow,
					FeedID:           benchmarkPriceFeedID,
					ABI: []lloevm.ABIEncoder{
						newSingleABIEncoder("int192", standardMultiplier),
					},
				}),
			},
			// Sample funding rate scheam
			4: {
				ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpacked,
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
						StreamID:   binanceFundingRateStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   binanceFundingTimeStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   binanceFundingIntervalHoursStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   deribitFundingRateStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   deribitFundingTimeStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   deribitFundingIntervalHoursStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
				},
				Opts: mustEncodeOpts(&lloevm.ReportFormatEVMABIEncodeOpts{
					BaseUSDFee:       decimal.NewFromFloat32(0.1),
					ExpirationWindow: expirationWindow,
					FeedID:           fundingRateFeedID,
					ABI: []lloevm.ABIEncoder{
						newSingleABIEncoder("int192", nil),
						newSingleABIEncoder("int192", nil),
						newSingleABIEncoder("int192", nil),
						newSingleABIEncoder("int192", nil),
						newSingleABIEncoder("int192", nil),
						newSingleABIEncoder("int192", nil),
					},
				}),
			},
			// Simple sample streamlined schema
			simpleStreamlinedChannelID: {
				ReportFormat: llotypes.ReportFormatEVMStreamlined,
				Streams: []llotypes.Stream{
					{
						StreamID:   ethStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
				},
				Opts: mustEncodeOpts(&lloevm.ReportFormatEVMStreamlinedOpts{
					ABI: []lloevm.ABIEncoder{
						newSingleABIEncoder("int128", standardMultiplier),
					},
				}),
			},
			// Complex sample streamlined schema
			complexStreamlinedChannelID: {
				ReportFormat: llotypes.ReportFormatEVMStreamlined,
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
						StreamID:   dexBasedAssetPriceStreamID,
						Aggregator: llotypes.AggregatorMode,
					},
				},
				Opts: mustEncodeOpts(&lloevm.ReportFormatEVMStreamlinedOpts{
					ABI: []lloevm.ABIEncoder{
						newSingleABIEncoder("int192", standardMultiplier),
						newSingleABIEncoder("int8", ubig.NewI(1)),
						newSingleABIEncoder("uint64", ubig.NewI(100)),
					},
				}),
			},
			// Sample timestamped stock price schema/RWA
			sampleTimestampsStockPriceChannelID: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{
						StreamID:   timestampedStonkPriceStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   nullTimestampPriceStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   missingTimestampPriceStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
				},
			},
		}
		url, sha := newChannelDefinitionsServer(t, channelDefinitions)

		// Set channel definitions
		_, err := configStore.SetChannelDefinitions(steve, donID, url, sha)
		require.NoError(t, err)
		backend.Commit()

		pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
donID = %d
channelDefinitionsContractAddress = "0x%x"
channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock)

		bridgeName := "superbridge"

		responseJSON := `{
	"data": {
		"benchmarkPrice": "111.22",
		"marketStatus": 1
	},
	"result": {
		"benchmarkPrice": "2976.39",
		"baseMarketDepth": "1000.1212",
		"quoteMarketDepth": "998.5431",
		"binanceFundingRate": "1234.5678",
		"binanceFundingTime": "1630000000",
		"binanceFundingIntervalHours": "8",
		"deribitFundingRate": "5432.2345",
		"deribitFundingTime": "1630000000",
		"deribitFundingIntervalHours": "8",
		"ethPrice": "3976.39",
		"linkPrice": "23.45"
	},
	"timestamps": {
		"providerIndicatedTimeUnixMs": 1742314713000,
		"providerIndicatedTimeUnixMs_TestNull": null,
		"providerDataReceivedUnixMs": 1742314713050
	}
}`

		pricePipeline := fmt.Sprintf(`
dp          [type=bridge name="%s" requestData="{\\"data\\":{\\"data\\":\\"foo\\"}}"];
eth_parse    					[type=jsonparse path="result,ethPrice"];
eth_decimal 					[type=multiply times=1 streamID=%d];
link_parse    				[type=jsonparse path="result,linkPrice"];
link_decimal 				[type=multiply times=1 streamID=%d];
dp -> eth_parse -> eth_decimal;
dp -> link_parse -> link_decimal;
`, bridgeName, ethStreamID, linkStreamID)

		dexBasedAssetPipeline := fmt.Sprintf(`
dp          [type=bridge name="%s" requestData="{\\"data\\":{\\"data\\":\\"foo\\"}}"];

bp_parse    				[type=jsonparse path="result,benchmarkPrice"];
base_market_depth_parse   	[type=jsonparse path="result,baseMarketDepth"];
quote_market_depth_parse    [type=jsonparse path="result,quoteMarketDepth"];

bp_decimal 					[type=multiply times=1 streamID=%d];
base_market_depth_decimal   [type=multiply times=1 streamID=%d];
quote_market_depth_decimal  [type=multiply times=1 streamID=%d];

dp -> bp_parse -> bp_decimal;
dp -> base_market_depth_parse -> base_market_depth_decimal;
dp -> quote_market_depth_parse -> quote_market_depth_decimal;
`, bridgeName, dexBasedAssetPriceStreamID, baseMarketDepthStreamID, quoteMarketDepthStreamID)

		// Don't use a multiply task so that the task result has int64 type.
		rwaPipeline := fmt.Sprintf(`
dp [type=bridge name="%s" requestData="{\\"data\\":{\\"data\\":\\"foo\\"}}"];

market_status_parse [type=jsonparse path="data,marketStatus" streamID=%d];
stonk_price_parse [type=jsonparse path="data,benchmarkPrice"];

dp -> market_status_parse;

provider_indicated_time_parse [type=jsonparse lax=true path="timestamps,providerIndicatedTimeUnixMs"];
provider_data_received_parse [type=jsonparse lax=true path="timestamps,providerDataReceivedUnixMs"];
provider_indicated_time [type=median lax=true];
provider_data_received [type=median lax=true];
stonk_price_timestamped [type=merge left="{}" right="{\\"streamValueType\\": %d, \\"timestamps\\":{\\"providerIndicatedTimeUnixMs\\":$(provider_indicated_time),\\"providerDataReceivedUnixMs\\":$(provider_data_received)}, \\"result\\": $(stonk_price_parse)}" streamID=%d];

dp -> provider_indicated_time_parse -> provider_indicated_time;
dp -> provider_data_received_parse -> provider_data_received;
dp -> stonk_price_parse -> stonk_price_timestamped;

# test null providerIndicatedTimeUnixMs
null_provider_indicated_time_parse [type=jsonparse lax=true path="timestamps,providerIndicatedTimeUnixMs_TestNull"];
null_provider_indicated_time [type=median lax=true];
stonk_price_timestamped_null_indicated_time [type=merge left="{}" right="{\\"streamValueType\\": %d, \\"timestamps\\":{\\"providerIndicatedTimeUnixMs\\":$(null_provider_indicated_time),\\"providerDataReceivedUnixMs\\":$(provider_data_received)}, \\"result\\": $(stonk_price_parse)}" streamID=%d];

dp -> null_provider_indicated_time_parse -> null_provider_indicated_time;
dp -> stonk_price_parse -> stonk_price_timestamped_null_indicated_time;

# test missing providerIndicatedTimeUnixMs
missing_provider_indicated_time_parse [type=jsonparse lax=true path="timestamps,providerIndicatedTimeUnixMs_TestMissing"];
missing_provider_indicated_time [type=median lax=true];
stonk_price_timestamped_missing_indicated_time [type=merge left="{}" right="{\\"streamValueType\\": %d, \\"timestamps\\":{\\"providerIndicatedTimeUnixMs\\":$(missing_provider_indicated_time),\\"providerDataReceivedUnixMs\\":$(provider_data_received)}, \\"result\\": $(stonk_price_parse)}" streamID=%d];

dp -> missing_provider_indicated_time_parse -> missing_provider_indicated_time;
dp -> stonk_price_parse -> stonk_price_timestamped_missing_indicated_time;
`, bridgeName, marketStatusStreamID,
			datastreamsllo.LLOStreamValue_TimestampedStreamValue, timestampedStonkPriceStreamID,
			datastreamsllo.LLOStreamValue_TimestampedStreamValue, nullTimestampPriceStreamID,
			datastreamsllo.LLOStreamValue_TimestampedStreamValue, missingTimestampPriceStreamID,
		)

		benchmarkPricePipeline := fmt.Sprintf(`
dp          [type=bridge name="%s" requestData="{\\"data\\":{\\"data\\":\\"foo\\"}}"];

bp_parse    				[type=jsonparse path="result,benchmarkPrice"];
bp_decimal 					[type=multiply times=1 streamID=%d];

dp -> bp_parse -> bp_decimal;
`, bridgeName, benchmarkPriceStreamID)

		fundingRatePipeline := fmt.Sprintf(`
dp          [type=bridge name="%s" requestData="{\\"data\\":{\\"data\\":\\"foo\\"}}"];

binance_funding_rate_parse   [type=jsonparse path="result,binanceFundingRate"];
binance_funding_rate_decimal [type=multiply times=1 streamID=%d];

binance_funding_time_parse   [type=jsonparse path="result,binanceFundingTime"];
binance_funding_time_decimal [type=multiply times=1 streamID=%d];

binance_funding_interval_hours_parse   [type=jsonparse path="result,binanceFundingIntervalHours"];
binance_funding_interval_hours_decimal [type=multiply times=1 streamID=%d];

deribit_funding_rate_parse   [type=jsonparse path="result,deribitFundingRate"];
deribit_funding_rate_decimal [type=multiply times=1 streamID=%d];

deribit_funding_time_parse   [type=jsonparse path="result,deribitFundingTime"];
deribit_funding_time_decimal [type=multiply times=1 streamID=%d];

deribit_funding_interval_hours_parse   [type=jsonparse path="result,deribitFundingIntervalHours"];
deribit_funding_interval_hours_decimal [type=multiply times=1 streamID=%d];

dp -> binance_funding_rate_parse -> binance_funding_rate_decimal;
dp -> binance_funding_time_parse -> binance_funding_time_decimal;
dp -> binance_funding_interval_hours_parse -> binance_funding_interval_hours_decimal;
dp -> deribit_funding_rate_parse -> deribit_funding_rate_decimal;
dp -> deribit_funding_time_parse -> deribit_funding_time_decimal;
dp -> deribit_funding_interval_hours_parse -> deribit_funding_interval_hours_decimal;

`, bridgeName, binanceFundingRateStreamID, binanceFundingTimeStreamID, binanceFundingIntervalHoursStreamID, deribitFundingRateStreamID, deribitFundingTimeStreamID, deribitFundingIntervalHoursStreamID)

		for i, node := range nodes {
			// superBridge returns a JSON with everything you want in it,
			// stream specs can just pick the individual fields they need
			createBridge(t, bridgeName, responseJSON, node.App.BridgeORM())
			addStreamSpec(t, node, "pricePipeline", nil, pricePipeline)
			addStreamSpec(t, node, "dexBasedAssetPipeline", nil, dexBasedAssetPipeline)
			addStreamSpec(t, node, "rwaPipeline", nil, rwaPipeline)
			addStreamSpec(t, node, "benchmarkPricePipeline", nil, benchmarkPricePipeline)
			addStreamSpec(t, node, "fundingRatePipeline", nil, fundingRatePipeline)

			addLLOJob(
				t,
				node,
				configuratorAddress,
				bootstrapPeerID,
				bootstrapNodePort,
				clientPubKeys[i],
				"llo-evm-abi-encode-unpacked-test",
				pluginConfig,
				relayType,
				relayConfig,
			)
		}

		// Set config on configurator
		digest := setProductionConfig(
			t, donID, steve, backend, configurator, configuratorAddress, nodes, WithOracles(oracles), WithOffchainConfig(offchainConfig),
		)

		// NOTE: Wait for one of each type of report
		feedIDs := map[[32]byte]struct{}{
			dexBasedAssetFeedID:              {},
			rwaFeedID:                        {},
			benchmarkPriceFeedID:             {},
			fundingRateFeedID:                {},
			simpleStreamlinedFeedID:          {},
			complexStreamlinedFeedID:         {},
			sampleTimestampsStockPriceFeedID: {},
		}

		for pckt := range packetCh {
			req := pckt.req
			switch req.ReportFormat {
			case uint32(llotypes.ReportFormatEVMABIEncodeUnpacked):
				v := make(map[string]interface{})
				err := mercury.PayloadTypes.UnpackIntoMap(v, req.Payload)
				require.NoError(t, err)
				report, exists := v["report"]
				if !exists {
					t.Fatalf("expected payload %#v to contain 'report'", v)
				}
				reportCtx, exists := v["reportContext"]
				if !exists {
					t.Fatalf("expected payload %#v to contain 'reportContext'", v)
				}

				// Check the report context
				assert.Equal(t, [32]byte(digest), reportCtx.([3][32]uint8)[0])                                                                      // config digest
				assert.Equal(t, "000000000000000000000000000000000000000000000000000d8e0d00000001", fmt.Sprintf("%x", reportCtx.([3][32]uint8)[2])) // extra hash

				reportElems := make(map[string]interface{})
				err = lloevm.BaseSchema.UnpackIntoMap(reportElems, report.([]byte))
				require.NoError(t, err)

				feedID := reportElems["feedId"].([32]uint8)
				delete(feedIDs, feedID)

				// Check headers
				assert.GreaterOrEqual(t, reportElems["validFromTimestamp"].(uint32), uint32(testStartTimeStamp.Unix())) //nolint:gosec // G115
				assert.GreaterOrEqual(t, int(reportElems["observationsTimestamp"].(uint32)), int(testStartTimeStamp.Unix()))
				// Zero fees since both eth/link stream specs are missing, don't
				// care about billing for purposes of this test
				assert.Equal(t, "25148438659186", reportElems["nativeFee"].(*big.Int).String())
				assert.Equal(t, "4264392324093817", reportElems["linkFee"].(*big.Int).String())
				assert.Equal(t, reportElems["observationsTimestamp"].(uint32)+expirationWindow, reportElems["expiresAt"].(uint32))

				// Check payload values
				payload := report.([]byte)[192:]
				switch hex.EncodeToString(feedID[:]) {
				case hex.EncodeToString(dexBasedAssetFeedID[:]):
					require.Len(t, payload, 96)
					args := abi.Arguments([]abi.Argument{
						{Name: "benchmarkPrice", Type: mustNewType("int192")},
						{Name: "baseMarketDepth", Type: mustNewType("int192")},
						{Name: "quoteMarketDepth", Type: mustNewType("int192")},
					})
					v := make(map[string]interface{})
					err := args.UnpackIntoMap(v, payload)
					require.NoError(t, err)

					assert.Equal(t, "2976390000000000000000", v["benchmarkPrice"].(*big.Int).String())
					assert.Equal(t, "1000", v["baseMarketDepth"].(*big.Int).String())
					assert.Equal(t, "998", v["quoteMarketDepth"].(*big.Int).String())
				case hex.EncodeToString(rwaFeedID[:]):
					require.Len(t, payload, 32)
					args := abi.Arguments([]abi.Argument{
						{Name: "marketStatus", Type: mustNewType("uint32")},
					})
					v := make(map[string]interface{})
					err := args.UnpackIntoMap(v, payload)
					require.NoError(t, err)

					assert.Equal(t, uint32(1), v["marketStatus"].(uint32))
				case hex.EncodeToString(benchmarkPriceFeedID[:]):
					require.Len(t, payload, 32)
					args := abi.Arguments([]abi.Argument{
						{Name: "benchmarkPrice", Type: mustNewType("int192")},
					})
					v := make(map[string]interface{})
					err := args.UnpackIntoMap(v, payload)
					require.NoError(t, err)

					assert.Equal(t, "2976390000000000000000", v["benchmarkPrice"].(*big.Int).String())
				case hex.EncodeToString(fundingRateFeedID[:]):
					require.Len(t, payload, 192)
					args := abi.Arguments([]abi.Argument{
						{Name: "binanceFundingRate", Type: mustNewType("int192")},
						{Name: "binanceFundingTime", Type: mustNewType("int192")},
						{Name: "binanceFundingIntervalHours", Type: mustNewType("int192")},
						{Name: "deribitFundingRate", Type: mustNewType("int192")},
						{Name: "deribitFundingTime", Type: mustNewType("int192")},
						{Name: "deribitFundingIntervalHours", Type: mustNewType("int192")},
					})
					v := make(map[string]interface{})
					err := args.UnpackIntoMap(v, payload)
					require.NoError(t, err)

					assert.Equal(t, "1234", v["binanceFundingRate"].(*big.Int).String())
					assert.Equal(t, "1630000000", v["binanceFundingTime"].(*big.Int).String())
					assert.Equal(t, "8", v["binanceFundingIntervalHours"].(*big.Int).String())
					assert.Equal(t, "5432", v["deribitFundingRate"].(*big.Int).String())
					assert.Equal(t, "1630000000", v["deribitFundingTime"].(*big.Int).String())
					assert.Equal(t, "8", v["deribitFundingIntervalHours"].(*big.Int).String())
				default:
					t.Fatalf("unexpected feedID: %x", feedID)
				}
			case uint32(llotypes.ReportFormatEVMStreamlined):
				p := &lloevm.LLOEVMStreamlinedReportWithContext{}
				require.NoError(t, proto.Unmarshal(req.Payload, p))
				// proto auxiliary fields
				assert.Equal(t, digest[:], p.ConfigDigest)
				assert.Greater(t, p.SeqNr, uint64(1))

				// payload check
				payload := p.PackedPayload
				assert.Equal(t, digest[:], payload[:32])
				lenReport := int(binary.BigEndian.Uint16(payload[32:34]))
				report := make([]byte, lenReport)
				copy(report, payload[34:])
				numSigs := payload[34+lenReport]
				assert.Equal(t, int(fNodes+1), int(numSigs))
				assert.Len(t, payload, 32+2+lenReport+1+int(numSigs)*65)

				// report contents check
				// uint32 report format
				// uint32 channel ID
				rfBytes := report[:4]
				rf := binary.BigEndian.Uint32(rfBytes)
				assert.Equal(t, uint32(llotypes.ReportFormatEVMStreamlined), rf)
				cidBytes := report[4:8]
				cid := binary.BigEndian.Uint32(cidBytes)
				switch cid {
				case simpleStreamlinedChannelID:
					assert.Len(t, report, 32)
					tsbytes := report[8:16]
					ts := binary.BigEndian.Uint64(tsbytes)
					assert.GreaterOrEqual(t, ts, uint64(testStartTimeStamp.Unix())) //nolint:gosec // g115
					// int128
					assert.Equal(t, "00000000000000d78f7f252ecf870000", hex.EncodeToString(report[16:]))
				case complexStreamlinedChannelID:
					assert.Len(t, report, 49)
					tsbytes := report[8:16]
					ts := binary.BigEndian.Uint64(tsbytes)
					assert.GreaterOrEqual(t, ts, uint64(testStartTimeStamp.Unix())) //nolint:gosec // g115
					// int192, int8, uint64 stream values
					assert.Equal(t, "000000000000000000000000000000d78f7f252ecf870000", hex.EncodeToString(report[16:40]))
					assert.Equal(t, "17", hex.EncodeToString(report[40:41]))
					assert.Equal(t, "0000000000048aa7", hex.EncodeToString(report[41:]))
				default:
					t.Fatalf("unexpected channel: %d", cid)
				}
				delete(feedIDs, pad32bytes(cid))
			case uint32(llotypes.ReportFormatJSON):
				v := make(map[string]interface{})
				err := json.Unmarshal(req.Payload, &v)
				require.NoError(t, err)
				report := v["report"].(map[string]interface{})
				cid := report["ChannelID"].(float64)
				delete(feedIDs, pad32bytes(uint32(cid)))
				assert.Len(t, report["Values"].([]interface{}), 3)
				// default uses provider indicated time
				tsv1 := report["Values"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, 2, int(tsv1["t"].(float64)))
				assert.Equal(t, `TSV{ObservedAtNanoseconds: 1742314713000000000, StreamValue: {"t":0,"v":"111.22"}}`, tsv1["v"].(string))
				// null provider indicated time - uses data received time fallback
				tsv2 := report["Values"].([]interface{})[1].(map[string]interface{})
				assert.Equal(t, 2, int(tsv2["t"].(float64)))
				assert.Equal(t, `TSV{ObservedAtNanoseconds: 1742314713050000000, StreamValue: {"t":0,"v":"111.22"}}`, tsv2["v"].(string))
				// missing provider indicated time - uses data received time fallback
				tsv3 := report["Values"].([]interface{})[2].(map[string]interface{})
				assert.Equal(t, 2, int(tsv3["t"].(float64)))
				assert.Equal(t, `TSV{ObservedAtNanoseconds: 1742314713050000000, StreamValue: {"t":0,"v":"111.22"}}`, tsv3["v"].(string))
			default:
				t.Fatalf("unexpected report format: %d", req.ReportFormat)
			}

			if len(feedIDs) == 0 {
				break
			}
		}
	})
}

func TestIntegration_LLO_stress_test_V1(t *testing.T) {
	t.Parallel()

	// logLevel: the log level to use for the nodes
	// setting a more verbose log level increases cpu usage significantly
	// const logLevel = toml.LogLevel(zapcore.DebugLevel)
	const logLevel = toml.LogLevel(zapcore.ErrorLevel)

	// NOTE: Tweak these values to increase or decrease the intensity of the
	// stress test
	//
	// nChannels: the total number of channels
	// nReports: the number of reports to expect per node
	// defaultMinReportInterval: minimum time between report emission (set to 1ns to produce as fast as possible)

	// STRESS TEST PARAMETERS

	// LOW STRESS
	const nChannels = 100
	const nReports = 250
	const defaultMinReportInterval = 5 * time.Millisecond

	// HIGHER STRESS
	// const nChannels = 2000
	// const nReports = 50_000
	// const defaultMinReportInterval = 1 * time.Nanosecond

	// PROTOCOL CONFIGURATION
	ocrConfigOpts := []OCRConfigOption{
		WithOffchainConfig(datastreamsllo.OffchainConfig{
			ProtocolVersion:                     1,
			DefaultMinReportIntervalNanoseconds: uint64(defaultMinReportInterval),
		}),
		func(cfg *OCRConfig) {
			// cfg.DeltaRound = 0 // Go as fast as possible
			cfg.DeltaRound = 50 * time.Millisecond
			cfg.DeltaGrace = 5 * time.Millisecond
			cfg.DeltaCertifiedCommitRequest = 5 * time.Millisecond
		},
	}

	// SETUP

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)

	const salt = 302

	for i := 0; i < nNodes; i++ {
		k := big.NewInt(int64(salt + i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	steve, backend, configurator, configuratorAddress, _, _, _, _, configStore, configStoreAddress, _, _, _, _ := setupBlockchain(t)
	fromBlock := 1

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo", backend, bootstrapCSAKey, func(c *chainlink.Config) {
		c.Log.Level = ptr(logLevel)
	})
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	t.Run("produces reports properly", func(t *testing.T) {
		packets := make(chan *packet, nReports*nNodes)
		serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
		serverPubKey := serverKey.PublicKey
		srv := NewMercuryServer(t, serverKey, packets)

		serverURL := startMercuryServer(t, srv, clientPubKeys)

		donID := uint32(888333)
		streams := []Stream{ethStream}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
			c.Mercury.Transmitter.Protocol = ptr(config.MercuryTransmitterProtocolGRPC)
			c.Log.Level = ptr(logLevel)
		})

		chainID := testutils.SimulatedChainID
		relayType := "evm"
		relayConfig := fmt.Sprintf(`
chainID = "%s"
fromBlock = %d
lloDonID = %d
lloConfigMode = "bluegreen"
`, chainID, fromBlock, donID)
		addBootstrapJob(t, bootstrapNode, configuratorAddress, "job-3", relayType, relayConfig)

		// Channel definitions
		channelDefinitions := llotypes.ChannelDefinitions{}
		for i := uint32(0); i < nChannels; i++ {
			channelDefinitions[i] = llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{
						StreamID:   ethStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
				},
			}
		}
		url, sha := newChannelDefinitionsServer(t, channelDefinitions)

		// Set channel definitions
		_, err := configStore.SetChannelDefinitions(steve, donID, url, sha)
		require.NoError(t, err)
		backend.Commit()

		// one working and one broken transmission server
		pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
donID = %d
channelDefinitionsContractAddress = "0x%x"
channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock)
		for i, node := range nodes {
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
			addMemoStreamSpecs(t, node, streams)
		}

		// Set config on configurator
		opts := []OCRConfigOption{WithOracles(oracles)}
		opts = append(opts, ocrConfigOpts...)
		blueDigest := setProductionConfig(
			t, donID, steve, backend, configurator, configuratorAddress, nodes, opts...,
		)

		// NOTE: Wait for nReports reports per node
		// transmitter addr => count of reports
		cnts := map[string]int{}
		// transmitter addr => channel ID => reports
		m := map[string]map[uint32][]datastreamsllo.Report{}
		stopOnce := sync.Once{}

		for pckt := range packets {
			pr, ok := peer.FromContext(pckt.ctx)
			require.True(t, ok)
			addr := pr.Addr
			req := pckt.req

			assert.Equal(t, uint32(llotypes.ReportFormatJSON), req.ReportFormat)
			_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.Payload)
			require.NoError(t, err)

			cm, exists := m[addr.String()]
			if !exists {
				cm = make(map[uint32][]datastreamsllo.Report)
				m[addr.String()] = cm
			}
			cm[r.ChannelID] = append(cm[r.ChannelID], r)

			cnts[addr.String()]++
			finished := 0
			for _, cnt := range cnts {
				if cnt >= nReports {
					finished++
				}
			}
			if finished >= nNodes {
				stopOnce.Do(func() {
					// Stop all nodes, close the channel
					// This helps transmissions have a chance to complete (but
					// doesn't ensure it; libocr cancels the transmit context
					// immediately on stop signal)
					// Loop will exit once all packets are consumed
					for _, node := range nodes {
						require.NoError(t, node.App.Stop())
					}
					close(packets)
				})
			}
		}

		// Transmissions can occur out of order when we go very fast, so sort by seqNr
		for _, cm := range m {
			for _, rs := range cm {
				sort.Slice(rs, func(i, j int) bool {
					return rs[i].SeqNr < rs[j].SeqNr
				})
			}
		}

		// Check reports
		for addr, cm := range m {
			spacings := []uint64{}
			for _, rs := range cm {
				var prevObsTsNanos uint64
				for i, r := range rs {
					assert.Equal(t, blueDigest, r.ConfigDigest)
					assert.False(t, r.Specimen)
					assert.Len(t, r.Values, 1)
					assert.Equal(t, "2976.39", r.Values[0].(*datastreamsllo.Decimal).String())

					if i > 0 {
						if rs[i-1].SeqNr+1 != r.SeqNr {
							// t.Logf("gap in SeqNr at index %d; %d!=%d: len(rs)=%d", i, rs[i-1].SeqNr, r.SeqNr, len(rs))
							// We actually expect a transmission every round; if there's a gap in seqNr it means that the transmissions were likely cut off due to the app being shut down. We are probably at the end of the usable reports list so just assume completion here.
							break
						}

						// No gaps
						require.Equal(t, prevObsTsNanos, r.ValidAfterNanoseconds, "gap in reports for transmitter %s at index %d; %d!=%d: prevReport=%s, thisReport=%s", addr, i, prevObsTsNanos, r.ValidAfterNanoseconds, mustMarshalJSON(rs[i-1]), mustMarshalJSON(r))
						// Timestamps are sane
						require.GreaterOrEqual(t, r.ObservationTimestampNanoseconds, r.ValidAfterNanoseconds, "observation timestamp is before valid after timestamp for transmitter %s at index %d: report=%s", addr, i, mustMarshalJSON(r))
						// Reports are separated by at least the minimum interval
						require.GreaterOrEqual(t, r.ObservationTimestampNanoseconds-uint64(defaultMinReportInterval), prevObsTsNanos, "reports are too close together for transmitter %s at index %d: prevReport=%s, thisReport=%s; expected at least %d nanoseconds of distance", addr, i, mustMarshalJSON(rs[i-1]), mustMarshalJSON(r), defaultMinReportInterval)

						spacings = append(spacings, r.ObservationTimestampNanoseconds-prevObsTsNanos)
					}
					prevObsTsNanos = r.ObservationTimestampNanoseconds
				}
			}
			avgSpacing := uint64(0)
			for _, spacing := range spacings {
				avgSpacing += spacing
			}
			avgSpacing /= uint64(len(spacings))
			t.Logf("transmitter %s: average spacing between reports: %d nanoseconds (%f seconds)", addr, avgSpacing, float64(avgSpacing)/1e9)
		}
	})
}

func TestIntegration_LLO_transmit_errors(t *testing.T) {
	t.Parallel()

	// logLevel: the log level to use for the nodes
	// setting a more verbose log level increases cpu usage significantly
	const logLevel = toml.LogLevel(zapcore.ErrorLevel)
	// const logLevel = toml.LogLevel(zapcore.ErrorLevel)

	// NOTE: Tweak these values to increase or decrease the intensity of the
	// stress test
	//
	// nChannels: the total number of channels
	// maxQueueSize: the maximum size of the transmit queue
	// nReports: the number of reports to expect per node

	// LESS STRESSFUL
	const nChannels = 200
	const maxQueueSize = 10
	const nReports = 1_000

	// MORE STRESSFUL
	// const nChannels = 2000
	// const maxQueueSize = 4_000
	// const nReports = 10_000

	// PROTOCOL CONFIGURATION
	// TODO: test both
	offchainConfig := datastreamsllo.OffchainConfig{
		ProtocolVersion:                     1,
		DefaultMinReportIntervalNanoseconds: uint64(50 * time.Millisecond),
	}

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)

	const salt = 301

	for i := 0; i < nNodes; i++ {
		k := big.NewInt(int64(salt + i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	steve, backend, configurator, configuratorAddress, _, _, _, _, configStore, configStoreAddress, _, _, _, _ := setupBlockchain(t)
	fromBlock := 1

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo", backend, bootstrapCSAKey, func(c *chainlink.Config) {
		c.Log.Level = ptr(logLevel)
	})
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	t.Run("transmit queue does not grow unbounded", func(t *testing.T) {
		packets := make(chan *packet, 100000)
		serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
		serverPubKey := serverKey.PublicKey
		srv := NewMercuryServer(t, serverKey, packets)

		serverURL := startMercuryServer(t, srv, clientPubKeys)

		donID := uint32(888333)
		streams := []Stream{ethStream, linkStream}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
			c.Mercury.Transmitter.Protocol = ptr(config.MercuryTransmitterProtocolGRPC)
			c.Mercury.Transmitter.TransmitQueueMaxSize = ptr(uint32(maxQueueSize)) // Test queue overflow
			c.Log.Level = ptr(logLevel)
		})

		chainID := testutils.SimulatedChainID
		relayType := "evm"
		relayConfig := fmt.Sprintf(`
chainID = "%s"
fromBlock = %d
lloDonID = %d
lloConfigMode = "bluegreen"
`, chainID, fromBlock, donID)
		addBootstrapJob(t, bootstrapNode, configuratorAddress, "job-3", relayType, relayConfig)

		// Channel definitions
		channelDefinitions := llotypes.ChannelDefinitions{}
		for i := uint32(0); i < nChannels; i++ {
			channelDefinitions[i] = llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{
						StreamID:   ethStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
				},
			}
		}
		url, sha := newChannelDefinitionsServer(t, channelDefinitions)

		// Set channel definitions
		_, err := configStore.SetChannelDefinitions(steve, donID, url, sha)
		require.NoError(t, err)
		backend.Commit()

		// one working and one broken transmission server
		pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x", "example.invalid" = "%x" }
donID = %d
channelDefinitionsContractAddress = "0x%x"
channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, serverPubKey, donID, configStoreAddress, fromBlock)
		addOCRJobsEVMPremiumLegacy(t, streams, serverPubKey, serverURL, configuratorAddress, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

		var blueDigest ocr2types.ConfigDigest

		{
			// Set config on configurator
			blueDigest = setProductionConfig(
				t, donID, steve, backend, configurator, configuratorAddress, nodes, WithOracles(oracles), WithOffchainConfig(offchainConfig),
			)

			// NOTE: Wait for nReports reports
			// count of packets received keyed by transmitter IP
			m := map[string]int{}
			for pckt := range packets {
				pr, ok := peer.FromContext(pckt.ctx)
				require.True(t, ok)
				addr := pr.Addr
				req := pckt.req

				assert.Equal(t, uint32(llotypes.ReportFormatJSON), req.ReportFormat)
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.Payload)
				require.NoError(t, err)

				assert.Equal(t, blueDigest, r.ConfigDigest)
				assert.False(t, r.Specimen)
				assert.Len(t, r.Values, 1)
				assert.Equal(t, "2976.39", r.Values[0].(*datastreamsllo.Decimal).String())

				m[addr.String()]++
				finished := 0
				for _, cnt := range m {
					if cnt >= nReports {
						finished++
					}
				}
				if finished == 4 {
					break
				}
			}
		}

		// Shut all nodes down
		for i, node := range nodes {
			require.NoError(t, node.App.Stop())
			// Ensure that the transmit queue was limited
			db := node.App.GetDB()
			cnt := 0

			// The failing server
			err := db.GetContext(t.Context(), &cnt, "SELECT count(*) FROM llo_mercury_transmit_queue WHERE server_url = 'example.invalid'")
			require.NoError(t, err)
			assert.LessOrEqual(t, cnt, maxQueueSize, "persisted transmit queue size too large for node %d for failing server", i)

			// The succeeding server
			err = db.GetContext(t.Context(), &cnt, "SELECT count(*) FROM llo_mercury_transmit_queue WHERE server_url = $1", serverURL)
			require.NoError(t, err)
			assert.LessOrEqual(t, cnt, maxQueueSize, "persisted transmit queue size too large for node %d for succeeding server", i)
		}
	})
}

func TestIntegration_LLO_blue_green_lifecycle(t *testing.T) {
	t.Parallel()

	offchainConfigs := []datastreamsllo.OffchainConfig{
		{
			ProtocolVersion:                     0,
			DefaultMinReportIntervalNanoseconds: 0,
		},
		{
			ProtocolVersion:                     1,
			DefaultMinReportIntervalNanoseconds: 1,
		},
	}
	for _, offchainConfig := range offchainConfigs {
		t.Run(fmt.Sprintf("offchainConfig=%+v", offchainConfig), func(t *testing.T) {
			t.Parallel()

			testIntegrationLLOBlueGreenLifecycle(t, offchainConfig)
		})
	}
}

func testIntegrationLLOBlueGreenLifecycle(t *testing.T, offchainConfig datastreamsllo.OffchainConfig) {
	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)

	const salt = 300

	for i := 0; i < nNodes; i++ {
		k := big.NewInt(int64(salt + i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	steve, backend, configurator, configuratorAddress, _, _, _, _, configStore, configStoreAddress, _, _, _, _ := setupBlockchain(t)
	fromBlock := 1

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo", backend, bootstrapCSAKey, nil)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	t.Run("Blue/Green lifecycle (using JSON report format)", func(t *testing.T) {
		packetCh := make(chan *packet, 100000)
		serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
		serverPubKey := serverKey.PublicKey
		srv := NewMercuryServer(t, serverKey, packetCh)

		serverURL := startMercuryServer(t, srv, clientPubKeys)

		donID := uint32(888333)
		streams := []Stream{ethStream, linkStream}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
			c.Mercury.Transmitter.Protocol = ptr(config.MercuryTransmitterProtocolGRPC)
		})

		chainID := testutils.SimulatedChainID
		relayType := "evm"
		relayConfig := fmt.Sprintf(`
chainID = "%s"
fromBlock = %d
lloDonID = %d
lloConfigMode = "bluegreen"
`, chainID, fromBlock, donID)
		addBootstrapJob(t, bootstrapNode, configuratorAddress, "job-3", relayType, relayConfig)

		// Channel definitions
		channelDefinitions := llotypes.ChannelDefinitions{
			1: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{
						StreamID:   ethStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
				},
			},
		}
		url, sha := newChannelDefinitionsServer(t, channelDefinitions)

		// Set channel definitions
		_, err := configStore.SetChannelDefinitions(steve, donID, url, sha)
		require.NoError(t, err)
		backend.Commit()

		pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
donID = %d
channelDefinitionsContractAddress = "0x%x"
channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock)
		addOCRJobsEVMPremiumLegacy(t, streams, serverPubKey, serverURL, configuratorAddress, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

		var blueDigest ocr2types.ConfigDigest
		var greenDigest ocr2types.ConfigDigest

		allReports := make(map[types.ConfigDigest][]datastreamsllo.Report)
		// start off with blue=production, green=staging (specimen reports)
		{
			// Set config on configurator
			blueDigest = setProductionConfig(
				t, donID, steve, backend, configurator, configuratorAddress, nodes, WithOracles(oracles), WithOffchainConfig(offchainConfig),
			)

			// NOTE: Wait until blue produces a report

			for pckt := range packetCh {
				req := pckt.req
				assert.Equal(t, uint32(llotypes.ReportFormatJSON), req.ReportFormat)
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.Payload)
				require.NoError(t, err)

				allReports[r.ConfigDigest] = append(allReports[r.ConfigDigest], r)

				assert.Equal(t, blueDigest, r.ConfigDigest)
				assert.False(t, r.Specimen)
				assert.Len(t, r.Values, 1)
				assert.Equal(t, "2976.39", r.Values[0].(*datastreamsllo.Decimal).String())
				break
			}
		}
		// setStagingConfig does not affect production
		{
			greenDigest = setStagingConfig(
				t, donID, steve, backend, configurator, configuratorAddress, nodes, WithPredecessorConfigDigest(blueDigest), WithOracles(oracles), WithOffchainConfig(offchainConfig),
			)

			// NOTE: Wait until green produces the first "specimen" report

			for pckt := range packetCh {
				req := pckt.req
				assert.Equal(t, uint32(llotypes.ReportFormatJSON), req.ReportFormat)
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.Payload)
				require.NoError(t, err)

				allReports[r.ConfigDigest] = append(allReports[r.ConfigDigest], r)
				if r.Specimen {
					assert.Len(t, r.Values, 1)
					assert.Equal(t, "2976.39", r.Values[0].(*datastreamsllo.Decimal).String())

					assert.Equal(t, greenDigest, r.ConfigDigest)
					break
				}
				assert.Equal(t, blueDigest, r.ConfigDigest)
			}
		}
		// promoteStagingConfig flow has clean and gapless hand off from old production to newly promoted staging instance, leaving old production instance in 'retired' state
		{
			promoteStagingConfig(t, donID, steve, backend, configurator, configuratorAddress, false)

			// NOTE: Wait for first non-specimen report for the newly promoted (green) instance

			for pckt := range packetCh {
				req := pckt.req
				assert.Equal(t, uint32(llotypes.ReportFormatJSON), req.ReportFormat)
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.Payload)
				require.NoError(t, err)

				allReports[r.ConfigDigest] = append(allReports[r.ConfigDigest], r)

				if !r.Specimen && r.ConfigDigest == greenDigest {
					break
				}
			}

			initialPromotedGreenReport := allReports[greenDigest][len(allReports[greenDigest])-1]
			finalBlueReport := allReports[blueDigest][len(allReports[blueDigest])-1]

			for _, digest := range []ocr2types.ConfigDigest{blueDigest, greenDigest} {
				// Transmissions are not guaranteed to be in order
				sort.Slice(allReports[digest], func(i, j int) bool {
					return allReports[digest][i].SeqNr < allReports[digest][j].SeqNr
				})
				seenSeqNr := uint64(0)
				highestObsTsNanos := uint64(0)
				highestValidAfterNanos := uint64(0)
				for i := 0; i < len(allReports[digest]); i++ {
					r := allReports[digest][i]
					switch digest {
					case greenDigest:
						if i == len(allReports[digest])-1 {
							assert.False(t, r.Specimen)
						} else {
							assert.True(t, r.Specimen)
						}
					case blueDigest:
						assert.False(t, r.Specimen)
					}
					if r.SeqNr > seenSeqNr {
						// skip first one
						if highestObsTsNanos > 0 {
							if digest == greenDigest && i == len(allReports[digest])-1 {
								// NOTE: This actually CHANGES on the staging
								// handover and can go backwards - the gapless
								// handover test is handled below
								break
							}
							if offchainConfig.ProtocolVersion == 0 {
								// validAfter is always truncated to 1s in v0
								// IMPORTANT: gapless handovers in v0 ONLY supported at 1s resolution!!
								assert.Equal(t, highestObsTsNanos/1e9*1e9, r.ValidAfterNanoseconds/1e9*1e9, "%d: (n-1)ObservationsTimestampSeconds->(n)ValidAfterNanoseconds should be gapless to within 1s resolution, got: %d vs %d", i, highestObsTsNanos, r.ValidAfterNanoseconds)
							} else {
								assert.Equal(t, highestObsTsNanos, r.ValidAfterNanoseconds, "%d: (n-1)ObservationsTimestampSeconds->(n)ValidAfterNanoseconds should be gapless, got: %d vs %d", i, highestObsTsNanos, r.ValidAfterNanoseconds)
							}
							assert.Greater(t, r.ObservationTimestampNanoseconds, highestObsTsNanos, "%d: overlapping/duplicate report ObservationTimestampNanoseconds, got: %d vs %d", i, r.ObservationTimestampNanoseconds, highestObsTsNanos)
							assert.Greater(t, r.ValidAfterNanoseconds, highestValidAfterNanos, "%d: overlapping/duplicate report ValidAfterNanoseconds, got: %d vs %d", i, r.ValidAfterNanoseconds, highestValidAfterNanos)
							assert.Less(t, r.ValidAfterNanoseconds, r.ObservationTimestampNanoseconds)
						}
						seenSeqNr = r.SeqNr
						highestObsTsNanos = r.ObservationTimestampNanoseconds
						highestValidAfterNanos = r.ValidAfterNanoseconds
					}
				}
			}

			// Gapless handover
			assert.Less(t, finalBlueReport.ValidAfterNanoseconds, finalBlueReport.ObservationTimestampNanoseconds)

			if offchainConfig.ProtocolVersion == 0 {
				// validAfter is always truncated to 1s in v0
				// IMPORTANT: gapless handovers in v0 ONLY supported at 1s resolution!!
				assert.Equal(t, finalBlueReport.ObservationTimestampNanoseconds/1e9*1e9, initialPromotedGreenReport.ValidAfterNanoseconds/1e9*1e9)
			} else {
				assert.Equal(t, finalBlueReport.ObservationTimestampNanoseconds, initialPromotedGreenReport.ValidAfterNanoseconds)
			}

			assert.Less(t, initialPromotedGreenReport.ValidAfterNanoseconds, initialPromotedGreenReport.ObservationTimestampNanoseconds)
		}
		// retired instance does not produce reports
		{
			// NOTE: Wait for five "green" reports to be produced and assert no "blue" reports

			i := 0
			for pckt := range packetCh {
				req := pckt.req
				i++
				if i == 5 {
					break
				}
				assert.Equal(t, uint32(llotypes.ReportFormatJSON), req.ReportFormat)
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.Payload)
				require.NoError(t, err)

				allReports[r.ConfigDigest] = append(allReports[r.ConfigDigest], r)
				assert.False(t, r.Specimen)
				assert.Equal(t, greenDigest, r.ConfigDigest)
			}
		}
		// setStagingConfig replaces 'retired' instance with new config and starts producing specimen reports again
		{
			blueDigest = setStagingConfig(
				t, donID, steve, backend, configurator, configuratorAddress, nodes, WithPredecessorConfigDigest(greenDigest), WithOracles(oracles), WithOffchainConfig(offchainConfig),
			)

			// NOTE: Wait until blue produces the first "specimen" report

			for pckt := range packetCh {
				req := pckt.req
				assert.Equal(t, uint32(llotypes.ReportFormatJSON), req.ReportFormat)
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.Payload)
				require.NoError(t, err)

				allReports[r.ConfigDigest] = append(allReports[r.ConfigDigest], r)
				if r.Specimen {
					assert.Equal(t, blueDigest, r.ConfigDigest)
					break
				}
				assert.Equal(t, greenDigest, r.ConfigDigest)
			}
		}
		// promoteStagingConfig swaps the instances again
		{
			// TODO: Check that once an instance enters 'retired' state, it
			// doesn't produce reports or bother making observations
			promoteStagingConfig(t, donID, steve, backend, configurator, configuratorAddress, true)

			// NOTE: Wait for first non-specimen report for the newly promoted (blue) instance

			for pckt := range packetCh {
				req := pckt.req
				assert.Equal(t, uint32(llotypes.ReportFormatJSON), req.ReportFormat)
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.Payload)
				require.NoError(t, err)

				allReports[r.ConfigDigest] = append(allReports[r.ConfigDigest], r)

				if !r.Specimen && r.ConfigDigest == blueDigest {
					break
				}
			}

			initialPromotedBlueReport := allReports[blueDigest][len(allReports[blueDigest])-1]
			finalGreenReport := allReports[greenDigest][len(allReports[greenDigest])-1]

			// Gapless handover
			assert.Less(t, finalGreenReport.ValidAfterNanoseconds, finalGreenReport.ObservationTimestampNanoseconds)
			if offchainConfig.ProtocolVersion == 0 {
				// validAfter is always truncated to 1s in v0
				// IMPORTANT: gapless handovers in v0 ONLY supported at 1s resolution!!
				assert.Equal(t, finalGreenReport.ObservationTimestampNanoseconds/1e9*1e9, initialPromotedBlueReport.ValidAfterNanoseconds/1e9*1e9, 1_000_000_000, "ObservationTimestampSeconds->ValidAfterNanoseconds should be gapless to within 1s resolution, got: %d vs %d", finalGreenReport.ObservationTimestampNanoseconds, initialPromotedBlueReport.ValidAfterNanoseconds)
			} else {
				assert.Equal(t, finalGreenReport.ObservationTimestampNanoseconds, initialPromotedBlueReport.ValidAfterNanoseconds, "ObservationTimestampSeconds->ValidAfterNanoseconds should be gapless, got: %d vs %d", finalGreenReport.ObservationTimestampNanoseconds, initialPromotedBlueReport.ValidAfterNanoseconds)
			}
			assert.Less(t, initialPromotedBlueReport.ValidAfterNanoseconds, initialPromotedBlueReport.ObservationTimestampNanoseconds)
		}
		// adding a new channel definition is picked up on the fly
		{
			channelDefinitions[2] = llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{
						StreamID:   linkStreamID,
						Aggregator: llotypes.AggregatorMedian,
					},
				},
			}

			url, sha := newChannelDefinitionsServer(t, channelDefinitions)

			// Set channel definitions
			_, err := configStore.SetChannelDefinitions(steve, donID, url, sha)
			require.NoError(t, err)
			backend.Commit()

			// NOTE: Wait until the first report for the new channel definition is produced

			for pckt := range packetCh {
				req := pckt.req
				assert.Equal(t, uint32(llotypes.ReportFormatJSON), req.ReportFormat)
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.Payload)
				require.NoError(t, err)

				allReports[r.ConfigDigest] = append(allReports[r.ConfigDigest], r)

				// Green is retired, it shouldn't be producing anything
				assert.Equal(t, blueDigest, r.ConfigDigest)
				assert.False(t, r.Specimen)

				if r.ChannelID == 2 {
					assert.Len(t, r.Values, 1)
					assert.Equal(t, "13.25", r.Values[0].(*datastreamsllo.Decimal).String())
					break
				}
				assert.Len(t, r.Values, 1)
				assert.Equal(t, "2976.39", r.Values[0].(*datastreamsllo.Decimal).String())
			}
		}
		t.Run("deleting the jobs turns off oracles and cleans up resources", func(t *testing.T) {
			t.Skip("TODO - MERC-3524")
		})
		t.Run("adding new jobs again picks up the correct configs", func(t *testing.T) {
			t.Skip("TODO - MERC-3524")
		})
	})
}

func setupNodes(t *testing.T, nNodes int, backend evmtypes.Backend, clientCSAKeys []csakey.KeyV2, f func(*chainlink.Config)) (oracles []confighelper.OracleIdentityExtra, nodes []Node) {
	ports := freeport.GetN(t, nNodes)
	for i := 0; i < nNodes; i++ {
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_streams_%d", i), backend, clientCSAKeys[i], f)

		nodes = append(nodes, Node{
			app, transmitter, kb, observedLogs,
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

func mustNewType(t string) abi.Type {
	result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
	if err != nil {
		panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
	}
	return result
}

func mustMarshalJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func pad32bytes(d uint32) [32]byte {
	var result [32]byte
	binary.BigEndian.PutUint32(result[28:], d)
	return result
}

func newSingleABIEncoder(typ string, multiplier *ubig.Big) (enc lloevm.ABIEncoder) {
	if multiplier == nil {
		err := json.Unmarshal([]byte(fmt.Sprintf(`{"type":"%s"}`, typ)), &enc)
		if err != nil {
			panic(err)
		}
		return
	}
	err := json.Unmarshal([]byte(fmt.Sprintf(`{"type":"%s","multiplier":"%s"}`, typ, multiplier.String())), &enc)
	if err != nil {
		panic(err)
	}
	return
}
