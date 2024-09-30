package common

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	vrf_common_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/common/vrf"
)

func CreateFundAndGetSendingKeys(
	l zerolog.Logger,
	client *seth.Client,
	node *VRFNode,
	chainlinkNodeFunding float64,
	numberOfTxKeysToCreate int,
	chainID *big.Int,
) ([]string, []common.Address, error) {
	newNativeTokenKeyAddresses, err := CreateAndFundSendingKeys(l, client, node, chainlinkNodeFunding, numberOfTxKeysToCreate, chainID)
	if err != nil {
		return nil, nil, err
	}
	nativeTokenPrimaryKeyAddress, err := node.CLNode.API.PrimaryEthAddress()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrNodePrimaryKey, err)
	}
	allNativeTokenKeyAddressStrings := append(newNativeTokenKeyAddresses, nativeTokenPrimaryKeyAddress)
	allNativeTokenKeyAddresses := make([]common.Address, len(allNativeTokenKeyAddressStrings))
	for _, addressString := range allNativeTokenKeyAddressStrings {
		allNativeTokenKeyAddresses = append(allNativeTokenKeyAddresses, common.HexToAddress(addressString))
	}
	return allNativeTokenKeyAddressStrings, allNativeTokenKeyAddresses, nil
}

func CreateAndFundSendingKeys(
	l zerolog.Logger,
	client *seth.Client,
	node *VRFNode,
	chainlinkNodeFunding float64,
	numberOfNativeTokenAddressesToCreate int,
	chainID *big.Int,
) ([]string, error) {
	var newNativeTokenKeyAddresses []string
	for i := 0; i < numberOfNativeTokenAddressesToCreate; i++ {
		newTxKey, response, err := node.CLNode.API.CreateTxKey("evm", chainID.String())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrNodeNewTxKey, err)
		}
		if response.StatusCode != 200 {
			return nil, fmt.Errorf("error creating transaction key - response code, err %d", response.StatusCode)
		}
		newNativeTokenKeyAddresses = append(newNativeTokenKeyAddresses, newTxKey.Data.Attributes.Address)
		_, err = actions.SendFunds(l, client, actions.FundsToSendPayload{
			ToAddress:  common.HexToAddress(newTxKey.Data.Attributes.Address),
			Amount:     conversions.EtherToWei(big.NewFloat(chainlinkNodeFunding)),
			PrivateKey: client.PrivateKeys[0],
		})
		if err != nil {
			return nil, err
		}
	}
	return newNativeTokenKeyAddresses, nil
}

func SetupBHSNode(
	sethClient *seth.Client,
	config *vrf_common_config.General,
	numberOfTxKeysToCreate int,
	chainID *big.Int,
	coordinatorAddress string,
	BHSAddress string,
	txKeyFunding float64,
	l zerolog.Logger,
	bhsNode *VRFNode,
) error {
	bhsTXKeyAddressStrings, _, err := CreateFundAndGetSendingKeys(
		l,
		sethClient,
		bhsNode,
		txKeyFunding,
		numberOfTxKeysToCreate,
		chainID,
	)
	if err != nil {
		return err
	}
	bhsNode.TXKeyAddressStrings = bhsTXKeyAddressStrings
	bhsSpec := client.BlockhashStoreJobSpec{
		ForwardingAllowed:        false,
		CoordinatorV2Address:     coordinatorAddress,
		CoordinatorV2PlusAddress: coordinatorAddress,
		BlockhashStoreAddress:    BHSAddress,
		FromAddresses:            bhsTXKeyAddressStrings,
		EVMChainID:               chainID.String(),
		WaitBlocks:               *config.BHSJobWaitBlocks,
		LookbackBlocks:           *config.BHSJobLookBackBlocks,
		PollPeriod:               config.BHSJobPollPeriod.Duration,
		RunTimeout:               config.BHSJobRunTimeout.Duration,
	}
	l.Info().Msg("Creating BHS Job")
	bhsJob, err := CreateBHSJob(
		bhsNode.CLNode.API,
		bhsSpec,
	)
	if err != nil {
		return fmt.Errorf("%s, err %w", "", err)
	}
	bhsNode.Job = bhsJob
	return nil
}

