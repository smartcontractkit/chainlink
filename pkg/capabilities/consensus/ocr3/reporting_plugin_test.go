package ocr3

import (
	"testing"

	"github.com/google/uuid"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	pbtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
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

	o := values.NewString("hello")
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

	_, err = rp.Observation(ctx, outcomeCtx, q)
	require.NoError(t, err)

	// obs := &pbtypes.Observations{
	// Id: &pbtypes.Id{
	// }
	// err = proto.Unmarshal(obspb, obs)
	// require.NoError(t, err)

	// fmt.Printf("%+v", obs)

	// assert.NotNil(t, nil)
}
