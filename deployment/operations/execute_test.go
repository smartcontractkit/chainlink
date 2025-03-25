package operations

import (
	"context"
	"errors"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func Test_ExecuteOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		options           []ExecuteOption[int, any]
		IsUnrecoverable   bool
		wantOpCalledTimes int
		wantOutput        int
		wantErr           string
	}{
		{
			name:              "DefaultRetry",
			wantOpCalledTimes: 3,
			wantOutput:        2,
		},
		{
			name:              "NoRetry",
			options:           []ExecuteOption[int, any]{WithRetryConfig[int, any](RetryConfig[int, any]{DisableRetry: true})},
			wantOpCalledTimes: 1,
			wantErr:           "test error",
		},
		{
			name: "NewInputHook",
			options: []ExecuteOption[int, any]{WithRetryConfig[int, any](RetryConfig[int, any]{InputHook: func(input int, deps any) int {
				// update input to 5 after first failed attempt
				return 5
			}})},
			wantOpCalledTimes: 3,
			wantOutput:        6,
		},
		{
			name:              "UnrecoverableError",
			IsUnrecoverable:   true,
			wantOpCalledTimes: 1,
			wantErr:           "fatal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			failTimes := 2
			handlerCalledTimes := 0
			handler := func(b Bundle, deps any, input int) (output int, err error) {
				handlerCalledTimes++
				if tt.IsUnrecoverable {
					return 0, NewUnrecoverableError(errors.New("fatal error"))
				}

				if failTimes > 0 {
					failTimes--
					return 0, errors.New("test error")
				}

				return input + 1, nil
			}
			op := NewOperation("plus1", semver.MustParse("1.0.0"), "test operation", handler)
			e := NewBundle(context.Background, logger.Test(t), NewMemoryReporter())

			res, err := ExecuteOperation(e, op, nil, 1, tt.options...)

			if tt.wantErr != "" {
				require.Error(t, res.Err)
				require.Error(t, err)
				require.ErrorContains(t, res.Err, tt.wantErr)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, res.Err)
				require.NoError(t, err)
				assert.Equal(t, tt.wantOutput, res.Output)
			}
			assert.Equal(t, tt.wantOpCalledTimes, handlerCalledTimes)
			// check report is added to reporter
			report, err := e.reporter.GetReport(res.ID)
			require.NoError(t, err)
			assert.NotNil(t, report)
		})
	}
}

func Test_ExecuteOperation_ErrorReporter(t *testing.T) {
	op := NewOperation("plus1", semver.MustParse("1.0.0"), "test operation",
		func(e Bundle, deps any, input int) (output int, err error) {
			return input + 1, nil
		})

	reportErr := errors.New("add report error")
	errReporter := errorReporter{
		AddReportError: reportErr,
	}
	e := NewBundle(context.Background, logger.Test(t), errReporter)

	res, err := ExecuteOperation(e, op, nil, 1)
	require.Error(t, err)
	require.ErrorContains(t, err, reportErr.Error())
	require.NoError(t, res.Err)
}

type errorReporter struct {
	GetReportError  error
	GetReportsError error
	AddReportError  error
}

func (e errorReporter) GetReport(id string) (Report[any, any], error) {
	return Report[any, any]{}, e.GetReportError
}

func (e errorReporter) GetReports() ([]Report[any, any], error) {
	return nil, e.GetReportsError
}

func (e errorReporter) AddReport(report Report[any, any]) error {
	return e.AddReportError
}
