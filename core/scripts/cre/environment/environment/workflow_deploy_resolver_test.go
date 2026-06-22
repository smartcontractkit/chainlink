package environment

import (
	"testing"

	"github.com/stretchr/testify/require"

	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func TestFinalizeWorkflowDonFamily(t *testing.T) {
	t.Parallel()

	t.Run("explicit family", func(t *testing.T) {
		t.Parallel()
		family, err := finalizeWorkflowDonFamily("feeds-zone-a", true)
		require.NoError(t, err)
		require.Equal(t, "feeds-zone-a", family)
	})

	t.Run("legacy default", func(t *testing.T) {
		t.Parallel()
		family, err := finalizeWorkflowDonFamily("", false)
		require.NoError(t, err)
		require.Equal(t, envconfig.DefaultDONFamily, family)
	})

	t.Run("pairing requires family", func(t *testing.T) {
		t.Parallel()
		_, err := finalizeWorkflowDonFamily("", true)
		require.Error(t, err)
	})
}
