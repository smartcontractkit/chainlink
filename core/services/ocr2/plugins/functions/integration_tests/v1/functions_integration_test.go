package functions_test

import (
	"testing"
	"time"

	functionsConfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	utils "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/integration_tests/v1/internal"
)

var (
	// a batch of 8 max-length results uses around 2M gas (assuming 70k gas per client callback - see FunctionsClientExample.sol)
	// TODO: revisit gas limit vs batch sizing once necessary contract changes are made
	nOracleNodes      = 4
	nClients          = 50
	requestLenBytes   = 1_000
	maxGas            = 1_700_000
	maxTotalReportGas = 560_000
	batchSize         = 8
)

func TestIntegration_Functions_MultipleRequests_Success(t *testing.T) {
	// simulated chain with all contracts
	owner, b, ticker, coordinatorContractAddress, coordinatorContract, clientContracts, routerAddress, routerContract, linkToken, allowListContractAddress, allowListContract := utils.StartNewChainWithContracts(t, nClients)
	defer ticker.Stop()

	utils.SetupRouterRoutes(t, owner, routerContract, coordinatorContractAddress, allowListContractAddress)
	b.Commit()

	_, _, oracleIdentities := utils.CreateFunctionsNodes(t, owner, b, routerAddress, coordinatorContractAddress, 39989, nOracleNodes, maxGas, nil, nil)

	pluginConfig := functionsConfig.ReportingPluginConfig{
		MaxQueryLengthBytes:       10_000,
		MaxObservationLengthBytes: 10_000,
		MaxReportLengthBytes:      15_000,
		MaxRequestBatchSize:       uint32(batchSize),
		MaxReportTotalCallbackGas: uint32(maxTotalReportGas),
		DefaultAggregationMethod:  functionsConfig.AggregationMethod_AGGREGATION_MODE,
		UniqueReports:             true,
	}

	// config for oracle contract
	utils.SetOracleConfig(t, owner, coordinatorContract, oracleIdentities, batchSize, &pluginConfig)
	utils.CommitWithFinality(b)

	// validate that all client contracts got correct responses to their requests
	utils.ClientTestRequests(t, owner, b, linkToken, routerAddress, routerContract, allowListContract, clientContracts, requestLenBytes, utils.DefaultSecretsBytes, 1*time.Minute)
}
