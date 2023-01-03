package directrequestocr

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/services/directrequestocr"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr/config"
)

type DirectRequestReportingPluginFactory struct {
	Logger    commontypes.Logger
	PluginORM directrequestocr.ORM
	JobID     uuid.UUID
}

var _ types.ReportingPluginFactory = (*DirectRequestReportingPluginFactory)(nil)

type directRequestReporting struct {
	logger         commontypes.Logger
	pluginORM      directrequestocr.ORM
	jobID          uuid.UUID
	reportCodec    *ReportCodec
	genericConfig  *types.ReportingPluginConfig
	specificConfig *config.ReportingPluginConfigWrapper
}

var _ types.ReportingPlugin = &directRequestReporting{}

func formatRequestId(requestId []byte) string {
	return fmt.Sprintf("0x%x", requestId)
}

// NewReportingPlugin complies with ReportingPluginFactory
func (f DirectRequestReportingPluginFactory) NewReportingPlugin(rpConfig types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
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
		Name:          "directRequestReporting",
		UniqueReports: pluginConfig.Config.GetUniqueReports(), // Enforces (N+F+1)/2 signatures. Must match setting in OCR2Base.sol.
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       int(pluginConfig.Config.GetMaxQueryLengthBytes()),
			MaxObservationLength: int(pluginConfig.Config.GetMaxObservationLengthBytes()),
			MaxReportLength:      int(pluginConfig.Config.GetMaxReportLengthBytes()),
		},
	}
	plugin := directRequestReporting{
		logger:         f.Logger,
		pluginORM:      f.PluginORM,
		jobID:          f.JobID,
		reportCodec:    codec,
		genericConfig:  &rpConfig,
		specificConfig: pluginConfig,
	}
	return &plugin, info, nil
}

// Query() complies with ReportingPlugin
func (r *directRequestReporting) Query(ctx context.Context, ts types.ReportTimestamp) (types.Query, error) {
	r.logger.Debug("directRequestReporting Query phase", commontypes.LogFields{
		"epoch":    ts.Epoch,
		"round":    ts.Round,
		"oracleID": r.genericConfig.OracleID,
	})
	maxBatchSize := r.specificConfig.Config.GetMaxRequestBatchSize()
	results, err := r.pluginORM.FindOldestEntriesByState(directrequestocr.RESULT_READY, maxBatchSize)
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
	r.logger.Debug("directRequestReporting Query phase done", commontypes.LogFields{
		"epoch":      ts.Epoch,
		"round":      ts.Round,
		"oracleID":   r.genericConfig.OracleID,
		"queryLen":   len(queryProto.RequestIDs),
		"requestIDs": idStrs,
	})
	return proto.Marshal(&queryProto)
}

