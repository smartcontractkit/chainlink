package functions_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	functionsConfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	utils "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/integration_tests/v1/internal"
)

var (
	// a batch of 8 max-length results uses around 2M gas (assuming 70k gas per client callback - see FunctionsClientExample.sol)
	nOracleNodes      = 4
	nClients          = 50
	requestLenBytes   = 1_000
	maxGas            = 1_700_000
	maxTotalReportGas = 560_000
	batchSize         = 8
)

func TestIntegration_Functions_MultipleV1Requests_Success(t *testing.T) {
	// simulated chain with all contracts
	owner, b, commit, stop, active, proposed, clientContracts, routerAddress, routerContract, linkToken, allowListContractAddress, allowListContract := utils.StartNewChainWithContracts(t, nClients)
	defer stop()

	utils.SetupRouterRoutes(t, b, owner, routerContract, active.Address, proposed.Address, allowListContractAddress)

	_, _, oracleIdentities := utils.CreateFunctionsNodes(t, owner, b, routerAddress, nOracleNodes, maxGas, nil, nil)

	pluginConfig := functionsConfig.ReportingPluginConfig{
		MaxQueryLengthBytes:       10_000,
		MaxObservationLengthBytes: 15_000,
		MaxReportLengthBytes:      15_000,
		MaxRequestBatchSize:       uint32(batchSize),
		MaxReportTotalCallbackGas: uint32(maxTotalReportGas),
		DefaultAggregationMethod:  functionsConfig.AggregationMethod_AGGREGATION_MODE,
		UniqueReports:             true,
	}

	// config for oracle contract
	utils.SetOracleConfig(t, b, owner, active.Contract, oracleIdentities, batchSize, &pluginConfig)

	subscriptionId := utils.CreateAndFundSubscriptions(t, b, owner, linkToken, routerAddress, routerContract, clientContracts, allowListContract)
	commit()
	utils.ClientTestRequests(t, owner, b, clientContracts, requestLenBytes, nil, subscriptionId, 1*time.Minute)
}

func TestIntegration_Functions_MultipleV1Requests_ThresholdDecryptionSuccess(t *testing.T) {
	// simulated chain with all contracts
	owner, b, commit, stop, active, proposed, clientContracts, routerAddress, routerContract, linkToken, allowListContractAddress, allowListContract := utils.StartNewChainWithContracts(t, nClients)
	defer stop()

	utils.SetupRouterRoutes(t, b, owner, routerContract, active.Address, proposed.Address, allowListContractAddress)

	_, _, oracleIdentities := utils.CreateFunctionsNodes(t, owner, b, routerAddress, nOracleNodes, maxGas, utils.ExportedOcr2Keystores, utils.MockThresholdKeyShares)

	pluginConfig := functionsConfig.ReportingPluginConfig{
		MaxQueryLengthBytes:       10_000,
		MaxObservationLengthBytes: 15_000,
		MaxReportLengthBytes:      15_000,
		MaxRequestBatchSize:       uint32(batchSize),
		MaxReportTotalCallbackGas: uint32(maxTotalReportGas),
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
	utils.SetOracleConfig(t, b, owner, active.Contract, oracleIdentities, batchSize, &pluginConfig)

	subscriptionId := utils.CreateAndFundSubscriptions(t, b, owner, linkToken, routerAddress, routerContract, clientContracts, allowListContract)
	commit()
	utils.ClientTestRequests(t, owner, b, clientContracts, requestLenBytes, utils.DefaultSecretsUrlsBytes, subscriptionId, 1*time.Minute)
}

func TestIntegration_Functions_MultipleV1Requests_WithUpgrade(t *testing.T) {
	// simulated chain with all contracts
	owner, b, commit, stop, active, proposed, clientContracts, routerAddress, routerContract, linkToken, allowListContractAddress, allowListContract := utils.StartNewChainWithContracts(t, nClients)
	defer stop()

	utils.SetupRouterRoutes(t, b, owner, routerContract, active.Address, proposed.Address, allowListContractAddress)

	_, _, oracleIdentities := utils.CreateFunctionsNodes(t, owner, b, routerAddress, nOracleNodes, maxGas, utils.ExportedOcr2Keystores, utils.MockThresholdKeyShares)

	pluginConfig := functionsConfig.ReportingPluginConfig{
		MaxQueryLengthBytes:       10_000,
		MaxObservationLengthBytes: 15_000,
		MaxReportLengthBytes:      15_000,
		MaxRequestBatchSize:       uint32(batchSize),
		MaxReportTotalCallbackGas: uint32(maxTotalReportGas),
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

	// set config for both coordinators
	utils.SetOracleConfig(t, b, owner, active.Contract, oracleIdentities, batchSize, &pluginConfig)
	utils.SetOracleConfig(t, b, owner, proposed.Contract, oracleIdentities, batchSize, &pluginConfig)

	subscriptionId := utils.CreateAndFundSubscriptions(t, b, owner, linkToken, routerAddress, routerContract, clientContracts, allowListContract)
	utils.ClientTestRequests(t, owner, b, clientContracts, requestLenBytes, utils.DefaultSecretsUrlsBytes, subscriptionId, 1*time.Minute)

	// upgrade and send requests again
	_, err := routerContract.UpdateContracts(owner)
	require.NoError(t, err)
	commit()
	utils.ClientTestRequests(t, owner, b, clientContracts, requestLenBytes, utils.DefaultSecretsUrlsBytes, subscriptionId, 1*time.Minute)
}
