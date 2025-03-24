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

// SequenceReport is a report for a sequence.
// It contains a report for the sequence itself and also a list of reports
// for all the operations executed as part of the sequence.
// The latter is useful when we want to return all the reports of the operations
// executed as part of the sequence in changeset output.
type SequenceReport[IN, OUT any] struct {
	Report[IN, OUT]

	// ExecutionReports is a list of report all the operations & sequence that was executed as part of this sequence.
	ExecutionReports []Report[any, any]
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
	GetExecutionReports(reportID string) ([]Report[any, any], error)
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

// GetExecutionReports returns all the reports that was executed as part of a sequence including itself.
// It does this by recursively fetching all the child reports.
// Useful when returning all the reports in a sequence to the changeset output.
func (e *MemoryReporter) GetExecutionReports(seqID string) ([]Report[any, any], error) {
	var allReports []Report[any, any]

	var getReportsRecursively func(id string) error
	getReportsRecursively = func(id string) error {
		report, err := e.GetReport(id)
		if err != nil {
			return err
		}
		for _, childID := range report.ChildOperationReports {
			if err := getReportsRecursively(childID); err != nil {
				return err
			}
		}
		allReports = append(allReports, report)

		return nil
	}

	if err := getReportsRecursively(seqID); err != nil {
		return nil, err
	}

	return allReports, nil
}

// RecentReporter is a wrapper around a Reporter that keeps track of the most recent reports.
// Useful when trying to get a list of reports that was recently added in a sequence.
type RecentReporter struct {
	Reporter
	recentReports []Report[any, any]
}

// AddReport adds a report to the recent reporter.
func (e *RecentReporter) AddReport(report Report[any, any]) error {
	err := e.Reporter.AddReport(report)
	if err != nil {
		return err
	}
	e.recentReports = append(e.recentReports, report)
	return nil
}

// GetRecentReports returns all the reports that was added since the construction of the RecentReporter.
func (e *RecentReporter) GetRecentReports() []Report[any, any] {
	return e.recentReports
}

// NewRecentMemoryReporter creates a new RecentReporter.
func NewRecentMemoryReporter(reporter Reporter) *RecentReporter {
	r := &RecentReporter{
		Reporter:      reporter,
		recentReports: []Report[any, any]{},
	}
	return r
}
