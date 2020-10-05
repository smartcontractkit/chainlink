package networking

import (
	"context"
	"sync"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	dhtrouter "github.com/smartcontractkit/chainlink/libocr/networking/dht-router"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

var (
	_ types.Bootstrapper = &bootstrapper{}
)

type bootstrapper struct {
	peer             *concretePeer
	peerAllowlist    map[p2ppeer.ID]struct{}
	bootstrappers    []p2ppeer.AddrInfo
	routing          dhtrouter.PeerDiscoveryRouter
	logger           types.Logger
	configDigest     types.ConfigDigest
	ctx              context.Context
	ctxCancel        context.CancelFunc
	state            bootstrapperState
	stateMu          *sync.Mutex
	failureThreshold int
}

type bootstrapperState int

const (
	bootstrapperUnstarted = iota
	bootstrapperStarted
	bootstrapperClosed
)

func newBootstrapper(logger types.Logger, configDigest types.ConfigDigest,
	peer *concretePeer, peerIDs []p2ppeer.ID, bootstrappers []p2ppeer.AddrInfo, F int,
) (*bootstrapper, error) {
	allowlist := make(map[p2ppeer.ID]struct{})
	for _, pid := range peerIDs {
		allowlist[pid] = struct{}{}
	}
	for _, b := range bootstrappers {
		allowlist[b.ID] = struct{}{}
	}

	ctx, cancel := context.WithCancel(context.Background())

	logger = loghelper.MakeLoggerWithContext(logger, types.LogFields{
		"id":           "OCREndpoint",
		"configDigest": configDigest.Hex(),
	})

	return &bootstrapper{
		peer,
		allowlist,
		bootstrappers,
		nil,
		logger,
		configDigest,
		ctx,
		cancel,
		bootstrapperUnstarted,
		new(sync.Mutex),
		F,
	}, nil
}

func (b *bootstrapper) Start() error {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()

	if b.state != bootstrapperUnstarted {
		panic("bootstrapper has already been started")
	}

	b.state = bootstrapperStarted

	if err := b.peer.register(b); err != nil {
		return err
	}

	if err := b.setupDHT(); err != nil {
		return errors.Wrap(err, "error setting up DHT")
	}

	b.logger.Info("Bootstrapper: Started listening", nil)

	return nil
}

func (b *bootstrapper) setupDHT() (err error) {
	config := dhtrouter.BuildConfig(
		b.bootstrappers,
		dhtPrefix,
		b.configDigest,
		b.logger,
		b.failureThreshold,
		true,
	)

	acl := dhtrouter.NewPermitListACL(b.logger)

	acl.Activate(config.ProtocolID(), b.allowlist()...)
	aclHost := dhtrouter.WrapACL(b.peer, acl, b.logger)

	b.routing, err = dhtrouter.NewDHTRouter(
		b.ctx,
		config,
		aclHost,
	)
	if err != nil {
		return errors.Wrap(err, "could not initialize DHTRouter")
	}

	b.routing.Start()

	return nil
}

func (b *bootstrapper) Close() error {
	b.stateMu.Lock()
	if b.state != bootstrapperStarted {
		b.stateMu.Unlock()
		panic("cannot close bootstrapper that is not started")
	}
	b.state = ocrEndpointClosed
	b.stateMu.Unlock()

	if err := b.routing.Close(); err != nil {
		return errors.Wrap(err, "could not close dht router")
	}

	b.ctxCancel()

	return errors.Wrap(b.peer.deregister(b), "could not unregister bootstrapper")
}

func (b *bootstrapper) isAllowed(id p2ppeer.ID) bool {
	_, ok := b.peerAllowlist[id]
	return ok
}

func (b *bootstrapper) allowlist() (allowlist []p2ppeer.ID) {
	for k := range b.peerAllowlist {
		allowlist = append(allowlist, k)
	}
	return
}

func (b *bootstrapper) getConfigDigest() types.ConfigDigest {
	return b.configDigest
}
