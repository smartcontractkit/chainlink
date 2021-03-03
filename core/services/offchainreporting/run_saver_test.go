package offchainreporting

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunSaver(t *testing.T) {
	pipelineRunner := new(mocks.Runner)
	rr := make(chan pipeline.RunWithResults, 100)
	rs := NewResultRunSaver(
		rr,
		pipelineRunner,
		make(chan struct{}),
		1,
	)
	require.NoError(t, rs.Start())
	for i := 0; i < 100; i++ {
		pipelineRunner.On("InsertFinishedRunWithResults", mock.Anything, mock.Anything, mock.Anything).
			Return(int64(i), nil).Once()
		rr <- pipeline.RunWithResults{Run: pipeline.Run{ID: int64(i)}}
	}
	require.NoError(t, rs.Close())
}
