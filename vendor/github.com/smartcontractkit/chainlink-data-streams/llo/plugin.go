package llo

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

// TODO: Split out this file and write unit tests: https://smartcontract-it.atlassian.net/browse/MERC-3524

// Additional limits so we can more effectively bound the size of observations
// NOTE: These are hardcoded because these exact values are relied upon as a
// property of coming to consensus, it's too dangerous to make these
// configurable on a per-node basis. It may be possible to add them to the
// OffchainConfig if they need to be changed dynamically and in a
// backwards-compatible way.
const (
	// Maximum amount of channels that can be added per round (if more than
	// this needs to be added, it will be added in batches until everything is
	// up-to-date)
	MaxObservationRemoveChannelIDsLength = 5
	// Maximum amount of channels that can be removed per round (if more than
	// this needs to be removed, it will be removed in batches until everything
	// is up-to-date)
	MaxObservationUpdateChannelDefinitionsLength = 5
	// Maximum number of streams that can be observed per round
	// TODO: This needs to be implemented on the Observation side so we don't
	// even generate an observation that fails this
	MaxObservationStreamValuesLength = 10_000
	// MaxOutcomeChannelDefinitionsLength is the maximum number of channels that
	// can be supported
	// TODO: This needs to be implemented on the Observation side so we don't
	// even generate an observation that fails this
	MaxOutcomeChannelDefinitionsLength = 10_000
)

type DSOpts interface {
	VerboseLogging() bool
	SeqNr() uint64
}

type dsOpts struct {
	verboseLogging bool
	seqNr          uint64
}

func (o dsOpts) VerboseLogging() bool {
	return o.verboseLogging
}

func (o dsOpts) SeqNr() uint64 {
	return o.seqNr
}

type DataSource interface {
	// For each known streamID, Observe should set the observed value in the
	// passed streamValues.
	// If an observation fails, or the stream is unknown, no value should be
	// set.
	Observe(ctx context.Context, streamValues StreamValues, opts DSOpts) error
}

// Protocol instances start in either the staging or production stage. They
// may later be retired and "hand over" their work to another protocol instance
// that will move from the staging to the production stage.
const (
	LifeCycleStageStaging    llotypes.LifeCycleStage = "staging"
	LifeCycleStageProduction llotypes.LifeCycleStage = "production"
	LifeCycleStageRetired    llotypes.LifeCycleStage = "retired"
)

type RetirementReport struct {
	// Carries validity time stamps between protocol instances to ensure there
	// are no gaps
	ValidAfterSeconds map[llotypes.ChannelID]uint32
}

type ShouldRetireCache interface { // reads asynchronously from onchain ConfigurationStore
	// Should the protocol instance retire according to the configuration
	// contract?
	// See: https://github.com/smartcontractkit/mercury-v1-sketch/blob/main/onchain/src/ConfigurationStore.sol#L18
	ShouldRetire() (bool, error)
}

// The predecessor protocol instance stores its attested retirement report in
// this cache (locally, offchain), so it can be fetched by the successor
// protocol instance.
//
// PredecessorRetirementReportCache is populated by the old protocol instance
// writing to it and the new protocol instance reading from it.
//
// The sketch envisions it being implemented as a single object that is shared
// between different protocol instances.
type PredecessorRetirementReportCache interface {
	AttestedRetirementReport(predecessorConfigDigest ocr2types.ConfigDigest) ([]byte, error)
	CheckAttestedRetirementReport(predecessorConfigDigest ocr2types.ConfigDigest, attestedRetirementReport []byte) (RetirementReport, error)
}

type ChannelDefinitionCache interface {
	Definitions() llotypes.ChannelDefinitions
}

// MakeChannelHash is used for mapping ChannelDefinitionWithIDs
func MakeChannelHash(cd ChannelDefinitionWithID) ChannelHash {
	h := sha256.New()
	merr := errors.Join(
		binary.Write(h, binary.BigEndian, cd.ChannelID),
		binary.Write(h, binary.BigEndian, cd.ReportFormat),
		binary.Write(h, binary.BigEndian, uint32(len(cd.Streams))),
	)
	for _, strm := range cd.Streams {
		merr = errors.Join(merr, binary.Write(h, binary.BigEndian, strm.StreamID))
		merr = errors.Join(merr, binary.Write(h, binary.BigEndian, strm.Aggregator))
	}
	if merr != nil {
		// This should never happen
		panic(merr)
	}
	h.Write(cd.Opts)
	var result [32]byte
	h.Sum(result[:0])
	return result
}

