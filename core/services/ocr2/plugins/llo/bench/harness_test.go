// Package bench contains a comparative micro-benchmark between the OCR3.0 LLO
// plugin (chainlink-data-streams/llo/v30) and the OCR3.1 LLO plugin
// (chainlink-data-streams/llo/dev/v31).
//
// The two plugins implement identical LLO application logic on top of different
// OCR protocols. The performance question this benchmark answers is what that
// protocol difference costs at steady state:
//
//   - v30 (OCR3.0) carries all of its state in a single Outcome blob that is
//     decoded from the previous round and re-encoded every round. Its per-round
//     cost therefore scales with the *total* state size (channels + streams).
//   - v31 (OCR3.1) keeps its state in a replicated KeyValueState (a pebble
//     database in production) and only reads/writes the keys that change. Its
//     per-round cost scales with per-round *churn*, not total state.
//
// Both plugins are driven exclusively through their exported ReportingPlugin
// APIs, so nothing in the read-only chainlink-data-streams or libocr modules is
// modified. v31's KeyValueState is backed by libocr's in-memory
// KeyValueDatabase (offchainreporting2plus/ocrintegrationtesthelpers) rather
// than the production pebble factory: pebble commits with fsync every round,
// which dwarfs and obscures the plugin's own work. The in-memory store isolates
// plugin CPU/allocation cost, making it directly comparable to v30's in-memory
// Outcome blob. Neither plugin's oracle-level OCR3 protocol Database
// (core/services/llo/delegate.go) is modeled here.
//
// Scope note: blob offloading (v31 offloads large observation payloads to
// libocr blobs above BlobThreshold) cannot be exercised outside libocr's oracle
// runtime — a BlobHandle cannot be constructed by application code (see the
// comment in v31/plugin_test.go). We therefore disable blob offloading here so
// observations always inline, which also makes the observation-transport cost
// directly comparable to v30 (which is always inline). Blob behavior is covered
// by the integration tests.
package bench

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	llodatasource "github.com/smartcontractkit/chainlink-data-streams/llo/datasource"
	llov31 "github.com/smartcontractkit/chainlink-data-streams/llo/dev/v31"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
	llov30 "github.com/smartcontractkit/chainlink-data-streams/llo/v30"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	memkvdb "github.com/smartcontractkit/libocr/offchainreporting2plus/ocrintegrationtesthelpers"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	corello "github.com/smartcontractkit/chainlink/v2/core/services/llo"
)

// benchConfigDigest is a fixed config digest shared by both plugins. v31 passes
// it to the KeyValueDatabase factory.
var benchConfigDigest = ocrtypes.ConfigDigest{'b', 'e', 'n', 'c', 'h'}

const (
	// maxDurationObservation bounds the DataSource.Observe context in both
	// plugins. It must be > 0 or the observation context is already expired.
	maxDurationObservation = 5 * time.Second
	// channelsPerRound is the protocol cap on how many channel definitions an
	// observation may vote to add per round (MaxObservationUpdateChannelDefinitionsLength).
	// Establishing C channels therefore takes ~ceil(C/channelsPerRound) rounds.
	channelsPerRound = 5
	// warmupRoundSlack is added on top of the minimum rounds needed to add every
	// channel, to allow the last batch to become reportable and to absorb the
	// bootstrap round.
	warmupRoundSlack = 32
)

// ---------------------------------------------------------------------------
// Workload
// ---------------------------------------------------------------------------

// workload describes the size of the report-production problem: numChannels
// channels, each observing streamsPerChannel distinct streams via the median
// aggregator, emitting a JSON report.
//
// Stream IDs are unique across the whole workload, so the total number of
// distinct observed streams is numChannels*streamsPerChannel. Report format is
// held constant (JSON) because it is orthogonal to the v30-vs-v31 delta: both
// plugins use the identical report codec, so its cost cancels out of the
// comparison.
type workload struct {
	numChannels       int
	streamsPerChannel int
}

func (w workload) String() string {
	return fmt.Sprintf("ch=%d_str=%d", w.numChannels, w.streamsPerChannel)
}

// channelDefinitions builds the workload's channel definitions and returns the
// full set of stream IDs referenced.
func (w workload) channelDefinitions() (llotypes.ChannelDefinitions, []llotypes.StreamID) {
	defs := make(llotypes.ChannelDefinitions, w.numChannels)
	var streamIDs []llotypes.StreamID
	var sid llotypes.StreamID
	for c := 0; c < w.numChannels; c++ {
		streams := make([]llotypes.Stream, 0, w.streamsPerChannel)
		for s := 0; s < w.streamsPerChannel; s++ {
			sid++
			streams = append(streams, llotypes.Stream{StreamID: sid, Aggregator: llotypes.AggregatorMedian})
			streamIDs = append(streamIDs, sid)
		}
		defs[llotypes.ChannelID(c+1)] = llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatJSON,
			Streams:      streams,
		}
	}
	return defs, streamIDs
}

