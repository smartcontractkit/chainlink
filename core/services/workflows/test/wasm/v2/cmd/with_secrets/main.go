//go:build wasip1

package main

import (
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basictrigger"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
)

type runtimeConfig struct {
	Name   string `yaml:"name"`
	Number int32  `yaml:"number"`
}

func CreateWorkflow(env *cre.Environment[*runtimeConfig]) (cre.Workflow[*runtimeConfig], error) {
	runnerCfg := env.Config
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

func onTrigger(env *cre.Environment[*runtimeConfig], _ cre.Runtime, _ *basictrigger.Outputs) (string, error) {
	secret, err := env.GetSecret(&pb.SecretRequest{Namespace: "Default", Id: "Foo"}).Await()
	if err != nil {
		return "", err
	}
	return secret.Value, nil
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
