//go:build wasip1

package main

import (
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"github.com/smartcontractkit/cre-sdk-go/internal_testing/capabilities/basictrigger"
)

func CreateWorkflow(_ string, _ *slog.Logger, _ cre.SecretsProvider) (cre.TeeWorkflow[string], error) {
	return cre.TeeWorkflow[string]{
		cre.HandlerInTee(
			basictrigger.Trigger(&basictrigger.Config{Name: "test", Number: 0}),
			func(_ string, _ cre.TeeRuntime, _ *basictrigger.Outputs) (string, error) {
				return "Hello, world!", nil
			},
		),
	}, nil
}

func main() {
	wasm.NewTeeRunner(cre.AnyTee{}, func(b []byte) (string, error) { return string(b), nil }).Run(CreateWorkflow)
}
