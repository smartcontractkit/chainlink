package llo_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/destination_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	lloevm "github.com/smartcontractkit/chainlink/v2/core/services/llo/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	mercuryverifier "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/verifier"
)

var (
	fNodes = uint8(1)
	nNodes = 4 // number of nodes (not including bootstrap)
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
	baseBenchmarkPrice decimal.Decimal
	baseBid            decimal.Decimal
	baseAsk            decimal.Decimal
}

var (
	btcStream = Stream{
		id:                 51,
		baseBenchmarkPrice: decimal.NewFromFloat32(56_114.41),
	}
	ethStream = Stream{
		id:                 52,
		baseBenchmarkPrice: decimal.NewFromFloat32(2_976.39),
	}
	linkStream = Stream{
		id:                 53,
		baseBenchmarkPrice: decimal.NewFromFloat32(13.25),
	}
	dogeStream = Stream{
		id:                 54,
		baseBenchmarkPrice: decimal.NewFromFloat32(0.10960935),
	}
	quoteStream = Stream{
		id:                 55,
		baseBenchmarkPrice: decimal.NewFromFloat32(1000.1212),
		baseBid:            decimal.NewFromFloat32(998.5431),
		baseAsk:            decimal.NewFromFloat32(1001.6999),
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

	offchainConfig = []byte{}

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

func setConfig(t *testing.T, steve *bind.TransactOpts, backend *backends.SimulatedBackend, verifierContract *channel_verifier.ChannelVerifier, verifierAddress common.Address, nodes []Node, oracles []confighelper.OracleIdentityExtra) ocr2types.ConfigDigest {
	signers, _, _, _, offchainConfigVersion, offchainConfig := generateConfig(t, oracles)

	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	offchainTransmitters := make([][32]byte, nNodes)
	for i := 0; i < nNodes; i++ {
		offchainTransmitters[i] = nodes[i].ClientPubKey
	}
	_, err = verifierContract.SetConfig(steve, signerAddresses, offchainTransmitters, fNodes, offchainConfig, offchainConfigVersion, offchainConfig, nil)
	require.NoError(t, err)

	backend.Commit()

	l, err := verifierContract.LatestConfigDigestAndEpoch(&bind.CallOpts{})
	require.NoError(t, err)

	return l.ConfigDigest
}

// On-chain format is not finalized yet so use the dummy relayer for testing
func TestIntegration_LLO_Dummy(t *testing.T) {
	testStartTimeStamp := time.Now()

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)
	for i := 0; i < nNodes; i++ {
		k := big.NewInt(int64(i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo", nil, bootstrapCSAKey)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	t.Run("produces reports in v0.3 format", func(t *testing.T) {
		streams := []Stream{ethStream, linkStream, quoteStream}
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
			app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_streams_%d", i), nil, clientCSAKeys[i])

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

		verifierAddress := common.Address{}
		chainID := "llo-dummy"
		relayType := "dummy"
		cd := ocr2types.ConfigDigest{0x01, 0x02, 0x03, 0x04}
		signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig := generateConfig(t, oracles)
		var signersMarshalled, transmittersMarshalled []byte
		{
			var err error
			signersHex := make([]string, len(signers))
			for i, signer := range signers {
				signersHex[i] = fmt.Sprintf("0x%x", signer)
			}
			signersMarshalled, err = json.Marshal(signersHex)
			require.NoError(t, err)

			transmittersMarshalled, err = json.Marshal(transmitters)
			require.NoError(t, err)
		}

		relayConfig := fmt.Sprintf(`chainID = "%s"
configTracker = {
	configDigest = "0x%x",
	configCount = 0,
	signers = %s,
	transmitters = %s,
	f = %d,
	onchainConfig = "0x%x",
	offchainConfigVersion = %d,
	offchainConfig = "0x%x",
	blockHeight = 10
}`, chainID, cd[:], string(signersMarshalled), string(transmittersMarshalled), f, onchainConfig, offchainConfigVersion, offchainConfig)
		addBootstrapJob(t, bootstrapNode, verifierAddress, "job-2", relayType, relayConfig)

		serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(-1))
		serverPubKey := serverKey.PublicKey
		serverURL := "foo"
		configStoreAddress := common.Address{}

		chainSelector := 4949039107694359620 // arbitrum mainnet

		feedID := [32]byte{00, 03, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
		expirationWindow := 3600
		multiplier := big.NewInt(1e18)
		baseUSDFee := 10
		// 52 = eth, 53 = link, 55 = quote
		channelDefinitions := fmt.Sprintf(`
{
	"42": {
		"reportFormat": %d,
		"chainSelector": %d,
		"streams": [{"streamId": 52, "aggregator": %d}, {"streamId": 53, "aggregator": %d}, {"streamId": 55, "aggregator": %d}],
		"opts": {
			"feedId": "0x%x",
			"expirationWindow": %d,
			"multiplier": "%s",
			"baseUSDFee": "%d"
		}
	}
}`, llotypes.ReportFormatEVMPremiumLegacy, chainSelector, llotypes.AggregatorMedian, llotypes.AggregatorMedian, llotypes.AggregatorQuote, feedID, expirationWindow, multiplier.String(), baseUSDFee)

		pluginConfig := fmt.Sprintf(`serverURL = "foo"
donID = 42
serverPubKey = "%x"
channelDefinitions = %q`, serverPubKey, channelDefinitions)
		jobIDs := addOCRJobsEVMPremiumLegacy(t, streams, serverPubKey, serverURL, verifierAddress, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

		steve, backend, verifier, verifierProxy, _ := setupV03Blockchain(t)

		// Set config
		recipientAddressesAndWeights := []destination_verifier.CommonAddressAndWeight{}
		signerAddresses := make([]common.Address, len(oracles))
		for i, oracle := range oracles {
			signerAddresses[i] = common.BytesToAddress(oracle.OracleIdentity.OnchainPublicKey)
		}

		_, err := verifier.SetConfig(steve, signerAddresses, f, recipientAddressesAndWeights)
		require.NoError(t, err)
		backend.Commit()

		for i, node := range nodes {
			le := testutils.WaitForLogMessage(t, node.ObservedLogs, "Transmit")
			fields := le.ContextMap()
			assert.Equal(t, hexutil.Encode(cd[:]), "0x"+fields["digest"].(string))
			assert.Equal(t, llotypes.ReportInfo{LifeCycleStage: "production", ReportFormat: llotypes.ReportFormatEVMPremiumLegacy}, fields["report.Info"])

			if fields["report.Report"] == nil {
				t.Fatal("FAIL: expected log fields to contain 'report.Report'")
			}
			binaryReport := fields["report.Report"].(types.Report)
			report, err := (lloevm.ReportCodecPremiumLegacy{}).Decode(binaryReport)
			require.NoError(t, err)
			assert.Equal(t, feedID, report.FeedId)
			assert.GreaterOrEqual(t, report.ObservationsTimestamp, uint32(testStartTimeStamp.Unix()))
			assert.Equal(t, quoteStream.baseBenchmarkPrice.Mul(decimal.NewFromBigInt(multiplier, 0)).String(), report.BenchmarkPrice.String())
			assert.Equal(t, quoteStream.baseBid.Mul(decimal.NewFromBigInt(multiplier, 0)).String(), report.Bid.String())
			assert.Equal(t, quoteStream.baseAsk.Mul(decimal.NewFromBigInt(multiplier, 0)).String(), report.Ask.String())
			assert.GreaterOrEqual(t, report.ValidFromTimestamp, uint32(testStartTimeStamp.Unix()))
			assert.Equal(t, report.ObservationsTimestamp+uint32(expirationWindow), report.ExpiresAt)
			assert.Equal(t, big.NewInt(754716981132075472), report.LinkFee)
			assert.Equal(t, big.NewInt(3359774760700043), report.NativeFee)

			seqNr := fields["seqNr"].(uint64)
			assert.Greater(t, int(seqNr), 0)

			sigs := fields["sigs"].([]types.AttributedOnchainSignature)

			t.Run(fmt.Sprintf("emulate mercury server verifying report (local verification) - node %d", i), func(t *testing.T) {
				var rs [][32]byte
				var ss [][32]byte
				var vs [32]byte
				for i, as := range sigs {
					r, s, v, err := evmutil.SplitSignature(as.Signature)
					if err != nil {
						panic("error in SplitSignature")
					}
					rs = append(rs, r)
					ss = append(ss, s)
					vs[i] = v
				}
				rc := lloevm.LegacyReportContext(cd, seqNr)
				rawReportCtx := evmutil.RawReportContext(rc)
				rv := mercuryverifier.NewVerifier()

				reportSigners, err := rv.Verify(mercuryverifier.SignedReport{
					RawRs:         rs,
					RawSs:         ss,
					RawVs:         vs,
					ReportContext: rawReportCtx,
					Report:        binaryReport,
				}, f, signerAddresses)
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(reportSigners), int(f+1))
				assert.Subset(t, signerAddresses, reportSigners)
			})

			t.Run(fmt.Sprintf("test on-chain verification - node %d", i), func(t *testing.T) {
				signedReport, err := lloevm.ReportCodecPremiumLegacy{}.Pack(cd, seqNr, binaryReport, sigs)
				require.NoError(t, err)

				_, err = verifierProxy.Verify(steve, signedReport, []byte{})
				require.NoError(t, err)

			})
		}

		t.Run("if link/eth stream specs start failing, uses 0 for the fee", func(t *testing.T) {
			t.Run("link/eth stream specs are missing", func(t *testing.T) {
				// delete eth/link stream specs
				for idx, strmIDs := range jobIDs {
					for strmID, jobID := range strmIDs {
						if strmID == ethStream.id || strmID == linkStream.id {
							nodes[idx].DeleteJob(t, jobID)
						}
					}
				}

				for _, node := range nodes {
					node.ObservedLogs.TakeAll()

					le := testutils.WaitForLogMessage(t, node.ObservedLogs, "Observation failed for streams")
					fields := le.ContextMap()
					assert.Equal(t, []interface{}{ethStream.id, linkStream.id}, fields["failedStreamIDs"])
					assert.Equal(t, []interface{}{"StreamID: 52; Reason: missing stream: 52", "StreamID: 53; Reason: missing stream: 53"}, fields["errors"])

					le = testutils.WaitForLogMessage(t, node.ObservedLogs, "Transmit")
					fields = le.ContextMap()
					assert.Equal(t, hexutil.Encode(cd[:]), "0x"+fields["digest"].(string))
					assert.Equal(t, llotypes.ReportInfo{LifeCycleStage: "production", ReportFormat: llotypes.ReportFormatEVMPremiumLegacy}, fields["report.Info"])

					if fields["report.Report"] == nil {
						t.Fatal("FAIL: expected log fields to contain 'report.Report'")
					}
					binaryReport := fields["report.Report"].(types.Report)
					report, err := (lloevm.ReportCodecPremiumLegacy{}).Decode(binaryReport)
					require.NoError(t, err)
					assert.Equal(t, feedID, report.FeedId)
					assert.GreaterOrEqual(t, report.ObservationsTimestamp, uint32(testStartTimeStamp.Unix()))
					assert.Equal(t, quoteStream.baseBenchmarkPrice.Mul(decimal.NewFromBigInt(multiplier, 0)).String(), report.BenchmarkPrice.String())
					assert.Equal(t, quoteStream.baseBid.Mul(decimal.NewFromBigInt(multiplier, 0)).String(), report.Bid.String())
					assert.Equal(t, quoteStream.baseAsk.Mul(decimal.NewFromBigInt(multiplier, 0)).String(), report.Ask.String())
					assert.GreaterOrEqual(t, report.ValidFromTimestamp, uint32(testStartTimeStamp.Unix()))
					assert.Equal(t, report.ObservationsTimestamp+uint32(expirationWindow), report.ExpiresAt)
					assert.Equal(t, "0", report.LinkFee.String())
					assert.Equal(t, "0", report.NativeFee.String())
				}
			})

			t.Run("link/eth stream specs have EAs that return error", func(t *testing.T) {
				// add new stream specs that will fail
				for i, node := range nodes {
					for j, strm := range streams {
						if strm.id == ethStream.id || strm.id == linkStream.id {
							var name string
							if j == 0 {
								name = "nativeprice"
							} else {
								name = "linkprice"
							}
							name = fmt.Sprintf("%s-%d-%d-erroring", name, strm.id, j)
							bmBridge := createErroringBridge(t, name, i, node.App.BridgeORM())
							addSingleDecimalStreamJob(
								t,
								node,
								strm.id,
								bmBridge,
							)
						}
					}
				}

				for _, node := range nodes {
					node.ObservedLogs.TakeAll()

					le := testutils.WaitForLogMessage(t, node.ObservedLogs, "Observation failed for streams")
					fields := le.ContextMap()
					assert.Equal(t, []interface{}{ethStream.id, linkStream.id}, fields["failedStreamIDs"])
					assert.Len(t, fields["errors"], 2)
					for _, err := range fields["errors"].([]interface{}) {
						assert.Contains(t, err.(string), "Reason: failed to extract big.Int")
						assert.Contains(t, err.(string), "status code 500")
					}

					le = testutils.WaitForLogMessage(t, node.ObservedLogs, "Transmit")
					fields = le.ContextMap()
					assert.Equal(t, hexutil.Encode(cd[:]), "0x"+fields["digest"].(string))
					assert.Equal(t, llotypes.ReportInfo{LifeCycleStage: "production", ReportFormat: llotypes.ReportFormatEVMPremiumLegacy}, fields["report.Info"])

					if fields["report.Report"] == nil {
						t.Fatal("FAIL: expected log fields to contain 'report.Report'")
					}
					binaryReport := fields["report.Report"].(types.Report)
					report, err := (lloevm.ReportCodecPremiumLegacy{}).Decode(binaryReport)
					require.NoError(t, err)
					assert.Equal(t, feedID, report.FeedId)
					assert.GreaterOrEqual(t, report.ObservationsTimestamp, uint32(testStartTimeStamp.Unix()))
					assert.Equal(t, quoteStream.baseBenchmarkPrice.Mul(decimal.NewFromBigInt(multiplier, 0)).String(), report.BenchmarkPrice.String())
					assert.Equal(t, quoteStream.baseBid.Mul(decimal.NewFromBigInt(multiplier, 0)).String(), report.Bid.String())
					assert.Equal(t, quoteStream.baseAsk.Mul(decimal.NewFromBigInt(multiplier, 0)).String(), report.Ask.String())
					assert.GreaterOrEqual(t, report.ValidFromTimestamp, uint32(testStartTimeStamp.Unix()))
					assert.Equal(t, int(report.ObservationsTimestamp+uint32(expirationWindow)), int(report.ExpiresAt))
					assert.Equal(t, "0", report.LinkFee.String())
					assert.Equal(t, "0", report.NativeFee.String())
				}
			})
		})

		t.Run("deleting LLO jobs cleans up resources", func(t *testing.T) {
			t.Skip("TODO - https://smartcontract-it.atlassian.net/browse/MERC-3653")
		})
	})
}
