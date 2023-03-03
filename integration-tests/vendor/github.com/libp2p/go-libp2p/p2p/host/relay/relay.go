package relay

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-discovery"

	ma "github.com/multiformats/go-multiaddr"
)

var (
	// this is purposefully long to require some node stability before advertising as a relay
	AdvertiseBootDelay = 15 * time.Minute
	AdvertiseTTL       = 30 * time.Minute
)

// Advertise advertises this node as a libp2p relay.
func Advertise(ctx context.Context, advertise discovery.Advertiser) {
	go func() {
		select {
		case <-time.After(AdvertiseBootDelay):
			discovery.Advertise(ctx, advertise, RelayRendezvous, discovery.TTL(AdvertiseTTL))
		case <-ctx.Done():
		}
	}()
}

// Filter filters out all relay addresses.
func Filter(addrs []ma.Multiaddr) []ma.Multiaddr {
	raddrs := make([]ma.Multiaddr, 0, len(addrs))
	for _, addr := range addrs {
		if isRelayAddr(addr) {
			continue
		}
		raddrs = append(raddrs, addr)
	}
	return raddrs
}
