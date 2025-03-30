package operations

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"sync"
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
				require.Nil(t, res.Err)
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
		Reporter:       NewMemoryReporter(),
		AddReportError: reportErr,
	}
	e := NewBundle(context.Background, logger.Test(t), errReporter)

	res, err := ExecuteOperation(e, op, nil, 1)
	require.Error(t, err)
	require.ErrorContains(t, err, reportErr.Error())
	require.Nil(t, res.Err)
}

func Test_ExecuteOperation_WithPreviousRun(t *testing.T) {
	t.Parallel()

	handlerCalledTimes := 0
	handler := func(b Bundle, deps any, input int) (output int, err error) {
		handlerCalledTimes++
		return input + 1, nil
	}
	handlerWithErrorCalledTimes := 0
	handlerWithError := func(b Bundle, deps any, input int) (output int, err error) {
		handlerWithErrorCalledTimes++
		return 0, NewUnrecoverableError(errors.New("test error"))
	}

	op := NewOperation("plus1", semver.MustParse("1.0.0"), "test operation", handler)
	opWithError := NewOperation("plus1-error", semver.MustParse("1.0.0"), "test operation error", handlerWithError)
	bundle := NewBundle(t.Context, logger.Test(t), NewMemoryReporter())

	// first run
	res, err := ExecuteOperation(bundle, op, nil, 1)
	require.NoError(t, err)
	require.Nil(t, res.Err)
	assert.Equal(t, 2, res.Output)
	assert.Equal(t, 1, handlerCalledTimes)

	// rerun should return previous report
	res, err = ExecuteOperation(bundle, op, nil, 1)
	require.NoError(t, err)
	require.Nil(t, res.Err)
	assert.Equal(t, 2, res.Output)
	assert.Equal(t, 1, handlerCalledTimes)

	// new run with different input, should perform execution
	res, err = ExecuteOperation(bundle, op, nil, 3)
	require.NoError(t, err)
	require.Nil(t, res.Err)
	assert.Equal(t, 4, res.Output)
	assert.Equal(t, 2, handlerCalledTimes)

	// new run with different op, should perform execution
	op = NewOperation("plus1-v2", semver.MustParse("2.0.0"), "test operation", handler)
	res, err = ExecuteOperation(bundle, op, nil, 1)
	require.NoError(t, err)
	require.Nil(t, res.Err)
	assert.Equal(t, 2, res.Output)
	assert.Equal(t, 3, handlerCalledTimes)

	// new run with op that returns error
	res, err = ExecuteOperation(bundle, opWithError, nil, 1)
	require.Error(t, err)
	require.ErrorContains(t, err, "test error")
	require.ErrorContains(t, res.Err, "test error")
	assert.Equal(t, 1, handlerWithErrorCalledTimes)

	// rerun with op that returns error, should attempt execution again
	res, err = ExecuteOperation(bundle, opWithError, nil, 1)
	require.Error(t, err)
	require.ErrorContains(t, err, "test error")
	require.ErrorContains(t, res.Err, "test error")
	assert.Equal(t, 2, handlerWithErrorCalledTimes)
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
				require.Nil(t, seqReport.Err)
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

func Test_ExecuteSequence_WithPreviousRun(t *testing.T) {
	t.Parallel()

	version := semver.MustParse("1.0.0")
	op := NewOperation("plus1", version, "plus 1",
		func(b Bundle, deps OpDeps, input int) (output int, err error) {
			return input + 1, nil
		})

	handlerCalledTimes := 0
	handler := func(b Bundle, deps any, input int) (int, error) {
		handlerCalledTimes++
		res, err := ExecuteOperation(b, op, OpDeps{}, input)
		if err != nil {
			return 0, err
		}
		return res.Output, nil
	}
	handlerWithErrorCalledTimes := 0
	handlerWithError := func(b Bundle, deps any, input int) (int, error) {
		handlerWithErrorCalledTimes++
		return 0, NewUnrecoverableError(errors.New("test error"))
	}
	sequence := NewSequence("seq-plus1", version, "plus 1", handler)
	sequenceWithError := NewSequence("seq-plus1-error", version, "plus 1 error", handlerWithError)

	bundle := NewBundle(context.Background, logger.Test(t), NewMemoryReporter())

	// first run
	res, err := ExecuteSequence(bundle, sequence, nil, 1)
	require.NoError(t, err)
	require.Nil(t, res.Err)
	assert.Equal(t, 2, res.Output)
	assert.Len(t, res.ExecutionReports, 2) // 1 seq report + 1 op report
	assert.Equal(t, 1, handlerCalledTimes)

	marshal, err := json.MarshalIndent(res.ExecutionReports, "", "  ")
	require.NoError(t, err)
	t.Log(string(marshal))
	// rerun should return previous report
	res, err = ExecuteSequence(bundle, sequence, nil, 1)
	require.NoError(t, err)
	require.Nil(t, res.Err)
	assert.Equal(t, 2, res.Output)
	assert.Len(t, res.ExecutionReports, 2) // 1 seq report + 1 op report
	assert.Equal(t, 1, handlerCalledTimes)

	// new run with different input, should perform execution
	res, err = ExecuteSequence(bundle, sequence, nil, 3)
	require.NoError(t, err)
	require.Nil(t, res.Err)
	assert.Equal(t, 4, res.Output)
	assert.Len(t, res.ExecutionReports, 2) // 1 seq report + 1 op report
	assert.Equal(t, 2, handlerCalledTimes)

	// new run with different sequence but same operation, should perform execution
	sequence = NewSequence("seq-plus1-v2", semver.MustParse("2.0.0"), "plus 1", handler)
	res, err = ExecuteSequence(bundle, sequence, nil, 1)
	require.NoError(t, err)
	require.Nil(t, res.Err)
	assert.Equal(t, 2, res.Output)
	// only 1 because the op was not executed due to previous execution found
	assert.Len(t, res.ExecutionReports, 1)
	assert.Equal(t, 3, handlerCalledTimes)

	// new run with sequence that returns error
	res, err = ExecuteSequence(bundle, sequenceWithError, nil, 1)
	require.Error(t, err)
	require.ErrorContains(t, err, "test error")
	require.ErrorContains(t, res.Err, "test error")
	assert.Equal(t, 1, handlerWithErrorCalledTimes)

	// rerun with sequence that returns error, should attempt execution again
	res, err = ExecuteSequence(bundle, sequenceWithError, nil, 1)
	require.Error(t, err)
	require.ErrorContains(t, err, "test error")
	require.ErrorContains(t, res.Err, "test error")
	assert.Equal(t, 2, handlerWithErrorCalledTimes)
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
		name          string
		setupReporter func() Reporter
		wantErr       string
	}{
		{
			name: "AddReport returns an error",
			setupReporter: func() Reporter {
				return errorReporter{
					Reporter:       NewMemoryReporter(),
					AddReportError: errors.New("add report error"),
				}
			},
			wantErr: "add report error",
		},
		{
			name: "GetExecutionReports returns an error",
			setupReporter: func() Reporter {
				return errorReporter{
					Reporter:                 NewMemoryReporter(),
					GetExecutionReportsError: errors.New("get execution reports error"),
				}
			},
			wantErr: "get execution reports error",
		},
		{
			name: "Loaded previous report but GetExecutionReports returns an error",
			setupReporter: func() Reporter {
				r := errorReporter{
					Reporter:                 NewMemoryReporter(),
					GetExecutionReportsError: errors.New("get execution reports error"),
				}
				err := r.AddReport(genericReport(
					NewReport(sequence.def, 1, 2, nil),
				))
				require.NoError(t, err)

				return r
			},
			wantErr: "get execution reports error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := NewBundle(context.Background, logger.Test(t), tt.setupReporter())
			_, err := ExecuteSequence(e, sequence, OpDeps{}, 1)
			require.Error(t, err)
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func Test_constructUniqueHashFrom(t *testing.T) {
	t.Parallel()

	type Input struct {
		A int
		B int
	}

	definition := Definition{
		ID:          "plus1",
		Version:     semver.MustParse("1.0.0"),
		Description: "plus 1",
	}
	tests := []struct {
		name    string
		def     Definition
		input   any
		want    string
		wantErr string
	}{
		{
			name: "Same def and input should always have the same hash (struct input)",
			def:  definition,
			input: Input{
				A: 1,
				B: 2,
			},
			want: "e6148d004b97353d8361d8cbcfbefe77da97dec220bd04449667f6ba6180d46c",
		},
		{
			name: "Same def and same input (different map order) should always have the same hash (struct input)",
			def:  definition,
			input: Input{
				B: 2,
				A: 1,
			},
			want: "e6148d004b97353d8361d8cbcfbefe77da97dec220bd04449667f6ba6180d46c",
		},
		{
			name:  "Same def and input should always have the same hash (literal input)",
			def:   definition,
			input: 1,
			want:  "3f409a93e9fe5507c0d4a902ba5cb6e80a3740c74dc6a4a9bca13e71f2d46ca5",
		},
		{
			name: "Different def, same input should have different hash",
			def: Definition{
				ID:          "plus2",
				Version:     semver.MustParse("1.0.0"),
				Description: "plus 2",
			},
			input: 1,
			want:  "9f9d42fb8ced1129a0071a30a301cc03e94f12f225f7649ab34df16e6891d37c",
		},
		{
			name: "Different input, same def should have different hash (struct input)",
			def:  definition,
			input: Input{
				B: 1,
				A: 2,
			},
			want: "6c8b5986bd08e115df5cbbd77c23fcd28ea31055d0cf6de8ec29ddd7fd056bca",
		},
		{
			name:  "Different input, same def should have different hash (literal input)",
			def:   definition,
			input: 2,
			want:  "468cd7f2b432fa59559c69a2db32fc0d44d829328a6db5c7d745889c53f4fff0",
		},
		{
			name:    "invalid input",
			def:     definition,
			input:   math.NaN(),
			wantErr: "json: unsupported value: NaN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := &sync.Map{}
			hash, err := constructUniqueHashFrom(cache, tt.def, tt.input)
			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)

				// should not store anything in cache on failure
				key := struct {
					Def   Definition
					Input any
				}{tt.def, tt.input}
				_, ok := cache.Load(key)
				require.False(t, ok)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, hash)

				// check cache
				key := struct {
					Def   Definition
					Input any
				}{tt.def, tt.input}
				cached, ok := cache.Load(key)
				require.True(t, ok)
				assert.Equal(t, tt.want, cached)

				// this call should use the cache
				hash2, err := constructUniqueHashFrom(cache, tt.def, tt.input)
				require.NoError(t, err)
				assert.Equal(t, tt.want, hash2)
			}
		})
	}
}

