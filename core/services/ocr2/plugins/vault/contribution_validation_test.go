package vault

import (
	"testing"

	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

func TestCombineObservationErrors(t *testing.T) {
	t.Parallel()
	t.Run("empty observations uses fallback", func(t *testing.T) {
		t.Parallel()
		result := consensusObservationError([]*vaultcommon.Observation{}, 0)
		require.Equal(t, "request is not valid", result)
	})

	t.Run("single observation reaches consensus with f=0", func(t *testing.T) {
		t.Parallel()
		obs := []*vaultcommon.Observation{
			{Error: &vaultcommon.ObservationError{Message: "node error"}},
		}
		result := consensusObservationError(obs, 0)
		require.Equal(t, "node error", result)
	})

	t.Run("multiple observations reach consensus with f+1 agreement", func(t *testing.T) {
		t.Parallel()
		f := 1
		obs := []*vaultcommon.Observation{
			{Error: &vaultcommon.ObservationError{Message: "timeout error"}},
			{Error: &vaultcommon.ObservationError{Message: "timeout error"}},
			{Error: &vaultcommon.ObservationError{Message: "different error"}},
		}
		result := consensusObservationError(obs, f)
		require.Equal(t, "timeout error", result)
	})

	t.Run("returns fallback when no error reaches f+1 consensus", func(t *testing.T) {
		t.Parallel()
		f := 2
		obs := []*vaultcommon.Observation{
			{Error: &vaultcommon.ObservationError{Message: "error1"}},
			{Error: &vaultcommon.ObservationError{Message: "error2"}},
			{Error: &vaultcommon.ObservationError{Message: "error3"}},
		}
		result := consensusObservationError(obs, f)
		require.Equal(t, "request is not valid", result)
	})

	t.Run("ignores empty error messages", func(t *testing.T) {
		t.Parallel()
		f := 0
		obs := []*vaultcommon.Observation{
			{Error: &vaultcommon.ObservationError{Message: ""}},
			{Error: &vaultcommon.ObservationError{Message: "valid error"}},
			{Error: &vaultcommon.ObservationError{Message: ""}},
		}
		result := consensusObservationError(obs, f)
		require.Equal(t, "valid error", result)
	})

	t.Run("picks most frequent error when multiple reach consensus", func(t *testing.T) {
		t.Parallel()
		f := 1
		obs := []*vaultcommon.Observation{
			{Error: &vaultcommon.ObservationError{Message: "error1"}},
			{Error: &vaultcommon.ObservationError{Message: "error1"}},
			{Error: &vaultcommon.ObservationError{Message: "error2"}},
			{Error: &vaultcommon.ObservationError{Message: "error2"}},
			{Error: &vaultcommon.ObservationError{Message: "error2"}},
		}
		result := consensusObservationError(obs, f)
		require.Equal(t, "error2", result)
	})

	t.Run("requires exactly f+1 for consensus", func(t *testing.T) {
		t.Parallel()
		f := 2
		obs := []*vaultcommon.Observation{
			{Error: &vaultcommon.ObservationError{Message: "error"}},
			{Error: &vaultcommon.ObservationError{Message: "error"}},
			{Error: &vaultcommon.ObservationError{Message: "other"}},
		}
		result := consensusObservationError(obs, f)
		require.Equal(t, "request is not valid", result)
	})

	t.Run("reaches consensus with exactly f+1 matching errors", func(t *testing.T) {
		t.Parallel()
		f := 2
		obs := []*vaultcommon.Observation{
			{Error: &vaultcommon.ObservationError{Message: "consensus error"}},
			{Error: &vaultcommon.ObservationError{Message: "consensus error"}},
			{Error: &vaultcommon.ObservationError{Message: "consensus error"}},
			{Error: &vaultcommon.ObservationError{Message: "other"}},
		}
		result := consensusObservationError(obs, f)
		require.Equal(t, "consensus error", result)
	})

	t.Run("all observations with same error reaches consensus", func(t *testing.T) {
		t.Parallel()
		f := 4
		obs := []*vaultcommon.Observation{
			{Error: &vaultcommon.ObservationError{Message: "all same"}},
			{Error: &vaultcommon.ObservationError{Message: "all same"}},
			{Error: &vaultcommon.ObservationError{Message: "all same"}},
			{Error: &vaultcommon.ObservationError{Message: "all same"}},
			{Error: &vaultcommon.ObservationError{Message: "all same"}},
		}
		result := consensusObservationError(obs, f)
		require.Equal(t, "all same", result)
	})
}