// A ReportingPlugin allows plugging custom logic into the OCR3 protocol. The OCR
// protocol handles cryptography, networking, ensuring that a sufficient number
// of nodes is in agreement about any report, transmitting the report to the
// contract, etc... The ReportingPlugin handles application-specific logic. To do so,
// the ReportingPlugin defines a number of callbacks that are called by the OCR
// protocol logic at certain points in the protocol's execution flow. The report
// generated by the ReportingPlugin must be in a format understood by contract that
// the reports are transmitted to.
//
// We assume that each correct node participating in the protocol instance will
// be running the same ReportingPlugin implementation. However, not all nodes may be
// correct; up to f nodes be faulty in arbitrary ways (aka byzantine faults).
// For example, faulty nodes could be down, have intermittent connectivity
// issues, send garbage messages, or be controlled by an adversary.
//
// For a protocol round where everything is working correctly, followers will
// call Observation, Outcome, and Reports. For each report,
// ShouldAcceptAttestedReport will be called as well. If
// ShouldAcceptAttestedReport returns true, ShouldTransmitAcceptedReport will
// be called. However, an ReportingPlugin must also correctly handle the case where
// faults occur.
//
// In particular, an ReportingPlugin must deal with cases where:
//
// - only a subset of the functions on the ReportingPlugin are invoked for a given
// round
//
// - an arbitrary number of seqnrs has been skipped between invocations of the
// ReportingPlugin
//
// - the observation returned by Observation is not included in the list of
// AttributedObservations passed to Report
//
// - a query or observation is malformed. (For defense in depth, it is also
// recommended that malformed outcomes are handled gracefully.)
//
// - instances of the ReportingPlugin run by different oracles have different call
// traces. E.g., the ReportingPlugin's Observation function may have been invoked on
// node A, but not on node B.
//
// All functions on an ReportingPlugin should be thread-safe.
//
// All functions that take a context as their first argument may still do cheap
// computations after the context expires, but should stop any blocking
// interactions with outside services (APIs, database, ...) and return as
// quickly as possible. (Rough rule of thumb: any such computation should not
// take longer than a few ms.) A blocking function may block execution of the
// entire protocol instance on its node!
//
// For a given OCR protocol instance, there can be many (consecutive) instances
// of an ReportingPlugin, e.g. due to software restarts. If you need ReportingPlugin state
// to survive across restarts, you should store it in the Outcome or persist it.
// A ReportingPlugin instance will only ever serve a single protocol instance.
var _ ocr3types.ReportingPluginFactory[llotypes.ReportInfo] = &PluginFactory{}

func NewPluginFactory(cfg Config, prrc PredecessorRetirementReportCache, src ShouldRetireCache, cdc ChannelDefinitionCache, ds DataSource, lggr logger.Logger, codecs map[llotypes.ReportFormat]ReportCodec) *PluginFactory {
	return &PluginFactory{
		cfg, prrc, src, cdc, ds, lggr, codecs,
	}
}

type Config struct {
	// Enables additional logging that might be expensive, e.g. logging entire
	// channel definitions on every round or other very large structs
	VerboseLogging bool
}

type PluginFactory struct {
	Config                           Config
	PredecessorRetirementReportCache PredecessorRetirementReportCache
	ShouldRetireCache                ShouldRetireCache
	ChannelDefinitionCache           ChannelDefinitionCache
	DataSource                       DataSource
	Logger                           logger.Logger
	Codecs                           map[llotypes.ReportFormat]ReportCodec
}

func (f *PluginFactory) NewReportingPlugin(cfg ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[llotypes.ReportInfo], ocr3types.ReportingPluginInfo, error) {
	offchainCfg, err := DecodeOffchainConfig(cfg.OffchainConfig)
	if err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, fmt.Errorf("NewReportingPlugin failed to decode offchain config; got: 0x%x (len: %d); %w", cfg.OffchainConfig, len(cfg.OffchainConfig), err)
	}

	return &Plugin{
			f.Config,
			offchainCfg.PredecessorConfigDigest,
			cfg.ConfigDigest,
			f.PredecessorRetirementReportCache,
			f.ShouldRetireCache,
			f.ChannelDefinitionCache,
			f.DataSource,
			f.Logger,
			cfg.F,
			protoObservationCodec{},
			protoOutcomeCodec{},
			f.Codecs,
		}, ocr3types.ReportingPluginInfo{
			Name: "LLO",
			Limits: ocr3types.ReportingPluginLimits{
				MaxQueryLength:       0,
				MaxObservationLength: ocr3types.MaxMaxObservationLength, // TODO: use tighter bound MERC-3524
				MaxOutcomeLength:     ocr3types.MaxMaxOutcomeLength,     // TODO: use tighter bound MERC-3524
				MaxReportLength:      ocr3types.MaxMaxReportLength,      // TODO: use tighter bound MERC-3524
				MaxReportCount:       ocr3types.MaxMaxReportCount,       // TODO: use tighter bound MERC-3524
			},
		}, nil
}

var _ ocr3types.ReportingPlugin[llotypes.ReportInfo] = &Plugin{}

type ReportCodec interface {
	// Encode may be lossy, so no Decode function is expected
	Encode(Report, llotypes.ChannelDefinition) ([]byte, error)
}

type Plugin struct {
	Config                           Config
	PredecessorConfigDigest          *types.ConfigDigest
	ConfigDigest                     types.ConfigDigest
	PredecessorRetirementReportCache PredecessorRetirementReportCache
	ShouldRetireCache                ShouldRetireCache
	ChannelDefinitionCache           ChannelDefinitionCache
	DataSource                       DataSource
	Logger                           logger.Logger
	F                                int
	ObservationCodec                 ObservationCodec
	OutcomeCodec                     OutcomeCodec
	Codecs                           map[llotypes.ReportFormat]ReportCodec
}

// Query creates a Query that is sent from the leader to all follower nodes
// as part of the request for an observation. Be careful! A malicious leader
// could equivocate (i.e. send different queries to different followers.)
// Many applications will likely be better off always using an empty query
// if the oracles don't need to coordinate on what to observe (e.g. in case
// of a price feed) or the underlying data source offers an (eventually)
// consistent view to different oracles (e.g. in case of observing a
// blockchain).
//
// You may assume that the outctx.SeqNr is increasing monotonically (though
// *not* strictly) across the lifetime of a protocol instance and that
// outctx.previousOutcome contains the consensus outcome with sequence
// number (outctx.SeqNr-1).
func (p *Plugin) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
	return nil, nil
}

