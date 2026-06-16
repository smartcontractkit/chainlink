//go:build wasip1

package main

import (
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basicaction"
	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basictrigger"
)

func CreateWorkflow(_ string, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[string], error) {
	return cre.Workflow[string]{
		cre.Handler(
			basictrigger.Trigger(&basictrigger.Config{Name: "test", Number: 0}),
			func(_ string, runtime cre.Runtime, _ *basictrigger.Outputs) (string, error) {
				basicAction := &basicaction.BasicAction{}
				basicAction.PerformAction(runtime, &basicaction.Inputs{InputThing: true})
				return "Hello, world!", nil
			},
		),
	}, nil
}

func main() {
	wasm.NewRunner(func(b []byte) (string, error) { return string(b), nil }).Run(CreateWorkflow)
}
