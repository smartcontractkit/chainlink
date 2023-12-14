package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

type FeedsManagerChainConfigResolver struct {
	cfg feeds.ChainConfig
}

func NewFeedsManagerChainConfig(cfg feeds.ChainConfig) *FeedsManagerChainConfigResolver {
	return &FeedsManagerChainConfigResolver{cfg: cfg}
}

func NewFeedsManagerChainConfigs(cfgs []feeds.ChainConfig) []*FeedsManagerChainConfigResolver {
	var resolvers []*FeedsManagerChainConfigResolver
	for _, cfg := range cfgs {
		resolvers = append(resolvers, NewFeedsManagerChainConfig(cfg))
	}

	return resolvers
}

// ID resolves the chain configs's unique identifier.
func (r *FeedsManagerChainConfigResolver) ID() graphql.ID {
	return int64GQLID(r.cfg.ID)
}

// ChainID resolves the chain configs's chain id.
func (r *FeedsManagerChainConfigResolver) ChainID() string {
	return r.cfg.ChainID
}

// ChainType resolves the chain configs's chain type.
func (r *FeedsManagerChainConfigResolver) ChainType() string {
	return string(r.cfg.ChainType)
}

// AccountAddr resolves the chain configs's account address.
func (r *FeedsManagerChainConfigResolver) AccountAddr() string {
	return r.cfg.AccountAddress
}

// AdminAddr resolves the chain configs's admin address.
func (r *FeedsManagerChainConfigResolver) AdminAddr() string {
	return r.cfg.AdminAddress
}

// FluxMonitorJobConfig resolves the chain configs's Flux Monitor Config.
func (r *FeedsManagerChainConfigResolver) FluxMonitorJobConfig() *FluxMonitorJobConfigResolver {
	return &FluxMonitorJobConfigResolver{cfg: r.cfg.FluxMonitorConfig}
}

// OCR1JobConfig resolves the chain configs's OCR1 Config.
func (r *FeedsManagerChainConfigResolver) OCR1JobConfig() *OCR1JobConfigResolver {
	return &OCR1JobConfigResolver{cfg: r.cfg.OCR1Config}
}

// OCR2JobConfig resolves the chain configs's OCR2 Config.
func (r *FeedsManagerChainConfigResolver) OCR2JobConfig() *OCR2JobConfigResolver {
	return &OCR2JobConfigResolver{cfg: r.cfg.OCR2Config}
}

type FluxMonitorJobConfigResolver struct {
	cfg feeds.FluxMonitorConfig
}

func (r *FluxMonitorJobConfigResolver) Enabled() bool {
	return r.cfg.Enabled
}

type OCR1JobConfigResolver struct {
	cfg feeds.OCR1Config
}

func (r *OCR1JobConfigResolver) Enabled() bool {
	return r.cfg.Enabled
}

func (r *OCR1JobConfigResolver) IsBootstrap() bool {
	return r.cfg.IsBootstrap
}

func (r *OCR1JobConfigResolver) Multiaddr() *string {
	return r.cfg.Multiaddr.Ptr()
}

func (r *OCR1JobConfigResolver) P2PPeerID() *string {
	return r.cfg.P2PPeerID.Ptr()
}

func (r *OCR1JobConfigResolver) KeyBundleID() *string {
	return r.cfg.KeyBundleID.Ptr()
}

type OCR2JobConfigResolver struct {
	cfg feeds.OCR2ConfigModel
}

func (r *OCR2JobConfigResolver) Enabled() bool {
	return r.cfg.Enabled
}

func (r *OCR2JobConfigResolver) IsBootstrap() bool {
	return r.cfg.IsBootstrap
}

func (r *OCR2JobConfigResolver) Multiaddr() *string {
	return r.cfg.Multiaddr.Ptr()
}

func (r *OCR2JobConfigResolver) ForwarderAddress() *string {
	return r.cfg.ForwarderAddress.Ptr()
}

func (r *OCR2JobConfigResolver) P2PPeerID() *string {
	return r.cfg.P2PPeerID.Ptr()
}

func (r *OCR2JobConfigResolver) KeyBundleID() *string {
	return r.cfg.KeyBundleID.Ptr()
}

func (r *OCR2JobConfigResolver) Plugins() *PluginsResolver {
	return &PluginsResolver{plugins: r.cfg.Plugins}
}

// -- CreateFeedsManagerChainConfig Mutation --

// CreateFeedsManagerChainConfigPayloadResolver resolves the response to
// CreateFeedsManagerChainConfig
type CreateFeedsManagerChainConfigPayloadResolver struct {
	cfg *feeds.ChainConfig
	// inputErrors maps an input path to a string
	inputErrs map[string]string
	NotFoundErrorUnionType
}

