package vrfv2_actions

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_config"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var (
	ErrNodePrimaryKey          = "error getting node's primary ETH key"
	ErrNodeNewTxKey            = "error creating node's EVM transaction key"
	ErrCreatingProvingKeyHash  = "error creating a keyHash from the proving key"
	ErrRegisteringProvingKey   = "error registering a proving key on Coordinator contract"
	ErrRegisterProvingKey      = "error registering proving keys"
	ErrEncodingProvingKey      = "error encoding proving key"
	ErrCreatingVRFv2Key        = "error creating VRFv2 key"
	ErrDeployBlockHashStore    = "error deploying blockhash store"
	ErrDeployCoordinator       = "error deploying VRF CoordinatorV2"
	ErrAdvancedConsumer        = "error deploying VRFv2 Advanced Consumer"
	ErrABIEncodingFunding      = "error Abi encoding subscriptionID"
	ErrSendingLinkToken        = "error sending Link token"
	ErrCreatingVRFv2Job        = "error creating VRFv2 job"
	ErrParseJob                = "error parsing job definition"
	ErrDeployVRFV2Contracts    = "error deploying VRFV2 contracts"
	ErrSetVRFCoordinatorConfig = "error setting config for VRF Coordinator contract"
	ErrCreateVRFSubscription   = "error creating VRF Subscription"
	ErrAddConsumerToSub        = "error adding consumer to VRF Subscription"
	ErrFundSubWithLinkToken    = "error funding subscription with Link tokens"
	ErrCreateVRFV2Jobs         = "error creating VRF V2 Jobs"
	ErrGetPrimaryKey           = "error getting primary ETH key address"
	ErrRestartCLNode           = "error restarting CL node"
	ErrWaitTXsComplete         = "error waiting for TXs to complete"
	ErrRequestRandomness       = "error requesting randomness"

	ErrWaitRandomWordsRequestedEvent = "error waiting for RandomWordsRequested event"
	ErrWaitRandomWordsFulfilledEvent = "error waiting for RandomWordsFulfilled event"
	ErrDeployWrapper                 = "error deploying VRFV2PlusWrapper"
)

func DeployVRFV2Contracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	linkTokenContract contracts.LinkToken,
	linkEthFeedContract contracts.MockETHLINKFeed,
	consumerContractsAmount int,
) (*VRFV2Contracts, error) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployBlockHashStore, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkTokenContract.Address(), bhs.Address(), linkEthFeedContract.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployCoordinator, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	consumers, err := DeployVRFV2Consumers(contractDeployer, coordinator, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	return &VRFV2Contracts{coordinator, bhs, consumers}, nil
}

