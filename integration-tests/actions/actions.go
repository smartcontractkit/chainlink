// Package actions enables common chainlink interactions
package actions

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// ContractDeploymentInterval After how many contract actions to wait before starting any more
// Example: When deploying 1000 contracts, stop every ContractDeploymentInterval have been deployed to wait before continuing
var ContractDeploymentInterval = 200

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodes(
	nodes []*client.ChainlinkK8sClient,
	client blockchain.EVMClient,
	amount *big.Float,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		recipient := common.HexToAddress(toAddress)
		msg := ethereum.CallMsg{
			From:  common.HexToAddress(client.GetDefaultWallet().Address()),
			To:    &recipient,
			Value: utils.EtherToWei(amount),
		}
		gasEstimates, err := client.EstimateGas(msg)
		if err != nil {
			return err
		}
		err = client.Fund(toAddress, amount, gasEstimates)
		if err != nil {
			return err
		}
	}
	return client.WaitForEvents()
}

// FundChainlinkNodesAddress will fund all of the provided Chainlink nodes address at given index with a set amount of native currency
func FundChainlinkNodesAddress(
	nodes []*client.ChainlinkK8sClient,
	client blockchain.EVMClient,
	amount *big.Float,
	keyIndex int,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.EthAddresses()
		if err != nil {
			return err
		}
		toAddr := common.HexToAddress(toAddress[keyIndex])
		gasEstimates, err := client.EstimateGas(ethereum.CallMsg{
			To: &toAddr,
		})
		if err != nil {
			return err
		}
		err = client.Fund(toAddress[keyIndex], amount, gasEstimates)
		if err != nil {
			return err
		}
	}
	return client.WaitForEvents()
}

// FundChainlinkNodesAddress will fund all of the provided Chainlink nodes addresses with a set amount of native currency
func FundChainlinkNodesAddresses(
	nodes []*client.ChainlinkClient,
	client blockchain.EVMClient,
	amount *big.Float,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.EthAddressesForChain(client.GetChainID().String())
		if err != nil {
			return err
		}
		for _, addr := range toAddress {
			toAddr := common.HexToAddress(addr)
			gasEstimates, err := client.EstimateGas(ethereum.CallMsg{
				To: &toAddr,
			})
			if err != nil {
				return err
			}
			err = client.Fund(addr, amount, gasEstimates)
			if err != nil {
				return err
			}
		}
	}
	return client.WaitForEvents()
}

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodesLink(
	nodes []*client.ChainlinkK8sClient,
	blockchain blockchain.EVMClient,
	linkToken contracts.LinkToken,
	linkAmount *big.Int,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		err = linkToken.Transfer(toAddress, linkAmount)
		if err != nil {
			return err
		}
	}
	return blockchain.WaitForEvents()
}

// ChainlinkNodeAddresses will return all the on-chain wallet addresses for a set of Chainlink nodes
func ChainlinkNodeAddresses(nodes []*client.ChainlinkK8sClient) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, node := range nodes {
		primaryAddress, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(primaryAddress))
	}
	return addresses, nil
}

// ChainlinkNodeAddressesAtIndex will return all the on-chain wallet addresses for a set of Chainlink nodes
func ChainlinkNodeAddressesAtIndex(nodes []*client.ChainlinkK8sClient, keyIndex int) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, node := range nodes {
		nodeAddresses, err := node.EthAddresses()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(nodeAddresses[keyIndex]))
	}
	return addresses, nil
}

// SetChainlinkAPIPageSize specifies the page size from the Chainlink API, useful for high volume testing
func SetChainlinkAPIPageSize(nodes []*client.ChainlinkK8sClient, pageSize int) {
	for _, n := range nodes {
		n.SetPageSize(pageSize)
	}
}

// ExtractRequestIDFromJobRun extracts RequestID from job runs response
func ExtractRequestIDFromJobRun(jobDecodeData client.RunsResponseData) ([]byte, error) {
	var taskRun client.TaskRun
	for _, tr := range jobDecodeData.Attributes.TaskRuns {
		if tr.Type == "ethabidecodelog" {
			taskRun = tr
		}
	}
	var decodeLogTaskRun *client.DecodeLogTaskRun
	if err := json.Unmarshal([]byte(taskRun.Output), &decodeLogTaskRun); err != nil {
		return nil, err
	}
	rqInts := decodeLogTaskRun.RequestID
	return rqInts, nil
}

