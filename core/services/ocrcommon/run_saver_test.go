package ocrcommon

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
)

func TestRunSaver(t *testing.T) {
	pipelineRunner := new(mocks.Runner)
	rr := make(chan pipeline.Run, 100)
	rs := NewResultRunSaver(
		rr,
		pipelineRunner,
		make(chan struct{}),
		logger.TestLogger(t),
	)
	require.NoError(t, rs.Start(testutils.Context(t)))
	for i := 0; i < 100; i++ {
		d := i
		pipelineRunner.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).(*pipeline.Run).ID = int64(d)
			}).
			Once()
		rr <- pipeline.Run{ID: int64(i)}
	}
	require.NoError(t, rs.Close())
}
