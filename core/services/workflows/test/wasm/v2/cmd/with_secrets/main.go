//go:build wasip1

package main

import (
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basictrigger"
	"github.com/smartcontractkit/cre-sdk-go/sdk"
	"github.com/smartcontractkit/cre-sdk-go/sdk/wasm"
	"gopkg.in/yaml.v3"
)

type runtimeConfig struct {
	Name   string `yaml:"name"`
	Number int32  `yaml:"number"`
}

func CreateWorkflow(env *sdk.Environment[*runtimeConfig]) (sdk.Workflow[*runtimeConfig], error) {
	runnerCfg := env.Config
	return sdk.Workflow[*runtimeConfig]{
		sdk.Handler(
			basictrigger.Trigger(&basictrigger.Config{
				Name:   runnerCfg.Name,
				Number: runnerCfg.Number,
			}),
			onTrigger,
		),
	}, nil
}

func onTrigger(env *sdk.Environment[*runtimeConfig], _ sdk.Runtime, _ *basictrigger.Outputs) (string, error) {
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
