package discovery

import (
	"context"
	"github.com/libp2p/go-libp2p-core/discovery"
	"time"

	cid "github.com/ipfs/go-cid"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"

	mh "github.com/multiformats/go-multihash"
)

// RoutingDiscovery is an implementation of discovery using ContentRouting.
// Namespaces are translated to Cids using the SHA256 hash.
type RoutingDiscovery struct {
	routing.ContentRouting
}

func NewRoutingDiscovery(router routing.ContentRouting) *RoutingDiscovery {
	return &RoutingDiscovery{router}
}

func (d *RoutingDiscovery) Advertise(ctx context.Context, ns string, opts ...Option) (time.Duration, error) {
	var options Options
	err := options.Apply(opts...)
	if err != nil {
		return 0, err
	}

	ttl := options.Ttl
	if ttl == 0 || ttl > 3*time.Hour {
		// the DHT provider record validity is 24hrs, but it is recommnded to republish at least every 6hrs
		// we go one step further and republish every 3hrs
		ttl = 3 * time.Hour
	}

	cid, err := nsToCid(ns)
	if err != nil {
		return 0, err
	}

	// this context requires a timeout; it determines how long the DHT looks for
	// closest peers to the key/CID before it goes on to provide the record to them.
	// Not setting a timeout here will make the DHT wander forever.
	pctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err = d.Provide(pctx, cid, true)
	if err != nil {
		return 0, err
	}

	return ttl, nil
}

func (d *RoutingDiscovery) FindPeers(ctx context.Context, ns string, opts ...Option) (<-chan peer.AddrInfo, error) {
	var options Options
	err := options.Apply(opts...)
	if err != nil {
		return nil, err
	}

	limit := options.Limit
	if limit == 0 {
		limit = 100 // that's just arbitrary, but FindProvidersAsync needs a count
	}

	cid, err := nsToCid(ns)
	if err != nil {
		return nil, err
	}

	return d.FindProvidersAsync(ctx, cid, limit), nil
}

func nsToCid(ns string) (cid.Cid, error) {
	h, err := mh.Sum([]byte(ns), mh.SHA2_256, -1)
	if err != nil {
		return cid.Undef, err
	}

	return cid.NewCidV1(cid.Raw, h), nil
}

func NewDiscoveryRouting(disc discovery.Discovery, opts ...discovery.Option) *DiscoveryRouting {
	return &DiscoveryRouting{disc, opts}
}

type DiscoveryRouting struct {
	discovery.Discovery
	opts []discovery.Option
}

func (r *DiscoveryRouting) Provide(ctx context.Context, c cid.Cid, bcast bool) error {
	if !bcast {
		return nil
	}

	_, err := r.Advertise(ctx, cidToNs(c), r.opts...)
	return err
}

func (r *DiscoveryRouting) FindProvidersAsync(ctx context.Context, c cid.Cid, limit int) <-chan peer.AddrInfo {
	ch, _ := r.FindPeers(ctx, cidToNs(c), append([]discovery.Option{discovery.Limit(limit)}, r.opts...)...)
	return ch
}

func cidToNs(c cid.Cid) string {
	return "/provider/" + c.String()
}
