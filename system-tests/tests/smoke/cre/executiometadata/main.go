//go:build wasip1

package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/executiometadata/config"
)

func main() {
	wasm.NewRunner(parseConfig).Run(initWorkflow)
}

func parseConfig(raw []byte) (config.Config, error) {
	cfg := config.Config{Schedule: config.DefaultSchedule}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return cfg, fmt.Errorf("parse workflow configuration: %w", err)
		}
	}
	cfg.Schedule = strings.TrimSpace(cfg.Schedule)
	if cfg.Schedule == "" {
		return cfg, errors.New("cron schedule must not be empty")
	}
	return cfg, nil
}

func initWorkflow(cfg config.Config, logger *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	logger.Info("Execution metadata workflow configured", "schedule", cfg.Schedule)
	return cre.Workflow[config.Config]{
		cre.Handler(cron.Trigger(&cron.Config{Schedule: cfg.Schedule}), onTrigger),
	}, nil
}

func onTrigger(cfg config.Config, runtime cre.Runtime, _ *cron.Payload) (cre.WorkflowExecutionMetadata, error) {
	executionMetadata := cre.ExecutionMetadata(runtime)
	logger := runtime.Logger().With(
		"workflowID", executionMetadata.WorkflowID,
		"executionID", executionMetadata.ExecutionID,
		"workflowOwner", executionMetadata.WorkflowOwner,
		"donID", executionMetadata.DonID,
		"donF", executionMetadata.DonF,
	)

	if err := validateExecutionMetadata(executionMetadata); err != nil {
		logger.Error("initial workflow execution metadata validation failed", "error", err)
		return executionMetadata, err
	}

	consensusMetadata, err := cre.RunInNodeMode(cfg, runtime, func(_ config.Config, node cre.NodeRuntime) (cre.WorkflowExecutionMetadata, error) {
		nodeMetadata := cre.ExecutionMetadata(node)
		if nodeMetadata != executionMetadata {
			return nodeMetadata, errors.New("node mode workflow execution metadata does not match expected execution metadata")
		}
		node.Logger().Info("Node mode execution metadata matches expected execution metadata",
			"workflowID", nodeMetadata.WorkflowID,
			"executionID", nodeMetadata.ExecutionID,
		)
		return nodeMetadata, nil
	}, cre.ConsensusIdenticalAggregation[cre.WorkflowExecutionMetadata]()).Await()
	if err != nil {
		err = fmt.Errorf("confirm node mode execution metadata through consensus: %w", err)
		logger.Error("node mode execution metadata validation failed", "error", err)
		return executionMetadata, err
	}
	if consensusMetadata != executionMetadata {
		err := errors.New("consensus metadata does not match execution metadata")
		logger.Error("Execution metadata consensus validation failed", "error", err)
		return executionMetadata, err
	}

	logger.Info("Workflow execution metadata fetched successfully",
		"event", "WORKFLOW_EXECUTION_METADATA",
	)
	return executionMetadata, nil
}

func validateExecutionMetadata(executionMetadata cre.WorkflowExecutionMetadata) error {
	for _, field := range []struct {
		name  string
		value string
		size  int
	}{
		{name: "workflow ID", value: executionMetadata.WorkflowID, size: 32},
		{name: "execution ID", value: executionMetadata.ExecutionID, size: 32},
		{name: "workflow owner", value: executionMetadata.WorkflowOwner, size: 20},
	} {
		decoded, err := hex.DecodeString(field.value)
		if err != nil {
			return fmt.Errorf("invalid field in execution metadata '%s': '%w'", field.name, err)
		}
		if len(decoded) != field.size || field.value != strings.ToLower(field.value) {
			return fmt.Errorf("%s must be %d-byte lowercase hexadecimal without 0x", field.name, field.size)
		}
	}
	return nil
}
