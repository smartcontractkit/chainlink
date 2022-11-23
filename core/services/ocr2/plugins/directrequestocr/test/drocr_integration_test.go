package directrequestocr_test

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

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	utils "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr/test/internal"
)

func TestIntegration_OCR2DR_MultipleRequests_Success(t *testing.T) {
	nNodes := 4
	nClients := 20

	// simulated chain with all contracts
	owner, b, ticker, oracleContractAddress, oracleContract, clientContracts, registryAddress, registryContract, linkToken := utils.StartNewChainWithContracts(t, nClients)
	defer ticker.Stop()

	// bootstrap node and job
	bootstrapNodePort := uint16(29999)
	bootstrapNode := utils.StartNewNode(t, owner, bootstrapNodePort, "bootstrap", b, nil)
	utils.AddBootstrapJob(t, bootstrapNode.App, oracleContractAddress)

	// oracle nodes with jobs, bridges and mock EAs
	var jobIds []int32
	var oracles []confighelper2.OracleIdentityExtra
	var apps []*cltest.TestApplication
	for i := 0; i < nNodes; i++ {
		oracleNode := utils.StartNewNode(t, owner, bootstrapNodePort+1+uint16(i), fmt.Sprintf("oracle%d", i), b, []commontypes.BootstrapperLocator{
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
	utils.SetOracleConfig(t, owner, oracleContract, oracles)
	utils.CommitWithFinality(b)

	// set up subscription
	subscriptionId := utils.CreateAndFundSubscriptions(t, owner, linkToken, registryAddress, registryContract, clientContracts)

	// send requests
	sent := make([][]byte, nClients)
	s := rand.NewSource(666)
	r := rand.New(s)
	for i := 0; i < nClients; i++ {
		sent[i] = []byte{byte(r.Uint32() % 256)}
		_, err := clientContracts[i].Contract.SendRequest(owner, hex.EncodeToString(sent[i]), []byte{}, []string{}, subscriptionId)
		require.NoError(t, err)
	}
	utils.CommitWithFinality(b)

	// validate that all DR-OCR jobs completed as many runs as sent requests
	var wg sync.WaitGroup
	for i := 0; i < nNodes; i++ {
		ic := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			cltest.WaitForPipelineComplete(t, ic, jobIds[ic], nClients, 5, apps[ic].JobORM(), 1*time.Minute, 1*time.Second)
		}()
	}
	wg.Wait()

	// validate that all client contracts got correct responses to their requests
	for i := 0; i < nClients; i++ {
		ic := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			gomega.NewGomegaWithT(t).Eventually(func() []byte {
				answer, err := clientContracts[ic].Contract.LastResponse(nil)
				require.NoError(t, err)
				return answer
			}, 1*time.Minute, 1*time.Second).Should(gomega.Equal(append([]byte{0xab}, sent[ic]...)))
		}()
	}
	wg.Wait()
}
