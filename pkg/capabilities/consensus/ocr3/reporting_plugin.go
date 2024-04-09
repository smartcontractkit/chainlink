package ocr3

import (
	"context"

	ocrcommon "github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	pbtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

var _ ocr3types.ReportingPlugin[[]byte] = (*reportingPlugin)(nil)

type capabilityIface interface {
	transmitResponse(ctx context.Context, resp *response) error
	getAggregator(workflowID string) (pbtypes.Aggregator, error)
	getEncoder(workflowID string) (pbtypes.Encoder, error)
}

type reportingPlugin struct {
	batchSize int
	s         *store
	r         capabilityIface
	config    ocr3types.ReportingPluginConfig
	lggr      logger.Logger
}

func newReportingPlugin(s *store, r capabilityIface, batchSize int, config ocr3types.ReportingPluginConfig, lggr logger.Logger) (*reportingPlugin, error) {
	// TODO: extract limits from OnchainConfig
	// and perform validation.

	return &reportingPlugin{
		s:         s,
		r:         r,
		batchSize: batchSize,
		config:    config,
		lggr:      logger.Named(lggr, "OCR3ConsensusReportingPlugin"),
	}, nil
}

func (r *reportingPlugin) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
	batch, err := r.s.firstN(ctx, r.batchSize)
	if err != nil {
		r.lggr.Errorw("could not retrieve batch", "error", err)
		return nil, err
	}

	ids := []*pbtypes.Id{}
	for _, r := range batch {
		ids = append(ids, &pbtypes.Id{
			WorkflowExecutionId: r.WorkflowExecutionID,
			WorkflowId:          r.WorkflowID,
		})
	}

	r.lggr.Debugw("Query complete", "len", len(ids))
	return proto.Marshal(&pbtypes.Query{
		Ids: ids,
	})
}

func (r *reportingPlugin) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query) (types.Observation, error) {
	queryReq := &pbtypes.Query{}
	err := proto.Unmarshal(query, queryReq)
	if err != nil {
		return nil, err
	}

	weids := []string{}
	for _, q := range queryReq.Ids {
		weids = append(weids, q.WorkflowExecutionId)
	}

	reqs, err := r.s.getN(ctx, weids)
	if err != nil {
		return nil, err
	}

	obs := &pbtypes.Observations{}
	for _, rq := range reqs {
		listProto := values.Proto(rq.Observations).GetListValue()
		if listProto == nil {
			r.lggr.Errorw("observations are not a list", "weID", rq.WorkflowExecutionID)
			continue
		}
		r := &pbtypes.Observation{
			Observations: listProto,
			Id: &pbtypes.Id{
				WorkflowExecutionId: rq.WorkflowExecutionID,
				WorkflowId:          rq.WorkflowID,
			},
		}

		obs.Observations = append(obs.Observations, r)
	}

	r.lggr.Debugw("Observation complete", "len", len(obs.Observations), "queryLen", len(queryReq.Ids))
	return proto.Marshal(obs)
}

func (r *reportingPlugin) ValidateObservation(outctx ocr3types.OutcomeContext, query types.Query, ao types.AttributedObservation) error {
	return nil
}

func (r *reportingPlugin) ObservationQuorum(outctx ocr3types.OutcomeContext, query types.Query) (ocr3types.Quorum, error) {
	return ocr3types.QuorumTwoFPlusOne, nil
}

