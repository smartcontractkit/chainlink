package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
)

// JobErrorResolver resolves a Job Error
type JobErrorResolver struct {
	// This is purposefully named Error instead of Err to differentiate it from
	// a standard golang error.
	specError job.SpecError
}

func NewJobError(specError job.SpecError) *JobErrorResolver {
	return &JobErrorResolver{specError: specError}
}

func NewJobErrors(specErrors []job.SpecError) []*JobErrorResolver {
	var resolvers []*JobErrorResolver
	for _, e := range specErrors {
		resolvers = append(resolvers, NewJobError(e))
	}

	return resolvers
}

// ID resolves the job error's id.
func (r *JobErrorResolver) ID() graphql.ID {
	return graphql.ID(stringutils.FromInt64(r.specError.ID))
}

// Description resolves the job error's description.
func (r *JobErrorResolver) Description() string {
	return r.specError.Description
}

// Occurrences resolves the job error's number of occurrences.
func (r *JobErrorResolver) Occurrences() int32 {
	return int32(r.specError.Occurrences)
}

// CreatedAt resolves the job error's created at timestamp.
func (r *JobErrorResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.specError.CreatedAt}
}

// UpdatedAt resolves the job error's updated at timestamp.
func (r *JobErrorResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.specError.UpdatedAt}
}

// -- DismissJobError Mutation --

type DismissJobErrorPayloadResolver struct {
	specError *job.SpecError
	NotFoundErrorUnionType
}

func NewDismissJobErrorPayload(specError *job.SpecError, err error) *DismissJobErrorPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "JobSpecError not found"}

	return &DismissJobErrorPayloadResolver{specError: specError, NotFoundErrorUnionType: e}
}

func (r *DismissJobErrorPayloadResolver) ToDismissJobErrorSuccess() (*DismissJobErrorSuccessResolver, bool) {
	if r.err != nil {
		return nil, false
	}

	return NewDismissJobErrorSuccess(r.specError), true
}

type DismissJobErrorSuccessResolver struct {
	specError *job.SpecError
}

func NewDismissJobErrorSuccess(specError *job.SpecError) *DismissJobErrorSuccessResolver {
	return &DismissJobErrorSuccessResolver{specError: specError}
}

func (r *DismissJobErrorSuccessResolver) JobError() *JobErrorResolver {
	return NewJobError(*r.specError)
}
