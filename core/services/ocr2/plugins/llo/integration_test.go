package llo_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc/peer"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/wsrpc/credentials"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/destination_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/destination_verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/fee_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/reward_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	lloevm "github.com/smartcontractkit/chainlink/v2/core/services/llo/evm"
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
	steve := testutils.MustNewSimTransactor(t) // config contract deployer and owner
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
		id:                 52,
		baseBenchmarkPrice: decimal.NewFromFloat32(2_976.39),
	}
	linkStream = Stream{
		id:                 53,
		baseBenchmarkPrice: decimal.NewFromFloat32(13.25),
	}
	quoteStream1 = Stream{
		id:                 55,
		baseBenchmarkPrice: decimal.NewFromFloat32(1000.1212),
		baseBid:            decimal.NewFromFloat32(998.5431),
		baseAsk:            decimal.NewFromFloat32(1001.6999),
	}
	quoteStream2 = Stream{
		id:                 56,
		baseBenchmarkPrice: decimal.NewFromFloat32(500.1212),
		baseBid:            decimal.NewFromFloat32(499.5431),
		baseAsk:            decimal.NewFromFloat32(502.6999),
	}
)

func generateBlueGreenConfig(t *testing.T, oracles []confighelper.OracleIdentityExtra, predecessorConfigDigest *ocr2types.ConfigDigest) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) {
	onchainConfig, err := (&datastreamsllo.EVMOnchainConfigCodec{}).Encode(datastreamsllo.OnchainConfig{
		Version:                 1,
		PredecessorConfigDigest: predecessorConfigDigest,
	})
	require.NoError(t, err)
	return generateConfig(t, oracles, onchainConfig)
}

func generateConfig(t *testing.T, oracles []confighelper.OracleIdentityExtra, inOnchainConfig []byte) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f uint8,
	outOnchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) {
	rawReportingPluginConfig := datastreamsllo.OffchainConfig{}
	reportingPluginConfig, err := rawReportingPluginConfig.Encode()
	require.NoError(t, err)

	signers, transmitters, f, outOnchainConfig, offchainConfigVersion, offchainConfig, err = ocr3confighelper.ContractSetConfigArgsForTests(
		2*time.Second,        // DeltaProgress
		20*time.Second,       // DeltaResend
		400*time.Millisecond, // DeltaInitial
		500*time.Millisecond, // DeltaRound
		250*time.Millisecond, // DeltaGrace
		300*time.Millisecond, // DeltaCertifiedCommitRequest
		1*time.Minute,        // DeltaStage
		100,                  // rMax
		[]int{len(oracles)},  // S
		oracles,
		reportingPluginConfig, // reportingPluginConfig []byte,
		nil,                   // maxDurationInitialization
		0,                     // maxDurationQuery
		250*time.Millisecond,  // maxDurationObservation
		0,                     // maxDurationShouldAcceptAttestedReport
		0,                     // maxDurationShouldTransmitAcceptedReport
		int(fNodes),           // f
		inOnchainConfig,       // encoded onchain config
	)

	require.NoError(t, err)

	return
}

func setLegacyConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, legacyVerifier *verifier.Verifier, legacyVerifierAddr common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra) ocr2types.ConfigDigest {
	onchainConfig, err := (&datastreamsllo.EVMOnchainConfigCodec{}).Encode(datastreamsllo.OnchainConfig{
		Version:                 1,
		PredecessorConfigDigest: nil,
	})
	require.NoError(t, err)

	signers, _, _, _, offchainConfigVersion, offchainConfig := generateConfig(t, oracles, onchainConfig)

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

func setStagingConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra, predecessorConfigDigest ocr2types.ConfigDigest) ocr2types.ConfigDigest {
	return setBlueGreenConfig(t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles, &predecessorConfigDigest)
}

func setProductionConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra) ocr2types.ConfigDigest {
	return setBlueGreenConfig(t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles, nil)
}

func setBlueGreenConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra, predecessorConfigDigest *ocr2types.ConfigDigest) ocr2types.ConfigDigest {
	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig := generateBlueGreenConfig(t, oracles, predecessorConfigDigest)

	var onchainPubKeys [][]byte
	for _, signer := range signers {
		onchainPubKeys = append(onchainPubKeys, signer)
	}
	offchainTransmitters := make([][32]byte, nNodes)
	for i := 0; i < nNodes; i++ {
		offchainTransmitters[i] = nodes[i].ClientPubKey
	}
	donIDPadded := llo.DonIDToBytes32(donID)
	isProduction := predecessorConfigDigest == nil
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
		srv := NewWSRPCMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), reqs)

		serverURL := startWSRPCMercuryServer(t, srv, clientPubKeys)

		donID := uint32(995544)
		streams := []Stream{ethStream, linkStream, quoteStream1, quoteStream2}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, streams, func(c *chainlink.Config) {
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
			t, donID, steve, backend, legacyVerifier, legacyVerifierAddr, nodes, oracles,
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

func TestIntegration_LLO_evm_abi_encode_unpacked(t *testing.T) {
	t.Parallel()

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

	t.Run("generates reports using go ReportFormatEVMABIEncodeUnpacked format", func(t *testing.T) {
		packetCh := make(chan *packet, 100000)
		serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
		serverPubKey := serverKey.PublicKey
		srv := NewMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), packetCh)

		serverURL := startMercuryServer(t, srv, clientPubKeys)

		donID := uint32(888333)
		streams := []Stream{ethStream, linkStream}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, streams, func(c *chainlink.Config) {
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

		mustEncodeOpts := func(opts *lloevm.ReportFormatEVMABIEncodeOpts) []byte {
			encoded, err := json.Marshal(opts)
			require.NoError(t, err)
			return encoded
		}

		standardMultiplier := ubig.NewI(1e18)

		dexBasedAssetFeedID := utils.NewHash()
		rwaFeedID := utils.NewHash()
		benchmarkPriceFeedID := utils.NewHash()
		fundingRateFeedID := utils.NewHash()
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
						lloevm.ABIEncoder{
							StreamID:   dexBasedAssetPriceStreamID,
							Type:       "int192",
							Multiplier: standardMultiplier,
						},
						lloevm.ABIEncoder{
							StreamID: baseMarketDepthStreamID,
							Type:     "int192",
						},
						lloevm.ABIEncoder{
							StreamID: quoteMarketDepthStreamID,
							Type:     "int192",
						},
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
						{
							StreamID: marketStatusStreamID,
							Type:     "uint32",
						},
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
						{
							StreamID:   benchmarkPriceStreamID,
							Type:       "int192",
							Multiplier: standardMultiplier,
						},
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
						{
							StreamID: binanceFundingRateStreamID,
							Type:     "int192",
						},
						{
							StreamID: binanceFundingTimeStreamID,
							Type:     "int192",
						},
						{
							StreamID: binanceFundingIntervalHoursStreamID,
							Type:     "int192",
						},
						{
							StreamID: deribitFundingRateStreamID,
							Type:     "int192",
						},
						{
							StreamID: deribitFundingTimeStreamID,
							Type:     "int192",
						},
						{
							StreamID: deribitFundingIntervalHoursStreamID,
							Type:     "int192",
						},
					},
				}),
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

		resultJSON := `{
	"benchmarkPrice": "2976.39",
	"baseMarketDepth": "1000.1212",
	"quoteMarketDepth": "998.5431",
	"marketStatus": "1",
	"binanceFundingRate": "1234.5678",
	"binanceFundingTime": "1630000000",
	"binanceFundingIntervalHours": "8",
	"deribitFundingRate": "5432.2345",
	"deribitFundingTime": "1630000000",
	"deribitFundingIntervalHours": "8"
}`

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

		rwaPipeline := fmt.Sprintf(`
dp          [type=bridge name="%s" requestData="{\\"data\\":{\\"data\\":\\"foo\\"}}"];

market_status_parse   [type=jsonparse path="result,marketStatus"];
market_status_decimal [type=multiply times=1 streamID=%d];

dp -> market_status_parse -> market_status_decimal;
`, bridgeName, marketStatusStreamID)

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
			createBridge(t, bridgeName, resultJSON, node.App.BridgeORM())
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
			t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles,
		)

		// NOTE: Wait for one of each type of report
		feedIDs := map[[32]byte]struct{}{
			dexBasedAssetFeedID:  {},
			rwaFeedID:            {},
			benchmarkPriceFeedID: {},
			fundingRateFeedID:    {},
		}

		for pckt := range packetCh {
			req := pckt.req
			assert.Equal(t, uint32(llotypes.ReportFormatEVMABIEncodeUnpacked), req.ReportFormat)
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
			assert.Equal(t, "0", reportElems["nativeFee"].(*big.Int).String())
			assert.Equal(t, "0", reportElems["linkFee"].(*big.Int).String())
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

			if len(feedIDs) == 0 {
				break
			}
		}
	})
}