func (r *reportingPlugin) Outcome(outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (ocr3types.Outcome, error) {
	// execution ID -> oracle ID -> list of observations
	m := map[string]map[ocrcommon.OracleID][]values.Value{}
	for _, o := range aos {
		obs := &pbtypes.Observations{}
		err := proto.Unmarshal(o.Observation, obs)
		if err != nil {
			r.lggr.Errorw("could not unmarshal observation", "error", err, "observation", obs)
			continue
		}

		for _, rq := range obs.Observations {
			weid := rq.Id.WorkflowExecutionId

			obsList := values.FromListValueProto(rq.Observations)
			if obsList == nil {
				r.lggr.Errorw("observations are not a list", "weID", weid, "oracleID", o.Observer)
				continue
			}

			if _, ok := m[weid]; !ok {
				m[weid] = make(map[ocrcommon.OracleID][]values.Value)
			}
			m[weid][o.Observer] = obsList.Underlying
		}
	}

	q := &pbtypes.Query{}
	err := proto.Unmarshal(query, q)
	if err != nil {
		return nil, err
	}

	o := &pbtypes.Outcome{}
	err = proto.Unmarshal(outctx.PreviousOutcome, o)
	if err != nil {
		return nil, err
	}
	if o.Outcomes == nil {
		o.Outcomes = map[string]*pbtypes.AggregationOutcome{}
	}

	// Wipe out the CurrentReports. This gets regenerated
	// every time since we only want to transmit reports that
	// are part of the current Query.
	o.CurrentReports = []*pbtypes.Report{}

	for _, weid := range q.Ids {
		obs, ok := m[weid.WorkflowExecutionId]
		if !ok {
			r.lggr.Debugw("could not find any observations matching weid requested in the query", "weid", weid.WorkflowExecutionId)
			continue
		}

		workflowOutcome, ok := o.Outcomes[weid.WorkflowId]
		if !ok {
			r.lggr.Debugw("could not find existing outcome for workflow, aggregator will create a new one", "workflowID", weid.WorkflowId)
		}

		if len(obs) < (2*r.config.F + 1) {
			r.lggr.Debugw("insufficient observations for workflow execution id", "weid", weid.WorkflowExecutionId)
			continue
		}

		agg, err := r.r.getAggregator(weid.WorkflowId)
		if err != nil {
			r.lggr.Errorw("could not retrieve aggregator for workflow", "error", err, "workflowID", weid.WorkflowId)
			continue
		}

		outcome, err := agg.Aggregate(workflowOutcome, obs)
		if err != nil {
			r.lggr.Errorw("error aggregating outcome", "error", err, "workflowID", weid.WorkflowId)
			return nil, err
		}

		report := &pbtypes.Report{
			Outcome: outcome,
			Id:      weid,
		}
		o.CurrentReports = append(o.CurrentReports, report)

		o.Outcomes[weid.WorkflowId] = outcome
	}

	r.lggr.Debugw("Outcome complete", "len", len(o.Outcomes), "nCurrentReports", len(o.CurrentReports))
	return proto.Marshal(o)
}

func (r *reportingPlugin) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[[]byte], error) {
	o := &pbtypes.Outcome{}
	err := proto.Unmarshal(outcome, o)
	if err != nil {
		return nil, err
	}

	reports := []ocr3types.ReportWithInfo[[]byte]{}

	for _, report := range o.CurrentReports {
		outcome, id := report.Outcome, report.Id

		info := &pbtypes.ReportInfo{
			Id:           id,
			ShouldReport: outcome.ShouldReport,
		}

		var report []byte
		if info.ShouldReport {
			newOutcome, err := pbtypes.AppendWorkflowIDs(outcome, id.WorkflowId, id.WorkflowExecutionId)
			if err != nil {
				r.lggr.Errorw("could not append IDs")
				continue
			}

			enc, err := r.r.getEncoder(id.WorkflowId)
			if err != nil {
				r.lggr.Errorw("could not retrieve encoder for workflow", "error", err, "workflowID", id.WorkflowId)
				continue
			}

			mv := values.FromMapValueProto(newOutcome.EncodableOutcome)
			report, err = enc.Encode(context.TODO(), *mv)
			if err != nil {
				r.lggr.Errorw("could not encode report for workflow", "error", err, "workflowID", id.WorkflowId)
				continue
			}
		}

		p, err := proto.Marshal(info)
		if err != nil {
			r.lggr.Errorw("could not marshal id into ReportWithInfo", "error", err, "workflowID", id.WorkflowId, "shouldReport", info.ShouldReport)
			continue
		}

		reports = append(reports, ocr3types.ReportWithInfo[[]byte]{
			Report: report,
			Info:   p,
		})
	}

	r.lggr.Debugw("Reports complete", "len", len(reports))
	return reports, nil
}

func (r *reportingPlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	info := &pbtypes.ReportInfo{}
	err := proto.Unmarshal(rwi.Info, info)
	if err != nil {
		r.lggr.Error("could not unmarshal info")
		return false, err
	}

	resp := map[string]any{}
	if info.ShouldReport {
		resp["report"] = []byte(rwi.Report)
	} else {
		resp["report"] = nil
	}

	v, err := values.Wrap(resp)
	if err != nil {
		r.lggr.Error("could not wrap report", "payload", resp)
		return false, err
	}

	r.lggr.Debugw("ShouldAcceptAttestedReport transmitting", "shouldReport", info.ShouldReport, "len", len(rwi.Report))
	err = r.r.transmitResponse(ctx, &response{
		CapabilityResponse: capabilities.CapabilityResponse{
			Value: v,
		},
		WorkflowExecutionID: info.Id.WorkflowExecutionId,
	})
	if err != nil {
		r.lggr.Errorw("could not transmit response", "error", err, "weid", info.Id.WorkflowExecutionId)
		return false, err
	}

	return false, nil
}

func (r *reportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, nil
}

func (r *reportingPlugin) Close() error {
	return nil
}