func CreateBHSJob(
	chainlinkNode *client.ChainlinkClient,
	bhsJobSpecConfig client.BlockhashStoreJobSpec,
) (*client.Job, error) {
	jobUUID := uuid.New()
	spec := &client.BlockhashStoreJobSpec{
		Name:                     fmt.Sprintf("bhs-%s", jobUUID),
		ForwardingAllowed:        bhsJobSpecConfig.ForwardingAllowed,
		CoordinatorV2Address:     bhsJobSpecConfig.CoordinatorV2Address,
		CoordinatorV2PlusAddress: bhsJobSpecConfig.CoordinatorV2PlusAddress,
		BlockhashStoreAddress:    bhsJobSpecConfig.BlockhashStoreAddress,
		FromAddresses:            bhsJobSpecConfig.FromAddresses,
		EVMChainID:               bhsJobSpecConfig.EVMChainID,
		ExternalJobID:            jobUUID.String(),
		WaitBlocks:               bhsJobSpecConfig.WaitBlocks,
		LookbackBlocks:           bhsJobSpecConfig.LookbackBlocks,
		PollPeriod:               bhsJobSpecConfig.PollPeriod,
		RunTimeout:               bhsJobSpecConfig.RunTimeout,
	}

	job, err := chainlinkNode.MustCreateJob(spec)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrCreatingBHSJob, err)
	}
	return job, nil
}

func SetupBHFNode(
	sethClient *seth.Client,
	config *vrf_common_config.General,
	numberOfTxKeysToCreate int,
	chainID *big.Int,
	coordinatorAddress string,
	BHSAddress string,
	batchBHSAddress string,
	txKeyFunding float64,
	l zerolog.Logger,
	bhfNode *VRFNode,
) error {
	bhfTXKeyAddressStrings, _, err := CreateFundAndGetSendingKeys(
		l,
		sethClient,
		bhfNode,
		txKeyFunding,
		numberOfTxKeysToCreate,
		chainID,
	)
	if err != nil {
		return err
	}
	bhfNode.TXKeyAddressStrings = bhfTXKeyAddressStrings
	bhfSpec := client.BlockHeaderFeederJobSpec{
		ForwardingAllowed:          false,
		CoordinatorV2Address:       coordinatorAddress,
		CoordinatorV2PlusAddress:   coordinatorAddress,
		BlockhashStoreAddress:      BHSAddress,
		BatchBlockhashStoreAddress: batchBHSAddress,
		FromAddresses:              bhfTXKeyAddressStrings,
		EVMChainID:                 chainID.String(),
		WaitBlocks:                 *config.BHFJobWaitBlocks,
		LookbackBlocks:             *config.BHFJobLookBackBlocks,
		PollPeriod:                 config.BHFJobPollPeriod.Duration,
		RunTimeout:                 config.BHFJobRunTimeout.Duration,
	}
	l.Info().Msg("Creating BHF Job")
	bhfJob, err := CreateBHFJob(
		bhfNode.CLNode.API,
		bhfSpec,
	)
	if err != nil {
		return fmt.Errorf("%s, err %w", "", err)
	}
	bhfNode.Job = bhfJob
	return nil
}

func CreateBHFJob(
	chainlinkNode *client.ChainlinkClient,
	bhfJobSpecConfig client.BlockHeaderFeederJobSpec,
) (*client.Job, error) {
	jobUUID := uuid.New()
	spec := &client.BlockHeaderFeederJobSpec{
		Name:                       fmt.Sprintf("bhf-%s", jobUUID),
		ForwardingAllowed:          bhfJobSpecConfig.ForwardingAllowed,
		CoordinatorV2Address:       bhfJobSpecConfig.CoordinatorV2Address,
		CoordinatorV2PlusAddress:   bhfJobSpecConfig.CoordinatorV2PlusAddress,
		BlockhashStoreAddress:      bhfJobSpecConfig.BlockhashStoreAddress,
		BatchBlockhashStoreAddress: bhfJobSpecConfig.BatchBlockhashStoreAddress,
		FromAddresses:              bhfJobSpecConfig.FromAddresses,
		EVMChainID:                 bhfJobSpecConfig.EVMChainID,
		ExternalJobID:              jobUUID.String(),
		WaitBlocks:                 bhfJobSpecConfig.WaitBlocks,
		LookbackBlocks:             bhfJobSpecConfig.LookbackBlocks,
		PollPeriod:                 bhfJobSpecConfig.PollPeriod,
		RunTimeout:                 bhfJobSpecConfig.RunTimeout,
	}

	job, err := chainlinkNode.MustCreateJob(spec)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrCreatingBHSJob, err)
	}
	return job, nil
}

