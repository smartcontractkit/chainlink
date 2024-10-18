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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/wsrpc/credentials"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
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
	*backends.SimulatedBackend,
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
	genesisData := core.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	backend.Commit()
	backend.Commit() // ensure starting block number at least 1

	// Configurator
	configuratorAddress, _, configurator, err := configurator.DeployConfigurator(steve, backend)
	require.NoError(t, err)

	// DestinationVerifierProxy
	destinationVerifierProxyAddr, _, verifierProxy, err := destination_verifier_proxy.DeployDestinationVerifierProxy(steve, backend)
	require.NoError(t, err)
	// DestinationVerifier
	destinationVerifierAddr, _, destinationVerifier, err := destination_verifier.DeployDestinationVerifier(steve, backend, destinationVerifierProxyAddr)
	require.NoError(t, err)
	// AddVerifier
	_, err = verifierProxy.SetVerifier(steve, destinationVerifierAddr)
	require.NoError(t, err)

	// Legacy mercury verifier
	legacyVerifier, legacyVerifierAddr, legacyVerifierProxy, legacyVerifierProxyAddr := setupLegacyMercuryVerifier(t, steve, backend)

	// ChannelConfigStore
	configStoreAddress, _, configStore, err := channel_config_store.DeployChannelConfigStore(steve, backend)
	require.NoError(t, err)

	backend.Commit()

	return steve, backend, configurator, configuratorAddress, destinationVerifier, destinationVerifierAddr, verifierProxy, destinationVerifierProxyAddr, configStore, configStoreAddress, legacyVerifier, legacyVerifierAddr, legacyVerifierProxy, legacyVerifierProxyAddr
}

func setupLegacyMercuryVerifier(t *testing.T, steve *bind.TransactOpts, backend *backends.SimulatedBackend) (*verifier.Verifier, common.Address, *verifier_proxy.VerifierProxy, common.Address) {
	linkTokenAddress, _, linkToken, err := link_token_interface.DeployLinkToken(steve, backend)
	require.NoError(t, err)
	_, err = linkToken.Transfer(steve, steve.From, big.NewInt(1000))
	require.NoError(t, err)
	nativeTokenAddress, _, nativeToken, err := link_token_interface.DeployLinkToken(steve, backend)
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

func setLegacyConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend *backends.SimulatedBackend, legacyVerifier *verifier.Verifier, legacyVerifierAddr common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra) ocr2types.ConfigDigest {
	onchainConfig, err := (&datastreamsllo.EVMOnchainConfigCodec{}).Encode(datastreamsllo.OnchainConfig{
		Version:                 1,
		PredecessorConfigDigest: nil,
	})
	require.NoError(t, err)

	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig := generateConfig(t, oracles, onchainConfig)

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

func setStagingConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend *backends.SimulatedBackend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra, predecessorConfigDigest ocr2types.ConfigDigest) ocr2types.ConfigDigest {
	return setBlueGreenConfig(t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles, &predecessorConfigDigest)
}

func setProductionConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend *backends.SimulatedBackend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra) ocr2types.ConfigDigest {
	return setBlueGreenConfig(t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles, nil)
}

func setBlueGreenConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend *backends.SimulatedBackend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra, predecessorConfigDigest *ocr2types.ConfigDigest) ocr2types.ConfigDigest {
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
	logs, err := backend.FilterLogs(testutils.Context(t), ethereum.FilterQuery{Addresses: []common.Address{configuratorAddress}, Topics: [][]common.Hash{[]common.Hash{topic, donIDPadded}}})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(logs), 1)

	cfg, err := mercury.ConfigFromLog(logs[len(logs)-1].Data)
	require.NoError(t, err)

	return cfg.ConfigDigest
}

func promoteStagingConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend *backends.SimulatedBackend, configurator *configurator.Configurator, configuratorAddress common.Address, isGreenProduction bool) {
	donIDPadded := llo.DonIDToBytes32(donID)
	_, err := configurator.PromoteStagingConfig(steve, donIDPadded, isGreenProduction)
	require.NoError(t, err)

	// libocr requires a few confirmations to accept the config
	backend.Commit()
	backend.Commit()
	backend.Commit()
	backend.Commit()
}

