package sharedpeerabstraction

import (
	"context"
	"crypto"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// SharedPeerAbstraction provides a unified interface for both OCR and non-OCR peer functionality
type SharedPeerAbstraction interface {
	services.Service

	// GetPeer returns the native P2P peer for direct communication
	GetPeer() p2ptypes.Peer

	// GetPeerGroupFactory returns the OCR PeerGroupFactory if available
	GetPeerGroupFactory() ocrnetworking.PeerGroupFactory

	// HasOCRCapability indicates if this abstraction supports OCR functionality
	HasOCRCapability() bool

	// GetSingletonWrapper returns the underlying SingletonPeerWrapper if available
	GetSingletonWrapper() *ocrcommon.SingletonPeerWrapper

	// Sign provides cryptographic signing capability
	Sign(data []byte) ([]byte, error)
}

// sharedPeerAbstraction implements SharedPeerAbstraction
type sharedPeerAbstraction struct {
	services.StateMachine

	// Configuration
	keyStore    keystore.Master
	keystoreP2P keystore.P2P
	p2pCfg      config.P2P
	ocrCfg      ocrcommon.PeerWrapperOCRConfig
	ds          sqlutil.DataSource
	lggr        logger.Logger

	// State
	mu               sync.RWMutex
	singletonWrapper *ocrcommon.SingletonPeerWrapper
	externalPeer     p2ptypes.Peer
	privateKey       crypto.Signer

	// Mode flags
	enableOCR     bool
	forceExternal bool
}

// SharedPeerConfig holds configuration for the shared peer abstraction
type SharedPeerConfig struct {
	KeyStore    keystore.Master
	KeystoreP2P keystore.P2P
	P2PConfig   config.P2P
	OCRConfig   ocrcommon.PeerWrapperOCRConfig
	DataSource  sqlutil.DataSource
	Logger      logger.Logger

	// EnableOCR determines if OCR functionality should be enabled
	EnableOCR bool

	// ForceExternal forces the use of external peer even when OCR is available
	ForceExternal bool
}

var _ SharedPeerAbstraction = &sharedPeerAbstraction{}

// NewSharedPeerAbstraction creates a new shared peer abstraction
func NewSharedPeerAbstraction(cfg SharedPeerConfig) *sharedPeerAbstraction {
	return &sharedPeerAbstraction{
		keyStore:      cfg.KeyStore,
		keystoreP2P:   cfg.KeystoreP2P,
		p2pCfg:        cfg.P2PConfig,
		ocrCfg:        cfg.OCRConfig,
		ds:            cfg.DataSource,
		lggr:          cfg.Logger.Named("SharedPeerAbstraction"),
		enableOCR:     cfg.EnableOCR,
		forceExternal: cfg.ForceExternal,
	}
}

func (s *sharedPeerAbstraction) Start(ctx context.Context) error {
	return s.StartOnce("SharedPeerAbstraction", func() error {
		s.mu.Lock()
		defer s.mu.Unlock()

		// Determine which peer implementation to use
		if s.enableOCR && !s.forceExternal {
			return s.startWithOCR(ctx)
		}
		return s.startWithExternalPeer(ctx)
	})
}

func (s *sharedPeerAbstraction) startWithOCR(ctx context.Context) error {
	s.lggr.Info("Starting shared peer abstraction with OCR support")

	// Validate OCR configuration
	if err := ocrcommon.ValidatePeerWrapperConfig(s.p2pCfg); err != nil {
		return fmt.Errorf("invalid P2P configuration for OCR: %w", err)
	}

	// Create and start SingletonPeerWrapper
	s.singletonWrapper = ocrcommon.NewSingletonPeerWrapper(
		s.keyStore,
		s.p2pCfg,
		s.ocrCfg,
		s.ds,
		s.lggr,
	)

	if err := s.singletonWrapper.Start(ctx); err != nil {
		return fmt.Errorf("failed to start singleton peer wrapper: %w", err)
	}

	// Extract private key for signing
	key, err := s.keyStore.P2P().GetOrFirst(s.p2pCfg.PeerID())
	if err != nil {
		return fmt.Errorf("failed to get P2P key: %w", err)
	}
	s.privateKey = key

	s.lggr.Info("Successfully started shared peer abstraction with OCR support")
	return nil
}

func (s *sharedPeerAbstraction) startWithExternalPeer(ctx context.Context) error {
	s.lggr.Info("Starting shared peer abstraction with external peer")

	// Get P2P key
	key, err := s.keystoreP2P.GetOrFirst(s.p2pCfg.PeerID())
	if err != nil {
		return fmt.Errorf("failed to get P2P key: %w", err)
	}
	s.privateKey = key

	// Create peer configuration
	peerConfig, err := s.createExternalPeerConfig(key)
	if err != nil {
		return fmt.Errorf("failed to create peer config: %w", err)
	}

	// Create and start external peer
	peer, err := p2p.NewPeer(peerConfig, s.lggr)
	if err != nil {
		return fmt.Errorf("failed to create external peer: %w", err)
	}

	if err := peer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start external peer: %w", err)
	}

	s.externalPeer = peer
	s.lggr.Info("Successfully started shared peer abstraction with external peer")
	return nil
}

