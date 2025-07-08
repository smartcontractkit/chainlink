//go:build wasip1

package main

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/cre-sdk-go/sdk"
	"github.com/smartcontractkit/cre-sdk-go/sdk/wasm"
	"gopkg.in/yaml.v3"
)

type runtimeConfig struct {
	Schedule string `yaml:"schedule"`
}

func RunSimpleCronWorkflow(env *sdk.Environment[*runtimeConfig]) (sdk.Workflow[*runtimeConfig], error) {
	cfg := &cron.Config{
		Schedule: env.Config.Schedule,
	}

	req := &pb.SecretRequest{
		Id: "DATA_SOURCE_API_KEY",
	}

	secret, err := env.GetSecret(req).Await()
	if err != nil {
		env.Logger.Error(fmt.Sprintf("failed to get secret: %v", err))
		return nil, err
	}

	return sdk.Workflow[*runtimeConfig]{
		sdk.Handler(
			cron.Trigger(cfg),
			makeCallback(secret.Value),
		),
	}, nil
}

func makeCallback(apiKey string) func(*sdk.Environment[*runtimeConfig], sdk.Runtime, *cron.Payload) (string, error) {
	onTrigger := func(env *sdk.Environment[*runtimeConfig], runtime sdk.Runtime, outputs *cron.Payload) (string, error) {
		return fmt.Sprintf("ping (Schedule: %s, API KEY: %s)", env.Config.Schedule, apiKey), nil
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
