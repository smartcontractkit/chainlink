package networking

import (
	"fmt"
	"net"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	nettypes "github.com/smartcontractkit/libocr/networking/types"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/libocr/internal/configdigesthelper"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/pkg/errors"
)

// PeerConfig configures the peer. A peer can operate with the v1 or v2 or both networking stacks, depending on
// the NetworkingStack set. The options for each stack are clearly marked, those for v1 start with V1 and those for v2
// start with V2. Only the options for the desired stack(s) need to be set.
type PeerConfig struct {
	// NetworkingStack declares which network stack will be used: v1, v2 or both (prefer v2).
	NetworkingStack NetworkingStack
	PrivKey         p2pcrypto.PrivKey
	Logger          commontypes.Logger

	V1ListenIP   net.IP
	V1ListenPort uint16
	// V1AnnounceIP _should not_ be an unspecified IP (https://pkg.go.dev/net#IP.IsUnspecified), otherwise V1 peer
	// discovery will not work. V1ListenIP might be an unspecified IP, so be careful!
	// For auto-detection of the V1 announce IP and port based on V1ListenIP and V1ListenPort, you should leave both
	// V1AnnounceIP and V1AnnouncePort unspecified!
	V1AnnounceIP net.IP
	// V1AnnouncePort should be set if and only if V1AnnounceIP is also set.
	V1AnnouncePort uint16
	V1Peerstore    p2ppeerstore.Peerstore

	// This should be 0 most of times, but when needed (eg when counter is somehow rolled back)
	// users can bump this value to manually bump the counter.
	V1DHTAnnouncementCounterUserPrefix uint32

	// V2ListenAddresses contains the addresses the peer will listen to on the network in <ip>:<port> form as
	// accepted by net.Listen.
	V2ListenAddresses []string

	// V2AnnounceAddresses contains the addresses the peer will advertise on the network in <host>:<port> form as
	// accepted by net.Dial. The addresses should be reachable by peers of interest.
	// May be left unspecified, in which case the announce addresses are auto-detected based on V2ListenAddresses.
	V2AnnounceAddresses []string

	// Every V2DeltaReconcile a Reconcile message is sent to every peer.
	V2DeltaReconcile time.Duration

	// Dial attempts will be at least V2DeltaDial apart.
	V2DeltaDial time.Duration

	V2DiscovererDatabase nettypes.DiscovererDatabase

	V1EndpointConfig EndpointConfigV1
	V2EndpointConfig EndpointConfigV2
}

// concretePeer represents a libp2p and/or ragep2p peer
type concretePeer struct {
	v1              *concretePeerV1
	v2              *concretePeerV2
	logger          loghelper.LoggerWithContext
	networkingStack NetworkingStack
}

// NewPeer constructs a new peer, consisting of the v1 and/or v2 sub-peers
// depending on the networking stack requested in PeerConfig. Specifically:
// NetworkingStackV1: only the v1 peer is started
// NetworkingStackV2: only the v2 peer is started
// NetworkingStackV1V2: both v1 and v2 are started, and NewPeer will fail if
// either fails to start.
//
// Users are expected to create (using the OCR*Factory() methods) and close endpoints and bootstrappers before calling
// Close() on the peer itself.
func NewPeer(c PeerConfig) (*concretePeer, error) {
	var (
		v1  *concretePeerV1
		v2  *concretePeerV2
		err error
	)

	if !c.NetworkingStack.needsv1() && !c.NetworkingStack.needsv2() {
		return nil, errors.New("networking stack must be v1, v2, or v1v2")
	}

	logger := loghelper.MakeRootLoggerWithContext(c.Logger)
	if c.NetworkingStack.needsv1() {
		v1, err = newPeerV1(c)
		if err != nil {
			return nil, fmt.Errorf("failed to make v1 peer: %w", err)
		}
		logger = logger.MakeChild(commontypes.LogFields{"v1peerID": v1.PeerID()})
	}
	if c.NetworkingStack.needsv2() {
		v2, err = newPeerV2(c)
		if err != nil {
			return nil, fmt.Errorf("failed to make v2 peer: %w", err)
		}
		logger = logger.MakeChild(commontypes.LogFields{"v2peerID": v2.PeerID()})
	}
	return &concretePeer{v1, v2, logger, c.NetworkingStack}, nil
}

func (p *concretePeer) newEndpointV1(
	configDigest ocr2types.ConfigDigest,
	peerIDs []string,
	v1bootstrappers []string,
	f int,
	limits BinaryNetworkEndpointLimits,
) (commontypes.BinaryNetworkEndpoint, error) {
	v1ConfigDigest, err := configdigesthelper.OCR2ToOCR1(configDigest)
	if err != nil {
		return nil, err
	}

	return p.v1.newEndpoint(v1ConfigDigest, peerIDs, v1bootstrappers, f, limits)
}

