package functions

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type FunctionsReportingPluginFactory struct {
	Logger    commontypes.Logger
	PluginORM functions.ORM
	JobID     uuid.UUID
}

var _ types.ReportingPluginFactory = (*FunctionsReportingPluginFactory)(nil)

type functionsReporting struct {
	logger         commontypes.Logger
	pluginORM      functions.ORM
	jobID          uuid.UUID
	reportCodec    *ReportCodec
	genericConfig  *types.ReportingPluginConfig
	specificConfig *config.ReportingPluginConfigWrapper
}

var _ types.ReportingPlugin = &functionsReporting{}

var (
	promReportingPlugins = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_reporting_plugin_restarts",
		Help: "Metric to track number of reporting plugin restarts",
	}, []string{"jobID"})

	promReportingPluginsQuery = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_reporting_plugin_query",
		Help: "Metric to track number of reporting plugin Query calls",
	}, []string{"jobID"})

	promReportingPluginsObservation = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_reporting_plugin_observation",
		Help: "Metric to track number of reporting plugin Observation calls",
	}, []string{"jobID"})

	promReportingPluginsReport = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_reporting_plugin_report",
		Help: "Metric to track number of reporting plugin Report calls",
	}, []string{"jobID"})

	promReportingAcceptReports = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_reporting_plugin_accept",
		Help: "Metric to track number of accepting reports",
	}, []string{"jobID"})

	promReportingTransmitReports = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_reporting_plugin_transmit",
		Help: "Metric to track number of transmiting reports",
	}, []string{"jobID"})

	promReportingTransmitBatchSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "functions_reporting_plugin_transmit_batch_size",
		Help:    "Metric to track batch size of transmitting reports",
		Buckets: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 100, 1000},
	}, []string{"jobID"})
)

func formatRequestId(requestId []byte) string {
	return fmt.Sprintf("0x%x", requestId)
}

// NewReportingPlugin complies with ReportingPluginFactory
func (f FunctionsReportingPluginFactory) NewReportingPlugin(rpConfig types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	pluginConfig, err := config.DecodeReportingPluginConfig(rpConfig.OffchainConfig)
	if err != nil {
		f.Logger.Error("unable to decode reporting plugin config", commontypes.LogFields{
			"digest": rpConfig.ConfigDigest.String(),
		})
		return nil, types.ReportingPluginInfo{}, err
	}
	codec, err := NewReportCodec()
	if err != nil {
		f.Logger.Error("unable to create a report codec object", commontypes.LogFields{})
		return nil, types.ReportingPluginInfo{}, err
	}
	info := types.ReportingPluginInfo{
		Name:          "functionsReporting",
		UniqueReports: pluginConfig.Config.GetUniqueReports(), // Enforces (N+F+1)/2 signatures. Must match setting in OCR2Base.sol.
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       int(pluginConfig.Config.GetMaxQueryLengthBytes()),
			MaxObservationLength: int(pluginConfig.Config.GetMaxObservationLengthBytes()),
			MaxReportLength:      int(pluginConfig.Config.GetMaxReportLengthBytes()),
		},
	}
	plugin := functionsReporting{
		logger:         f.Logger,
		pluginORM:      f.PluginORM,
		jobID:          f.JobID,
		reportCodec:    codec,
		genericConfig:  &rpConfig,
		specificConfig: pluginConfig,
	}
	promReportingPlugins.WithLabelValues(f.JobID.String()).Inc()
	return &plugin, info, nil
}

// Query() complies with ReportingPlugin
func (r *functionsReporting) Query(ctx context.Context, ts types.ReportTimestamp) (types.Query, error) {
	r.logger.Debug("FunctionsReporting Query start", commontypes.LogFields{
		"epoch":    ts.Epoch,
		"round":    ts.Round,
		"oracleID": r.genericConfig.OracleID,
	})
	maxBatchSize := r.specificConfig.Config.GetMaxRequestBatchSize()
	results, err := r.pluginORM.FindOldestEntriesByState(functions.RESULT_READY, maxBatchSize, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}

	queryProto := Query{}
	var idStrs []string
	for _, result := range results {
		result := result
		queryProto.RequestIDs = append(queryProto.RequestIDs, result.RequestID[:])
		idStrs = append(idStrs, formatRequestId(result.RequestID[:]))
	}
	r.logger.Debug("FunctionsReporting Query end", commontypes.LogFields{
		"epoch":      ts.Epoch,
		"round":      ts.Round,
		"oracleID":   r.genericConfig.OracleID,
		"queryLen":   len(queryProto.RequestIDs),
		"requestIDs": idStrs,
	})
	promReportingPluginsQuery.WithLabelValues(r.jobID.String()).Inc()
	return proto.Marshal(&queryProto)
}

