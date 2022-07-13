// Package actions enables common chainlink interactions
package actions

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/ginkgo/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"golang.org/x/sync/errgroup"
)

// ContractDeploymentInterval After how many contract actions to wait before starting any more
// Example: When deploying 1000 contracts, stop every ContractDeploymentInterval have been deployed to wait before continuing
var ContractDeploymentInterval = 200

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodes(
	nodes []client.Chainlink,
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

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodesLink(
	nodes []client.Chainlink,
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
func ChainlinkNodeAddresses(nodes []client.Chainlink) ([]common.Address, error) {
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

// SetChainlinkAPIPageSize specifies the page size from the Chainlink API, useful for high volume testing
func SetChainlinkAPIPageSize(nodes []client.Chainlink, pageSize int) {
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
	chainlinkNodes []client.Chainlink,
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
		ocrKeys, err := chainlink.ReadOCRKeys()
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
	env *environment.Environment,
	logsFolderPath string,
	chainlinkNodes []client.Chainlink,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	c blockchain.EVMClient,
) error {
	if err := actions.WriteTeardownLogs(env, optionalTestReporter); err != nil {
		return errors.Wrap(err, "Error dumping environment logs, leaving environment running for manual retrieval")
	}
	if c != nil && chainlinkNodes != nil && len(chainlinkNodes) > 0 {
		if err := returnFunds(chainlinkNodes, c); err != nil {
			log.Error().Err(err).Str("Namespace", env.Cfg.Namespace).
				Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
					"Environment is left running so you can try manually!")
		}
	} else {
		log.Info().Msg("Successfully returned funds from chainlink nodes to default network wallets")
	}
	// nolint
	if c != nil {
		c.Close()
	}

	keepEnvs := os.Getenv("KEEP_ENVIRONMENTS")
	if keepEnvs == "" {
		keepEnvs = "NEVER"
	}

	switch strings.ToUpper(keepEnvs) {
	case "ALWAYS":
	case "ONFAIL":
		if ginkgo.CurrentSpecReport().Failed() {
			return env.Shutdown()
		}
	case "NEVER":
		return env.Shutdown()
	default:
		log.Warn().Str("Invalid Keep Value", keepEnvs).
			Msg("Invalid 'keep_environments' value, see the 'framework.yaml' file")
	}
	return nil
}

// TeardownRemoteSuite is used when running a test within a remote-test-runner, like for long-running performance and
// soak tests
func TeardownRemoteSuite(
	env *environment.Environment,
	chainlinkNodes []client.Chainlink,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	client blockchain.EVMClient,
) error {
	var err error
	if err = actions.SendReport(env, "./", optionalTestReporter); err != nil {
		log.Warn().Err(err).Msg("Error writing test report")
	}
	if err = returnFunds(chainlinkNodes, client); err != nil {
		log.Error().Err(err).Str("Namespace", env.Cfg.Namespace).
			Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
				"Environment is left running so you can try manually!")
	}
	return err
}

// Returns all the funds from the chainlink nodes to the networks default address
func returnFunds(chainlinkNodes []client.Chainlink, client blockchain.EVMClient) error {
	if client == nil {
		log.Warn().Msg("No blockchain client found, unable to return funds from chainlink nodes.")
	}
	for _, node := range chainlinkNodes {
		if err := node.SetSessionCookie(); err != nil {
			return err
		}
	}
	log.Info().Msg("Attempting to return Chainlink node funds to default network wallets")
	if client.NetworkSimulated() {
		log.Info().Str("Network Name", client.GetNetworkName()).
			Msg("Network is a simulated network. Skipping fund return.")
		return nil
	}

	addressMap, err := sendFunds(chainlinkNodes, client)
	if err != nil {
		return err
	}

	err = checkFunds(chainlinkNodes, addressMap, strings.ToLower(client.GetDefaultWallet().Address()))
	if err != nil {
		return err
	}
	addressMap, err = sendFunds(chainlinkNodes, client)
	if err != nil {
		return err
	}
	return checkFunds(chainlinkNodes, addressMap, strings.ToLower(client.GetDefaultWallet().Address()))
}

// Requests that all the chainlink nodes send their funds back to the network's default wallet
// This is surprisingly tricky, and fairly annoying due to Go's lack of syntactic sugar and how chainlink nodes handle txs
func sendFunds(chainlinkNodes []client.Chainlink, network blockchain.EVMClient) (map[int]string, error) {
	chainlinkTransactionAddresses := make(map[int]string)
	sendFundsErrGroup := new(errgroup.Group)
	for ni, n := range chainlinkNodes {
		nodeIndex := ni // https://golang.org/doc/faq#closures_and_goroutines
		node := n
		// Send async request to each chainlink node to send a transaction back to the network default wallet
		sendFundsErrGroup.Go(
			func() error {
				primaryEthKeyData, err := node.ReadPrimaryETHKey()
				if err != nil {
					// TODO: Support non-EVM chain fund returns
					if strings.Contains(err.Error(), "No ETH keys present") {
						log.Warn().Msg("Not returning any funds. Only support EVM chains for fund returns at the moment")
						return nil
					}
					return err
				}

				nodeBalanceString := primaryEthKeyData.Attributes.ETHBalance
				if nodeBalanceString != "0" { // If key has a non-zero balance, attempt to transfer it back
					gasCost, err := network.EstimateTransactionGasCost()
					if err != nil {
						return err
					}

					// TODO: Imperfect gas calculation buffer of 50 Gwei. Seems to be the result of differences in chainlink
					// gas handling. Working with core team on a better solution
					gasCost = gasCost.Add(gasCost, big.NewInt(50000000000))
					nodeBalance, _ := big.NewInt(0).SetString(nodeBalanceString, 10)
					transferAmount := nodeBalance.Sub(nodeBalance, gasCost)
					_, err = node.SendNativeToken(transferAmount, primaryEthKeyData.Attributes.Address, network.GetDefaultWallet().Address())
					if err != nil {
						return err
					}
					// Add the address to our map to check for later (hashes aren't returned, sadly)
					chainlinkTransactionAddresses[nodeIndex] = strings.ToLower(primaryEthKeyData.Attributes.Address)
				}
				return nil
			},
		)

	}
	return chainlinkTransactionAddresses, sendFundsErrGroup.Wait()
}

// checks that the funds made it from the chainlink node to the network address
// this turns out to be tricky to do, given how chainlink handles pending transactions, thus the complexity
func checkFunds(chainlinkNodes []client.Chainlink, sentFromAddressesMap map[int]string, toAddress string) error {
	successfulConfirmations := make(map[int]bool)
	err := retry.Do( // Might take some time for txs to confirm, check up on the nodes a few times
		func() error {
			log.Debug().Msg("Attempting to confirm chainlink nodes transferred back funds")
			transactionErrGroup := new(errgroup.Group)
			for i, n := range chainlinkNodes {
				nodeIndex := i
				node := n // https://golang.org/doc/faq#closures_and_goroutines
				sentFromAddress, nodeHasFunds := sentFromAddressesMap[nodeIndex]
				successfulConfirmation := successfulConfirmations[nodeIndex]
				// Async check on all the nodes if their transactions are confirmed
				if nodeHasFunds && !successfulConfirmation { // Only if node has funds and hasn't already sent them
					transactionErrGroup.Go(func() error {
						err := confirmTransaction(node, sentFromAddress, toAddress)
						if err == nil {
							successfulConfirmations[nodeIndex] = true
						}
						return err
					})
				} else {
					log.Debug().Int("Node Number", nodeIndex).Msg("Chainlink node had no funds to return")
				}
			}

			return transactionErrGroup.Wait()
		},
		retry.Delay(time.Second*5),
		retry.MaxDelay(time.Second*5),
		retry.Attempts(20),
	)

	return err
}

// helper to confirm that the latest attempted transaction on the chainlink node with the expected from and to addresses
// has been confirmed
func confirmTransaction(
	chainlinkNode client.Chainlink,
	fromAddress string,
	toAddress string,
) error {
	transactionAttempts, err := chainlinkNode.ReadTransactionAttempts()
	if err != nil {
		return err
	}
	log.Debug().Str("From", fromAddress).
		Str("To", toAddress).
		Msg("Attempting to confirm node returned funds")
	// Loop through all transactions on the node
	for _, tx := range transactionAttempts.Data {
		if tx.Attributes.From == fromAddress && strings.ToLower(tx.Attributes.To) == toAddress {
			if tx.Attributes.State == "confirmed" {
				return nil
			}
			return fmt.Errorf("Expected transaction to be confirmed. From: %s To: %s State: %s", fromAddress, toAddress, tx.Attributes.State)
		}
	}
	return fmt.Errorf("Did not find expected transaction on node. From: %s To: %s", fromAddress, toAddress)
}
