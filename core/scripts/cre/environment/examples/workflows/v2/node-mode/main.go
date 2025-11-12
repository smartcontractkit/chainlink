//go:build wasip1

package main

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"

	"github.com/google/uuid"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type None struct{}

func main() {
	wasm.NewRunner(func(configBytes []byte) (None, error) {
		return None{}, nil
	}).Run(RunSimpleCronWorkflow)
}

func RunSimpleCronWorkflow(_ None, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[None], error) {
	workflows := cre.Workflow[None]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onTrigger,
		),
	}
	return workflows, nil
}

func onTrigger(cfg None, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	runtime.Logger().Info("Triggered fetch of value")

	err := testIdenticalConsensus(cfg, runtime)
	if err != nil {
		return "testIdenticalConsensus failed", err
	}

	err = testIdenticalConsensusFailure(cfg, runtime)
	if err != nil {
		return "testIdenticalConsensusFailure failed", err
	}

	err = testIdenticalConsensusFailureWithDefault(cfg, runtime)
	if err != nil {
		return "testIdenticalConsensusFailureWithDefault failed", err
	}

	err = testMedianConsensus(cfg, runtime)
	if err != nil {
		return "testMedianConsensus failed", err
	}

	err = testMedianConsensusWithErrors(cfg, runtime)
	if err != nil {
		return "testMedianConsensusWithErrors failed", err
	}

	err = testMedianConsensusWithErrorsAndDefault(cfg, runtime)
	if err != nil {
		return "testMedianConsensusWithErrors failed", err
	}

	runtime.Logger().Info("Successfully passed all consensus tests")

	return "success", nil
}

func testMedianConsensus(cfg None, runtime cre.Runtime) error {
	runtime.Logger().Info("Starting testMedianConsensus")
	mathPromise := cre.RunInNodeMode(cfg, runtime, fetchData, cre.ConsensusMedianAggregation[int]())
	offchainValue, err := mathPromise.Await()
	if err != nil {
		runtime.Logger().Warn("Median consensus error", "error", err)
		return err
	}
	runtime.Logger().Info("Successfully fetched offchain value and reached consensus", "result", offchainValue)
	return nil
}

func testMedianConsensusWithErrors(cfg None, runtime cre.Runtime) error {
	runtime.Logger().Info("Starting testMedianConsensusWithErrors")
	mathPromise := cre.RunInNodeMode(cfg, runtime, func(config None, nodeRuntime cre.NodeRuntime) (int, error) {
		return 0, errors.New("simulated error")
	}, cre.ConsensusMedianAggregation[int]())
	offchainValue, err := mathPromise.Await()
	if err == nil {
		runtime.Logger().Warn("expected median consensus error", "error", err)
		return err
	} else {
		expectedInMessage := "simulated error"
		if !strings.Contains(err.Error(), expectedInMessage) {
			runtime.Logger().Warn("expected median consensus error", "error", err)
			return fmt.Errorf("expected error to contain '%s', got '%s'", expectedInMessage, err.Error())
		}
	}

	runtime.Logger().Info("Successfully tested consensus errors", "result", offchainValue)
	return nil
}

func testMedianConsensusWithErrorsAndDefault(cfg None, runtime cre.Runtime) error {
	runtime.Logger().Info("Starting testMedianConsensusWithErrorsAndDefault")
	mathPromise := cre.RunInNodeMode(cfg, runtime, func(config None, nodeRuntime cre.NodeRuntime) (int, error) {
		return 0, errors.New("simulated error")
	}, cre.ConsensusMedianAggregation[int]().WithDefault(42))
	offchainValue, err := mathPromise.Await()
	if err != nil {
		runtime.Logger().Warn("Median consensus with errors and default error error", "error", err)
		return err
	}
	runtime.Logger().Info("Successfully used default when errors", "result", offchainValue)
	return nil
}

func testIdenticalConsensus(cfg None, runtime cre.Runtime) error {
	runtime.Logger().Info("Starting testIdenticalConsensus")
	sameValueStr := "samevalue"
	byteSlicePromise := cre.RunInNodeMode(cfg, runtime, func(config None, nodeRuntime cre.NodeRuntime) ([]byte, error) {
		return []byte(sameValueStr), nil
	}, cre.ConsensusIdenticalAggregation[[]byte]())
	resultBytes, err := byteSlicePromise.Await()
	if err != nil {
		runtime.Logger().Warn("Consensus error on identical consensus test", "error", err)
		return err
	}

	if string(resultBytes) == sameValueStr {
		runtime.Logger().Info("Identical consensus on a []byte succeeded", "result", resultBytes)
	} else {
		msg := fmt.Sprintf("Identical consensus on a []byte failed, expected '%s', got '%s'", sameValueStr, resultBytes)
		runtime.Logger().Error(msg)
		return errors.New(msg)
	}

	return nil
}

// testIdenticalConsensus tests identical consensus on a []byte value that is different on each node.
func testIdenticalConsensusFailure(cfg None, runtime cre.Runtime) error {
	runtime.Logger().Info("Starting testIdenticalConsensusFailure")
	sameValueStr := uuid.New().String()
	byteSlicePromise := cre.RunInNodeMode(cfg, runtime, func(config None, nodeRuntime cre.NodeRuntime) ([]byte, error) {
		return []byte(sameValueStr), nil
	}, cre.ConsensusIdenticalAggregation[[]byte]())
	resultBytes, err := byteSlicePromise.Await()
	if err == nil {
		runtime.Logger().Warn("Expected consensus to fail but it succeeded", "result", resultBytes)
		return errors.New("Expected consensus to fail but it succeeded")
	}

	errString := err.Error()
	if !strings.Contains(errString, "no values met f+1 threshold") {
		runtime.Logger().Warn("Unexpected error message", "error", errString)
		return fmt.Errorf("unexpected error message: %s", errString)
	}

	return nil
}

func testIdenticalConsensusFailureWithDefault(cfg None, runtime cre.Runtime) error {
	runtime.Logger().Info("Starting testIdenticalConsensusFailureWithDefault")
	sameValueStr := uuid.New().String()
	defaultStr := "adefault"
	byteSlicePromise := cre.RunInNodeMode(cfg, runtime, func(config None, nodeRuntime cre.NodeRuntime) ([]byte, error) {
		return []byte(sameValueStr), nil
	}, cre.ConsensusIdenticalAggregation[[]byte]().WithDefault([]byte(defaultStr)))
	resultBytes, err := byteSlicePromise.Await()
	if err != nil {
		runtime.Logger().Warn("Consensus error on identical consensus failoure with default test", "error", err)
		return err
	}

	if string(resultBytes) == defaultStr {
		runtime.Logger().Info("Identical consensus on a []byte with default succeeded", "result", resultBytes)
	} else {
		msg := fmt.Sprintf("Identical consensus on a []byte with default failed, expected '%s', got '%s'", defaultStr, resultBytes)
		runtime.Logger().Error(msg)
		return errors.New(msg)
	}

	return nil
}

func fetchData(cfg None, nodeRuntime cre.NodeRuntime) (int, error) {

	randomValue := rand.Intn(10000)
	nodeRuntime.Logger().Info("Generate random value", "randomValue", randomValue)

	// Generate a random int64
	return randomValue, nil
}
