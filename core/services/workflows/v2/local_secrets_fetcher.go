package v2

import (
	"context"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
)

type localSecretsFetcher struct {
	secrets map[string]string
}

func NewLocalSecretsFetcher(secrets map[string]string) SecretsFetcher {
	return &localSecretsFetcher{secrets: secrets}
}

func (f *localSecretsFetcher) GetSecrets(_ context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	responses := make([]*sdkpb.SecretResponse, 0, len(request.Requests))
	for _, req := range request.Requests {
		value, ok := f.secrets[req.Id]
		if !ok {
			responses = append(responses, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Id:    req.Id,
						Error: "secret not found in local secrets",
					},
				},
			})
			continue
		}

		responses = append(responses, &sdkpb.SecretResponse{
			Response: &sdkpb.SecretResponse_Secret{
				Secret: &sdkpb.Secret{
					Id:        req.Id,
					Namespace: req.Namespace,
					Value:     value,
				},
			},
		})
	}
	return responses, nil
}
