package discovery

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

// BackoffConnector is a utility to connect to peers, but only if we have not recently tried connecting to them already
type BackoffConnector struct {
	cache      *lru.TwoQueueCache
	host       host.Host
	connTryDur time.Duration
	backoff    BackoffFactory
	mux        sync.Mutex
}

// NewBackoffConnector creates a utility to connect to peers, but only if we have not recently tried connecting to them already
// cacheSize is the size of a TwoQueueCache
// connectionTryDuration is how long we attempt to connect to a peer before giving up
// backoff describes the strategy used to decide how long to backoff after previously attempting to connect to a peer
func NewBackoffConnector(h host.Host, cacheSize int, connectionTryDuration time.Duration, backoff BackoffFactory) (*BackoffConnector, error) {
	cache, err := lru.New2Q(cacheSize)
	if err != nil {
		return nil, err
	}

	return &BackoffConnector{
		cache:      cache,
		host:       h,
		connTryDur: connectionTryDuration,
		backoff:    backoff,
	}, nil
}

type connCacheData struct {
	nextTry time.Time
	strat   BackoffStrategy
}

// Connect attempts to connect to the peers passed in by peerCh. Will not connect to peers if they are within the backoff period.
// As Connect will attempt to dial peers as soon as it learns about them, the caller should try to keep the number,
// and rate, of inbound peers manageable.
func (c *BackoffConnector) Connect(ctx context.Context, peerCh <-chan peer.AddrInfo) {
	for {
		select {
		case pi, ok := <-peerCh:
			if !ok {
				return
			}

			if pi.ID == c.host.ID() || pi.ID == "" {
				continue
			}

			c.mux.Lock()
			val, ok := c.cache.Get(pi.ID)
			var cachedPeer *connCacheData
			if ok {
				tv := val.(*connCacheData)
				now := time.Now()
				if now.Before(tv.nextTry) {
					c.mux.Unlock()
					continue
				}

				tv.nextTry = now.Add(tv.strat.Delay())
			} else {
				cachedPeer = &connCacheData{strat: c.backoff()}
				cachedPeer.nextTry = time.Now().Add(cachedPeer.strat.Delay())
				c.cache.Add(pi.ID, cachedPeer)
			}
			c.mux.Unlock()

			go func(pi peer.AddrInfo) {
				ctx, cancel := context.WithTimeout(ctx, c.connTryDur)
				defer cancel()

				err := c.host.Connect(ctx, pi)
				if err != nil {
					log.Debugf("Error connecting to pubsub peer %s: %s", pi.ID, err.Error())
					return
				}
			}(pi)

		case <-ctx.Done():
			log.Infof("discovery: backoff connector context error %v", ctx.Err())
			return
		}
	}
}
