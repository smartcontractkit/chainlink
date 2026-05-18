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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/destination_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/destination_verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	lloevm "github.com/smartcontractkit/chainlink/v2/core/services/llo/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	reportcodecv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
	mercuryverifier "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/verifier"
)

var (
	fNodes = uint8(1)
	nNodes = 4 // number of nodes (not including bootstrap)
)

func setupBlockchain(t *testing.T) (*bind.TransactOpts, *backends.SimulatedBackend, *configurator.Configurator, common.Address, *destination_verifier.DestinationVerifier, common.Address, *destination_verifier_proxy.DestinationVerifierProxy, common.Address, *channel_config_store.ChannelConfigStore, common.Address) {
	steve := testutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := core.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	backend.Commit()
	backend.Commit() // ensure starting block number at least 1

	// Configurator
	configuratorAddress, _, configurator, err := configurator.DeployConfigurator(steve, backend)
	require.NoError(t, err)

	// DestinationVerifierProxy
	verifierProxyAddr, _, verifierProxy, err := destination_verifier_proxy.DeployDestinationVerifierProxy(steve, backend)
	require.NoError(t, err)
	// DestinationVerifier
	verifierAddr, _, verifier, err := destination_verifier.DeployDestinationVerifier(steve, backend, verifierProxyAddr)
	require.NoError(t, err)
	// AddVerifier
	_, err = verifierProxy.SetVerifier(steve, verifierAddr)
	require.NoError(t, err)

	// ChannelConfigStore
	configStoreAddress, _, configStore, err := channel_config_store.DeployChannelConfigStore(steve, backend)
	require.NoError(t, err)

	backend.Commit()

	return steve, backend, configurator, configuratorAddress, verifier, verifierAddr, verifierProxy, verifierProxyAddr, configStore, configStoreAddress
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

func generateConfig(t *testing.T, oracles []confighelper.OracleIdentityExtra) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) {
	rawReportingPluginConfig := datastreamsllo.OffchainConfig{}
	reportingPluginConfig, err := rawReportingPluginConfig.Encode()
	require.NoError(t, err)

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err = ocr3confighelper.ContractSetConfigArgsForTests(
		2*time.Second,         // DeltaProgress
		20*time.Second,        // DeltaResend
		400*time.Millisecond,  // DeltaInitial
		1000*time.Millisecond, // DeltaRound
		500*time.Millisecond,  // DeltaGrace
		300*time.Millisecond,  // DeltaCertifiedCommitRequest
		1*time.Minute,         // DeltaStage
		100,                   // rMax
		[]int{len(oracles)},   // S
		oracles,
		reportingPluginConfig, // reportingPluginConfig []byte,
		0,                     // maxDurationQuery
		250*time.Millisecond,  // maxDurationObservation
		0,                     // maxDurationShouldAcceptAttestedReport
		0,                     // maxDurationShouldTransmitAcceptedReport
		int(fNodes),           // f
		onchainConfig,         // encoded onchain config
	)

	require.NoError(t, err)

	return
}

func setConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend *backends.SimulatedBackend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra) ocr2types.ConfigDigest {
	signers, _, _, _, offchainConfigVersion, offchainConfig := generateConfig(t, oracles)

	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	offchainTransmitters := make([][32]byte, nNodes)
	for i := 0; i < nNodes; i++ {
		offchainTransmitters[i] = nodes[i].ClientPubKey
	}
	donIDPadded := evm.DonIDToBytes32(donID)
	_, err = configurator.SetConfig(steve, donIDPadded, signerAddresses, offchainTransmitters, fNodes, offchainConfig, offchainConfigVersion, offchainConfig)
	require.NoError(t, err)

	// libocr requires a few confirmations to accept the config
	backend.Commit()
	backend.Commit()
	backend.Commit()
	backend.Commit()

	logs, err := backend.FilterLogs(testutils.Context(t), ethereum.FilterQuery{Addresses: []common.Address{configuratorAddress}, Topics: [][]common.Hash{[]common.Hash{mercury.FeedScopedConfigSet, donIDPadded}}})
	require.NoError(t, err)
	require.Len(t, logs, 1)

	cfg, err := mercury.ConfigFromLog(logs[0].Data)
	require.NoError(t, err)

	return cfg.ConfigDigest
}

