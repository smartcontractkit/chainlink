package smoke

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomFlaky_JUST_FOR_TESTING_FLAKEGUARD(t *testing.T) {
	require.True(t, false, "This test is just for testing the flakeguard")
}
