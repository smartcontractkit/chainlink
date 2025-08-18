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

const maxBatchSize = 10

type Service struct {
	lggr         logger.Logger
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

	resp, err := s.handleRequest(ctx, id, r)
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

func (s *Service) handleRequest(ctx context.Context, id string, request proto.Message) (*Response, error) {
	respCh := make(chan *Response, 1)
	s.handler.SendRequest(ctx, &Request{
		Payload:      request,
		ResponseChan: respCh,

		expiryTime: s.clock.Now().Add(s.expiresAfter),
		id:         id,
	})
	s.lggr.Debugw("sent request to handler", "requestId", id)

	select {
	case <-ctx.Done():
		s.lggr.Debugw("request timed out", "requestId", id, "error", ctx.Err())
		return nil, ctx.Err()
	case resp := <-respCh:
		s.lggr.Debugw("received response for request", "requestId", id, "error", resp.Error)
		if resp.Error != "" {
			return nil, fmt.Errorf("error processing request %s: %w", id, errors.New(resp.Error))
		}

		return resp, nil
	}
}

func (s *Service) CreateSecrets(ctx context.Context, request *vault.CreateSecretsRequest) (*Response, error) {
	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Service) UpdateSecrets(ctx context.Context, request *vault.UpdateSecretsRequest) (*Response, error) {
	if request.RequestId == "" {
		return nil, errors.New("request ID must not be empty")
	}

	if len(request.EncryptedSecrets) >= maxBatchSize {
		return nil, fmt.Errorf("request batch size exceeds maximum of %d", maxBatchSize)
	}

	for _, req := range request.EncryptedSecrets {
		if req.Id == nil {
			return nil, errors.New("secret ID must not be nil")
		}

		if req.Id.Key == "" || req.Id.Owner == "" {
			return nil, fmt.Errorf("secret ID must have both key and owner set: %v", req.Id)
		}
	}

	// TODO: secrets should be encrypted with the correct key
	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Service) ListSecretIdentifiers(ctx context.Context, request *vault.ListSecretIdentifiersRequest) (*Response, error) {
	if request.RequestId == "" {
		return nil, errors.New("request ID must not be empty")
	}

	if request.Owner == "" {
		return nil, errors.New("owner must not be empty")
	}

	return s.handleRequest(ctx, request.RequestId, request)
}

func NewService(
	lggr logger.Logger,
	clock clockwork.Clock,
	expiresAfter time.Duration,
	handler *requests.Handler[*Request, *Response],
) *Service {
	return &Service{
		lggr:         lggr.Named("VaultService"),
		clock:        clock,
		expiresAfter: expiresAfter,
		handler:      handler,
	}
}
