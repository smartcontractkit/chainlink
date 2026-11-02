package llo_test

import (
	"crypto/ed25519"
	"fmt"
	"math/big"
	"net"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/peer"

	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/libocr/commontypes"
	memkvdb "github.com/smartcontractkit/libocr/offchainreporting2plus/ocrintegrationtesthelpers"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/csakey"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	lloconfig "github.com/smartcontractkit/chainlink-data-streams/llo/pluginconfig"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
	lloreportcodec "github.com/smartcontractkit/chainlink-data-streams/llo/reportcodec"
	mercurytransmitter "github.com/smartcontractkit/chainlink-data-streams/llo/transmitter/dataengine"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2"
)

// DON benchmark parameters. Each channel reads its own memo stream, so every
// round observes donBenchStreams values and emits donBenchChannels JSON
// reports per oracle.
const (
	donBenchChannels = 2000
	donBenchStreams  = 2000
	// donBenchRounds is how many full rounds, every channel delivered, each
	// oracle must reach before the run stops, unless LLO_DON_BENCH_ROUNDS
	// overrides it. Channels are voted in five per round, so the first full
	// round comes after donBenchChannels/5 rounds.
	donBenchRounds    = 20
	donBenchRoundsEnv = "LLO_DON_BENCH_ROUNDS"
	// donBenchDeltaRound is the round cadence floor libocr enforces, unless
	// LLO_DON_BENCH_DELTA_ROUND overrides it. The lowest the OCR3.1 production
	// sanity check accepts with the default snapshot retention, and the same
	// for v30 so the two runs compare.
	donBenchDeltaRound    = 40 * time.Millisecond
	donBenchDeltaRoundEnv = "LLO_DON_BENCH_DELTA_ROUND"
	// donBenchDeltaRequest is how long a follower waits for the leader before
	// asking a peer for the certified commit (OCR3) or the reports plus
	// precursor (OCR3.1). Several shaped round trips, so a follower whose
	// state transition still waits on blob fetches does not ask a peer for a
	// precursor it is about to compute itself.
	donBenchDeltaRequest = 400 * time.Millisecond
	donBenchEnv          = "LLO_DON_BENCH"
	// donBenchP2PHostEnv moves the oracles' p2p traffic to another loopback
	// address, so it can be shaped apart from Postgres and the transmission
	// server, which stay on 127.0.0.1. Default 127.0.0.1.
	donBenchP2PHostEnv = "LLO_DON_BENCH_P2P_HOST"
)

