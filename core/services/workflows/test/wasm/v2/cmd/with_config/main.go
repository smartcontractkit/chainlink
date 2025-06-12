//go:build wasip1

package main

import (
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/protoc/pkg/test_capabilities/basictrigger"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2"
	"gopkg.in/yaml.v3"
)

type runtimeConfig struct {
	Name   string `yaml:"name"`
	Number int32  `yaml:"number"`
}

func RunConfiguredWorkflow(runner sdk.DonRunner) {
	basic := &basictrigger.Basic{}
	b := runner.Config()
	var runnerCfg runtimeConfig
	if err := yaml.Unmarshal(b, &runnerCfg); err != nil {
		panic(err)
	}

	runner.Run(&sdk.WorkflowArgs[sdk.DonRuntime]{
		Handlers: []sdk.Handler[sdk.DonRuntime]{
			sdk.NewDonHandler(
				basic.Trigger(&basictrigger.Config{
					Name:   runnerCfg.Name,
					Number: runnerCfg.Number,
				}),
				onTrigger,
			),
		},
	})
}

func onTrigger(runtime sdk.DonRuntime, outputs *basictrigger.Outputs) (string, error) {
	b := runtime.Config()

	var cfg runtimeConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return "", err
	}

	return string(b), nil
}

func main() {
	RunConfiguredWorkflow(wasm.NewDonRunner())
}
