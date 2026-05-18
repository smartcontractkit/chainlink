package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidWorkflowID(t *testing.T) {
	require.NotNil(t, ValidateWorkflowOrExecutionID("too_short"))
	require.NotNil(t, ValidateWorkflowOrExecutionID("nothex--95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"))
	require.NoError(t, ValidateWorkflowOrExecutionID("15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"))
}

func TestIsValidTriggerEventID(t *testing.T) {
	require.False(t, IsValidID(""))
	require.False(t, IsValidID("\n\n"))
	require.True(t, IsValidID("id_id_2"))
}
