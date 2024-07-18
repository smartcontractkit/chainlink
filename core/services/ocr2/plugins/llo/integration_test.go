package llo_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

var (
	fNodes           = uint8(1)
	nNodes           = 4 // number of nodes (not including bootstrap)
	multiplier int64 = 100000000
)

func setupBlockchain(t *testing.T) (*bind.TransactOpts, *backends.SimulatedBackend, *channel_verifier.ChannelVerifier, common.Address, *channel_config_store.ChannelConfigStore, common.Address) {
	steve := testutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := core.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	backend.Commit()
	backend.Commit() // ensure starting block number at least 1

	// Deploy contracts
	verifierProxyAddr, _, _, err := verifier_proxy.DeployVerifierProxy(steve, backend, common.Address{}) // zero address for access controller disables access control
	require.NoError(t, err)

	verifierAddress, _, verifierContract, err := channel_verifier.DeployChannelVerifier(steve, backend, verifierProxyAddr)
	require.NoError(t, err)
	configStoreAddress, _, configStoreContract, err := channel_config_store.DeployChannelConfigStore(steve, backend)
	require.NoError(t, err)

	backend.Commit()

	return steve, backend, verifierContract, verifierAddress, configStoreContract, configStoreAddress
}

type Stream struct {
	id                 uint32
	baseBenchmarkPrice *big.Int
}

var (
	btcStream = Stream{
		id:                 51,
		baseBenchmarkPrice: big.NewInt(20_000 * multiplier),
	}
	ethStream = Stream{
		id:                 52,
		baseBenchmarkPrice: big.NewInt(1_568 * multiplier),
	}
	linkStream = Stream{
		id:                 53,
		baseBenchmarkPrice: big.NewInt(7150 * multiplier / 1000),
	}
	dogeStream = Stream{
		id:                 54,
		baseBenchmarkPrice: big.NewInt(2_020 * multiplier),
	}
)

