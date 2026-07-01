package vaultutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

func TestBuildWorkflowGetSecretsRequestID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		md   capabilities.RequestMetadata
		want string
	}{
		{
			name: "execution path",
			md: capabilities.RequestMetadata{
				WorkflowID:          "wf-1",
				WorkflowExecutionID: "abc123sha",
				ReferenceID:         "42",
			},
			want: "wf-1::abc123sha::42",
		},
		{
			name: "subscription path",
			md: capabilities.RequestMetadata{
				WorkflowID:  "wf-1",
				ReferenceID: "7",
			},
			want: "wf-1::subscription::7",
		},
		{
			name: "gate off empty workflow ID",
			md: capabilities.RequestMetadata{
				WorkflowExecutionID: "abc123sha",
				ReferenceID:         "42",
			},
			want: "::abc123sha::42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, BuildWorkflowGetSecretsRequestID(tt.md))
		})
	}
}
