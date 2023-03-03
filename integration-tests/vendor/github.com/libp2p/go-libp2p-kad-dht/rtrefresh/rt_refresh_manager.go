package rtrefresh

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"

	kbucket "github.com/libp2p/go-libp2p-kbucket"

	"github.com/hashicorp/go-multierror"
	logging "github.com/ipfs/go-log"
	"github.com/multiformats/go-base32"
)

var logger = logging.Logger("dht/RtRefreshManager")

const (
	peerPingTimeout = 10 * time.Second
)

type triggerRefreshReq struct {
	respCh          chan error
	forceCplRefresh bool
}

type RtRefreshManager struct {
	ctx       context.Context
	cancel    context.CancelFunc
	refcount  sync.WaitGroup
	closeOnce sync.Once

	// peerId of this DHT peer i.e. self peerId.
	h         host.Host
	dhtPeerId peer.ID
	rt        *kbucket.RoutingTable

	enableAutoRefresh   bool                                        // should run periodic refreshes ?
	refreshKeyGenFnc    func(cpl uint) (string, error)              // generate the key for the query to refresh this cpl
	refreshQueryFnc     func(ctx context.Context, key string) error // query to run for a refresh.
	refreshQueryTimeout time.Duration                               // timeout for one refresh query

	// interval between two periodic refreshes.
	// also, a cpl wont be refreshed if the time since it was last refreshed
	// is below the interval..unless a "forced" refresh is done.
	refreshInterval                    time.Duration
	successfulOutboundQueryGracePeriod time.Duration

	triggerRefresh chan *triggerRefreshReq // channel to write refresh requests to.

	refreshDoneCh chan struct{} // write to this channel after every refresh
}

func NewRtRefreshManager(h host.Host, rt *kbucket.RoutingTable, autoRefresh bool,
	refreshKeyGenFnc func(cpl uint) (string, error),
	refreshQueryFnc func(ctx context.Context, key string) error,
	refreshQueryTimeout time.Duration,
	refreshInterval time.Duration,
	successfulOutboundQueryGracePeriod time.Duration,
	refreshDoneCh chan struct{}) (*RtRefreshManager, error) {

	ctx, cancel := context.WithCancel(context.Background())
	return &RtRefreshManager{
		ctx:       ctx,
		cancel:    cancel,
		h:         h,
		dhtPeerId: h.ID(),
		rt:        rt,

		enableAutoRefresh: autoRefresh,
		refreshKeyGenFnc:  refreshKeyGenFnc,
		refreshQueryFnc:   refreshQueryFnc,

		refreshQueryTimeout:                refreshQueryTimeout,
		refreshInterval:                    refreshInterval,
		successfulOutboundQueryGracePeriod: successfulOutboundQueryGracePeriod,

		triggerRefresh: make(chan *triggerRefreshReq),
		refreshDoneCh:  refreshDoneCh,
	}, nil
}

func (r *RtRefreshManager) Start() error {
	r.refcount.Add(1)
	go r.loop()
	return nil
}

func (r *RtRefreshManager) Close() error {
	r.closeOnce.Do(func() {
		r.cancel()
		r.refcount.Wait()
	})
	return nil
}

// RefreshRoutingTable requests the refresh manager to refresh the Routing Table.
// If the force parameter is set to true true, all buckets will be refreshed irrespective of when they were last refreshed.
//
// The returned channel will block until the refresh finishes, then yield the
// error and close. The channel is buffered and safe to ignore.
// FIXME: this can block. Ideally, we'd return a channel without blocking.
// https://github.com/libp2p/go-libp2p-kad-dht/issues/609
func (r *RtRefreshManager) Refresh(force bool) <-chan error {
	resp := make(chan error, 1)
	select {
	case r.triggerRefresh <- &triggerRefreshReq{respCh: resp, forceCplRefresh: force}:
	case <-r.ctx.Done():
		resp <- r.ctx.Err()
	}
	return resp
}

// RefreshNoWait requests the refresh manager to refresh the Routing Table.
// However, it moves on without blocking if it's request can't get through.
func (r *RtRefreshManager) RefreshNoWait() {
	select {
	case r.triggerRefresh <- &triggerRefreshReq{}:
	default:
	}
}