// TestIntegration_LLO_don_benchmark runs a 4 oracle DON with real libocr over
// loopback TCP once per plugin version, on the same channel set, and
// reports what the in-process bench package cannot: round cadence under the
// real message pattern, observation to transmission latency and report
// throughput as delivered to the transmission server.
//
// It is opt-in through LLO_DON_BENCH=1, since each version starts a node set
// and runs for a minute or more. Versions run one after the other so they do
// not contend for CPU.
//
// The oracles talk over loopback TCP, where bytes are close to free. To make
// the volume each version moves cost time, put the p2p traffic on its own
// loopback address with LLO_DON_BENCH_P2P_HOST=127.0.0.2 and shape that
// address only, so Postgres and the transmission server on 127.0.0.1 stay
// unshaped. macOS, where lo0 needs the alias first:
//
//	sudo ifconfig lo0 alias 127.0.0.2 up
//	sudo dnctl pipe 1 config delay 20 bw 50Mbit/s
//	echo "dummynet out proto tcp from 127.0.0.2 to 127.0.0.2 pipe 1" | sudo pfctl -f -
//	sudo pfctl -e
//
// and to undo, restoring the default rule set:
//
//	sudo pfctl -d
//	sudo pfctl -f /etc/pf.conf
//	sudo dnctl -q flush
//	sudo ifconfig lo0 -alias 127.0.0.2
//
// Linux, where 127.0.0.2 is already up:
//
//	sudo tc qdisc add dev lo root handle 1: prio
//	sudo tc qdisc add dev lo parent 1:3 handle 30: netem delay 20ms rate 50mbit
//	sudo tc filter add dev lo parent 1: protocol ip prio 3 u32 match ip dst 127.0.0.2/32 flowid 1:3
//
// and to undo:
//
//	sudo tc qdisc del dev lo root
//
// Wire bytes are counted on the endpoint wrapper (netCounter), so the volume
// comparison holds with or without shaping.
func TestIntegration_LLO_don_benchmark(t *testing.T) { //nolint:paralleltest // runs alone so the timings are not skewed
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

	t.Logf("DON benchmark: %d oracles, %d channels, %d streams, %d rounds, DeltaRound %s, p2p host %s", nNodes, donBenchChannels, donBenchStreams, donBenchRoundsValue(t), donBenchDeltaRoundValue(t), donBenchP2PHost())
	for _, r := range results { //nolint:paralleltest // logs results, runs no subtests
		t.Log(r.String())
		t.Logf("%s full round message sizes by bytes: %s", r.version, r.sizeProfile)
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
	// wireBytesPerRound is the bytes one oracle sends per round over the full
	// rounds, and wireBytesTotal the DON total over the whole run including
	// the channel voting rounds.
	wireBytesPerRound uint64
	wireBytesTotal    uint64
	// requestMessages and requestBytes count OCR3.1 request and response
	// messages over the whole run: reports plus precursor requests and blob
	// exchange. Zero for v30, whose endpoint has no typed messages.
	requestMessages, requestBytes uint64
	// sizeProfile is the message size histogram over the full rounds,
	// largest byte totals first, which fingerprints message types: at 2000
	// channels a report signatures message is about 130 kB, a zstd framed
	// precursor about 16 kB, a v30 observation 29 to 39 kB and a blob about
	// 3 kB.
	sizeProfile string
}

func (r donBenchResult) String() string {
	return fmt.Sprintf("%s: full_rounds=%d reports/oracle=%d first_report=%s full_round=%s round_interval median=%s p90=%s transmit_latency median=%s p90=%s reports/s/oracle=%.0f wire_B/round/oracle=%d wire_B_total=%d reqresp_msgs=%d reqresp_B=%d",
		r.version, r.rounds, r.reportsPerOracle, r.timeToFirstReport.Round(time.Millisecond), r.timeToFullRound.Round(time.Millisecond),
		r.roundIntervalMedian.Round(time.Microsecond), r.roundIntervalP90.Round(time.Microsecond),
		r.transmitLatencyMedian.Round(time.Microsecond), r.transmitLatencyP90.Round(time.Microsecond),
		r.reportsPerSecond, r.wireBytesPerRound, r.wireBytesTotal, r.requestMessages, r.requestBytes)
}

// donBenchReport is one delivered report with the wall-clock time the
// transmission server received it.
type donBenchReport struct {
	report     lloprotocol.Report
	receivedAt time.Time
}

// donBenchDeltaRoundValue is the DeltaRound for the run.
func donBenchDeltaRoundValue(t *testing.T) time.Duration {
	t.Helper()
	v := os.Getenv(donBenchDeltaRoundEnv)
	if v == "" {
		return donBenchDeltaRound
	}
	d, err := time.ParseDuration(v)
	require.NoError(t, err, "%s", donBenchDeltaRoundEnv)
	return d
}

// donBenchRoundsValue is the full round count a run waits for.
func donBenchRoundsValue(t *testing.T) int {
	t.Helper()
	v := os.Getenv(donBenchRoundsEnv)
	if v == "" {
		return donBenchRounds
	}
	n, err := strconv.Atoi(v)
	require.NoError(t, err, "%s", donBenchRoundsEnv)
	require.Positive(t, n, "%s", donBenchRoundsEnv)
	return n
}

// donBenchP2PHost is the loopback address the oracles listen and announce on.
func donBenchP2PHost() string {
	if h := os.Getenv(donBenchP2PHostEnv); h != "" {
		return h
	}
	return "127.0.0.1"
}

// p2pOnHost rewrites the node's p2p listen and announce addresses, which the
// node helper sets on 127.0.0.1, onto host.
func p2pOnHost(host string) func(c *chainlink.Config) {
	return func(c *chainlink.Config) {
		rewrite := func(addrs *[]string) {
			if addrs == nil {
				return
			}
			out := make([]string, len(*addrs))
			for i, a := range *addrs {
				_, port, err := net.SplitHostPort(a)
				if err != nil {
					out[i] = a
					continue
				}
				out[i] = net.JoinHostPort(host, port)
			}
			*addrs = out
		}
		rewrite(c.P2P.V2.ListenAddresses)
		rewrite(c.P2P.V2.AnnounceAddresses)
	}
}

func runDONBench(t *testing.T, version lloconfig.PluginVersion, salt int64, donID uint32) donBenchResult {
	const logLevel = toml.LogLevel(zapcore.ErrorLevel)
	p2pHost := donBenchP2PHost()
	onP2PHost := p2pOnHost(p2pHost)

	offchainConfig := lloprotocol.OffchainConfig{
		ProtocolVersion:                     1,
		DefaultMinReportIntervalNanoseconds: 1,
		AggregationFaultTolerance:           aggregationFaultTolerance,
		EnableObservationCompression:        true,
	}
	ocrConfigOpts := []OCRConfigOption{
		WithOffchainConfig(offchainConfig),
		func(cfg *OCRConfig) {
			cfg.DeltaRound = donBenchDeltaRoundValue(t)
			cfg.DeltaGrace = 5 * time.Millisecond
			cfg.DeltaCertifiedCommitRequest = donBenchDeltaRequest
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
		onP2PHost(c)
	})
	fullRounds := donBenchRoundsValue(t)
	bootstrapper := fmt.Sprintf("%s@%s", bootstrapPeerID, net.JoinHostPort(p2pHost, strconv.Itoa(bootstrapNodePort)))
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	// The server drops packets once this buffer is full, so size it for the
	// whole run plus the rounds in flight when it ends.
	packets := make(chan *packet, 2*donBenchChannels*fullRounds*nNodes)
	serverKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 2))
	serverPubKey := serverKey.PublicKey
	srv := NewMercuryServer(t, serverKey, packets)
	serverURL := startMercuryServer(t, srv, clientPubKeys)

	streams := make([]Stream, donBenchStreams)
	for i := range streams {
		streams[i] = Stream{id: uint32(1000 + i), baseBenchmarkPrice: decimal.NewFromInt(int64(1000 + i))}
	}

	// v31 state lives in memory so the round cost excludes disk writes, which
	// the in-process bench package does not pay either.
	wire := &netCounter{peers: nNodes}
	oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
		c.Mercury.Transmitter.Protocol = new(mercurytransmitter.MercuryTransmitterProtocolGRPC)
		c.Log.Level = new(logLevel)
		onP2PHost(c)
	}, memkvdb.NewStatelessInMemoryKeyValueDatabaseFactory(), wire)

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
		addLLOJobWithBootstrapper(t, node, configuratorAddress, bootstrapper, clientPubKeys[i], fmt.Sprintf("bench-%s", version), pluginConfig, relayType, relayConfig)
		addMemoStreamSpecs(t, node, streams)
	}

	opts := append([]OCRConfigOption{WithOracles(oracles)}, ocrConfigOpts...)
	configSetAt := time.Now()
	setProductionConfig(t, donID, steve, backend, configurator, configuratorAddress, nodes, opts...)

	// Collect until every oracle has delivered fullRounds full rounds.
	// Transmissions arrive out of order, so rounds are counted by seqNr.
	byOracle := map[string][]donBenchReport{}
	// oracle => seqNr => reports delivered for that round
	roundsByOracle := map[string]map[uint64]int{}
	var firstReportAt, fullRoundAt time.Time
	// Wire bytes over the full rounds: sampled when the first full round
	// lands, with the seqNrs seen from then on as the round count.
	var fullRoundBytes uint64
	var firstFullSeqNr uint64
	seqNrsSeen := map[uint64]struct{}{}
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
			fullRoundBytes = wire.bytes()
			firstFullSeqNr = r.SeqNr
			wire.resetSizes()
		}
		if r.SeqNr >= firstFullSeqNr && firstFullSeqNr != 0 {
			seqNrsSeen[r.SeqNr] = struct{}{}
		}

		finished := 0
		for _, rounds := range roundsByOracle {
			full := 0
			for _, n := range rounds {
				if n >= donBenchChannels {
					full++
				}
			}
			if full >= fullRounds {
				finished++
			}
		}
		if finished >= nNodes {
			break
		}
	}
	endBytes := wire.bytes()
	for _, node := range nodes {
		require.NoError(t, node.App.Stop())
	}
	require.NoError(t, bootstrapNode.App.Stop())

	res := summarizeDONBench(t, version, byOracle, roundsByOracle)
	res.timeToFirstReport = firstReportAt.Sub(configSetAt)
	res.timeToFullRound = fullRoundAt.Sub(configSetAt)
	res.wireBytesTotal = endBytes
	res.requestMessages, res.requestBytes = wire.reqCount.Load(), wire.reqBytes.Load()
	res.sizeProfile = wire.profile(8)
	if n := uint64(len(seqNrsSeen)); n > 0 {
		res.wireBytesPerRound = (endBytes - fullRoundBytes) / n / uint64(nNodes) //nolint:gosec // G115: nNodes is a small positive count
	}
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

