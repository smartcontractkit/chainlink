package vault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ capabilities.ExecutableCapability = (*Service)(nil)

type Service struct {
	clock        clockwork.Clock
	expiresAfter time.Duration
	handler      *requests.Handler[*Request, *Response]
	store        *requests.Store[*Request, *Response]
}

func (s *Service) Start(ctx context.Context) error {
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.NewCapabilityInfo(vault.CapabilityID, capabilities.CapabilityTypeAction, "Vault Service")
}

func (s *Service) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	// Left unimplemented as this method will never be called
	// for this capability
	return nil
}

func (s *Service) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	// Left unimplemented as this method will never be called
	// for this capability
	return nil
}

func (s *Service) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	if request.Payload == nil {
		return capabilities.CapabilityResponse{}, errors.New("capability does not support v1 requests")
	}

	if request.Method != vault.MethodGetSecrets {
		return capabilities.CapabilityResponse{}, errors.New("unsupported method: can only call GetSecrets via capability interface")
	}

	r := &vault.GetSecretsRequest{}
	err := request.Payload.UnmarshalTo(r)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("could not unmarshal payload to GetSecretsRequest: %w", err)
	}

	respCh := make(chan *Response, 1)
	s.handler.SendRequest(ctx, &Request{
		Payload:      r,
		ResponseChan: respCh,

		expiryTime: s.clock.Now().Add(s.expiresAfter),
		id:         "", // TODO
	})

	select {
	case <-ctx.Done():
		return capabilities.CapabilityResponse{}, ctx.Err()
	case resp := <-respCh:
		// TODO: handle response
		_ = resp
		return capabilities.CapabilityResponse{}, nil
	}
}

func NewService(
	lggr logger.Logger,
	clock clockwork.Clock,
	expiresAfter time.Duration,
) *Service {
	store := requests.NewStore[*Request, *Response]()
	return &Service{
		clock:        clock,
		expiresAfter: expiresAfter,
		store:        store,
		handler:      requests.NewHandler(lggr, store, clock, expiresAfter),
	}
}
