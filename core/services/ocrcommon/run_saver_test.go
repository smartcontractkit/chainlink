package ocrcommon

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
)

func TestRunSaver(t *testing.T) {
	pipelineRunner := mocks.NewRunner(t)
	rs := NewResultRunSaver(
		pipelineRunner,
		logger.TestLogger(t),
		1000,
		100,
	)
	servicetest.Run(t, rs)
	for i := 0; i < 100; i++ {
		d := i
		pipelineRunner.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(2).(*pipeline.Run).ID = int64(d)
			}).
			Once()
		rs.Save(&pipeline.Run{ID: int64(i)})
	}
}