// ---------------------------------------------------------------------------
// Mocks (shared, version-agnostic)
// ---------------------------------------------------------------------------

type mockChannelDefinitionCache struct{ defs llotypes.ChannelDefinitions }

func (m *mockChannelDefinitionCache) Definitions(llotypes.ChannelDefinitions) llotypes.ChannelDefinitions {
	return m.defs
}
func (m *mockChannelDefinitionCache) Start(context.Context) error    { return nil }
func (m *mockChannelDefinitionCache) Close() error                   { return nil }
func (m *mockChannelDefinitionCache) Ready() error                   { return nil }
func (m *mockChannelDefinitionCache) HealthReport() map[string]error { return nil }
func (m *mockChannelDefinitionCache) Name() string                   { return "benchChannelDefinitionCache" }

// staticDataSource fills every requested stream with a fixed decimal value.
// This keeps the DataSource out of the measured critical path (no I/O, no
// allocation-heavy pipeline) so the benchmark isolates plugin cost.
type staticDataSource struct{ value *lloprotocol.Decimal }

func newStaticDataSource() *staticDataSource {
	return &staticDataSource{value: lloprotocol.ToDecimal(decimal.NewFromInt(123456))}
}

func (d *staticDataSource) Observe(_ context.Context, sv lloprotocol.StreamValues, _ llodatasource.DSOpts) error {
	for k := range sv {
		sv[k] = d.value
	}
	return nil
}

type mockShouldRetireCache struct{}

func (mockShouldRetireCache) ShouldRetire(ocrtypes.ConfigDigest) (bool, error) { return false, nil }

type mockOnchainConfigCodec struct{}

func (mockOnchainConfigCodec) Decode([]byte) (lloprotocol.OnchainConfig, error) {
	return lloprotocol.OnchainConfig{}, nil
}
func (mockOnchainConfigCodec) Encode(lloprotocol.OnchainConfig) ([]byte, error) { return nil, nil }

// ---------------------------------------------------------------------------
// Plugin construction
// ---------------------------------------------------------------------------

func reportCodecs() map[llotypes.ReportFormat]lloprotocol.ReportCodec {
	// The same production codec set both plugins use (delegate.go). Only the
	// JSON codec is exercised by this workload.
	return corello.NewReportCodecs(logger.Nop(), 0)
}

// benchOffchainConfig selects protocol version 1 with a 1ns minimum report
// interval. Version 1 makes both plugins use full nanosecond timestamp
// resolution for the JSON report format (version 0 truncates v30's timestamps
// to whole seconds, which would prevent reporting within a single wall-clock
// second and diverge from v31). The 1ns interval effectively reports every
// round while keeping both plugins on identical reportability rules.
func benchOffchainConfig() []byte {
	b, err := lloprotocol.OffchainConfig{
		ProtocolVersion:                     1,
		DefaultMinReportIntervalNanoseconds: 1,
	}.Encode()
	if err != nil {
		panic(err)
	}
	return b
}

func pluginConfig(n, f int) ocr3types.ReportingPluginConfig {
	return ocr3types.ReportingPluginConfig{
		ConfigDigest:           benchConfigDigest,
		N:                      n,
		F:                      f,
		MaxDurationObservation: maxDurationObservation,
		OffchainConfig:         benchOffchainConfig(),
	}
}

func buildV30(tb testing.TB, defs llotypes.ChannelDefinitions, n, f int) ocr3types.ReportingPlugin[llotypes.ReportInfo] {
	tb.Helper()
	factory := llov30.NewPluginFactory(llov30.PluginFactoryParams{
		Config:                 llov30.Config{VerboseLogging: false},
		ShouldRetireCache:      mockShouldRetireCache{},
		RetirementReportCodec:  lloprotocol.StandardRetirementReportCodec{},
		ChannelDefinitionCache: &mockChannelDefinitionCache{defs: defs},
		DataSource:             newStaticDataSource(),
		Logger:                 logger.Nop(),
		OnchainConfigCodec:     mockOnchainConfigCodec{},
		ReportCodecs:           reportCodecs(),
	})
	p, _, err := factory.NewReportingPlugin(context.Background(), pluginConfig(n, f))
	require.NoError(tb, err)
	return p
}

