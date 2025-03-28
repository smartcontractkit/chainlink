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

	// Generate a random number between 0 and 9
	// If the number is < 6, the test will fail (60% failure rate)
	randomValue := rand.Intn(10)

	require.True(t, randomValue >= 3, "This test is designed to fail 30%% of the time. Got value: %d", randomValue)
}