func TestIntegration_LLO(t *testing.T) {
	testStartTimeStamp := time.Now()
	multiplier := decimal.New(1, 18)
	expirationWindow := time.Hour / time.Second

	reqs := make(chan request, 100000)
	serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	serverPubKey := serverKey.PublicKey
	srv := NewMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), reqs)

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)
	for i := 0; i < nNodes; i++ {
		k := big.NewInt(int64(i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	steve, backend, configurator, configuratorAddress, verifier, _, verifierProxy, _, configStore, configStoreAddress, legacyVerifier, legacyVerifierAddr, _, _ := setupBlockchain(t)
	fromBlock := 1

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo", backend, bootstrapCSAKey)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	t.Run("using legacy verifier configuration contract, produces reports in v0.3 format", func(t *testing.T) {
		serverURL := startMercuryServer(t, srv, clientPubKeys)

		donID := uint32(995544)
		streams := []Stream{ethStream, linkStream, quoteStream1, quoteStream2}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, streams)

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

				t.Run(fmt.Sprintf("emulate mercury server verifying report (local verification) - node %x", req.pk), func(t *testing.T) {
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
				})

				t.Run(fmt.Sprintf("test on-chain verification - node %x", req.pk), func(t *testing.T) {
					t.Run("destination verifier", func(t *testing.T) {
						_, err = verifierProxy.Verify(steve, req.req.Payload, []byte{})
						require.NoError(t, err)
					})
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

	t.Run("Blue/Green lifecycle (using JSON report format)", func(t *testing.T) {
		serverURL := startMercuryServer(t, srv, clientPubKeys)

		donID := uint32(888333)
		streams := []Stream{ethStream, linkStream}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, streams)

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
		t.Run("start off with blue=production, green=staging (specimen reports)", func(t *testing.T) {
			// Set config on configurator
			blueDigest = setProductionConfig(
				t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles,
			)

			// NOTE: Wait until blue produces a report

			for req := range reqs {
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.req.Payload)
				require.NoError(t, err)

				allReports[r.ConfigDigest] = append(allReports[r.ConfigDigest], r)

				assert.Equal(t, blueDigest, r.ConfigDigest)
				assert.False(t, r.Specimen)
				assert.Len(t, r.Values, 1)
				assert.Equal(t, "2976.39", r.Values[0].(*datastreamsllo.Decimal).String())
				break
			}
		})
		t.Run("setStagingConfig does not affect production", func(t *testing.T) {
			greenDigest = setStagingConfig(
				t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles, blueDigest,
			)

			// NOTE: Wait until green produces the first "specimen" report

			for req := range reqs {
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.req.Payload)
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
		})
		t.Run("promoteStagingConfig flow has clean and gapless hand off from old production to newly promoted staging instance, leaving old production instance in 'retired' state", func(t *testing.T) {
			promoteStagingConfig(t, donID, steve, backend, configurator, configuratorAddress, false)

			// NOTE: Wait for first non-specimen report for the newly promoted (green) instance

			for req := range reqs {
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.req.Payload)
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
		})
		t.Run("retired instance does not produce reports", func(t *testing.T) {
			// NOTE: Wait for five "green" reports to be produced and assert no "blue" reports

			i := 0
			for req := range reqs {
				i++
				if i == 5 {
					break
				}
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.req.Payload)
				require.NoError(t, err)

				allReports[r.ConfigDigest] = append(allReports[r.ConfigDigest], r)
				assert.False(t, r.Specimen)
				assert.Equal(t, greenDigest, r.ConfigDigest)
			}
		})
		t.Run("setStagingConfig replaces 'retired' instance with new config and starts producing specimen reports again", func(t *testing.T) {
			blueDigest = setStagingConfig(
				t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles, greenDigest,
			)

			// NOTE: Wait until blue produces the first "specimen" report

			for req := range reqs {
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.req.Payload)
				require.NoError(t, err)

				allReports[r.ConfigDigest] = append(allReports[r.ConfigDigest], r)
				if r.Specimen {
					assert.Equal(t, blueDigest, r.ConfigDigest)
					break
				}
				assert.Equal(t, greenDigest, r.ConfigDigest)
			}
		})
		t.Run("promoteStagingConfig swaps the instances again", func(t *testing.T) {
			// TODO: Check that once an instance enters 'retired' state, it
			// doesn't produce reports or bother making observations
			promoteStagingConfig(t, donID, steve, backend, configurator, configuratorAddress, true)

			// NOTE: Wait for first non-specimen report for the newly promoted (blue) instance

			for req := range reqs {
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.req.Payload)
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
		})
		t.Run("adding a new channel definition is picked up on the fly", func(t *testing.T) {
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

			for req := range reqs {
				_, _, r, _, err := (datastreamsllo.JSONReportCodec{}).UnpackDecode(req.req.Payload)
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
		})
		t.Run("deleting the jobs turns off oracles and cleans up resources", func(t *testing.T) {
			t.Skip("TODO - MERC-3524")
		})
		t.Run("adding new jobs again picks up the correct configs", func(t *testing.T) {
			t.Skip("TODO - MERC-3524")
		})
	})
}

func setupNodes(t *testing.T, nNodes int, backend *backends.SimulatedBackend, clientCSAKeys []csakey.KeyV2, streams []Stream) (oracles []confighelper.OracleIdentityExtra, nodes []Node) {
	ports := freeport.GetN(t, nNodes)
	for i := 0; i < nNodes; i++ {
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_streams_%d", i), backend, clientCSAKeys[i])

		nodes = append(nodes, Node{
			app, transmitter, kb, observedLogs,
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
