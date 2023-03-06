package networking

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/multierr"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	dhtrouter "github.com/smartcontractkit/libocr/networking/dht-router"
	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
)

var (
	_ commontypes.Bootstrapper = &bootstrapperV1{}
)

type bootstrapperV1 struct {
	peer            *concretePeerV1
	peerAllowlist   map[p2ppeer.ID]struct{}
	v1bootstrappers []p2ppeer.AddrInfo
	routing         dhtrouter.PeerDiscoveryRouter
	logger          loghelper.LoggerWithContext
	configDigest    ocr1types.ConfigDigest
	ctx             context.Context
	ctxCancel       context.CancelFunc
	registered      bool
	state           bootstrapperState

	stateMu              *sync.Mutex
	f                    int
	lowerBandwidthLimits func()
}

type bootstrapperState int

const (
	bootstrapperUnstarted = iota
	bootstrapperStarted
	bootstrapperClosed
)

// Bandwidth rate limiter parameters for the bootstrapperV1.
// bootstrappers are contacted to fetch the mapping between peer IDs and peer IPs.
// This bootstrapping is supposed to happen relatively rarely. Also, the full mapping is only a few KiB.
const (
	bootstrapperTokenBucketRefillRate = 20 * 1024 // 20 KiB/s
	bootstrapperTokenBucketSize       = 50 * 1024 // 50 KiB/s
)

func newBootstrapperV1(
	logger loghelper.LoggerWithContext,
	configDigest ocr1types.ConfigDigest,
	peer *concretePeerV1,
	v1peerIDs []p2ppeer.ID,
	v1bootstrappers []p2ppeer.AddrInfo,
	f int,
) (*bootstrapperV1, error) {
	lowerBandwidthLimits := increaseBandwidthLimits(peer.bandwidthLimiters, v1peerIDs, v1bootstrappers,
		bootstrapperTokenBucketRefillRate, bootstrapperTokenBucketSize, logger)

	allowlist := make(map[p2ppeer.ID]struct{})
	for _, pid := range v1peerIDs {
		allowlist[pid] = struct{}{}
	}
	for _, b := range v1bootstrappers {
		allowlist[b.ID] = struct{}{}
	}

	ctx, cancel := context.WithCancel(context.Background())

	logger = logger.MakeChild(commontypes.LogFields{
		"id":           "bootstrapperV1",
		"configDigest": configDigest.Hex(),
	})
	return &bootstrapperV1{
		peer,
		allowlist,
		v1bootstrappers,
		nil,
		logger,
		configDigest,
		ctx,
		cancel,
		false,
		bootstrapperUnstarted,
		new(sync.Mutex),
		f,
		lowerBandwidthLimits,
	}, nil
}

// Start the bootstrapperV1. Should only be called once. Even in case of error Close() _should_ be called afterwards for cleanup.
func (b *bootstrapperV1) Start() error {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()

	if b.state != bootstrapperUnstarted {
		return fmt.Errorf("cannot start bootstrapperV1 that is not unstarted, state was: %d", b.state)
	}

	b.state = bootstrapperStarted

	if err := b.peer.register(b); err != nil {
		return err
	}
	b.registered = true
	if err := b.setupDHT(); err != nil {
		return fmt.Errorf("error setting up DHT: %w", err)
	}

	b.logger.Info("BootstrapperV1: Started listening", nil)

	return nil
}

func (b *bootstrapperV1) setupDHT() (err error) {
	config := dhtrouter.BuildConfig(
		b.v1bootstrappers,
		dhtPrefix,
		b.configDigest,
		b.logger,
		b.peer.endpointConfig.BootstrapCheckInterval,
		b.f,
		true,
		b.peer.dhtAnnouncementCounterUserPrefix,
	)

	acl := dhtrouter.NewPermitListACL(b.logger)

	acl.Activate(config.ProtocolID(), b.allowlist()...)
	aclHost := dhtrouter.WrapACL(b.peer.host, acl, b.logger)

	routing, err := dhtrouter.NewDHTRouter(
		b.ctx,
		config,
		aclHost,
	)
	if err != nil {
		return fmt.Errorf("could not initialize DHTRouter: %w", err)
	}

	b.routing = routing

	// Async
	b.routing.Start()

	return nil
}

func (b *bootstrapperV1) Close() error {
	b.stateMu.Lock()
	if b.state != bootstrapperStarted {
		defer b.stateMu.Unlock()
		return fmt.Errorf("cannot close bootstrapperV1 that is not started, state was: %d", b.state)
	}
	b.state = bootstrapperClosed
	b.stateMu.Unlock()

	b.ctxCancel()

	var allErrors error
	b.logger.Debug("BootstrapperV1: lowering v1 bandwidth limits when closing the bootstrapperV1", nil)
	b.lowerBandwidthLimits()

	if b.routing != nil {
		if err := b.routing.Close(); err != nil {
			allErrors = multierr.Append(allErrors, fmt.Errorf("could not close dht router: %w", err))
		}
	}
	if b.registered {
		if err := b.peer.deregister(b); err != nil {
			allErrors = multierr.Append(allErrors, fmt.Errorf("could not unregister bootstrapperV1: %w", err))
		}
	}
	return allErrors
}

// Conform to allower interface
func (b *bootstrapperV1) isAllowed(id p2ppeer.ID) bool {
	_, ok := b.peerAllowlist[id]
	return ok
}

// Conform to allower interface
func (b *bootstrapperV1) allowlist() (allowlist []p2ppeer.ID) {
	for k := range b.peerAllowlist {
		allowlist = append(allowlist, k)
	}
	return
}

func (b *bootstrapperV1) getConfigDigest() ocr1types.ConfigDigest {
	return b.configDigest
}
