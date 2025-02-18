package deployment

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Context provide common utilities to the operations
type OpContext struct {
	Log logger.Logger
}

// TODO: We might want to enforce Validation on inputs
type Input interface {
	Validate() error
}

type ExecuteFunc[I, O, D any] func(ctx OpContext, deps D, input I) (output O, err error)

// An Operation / Sequence is defined by a unique ID and version. It has a description explaining what it does.
type Definition struct {
	ID          string
	Version     string
	Description string
}

// Operations are the low level building blocks of the system. Developers are completely free to define their own operation with their own input and output types.
// They execute one operation, which can perform max 1 side effect (e.g. send a transaction, post a job spec...)
// TODO: There should be some constraint on Input, Output and Deps
type Operation[Input, Output, Deps any] struct {
	def      Definition
	execFunc ExecuteFunc[Input, Output, Deps]
}

// TODO: Add std context.Context
func NewOperation[I, O, D any](version string, description string, execFunc ExecuteFunc[I, O, D]) *Operation[I, O, D] {
	return &Operation[I, O, D]{
		def: Definition{
			// Id and version are useful to identify the operation
			ID:          "__placeholder__",
			Version:     version,
			Description: description,
		},
		execFunc: execFunc,
	}
}

func (o *Operation[I, O, D]) Execute(ctx OpContext, deps D, input I) (output O, err error) {
	ctx.Log.Infow("Executing operation", "id", o.def.ID, "version", o.def.Version, "description", o.def.Description)
	return o.execFunc(ctx, deps, input)
}

func (o *Operation[I, O, D]) ID() string {
	return o.def.ID
}

// TODO: Version should be a standard semver
func (o *Operation[I, O, D]) Version() string {
	return o.def.Version
}

func (o *Operation[I, O, D]) Description() string {
	return o.def.Description
}

type EmptyInput struct{}

// Reports
type Report[I, O, D any] struct {
	ID        uuid.UUID
	Def       Definition
	Output    O
	Input     I
	Timestamp time.Time
	// Internal report IDs
	SubOps []uuid.UUID
	Err    error
}

func NewReport[I, O, D any](operation Operation[I, O, D], input I, output O, err error) Report[I, O, D] {
	return Report[I, O, D]{
		ID:        uuid.New(),
		Def:       operation.def,
		Output:    output,
		Input:     input,
		Timestamp: time.Now(),
		Err:       err,
	}
}

type ReportAny Report[any, any, any]

// Reprter manages reports. It can store them in memory, in the FS, etc.
type IReporter interface {
	GetReport(id uuid.UUID) (ReportAny, error)
	GetReports() []ReportAny
	AddReport(report ReportAny)
}

// In memory reporter
type MemoryReporter struct {
	// cache of operations
	reports []ReportAny
}

func NewMemoryReporter(reports []ReportAny) *MemoryReporter {
	reporter := &MemoryReporter{}
	reporter.reports = reports
	return reporter
}

func (e *MemoryReporter) AddReport(report ReportAny) {
	// Add to storage
	e.reports = append(e.reports, report)
}

func (e *MemoryReporter) GetReports() []ReportAny {
	return e.reports
}

func (e *MemoryReporter) GetReport(id uuid.UUID) (ReportAny, error) {
	for _, report := range e.reports {
		if report.ID == id {
			return report, nil
		}
	}
	return ReportAny{}, errors.New("report not found")
}

func typeReport[I, O, D any](r ReportAny) (Report[I, O, D], bool) {
	// Attempt to type assert the Input and Output fields
	input, ok1 := r.Input.(I)
	output, ok2 := r.Output.(O)

	if !ok1 || !ok2 {
		return Report[I, O, D]{}, false
	}

	return Report[I, O, D]{
		ID:        r.ID,
		Def:       r.Def,
		Output:    output,
		Input:     input,
		Timestamp: r.Timestamp,
		SubOps:    r.SubOps,
		Err:       r.Err,
	}, true
}

// OpEnv holds the utilities to use the Operations Client. This env is not exposed directly to operations
type OpEnv struct {
	Reporter IReporter
	Log      logger.Logger
}

// Operations are low level, and should rarely be used directly.
// Execute is the main function to interact with the operations. Standarizes the execution API and experience.
// It accepts an environment with logger and reporter:
// - Reporter holds previous reports of executed operations. The user of the function should be reponsible of the reports exposed to this function
// - The function will look up on the Reporter previous reports for previous executions. If match, it will return the op report
// Could be expanded to accept middlewares to support default logging, tracing, etc.
// TODO: We might want to split the previois reporte check into a separate function destined to be used only by sequences. Can get messy if ExecuteOp is used outside Sequences (ExecuteOpWithRetry or similar)
func ExecuteOp[I, O, D any](
	env OpEnv,
	operation *Operation[I, O, D],
	deps D,
	input I,
) (Report[I, O, D], error) {
	prevReports := env.Reporter.GetReports()
	thisOpID := buildStepIdFromReport(operation.def, input)

	// Check if operation was run previously and return the report if successful
	for _, report := range prevReports {
		if buildStepIdFromReport(report.Def, report.Input) == thisOpID && report.Err == nil {
			typedReport, ok := typeReport[I, O, D](report)
			if !ok {
				env.Log.Debugw("Operation already executed but couldn't type its Report", "id", report.ID)
				// We couldn't type the report, something is off. We execute the operation again
				continue
			}
			env.Log.Debugw("Operation already executed. Returning its result from Report storage", "id", report.ID)
			return typedReport, nil
		}
	}

	// Build the context utilities needed fo the operation
	ctx := OpContext{
		Log: env.Log,
	}

	output, err := operation.Execute(ctx, deps, input)
	report := NewReport(*operation, input, output, err)

	// We store a generic report. As is only for storing, we don't mind losing types there
	genericReport := ReportAny{
		ID: report.ID,
		Def: Definition{
			ID:          operation.ID(),
			Version:     operation.Version(),
			Description: operation.Description(),
		},
		Output:    report.Output,
		Input:     report.Input,
		Timestamp: report.Timestamp,
		Err:       report.Err,
	}

	env.Reporter.AddReport(genericReport)
	return report, err
}