// Observation() complies with ReportingPlugin
func (r *functionsReporting) Observation(ctx context.Context, ts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	r.logger.Debug("FunctionsReporting Observation start", commontypes.LogFields{
		"epoch":    ts.Epoch,
		"round":    ts.Round,
		"oracleID": r.genericConfig.OracleID,
	})

	queryProto := &Query{}
	err := proto.Unmarshal(query, queryProto)
	if err != nil {
		return nil, err
	}

	observationProto := Observation{}
	processedIds := make(map[[32]byte]bool)
	var idStrs []string
	for _, id := range queryProto.RequestIDs {
		id, err := sliceToByte32(id)
		if err != nil {
			r.logger.Error("FunctionsReporting Observation invalid ID", commontypes.LogFields{
				"requestID": formatRequestId(id[:]),
				"err":       err,
			})
			continue
		}
		if _, ok := processedIds[id]; ok {
			r.logger.Error("FunctionsReporting Observation duplicate ID in query", commontypes.LogFields{
				"requestID": formatRequestId(id[:]),
			})
			continue
		}
		processedIds[id] = true
		localResult, err2 := r.pluginORM.FindById(id, pg.WithParentCtx(ctx))
		if err2 != nil {
			r.logger.Debug("FunctionsReporting Observation can't find request from query", commontypes.LogFields{
				"requestID": formatRequestId(id[:]),
				"err":       err2,
			})
			continue
		}
		// NOTE: ignoring TIMED_OUT requests, which potentially had ready results
		if localResult.State == functions.RESULT_READY {
			resultProto := ProcessedRequest{
				RequestID: localResult.RequestID[:],
				Result:    localResult.Result,
				Error:     localResult.Error,
			}
			observationProto.ProcessedRequests = append(observationProto.ProcessedRequests, &resultProto)
			idStrs = append(idStrs, formatRequestId(localResult.RequestID[:]))
		}
	}
	r.logger.Debug("FunctionsReporting Observation end", commontypes.LogFields{
		"epoch":          ts.Epoch,
		"round":          ts.Round,
		"oracleID":       r.genericConfig.OracleID,
		"nReadyRequests": len(observationProto.ProcessedRequests),
		"requestIDs":     idStrs,
	})

	promReportingPluginsObservation.WithLabelValues(r.jobID.String()).Inc()
	return proto.Marshal(&observationProto)
}

