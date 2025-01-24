package deployment

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Context provide common utilities and dependencies to the operations
type Context[Deps any] struct {
	Log logger.Logger
	// Operations can require dependencies. Dependencies are interfaces to be implemented and passed in by the caller.
	Deps Deps
}

type ExecuteFunc[I, O, Deps any] func(context Context[Deps], input I) (output O, err error)

// An Operation is defined by a unique ID and version. It has a description explaining what it does.
type OperationDefinition struct {
	Id          string
	Version     string
	Description string
}

// TODO: We might want to enforce Validation on inputs
type Input interface {
	Validate() error
}

// Operations are the low level building blocks of the system. Developers are completely free to define their own operation with their own input and output types.
// They execute one operation, which can perform max 1 side effect (e.g. send a transaction, post a job spec...)
type Operation[I, O, Deps any] struct {
	def      OperationDefinition
	input    I // might not be needed
	execFunc ExecuteFunc[I, O, Deps]
}

// TODO: Add std context.Context
func NewOperation[I, O, Deps any](version string, description string, execFunc ExecuteFunc[I, O, Deps]) *Operation[I, O, Deps] {
	return &Operation[I, O, Deps]{
		def: OperationDefinition{
			// Id and version are useful to identify the operation
			Id:          "__placeholder__",
			Version:     version,
			Description: description,
		},
		execFunc: execFunc,
	}
}

func (o *Operation[I, O, Deps]) Execute(ctx Context[Deps], input I) (output O, err error) {
	ctx.Log.Infow("Executing operation", "id", o.def.Id, "version", o.def.Version, "description", o.def.Description)
	return o.execFunc(ctx, input)
}

func (o *Operation[I, O, Deps]) Inspect(ctx Context[Deps], input I) (err error) {
	// TODO: Inspection returns the payload the execute function will send (e.g. the transaction data, the job spec, etc) Useful for composition (e.g. generate MCMS proposals) and debugging
	return nil
}

func (o *Operation[I, O, Deps]) ID() string {
	return o.def.Id
}

// TODO: Version should be a standard semver
func (o *Operation[I, O, Deps]) Version() string {
	return o.def.Version
}

func (o *Operation[I, O, Deps]) Description() string {
	return o.def.Description
}

type EmptyInput struct{}

// Reports
type Report[I, O, Deps any] struct {
	OpDef     OperationDefinition
	Output    O
	Input     I
	Timestamp time.Time
	err       error
}

func NewReport[I, O, Deps any](operation Operation[I, O, Deps], input I, output O, err error) Report[I, O, Deps] {
	return Report[I, O, Deps]{
		OpDef:     operation.def,
		Output:    output,
		Input:     input,
		Timestamp: time.Now(),
		err:       err,
	}
}

// Reprter manages reports. It can store them in memory, in the FS, etc.
type IReporter interface {
	GetReports() []Report[any, any, any]
	AddReport(report Report[any, any, any])
}

// In memory reporter
type MemoryReporter struct {
	// cache of operations
	reports []Report[any, any, any]
}

func NewMemoryReporter(reports []Report[any, any, any]) *MemoryReporter {
	reporter := &MemoryReporter{}
	reporter.reports = reports
	return reporter
}

func (e *MemoryReporter) AddReport(report Report[any, any, any]) {
	// Add to cache
	e.reports = append(e.reports, report)
}

func (e *MemoryReporter) GetReports() []Report[any, any, any] {
	return e.reports
}

// Operations are low level, and should rarely be used directly.
// Execute is the main function to interact with the operations. Standarizes the execution API and experience.
// Could be expanded to accept middlewares to support default logging, tracing, etc.
func Execute[I, O, Deps any](
	reporter IReporter,
	operation *Operation[I, O, Deps],
	ctx Context[Deps],
	input I,
) (Report[I, O, Deps], error) {
	// TODO: Check if report is in cache. Return if exists
	output, err := operation.Execute(ctx, input)
	report := NewReport(*operation, input, output, err)

	// We store a generic report as is only for storing, we don't mind losing types there
	genericReport := Report[any, any, any]{
		OpDef: OperationDefinition{
			Id:          operation.ID(),
			Version:     operation.Version(),
			Description: operation.Description(),
		},
		Output:    report.Output,
		Input:     report.Input,
		Timestamp: report.Timestamp,
		err:       report.err,
	}

	reporter.AddReport(genericReport)
	return report, err
}
