package ccip

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/stretchr/testify/require"
)

type LoadConfig struct {
	LoadDuration         *string
	MessageTypeWeights   *[]int
	RequestFrequency     *string
	CribEnvDirectory     *string
	NumDestinationChains *int
	TimeoutDuration      *string
}

func (l *LoadConfig) Validate(t *testing.T, e *deployment.Environment) {
	_, err := time.ParseDuration(*l.LoadDuration)
	require.NoError(t, err, "LoadDuration must be a valid duration")

	_, err = time.ParseDuration(*l.TimeoutDuration)
	require.NoError(t, err, "TimeoutDuration must be a valid duration")

	agg := 0
	for _, w := range *l.MessageTypeWeights {
		agg += w
	}
	require.Equal(t, 100, agg, "Sum of MessageTypeWeights must be 100")

	require.GreaterOrEqual(t, *l.NumDestinationChains, 1, "NumDestinationChains must be greater than or equal to 1")
	require.GreaterOrEqual(t, len(e.Chains), *l.NumDestinationChains, "NumDestinationChains must be less than or equal to the number of chains in the environment")
}

func (l *LoadConfig) GetLoadDuration() time.Duration {
	ld, _ := time.ParseDuration(*l.LoadDuration)
	return ld
}

func (l *LoadConfig) GetTimeoutDuration() time.Duration {
	ld, _ := time.ParseDuration(*l.TimeoutDuration)
	if ld == 0 {
		return 30 * time.Minute
	}
	return ld
}
