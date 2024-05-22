package datafeeds

import (
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"

	ocrcommon "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const (
	// Aggregator outputs reports in the following format:
	//   []Reports{FeedID []byte, RawReport []byte, Price *big.Int, Timestamp int64}
	// Example of a compatible EVM encoder ABI config:
	//   (bytes32 FeedID, bytes RawReport, uint256 Price, uint64 Timestamp)[] Reports
	TopLevelListOutputFieldName = "Reports"
	FeedIDOutputFieldName       = "FeedID"
	RawReportOutputFieldName    = "RawReport"
	PriceOutputFieldName        = "Price"
	TimestampOutputFieldName    = "Timestamp"

	addrLen = 20
)

type aggregatorConfig struct {
	Feeds map[datastreams.FeedID]feedConfig
}

type feedConfig struct {
	Deviation       decimal.Decimal `mapstructure:"-"`
	Heartbeat       int
	DeviationString string `mapstructure:"deviation"`
}

type dataFeedsAggregator struct {
	config      aggregatorConfig
	reportCodec datastreams.ReportCodec
	lggr        logger.Logger
}

var _ types.Aggregator = (*dataFeedsAggregator)(nil)

// This Aggregator has two phases:
//  1. Agree on valid trigger signers by extracting them from event metadata and aggregating using MODE (at least F+1 copies needed).
//  2. For each FeedID, select latest valid report, using signers list obtained in phase 1.
//
// EncodableOutcome is a list of aggregated price points.
// Metadata is a map of feedID -> (timestamp, price) representing onchain state (see DataFeedsOutcomeMetadata proto)
func (a *dataFeedsAggregator) Aggregate(previousOutcome *types.AggregationOutcome, observations map[ocrcommon.OracleID][]values.Value, f int) (*types.AggregationOutcome, error) {
	allowedSigners, minRequiredSignatures, payloads := a.extractSignersAndPayloads(observations, f)
	if len(payloads) > 0 && minRequiredSignatures == 0 {
		return nil, fmt.Errorf("cannot process non-empty observation payloads with minRequiredSignatures set to 0")
	}
	a.lggr.Debugw("extracted signers", "nAllowedSigners", len(allowedSigners), "minRequired", minRequiredSignatures, "nPayloads", len(payloads))
	// find latest valid report for each feed ID
	latestReportPerFeed := make(map[datastreams.FeedID]datastreams.FeedReport)
	for nodeID, payload := range payloads {
		mercuryReports, err := a.reportCodec.UnwrapValid(payload, allowedSigners, minRequiredSignatures)
		if err != nil {
			a.lggr.Errorf("node %d contributed with invalid reports: %v", nodeID, err)
			continue
		}
		for _, report := range mercuryReports {
			latest, ok := latestReportPerFeed[datastreams.FeedID(report.FeedID)]
			if !ok || report.ObservationTimestamp > latest.ObservationTimestamp {
				latestReportPerFeed[datastreams.FeedID(report.FeedID)] = report
			}
		}
	}
	a.lggr.Debugw("collected latestReportPerFeed", "len", len(latestReportPerFeed))

	currentState := &DataFeedsOutcomeMetadata{}
	if previousOutcome != nil {
		err := proto.Unmarshal(previousOutcome.Metadata, currentState)
		if err != nil {
			return nil, err
		}
	}
	// initialize empty state for missing feeds
	if currentState.FeedInfo == nil {
		currentState.FeedInfo = make(map[string]*DataFeedsMercuryReportInfo)
	}
	for feedID := range a.config.Feeds {
		if _, ok := currentState.FeedInfo[feedID.String()]; !ok {
			currentState.FeedInfo[feedID.String()] = &DataFeedsMercuryReportInfo{
				ObservationTimestamp: 0, // will always trigger an update
				BenchmarkPrice:       big.NewInt(0).Bytes(),
			}
			a.lggr.Debugw("initializing empty onchain state for feed", "feedID", feedID.String())
		}
	}
	// remove obsolete feeds from state
	for feedID := range currentState.FeedInfo {
		if _, ok := a.config.Feeds[datastreams.FeedID(feedID)]; !ok {
			delete(currentState.FeedInfo, feedID)
		}
		a.lggr.Debugw("removed obsolete feedID from state", "feedID", feedID)
	}

	reportsNeedingUpdate := []datastreams.FeedReport{}
	allIDs := []string{}
	for feedID := range currentState.FeedInfo {
		allIDs = append(allIDs, feedID)
	}
	// ensure deterministic order of reportsNeedingUpdate
	sort.Slice(allIDs, func(i, j int) bool { return allIDs[i] < allIDs[j] })
	for _, feedID := range allIDs {
		previousReportInfo := currentState.FeedInfo[feedID]
		feedID, err := datastreams.NewFeedID(feedID)
		if err != nil {
			a.lggr.Errorf("could not convert %s to feedID", feedID)
			continue
		}
		latestReport, ok := latestReportPerFeed[feedID]
		if !ok {
			a.lggr.Errorf("no new Mercury report for feed: %v", feedID)
			continue
		}
		config := a.config.Feeds[feedID]
		oldPrice := big.NewInt(0).SetBytes(previousReportInfo.BenchmarkPrice)
		newPrice := big.NewInt(0).SetBytes(latestReport.BenchmarkPrice)
		if latestReport.ObservationTimestamp-previousReportInfo.ObservationTimestamp > int64(config.Heartbeat) ||
			deviation(oldPrice, newPrice) > config.Deviation.InexactFloat64() {
			previousReportInfo.ObservationTimestamp = latestReport.ObservationTimestamp
			previousReportInfo.BenchmarkPrice = latestReport.BenchmarkPrice
			reportsNeedingUpdate = append(reportsNeedingUpdate, latestReport)
		}
	}

	marshalledState, err := proto.MarshalOptions{Deterministic: true}.Marshal(currentState)
	if err != nil {
		return nil, err
	}

	toWrap := []any{}
	for _, report := range reportsNeedingUpdate {
		feedID := datastreams.FeedID(report.FeedID).Bytes()
		toWrap = append(toWrap,
			map[string]any{
				FeedIDOutputFieldName:    feedID[:],
				RawReportOutputFieldName: report.FullReport,
				PriceOutputFieldName:     big.NewInt(0).SetBytes(report.BenchmarkPrice),
				TimestampOutputFieldName: report.ObservationTimestamp,
			})
	}

	wrappedReportsNeedingUpdates, err := values.NewMap(map[string]any{
		TopLevelListOutputFieldName: toWrap,
	})
	if err != nil {
		return nil, err
	}
	reportsProto := values.Proto(wrappedReportsNeedingUpdates)

	a.lggr.Debugw("Aggregate complete", "nReportsNeedingUpdate", len(reportsNeedingUpdate))
	return &types.AggregationOutcome{
		EncodableOutcome: reportsProto.GetMapValue(),
		Metadata:         marshalledState,
		ShouldReport:     len(reportsNeedingUpdate) > 0,
	}, nil
}

