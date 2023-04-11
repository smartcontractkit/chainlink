package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

// SpecStatus defines the enum values for GQL
type SpecStatus string

const (
	// revive:disable
	SpecStatusUnknown   SpecStatus = "UNKNOWN"
	SpecStatusPending   SpecStatus = "PENDING"
	SpecStatusApproved  SpecStatus = "APPROVED"
	SpecStatusRejected  SpecStatus = "REJECTED"
	SpecStatusCancelled SpecStatus = "CANCELLED"
	SpecStatusRevoked   SpecStatus = "REVOKED"
	// revive:enable
)

// ToSpecStatus converts the feeds status into the enum value.
func ToSpecStatus(s feeds.SpecStatus) SpecStatus {
	switch s {
	case feeds.SpecStatusApproved:
		return SpecStatusApproved
	case feeds.SpecStatusPending:
		return SpecStatusPending
	case feeds.SpecStatusRejected:
		return SpecStatusRejected
	case feeds.SpecStatusCancelled:
		return SpecStatusCancelled
	case feeds.SpecStatusRevoked:
		return SpecStatusRevoked
	default:
		return SpecStatusUnknown
	}
}

// JobProposalSpecResolver resolves the Job Proposal Spec type.
type JobProposalSpecResolver struct {
	spec *feeds.JobProposalSpec
}

// NewJobProposalSpec creates a new JobProposalSpecResolver.
func NewJobProposalSpec(spec *feeds.JobProposalSpec) *JobProposalSpecResolver {
	return &JobProposalSpecResolver{spec: spec}
}

// NewJobProposalSpecs creates a slice of JobProposalSpecResolvers.
func NewJobProposalSpecs(specs []feeds.JobProposalSpec) []*JobProposalSpecResolver {
	var resolvers []*JobProposalSpecResolver

	for i := range specs {
		resolvers = append(resolvers, NewJobProposalSpec(&specs[i]))
	}

	return resolvers
}

// ID resolves to the job proposal spec ID
func (r *JobProposalSpecResolver) ID() graphql.ID {
	return int64GQLID(r.spec.ID)
}

// Definition resolves to the job proposal spec definition
func (r *JobProposalSpecResolver) Definition() string {
	return r.spec.Definition
}

// Version resolves to the job proposal spec version
func (r *JobProposalSpecResolver) Version() int32 {
	return r.spec.Version
}

// Status resolves to the job proposal spec's status
func (r *JobProposalSpecResolver) Status() SpecStatus {
	return ToSpecStatus(r.spec.Status)
}

// StatusUpdatedAt resolves to the the last timestamp that the spec status was
// updated.
func (r *JobProposalSpecResolver) StatusUpdatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.StatusUpdatedAt}
}

// CreatedAt resolves to the job proposal spec's created at timestamp
func (r *JobProposalSpecResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.CreatedAt}
}

// UpdatedAt resolves to the job proposal spec's updated at timestamp
func (r *JobProposalSpecResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.spec.UpdatedAt}
}

// -- ApproveJobProposal Mutation --

// ApproveJobProposalSpecPayloadResolver resolves the spec payload.
type ApproveJobProposalSpecPayloadResolver struct {
	spec *feeds.JobProposalSpec
	NotFoundErrorUnionType
}

// NewApproveJobProposalSpecPayload generates the spec payload resolver.
func NewApproveJobProposalSpecPayload(spec *feeds.JobProposalSpec, err error) *ApproveJobProposalSpecPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: notFoundErrorMessage}

	return &ApproveJobProposalSpecPayloadResolver{spec: spec, NotFoundErrorUnionType: e}
}

// ToApproveJobProposalSpecSuccess resolves to the approval job proposal success
// resolver.
func (r *ApproveJobProposalSpecPayloadResolver) ToApproveJobProposalSpecSuccess() (*ApproveJobProposalSpecSuccessResolver, bool) {
	if r.spec != nil {
		return NewApproveJobProposalSpecSuccess(r.spec), true
	}

	return nil, false
}

// ToJobAlreadyExistsError -
func (r *ApproveJobProposalSpecPayloadResolver) ToJobAlreadyExistsError() (*JobAlreadyExistsErrorResolver, bool) {
	if r.err != nil && errors.Is(r.err, feeds.ErrJobAlreadyExists) {
		return NewJobAlreadyExistsError(r.err.Error()), true
	}

	return nil, false
}

// JobAlreadyExistsErrorResolver -
type JobAlreadyExistsErrorResolver struct {
	message string
}

// NewJobAlreadyExistsError -
func NewJobAlreadyExistsError(message string) *JobAlreadyExistsErrorResolver {
	return &JobAlreadyExistsErrorResolver{
		message: message,
	}
}

// Message -
func (r *JobAlreadyExistsErrorResolver) Message() string {
	return r.message
}

// Code -
func (r *JobAlreadyExistsErrorResolver) Code() ErrorCode {
	return ErrorCodeUnprocessable
}

// ApproveJobProposalSpecSuccessResolver resolves the approval success response.
type ApproveJobProposalSpecSuccessResolver struct {
	spec *feeds.JobProposalSpec
}

// NewApproveJobProposalSpecSuccess generates the resolver.
func NewApproveJobProposalSpecSuccess(spec *feeds.JobProposalSpec) *ApproveJobProposalSpecSuccessResolver {
	return &ApproveJobProposalSpecSuccessResolver{spec: spec}
}

