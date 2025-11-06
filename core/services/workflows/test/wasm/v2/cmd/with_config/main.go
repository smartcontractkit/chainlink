//go:build wasip1

package main

import (
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basictrigger"
	"gopkg.in/yaml.v3"
)

type runtimeConfig struct {
	Name   string `yaml:"name"`
	Number int32  `yaml:"number"`
}

func CreateWorkflow(runnerCfg *runtimeConfig, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*runtimeConfig], error) {
	return cre.Workflow[*runtimeConfig]{
		cre.Handler(
			basictrigger.Trigger(&basictrigger.Config{
				Name:   runnerCfg.Name,
				Number: runnerCfg.Number,
			}),
			onTrigger,
		),
	}, nil
}

func onTrigger(config *runtimeConfig, runtime cre.Runtime, _ *basictrigger.Outputs) (string, error) {
	runtime.Logger().Info("onTrigger called")
	b, err := yaml.Marshal(config)
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