// EncodeOnChainVRFProvingKey encodes uncompressed public VRF key to on-chain representation
func EncodeOnChainVRFProvingKey(vrfKey client.VRFKey) ([2]*big.Int, error) {
	uncompressed := vrfKey.Data.Attributes.Uncompressed
	provingKey := [2]*big.Int{}
	var set1 bool
	var set2 bool
	// strip 0x to convert to int
	provingKey[0], set1 = new(big.Int).SetString(uncompressed[2:66], 16)
	if !set1 {
		return [2]*big.Int{}, fmt.Errorf("can not convert VRF key to *big.Int")
	}
	provingKey[1], set2 = new(big.Int).SetString(uncompressed[66:], 16)
	if !set2 {
		return [2]*big.Int{}, fmt.Errorf("can not convert VRF key to *big.Int")
	}
	return provingKey, nil
}

// GetMockserverInitializerDataForOTPE creates mocked weiwatchers data needed for otpe
func GetMockserverInitializerDataForOTPE(
	OCRInstances []contracts.OffchainAggregator,
	chainlinkNodes []*client.ChainlinkK8sClient,
) (interface{}, error) {
	var contractsInfo []ctfClient.ContractInfoJSON

	for index, OCRInstance := range OCRInstances {
		contractInfo := ctfClient.ContractInfoJSON{
			ContractVersion: 4,
			Path:            fmt.Sprintf("contract_%d", index),
			Status:          "live",
			ContractAddress: OCRInstance.Address(),
		}

		contractsInfo = append(contractsInfo, contractInfo)
	}

	contractsInitializer := ctfClient.HttpInitializer{
		Request:  ctfClient.HttpRequest{Path: "/contracts.json"},
		Response: ctfClient.HttpResponse{Body: contractsInfo},
	}

	var nodesInfo []ctfClient.NodeInfoJSON

	for _, chainlink := range chainlinkNodes {
		ocrKeys, err := chainlink.MustReadOCRKeys()
		if err != nil {
			return nil, err
		}
		nodeInfo := ctfClient.NodeInfoJSON{
			NodeAddress: []string{ocrKeys.Data[0].Attributes.OnChainSigningAddress},
			ID:          ocrKeys.Data[0].ID,
		}
		nodesInfo = append(nodesInfo, nodeInfo)
	}

	nodesInitializer := ctfClient.HttpInitializer{
		Request:  ctfClient.HttpRequest{Path: "/nodes.json"},
		Response: ctfClient.HttpResponse{Body: nodesInfo},
	}
	initializers := []ctfClient.HttpInitializer{contractsInitializer, nodesInitializer}
	return initializers, nil
}

// TeardownSuite tears down networks/clients and environment and creates a logs folder for failed tests in the
// specified path. Can also accept a testreporter (if one was used) to log further results
func TeardownSuite(
	t *testing.T,
	env *environment.Environment,
	logsFolderPath string,
	chainlinkNodes []*client.ChainlinkK8sClient,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	failingLogLevel zapcore.Level, // Examines logs after the test, and fails the test if any Chainlink logs are found at or above provided level
	clients ...blockchain.EVMClient,
) error {
	l := logging.GetTestLogger(t)
	if err := testreporters.WriteTeardownLogs(t, env, optionalTestReporter, failingLogLevel); err != nil {
		return errors.Wrap(err, "Error dumping environment logs, leaving environment running for manual retrieval")
	}
	// Delete all jobs to stop depleting the funds
	err := DeleteAllJobs(chainlinkNodes)
	if err != nil {
		l.Warn().Msgf("Error deleting jobs %+v", err)
	}

	for _, c := range clients {
		if c != nil && chainlinkNodes != nil && len(chainlinkNodes) > 0 {
			if err := ReturnFunds(chainlinkNodes, c); err != nil {
				// This printed line is required for tests that use real funds to propagate the failure
				// out to the system running the test. Do not remove
				fmt.Println(environment.FAILED_FUND_RETURN)
				l.Error().Err(err).Str("Namespace", env.Cfg.Namespace).
					Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
						"Environment is left running so you can try manually!")
			}
		} else {
			l.Info().Msg("Successfully returned funds from chainlink nodes to default network wallets")
		}
		// nolint
		if c != nil {
			err := c.Close()
			if err != nil {
				return err
			}
		}
	}

	return env.Shutdown()
}

// TeardownRemoteSuite is used when running a test within a remote-test-runner, like for long-running performance and
// soak tests
func TeardownRemoteSuite(
	t *testing.T,
	namespace string,
	chainlinkNodes []*client.ChainlinkK8sClient,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	client blockchain.EVMClient,
) error {
	l := logging.GetTestLogger(t)
	var err error
	if err = testreporters.SendReport(t, namespace, "./", optionalTestReporter); err != nil {
		l.Warn().Err(err).Msg("Error writing test report")
	}
	// Delete all jobs to stop depleting the funds
	err = DeleteAllJobs(chainlinkNodes)
	if err != nil {
		l.Warn().Msgf("Error deleting jobs %+v", err)
	}

	if err = ReturnFunds(chainlinkNodes, client); err != nil {
		l.Error().Err(err).Str("Namespace", namespace).
			Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
				"Environment is left running so you can try manually!")
	}
	return err
}

