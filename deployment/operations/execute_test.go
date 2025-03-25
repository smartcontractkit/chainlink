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

func Test_ExecuteSequence(t *testing.T) {
	t.Parallel()

	version := semver.MustParse("1.0.0")

	tests := []struct {
		name            string
		simulateOpError bool
		wantOutput      int
		wantErr         string
	}{
		{
			name:       "Success Execution",
			wantOutput: 3,
		},
		{
			name:            "Error Execution",
			simulateOpError: true,
			wantErr:         "fatal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			op := NewOperation("plus1", version, "plus 1",
				func(e Bundle, deps OpDeps, input int) (output int, err error) {
					if tt.simulateOpError {
						return 0, NewUnrecoverableError(errors.New("fatal error"))
					}
					return input + 1, nil
				})

			var opID string
			sequence := NewSequence("seq-plus1", version, "plus 1",
				func(env Bundle, deps any, input int) (int, error) {
					res, err := ExecuteOperation(env, op, OpDeps{}, input)
					// capture for verification later
					opID = res.ID
					if err != nil {
						return 0, err
					}

					return res.Output + 1, nil
				})

			e := NewBundle(context.Background, logger.Test(t), NewMemoryReporter())

			seqReport, err := ExecuteSequence(e, sequence, nil, 1)

			if tt.simulateOpError {
				require.Error(t, seqReport.Err)
				require.Error(t, err)
				require.ErrorContains(t, seqReport.Err, tt.wantErr)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, seqReport.Err)
				require.NoError(t, err)
				assert.Equal(t, tt.wantOutput, seqReport.Output)
			}
			assert.Equal(t, []string{opID}, seqReport.ChildOperationReports)
			// check report is added to reporter
			report, err := e.reporter.GetReport(seqReport.ID)
			require.NoError(t, err)
			assert.NotNil(t, report)
			assert.Len(t, seqReport.ExecutionReports, 2) // 1 seq report + 1 op report

			// check allReports contain the parent and child reports
			childReport, err := e.reporter.GetReport(opID)
			require.NoError(t, err)
			assert.Equal(t, seqReport.ExecutionReports[0], childReport)
			assert.Equal(t, seqReport.ExecutionReports[1], report)
		})
	}
}

func Test_ExecuteSequence_ErrorReporter(t *testing.T) {
	t.Parallel()

	version := semver.MustParse("1.0.0")
	op := NewOperation("plus1", version, "plus 1",
		func(e Bundle, deps OpDeps, input int) (output int, err error) {
			return input + 1, nil
		})

	sequence := NewSequence("seq-plus1", version, "plus 1",
		func(env Bundle, deps OpDeps, input int) (int, error) {
			res, err := ExecuteOperation(env, op, OpDeps{}, input)
			if err != nil {
				return 0, err
			}

			return res.Output + 1, nil
		})

	tests := []struct {
		name           string
		setupErrorFunc func() errorReporter
		wantErr        string
	}{
		{
			name: "AddReport returns an error",
			setupErrorFunc: func() errorReporter {
				return errorReporter{
					AddReportError: errors.New("add report error"),
				}
			},
			wantErr: "add report error",
		},
		{
			name: "GetExecutionReports returns an error",
			setupErrorFunc: func() errorReporter {
				return errorReporter{
					GetExecutionReportsError: errors.New("get execution reports error"),
				}
			},
			wantErr: "get execution reports error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := NewBundle(context.Background, logger.Test(t), tt.setupErrorFunc())
			_, err := ExecuteSequence(e, sequence, OpDeps{}, 1)
			require.Error(t, err)
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

type errorReporter struct {
	GetReportError           error
	GetReportsError          error
	AddReportError           error
	GetExecutionReportsError error
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

func (e errorReporter) GetExecutionReports(id string) ([]Report[any, any], error) {
	return nil, e.GetExecutionReportsError
}
