package discovery

import (
	"time"

	core "github.com/libp2p/go-libp2p-core/discovery"
)

// Deprecated: use skel.Advertiser instead.
type Advertiser = core.Advertiser

// Deprecated: use skel.Discoverer instead.
type Discoverer = core.Discoverer

// Deprecated: use skel.Discovery instead.
type Discovery = core.Discovery

// Deprecated: use github.com/libp2p/go-libp2p-core/discovery.Option instead.
type Option = core.Option

// Deprecated: use github.com/libp2p/go-libp2p-core/discovery.Options instead.
type Options = core.Options

// Deprecated: use github.com/libp2p/go-libp2p-core/discovery.TTL instead.
func TTL(ttl time.Duration) core.Option {
	return core.TTL(ttl)
}

// Deprecated: use github.com/libp2p/go-libp2p-core/discovery.Limit instead.
func Limit(limit int) core.Option {
	return core.Limit(limit)
}
