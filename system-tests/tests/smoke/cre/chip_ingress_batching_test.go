package cre

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/stretchr/testify/require"
)

func ExecuteChipIngressBatchingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()

	testLogger := framework.L
	const workflowFileLocation = "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	var publishCalls atomic.Int64
	var publishBatchCalls atomic.Int64
	var publishBatchEventCount atomic.Int64

	basePublishFn := t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh)

	// PublishFunc handles individual Publish RPCs.
	publishFn := func(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
		publishCalls.Add(1)
		return basePublishFn(ctx, event)
	}

	// PublishBatchFunc handles PublishBatch RPCs directly, then fans out to basePublishFn per event.
	publishBatchFn := func(ctx context.Context, batch *chippb.CloudEventBatch) (*chippb.PublishResponse, error) {
		publishBatchCalls.Add(1)
		publishBatchEventCount.Add(int64(len(batch.Events)))

		for _, event := range batch.Events {
			if _, err := basePublishFn(ctx, event); err != nil {
				return nil, fmt.Errorf("publish batch: event %s: %w", event.GetId(), err)
			}
		}
		return &chippb.PublishResponse{}, nil
	}

	server := t_helpers.StartChipTestSink(t, publishFn, t_helpers.WithPublishBatchFunc(publishBatchFn))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	const workflowCount = 3
	workflowIDs := make([]string, 0, workflowCount)
	for i := range workflowCount {
		workflowName := fmt.Sprintf("chip-ingress-batching-%d-%d", i, time.Now().UnixNano())
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &crontypes.WorkflowConfig{
			Schedule: "*/30 * * * * *",
		}, workflowFileLocation)
		workflowIDs = append(workflowIDs, workflowID)
	}

	for _, workflowID := range workflowIDs {
		t_helpers.WatchWorkflowLogs(
			t,
			testLogger,
			userLogsCh,
			baseMessageCh,
			t_helpers.WorkflowEngineInitErrorLog,
			"Amazing workflow user log",
			4*time.Minute,
			t_helpers.WithUserLogWorkflowID(workflowID),
		)
	}

	// Verify batching occurred.
	batchCalls := publishBatchCalls.Load()
	batchEvents := publishBatchEventCount.Load()
	singleCalls := publishCalls.Load()

	require.Greater(t, batchCalls, int64(0),
		"expected at least one PublishBatch RPC call, got 0 (single Publish calls: %d)", singleCalls)

	testLogger.Info().
		Int64("publish_calls", singleCalls).
		Int64("publish_batch_calls", batchCalls).
		Int64("publish_batch_events", batchEvents).
		Msg("CHiP ingress batching verified")
}