type Observation struct {
	// Attested (i.e. signed by f+1 oracles) retirement report from predecessor
	// protocol instance
	AttestedPredecessorRetirement []byte
	// Should this protocol instance be retired?
	ShouldRetire bool
	// Timestamp from when observation is made
	// Note that this is the timestamp immediately before we initiate any
	// observations
	UnixTimestampNanoseconds int64
	// Votes to remove/add channels. Subject to MAX_OBSERVATION_*_LENGTH limits
	RemoveChannelIDs map[llotypes.ChannelID]struct{}
	// Votes to add or replace channel definitions
	UpdateChannelDefinitions llotypes.ChannelDefinitions
	// Observed (numeric) stream values. Subject to
	// MaxObservationStreamValuesLength limit
	StreamValues StreamValues
}

// Observation gets an observation from the underlying data source. Returns
// a value or an error.
//
// You may assume that the outctx.SeqNr is increasing monotonically (though
// *not* strictly) across the lifetime of a protocol instance and that
// outctx.previousOutcome contains the consensus outcome with sequence
// number (outctx.SeqNr-1).
//
// Should return a serialized Observation struct.
func (p *Plugin) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query) (types.Observation, error) {
	// NOTE: First sequence number is always 1 (0 is invalid)
	if outctx.SeqNr < 1 {
		return types.Observation{}, fmt.Errorf("got invalid seqnr=%d, must be >=1", outctx.SeqNr)
	} else if outctx.SeqNr == 1 {
		// First round always has empty PreviousOutcome
		// Don't bother observing on the first ever round, because the result
		// will never be used anyway.
		// See case at the top of Outcome()
		return types.Observation{}, nil
	}
	// Second round will have no channel definitions yet, but may vote to add
	// them

	// QUESTION: is there a way to have this captured in EAs so we get something
	// closer to the source?
	nowNanoseconds := time.Now().UnixNano()

	previousOutcome, err := p.OutcomeCodec.Decode(outctx.PreviousOutcome)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling previous outcome: %w", err)
	}

	var attestedRetirementReport []byte
	// Only try to fetch this from the cache if this instance if configured
	// with a predecessor and we're still in the staging stage.
	if p.PredecessorConfigDigest != nil && previousOutcome.LifeCycleStage == LifeCycleStageStaging {
		var err2 error
		attestedRetirementReport, err2 = p.PredecessorRetirementReportCache.AttestedRetirementReport(*p.PredecessorConfigDigest)
		if err2 != nil {
			return nil, fmt.Errorf("error fetching attested retirement report from cache: %w", err2)
		}
	}

	shouldRetire, err := p.ShouldRetireCache.ShouldRetire()
	if err != nil {
		return nil, fmt.Errorf("error fetching shouldRetire from cache: %w", err)
	}

	// vote to remove channel ids if they're in the previous outcome
	// ChannelDefinitions or ValidAfterSeconds
	removeChannelIDs := map[llotypes.ChannelID]struct{}{}
	// vote to add channel definitions that aren't present in the previous
	// outcome ChannelDefinitions
	// FIXME: Why care about ValidAfterSeconds here?
	var updateChannelDefinitions llotypes.ChannelDefinitions
	{
		// NOTE: Be careful using maps, since key ordering is randomized! All
		// addition/removal lists must be built deterministically so that nodes
		// can agree on the same set of changes.
		//
		// ChannelIDs should always be sorted the same way (channel ID ascending).
		expectedChannelDefs := p.ChannelDefinitionCache.Definitions()
		if err := VerifyChannelDefinitions(expectedChannelDefs); err != nil {
			return nil, fmt.Errorf("ChannelDefinitionCache.Definitions is invalid: %w", err)
		}

		removeChannelDefinitions := subtractChannelDefinitions(previousOutcome.ChannelDefinitions, expectedChannelDefs, MaxObservationRemoveChannelIDsLength)
		for channelID := range removeChannelDefinitions {
			removeChannelIDs[channelID] = struct{}{}
		}

		// TODO: needs testing
		validAfterSecondsChannelIDs := maps.Keys(previousOutcome.ValidAfterSeconds)
		// Sort so we cut off deterministically
		sortChannelIDs(validAfterSecondsChannelIDs)
		for _, channelID := range validAfterSecondsChannelIDs {
			if len(removeChannelIDs) >= MaxObservationRemoveChannelIDsLength {
				break
			}
			if _, ok := expectedChannelDefs[channelID]; !ok {
				removeChannelIDs[channelID] = struct{}{}
			}
		}

		// NOTE: This is slow because it deeply compares every value in the map.
		// To improve performance, consider changing channel voting to happen
		// every N rounds instead of every round. Or, alternatively perhaps the
		// first e.g. 100 rounds could check every round to allow for fast feed
		// spinup, then after that every 10 or 100 rounds.
		updateChannelDefinitions = make(llotypes.ChannelDefinitions)
		expectedChannelIDs := maps.Keys(expectedChannelDefs)
		// Sort so we cut off deterministically
		sortChannelIDs(expectedChannelIDs)
		for _, channelID := range expectedChannelIDs {
			prev, exists := previousOutcome.ChannelDefinitions[channelID]
			channelDefinition := expectedChannelDefs[channelID]
			if exists && prev.Equals(channelDefinition) {
				continue
			}
			// Add or replace channel
			updateChannelDefinitions[channelID] = channelDefinition
			if len(updateChannelDefinitions) >= MaxObservationUpdateChannelDefinitionsLength {
				// Never add more than MaxObservationUpdateChannelDefinitionsLength
				break
			}
		}

		if len(updateChannelDefinitions) > 0 {
			p.Logger.Debugw("Voting to update channel definitions",
				"updateChannelDefinitions", updateChannelDefinitions,
				"seqNr", outctx.SeqNr,
				"stage", "Observation")
		}
		if len(removeChannelIDs) > 0 {
			p.Logger.Debugw("Voting to remove channel definitions",
				"removeChannelIDs", removeChannelIDs,
				"seqNr", outctx.SeqNr,
				"stage", "Observation",
			)
		}
	}

	var streamValues StreamValues
	if len(previousOutcome.ChannelDefinitions) == 0 {
		p.Logger.Debugw("ChannelDefinitions is empty, will not generate any observations", "stage", "Observation", "seqNr", outctx.SeqNr)
	} else {
		streamValues = make(StreamValues)
		for _, channelDefinition := range previousOutcome.ChannelDefinitions {
			for _, strm := range channelDefinition.Streams {
				streamValues[strm.StreamID] = nil
			}
		}

		if err := p.DataSource.Observe(ctx, streamValues, dsOpts{p.Config.VerboseLogging, outctx.SeqNr}); err != nil {
			return nil, fmt.Errorf("DataSource.Observe error: %w", err)
		}
	}

	var rawObservation []byte
	{
		var err error
		rawObservation, err = p.ObservationCodec.Encode(Observation{
			attestedRetirementReport,
			shouldRetire,
			nowNanoseconds,
			removeChannelIDs,
			updateChannelDefinitions,
			streamValues,
		})
		if err != nil {
			return nil, fmt.Errorf("Observation encode error: %w", err)
		}
	}

	return rawObservation, nil
}