func (s *sharedPeerAbstraction) createExternalPeerConfig(key crypto.Signer) (p2p.PeerConfig, error) {
	// Create discoverer database
	discovererDB := ocrcommon.NewDON2DONDiscovererDatabase(s.ds, s.p2pCfg.PeerID().Raw())

	// Convert bootstrappers
	bootstrappers, err := s.convertBootstrappers()
	if err != nil {
		return p2p.PeerConfig{}, fmt.Errorf("failed to convert bootstrappers: %w", err)
	}

	return p2p.PeerConfig{
		PrivateKey:         key,
		ListenAddresses:    s.p2pCfg.V2().ListenAddresses(),
		AnnounceAddresses:  s.p2pCfg.V2().AnnounceAddresses(),
		Bootstrappers:      bootstrappers,
		DeltaReconcile:     s.p2pCfg.V2().DeltaReconcile().Duration(),
		DeltaDial:          s.p2pCfg.V2().DeltaDial().Duration(),
		DiscovererDatabase: discovererDB,
		MetricsRegisterer:  prometheus.NewRegistry(), // Use separate registerer to avoid conflicts
	}, nil
}

func (s *sharedPeerAbstraction) convertBootstrappers() ([]ragetypes.PeerInfo, error) {
	bootstrappers := s.p2pCfg.V2().DefaultBootstrappers()
	infos := make([]ragetypes.PeerInfo, 0, len(bootstrappers))

	for _, b := range bootstrappers {
		addrs := make([]ragetypes.Address, len(b.Addrs))
		for i, a := range b.Addrs {
			addrs[i] = ragetypes.Address(a)
		}

		var rageID p2ptypes.PeerID
		if err := rageID.UnmarshalText([]byte(b.PeerID)); err != nil {
			return nil, fmt.Errorf("failed to unmarshal peer ID (%q): %w", b.PeerID, err)
		}

		infos = append(infos, ragetypes.PeerInfo{
			ID:    rageID,
			Addrs: addrs,
		})
	}

	return infos, nil
}

func (s *sharedPeerAbstraction) Close() error {
	return s.StopOnce("SharedPeerAbstraction", func() error {
		s.mu.Lock()
		defer s.mu.Unlock()

		var err error

		if s.singletonWrapper != nil {
			if closeErr := s.singletonWrapper.Close(); closeErr != nil {
				err = errors.Wrap(closeErr, "failed to close singleton wrapper")
			}
			s.singletonWrapper = nil
		}

		if s.externalPeer != nil {
			if closeErr := s.externalPeer.Close(); closeErr != nil {
				if err != nil {
					err = errors.Wrap(err, closeErr.Error())
				} else {
					err = errors.Wrap(closeErr, "failed to close external peer")
				}
			}
			s.externalPeer = nil
		}

		s.privateKey = nil
		s.lggr.Info("Shared peer abstraction closed")
		return err
	})
}

func (s *sharedPeerAbstraction) GetPeer() p2ptypes.Peer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// If we have an external peer, return it directly
	if s.externalPeer != nil {
		return s.externalPeer
	}

	// For OCR-based peers, we don't have a direct p2ptypes.Peer
	// This is a limitation of the current architecture
	return nil
}

func (s *sharedPeerAbstraction) GetPeerGroupFactory() ocrnetworking.PeerGroupFactory {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.singletonWrapper != nil {
		return s.singletonWrapper.PeerGroupFactory
	}

	return nil
}

func (s *sharedPeerAbstraction) HasOCRCapability() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.singletonWrapper != nil
}

func (s *sharedPeerAbstraction) GetSingletonWrapper() *ocrcommon.SingletonPeerWrapper {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.singletonWrapper
}

func (s *sharedPeerAbstraction) Sign(data []byte) ([]byte, error) {
	s.mu.RLock()
	privateKey := s.privateKey
	s.mu.RUnlock()

	if privateKey == nil {
		return nil, errors.New("private key not available")
	}

	return privateKey.Sign(nil, data, crypto.Hash(0))
}

func (s *sharedPeerAbstraction) Ready() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.singletonWrapper != nil {
		return s.singletonWrapper.Ready()
	}

	if s.externalPeer != nil {
		return s.externalPeer.Ready()
	}

	return errors.New("no peer implementation available")
}

func (s *sharedPeerAbstraction) HealthReport() map[string]error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.singletonWrapper != nil {
		return s.singletonWrapper.HealthReport()
	}

	if s.externalPeer != nil {
		return s.externalPeer.HealthReport()
	}

	return map[string]error{s.Name(): errors.New("no peer implementation available")}
}

func (s *sharedPeerAbstraction) Name() string {
	return s.lggr.Name()
}