func (a *dataFeedsAggregator) extractSignersAndPayloads(observations map[ocrcommon.OracleID][]values.Value, fConsensus int) ([][]byte, int, map[ocrcommon.OracleID]values.Value) {
	payloads := make(map[ocrcommon.OracleID]values.Value)
	signers := make(map[[addrLen]byte]int)
	mins := make(map[int]int)
	for nodeID, nodeObservations := range observations {
		// we only expect a single observation per node - a Streams trigger event
		if len(nodeObservations) == 0 || nodeObservations[0] == nil {
			a.lggr.Warnf("node %d contributed with empty observations", nodeID)
			continue
		}
		if len(nodeObservations) > 1 {
			a.lggr.Warnf("node %d contributed with more than one observation", nodeID)
			continue
		}
		triggerEvent := &capabilities.TriggerEvent{}
		if err := nodeObservations[0].UnwrapTo(triggerEvent); err != nil {
			a.lggr.Warnf("could not parse observations from node %d: %v", nodeID, err)
			continue
		}
		meta := &datastreams.SignersMetadata{}
		if err := triggerEvent.Metadata.UnwrapTo(meta); err != nil {
			a.lggr.Warnf("could not parse trigger metadata from node %d: %v", nodeID, err)
			continue
		}
		currentNodeSigners, err := extractUniqueSigners(meta.Signers)
		if err != nil {
			a.lggr.Warnf("could not extract signers from node %d: %v", nodeID, err)
			continue
		}
		for signer := range currentNodeSigners {
			signers[signer]++
		}
		mins[meta.MinRequiredSignatures]++
		payloads[nodeID] = triggerEvent.Payload
	}
	// Agree on signers list and min-required. It's technically possible to have F+1 valid values from one trigger DON and F+1 from another trigger DON.
	// In that case both values are legitimate and signers list will contain nodes from both DONs. However, min-required value will be the higher one (if different).
	allowedSigners := [][]byte{}
	for signer, count := range signers {
		signer := signer
		if count >= fConsensus+1 {
			allowedSigners = append(allowedSigners, signer[:])
		}
	}
	minRequired := 0
	for minCandidate, count := range mins {
		if count >= fConsensus+1 && minCandidate > minRequired {
			minRequired = minCandidate
		}
	}
	return allowedSigners, minRequired, payloads
}

