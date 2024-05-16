package triggers

import (
	"errors"
	"sort"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type mercuryRemoteAggregator struct {
	codec datastreams.ReportCodec
	lggr  logger.Logger
}

// This aggregator is used by TriggerSubscriber to aggregate trigger events from multiple remote nodes.
// NOTE: Once Mercury supports parallel composition (and thus guarantee identical sets of reports),
// this will be replaced by the default MODE aggregator.
func NewMercuryRemoteAggregator(codec datastreams.ReportCodec, lggr logger.Logger) *mercuryRemoteAggregator {
	return &mercuryRemoteAggregator{
		codec: codec,
		lggr:  lggr,
	}
}

func (a *mercuryRemoteAggregator) Aggregate(triggerEventID string, responses [][]byte) (capabilities.CapabilityResponse, error) {
	latestReports := make(map[datastreams.FeedID]datastreams.FeedReport)
	latestTimestamps := make(map[datastreams.FeedID]int64)
	latestGlobalTs := int64(0)
	for _, response := range responses {
		unmarshaled, err := pb.UnmarshalCapabilityResponse(response)
		if err != nil {
			a.lggr.Errorw("could not unmarshal one of capability responses (faulty sender?)", "error", err)
			continue
		}
		triggerEvent := &capabilities.TriggerEvent{}
		if err = unmarshaled.Value.UnwrapTo(triggerEvent); err != nil {
			a.lggr.Errorw("could not unwrap one of trigger events", "error", err)
			continue
		}
		feedReports, err := a.codec.Unwrap(triggerEvent.Payload)
		if err != nil {
			a.lggr.Errorw("could not unwrap one of capability responses", "error", err)
			continue
		}
		// save latest valid report for each feed ID
		for _, report := range feedReports {
			latestTs, ok := latestTimestamps[datastreams.FeedID(report.FeedID)]
			if !ok || report.ObservationTimestamp > latestTs {
				latestReports[datastreams.FeedID(report.FeedID)] = report
				latestTimestamps[datastreams.FeedID(report.FeedID)] = report.ObservationTimestamp
			}
			if report.ObservationTimestamp > latestGlobalTs {
				latestGlobalTs = report.ObservationTimestamp
			}
		}
	}
	if len(latestReports) == 0 {
		return capabilities.CapabilityResponse{}, errors.New("no valid reports found")
	}
	reportList := []datastreams.FeedReport{}
	allIDs := []string{}
	for _, report := range latestReports {
		allIDs = append(allIDs, report.FeedID)
	}
	sort.Strings(allIDs)
	for _, feedID := range allIDs {
		reportList = append(reportList, latestReports[datastreams.FeedID(feedID)])
	}
	return wrapReports(reportList, triggerEventID, latestGlobalTs)
}