func DeployVRFV2Consumers(contractDeployer contracts.ContractDeployer, coordinator contracts.VRFCoordinatorV2, consumerContractsAmount int) ([]contracts.VRFv2LoadTestConsumer, error) {
	var consumers []contracts.VRFv2LoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFv2LoadTestConsumer(coordinator.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func DeployVRFV2WrapperConsumers(contractDeployer contracts.ContractDeployer, linkTokenAddress string, vrfV2Wrapper contracts.VRFV2Wrapper, consumerContractsAmount int) ([]contracts.VRFv2WrapperLoadTestConsumer, error) {
	var consumers []contracts.VRFv2WrapperLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFV2WrapperLoadTestConsumer(linkTokenAddress, vrfV2Wrapper.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func DeployVRFV2DirectFundingContracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	linkTokenAddress string,
	linkEthFeedAddress string,
	coordinator contracts.VRFCoordinatorV2,
	consumerContractsAmount int,
) (*VRFV2WrapperContracts, error) {
	vrfv2Wrapper, err := contractDeployer.DeployVRFV2Wrapper(linkTokenAddress, linkEthFeedAddress, coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployWrapper, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	consumers, err := DeployVRFV2WrapperConsumers(contractDeployer, linkTokenAddress, vrfv2Wrapper, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	return &VRFV2WrapperContracts{vrfv2Wrapper, consumers}, nil
}

func CreateVRFV2Job(
	chainlinkNode *client.ChainlinkClient,
	coordinatorAddress string,
	nativeTokenKeyAddresses []string,
	pubKeyCompressed string,
	chainID string,
	minIncomingConfirmations uint16,
) (*client.Job, error) {
	jobUUID := uuid.New()
	os := &client.VRFV2TxPipelineSpec{
		Address: coordinatorAddress,
	}
	ost, err := os.String()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrParseJob, err)
	}

	job, err := chainlinkNode.MustCreateJob(&client.VRFV2JobSpec{
		Name:                     fmt.Sprintf("vrf-v2-%s", jobUUID),
		CoordinatorAddress:       coordinatorAddress,
		FromAddresses:            nativeTokenKeyAddresses,
		EVMChainID:               chainID,
		MinIncomingConfirmations: int(minIncomingConfirmations),
		PublicKey:                pubKeyCompressed,
		ExternalJobID:            jobUUID.String(),
		ObservationSource:        ost,
		BatchFulfillmentEnabled:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrCreatingVRFv2Job, err)
	}

	return job, nil
}

func VRFV2RegisterProvingKey(
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2,
) (VRFV2EncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return VRFV2EncodedProvingKey{}, fmt.Errorf("%s, err %w", ErrEncodingProvingKey, err)
	}
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	if err != nil {
		return VRFV2EncodedProvingKey{}, fmt.Errorf("%s, err %w", ErrRegisterProvingKey, err)
	}
	return provingKey, nil
}

func FundVRFCoordinatorV2Subscription(
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	chainClient blockchain.EVMClient,
	subscriptionID uint64,
	linkFundingAmountJuels *big.Int,
) error {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subscriptionID)
	if err != nil {
		return fmt.Errorf("%s, err %w", ErrABIEncodingFunding, err)
	}
	_, err = linkToken.TransferAndCall(coordinator.Address(), linkFundingAmountJuels, encodedSubId)
	if err != nil {
		return fmt.Errorf("%s, err %w", ErrSendingLinkToken, err)
	}
	return chainClient.WaitForEvents()
}

// SetupVRFV2Environment will create specified number of subscriptions and add the same conumer/s to each of them
func SetupVRFV2Environment(
	env *test_env.CLClusterTestEnv,
	vrfv2Config vrfv2_config.VRFV2Config,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	registerProvingKeyAgainstAddress string,
	numberOfTxKeysToCreate int,
	numberOfConsumers int,
	numberOfSubToCreate int,
	l zerolog.Logger,
) (*VRFV2Contracts, []uint64, *VRFV2Data, error) {
	l.Info().Msg("Starting VRFV2 environment setup")
	l.Info().Msg("Deploying VRFV2 contracts")
	vrfv2Contracts, err := DeployVRFV2Contracts(
		env.ContractDeployer,
		env.EVMClient,
		linkToken,
		mockNativeLINKFeed,
		numberOfConsumers,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrDeployVRFV2Contracts, err)
	}
	vrfCoordinatorV2FeeConfig := vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
		FulfillmentFlatFeeLinkPPMTier1: vrfv2Config.FulfillmentFlatFeeLinkPPMTier1,
		FulfillmentFlatFeeLinkPPMTier2: vrfv2Config.FulfillmentFlatFeeLinkPPMTier2,
		FulfillmentFlatFeeLinkPPMTier3: vrfv2Config.FulfillmentFlatFeeLinkPPMTier3,
		FulfillmentFlatFeeLinkPPMTier4: vrfv2Config.FulfillmentFlatFeeLinkPPMTier4,
		FulfillmentFlatFeeLinkPPMTier5: vrfv2Config.FulfillmentFlatFeeLinkPPMTier5,
		ReqsForTier2:                   big.NewInt(vrfv2Config.ReqsForTier2),
		ReqsForTier3:                   big.NewInt(vrfv2Config.ReqsForTier3),
		ReqsForTier4:                   big.NewInt(vrfv2Config.ReqsForTier4),
		ReqsForTier5:                   big.NewInt(vrfv2Config.ReqsForTier5)}

	l.Info().Str("Coordinator", vrfv2Contracts.Coordinator.Address()).Msg("Setting Coordinator Config")
	err = vrfv2Contracts.Coordinator.SetConfig(
		vrfv2Config.MinimumConfirmations,
		vrfv2Config.MaxGasLimitCoordinatorConfig,
		vrfv2Config.StalenessSeconds,
		vrfv2Config.GasAfterPaymentCalculation,
		big.NewInt(vrfv2Config.LinkNativeFeedResponse),
		vrfCoordinatorV2FeeConfig,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrSetVRFCoordinatorConfig, err)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	l.Info().
		Str("Coordinator", vrfv2Contracts.Coordinator.Address()).
		Int("Number of Subs to create", numberOfSubToCreate).
		Msg("Creating and funding subscriptions, adding consumers")
	subIDs, err := CreateFundSubsAndAddConsumers(
		env,
		vrfv2Config,
		linkToken,
		vrfv2Contracts.Coordinator, vrfv2Contracts.LoadTestConsumers, numberOfSubToCreate)
	if err != nil {
		return nil, nil, nil, err
	}
	l.Info().Str("Node URL", env.ClCluster.NodeAPIs()[0].URL()).Msg("Creating VRF Key on the Node")
	vrfKey, err := env.ClCluster.NodeAPIs()[0].MustCreateVRFKey()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrCreatingVRFv2Key, err)
	}
	pubKeyCompressed := vrfKey.Data.ID

	l.Info().Str("Coordinator", vrfv2Contracts.Coordinator.Address()).Msg("Registering Proving Key")
	provingKey, err := VRFV2RegisterProvingKey(vrfKey, registerProvingKeyAgainstAddress, vrfv2Contracts.Coordinator)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrRegisteringProvingKey, err)
	}
	keyHash, err := vrfv2Contracts.Coordinator.HashOfKey(context.Background(), provingKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrCreatingProvingKeyHash, err)
	}

	chainID := env.EVMClient.GetChainID()
	newNativeTokenKeyAddresses, err := CreateAndFundSendingKeys(env, vrfv2Config, numberOfTxKeysToCreate, chainID)
	if err != nil {
		return nil, nil, nil, err
	}
	nativeTokenPrimaryKeyAddress, err := env.ClCluster.NodeAPIs()[0].PrimaryEthAddress()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrNodePrimaryKey, err)
	}
	allNativeTokenKeyAddresses := append(newNativeTokenKeyAddresses, nativeTokenPrimaryKeyAddress)

	l.Info().Msg("Creating VRFV2  Job")
	vrfV2job, err := CreateVRFV2Job(
		env.ClCluster.NodeAPIs()[0],
		vrfv2Contracts.Coordinator.Address(),
		allNativeTokenKeyAddresses,
		pubKeyCompressed,
		chainID.String(),
		vrfv2Config.MinimumConfirmations,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrCreateVRFV2Jobs, err)
	}

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	nodeConfig := node.NewConfig(env.ClCluster.Nodes[0].NodeConfig,
		node.WithVRFv2EVMEstimator(allNativeTokenKeyAddresses, vrfv2Config.CLNodeMaxGasPriceGWei),
	)
	l.Info().Msg("Restarting Node with new sending key PriceMax configuration")
	err = env.ClCluster.Nodes[0].Restart(nodeConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrRestartCLNode, err)
	}

	vrfv2KeyData := VRFV2KeyData{
		VRFKey:            vrfKey,
		EncodedProvingKey: provingKey,
		KeyHash:           keyHash,
	}

	data := VRFV2Data{
		vrfv2KeyData,
		vrfV2job,
		nativeTokenPrimaryKeyAddress,
		chainID,
	}

	l.Info().Msg("VRFV2  environment setup is finished")
	return vrfv2Contracts, subIDs, &data, nil
}