// ---------------------------------------------------------------------------
// Wire byte accounting
// ---------------------------------------------------------------------------

// netCounter wraps the network endpoint factories of every oracle and counts
// the bytes they put on the wire. Only sends are counted, so the total over
// the DON is the volume moved once. A broadcast is charged once per peer.
// Receives are not wrapped, which would add a goroutine hop to every message
// on the round path.
type netCounter struct {
	peers      int
	sentBytes  atomic.Uint64
	sentCount  atomic.Uint64
	bcastBytes atomic.Uint64
	// reqCount and reqBytes track OCR3.1 request and response messages,
	// each charged once per receiving peer.
	reqCount atomic.Uint64
	reqBytes atomic.Uint64
	// sizes histograms messages by kB bucket, message count and bytes, each
	// charged once per receiving peer.
	mu    sync.Mutex
	sizes map[int]*sizeBucket
}

type sizeBucket struct{ count, bytes uint64 }

func (c *netCounter) record(n, copies int) {
	kb := n >> 10
	c.mu.Lock()
	if c.sizes == nil {
		c.sizes = map[int]*sizeBucket{}
	}
	b := c.sizes[kb]
	if b == nil {
		b = &sizeBucket{}
		c.sizes[kb] = b
	}
	b.count += uint64(copies)             //nolint:gosec // G115: small peer count
	b.bytes += uint64(n) * uint64(copies) //nolint:gosec // G115: payload length
	c.mu.Unlock()
}