// Should return an error if an observation isn't well-formed.
// Non-well-formed  observations will be discarded by the protocol. This is
// called for each observation, don't do anything slow in here.
//
// You may assume that the outctx.SeqNr is increasing monotonically (though
// *not* strictly) across the lifetime of a protocol instance and that
// outctx.previousOutcome contains the consensus outcome with sequence
// number (outctx.SeqNr-1).
func (p *Plugin) ValidateObservation(outctx ocr3types.OutcomeContext, query types.Query, ao types.AttributedObservation) error {
	if outctx.SeqNr < 1 {
		return fmt.Errorf("Invalid SeqNr: %d", outctx.SeqNr)
	} else if outctx.SeqNr == 1 {
		if len(ao.Observation) != 0 {
			return fmt.Errorf("Expected empty observation for first round, got: 0x%x", ao.Observation)
		}
	}

	observation, err := p.ObservationCodec.Decode(ao.Observation)
	if err != nil {
		// Critical error
		// If the previous outcome cannot be decoded for whatever reason, the
		// protocol will become permanently stuck at this point
		return fmt.Errorf("Observation decode error (got: 0x%x): %w", ao.Observation, err)
	}

	if p.PredecessorConfigDigest == nil && len(observation.AttestedPredecessorRetirement) != 0 {
		return fmt.Errorf("AttestedPredecessorRetirement is not empty even though this instance has no predecessor")
	}

	if len(observation.UpdateChannelDefinitions) > MaxObservationUpdateChannelDefinitionsLength {
		return fmt.Errorf("UpdateChannelDefinitions is too long: %v vs %v", len(observation.UpdateChannelDefinitions), MaxObservationUpdateChannelDefinitionsLength)
	}

	if len(observation.RemoveChannelIDs) > MaxObservationRemoveChannelIDsLength {
		return fmt.Errorf("RemoveChannelIDs is too long: %v vs %v", len(observation.RemoveChannelIDs), MaxObservationRemoveChannelIDsLength)
	}

	if err := VerifyChannelDefinitions(observation.UpdateChannelDefinitions); err != nil {
		return fmt.Errorf("UpdateChannelDefinitions is invalid: %w", err)
	}

	if len(observation.StreamValues) > MaxObservationStreamValuesLength {
		return fmt.Errorf("StreamValues is too long: %v vs %v", len(observation.StreamValues), MaxObservationStreamValuesLength)
	}

	return nil
}

type Outcome struct {
	// LifeCycleStage the protocol is in
	LifeCycleStage llotypes.LifeCycleStage
	// ObservationsTimestampNanoseconds is the median timestamp from the
	// latest set of observations
	ObservationsTimestampNanoseconds int64
	// ChannelDefinitions defines the set & structure of channels for which we
	// generate reports
	ChannelDefinitions llotypes.ChannelDefinitions
	// Latest ValidAfterSeconds value for each channel, reports for each channel
	// span from ValidAfterSeconds to ObservationTimestampSeconds
	ValidAfterSeconds map[llotypes.ChannelID]uint32
	// StreamAggregates contains stream IDs mapped to various aggregations.
	// Usually you will only have one aggregation type per stream but since
	// channels can define different aggregation methods, sometimes we will
	// need multiple.
	StreamAggregates StreamAggregates
}

// The Outcome's ObservationsTimestamp rounded down to seconds precision
func (out *Outcome) ObservationsTimestampSeconds() (uint32, error) {
	result := time.Unix(0, out.ObservationsTimestampNanoseconds).Unix()
	if int64(uint32(result)) != result {
		return 0, fmt.Errorf("timestamp doesn't fit into uint32: %v", result)
	}
	return uint32(result), nil
}