func SetupVRFV2WrapperEnvironment(
	env *test_env.CLClusterTestEnv,
	vrfv2Config vrfv2_config.VRFV2Config,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	coordinator contracts.VRFCoordinatorV2,
	keyHash [32]byte,
	wrapperConsumerContractsAmount int,
) (*VRFV2WrapperContracts, *uint64, error) {
	// Deploy VRF v2 direct funding contracts
	wrapperContracts, err := DeployVRFV2DirectFundingContracts(
		env.ContractDeployer,
		env.EVMClient,
		linkToken.Address(),
		mockNativeLINKFeed.Address(),
		coordinator,
		wrapperConsumerContractsAmount,
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	// Configure VRF v2 wrapper contract
	err = wrapperContracts.VRFV2Wrapper.SetConfig(
		vrfv2Config.WrapperGasOverhead,
		vrfv2Config.CoordinatorGasOverhead,
		vrfv2Config.WrapperPremiumPercentage,
		keyHash,
		vrfv2Config.WrapperMaxNumberOfWords,
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	// Fetch wrapper subscription ID
	wrapperSubID, err := wrapperContracts.VRFV2Wrapper.GetSubID(context.Background())
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	// Fund wrapper subscription
	err = FundSubscriptions(env, vrfv2Config, linkToken, coordinator, []uint64{wrapperSubID})
	if err != nil {
		return nil, nil, err
	}

	// Fund consumer with LINK
	err = linkToken.Transfer(
		wrapperContracts.LoadTestConsumers[0].Address(),
		big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(vrfv2Config.WrapperConsumerFundingAmountLink)),
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	return wrapperContracts, &wrapperSubID, nil
}

func CreateAndFundSendingKeys(env *test_env.CLClusterTestEnv, vrfv2Config vrfv2_config.VRFV2Config, numberOfNativeTokenAddressesToCreate int, chainID *big.Int) ([]string, error) {
	var newNativeTokenKeyAddresses []string
	for i := 0; i < numberOfNativeTokenAddressesToCreate; i++ {
		newTxKey, response, err := env.ClCluster.NodeAPIs()[0].CreateTxKey("evm", chainID.String())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrNodeNewTxKey, err)
		}
		if response.StatusCode != 200 {
			return nil, fmt.Errorf("error creating transaction key - response code, err %d", response.StatusCode)
		}
		newNativeTokenKeyAddresses = append(newNativeTokenKeyAddresses, newTxKey.Data.ID)
		err = actions.FundAddress(env.EVMClient, newTxKey.Data.ID, big.NewFloat(vrfv2Config.ChainlinkNodeFunding))
		if err != nil {
			return nil, err
		}
	}
	return newNativeTokenKeyAddresses, nil
}