// Observation() complies with ReportingPlugin
func (r *directRequestReporting) Observation(ctx context.Context, ts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	r.logger.Debug("directRequestReporting Observation phase", commontypes.LogFields{
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
		id := sliceToByte32(id)
		if _, ok := processedIds[id]; ok {
			r.logger.Error("directRequestReporting Observation phase duplicate ID in query", commontypes.LogFields{
				"requestID": formatRequestId(id[:]),
			})
			continue
		}
		processedIds[id] = true
		localResult, err2 := r.pluginORM.FindById(id)
		if err2 != nil {
			r.logger.Debug("directRequestReporting Observation phase can't find request from query", commontypes.LogFields{
				"requestID": formatRequestId(id[:]),
				"err":       err2,
			})
			continue
		}
		// NOTE: ignoring TIMED_OUT requests, which potentially had ready results
		if localResult.State == directrequestocr.RESULT_READY {
			resultProto := ProcessedRequest{
				RequestID: localResult.RequestID[:],
				Result:    localResult.Result,
				Error:     localResult.Error,
			}
			observationProto.ProcessedRequests = append(observationProto.ProcessedRequests, &resultProto)
			idStrs = append(idStrs, formatRequestId(localResult.RequestID[:]))
		}
	}
	r.logger.Debug("directRequestReporting Observation phase done", commontypes.LogFields{
		"epoch":          ts.Epoch,
		"round":          ts.Round,
		"oracleID":       r.genericConfig.OracleID,
		"nReadyRequests": len(observationProto.ProcessedRequests),
		"requestIDs":     idStrs,
	})

	return proto.Marshal(&observationProto)
}

// Report() complies with ReportingPlugin
func (r *directRequestReporting) Report(ctx context.Context, ts types.ReportTimestamp, query types.Query, obs []types.AttributedObservation) (bool, types.Report, error) {
	r.logger.Debug("directRequestReporting Report phase", commontypes.LogFields{
		"epoch":         ts.Epoch,
		"round":         ts.Round,
		"oracleID":      r.genericConfig.OracleID,
		"nObservations": len(obs),
	})

	queryProto := &Query{}
	err := proto.Unmarshal(query, queryProto)
	if err != nil {
		r.logger.Error("directRequestReporting Report phase unable to decode query!",
			commontypes.LogFields{"err": err})
		return false, nil, err
	}

	reqIdToObservationList := make(map[string][]*ProcessedRequest)
	for _, id := range queryProto.RequestIDs {
		reqIdToObservationList[formatRequestId(id)] = []*ProcessedRequest{}
	}

	for _, ob := range obs {
		observationProto := &Observation{}
		err = proto.Unmarshal(ob.Observation, observationProto)
		if err != nil {
			r.logger.Error("directRequestReporting Report phase unable to decode observation!",
				commontypes.LogFields{"err": err, "observer": ob.Observer})
			continue
		}
		for _, processedReq := range observationProto.ProcessedRequests {
			id := formatRequestId(processedReq.RequestID)
			if val, ok := reqIdToObservationList[id]; ok {
				reqIdToObservationList[id] = append(val, processedReq)
			}
		}
	}

	defaultAggMethod := r.specificConfig.Config.GetDefaultAggregationMethod()
	var allAggregated []*ProcessedRequest
	var allIdStrs []string
	for reqId, observations := range reqIdToObservationList {
		if !CanAggregate(r.genericConfig.N, r.genericConfig.F, observations) {
			r.logger.Debug("directRequestReporting unable to aggregate request in current round", commontypes.LogFields{
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
			r.logger.Error("directRequestReporting error when aggregating reqId", commontypes.LogFields{
				"epoch":     ts.Epoch,
				"round":     ts.Round,
				"requestID": reqId,
				"err":       errAgg,
			})
			continue
		}
		allAggregated = append(allAggregated, aggregated)
		allIdStrs = append(allIdStrs, reqId)
	}

	r.logger.Debug("directRequestReporting Report phase done", commontypes.LogFields{
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
	return true, reportBytes, nil
}

// ShouldAcceptFinalizedReport() complies with ReportingPlugin
func (r *directRequestReporting) ShouldAcceptFinalizedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	r.logger.Debug("directRequestReporting ShouldAcceptFinalizedReport phase", commontypes.LogFields{
		"epoch":    ts.Epoch,
		"round":    ts.Round,
		"oracleID": r.genericConfig.OracleID,
	})

	// NOTE: The output of the Report() phase needs to be later decoded by the contract. So unfortunately we
	// can't use anything more convenient like protobufs but we need to ABI-decode here instead.
	decoded, err := r.reportCodec.DecodeReport(report)
	if err != nil {
		r.logger.Error("directRequestReporting unable to decode report built in reporting phase", commontypes.LogFields{"err": err})
		return false, err
	}

	allIds := []string{}
	needTransmissionIds := []string{}
	for _, item := range decoded {
		reqIdStr := formatRequestId(item.RequestID)
		allIds = append(allIds, reqIdStr)
		_, err := r.pluginORM.FindById(sliceToByte32(item.RequestID))
		if err != nil {
			r.logger.Warn("directRequestReporting request doesn't exist locally! Accepting anyway.", commontypes.LogFields{"requestID": reqIdStr})
			needTransmissionIds = append(needTransmissionIds, reqIdStr)
			continue
		}
		err = r.pluginORM.SetFinalized(sliceToByte32(item.RequestID), item.Result, item.Error) // validates state transition
		if err != nil {
			r.logger.Debug("directRequestReporting state couldn't be changed to FINALIZED. Not transmitting.", commontypes.LogFields{"requestID": reqIdStr})
			continue
		}
		needTransmissionIds = append(needTransmissionIds, reqIdStr)
	}
	r.logger.Debug("directRequestReporting ShouldAcceptFinalizedReport phase done", commontypes.LogFields{
		"epoch":               ts.Epoch,
		"round":               ts.Round,
		"oracleID":            r.genericConfig.OracleID,
		"allIds":              allIds,
		"needTransmissionIds": needTransmissionIds,
		"reporting":           len(needTransmissionIds) > 0,
	})
	return len(needTransmissionIds) > 0, nil
}

// ShouldTransmitAcceptedReport() complies with ReportingPlugin
func (r *directRequestReporting) ShouldTransmitAcceptedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	r.logger.Debug("directRequestReporting ShouldTransmitAcceptedReport phase", commontypes.LogFields{
		"epoch":    ts.Epoch,
		"round":    ts.Round,
		"oracleID": r.genericConfig.OracleID,
	})

	decoded, err := r.reportCodec.DecodeReport(report)
	if err != nil {
		r.logger.Error("directRequestReporting unable to decode report built in reporting phase", commontypes.LogFields{"err": err})
		return false, err
	}

	allIds := []string{}
	needTransmissionIds := []string{}
	for _, item := range decoded {
		reqIdStr := formatRequestId(item.RequestID)
		allIds = append(allIds, reqIdStr)
		request, err := r.pluginORM.FindById(sliceToByte32(item.RequestID))
		if err != nil {
			r.logger.Warn("directRequestReporting request doesn't exist locally! Transmitting anyway.", commontypes.LogFields{"requestID": reqIdStr})
			needTransmissionIds = append(needTransmissionIds, reqIdStr)
			continue
		}
		if request.State == directrequestocr.TIMED_OUT || request.State == directrequestocr.CONFIRMED {
			r.logger.Debug("directRequestReporting request is not FINALIZED any more. Not transmitting.",
				commontypes.LogFields{
					"requestID": reqIdStr,
					"state":     request.State.String(),
				})
			continue
		}
		if request.State == directrequestocr.IN_PROGRESS || request.State == directrequestocr.RESULT_READY {
			r.logger.Warn("directRequestReporting unusual request state. Still transmitting.",
				commontypes.LogFields{
					"requestID": reqIdStr,
					"state":     request.State.String(),
				})
		}
		needTransmissionIds = append(needTransmissionIds, reqIdStr)
	}
	r.logger.Debug("directRequestReporting ShouldTransmitAcceptedReport phase done", commontypes.LogFields{
		"epoch":               ts.Epoch,
		"round":               ts.Round,
		"oracleID":            r.genericConfig.OracleID,
		"allIds":              allIds,
		"needTransmissionIds": needTransmissionIds,
		"reporting":           len(needTransmissionIds) > 0,
	})
	return len(needTransmissionIds) > 0, nil
}

// Close() complies with ReportingPlugin
func (r *directRequestReporting) Close() error {
	r.logger.Debug("directRequestReporting Close", commontypes.LogFields{
		"oracleID": r.genericConfig.OracleID,
	})
	return nil
}
