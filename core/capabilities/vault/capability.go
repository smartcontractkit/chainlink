package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jonboulle/clockwork"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	vaultapi "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
	syncerV2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
)

var _ capabilities.ExecutableCapability = (*Capability)(nil)

type Capability struct {
	lggr              logger.Logger
	clock             clockwork.Clock
	expiresAfter      time.Duration
	handler           *requests.Handler[*Request, *Response]
	requestAuthorizer vaultapi.RequestAuthorizer
}

func (s *Capability) Start(ctx context.Context) error {
	return s.handler.Start(ctx)
}

func (s *Capability) Close() error {
	return s.handler.Close()
}

func (s *Capability) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.NewCapabilityInfo(vaultcommon.CapabilityID, capabilities.CapabilityTypeAction, "Vault Capability")
}

func (s *Capability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	// Left unimplemented as this method will never be called
	// for this capability
	return nil
}

func (s *Capability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	// Left unimplemented as this method will never be called
	// for this capability
	return nil
}

func (s *Capability) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	if request.Payload == nil {
		return capabilities.CapabilityResponse{}, errors.New("capability does not support v1 requests")
	}

	if request.Method != vaultapi.MethodSecretsGet {
		return capabilities.CapabilityResponse{}, errors.New("unsupported method: can only call GetSecrets via capability interface")
	}

	r := &vaultcommon.GetSecretsRequest{}
	err := request.Payload.UnmarshalTo(r)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("could not unmarshal payload to GetSecretsRequest: %w", err)
	}

	// Validate the request: we only check that the request contains at least one secret request.
	// All other validations are done in the plugin and subject to consensus.
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
	resppb := &vaultcommon.GetSecretsResponse{}
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

