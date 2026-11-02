package v1_6

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHybridGroupUpdateToContractUpdate(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		group uint8
		valid bool
	}{
		{"LOCK_AND_RELEASE", 0, true},
		{"BURN_AND_MNIT", 0, false},
		{"BURN_AND_MINT", 1, true},
		{"0", 0, false},
		{"1", 0, false},
		{"", 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := HybridGroupUpdate{
				RemoteChainSelector: 2308837218439511688,
				RemoteChainSupply:   big.NewInt(123),
				Group:               tc.name,
			}
			update, err := input.toContractUpdate()
			if !tc.valid {
				require.ErrorContains(t, err, "must be LOCK_AND_RELEASE or BURN_AND_MINT")
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.group, update.Group)
			require.Equal(t, input.RemoteChainSelector, update.RemoteChainSelector)
			require.Equal(t, input.RemoteChainSupply, update.RemoteChainSupply)
		})
	}
}

func TestHybridGroupConfigRejectsNumericGroups(t *testing.T) {
	t.Parallel()
	for _, payload := range []string{
		`{"updates":[{"group":0}]}`,
		`{"updates":[{"group":1}]}`,
	} {
		var config HybridGroupConfig
		require.Error(t, json.Unmarshal([]byte(payload), &config))
	}
}
