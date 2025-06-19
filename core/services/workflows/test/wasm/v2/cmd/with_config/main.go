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

func CreateWorkflow(wcx *sdk.WorkflowContext[*runtimeConfig]) (sdk.Workflow[*runtimeConfig], error) {
	runnerCfg := wcx.Config
	return sdk.Workflow[*runtimeConfig]{
		sdk.On(
			basictrigger.Trigger(&basictrigger.Config{
				Name:   runnerCfg.Name,
				Number: runnerCfg.Number,
			}),
			onTrigger,
		),
	}, nil
}

func onTrigger(wcx *sdk.WorkflowContext[*runtimeConfig], _ sdk.Runtime, _ *basictrigger.Outputs) (string, error) {
	wcx.Logger.Info("onTrigger called")
	b, err := yaml.Marshal(wcx.Config)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func main() {
	wasm.NewRunner(func(b []byte) (*runtimeConfig, error) {
		tmp := &runtimeConfig{}
		if err := yaml.Unmarshal(b, tmp); err != nil {
			return nil, err
		}
		return tmp, nil
	}).Run(CreateWorkflow)
}
