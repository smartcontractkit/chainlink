package deployment

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Ideally the handler receives an ExecuteSeq and ExecuteOp functions with the env loaded (reporter, etc) as param instead of the whole OpEnv
// Difficult to get it while keeping generics
type SeqHandler[I, O, D any] func(env OpEnv, deps D, input I) (output O, err error)

type Sequence[I, O, D any] struct {
	def     Definition
	handler SeqHandler[I, O, D]
}

func NewSequence[I, O, D any](version string, description string, handler SeqHandler[I, O, D]) *Sequence[I, O, D] {
	return &Sequence[I, O, D]{
		def: Definition{
			ID:          "__placeholder_	_",
			Version:     version,
			Description: description,
		},
		handler: handler,
	}
}

func NewSequenceReport[I, O, D any](sequence Sequence[I, O, D], input I, output O, err error, steps ...uuid.UUID) Report[I, O, D] {
	return Report[I, O, D]{
		ID:        uuid.New(),
		Def:       sequence.def,
		Output:    output,
		Input:     input,
		Timestamp: time.Now(),
		Err:       err,
		SubOps:    steps,
	}
}

// The cache reporter wraps the original reporter and for every operation execution, it caches the report internally
// The cache is passed down to the operation executions that are part of the sequence. Operations can only access reports within the sequence
type cacheReporter struct {
	reporter IReporter
	cache    map[string]ReportAny
}

func newCacheReporter(reporter IReporter, prevReports ...ReportAny) IReporter {
	cache := make(map[string]ReportAny)
	// if previous reporter is of cacheReporter type, we reuse its cache here too. This is for cases of sequences using another sequences
	if c, ok := reporter.(cacheReporter); ok {
		cache = c.cache
	}

	for _, r := range prevReports {
		if r.ID == uuid.Nil {
			continue
		}
		for _, subOp := range r.SubOps {
			internal, err := reporter.GetReport(subOp)
			if err == nil {
				cache[buildStepIdFromReport(internal.Def, internal.Input)] = internal
			}
		}
	}
	return cacheReporter{
		reporter: reporter,
		cache:    cache,
	}
}

func (r cacheReporter) AddReport(report ReportAny) {
	r.reporter.AddReport(report)
	r.cache[buildStepIdFromReport(report.Def, report.Input)] = report
}

func buildStepIdFromReport(def Definition, input any) string {
	// Convert struct to JSON for a consistent string representation
	// TODO: Should handle errors
	defData, _ := json.Marshal(def)

	// Convert input to JSON (to handle different input types consistently)
	inputData, _ := json.Marshal(input)

	// Combine Definition and Input data for hashing
	combinedData := append(defData, inputData...)

	return string(combinedData)
}

// We dont want to expose every report to internal operations
func (r cacheReporter) GetReports() []ReportAny {
	// Only return the cache
	reports := make([]ReportAny, 0, len(r.cache))
	for _, report := range r.cache {
		reports = append(reports, report)
	}
	return reports
}

func (r cacheReporter) GetReport(id uuid.UUID) (ReportAny, error) {
	// Navigate through cache and compare Id
	for _, report := range r.cache {
		if report.ID == id {
			return report, nil
		}
	}
	return ReportAny{}, errors.New("report not found")
}

// Main function to execute sequences.
func ExecuteSeq[I, O, D any](env OpEnv, sequence Sequence[I, O, D], deps D, input I, reportId ...uuid.UUID) Report[I, O, D] {
	var prevReport ReportAny
	var err error
	if len(reportId) == 1 {
		prevReport, err = env.Reporter.GetReport(reportId[0])
		if err != nil {
			env.Log.Errorw("Failed to retrieve report for retry", "id", reportId[0], "error", err)
		}
	}

	// We intercept the reporter to cache the reports and limit report access to operations. If resuming, we cache the previous reports
	reporter := newCacheReporter(env.Reporter, prevReport)

	env.Log.Infow("Executing sequence", "id", sequence.def.ID, "version", sequence.def.Version, "description", sequence.def.Description)
	ret, err := sequence.handler(OpEnv{Log: env.Log, Reporter: reporter}, deps, input)

	// We return the report, with sequence operations as steps
	internalReps := reporter.GetReports()
	steps := make([]uuid.UUID, 0, len(internalReps))
	for _, rep := range internalReps {
		steps = append(steps, rep.ID)
	}

	seqReport := NewSequenceReport(
		sequence,
		input,
		ret,
		err,
		steps...,
	)

	// Add the report to the reporter
	genericReport := ReportAny{
		ID: seqReport.ID,
		Def: Definition{
			ID:          sequence.def.ID,
			Version:     sequence.def.Version,
			Description: sequence.def.Description,
		},
		Output:    seqReport.Output,
		Input:     seqReport.Input,
		Timestamp: seqReport.Timestamp,
		Err:       seqReport.Err,
		SubOps:    seqReport.SubOps,
	}

	env.Reporter.AddReport(genericReport)
	return seqReport
}