func TestIntegration_LLO_stress_test_and_transmit_errors(t *testing.T) {
	t.Parallel()

	// logLevel: the log level to use for the nodes
	// setting a more verbose log level increases cpu usage significantly
	const logLevel = toml.LogLevel(zapcore.ErrorLevel)

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
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo", backend, bootstrapCSAKey, nil)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	t.Run("transmit queue does not grow unbounded", func(t *testing.T) {
		packets := make(chan *packet, 100000)
		serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
		serverPubKey := serverKey.PublicKey
		srv := NewMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), packets)

		serverURL := startMercuryServer(t, srv, clientPubKeys)

		donID := uint32(888333)
		streams := []Stream{ethStream, linkStream}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, streams, func(c *chainlink.Config) {
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
		// 2,000 channels should produce 2,000 reports per second
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
				t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles,
			)

			// NOTE: Wait for 40,000 reports (should take about 5 seconds) - 2,000 reports per second * 4 transmitters * 5 seconds
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
			err := db.GetContext(tests.Context(t), &cnt, "SELECT count(*) FROM llo_mercury_transmit_queue WHERE server_url = 'example.invalid'")
			require.NoError(t, err)
			assert.LessOrEqual(t, cnt, maxQueueSize, "persisted transmit queue size too large for node %d for failing server", i)

			// The succeeding server
			err = db.GetContext(tests.Context(t), &cnt, "SELECT count(*) FROM llo_mercury_transmit_queue WHERE server_url = $1", serverURL)
			require.NoError(t, err)
			assert.LessOrEqual(t, cnt, maxQueueSize, "persisted transmit queue size too large for node %d for succeeding server", i)
		}
	})
}

func TestIntegration_LLO_blue_green_lifecycle(t *testing.T) {
	t.Parallel()

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
		srv := NewMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), packetCh)

		serverURL := startMercuryServer(t, srv, clientPubKeys)

		donID := uint32(888333)
		streams := []Stream{ethStream, linkStream}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, streams, func(c *chainlink.Config) {
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
				t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles,
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
				t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles, blueDigest,
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
				highestObservationTs := uint32(0)
				highestValidAfterSeconds := uint32(0)
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
						if highestObservationTs > 0 {
							if digest == greenDigest && i == len(allReports[digest])-1 {
								// NOTE: This actually CHANGES on the staging
								// handover and can go backwards - the gapless
								// handover test is handled below
								break
							}
							assert.Equal(t, highestObservationTs, r.ValidAfterSeconds, "%d: (n-1)ObservationsTimestampSeconds->(n)ValidAfterSeconds should be gapless, got: %d vs %d", i, highestObservationTs, r.ValidAfterSeconds)
							assert.Greater(t, r.ObservationTimestampSeconds, highestObservationTs, "%d: overlapping/duplicate report ObservationTimestampSeconds, got: %d vs %d", i, r.ObservationTimestampSeconds, highestObservationTs)
							assert.Greater(t, r.ValidAfterSeconds, highestValidAfterSeconds, "%d: overlapping/duplicate report ValidAfterSeconds, got: %d vs %d", i, r.ValidAfterSeconds, highestValidAfterSeconds)
							assert.Less(t, r.ValidAfterSeconds, r.ObservationTimestampSeconds)
						}
						seenSeqNr = r.SeqNr
						highestObservationTs = r.ObservationTimestampSeconds
						highestValidAfterSeconds = r.ValidAfterSeconds
					}
				}
			}

			// Gapless handover
			assert.Less(t, finalBlueReport.ValidAfterSeconds, finalBlueReport.ObservationTimestampSeconds)
			assert.Equal(t, finalBlueReport.ObservationTimestampSeconds, initialPromotedGreenReport.ValidAfterSeconds)
			assert.Less(t, initialPromotedGreenReport.ValidAfterSeconds, initialPromotedGreenReport.ObservationTimestampSeconds)
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
				t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles, greenDigest,
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
			assert.Less(t, finalGreenReport.ValidAfterSeconds, finalGreenReport.ObservationTimestampSeconds)
			assert.Equal(t, finalGreenReport.ObservationTimestampSeconds, initialPromotedBlueReport.ValidAfterSeconds)
			assert.Less(t, initialPromotedBlueReport.ValidAfterSeconds, initialPromotedBlueReport.ObservationTimestampSeconds)
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

func setupNodes(t *testing.T, nNodes int, backend evmtypes.Backend, clientCSAKeys []csakey.KeyV2, streams []Stream, f func(*chainlink.Config)) (oracles []confighelper.OracleIdentityExtra, nodes []Node) {
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
				TransmitAccount:   ocr2types.Account(fmt.Sprintf("%x", transmitter[:])),
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