func (s *Capability) CreateSecrets(ctx context.Context, request *vaultcommon.CreateSecretsRequest) (*Response, error) {
	s.lggr.Infof("Received Request: %s", request.String())
	err := s.validateRequest(request.RequestId, request.EncryptedSecrets, nil)
	if err != nil {
		s.lggr.Infof("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}
	authorized, owner, err := s.isAuthorizedRequest(ctx, request)
	if authorized == false || err != nil {
		s.lggr.Infof("Request [%s] not authorized for owner: %s", request.String(), owner)
		return nil, errors.New("request not authorized: " + err.Error())
	}
	if !strings.HasPrefix(request.RequestId, owner) {
		// Gateway should ensure it prefixes request ids with the owner, to ensure request uniqueness
		s.lggr.Infof("Request ID: [%s] must start with owner address: [%s]", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " must start with owner address: " + owner)
	}
	for _, req := range request.EncryptedSecrets {
		// Right owner for secrets can only be set here, after authorization
		// This ensures that users cannot access secrets belonging to other owners
		req.Id.Owner = owner
	}
	s.lggr.Infof("Processing authorized and normalized request [%s]", request.String())
	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Capability) UpdateSecrets(ctx context.Context, request *vaultcommon.UpdateSecretsRequest) (*Response, error) {
	s.lggr.Infof("Received Request: %s", request.String())
	err := s.validateRequest(request.RequestId, request.EncryptedSecrets, nil)
	if err != nil {
		s.lggr.Infof("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}
	authorized, owner, err := s.isAuthorizedRequest(ctx, request)
	if authorized == false || err != nil {
		s.lggr.Infof("Request [%s] not authorized for owner: %s", request.String(), owner)
		return nil, errors.New("request not authorized: " + err.Error())
	}
	if !strings.HasPrefix(request.RequestId, owner) {
		// Gateway should ensure it prefixes request ids with the owner, to ensure request uniqueness
		s.lggr.Infof("Request ID: [%s] must start with owner address: [%s]", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " must start with owner address: " + owner)
	}
	for _, req := range request.EncryptedSecrets {
		// Right owner for secrets can only be set here, after authorization
		// This ensures that users cannot access secrets belonging to other owners
		req.Id.Owner = owner
	}
	s.lggr.Infof("Processing authorized and normalized request [%s]", request.String())
	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Capability) DeleteSecrets(ctx context.Context, request *vaultcommon.DeleteSecretsRequest) (*Response, error) {
	s.lggr.Infof("Received Request: %s", request.String())
	err := s.validateRequest(request.RequestId, nil, request.Ids)
	if err != nil {
		s.lggr.Infof("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}

	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Capability) GetSecrets(ctx context.Context, requestID string, request *vaultcommon.GetSecretsRequest) (*Response, error) {
	s.lggr.Infof("Received Request: %s", request.String())
	if len(request.Requests) == 0 {
		return nil, errors.New("no GetSecret request specified in request")
	}
	if len(request.Requests) >= vaultapi.MaxBatchSize {
		return nil, fmt.Errorf("request batch size exceeds maximum of %d", vaultapi.MaxBatchSize)
	}

	// No auth needed, as this method is not exposed externally
	return s.handleRequest(ctx, requestID, request)
}

func (s *Capability) handleRequest(ctx context.Context, requestID string, request proto.Message) (*Response, error) {
	respCh := make(chan *Response, 1)
	s.handler.SendRequest(ctx, &Request{
		Payload:      request,
		ResponseChan: respCh,

		ExpiryTimeVal: s.clock.Now().Add(s.expiresAfter),
		IDVal:         requestID,
	})
	s.lggr.Debugw("sent request to OCR handler", "requestID", requestID)
	select {
	case <-ctx.Done():
		s.lggr.Debugw("request timed out", "requestID", requestID, "error", ctx.Err())
		return nil, ctx.Err()
	case resp := <-respCh:
		s.lggr.Debugw("received response for request", "requestID", requestID, "error", resp.Error)
		if resp.Error != "" {
			return nil, fmt.Errorf("error processing request %s: %w", requestID, errors.New(resp.Error))
		}

		return resp, nil
	}
}

func (s *Capability) isAuthorizedRequest(ctx context.Context, request any) (bool, string, error) {
	var params json.RawMessage
	params, err := json.Marshal(request)
	if err != nil {
		return false, "", fmt.Errorf("could not marshal CreateSecretsRequest: %w", err)
	}
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  vaultapi.MethodSecretsCreate,
		Params:  &params,
	}
	return s.requestAuthorizer.AuthorizeRequest(ctx, jsonRequest)
}

func (s *Capability) validateRequest(id string, encryptedSecrets []*vaultcommon.EncryptedSecret, ids []*vaultcommon.SecretIdentifier) error {
	if id == "" {
		return errors.New("request ID must not be empty")
	}
	if len(encryptedSecrets) >= vaultapi.MaxBatchSize {
		return errors.New("request batch size exceeds maximum of " + strconv.Itoa(vaultapi.MaxBatchSize))
	}
	uniqueIDs := map[string]bool{}
	for _, req := range encryptedSecrets {
		if req.Id == nil {
			errors.New("secret ID must not be nil")
		}

		if req.Id.Key == "" || req.Id.Owner == "" {
			return errors.New("secret ID must have both key and owner set: " + req.Id.String())
		}

		_, ok := uniqueIDs[KeyFor(req.Id)]
		if ok {
			return errors.New("duplicate secret ID found: " + req.Id.String())
		}

		uniqueIDs[KeyFor(req.Id)] = true
	}
	// TODO(https://smartcontract-it.atlassian.net/browse/PRIV-155): encryptedSecrets should be encrypted by the right public key
	for _, id := range ids {
		if id.Key == "" || id.Owner == "" {
			return errors.New("secret ID must have both key and owner set: " + id.String())
		}

		_, ok := uniqueIDs[KeyFor(id)]
		if ok {
			return errors.New("duplicate secret ID found: " + id.String())
		}

		uniqueIDs[KeyFor(id)] = true
	}
	return nil
}

func NewCapability(
	lggr logger.Logger,
	clock clockwork.Clock,
	expiresAfter time.Duration,
	handler *requests.Handler[*Request, *Response],
	workflowRegistrySyncer syncerV2.WorkflowRegistrySyncer,
) *Capability {
	return &Capability{
		lggr:              lggr.Named("VaultCapability"),
		clock:             clock,
		expiresAfter:      expiresAfter,
		handler:           handler,
		requestAuthorizer: vaultapi.NewRequestAuthorizer(lggr, workflowRegistrySyncer),
	}
}
