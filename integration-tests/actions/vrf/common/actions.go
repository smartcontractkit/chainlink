package common

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/seth"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
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
		_, err = actions_seth.SendFunds(l, client, actions_seth.FundsToSendPayload{
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
	env *test_env.CLClusterTestEnv,
	config *vrf_common_config.General,
	numberOfTxKeysToCreate int,
	chainID *big.Int,
	coordinatorAddress string,
	BHSAddress string,
	txKeyFunding float64,
	l zerolog.Logger,
	bhsNode *VRFNode,
) error {
	sethClient, err := env.GetSethClient(chainID.Int64())
	if err != nil {
		return err
	}

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
	env *test_env.CLClusterTestEnv,
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
	sethClient, err := env.GetSethClient(chainID.Int64())
	if err != nil {
		return err
	}
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
				_, err := actions_seth.SendFunds(l, client, actions_seth.FundsToSendPayload{
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

func BuildNewCLEnvForVRF(t *testing.T, envConfig VRFEnvConfig, newEnvConfig NewEnvConfig, network ctf_test_env.EthereumNetwork) (*test_env.CLClusterTestEnv, error) {
	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&envConfig.TestConfig).
		WithPrivateEthereumNetwork(network.EthereumNetworkConfig).
		WithCLNodes(len(newEnvConfig.NodesToCreate)).
		WithFunding(big.NewFloat(*envConfig.TestConfig.Common.ChainlinkNodeFunding)).
		WithChainlinkNodeLogScanner(newEnvConfig.ChainlinkNodeLogScannerSettings).
		WithCustomCleanup(envConfig.CleanupFn).
		WithSeth().
		Build()
	if err != nil {
		return nil, fmt.Errorf("%s, err: %w", "error creating test env", err)
	}
	return env, nil
}
