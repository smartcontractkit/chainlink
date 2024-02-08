package common

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	testconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2"
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
		newNativeTokenKeyAddresses = append(newNativeTokenKeyAddresses, newTxKey.Data.ID)
		err = actions.FundAddress(client, newTxKey.Data.ID, big.NewFloat(chainlinkNodeFunding))
		if err != nil {
			return nil, err
		}
	}
	return newNativeTokenKeyAddresses, nil
}

func SetupBHSNode(
	env *test_env.CLClusterTestEnv,
	config *testconfig.General,
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
