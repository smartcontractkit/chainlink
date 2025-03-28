package smoke

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRandomFlaky_JUST_FOR_TESTING_FLAKEGUARD(t *testing.T) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 9 (inclusive)
	randomValue := rand.Intn(10)

	// Fails when randomValue is 0, 1, or 2
	require.True(t, randomValue >= 5, "Got value: %d", randomValue)
}
