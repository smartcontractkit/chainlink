package fluxmonitorv2_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPipelineRun_Execute(t *testing.T) {
	dec := decimal.NewFromInt(100)

	testCases := []struct {
		name       string
		runID      int64
		results    pipeline.FinalResult
		err        error
		expRunID   int64
		expDecimal *decimal.Decimal
		expErr     bool
	}{
		{
			name:  "success",
			runID: 1,
			results: pipeline.FinalResult{
				Values: []interface{}{dec},
				Errors: []error{nil},
			},
			err:        nil,
			expRunID:   1,
			expDecimal: &dec,
			expErr:     false,
		},
		{
			name:  "pipeline error",
			runID: 0,
			results: pipeline.FinalResult{
				Values: []interface{}{},
				Errors: []error{nil},
			},
			err:        errors.New("pipeline error"),
			expRunID:   0,
			expDecimal: nil,
			expErr:     true,
		},
		{
			name:  "error extracting singular result",
			runID: 1,
			results: pipeline.FinalResult{
				Values: []interface{}{},
				Errors: []error{nil},
			},
			err:        nil,
			expRunID:   1,
			expDecimal: nil,
			expErr:     true,
		},
		{
			name:  "error converting result to decimal",
			runID: 1,
			results: pipeline.FinalResult{
				Values: []interface{}{"str"},
				Errors: []error{nil},
			},
			err:        nil,
			expRunID:   1,
			expDecimal: nil,
			expErr:     true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var (
				runner      = new(mocks.Runner)
				spec        = pipeline.Spec{}
				jobID       = int32(1)
				l           = *logger.Default
				pipelineRun = fluxmonitorv2.NewPipelineRun(runner, spec, jobID, l)
			)

			runner.
				On("ExecuteAndInsertNewRun", context.Background(), spec, l).
				Return(tc.runID, tc.results, tc.err)

			aRunID, aDecimal, aErr := pipelineRun.Execute()

			assert.Equal(t, tc.expRunID, aRunID)
			assert.Equal(t, tc.expDecimal, aDecimal)
			if tc.expErr {
				assert.Error(t, aErr)
			} else {
				assert.NoError(t, aErr)
			}

			runner.AssertExpectations(t)
		})
	}

}