func TestIntegration_LLO(t *testing.T) {
	testStartTimeStamp := time.Now()
	donID := uint32(995544)
	multiplier, err := decimal.NewFromString("1e18")
	require.NoError(t, err)
	expirationWindow := time.Hour / time.Second

	reqs := make(chan request)
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
	serverURL := startMercuryServer(t, srv, clientPubKeys)

	steve, backend, configurator, configuratorAddress, verifier, _, verifierProxy, _, configStore, configStoreAddress := setupBlockchain(t)
	fromBlock := 1

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo", backend, bootstrapCSAKey)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	t.Run("produces reports in v0.3 format", func(t *testing.T) {
		streams := []Stream{ethStream, linkStream, quoteStream1, quoteStream2}
		streamMap := make(map[uint32]Stream)
		for _, strm := range streams {
			streamMap[strm.id] = strm
		}

		// Setup oracle nodes
		var (
			oracles []confighelper.OracleIdentityExtra
			nodes   []Node
		)
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

		chainID := testutils.SimulatedChainID
		relayType := "evm"
		relayConfig := fmt.Sprintf(`
chainID = "%s"
fromBlock = %d
lloDonID = %d
`, chainID, fromBlock, donID)
		addBootstrapJob(t, bootstrapNode, configuratorAddress, "job-2", relayType, relayConfig)

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

		channelDefinitionsJSON, err := json.MarshalIndent(channelDefinitions, "", "  ")
		require.NoError(t, err)
		channelDefinitionsSHA := sha3.Sum256(channelDefinitionsJSON)

		// Set up channel definitions server
		channelDefinitionsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/channel-definitions", r.URL.Path)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err2 := w.Write(channelDefinitionsJSON)
			require.NoError(t, err2)
		}))
		t.Cleanup(channelDefinitionsServer.Close)

		// Set channel definitions
		_, err = configStore.SetChannelDefinitions(steve, donID, channelDefinitionsServer.URL+"/channel-definitions", channelDefinitionsSHA)
		require.NoError(t, err)
		backend.Commit()

		pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
donID = %d
channelDefinitionsContractAddress = "0x%x"
channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock)
		addOCRJobsEVMPremiumLegacy(t, streams, serverPubKey, serverURL, configuratorAddress, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

		// Set config on configurator
		setConfig(
			t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles,
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

					reportSigners, err2 := rv.Verify(mercuryverifier.SignedReport{
						RawRs:         v["rawRs"].([][32]byte),
						RawSs:         v["rawSs"].([][32]byte),
						RawVs:         v["rawVs"].([32]byte),
						ReportContext: v["reportContext"].([3][32]byte),
						Report:        v["report"].([]byte),
					}, fNodes, signerAddresses)
					require.NoError(t, err2)
					assert.GreaterOrEqual(t, len(reportSigners), int(fNodes+1))
					assert.Subset(t, signerAddresses, reportSigners)
				})

				t.Run(fmt.Sprintf("test on-chain verification - node %x", req.pk), func(t *testing.T) {
					_, err = verifierProxy.Verify(steve, req.req.Payload, []byte{})
					require.NoError(t, err)
				})

				t.Logf("oracle %x reported for 0x%x", req.pk, feedID)

				seen[feedID][req.pk] = struct{}{}
				if len(seen[feedID]) == nNodes {
					t.Logf("all oracles reported for 0x%x", feedID)
					delete(seen, feedID)
					if len(seen) == 0 {
						break // saw all oracles; success!
					}
				}
			}
		})

		t.Run("deleting LLO jobs cleans up resources", func(t *testing.T) {
			t.Skip("TODO - https://smartcontract-it.atlassian.net/browse/MERC-3653")
		})
	})
}
