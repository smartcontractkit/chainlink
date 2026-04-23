package capabilities

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestMetadataRegistry_LocalNode_UsesConfiguredWorkflowDONF(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		registry     TestMetadataRegistry
		expectedDonF uint8
	}{
		{
			name:         "default workflow DON fault tolerance",
			registry:     TestMetadataRegistry{},
			expectedDonF: 0,
		},
		{
			name:         "mock trigger workflow DON fault tolerance",
			registry:     TestMetadataRegistry{WorkflowDONF: 1},
			expectedDonF: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			node, err := tt.registry.LocalNode(context.Background())
			require.NoError(t, err)
			require.Equal(t, tt.expectedDonF, node.WorkflowDON.F)
		})
	}
}
