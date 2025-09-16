//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type WorkflowConfig struct {
	Schedule string `yaml:"schedule,omitempty"`
}

func main() {
	wasm.NewRunner(func(configBytes []byte) (WorkflowConfig, error) {
		cfg := WorkflowConfig{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return WorkflowConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return cfg, nil
	}).Run(RunSimpleCronWorkflow)
}

func RunSimpleCronWorkflow(config WorkflowConfig, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[WorkflowConfig], error) {
	workflows := cre.Workflow[WorkflowConfig]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: config.Schedule}),
			onTrigger,
		),
	}
	return workflows, nil
}

func onTrigger(_ WorkflowConfig, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	runtime.Logger().Info("Amazing workflow user log")
	return "such a lovely disaster", nil
}