func (r *RtRefreshManager) loop() {
	defer r.refcount.Done()

	var refreshTickrCh <-chan time.Time
	if r.enableAutoRefresh {
		err := r.doRefresh(true)
		if err != nil {
			logger.Warn("failed when refreshing routing table", err)
		}
		t := time.NewTicker(r.refreshInterval)
		defer t.Stop()
		refreshTickrCh = t.C
	}

	for {
		var waiting []chan<- error
		var forced bool
		select {
		case <-refreshTickrCh:
		case triggerRefreshReq := <-r.triggerRefresh:
			if triggerRefreshReq.respCh != nil {
				waiting = append(waiting, triggerRefreshReq.respCh)
			}
			forced = forced || triggerRefreshReq.forceCplRefresh
		case <-r.ctx.Done():
			return
		}

		// Batch multiple refresh requests if they're all waiting at the same time.
	OuterLoop:
		for {
			select {
			case triggerRefreshReq := <-r.triggerRefresh:
				if triggerRefreshReq.respCh != nil {
					waiting = append(waiting, triggerRefreshReq.respCh)
				}
				forced = forced || triggerRefreshReq.forceCplRefresh
			default:
				break OuterLoop
			}
		}

		// EXECUTE the refresh

		// ping Routing Table peers that haven't been heard of/from in the interval they should have been.
		// and evict them if they don't reply.
		var wg sync.WaitGroup
		for _, ps := range r.rt.GetPeerInfos() {
			if time.Since(ps.LastSuccessfulOutboundQueryAt) > r.successfulOutboundQueryGracePeriod {
				wg.Add(1)
				go func(ps kbucket.PeerInfo) {
					defer wg.Done()
					livelinessCtx, cancel := context.WithTimeout(r.ctx, peerPingTimeout)
					if err := r.h.Connect(livelinessCtx, peer.AddrInfo{ID: ps.Id}); err != nil {
						logger.Debugw("evicting peer after failed ping", "peer", ps.Id, "error", err)
						r.rt.RemovePeer(ps.Id)
					}
					cancel()
				}(ps)
			}
		}
		wg.Wait()

		// Query for self and refresh the required buckets
		err := r.doRefresh(forced)
		for _, w := range waiting {
			w <- err
			close(w)
		}
		if err != nil {
			logger.Warnw("failed when refreshing routing table", "error", err)
		}
	}
}

func (r *RtRefreshManager) doRefresh(forceRefresh bool) error {
	var merr error

	if err := r.queryForSelf(); err != nil {
		merr = multierror.Append(merr, err)
	}

	refreshCpls := r.rt.GetTrackedCplsForRefresh()

	rfnc := func(cpl uint) (err error) {
		if forceRefresh {
			err = r.refreshCpl(cpl)
		} else {
			err = r.refreshCplIfEligible(cpl, refreshCpls[cpl])
		}
		return
	}

	for c := range refreshCpls {
		cpl := uint(c)
		if err := rfnc(cpl); err != nil {
			merr = multierror.Append(merr, err)
		} else {
			// If we see a gap at a Cpl in the Routing table, we ONLY refresh up until the maximum cpl we
			// have in the Routing Table OR (2 * (Cpl+ 1) with the gap), whichever is smaller.
			// This is to prevent refreshes for Cpls that have no peers in the network but happen to be before a very high max Cpl
			// for which we do have peers in the network.
			// The number of 2 * (Cpl + 1) can be proved and a proof would have been written here if the programmer
			// had paid more attention in the Math classes at university.
			// So, please be patient and a doc explaining it will be published soon.
			if r.rt.NPeersForCpl(cpl) == 0 {
				lastCpl := min(2*(c+1), len(refreshCpls)-1)
				for i := c + 1; i < lastCpl+1; i++ {
					if err := rfnc(uint(i)); err != nil {
						merr = multierror.Append(merr, err)
					}
				}
				return merr
			}
		}
	}

	select {
	case r.refreshDoneCh <- struct{}{}:
	case <-r.ctx.Done():
		return r.ctx.Err()
	}

	return merr
}

func min(a int, b int) int {
	if a <= b {
		return a
	}

	return b
}

func (r *RtRefreshManager) refreshCplIfEligible(cpl uint, lastRefreshedAt time.Time) error {
	if time.Since(lastRefreshedAt) <= r.refreshInterval {
		logger.Debugf("not running refresh for cpl %d as time since last refresh not above interval", cpl)
		return nil
	}

	return r.refreshCpl(cpl)
}

func (r *RtRefreshManager) refreshCpl(cpl uint) error {
	// gen a key for the query to refresh the cpl
	key, err := r.refreshKeyGenFnc(cpl)
	if err != nil {
		return fmt.Errorf("failed to generated query key for cpl=%d, err=%s", cpl, err)
	}

	logger.Infof("starting refreshing cpl %d with key %s (routing table size was %d)",
		cpl, loggableRawKeyString(key), r.rt.Size())

	if err := r.runRefreshDHTQuery(key); err != nil {
		return fmt.Errorf("failed to refresh cpl=%d, err=%s", cpl, err)
	}

	logger.Infof("finished refreshing cpl %d, routing table size is now %d", cpl, r.rt.Size())
	return nil
}

func (r *RtRefreshManager) queryForSelf() error {
	if err := r.runRefreshDHTQuery(string(r.dhtPeerId)); err != nil {
		return fmt.Errorf("failed to query for self, err=%s", err)
	}
	return nil
}

func (r *RtRefreshManager) runRefreshDHTQuery(key string) error {
	queryCtx, cancel := context.WithTimeout(r.ctx, r.refreshQueryTimeout)
	defer cancel()

	err := r.refreshQueryFnc(queryCtx, key)

	if err == nil || (err == context.DeadlineExceeded && queryCtx.Err() == context.DeadlineExceeded) {
		return nil
	}

	return err
}

type loggableRawKeyString string

func (lk loggableRawKeyString) String() string {
	k := string(lk)

	if len(k) == 0 {
		return k
	}

	encStr := base32.RawStdEncoding.EncodeToString([]byte(k))

	return encStr
}
