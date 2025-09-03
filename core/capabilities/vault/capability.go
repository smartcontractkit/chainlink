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
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

var _ capabilities.ExecutableCapability = (*Capability)(nil)

type Capability struct {
	lggr              logger.Logger
	clock             clockwork.Clock
	expiresAfter      time.Duration
	handler           *requests.Handler[*vaulttypes.Request, *vaulttypes.Response]
	requestAuthorizer RequestAuthorizer
}

func (s *Capability) Start(ctx context.Context) error {
	return s.handler.Start(ctx)
}

func (s *Capability) Close() error {
	return s.handler.Close()
}

func (s *Capability) Info(_ context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.NewCapabilityInfo(vaultcommon.CapabilityID, capabilities.CapabilityTypeAction, "Vault Capability")
}

func (s *Capability) RegisterToWorkflow(_ context.Context, _ capabilities.RegisterToWorkflowRequest) error {
	// Left unimplemented as this method will never be called
	// for this capability
	return nil
}

func (s *Capability) UnregisterFromWorkflow(_ context.Context, _ capabilities.UnregisterFromWorkflowRequest) error {
	// Left unimplemented as this method will never be called
	// for this capability
	return nil
}

func (s *Capability) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	if request.Payload == nil {
		return capabilities.CapabilityResponse{}, errors.New("capability does not support v1 requests")
	}

	if request.Method != vaulttypes.MethodSecretsGet {
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
	respPB := &vaultcommon.GetSecretsResponse{}
	err = proto.Unmarshal(resp.Payload, respPB)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("could not unmarshal response to GetSecretsResponse: %w", err)
	}

	anyProto, err := anypb.New(respPB)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("could not marshal response to anypb: %w", err)
	}

	return capabilities.CapabilityResponse{
		Payload: anyProto,
	}, nil
}

func ValidateCreateSecretsRequest(request *vaultcommon.CreateSecretsRequest) error {
	return validateWriteRequest(request.RequestId, request.EncryptedSecrets)
}

func validateWriteRequest(id string, encryptedSecrets []*vaultcommon.EncryptedSecret) error {
	if id == "" {
		return errors.New("request ID must not be empty")
	}
	if len(encryptedSecrets) >= vaulttypes.MaxBatchSize {
		return errors.New("request batch size exceeds maximum of " + strconv.Itoa(vaulttypes.MaxBatchSize))
	}
	if len(encryptedSecrets) == 0 {
		return errors.New("request batch must contain at least 1 item")
	}

	uniqueIDs := map[string]bool{}
	for idx, req := range encryptedSecrets {
		if req.Id == nil {
			return errors.New("secret ID must not be nil at index " + strconv.Itoa(idx))
		}

		if req.Id.Key == "" || req.Id.Owner == "" {
			return errors.New("secret ID must have both key and owner set at index " + strconv.Itoa(idx) + ":" + req.Id.String())
		}

		if req.EncryptedValue == "" {
			return errors.New("secret must have encrypted value set at index " + strconv.Itoa(idx) + ":" + req.Id.String())
		}

		_, ok := uniqueIDs[vaulttypes.KeyFor(req.Id)]
		if ok {
			return errors.New("duplicate secret ID found at index " + strconv.Itoa(idx) + ": " + req.Id.String())
		}

		uniqueIDs[vaulttypes.KeyFor(req.Id)] = true
	}

	// TODO(https://smartcontract-it.atlassian.net/browse/PRIV-155): encryptedSecrets should be encrypted by the right public key
	return nil
}

