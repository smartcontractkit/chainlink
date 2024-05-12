// Package actions enables common chainlink interactions
package actions

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

// ContractDeploymentInterval After how many contract actions to wait before starting any more
// Example: When deploying 1000 contracts, stop every ContractDeploymentInterval have been deployed to wait before continuing
var ContractDeploymentInterval = 200

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amountCreateOCRv2Jobs of native currency
// Deprecated: we are moving away from blockchain.EVMClient, use actions_seth.FundChainlinkNodes
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
			Value: conversions.EtherToWei(amount),
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
	chainlinkNodes []*client.ChainlinkK8sClient,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	failingLogLevel zapcore.Level, // Examines logs after the test, and fails the test if any Chainlink logs are found at or above provided level
	grafnaUrlProvider testreporters.GrafanaURLProvider,
	clients ...blockchain.EVMClient,
) error {
	l := logging.GetTestLogger(t)
	if err := testreporters.WriteTeardownLogs(t, env, optionalTestReporter, failingLogLevel, grafnaUrlProvider); err != nil {
		return fmt.Errorf("Error dumping environment logs, leaving environment running for manual retrieval, err: %w", err)
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
// Deprecated: we are moving away from blockchain.EVMClient, use actions_seth.TeardownRemoteSuite
func TeardownRemoteSuite(
	t *testing.T,
	namespace string,
	chainlinkNodes []*client.ChainlinkK8sClient,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	grafnaUrlProvider testreporters.GrafanaURLProvider,
	client blockchain.EVMClient,
) error {
	l := logging.GetTestLogger(t)
	var err error
	if err = testreporters.SendReport(t, namespace, "./", optionalTestReporter, grafnaUrlProvider); err != nil {
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
			return fmt.Errorf("error reading jobs from chainlink node, err: %w", err)
		}
		for _, maps := range jobs.Data {
			if _, ok := maps["id"]; !ok {
				return fmt.Errorf("error reading job id from chainlink node's jobs %+v", jobs.Data)
			}
			id := maps["id"].(string)
			_, err := node.DeleteJob(id)
			if err != nil {
				return fmt.Errorf("error deleting job from chainlink node, err: %w", err)
			}
		}
	}
	return nil
}

// ReturnFunds attempts to return all the funds from the chainlink nodes to the network's default address
// all from a remote, k8s style environment
func ReturnFunds(chainlinkNodes []*client.ChainlinkK8sClient, blockchainClient blockchain.EVMClient) error {
	if blockchainClient == nil {
		return fmt.Errorf("blockchain client is nil, unable to return funds from chainlink nodes")
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
	if newImage == "" || newVersion == "" {
		return errors.New("New image and new version is needed to upgrade the node")
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

// todo - move to CTF
func GenerateWallet() (common.Address, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

// todo - move to CTF
func GetTxFromAddress(tx *types.Transaction) (string, error) {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	return from.String(), err
}

// todo - move to CTF
func DecodeTxInputData(abiString string, data []byte) (map[string]interface{}, error) {
	jsonABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}
	methodSigData := data[:4]
	inputsSigData := data[4:]
	method, err := jsonABI.MethodById(methodSigData)
	if err != nil {
		return nil, err
	}
	inputsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
		return nil, err
	}
	return inputsMap, nil
}

// todo - move to CTF
func WaitForBlockNumberToBe(
	waitForBlockNumberToBe uint64,
	client *seth.Client,
	wg *sync.WaitGroup,
	timeout time.Duration,
	t testing.TB,
	l zerolog.Logger,
) (uint64, error) {
	blockNumberChannel := make(chan uint64)
	errorChannel := make(chan error)
	testContext, testCancel := context.WithTimeout(context.Background(), timeout)
	defer testCancel()
	ticker := time.NewTicker(time.Second * 5)
	var latestBlockNumber uint64
	for {
		select {
		case <-testContext.Done():
			ticker.Stop()
			wg.Done()
			return latestBlockNumber,
				fmt.Errorf("timeout waiting for Block Number to be: %d. Last recorded block number was: %d",
					waitForBlockNumberToBe, latestBlockNumber)
		case <-ticker.C:
			go func() {
				currentBlockNumber, err := client.Client.BlockNumber(testcontext.Get(t))
				if err != nil {
					errorChannel <- err
				}
				l.Info().
					Uint64("Latest Block Number", currentBlockNumber).
					Uint64("Desired Block Number", waitForBlockNumberToBe).
					Msg("Waiting for Block Number to be")
				blockNumberChannel <- currentBlockNumber
			}()
		case latestBlockNumber = <-blockNumberChannel:
			if latestBlockNumber >= waitForBlockNumberToBe {
				ticker.Stop()
				wg.Done()
				l.Info().
					Uint64("Latest Block Number", latestBlockNumber).
					Uint64("Desired Block Number", waitForBlockNumberToBe).
					Msg("Desired Block Number reached!")
				return latestBlockNumber, nil
			}
		case err := <-errorChannel:
			ticker.Stop()
			wg.Done()
			return 0, err
		}
	}
}

// todo - move to EVMClient
func RewindSimulatedChainToBlockNumber(
	ctx context.Context,
	client *seth.Client,
	rpcURL string,
	rewindChainToBlockNumber uint64,
	l zerolog.Logger,
) (uint64, error) {
	latestBlockNumberBeforeReorg, err := client.Client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting latest block number: %w", err)
	}

	l.Info().
		Str("RPC URL", rpcURL).
		Uint64("Latest Block Number before Reorg", latestBlockNumberBeforeReorg).
		Uint64("Rewind Chain to Block Number", rewindChainToBlockNumber).
		Msg("Performing Reorg on chain by rewinding chain to specific block number")

	_, err = NewRPCRawClient(rpcURL).SetHeadForSimulatedChain(rewindChainToBlockNumber)

	if err != nil {
		return 0, fmt.Errorf("error making reorg: %w", err)
	}

	latestBlockNumberAfterReorg, err := client.Client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting latest block number: %w", err)
	}

	l.Info().
		Uint64("Block Number", latestBlockNumberAfterReorg).
		Msg("Latest Block Number after Reorg")
	return latestBlockNumberAfterReorg, nil
}

func GetRPCUrl(env *test_env.CLClusterTestEnv, chainID int64) (string, error) {
	provider, err := env.GetRpcProvider(chainID)
	if err != nil {
		return "", err
	}
	return provider.PublicHttpUrls()[0], nil
}

// RPCRawClient
// created separate client since method evmClient.RawJsonRPCCall fails on "invalid argument 0: json: cannot unmarshal non-string into Go value of type hexutil.Uint64"
type RPCRawClient struct {
	resty *resty.Client
}

func NewRPCRawClient(url string) *RPCRawClient {
	isDebug := os.Getenv("DEBUG_RESTY") == "true"
	restyClient := resty.New().SetDebug(isDebug).SetBaseURL(url)
	return &RPCRawClient{
		resty: restyClient,
	}
}

func (g *RPCRawClient) SetHeadForSimulatedChain(setHeadToBlockNumber uint64) (JsonRPCResponse, error) {
	var responseObject JsonRPCResponse
	postBody, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "debug_setHead",
		"params":  []string{hexutil.EncodeUint64(setHeadToBlockNumber)},
	})
	resp, err := g.resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(postBody).
		SetResult(&responseObject).
		Post("")

	if err != nil {
		return JsonRPCResponse{}, fmt.Errorf("error making API request: %w", err)
	}
	statusCode := resp.StatusCode()
	if statusCode != 200 && statusCode != 201 {
		return JsonRPCResponse{}, fmt.Errorf("error invoking debug_setHead method, received unexpected status code %d: %s", statusCode, resp.String())
	}
	if responseObject.Error != "" {
		return JsonRPCResponse{}, fmt.Errorf("received non-empty error field: %v", responseObject.Error)
	}
	return responseObject, nil
}

type JsonRPCResponse struct {
	Version string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
}