func CreateFundSubsAndAddConsumers(
	env *test_env.CLClusterTestEnv,
	vrfv2Config vrfv2_config.VRFV2Config,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	consumers []contracts.VRFv2LoadTestConsumer,
	numberOfSubToCreate int,
) ([]uint64, error) {
	subIDs, err := CreateSubsAndFund(env, vrfv2Config, linkToken, coordinator, numberOfSubToCreate)
	if err != nil {
		return nil, err
	}
	subToConsumersMap := map[uint64][]contracts.VRFv2LoadTestConsumer{}

	//each subscription will have the same consumers
	for _, subID := range subIDs {
		subToConsumersMap[subID] = consumers
	}

	err = AddConsumersToSubs(
		subToConsumersMap,
		coordinator,
	)
	if err != nil {
		return nil, err
	}

	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	return subIDs, nil
}

func CreateSubsAndFund(
	env *test_env.CLClusterTestEnv,
	vrfv2Config vrfv2_config.VRFV2Config,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	subAmountToCreate int,
) ([]uint64, error) {
	subs, err := CreateSubs(env, coordinator, subAmountToCreate)
	if err != nil {
		return nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	err = FundSubscriptions(env, vrfv2Config, linkToken, coordinator, subs)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func CreateSubs(
	env *test_env.CLClusterTestEnv,
	coordinator contracts.VRFCoordinatorV2,
	subAmountToCreate int,
) ([]uint64, error) {
	var subIDArr []uint64

	for i := 0; i < subAmountToCreate; i++ {
		subID, err := CreateSubAndFindSubID(env, coordinator)
		if err != nil {
			return nil, err
		}
		subIDArr = append(subIDArr, subID)
	}
	return subIDArr, nil
}

func AddConsumersToSubs(
	subToConsumerMap map[uint64][]contracts.VRFv2LoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
) error {
	for subID, consumers := range subToConsumerMap {
		for _, consumer := range consumers {
			err := coordinator.AddConsumer(subID, consumer.Address())
			if err != nil {
				return fmt.Errorf("%s, err %w", ErrAddConsumerToSub, err)
			}
		}
	}
	return nil
}

func CreateSubAndFindSubID(env *test_env.CLClusterTestEnv, coordinator contracts.VRFCoordinatorV2) (uint64, error) {
	tx, err := coordinator.CreateSubscription()
	if err != nil {
		return 0, fmt.Errorf("%s, err %w", ErrCreateVRFSubscription, err)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return 0, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	receipt, err := env.EVMClient.GetTxReceipt(tx.Hash())
	if err != nil {
		return 0, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	//SubscriptionsCreated Log should be emitted with the subscription ID
	subID := receipt.Logs[0].Topics[1].Big().Uint64()

	return subID, nil
}

func FundSubscriptions(
	env *test_env.CLClusterTestEnv,
	vrfv2Config vrfv2_config.VRFV2Config,
	linkAddress contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	subIDs []uint64,
) error {
	for _, subID := range subIDs {
		//Link Billing
		amountJuels := conversions.EtherToWei(big.NewFloat(vrfv2Config.SubscriptionFundingAmountLink))
		err := FundVRFCoordinatorV2Subscription(linkAddress, coordinator, env.EVMClient, subID, amountJuels)
		if err != nil {
			return fmt.Errorf("%s, err %w", ErrFundSubWithLinkToken, err)
		}
	}
	err := env.EVMClient.WaitForEvents()
	if err != nil {
		return fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	return nil
}

func DirectFundingRequestRandomnessAndWaitForFulfillment(
	consumer contracts.VRFv2WrapperLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
	vrfv2Data *VRFV2Data,
	subID uint64,
	randomnessRequestCountPerRequest uint16,
	vrfv2Config vrfv2_config.VRFV2Config,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, error) {
	logRandRequest(consumer.Address(), coordinator.Address(), subID, vrfv2Config, l)
	_, err := consumer.RequestRandomness(
		vrfv2Config.MinimumConfirmations,
		vrfv2Config.CallbackGasLimit,
		vrfv2Config.NumberOfWords,
		randomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrRequestRandomness, err)
	}
	wrapperAddress, err := consumer.GetWrapper(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting wrapper address, err: %w", err)
	}
	fulfillmentEvents, err := WaitForRequestAndFulfillmentEvents(
		wrapperAddress.String(),
		coordinator,
		vrfv2Data,
		subID,
		randomWordsFulfilledEventTimeout,
		l,
	)
	return fulfillmentEvents, err
}

func RequestRandomnessAndWaitForFulfillment(
	consumer contracts.VRFv2LoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
	vrfv2Data *VRFV2Data,
	subID uint64,
	randomnessRequestCountPerRequest uint16,
	vrfv2Config vrfv2_config.VRFV2Config,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, error) {
	logRandRequest(consumer.Address(), coordinator.Address(), subID, vrfv2Config, l)
	_, err := consumer.RequestRandomness(
		vrfv2Data.KeyHash,
		subID,
		vrfv2Config.MinimumConfirmations,
		vrfv2Config.CallbackGasLimit,
		vrfv2Config.NumberOfWords,
		randomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrRequestRandomness, err)
	}

	fulfillmentEvents, err := WaitForRequestAndFulfillmentEvents(
		consumer.Address(),
		coordinator,
		vrfv2Data,
		subID,
		randomWordsFulfilledEventTimeout,
		l,
	)
	return fulfillmentEvents, err
}

func WaitForRequestAndFulfillmentEvents(
	consumerAddress string,
	coordinator contracts.VRFCoordinatorV2,
	vrfv2Data *VRFV2Data,
	subID uint64,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, error) {
	randomWordsRequestedEvent, err := coordinator.WaitForRandomWordsRequestedEvent(
		[][32]byte{vrfv2Data.KeyHash},
		[]uint64{subID},
		[]common.Address{common.HexToAddress(consumerAddress)},
		time.Minute*1,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitRandomWordsRequestedEvent, err)
	}

	LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent)
	randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
		[]*big.Int{randomWordsRequestedEvent.RequestId},
		randomWordsFulfilledEventTimeout,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitRandomWordsFulfilledEvent, err)
	}

	LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent)
	return randomWordsFulfilledEvent, err
}

