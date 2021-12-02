package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_ChainSpecificConfig_BSCDefaults(t *testing.T) {
	t.Parallel()

	c := chainSpecificConfigDefaultSets[56] // BSC
	require.Equal(t, 2*time.Second, c.ocrDatabaseTimeout)
	require.Equal(t, 2*time.Second, c.ocrContractTransmitterTransmitTimeout)
	require.Equal(t, 500*time.Millisecond, c.ocrObservationGracePeriod)
}

func Test_ChainSpecificConfig_AllComplete(t *testing.T) {
	t.Parallel()

	for _, c := range chainSpecificConfigDefaultSets {
		require.True(t, c.complete)
	}
}
