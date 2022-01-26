package resolver

import (
	"strconv"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/job"
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
	return graphql.ID(strconv.FormatInt(r.specError.ID, 10))
}

// Description resolves the job error's description.
func (r *JobErrorResolver) Description() string {
	return r.specError.Description
}

// Occurrences resolves the job error's number of occurances.
func (r *JobErrorResolver) Occurrences() int32 {
	return int32(r.specError.Occurrences)
}

// CreatedAt resolves the job error's created at timestamp.
func (r *JobErrorResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.specError.CreatedAt}
}
