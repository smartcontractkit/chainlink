//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"
)

func main() {
	wasm.NewRunner(func(configBytes []byte) (types.WorkflowConfig, error) {
		cfg := types.WorkflowConfig{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return types.WorkflowConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return cfg, nil
	}).Run(RunSimpleCronWorkflow)
}

func RunSimpleCronWorkflow(config types.WorkflowConfig, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[types.WorkflowConfig], error) {
	workflows := cre.Workflow[types.WorkflowConfig]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: config.Schedule}),
			onTrigger,
		),
	}
	return workflows, nil
}

func onTrigger(_ types.WorkflowConfig, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	runtime.Logger().Info("Amazing workflow user log")
	return "such a lovely disaster", nil
}
