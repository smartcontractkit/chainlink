package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
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
