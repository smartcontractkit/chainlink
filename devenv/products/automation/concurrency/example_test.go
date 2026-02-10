package concurrency_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/devenv/products/automation/concurrency"
)

type client struct{}

func (c *client) getConcurrency() int {
	return 1
}

func (c *client) deployContractConfigurableFromKey(_ int, _ contractConfiguration) (ContractIntstance, error) {
	return ContractIntstance{}, nil
}

func (c *client) deployContractFromKey(_ int) (ContractIntstance, error) {
	return ContractIntstance{}, nil
}

type ContractIntstance struct{}

type contractConfiguration struct{}

type contractResult struct {
	instance ContractIntstance
}

func (k contractResult) GetResult() ContractIntstance {
	return k.instance
}

func TestExampleContractsWithConfiguration(t *testing.T) {
	instances, err := DeployContractsWithConfiguration(&client{}, []contractConfiguration{{}, {}})
	require.NoError(t, err, "failed to deploy contract instances")
	require.Len(t, instances, 2, "expected 2 contract instances")
}

// DeployContractsWithConfiguration shows a very simplified method that deploys concurrently contract instances with given configurations
func DeployContractsWithConfiguration(client *client, contractConfigs []contractConfiguration) ([]ContractIntstance, error) {
	l := framework.L

	executor := concurrency.NewConcurrentExecutor[ContractIntstance, contractResult, contractConfiguration](l)

	var deployContractFn = func(channel chan contractResult, errorCh chan error, executorNum int, payload contractConfiguration) {
		keyNum := executorNum + 1 // key 0 is the root key

		instance, err := client.deployContractConfigurableFromKey(keyNum, payload)
		if err != nil {
			errorCh <- err
			return
		}

		channel <- contractResult{instance: instance}
	}

	results, err := executor.Execute(client.getConcurrency(), contractConfigs, deployContractFn)
	if err != nil {
		return []ContractIntstance{}, err
	}

	if len(results) != len(contractConfigs) {
		return []ContractIntstance{}, fmt.Errorf("expected %v results, got %v", len(contractConfigs), len(results))
	}

	return results, nil
}

func TestExampleContractsWithoutConfiguration(t *testing.T) {
	instances, err := DeployIdenticalContracts(&client{}, 2)
	require.NoError(t, err, "failed to deploy contract instances")
	require.Len(t, instances, 2, "expected 2 contract instances")
}

// DeployIdenticalContracts shows a very simplified method that deploys concurrently identical contract instances
// which require no configuration, just need to be executed N amount of times
func DeployIdenticalContracts(client *client, numberOfContracts int) ([]ContractIntstance, error) {
	l := framework.L

	executor := concurrency.NewConcurrentExecutor[ContractIntstance, contractResult, concurrency.NoTaskType](l)

	var deployContractFn = func(channel chan contractResult, errorCh chan error, executorNum int) {
		keyNum := executorNum + 1 // key 0 is the root key

		instance, err := client.deployContractFromKey(keyNum)
		if err != nil {
			errorCh <- err
			return
		}

		channel <- contractResult{instance: instance}
	}

	results, err := executor.ExecuteSimple(client.getConcurrency(), numberOfContracts, deployContractFn)
	if err != nil {
		return []ContractIntstance{}, err
	}

	if len(results) != numberOfContracts {
		return []ContractIntstance{}, fmt.Errorf("expected %v results, got %v", numberOfContracts, len(results))
	}

	return results, nil
}
