package resolver

import (
	"strconv"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
)

// FeedsManagerResolver resolves the FeedsManager type.
type FeedsManagerResolver struct {
	mgr feeds.FeedsManager
}

func NewFeedsManager(mgr feeds.FeedsManager) *FeedsManagerResolver {
	return &FeedsManagerResolver{mgr: mgr}
}

func NewFeedsManagers(mgrs []feeds.FeedsManager) []*FeedsManagerResolver {
	resolvers := []*FeedsManagerResolver{}
	for _, mgr := range mgrs {
		resolvers = append(resolvers, NewFeedsManager(mgr))
	}

	return resolvers
}

// ID resolves the feed managers's unique identifier.
func (r *FeedsManagerResolver) ID() graphql.ID {
	return graphql.ID(strconv.FormatInt(r.mgr.ID, 10))
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

// JobTypes resolves the feed managers's jobTypes field.
func (r *FeedsManagerResolver) JobTypes() []string {
	return r.mgr.JobTypes
}

// IsBootstrapPeer resolves the feed managers's isBootstrapPeer field.
func (r *FeedsManagerResolver) IsBootstrapPeer() bool {
	return r.mgr.IsOCRBootstrapPeer
}

// IsConnectionActive resolves the feed managers's isConnectionActive field.
func (r *FeedsManagerResolver) IsConnectionActive() bool {
	return r.mgr.IsConnectionActive
}

// BootstrapPeer resolves the feed managers's isConnectionActive field.
func (r *FeedsManagerResolver) BootstrapPeerMultiaddr() *string {
	return r.mgr.OCRBootstrapPeerMultiaddr.Ptr()
}

// CreatedAt resolves the chains's created at field.
func (r *FeedsManagerResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.mgr.CreatedAt}
}
