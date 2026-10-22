package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkflowName_Validate(t *testing.T) {
	t.Parallel()
	_, err := NewWorkflowName("")
	require.Error(t, err)

	_, err = NewWorkflowName(string(make([]byte, maxWorkflowNameLength+1)))
	require.Error(t, err)

	_, err = NewWorkflowName("a")
	require.NoError(t, err)
}

// TestWorkflowID_FromHex covers parsing a workflow ID. Only the canonical
// (i.e., lowercase, no "0x" prefix) workflow ID should match the value returned
// by calling WorkflowID.Hex().
func TestWorkflowID_FromHex(t *testing.T) {
	t.Parallel()

	const canonical = "aabbccddeeff00112233aabbccddeeff00112233aabbccddeeff00112233ffff"

	tests := []struct {
		name    string
		h       string
		wantErr bool
	}{
		{name: "canonical lowercase hex", h: canonical},
		{name: "uppercase hex", h: strings.ToUpper(canonical), wantErr: true},
		{name: "mixed case hex", h: "AAbbccddeeff00112233aabbccddeeff00112233aabbccddeeff00112233ffff", wantErr: true},
		{name: "0x prefix", h: "0x" + canonical, wantErr: true},
		{name: "wrong chars", h: "wrong chars", wantErr: true},
		{name: "wrong length", h: "00112233", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := WorkflowIDFromHex(tc.h)
			if tc.wantErr {
				require.Error(t, err, "WorkflowIDFromHex(%q)", tc.h)
				return
			}
			require.NoError(t, err, "WorkflowIDFromHex(%q)", tc.h)
			require.Equal(t, tc.h, got.Hex(), "WorkflowIDFromHex(%q).Hex()", tc.h)
		})
	}
}

func TestValidateWorkflowOwner(t *testing.T) {
	t.Parallel()
	require.Error(t, ValidateWorkflowOwner("wrong chars"))
	require.Error(t, ValidateWorkflowOwner("00112233")) // wrong length
	require.NoError(t, ValidateWorkflowOwner("aabbccddeeff00112233aabbccddeeff00112233"))
}