// resetSizes clears the histogram, so it covers the rounds that follow.
func (c *netCounter) resetSizes() {
	c.mu.Lock()
	c.sizes = nil
	c.mu.Unlock()
}

// profile lists the top buckets by bytes as "kB:count/MB".
func (c *netCounter) profile(top int) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	kbs := make([]int, 0, len(c.sizes))
	for kb := range c.sizes {
		kbs = append(kbs, kb)
	}
	sort.Slice(kbs, func(i, j int) bool { return c.sizes[kbs[i]].bytes > c.sizes[kbs[j]].bytes })
	var parts []string
	for i, kb := range kbs {
		if i == top {
			break
		}
		b := c.sizes[kb]
		parts = append(parts, fmt.Sprintf("%dkB:%d/%.1fMB", kb, b.count, float64(b.bytes)/1e6))
	}
	return strings.Join(parts, " ")
}

var _ ocr2.LLONetworkEndpointWrapper = (*netCounter)(nil)

func (c *netCounter) WrapEndpointFactory(f ocrtypes.BinaryNetworkEndpointFactory) ocrtypes.BinaryNetworkEndpointFactory {
	return countingEndpointFactory{f, c}
}

func (c *netCounter) WrapEndpoint2Factory(f ocrtypes.BinaryNetworkEndpoint2Factory) ocrtypes.BinaryNetworkEndpoint2Factory {
	return countingEndpoint2Factory{f, c}
}

