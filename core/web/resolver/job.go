package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/web/loader"
)

// JobResolver resolves the Job type.
type JobResolver struct {
	app chainlink.Application
	j   job.Job
}

func NewJob(app chainlink.Application, j job.Job) *JobResolver {
	return &JobResolver{app: app, j: j}
}

func NewJobs(app chainlink.Application, jobs []job.Job) []*JobResolver {
	var resolvers []*JobResolver
	for _, j := range jobs {
		resolvers = append(resolvers, NewJob(app, j))
	}

	return resolvers
}

// ID resolves the job's id.
func (r *JobResolver) ID() graphql.ID {
	return int32GQLID(r.j.ID)
}

// CreatedAt resolves the job's created at timestamp.
func (r *JobResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.j.CreatedAt}
}

// Errors resolves the job's top level errors.
func (r *JobResolver) Errors(ctx context.Context) ([]*JobErrorResolver, error) {
	specErrs, err := loader.GetJobSpecErrorsByJobID(ctx, r.j.ID)
	if err != nil {
		return nil, err
	}

	return NewJobErrors(specErrs), nil
}

// ExternalJobID resolves the job's external job id.
func (r *JobResolver) ExternalJobID() string {
	return r.j.ExternalJobID.String()
}

// MaxTaskDuration resolves the job's max task duration.
func (r *JobResolver) MaxTaskDuration() string {
	return r.j.MaxTaskDuration.Duration().String()
}

// Name resolves the job's name.
func (r *JobResolver) Name() string {
	return r.j.Name.ValueOrZero()
}

// ObservationSource resolves the job's observation source.
//
// This could potentially be moved to a dataloader in the future as we are
// fetching it from a relationship.
func (r *JobResolver) ObservationSource() string {
	return r.j.PipelineSpec.DotDagSource
}

// SchemaVersion resolves the job's schema version.
func (r *JobResolver) SchemaVersion() int32 {
	return int32(r.j.SchemaVersion)
}

// GasLimit resolves the job's gas limit.
func (r *JobResolver) GasLimit() *int32 {
	if !r.j.GasLimit.Valid {
		return nil
	}
	v := int32(r.j.GasLimit.Uint32)
	return &v
}

// ForwardingAllowed sets whether txs submitted by this job should be forwarded when possible.
func (r *JobResolver) ForwardingAllowed() *bool {
	return &r.j.ForwardingAllowed
}

// Type resolves the job's type.
func (r *JobResolver) Type() string {
	return string(r.j.Type)
}

// Spec resolves the job's spec.
func (r *JobResolver) Spec() *SpecResolver {
	return NewSpec(r.j)
}

// Runs fetches the runs for a Job.
func (r *JobResolver) Runs(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*JobRunsPayloadResolver, error) {
	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	if limit > 100 {
		limit = 100
	}

	ids, err := r.app.JobORM().FindPipelineRunIDsByJobID(ctx, r.j.ID, offset, limit)
	if err != nil {
		return nil, err
	}

	runs, err := loader.GetJobRunsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	count, err := r.app.JobORM().CountPipelineRunsByJobID(ctx, r.j.ID)
	if err != nil {
		return nil, err
	}

	return NewJobRunsPayload(runs, count, r.app), nil
}

// JobsPayloadResolver resolves a page of jobs
type JobsPayloadResolver struct {
	app   chainlink.Application
	jobs  []job.Job
	total int32
}

func NewJobsPayload(app chainlink.Application, jobs []job.Job, total int32) *JobsPayloadResolver {
	return &JobsPayloadResolver{
		app:   app,
		jobs:  jobs,
		total: total,
	}
}

// Results returns the jobs.
func (r *JobsPayloadResolver) Results() []*JobResolver {
	return NewJobs(r.app, r.jobs)
}

// Metadata returns the pagination metadata.
func (r *JobsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

type JobPayloadResolver struct {
	app chainlink.Application
	job *job.Job
	NotFoundErrorUnionType
}

func NewJobPayload(app chainlink.Application, j *job.Job, err error) *JobPayloadResolver {
	e := NotFoundErrorUnionType{err, "job not found", nil}

	return &JobPayloadResolver{app: app, job: j, NotFoundErrorUnionType: e}
}

// ToJob implements the JobPayload union type of the payload
func (r *JobPayloadResolver) ToJob() (*JobResolver, bool) {
	if r.job != nil {
		return NewJob(r.app, *r.job), true
	}

	return nil, false
}

// -- CreateJob Mutation --

type CreateJobPayloadResolver struct {
	app       chainlink.Application
	j         *job.Job
	inputErrs map[string]string
}

func NewCreateJobPayload(app chainlink.Application, job *job.Job, inputErrs map[string]string) *CreateJobPayloadResolver {
	return &CreateJobPayloadResolver{app: app, j: job, inputErrs: inputErrs}
}

func (r *CreateJobPayloadResolver) ToCreateJobSuccess() (*CreateJobSuccessResolver, bool) {
	if r.inputErrs != nil {
		return nil, false
	}

	return NewCreateJobSuccess(r.app, r.j), true
}

func (r *CreateJobPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs == nil {
		return nil, false
	}

	var errs []*InputErrorResolver

	for path, message := range r.inputErrs {
		errs = append(errs, NewInputError(path, message))
	}

	return NewInputErrors(errs), true
}

type CreateJobSuccessResolver struct {
	app chainlink.Application
	j   *job.Job
}

func NewCreateJobSuccess(app chainlink.Application, job *job.Job) *CreateJobSuccessResolver {
	return &CreateJobSuccessResolver{app: app, j: job}
}

func (r *CreateJobSuccessResolver) Job() *JobResolver {
	return NewJob(r.app, *r.j)
}

// -- DeleteJob Mutation --

type DeleteJobPayloadResolver struct {
	app chainlink.Application
	j   *job.Job
	NotFoundErrorUnionType
}

func NewDeleteJobPayload(app chainlink.Application, j *job.Job, err error) *DeleteJobPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "job not found"}

	return &DeleteJobPayloadResolver{app: app, j: j, NotFoundErrorUnionType: e}
}

func (r *DeleteJobPayloadResolver) ToDeleteJobSuccess() (*DeleteJobSuccessResolver, bool) {
	if r.j == nil {
		return nil, false
	}

	return NewDeleteJobSuccess(r.app, r.j), true
}

type DeleteJobSuccessResolver struct {
	app chainlink.Application
	j   *job.Job
}

func NewDeleteJobSuccess(app chainlink.Application, job *job.Job) *DeleteJobSuccessResolver {
	return &DeleteJobSuccessResolver{app: app, j: job}
}

func (r *DeleteJobSuccessResolver) Job() *JobResolver {
	return NewJob(r.app, *r.j)
}
