package resolver

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/job"
)

// JobResolver resolves the Job type.
type JobResolver struct {
	j job.Job
}

func NewJob(j job.Job) *JobResolver {
	return &JobResolver{j: j}
}

func NewJobs(jobs []job.Job) []*JobResolver {
	resolvers := []*JobResolver{}
	for _, j := range jobs {
		resolvers = append(resolvers, NewJob(j))
	}

	return resolvers
}

// ID resolves the job's id.
func (r *JobResolver) ID() graphql.ID {
	return graphql.ID(strconv.FormatInt(int64(r.j.ID), 10))
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
	if r.j.Name.IsZero() {
		return "No name"

	}

	return r.j.Name.ValueOrZero()
}

// ObservationSource resolves the job's observation source.
//
// This could potentially by moved to a dataloader in the future as we are
// fetching it from a relationship.
func (r *JobResolver) ObservationSource() string {
	return r.j.PipelineSpec.DotDagSource
}

// SchemaVersion resolves the job's schema version.
func (r *JobResolver) SchemaVersion() int32 {
	return int32(r.j.SchemaVersion)
}

// Spec resolves the job's spec.
func (r *JobResolver) Spec() *SpecResolver {
	return NewSpec(r.j)
}

// CreatedAt resolves the job's created at timestamp.
func (r *JobResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.j.CreatedAt}
}

// JobsPayloadResolver resolves a page of jobs
type JobsPayloadResolver struct {
	jobs  []job.Job
	total int32
}

func NewJobsPayload(jobs []job.Job, total int32) *JobsPayloadResolver {
	return &JobsPayloadResolver{
		jobs:  jobs,
		total: total,
	}
}

// Results returns the bridges.
func (r *JobsPayloadResolver) Results() []*JobResolver {
	return NewJobs(r.jobs)
}

// Metadata returns the pagination metadata.
func (r *JobsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

type JobPayloadResolver struct {
	job *job.Job
	err error
}

func NewJobPayload(j *job.Job, err error) *JobPayloadResolver {
	return &JobPayloadResolver{
		job: j,
		err: err,
	}
}

// ToJob implements the JobPayload union type of the payload
func (r *JobPayloadResolver) ToJob() (*JobResolver, bool) {
	if r.job != nil {
		return NewJob(*r.job), true
	}

	return nil, false
}

// ToNotFoundError implements the NotFoundError union type of the payload
func (r *JobPayloadResolver) ToNotFoundError() (*NotFoundErrorResolver, bool) {
	if r.err != nil && errors.Is(r.err, sql.ErrNoRows) {
		return NewNotFoundError("job not found"), true
	}

	return nil, false
}