func WaitForRequestCountEqualToFulfilmentCount(consumer contracts.VRFv2LoadTestConsumer, timeout time.Duration, wg *sync.WaitGroup) (*big.Int, *big.Int, error) {
	metricsChannel := make(chan *contracts.VRFLoadTestMetrics)
	metricsErrorChannel := make(chan error)

	testContext, testCancel := context.WithTimeout(context.Background(), timeout)
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
			go retrieveLoadTestMetrics(consumer, metricsChannel, metricsErrorChannel)
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
	consumer contracts.VRFv2LoadTestConsumer,
	metricsChannel chan *contracts.VRFLoadTestMetrics,
	metricsErrorChannel chan error,
) {
	metrics, err := consumer.GetLoadTestMetrics(context.Background())
	if err != nil {
		metricsErrorChannel <- err
	}
	metricsChannel <- metrics
}

func LogSubDetails(l zerolog.Logger, subscription vrf_coordinator_v2.GetSubscription, subID uint64, coordinator contracts.VRFCoordinatorV2) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Str("Link Balance", (*commonassets.Link)(subscription.Balance).Link()).
		Uint64("Subscription ID", subID).
		Str("Subscription Owner", subscription.Owner.String()).
		Interface("Subscription Consumers", subscription.Consumers).
		Msg("Subscription Data")
}

func LogRandomnessRequestedEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2,
	randomWordsRequestedEvent *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested,
) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Str("Request ID", randomWordsRequestedEvent.RequestId.String()).
		Uint64("Subscription ID", randomWordsRequestedEvent.SubId).
		Str("Sender Address", randomWordsRequestedEvent.Sender.String()).
		Interface("Keyhash", randomWordsRequestedEvent.KeyHash).
		Uint32("Callback Gas Limit", randomWordsRequestedEvent.CallbackGasLimit).
		Uint32("Number of Words", randomWordsRequestedEvent.NumWords).
		Uint16("Minimum Request Confirmations", randomWordsRequestedEvent.MinimumRequestConfirmations).
		Msg("RandomnessRequested Event")
}

func LogRandomWordsFulfilledEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2,
	randomWordsFulfilledEvent *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled,
) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Str("Total Payment", randomWordsFulfilledEvent.Payment.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Str("Request ID", randomWordsFulfilledEvent.RequestId.String()).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Msg("RandomWordsFulfilled Event (TX metadata)")
}

func logRandRequest(
	consumer string,
	coordinator string,
	subID uint64,
	vrfv2Config vrfv2_config.VRFV2Config,
	l zerolog.Logger,
) {
	l.Debug().
		Str("Consumer", consumer).
		Str("Coordinator", coordinator).
		Uint64("SubID", subID).
		Uint16("MinimumConfirmations", vrfv2Config.MinimumConfirmations).
		Uint32("CallbackGasLimit", vrfv2Config.CallbackGasLimit).
		Uint16("RandomnessRequestCountPerRequest", vrfv2Config.RandomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", vrfv2Config.RandomnessRequestCountPerRequestDeviation).
		Msg("Requesting randomness")
}
