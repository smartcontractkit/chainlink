package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/bridges"
)

// BridgeResolver resolves the Bridge type.
type BridgeResolver struct {
	bridge bridges.BridgeType
}

func NewBridge(bridge bridges.BridgeType) *BridgeResolver {
	return &BridgeResolver{bridge: bridge}
}

func NewBridges(bridges []bridges.BridgeType) []*BridgeResolver {
	var resolvers []*BridgeResolver
	for _, b := range bridges {
		resolvers = append(resolvers, NewBridge(b))
	}

	return resolvers
}

// ID resolves the bridge's name as the id (Bridge does not have an id).
func (r *BridgeResolver) ID() graphql.ID {
	return graphql.ID(r.bridge.Name.String())
}

// Name resolves the bridge's name.
func (r *BridgeResolver) Name() string {
	return r.bridge.Name.String()
}

// URL resolves the bridge's url.
func (r *BridgeResolver) URL() string {
	return r.bridge.URL.String()
}

// Confirmations resolves the bridge's url.
func (r *BridgeResolver) Confirmations() int32 {
	return int32(r.bridge.Confirmations)
}

// OutgoingToken resolves the bridge's outgoing token.
func (r *BridgeResolver) OutgoingToken() string {
	return r.bridge.OutgoingToken
}

// MinimumContractPayment resolves the bridge's minimum contract payment.
func (r *BridgeResolver) MinimumContractPayment() string {
	return r.bridge.MinimumContractPayment.String()
}

// CreatedAt resolves the bridge's created at field.
func (r *BridgeResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.bridge.CreatedAt}
}

// BridgePayloadResolver resolves a single bridge response
type BridgePayloadResolver struct {
	bridge bridges.BridgeType
	NotFoundErrorUnionType
}

func NewBridgePayload(bridge bridges.BridgeType, err error) *BridgePayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "bridge not found"}

	return &BridgePayloadResolver{bridge: bridge, NotFoundErrorUnionType: e}
}

// ToBridge implements the Bridge union type of the payload
func (r *BridgePayloadResolver) ToBridge() (*BridgeResolver, bool) {
	if r.err == nil {
		return NewBridge(r.bridge), true
	}

	return nil, false
}

// BridgesPayloadResolver resolves a page of bridges
type BridgesPayloadResolver struct {
	bridges []bridges.BridgeType
	total   int32
}

func NewBridgesPayload(bridges []bridges.BridgeType, total int32) *BridgesPayloadResolver {
	return &BridgesPayloadResolver{
		bridges: bridges,
		total:   total,
	}
}

// Results returns the bridges.
func (r *BridgesPayloadResolver) Results() []*BridgeResolver {
	return NewBridges(r.bridges)
}

// Metadata returns the pagination metadata.
func (r *BridgesPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

// CreateBridgePayloadResolver
type CreateBridgePayloadResolver struct {
	bridge        bridges.BridgeType
	incomingToken string
}

func NewCreateBridgePayload(bridge bridges.BridgeType, incomingToken string) *CreateBridgePayloadResolver {
	return &CreateBridgePayloadResolver{
		bridge:        bridge,
		incomingToken: incomingToken,
	}
}

func (r *CreateBridgePayloadResolver) ToCreateBridgeSuccess() (*CreateBridgeSuccessResolver, bool) {
	return NewCreateBridgeSuccessResolver(r.bridge, r.incomingToken), true
}

type CreateBridgeSuccessResolver struct {
	bridge        bridges.BridgeType
	incomingToken string
}

func NewCreateBridgeSuccessResolver(bridge bridges.BridgeType, incomingToken string) *CreateBridgeSuccessResolver {
	return &CreateBridgeSuccessResolver{
		bridge:        bridge,
		incomingToken: incomingToken,
	}
}

// Bridge resolves the bridge.
func (r *CreateBridgeSuccessResolver) Bridge() *BridgeResolver {
	return NewBridge(r.bridge)
}

// Token resolves the bridge's incoming token.
func (r *CreateBridgeSuccessResolver) IncomingToken() string {
	return r.incomingToken
}

type UpdateBridgePayloadResolver struct {
	bridge *bridges.BridgeType
	NotFoundErrorUnionType
}

func NewUpdateBridgePayload(bridge *bridges.BridgeType, err error) *UpdateBridgePayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "bridge not found"}

	return &UpdateBridgePayloadResolver{bridge: bridge, NotFoundErrorUnionType: e}
}