func WaitForRequestCountEqualToFulfilmentCount(
	ctx context.Context,
	consumer VRFLoadTestConsumer,
	timeout time.Duration,
	wg *sync.WaitGroup,
) (*big.Int, *big.Int, error) {
	metricsChannel := make(chan *contracts.VRFLoadTestMetrics)
	metricsErrorChannel := make(chan error)

	testContext, testCancel := context.WithTimeout(ctx, timeout)
	defer testCancel()

	ticker := time.NewTicker(time.Second * 1)
	var metrics *contracts.VRFLoadTestMetrics
	for {
		select {
		case <-testContext.Done():
			ticker.Stop()
			wg.Done()
			return metrics.RequestCount, metrics.FulfilmentCount,
				fmt.Errorf("timeout waiting for rand request and fulfilments to be equal AFTER performance test was executed. Request Count: %d, Fulfilment Count: %d",
					metrics.RequestCount.Uint64(), metrics.FulfilmentCount.Uint64())
		case <-ticker.C:
			go retrieveLoadTestMetrics(ctx, consumer, metricsChannel, metricsErrorChannel)
		case metrics = <-metricsChannel:
			if metrics.RequestCount.Cmp(metrics.FulfilmentCount) == 0 {
				ticker.Stop()
				wg.Done()
				return metrics.RequestCount, metrics.FulfilmentCount, nil
			}
		case err := <-metricsErrorChannel:
			ticker.Stop()
			wg.Done()
			return nil, nil, err
		}
	}
}

func retrieveLoadTestMetrics(
	ctx context.Context,
	consumer VRFLoadTestConsumer,
	metricsChannel chan *contracts.VRFLoadTestMetrics,
	metricsErrorChannel chan error,
) {
	metrics, err := consumer.GetLoadTestMetrics(ctx)
	if err != nil {
		metricsErrorChannel <- err
	}
	metricsChannel <- metrics
}

func CreateNodeTypeToNodeMap(cluster *test_env.ClCluster, nodesToCreate []VRFNodeType) (map[VRFNodeType]*VRFNode, error) {
	var nodesMap = make(map[VRFNodeType]*VRFNode)
	if len(cluster.Nodes) < len(nodesToCreate) {
		return nil, fmt.Errorf("not enough nodes in the cluster (cluster size is %d nodes) to create %d nodes", len(cluster.Nodes), len(nodesToCreate))
	}
	for i, nodeType := range nodesToCreate {
		nodesMap[nodeType] = &VRFNode{
			CLNode: cluster.Nodes[i],
		}
	}
	return nodesMap, nil
}

func CreateVRFKeyOnVRFNode(vrfNode *VRFNode, l zerolog.Logger) (*client.VRFKey, string, error) {
	l.Info().Str("Node URL", vrfNode.CLNode.API.URL()).Msg("Creating VRF Key on the Node")
	vrfKey, err := vrfNode.CLNode.API.MustCreateVRFKey()
	if err != nil {
		return nil, "", fmt.Errorf("%s, err %w", ErrCreatingVRFKey, err)
	}
	pubKeyCompressed := vrfKey.Data.ID
	l.Info().
		Str("Node URL", vrfNode.CLNode.API.URL()).
		Str("Keyhash", vrfKey.Data.Attributes.Hash).
		Str("VRF Compressed Key", vrfKey.Data.Attributes.Compressed).
		Str("VRF Uncompressed Key", vrfKey.Data.Attributes.Uncompressed).
		Msg("VRF Key created on the Node")
	return vrfKey, pubKeyCompressed, nil
}