// Indicates whether a report can be generated for the given channel.
// Returns nil if channel is reportable
// TODO: Test this function
func (out *Outcome) IsReportable(channelID llotypes.ChannelID) *ErrUnreportableChannel {
	if out.LifeCycleStage == LifeCycleStageRetired {
		return &ErrUnreportableChannel{nil, "IsReportable=false; retired channel", channelID}
	}

	observationsTimestampSeconds, err := out.ObservationsTimestampSeconds()
	if err != nil {
		return &ErrUnreportableChannel{err, "IsReportable=false; invalid observations timestamp", channelID}
	}

	channelDefinition, exists := out.ChannelDefinitions[channelID]
	if !exists {
		return &ErrUnreportableChannel{nil, "IsReportable=false; no channel definition with this ID", channelID}
	}

	for _, strm := range channelDefinition.Streams {
		if out.StreamAggregates[strm.StreamID] == nil {
			// FIXME: Is this comment actually correct?
			// This can happen in normal operation, because in Report() we use
			// the ChannelDefinitions in the generated Outcome. But that was
			// compiled with Observations made using the ChannelDefinitions
			// from the PREVIOUS outcome. So if channel definitions have been
			// added in this round, we would not expect there to be
			// observations present for new streams in those channels.
			return &ErrUnreportableChannel{nil, fmt.Sprintf("IsReportable=false; median was nil for stream %d", strm.StreamID), channelID}
		}
	}

	if _, ok := out.ValidAfterSeconds[channelID]; !ok {
		// No validAfterSeconds entry yet, this must be a new channel.
		// validAfterSeconds will be populated in Outcome() so the channel
		// becomes reportable in later protocol rounds.
		// TODO: Test this case, haven't seen it in prod logs even though it would be expected
		return &ErrUnreportableChannel{nil, "IsReportable=false; no validAfterSeconds entry yet, this must be a new channel", channelID}
	}

	if validAfterSeconds := out.ValidAfterSeconds[channelID]; validAfterSeconds >= observationsTimestampSeconds {
		return &ErrUnreportableChannel{nil, fmt.Sprintf("IsReportable=false; not valid yet (observationsTimestampSeconds=%d < validAfterSeconds=%d)", observationsTimestampSeconds, validAfterSeconds), channelID}
	}

	return nil
}

type ErrUnreportableChannel struct {
	Inner     error
	Reason    string
	ChannelID llotypes.ChannelID
}

func (e *ErrUnreportableChannel) Error() string {
	s := fmt.Sprintf("ChannelID: %d; Reason: %s", e.ChannelID, e.Reason)
	if e.Inner != nil {
		s += fmt.Sprintf("; Err: %v", e.Inner)
	}
	return s
}

func (e *ErrUnreportableChannel) String() string {
	return e.Error()
}

func (e *ErrUnreportableChannel) Unwrap() error {
	return e.Inner
}

// List of reportable channels (according to IsReportable), sorted according
// to a canonical ordering
// TODO: test this
func (out *Outcome) ReportableChannels() (reportable []llotypes.ChannelID, unreportable []*ErrUnreportableChannel) {
	for channelID := range out.ChannelDefinitions {
		if err := out.IsReportable(channelID); err != nil {
			unreportable = append(unreportable, err)
		} else {
			reportable = append(reportable, channelID)
		}
	}

	sort.Slice(reportable, func(i, j int) bool {
		return reportable[i] < reportable[j]
	})

	return
}