func buildV31(tb testing.TB, defs llotypes.ChannelDefinitions, n, f int) (ocr3_1types.ReportingPlugin[llotypes.ReportInfo], ocr3_1types.KeyValueDatabase) {
	tb.Helper()
	factory := llov31.NewPluginFactory(llov31.PluginFactoryParams{
		Config:                 llov31.Config{VerboseLogging: false},
		ShouldRetireCache:      mockShouldRetireCache{},
		RetirementReportCodec:  lloprotocol.StandardRetirementReportCodec{},
		ChannelDefinitionCache: &mockChannelDefinitionCache{defs: defs},
		DataSource:             newStaticDataSource(),
		Logger:                 logger.Nop(),
		OnchainConfigCodec:     mockOnchainConfigCodec{},
		ReportCodecs:           reportCodecs(),
		// Negative disables blob offloading; observations always inline. See the
		// package-level scope note.
		BlobThreshold: -1,
	})
	p, _, err := factory.NewReportingPlugin(context.Background(), pluginConfig(n, f), nil)
	require.NoError(tb, err)

	// libocr's in-memory KeyValueDatabase (the same helper v31's integration
	// tests use): a btree behind the production KeyValueDatabaseFactory
	// interface, whose Commit applies to memory with no WAL/fsync. This isolates
	// the plugin's CPU/allocation cost from storage-engine cost, so the v31
	// numbers are directly comparable to v30's in-memory Outcome blob.
	dbFactory := memkvdb.NewStatelessInMemoryKeyValueDatabaseFactory()
	db, err := dbFactory.NewKeyValueDatabase(benchConfigDigest)
	require.NoError(tb, err)
	tb.Cleanup(func() { _ = db.Close() })
	return p, db
}

// ---------------------------------------------------------------------------
// Round drivers
// ---------------------------------------------------------------------------

func attributedObservation(observer int, obs []byte) ocrtypes.AttributedObservation {
	return ocrtypes.AttributedObservation{Observer: commontypes.OracleID(observer), Observation: obs} //nolint:gosec // G115: observer is a small oracle index
}

// replicate builds n AttributedObservations from a single serialized
// observation. Real oracles observe slightly different values; identical copies
// are representative for a benchmark and keep the workload deterministic (the
// aggregation still processes n observations per stream).
func replicate(obs []byte, n int) []ocrtypes.AttributedObservation {
	aos := make([]ocrtypes.AttributedObservation, 0, n)
	for i := range n {
		aos = append(aos, attributedObservation(i, obs))
	}
	return aos
}

// v30Round drives one full v30 round (Observation → Outcome → Reports) against
// a fixed previous outcome. It returns the produced outcome, the reports, and
// the serialized observation.
func v30Round(tb testing.TB, p ocr3types.ReportingPlugin[llotypes.ReportInfo], seqNr uint64, prevOutcome []byte, n int) (ocr3types.Outcome, []ocr3types.ReportPlus[llotypes.ReportInfo], []byte) {
	ctx := context.Background()
	outctx := ocr3types.OutcomeContext{SeqNr: seqNr, PreviousOutcome: prevOutcome}
	obs, err := p.Observation(ctx, outctx, nil)
	require.NoError(tb, err)
	aos := replicate(obs, n)
	outcome, err := p.Outcome(ctx, outctx, nil, aos)
	require.NoError(tb, err)
	reports, err := p.Reports(ctx, seqNr, outcome)
	require.NoError(tb, err)
	return outcome, reports, obs
}

// v31Round drives one full v31 round (Observation → StateTransition → Reports)
// against the KeyValueDatabase, mirroring libocr's per-seqNr transaction
// lifecycle: Observation reads a snapshot of the state committed after seqNr-1;
// StateTransition mutates a batch that is committed to advance the state to
// seqNr. It mutates db.
func v31Round(tb testing.TB, p ocr3_1types.ReportingPlugin[llotypes.ReportInfo], db ocr3_1types.KeyValueDatabase, seqNr uint64, n int) ([]ocr3types.ReportPlus[llotypes.ReportInfo], []byte) {
	ctx := context.Background()

	var obs []byte
	if seqNr > 1 {
		rtx, err := db.NewReadTransaction()
		require.NoError(tb, err)
		obs, err = p.Observation(ctx, seqNr, ocrtypes.AttributedQuery{}, rtx, nil)
		rtx.Discard()
		require.NoError(tb, err)
	}
	aos := bootOrReplicate(obs, n, seqNr)

	wtx, err := db.NewReadWriteTransaction()
	require.NoError(tb, err)
	prec, err := p.StateTransition(ctx, seqNr, ocrtypes.AttributedQuery{}, aos, wtx, nil)
	if err != nil {
		wtx.Discard()
		require.NoError(tb, err)
	}
	require.NoError(tb, wtx.Commit())

	reports, err := p.Reports(ctx, seqNr, prec)
	require.NoError(tb, err)
	return reports, obs
}

