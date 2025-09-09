//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"

	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
)

type runtimeConfig struct {
	Schedule string `yaml:"schedule"`
}

func RunSimpleCronWorkflow(config *runtimeConfig, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*runtimeConfig], error) {
	cfg := &cron.Config{
		Schedule: config.Schedule,
	}

	req := &sdk.SecretRequest{
		Id: "DATA_SOURCE_API_KEY",
	}

	secret, err := secretsProvider.GetSecret(req).Await()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get secret: %v", err))
		return nil, err
	}

	return cre.Workflow[*runtimeConfig]{
		cre.Handler(
			cron.Trigger(cfg),
			makeCallback(secret.Value),
		),
	}, nil
}

func makeCallback(apiKey string) func(*runtimeConfig, cre.Runtime, *cron.Payload) (string, error) {
	onTrigger := func(config *runtimeConfig, runtime cre.Runtime, outputs *cron.Payload) (string, error) {
		return fmt.Sprintf("ping (Schedule: %s, API KEY: %s)", config.Schedule, apiKey), nil
	}
	return onTrigger
}

func main() {
	wasm.NewRunner(func(b []byte) (*runtimeConfig, error) {
		cfg := &runtimeConfig{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	}).Run(RunSimpleCronWorkflow)
}
