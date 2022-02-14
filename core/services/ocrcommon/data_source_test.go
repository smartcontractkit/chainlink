package ocrcommon_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	pipelinemocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	mockValue = "100000000"
)

func Test_InMemoryDataSource(t *testing.T) {
	runner := new(pipelinemocks.Runner)
	runner.On("ExecuteRun", mock.Anything, mock.AnythingOfType("pipeline.Spec"), mock.Anything, mock.Anything).
		Return(pipeline.Run{}, pipeline.TaskRunResults{
			{
				Result: pipeline.Result{
					Value: mockValue,
					Error: nil,
				},
				Task: &pipeline.HTTPTask{},
			},
		}, nil)

	ds := ocrcommon.NewInMemoryDataSource(runner, job.Job{}, pipeline.Spec{}, logger.TestLogger(t))
	val, err := ds.Observe(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, mockValue, val.String()) // returns expected value after pipeline run
}

func Test_NewDataSourceV2(t *testing.T) {
	runner := new(pipelinemocks.Runner)
	runner.On("ExecuteRun", mock.Anything, mock.AnythingOfType("pipeline.Spec"), mock.Anything, mock.Anything).
		Return(pipeline.Run{}, pipeline.TaskRunResults{
			{
				Result: pipeline.Result{
					Value: mockValue,
					Error: nil,
				},
				Task: &pipeline.HTTPTask{},
			},
		}, nil)

	resChan := make(chan pipeline.Run, 100)
	ds := ocrcommon.NewDataSourceV2(runner, job.Job{}, pipeline.Spec{}, logger.TestLogger(t), resChan)
	val, err := ds.Observe(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, mockValue, val.String())   // returns expected value after pipeline run
	assert.Equal(t, pipeline.Run{}, <-resChan) // expected data properly passed to channel
}
