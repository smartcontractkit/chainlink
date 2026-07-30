package main

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"

	"proj/contracts/evm/src/generated/capabilities_registry"
)

type ExecutionResult struct {
	RegistryTypeAndVersion string
	RandomNumber           int64
}

// Workflow configuration loaded from the config.json file.
type Config struct {
	// ChainSelector identifies the EVM chain to read from.
	// Local CRE anvil chain 1337 => 3379446385462418246.
	ChainSelector uint64 `json:"chainSelector"`
	// RegistryAddress is the CapabilitiesRegistry contract address on that chain.
	RegistryAddress string `json:"registryAddress"`
}

// Workflow implementation with a list of capability triggers
func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	// Create the trigger
	cronTrigger := cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}) // Fires every 30 seconds

	// Register a handler with the trigger and a callback function
	return cre.Workflow[*Config]{
		cre.Handler(cronTrigger, onCronTrigger),
	}, nil
}

func onCronTrigger(config *Config, runtime cre.Runtime, trigger *cron.Payload) (*ExecutionResult, error) {
	logger := runtime.Logger()

	// 1) EVM read: CapabilitiesRegistry.typeAndVersion()
	version, err := readRegistryTypeAndVersion(config, runtime)
	if err != nil {
		return nil, err
	}
	logger.Info("CapabilitiesRegistry typeAndVersion",
		"chainSelector", config.ChainSelector,
		"address", config.RegistryAddress,
		"typeAndVersion", version,
	)

	// 2) Random number with median consensus across the DON nodes.
	// Each node draws from its own NodeRuntime random source; the DON agrees on the median.
	randomMedian, err := cre.RunInNodeMode(
		config,
		runtime,
		nextRandomNumber,
		cre.ConsensusMedianAggregation[int64](),
	).Await()
	if err != nil {
		return nil, fmt.Errorf("failed to get median random number: %w", err)
	}
	logger.Info("Random number (median consensus)", "median", randomMedian)

	return &ExecutionResult{
		RegistryTypeAndVersion: version,
		RandomNumber:           randomMedian,
	}, nil
}

func readRegistryTypeAndVersion(config *Config, runtime cre.Runtime) (string, error) {
	evmClient := &evm.Client{ChainSelector: config.ChainSelector}
	registry, err := capabilities_registry.NewCapabilitiesRegistry(
		evmClient,
		common.HexToAddress(config.RegistryAddress),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create CapabilitiesRegistry binding: %w", err)
	}

	// Passing nil for blockNumber reads at the latest finalized block.
	version, err := registry.TypeAndVersion(runtime, nil).Await()
	if err != nil {
		return "", fmt.Errorf("failed to call typeAndVersion: %w", err)
	}
	return version, nil
}

// nextRandomNumber runs per-node and returns this node's next random int in [1,100]
// drawn from the NodeRuntime's random source.
func nextRandomNumber(config *Config, nodeRuntime cre.NodeRuntime) (int64, error) {
	rng, err := nodeRuntime.Rand()
	if err != nil {
		return 0, fmt.Errorf("failed to get node random source: %w", err)
	}

	n := int64(rng.Intn(100) + 1)
	nodeRuntime.Logger().Info("Generated random number (per node)", "value", n)
	return n, nil
}
