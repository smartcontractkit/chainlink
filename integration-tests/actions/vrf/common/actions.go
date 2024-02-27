package common

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	vrf_common_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/common/vrf"
)

func CreateFundAndGetSendingKeys(
	client blockchain.EVMClient,
	node *VRFNode,
	chainlinkNodeFunding float64,
	numberOfTxKeysToCreate int,
	chainID *big.Int,
) ([]string, []common.Address, error) {
	newNativeTokenKeyAddresses, err := CreateAndFundSendingKeys(client, node, chainlinkNodeFunding, numberOfTxKeysToCreate, chainID)
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
	client blockchain.EVMClient,
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
		err = actions.FundAddress(client, newTxKey.Data.Attributes.Address, big.NewFloat(chainlinkNodeFunding))
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
	bhsTXKeyAddressStrings, _, err := CreateFundAndGetSendingKeys(
		env.EVMClient,
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

func CreateNodeTypeToNodeMap(cluster *test_env.ClCluster, nodesToCreate []VRFNodeType) map[VRFNodeType]*VRFNode {
	var nodesMap = make(map[VRFNodeType]*VRFNode)
	for i, nodeType := range nodesToCreate {
		nodesMap[nodeType] = &VRFNode{
			CLNode: cluster.Nodes[i],
		}
	}
	return nodesMap
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
