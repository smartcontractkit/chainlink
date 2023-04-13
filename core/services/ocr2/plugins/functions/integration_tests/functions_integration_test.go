package functions_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/libocr/commontypes"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	utils "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/integration_tests/internal"
)

func TestIntegration_Functions_MultipleRequests_Success(t *testing.T) {
	// a batch of 8 max-length results uses around 1M gas (assuming 70k gas per client callback - see FunctionsClientExample.sol)
	nOracleNodes := 4
	nClients := 50
	requestLenBytes := 1000
	maxGas := 1_300_000
	batchSize := 8

	// simulated chain with all contracts
	owner, b, ticker, oracleContractAddress, oracleContract, clientContracts, registryAddress, registryContract, linkToken := utils.StartNewChainWithContracts(t, nClients)
	defer ticker.Stop()

	// bootstrap node and job
	bootstrapNodePort := uint16(39999)
	bootstrapNode := utils.StartNewNode(t, owner, bootstrapNodePort, "bootstrap", b, uint32(maxGas), nil)
	utils.AddBootstrapJob(t, bootstrapNode.App, oracleContractAddress)

	// oracle nodes with jobs, bridges and mock EAs
	var jobIds []int32
	var oracles []confighelper2.OracleIdentityExtra
	var apps []*cltest.TestApplication
	for i := 0; i < nOracleNodes; i++ {
		oracleNode := utils.StartNewNode(t, owner, bootstrapNodePort+1+uint16(i), fmt.Sprintf("oracle%d", i), b, uint32(maxGas), []commontypes.BootstrapperLocator{
			{PeerID: bootstrapNode.PeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
		})
		apps = append(apps, oracleNode.App)
		oracles = append(oracles, oracleNode.OracleIdentity)

		ea := utils.StartNewMockEA(t)
		defer ea.Close()

		ocrJob := utils.AddOCR2Job(t, apps[i], oracleContractAddress, oracleNode.Keybundle.ID(), oracleNode.Transmitter, ea.URL)
		jobIds = append(jobIds, ocrJob.ID)
	}

	// config for registry contract
	utils.SetRegistryConfig(t, owner, registryContract, oracleContractAddress)

	// config for oracle contract
	utils.SetOracleConfig(t, owner, oracleContract, oracles, batchSize)
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
		_, err := clientContracts[i].Contract.SendRequest(owner, hex.EncodeToString(requestSources[i]), []byte{}, []string{}, subscriptionId)
		require.NoError(t, err)
	}
	utils.CommitWithFinality(b)

	// validate that all pipeline jobs completed as many runs as sent requests
	const tasksPerRun = 3
	var wg sync.WaitGroup
	for i := 0; i < nOracleNodes; i++ {
		ic := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			cltest.WaitForPipelineComplete(t, ic, jobIds[ic], nClients, tasksPerRun, apps[ic].JobORM(), 1*time.Minute, 1*time.Second)
		}()
	}
	wg.Wait()

	// validate that all client contracts got correct responses to their requests
	for i := 0; i < nClients; i++ {
		ic := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			gomega.NewGomegaWithT(t).Eventually(func() [32]byte {
				answer, err := clientContracts[ic].Contract.LastResponse(nil)
				require.NoError(t, err)
				return answer
			}, 1*time.Minute, 1*time.Second).Should(gomega.Equal(utils.GetExpectedResponse(requestSources[ic])))
		}()
	}
	wg.Wait()
}
