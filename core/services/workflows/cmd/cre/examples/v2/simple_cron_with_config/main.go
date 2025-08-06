//go:build wasip1

package main

import (
	"fmt"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"
)

type runtimeConfig struct {
	Schedule string `yaml:"schedule"`
}

func RunSimpleCronWorkflow(env *cre.Environment[*runtimeConfig]) (cre.Workflow[*runtimeConfig], error) {
	cfg := &cron.Config{
		Schedule: env.Config.Schedule,
	}

	return cre.Workflow[*runtimeConfig]{
		cre.Handler(
			cron.Trigger(cfg),
			onTrigger,
		),
	}, nil
}

func onTrigger(env *cre.Environment[*runtimeConfig], runtime cre.Runtime, outputs *cron.Payload) (string, error) {
	env.Logger.Info("inside onTrigger handler")
	return fmt.Sprintf("success (Schedule: %s)", env.Config.Schedule), nil
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