func Test_loadPreviousSuccessfulReport(t *testing.T) {
	t.Parallel()

	version := semver.MustParse("1.0.0")
	definition := Definition{
		ID:          "plus1",
		Version:     version,
		Description: "plus 1",
	}

	tests := []struct {
		name          string
		setupReporter func() Reporter
		input         float64
		wantDef       Definition
		wantInput     float64
		wantFound     bool
	}{
		{
			name: "Failed to GetReports",
			setupReporter: func() Reporter {
				return errorReporter{
					GetReportsError: errors.New("failed to get reports"),
				}
			},
			input:     1,
			wantFound: false,
		},
		{
			name: "Successful Report found - return report",
			setupReporter: func() Reporter {
				r := NewMemoryReporter()
				err := r.AddReport(genericReport(
					NewReport(definition, 1, 2, nil),
				))
				require.NoError(t, err)

				return r
			},
			input:     1,
			wantDef:   definition,
			wantInput: 1,
			wantFound: true,
		},
		{
			name: "Report with error found - ignore report",
			setupReporter: func() Reporter {
				r := NewMemoryReporter()
				err := r.AddReport(genericReport(
					NewReport(definition, 1, 2, errors.New("failed")),
				))
				require.NoError(t, err)

				return r
			},
			input:     1,
			wantFound: false,
		},
		{
			name:      "Report not found",
			input:     1,
			wantFound: false,
		},
		{
			name:      "Current report with bad hash",
			input:     math.NaN(),
			wantFound: false,
		},
		{
			name: "Previous report with bad hash",
			setupReporter: func() Reporter {
				r := NewMemoryReporter()
				err := r.AddReport(genericReport(
					NewReport(definition, math.NaN(), 2, nil),
				))
				require.NoError(t, err)

				return r
			},
			input:     1,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bundle := NewBundle(context.Background, logger.Test(t), NewMemoryReporter())
			if tt.setupReporter != nil {
				bundle.reporter = tt.setupReporter()
			}

			report, found := loadPreviousSuccessfulReport[float64, int](bundle, definition, tt.input)
			assert.Equal(t, tt.wantFound, found)

			if tt.wantFound {
				assert.Equal(t, tt.wantDef, report.Def)
				assert.InDelta(t, tt.wantInput, report.Input, 0)
			}
		})
	}
}

type errorReporter struct {
	Reporter
	GetReportError           error
	GetReportsError          error
	AddReportError           error
	GetExecutionReportsError error
}

func (e errorReporter) GetReport(id string) (Report[any, any], error) {
	if e.GetReportError != nil {
		return Report[any, any]{}, e.GetReportError
	}
	return e.Reporter.GetReport(id)
}

func (e errorReporter) GetReports() ([]Report[any, any], error) {
	if e.GetReportsError != nil {
		return nil, e.GetReportsError
	}
	return e.Reporter.GetReports()
}

func (e errorReporter) AddReport(report Report[any, any]) error {
	if e.AddReportError != nil {
		return e.AddReportError
	}
	return e.Reporter.AddReport(report)
}

func (e errorReporter) GetExecutionReports(id string) ([]Report[any, any], error) {
	if e.GetExecutionReportsError != nil {
		return nil, e.GetExecutionReportsError
	}
	return e.Reporter.GetExecutionReports(id)
}
