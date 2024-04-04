package datafeeds

import (
	"fmt"
	"math"

	"github.com/shopspring/decimal"
	ocrcommon "github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const OutputFieldName = "mercury_reports"

type aggregatorConfig struct {
	Feeds map[mercury.FeedID]feedConfig
}

type feedConfig struct {
	Deviation       decimal.Decimal `mapstructure:"-"`
	Heartbeat       int
	DeviationString string `mapstructure:"deviation"`
}

//go:generate mockery --quiet --name MercuryCodec --output ./mocks/ --case=underscore
type MercuryCodec interface {
	// validate each report and convert to ReportSet struct
	Unwrap(raw values.Value) (mercury.ReportSet, error)

	// validate each report and convert to Value
	Wrap(reportSet mercury.ReportSet) (values.Value, error)
}

type dataFeedsAggregator struct {
	config       aggregatorConfig
	mercuryCodec MercuryCodec
	lggr         logger.Logger
}

var _ types.Aggregator = (*dataFeedsAggregator)(nil)

// EncodableOutcome is a list of AggregatedPricePoints
// Metadata is a map of feedID -> (timestamp, price) representing onchain state (see DataFeedsOutcomeMetadata proto)
func (a *dataFeedsAggregator) Aggregate(previousOutcome *types.AggregationOutcome, observations map[ocrcommon.OracleID][]values.Value) (*types.AggregationOutcome, error) {
	// find latest valid Mercury report for each feed ID
	latestReportPerFeed := make(map[mercury.FeedID]mercury.Report)
	for nodeID, nodeObservations := range observations {
		// we only expect a single observation per node - new Mercury data
		if len(nodeObservations) == 0 {
			a.lggr.Warnf("node %d contributed with empty observations", nodeID)
			continue
		}
		if len(nodeObservations) > 1 {
			a.lggr.Warnf("node %d contributed with more than one observation", nodeID)
		}
		mercuryReportSet, err := a.mercuryCodec.Unwrap(nodeObservations[0])
		if err != nil {
			a.lggr.Errorf("node %d contributed with invalid Mercury reports: %v", nodeID, err)
			continue
		}
		for feedID, report := range mercuryReportSet.Reports {
			latest, ok := latestReportPerFeed[feedID]
			if !ok || report.Info.Timestamp > latest.Info.Timestamp {
				latestReportPerFeed[feedID] = report
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
	} else {
		a.lggr.Debug("empty previous outcome - initializing empty onchain state")
		currentState.FeedInfo = make(map[string]*DataFeedsMercuryReportInfo)
		for feedID := range a.config.Feeds {
			currentState.FeedInfo[feedID.String()] = &DataFeedsMercuryReportInfo{
				Timestamp: 0, // will always trigger an update
				Price:     0.0,
			}
		}
	}

	reportsNeedingUpdate := []any{} // [][]byte
	for feedID, previousReportInfo := range currentState.FeedInfo {
		feedID := mercury.FeedID(feedID)
		err := feedID.Validate()
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
		if latestReport.Info.Timestamp-previousReportInfo.Timestamp > uint32(config.Heartbeat) ||
			deviation(previousReportInfo.Price, latestReport.Info.Price) > config.Deviation.InexactFloat64() {
			previousReportInfo.Timestamp = latestReport.Info.Timestamp
			previousReportInfo.Price = latestReport.Info.Price
			reportsNeedingUpdate = append(reportsNeedingUpdate, latestReport.FullReport)
		}
	}

	marshalledState, err := proto.Marshal(currentState)
	if err != nil {
		return nil, err
	}

	wrappedReportsNeedingUpdates, err := values.NewMap(map[string]any{OutputFieldName: reportsNeedingUpdate})
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

func deviation(old, new float64) float64 {
	diff := math.Abs(new - old)
	if old == 0.0 {
		if diff == 0.0 {
			return 0.0
		}
		return math.MaxFloat64
	}
	return diff / old
}

func NewDataFeedsAggregator(config values.Map, mercuryCodec MercuryCodec, lggr logger.Logger) (types.Aggregator, error) {
	parsedConfig, err := ParseConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config (%+v): %w", config, err)
	}
	return &dataFeedsAggregator{
		config:       parsedConfig,
		mercuryCodec: mercuryCodec,
		lggr:         logger.Named(lggr, "DataFeedsAggregator"),
	}, nil
}

func ParseConfig(config values.Map) (aggregatorConfig, error) {
	parsedConfig := aggregatorConfig{
		Feeds: make(map[mercury.FeedID]feedConfig),
	}
	for feedIDStr, feedCfg := range config.Underlying {
		feedID := mercury.FeedID(feedIDStr)
		err := feedID.Validate()
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
