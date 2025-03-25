package operations

import (
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
	return report, report.Err
}

// ExecuteSequence executes a sequence and returns a SequenceReport.
func ExecuteSequence[IN, OUT, DEP any](
	b Bundle, sequence *Sequence[IN, OUT, DEP], deps DEP, input IN,
) (SequenceReport[IN, OUT], error) {
	b.Logger.Infow("Executing sequence", "id", sequence.def.ID,
		"version", sequence.def.Version, "description", sequence.def.Description)
	recentReporter := NewRecentMemoryReporter(b.reporter)
	newBundle := Bundle{
		Logger:     b.Logger,
		GetContext: b.GetContext,
		reporter:   recentReporter,
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
	return SequenceReport[IN, OUT]{report, executionReports}, report.Err
}

func genericReport[IN, OUT any](r Report[IN, OUT]) Report[any, any] {
	return Report[any, any]{
		ID: r.ID,
		Def: Definition{
			ID:          r.Def.ID,
			Version:     r.Def.Version,
			Description: r.Def.Description,
		},
		Output:                r.Output,
		Input:                 r.Input,
		Timestamp:             r.Timestamp,
		Err:                   r.Err,
		ChildOperationReports: r.ChildOperationReports,
	}
}

// NewUnrecoverableError creates an error that indicates an unrecoverable error.
// If this error is returned inside an operation, the operation will no longer retry.
// This allows the operation to fail fast if it encounters an unrecoverable error.
func NewUnrecoverableError(err error) error {
	return retry.Unrecoverable(err)
}
