package v2

import (
	"context"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
)

type localSecretsFetcher struct {
	secrets      map[string]string
	owner        string
	callsCounter metric.Int64Counter
}

func NewLocalSecretsFetcher(owner string, secrets map[string]string) SecretsFetcher {
	var callsCounter metric.Int64Counter
	if c, err := beholder.GetMeter().Int64Counter(
		"platform_engine_local_secrets_getsecrets_total",
		metric.WithUnit("1"),
		metric.WithDescription("Workflow local secrets override fetcher batch GetSecrets invocations"),
	); err == nil {
		callsCounter = c
	}
	resolvedOwner := owner
	if n, err := normalizeOwner(owner); err == nil {
		resolvedOwner = n
	}
	return &localSecretsFetcher{
		secrets:      secrets,
		owner:        resolvedOwner,
		callsCounter: callsCounter,
	}
}

func (f *localSecretsFetcher) GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	if f.callsCounter != nil {
		f.callsCounter.Add(ctx, 1)
	}
	responses := make([]*sdkpb.SecretResponse, 0, len(request.Requests))
	for _, req := range request.Requests {
		value, ok := f.secrets[req.Id]
		if !ok {
			responses = append(responses, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Id:        req.Id,
						Namespace: req.Namespace,
						Owner:     f.owner,
						Error:     "secret not found in local secrets",
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
					Owner:     f.owner,
					Value:     value,
				},
			},
		})
	}
	return responses, nil
}
