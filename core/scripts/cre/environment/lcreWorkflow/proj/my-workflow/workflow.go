package main

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"

	"proj/contracts/evm/src/generated/capabilities_registry"
)

// randomNumberURL returns a single random integer in [1,100] as plain text.
const randomNumberURL = "https://www.random.org/integers/?num=1&min=1&max=100&col=1&base=10&format=plain&rnd=new"

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

	version, err := readRegistryTypeAndVersion(config, runtime)
	if err != nil {
		return nil, err
	}
	logger.Info("CapabilitiesRegistry typeAndVersion",
		"chainSelector", config.ChainSelector,
		"address", config.RegistryAddress,
		"typeAndVersion", version,
	)

	randomMedian, err := http.SendRequest(
		config,
		runtime,
		&http.Client{},
		fetchRandomNumber,
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

// fetchRandomNumber runs per-node and returns this node's random int in [1,100]
// fetched from random.org over the HTTP capability.
func fetchRandomNumber(config *Config, logger *slog.Logger, sendRequester *http.SendRequester) (int64, error) {
	resp, err := sendRequester.SendRequest(&http.Request{
		Url:     randomNumberURL,
		Method:  "GET",
		Timeout: &durationpb.Duration{Seconds: 10},
	}).Await()
	if err != nil {
		return 0, fmt.Errorf("failed to fetch random number: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("unexpected status code %d fetching random number: %s", resp.StatusCode, string(resp.Body))
	}

	n, err := strconv.ParseInt(strings.TrimSpace(string(resp.Body)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse random number %q: %w", string(resp.Body), err)
	}

	logger.Info("Fetched random number (per node)", "value", n)
	return n, nil
}