func TestIntegration_LLO(t *testing.T) {
	t.Skip("waiting on https://github.com/smartcontractkit/chainlink/pull/13780")
	//     testStartTimeStamp := uint32(time.Now().Unix())

	//     const fromBlock = 1 // cannot use zero, start from block 1

	//     // streams
	//     streams := []Stream{btcStream, ethStream, linkStream, dogeStream}
	//     streamMap := make(map[uint32]Stream)
	//     for _, strm := range streams {
	//         streamMap[strm.id] = strm
	//     }

	//     reqs := make(chan request)
	//     serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	//     serverPubKey := serverKey.PublicKey
	//     srv := NewMercuryServer(t, ed25519.PrivateKey(serverKey.Raw()), reqs, nil)

	//     clientCSAKeys := make([]csakey.KeyV2, nNodes)
	//     clientPubKeys := make([]ed25519.PublicKey, nNodes)
	//     for i := 0; i < nNodes; i++ {
	//         k := big.NewInt(int64(i))
	//         key := csakey.MustNewV2XXXTestingOnly(k)
	//         clientCSAKeys[i] = key
	//         clientPubKeys[i] = key.PublicKey
	//     }
	//     serverURL := startMercuryServer(t, srv, clientPubKeys)
	//     chainID := testutils.SimulatedChainID

	//     steve, backend, verifierContract, verifierAddress, configStoreContract, configStoreAddress := setupBlockchain(t)

	//     // Setup bootstrap
	//     bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	//     bootstrapNodePort := freeport.GetOne(t)
	//     appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_mercury", backend, bootstrapCSAKey)
	//     bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	//     // Setup oracle nodes
	//     var (
	//         oracles []confighelper.OracleIdentityExtra
	//         nodes   []Node
	//     )
	//     ports := freeport.GetN(t, nNodes)
	//     for i := 0; i < nNodes; i++ {
	//         app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_streams_%d", i), backend, clientCSAKeys[i])

	//         nodes = append(nodes, Node{
	//             app, transmitter, kb, observedLogs,
	//         })
	//         offchainPublicKey, _ := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
	//         oracles = append(oracles, confighelper.OracleIdentityExtra{
	//             OracleIdentity: confighelper.OracleIdentity{
	//                 OnchainPublicKey:  offchainPublicKey,
	//                 TransmitAccount:   ocr2types.Account(fmt.Sprintf("%x", transmitter[:])),
	//                 OffchainPublicKey: kb.OffchainPublicKey(),
	//                 PeerID:            peerID,
	//             },
	//             ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
	//         })
	//     }

	//     // Commit blocks to finality depth to ensure LogPoller has finalized blocks to read from
	//     ch, err := nodes[0].App.GetRelayers().LegacyEVMChains().Get(testutils.SimulatedChainID.String())
	//     require.NoError(t, err)
	//     finalityDepth := ch.Config().EVM().FinalityDepth()
	//     for i := 0; i < int(finalityDepth); i++ {
	//         backend.Commit()
	//     }

	//     configDigest := setConfig(t, steve, backend, verifierContract, verifierAddress, nodes, oracles)
	//     channelDefinitions := setChannelDefinitions(t, steve, backend, configStoreContract, streams)

	//     relayType := "evm"
	//     relayConfig := fmt.Sprintf(`chainID = %s
	// fromBlock = %d`, chainID.String(), fromBlock)
	//     addBootstrapJob(t, bootstrapNode, verifierAddress, "job-1", relayType, relayConfig)

	//     pluginConfig := fmt.Sprintf(`serverURL = "%s"
	// serverPubKey = "%x"
	// channelDefinitionsContractFromBlock = %d
	// channelDefinitionsContractAddress = "%s"`, serverURL, serverPubKey, fromBlock, configStoreAddress.String())
	//     addOCRJobs(t, streams, serverPubKey, serverURL, verifierAddress, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)
	//     t.Run("receives at least one report per feed from each oracle when EAs are at 100% reliability", func(t *testing.T) {
	//         // Expect at least one report per channel from each oracle (keyed by transmitter ID)
	//         seen := make(map[ocr2types.Account]map[llotypes.ChannelID]struct{})

	//         for channelID, defn := range channelDefinitions {
	//             t.Logf("Expect report for channel ID %x (definition: %#v)", channelID, defn)
	//         }
	//         for _, o := range oracles {
	//             t.Logf("Expect report from oracle %s", o.OracleIdentity.TransmitAccount)
	//             seen[o.OracleIdentity.TransmitAccount] = make(map[llotypes.ChannelID]struct{})
	//         }
	//         for req := range reqs {
	//             if _, exists := seen[req.TransmitterID()]; !exists {
	//                 // oracle already reported on all channels; discard
	//                 // if this test timeouts, check for expected transmitter ID
	//                 continue
	//             }

	//             v := make(map[string]interface{})
	//             err := llo.PayloadTypes.UnpackIntoMap(v, req.req.Payload)
	//             require.NoError(t, err)
	//             report, exists := v["report"]
	//             if !exists {
	//                 t.Fatalf("FAIL: expected payload %#v to contain 'report'", v)
	//             }

	//             t.Logf("Got report from oracle %s with format: %d", req.pk, req.req.ReportFormat)

	//             var r datastreamsllo.Report

	//             switch req.req.ReportFormat {
	//             case uint32(llotypes.ReportFormatJSON):
	//                 t.Logf("Got report (JSON) from oracle %x: %s", req.pk, string(report.([]byte)))
	//                 var err error
	//                 r, err = (datastreamsllo.JSONReportCodec{}).Decode(report.([]byte))
	//                 require.NoError(t, err, "expected valid JSON")
	//             case uint32(llotypes.ReportFormatEVM):
	//                 t.Logf("Got report (EVM) from oracle %s: 0x%x", req.pk, report.([]byte))
	//                 var err error
	//                 r, err = (lloevm.ReportCodec{}).Decode(report.([]byte))
	//                 require.NoError(t, err, "expected valid EVM encoding")
	//             default:
	//                 t.Fatalf("FAIL: unexpected report format: %q", req.req.ReportFormat)
	//             }

	//             assert.Equal(t, configDigest, r.ConfigDigest)
	//             assert.Equal(t, uint64(0x2ee634951ef71b46), r.ChainSelector)
	//             assert.GreaterOrEqual(t, r.SeqNr, uint64(1))
	//             assert.GreaterOrEqual(t, r.ValidAfterSeconds, testStartTimeStamp)
	//             assert.Equal(t, r.ValidAfterSeconds+1, r.ValidUntilSeconds)

	//             // values
	//             defn, exists := channelDefinitions[r.ChannelID]
	//             require.True(t, exists, "expected channel ID to be in channelDefinitions")

	//             require.Equal(t, len(defn.StreamIDs), len(r.Values))

	//             for i, strmID := range defn.StreamIDs {
	//                 strm, exists := streamMap[strmID]
	//                 require.True(t, exists, "invariant violation: expected stream ID to be present")
	//                 assert.InDelta(t, strm.baseBenchmarkPrice.Int64(), r.Values[i].Int64(), 5000000)
	//             }

	//             assert.False(t, r.Specimen)

	//             seen[req.TransmitterID()][r.ChannelID] = struct{}{}
	//             t.Logf("Got report from oracle %s with channel: %x)", req.TransmitterID(), r.ChannelID)

	//             if _, exists := seen[req.TransmitterID()]; exists && len(seen[req.TransmitterID()]) == len(channelDefinitions) {
	//                 t.Logf("All channels reported for oracle with transmitterID %s", req.TransmitterID())
	//                 delete(seen, req.TransmitterID())
	//             }
	//             if len(seen) == 0 {
	//                 break // saw all oracles; success!
	//             }

	//             // bit of a hack here but shouldn't hurt anything, we wanna dump
	//             // `seen` before the test ends to aid in debugging test failures
	//             if d, ok := t.Deadline(); ok {
	//                 select {
	//                 case <-time.After(time.Until(d.Add(-100 * time.Millisecond))):
	//                     if len(seen) > 0 {
	//                         t.Fatalf("FAILED: ERROR: missing expected reports: %#v\n", seen)
	//                     }
	//                 default:
	//                 }
	//             }
	//         }
	//     })

	//     // TODO: test verification
	// }

	// func generateConfig(t *testing.T, nodes []Node, oracles []confighelper.OracleIdentityExtra) (
	//     signers []types.OnchainPublicKey,
	//     transmitters []types.Account,
	//     f uint8,
	//     onchainConfig []byte,
	//     offchainConfigVersion uint64,
	//     offchainConfig []byte,
	// ) {
	//     // Setup config on contract
	//     rawOnchainConfig := llo.OnchainConfig{}
	//     onchainConfig, err := (&llo.JSONOnchainConfigCodec{}).Encode(rawOnchainConfig)
	//     require.NoError(t, err)

	//     rawReportingPluginConfig := datastreamsllo.OffchainConfig{}
	//     reportingPluginConfig, err := rawReportingPluginConfig.Encode()
	//     require.NoError(t, err)

	//     signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err = ocr3confighelper.ContractSetConfigArgsForTests(
	//         2*time.Second,         // DeltaProgress
	//         20*time.Second,        // DeltaResend
	//         400*time.Millisecond,  // DeltaInitial
	//         1000*time.Millisecond, // DeltaRound
	//         500*time.Millisecond,  // DeltaGrace
	//         300*time.Millisecond,  // DeltaCertifiedCommitRequest
	//         1*time.Minute,         // DeltaStage
	//         100,                   // rMax
	//         []int{len(nodes)},     // S
	//         oracles,
	//         reportingPluginConfig, // reportingPluginConfig []byte,
	//         0,                     // maxDurationQuery
	//         250*time.Millisecond,  // maxDurationObservation
	//         0,                     // maxDurationShouldAcceptAttestedReport
	//         0,                     // maxDurationShouldTransmitAcceptedReport
	//         int(fNodes),           // f
	//         onchainConfig,
	//     )

	//     require.NoError(t, err)

	//     return
	// }

	// func setConfig(t *testing.T, steve *bind.TransactOpts, backend *backends.SimulatedBackend, verifierContract *channel_verifier.ChannelVerifier, verifierAddress common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra) ocr2types.ConfigDigest {
	//     signers, _, _, _, offchainConfigVersion, offchainConfig := generateConfig(t, nodes, oracles)

	//     signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	//     require.NoError(t, err)
	//     offchainTransmitters := make([][32]byte, nNodes)
	//     for i := 0; i < nNodes; i++ {
	//         offchainTransmitters[i] = nodes[i].ClientPubKey
	//     }
	//     _, err = verifierContract.SetConfig(steve, signerAddresses, offchainTransmitters, fNodes, offchainConfig, offchainConfigVersion, offchainConfig, nil)
	//     require.NoError(t, err)

	//     backend.Commit()

	//     l, err := verifierContract.LatestConfigDigestAndEpoch(&bind.CallOpts{})
	//     require.NoError(t, err)

	//     return l.ConfigDigest
	// }

	// func setChannelDefinitions(t *testing.T, steve *bind.TransactOpts, backend *backends.SimulatedBackend, configStoreContract *channel_config_store.ChannelConfigStore, streams []Stream) map[llotypes.ChannelID]channel_config_store.IChannelConfigStoreChannelDefinition {
	//     channels := []llotypes.ChannelID{
	//         rand.Uint32(),
	//         rand.Uint32(),
	//         rand.Uint32(),
	//         rand.Uint32(),
	//     }

	//     chainSelector, err := chainselectors.SelectorFromChainId(testutils.SimulatedChainID.Uint64())
	//     require.NoError(t, err)

	//     streamIDs := make([]uint32, len(streams))
	//     for i := 0; i < len(streams); i++ {
	//         streamIDs[i] = streams[i].id
	//     }

	//     // First set contains [1,len(streams)]
	//     channel0Def := channel_config_store.IChannelConfigStoreChannelDefinition{
	//         ReportFormat:  uint32(llotypes.ReportFormatJSON),
	//         ChainSelector: chainSelector,
	//         StreamIDs:     streamIDs[1:len(streams)],
	//     }
	//     channel1Def := channel_config_store.IChannelConfigStoreChannelDefinition{
	//         ReportFormat:  uint32(llotypes.ReportFormatEVM),
	//         ChainSelector: chainSelector,
	//         StreamIDs:     streamIDs[1:len(streams)],
	//     }

	//     // Second set contains [0,len(streams)-1]
	//     channel2Def := channel_config_store.IChannelConfigStoreChannelDefinition{
	//         ReportFormat:  uint32(llotypes.ReportFormatJSON),
	//         ChainSelector: chainSelector,
	//         StreamIDs:     streamIDs[0 : len(streams)-1],
	//     }
	//     channel3Def := channel_config_store.IChannelConfigStoreChannelDefinition{
	//         ReportFormat:  uint32(llotypes.ReportFormatEVM),
	//         ChainSelector: chainSelector,
	//         StreamIDs:     streamIDs[0 : len(streams)-1],
	//     }

	//     require.NoError(t, utils.JustError(configStoreContract.AddChannel(steve, channels[0], channel0Def)))
	//     require.NoError(t, utils.JustError(configStoreContract.AddChannel(steve, channels[1], channel1Def)))
	//     require.NoError(t, utils.JustError(configStoreContract.AddChannel(steve, channels[2], channel2Def)))
	//     require.NoError(t, utils.JustError(configStoreContract.AddChannel(steve, channels[3], channel3Def)))

	//     backend.Commit()

	//     channelDefinitions := make(map[llotypes.ChannelID]channel_config_store.IChannelConfigStoreChannelDefinition)

	//     channelDefinitions[channels[0]] = channel0Def
	//     channelDefinitions[channels[1]] = channel1Def
	//     channelDefinitions[channels[2]] = channel2Def
	//     channelDefinitions[channels[3]] = channel3Def

	//     backend.Commit()

	//     return channelDefinitions
	// }

	// func TestIntegration_LLO_Dummy(t *testing.T) {
	//     testStartTimeStamp := time.Now()

	//     streams := []Stream{btcStream, ethStream, linkStream, dogeStream}
	//     streamMap := make(map[uint32]Stream)
	//     for _, strm := range streams {
	//         streamMap[strm.id] = strm
	//     }

	//     clientCSAKeys := make([]csakey.KeyV2, nNodes)
	//     clientPubKeys := make([]ed25519.PublicKey, nNodes)
	//     for i := 0; i < nNodes; i++ {
	//         k := big.NewInt(int64(i))
	//         key := csakey.MustNewV2XXXTestingOnly(k)
	//         clientCSAKeys[i] = key
	//         clientPubKeys[i] = key.PublicKey
	//     }

	//     // Setup bootstrap
	//     bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	//     bootstrapNodePort := freeport.GetOne(t)
	//     appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_mercury", nil, bootstrapCSAKey)
	//     bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	//     t.Run("with at least one channel", func(t *testing.T) {
	//         // Setup oracle nodes
	//         var (
	//             oracles []confighelper.OracleIdentityExtra
	//             nodes   []Node
	//         )
	//         ports := freeport.GetN(t, nNodes)
	//         for i := 0; i < nNodes; i++ {
	//             app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_streams_%d", i), nil, clientCSAKeys[i])

	//             nodes = append(nodes, Node{
	//                 app, transmitter, kb, observedLogs,
	//             })
	//             offchainPublicKey, _ := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
	//             oracles = append(oracles, confighelper.OracleIdentityExtra{
	//                 OracleIdentity: confighelper.OracleIdentity{
	//                     OnchainPublicKey:  offchainPublicKey,
	//                     TransmitAccount:   ocr2types.Account(fmt.Sprintf("%x", transmitter[:])),
	//                     OffchainPublicKey: kb.OffchainPublicKey(),
	//                     PeerID:            peerID,
	//                 },
	//                 ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
	//             })
	//         }

	//         verifierAddress := common.Address{}
	//         chainID := "llo-dummy"
	//         relayType := "dummy"
	//         cd := "0x0102030405060708010203040506070801020304050607080102030405060708"
	//         signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig := generateConfig(t, nodes, oracles)
	//         var signersMarshalled, transmittersMarshalled []byte
	//         {
	//             var err error
	//             signersHex := make([]string, len(signers))
	//             for i, signer := range signers {
	//                 signersHex[i] = fmt.Sprintf("0x%x", signer)
	//             }
	//             signersMarshalled, err = json.Marshal(signersHex)
	//             require.NoError(t, err)

	//             transmittersMarshalled, err = json.Marshal(transmitters)
	//             require.NoError(t, err)
	//         }

	//         relayConfig := fmt.Sprintf(`chainID = "%s"
	// configTracker = {
	//     configDigest = "%s",
	//     configCount = 0,
	//     signers = %s,
	//     transmitters = %s,
	//     f = %d,
	//     onchainConfig = "0x%x",
	//     offchainConfigVersion = %d,
	//     offchainConfig = "0x%x",
	//     blockHeight = 10
	// }`, chainID, cd, string(signersMarshalled), string(transmittersMarshalled), f, onchainConfig, offchainConfigVersion, offchainConfig)
	//         addBootstrapJob(t, bootstrapNode, verifierAddress, "job-1", relayType, relayConfig)

	//         serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	//         serverPubKey := serverKey.PublicKey
	//         serverURL := "foo"
	//         configStoreAddress := common.Address{}

	//         // NOTE: Don't actually care about the chain ID, it just needs to be
	//         // a valid chainSelector
	//         chainSelector, err := chainselectors.SelectorFromChainId(testutils.SimulatedChainID.Uint64())
	//         require.NoError(t, err)

	//         channelDefinitions := fmt.Sprintf(`{
	// "42": {
	//     "reportFormat": %d,
	//     "chainSelector": %d,
	//     "streamIds": [51, 52]
	//     }
	// }`, llotypes.ReportFormatJSON, chainSelector)

	//         pluginConfig := fmt.Sprintf(`serverURL = "foo"
	// serverPubKey = "%x"
	// channelDefinitions = %q`, serverPubKey, channelDefinitions)
	//         addOCRJobs(t, streams, serverPubKey, serverURL, verifierAddress, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

	//         for _, node := range nodes {
	//             le := testutils.WaitForLogMessage(t, node.ObservedLogs, "Transmit")
	//             fields := le.ContextMap()
	//             assert.Equal(t, cd[2:], fields["digest"])
	//             assert.Equal(t, llotypes.ReportInfo{LifeCycleStage: "production", ReportFormat: llotypes.ReportFormatJSON}, fields["report.Info"])

	//             if fields["report.Report"] == nil {
	//                 t.Fatal("FAIL: expected log fields to contain 'report.Report'")
	//             }
	//             binaryReport := fields["report.Report"].(types.Report)
	//             report, err := (datastreamsllo.JSONReportCodec{}).Decode(binaryReport)
	//             require.NoError(t, err)
	//             assert.Equal(t, datastreamsllo.Report{
	//                 ConfigDigest:      types.ConfigDigest{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8},
	//                 ChainSelector:     0x2ee634951ef71b46,
	//                 SeqNr:             fields["seqNr"].(uint64),
	//                 ChannelID:         0x2a,
	//                 ValidAfterSeconds: report.ValidAfterSeconds, // tested separately below
	//                 ValidUntilSeconds: report.ValidUntilSeconds, // tested separately below
	//                 Values:            []*big.Int{big.NewInt(2000002000000), big.NewInt(156802000000)},
	//                 Specimen:          false,
	//             }, report)
	//             assert.GreaterOrEqual(t, report.ValidUntilSeconds, uint32(testStartTimeStamp.Unix()))
	//             assert.GreaterOrEqual(t, report.ValidAfterSeconds, uint32(testStartTimeStamp.Unix()))
	//             assert.GreaterOrEqual(t, report.ValidUntilSeconds, report.ValidAfterSeconds)

	//	        assert.GreaterOrEqual(t, int(fields["seqNr"].(uint64)), 0)
	//	    }
	//	})
}
