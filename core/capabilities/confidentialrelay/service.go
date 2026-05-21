package confidentialrelay

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

// Service is a thin lifecycle wrapper around the confidential relay handler.
// The relay handler needs the gateway connector, which isn't available until
// the ServiceWrapper starts. This wrapper defers handler creation to Start().
type Service struct {
	services.Service
	eng *services.Engine

	wrapper       *gatewayconnector.ServiceWrapper
	capRegistry   core.CapabilitiesRegistry
	p2pKeystore   keystore.P2P
	peerID        p2pkey.PeerID
	lggr          logger.Logger
	limitsFactory limits.Factory

	handler *Handler
}

func NewService(
	wrapper *gatewayconnector.ServiceWrapper,
	capRegistry core.CapabilitiesRegistry,
	p2pKeystore keystore.P2P,
	peerID p2pkey.PeerID,
	lggr logger.Logger,
	limitsFactory limits.Factory,
) *Service {
	s := &Service{
		wrapper:       wrapper,
		capRegistry:   capRegistry,
		p2pKeystore:   p2pKeystore,
		peerID:        peerID,
		lggr:          lggr,
		limitsFactory: limitsFactory,
	}
	s.Service, s.eng = services.Config{
		Name:  "ConfidentialRelayService",
		Start: s.start,
		Close: s.close,
	}.NewServiceEngine(lggr)
	return s
}

func (s *Service) start(ctx context.Context) error {
	conn := s.wrapper.GetGatewayConnector()
	if conn == nil {
		return errors.New("gateway connector not available")
	}
	key, err := s.p2pKeystore.GetOrFirst(s.peerID)
	if err != nil {
		return fmt.Errorf("failed to get p2p key for confidential relay signing: %w", err)
	}
	h, err := NewHandler(s.capRegistry, conn, newRelayResponseSigner(key), s.lggr, s.limitsFactory)
	if err != nil {
		return err
	}
	s.handler = h
	return h.Start(ctx)
}

func (s *Service) close() error {
	if s.handler != nil {
		return s.handler.Close()
	}
	return nil
}

type relayResponseSigner interface {
	PublicKey() []byte
	Sign(payload []byte) ([]byte, error)
}

type relayResponseSignerImpl struct {
	key p2pkey.KeyV2
}

func newRelayResponseSigner(key p2pkey.KeyV2) relayResponseSigner {
	return relayResponseSignerImpl{key: key}
}

func (s relayResponseSignerImpl) PublicKey() []byte {
	pub := s.key.Public().(ed25519.PublicKey)
	return append([]byte(nil), pub...)
}

func (s relayResponseSignerImpl) Sign(payload []byte) ([]byte, error) {
	return s.key.Sign(rand.Reader, payload, crypto.Hash(0))
}
