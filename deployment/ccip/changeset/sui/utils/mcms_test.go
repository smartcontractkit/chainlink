package utils

import (
	"testing"

	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSuiRoleFromAction(t *testing.T) {
	tests := []struct {
		name      string
		action    mcmstypes.TimelockAction
		expected  suisdk.TimelockRole
		expectErr bool
	}{
		{"schedule maps to proposer", mcmstypes.TimelockActionSchedule, suisdk.TimelockRoleProposer, false},
		{"bypass maps to bypasser", mcmstypes.TimelockActionBypass, suisdk.TimelockRoleBypasser, false},
		{"cancel maps to canceller", mcmstypes.TimelockActionCancel, suisdk.TimelockRoleCanceller, false},
		{"empty defaults to proposer", "", suisdk.TimelockRoleProposer, false},
		{"unsupported action errors", mcmstypes.TimelockAction("invalid"), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := GetSuiRoleFromAction(tt.action)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, role)
			}
		})
	}
}
