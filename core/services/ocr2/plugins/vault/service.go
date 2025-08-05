package vault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonboulle/clockwork"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

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
}

func (s *Service) Start(ctx context.Context) error {
	return s.handler.Start(ctx)
}

func (s *Service) Close() error {
	return s.handler.Close()
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

	// Validate the request: we only check that the request contains at least one secret request.
	// All other validation is done in the plugin and subject to consensus.
	if len(r.Requests) == 0 {
		return capabilities.CapabilityResponse{}, errors.New("no secret request specified in request")
	}

	// We need to generate sufficiently unique IDs accounting for two cases:
	// 1. called during the subscription phase, in which case the executionID will be blank
	// 2. called during execution, in which case it'll be present.
	// The reference ID is unique per phase, so we need to differentiate when generating
	// an ID.
	md := request.Metadata
	phaseOrExecution := md.WorkflowExecutionID
	if phaseOrExecution == "" {
		phaseOrExecution = "subscription"
	}
	id := fmt.Sprintf("%s::%s::%s", md.WorkflowID, phaseOrExecution, md.ReferenceID)

	resp, err := handleRequest(ctx, s, id, r)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	// Note: we can drop the signatures from the response above here
	// since only a valid report will be successfully decryptable by the workflow DON.
	resppb := &vault.GetSecretsResponse{}
	err = proto.Unmarshal(resp.Payload, resppb)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("could not unmarshal response to GetSecretsResponse: %w", err)
	}

	anyproto, err := anypb.New(resppb)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("could not marshal response to anypb: %w", err)
	}

	return capabilities.CapabilityResponse{
		Payload: anyproto,
	}, nil
}

func handleRequest(ctx context.Context, s *Service, id string, request proto.Message) (*Response, error) {
	respCh := make(chan *Response, 1)
	s.handler.SendRequest(ctx, &Request{
		Payload:      request,
		ResponseChan: respCh,

		expiryTime: s.clock.Now().Add(s.expiresAfter),
		id:         id,
	})

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp := <-respCh:
		if resp.Error != "" {
			return nil, fmt.Errorf("error processing request %s: %w", id, errors.New(resp.Error))
		}

		return resp, nil
	}
}

func (s *Service) CreateSecrets(ctx context.Context, request *vault.CreateSecretsRequest) (*Response, error) {
	return handleRequest(ctx, s, request.RequestId, request)
}

func NewService(
	lggr logger.Logger,
	store *requests.Store[*Request],
	clock clockwork.Clock,
	expiresAfter time.Duration,
) *Service {
	return &Service{
		clock:        clock,
		expiresAfter: expiresAfter,
		handler:      requests.NewHandler(lggr, store, clock, expiresAfter),
	}
}
