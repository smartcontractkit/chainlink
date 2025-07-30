//go:build wasip1

package main

import (
	"time"

	sdk "github.com/smartcontractkit/cre-sdk-go/sdk"
	"github.com/smartcontractkit/cre-sdk-go/sdk/wasm"
)

type None struct{}

func main() {
	wasm.NewRunner(func(configBytes []byte) (None, error) {
		return None{}, nil
	}).Run(RunSimpleTimeWorkflow)
}

func RunSimpleTimeWorkflow(wcx *sdk.Environment[None]) (sdk.Workflow[None], error) {
	donTime := "donTime=" + time.Now().Format("2006-01-02 15:04:05")
	wcx.Logger.Info(donTime)
	workflows := sdk.Workflow[None]{}
	return workflows, nil
}
