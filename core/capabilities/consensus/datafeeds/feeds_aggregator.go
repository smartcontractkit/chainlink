package datafeeds

import (
	"math"
	"time"

	ocrcommon "github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/consensus/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type DataFeedsAggregatorConfig struct {
	PricePoints map[string]PricePointConfig
}

type PricePointConfig struct {
	Deviation float64
	Heartbeat time.Duration
}

type dataFeedsAggregator struct {
	config       DataFeedsAggregatorConfig
	mercuryCodec mercury.MercuryCodec
	lggr         logger.Logger
}

// EncodableOutcome is a list of AggregatedPricePoints
// Metadata is a map of feedID -> (timestamp, price) representing onchain state (see DataFeedsOutcomeMetadata proto)
func (a *dataFeedsAggregator) Aggregate(previousOutcome *types.AggregationOutcome, observations map[ocrcommon.OracleID][]values.Value) (*types.AggregationOutcome, error) {
	// find latest valid Mercury report for each feed ID
	latestReportPerFeed := make(map[string]mercury.MercuryReport)
	for nodeId, nodeObservations := range observations {
		// we only expect a single observation per node - new Mercury data
		if len(nodeObservations) == 0 {
			a.lggr.Warnf("node %d contributed with empty observations", nodeId)
			continue
		}
		if len(nodeObservations) > 1 {
			a.lggr.Warnf("node %d contributed with more than one observation", nodeId)
		}
		mercuryReportSet, err := a.mercuryCodec.Unwrap(nodeObservations[0])
		if err != nil {
			a.lggr.Errorf("node %d contributed with invalid Mercury reports: %v", nodeId, err)
			continue
		}
		for feedID, report := range mercuryReportSet.Reports {
			latest, ok := latestReportPerFeed[feedID]
			if !ok || report.Info.Timestamp > latest.Info.Timestamp {
				latestReportPerFeed[feedID] = report
			}
		}
	}

	// TODO: handle empty previousOutcome

	onchainState := &DataFeedsOutcomeMetadata{}
	err := proto.Unmarshal(previousOutcome.Metadata, onchainState)
	if err != nil {
		return nil, err
	}

	needUpdate := []string{}
	for feedId, onchainFeedInfo := range onchainState.FeedInfo {
		latestReport, ok := latestReportPerFeed[feedId]
		if !ok {
			a.lggr.Errorf("no new Mercury report for feed: %v", feedId)
			continue
		}
		config := a.config.PricePoints[feedId]
		if latestReport.Info.Timestamp-onchainFeedInfo.Timestamp > uint32(config.Heartbeat) ||
			deviation(onchainFeedInfo.Price, latestReport.Info.Price) > config.Deviation {
			onchainFeedInfo.Timestamp = latestReport.Info.Timestamp
			onchainFeedInfo.Price = latestReport.Info.Price
			needUpdate = append(needUpdate, feedId)
		}
	}

	marshalledOnchainState, err := proto.Marshal(onchainState)
	if err != nil {
		return nil, err
	}

	return &types.AggregationOutcome{
		// TODO: set EncodableOutcome
		Metadata:     marshalledOnchainState,
		ShouldReport: len(needUpdate) > 0,
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

func NewDataFeedsAggregator(config values.Map, mercuryCodec mercury.MercuryCodec, lggr logger.Logger) (types.Aggregator, error) {
	parsedConfig, err := parseConfig(config)
	if err != nil {
		return nil, err
	}
	return &dataFeedsAggregator{
		config:       parsedConfig,
		mercuryCodec: mercuryCodec,
		lggr:         lggr,
	}, nil
}

func parseConfig(config values.Map) (DataFeedsAggregatorConfig, error) {
	// TODO: implement parsing; remap feed names into hex-encoded feed IDs
	return DataFeedsAggregatorConfig{
		PricePoints: map[string]PricePointConfig{
			"ETHUSD": {
				Deviation: 0.05,
				Heartbeat: 5 * time.Minute,
			},
		},
	}, nil
}
