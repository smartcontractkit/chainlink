package deployment

import (
	"time"

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

// An Operation is defined by a unique ID and version. It has a description explaining what it does.
type OperationDefinition struct {
	ID          string
	Version     string
	Description string
}

// Operations are the low level building blocks of the system. Developers are completely free to define their own operation with their own input and output types.
// They execute one operation, which can perform max 1 side effect (e.g. send a transaction, post a job spec...)
// TODO: There should be some constraint on Input, Output and Deps
type Operation[Input, Output, Deps any] struct {
	def      OperationDefinition
	execFunc ExecuteFunc[Input, Output, Deps]
}

// TODO: Add std context.Context
func NewOperation[I, O, D any](version string, description string, execFunc ExecuteFunc[I, O, D]) *Operation[I, O, D] {
	return &Operation[I, O, D]{
		def: OperationDefinition{
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

func (o *Operation[I, O, D]) Inspect(ctx OpContext, deps D, input I) (err error) {
	// TODO: Inspection returns the payload the execute function will send (e.g. the transaction data, the job spec, etc) Useful for composition (e.g. generate MCMS proposals) and debugging
	return nil
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
	OpDef     OperationDefinition
	Output    O
	Input     I
	Timestamp time.Time
	err       error
}

func NewReport[I, O, D any](operation Operation[I, O, D], input I, output O, err error) Report[I, O, D] {
	return Report[I, O, D]{
		OpDef:     operation.def,
		Output:    output,
		Input:     input,
		Timestamp: time.Now(),
		err:       err,
	}
}

type ReportAny Report[any, any, any]

// Reprter manages reports. It can store them in memory, in the FS, etc.
type IReporter interface {
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

// OpEnv holds the utilities to use the Operations Client. This env is not exposed directly to operations
type OpEnv struct {
	Reporter IReporter
	Log      logger.Logger
}

// Operations are low level, and should rarely be used directly.
// Execute is the main function to interact with the operations. Standarizes the execution API and experience.
// Could be expanded to accept middlewares to support default logging, tracing, etc.
func ExecuteOp[I, O, D any](
	env OpEnv,
	operation *Operation[I, O, D],
	deps D,
	input I,
) (Report[I, O, D], error) {
	// TODO: Check if report is in cache. Return if exists

	// Build the context utilities needed fo the operation
	ctx := OpContext{
		Log: env.Log,
	}

	output, err := operation.Execute(ctx, deps, input)
	report := NewReport(*operation, input, output, err)

	// We store a generic report. As is only for storing, we don't mind losing types there
	genericReport := ReportAny{
		OpDef: OperationDefinition{
			ID:          operation.ID(),
			Version:     operation.Version(),
			Description: operation.Description(),
		},
		Output:    report.Output,
		Input:     report.Input,
		Timestamp: report.Timestamp,
		err:       report.err,
	}

	env.Reporter.AddReport(genericReport)
	return report, err
}
