package llo_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/peer"

	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/csakey"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	mercurytransmitter "github.com/smartcontractkit/chainlink-data-streams/llo/transmitter/de"
	"github.com/smartcontractkit/chainlink-data-streams/mercury"
	reportcodecv3 "github.com/smartcontractkit/chainlink-data-streams/mercury/v3/reportcodec"
	mercuryverifier "github.com/smartcontractkit/chainlink-data-streams/mercury/verifier"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/destination_verifier"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// historyBackfillOptsJSON builds ChannelOpts JSON for ReportFormatHistoryBackfill.
// Map iteration order is non-deterministic; OCR requires identical proposal bytes
// across oracles, so keys are emitted in sorted order.
func historyBackfillOptsJSON(targetChannelID llotypes.ChannelID, observations map[uint64]map[llotypes.StreamID]string) llotypes.ChannelOpts {
	tsKeys := make([]uint64, 0, len(observations))
	for ts := range observations {
		tsKeys = append(tsKeys, ts)
	}
	slices.Sort(tsKeys)

	var b strings.Builder
	// Field order must match json.Marshal / channel cache (lexicographic): observations before targetChannelId.
	b.WriteString(`{"observations":{`)
	for i, ts := range tsKeys {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(`"`)
		b.WriteString(strconv.FormatUint(ts, 10))
		b.WriteString(`":{`)
		row := observations[ts]
		sids := make([]llotypes.StreamID, 0, len(row))
		for sid := range row {
			sids = append(sids, sid)
		}
		slices.Sort(sids)
		for j, sid := range sids {
			if j > 0 {
				b.WriteByte(',')
			}
			b.WriteString(`"`)
			b.WriteString(strconv.FormatUint(uint64(sid), 10))
			b.WriteString(`":`)
			val, err := json.Marshal(row[sid])
			if err != nil {
				panic(err)
			}
			b.Write(val)
		}
		b.WriteString(`}`)
	}
	b.WriteString(`},"targetChannelId":`)
	b.WriteString(strconv.FormatUint(uint64(targetChannelID), 10))
	b.WriteString(`}`)
	return []byte(b.String())
}

func decodeEVMPremiumLegacyReport(t *testing.T, pckt *packet) (feedID [32]byte, obsTs uint32, benchmark *big.Int) {
	t.Helper()
	require.Equal(t, uint32(llotypes.ReportFormatEVMPremiumLegacy), pckt.req.ReportFormat)
	v := make(map[string]any)
	err := mercury.PayloadTypes.UnpackIntoMap(v, pckt.req.Payload)
	require.NoError(t, err)
	report, exists := v["report"]
	require.True(t, exists, "payload should contain report")
	reportElems := make(map[string]any)
	err = reportcodecv3.ReportTypes.UnpackIntoMap(reportElems, report.([]byte))
	require.NoError(t, err)
	feedID = reportElems["feedId"].([32]uint8)
	obsTs = reportElems["observationsTimestamp"].(uint32)
	benchmark = reportElems["benchmarkPrice"].(*big.Int)
	return feedID, obsTs, benchmark
}

func quoteBackfillString(benchmark float64) string {
	bm := decimal.NewFromFloat(benchmark)
	bid := bm.Sub(decimal.NewFromFloat(1))
	ask := bm.Add(decimal.NewFromFloat(1))
	return fmt.Sprintf("Q{Bid: %s, Benchmark: %s, Ask: %s}", bid.String(), bm.String(), ask.String())
}