func (s *Capability) CreateSecrets(ctx context.Context, request *vaultcommon.CreateSecretsRequest) (*vaulttypes.Response, error) {
	s.lggr.Infof("Received Request: %s", request.String())
	err := ValidateCreateSecretsRequest(request)
	if err != nil {
		s.lggr.Infof("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}
	authorized, owner, err := s.isAuthorizedRequest(ctx, request, vaulttypes.MethodSecretsCreate)
	if !authorized || err != nil {
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

func ValidateUpdateSecretsRequest(request *vaultcommon.UpdateSecretsRequest) error {
	return validateWriteRequest(request.RequestId, request.EncryptedSecrets)
}

func (s *Capability) UpdateSecrets(ctx context.Context, request *vaultcommon.UpdateSecretsRequest) (*vaulttypes.Response, error) {
	s.lggr.Infof("Received Request: %s", request.String())
	err := ValidateUpdateSecretsRequest(request)
	if err != nil {
		s.lggr.Infof("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}
	authorized, owner, err := s.isAuthorizedRequest(ctx, request, vaulttypes.MethodSecretsUpdate)
	if !authorized || err != nil {
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

func ValidateDeleteSecretsRequest(request *vaultcommon.DeleteSecretsRequest) error {
	if request.RequestId == "" {
		return errors.New("request ID must not be empty")
	}
	if len(request.Ids) >= vaulttypes.MaxBatchSize {
		return errors.New("request batch size exceeds maximum of " + strconv.Itoa(vaulttypes.MaxBatchSize))
	}

	uniqueIDs := map[string]bool{}
	for idx, id := range request.Ids {
		if id.Key == "" || id.Owner == "" {
			return errors.New("secret ID must have both key and owner set at index " + strconv.Itoa(idx) + ": " + id.String())
		}

		_, ok := uniqueIDs[vaulttypes.KeyFor(id)]
		if ok {
			return errors.New("duplicate secret ID found at index " + strconv.Itoa(idx) + ": " + id.String())
		}

		uniqueIDs[vaulttypes.KeyFor(id)] = true
	}
	return nil
}

func (s *Capability) DeleteSecrets(ctx context.Context, request *vaultcommon.DeleteSecretsRequest) (*vaulttypes.Response, error) {
	s.lggr.Infof("Received Request: %s", request.String())
	err := ValidateDeleteSecretsRequest(request)
	if err != nil {
		s.lggr.Infof("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}

	authorized, owner, err := s.isAuthorizedRequest(ctx, request, vaulttypes.MethodSecretsDelete)
	if !authorized || err != nil {
		s.lggr.Infof("Request [%s] not authorized for owner: %s", request.String(), owner)
		return nil, errors.New("request not authorized: " + err.Error())
	}
	if !strings.HasPrefix(request.RequestId, owner) {
		// Gateway should ensure it prefixes request ids with the owner, to ensure request uniqueness
		s.lggr.Infof("Request ID: [%s] must start with owner address: [%s]", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " must start with owner address: " + owner)
	}
	for _, req := range request.Ids {
		// Right owner for secrets can only be set here, after authorization
		// This ensures that users cannot access secrets belonging to other owners
		req.Owner = owner
	}
	s.lggr.Infof("Processing authorized and normalized request [%s]", request.String())
	return s.handleRequest(ctx, request.RequestId, request)
}

func ValidateGetSecretsRequest(request *vaultcommon.GetSecretsRequest) error {
	if len(request.Requests) == 0 {
		return errors.New("no GetSecret request specified in request")
	}
	if len(request.Requests) >= vaulttypes.MaxBatchSize {
		return fmt.Errorf("request batch size exceeds maximum of %d", vaulttypes.MaxBatchSize)
	}

	for idx, req := range request.Requests {
		if req.Id.Key == "" || req.Id.Owner == "" {
			return errors.New("secret ID must have both key and owner set at index " + strconv.Itoa(idx) + ": " + req.Id.String())
		}
	}

	return nil
}

func (s *Capability) GetSecrets(ctx context.Context, requestID string, request *vaultcommon.GetSecretsRequest) (*vaulttypes.Response, error) {
	s.lggr.Infof("Received Request: %s", request.String())
	if err := ValidateGetSecretsRequest(request); err != nil {
		s.lggr.Infof("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}

	// No auth needed, as this method is not exposed externally
	return s.handleRequest(ctx, requestID, request)
}

func ValidateListSecretIdentifiersRequest(request *vaultcommon.ListSecretIdentifiersRequest) error {
	if request.RequestId == "" {
		return errors.New("request ID must not be empty")
	}
	if request.Owner == "" {
		return errors.New("owner must not be empty")
	}
	return nil
}

func (s *Capability) ListSecretIdentifiers(ctx context.Context, request *vaultcommon.ListSecretIdentifiersRequest) (*vaulttypes.Response, error) {
	s.lggr.Infof("Received Request: %s", request.String())
	err := ValidateListSecretIdentifiersRequest(request)
	if err != nil {
		s.lggr.Infof("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}

	authorized, owner, err := s.isAuthorizedRequest(ctx, request, vaulttypes.MethodSecretsList)
	if !authorized || err != nil {
		s.lggr.Infof("Request [%s] not authorized for owner: %s", request.String(), owner)
		return nil, errors.New("request not authorized: " + err.Error())
	}
	if !strings.HasPrefix(request.RequestId, owner) {
		// Gateway should ensure it prefixes request ids with the owner, to ensure request uniqueness
		s.lggr.Infof("Request ID: [%s] must start with owner address: [%s]", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " must start with owner address: " + owner)
	}
	// Right owner for secrets can only be set here, after authorization
	// This ensures that users cannot access secrets belonging to other owners
	request.Owner = owner

	s.lggr.Infof("Processing authorized and normalized request [%s]", request.String())
	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Capability) handleRequest(ctx context.Context, requestID string, request proto.Message) (*vaulttypes.Response, error) {
	respCh := make(chan *vaulttypes.Response, 1)
	s.handler.SendRequest(ctx, &vaulttypes.Request{
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

func (s *Capability) isAuthorizedRequest(ctx context.Context, request any, method string) (bool, string, error) {
	var params json.RawMessage
	params, err := json.Marshal(request)
	if err != nil {
		return false, "", fmt.Errorf("could not marshal CreateSecretsRequest: %w", err)
	}
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  method,
		Params:  &params,
	}
	return s.requestAuthorizer.AuthorizeRequest(ctx, jsonRequest)
}

func NewCapability(
	lggr logger.Logger,
	clock clockwork.Clock,
	expiresAfter time.Duration,
	handler *requests.Handler[*vaulttypes.Request, *vaulttypes.Response],
	requestAuthorizer RequestAuthorizer,
) *Capability {
	return &Capability{
		lggr:              logger.Named(lggr, "VaultCapability"),
		clock:             clock,
		expiresAfter:      expiresAfter,
		handler:           handler,
		requestAuthorizer: requestAuthorizer,
	}
}
