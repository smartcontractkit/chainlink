package resolver

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/web/loader"
)

var notFoundErrorMessage = "job proposal not found"

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
func (r *JobProposalResolver) Status() JobProposalStatus {
	if status, err := ToJobProposalStatus(r.jp.Status); err == nil {
		return status
	}
	return ""
}

// ExternalJobID resolves to the job proposal ExternalJobID
func (r *JobProposalResolver) ExternalJobID() string {
	if r.jp.ExternalJobID.Valid {
		return r.jp.ExternalJobID.UUID.String()
	}

	return "no valid"
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
	if r.err != nil && errors.Is(r.err, sql.ErrNoRows) {
		return NewNotFoundError("job proposal not found"), true
	}

	return nil, false
}

// -- Mutations shared interface --

type JobProposalPayload interface {
	ToNotFoundError() (*NotFoundErrorResolver, bool)
}

type JobProposalAction string

const (
	approve JobProposalAction = "approve"
	cancel  JobProposalAction = "cancel"
)

func ToJobProposalPayload(a JobProposalAction, jp *feeds.JobProposal, err error) (JobProposalPayload, error) {
	switch a {
	case approve:
		return NewApproveJobProposalPayload(jp, err), nil
	case cancel:
		return NewCancelJobProposalPayload(jp, err), nil
	default:
		return nil, errors.New("invalid job proposal action")
	}
}

// -- ApproveJobProposal Mutation --

type ApproveJobProposalPayloadResolver struct {
	jp *feeds.JobProposal
	ResolvesNotFoundError
}

func NewApproveJobProposalPayload(jp *feeds.JobProposal, err error) *ApproveJobProposalPayloadResolver {
	e := ResolvesNotFoundError{err: err, message: notFoundErrorMessage}

	return &ApproveJobProposalPayloadResolver{jp, e}
}

// ToApproveJobProposalSuccess resolves to the approval job proposal success resolver
func (r *ApproveJobProposalPayloadResolver) ToApproveJobProposalSuccess() (*ApproveJobProposalSuccessResolver, bool) {
	if r.jp != nil {
		return NewApproveJobProposalSuccess(r.jp), true
	}

	return nil, false
}

type ApproveJobProposalSuccessResolver struct {
	jp *feeds.JobProposal
}

func NewApproveJobProposalSuccess(jp *feeds.JobProposal) *ApproveJobProposalSuccessResolver {
	return &ApproveJobProposalSuccessResolver{jp}
}

func (r *ApproveJobProposalSuccessResolver) JobProposal() *JobProposalResolver {
	return NewJobProposal(r.jp)
}

// -- CancelJobProposal Mutation --

type CancelJobProposalPayloadResolver struct {
	jp *feeds.JobProposal
	ResolvesNotFoundError
}

func NewCancelJobProposalPayload(jp *feeds.JobProposal, err error) *CancelJobProposalPayloadResolver {
	e := ResolvesNotFoundError{err: err, message: notFoundErrorMessage}

	return &CancelJobProposalPayloadResolver{jp, e}
}

// ToCancelJobProposalSuccess resolves to the approval job proposal success resolver
func (r *CancelJobProposalPayloadResolver) ToCancelJobProposalSuccess() (*CancelJobProposalSuccessResolver, bool) {
	if r.jp != nil {
		return NewCancelJobProposalSuccess(r.jp), true
	}

	return nil, false
}

type CancelJobProposalSuccessResolver struct {
	jp *feeds.JobProposal
}

func NewCancelJobProposalSuccess(jp *feeds.JobProposal) *CancelJobProposalSuccessResolver {
	return &CancelJobProposalSuccessResolver{jp}
}

func (r *CancelJobProposalSuccessResolver) JobProposal() *JobProposalResolver {
	return NewJobProposal(r.jp)
}
