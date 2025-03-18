package operations

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Report is the result of an operation.
// It contains the inputs and other metadata that was used to execute the operation.
type Report[IN, OUT any] struct {
	ID        string     `json:"ID"`
	Def       Definition `json:"Definition"`
	Output    OUT        `json:"Output"`
	Input     IN         `json:"Input"`
	Timestamp time.Time  `json:"Timestamp"`
	Err       error      `json:"Error"`
	// stores a list of report ID for an operation that was executed as part of a sequence.
	ChildOperationReports []string `json:"ChildOperationReports"`
}

// NewReport creates a new report.
// ChildOperationReports is applicable only for Sequence.
func NewReport[IN, OUT any](
	def Definition, input IN, output OUT, err error, childReportsID ...string,
) Report[IN, OUT] {
	return Report[IN, OUT]{
		ID:                    uuid.New().String(),
		Def:                   def,
		Output:                output,
		Input:                 input,
		Timestamp:             time.Now(),
		Err:                   err,
		ChildOperationReports: childReportsID,
	}
}

var ErrReportNotFound = errors.New("report not found")

// Reporter manages reports. It can store them in memory, in the FS, etc.
type Reporter interface {
	GetReport(id string) (Report[any, any], error)
	GetReports() ([]Report[any, any], error)
	AddReport(report Report[any, any]) error
}

// MemoryReporter stores reports in memory.
type MemoryReporter struct {
	reports []Report[any, any]
}

type MemoryReporterOption func(*MemoryReporter)

// WithReports is an option to initialize the MemoryReporter with a list of reports.
func WithReports(reports []Report[any, any]) MemoryReporterOption {
	return func(mr *MemoryReporter) {
		mr.reports = reports
	}
}

// NewMemoryReporter creates a new MemoryReporter.
// It can be initialized with a list of reports using the WithReports option.
func NewMemoryReporter(options ...MemoryReporterOption) *MemoryReporter {
	reporter := &MemoryReporter{}
	for _, opt := range options {
		opt(reporter)
	}
	return reporter
}

// AddReport adds a report to the memory reporter.
func (e *MemoryReporter) AddReport(report Report[any, any]) error {
	e.reports = append(e.reports, report)
	return nil
}

// GetReports returns all reports.
func (e *MemoryReporter) GetReports() ([]Report[any, any], error) {
	return e.reports, nil
}

// GetReport returns a report by ID.
// Returns ErrReportNotFound if the report is not found.
func (e *MemoryReporter) GetReport(id string) (Report[any, any], error) {
	for _, report := range e.reports {
		if report.ID == id {
			return report, nil
		}
	}
	return Report[any, any]{}, fmt.Errorf("report_id %s: %w", id, ErrReportNotFound)
}