func TestIntegration_LLO_history_backfill(t *testing.T) {
	t.Parallel()

	const (
		salt              = 600
		donID             = uint32(776655)
		targetChannelID   = llotypes.ChannelID(1)
		backfillChannelID = llotypes.ChannelID(100)
		streamNative      = uint32(290)
		streamLink        = uint32(291)
		streamQuote       = uint32(292)
	)

	targetFeedID := common.HexToHash("0x0004444444444444444444444444444444444444444444444444444444444444")
	multiplier := decimal.New(1, 18)
	expirationWindow := uint32(3600)

	offchainConfig := datastreamsllo.OffchainConfig{
		ProtocolVersion:                     1,
		DefaultMinReportIntervalNanoseconds: uint64(1 * time.Second),
		EnableObservationCompression:        true,
	}

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)
	for i := range nNodes {
		k := big.NewInt(int64(salt + i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	steve, backend, configurator, configuratorAddress, destinationVerifier, _, _, _, configStore, configStoreAddress, _, _, _, _ := setupBlockchain(t)
	fromBlock := 1

	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_llo_backfill", backend, bootstrapCSAKey, nil)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	packetCh := make(chan *packet, 100000)
	serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
	serverPubKey := serverKey.PublicKey
	srv := NewMercuryServer(t, serverKey, packetCh)
	serverURL := startMercuryServer(t, srv, clientPubKeys)

	oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
		c.Mercury.Transmitter.Protocol = new(mercurytransmitter.MercuryTransmitterProtocolGRPC)
	})

	chainID := testutils.SimulatedChainID
	relayType := "evm"
	relayConfig := fmt.Sprintf(`
chainID = "%s"
fromBlock = %d
lloDonID = %d
lloConfigMode = "bluegreen"
`, chainID, fromBlock, donID)
	addBootstrapJob(t, bootstrapNode, configuratorAddress, "job-history-backfill", relayType, relayConfig)

	pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
donID = %d
channelDefinitionsContractAddress = "0x%x"
channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock)

	nativeStrm := Stream{
		id:                 streamNative,
		baseBenchmarkPrice: decimal.NewFromFloat(2976.39),
	}
	linkStrm := Stream{
		id:                 streamLink,
		baseBenchmarkPrice: decimal.NewFromFloat(13.25),
	}
	quoteStrm := Stream{
		id:                 streamQuote,
		baseBenchmarkPrice: decimal.NewFromFloat(1000.1212),
		baseBid:            decimal.NewFromFloat(998.5431),
		baseAsk:            decimal.NewFromFloat(1001.6999),
	}
	streams := []Stream{nativeStrm, linkStrm, quoteStrm}

	addOCRJobsEVMPremiumLegacy(t, streams, serverPubKey, serverURL, configuratorAddress, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

	targetStreams := []llotypes.Stream{
		{StreamID: streamNative, Aggregator: llotypes.AggregatorMedian},
		{StreamID: streamLink, Aggregator: llotypes.AggregatorMedian},
		{StreamID: streamQuote, Aggregator: llotypes.AggregatorQuote},
	}
	targetOpts := llotypes.ChannelOpts(fmt.Appendf(nil,
		`{"baseUSDFee":"0.1","expirationWindow":%d,"feedId":"0x%x","multiplier":"%s"}`,
		expirationWindow, targetFeedID, multiplier.String(),
	))

	phase1Defs := llotypes.ChannelDefinitions{
		targetChannelID: {
			ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
			Streams:      targetStreams,
			Opts:         targetOpts,
		},
	}
	url, sha := newChannelDefinitionsServer(t, phase1Defs)
	_, err := configStore.SetChannelDefinitions(steve, donID, url, sha)
	require.NoError(t, err)
	backend.Commit()

	setProductionConfig(
		t, donID, steve, backend, configurator, configuratorAddress, nodes,
		WithOracles(oracles), WithOffchainConfig(offchainConfig),
	)

	signerAddresses := make([]common.Address, len(oracles))
	for i, oracle := range oracles {
		signerAddresses[i] = common.BytesToAddress(oracle.OnchainPublicKey)
	}
	_, err = destinationVerifier.SetConfig(steve, signerAddresses, fNodes, []destination_verifier.CommonAddressAndWeight{})
	require.NoError(t, err)
	backend.Commit()

	testStart := time.Now()
	expectedLiveBenchmark := quoteStrm.baseBenchmarkPrice.Mul(multiplier).BigInt()

	var sawLive bool
	require.Eventually(t, func() bool {
		pckt, errReceive := receiveWithTimeout(t, packetCh, 2*time.Second)
		if errReceive != nil {
			return sawLive
		}
		if pckt.req.ReportFormat != uint32(llotypes.ReportFormatEVMPremiumLegacy) {
			return sawLive
		}
		feedID, obsTs, bm := decodeEVMPremiumLegacyReport(t, pckt)
		if feedID != targetFeedID {
			return sawLive
		}
		if int64(obsTs) < testStart.Unix()-5 {
			return sawLive
		}
		if bm.String() != expectedLiveBenchmark.String() {
			return sawLive
		}
		{
			v := make(map[string]any)
			require.NoError(t, mercury.PayloadTypes.UnpackIntoMap(v, pckt.req.Payload))
			rv := mercuryverifier.NewVerifier()
			_, errVerify := rv.Verify(mercuryverifier.SignedReport{
				RawRs:         v["rawRs"].([][32]byte),
				RawSs:         v["rawSs"].([][32]byte),
				RawVs:         v["rawVs"].([32]byte),
				ReportContext: v["reportContext"].([3][32]byte),
				Report:        v["report"].([]byte),
			}, fNodes, signerAddresses)
			if errVerify != nil {
				return sawLive
			}
		}
		pr, ok := peer.FromContext(pckt.ctx)
		require.True(t, ok)
		t.Logf("live report from %s feed=%s obsTs=%d", pr.String(), hex.EncodeToString(feedID[:]), obsTs)
		sawLive = true
		return true
	}, reportTimeout, 200*time.Millisecond, "expected live EVMPremiumLegacy report for target feedId")

	// Phase 2: add history_backfill channel (same streams as target; observations in the past).
	histBenchmarks := map[uint64]float64{100: 100.0, 150: 150.0, 200: 200.0}
	observations := make(map[uint64]map[llotypes.StreamID]string, len(histBenchmarks))
	for ts, bench := range histBenchmarks {
		observations[ts] = map[llotypes.StreamID]string{
			streamNative: "2976.39",
			streamLink:   "13.25",
			streamQuote:  quoteBackfillString(bench),
		}
	}
	phase2Defs := llotypes.ChannelDefinitions{
		targetChannelID: {
			ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
			Streams:      targetStreams,
			Opts:         targetOpts,
		},
		backfillChannelID: {
			ReportFormat: llotypes.ReportFormatHistoryBackfill,
			Streams:      targetStreams,
			Opts:         historyBackfillOptsJSON(targetChannelID, observations),
		},
	}
	url2, sha2 := newChannelDefinitionsServer(t, phase2Defs)
	_, err = configStore.SetChannelDefinitions(steve, donID, url2, sha2)
	require.NoError(t, err)
	backend.Commit()

	histByObsTs := make(map[uint32]*big.Int)
	require.Eventually(t, func() bool {
		for len(histByObsTs) < 3 {
			pckt, err := receiveWithTimeout(t, packetCh, 2*time.Second)
			if err != nil {
				break
			}
			if pckt.req.ReportFormat != uint32(llotypes.ReportFormatEVMPremiumLegacy) {
				continue
			}
			feedID, obsTs, bm := decodeEVMPremiumLegacyReport(t, pckt)
			if feedID != targetFeedID {
				continue
			}
			switch obsTs {
			case 100, 150, 200:
				histByObsTs[obsTs] = new(big.Int).Set(bm)
			default:
				// live reports continue; ignore
			}
		}
		return len(histByObsTs) == 3
	}, reportTimeout, 200*time.Millisecond, "expected three historical backfill reports at obsTs 100, 150, 200")

	histTs := make([]uint32, 0, len(histByObsTs))
	for ts := range histByObsTs {
		histTs = append(histTs, ts)
	}
	slices.Sort(histTs)
	require.Equal(t, []uint32{100, 150, 200}, histTs)
	for _, ts := range histTs {
		bench := histBenchmarks[uint64(ts)]
		want := decimal.NewFromFloat(bench).Mul(multiplier).BigInt()
		require.Equal(t, want.String(), histByObsTs[ts].String(), "benchmark for obsTs %d", ts)
	}

	// Several oracles may transmit the same attested round; flush duplicate historical
	// packets so phase 3 is not tripped by the same obsTs delivered multiple times.
	flushUntil := time.Now().Add(4 * time.Second)
	for time.Now().Before(flushUntil) {
		_, err := receiveWithTimeout(t, packetCh, 150*time.Millisecond)
		if err != nil {
			continue
		}
	}

	// Phase 3: no further reports with historical observation timestamps.
	checkWindow := 8 * time.Second
	deadline := time.Now().Add(checkWindow)
	for time.Now().Before(deadline) {
		pckt, err := receiveWithTimeout(t, packetCh, 1*time.Second)
		if err != nil {
			continue
		}
		if pckt.req.ReportFormat != uint32(llotypes.ReportFormatEVMPremiumLegacy) {
			continue
		}
		feedID, obsTs, _ := decodeEVMPremiumLegacyReport(t, pckt)
		if feedID != targetFeedID {
			continue
		}
		switch obsTs {
		case 100, 150, 200:
			require.Fail(t, "unexpected additional historical backfill report", "obsTs=%d", obsTs)
		default:
			require.GreaterOrEqual(t, int64(obsTs), testStart.Unix()-30,
				"only live-scale observation timestamps expected after backfill completes")
		}
	}
}
