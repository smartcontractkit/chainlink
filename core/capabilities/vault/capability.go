package vault

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jonboulle/clockwork"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

var _ capabilities.ExecutableCapability = (*Capability)(nil)

type Capability struct {
	lggr                 logger.Logger
	clock                clockwork.Clock
	expiresAfter         time.Duration
	handler              *requests.Handler[*vaulttypes.Request, *vaulttypes.Response]
	capabilitiesRegistry core.CapabilitiesRegistry
	publicKey            *LazyPublicKey
	linker               *OrgIDToWorkflowOwnerLinker
	*RequestValidator
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

	if lerr := s.MaxRequestBatchSizeLimiter.Close(); lerr != nil {
		err = errors.Join(err, fmt.Errorf("error closing request batch size limiter: %w", lerr))
	}
	if lerr := s.linker.Close(); lerr != nil {
		err = errors.Join(err, fmt.Errorf("error closing org_id linker: %w", lerr))
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

	for idx, req := range r.Requests {
		if req == nil { // defensive: protobuf strips nil elements, but guard against in-process callers
			s.lggr.Errorw("get secrets request contains nil secret request", "index", idx)
			return capabilities.CapabilityResponse{}, fmt.Errorf("nil secret request at index %d", idx)
		}
		if req.Id != nil && normalizeOwner(req.Id.Owner) != normalizeOwner(request.Metadata.WorkflowOwner) {
			s.lggr.Errorw("get secrets request owner mismatch", "index", idx, "secretOwner", req.Id.Owner, "workflowOwner", request.Metadata.WorkflowOwner)
			return capabilities.CapabilityResponse{}, fmt.Errorf("secret identifier owner %q does not match workflow owner %q at index %d", req.Id.Owner, request.Metadata.WorkflowOwner, idx)
		}
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

func (s *Capability) CreateSecrets(ctx context.Context, request *vaultcommon.CreateSecretsRequest) (*vaulttypes.Response, error) {
	s.lggr.Debugf("Received Request: %s", request.String())
	err := s.ValidateCreateSecretsRequest(s.publicKey.Get(), request)
	if err != nil {
		s.lggr.Debugf("RequestId: [%s] failed validation checks: %s", request.RequestId, err.Error())
		return nil, err
	}
	resolvedIdentity, err := s.resolveRequestIdentity(ctx, request.OrgId, request.WorkflowOwner)
	if err != nil {
		return nil, err
	}
	request.OrgId = resolvedIdentity.OrgID
	request.WorkflowOwner = resolvedIdentity.WorkflowOwner
	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Capability) UpdateSecrets(ctx context.Context, request *vaultcommon.UpdateSecretsRequest) (*vaulttypes.Response, error) {
	s.lggr.Debugf("Received Request: %s", request.String())
	err := s.ValidateUpdateSecretsRequest(s.publicKey.Get(), request)
	if err != nil {
		s.lggr.Debugf("RequestId: [%s] failed validation checks: %s", request.RequestId, err.Error())
		return nil, err
	}
	resolvedIdentity, err := s.resolveRequestIdentity(ctx, request.OrgId, request.WorkflowOwner)
	if err != nil {
		return nil, err
	}
	request.OrgId = resolvedIdentity.OrgID
	request.WorkflowOwner = resolvedIdentity.WorkflowOwner
	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Capability) DeleteSecrets(ctx context.Context, request *vaultcommon.DeleteSecretsRequest) (*vaulttypes.Response, error) {
	s.lggr.Debugf("Received Request: %s", request.String())
	err := s.ValidateDeleteSecretsRequest(request)
	if err != nil {
		s.lggr.Debugf("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}
	resolvedIdentity, err := s.resolveRequestIdentity(ctx, request.OrgId, request.WorkflowOwner)
	if err != nil {
		return nil, err
	}
	request.OrgId = resolvedIdentity.OrgID
	request.WorkflowOwner = resolvedIdentity.WorkflowOwner
	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Capability) GetSecrets(ctx context.Context, requestID string, request *vaultcommon.GetSecretsRequest) (*vaulttypes.Response, error) {
	s.lggr.Debugf("Received Request: %s", request.String())
	if err := s.ValidateGetSecretsRequest(request); err != nil {
		s.lggr.Debugf("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}

	// No auth needed, as this method is not exposed externally
	return s.handleRequest(ctx, requestID, request)
}

func (s *Capability) ListSecretIdentifiers(ctx context.Context, request *vaultcommon.ListSecretIdentifiersRequest) (*vaulttypes.Response, error) {
	s.lggr.Debugf("Received Request: %s", request.String())
	err := s.ValidateListSecretIdentifiersRequest(request)
	if err != nil {
		s.lggr.Debugf("Request: [%s] failed validation checks: %s", request.String(), err.Error())
		return nil, err
	}
	resolvedIdentity, err := s.resolveRequestIdentity(ctx, request.OrgId, request.WorkflowOwner)
	if err != nil {
		return nil, err
	}
	request.OrgId = resolvedIdentity.OrgID
	request.WorkflowOwner = resolvedIdentity.WorkflowOwner
	return s.handleRequest(ctx, request.RequestId, request)
}

func (s *Capability) GetPublicKey(ctx context.Context, request *vaultcommon.GetPublicKeyRequest) (*vaultcommon.GetPublicKeyResponse, error) {
	l := logger.With(s.lggr, "method", "GetPublicKey")
	l.Debugf("Received Request: GetPublicKeyRequest")

	pubKey := s.publicKey.Get()
	if pubKey == nil {
		l.Debug("could not get public key: is the plugin initialized?")
		return nil, errors.New("could not get public key: is the plugin initialized?")
	}

	pkb, err := pubKey.Marshal()
	if err != nil {
		l.Debugf("could not marshal public key: %s", err.Error())
		return nil, fmt.Errorf("could not marshal public key: %w", err)
	}

	return &vaultcommon.GetPublicKeyResponse{
		PublicKey: hex.EncodeToString(pkb),
	}, nil
}

func normalizeOwner(owner string) string {
	return strings.ToLower(strings.TrimPrefix(owner, "0x"))
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

// resolveRequestIdentity validates and normalizes the org/workflow-owner pair that the vault plugin consumes.
func (s *Capability) resolveRequestIdentity(ctx context.Context, orgID string, workflowOwner string) (LinkedVaultRequestIdentity, error) {
	s.lggr.Debugw("resolving request identity", "orgID", orgID, "workflowOwner", workflowOwner)
	linked, err := s.linker.Link(ctx, orgID, workflowOwner)
	if err != nil {
		s.lggr.Errorw("failed to resolve request identity", "orgID", orgID, "workflowOwner", workflowOwner, "err", err)
		return LinkedVaultRequestIdentity{}, err
	}
	s.lggr.Debugw("resolved request identity", "orgID", linked.OrgID, "workflowOwner", linked.WorkflowOwner)

	return linked, nil
}

func NewCapability(
	lggr logger.Logger,
	clock clockwork.Clock,
	expiresAfter time.Duration,
	handler *requests.Handler[*vaulttypes.Request, *vaulttypes.Response],
	capabilitiesRegistry core.CapabilitiesRegistry,
	publicKey *LazyPublicKey,
	orgResolver orgresolver.OrgResolver,
	limitsFactory limits.Factory,
) (*Capability, error) {
	limiter, err := limits.MakeBoundLimiter(limitsFactory, cresettings.Default.VaultRequestBatchSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("could not create request batch size limiter: %w", err)
	}
	linker, err := NewOrgIDToWorkflowOwnerLinker(orgResolver, limitsFactory)
	if err != nil {
		return nil, err
	}
	ciphertextLimiter, err := limits.MakeUpperBoundLimiter(limitsFactory, cresettings.Default.VaultCiphertextSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("could not create ciphertext size limiter: %w", err)
	}
	return &Capability{
		lggr:                 logger.Named(lggr, "VaultCapability"),
		clock:                clock,
		expiresAfter:         expiresAfter,
		handler:              handler,
		capabilitiesRegistry: capabilitiesRegistry,
		publicKey:            publicKey,
		linker:               linker,
		RequestValidator:     NewRequestValidator(limiter, ciphertextLimiter),
	}, nil
}