// newEndpoint returns an appropriate OCR endpoint depending on the networking stack used
func (p *concretePeer) newEndpoint(
	networkingStack NetworkingStack,
	configDigest ocr2types.ConfigDigest,
	peerIDs []string,
	v1bootstrappers []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
	limits BinaryNetworkEndpointLimits,
) (commontypes.BinaryNetworkEndpoint, error) {
	if !networkingStack.subsetOf(p.networkingStack) {
		return nil, fmt.Errorf("unsupported networking stack %s for peer which has %s", networkingStack, p.networkingStack)
	}
	var (
		v1, v2       commontypes.BinaryNetworkEndpoint
		v1err, v2err error
	)
	if networkingStack.needsv1() {
		v1, v1err = p.newEndpointV1(
			configDigest,
			peerIDs,
			v1bootstrappers,
			f,
			limits,
		)
		if v1err != nil || networkingStack == NetworkingStackV1 {
			return v1, v1err
		}
	}
	if networkingStack.needsv2() {
		v2, v2err = p.v2.newEndpoint(configDigest, peerIDs, v2bootstrappers, f, limits)
		if networkingStack == NetworkingStackV2 {
			return v2, v2err
		}
	}

	if v2err != nil {
		p.logger.Critical("PeerV1V2: NewEndpoint: Failed for v2, operating only with v1", commontypes.LogFields{"error": v2err})
		return v1, nil
	}

	logger := p.logger.MakeChild(commontypes.LogFields{"id": "OCREndpointV1V2", "configDigest": configDigest})
	return newOCREndpointV1V2(logger, peerIDs, v1, v2)
}

func (p *concretePeer) newBootstrapperV1(
	configDigest ocr2types.ConfigDigest,
	peerIDs []string,
	v1bootstrappers []string,
	f int,
) (commontypes.Bootstrapper, error) {
	v1ConfigDigest, err := configdigesthelper.OCR2ToOCR1(configDigest)
	if err != nil {
		return nil, err
	}

	return p.v1.newBootstrapper(v1ConfigDigest, peerIDs, v1bootstrappers, f)
}

func (p *concretePeer) newBootstrapper(
	networkingStack NetworkingStack,
	configDigest ocr2types.ConfigDigest,
	peerIDs []string,
	v1bootstrappers []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
) (commontypes.Bootstrapper, error) {
	if !networkingStack.subsetOf(p.networkingStack) {
		return nil, fmt.Errorf("unsupported networking stack %s for peer which has %s", networkingStack, p.networkingStack)
	}
	var (
		v1, v2       commontypes.Bootstrapper
		v1err, v2err error
	)
	if networkingStack.needsv1() {
		v1, v1err = p.newBootstrapperV1(configDigest, peerIDs, v1bootstrappers, f)
		if v1err != nil || networkingStack == NetworkingStackV1 {
			return v1, v1err
		}
	}
	if networkingStack.needsv2() {
		v2, v2err = p.v2.newBootstrapper(configDigest, peerIDs, v2bootstrappers, f)
		if networkingStack == NetworkingStackV2 {
			return v2, v2err
		}
	}

	if v2err != nil {
		p.logger.Critical("PeerV1V2: NewBootstrapper: Failed for v2, operating only with v1", commontypes.LogFields{"error": v2err})
		return v1, nil
	}

	logger := p.logger.MakeChild(commontypes.LogFields{"id": "BootstrapperV1V2", "configDigest": configDigest})
	return newBootstrapperV1V2(logger, v1, v2)
}

func (p *concretePeer) PeerID() string {
	if p.networkingStack.needsv1() {
		return p.v1.PeerID()
	}
	return p.v2.PeerID()
}

func (p *concretePeer) Close() error {
	var allErrors error
	if p.networkingStack.needsv1() {
		allErrors = multierr.Append(allErrors, p.v1.Close())
	}
	if p.networkingStack.needsv2() {
		allErrors = multierr.Append(allErrors, p.v2.Close())
	}
	return allErrors
}

func (p *concretePeer) OCR1BinaryNetworkEndpointFactory() *ocr1BinaryNetworkEndpointFactory {
	return &ocr1BinaryNetworkEndpointFactory{p}
}

func (p *concretePeer) OCR2BinaryNetworkEndpointFactory() *ocr2BinaryNetworkEndpointFactory {
	return &ocr2BinaryNetworkEndpointFactory{p}
}

func (p *concretePeer) OCR1BootstrapperFactory() *ocr1BootstrapperFactory {
	return &ocr1BootstrapperFactory{p}
}

func (p *concretePeer) OCR2BootstrapperFactory() *ocr2BootstrapperFactory {
	return &ocr2BootstrapperFactory{p}
}
