package llo_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
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
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	evmtestutils "github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/destination_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/destination_verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
)

// Copied from core/services/ocr2/plugins/llo/integration_test.go

var (
	fNodes = uint8(1)
	nNodes = 1 // number of nodes (not including bootstrap)
)

func TestIntegration_secondary_feed_transmission(t *testing.T) {
	// testStartTimeStamp := time.Now()
	// multiplier := decimal.New(1, 18)
	// expirationWindow := time.Hour / time.Second

	const salt = 100

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)
	for i := 0; i < nNodes; i++ {
		k := big.NewInt(int64(salt + i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	steve, backend, _, _, verifier, _, verifierProxy, _, configStore, configStoreAddress := setupBlockchain(t)
	// fromBlock := 1
	fmt.Printf("here: %s\n", verifierProxy.Address().Hex())
	fmt.Printf("here: %s\n", configStore.Address().Hex())
	fmt.Printf("here: %s\n", configStoreAddress.Hex())
	fmt.Printf("steve: %s\n", steve.From.String())
	fmt.Printf("verifier: %s\n", verifier.Address().Hex())

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_svr", backend, bootstrapCSAKey, nil)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	fmt.Printf("Bootstrap node: %s\n", bootstrapPeerID)
	fmt.Printf("Bootstrap node: %s\n", bootstrapNode.ClientPubKey.String())

	// t.Run("sends a transmission to the secondary feed", func(t *testing.T) {

	// 	donID := uint32(995544)

	// 	// Setup oracle nodes
	// 	oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {

	// 	})

	// 	chainID := testutils.SimulatedChainID
	// 	relayType := "evm"
	// 	relayConfig := fmt.Sprintf(`
	// 		chainID = "%s"
	// 		fromBlock = %d
	// 		lloDonID = %d
	// 		lloConfigMode = "mercury"
	// `, chainID, fromBlock, donID)
	// 	addBootstrapJob(t, bootstrapNode, legacyVerifierAddr, "job-2", relayType, relayConfig)

	// 	// Channel definitions
	// 	channelDefinitions := llotypes.ChannelDefinitions{
	// 		1: {
	// 			ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
	// 			Streams: []llotypes.Stream{
	// 				{
	// 					StreamID:   ethStreamID,
	// 					Aggregator: llotypes.AggregatorMedian,
	// 				},
	// 			},
	// 			Opts: llotypes.ChannelOpts([]byte(fmt.Sprintf(`{"baseUSDFee":"0.1","expirationWindow":%d,"feedId":"0x%x","multiplier":"%s"}`, expirationWindow, quoteStreamFeedID1, multiplier.String()))),
	// 		},
	// 		2: {
	// 			ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
	// 			Streams: []llotypes.Stream{
	// 				{
	// 					StreamID:   ethStreamID,
	// 					Aggregator: llotypes.AggregatorMedian,
	// 				},
	// 			},
	// 			Opts: llotypes.ChannelOpts([]byte(fmt.Sprintf(`{"baseUSDFee":"0.1","expirationWindow":%d,"feedId":"0x%x","multiplier":"%s"}`, expirationWindow, quoteStreamFeedID2, multiplier.String()))),
	// 		},
	// 	}

	// 	pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
	// donID = %d
	// channelDefinitionsContractAddress = "0x%x"
	// channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock)
	// 	addOCRJobsEVMPremiumLegacy(t, streams, serverPubKey, serverURL, legacyVerifierAddr, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

	// 	// Set config on configurator
	// 	setLegacyConfig(
	// 		t, donID, steve, backend, legacyVerifier, legacyVerifierAddr, nodes, oracles,
	// 	)

	// 	// Set config on the destination verifier
	// 	signerAddresses := make([]common.Address, len(oracles))
	// 	for i, oracle := range oracles {
	// 		signerAddresses[i] = common.BytesToAddress(oracle.OracleIdentity.OnchainPublicKey)
	// 	}
	// 	{
	// 		recipientAddressesAndWeights := []destination_verifier.CommonAddressAndWeight{}

	// 		_, err := verifier.SetConfig(steve, signerAddresses, fNodes, recipientAddressesAndWeights)
	// 		require.NoError(t, err)
	// 		backend.Commit()
	// 	}

	// 	// Expect at least one report per feed from each oracle
	// 	seen := make(map[[32]byte]map[credentials.StaticSizedPublicKey]struct{})
	// 	for _, cd := range channelDefinitions {
	// 		var opts lloevm.ReportFormatEVMPremiumLegacyOpts
	// 		err := json.Unmarshal(cd.Opts, &opts)
	// 		require.NoError(t, err)
	// 		// feedID will be deleted when all n oracles have reported
	// 		seen[opts.FeedID] = make(map[credentials.StaticSizedPublicKey]struct{}, nNodes)
	// 	}
	// 	for req := range reqs {
	// 		assert.Equal(t, uint32(llotypes.ReportFormatEVMPremiumLegacy), req.req.ReportFormat)
	// 		v := make(map[string]interface{})
	// 		err := mercury.PayloadTypes.UnpackIntoMap(v, req.req.Payload)
	// 		require.NoError(t, err)
	// 		report, exists := v["report"]
	// 		if !exists {
	// 			t.Fatalf("expected payload %#v to contain 'report'", v)
	// 		}
	// 		reportElems := make(map[string]interface{})
	// 		err = reportcodecv3.ReportTypes.UnpackIntoMap(reportElems, report.([]byte))
	// 		require.NoError(t, err)

	// 		feedID := reportElems["feedId"].([32]uint8)

	// 		if _, exists := seen[feedID]; !exists {
	// 			continue // already saw all oracles for this feed
	// 		}

	// 		var expectedBm, expectedBid, expectedAsk *big.Int
	// 		if feedID == quoteStreamFeedID1 {
	// 			expectedBm = quoteStream1.baseBenchmarkPrice.Mul(multiplier).BigInt()
	// 			expectedBid = quoteStream1.baseBid.Mul(multiplier).BigInt()
	// 			expectedAsk = quoteStream1.baseAsk.Mul(multiplier).BigInt()
	// 		} else if feedID == quoteStreamFeedID2 {
	// 			expectedBm = quoteStream2.baseBenchmarkPrice.Mul(multiplier).BigInt()
	// 			expectedBid = quoteStream2.baseBid.Mul(multiplier).BigInt()
	// 			expectedAsk = quoteStream2.baseAsk.Mul(multiplier).BigInt()
	// 		} else {
	// 			t.Fatalf("unrecognized feedID: 0x%x", feedID)
	// 		}

	// 		assert.GreaterOrEqual(t, reportElems["validFromTimestamp"].(uint32), uint32(testStartTimeStamp.Unix()))
	// 		assert.GreaterOrEqual(t, int(reportElems["observationsTimestamp"].(uint32)), int(testStartTimeStamp.Unix()))
	// 		assert.Equal(t, "33597747607000", reportElems["nativeFee"].(*big.Int).String())
	// 		assert.Equal(t, "7547169811320755", reportElems["linkFee"].(*big.Int).String())
	// 		assert.Equal(t, reportElems["observationsTimestamp"].(uint32)+uint32(expirationWindow), reportElems["expiresAt"].(uint32))
	// 		assert.Equal(t, expectedBm.String(), reportElems["benchmarkPrice"].(*big.Int).String())
	// 		assert.Equal(t, expectedBid.String(), reportElems["bid"].(*big.Int).String())
	// 		assert.Equal(t, expectedAsk.String(), reportElems["ask"].(*big.Int).String())

	// 		// emulate mercury server verifying report (local verification)
	// 		{
	// 			rv := mercuryverifier.NewVerifier()

	// 			reportSigners, err := rv.Verify(mercuryverifier.SignedReport{
	// 				RawRs:         v["rawRs"].([][32]byte),
	// 				RawSs:         v["rawSs"].([][32]byte),
	// 				RawVs:         v["rawVs"].([32]byte),
	// 				ReportContext: v["reportContext"].([3][32]byte),
	// 				Report:        v["report"].([]byte),
	// 			}, fNodes, signerAddresses)
	// 			require.NoError(t, err)
	// 			assert.GreaterOrEqual(t, len(reportSigners), int(fNodes+1))
	// 			assert.Subset(t, signerAddresses, reportSigners)
	// 		}

	// 	}
	// })
}

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

	// ChannelConfigStore
	configStoreAddress, _, configStore, err := channel_config_store.DeployChannelConfigStore(steve, backend.Client())
	require.NoError(t, err)

	backend.Commit()

	return steve, backend, configurator, configuratorAddress, destinationVerifier, destinationVerifierAddr, verifierProxy, destinationVerifierProxyAddr, configStore, configStoreAddress
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

func setStagingConfig(t *testing.T, donID uint32, steve *bind.TransactOpts, backend evmtypes.Backend, configurator *configurator.Configurator, configuratorAddress common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra, predecessorConfigDigest ocr2types.ConfigDigest) ocr2types.ConfigDigest {
	return setBlueGreenConfig(t, donID, steve, backend, configurator, configuratorAddress, nodes, oracles, &predecessorConfigDigest)
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
				TransmitAccount:   ocr2types.Account(fmt.Sprintf("%x", transmitter[:])),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}
	return
}

func mustNewType(t string) abi.Type {
	result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
	if err != nil {
		panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
	}
	return result
}
