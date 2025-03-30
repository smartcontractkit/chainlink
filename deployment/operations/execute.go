package operations

import (
	"fmt"

	"github.com/avast/retry-go/v4"
)

// ExecuteConfig is the configuration for the ExecuteOperation function.
type ExecuteConfig[IN, DEP any] struct {
	retryConfig RetryConfig[IN, DEP]
}

type ExecuteOption[IN, DEP any] func(*ExecuteConfig[IN, DEP])

type RetryConfig[IN, DEP any] struct {
	// DisableRetry disables the retry mechanism if set to true.
	DisableRetry bool
	// InputHook is a function that returns an updated input before retrying the operation.
	// The operation when retried will use the input returned by this function.
	// This is useful for scenarios like updating the gas limit.
	// This will be ignored if DisableRetry is set to true.
	InputHook func(input IN, deps DEP) IN
}

// WithRetryConfig is an ExecuteOption that sets the retry configuration.
func WithRetryConfig[IN, DEP any](config RetryConfig[IN, DEP]) ExecuteOption[IN, DEP] {
	return func(c *ExecuteConfig[IN, DEP]) {
		c.retryConfig = config
	}
}

// ExecuteOperation executes an operation with the given input and dependencies.
// Execution will return the previous successful execution result and skip execution if there was a
// previous successful run found in the Reports.
// If previous unsuccessful execution was found, the execution will not be skipped.
//
// Note:
// Operations that were skipped will not be added to the reporter.
//
// Retry:
// By default, it retries the operation up to 10 times with exponential backoff if it fails.
// Use WithRetryConfig to customize the retry behavior.
// To cancel the retry early, return an error with NewUnrecoverableError.
func ExecuteOperation[IN, OUT, DEP any](
	b Bundle,
	operation *Operation[IN, OUT, DEP],
	deps DEP,
	input IN,
	opts ...ExecuteOption[IN, DEP],
) (Report[IN, OUT], error) {
	if previousReport, found := loadPreviousSuccessfulReport[IN, OUT](b, operation.def, input); found {
		b.Logger.Infow("Operation already executed. Returning previous result", "id", operation.def.ID,
			"version", operation.def.Version, "description", operation.def.Description)
		return previousReport, nil
	}

	executeConfig := &ExecuteConfig[IN, DEP]{retryConfig: RetryConfig[IN, DEP]{}}
	for _, opt := range opts {
		opt(executeConfig)
	}

	var output OUT
	var err error

	if executeConfig.retryConfig.DisableRetry {
		output, err = operation.execute(b, deps, input)
	} else {
		var inputTemp = input
		output, err = retry.DoWithData(func() (OUT, error) {
			return operation.execute(b, deps, inputTemp)
		}, retry.OnRetry(func(attempt uint, err error) {
			b.Logger.Infow("Operation failed. Retrying...",
				"operation", operation.def.ID, "attempt", attempt, "error", err)

			if executeConfig.retryConfig.InputHook != nil {
				inputTemp = executeConfig.retryConfig.InputHook(inputTemp, deps)
			}
		}))
	}

	report := NewReport(operation.def, input, output, err)
	err = b.reporter.AddReport(genericReport(report))
	if err != nil {
		return Report[IN, OUT]{}, err
	}
	if report.Err != nil {
		return report, report.Err
	}
	return report, nil
}

// ExecuteSequence executes a Sequence and returns a SequenceReport.
// The SequenceReport contains a report for the Sequence and also the execution reports which are all
// the operations that were executed as part of this sequence.
// The latter is useful when we want to return all the executed reports to the changeset output.
// Execution will return the previous successful execution result and skip execution if there was a
// previous successful run found in the Reports.
// If previous unsuccessful execution was found, the execution will not be skipped.
//
// Note:
// Sequences or Operations that were skipped will not be added to the reporter.
// THe ExecutionReports do not include Sequences or Operations that were skipped.
func ExecuteSequence[IN, OUT, DEP any](
	b Bundle, sequence *Sequence[IN, OUT, DEP], deps DEP, input IN,
) (SequenceReport[IN, OUT], error) {
	if previousReport, found := loadPreviousSuccessfulReport[IN, OUT](b, sequence.def, input); found {
		executionReports, err := b.reporter.GetExecutionReports(previousReport.ID)
		if err != nil {
			return SequenceReport[IN, OUT]{}, err
		}
		b.Logger.Infow("Sequence already executed. Returning previous result", "id", sequence.def.ID,
			"version", sequence.def.Version, "description", sequence.def.Description)
		return SequenceReport[IN, OUT]{previousReport, executionReports}, nil
	}

	b.Logger.Infow("Executing sequence", "id", sequence.def.ID,
		"version", sequence.def.Version, "description", sequence.def.Description)
	recentReporter := NewRecentMemoryReporter(b.reporter)
	newBundle := Bundle{
		Logger:          b.Logger,
		GetContext:      b.GetContext,
		reporter:        recentReporter,
		reportHashCache: b.reportHashCache,
	}
	ret, err := sequence.handler(newBundle, deps, input)

	recentReports := recentReporter.GetRecentReports()
	childReports := make([]string, 0, len(recentReports))
	for _, rep := range recentReports {
		childReports = append(childReports, rep.ID)
	}

	report := NewReport(
		sequence.def,
		input,
		ret,
		err,
		childReports...,
	)

	err = b.reporter.AddReport(genericReport(report))
	if err != nil {
		return SequenceReport[IN, OUT]{}, err
	}
	executionReports, err := b.reporter.GetExecutionReports(report.ID)
	if err != nil {
		return SequenceReport[IN, OUT]{}, err
	}
	if report.Err != nil {
		return SequenceReport[IN, OUT]{report, executionReports}, report.Err
	}
	return SequenceReport[IN, OUT]{report, executionReports}, nil
}

// NewUnrecoverableError creates an error that indicates an unrecoverable error.
// If this error is returned inside an operation, the operation will no longer retry.
// This allows the operation to fail fast if it encounters an unrecoverable error.
func NewUnrecoverableError(err error) error {
	return retry.Unrecoverable(err)
}

func loadPreviousSuccessfulReport[IN, OUT any](
	b Bundle, def Definition, input IN,
) (Report[IN, OUT], bool) {
	prevReports, err := b.reporter.GetReports()
	if err != nil {
		b.Logger.Errorw("Failed to get reports", "error", err)
		return Report[IN, OUT]{}, false
	}
	currentHash, err := constructUniqueHashFrom(b.reportHashCache, def, input)
	if err != nil {
		b.Logger.Errorw("Failed to construct unique hash", "error", err)
		return Report[IN, OUT]{}, false
	}

	for _, report := range prevReports {
		// Check if operation/sequence was run previously and return the report if successful
		reportHash, err := constructUniqueHashFrom(b.reportHashCache, report.Def, report.Input)
		if err != nil {
			b.Logger.Errorw("Failed to construct unique hash for previous report", "error", err)
			continue
		}
		if reportHash == currentHash && report.Err == nil {
			typedReport, ok := typeReport[IN, OUT](report)
			if !ok {
				b.Logger.Debugw(fmt.Sprintf("Previous %s execution found but couldn't find its matching Report", def.ID), "report_id", report.ID)
				continue
			}
			b.Logger.Debugw(fmt.Sprintf("Previous %s execution found. Returning its result from Report storage", def.ID), "report_id", report.ID)
			return typedReport, true
		}
	}
	// No previous execution was found
	return Report[IN, OUT]{}, false
}