// Generates an outcome for a seqNr, typically based on the previous
// outcome, the current query, and the current set of attributed
// observations.
//
// This function should be pure. Don't do anything slow in here.
//
// You may assume that the outctx.SeqNr is increasing monotonically (though
// *not* strictly) across the lifetime of a protocol instance and that
// outctx.previousOutcome contains the consensus outcome with sequence
// number (outctx.SeqNr-1).
//
// libocr guarantees that this will always be called with at least 2f+1
// AttributedObservations
func (p *Plugin) Outcome(outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (ocr3types.Outcome, error) {
	if len(aos) < 2*p.F+1 {
		return nil, fmt.Errorf("invariant violation: expected at least 2f+1 attributed observations, got %d (f: %d)", len(aos), p.F)
	}

	// Initial outcome is kind of a "keystone" with minimum extra information
	if outctx.SeqNr <= 1 {
		// Initial Outcome
		var lifeCycleStage llotypes.LifeCycleStage
		if p.PredecessorConfigDigest == nil {
			// Start straight in production if we have no predecessor
			lifeCycleStage = LifeCycleStageProduction
		} else {
			lifeCycleStage = LifeCycleStageStaging
		}
		outcome := Outcome{
			lifeCycleStage,
			0,
			nil,
			nil,
			nil,
		}
		return p.OutcomeCodec.Encode(outcome)
	}

	/////////////////////////////////
	// Decode previousOutcome
	/////////////////////////////////
	previousOutcome, err := p.OutcomeCodec.Decode(outctx.PreviousOutcome)
	if err != nil {
		return nil, fmt.Errorf("error decoding previous outcome: %v", err)
	}

	/////////////////////////////////
	// Decode observations
	/////////////////////////////////

	// a single valid retirement report is enough
	var validPredecessorRetirementReport *RetirementReport

	shouldRetireVotes := 0

	timestampsNanoseconds := []int64{}

	removeChannelVotesByID := map[llotypes.ChannelID]int{}

	// for each channelId count number of votes that mention it and count number of votes that include it.
	updateChannelVotesByHash := map[ChannelHash]int{}
	updateChannelDefinitionsByHash := map[ChannelHash]ChannelDefinitionWithID{}

	streamObservations := make(map[llotypes.StreamID][]StreamValue)

	for _, ao := range aos {
		// TODO: Put in a function
		// MERC-3524
		observation, err2 := p.ObservationCodec.Decode(ao.Observation)
		if err2 != nil {
			p.Logger.Warnw("ignoring invalid observation", "oracleID", ao.Observer, "error", err2)
			continue
		}

		if len(observation.AttestedPredecessorRetirement) != 0 && validPredecessorRetirementReport == nil {
			pcd := *p.PredecessorConfigDigest
			retirementReport, err3 := p.PredecessorRetirementReportCache.CheckAttestedRetirementReport(pcd, observation.AttestedPredecessorRetirement)
			if err3 != nil {
				p.Logger.Warnw("ignoring observation with invalid attested predecessor retirement", "oracleID", ao.Observer, "error", err3, "predecessorConfigDigest", pcd)
				continue
			}
			validPredecessorRetirementReport = &retirementReport
		}

		if observation.ShouldRetire {
			shouldRetireVotes++
		}

		timestampsNanoseconds = append(timestampsNanoseconds, observation.UnixTimestampNanoseconds)

		for channelID := range observation.RemoveChannelIDs {
			removeChannelVotesByID[channelID]++
		}

		for channelID, channelDefinition := range observation.UpdateChannelDefinitions {
			defWithID := ChannelDefinitionWithID{channelDefinition, channelID}
			channelHash := MakeChannelHash(defWithID)
			updateChannelVotesByHash[channelHash]++
			updateChannelDefinitionsByHash[channelHash] = defWithID
		}

		var missingObservations []llotypes.StreamID
		for id, sv := range observation.StreamValues {
			if sv != nil { // FIXME: nil checks don't work here. Test this and figure out what to do (also, are there other cases?)
				streamObservations[id] = append(streamObservations[id], sv)
			} else {
				missingObservations = append(missingObservations, id)
			}
		}
		if p.Config.VerboseLogging {
			if len(missingObservations) > 0 {
				sort.Slice(missingObservations, func(i, j int) bool { return missingObservations[i] < missingObservations[j] })
				p.Logger.Debugw("Peer was missing observations", "streamIDs", missingObservations, "oracleID", ao.Observer, "stage", "Outcome", "seqNr", outctx.SeqNr)
			}
			p.Logger.Debugw("Got observations from peer", "stage", "Outcome", "sv", streamObservations, "oracleID", ao.Observer, "seqNr", outctx.SeqNr)
		}
	}

	if len(timestampsNanoseconds) == 0 {
		return nil, errors.New("no valid observations")
	}

	var outcome Outcome

	/////////////////////////////////
	// outcome.LifeCycleStage
	/////////////////////////////////
	if previousOutcome.LifeCycleStage == LifeCycleStageStaging && validPredecessorRetirementReport != nil {
		// Promote this protocol instance to the production stage! 🚀

		// override ValidAfterSeconds with the value from the retirement report
		// so that we have no gaps in the validity time range.
		outcome.ValidAfterSeconds = validPredecessorRetirementReport.ValidAfterSeconds
		outcome.LifeCycleStage = LifeCycleStageProduction
	} else {
		outcome.LifeCycleStage = previousOutcome.LifeCycleStage
	}

	if outcome.LifeCycleStage == LifeCycleStageProduction && shouldRetireVotes > p.F {
		outcome.LifeCycleStage = LifeCycleStageRetired
	}

	/////////////////////////////////
	// outcome.ObservationsTimestampNanoseconds
	// TODO: Refactor this into an aggregate function
	// MERC-3524
	sort.Slice(timestampsNanoseconds, func(i, j int) bool { return timestampsNanoseconds[i] < timestampsNanoseconds[j] })
	outcome.ObservationsTimestampNanoseconds = timestampsNanoseconds[len(timestampsNanoseconds)/2]

	/////////////////////////////////
	// outcome.ChannelDefinitions
	/////////////////////////////////
	outcome.ChannelDefinitions = previousOutcome.ChannelDefinitions
	if outcome.ChannelDefinitions == nil {
		outcome.ChannelDefinitions = llotypes.ChannelDefinitions{}
	}

	// if retired, stop updating channel definitions
	if outcome.LifeCycleStage == LifeCycleStageRetired {
		removeChannelVotesByID, updateChannelDefinitionsByHash = nil, nil
	}

	var removedChannelIDs []llotypes.ChannelID
	for channelID, voteCount := range removeChannelVotesByID {
		if voteCount <= p.F {
			continue
		}
		removedChannelIDs = append(removedChannelIDs, channelID)
		delete(outcome.ChannelDefinitions, channelID)
	}

	type hashWithID struct {
		ChannelHash
		ChannelDefinitionWithID
	}
	orderedHashes := make([]hashWithID, 0, len(updateChannelDefinitionsByHash))
	for channelHash, dfnWithID := range updateChannelDefinitionsByHash {
		orderedHashes = append(orderedHashes, hashWithID{channelHash, dfnWithID})
	}
	// Use predictable order for adding channels (id asc) so that extras that
	// exceed the max are consistent across all nodes
	sort.Slice(orderedHashes, func(i, j int) bool { return orderedHashes[i].ChannelID < orderedHashes[j].ChannelID })
	for _, hwid := range orderedHashes {
		voteCount := updateChannelVotesByHash[hwid.ChannelHash]
		if voteCount <= p.F {
			continue
		}
		defWithID := hwid.ChannelDefinitionWithID
		if original, exists := outcome.ChannelDefinitions[defWithID.ChannelID]; exists {
			p.Logger.Debugw("Adding channel (replacement)",
				"channelID", defWithID.ChannelID,
				"originalChannelDefinition", original,
				"replaceChannelDefinition", defWithID,
				"seqNr", outctx.SeqNr,
				"stage", "Outcome",
			)
			outcome.ChannelDefinitions[defWithID.ChannelID] = defWithID.ChannelDefinition
		} else if len(outcome.ChannelDefinitions) >= MaxOutcomeChannelDefinitionsLength {
			p.Logger.Warnw("Adding channel FAILED. Cannot add channel, outcome already contains maximum number of channels",
				"maxOutcomeChannelDefinitionsLength", MaxOutcomeChannelDefinitionsLength,
				"addChannelDefinition", defWithID,
				"seqNr", outctx.SeqNr,
				"stage", "Outcome",
			)
			// continue, don't break here because remaining channels might be a
			// replacement rather than an addition, and this is still ok
			continue
		}
		p.Logger.Debugw("Adding channel (new)",
			"channelID", defWithID.ChannelID,
			"addChannelDefinition", defWithID,
			"seqNr", outctx.SeqNr,
			"stage", "Outcome",
		)
		outcome.ChannelDefinitions[defWithID.ChannelID] = defWithID.ChannelDefinition
	}

	/////////////////////////////////
	// outcome.ValidAfterSeconds
	/////////////////////////////////

	// ValidAfterSeconds can be non-nil here if earlier code already
	// populated ValidAfterSeconds during promotion to production. In this
	// case, nothing to do.
	if outcome.ValidAfterSeconds == nil {
		previousObservationsTimestampSeconds, err2 := previousOutcome.ObservationsTimestampSeconds()
		if err2 != nil {
			return nil, fmt.Errorf("error getting previous outcome's observations timestamp: %v", err2)
		}

		outcome.ValidAfterSeconds = map[llotypes.ChannelID]uint32{}
		for channelID, previousValidAfterSeconds := range previousOutcome.ValidAfterSeconds {
			if err3 := previousOutcome.IsReportable(channelID); err3 != nil {
				if p.Config.VerboseLogging {
					p.Logger.Debugw("Channel is not reportable", "channelID", channelID, "err", err3, "stage", "Outcome", "seqNr", outctx.SeqNr)
				}
				// was reported based on previous outcome
				outcome.ValidAfterSeconds[channelID] = previousObservationsTimestampSeconds
			} else {
				// was skipped based on previous outcome
				outcome.ValidAfterSeconds[channelID] = previousValidAfterSeconds
			}
		}
	}

	observationsTimestampSeconds, err := outcome.ObservationsTimestampSeconds()
	if err != nil {
		return nil, fmt.Errorf("error getting outcome's observations timestamp: %w", err)
	}

	for channelID := range outcome.ChannelDefinitions {
		if _, ok := outcome.ValidAfterSeconds[channelID]; !ok {
			// new channel, set validAfterSeconds to observations timestamp
			outcome.ValidAfterSeconds[channelID] = observationsTimestampSeconds
		}
	}

	// One might think that we should simply delete any channel from
	// ValidAfterSeconds that is not mentioned in the ChannelDefinitions. This
	// could, however, lead to gaps being created if this protocol instance is
	// promoted from staging to production while we're still "ramping up" the
	// full set of channels. We do the "safe" thing (i.e. minimizing occurrence
	// of gaps) here and only remove channels if there has been an explicit vote
	// to remove them.
	for _, channelID := range removedChannelIDs {
		delete(outcome.ValidAfterSeconds, channelID)
	}

	/////////////////////////////////
	// outcome.StreamAggregates
	/////////////////////////////////
	outcome.StreamAggregates = make(map[llotypes.StreamID]map[llotypes.Aggregator]StreamValue, len(streamObservations))
	// Aggregation methods are defined on a per-channel basis, but we only want
	// to do the minimum necessary number of aggregations (one per stream/aggregator
	// pair) and re-use the same result, in case multiple channels share the
	// same stream/aggregator pair.
	for cid, cd := range outcome.ChannelDefinitions {
		for _, strm := range cd.Streams {
			sid, agg := strm.StreamID, strm.Aggregator
			if _, exists := outcome.StreamAggregates[sid][agg]; exists {
				// Should only happen in the case of duplicate streams, no
				// need to aggregate twice
				continue
			}
			aggF := GetAggregatorFunc(agg)
			if aggF == nil {
				return nil, fmt.Errorf("no aggregator function defined for aggregator of type %v", agg)
			}
			m, exists := outcome.StreamAggregates[sid]
			if !exists {
				m = make(map[llotypes.Aggregator]StreamValue)
				outcome.StreamAggregates[sid] = m
			}
			result, err := aggF(streamObservations[sid], p.F)
			if err != nil {
				if p.Config.VerboseLogging {
					p.Logger.Warnw("Aggregation failed", "aggregator", agg, "channelID", cid, "f", p.F, "streamID", sid, "observations", streamObservations[sid], "stage", "Outcome", "seqNr", outctx.SeqNr, "err", err)
				}
				// FIXME: Is this a complete failure?
				// MERC-3524
				continue
			}
			m[agg] = result
		}
	}

	if p.Config.VerboseLogging {
		p.Logger.Debugw("Generated outcome", "outcome", outcome, "stage", "Outcome", "seqNr", outctx.SeqNr)
	}
	return p.OutcomeCodec.Encode(outcome)
}

type Report struct {
	ConfigDigest types.ConfigDigest
	// OCR sequence number of this report
	SeqNr uint64
	// Channel that is being reported on
	ChannelID llotypes.ChannelID
	// Report is only valid at t > ValidAfterSeconds
	ValidAfterSeconds uint32
	// ObservationTimestampSeconds is the median of all observation timestamps
	// (note that this timestamp is taken immediately before we initiate any
	// observations)
	ObservationTimestampSeconds uint32
	// Values for every stream in the channel
	Values []StreamValue
	// The contract onchain will only validate non-specimen reports. A staging
	// protocol instance will generate specimen reports so we can validate it
	// works properly without any risk of misreports landing on chain.
	Specimen bool
}

func (p *Plugin) encodeReport(r Report, cd llotypes.ChannelDefinition) (types.Report, error) {
	codec, exists := p.Codecs[cd.ReportFormat]
	if !exists {
		return nil, fmt.Errorf("codec missing for ReportFormat=%q", cd.ReportFormat)
	}
	return codec.Encode(r, cd)
}

// Generates a (possibly empty) list of reports from an outcome. Each report
// will be signed and possibly be transmitted to the contract. (Depending on
// ShouldAcceptAttestedReport & ShouldTransmitAcceptedReport)
//
// This function should be pure. Don't do anything slow in here.
//
// This is likely to change in the future. It will likely be returning a
// list of report batches, where each batch goes into its own Merkle tree.
//
// You may assume that the outctx.SeqNr is increasing monotonically (though
// *not* strictly) across the lifetime of a protocol instance and that
// outctx.previousOutcome contains the consensus outcome with sequence
// number (outctx.SeqNr-1).
func (p *Plugin) Reports(seqNr uint64, rawOutcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[llotypes.ReportInfo], error) {
	if seqNr <= 1 {
		// no reports for initial round
		return nil, nil
	}

	outcome, err := p.OutcomeCodec.Decode(rawOutcome)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling outcome: %w", err)
	}

	observationsTimestampSeconds, err := outcome.ObservationsTimestampSeconds()
	if err != nil {
		return nil, fmt.Errorf("error getting observations timestamp: %w", err)
	}

	rwis := []ocr3types.ReportWithInfo[llotypes.ReportInfo]{}

	if outcome.LifeCycleStage == LifeCycleStageRetired {
		// if we're retired, emit special retirement report to transfer
		// ValidAfterSeconds part of state to the new protocol instance for a
		// "gapless" handover
		retirementReport := RetirementReport{
			outcome.ValidAfterSeconds,
		}

		rwis = append(rwis, ocr3types.ReportWithInfo[llotypes.ReportInfo]{
			// TODO: Needs retirement report codec
			Report: must(json.Marshal(retirementReport)),
			Info: llotypes.ReportInfo{
				LifeCycleStage: outcome.LifeCycleStage,
				ReportFormat:   llotypes.ReportFormatJSON,
			},
		})
	}

	reportableChannels, unreportableChannels := outcome.ReportableChannels()
	if p.Config.VerboseLogging {
		p.Logger.Debugw("Reportable channels", "reportableChannels", reportableChannels, "unreportableChannels", unreportableChannels, "stage", "Report", "seqNr", seqNr)
	}

	for _, cid := range reportableChannels {
		cd := outcome.ChannelDefinitions[cid]
		values := make([]StreamValue, 0, len(cd.Streams))
		for _, strm := range cd.Streams {
			// TODO: Can you ever get nil values (i.e. missing from the
			// StreamAggregates) here? What happens if you do?
			// MERC-3524
			values = append(values, outcome.StreamAggregates[strm.StreamID][strm.Aggregator])
		}

		report := Report{
			p.ConfigDigest,
			seqNr,
			cid,
			outcome.ValidAfterSeconds[cid],
			observationsTimestampSeconds,
			values,
			outcome.LifeCycleStage != LifeCycleStageProduction,
		}

		encoded, err := p.encodeReport(report, cd)
		if err != nil {
			return nil, err
		}
		rwis = append(rwis, ocr3types.ReportWithInfo[llotypes.ReportInfo]{
			Report: encoded,
			Info: llotypes.ReportInfo{
				LifeCycleStage: outcome.LifeCycleStage,
				ReportFormat:   cd.ReportFormat,
			},
		})
	}

	if p.Config.VerboseLogging && len(rwis) == 0 {
		p.Logger.Debugw("No reports, will not transmit anything", "reportableChannels", reportableChannels, "stage", "Report", "seqNr", seqNr)
	}

	return rwis, nil
}

