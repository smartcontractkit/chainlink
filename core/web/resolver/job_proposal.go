package resolver

import (
	"context"
	"strconv"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/web/loader"
)

var notFoundErrorMessage = "spec not found"

type JobProposalStatus string

const (
	PENDING   JobProposalStatus = "PENDING"
	APPROVED  JobProposalStatus = "APPROVED"
	REJECTED  JobProposalStatus = "REJECTED"
	CANCELLED JobProposalStatus = "CANCELLED"
)

func ToJobProposalStatus(s feeds.JobProposalStatus) (JobProposalStatus, error) {
	switch s {
	case feeds.JobProposalStatusApproved:
		return APPROVED, nil
	case feeds.JobProposalStatusPending:
		return PENDING, nil
	case feeds.JobProposalStatusRejected:
		return REJECTED, nil
	case feeds.JobProposalStatusCancelled:
		return CANCELLED, nil
	default:
		return "", errors.New("invalid job proposal status")
	}
}

// JobProposalResolver resolves the Job Proposal type
type JobProposalResolver struct {
	jp *feeds.JobProposal
}

// NewJobProposal creates a new JobProposalResolver
func NewJobProposal(jp *feeds.JobProposal) *JobProposalResolver {
	return &JobProposalResolver{jp: jp}
}

func NewJobProposals(jps []feeds.JobProposal) []*JobProposalResolver {
	var resolvers []*JobProposalResolver

	for i := range jps {
		resolvers = append(resolvers, NewJobProposal(&jps[i]))
	}

	return resolvers
}

// ID resolves to the job proposal ID
func (r *JobProposalResolver) ID() graphql.ID {
	return int64GQLID(r.jp.ID)
}

// Name resolves to the job proposal name
func (r *JobProposalResolver) Name() *string {
	return r.jp.Name.Ptr()
}

// Status resolves to the job proposal Status
func (r *JobProposalResolver) Status() JobProposalStatus {
	if status, err := ToJobProposalStatus(r.jp.Status); err == nil {
		return status
	}
	return ""
}

// ExternalJobID resolves to the job proposal ExternalJobID
func (r *JobProposalResolver) ExternalJobID() *string {
	if r.jp.ExternalJobID.Valid {
		id := r.jp.ExternalJobID.UUID.String()
		return &id
	}

	return nil
}

// JobID resolves to the job proposal JobID if it has an ExternalJobID
func (r *JobProposalResolver) JobID(ctx context.Context) (*string, error) {
	if !r.jp.ExternalJobID.Valid {
		return nil, nil
	}

	job, err := loader.GetJobByExternalJobID(ctx, r.jp.ExternalJobID.UUID.String())
	if err != nil {
		return nil, err
	}

	id := strconv.FormatInt(int64(job.ID), 10)

	return &id, err
}

// FeedsManager resolves the job proposal's feeds manager object field.
func (r *JobProposalResolver) FeedsManager(ctx context.Context) (*FeedsManagerResolver, error) {
	mgr, err := loader.GetFeedsManagerByID(ctx, strconv.FormatInt(r.jp.FeedsManagerID, 10))
	if err != nil {
		return nil, err
	}

	return NewFeedsManager(*mgr), nil
}

// MultiAddrs resolves to the job proposal MultiAddrs
func (r *JobProposalResolver) MultiAddrs() []string {
	return r.jp.Multiaddrs
}

// PendingUpdate resolves to whether the job proposal has a pending update.
func (r *JobProposalResolver) PendingUpdate() bool {
	return r.jp.PendingUpdate
}

// Specs returns all spec proposals associated with the proposal.
func (r *JobProposalResolver) Specs(ctx context.Context) ([]*JobProposalSpecResolver, error) {
	specs, err := loader.GetSpecsByJobProposalID(ctx, strconv.FormatInt(r.jp.ID, 10))
	if err != nil {
		return nil, err
	}

	return NewJobProposalSpecs(specs), nil
}

// LatestSpec returns the spec with the highest version number.
func (r *JobProposalResolver) LatestSpec(ctx context.Context) (*JobProposalSpecResolver, error) {
	spec, err := loader.GetLatestSpecByJobProposalID(ctx, strconv.FormatInt(r.jp.ID, 10))
	if err != nil {
		return nil, err
	}

	return NewJobProposalSpec(spec), nil
}

// RemoteUUID returns the remote FMS UUID of the proposal.
func (r *JobProposalResolver) RemoteUUID(ctx context.Context) string {
	return r.jp.RemoteUUID.String()
}

// -- GetJobProposal Query --

// JobProposalPayloadResolver resolves the job proposal payload type
type JobProposalPayloadResolver struct {
	jp *feeds.JobProposal
	NotFoundErrorUnionType
}

// NewJobProposalPayload creates a new job proposal payload
func NewJobProposalPayload(jp *feeds.JobProposal, err error) *JobProposalPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "job proposal not found"}

	return &JobProposalPayloadResolver{jp: jp, NotFoundErrorUnionType: e}
}

// ToJobProposal resolves to the job proposal resolver
func (r *JobProposalPayloadResolver) ToJobProposal() (*JobProposalResolver, bool) {
	if r.err == nil {
		return NewJobProposal(r.jp), true
	}

	return nil, false
}
