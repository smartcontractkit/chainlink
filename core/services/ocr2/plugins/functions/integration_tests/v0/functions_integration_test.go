package functions_test

import (
	"testing"
	"time"

	functionsConfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	utils "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/integration_tests/v0/internal"
)

var (
	// a batch of 8 max-length results uses around 1M gas (assuming 70k gas per client callback - see FunctionsClientExample.sol)
	nOracleNodes    = 4
	nClients        = 50
	requestLenBytes = 1000
	maxGas          = 1_300_000
	batchSize       = 8
)

func TestIntegration_Functions_MultipleRequests_Success(t *testing.T) {
	// simulated chain with all contracts
	owner, b, ticker, oracleContractAddress, oracleContract, clientContracts, registryAddress, registryContract, linkToken := utils.StartNewChainWithContracts(t, nClients)
	defer ticker.Stop()

	_, _, oracleIdentities := utils.CreateFunctionsNodes(t, owner, b, oracleContractAddress, nOracleNodes, maxGas, nil, nil)

	// config for registry contract
	utils.SetRegistryConfig(t, owner, registryContract, oracleContractAddress)

	pluginConfig := functionsConfig.ReportingPluginConfig{
		MaxQueryLengthBytes:       10_000,
		MaxObservationLengthBytes: 10_000,
		MaxReportLengthBytes:      15_000,
		MaxRequestBatchSize:       uint32(batchSize),
		DefaultAggregationMethod:  functionsConfig.AggregationMethod_AGGREGATION_MODE,
		UniqueReports:             true,
	}

	// config for oracle contract
	utils.SetOracleConfig(t, owner, oracleContract, oracleIdentities, batchSize, &pluginConfig)
	utils.CommitWithFinality(b)

	// validate that all client contracts got correct responses to their requests
	utils.ClientTestRequests(t, owner, b, linkToken, registryAddress, registryContract, clientContracts, requestLenBytes, utils.DefaultSecretsBytes, 1*time.Minute)
}

func TestIntegration_Functions_MultipleRequests_ThresholdDecryptionSuccess(t *testing.T) {
	// simulated chain with all contracts
	owner, b, ticker, oracleContractAddress, oracleContract, clientContracts, registryAddress, registryContract, linkToken := utils.StartNewChainWithContracts(t, nClients)
	defer ticker.Stop()

	_, _, oracleIdentities := utils.CreateFunctionsNodes(t, owner, b, oracleContractAddress, nOracleNodes, maxGas, utils.ExportedOcr2Keystores, utils.MockThresholdKeyShares)

	// config for registry contract
	utils.SetRegistryConfig(t, owner, registryContract, oracleContractAddress)

	pluginConfig := functionsConfig.ReportingPluginConfig{
		MaxQueryLengthBytes:       10_000,
		MaxObservationLengthBytes: 10_000,
		MaxReportLengthBytes:      15_000,
		MaxRequestBatchSize:       uint32(batchSize),
		DefaultAggregationMethod:  functionsConfig.AggregationMethod_AGGREGATION_MODE,
		UniqueReports:             true,
		ThresholdPluginConfig: &functionsConfig.ThresholdReportingPluginConfig{
			// approximately 750 bytes per test ciphertext + overhead
			MaxQueryLengthBytes:       70_000,
			MaxObservationLengthBytes: 70_000,
			MaxReportLengthBytes:      70_000,
			RequestCountLimit:         50,
			RequestTotalBytesLimit:    50_000,
			RequireLocalRequestCheck:  true,
			K:                         2,
		},
	}

	// config for oracle contract
	utils.SetOracleConfig(t, owner, oracleContract, oracleIdentities, batchSize, &pluginConfig)
	utils.CommitWithFinality(b)

	// validate that all client contracts got correct responses to their requests
	utils.ClientTestRequests(t, owner, b, linkToken, registryAddress, registryContract, clientContracts, requestLenBytes, utils.DefaultSecretsUrlsBytes, 1*time.Minute)
}