func (p *Plugin) ShouldAcceptAttestedReport(context.Context, uint64, ocr3types.ReportWithInfo[llotypes.ReportInfo]) (bool, error) {
	// Transmit it all to the Mercury server
	return true, nil
}

func (p *Plugin) ShouldTransmitAcceptedReport(context.Context, uint64, ocr3types.ReportWithInfo[llotypes.ReportInfo]) (bool, error) {
	// Transmit it all to the Mercury server
	return true, nil
}

// ObservationQuorum returns the minimum number of valid (according to
// ValidateObservation) observations needed to construct an outcome.
//
// This function should be pure. Don't do anything slow in here.
//
// This is an advanced feature. The "default" approach (what OCR1 & OCR2
// did) is to have an empty ValidateObservation function and return
// QuorumTwoFPlusOne from this function.
func (p *Plugin) ObservationQuorum(outctx ocr3types.OutcomeContext, query types.Query) (ocr3types.Quorum, error) {
	return ocr3types.QuorumTwoFPlusOne, nil
}

func (p *Plugin) Close() error {
	return nil
}

func subtractChannelDefinitions(minuend llotypes.ChannelDefinitions, subtrahend llotypes.ChannelDefinitions, limit int) llotypes.ChannelDefinitions {
	differenceList := []ChannelDefinitionWithID{}
	for channelID, channelDefinition := range minuend {
		if _, ok := subtrahend[channelID]; !ok {
			differenceList = append(differenceList, ChannelDefinitionWithID{channelDefinition, channelID})
		}
	}

	// Sort so we return deterministic result
	sort.Slice(differenceList, func(i, j int) bool {
		return differenceList[i].ChannelID < differenceList[j].ChannelID
	})

	if len(differenceList) > limit {
		differenceList = differenceList[:limit]
	}

	difference := llotypes.ChannelDefinitions{}
	for _, defWithID := range differenceList {
		difference[defWithID.ChannelID] = defWithID.ChannelDefinition
	}

	return difference
}

// deterministic sort of channel IDs
func sortChannelIDs(cids []llotypes.ChannelID) {
	sort.Slice(cids, func(i, j int) bool {
		return cids[i] < cids[j]
	})
}
