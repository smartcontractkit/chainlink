package directrequestocr

import (
	"context"
	"encoding/hex"

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
	reportCodec    *reportCodec
	genericConfig  *types.ReportingPluginConfig
	specificConfig *config.ReportingPluginConfigWrapper
}

var _ types.ReportingPlugin = &directRequestReporting{}

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
		"epoch": ts.Epoch,
		"round": ts.Round,
	})
	maxBatchSize := r.specificConfig.Config.GetMaxRequestBatchSize()
	results, err := r.pluginORM.FindOldestEntriesByState(directrequestocr.RESULT_READY, maxBatchSize)
	if err != nil {
		return nil, err
	}

	queryProto := Query{}
	for _, result := range results {
		queryProto.RequestIDs = append(queryProto.RequestIDs, result.ContractRequestID[:])
	}
	return proto.Marshal(&queryProto)
}

// Observation() complies with ReportingPlugin
func (r *directRequestReporting) Observation(ctx context.Context, ts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	r.logger.Debug("directRequestReporting Observation phase", commontypes.LogFields{
		"epoch": ts.Epoch,
		"round": ts.Round,
	})

	queryProto := &Query{}
	err := proto.Unmarshal(query, queryProto)
	if err != nil {
		return nil, err
	}

	observationProto := Observation{}
	for _, id := range queryProto.RequestIDs {
		localResult, _ := r.pluginORM.FindById(sliceToByte32(id))
		if localResult.State == directrequestocr.RESULT_READY {
			resultProto := ProcessedRequest{
				RequestID: localResult.ContractRequestID[:],
				Result:    localResult.Result,
				Error:     []byte(localResult.Error),
			}
			observationProto.ProcessedRequests = append(observationProto.ProcessedRequests, &resultProto)
		}
	}
	r.logger.Debug("directRequestReporting Observation phase done", commontypes.LogFields{
		"nReadyRequests": len(observationProto.ProcessedRequests),
	})

	return proto.Marshal(&observationProto)
}

// Report() complies with ReportingPlugin
func (r *directRequestReporting) Report(ctx context.Context, ts types.ReportTimestamp, query types.Query, obs []types.AttributedObservation) (bool, types.Report, error) {
	r.logger.Debug("directRequestReporting Report phase", commontypes.LogFields{
		"epoch":         ts.Epoch,
		"round":         ts.Round,
		"nObservations": len(obs),
	})

	queryProto := &Query{}
	err := proto.Unmarshal(query, queryProto)
	if err != nil {
		r.logger.Error("directRequestReporting Report phase unable to decode query!",
			commontypes.LogFields{"err": err})
		return false, nil, err
	}

	reqIdToResultList := make(map[string][]*ProcessedRequest)
	for _, id := range queryProto.RequestIDs {
		reqIdToResultList[string(id)] = []*ProcessedRequest{}
	}

	for _, ob := range obs {
		observationProto := &Observation{}
		err = proto.Unmarshal(ob.Observation, observationProto)
		if err != nil {
			r.logger.Error("directRequestReporting Report phase unable to decode observation!",
				commontypes.LogFields{"err": err, "observer": ob.Observer})
			continue
		}
		for _, res := range observationProto.ProcessedRequests {
			id := string(res.RequestID)
			if val, ok := reqIdToResultList[id]; ok {
				reqIdToResultList[id] = append(val, res)
			}
		}
	}

	// TODO make aggregation modular and configurable with Median as default.
	// https://app.shortcut.com/chainlinklabs/story/56740/modular-aggregation
	const minRequiredObservations = 3
	var aggregated []*ProcessedRequest
	for _, obsArr := range reqIdToResultList {
		if len(obsArr) >= minRequiredObservations {
			aggregated = append(aggregated, obsArr[0])
		}
	}

	r.logger.Debug("directRequestReporting Report phase done", commontypes.LogFields{
		"nAggregatedRequests": len(aggregated),
		"reporting":           len(aggregated) > 0,
	})
	if len(aggregated) == 0 {
		return false, nil, nil
	}
	reportBytes, err := r.reportCodec.EncodeReport(aggregated)
	if err != nil {
		return false, nil, err
	}
	return true, reportBytes, nil

}

// ShouldAcceptFinalizedReport() complies with ReportingPlugin
func (r *directRequestReporting) ShouldAcceptFinalizedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	r.logger.Debug("directRequestReporting ShouldAcceptFinalizedReport phase", commontypes.LogFields{
		"epoch": ts.Epoch,
		"round": ts.Round,
	})
	return true, nil
}

// ShouldTransmitAcceptedReport() complies with ReportingPlugin
func (r *directRequestReporting) ShouldTransmitAcceptedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	r.logger.Debug("directRequestReporting ShouldTransmitAcceptedReport phase", commontypes.LogFields{
		"epoch": ts.Epoch,
		"round": ts.Round,
	})

	// NOTE: The output of the Report() phase needs to be later decoded by the contract. So unfortunately we
	// can't use anything more convenient like protobufs but we need to ABI-decode here instead.
	decoded, err := r.reportCodec.DecodeReport(report)
	if err != nil {
		r.logger.Error("unable to decode report built in reporting phase", commontypes.LogFields{"err": err})
		return false, err
	}

	allIds := []string{}
	needTransmissionIds := []string{}
	for _, item := range decoded {
		reqIdStr := hex.EncodeToString(item.RequestID)
		allIds = append(allIds, reqIdStr)
		prevState, err := r.pluginORM.SetState(sliceToByte32(item.RequestID), directrequestocr.TRANSMITTED)
		if err != nil {
			// TODO handle state-transition errors inside the ORM (e.g. can't move from CONFIRMED back to TRANSMITTED)
			// TODO it's possible that report will have results for requests that I never received
			// https://app.shortcut.com/chainlinklabs/story/54049/database-table-in-core-node
			r.logger.Debug("directRequestReporting unable to set state to TRANSMITTED", commontypes.LogFields{"requestID": item.RequestID})
			needTransmissionIds = append(needTransmissionIds, reqIdStr)
		} else if prevState != directrequestocr.TRANSMITTED && prevState != directrequestocr.CONFIRMED {
			needTransmissionIds = append(needTransmissionIds, reqIdStr)
		}
	}
	r.logger.Debug("directRequestReporting ShouldTransmitAcceptedReport phase done", commontypes.LogFields{
		"allIds":              allIds,
		"needTransmissionIds": needTransmissionIds,
		"reporting":           len(needTransmissionIds) > 0,
	})
	return len(needTransmissionIds) > 0, nil
}

// Close() complies with ReportingPlugin
func (r *directRequestReporting) Close() error {
	r.logger.Debug("directRequestReporting Close", commontypes.LogFields{})
	return nil
}
