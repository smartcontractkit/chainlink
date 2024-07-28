package ocrcommon

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/networking/rageping"

	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
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

	// SingletonPeerWrapper manages all libocr peers for the application
	SingletonPeerWrapper struct {
		services.StateMachine
		keyStore keystore.Master
		p2pCfg   config.P2P
		ocrCfg   PeerWrapperOCRConfig
		ds       sqlutil.DataSource
		lggr     logger.Logger
		PeerID   p2pkey.PeerID

		// Used at shutdown to stop all of this peer's goroutines
		peerCloser io.Closer

		// OCR1 peer adapter
		Peer1 *peerAdapterOCR1

		// OCR2 peer adapter
		Peer2 *peerAdapterOCR2
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
func NewSingletonPeerWrapper(keyStore keystore.Master, p2pCfg config.P2P, ocrCfg PeerWrapperOCRConfig, ds sqlutil.DataSource, lggr logger.Logger) *SingletonPeerWrapper {
	return &SingletonPeerWrapper{
		keyStore: keyStore,
		p2pCfg:   p2pCfg,
		ocrCfg:   ocrCfg,
		ds:       ds,
		lggr:     lggr.Named("SingletonPeerWrapper"),
	}
}

func (p *SingletonPeerWrapper) IsStarted() bool { return p.Ready() == nil }

// Start starts SingletonPeerWrapper.
func (p *SingletonPeerWrapper) Start(context.Context) error {
	return p.StartOnce("SingletonPeerWrapper", func() error {
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
		p.peerCloser = peer
		return nil
	})
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

	config := p.p2pCfg
	peerConfig := ocrnetworking.PeerConfig{
		PrivKey: key.PrivKey,
		Logger:  commonlogger.NewOCRWrapper(p.lggr, p.ocrCfg.TraceLogging(), func(string) {}),

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

	return peerConfig, nil
}

// Close closes the peer and peerstore
func (p *SingletonPeerWrapper) Close() error {
	return p.StopOnce("SingletonPeerWrapper", func() (err error) {
		if p.peerCloser != nil {
			err = p.peerCloser.Close()
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
