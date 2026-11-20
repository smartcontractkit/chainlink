package v2

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Synthetic identifiers use production encodings.
const (
	metadataWorkflowIDFixture       = "0020350fb4f89856a031268826d6e9b9340996637aa27859baa7b07db58bb020"
	metadataExecutionIDFixture      = "9f38bdece8892b3452c401a7a2d5a170f078ab9b2c98c591e24028ce94b33ed2"
	metadataOwnerFixture            = "5366890c26e8737b89d9084b4c632f7de73d1674"
	metadataOtherWorkflowIDFixture  = "0052457bc477e58cad894372b8a80d1f0b03322c8a070d93bcf9bbd1f21cbd09"
	metadataOtherExecutionIDFixture = "e48ddcb9c9fdca45b68f7d37d1dc8bcb29df3777b3ece2f01c14d54872d44f01"
)

func TestEngine_ExecutionMetadata(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name        string
		executionID string
		f           uint8
	}{
		{name: "trigger", executionID: metadataExecutionIDFixture, f: 3},
		{name: "subscribe", f: 3},
		{name: "zero fault tolerance", executionID: metadataExecutionIDFixture},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup.
			engine := metadataEngineFixture()
			engine.localNode.Store(&capabilities.Node{WorkflowDON: capabilities.DON{ID: 17, F: tc.f}})

			// Execute.
			executionMetadata := engine.getExecutionMetadata(tc.executionID)

			// Verify.
			want := &sdk.WorkflowExecutionMetadata{
				WorkflowId:    metadataWorkflowIDFixture,
				ExecutionId:   tc.executionID,
				WorkflowOwner: metadataOwnerFixture,
				DonId:         17,
				DonF:          uint32(tc.f),
			}
			require.True(t, proto.Equal(want, executionMetadata), "unexpected execution metadata: %v", executionMetadata)
		})
	}
}

func TestEngine_ExecutionMetadataSnapshot(t *testing.T) {
	t.Parallel()

	// Setup.
	engine := metadataEngineFixture()
	first := engine.getExecutionMetadata(metadataExecutionIDFixture)
	wantFirst := proto.Clone(first)

	// Refresh the registry snapshot.
	engine.localNode.Store(&capabilities.Node{WorkflowDON: capabilities.DON{ID: 19, F: 0, ConfigVersion: 2}})
	second := engine.getExecutionMetadata(metadataExecutionIDFixture)

	// Verify both snapshots.
	wantSecond := &sdk.WorkflowExecutionMetadata{
		WorkflowId:    metadataWorkflowIDFixture,
		ExecutionId:   metadataExecutionIDFixture,
		WorkflowOwner: metadataOwnerFixture,
		DonId:         19,
	}
	require.True(t, proto.Equal(wantSecond, second), "unexpected execution metadata: %v", second)
	require.True(t, proto.Equal(wantFirst, first))
}

func TestEngine_ExecutionMetadataReturnsIndependentValues(t *testing.T) {
	t.Parallel()

	// Setup.
	engine := metadataEngineFixture()
	first := engine.getExecutionMetadata("")
	want := proto.Clone(first)

	// Mutate a returned value.
	first.WorkflowId = metadataOtherWorkflowIDFixture
	first.ExecutionId = metadataExecutionIDFixture
	first.DonF = 0
	executionMetadata := engine.getExecutionMetadata("")

	// Verify engine state is unchanged.
	require.True(t, proto.Equal(want, executionMetadata))
}

func TestEngine_ExecutionMetadataDoesNotReuseExecutionID(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name        string
		executionID string
	}{
		{name: "next trigger", executionID: metadataOtherExecutionIDFixture},
		{name: "subscribe"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup a previous execution.
			engine := metadataEngineFixture()
			previous := engine.getExecutionMetadata(metadataExecutionIDFixture)
			want := proto.Clone(previous).(*sdk.WorkflowExecutionMetadata)
			want.ExecutionId = tc.executionID

			// Build metadata for the next invocation.
			executionMetadata := engine.getExecutionMetadata(tc.executionID)

			// Verify execution IDs are independent.
			require.True(t, proto.Equal(want, executionMetadata), "unexpected execution metadata: %v", executionMetadata)
			require.Equal(t, metadataExecutionIDFixture, previous.ExecutionId)
		})
	}
}

func metadataEngineFixture() *baseEngine {
	engine := &baseEngine{cfg: &EngineConfig{
		WorkflowID:    metadataWorkflowIDFixture,
		WorkflowOwner: metadataOwnerFixture,
	}}
	engine.localNode.Store(&capabilities.Node{WorkflowDON: capabilities.DON{ID: 17, F: 3, ConfigVersion: 1}})
	return engine
}

// The metadata is part of the ExecuteRequest bytes that confidential workflows hash across nodes,
// so it is gated by an activation window checked against the execution timestamp: it must stay off
// by default and switch on at the same execution for every node.
func TestEngine_ExecutionMetadataGate(t *testing.T) {
	t.Parallel()

	window := func(from, to time.Time) func(*cresettings.Workflows) {
		return func(w *cresettings.Workflows) {
			w.FeatureExecutionMetadataActivePeriod.DefaultValue = settings.TimeRange(from, to).DefaultValue
		}
	}
	start := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 12, 2, 0, 0, 0, 0, time.UTC)

	for _, tc := range []struct {
		name  string
		cfgFn func(*cresettings.Workflows)
		at    time.Time
		want  bool
	}{
		{name: "default window is off", at: start},
		{name: "inside the window", cfgFn: window(start, end), at: start.Add(time.Hour), want: true},
		{name: "before the window", cfgFn: window(start, end), at: start.Add(-time.Second)},
		{name: "after the window", cfgFn: window(start, end), at: end.Add(time.Second)},
		{name: "no execution timestamp stays off", cfgFn: window(start, end), at: time.Time{}},
		{name: "window reaching year 1 covers a zero timestamp", cfgFn: window(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), end), at: time.Time{}, want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup.
			engine := metadataEngineFixture()
			limiters, err := NewLimiters(limits.Factory{Logger: logger.TestLogger(t)}, tc.cfgFn)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, limiters.Close()) })
			engine.cfg.LocalLimiters = limiters
			ctx := contexts.WithCRE(context.Background(), contexts.CRE{Workflow: metadataWorkflowIDFixture})

			// Execute.
			got := engine.executionMetadata(ctx, tc.at, metadataExecutionIDFixture)

			// Verify.
			if !tc.want {
				require.Nil(t, got)
				return
			}
			require.True(t, proto.Equal(engine.getExecutionMetadata(metadataExecutionIDFixture), got), "unexpected execution metadata: %v", got)
		})
	}
}
