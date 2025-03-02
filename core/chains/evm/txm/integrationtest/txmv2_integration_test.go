package llo_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/sdk/freeport"
	clcommonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

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
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// Copied from core/services/ocr2/plugins/llo/integration_test.go + svr-contracts/test/e2e/svr_test.go

/*
	Steps to run:
	* `docker run --name cl-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dbname -p 5432:5432 -d postgres`
	* `make setup-testdb`
	* `CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -run ^TestIntegration_secondary_feed_transmission$ github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/integrationtest -v`
*/

var (
	fNodes = uint8(1)
)

func TestIntegration_secondary_feed_transmission(t *testing.T) {
	// testStartTimeStamp := time.Now()
	// multiplier := decimal.New(1, 18)
	// expirationWindow := time.Hour / time.Second

	const salt = 100

	k := big.NewInt(int64(salt))
	key := csakey.MustNewV2XXXTestingOnly(k)
	clientCSAKey := key

	steve, backend, _, _, verifier, _, verifierProxy, _, configStore, configStoreAddress := setupBlockchain(t)
	fromBlock := 1
	fmt.Printf("here: %s\n", verifierProxy.Address().Hex())
	fmt.Printf("here: %s\n", configStore.Address().Hex())
	fmt.Printf("here: %s\n", configStoreAddress.Hex())
	fmt.Printf("steve: %s\n", steve.From.String())
	fmt.Printf("verifier: %s\n", verifier.Address().Hex())

	// donID := uint32(995544)

	// Setup the node
	port := freeport.GetOne(t)
	app, _, transmitter, kb, observedLogs := setupNode(t, port, "oracle_svr", backend, clientCSAKey, nil) // TODO(gg): fix db name?
	node := Node{app, transmitter, kb, observedLogs}
	fmt.Printf("created node with transmitter %#v\n", transmitter.String())

	// CreateTxKey creates a tx key on the Chainlink node
	// func (c *ChainlinkClient) CreateTxKey(chain string, chainId string) (*TxKey, *http.Response, error) {
	// 	txKey := &TxKey{}
	// 	framework.L.Info().Str(NodeURL, c.Config.URL).Msg("Creating Tx Key")
	// 	resp, err := c.APIClient.R().
	// 		SetPathParams(map[string]string{
	// 			"chain": chain,
	// 		}).
	// 		SetQueryParam("evmChainID", chainId).
	// 		SetResult(txKey).
	// 		Post("/v2/keys/{chain}")
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// 	return txKey, resp.RawResponse, err
	// }

	primaryTransmitterKey, err := node.App.GetKeyStore().Eth().Create(context.Background(), big.NewInt(int64(1337)))
	require.NoErrorf(t, err, "could not create primary transmitter key")
	secondaryTransmitterKey, err := node.App.GetKeyStore().Eth().Create(context.Background(), big.NewInt(int64(1337)))
	require.NoErrorf(t, err, "could not create secondary transmitter key")

	keys, err := node.App.GetKeyStore().Eth().GetAll(context.Background())
	require.NoError(t, err, "could not get node's eth keystest")

	fmt.Printf("Keys are %#v\n", keys)

	// offchainPublicKey, err := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
	// require.NoError(t, err)
	// oracles = append(oracles, confighelper.OracleIdentityExtra{
	// 	OracleIdentity: confighelper.OracleIdentity{
	// 		OnchainPublicKey:  offchainPublicKey,
	// 		TransmitAccount:   ocr2types.Account(fmt.Sprintf("%x", transmitter[:])),
	// 		OffchainPublicKey: kb.OffchainPublicKey(),
	// 		PeerID:            peerID,
	// 	},
	// 	ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
	// })

	// chainID := testutils.SimulatedChainID

	// TOOD(gg) deploy dualAggContracts
	dualAggContractsAddresses := []string{"0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C"}

	for _, contractAddress := range dualAggContractsAddresses {
		fmt.Printf("Creating feed for %s\n", contractAddress)
		firstKey := node.KeyBundle.Raw().Key().ID()

		// TODO(gg): put this into the actual job
		// job := fmt.Sprintf(oevJobSpec, feedNr, uuid.New().String(), contractAddress.Addresses[0].String(), "evm", firstKey, primaryAddresses[i], bootstrapPeerID.Data[0].Attributes.PeerID,
		// strings.TrimPrefix(node.DockerP2PUrl, "http://"), chainID, fromBlock, contractAddress.Addresses[0].String(), secondaryAddresses[i])

		// [relayConfig]
		// chainID = %s
		// fromBlock = %d
		// enableDualTransmission = true

		// [relayConfig.dualTransmission]
		// contractAddress = "%s"
		// transmitterAddress = "%s"

		// [relayConfig.dualTransmission.meta]
		// hint = [ "calldata" ]
		// refund = [ "0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C:90" ]

		// [pluginConfig]
		// juelsPerFeeCoinSource = """
		// juels_per_fee_coin [type="sum" values=<[0]>];
		// """
		// `

		jb := &job.Job{
			Type:              job.OffchainReporting2,
			SchemaVersion:     1,
			Name:              null.StringFrom("SVR job 1"),
			CronSpec:          &job.CronSpec{CronSchedule: "@every 1s"},
			PipelineSpec:      &pipeline.Spec{},
			ExternalJobID:     uuid.New(),
			ForwardingAllowed: true,
			MaxTaskDuration:   *models.NewInterval(0 * time.Second),
			OCR2OracleSpec: &job.OCR2OracleSpec{
				ContractID:           contractAddress,
				Relay:                "evm",
				OCRKeyBundleID:       null.StringFrom(firstKey),
				PluginType:           clcommonTypes.Median,
				TransmitterID:        null.StringFrom(primaryTransmitterKey.Address.Hex()),
				AllowNoBootstrappers: true,
				P2PV2Bootstrappers:   []string{}, // bootstrapPeerID.Data[0].Attributes.PeerID, needed?
				RelayConfig: map[string]any{
					"chainID":                "1337",
					"fromBlock":              fromBlock,
					"enableDualTransmission": true,
					"dualTransmission": map[string]any{
						"contractAddress":    contractAddress,
						"transmitterAddress": secondaryTransmitterKey.Address.Hex(),
						"meta": map[string]any{
							"hint":   []any{"calldata"},
							"refund": []any{"0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C:90"},
						},
					},
				},
				PluginConfig: map[string]any{
					"juelsPerFeeCoinSource": "juels_per_fee_coin [type=\"sum\" values=<[0]>]",
				},
			},
		}
		// err := helper.pipelineHelper.Jrm.CreateJob(testutils.Context(t), jb)
		err := node.App.AddJobV2(context.Background(), jb)
		require.NoError(t, err, "Failed to create feed job")
	}

	// relayType := "evm"
	// relayConfig := fmt.Sprintf(`
	// 			chainID = "%s"
	// 			fromBlock = %d
	// 	`, chainID, fromBlock, donID)
	// addBootstrapJob(t, bootstrapNode, legacyVerifierAddr, "job-2", relayType, relayConfig)

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
	offchainTransmitters := make([][32]byte, 1)
	for i := 0; i < 1; i++ {
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

func mustNewType(t string) abi.Type {
	result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
	if err != nil {
		panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
	}
	return result
}

var oevJobSpec = `
type = "offchainreporting2"
schemaVersion = 1
name = "OEV job %d"
externalJobID = "%s"
forwardingAllowed = true
maxTaskDuration = "0s"
contractID = "%s"
relay = "%s"
ocrKeyBundleID = "%s"
pluginType = "median"
transmitterID = "%s"
p2pv2Bootstrappers = ["%s@%s"]

observationSource = """
 //randomness
    val1 [type="memo" value="10"]
    val2 [type="memo" value="20"]
    val3 [type="memo" value="30"]
    val4 [type="memo" value="40"]
    val5 [type="memo" value="50"]
    val6 [type="memo" value="60"]
    val7 [type="memo" value="70"]
    val8 [type="memo" value="80"]
    val9 [type="memo" value="90"]

    random1 [type="any"]
    random2 [type="any"]
    random3 [type="any"]

    val1 -> random1
    val2 -> random2
    val3 -> random3
    val4 -> random1
    val5 -> random2
    val6 -> random3
    val7 -> random1
    val8 -> random2
    val9 -> random3


    // data source 1
    ds1_multiply [type="multiply" times=100]

     // data source 2
    ds2_multiply [type="multiply" times=100]


    // data source 3
    ds3_multiply [type="multiply" times=100]


    random1 -> ds1_multiply -> answer
    random2 -> ds2_multiply -> answer
    random3 -> ds3_multiply -> answer

    answer [type=median]
"""

[relayConfig]
chainID = %s
fromBlock = %d
enableDualTransmission = true

[relayConfig.dualTransmission]
contractAddress = "%s"
transmitterAddress = "%s"

[relayConfig.dualTransmission.meta]
hint = [ "calldata" ]
refund = [ "0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C:90" ]

[pluginConfig]
juelsPerFeeCoinSource = """
juels_per_fee_coin [type="sum" values=<[0]>];
"""
`
