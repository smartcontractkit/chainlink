package v2

import (
	"context"
	"errors"

	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
)

type SecretsFetcher interface {
	GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error)
}

type unimplementedSecretsFetcher struct{}

func (u unimplementedSecretsFetcher) GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	return nil, errors.New("secrets fetching is not implemented")
}
