package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/core/web/loader"
)

type JobRunStatus string

const (
	JobRunStatusUnknown   JobRunStatus = "UNKNOWN"
	JobRunStatusRunning   JobRunStatus = "RUNNING"
	JobRunStatusSuspended JobRunStatus = "SUSPENDED"
	JobRunStatusErrored   JobRunStatus = "ERRORED"
	JobRunStatusCompleted JobRunStatus = "COMPLETED"
)

func NewJobRunStatus(status pipeline.RunStatus) JobRunStatus {
	switch status {
	case pipeline.RunStatusRunning:
		return JobRunStatusRunning
	case pipeline.RunStatusSuspended:
		return JobRunStatusSuspended
	case pipeline.RunStatusErrored:
		return JobRunStatusErrored
	case pipeline.RunStatusCompleted:
		return JobRunStatusCompleted
	default:
		return JobRunStatusUnknown
	}
}

var outputRetrievalErrorStr = "error: unable to retrieve outputs"

type JobRunResolver struct {
	run pipeline.Run
	app chainlink.Application
}

func NewJobRun(run pipeline.Run, app chainlink.Application) *JobRunResolver {
	return &JobRunResolver{run: run, app: app}
}

func NewJobRuns(runs []pipeline.Run, app chainlink.Application) []*JobRunResolver {
	var resolvers []*JobRunResolver

	for _, run := range runs {
		resolvers = append(resolvers, NewJobRun(run, app))
	}

	return resolvers
}

func (r *JobRunResolver) ID() graphql.ID {
	return int64GQLID(r.run.ID)
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
	return int32GQLID(r.run.PipelineSpecID)
}

func (r *JobRunResolver) FatalErrors() []string {
	var errs []string

	for _, err := range r.run.StringFatalErrors() {
		if err != nil {
			errs = append(errs, *err)
		}
	}

	return errs
}

func (r *JobRunResolver) AllErrors() []string {
	var errs []string

	for _, err := range r.run.StringAllErrors() {
		if err != nil {
			errs = append(errs, *err)
		}
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

func (r *JobRunResolver) Status() JobRunStatus {
	return NewJobRunStatus(r.run.State)
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

func (r *JobRunResolver) Job(ctx context.Context) (*JobResolver, error) {
	plnSpecID := stringutils.FromInt32(r.run.PipelineSpecID)
	job, err := loader.GetJobByPipelineSpecID(ctx, plnSpecID)
	if err != nil {
		return nil, err
	}

	return NewJob(r.app, *job), nil
}

func (r *JobRunResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.run.CreatedAt}
}

func (r *JobRunResolver) FinishedAt() *graphql.Time {
	return &graphql.Time{Time: r.run.FinishedAt.ValueOrZero()}
}

// -- JobRun query --

type JobRunPayloadResolver struct {
	jr  *pipeline.Run
	app chainlink.Application
	NotFoundErrorUnionType
}

func NewJobRunPayload(jr *pipeline.Run, app chainlink.Application, err error) *JobRunPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "job run not found", isExpectedErrorFn: nil}

	return &JobRunPayloadResolver{jr: jr, app: app, NotFoundErrorUnionType: e}
}

func (r *JobRunPayloadResolver) ToJobRun() (*JobRunResolver, bool) {
	if r.err != nil {
		return nil, false
	}

	return NewJobRun(*r.jr, r.app), true
}

// JobRunsPayloadResolver resolves a page of job runs
type JobRunsPayloadResolver struct {
	runs  []pipeline.Run
	total int32
	app   chainlink.Application
}

func NewJobRunsPayload(runs []pipeline.Run, total int32, app chainlink.Application) *JobRunsPayloadResolver {
	return &JobRunsPayloadResolver{
		runs:  runs,
		total: total,
		app:   app,
	}
}

// Results returns the job runs.
func (r *JobRunsPayloadResolver) Results() []*JobRunResolver {
	return NewJobRuns(r.runs, r.app)
}

// Metadata returns the pagination metadata.
func (r *JobRunsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

// -- RunJob Mutation --

type RunJobPayloadResolver struct {
	run *pipeline.Run
	app chainlink.Application
	NotFoundErrorUnionType
}

func NewRunJobPayload(run *pipeline.Run, app chainlink.Application, err error) *RunJobPayloadResolver {
	var e NotFoundErrorUnionType

	if err != nil {
		e = NotFoundErrorUnionType{err: err, message: err.Error(), isExpectedErrorFn: func(err error) bool {
			return errors.Is(err, webhook.ErrJobNotExists)
		}}
	}

	return &RunJobPayloadResolver{run: run, app: app, NotFoundErrorUnionType: e}
}

func (r *RunJobPayloadResolver) ToRunJobSuccess() (*RunJobSuccessResolver, bool) {
	if r.err != nil {
		return nil, false
	}

	return NewRunJobSuccess(*r.run, r.app), true
}

func (r *RunJobPayloadResolver) ToRunJobCannotRunError() (*RunJobCannotRunErrorResolver, bool) {
	if r.err == nil {
		return nil, false
	}

	if errors.Is(r.err, webhook.ErrJobNotExists) {
		return nil, false
	}

	return NewRunJobCannotRunError(r.err), true
}

type RunJobSuccessResolver struct {
	run pipeline.Run
	app chainlink.Application
}

func NewRunJobSuccess(run pipeline.Run, app chainlink.Application) *RunJobSuccessResolver {
	return &RunJobSuccessResolver{run: run, app: app}
}

func (r *RunJobSuccessResolver) JobRun() *JobRunResolver {
	return NewJobRun(r.run, r.app)
}

type RunJobCannotRunErrorResolver struct {
	message string
	code    ErrorCode
}

func NewRunJobCannotRunError(err error) *RunJobCannotRunErrorResolver {
	return &RunJobCannotRunErrorResolver{message: "", code: ErrorCodeUnprocessable}
}

func (r *RunJobCannotRunErrorResolver) Code() ErrorCode {
	return r.code
}

func (r *RunJobCannotRunErrorResolver) Message() string {
	return r.message
}
