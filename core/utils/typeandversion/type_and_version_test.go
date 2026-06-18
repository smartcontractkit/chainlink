package typeandversion_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/utils/typeandversion"
)

func TestParseTypeAndVersion(t *testing.T) {
	t.Parallel()

	t.Run("valid string", func(t *testing.T) {
		t.Parallel()
		contractType, version, err := typeandversion.ParseTypeAndVersion("SomeType 1.2.0")
		require.NoError(t, err)
		require.Equal(t, "SomeType", contractType)
		require.Equal(t, "1.2.0", version)
	})

	t.Run("invalid string - too short", func(t *testing.T) {
		t.Parallel()
		_, _, err := typeandversion.ParseTypeAndVersion("v1.2.0")
		require.ErrorContains(t, err, "invalid type and version v1.2.0")
	})

	t.Run("invalid string - too long", func(t *testing.T) {
		t.Parallel()
		_, _, err := typeandversion.ParseTypeAndVersion("SomeType WithMoreWords vv1.2.0")
		require.ErrorContains(t, err, "invalid type and version SomeType WithMoreWords vv1.2.0")
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		contractType, version, err := typeandversion.ParseTypeAndVersion("")
		require.NoError(t, err)
		require.Equal(t, typeandversion.UnknownContractType, contractType)
		require.Equal(t, typeandversion.DefaultVersion, version)
	})
}
