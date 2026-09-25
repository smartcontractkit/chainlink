package triggers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validWorkflowID is a syntactically valid workflow ID: a 64-character hex string (32 bytes).
var validWorkflowID = strings.Repeat("ab", 32)

func TestRegistrationID(t *testing.T) {
	tests := []struct {
		name         string
		workflowID   string
		triggerIndex int
		want         string
	}{
		{
			name:         "simple workflow ID",
			workflowID:   "abc123",
			triggerIndex: 0,
			want:         "trigger_reg_abc123_0",
		},
		{
			name:         "workflow ID containing underscores",
			workflowID:   "wf_with_underscores",
			triggerIndex: 5,
			want:         "trigger_reg_wf_with_underscores_5",
		},
		{
			name:         "empty workflow ID",
			workflowID:   "",
			triggerIndex: 2,
			want:         "trigger_reg__2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RegistrationID(tt.workflowID, tt.triggerIndex)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseWorkflowID(t *testing.T) {
	tests := []struct {
		name             string
		registrationID   string
		wantWorkflowID   string
		wantErrSubstring string
	}{
		{
			name:           "valid workflow ID",
			registrationID: "trigger_reg_" + validWorkflowID + "_0",
			wantWorkflowID: validWorkflowID,
		},
		{
			name:           "valid workflow ID with large trigger index",
			registrationID: "trigger_reg_" + validWorkflowID + "_12345",
			wantWorkflowID: validWorkflowID,
		},
		{
			name:             "missing prefix",
			registrationID:   "not_a_registration_id_0",
			wantErrSubstring: "missing prefix",
		},
		{
			name:             "missing trigger index",
			registrationID:   "trigger_reg_" + validWorkflowID,
			wantErrSubstring: "missing trigger index",
		},
		{
			name:             "non-numeric trigger index",
			registrationID:   "trigger_reg_" + validWorkflowID + "_notanumber",
			wantErrSubstring: "invalid trigger index",
		},
		{
			name:             "empty workflow ID",
			registrationID:   "trigger_reg__2",
			wantErrSubstring: "invalid workflow ID",
		},
		{
			name:             "workflow ID too short",
			registrationID:   "trigger_reg_abc123_0",
			wantErrSubstring: "invalid workflow ID",
		},
		{
			name:             "workflow ID not hex-encoded",
			registrationID:   "trigger_reg_" + strings.Repeat("z", 64) + "_0",
			wantErrSubstring: "invalid workflow ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseWorkflowID(tt.registrationID)
			if tt.wantErrSubstring != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrSubstring)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantWorkflowID, got)
		})
	}
}

func TestRegistrationID_ParseWorkflowID_RoundTrip(t *testing.T) {
	tests := []struct {
		name         string
		workflowID   string
		triggerIndex int
	}{
		{name: "trigger index zero", workflowID: validWorkflowID, triggerIndex: 0},
		{name: "large trigger index", workflowID: validWorkflowID, triggerIndex: 12345},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registrationID := RegistrationID(tt.workflowID, tt.triggerIndex)

			got, err := ParseWorkflowID(registrationID)
			require.NoError(t, err)
			assert.Equal(t, tt.workflowID, got)
		})
	}
}