// bootOrReplicate returns the bootstrap observation set (empty observations)
// for the first round, otherwise n copies of obs.
func bootOrReplicate(obs []byte, n int, seqNr uint64) []ocrtypes.AttributedObservation {
	if seqNr == 1 {
		// First round: empty observations, only 2f+1 are needed; n is fine.
		return replicate(nil, n)
	}
	return replicate(obs, n)
}

// ---------------------------------------------------------------------------
// Warmup
// ---------------------------------------------------------------------------

// warmupRounds returns the round budget needed to establish `channels`
// channels (5 per round) plus slack for reportability and bootstrap.
func warmupRounds(channels int) uint64 {
	return uint64(channels/channelsPerRound + warmupRoundSlack) //nolint:gosec // G115: small non-negative round budget
}

// warmV30 drives rounds until all `channels` channels are established and
// reportable, then returns the previous-outcome and seqNr that reliably yield a
// full report set. That (prevOutcome, seqNr) pair is reused for every measured
// iteration: the stored watermarks are in the past, so each measured round
// (with a fresh wall-clock observation timestamp) reports every channel.
func warmV30(tb testing.TB, p ocr3types.ReportingPlugin[llotypes.ReportInfo], n, channels int) (steadyPrev ocr3types.Outcome, steadySeq uint64) {
	tb.Helper()
	// seqNr 1 bootstraps with empty observations and no previous outcome.
	prev, _, _ := v30Round(tb, p, 1, nil, n)
	maxRounds := warmupRounds(channels)
	for seqNr := uint64(2); seqNr <= maxRounds; seqNr++ {
		outcome, reports, _ := v30Round(tb, p, seqNr, prev, n)
		if len(reports) >= channels {
			return prev, seqNr
		}
		prev = outcome
	}
	tb.Fatalf("v30 did not reach %d reportable channels within %d warmup rounds", channels, maxRounds)
	return nil, 0
}

// warmV31 drives rounds until all `channels` channels are established and
// reportable, returning the next seqNr to use. Unlike v30, v31 must chain
// (state lives in the mutated KeyValueDatabase), so measured iterations continue from
// this seqNr.
func warmV31(tb testing.TB, p ocr3_1types.ReportingPlugin[llotypes.ReportInfo], db ocr3_1types.KeyValueDatabase, n, channels int) (nextSeq uint64) {
	tb.Helper()
	maxRounds := warmupRounds(channels)
	for seqNr := uint64(1); seqNr <= maxRounds; seqNr++ {
		reports, _ := v31Round(tb, p, db, seqNr, n)
		if len(reports) >= channels {
			return seqNr + 1
		}
	}
	tb.Fatalf("v31 did not reach %d reportable channels within %d warmup rounds", channels, maxRounds)
	return 0
}

// reportBytes sums the serialized report payloads.
func reportBytes(reports []ocr3types.ReportPlus[llotypes.ReportInfo]) int {
	var total int
	for _, r := range reports {
		total += len(r.ReportWithInfo.Report)
	}
	return total
}

// ---------------------------------------------------------------------------
// Counting KV wrappers (size instrumentation, not used on the hot path)
// ---------------------------------------------------------------------------

// countingRW wraps a KeyValueState read-write transaction and tallies the volume of
// KeyValueState I/O a single v31 round performs. It is used only for the
// one-off size probe, never inside the timed loop, so its counter overhead
// does not pollute latency measurements.
type countingRW struct {
	inner                                               ocr3_1types.KeyValueStateReadWriter
	reads, readBytes, writes, writeBytes, deletes, keys int
}

func (c *countingRW) Read(key []byte) ([]byte, error) {
	v, err := c.inner.Read(key)
	c.reads++
	c.readBytes += len(v)
	return v, err
}

func (c *countingRW) Write(key, value []byte) error {
	c.writes++
	c.keys++
	c.writeBytes += len(key) + len(value)
	return c.inner.Write(key, value)
}

func (c *countingRW) Delete(key []byte) error {
	c.deletes++
	c.keys++
	return c.inner.Delete(key)
}

var _ ocr3_1types.KeyValueStateReadWriter = (*countingRW)(nil)

// countingReader wraps a read transaction and tallies read volume.
type countingReader struct {
	inner            ocr3_1types.KeyValueStateReader
	reads, readBytes int
}

func (c *countingReader) Read(key []byte) ([]byte, error) {
	v, err := c.inner.Read(key)
	c.reads++
	c.readBytes += len(v)
	return v, err
}

var _ ocr3_1types.KeyValueStateReader = (*countingReader)(nil)
