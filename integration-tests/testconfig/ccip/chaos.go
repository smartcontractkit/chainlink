package ccip

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/stretchr/testify/require"
)

type ChaosConfig struct {
	Namespace                   string
	ExperimentFullInterval      string
	ExperimentInjectionInterval string
	ReorgDepthBelowFinality     int
	ReorgDepthAboveFinality     int
	SrcChainURL                 string
	DstChainURL                 string
}

func (l *ChaosConfig) Validate(t *testing.T, e *deployment.Environment) {
	require.NotEmpty(t, l.Namespace, "k8s namespace can't be empty")
	require.NotEmpty(t, l.ExperimentFullInterval, "experiment full interval can't be null, use Go time format 1h2m3s")
	require.NotEmpty(t, l.ExperimentInjectionInterval, "experiment injection interval can't be null, use Go time format 1h2m3s")
	require.NotEmpty(t, l.ReorgDepthBelowFinality, "reorg depth below finality can't be 0")
	require.NotEmpty(t, l.ReorgDepthAboveFinality, "reorg depth above finality can't be 0")
	require.NotEmpty(t, l.SrcChainURL, "src chain URL can't be null")
	require.NotEmpty(t, l.DstChainURL, "src chain URL can't be null")
}

func (l *ChaosConfig) GetExperimentInterval() time.Duration {
	ld, _ := time.ParseDuration(l.ExperimentFullInterval)
	return ld
}

func (l *ChaosConfig) GetExperimentInjectionInterval() time.Duration {
	ld, _ := time.ParseDuration(l.ExperimentInjectionInterval)
	return ld
}
