package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	ccip "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestJobSpecChangeset(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  4,
	})
	output, err := changeset.CCIPCapabilityJobspecChangeset(e, nil)
	require.NoError(t, err)
	require.NotNil(t, output.JobSpecs)
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	for _, node := range nodes {
		jobs, exists := output.JobSpecs[node.NodeID]
		require.True(t, exists)
		require.NotNil(t, jobs)
		for _, job := range jobs {
			_, err = ccip.ValidatedCCIPSpec(job)
			require.NoError(t, err)
		}
	}
}