// Report() complies with ReportingPlugin
func (r *functionsReporting) Report(ctx context.Context, ts types.ReportTimestamp, query types.Query, obs []types.AttributedObservation) (bool, types.Report, error) {
	r.logger.Debug("FunctionsReporting Report start", commontypes.LogFields{
		"epoch":         ts.Epoch,
		"round":         ts.Round,
		"oracleID":      r.genericConfig.OracleID,
		"nObservations": len(obs),
	})

	queryProto := &Query{}
	err := proto.Unmarshal(query, queryProto)
	if err != nil {
		r.logger.Error("FunctionsReporting Report: unable to decode query!",
			commontypes.LogFields{"err": err})
		return false, nil, err
	}

	reqIdToObservationList := make(map[string][]*ProcessedRequest)
	var uniqueQueryIds []string
	for _, id := range queryProto.RequestIDs {
		reqId := formatRequestId(id)
		if _, ok := reqIdToObservationList[reqId]; ok {
			r.logger.Error("FunctionsReporting Report: duplicate ID in query", commontypes.LogFields{
				"requestID": reqId,
			})
			continue
		}
		uniqueQueryIds = append(uniqueQueryIds, reqId)
		reqIdToObservationList[reqId] = []*ProcessedRequest{}
	}

	for _, ob := range obs {
		observationProto := &Observation{}
		err = proto.Unmarshal(ob.Observation, observationProto)
		if err != nil {
			r.logger.Error("FunctionsReporting Report: unable to decode observation!",
				commontypes.LogFields{"err": err, "observer": ob.Observer})
			continue
		}
		seenReqIds := make(map[string]struct{})
		for _, processedReq := range observationProto.ProcessedRequests {
			id := formatRequestId(processedReq.RequestID)
			if _, seen := seenReqIds[id]; seen {
				r.logger.Error("FunctionsReporting Report: observation contains duplicate IDs!",
					commontypes.LogFields{"requestID": id, "observer": ob.Observer})
				continue
			}
			if val, ok := reqIdToObservationList[id]; ok {
				reqIdToObservationList[id] = append(val, processedReq)
				seenReqIds[id] = struct{}{}
			} else {
				r.logger.Error("FunctionsReporting Report: observation contains ID that's not the query!",
					commontypes.LogFields{"requestID": id, "observer": ob.Observer})
			}
		}
	}

	defaultAggMethod := r.specificConfig.Config.GetDefaultAggregationMethod()
	var allAggregated []*ProcessedRequest
	var allIdStrs []string
	for _, reqId := range uniqueQueryIds {
		observations := reqIdToObservationList[reqId]
		if !CanAggregate(r.genericConfig.N, r.genericConfig.F, observations) {
			r.logger.Debug("FunctionsReporting Report: unable to aggregate request in current round", commontypes.LogFields{
				"epoch":         ts.Epoch,
				"round":         ts.Round,
				"requestID":     reqId,
				"nObservations": len(observations),
			})
			continue
		}

		// TODO: support per-request aggregation method
		// https://app.shortcut.com/chainlinklabs/story/57701/per-request-plugin-config
		aggregated, errAgg := Aggregate(defaultAggMethod, observations)
		if errAgg != nil {
			r.logger.Error("FunctionsReporting Report: error when aggregating reqId", commontypes.LogFields{
				"epoch":     ts.Epoch,
				"round":     ts.Round,
				"requestID": reqId,
				"err":       errAgg,
			})
			continue
		}
		r.logger.Debug("FunctionsReporting Report: aggregated successfully", commontypes.LogFields{
			"epoch":         ts.Epoch,
			"round":         ts.Round,
			"requestID":     reqId,
			"nObservations": len(observations),
		})
		allAggregated = append(allAggregated, aggregated)
		allIdStrs = append(allIdStrs, reqId)
	}

	r.logger.Debug("FunctionsReporting Report end", commontypes.LogFields{
		"epoch":               ts.Epoch,
		"round":               ts.Round,
		"oracleID":            r.genericConfig.OracleID,
		"nAggregatedRequests": len(allAggregated),
		"reporting":           len(allAggregated) > 0,
		"requestIDs":          allIdStrs,
	})
	if len(allAggregated) == 0 {
		return false, nil, nil
	}
	reportBytes, err := r.reportCodec.EncodeReport(allAggregated)
	if err != nil {
		return false, nil, err
	}
	promReportingPluginsReport.WithLabelValues(r.jobID.String()).Inc()
	return true, reportBytes, nil
}

// ShouldAcceptFinalizedReport() complies with ReportingPlugin
func (r *functionsReporting) ShouldAcceptFinalizedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	r.logger.Debug("FunctionsReporting ShouldAcceptFinalizedReport start", commontypes.LogFields{
		"epoch":    ts.Epoch,
		"round":    ts.Round,
		"oracleID": r.genericConfig.OracleID,
	})

	// NOTE: The output of the Report() phase needs to be later decoded by the contract. So unfortunately we
	// can't use anything more convenient like protobufs but we need to ABI-decode here instead.
	decoded, err := r.reportCodec.DecodeReport(report)
	if err != nil {
		r.logger.Error("FunctionsReporting ShouldAcceptFinalizedReport: unable to decode report built in reporting phase", commontypes.LogFields{"err": err})
		return false, err
	}

	allIds := []string{}
	needTransmissionIds := []string{}
	for _, item := range decoded {
		reqIdStr := formatRequestId(item.RequestID)
		allIds = append(allIds, reqIdStr)
		id, err := sliceToByte32(item.RequestID)
		if err != nil {
			r.logger.Error("FunctionsReporting ShouldAcceptFinalizedReport: invalid ID", commontypes.LogFields{"requestID": reqIdStr, "err": err})
			continue
		}
		_, err = r.pluginORM.FindById(id, pg.WithParentCtx(ctx))
		if err != nil {
			// TODO: Differentiate between ID not found and other ORM errors (https://smartcontract-it.atlassian.net/browse/DRO-215)
			r.logger.Warn("FunctionsReporting ShouldAcceptFinalizedReport: request doesn't exist locally! Accepting anyway.", commontypes.LogFields{"requestID": reqIdStr})
			needTransmissionIds = append(needTransmissionIds, reqIdStr)
			continue
		}
		err = r.pluginORM.SetFinalized(id, item.Result, item.Error, pg.WithParentCtx(ctx)) // validates state transition
		if err != nil {
			r.logger.Debug("FunctionsReporting ShouldAcceptFinalizedReport: state couldn't be changed to FINALIZED. Not transmitting.", commontypes.LogFields{"requestID": reqIdStr, "err": err})
			continue
		}
		needTransmissionIds = append(needTransmissionIds, reqIdStr)
	}
	r.logger.Debug("FunctionsReporting ShouldAcceptFinalizedReport end", commontypes.LogFields{
		"epoch":               ts.Epoch,
		"round":               ts.Round,
		"oracleID":            r.genericConfig.OracleID,
		"allIds":              allIds,
		"needTransmissionIds": needTransmissionIds,
		"accepting":           len(needTransmissionIds) > 0,
	})
	shouldAccept := len(needTransmissionIds) > 0
	if shouldAccept {
		promReportingAcceptReports.WithLabelValues(r.jobID.String()).Inc()
	}
	return shouldAccept, nil
}

