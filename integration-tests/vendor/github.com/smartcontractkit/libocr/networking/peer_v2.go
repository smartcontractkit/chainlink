package networking

import (
	"crypto/ed25519"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/networking/ragedisco"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"go.uber.org/multierr"
)

// concretePeerV2 represents a ragep2p peer with one peer ID listening on one port
type concretePeerV2 struct {
	peerID         ragetypes.PeerID
	host           *ragep2p.Host
	discoverer     *ragedisco.Ragep2pDiscoverer
	logger         loghelper.LoggerWithContext
	endpointConfig EndpointConfigV2
}

func newPeerV2(c PeerConfig) (*concretePeerV2, error) {

	rawPriv, err := c.PrivKey.Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to get raw private key to use for v2: %w", err)
	}
	ed25519Priv := ed25519.PrivateKey(rawPriv)
	if err := ed25519SanityCheck(ed25519Priv); err != nil {
		return nil, fmt.Errorf("ed25519 sanity check failed: %w", err)
	}

	peerID, err := ragetypes.PeerIDFromPrivateKey(ed25519Priv)
	if err != nil {
		return nil, fmt.Errorf("error extracting v2 peer ID from private key: %w", err)
	}

	logger := loghelper.MakeRootLoggerWithContext(c.Logger).MakeChild(commontypes.LogFields{
		"id":     "PeerV2",
		"peerID": peerID.String(),
	})

	announceAddresses := c.V2AnnounceAddresses
	if len(c.V2AnnounceAddresses) == 0 {
		announceAddresses = c.V2ListenAddresses
	}
	discoverer := ragedisco.NewRagep2pDiscoverer(c.V2DeltaReconcile, announceAddresses, c.V2DiscovererDatabase)
	host, err := ragep2p.NewHost(
		ragep2p.HostConfig{c.V2DeltaDial},
		ed25519Priv,
		c.V2ListenAddresses,
		discoverer,
		c.Logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct ragep2p host: %w", err)
	}
	err = host.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start ragep2p host: %w", err)
	}

	logger.Info("PeerV2: ragep2p host booted", nil)

	return &concretePeerV2{
		peerID,
		host,
		discoverer,
		logger,
		c.V2EndpointConfig,
	}, nil
}

// An endpointRegistration is held by an endpoint which services a particular configDigest. The invariant is that only
// there can be at most a single active (ie. not closed) endpointRegistration for some configDigest, and thus only at
// most one endpoint can service a particular configDigest at any given point in time. The endpoint is responsible for
// calling Close on the registration.
type endpointRegistration struct {
	deregisterFunc func() error
	once           sync.Once
}

func newEndpointRegistration(deregisterFunc func() error) *endpointRegistration {
	return &endpointRegistration{deregisterFunc, sync.Once{}}
}

func (r *endpointRegistration) Close() (err error) {
	r.once.Do(func() {
		err = r.deregisterFunc()
	})
	return err
}

func (p2 *concretePeerV2) register(configDigest ocr2types.ConfigDigest, oracles []ragetypes.PeerID, bootstrappers []ragetypes.PeerInfo) (*endpointRegistration, error) {
	if err := p2.discoverer.AddGroup(configDigest, oracles, bootstrappers); err != nil {
		p2.logger.Warn("PeerV2: Failed to register endpoint", commontypes.LogFields{"configDigest": configDigest})
		return nil, err
	}

	return newEndpointRegistration(func() error {
		// Discoverer will not be closed until concretePeerV2.Close() is called.
		// By the time concretePeerV2.Close() is called all endpoints/bootstrappers should have already been closed.
		// Even if this weren't true, RemoveGroup() is a no-op if the discoverer is closed.
		return p2.discoverer.RemoveGroup(configDigest)
	}), nil
}

func (p2 *concretePeerV2) PeerID() string {
	return p2.peerID.String()
}

func (p2 *concretePeerV2) Close() error {
	return p2.host.Close()
}
func decodev2Bootstrappers(v2bootstrappers []commontypes.BootstrapperLocator) (infos []ragetypes.PeerInfo, err error) {
	for _, b := range v2bootstrappers {
		addrs := make([]ragetypes.Address, len(b.Addrs))
		for i, a := range b.Addrs {
			addrs[i] = ragetypes.Address(a)
		}
		var rageID ragetypes.PeerID
		err := rageID.UnmarshalText([]byte(b.PeerID))
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal v2 peer ID (%q) from BootstrapperLocator: %w", b.PeerID, err)
		}
		infos = append(infos, ragetypes.PeerInfo{
			rageID,
			addrs,
		})
	}
	return
}

func decodev2PeerIDs(pids []string) ([]ragetypes.PeerID, error) {
	peerIDs := make([]ragetypes.PeerID, len(pids))
	for i, pid := range pids {
		var rid ragetypes.PeerID
		err := rid.UnmarshalText([]byte(pid))
		if err != nil {
			return nil, fmt.Errorf("error decoding v2 peer ID (%q): %w", pid, err)
		}
		peerIDs[i] = rid
	}
	return peerIDs, nil
}

func (p2 *concretePeerV2) newEndpoint(
	configDigest ocr2types.ConfigDigest,
	v2peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
	limits BinaryNetworkEndpointLimits,
) (commontypes.BinaryNetworkEndpoint, error) {
	if f <= 0 {
		return nil, errors.New("can't set F to 0 or smaller")
	}

	if len(v2bootstrappers) < 1 {
		return nil, errors.New("requires at least one v2 bootstrapper")
	}

	decodedv2PeerIDs, err := decodev2PeerIDs(v2peerIDs)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 peer IDs: %w", err)
	}

	decodedv2Bootstrappers, err := decodev2Bootstrappers(v2bootstrappers)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 bootstrappers: %w", err)
	}

	registration, err := p2.register(configDigest, decodedv2PeerIDs, decodedv2Bootstrappers)
	if err != nil {
		return nil, err
	}

	endpoint, err := newOCREndpointV2(
		p2.logger,
		configDigest,
		p2,
		decodedv2PeerIDs,
		decodedv2Bootstrappers,
		EndpointConfigV2{
			p2.endpointConfig.IncomingMessageBufferSize,
			p2.endpointConfig.OutgoingMessageBufferSize,
		},
		f,
		limits,
		registration,
	)
	if err != nil {
		// Important: we close registration in case newOCREndpointV2 failed to prevent zombie registrations.
		return nil, multierr.Combine(err, registration.Close())
	}
	return endpoint, nil
}

func (p2 *concretePeerV2) newBootstrapper(
	configDigest ocr2types.ConfigDigest,
	v2peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
) (commontypes.Bootstrapper, error) {
	if f <= 0 {
		return nil, errors.New("can't set f to zero or smaller")
	}

	decodedv2PeerIDs, err := decodev2PeerIDs(v2peerIDs)
	if err != nil {
		return nil, err
	}

	decodedv2Bootstrappers, err := decodev2Bootstrappers(v2bootstrappers)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 bootstrappers: %w", err)
	}

	registration, err := p2.register(configDigest, decodedv2PeerIDs, decodedv2Bootstrappers)
	if err != nil {
		return nil, err
	}

	bootstrapper, err := newBootstrapperV2(p2.logger, configDigest, p2, decodedv2PeerIDs, decodedv2Bootstrappers, f, registration)
	if err != nil {
		// Important: we close registration in case newBootstrapperV2 failed to prevent zombie registrations.
		return nil, multierr.Combine(err, registration.Close())
	}
	return bootstrapper, nil
}
