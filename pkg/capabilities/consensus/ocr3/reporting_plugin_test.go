package ocr3

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	pbtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

func TestReportingPlugin_Query_ErrorInQueueCall(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.Test(t)
	s := newStore()
	batchSize := 0
	rp, err := newReportingPlugin(s, nil, batchSize, ocr3types.ReportingPluginConfig{}, lggr)
	require.NoError(t, err)

	outcomeCtx := ocr3types.OutcomeContext{
		PreviousOutcome: []byte(""),
	}
	_, err = rp.Query(ctx, outcomeCtx)
	assert.Error(t, err)
}

func TestReportingPlugin_Query(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.Test(t)
	s := newStore()
	rp, err := newReportingPlugin(s, nil, defaultBatchSize, ocr3types.ReportingPluginConfig{}, lggr)
	require.NoError(t, err)

	eid := uuid.New().String()
	err = s.add(ctx, &request{
		WorkflowID:          workflowTestID,
		WorkflowExecutionID: eid,
	})
	require.NoError(t, err)
	outcomeCtx := ocr3types.OutcomeContext{
		PreviousOutcome: []byte(""),
	}

	q, err := rp.Query(ctx, outcomeCtx)
	require.NoError(t, err)

	qry := &pbtypes.Query{}
	err = proto.Unmarshal(q, qry)
	require.NoError(t, err)

	assert.Len(t, qry.Ids, 1)
	assert.Equal(t, qry.Ids[0].WorkflowId, workflowTestID)
	assert.Equal(t, qry.Ids[0].WorkflowExecutionId, eid)
}

func TestReportingPlugin_Observation(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.Test(t)
	s := newStore()
	rp, err := newReportingPlugin(s, nil, defaultBatchSize, ocr3types.ReportingPluginConfig{}, lggr)
	require.NoError(t, err)

	o, err := values.NewList([]any{"hello"})
	require.NoError(t, err)

	eid := uuid.New().String()
	err = s.add(ctx, &request{
		WorkflowID:          workflowTestID,
		WorkflowExecutionID: eid,
		Observations:        o,
	})
	require.NoError(t, err)
	outcomeCtx := ocr3types.OutcomeContext{
		PreviousOutcome: []byte(""),
	}

	q, err := rp.Query(ctx, outcomeCtx)
	require.NoError(t, err)

	obs, err := rp.Observation(ctx, outcomeCtx, q)
	require.NoError(t, err)

	obspb := &pbtypes.Observations{}
	err = proto.Unmarshal(obs, obspb)
	require.NoError(t, err)

	assert.Len(t, obspb.Observations, 1)
	fo := obspb.Observations[0]
	assert.Equal(t, fo.Id.WorkflowExecutionId, eid)
	assert.Equal(t, fo.Id.WorkflowId, workflowTestID)
	assert.Equal(t, o, values.FromListValueProto(fo.Observations))
}

func TestReportingPlugin_Observation_NoResults(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.Test(t)
	s := newStore()
	rp, err := newReportingPlugin(s, nil, defaultBatchSize, ocr3types.ReportingPluginConfig{}, lggr)
	require.NoError(t, err)

	outcomeCtx := ocr3types.OutcomeContext{
		PreviousOutcome: []byte(""),
	}

	q, err := rp.Query(ctx, outcomeCtx)
	require.NoError(t, err)

	obs, err := rp.Observation(ctx, outcomeCtx, q)
	require.NoError(t, err)

	obspb := &pbtypes.Observations{}
	err = proto.Unmarshal(obs, obspb)
	require.NoError(t, err)

	assert.Len(t, obspb.Observations, 0)
}

type mockCapability struct {
	gotResponse *outputs
	aggregator  *aggregator
	encoder     *enc
}

func (mc *mockCapability) transmitResponse(ctx context.Context, resp *outputs) error {
	mc.gotResponse = resp
	return nil
}

type aggregator struct {
	gotObs  map[commontypes.OracleID][]values.Value
	outcome *pbtypes.AggregationOutcome
}

func (a *aggregator) Aggregate(pout *pbtypes.AggregationOutcome, observations map[commontypes.OracleID][]values.Value, _ int) (*pbtypes.AggregationOutcome, error) {
	a.gotObs = observations
	nm, err := values.NewMap(
		map[string]any{
			"aggregated": "outcome",
		},
	)
	if err != nil {
		return nil, err
	}
	a.outcome = &pbtypes.AggregationOutcome{
		EncodableOutcome: values.Proto(nm).GetMapValue(),
	}
	return a.outcome, nil
}

type enc struct {
	gotInput values.Map
}

func (e *enc) Encode(ctx context.Context, input values.Map) ([]byte, error) {
	e.gotInput = input
	return proto.Marshal(values.Proto(&input))
}

func (mc *mockCapability) getAggregator(workflowID string) (pbtypes.Aggregator, error) {
	return mc.aggregator, nil
}

func (mc *mockCapability) getEncoder(workflowID string) (pbtypes.Encoder, error) {
	return mc.encoder, nil
}

