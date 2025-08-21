//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"
)

type runtimeConfig struct {
	Schedule string `yaml:"schedule"`
}

func RunSimpleCronWorkflow(config *runtimeConfig, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*runtimeConfig], error) {
	cfg := &cron.Config{
		Schedule: config.Schedule,
	}

	return cre.Workflow[*runtimeConfig]{
		cre.Handler(
			cron.Trigger(cfg),
			onTrigger,
		),
	}, nil
}

func onTrigger(config *runtimeConfig, runtime cre.Runtime, outputs *cron.Payload) (string, error) {
	runtime.Logger().Info("inside onTrigger handler")
	return fmt.Sprintf("success (Schedule: %s)", config.Schedule), nil
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
