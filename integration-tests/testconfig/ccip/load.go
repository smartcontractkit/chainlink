package ccip

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type LoadConfig struct {
	LoadDuration              *string
	LokiEndpoint              *string
	MessageTypeWeights        *[]int
	RequestFrequency          *string
	EnabledDestionationChains *[]uint64
	CribEnvDirectory          *string
	TestDuration              *string
}

func (l *LoadConfig) Validate(t *testing.T) {
	ld, err := time.ParseDuration(*l.LoadDuration)
	require.NoError(t, err, "LoadDuration must be a valid duration")

	td, err := time.ParseDuration(*l.TestDuration)
	require.NoError(t, err, "TestDuration must be a valid duration")

	require.GreaterOrEqual(t, td, ld, "TestDuration must be greater than or equal to LoadDuration")

	agg := 0
	for _, w := range *l.MessageTypeWeights {
		agg += w
	}
	require.Equal(t, 100, agg, "Sum of MessageTypeWeights must be 100")
}

func (l *LoadConfig) GetLoadDuration() time.Duration {
	ld, _ := time.ParseDuration(*l.LoadDuration)
	return ld
}

func (l *LoadConfig) GetTestDuration() time.Duration {
	td, _ := time.ParseDuration(*l.TestDuration)
	return td
}
