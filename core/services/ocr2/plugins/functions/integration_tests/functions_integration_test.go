package functions_test

import (
	"encoding/hex"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/onsi/gomega"
	functionsConfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	utils "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/integration_tests/internal"
	"github.com/test-go/testify/require"
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

	_, _, oracleIdentities := utils.CreateFunctionsNodes(t, owner, b, oracleContractAddress, 39999, nOracleNodes, maxGas, nil, nil)

	// config for registry contract
	utils.SetRegistryConfig(t, owner, registryContract, oracleContractAddress)

	pluginConfig := functionsConfig.ReportingPluginConfig{
		MaxQueryLengthBytes:       10_000,
		MaxObservationLengthBytes: 10_000,
		MaxReportLengthBytes:      10_000,
		MaxRequestBatchSize:       uint32(batchSize),
		DefaultAggregationMethod:  functionsConfig.AggregationMethod_AGGREGATION_MODE,
		UniqueReports:             true,
	}

	// config for oracle contract
	utils.SetOracleConfig(t, owner, oracleContract, oracleIdentities, batchSize, &pluginConfig)
	utils.CommitWithFinality(b)

	// // set up subscription
	// subscriptionId := utils.CreateAndFundSubscriptions(t, owner, linkToken, registryAddress, registryContract, clientContracts)

	// // send requests
	// requestSources := make([][]byte, nClients)
	// rnd := rand.New(rand.NewSource(666))
	// for i := 0; i < nClients; i++ {
	// 	requestSources[i] = make([]byte, requestLenBytes)
	// 	for j := 0; j < requestLenBytes; j++ {
	// 		requestSources[i][j] = byte(rnd.Uint32() % 256)
	// 	}
	// 	_, err := clientContracts[i].Contract.SendRequest(
	// 		owner,
	// 		hex.EncodeToString(requestSources[i]),
	// 		utils.DefaultSecretsUrlsBytes,
	// 		[]string{utils.DefaultArg1, utils.DefaultArg2},
	// 		subscriptionId)
	// 	require.NoError(t, err)
	// }
	// utils.CommitWithFinality(b)

	// // validate that all client contracts got correct responses to their requests
	// var wg sync.WaitGroup
	// for i := 0; i < nClients; i++ {
	// 	ic := i
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		gomega.NewGomegaWithT(t).Eventually(func() [32]byte {
	// 			answer, err := clientContracts[ic].Contract.LastResponse(nil)
	// 			require.NoError(t, err)
	// 			return answer
	// 		}, 3*time.Minute, 1*time.Second).Should(gomega.Equal(utils.GetExpectedResponse(requestSources[ic])))
	// 	}()
	// }
	// wg.Wait()

	// validate that all client contracts got correct responses to their requests
	utils.ClientTestRequests(t, owner, b, linkToken, registryAddress, registryContract, clientContracts, requestLenBytes, time.Duration(1)*time.Minute)
}

func TestIntegration_Functions_MultipleRequests_ThresholdDecryptionSuccess(t *testing.T) {
	// simulated chain with all contracts
	owner, b, ticker, oracleContractAddress, oracleContract, clientContracts, registryAddress, registryContract, linkToken := utils.StartNewChainWithContracts(t, nClients)
	defer ticker.Stop()

	_, _, oracleIdentities := utils.CreateFunctionsNodes(t, owner, b, oracleContractAddress, 49999, nOracleNodes, maxGas, utils.ExportedOcr2Keystores, utils.MockThresholdKeyShares)

	// config for registry contract
	utils.SetRegistryConfig(t, owner, registryContract, oracleContractAddress)

	pluginConfig := functionsConfig.ReportingPluginConfig{
		MaxQueryLengthBytes:       10_000,
		MaxObservationLengthBytes: 10_000,
		MaxReportLengthBytes:      10_000,
		MaxRequestBatchSize:       uint32(batchSize),
		DefaultAggregationMethod:  functionsConfig.AggregationMethod_AGGREGATION_MODE,
		UniqueReports:             true,
		ThresholdPluginConfig: &functionsConfig.ThresholdReportingPluginConfig{
			MaxQueryLengthBytes:       10_000,
			MaxObservationLengthBytes: 10_000,
			MaxReportLengthBytes:      10_000,
			RequestCountLimit:         100,
			RequestTotalBytesLimit:    1_000,
			RequireLocalRequestCheck:  true,
		},
	}

	// config for oracle contract
	utils.SetOracleConfig(t, owner, oracleContract, oracleIdentities, batchSize, &pluginConfig)
	utils.CommitWithFinality(b)

	// set up subscription
	subscriptionId := utils.CreateAndFundSubscriptions(t, owner, linkToken, registryAddress, registryContract, clientContracts)

	// send requests
	requestSources := make([][]byte, nClients)
	rnd := rand.New(rand.NewSource(666))
	for i := 0; i < nClients; i++ {
		requestSources[i] = make([]byte, requestLenBytes)
		for j := 0; j < requestLenBytes; j++ {
			requestSources[i][j] = byte(rnd.Uint32() % 256)
		}
		_, err := clientContracts[i].Contract.SendRequest(
			owner,
			hex.EncodeToString(requestSources[i]),
			utils.DefaultSecretsUrlsBytes,
			[]string{utils.DefaultArg1, utils.DefaultArg2},
			subscriptionId)
		require.NoError(t, err)
	}
	utils.CommitWithFinality(b)

	// validate that all client contracts got correct responses to their requests
	var wg sync.WaitGroup
	for i := 0; i < nClients; i++ {
		ic := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			gomega.NewGomegaWithT(t).Eventually(func() [32]byte {
				answer, err := clientContracts[ic].Contract.LastResponse(nil)
				require.NoError(t, err)
				return answer
			}, 3*time.Minute, 1*time.Second).Should(gomega.Equal(utils.GetExpectedResponse(requestSources[ic])))
		}()
	}
	wg.Wait()

	// validate that all client contracts got correct responses to their requests
	// utils.ClientTestRequests(t, owner, b, linkToken, registryAddress, registryContract, clientContracts, requestLenBytes, time.Duration(3)*time.Minute)
}
