package chainlink

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilePersistedSecretGenerator(t *testing.T) {
	t.Parallel()
	rootDir := t.TempDir()
	var secretGenerator FilePersistedSecretGenerator

	initial, err := secretGenerator.Generate(rootDir)
	require.NoError(t, err)
	require.NotEqual(t, "", initial)
	require.NotEqual(t, "clsession_test_secret", initial)

	second, err := secretGenerator.Generate(rootDir)
	require.NoError(t, err)
	require.Equal(t, initial, second)
}
