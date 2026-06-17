package vault

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

type gatewaySecretsMethodDescriptor struct {
	errorPrefix                   string
	newRequest                    func() proto.Message
	validate                      func(ctx context.Context, v *RequestValidator, opts UserJSONRPCValidationOptions, requestID string, req proto.Message) error
	validateOwnersMatchAuthorized func(req proto.Message, workflowOwner string) error
	applyDefaults                 func(req proto.Message)
}

func lookupGatewaySecretsMethodDescriptor(method string) (gatewaySecretsMethodDescriptor, error) {
	desc, ok := gatewaySecretsMethodRegistry[method]
	if !ok {
		return gatewaySecretsMethodDescriptor{}, fmt.Errorf("unsupported gateway secrets method for validation: %s", method)
	}
	return desc, nil
}

var gatewaySecretsMethodRegistry = map[string]gatewaySecretsMethodDescriptor{
	vaulttypes.MethodSecretsCreate: {
		errorPrefix: "failed to validate create secrets request: ",
		newRequest: func() proto.Message {
			return &vaultcommon.CreateSecretsRequest{}
		},
		validate: func(ctx context.Context, v *RequestValidator, opts UserJSONRPCValidationOptions, requestID string, req proto.Message) error {
			request := req.(*vaultcommon.CreateSecretsRequest)
			request.RequestId = coalesceRequestID(request.RequestId, requestID)
			return v.ValidateCreateSecretsRequest(ctx, opts.PublicKey, request, opts.SkipLabelValidation)
		},
		validateOwnersMatchAuthorized: func(req proto.Message, workflowOwner string) error {
			return validateEncryptedSecretOwnerMismatch(req.(*vaultcommon.CreateSecretsRequest).EncryptedSecrets, workflowOwner)
		},
		applyDefaults: func(req proto.Message) {
			vaultutils.ApplyEncryptedSecretNamespaceDefaults(req.(*vaultcommon.CreateSecretsRequest).EncryptedSecrets)
		},
	},
	vaulttypes.MethodSecretsUpdate: {
		errorPrefix: "failed to validate update secrets request: ",
		newRequest: func() proto.Message {
			return &vaultcommon.UpdateSecretsRequest{}
		},
		validate: func(ctx context.Context, v *RequestValidator, opts UserJSONRPCValidationOptions, requestID string, req proto.Message) error {
			request := req.(*vaultcommon.UpdateSecretsRequest)
			request.RequestId = coalesceRequestID(request.RequestId, requestID)
			return v.ValidateUpdateSecretsRequest(ctx, opts.PublicKey, request, opts.SkipLabelValidation)
		},
		validateOwnersMatchAuthorized: func(req proto.Message, workflowOwner string) error {
			return validateEncryptedSecretOwnerMismatch(req.(*vaultcommon.UpdateSecretsRequest).EncryptedSecrets, workflowOwner)
		},
		applyDefaults: func(req proto.Message) {
			vaultutils.ApplyEncryptedSecretNamespaceDefaults(req.(*vaultcommon.UpdateSecretsRequest).EncryptedSecrets)
		},
	},
	vaulttypes.MethodSecretsDelete: {
		errorPrefix: "failed to validate delete secrets request: ",
		newRequest: func() proto.Message {
			return &vaultcommon.DeleteSecretsRequest{}
		},
		validate: func(ctx context.Context, v *RequestValidator, opts UserJSONRPCValidationOptions, requestID string, req proto.Message) error {
			request := req.(*vaultcommon.DeleteSecretsRequest)
			request.RequestId = coalesceRequestID(request.RequestId, requestID)
			return v.ValidateDeleteSecretsRequest(ctx, request)
		},
		validateOwnersMatchAuthorized: func(req proto.Message, workflowOwner string) error {
			return validateSecretIdentifierOwnerMismatch(req.(*vaultcommon.DeleteSecretsRequest).Ids, workflowOwner)
		},
		applyDefaults: func(req proto.Message) {
			vaultutils.ApplySecretIdentifierNamespaceDefaults(req.(*vaultcommon.DeleteSecretsRequest).Ids)
		},
	},
	vaulttypes.MethodSecretsList: {
		errorPrefix: "failed to validate list secret identifiers request: ",
		newRequest: func() proto.Message {
			return &vaultcommon.ListSecretIdentifiersRequest{}
		},
		validate: func(ctx context.Context, v *RequestValidator, opts UserJSONRPCValidationOptions, requestID string, req proto.Message) error {
			request := req.(*vaultcommon.ListSecretIdentifiersRequest)
			request.RequestId = coalesceRequestID(request.RequestId, requestID)
			return v.ValidateListSecretIdentifiersRequest(ctx, request)
		},
		validateOwnersMatchAuthorized: func(req proto.Message, workflowOwner string) error {
			listRequest := req.(*vaultcommon.ListSecretIdentifiersRequest)
			if vaultutils.NormalizeOwner(listRequest.Owner) != vaultutils.NormalizeOwner(workflowOwner) {
				return fmt.Errorf("list secrets owner %q does not match authorized workflow owner %q", listRequest.Owner, workflowOwner)
			}
			return nil
		},
		applyDefaults: func(req proto.Message) {
			request := req.(*vaultcommon.ListSecretIdentifiersRequest)
			if request.Namespace == "" {
				request.Namespace = vaulttypes.DefaultNamespace
			}
		},
	},
}

func setUserJSONRPCRequestID(msg proto.Message, requestID string) {
	switch m := msg.(type) {
	case *vaultcommon.CreateSecretsRequest:
		m.RequestId = requestID
	case *vaultcommon.UpdateSecretsRequest:
		m.RequestId = requestID
	case *vaultcommon.DeleteSecretsRequest:
		m.RequestId = requestID
	case *vaultcommon.ListSecretIdentifiersRequest:
		m.RequestId = requestID
	}
}

func invalidVaultParamsErrorPrefix(method string) string {
	desc, err := lookupGatewaySecretsMethodDescriptor(method)
	if err != nil {
		return ""
	}
	return desc.errorPrefix
}

func validateEncryptedSecretOwnerMismatch(encryptedSecrets []*vaultcommon.EncryptedSecret, workflowOwner string) error {
	if len(encryptedSecrets) == 0 {
		return errors.New("request batch must contain at least 1 item")
	}
	for idx, encryptedSecret := range encryptedSecrets {
		if encryptedSecret == nil {
			return fmt.Errorf("encrypted secret must not be nil at index %d", idx)
		}
		if encryptedSecret.Id == nil {
			return fmt.Errorf("secret ID must not be nil at index %d", idx)
		}
		if vaultutils.NormalizeOwner(encryptedSecret.Id.Owner) != vaultutils.NormalizeOwner(workflowOwner) {
			return fmt.Errorf("encrypted secret owner at index %d %q does not match authorized workflow owner %q", idx, encryptedSecret.Id.Owner, workflowOwner)
		}
	}
	return nil
}

func validateSecretIdentifierOwnerMismatch(ids []*vaultcommon.SecretIdentifier, workflowOwner string) error {
	if len(ids) == 0 {
		return errors.New("request batch must contain at least 1 item")
	}
	for idx, id := range ids {
		if id == nil {
			return fmt.Errorf("secret ID must not be nil at index %d", idx)
		}
		if vaultutils.NormalizeOwner(id.Owner) != vaultutils.NormalizeOwner(workflowOwner) {
			return fmt.Errorf("secret identifier owner at index %d %q does not match authorized workflow owner %q", idx, id.Owner, workflowOwner)
		}
	}
	return nil
}