// ShouldTransmitAcceptedReport() complies with ReportingPlugin
func (r *functionsReporting) ShouldTransmitAcceptedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	r.logger.Debug("FunctionsReporting ShouldTransmitAcceptedReport start", commontypes.LogFields{
		"epoch":    ts.Epoch,
		"round":    ts.Round,
		"oracleID": r.genericConfig.OracleID,
	})

	decoded, err := r.reportCodec.DecodeReport(report)
	if err != nil {
		r.logger.Error("FunctionsReporting ShouldTransmitAcceptedReport: unable to decode report built in reporting phase", commontypes.LogFields{"err": err})
		return false, err
	}

	allIds := []string{}
	needTransmissionIds := []string{}
	for _, item := range decoded {
		reqIdStr := formatRequestId(item.RequestID)
		allIds = append(allIds, reqIdStr)
		id, err := sliceToByte32(item.RequestID)
		if err != nil {
			r.logger.Error("FunctionsReporting ShouldAcceptFinalizedReport: invalid ID", commontypes.LogFields{"requestID": reqIdStr, "err": err})
			continue
		}
		request, err := r.pluginORM.FindById(id, pg.WithParentCtx(ctx))
		if err != nil {
			r.logger.Warn("FunctionsReporting ShouldTransmitAcceptedReport: request doesn't exist locally! Transmitting anyway.", commontypes.LogFields{"requestID": reqIdStr, "err": err})
			needTransmissionIds = append(needTransmissionIds, reqIdStr)
			continue
		}
		if request.State == functions.TIMED_OUT || request.State == functions.CONFIRMED {
			r.logger.Debug("FunctionsReporting ShouldTransmitAcceptedReport: request is not FINALIZED any more. Not transmitting.",
				commontypes.LogFields{
					"requestID": reqIdStr,
					"state":     request.State.String(),
				})
			continue
		}
		if request.State == functions.IN_PROGRESS || request.State == functions.RESULT_READY {
			r.logger.Warn("FunctionsReporting ShouldTransmitAcceptedReport: unusual request state. Still transmitting.",
				commontypes.LogFields{
					"requestID": reqIdStr,
					"state":     request.State.String(),
				})
		}
		needTransmissionIds = append(needTransmissionIds, reqIdStr)
	}
	r.logger.Debug("FunctionsReporting ShouldTransmitAcceptedReport end", commontypes.LogFields{
		"epoch":               ts.Epoch,
		"round":               ts.Round,
		"oracleID":            r.genericConfig.OracleID,
		"allIds":              allIds,
		"needTransmissionIds": needTransmissionIds,
		"transmitting":        len(needTransmissionIds) > 0,
	})
	shouldTransmit := len(needTransmissionIds) > 0
	if shouldTransmit {
		promReportingTransmitReports.WithLabelValues(r.jobID.String()).Inc()
		promReportingTransmitBatchSize.WithLabelValues(r.jobID.String()).Observe(float64(len(allIds)))
	}
	return shouldTransmit, nil
}

// Close() complies with ReportingPlugin
func (r *functionsReporting) Close() error {
	r.logger.Debug("FunctionsReporting Close", commontypes.LogFields{
		"oracleID": r.genericConfig.OracleID,
	})
	return nil
}
