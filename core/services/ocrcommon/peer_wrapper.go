package ocrcommon

import (
	"context"
	"crypto"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	"github.com/smartcontractkit/libocr/networking/rageping"
	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	creproxy "github.com/smartcontractkit/chainlink-protos/cre/impl/proxy"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

type PeerWrapperOCRConfig interface {
	TraceLogging() bool
}

type (
	peerAdapterOCR1 struct {
		ocr1types.BinaryNetworkEndpointFactory
		ocr1types.BootstrapperFactory
	}

	peerAdapterOCR2 struct {
		ocr2types.BinaryNetworkEndpointFactory
		ocr2types.BootstrapperFactory
	}

	peerAdapterOCR3_1 struct {
		ocr2types.BinaryNetworkEndpoint2Factory
		ocr2types.BootstrapperFactory
	}

	// SingletonPeerWrapper manages all libocr peers for the application
	SingletonPeerWrapper struct {
		services.StateMachine
		keyStore keystore.Master
		p2pCfg   config.P2P
		ocrCfg   PeerWrapperOCRConfig
		ds       sqlutil.DataSource
		lggr     logger.Logger
		PeerID   p2pkey.PeerID

		// proxyAddr, when non-empty, makes the wrapper delegate all rage
		// (OCR + PeerGroup) to an out-of-process proxy at this gRPC address
		// instead of running a local libocr peer.
		proxyAddr    string
		proxyClosers []io.Closer

		// Used at shutdown to stop all of this peer's goroutines
		peerCloser io.Closer

		// OCR1 peer adapter
		Peer1 *peerAdapterOCR1

		// OCR2 peer adapter
		Peer2 *peerAdapterOCR2

		// OCR3_1 peer adapter
		Peer3_1 *peerAdapterOCR3_1

		// PeerGroupFactory can be used to create PeerGroup instances
		PeerGroupFactory ocrnetworking.PeerGroupFactory
	}
)

func ValidatePeerWrapperConfig(config config.P2P) error {
	if len(config.V2().ListenAddresses()) == 0 {
		return errors.New("no P2P.V2.ListenAddresses specified")
	}
	return nil
}

// NewSingletonPeerWrapper creates a new peer based on the p2p keys in the keystore
// It currently only supports one peerID/key
// It should be fairly easy to modify it to support multiple peerIDs/keys using e.g. a map
// proxyAddr, when non-empty, makes the wrapper delegate all rage networking to
// an out-of-process proxy at that gRPC address instead of running a local peer.
func NewSingletonPeerWrapper(keyStore keystore.Master, p2pCfg config.P2P, ocrCfg PeerWrapperOCRConfig, ds sqlutil.DataSource, proxyAddr string, lggr logger.Logger) *SingletonPeerWrapper {
	return &SingletonPeerWrapper{
		keyStore:  keyStore,
		p2pCfg:    p2pCfg,
		ocrCfg:    ocrCfg,
		ds:        ds,
		proxyAddr: proxyAddr,
		lggr:      lggr.Named("SingletonPeerWrapper"),
	}
}

func (p *SingletonPeerWrapper) IsStarted() bool { return p.Ready() == nil }

// Start starts SingletonPeerWrapper.
func (p *SingletonPeerWrapper) Start(context.Context) error {
	return p.StartOnce("SingletonPeerWrapper", func() error {
		if p.proxyAddr != "" {
			return p.startProxy()
		}

		peerConfig, err := p.peerConfig()
		if err != nil {
			return err
		}

		p.lggr.Debugw("Creating OCR/OCR2 Peer", "config", peerConfig)
		// Note: creates and starts the peer
		peer, err := ocrnetworking.NewPeer(peerConfig)
		if err != nil {
			return errors.Wrap(err, "error calling NewPeer")
		}
		p.Peer1 = &peerAdapterOCR1{
			peer.OCR1BinaryNetworkEndpointFactory(),
			peer.OCR1BootstrapperFactory(),
		}
		p.Peer2 = &peerAdapterOCR2{
			peer.OCR2BinaryNetworkEndpointFactory(),
			peer.OCR2BootstrapperFactory(),
		}
		p.Peer3_1 = &peerAdapterOCR3_1{
			peer.OCR3_1BinaryNetworkEndpointFactory(),
			peer.OCR2BootstrapperFactory(),
		}

		p.PeerGroupFactory = peer.PeerGroupFactory()

		p.peerCloser = peer
		return nil
	})
}

// startProxy delegates all rage networking to an out-of-process proxy: it does
// not create a local libocr peer, and instead exposes Peer2 (OCR) and
// PeerGroupFactory (DON-to-DON) backed by proxy clients that connect to the
// proxy's gRPC. Peer1/Peer3_1 are not proxied yet.
func (p *SingletonPeerWrapper) startProxy() error {
	if ks, err := p.keyStore.P2P().GetAll(); err == nil && len(ks) == 0 {
		return errors.New("No P2P keys found in keystore. Peer wrapper will not be fully initialized")
	}
	key, err := p.keyStore.P2P().GetOrFirst(p.p2pCfg.PeerID())
	if err != nil {
		return err
	}
	p.PeerID = key.PeerID()

	// NOTE: libocr expects the raw peer ID (base58, no core-specific "p2p_"
	// prefix) — it is compared against peer IDs in the OCR config.
	endpointFactory, err := creproxy.NewProxyEndpointFactory(p.PeerID.Raw(), p.proxyAddr)
	if err != nil {
		return errors.Wrap(err, "failed to create proxy OCR endpoint factory")
	}
	endpoint2Factory, err := creproxy.NewProxyEndpoint2Factory(p.PeerID.Raw(), p.proxyAddr)
	if err != nil {
		_ = endpointFactory.Close()
		return errors.Wrap(err, "failed to create proxy OCR3.1 endpoint factory")
	}
	pgClient, err := creproxy.NewProxyPeerGroupFactory(p.proxyAddr)
	if err != nil {
		_ = endpointFactory.Close()
		_ = endpoint2Factory.Close()
		return errors.Wrap(err, "failed to create proxy peer group factory")
	}

	// The proxy only serves the endpoint side; the bootstrapper factory is unused
	// on non-bootstrap nodes. Peer1 (OCR1) is not proxied.
	p.Peer2 = &peerAdapterOCR2{endpointFactory, nil}
	p.Peer3_1 = &peerAdapterOCR3_1{endpoint2Factory, nil}
	p.PeerGroupFactory = newProxyBackedPeerGroupFactory(pgClient)
	p.proxyClosers = []io.Closer{endpointFactory, endpoint2Factory, pgClient}

	p.lggr.Infow("SingletonPeerWrapper delegating rage to proxy", "proxyAddr", p.proxyAddr, "peerID", p.PeerID.String())
	return nil
}

func (p *SingletonPeerWrapper) peerConfig() (ocrnetworking.PeerConfig, error) {
	// Peer wrapper panics if no p2p keys are present.
	if ks, err := p.keyStore.P2P().GetAll(); err == nil && len(ks) == 0 {
		return ocrnetworking.PeerConfig{}, errors.Errorf("No P2P keys found in keystore. Peer wrapper will not be fully initialized")
	}
	key, err := p.keyStore.P2P().GetOrFirst(p.p2pCfg.PeerID())
	if err != nil {
		return ocrnetworking.PeerConfig{}, err
	}
	p.PeerID = key.PeerID()

	discovererDB := NewOCRDiscovererDatabase(p.ds, p.PeerID.Raw())

	peerKeyring, err := NewSignerPeerKeyring(key)
	if err != nil {
		return ocrnetworking.PeerConfig{}, err
	}
	config := p.p2pCfg
	peerConfig := ocrnetworking.PeerConfig{
		PeerKeyring: peerKeyring,
		Logger:      commonlogger.NewOCRWrapper(p.lggr, p.ocrCfg.TraceLogging(), func(string) {}),

		// V2 config
		V2ListenAddresses:    config.V2().ListenAddresses(),
		V2AnnounceAddresses:  config.V2().AnnounceAddresses(), // NewPeer will handle the fallback to listen addresses for us.
		V2DeltaReconcile:     config.V2().DeltaReconcile().Duration(),
		V2DeltaDial:          config.V2().DeltaDial().Duration(),
		V2DiscovererDatabase: discovererDB,

		V2EndpointConfig: ocrnetworking.EndpointConfigV2{
			IncomingMessageBufferSize: config.IncomingMessageBufferSize(),
			OutgoingMessageBufferSize: config.OutgoingMessageBufferSize(),
		},
		MetricsRegisterer:            prometheus.DefaultRegisterer,
		LatencyMetricsServiceConfigs: rageping.DefaultConfigs(),
	}
	if config.EnableExperimentalRageP2P() {
		peerConfig.EnableExperimentalRageP2P = ocrnetworking.DangerDangerEnableExperimentalRageP2P
	}

	return peerConfig, nil
}

// Close closes the peer and peerstore
func (p *SingletonPeerWrapper) Close() error {
	return p.StopOnce("SingletonPeerWrapper", func() (err error) {
		if p.peerCloser != nil {
			err = p.peerCloser.Close()
		}
		for _, c := range p.proxyClosers {
			if cerr := c.Close(); cerr != nil && err == nil {
				err = cerr
			}
		}
		return err
	})
}

func (p *SingletonPeerWrapper) Name() string {
	return p.lggr.Name()
}

func (p *SingletonPeerWrapper) HealthReport() map[string]error {
	return map[string]error{p.Name(): p.Healthy()}
}

func (p *SingletonPeerWrapper) P2PConfig() config.P2P {
	return p.p2pCfg
}

type signerPeerKeyring struct {
	signer        crypto.Signer
	peerPublicKey ragetypes.PeerPublicKey
}

func NewSignerPeerKeyring(signer crypto.Signer) (ragetypes.PeerKeyring, error) {
	peerPublicKey, err := ragetypes.PeerPublicKeyFromGenericPublicKey(signer.Public())
	if err != nil {
		return nil, err
	}
	return &signerPeerKeyring{signer, peerPublicKey}, nil
}

func (s *signerPeerKeyring) PublicKey() ragetypes.PeerPublicKey {
	return s.peerPublicKey
}

func (s *signerPeerKeyring) Sign(msg []byte) (signature []byte, err error) {
	return s.signer.Sign(rand.Reader, msg, crypto.Hash(0))
}