func DeleteAllJobs(chainlinkNodes []*client.ChainlinkK8sClient) error {
	for _, node := range chainlinkNodes {
		if node == nil {
			return fmt.Errorf("found a nil chainlink node in the list of chainlink nodes while tearing down: %v", chainlinkNodes)
		}
		jobs, _, err := node.ReadJobs()
		if err != nil {
			return errors.Wrap(err, "error reading jobs from chainlink node")
		}
		for _, maps := range jobs.Data {
			if _, ok := maps["id"]; !ok {
				return errors.Errorf("error reading job id from chainlink node's jobs %+v", jobs.Data)
			}
			id := maps["id"].(string)
			_, err := node.DeleteJob(id)
			if err != nil {
				return errors.Wrap(err, "error deleting job from chainlink node")
			}
		}
	}
	return nil
}

// ReturnFunds attempts to return all the funds from the chainlink nodes to the network's default address
// all from a remote, k8s style environment
func ReturnFunds(chainlinkNodes []*client.ChainlinkK8sClient, blockchainClient blockchain.EVMClient) error {
	if blockchainClient == nil {
		return errors.New("blockchain client is nil, unable to return funds from chainlink nodes")
	}
	log.Info().Msg("Attempting to return Chainlink node funds to default network wallets")
	if blockchainClient.NetworkSimulated() {
		log.Info().Str("Network Name", blockchainClient.GetNetworkName()).
			Msg("Network is a simulated network. Skipping fund return.")
		return nil
	}

	for _, chainlinkNode := range chainlinkNodes {
		fundedKeys, err := chainlinkNode.ExportEVMKeysForChain(blockchainClient.GetChainID().String())
		if err != nil {
			return err
		}
		for _, key := range fundedKeys {
			keyToDecrypt, err := json.Marshal(key)
			if err != nil {
				return err
			}
			// This can take up a good bit of RAM and time. When running on the remote-test-runner, this can lead to OOM
			// issues. So we avoid running in parallel; slower, but safer.
			decryptedKey, err := keystore.DecryptKey(keyToDecrypt, client.ChainlinkKeyPassword)
			if err != nil {
				return err
			}
			err = blockchainClient.ReturnFunds(decryptedKey.PrivateKey)
			if err != nil {
				log.Error().Err(err).Str("Address", fundedKeys[0].Address).Msg("Error returning funds from Chainlink node")
			}
		}
	}
	return blockchainClient.WaitForEvents()
}

// FundAddresses will fund a list of addresses with an amount of native currency
func FundAddresses(blockchain blockchain.EVMClient, amount *big.Float, addresses ...string) error {
	for _, address := range addresses {
		toAddr := common.HexToAddress(address)
		gasEstimates, err := blockchain.EstimateGas(ethereum.CallMsg{
			To: &toAddr,
		})
		if err != nil {
			return err
		}
		if err := blockchain.Fund(address, amount, gasEstimates); err != nil {
			return err
		}
	}
	return blockchain.WaitForEvents()
}

// EncodeOnChainExternalJobID encodes external job uuid to on-chain representation
func EncodeOnChainExternalJobID(jobID uuid.UUID) [32]byte {
	var ji [32]byte
	copy(ji[:], strings.Replace(jobID.String(), "-", "", 4))
	return ji
}

// UpgradeChainlinkNodeVersions upgrades all Chainlink nodes to a new version, and then runs the test environment
// to apply the upgrades
func UpgradeChainlinkNodeVersions(
	testEnvironment *environment.Environment,
	newImage, newVersion string,
	nodes ...*client.ChainlinkK8sClient,
) error {
	if newImage == "" && newVersion == "" {
		return errors.New("unable to upgrade node version, found empty image and version, must provide either a new image or a new version")
	}
	for _, node := range nodes {
		if err := node.UpgradeVersion(testEnvironment, newImage, newVersion); err != nil {
			return err
		}
	}
	err := testEnvironment.RunUpdated(len(nodes))
	if err != nil { // Run the new environment and wait for changes to show
		return err
	}
	return client.ReconnectChainlinkNodes(testEnvironment, nodes)
}

func DeployLINKToken(cd contracts.ContractDeployer) (contracts.LinkToken, error) {
	linkToken, err := cd.DeployLinkTokenContract()
	if err != nil {
		return nil, err
	}
	return linkToken, err
}

func DeployMockETHLinkFeed(cd contracts.ContractDeployer, answer *big.Int) (contracts.MockETHLINKFeed, error) {
	mockETHLINKFeed, err := cd.DeployMockETHLINKFeed(answer)
	if err != nil {
		return nil, err
	}
	return mockETHLINKFeed, err
}
