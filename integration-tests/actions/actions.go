// Package actions enables common chainlink interactions
package actions

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// ContractDeploymentInterval After how many contract actions to wait before starting any more
// Example: When deploying 1000 contracts, stop every ContractDeploymentInterval have been deployed to wait before continuing
var ContractDeploymentInterval = 200

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodes(
	nodes []*client.Chainlink,
	client blockchain.EVMClient,
	amount *big.Float,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		err = client.Fund(toAddress, amount)
		if err != nil {
			return err
		}
	}
	return client.WaitForEvents()
}

// FundChainlinkNodesAddress will fund all of the provided Chainlink nodes address at given index with a set amount of native currency
func FundChainlinkNodesAddress(
	nodes []*client.Chainlink,
	client blockchain.EVMClient,
	amount *big.Float,
	keyIndex int,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.EthAddresses()
		if err != nil {
			return err
		}
		err = client.Fund(toAddress[keyIndex], amount)
		if err != nil {
			return err
		}
	}
	return client.WaitForEvents()
}

// FundChainlinkNodesAddress will fund all of the provided Chainlink nodes addresses with a set amount of native currency
func FundChainlinkNodesAddresses(
	nodes []*client.Chainlink,
	client blockchain.EVMClient,
	amount *big.Float,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.EthAddressesForChain(client.GetChainID().String())
		if err != nil {
			return err
		}
		for _, addr := range toAddress {
			err = client.Fund(addr, amount)
			if err != nil {
				return err
			}
		}
	}
	return client.WaitForEvents()
}

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodesLink(
	nodes []*client.Chainlink,
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
func ChainlinkNodeAddresses(nodes []*client.Chainlink) ([]common.Address, error) {
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
func ChainlinkNodeAddressesAtIndex(nodes []*client.Chainlink, keyIndex int) ([]common.Address, error) {
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
func SetChainlinkAPIPageSize(nodes []*client.Chainlink, pageSize int) {
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
	chainlinkNodes []*client.Chainlink,
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
	chainlinkNodes []*client.Chainlink,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	failingLogLevel zapcore.Level, // Examines logs after the test, and fails the test if any Chainlink logs are found at or above provided level
	clients ...blockchain.EVMClient,
) error {
	l := utils.GetTestLogger(t)
	if err := testreporters.WriteTeardownLogs(t, env, optionalTestReporter, failingLogLevel); err != nil {
		return errors.Wrap(err, "Error dumping environment logs, leaving environment running for manual retrieval")
	}
	for _, c := range clients {
		if c != nil && chainlinkNodes != nil && len(chainlinkNodes) > 0 {
			if err := returnFunds(chainlinkNodes, c); err != nil {
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
	env *environment.Environment,
	chainlinkNodes []*client.Chainlink,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	client blockchain.EVMClient,
) error {
	l := utils.GetTestLogger(t)
	var err error
	if err = testreporters.SendReport(t, env, "./", optionalTestReporter); err != nil {
		l.Warn().Err(err).Msg("Error writing test report")
	}
	if err = returnFunds(chainlinkNodes, client); err != nil {
		l.Error().Err(err).Str("Namespace", env.Cfg.Namespace).
			Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
				"Environment is left running so you can try manually!")
	}
	return err
}

// Returns all the funds from the chainlink nodes to the networks default address
func returnFunds(chainlinkNodes []*client.Chainlink, blockchainClient blockchain.EVMClient) error {
	if blockchainClient == nil {
		log.Warn().Msg("No blockchain client found, unable to return funds from chainlink nodes.")
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
				return err
			}
		}
	}
	return blockchainClient.WaitForEvents()
}

// FundAddresses will fund a list of addresses with an amount of native currency
func FundAddresses(blockchain blockchain.EVMClient, amount *big.Float, addresses ...string) error {
	for _, address := range addresses {
		if err := blockchain.Fund(address, amount); err != nil {
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
