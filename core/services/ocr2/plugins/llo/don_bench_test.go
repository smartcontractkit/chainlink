package llo_test

import (
	"crypto/ed25519"
	"fmt"
	"math/big"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/peer"

	"github.com/smartcontractkit/freeport"
	memkvdb "github.com/smartcontractkit/libocr/offchainreporting2plus/ocrintegrationtesthelpers"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/csakey"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	lloconfig "github.com/smartcontractkit/chainlink-data-streams/llo/pluginconfig"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
	lloreportcodec "github.com/smartcontractkit/chainlink-data-streams/llo/reportcodec"
	mercurytransmitter "github.com/smartcontractkit/chainlink-data-streams/llo/transmitter/dataengine"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// DON benchmark parameters. Each channel reads its own memo stream, so every
// round observes donBenchStreams values and emits donBenchChannels JSON
// reports per oracle.
const (
	donBenchChannels = 2000
	donBenchStreams  = 2000
	// donBenchRounds is how many full rounds, every channel delivered, each
	// oracle must reach before the run stops. Channels are voted in five per
	// round, so the first full round comes after donBenchChannels/5 rounds.
	donBenchRounds = 20
	// donBenchDeltaRound is the round cadence floor libocr enforces. The lowest
	// the OCR3.1 production sanity check accepts with the default snapshot
	// retention, and the same for v30 so the two runs compare.
	donBenchDeltaRound = 40 * time.Millisecond
	donBenchEnv        = "LLO_DON_BENCH"
)

// TestIntegration_LLO_don_benchmark runs a 4 oracle DON with real libocr over
// in-memory networking once per plugin version, on the same channel set, and
// reports what the in-process bench package cannot: round cadence under the
// real message pattern, observation to transmission latency and report
// throughput as delivered to the transmission server.
//
// It is opt-in through LLO_DON_BENCH=1, since each version starts a node set
// and runs for a minute or more. Versions run one after the other so they do
// not contend for CPU.
func TestIntegration_LLO_don_benchmark(t *testing.T) {
	if os.Getenv(donBenchEnv) == "" {
		t.Skipf("set %s=1 to run the DON benchmark", donBenchEnv)
	}

	versions := []struct {
		version lloconfig.PluginVersion
		salt    int64
		donID   uint32
	}{
		{version: lloconfig.PluginVersionV30, salt: 900, donID: 777300},
		{version: lloconfig.PluginVersionV31, salt: 950, donID: 777310},
	}
	results := make([]donBenchResult, 0, len(versions))
	for _, v := range versions {
		t.Run(string(v.version), func(t *testing.T) {
			results = append(results, runDONBench(t, v.version, v.salt, v.donID))
		})
	}

	t.Logf("DON benchmark: %d oracles, %d channels, %d streams, %d rounds, DeltaRound %s", nNodes, donBenchChannels, donBenchStreams, donBenchRounds, donBenchDeltaRound)
	for _, r := range results {
		t.Log(r.String())
	}
}

// donBenchResult is one version's run, aggregated over every oracle.
type donBenchResult struct {
	version lloconfig.PluginVersion
	// timeToFirstReport is from config set to the first delivered report,
	// timeToFullRound to the first round that delivered every channel.
	timeToFirstReport, timeToFullRound time.Duration
	// roundInterval is the time between consecutive reporting rounds, from the
	// agreed observation timestamps, which are the same on every oracle.
	roundIntervalMedian, roundIntervalP90 time.Duration
	// transmitLatency is from the agreed observation timestamp to the report
	// reaching the transmission server, per report, over every oracle.
	transmitLatencyMedian, transmitLatencyP90 time.Duration
	// reportsPerSecond is delivered reports per oracle over the measured span.
	reportsPerSecond float64
	rounds           int
	reportsPerOracle int
}

func (r donBenchResult) String() string {
	return fmt.Sprintf("%s: full_rounds=%d reports/oracle=%d first_report=%s full_round=%s round_interval median=%s p90=%s transmit_latency median=%s p90=%s reports/s/oracle=%.0f",
		r.version, r.rounds, r.reportsPerOracle, r.timeToFirstReport.Round(time.Millisecond), r.timeToFullRound.Round(time.Millisecond),
		r.roundIntervalMedian.Round(time.Microsecond), r.roundIntervalP90.Round(time.Microsecond),
		r.transmitLatencyMedian.Round(time.Microsecond), r.transmitLatencyP90.Round(time.Microsecond),
		r.reportsPerSecond)
}

// donBenchReport is one delivered report with the wall-clock time the
// transmission server received it.
type donBenchReport struct {
	report     lloprotocol.Report
	receivedAt time.Time
}

func runDONBench(t *testing.T, version lloconfig.PluginVersion, salt int64, donID uint32) donBenchResult {
	const logLevel = toml.LogLevel(zapcore.ErrorLevel)

	offchainConfig := lloprotocol.OffchainConfig{
		ProtocolVersion:                     1,
		DefaultMinReportIntervalNanoseconds: 1,
		AggregationFaultTolerance:           aggregationFaultTolerance,
		EnableObservationCompression:        true,
	}
	ocrConfigOpts := []OCRConfigOption{
		WithOffchainConfig(offchainConfig),
		func(cfg *OCRConfig) {
			cfg.DeltaRound = donBenchDeltaRound
			cfg.DeltaGrace = 5 * time.Millisecond
			cfg.DeltaCertifiedCommitRequest = 10 * time.Millisecond
		},
	}
	if version == lloconfig.PluginVersionV31 {
		ocrConfigOpts = append(ocrConfigOpts, WithV31())
	}

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)
	for i := range nNodes {
		key := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt + int64(i)))
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	steve, backend, configurator, configuratorAddress, _, _, _, _, configStore, configStoreAddress, _, _, _, _ := setupBlockchain(t)
	fromBlock := 1

	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, fmt.Sprintf("bootstrap_llo_bench_%s", version), backend, bootstrapCSAKey, func(c *chainlink.Config) {
		c.Log.Level = new(logLevel)
	})
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	// The server drops packets once this buffer is full, so size it for the
	// whole run plus the rounds in flight when it ends.
	packets := make(chan *packet, 2*donBenchChannels*donBenchRounds*nNodes)
	serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
	serverPubKey := serverKey.PublicKey
	srv := NewMercuryServer(t, serverKey, packets)
	serverURL := startMercuryServer(t, srv, clientPubKeys)

	streams := make([]Stream, donBenchStreams)
	for i := range streams {
		streams[i] = Stream{id: uint32(1000 + i), baseBenchmarkPrice: decimal.NewFromInt(int64(1000 + i))} //nolint:gosec // G115: small test counter
	}

	// v31 state lives in memory so the round cost excludes disk writes, which
	// the in-process bench package does not pay either.
	oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
		c.Mercury.Transmitter.Protocol = new(mercurytransmitter.MercuryTransmitterProtocolGRPC)
		c.Log.Level = new(logLevel)
	}, memkvdb.NewStatelessInMemoryKeyValueDatabaseFactory())

	relayType := "evm"
	relayConfig := fmt.Sprintf(`
chainID = "%s"
fromBlock = %d
lloDonID = %d
lloConfigMode = "bluegreen"
`, testutils.SimulatedChainID, fromBlock, donID)
	addBootstrapJob(t, bootstrapNode, configuratorAddress, fmt.Sprintf("job-bench-%s", version), relayType, relayConfig)

	channelDefinitions := make(llotypes.ChannelDefinitions, donBenchChannels)
	for i := range uint32(donBenchChannels) {
		channelDefinitions[i] = llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatJSON,
			Streams: []llotypes.Stream{{
				StreamID:   streams[int(i)%len(streams)].id,
				Aggregator: llotypes.AggregatorMedian,
			}},
		}
	}
	url, sha := newChannelDefinitionsServer(t, channelDefinitions)
	_, err := configStore.SetChannelDefinitions(steve, donID, url, sha)
	require.NoError(t, err)
	backend.Commit()

	pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
