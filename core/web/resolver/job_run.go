package resolver

import (
	"strconv"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

var outputRetrievalErrorStr = "error: unable to retrieve outputs"

type JobRunResolver struct {
	run pipeline.Run
}

func NewJobRun(run pipeline.Run) *JobRunResolver {
	return &JobRunResolver{run}
}

func NewJobRuns(runs []pipeline.Run) []*JobRunResolver {
	var resolvers []*JobRunResolver

	for _, run := range runs {
		resolvers = append(resolvers, NewJobRun(run))
	}

	return resolvers
}

func (r *JobRunResolver) ID() graphql.ID {
	return graphql.ID(strconv.Itoa(int(r.run.ID)))
}

func (r *JobRunResolver) Outputs() []*string {
	if !r.run.Outputs.Valid {
		return []*string{&outputRetrievalErrorStr}
	}

	outputs, err := r.run.StringOutputs()
	if err != nil {
		errMsg := err.Error()
		return []*string{&errMsg}
	}

	return outputs
}

func (r *JobRunResolver) PipelineSpecID() graphql.ID {
	return graphql.ID(strconv.Itoa(int(r.run.PipelineSpecID)))
}

func (r *JobRunResolver) FatalErrors() []string {
	var errs []string

	for _, err := range r.run.StringFatalErrors() {
		errs = append(errs, *err)
	}

	return errs
}

func (r *JobRunResolver) AllErrors() []string {
	var errs []string

	for _, err := range r.run.StringAllErrors() {
		errs = append(errs, *err)
	}

	return errs
}

func (r *JobRunResolver) Inputs() string {
	val, err := r.run.Inputs.MarshalJSON()
	if err != nil {
		return "error: unable to retrieve inputs"
	}

	return string(val)
}

// TaskRuns resolves the job run's task runs
//
// This could be moved to a data loader later, which means also modifying to ORM
// to not get everything at once
func (r *JobRunResolver) TaskRuns() []*TaskRunResolver {
	if len(r.run.PipelineTaskRuns) > 0 {
		return NewTaskRuns(r.run.PipelineTaskRuns)
	}

	return []*TaskRunResolver{}
}

func (r *JobRunResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.run.CreatedAt}
}

func (r *JobRunResolver) FinishedAt() *graphql.Time {
	return &graphql.Time{Time: r.run.FinishedAt.ValueOrZero()}
}
