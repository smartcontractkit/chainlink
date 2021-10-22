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
	resolvers := []*BridgeResolver{}
	for _, b := range bridges {
		resolvers = append(resolvers, NewBridge(b))
	}

	return resolvers
}

// ID resolves the bridge's name.
func (r *BridgeResolver) Name() string {
	return string(r.bridge.Name)
}

// URL resolves the bridge's url.
func (r *BridgeResolver) URL() string {
	return string(r.bridge.URL.String())
}

// Confirmations resolves the bridge's url.
//
// TODO - Fix typing and allow a uint to be returned
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

type CreateBridgePayloadResolver struct {
	bridge bridges.BridgeType
	token  string
}

func NewCreateBridgePayload(bridge bridges.BridgeType, token string) *CreateBridgePayloadResolver {
	return &CreateBridgePayloadResolver{
		bridge: bridge,
		token:  token,
	}
}

// Bridge resolves the payload's bridge.
func (r *CreateBridgePayloadResolver) Bridge() *BridgeResolver {
	return NewBridge(r.bridge)
}

// URL resolves the bridge's url.
func (r *CreateBridgePayloadResolver) Token() string {
	return r.token
}

type UpdateBridgePayloadResolver struct {
	bridge bridges.BridgeType
}

func NewUpdateBridgePayload(bridge bridges.BridgeType) *UpdateBridgePayloadResolver {
	return &UpdateBridgePayloadResolver{
		bridge: bridge,
	}
}

// Bridge resolves the payload's bridge.
func (r *UpdateBridgePayloadResolver) Bridge() *BridgeResolver {
	return NewBridge(r.bridge)
}
