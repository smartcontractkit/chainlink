package vault

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	vault2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/vault"
)

type SecretsService interface {
	CreateSecrets(ctx context.Context, request *vault.CreateSecretsRequest) (*vault2.Response, error)
	UpdateSecrets(ctx context.Context, request *vault.UpdateSecretsRequest) (*vault2.Response, error)
	GetSecrets(ctx context.Context, requestID string, request *vault.GetSecretsRequest) (*vault2.Response, error)
}