func (r *UpdateBridgePayloadResolver) ToUpdateBridgeSuccess() (*UpdateBridgeSuccessResolver, bool) {
	if r.bridge != nil {
		return NewUpdateBridgeSuccess(*r.bridge), true
	}

	return nil, false
}

// UpdateBridgePayloadResolver resolves
type UpdateBridgeSuccessResolver struct {
	bridge bridges.BridgeType
}

func NewUpdateBridgeSuccess(bridge bridges.BridgeType) *UpdateBridgeSuccessResolver {
	return &UpdateBridgeSuccessResolver{
		bridge: bridge,
	}
}

// Bridge resolves the success payload's bridge.
func (r *UpdateBridgeSuccessResolver) Bridge() *BridgeResolver {
	return NewBridge(r.bridge)
}

// -- DeleteBridge mutation --

type DeleteBridgePayloadResolver struct {
	bridge *bridges.BridgeType
	NotFoundErrorUnionType
}

func NewDeleteBridgePayload(bridge *bridges.BridgeType, err error) *DeleteBridgePayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "bridge not found"}

	return &DeleteBridgePayloadResolver{bridge: bridge, NotFoundErrorUnionType: e}
}

func (r *DeleteBridgePayloadResolver) ToDeleteBridgeSuccess() (*DeleteBridgeSuccessResolver, bool) {
	if r.bridge != nil {
		return NewDeleteBridgeSuccess(r.bridge), true
	}

	return nil, false
}

func (r *DeleteBridgePayloadResolver) ToDeleteBridgeConflictError() (*DeleteBridgeConflictErrorResolver, bool) {
	if r.err != nil {
		return NewDeleteBridgeConflictError(r.err.Error()), true
	}

	return nil, false
}

func (r *DeleteBridgePayloadResolver) ToDeleteBridgeInvalidNameError() (*DeleteBridgeInvalidNameErrorResolver, bool) {
	if r.err != nil {
		return NewDeleteBridgeInvalidNameError(r.err.Error()), true
	}

	return nil, false
}

type DeleteBridgeSuccessResolver struct {
	bridge *bridges.BridgeType
}

func NewDeleteBridgeSuccess(bridge *bridges.BridgeType) *DeleteBridgeSuccessResolver {
	return &DeleteBridgeSuccessResolver{bridge: bridge}
}

func (r *DeleteBridgeSuccessResolver) Bridge() *BridgeResolver {
	return NewBridge(*r.bridge)
}

type DeleteBridgeConflictErrorResolver struct {
	message string
}

func NewDeleteBridgeConflictError(message string) *DeleteBridgeConflictErrorResolver {
	return &DeleteBridgeConflictErrorResolver{message: message}
}

func (r *DeleteBridgeConflictErrorResolver) Message() string {
	return r.message
}

func (r *DeleteBridgeConflictErrorResolver) Code() ErrorCode {
	return ErrorCodeStatusConflict
}

type DeleteBridgeInvalidNameErrorResolver struct {
	message string
}

func NewDeleteBridgeInvalidNameError(message string) *DeleteBridgeInvalidNameErrorResolver {
	return &DeleteBridgeInvalidNameErrorResolver{message: message}
}

func (r *DeleteBridgeInvalidNameErrorResolver) Message() string {
	return r.message
}

func (r *DeleteBridgeInvalidNameErrorResolver) Code() ErrorCode {
	return ErrorCodeUnprocessable
}