func extractUniqueSigners(signers [][]byte) (map[[addrLen]byte]struct{}, error) {
	uniqueSigners := make(map[[addrLen]byte]struct{})
	for _, signer := range signers {
		if len(signer) != addrLen {
			return nil, fmt.Errorf("invalid signer length: %d", len(signer))
		}
		var signerBytes [addrLen]byte
		copy(signerBytes[:], signer)
		uniqueSigners[signerBytes] = struct{}{}
	}
	return uniqueSigners, nil
}

func deviation(oldPrice, newPrice *big.Int) float64 {
	diff := &big.Int{}
	diff.Sub(oldPrice, newPrice)
	diff.Abs(diff)
	if oldPrice.Cmp(big.NewInt(0)) == 0 {
		if diff.Cmp(big.NewInt(0)) == 0 {
			return 0.0
		}
		return math.MaxFloat64
	}
	diffFl, _ := diff.Float64()
	oldFl, _ := oldPrice.Float64()
	return diffFl / oldFl
}

func NewDataFeedsAggregator(config values.Map, reportCodec datastreams.ReportCodec, lggr logger.Logger) (types.Aggregator, error) {
	parsedConfig, err := ParseConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config (%+v): %w", config, err)
	}
	return &dataFeedsAggregator{
		config:      parsedConfig,
		reportCodec: reportCodec,
		lggr:        logger.Named(lggr, "DataFeedsAggregator"),
	}, nil
}

func ParseConfig(config values.Map) (aggregatorConfig, error) {
	parsedConfig := aggregatorConfig{
		Feeds: make(map[datastreams.FeedID]feedConfig),
	}
	for feedIDStr, feedCfg := range config.Underlying {
		feedID, err := datastreams.NewFeedID(feedIDStr)
		if err != nil {
			return aggregatorConfig{}, err
		}
		var parsedFeedConfig feedConfig
		err = feedCfg.UnwrapTo(&parsedFeedConfig)
		if err != nil {
			return aggregatorConfig{}, err
		}

		if parsedFeedConfig.DeviationString != "" {
			dec, err := decimal.NewFromString(parsedFeedConfig.DeviationString)
			if err != nil {
				return aggregatorConfig{}, err
			}

			parsedFeedConfig.Deviation = dec
		}
		parsedConfig.Feeds[feedID] = parsedFeedConfig
	}
	return parsedConfig, nil
}