func NewCreateFeedsManagerChainConfigPayload(cfg *feeds.ChainConfig, err error, inputErrs map[string]string) *CreateFeedsManagerChainConfigPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "chain config not found", isExpectedErrorFn: nil}

	return &CreateFeedsManagerChainConfigPayloadResolver{
		cfg:                    cfg,
		inputErrs:              inputErrs,
		NotFoundErrorUnionType: e,
	}
}

func (r *CreateFeedsManagerChainConfigPayloadResolver) ToCreateFeedsManagerChainConfigSuccess() (*CreateFeedsManagerChainConfigSuccessResolver, bool) {
	if r.cfg != nil {
		return NewCreateFeedsManagerChainConfigSuccessResolver(r.cfg), true
	}

	return nil, false
}

func (r *CreateFeedsManagerChainConfigPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type CreateFeedsManagerChainConfigSuccessResolver struct {
	cfg *feeds.ChainConfig
}

func NewCreateFeedsManagerChainConfigSuccessResolver(cfg *feeds.ChainConfig) *CreateFeedsManagerChainConfigSuccessResolver {
	return &CreateFeedsManagerChainConfigSuccessResolver{
		cfg: cfg,
	}
}

func (r *CreateFeedsManagerChainConfigSuccessResolver) ChainConfig() *FeedsManagerChainConfigResolver {
	return NewFeedsManagerChainConfig(*r.cfg)
}

// -- Delete FMS Chain Config --

type DeleteFeedsManagerChainConfigPayloadResolver struct {
	cfg *feeds.ChainConfig
	NotFoundErrorUnionType
}

func NewDeleteFeedsManagerChainConfigPayload(cfg *feeds.ChainConfig, err error) *DeleteFeedsManagerChainConfigPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "chain config not found", isExpectedErrorFn: nil}

	return &DeleteFeedsManagerChainConfigPayloadResolver{cfg: cfg, NotFoundErrorUnionType: e}
}

func (r *DeleteFeedsManagerChainConfigPayloadResolver) ToDeleteFeedsManagerChainConfigSuccess() (*DeleteFeedsManagerChainConfigSuccessResolver, bool) {
	if r.cfg == nil {
		return nil, false
	}

	return NewDeleteFeedsManagerChainConfigSuccess(*r.cfg), true
}

type DeleteFeedsManagerChainConfigSuccessResolver struct {
	cfg feeds.ChainConfig
}

func NewDeleteFeedsManagerChainConfigSuccess(cfg feeds.ChainConfig) *DeleteFeedsManagerChainConfigSuccessResolver {
	return &DeleteFeedsManagerChainConfigSuccessResolver{cfg: cfg}
}

func (r *DeleteFeedsManagerChainConfigSuccessResolver) ChainConfig() *FeedsManagerChainConfigResolver {
	return NewFeedsManagerChainConfig(r.cfg)
}

// -- UpdateFeedsManagerChainConfig Mutation --

// UpdateFeedsManagerChainConfigPayloadResolver resolves the response to
// UpdateFeedsManagerChainConfig
type UpdateFeedsManagerChainConfigPayloadResolver struct {
	cfg *feeds.ChainConfig
	// inputErrors maps an input path to a string
	inputErrs map[string]string
	NotFoundErrorUnionType
}

func NewUpdateFeedsManagerChainConfigPayload(cfg *feeds.ChainConfig, err error, inputErrs map[string]string) *UpdateFeedsManagerChainConfigPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "chain config not found", isExpectedErrorFn: nil}

	return &UpdateFeedsManagerChainConfigPayloadResolver{
		cfg:                    cfg,
		inputErrs:              inputErrs,
		NotFoundErrorUnionType: e,
	}
}

func (r *UpdateFeedsManagerChainConfigPayloadResolver) ToUpdateFeedsManagerChainConfigSuccess() (*UpdateFeedsManagerChainConfigSuccessResolver, bool) {
	if r.cfg != nil {
		return NewUpdateFeedsManagerChainConfigSuccessResolver(r.cfg), true
	}

	return nil, false
}

func (r *UpdateFeedsManagerChainConfigPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type UpdateFeedsManagerChainConfigSuccessResolver struct {
	cfg *feeds.ChainConfig
}

func NewUpdateFeedsManagerChainConfigSuccessResolver(cfg *feeds.ChainConfig) *UpdateFeedsManagerChainConfigSuccessResolver {
	return &UpdateFeedsManagerChainConfigSuccessResolver{
		cfg: cfg,
	}
}

func (r *UpdateFeedsManagerChainConfigSuccessResolver) ChainConfig() *FeedsManagerChainConfigResolver {
	return NewFeedsManagerChainConfig(*r.cfg)
}
