package chainlink

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestFilePersistedSecretGenerator(t *testing.T) {
	t.Parallel()
	tests.BelongsToCISuite(t, "with-db")
	rootDir := t.TempDir()
	var secretGenerator FilePersistedSecretGenerator

	initial, err := secretGenerator.Generate(rootDir)
	require.NoError(t, err)
	require.NotEmpty(t, initial)
	require.NotEqual(t, "clsession_test_secret", initial)

	second, err := secretGenerator.Generate(rootDir)
	require.NoError(t, err)
	require.Equal(t, initial, second)
}
