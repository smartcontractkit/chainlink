package wrapper

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/libocr/commontypes"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type peerWrapper struct {
	peer        types.Peer
	keystoreP2P keystore.P2P
	p2pConfig   config.P2P
	privateKey  ed25519.PrivateKey
	lggr        logger.Logger
}

var _ types.PeerWrapper = &peerWrapper{}
var _ types.Signer = &peerWrapper{}

func NewExternalPeerWrapper(keystoreP2P keystore.P2P, p2pConfig config.P2P, lggr logger.Logger) *peerWrapper {
	return &peerWrapper{
		keystoreP2P: keystoreP2P,
		p2pConfig:   p2pConfig,
		lggr:        lggr,
	}
}

func (e *peerWrapper) GetPeer() types.Peer {
	return e.peer
}

// convert to "external" P2P PeerConfig, which is independent of OCR
// this has to be done in Start() because keystore is not unlocked at construction time
func convertPeerConfig(keystoreP2P keystore.P2P, p2pConfig config.P2P) (p2p.PeerConfig, error) {
	key, err := keystoreP2P.GetOrFirst(p2pConfig.PeerID())
	if err != nil {
		return p2p.PeerConfig{}, err
	}

	// TODO(KS-106): use real DB
	discovererDB := p2p.NewInMemoryDiscovererDatabase()
	bootstrappers, err := convertBootstrapperLocators(p2pConfig.V2().DefaultBootstrappers())
	if err != nil {
		return p2p.PeerConfig{}, err
	}

	peerConfig := p2p.PeerConfig{
		PrivateKey: key.PrivKey,

		ListenAddresses:   p2pConfig.V2().ListenAddresses(),
		AnnounceAddresses: p2pConfig.V2().AnnounceAddresses(),
		Bootstrappers:     bootstrappers,

		DeltaReconcile:     p2pConfig.V2().DeltaReconcile().Duration(),
		DeltaDial:          p2pConfig.V2().DeltaDial().Duration(),
		DiscovererDatabase: discovererDB,

		// NOTE: this is equivalent to prometheus.DefaultRegisterer, but we need to use a separate
		// object to avoid conflicts with the OCR registerer
		MetricsRegisterer: prometheus.NewRegistry(),
	}

	return peerConfig, nil
}

func convertBootstrapperLocators(bootstrappers []commontypes.BootstrapperLocator) ([]ragetypes.PeerInfo, error) {
	infos := []ragetypes.PeerInfo{}
	for _, b := range bootstrappers {
		addrs := make([]ragetypes.Address, len(b.Addrs))
		for i, a := range b.Addrs {
			addrs[i] = ragetypes.Address(a)
		}
		var rageID types.PeerID
		err := rageID.UnmarshalText([]byte(b.PeerID))
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal v2 peer ID (%q) from BootstrapperLocator: %w", b.PeerID, err)
		}
		infos = append(infos, ragetypes.PeerInfo{
			ID:    rageID,
			Addrs: addrs,
		})
	}
	return infos, nil
}

func (e *peerWrapper) Start(ctx context.Context) error {
	cfg, err := convertPeerConfig(e.keystoreP2P, e.p2pConfig)
	if err != nil {
		return err
	}
	e.privateKey = cfg.PrivateKey
	e.lggr.Info("Starting external P2P peer")
	peer, err := p2p.NewPeer(cfg, e.lggr)
	if err != nil {
		return err
	}
	e.peer = peer
	return e.peer.Start(ctx)
}

func (e *peerWrapper) Close() error {
	return e.peer.Close()
}

func (e *peerWrapper) Ready() error {
	return nil
}

func (e *peerWrapper) HealthReport() map[string]error {
	return nil
}

func (e *peerWrapper) Name() string {
	return "PeerWrapper"
}

func (e *peerWrapper) Sign(msg []byte) ([]byte, error) {
	if e.privateKey == nil {
		return nil, fmt.Errorf("private key not set")
	}
	return ed25519.Sign(e.privateKey, msg), nil
}
