package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
)

// JobProposalResolver resolves the Job Proposal type
type JobProposalResolver struct {
	jp *feeds.JobProposal
}

// NewJobProposal creates a new JobProposalResolver
func NewJobProposal(jp *feeds.JobProposal) *JobProposalResolver {
	return &JobProposalResolver{jp}
}

// ID resolves to the job proposal ID
func (r *JobProposalResolver) ID() graphql.ID {
	return int32GQLID(int32(r.jp.ID))
}

// Spec resolves to the job proposal Spec
func (r *JobProposalResolver) Spec() string {
	return r.jp.Spec
}

// Status resolves to the job proposal Status
func (r *JobProposalResolver) Status() string {
	return string(r.jp.Status)
}

// ExternalJobID resolves to the job proposal ExternalJobID
func (r *JobProposalResolver) ExternalJobID() string {
	if r.jp.ExternalJobID.Valid {
		return r.jp.ExternalJobID.UUID.String()
	}

	return "no valid"
}

// MultiAddrs resolves to the job proposal MultiAddrs
func (r *JobProposalResolver) MultiAddrs() []string {
	return r.jp.Multiaddrs
}

// ProposedAt resolves to the job proposal ProposedAt date
func (r *JobProposalResolver) ProposedAt() graphql.Time {
	return graphql.Time{Time: r.jp.ProposedAt}
}

// -- GetJobProposal Query --

// JobProposalPayloadResolver resolves the job proposal payload type
type JobProposalPayloadResolver struct {
	jp  *feeds.JobProposal
	err error
}

// NewJobProposalPayload creates a new job proposal payload
func NewJobProposalPayload(jp *feeds.JobProposal, err error) *JobProposalPayloadResolver {
	return &JobProposalPayloadResolver{jp, err}
}

// ToJobProposal resolves to the job proposal resolver
func (r *JobProposalPayloadResolver) ToJobProposal() (*JobProposalResolver, bool) {
	if r.err == nil {
		return NewJobProposal(r.jp), true
	}

	return nil, false
}

// ToNotFoundError resolves to the not found error resolver
func (r *JobProposalPayloadResolver) ToNotFoundError() (*NotFoundErrorResolver, bool) {
	if r.err != nil {
		return NewNotFoundError("job proposal not found"), true
	}

	return nil, false
}
