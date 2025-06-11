//go:build wasip1

package main

import (
	"fmt"

	croncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2"
	"gopkg.in/yaml.v2"
)

type runtimeConfig struct {
	Schedule string `yaml:"schedule"`
}

func RunSimpleCronWorkflow(runner sdk.DonRunner) {
	cron := &croncap.Cron{}
	cfg := &croncap.Config{
		Schedule: "*/3 * * * * *", // every 3 seconds
	}

	runner.Run(&sdk.WorkflowArgs[sdk.DonRuntime]{
		Handlers: []sdk.Handler[sdk.DonRuntime]{
			sdk.NewDonHandler(
				cron.Trigger(cfg),
				onTrigger,
			),
		},
	})
}

func onTrigger(runtime sdk.DonRuntime, outputs *croncap.Payload) (string, error) {
	b := runtime.Config()

	var cfg runtimeConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return "", err
	}

	return fmt.Sprintf("ping (Schedule: %s)", cfg.Schedule), nil
}

func main() {
	RunSimpleCronWorkflow(wasm.NewDonRunner())
}
