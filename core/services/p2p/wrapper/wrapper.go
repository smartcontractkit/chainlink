package wrapper

import (
	"context"
	"crypto"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/libocr/commontypes"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type peerWrapper struct {
	peer        types.Peer
	keystoreP2P keystore.P2P
	p2pConfig   config.P2P
	privateKey  crypto.Signer
	lggr        logger.Logger
	ds          sqlutil.DataSource
}

var _ types.PeerWrapper = &peerWrapper{}
var _ types.Signer = &peerWrapper{}

func NewExternalPeerWrapper(keystoreP2P keystore.P2P, p2pConfig config.P2P, ds sqlutil.DataSource, lggr logger.Logger) *peerWrapper {
	return &peerWrapper{
		keystoreP2P: keystoreP2P,
		p2pConfig:   p2pConfig,
		lggr:        lggr.Named("PeerWrapper"),
		ds:          ds,
	}
}

func (e *peerWrapper) GetPeer() types.Peer {
	return e.peer
}

// convert to "external" P2P PeerConfig, which is independent of OCR
// this has to be done in Start() because keystore is not unlocked at construction time
func (e *peerWrapper) convertPeerConfig() (p2p.PeerConfig, error) {
	key, err := e.keystoreP2P.GetOrFirst(e.p2pConfig.PeerID())
	if err != nil {
		return p2p.PeerConfig{}, err
	}

	discovererDB := ocrcommon.NewDON2DONDiscovererDatabase(e.ds, key.PeerID().Raw())
	bootstrappers, err := convertBootstrapperLocators(e.p2pConfig.V2().DefaultBootstrappers())
	if err != nil {
		return p2p.PeerConfig{}, err
	}

	peerConfig := p2p.PeerConfig{
		PrivateKey: key,

		ListenAddresses:   e.p2pConfig.V2().ListenAddresses(),
		AnnounceAddresses: e.p2pConfig.V2().AnnounceAddresses(),
		Bootstrappers:     bootstrappers,

		DeltaReconcile:     e.p2pConfig.V2().DeltaReconcile().Duration(),
		DeltaDial:          e.p2pConfig.V2().DeltaDial().Duration(),
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
	cfg, err := e.convertPeerConfig()
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
	return e.lggr.Name()
}

func (e *peerWrapper) Sign(msg []byte) ([]byte, error) {
	if e.privateKey == nil {
		return nil, errors.New("private key not set")
	}
	return e.privateKey.Sign(rand.Reader, msg, crypto.Hash(0))
}