func FundNodesIfNeeded(ctx context.Context, existingEnvConfig *vrf_common_config.ExistingEnvConfig, client *seth.Client, l zerolog.Logger) error {
	if *existingEnvConfig.NodeSendingKeyFundingMin > 0 {
		for _, sendingKey := range existingEnvConfig.NodeSendingKeys {
			address := common.HexToAddress(sendingKey)
			sendingKeyBalance, err := client.Client.BalanceAt(ctx, address, nil)
			if err != nil {
				return err
			}
			fundingAtLeast := conversions.EtherToWei(big.NewFloat(*existingEnvConfig.NodeSendingKeyFundingMin))
			fundingToSendWei := new(big.Int).Sub(fundingAtLeast, sendingKeyBalance)
			log := l.Info().
				Str("Sending Key", sendingKey).
				Str("Sending Key Current Balance", sendingKeyBalance.String()).
				Str("Should have at least", fundingAtLeast.String())
			if fundingToSendWei.Cmp(big.NewInt(0)) == 1 {
				log.
					Str("Funding Amount in wei", fundingToSendWei.String()).
					Str("Funding Amount in ETH", conversions.WeiToEther(fundingToSendWei).String()).
					Msg("Funding Node's Sending Key")
				_, err := actions.SendFunds(l, client, actions.FundsToSendPayload{
					ToAddress:  common.HexToAddress(sendingKey),
					Amount:     fundingToSendWei,
					PrivateKey: client.PrivateKeys[0],
				})
				if err != nil {
					return err
				}
			} else {
				log.
					Msg("Skipping Node's Sending Key funding as it has enough funds")
			}
		}
	}
	return nil
}

func BuildNewCLEnvForVRF(l zerolog.Logger, t *testing.T, envConfig VRFEnvConfig, newEnvConfig NewEnvConfig, network ctf_test_env.EthereumNetwork) (*test_env.CLClusterTestEnv, *seth.Client, error) {
	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&envConfig.TestConfig).
		WithPrivateEthereumNetwork(network.EthereumNetworkConfig).
		WithCLNodes(len(newEnvConfig.NodesToCreate)).
		WithChainlinkNodeLogScanner(newEnvConfig.ChainlinkNodeLogScannerSettings).
		WithCustomCleanup(envConfig.CleanupFn).
		Build()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err: %w", "error creating test env", err)
	}

	evmNetwork, err := env.GetFirstEvmNetwork()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err: %w", "error getting first evm network", err)
	}
	sethClient, err := utils.TestAwareSethClient(t, envConfig.TestConfig, evmNetwork)
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err: %w", "error getting seth client", err)
	}

	err = actions.FundChainlinkNodesFromRootAddress(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()), big.NewFloat(*envConfig.TestConfig.Common.ChainlinkNodeFunding))
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err: %w", "failed to fund the nodes", err)
	}

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		_ = actions.ReturnFundsFromNodes(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()))
	})

	return env, sethClient, nil
}

func LoadExistingCLEnvForVRF(
	t *testing.T,
	envConfig VRFEnvConfig,
	commonExistingEnvConfig *vrf_common_config.ExistingEnvConfig,
	l zerolog.Logger,
) (*test_env.CLClusterTestEnv, *seth.Client, error) {
	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&envConfig.TestConfig).
		WithCustomCleanup(envConfig.CleanupFn).
		Build()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err: %w", "error creating test env", err)
	}
	evmNetwork, err := env.GetFirstEvmNetwork()
	if err != nil {
		return nil, nil, err
	}
	sethClient, err := utils.TestAwareSethClient(t, envConfig.TestConfig, evmNetwork)
	if err != nil {
		return nil, nil, err
	}
	err = FundNodesIfNeeded(testcontext.Get(t), commonExistingEnvConfig, sethClient, l)
	if err != nil {
		return nil, nil, err
	}
	return env, sethClient, nil
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
	isDebug := os.Getenv("RESTY_DEBUG") == "true"
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

// todo - move to CTF
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