func TestReportingPlugin_Outcome(t *testing.T) {
	lggr := logger.Test(t)
	s := newStore()
	cap := &mockCapability{
		aggregator: &aggregator{},
		encoder:    &enc{},
	}
	rp, err := newReportingPlugin(s, cap, defaultBatchSize, ocr3types.ReportingPluginConfig{}, lggr)
	require.NoError(t, err)

	weid := uuid.New().String()
	id := &pbtypes.Id{
		WorkflowExecutionId: weid,
		WorkflowId:          workflowTestID,
	}
	q := &pbtypes.Query{
		Ids: []*pbtypes.Id{id},
	}
	qb, err := proto.Marshal(q)
	require.NoError(t, err)
	o, err := values.NewList([]any{"hello"})
	require.NoError(t, err)
	obs := &pbtypes.Observations{
		Observations: []*pbtypes.Observation{
			{
				Id:           id,
				Observations: values.Proto(o).GetListValue(),
			},
		},
	}

	rawObs, err := proto.Marshal(obs)
	require.NoError(t, err)
	aos := []types.AttributedObservation{
		{
			Observation: rawObs,
			Observer:    commontypes.OracleID(1),
		},
	}

	outcome, err := rp.Outcome(ocr3types.OutcomeContext{}, qb, aos)
	require.NoError(t, err)

	opb := &pbtypes.Outcome{}
	err = proto.Unmarshal(outcome, opb)
	require.NoError(t, err)

	assert.Len(t, opb.CurrentReports, 1)

	cr := opb.CurrentReports[0]
	assert.EqualExportedValues(t, cr.Id, id)
	assert.EqualExportedValues(t, cr.Outcome, cap.aggregator.outcome)
	assert.EqualExportedValues(t, opb.Outcomes[workflowTestID], cap.aggregator.outcome)
}

func TestReportingPlugin_Reports_ShouldReportFalse(t *testing.T) {
	lggr := logger.Test(t)
	s := newStore()
	cap := &mockCapability{
		aggregator: &aggregator{},
		encoder:    &enc{},
	}
	rp, err := newReportingPlugin(s, cap, defaultBatchSize, ocr3types.ReportingPluginConfig{}, lggr)
	require.NoError(t, err)

	var sqNr uint64
	weid := uuid.New().String()
	id := &pbtypes.Id{
		WorkflowExecutionId: weid,
		WorkflowId:          workflowTestID,
	}
	nm, err := values.NewMap(
		map[string]any{
			"our": "aggregation",
		},
	)
	require.NoError(t, err)
	outcome := &pbtypes.Outcome{
		CurrentReports: []*pbtypes.Report{
			{
				Id: id,
				Outcome: &pbtypes.AggregationOutcome{
					EncodableOutcome: values.Proto(nm).GetMapValue(),
				},
			},
		},
	}
	pl, err := proto.Marshal(outcome)
	require.NoError(t, err)
	reports, err := rp.Reports(sqNr, pl)
	require.NoError(t, err)

	assert.Len(t, reports, 1)
	gotRep := reports[0]
	assert.Len(t, gotRep.Report, 0)

	ib := gotRep.Info
	info := &pbtypes.ReportInfo{}
	err = proto.Unmarshal(ib, info)
	require.NoError(t, err)

	assert.EqualExportedValues(t, info.Id, id)
	assert.False(t, info.ShouldReport)
}

func TestReportingPlugin_Reports_ShouldReportTrue(t *testing.T) {
	lggr := logger.Test(t)
	s := newStore()
	cap := &mockCapability{
		aggregator: &aggregator{},
		encoder:    &enc{},
	}
	rp, err := newReportingPlugin(s, cap, defaultBatchSize, ocr3types.ReportingPluginConfig{}, lggr)
	require.NoError(t, err)

	var sqNr uint64
	weid := uuid.New().String()
	id := &pbtypes.Id{
		WorkflowExecutionId: weid,
		WorkflowId:          workflowTestID,
	}
	nm, err := values.NewMap(
		map[string]any{
			"our": "aggregation",
		},
	)
	nmp := values.Proto(nm).GetMapValue()
	require.NoError(t, err)
	outcome := &pbtypes.Outcome{
		CurrentReports: []*pbtypes.Report{
			{
				Id: id,
				Outcome: &pbtypes.AggregationOutcome{
					EncodableOutcome: nmp,
					ShouldReport:     true,
				},
			},
		},
	}
	pl, err := proto.Marshal(outcome)
	require.NoError(t, err)
	reports, err := rp.Reports(sqNr, pl)
	require.NoError(t, err)

	assert.Len(t, reports, 1)
	gotRep := reports[0]

	rep := &pb.Value{}
	err = proto.Unmarshal(gotRep.Report, rep)
	require.NoError(t, err)

	// The workflow ID and execution ID get added to the report.
	nm.Underlying[pbtypes.WorkflowIDFieldName] = values.NewString(workflowTestID)
	nm.Underlying[pbtypes.ExecutionIDFieldName] = values.NewString(weid)
	fp := values.FromProto(rep)
	assert.Equal(t, nm, fp)

	ib := gotRep.Info
	info := &pbtypes.ReportInfo{}
	err = proto.Unmarshal(ib, info)
	require.NoError(t, err)

	assert.EqualExportedValues(t, info.Id, id)
	assert.True(t, info.ShouldReport)
}