donID = %d
channelDefinitionsContractAddress = "0x%x"
channelDefinitionsContractFromBlock = %d
pluginVersion = "%s"`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock, version)
	for i, node := range nodes {
		addLLOJob(t, node, configuratorAddress, bootstrapPeerID, bootstrapNodePort, clientPubKeys[i], fmt.Sprintf("bench-%s", version), pluginConfig, relayType, relayConfig)
		addMemoStreamSpecs(t, node, streams)
	}

	opts := append([]OCRConfigOption{WithOracles(oracles)}, ocrConfigOpts...)
	configSetAt := time.Now()
	setProductionConfig(t, donID, steve, backend, configurator, configuratorAddress, nodes, opts...)

	// Collect until every oracle has delivered donBenchRounds full rounds.
	// Transmissions arrive out of order, so rounds are counted by seqNr.
	byOracle := map[string][]donBenchReport{}
	// oracle => seqNr => reports delivered for that round
	roundsByOracle := map[string]map[uint64]int{}
	var firstReportAt, fullRoundAt time.Time
	for {
		pckt, err := receiveWithTimeout(t, packets, reportTimeout)
		require.NoError(t, err)
		receivedAt := time.Now()
		if firstReportAt.IsZero() {
			firstReportAt = receivedAt
		}
		pr, ok := peer.FromContext(pckt.ctx)
		require.True(t, ok)
		addr := pr.Addr.String()
		_, _, r, _, err := (lloreportcodec.JSONReportCodec{}).UnpackDecode(pckt.req.Payload)
		require.NoError(t, err)

		byOracle[addr] = append(byOracle[addr], donBenchReport{report: r, receivedAt: receivedAt})
		if roundsByOracle[addr] == nil {
			roundsByOracle[addr] = map[uint64]int{}
		}
		roundsByOracle[addr][r.SeqNr]++
		if fullRoundAt.IsZero() && roundsByOracle[addr][r.SeqNr] >= donBenchChannels {
			fullRoundAt = receivedAt
		}

		finished := 0
		for _, rounds := range roundsByOracle {
			full := 0
			for _, n := range rounds {
				if n >= donBenchChannels {
					full++
				}
			}
			if full >= donBenchRounds {
				finished++
			}
		}
		if finished >= nNodes {
			break
		}
	}
	for _, node := range nodes {
		require.NoError(t, node.App.Stop())
	}
	require.NoError(t, bootstrapNode.App.Stop())

	res := summarizeDONBench(t, version, byOracle, roundsByOracle)
	res.timeToFirstReport = firstReportAt.Sub(configSetAt)
	res.timeToFullRound = fullRoundAt.Sub(configSetAt)
	return res
}

// summarizeDONBench reduces the delivered reports to the result figures, over
// full rounds only: those every oracle delivered every channel for, so the
// channel voting rounds and the rounds cut off by shutdown stay out. Round
// intervals come from the agreed timestamps, identical across the DON; latency
// and throughput pool every oracle.
func summarizeDONBench(t *testing.T, version lloconfig.PluginVersion, byOracle map[string][]donBenchReport, roundsByOracle map[string]map[uint64]int) donBenchResult {
	t.Helper()
	res := donBenchResult{version: version}

	fullRound := func(seqNr uint64) bool {
		for _, rounds := range roundsByOracle {
			if rounds[seqNr] < donBenchChannels {
				return false
			}
		}
		return true
	}

	// Agreed observation timestamp per seqNr, same on every oracle.
	obsTsBySeqNr := map[uint64]uint64{}
	var latencies []time.Duration
	var spanStart, spanEnd time.Time
	total := 0
	for _, reports := range byOracle {
		for _, dr := range reports {
			if !fullRound(dr.report.SeqNr) {
				continue
			}
			total++
			obsTsBySeqNr[dr.report.SeqNr] = dr.report.ObservationTimestampNanoseconds
			obsAt := time.Unix(0, int64(dr.report.ObservationTimestampNanoseconds)) //nolint:gosec // G115: wall-clock nanoseconds fit int64
			latencies = append(latencies, dr.receivedAt.Sub(obsAt))
			if spanStart.IsZero() || dr.receivedAt.Before(spanStart) {
				spanStart = dr.receivedAt
			}
			if dr.receivedAt.After(spanEnd) {
				spanEnd = dr.receivedAt
			}
		}
	}
	require.NotEmpty(t, obsTsBySeqNr, "no full round delivered")

	seqNrs := make([]uint64, 0, len(obsTsBySeqNr))
	for seqNr := range obsTsBySeqNr {
		seqNrs = append(seqNrs, seqNr)
	}
	slices.Sort(seqNrs)
	var intervals []time.Duration
	for i := 1; i < len(seqNrs); i++ {
		// Only consecutive rounds measure cadence; a gap is a round that
		// reported nothing or whose transmissions were still queued.
		if seqNrs[i] != seqNrs[i-1]+1 {
			continue
		}
		intervals = append(intervals, time.Duration(obsTsBySeqNr[seqNrs[i]]-obsTsBySeqNr[seqNrs[i-1]])) //nolint:gosec // G115: agreed timestamps are monotonic
	}

	res.rounds = len(seqNrs)
	res.reportsPerOracle = total / len(byOracle)
	res.roundIntervalMedian, res.roundIntervalP90 = percentiles(intervals)
	res.transmitLatencyMedian, res.transmitLatencyP90 = percentiles(latencies)
	if span := spanEnd.Sub(spanStart); span > 0 {
		res.reportsPerSecond = float64(res.reportsPerOracle) / span.Seconds()
	}
	return res
}

// percentiles returns the median and p90 of ds, zero when empty.
func percentiles(ds []time.Duration) (p50, p90 time.Duration) {
	if len(ds) == 0 {
		return 0, 0
	}
	sorted := slices.Clone(ds)
	slices.Sort(sorted)
	return sorted[len(sorted)/2], sorted[len(sorted)*9/10]
}
