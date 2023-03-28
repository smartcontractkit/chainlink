package resolver

import (
	"context"
	"errors"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/v2/core/web/loader"
)

// FeedsManagerResolver resolves the FeedsManager type.
type FeedsManagerResolver struct {
	mgr feeds.FeedsManager
}

func NewFeedsManager(mgr feeds.FeedsManager) *FeedsManagerResolver {
	return &FeedsManagerResolver{mgr: mgr}
}

func NewFeedsManagers(mgrs []feeds.FeedsManager) []*FeedsManagerResolver {
	var resolvers []*FeedsManagerResolver
	for _, mgr := range mgrs {
		resolvers = append(resolvers, NewFeedsManager(mgr))
	}

	return resolvers
}

// ID resolves the feed managers's unique identifier.
func (r *FeedsManagerResolver) ID() graphql.ID {
	return int64GQLID(r.mgr.ID)
}

// Name resolves the feed managers's name field.
func (r *FeedsManagerResolver) Name() string {
	return r.mgr.Name
}

// URI resolves the feed managers's uri field.
func (r *FeedsManagerResolver) URI() string {
	return r.mgr.URI
}

// PublicKey resolves the feed managers's public key field.
func (r *FeedsManagerResolver) PublicKey() string {
	return r.mgr.PublicKey.String()
}

func (r *FeedsManagerResolver) JobProposals(ctx context.Context) ([]*JobProposalResolver, error) {
	jps, err := loader.GetJobProposalsByFeedsManagerID(ctx, stringutils.FromInt64(r.mgr.ID))
	if err != nil {
		return nil, err
	}

	return NewJobProposals(jps), nil
}

// IsConnectionActive resolves the feed managers's isConnectionActive field.
func (r *FeedsManagerResolver) IsConnectionActive() bool {
	return r.mgr.IsConnectionActive
}

func (r *FeedsManagerResolver) ChainConfigs(ctx context.Context) ([]*FeedsManagerChainConfigResolver, error) {
	cfgs, err := loader.GetFeedsManagerChainConfigsByManagerID(ctx, r.mgr.ID)
	if err != nil {
		return nil, err
	}

	return NewFeedsManagerChainConfigs(cfgs), nil
}

// CreatedAt resolves the chains's created at field.
func (r *FeedsManagerResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.mgr.CreatedAt}
}

// -- FeedsManager Query --

type FeedsManagerPayloadResolver struct {
	mgr *feeds.FeedsManager
	NotFoundErrorUnionType
}

func NewFeedsManagerPayload(mgr *feeds.FeedsManager, err error) *FeedsManagerPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "feeds manager not found", isExpectedErrorFn: nil}

	return &FeedsManagerPayloadResolver{mgr: mgr, NotFoundErrorUnionType: e}
}

// ToFeedsManager implements the FeedsManager union type of the payload
func (r *FeedsManagerPayloadResolver) ToFeedsManager() (*FeedsManagerResolver, bool) {
	if r.mgr != nil {
		return NewFeedsManager(*r.mgr), true
	}

	return nil, false
}

// -- FeedsManagers Query --

// FeedsManagersPayloadResolver resolves a list of feeds managers
type FeedsManagersPayloadResolver struct {
	feedsManagers []feeds.FeedsManager
}

func NewFeedsManagersPayload(feedsManagers []feeds.FeedsManager) *FeedsManagersPayloadResolver {
	return &FeedsManagersPayloadResolver{
		feedsManagers: feedsManagers,
	}
}

// Results returns the feeds managers.
func (r *FeedsManagersPayloadResolver) Results() []*FeedsManagerResolver {
	return NewFeedsManagers(r.feedsManagers)
}

// -- CreateFeedsManager Mutation --

// CreateFeedsManagerPayloadResolver -
type CreateFeedsManagerPayloadResolver struct {
	mgr *feeds.FeedsManager
	// inputErrors maps an input path to a string
	inputErrs map[string]string
	NotFoundErrorUnionType
}

func NewCreateFeedsManagerPayload(mgr *feeds.FeedsManager, err error, inputErrs map[string]string) *CreateFeedsManagerPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "feeds manager not found", isExpectedErrorFn: nil}

	return &CreateFeedsManagerPayloadResolver{
		mgr:                    mgr,
		inputErrs:              inputErrs,
		NotFoundErrorUnionType: e,
	}
}

func (r *CreateFeedsManagerPayloadResolver) ToCreateFeedsManagerSuccess() (*CreateFeedsManagerSuccessResolver, bool) {
	if r.mgr != nil {
		return NewCreateFeedsManagerSuccessResolver(*r.mgr), true
	}

	return nil, false
}

func (r *CreateFeedsManagerPayloadResolver) ToSingleFeedsManagerError() (*SingleFeedsManagerErrorResolver, bool) {
	if r.err != nil && errors.Is(r.err, feeds.ErrSingleFeedsManager) {
		return NewSingleFeedsManagerError(r.err.Error()), true
	}

	return nil, false
}

func (r *CreateFeedsManagerPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type CreateFeedsManagerSuccessResolver struct {
	mgr feeds.FeedsManager
}

func NewCreateFeedsManagerSuccessResolver(mgr feeds.FeedsManager) *CreateFeedsManagerSuccessResolver {
	return &CreateFeedsManagerSuccessResolver{
		mgr: mgr,
	}
}

func (r *CreateFeedsManagerSuccessResolver) FeedsManager() *FeedsManagerResolver {
	return NewFeedsManager(r.mgr)
}

// SingleFeedsManagerErrorResolver -
type SingleFeedsManagerErrorResolver struct {
	message string
}

func NewSingleFeedsManagerError(message string) *SingleFeedsManagerErrorResolver {
	return &SingleFeedsManagerErrorResolver{
		message: message,
	}
}

func (r *SingleFeedsManagerErrorResolver) Message() string {
	return r.message
}

func (r *SingleFeedsManagerErrorResolver) Code() ErrorCode {
	return ErrorCodeUnprocessable
}

// -- UpdateFeedsManager Mutation --

// UpdateFeedsManagerPayloadResolver -
type UpdateFeedsManagerPayloadResolver struct {
	mgr       *feeds.FeedsManager
	inputErrs map[string]string
	NotFoundErrorUnionType
}

func NewUpdateFeedsManagerPayload(mgr *feeds.FeedsManager, err error, inputErrs map[string]string) *UpdateFeedsManagerPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "feeds manager not found", isExpectedErrorFn: nil}

	return &UpdateFeedsManagerPayloadResolver{
		mgr:                    mgr,
		inputErrs:              inputErrs,
		NotFoundErrorUnionType: e,
	}
}

func (r *UpdateFeedsManagerPayloadResolver) ToUpdateFeedsManagerSuccess() (*UpdateFeedsManagerSuccessResolver, bool) {
	if r.mgr != nil {
		return NewUpdateFeedsManagerSuccessResolver(*r.mgr), true
	}

	return nil, false
}

func (r *UpdateFeedsManagerPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type UpdateFeedsManagerSuccessResolver struct {
	mgr feeds.FeedsManager
}

func NewUpdateFeedsManagerSuccessResolver(mgr feeds.FeedsManager) *UpdateFeedsManagerSuccessResolver {
	return &UpdateFeedsManagerSuccessResolver{
		mgr: mgr,
	}
}

func (r *UpdateFeedsManagerSuccessResolver) FeedsManager() *FeedsManagerResolver {
	return NewFeedsManager(r.mgr)
}