func (c *netCounter) sent(n int) {
	c.sentBytes.Add(uint64(n)) //nolint:gosec // G115: payload length
	c.sentCount.Add(1)
	c.record(n, 1)
}

func (c *netCounter) broadcast(n int) {
	c.bcastBytes.Add(uint64(n) * uint64(c.peers-1)) //nolint:gosec // G115: payload length
	c.sentCount.Add(uint64(c.peers - 1))            //nolint:gosec // G115: small peer count
	c.record(n, c.peers-1)
}

// request charges an OCR3.1 message to the request and response counters
// when it is one.
func (c *netCounter) request(msg ocrtypes.OutboundBinaryMessage, copies int) {
	switch msg.(type) {
	case ocrtypes.OutboundBinaryMessageRequest, ocrtypes.OutboundBinaryMessageResponse:
		c.reqCount.Add(uint64(copies))                                 //nolint:gosec // G115: small peer count
		c.reqBytes.Add(uint64(len(msg.GetPayload())) * uint64(copies)) //nolint:gosec // G115: payload length
	}
}

// bytes returns the total bytes sent across the DON so far.
func (c *netCounter) bytes() uint64 { return c.sentBytes.Load() + c.bcastBytes.Load() }

type countingEndpointFactory struct {
	ocrtypes.BinaryNetworkEndpointFactory
	c *netCounter
}

func (f countingEndpointFactory) NewEndpoint(cd ocrtypes.ConfigDigest, peerIDs []string, v2bootstrappers []commontypes.BootstrapperLocator, fLimit int, limits ocrtypes.BinaryNetworkEndpointLimits) (commontypes.BinaryNetworkEndpoint, error) {
	ep, err := f.BinaryNetworkEndpointFactory.NewEndpoint(cd, peerIDs, v2bootstrappers, fLimit, limits)
	if err != nil {
		return nil, err
	}
	return countingEndpoint{ep, f.c}, nil
}

type countingEndpoint struct {
	commontypes.BinaryNetworkEndpoint
	c *netCounter
}

func (e countingEndpoint) SendTo(payload []byte, to commontypes.OracleID) {
	e.c.sent(len(payload))
	e.BinaryNetworkEndpoint.SendTo(payload, to)
}

func (e countingEndpoint) Broadcast(payload []byte) {
	e.c.broadcast(len(payload))
	e.BinaryNetworkEndpoint.Broadcast(payload)
}

type countingEndpoint2Factory struct {
	ocrtypes.BinaryNetworkEndpoint2Factory
	c *netCounter
}

func (f countingEndpoint2Factory) NewEndpoint(cd ocrtypes.ConfigDigest, peerIDs []string, v2bootstrappers []commontypes.BootstrapperLocator, defaultPriorityConfig, lowPriorityConfig ocrtypes.BinaryNetworkEndpoint2Config) (ocrtypes.BinaryNetworkEndpoint2, error) {
	ep, err := f.BinaryNetworkEndpoint2Factory.NewEndpoint(cd, peerIDs, v2bootstrappers, defaultPriorityConfig, lowPriorityConfig)
	if err != nil {
		return nil, err
	}
	return countingEndpoint2{ep, f.c}, nil
}

type countingEndpoint2 struct {
	ocrtypes.BinaryNetworkEndpoint2
	c *netCounter
}

func (e countingEndpoint2) SendTo(msg ocrtypes.OutboundBinaryMessage, to commontypes.OracleID) {
	e.c.sent(len(msg.GetPayload()))
	e.c.request(msg, 1)
	e.BinaryNetworkEndpoint2.SendTo(msg, to)
}

func (e countingEndpoint2) Broadcast(msg ocrtypes.OutboundBinaryMessage) {
	e.c.broadcast(len(msg.GetPayload()))
	e.c.request(msg, e.c.peers-1)
	e.BinaryNetworkEndpoint2.Broadcast(msg)
}
