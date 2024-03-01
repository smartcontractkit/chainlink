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
		r := &pbtypes.Observation{
			Observation: values.Proto(rq.Observations),
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
			if _, ok := m[weid]; !ok {
				m[weid] = make(map[ocrcommon.OracleID][]values.Value)
			}

			m[weid][o.Observer] = append(m[weid][o.Observer], values.FromProto(rq.Observation))
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

	// Wipe out the ReportsToGenerate. This gets regenerated
	// every time since we only want to transmit reports that
	// are part of the current Query.
	o.ReportsToGenerate = []*pbtypes.Report{}

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

		if outcome.ShouldReport {
			report := &pbtypes.Report{
				Outcome: outcome,
				Id:      weid,
			}
			o.ReportsToGenerate = append(o.ReportsToGenerate, report)
		}

		o.Outcomes[weid.WorkflowId] = outcome
	}

	r.lggr.Debugw("Outcome complete", "len", len(o.Outcomes), "nReportsToGenerate", len(o.ReportsToGenerate))
	return proto.Marshal(o)
}

func (r *reportingPlugin) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[[]byte], error) {
	o := &pbtypes.Outcome{}
	err := proto.Unmarshal(outcome, o)
	if err != nil {
		return nil, err
	}

	reports := []ocr3types.ReportWithInfo[[]byte]{}

	// This doesn't handle a query which contains the same workflowId multiple times.
	for _, report := range o.ReportsToGenerate {
		outcome, id := report.Outcome, report.Id
		outcome, err := pbtypes.AppendWorkflowIDs(outcome, id.WorkflowId, id.WorkflowExecutionId)
		if err != nil {
			r.lggr.Errorw("could not append IDs")
			continue
		}

		enc, err := r.r.getEncoder(id.WorkflowId)
		if err != nil {
			r.lggr.Errorw("could not retrieve encoder for workflow", "error", err, "workflowID", id.WorkflowId)
			continue
		}

		mv := values.FromMapValueProto(outcome.EncodableOutcome)
		report, err := enc.Encode(context.Background(), *mv)
		if err != nil {
			r.lggr.Errorw("could not encode report for workflow", "error", err, "workflowID", id.WorkflowId)
			continue
		}

		p, err := proto.Marshal(id)
		if err != nil {
			r.lggr.Errorw("could not marshal id into ReportWithInfo", "error", err, "workflowID", id.WorkflowId)
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
	id := &pbtypes.Id{}
	err := proto.Unmarshal(rwi.Info, id)
	if err != nil {
		r.lggr.Error("could not unmarshal id")
		return false, err
	}

	b := values.NewBytes(rwi.Report)
	r.lggr.Debugw("ShouldAcceptAttestedReport transmitting", "len", len(b.Underlying))
	err = r.r.transmitResponse(ctx, &response{
		CapabilityResponse: capabilities.CapabilityResponse{
			Value: b,
		},
		WorkflowExecutionID: id.WorkflowExecutionId,
	})
	if err != nil {
		r.lggr.Errorw("could not transmit response", "error", err, "weid", id.WorkflowExecutionId)
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
