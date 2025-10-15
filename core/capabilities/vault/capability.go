package vault

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

var _ capabilities.ExecutableCapability = (*Capability)(nil)

type Capability struct {
	lggr                 logger.Logger
	clock                clockwork.Clock
	expiresAfter         time.Duration
	handler              *requests.Handler[*vaulttypes.Request, *vaulttypes.Response]
	requestAuthorizer    RequestAuthorizer
	capabilitiesRegistry core.CapabilitiesRegistry
	publicKey            *LazyPublicKey
}

func (s *Capability) Start(ctx context.Context) error {
	if err := s.handler.Start(ctx); err != nil {
		return fmt.Errorf("error starting vault DON request handler: %w", err)
	}

	closeHandler := func() {
		ierr := s.handler.Close()
		if ierr != nil {
			s.lggr.Errorf("error closing vault DON request handler after failed registration: %v", ierr)
		}
	}

	err := s.capabilitiesRegistry.Add(ctx, s)
	if err != nil {
		closeHandler()
		return fmt.Errorf("error registering vault capability: %w", err)
	}

	return nil
}

func (s *Capability) Close() error {
	err := s.capabilitiesRegistry.Remove(context.Background(), vaultcommon.CapabilityID)
	if err != nil {
		err = fmt.Errorf("error unregistering vault capability: %w", err)
	}

	ierr := s.handler.Close()
	if ierr != nil {
		err = errors.Join(err, fmt.Errorf("error closing vault DON request handler: %w", ierr))
	}

	return err
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

func ValidateCreateSecretsRequest(publicKey *tdh2easy.PublicKey, request *vaultcommon.CreateSecretsRequest) error {
	return validateWriteRequest(publicKey, request.RequestId, request.EncryptedSecrets)
}

// validateWriteRequest performs common validation for CreateSecrets and UpdateSecrets requests
// It treats publicKey as optional, since it can be nil if the gateway nodes don't have the public key cached yet
func validateWriteRequest(publicKey *tdh2easy.PublicKey, id string, encryptedSecrets []*vaultcommon.EncryptedSecret) error {
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
	cipherText := &tdh2easy.Ciphertext{}
	for idx, req := range encryptedSecrets {
		if req.Id == nil {
			return errors.New("secret ID must not be nil at index " + strconv.Itoa(idx))
		}

		if req.Id.Key == "" || req.Id.Namespace == "" || req.Id.Owner == "" {
			return errors.New("secret ID must have key, namespace and owner set at index " + strconv.Itoa(idx) + ":" + req.Id.String())
		}

		if req.EncryptedValue == "" {
			return errors.New("secret must have encrypted value set at index " + strconv.Itoa(idx) + ":" + req.Id.String())
		}

		// Validate that the encrypted value was indeed encrypted by the Vault public key
		cipherBytes, err := hex.DecodeString(req.EncryptedValue)
		if err != nil {
			return errors.New("failed to decode encrypted value at index " + strconv.Itoa(idx) + ":" + err.Error())
		}
		if publicKey != nil { // Public key can be nil if gateway cache isn't populated yet
			err = cipherText.UnmarshalVerify(cipherBytes, publicKey)
			if err != nil {
				return errors.New("failed to verify encrypted value at index " + strconv.Itoa(idx) + ":" + err.Error())
			}
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
	err := ValidateCreateSecretsRequest(s.publicKey.Get(), request)
	if err != nil {
		s.lggr.Infof("RequestId: [%s] failed validation checks: %s", request.RequestId, err.Error())
		return nil, err
	}
	authorized, owner, err := s.authorizeCreateSecrets(ctx, *request) //nolint:govet // The mutex isn't used
	if !authorized || err != nil {
		s.lggr.Infof("Request Id[%s] not authorized for owner: %s", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " not authorized: " + err.Error())
	}
	if !strings.HasPrefix(request.RequestId, owner) {
		// Gateway should ensure it prefixes request ids with the owner, to ensure request uniqueness
		s.lggr.Infof("Request ID: [%s] must start with owner address: [%s]", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " must start with owner address: " + owner)
	}
	for idx, req := range request.EncryptedSecrets {
		// Ensure that users cannot access secrets belonging to other owners
		if req.Id.Owner != owner {
			s.lggr.Infof("Secret ID owner: [%s] does not match authorized owner: [%s]", req.Id.Owner, owner)
			return nil, errors.New("secret ID owner: " + req.Id.Owner + " does not match authorized owner: " + owner + " at index " + strconv.Itoa(idx))
		}
	}
	s.lggr.Infof("Processing authorized and normalized request [%s]", request.String())
	return s.handleRequest(ctx, request.RequestId, request)
}

func ValidateUpdateSecretsRequest(publicKey *tdh2easy.PublicKey, request *vaultcommon.UpdateSecretsRequest) error {
	return validateWriteRequest(publicKey, request.RequestId, request.EncryptedSecrets)
}

func (s *Capability) UpdateSecrets(ctx context.Context, request *vaultcommon.UpdateSecretsRequest) (*vaulttypes.Response, error) {
	s.lggr.Infof("Received Request: %s", request.String())
	err := ValidateUpdateSecretsRequest(s.publicKey.Get(), request)
	if err != nil {
		s.lggr.Infof("RequestId: [%s] failed validation checks: %s", request.RequestId, err.Error())
		return nil, err
	}
	authorized, owner, err := s.authorizeUpdateSecrets(ctx, *request) //nolint:govet // The mutex isn't used
	if !authorized || err != nil {
		s.lggr.Infof("Request Id[%s] not authorized for owner: %s", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " not authorized: " + err.Error())
	}
	if !strings.HasPrefix(request.RequestId, owner) {
		// Gateway should ensure it prefixes request ids with the owner, to ensure request uniqueness
		s.lggr.Infof("Request ID: [%s] must start with owner address: [%s]", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " must start with owner address: " + owner)
	}
	for idx, req := range request.EncryptedSecrets {
		// Ensure that users cannot access secrets belonging to other owners
		if req.Id.Owner != owner {
			s.lggr.Infof("Secret ID owner: [%s] does not match authorized owner: [%s]", req.Id.Owner, owner)
			return nil, errors.New("secret ID owner: " + req.Id.Owner + " does not match authorized owner: " + owner + " at index " + strconv.Itoa(idx))
		}
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
		if id.Key == "" || id.Namespace == "" || id.Owner == "" {
			return errors.New("secret ID must have key, namespace and owner set at index " + strconv.Itoa(idx) + ": " + id.String())
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

	authorized, owner, err := s.authorizeDeleteSecrets(ctx, *request) //nolint:govet // The mutex isn't used
	if !authorized || err != nil {
		s.lggr.Infof("Request Id[%s] not authorized for owner: %s", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " not authorized: " + err.Error())
	}
	if !strings.HasPrefix(request.RequestId, owner) {
		// Gateway should ensure it prefixes request ids with the owner, to ensure request uniqueness
		s.lggr.Infof("Request ID: [%s] must start with owner address: [%s]", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " must start with owner address: " + owner)
	}
	for idx, req := range request.Ids {
		// Ensure that users cannot access secrets belonging to other owners
		if req.Owner != owner {
			s.lggr.Infof("Secret ID owner: [%s] does not match authorized owner: [%s]", req.Owner, owner)
			return nil, errors.New("secret ID owner: " + req.Owner + " does not match authorized owner: " + owner + " at index " + strconv.Itoa(idx))
		}
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
		if req.Id == nil {
			return errors.New("secret ID must have id set at index " + strconv.Itoa(idx))
		}
		if req.Id.Key == "" {
			return errors.New("secret ID must have key set at index " + strconv.Itoa(idx) + ": " + req.Id.String())
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
	if request.RequestId == "" || request.Owner == "" || request.Namespace == "" {
		return errors.New("requestID, owner or namespace must not be empty")
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

	authorized, owner, err := s.authorizeListSecrets(ctx, *request) //nolint:govet // The mutex isn't used
	if !authorized || err != nil {
		s.lggr.Infof("Request ID[%s] not authorized for owner: %s", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " not authorized: " + err.Error())
	}
	if !strings.HasPrefix(request.RequestId, owner) {
		// Gateway should ensure it prefixes request ids with the owner, to ensure request uniqueness
		s.lggr.Infof("Request ID: [%s] must start with owner address: [%s]", request.RequestId, owner)
		return nil, errors.New("request ID: " + request.RequestId + " must start with owner address: " + owner)
	}
	// Ensures that users cannot access secrets belonging to other owners
	request.Owner = owner
	if request.Owner != owner {
		s.lggr.Infof("Secret ID owner: [%s] does not match authorized owner: [%s]", request.Owner, owner)
		return nil, errors.New("secret ID owner: " + request.Owner + " does not match authorized owner: " + owner)
	}

	s.lggr.Infof("Processing authorized and normalized request [%s]", request.String())
	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Capability) GetPublicKey(ctx context.Context, request *vaultcommon.GetPublicKeyRequest) (*vaultcommon.GetPublicKeyResponse, error) {
	l := logger.With(s.lggr, "method", "GetPublicKey")
	l.Infof("Received Request: GetPublicKeyRequest")

	pubKey := s.publicKey.Get()
	if pubKey == nil {
		l.Info("could not get public key: is the plugin initialized?")
		return nil, errors.New("could not get public key: is the plugin initialized?")
	}

	pkb, err := pubKey.Marshal()
	if err != nil {
		l.Infof("could not marshal public key: %s", err.Error())
		return nil, fmt.Errorf("could not marshal public key: %w", err)
	}

	return &vaultcommon.GetPublicKeyResponse{
		PublicKey: hex.EncodeToString(pkb),
	}, nil
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

func (s *Capability) getOriginalRequestID(transformedRequestID string) (string, error) {
	// The transformed RequestID provided to Vault Nodes is of format <owner>::<user-provided-id>.
	// However, the RequestAuthorizer expects just the <user-provided-id> as the JSONRequest's ID fields,
	// since that's what was used by the caller when generating the request digest.
	requestIDParts := strings.Split(transformedRequestID, vaulttypes.RequestIDSeparator)
	if len(requestIDParts) != 2 {
		return "", errors.New("internal error: request ID must be in format <owner>::<user-provided-id>")
	}
	return requestIDParts[1], nil
}

func (s *Capability) authorizeCreateSecrets(ctx context.Context, request vaultcommon.CreateSecretsRequest) (bool, string, error) { //nolint:govet // The mutex isn't used
	originalRequestID, err := s.getOriginalRequestID(request.RequestId)
	if err != nil {
		return false, "", err
	}
	request.RequestId = originalRequestID

	return s.isAuthorizedRequest(ctx, &request, originalRequestID, vaulttypes.MethodSecretsCreate)
}

func (s *Capability) authorizeUpdateSecrets(ctx context.Context, request vaultcommon.UpdateSecretsRequest) (bool, string, error) { //nolint:govet // The mutex isn't used
	originalRequestID, err := s.getOriginalRequestID(request.RequestId)
	if err != nil {
		return false, "", err
	}
	request.RequestId = originalRequestID
	return s.isAuthorizedRequest(ctx, &request, originalRequestID, vaulttypes.MethodSecretsUpdate)
}

func (s *Capability) authorizeDeleteSecrets(ctx context.Context, request vaultcommon.DeleteSecretsRequest) (bool, string, error) { //nolint:govet // The mutex isn't used
	originalRequestID, err := s.getOriginalRequestID(request.RequestId)
	if err != nil {
		return false, "", err
	}
	request.RequestId = originalRequestID
	return s.isAuthorizedRequest(ctx, &request, originalRequestID, vaulttypes.MethodSecretsDelete)
}

func (s *Capability) authorizeListSecrets(ctx context.Context, request vaultcommon.ListSecretIdentifiersRequest) (bool, string, error) { //nolint:govet // The mutex isn't used
	originalRequestID, err := s.getOriginalRequestID(request.RequestId)
	if err != nil {
		return false, "", err
	}
	request.RequestId = originalRequestID
	return s.isAuthorizedRequest(ctx, &request, originalRequestID, vaulttypes.MethodSecretsList)
}

func (s *Capability) isAuthorizedRequest(ctx context.Context, request any, requestID, method string) (bool, string, error) {
	var params json.RawMessage
	params, err := json.Marshal(request)
	if err != nil {
		return false, "", fmt.Errorf("could not marshal CreateSecretsRequest: %w", err)
	}
	s.lggr.Debugw("Authorizing request", "method", method, "requestID", requestID)
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      requestID,
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
	capabilitiesRegistry core.CapabilitiesRegistry,
	publicKey *LazyPublicKey,
) *Capability {
	return &Capability{
		lggr:                 logger.Named(lggr, "VaultCapability"),
		clock:                clock,
		expiresAfter:         expiresAfter,
		handler:              handler,
		requestAuthorizer:    requestAuthorizer,
		capabilitiesRegistry: capabilitiesRegistry,
		publicKey:            publicKey,
	}
}