// Spec returns the job proposal spec.
func (r *ApproveJobProposalSpecSuccessResolver) Spec() *JobProposalSpecResolver {
	return NewJobProposalSpec(r.spec)
}

// -- CancelJobProposal Mutation --

// CancelJobProposalSpecPayloadResolver resolves the cancel payload response.
type CancelJobProposalSpecPayloadResolver struct {
	spec *feeds.JobProposalSpec
	NotFoundErrorUnionType
}

// NewCancelJobProposalSpecPayload generates the resolver.
func NewCancelJobProposalSpecPayload(spec *feeds.JobProposalSpec, err error) *CancelJobProposalSpecPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: notFoundErrorMessage}

	return &CancelJobProposalSpecPayloadResolver{spec: spec, NotFoundErrorUnionType: e}
}

// ToCancelJobProposalSpecSuccess resolves to the cancel job proposal spec
// success resolver.
func (r *CancelJobProposalSpecPayloadResolver) ToCancelJobProposalSpecSuccess() (*CancelJobProposalSpecSuccessResolver, bool) {
	if r.spec != nil {
		return NewCancelJobProposalSpecSuccess(r.spec), true
	}

	return nil, false
}

// CancelJobProposalSpecSuccessResolver resolves the cancellation success
// response.
type CancelJobProposalSpecSuccessResolver struct {
	spec *feeds.JobProposalSpec
}

// NewCancelJobProposalSpecSuccess generates the resolver.
func NewCancelJobProposalSpecSuccess(spec *feeds.JobProposalSpec) *CancelJobProposalSpecSuccessResolver {
	return &CancelJobProposalSpecSuccessResolver{spec: spec}
}

// Spec returns the job proposal spec.
func (r *CancelJobProposalSpecSuccessResolver) Spec() *JobProposalSpecResolver {
	return NewJobProposalSpec(r.spec)
}

// -- RejectJobProposalSpec Mutation --

// RejectJobProposalSpecPayloadResolver resolves the reject payload response.
type RejectJobProposalSpecPayloadResolver struct {
	spec *feeds.JobProposalSpec
	NotFoundErrorUnionType
}

// NewRejectJobProposalSpecPayload constructs a RejectJobProposalSpecPayloadResolver.
func NewRejectJobProposalSpecPayload(spec *feeds.JobProposalSpec, err error) *RejectJobProposalSpecPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: notFoundErrorMessage}

	return &RejectJobProposalSpecPayloadResolver{spec: spec, NotFoundErrorUnionType: e}
}

// ToRejectJobProposalSpecSuccess resolves to the reject job proposal spec
// success resolver.
func (r *RejectJobProposalSpecPayloadResolver) ToRejectJobProposalSpecSuccess() (*RejectJobProposalSpecSuccessResolver, bool) {
	if r.spec != nil {
		return NewRejectJobProposalSpecSuccess(r.spec), true
	}

	return nil, false
}

// RejectJobProposalSpecSuccessResolver resolves the rejection success response.
type RejectJobProposalSpecSuccessResolver struct {
	spec *feeds.JobProposalSpec
}

// NewRejectJobProposalSpecSuccess generates the resolver.
func NewRejectJobProposalSpecSuccess(spec *feeds.JobProposalSpec) *RejectJobProposalSpecSuccessResolver {
	return &RejectJobProposalSpecSuccessResolver{spec: spec}
}

// Spec returns the job proposal spec.
func (r *RejectJobProposalSpecSuccessResolver) Spec() *JobProposalSpecResolver {
	return NewJobProposalSpec(r.spec)
}

// -- UpdateJobProposalSpecDefinition Mutation --

// UpdateJobProposalSpecDefinitionPayloadResolver generates the update spec
// definition payload.
type UpdateJobProposalSpecDefinitionPayloadResolver struct {
	spec *feeds.JobProposalSpec
	NotFoundErrorUnionType
}

// NewUpdateJobProposalSpecDefinitionPayload constructs UpdateJobProposalSpecDefinitionPayloadResolver.
func NewUpdateJobProposalSpecDefinitionPayload(spec *feeds.JobProposalSpec, err error) *UpdateJobProposalSpecDefinitionPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: notFoundErrorMessage}

	return &UpdateJobProposalSpecDefinitionPayloadResolver{spec: spec, NotFoundErrorUnionType: e}
}

// ToUpdateJobProposalSpecDefinitionSuccess resolves to the update job proposal
// definition success resolver.
func (r *UpdateJobProposalSpecDefinitionPayloadResolver) ToUpdateJobProposalSpecDefinitionSuccess() (*UpdateJobProposalSpecDefinitionSuccessResolver, bool) {
	if r.spec != nil {
		return NewUpdateJobProposalSpecDefinitionSuccess(r.spec), true
	}

	return nil, false
}

// UpdateJobProposalSpecDefinitionSuccessResolver resolves the update success
// response.
type UpdateJobProposalSpecDefinitionSuccessResolver struct {
	spec *feeds.JobProposalSpec
}

// NewUpdateJobProposalSpecDefinitionSuccess constructs UpdateJobProposalSpecDefinitionSuccessResolver.
func NewUpdateJobProposalSpecDefinitionSuccess(spec *feeds.JobProposalSpec) *UpdateJobProposalSpecDefinitionSuccessResolver {
	return &UpdateJobProposalSpecDefinitionSuccessResolver{spec: spec}
}

// Spec returns the job proposal spec.
func (r *UpdateJobProposalSpecDefinitionSuccessResolver) Spec() *JobProposalSpecResolver {
	return NewJobProposalSpec(r.spec)
}
