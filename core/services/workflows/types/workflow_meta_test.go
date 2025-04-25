package types

import (
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

func TestWorkflowID_FromHex(t *testing.T) {
	t.Parallel()
	_, err := WorkflowIDFromHex("wrong chars")
	require.Error(t, err)

	_, err = WorkflowIDFromHex("00112233")
	require.Error(t, err) // wrong length

	_, err = NewWorkflowName("aabbccddeeff00112233aabbccddeeff00112233aabbccddeeff00112233ffff")
	require.NoError(t, err)
}

func TestValidateWorkflowOwner(t *testing.T) {
	t.Parallel()
	require.Error(t, ValidateWorkflowOwner("wrong chars"))
	require.Error(t, ValidateWorkflowOwner("00112233")) // wrong length
	require.NoError(t, ValidateWorkflowOwner("aabbccddeeff00112233aabbccddeeff00112233"))
}
