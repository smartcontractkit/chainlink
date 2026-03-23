package confidentialrelay

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
)

// Service is a thin lifecycle wrapper around the confidential relay handler.
// The relay handler needs the gateway connector, which isn't available until
// the ServiceWrapper starts. This wrapper defers handler creation to Start().
type Service struct {
	services.Service
	eng *services.Engine

	wrapper     *gatewayconnector.ServiceWrapper
	capRegistry core.CapabilitiesRegistry
	trustedPCRs []byte
	caRootsPEM  string
	lggr        logger.Logger

	handler *Handler
}

func NewService(
	wrapper *gatewayconnector.ServiceWrapper,
	capRegistry core.CapabilitiesRegistry,
	trustedPCRs []byte,
	caRootsPEM string,
	lggr logger.Logger,
) *Service {
	s := &Service{
		wrapper:     wrapper,
		capRegistry: capRegistry,
		trustedPCRs: trustedPCRs,
		caRootsPEM:  caRootsPEM,
		lggr:        lggr,
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
	h, err := NewHandler(s.capRegistry, conn, s.trustedPCRs, s.lggr, s.caRootsPEM)
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
