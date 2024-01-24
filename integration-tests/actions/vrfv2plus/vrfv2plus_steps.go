package vrfv2plus

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper_load_test_consumer"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	vrfv2plus_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"
)

var (
	ErrNodePrimaryKey                              = "error getting node's primary ETH key"
	ErrNodeNewTxKey                                = "error creating node's EVM transaction key"
	ErrCreatingProvingKeyHash                      = "error creating a keyHash from the proving key"
	ErrRegisteringProvingKey                       = "error registering a proving key on Coordinator contract"
	ErrRegisterProvingKey                          = "error registering proving keys"
	ErrEncodingProvingKey                          = "error encoding proving key"
	ErrCreatingVRFv2PlusKey                        = "error creating VRFv2Plus key"
	ErrDeployBlockHashStore                        = "error deploying blockhash store"
	ErrDeployCoordinator                           = "error deploying VRF CoordinatorV2Plus"
	ErrAdvancedConsumer                            = "error deploying VRFv2Plus Advanced Consumer"
	ErrABIEncodingFunding                          = "error Abi encoding subscriptionID"
	ErrSendingLinkToken                            = "error sending Link token"
	ErrCreatingVRFv2PlusJob                        = "error creating VRFv2Plus job"
	ErrParseJob                                    = "error parsing job definition"
	ErrDeployVRFV2_5Contracts                      = "error deploying VRFV2_5 contracts"
	ErrSetVRFCoordinatorConfig                     = "error setting config for VRF Coordinator contract"
	ErrCreateVRFSubscription                       = "error creating VRF Subscription"
	ErrAddConsumerToSub                            = "error adding consumer to VRF Subscription"
	ErrFundSubWithNativeToken                      = "error funding subscription with native token"
	ErrSetLinkNativeLinkFeed                       = "error setting Link and ETH/LINK feed for VRF Coordinator contract"
	ErrFundSubWithLinkToken                        = "error funding subscription with Link tokens"
	ErrCreateVRFV2PlusJobs                         = "error creating VRF V2 Plus Jobs"
	ErrGetPrimaryKey                               = "error getting primary ETH key address"
	ErrRestartCLNode                               = "error restarting CL node"
	ErrWaitTXsComplete                             = "error waiting for TXs to complete"
	ErrRequestRandomness                           = "error requesting randomness"
	ErrRequestRandomnessDirectFundingLinkPayment   = "error requesting randomness with direct funding and link payment"
	ErrRequestRandomnessDirectFundingNativePayment = "error requesting randomness with direct funding and native payment"

	ErrWaitRandomWordsRequestedEvent = "error waiting for RandomWordsRequested event"
	ErrWaitRandomWordsFulfilledEvent = "error waiting for RandomWordsFulfilled event"
	ErrLinkTotalBalance              = "error waiting for RandomWordsFulfilled event"
	ErrNativeTokenBalance            = "error waiting for RandomWordsFulfilled event"
	ErrDeployWrapper                 = "error deploying VRFV2PlusWrapper"
)

type VRFJobSpecConfig struct {
	ForwardingAllowed             bool
	CoordinatorAddress            string
	FromAddresses                 []string
	EVMChainID                    string
	MinIncomingConfirmations      int
	PublicKey                     string
	BatchFulfillmentEnabled       bool
	BatchFulfillmentGasMultiplier float64
	EstimateGasMultiplier         float64
	PollPeriod                    time.Duration
	RequestTimeout                time.Duration
}

func DeployVRFV2_5Contracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	consumerContractsAmount int,
) (*VRFV2_5Contracts, error) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployBlockHashStore, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2_5(bhs.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployCoordinator, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	consumers, err := DeployVRFV2PlusConsumers(contractDeployer, coordinator, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	return &VRFV2_5Contracts{coordinator, bhs, consumers}, nil
}

func DeployVRFV2PlusConsumers(contractDeployer contracts.ContractDeployer, coordinator contracts.VRFCoordinatorV2_5, consumerContractsAmount int) ([]contracts.VRFv2PlusLoadTestConsumer, error) {
	var consumers []contracts.VRFv2PlusLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFv2PlusLoadTestConsumer(coordinator.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func CreateVRFV2PlusJob(
	chainlinkNode *client.ChainlinkClient,
	vrfJobSpecConfig VRFJobSpecConfig,
) (*client.Job, error) {
	jobUUID := uuid.New()
	os := &client.VRFV2PlusTxPipelineSpec{
		Address:               vrfJobSpecConfig.CoordinatorAddress,
		EstimateGasMultiplier: vrfJobSpecConfig.EstimateGasMultiplier,
		FromAddress:           vrfJobSpecConfig.FromAddresses[0],
	}
	ost, err := os.String()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrParseJob, err)
	}

	job, err := chainlinkNode.MustCreateJob(&client.VRFV2PlusJobSpec{
		Name:                          fmt.Sprintf("vrf-v2-plus-%s", jobUUID),
		CoordinatorAddress:            vrfJobSpecConfig.CoordinatorAddress,
		FromAddresses:                 vrfJobSpecConfig.FromAddresses,
		EVMChainID:                    vrfJobSpecConfig.EVMChainID,
		MinIncomingConfirmations:      vrfJobSpecConfig.MinIncomingConfirmations,
		PublicKey:                     vrfJobSpecConfig.PublicKey,
		ExternalJobID:                 jobUUID.String(),
		ObservationSource:             ost,
		BatchFulfillmentEnabled:       vrfJobSpecConfig.BatchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: vrfJobSpecConfig.BatchFulfillmentGasMultiplier,
		PollPeriod:                    vrfJobSpecConfig.PollPeriod,
		RequestTimeout:                vrfJobSpecConfig.RequestTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrCreatingVRFv2PlusJob, err)
	}

	return job, nil
}

func VRFV2_5RegisterProvingKey(
	vrfKey *client.VRFKey,
	coordinator contracts.VRFCoordinatorV2_5,
) (VRFV2PlusEncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return VRFV2PlusEncodedProvingKey{}, fmt.Errorf("%s, err %w", ErrEncodingProvingKey, err)
	}
	err = coordinator.RegisterProvingKey(
		provingKey,
	)
	if err != nil {
		return VRFV2PlusEncodedProvingKey{}, fmt.Errorf("%s, err %w", ErrRegisterProvingKey, err)
	}
	return provingKey, nil
}

func VRFV2PlusUpgradedVersionRegisterProvingKey(
	vrfKey *client.VRFKey,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
) (VRFV2PlusEncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return VRFV2PlusEncodedProvingKey{}, fmt.Errorf("%s, err %w", ErrEncodingProvingKey, err)
	}
	err = coordinator.RegisterProvingKey(
		provingKey,
	)
	if err != nil {
		return VRFV2PlusEncodedProvingKey{}, fmt.Errorf("%s, err %w", ErrRegisterProvingKey, err)
	}
	return provingKey, nil
}

func FundVRFCoordinatorV2_5Subscription(
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	chainClient blockchain.EVMClient,
	subscriptionID *big.Int,
	linkFundingAmountJuels *big.Int,
) error {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint256"}]`, subscriptionID)
	if err != nil {
		return fmt.Errorf("%s, err %w", ErrABIEncodingFunding, err)
	}
	_, err = linkToken.TransferAndCall(coordinator.Address(), linkFundingAmountJuels, encodedSubId)
	if err != nil {
		return fmt.Errorf("%s, err %w", ErrSendingLinkToken, err)
	}
	return chainClient.WaitForEvents()
}

// SetupVRFV2_5Environment will create specified number of subscriptions and add the same conumer/s to each of them
func SetupVRFV2_5Environment(
	env *test_env.CLClusterTestEnv,
	vrfv2PlusTestConfig types.VRFv2PlusTestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	numberOfTxKeysToCreate int,
	numberOfConsumers int,
	numberOfSubToCreate int,
	l zerolog.Logger,
) (*VRFV2_5Contracts, []*big.Int, *VRFV2PlusData, error) {
	l.Info().Msg("Starting VRFV2 Plus environment setup")
	l.Info().Msg("Deploying VRFV2 Plus contracts")
	vrfv2_5Contracts, err := DeployVRFV2_5Contracts(env.ContractDeployer, env.EVMClient, numberOfConsumers)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrDeployVRFV2_5Contracts, err)
	}

	l.Info().Str("Coordinator", vrfv2_5Contracts.Coordinator.Address()).Msg("Setting Coordinator Config")
	vrfv2PlusConfig := vrfv2PlusTestConfig.GetVRFv2PlusConfig().General
	err = vrfv2_5Contracts.Coordinator.SetConfig(
		*vrfv2PlusConfig.MinimumConfirmations,
		*vrfv2PlusConfig.MaxGasLimitCoordinatorConfig,
		*vrfv2PlusConfig.StalenessSeconds,
		*vrfv2PlusConfig.GasAfterPaymentCalculation,
		big.NewInt(*vrfv2PlusConfig.FallbackWeiPerUnitLink),
		vrf_coordinator_v2_5.VRFCoordinatorV25FeeConfig{
			FulfillmentFlatFeeLinkPPM:   *vrfv2PlusConfig.FulfillmentFlatFeeLinkPPM,
			FulfillmentFlatFeeNativePPM: *vrfv2PlusConfig.FulfillmentFlatFeeNativePPM,
		},
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrSetVRFCoordinatorConfig, err)
	}

	l.Info().Str("Coordinator", vrfv2_5Contracts.Coordinator.Address()).Msg("Setting Link and ETH/LINK feed")
	err = vrfv2_5Contracts.Coordinator.SetLINKAndLINKNativeFeed(linkToken.Address(), mockNativeLINKFeed.Address())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrSetLinkNativeLinkFeed, err)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	l.Info().
		Str("Coordinator", vrfv2_5Contracts.Coordinator.Address()).
		Int("Number of Subs to create", numberOfSubToCreate).
		Msg("Creating and funding subscriptions, adding consumers")
	subIDs, err := CreateFundSubsAndAddConsumers(
		env,
		big.NewFloat(*vrfv2PlusConfig.SubscriptionFundingAmountNative),
		big.NewFloat(*vrfv2PlusConfig.SubscriptionFundingAmountLink),
		linkToken,
		vrfv2_5Contracts.Coordinator, vrfv2_5Contracts.LoadTestConsumers,
		numberOfSubToCreate,
		vrfv2plus_config.BillingType(*vrfv2PlusConfig.SubscriptionBillingType))
	if err != nil {
		return nil, nil, nil, err
	}
	l.Info().Str("Node URL", env.ClCluster.NodeAPIs()[0].URL()).Msg("Creating VRF Key on the Node")
	vrfKey, err := env.ClCluster.NodeAPIs()[0].MustCreateVRFKey()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrCreatingVRFv2PlusKey, err)
	}
	pubKeyCompressed := vrfKey.Data.ID

	l.Info().Str("Coordinator", vrfv2_5Contracts.Coordinator.Address()).Msg("Registering Proving Key")
	provingKey, err := VRFV2_5RegisterProvingKey(vrfKey, vrfv2_5Contracts.Coordinator)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrRegisteringProvingKey, err)
	}
	keyHash, err := vrfv2_5Contracts.Coordinator.HashOfKey(context.Background(), provingKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrCreatingProvingKeyHash, err)
	}

	chainID := env.EVMClient.GetChainID()
	newNativeTokenKeyAddresses, err := CreateAndFundSendingKeys(env, vrfv2PlusTestConfig, numberOfTxKeysToCreate, chainID)
	if err != nil {
		return nil, nil, nil, err
	}
	nativeTokenPrimaryKeyAddress, err := env.ClCluster.NodeAPIs()[0].PrimaryEthAddress()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrNodePrimaryKey, err)
	}
	allNativeTokenKeyAddresses := append(newNativeTokenKeyAddresses, nativeTokenPrimaryKeyAddress)

	vrfJobSpecConfig := VRFJobSpecConfig{
		ForwardingAllowed:             false,
		CoordinatorAddress:            vrfv2_5Contracts.Coordinator.Address(),
		FromAddresses:                 allNativeTokenKeyAddresses,
		EVMChainID:                    chainID.String(),
		MinIncomingConfirmations:      int(*vrfv2PlusConfig.MinimumConfirmations),
		PublicKey:                     pubKeyCompressed,
		EstimateGasMultiplier:         1,
		BatchFulfillmentEnabled:       false,
		BatchFulfillmentGasMultiplier: 1.15,
		PollPeriod:                    time.Second * 1,
		RequestTimeout:                time.Hour * 24,
	}

	l.Info().Msg("Creating VRFV2 Plus Job")
	job, err := CreateVRFV2PlusJob(
		env.ClCluster.NodeAPIs()[0],
		vrfJobSpecConfig,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrCreateVRFV2PlusJobs, err)
	}

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	nodeConfig := node.NewConfig(env.ClCluster.Nodes[0].NodeConfig,
		node.WithLogPollInterval(1*time.Second),
		node.WithVRFv2EVMEstimator(allNativeTokenKeyAddresses, *vrfv2PlusConfig.CLNodeMaxGasPriceGWei),
	)
	l.Info().Msg("Restarting Node with new sending key PriceMax configuration and log poll period configuration")
	err = env.ClCluster.Nodes[0].Restart(nodeConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", ErrRestartCLNode, err)
	}

	vrfv2PlusKeyData := VRFV2PlusKeyData{
		VRFKey:            vrfKey,
		EncodedProvingKey: provingKey,
		KeyHash:           keyHash,
	}

	data := VRFV2PlusData{
		vrfv2PlusKeyData,
		job,
		nativeTokenPrimaryKeyAddress,
		chainID,
	}

	l.Info().Msg("VRFV2 Plus environment setup is finished")
	return vrfv2_5Contracts, subIDs, &data, nil
}

func CreateAndFundSendingKeys(env *test_env.CLClusterTestEnv, commonTestConfig tc.CommonTestConfig, numberOfNativeTokenAddressesToCreate int, chainID *big.Int) ([]string, error) {
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
		err = actions.FundAddress(env.EVMClient, newTxKey.Data.ID, big.NewFloat(*commonTestConfig.GetCommonConfig().ChainlinkNodeFunding))
		if err != nil {
			return nil, err
		}
	}
	return newNativeTokenKeyAddresses, nil
}

func CreateFundSubsAndAddConsumers(
	env *test_env.CLClusterTestEnv,
	subscriptionFundingAmountNative *big.Float,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	consumers []contracts.VRFv2PlusLoadTestConsumer,
	numberOfSubToCreate int,
	subscriptionBillingType vrfv2plus_config.BillingType,
) ([]*big.Int, error) {
	subIDs, err := CreateSubsAndFund(
		env,
		subscriptionFundingAmountNative,
		subscriptionFundingAmountLink,
		linkToken,
		coordinator,
		numberOfSubToCreate,
		subscriptionBillingType,
	)
	if err != nil {
		return nil, err
	}
	subToConsumersMap := map[*big.Int][]contracts.VRFv2PlusLoadTestConsumer{}

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
	subscriptionFundingAmountNative *big.Float,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	subAmountToCreate int,
	subscriptionBillingType vrfv2plus_config.BillingType,
) ([]*big.Int, error) {
	subs, err := CreateSubs(env, coordinator, subAmountToCreate)
	if err != nil {
		return nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	err = FundSubscriptions(
		env,
		subscriptionFundingAmountNative,
		subscriptionFundingAmountLink,
		linkToken,
		coordinator,
		subs,
		subscriptionBillingType,
	)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func CreateSubs(
	env *test_env.CLClusterTestEnv,
	coordinator contracts.VRFCoordinatorV2_5,
	subAmountToCreate int,
) ([]*big.Int, error) {
	var subIDArr []*big.Int

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
	subToConsumerMap map[*big.Int][]contracts.VRFv2PlusLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2_5,
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

func CreateSubAndFindSubID(env *test_env.CLClusterTestEnv, coordinator contracts.VRFCoordinatorV2_5) (*big.Int, error) {
	tx, err := coordinator.CreateSubscription()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrCreateVRFSubscription, err)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	receipt, err := env.EVMClient.GetTxReceipt(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	//SubscriptionsCreated Log should be emitted with the subscription ID
	subID := receipt.Logs[0].Topics[1].Big()

	return subID, nil
}

func FundSubscriptions(
	env *test_env.CLClusterTestEnv,
	subscriptionFundingAmountNative *big.Float,
	subscriptionFundingAmountLink *big.Float,
	linkAddress contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	subIDs []*big.Int,
	subscriptionBillingType vrfv2plus_config.BillingType,
) error {
	for _, subID := range subIDs {
		switch subscriptionBillingType {
		case vrfv2plus_config.BillingType_Native:
			//Native Billing
			amountWei := conversions.EtherToWei(subscriptionFundingAmountNative)
			err := coordinator.FundSubscriptionWithNative(
				subID,
				amountWei,
			)
			if err != nil {
				return fmt.Errorf("%s, err %w", ErrFundSubWithNativeToken, err)
			}
		case vrfv2plus_config.BillingType_Link:
			//Link Billing
			amountJuels := conversions.EtherToWei(subscriptionFundingAmountLink)
			err := FundVRFCoordinatorV2_5Subscription(linkAddress, coordinator, env.EVMClient, subID, amountJuels)
			if err != nil {
				return fmt.Errorf("%s, err %w", ErrFundSubWithLinkToken, err)
			}
		case vrfv2plus_config.BillingType_Link_and_Native:
			//Native Billing
			amountWei := conversions.EtherToWei(subscriptionFundingAmountNative)
			err := coordinator.FundSubscriptionWithNative(
				subID,
				amountWei,
			)
			if err != nil {
				return fmt.Errorf("%s, err %w", ErrFundSubWithNativeToken, err)
			}
			//Link Billing
			amountJuels := conversions.EtherToWei(subscriptionFundingAmountLink)
			err = FundVRFCoordinatorV2_5Subscription(linkAddress, coordinator, env.EVMClient, subID, amountJuels)
			if err != nil {
				return fmt.Errorf("%s, err %w", ErrFundSubWithLinkToken, err)
			}
		default:
			return fmt.Errorf("invalid billing type: %s", subscriptionBillingType)
		}
	}
	err := env.EVMClient.WaitForEvents()
	if err != nil {
		return fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	return nil
}

func GetUpgradedCoordinatorTotalBalance(coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion) (linkTotalBalance *big.Int, nativeTokenTotalBalance *big.Int, err error) {
	linkTotalBalance, err = coordinator.GetLinkTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrLinkTotalBalance, err)
	}
	nativeTokenTotalBalance, err = coordinator.GetNativeTokenTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrNativeTokenBalance, err)
	}
	return
}

func GetCoordinatorTotalBalance(coordinator contracts.VRFCoordinatorV2_5) (linkTotalBalance *big.Int, nativeTokenTotalBalance *big.Int, err error) {
	linkTotalBalance, err = coordinator.GetLinkTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrLinkTotalBalance, err)
	}
	nativeTokenTotalBalance, err = coordinator.GetNativeTokenTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrNativeTokenBalance, err)
	}
	return
}

func RequestRandomnessAndWaitForFulfillment(
	consumer contracts.VRFv2PlusLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2_5,
	vrfv2PlusData *VRFV2PlusData,
	subID *big.Int,
	isNativeBilling bool,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled, error) {
	logRandRequest(l, consumer.Address(), coordinator.Address(), subID, isNativeBilling, minimumConfirmations, callbackGasLimit, numberOfWords, randomnessRequestCountPerRequest, randomnessRequestCountPerRequestDeviation)
	_, err := consumer.RequestRandomness(
		vrfv2PlusData.KeyHash,
		subID,
		minimumConfirmations,
		callbackGasLimit,
		isNativeBilling,
		numberOfWords,
		randomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrRequestRandomness, err)
	}

	return WaitForRequestAndFulfillmentEvents(
		consumer.Address(),
		coordinator,
		vrfv2PlusData,
		subID,
		isNativeBilling,
		randomWordsFulfilledEventTimeout,
		l,
	)
}

func RequestRandomnessAndWaitForFulfillmentUpgraded(
	consumer contracts.VRFv2PlusLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	vrfv2PlusData *VRFV2PlusData,
	subID *big.Int,
	isNativeBilling bool,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	l zerolog.Logger,
) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, error) {
	logRandRequest(l, consumer.Address(), coordinator.Address(), subID, isNativeBilling, minimumConfirmations, callbackGasLimit, numberOfWords, randomnessRequestCountPerRequest, randomnessRequestCountPerRequestDeviation)
	_, err := consumer.RequestRandomness(
		vrfv2PlusData.KeyHash,
		subID,
		minimumConfirmations,
		callbackGasLimit,
		isNativeBilling,
		numberOfWords,
		randomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrRequestRandomness, err)
	}

	randomWordsRequestedEvent, err := coordinator.WaitForRandomWordsRequestedEvent(
		[][32]byte{vrfv2PlusData.KeyHash},
		[]*big.Int{subID},
		[]common.Address{common.HexToAddress(consumer.Address())},
		time.Minute*1,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitRandomWordsRequestedEvent, err)
	}

	LogRandomnessRequestedEventUpgraded(l, coordinator, randomWordsRequestedEvent)

	randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
		[]*big.Int{subID},
		[]*big.Int{randomWordsRequestedEvent.RequestId},
		time.Minute*2,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitRandomWordsFulfilledEvent, err)
	}
	LogRandomWordsFulfilledEventUpgraded(l, coordinator, randomWordsFulfilledEvent)

	return randomWordsFulfilledEvent, err
}

func SetupVRFV2PlusWrapperEnvironment(
	env *test_env.CLClusterTestEnv,
	vrfv2PlusTestConfig types.VRFv2PlusTestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	coordinator contracts.VRFCoordinatorV2_5,
	keyHash [32]byte,
	wrapperConsumerContractsAmount int,
) (*VRFV2PlusWrapperContracts, *big.Int, error) {
	vrfv2PlusConfig := vrfv2PlusTestConfig.GetVRFv2PlusConfig().General
	wrapperContracts, err := DeployVRFV2PlusDirectFundingContracts(
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
	err = wrapperContracts.VRFV2PlusWrapper.SetConfig(
		*vrfv2PlusConfig.WrapperGasOverhead,
		*vrfv2PlusConfig.CoordinatorGasOverhead,
		*vrfv2PlusConfig.WrapperPremiumPercentage,
		keyHash,
		*vrfv2PlusConfig.WrapperMaxNumberOfWords,
		*vrfv2PlusConfig.StalenessSeconds,
		big.NewInt(*vrfv2PlusConfig.FallbackWeiPerUnitLink),
		*vrfv2PlusConfig.FulfillmentFlatFeeLinkPPM,
		*vrfv2PlusConfig.FulfillmentFlatFeeNativePPM,
	)
	if err != nil {
		return nil, nil, err
	}

	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	//fund sub
	wrapperSubID, err := wrapperContracts.VRFV2PlusWrapper.GetSubID(context.Background())
	if err != nil {
		return nil, nil, err
	}

	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	err = FundSubscriptions(env, big.NewFloat(*vrfv2PlusTestConfig.GetVRFv2PlusConfig().General.SubscriptionFundingAmountNative), big.NewFloat(*vrfv2PlusTestConfig.GetVRFv2PlusConfig().General.SubscriptionFundingAmountLink), linkToken, coordinator, []*big.Int{wrapperSubID}, vrfv2plus_config.BillingType(*vrfv2PlusTestConfig.GetVRFv2PlusConfig().General.SubscriptionBillingType))
	if err != nil {
		return nil, nil, err
	}

	//fund consumer with Link
	err = linkToken.Transfer(
		wrapperContracts.LoadTestConsumers[0].Address(),
		big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(*vrfv2PlusConfig.WrapperConsumerFundingAmountLink)),
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	//fund consumer with Eth
	err = wrapperContracts.LoadTestConsumers[0].Fund(big.NewFloat(*vrfv2PlusConfig.WrapperConsumerFundingAmountNativeToken))
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	return wrapperContracts, wrapperSubID, nil
}

func DeployVRFV2PlusWrapperConsumers(contractDeployer contracts.ContractDeployer, linkTokenAddress string, vrfV2PlusWrapper contracts.VRFV2PlusWrapper, consumerContractsAmount int) ([]contracts.VRFv2PlusWrapperLoadTestConsumer, error) {
	var consumers []contracts.VRFv2PlusWrapperLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFV2PlusWrapperLoadTestConsumer(linkTokenAddress, vrfV2PlusWrapper.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func DeployVRFV2PlusDirectFundingContracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	linkTokenAddress string,
	linkEthFeedAddress string,
	coordinator contracts.VRFCoordinatorV2_5,
	consumerContractsAmount int,
) (*VRFV2PlusWrapperContracts, error) {

	vrfv2PlusWrapper, err := contractDeployer.DeployVRFV2PlusWrapper(linkTokenAddress, linkEthFeedAddress, coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployWrapper, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}

	consumers, err := DeployVRFV2PlusWrapperConsumers(contractDeployer, linkTokenAddress, vrfv2PlusWrapper, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitTXsComplete, err)
	}
	return &VRFV2PlusWrapperContracts{vrfv2PlusWrapper, consumers}, nil
}

func DirectFundingRequestRandomnessAndWaitForFulfillment(
	consumer contracts.VRFv2PlusWrapperLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2_5,
	vrfv2PlusData *VRFV2PlusData,
	subID *big.Int,
	isNativeBilling bool,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled, error) {
	logRandRequest(l, consumer.Address(), coordinator.Address(), subID, isNativeBilling, minimumConfirmations, callbackGasLimit, numberOfWords, randomnessRequestCountPerRequest, randomnessRequestCountPerRequestDeviation)
	if isNativeBilling {
		_, err := consumer.RequestRandomnessNative(
			minimumConfirmations,
			callbackGasLimit,
			numberOfWords,
			randomnessRequestCountPerRequest,
		)
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrRequestRandomnessDirectFundingNativePayment, err)
		}
	} else {
		_, err := consumer.RequestRandomness(
			minimumConfirmations,
			callbackGasLimit,
			numberOfWords,
			randomnessRequestCountPerRequest,
		)
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrRequestRandomnessDirectFundingLinkPayment, err)
		}
	}
	wrapperAddress, err := consumer.GetWrapper(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting wrapper address, err: %w", err)
	}
	return WaitForRequestAndFulfillmentEvents(
		wrapperAddress.String(),
		coordinator,
		vrfv2PlusData,
		subID,
		isNativeBilling,
		randomWordsFulfilledEventTimeout,
		l,
	)
}

func WaitForRequestAndFulfillmentEvents(
	consumerAddress string,
	coordinator contracts.VRFCoordinatorV2_5,
	vrfv2PlusData *VRFV2PlusData,
	subID *big.Int,
	isNativeBilling bool,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled, error) {
	randomWordsRequestedEvent, err := coordinator.WaitForRandomWordsRequestedEvent(
		[][32]byte{vrfv2PlusData.KeyHash},
		[]*big.Int{subID},
		[]common.Address{common.HexToAddress(consumerAddress)},
		time.Minute*1,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitRandomWordsRequestedEvent, err)
	}

	LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent, isNativeBilling)

	randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
		[]*big.Int{subID},
		[]*big.Int{randomWordsRequestedEvent.RequestId},
		randomWordsFulfilledEventTimeout,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrWaitRandomWordsFulfilledEvent, err)
	}

	LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent, isNativeBilling)
	return randomWordsFulfilledEvent, err
}

func WaitForRequestCountEqualToFulfilmentCount(consumer contracts.VRFv2PlusLoadTestConsumer, timeout time.Duration, wg *sync.WaitGroup) (*big.Int, *big.Int, error) {
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

func ReturnFundsForFulfilledRequests(client blockchain.EVMClient, coordinator contracts.VRFCoordinatorV2_5, l zerolog.Logger) error {
	linkTotalBalance, err := coordinator.GetLinkTotalBalance(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting LINK total balance, err: %w", err)
	}
	defaultWallet := client.GetDefaultWallet().Address()
	l.Info().
		Str("LINK amount", linkTotalBalance.String()).
		Str("Returning to", defaultWallet).
		Msg("Returning LINK for fulfilled requests")
	err = coordinator.Withdraw(
		common.HexToAddress(defaultWallet),
	)
	if err != nil {
		return fmt.Errorf("Error withdrawing LINK from coordinator to default wallet, err: %w", err)
	}
	nativeTotalBalance, err := coordinator.GetNativeTokenTotalBalance(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting NATIVE total balance, err: %w", err)
	}
	l.Info().
		Str("Native Token amount", nativeTotalBalance.String()).
		Str("Returning to", defaultWallet).
		Msg("Returning Native Token for fulfilled requests")
	err = coordinator.WithdrawNative(
		common.HexToAddress(defaultWallet),
	)
	if err != nil {
		return fmt.Errorf("Error withdrawing NATIVE from coordinator to default wallet, err: %w", err)
	}
	return nil
}

func retrieveLoadTestMetrics(
	consumer contracts.VRFv2PlusLoadTestConsumer,
	metricsChannel chan *contracts.VRFLoadTestMetrics,
	metricsErrorChannel chan error,
) {
	metrics, err := consumer.GetLoadTestMetrics(context.Background())
	if err != nil {
		metricsErrorChannel <- err
	}
	metricsChannel <- metrics
}

func LogSubDetails(l zerolog.Logger, subscription vrf_coordinator_v2_5.GetSubscription, subID *big.Int, coordinator contracts.VRFCoordinatorV2_5) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Str("Link Balance", (*commonassets.Link)(subscription.Balance).Link()).
		Str("Native Token Balance", assets.FormatWei(subscription.NativeBalance)).
		Str("Subscription ID", subID.String()).
		Str("Subscription Owner", subscription.Owner.String()).
		Interface("Subscription Consumers", subscription.Consumers).
		Msg("Subscription Data")
}

func LogRandomnessRequestedEventUpgraded(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	randomWordsRequestedEvent *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested,
) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Str("Request ID", randomWordsRequestedEvent.RequestId.String()).
		Str("Subscription ID", randomWordsRequestedEvent.SubId.String()).
		Str("Sender Address", randomWordsRequestedEvent.Sender.String()).
		Interface("Keyhash", randomWordsRequestedEvent.KeyHash).
		Uint32("Callback Gas Limit", randomWordsRequestedEvent.CallbackGasLimit).
		Uint32("Number of Words", randomWordsRequestedEvent.NumWords).
		Uint16("Minimum Request Confirmations", randomWordsRequestedEvent.MinimumRequestConfirmations).
		Msg("RandomnessRequested Event")
}

func LogRandomWordsFulfilledEventUpgraded(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	randomWordsFulfilledEvent *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled,
) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Str("Total Payment in Juels", randomWordsFulfilledEvent.Payment.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Str("Subscription ID", randomWordsFulfilledEvent.SubID.String()).
		Str("Request ID", randomWordsFulfilledEvent.RequestId.String()).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Msg("RandomWordsFulfilled Event (TX metadata)")
}

func LogRandomnessRequestedEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2_5,
	randomWordsRequestedEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested,
	isNativeBilling bool,
) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Bool("Native Billing", isNativeBilling).
		Str("Request ID", randomWordsRequestedEvent.RequestId.String()).
		Str("Subscription ID", randomWordsRequestedEvent.SubId.String()).
		Str("Sender Address", randomWordsRequestedEvent.Sender.String()).
		Interface("Keyhash", randomWordsRequestedEvent.KeyHash).
		Uint32("Callback Gas Limit", randomWordsRequestedEvent.CallbackGasLimit).
		Uint32("Number of Words", randomWordsRequestedEvent.NumWords).
		Uint16("Minimum Request Confirmations", randomWordsRequestedEvent.MinimumRequestConfirmations).
		Msg("RandomnessRequested Event")
}

func LogRandomWordsFulfilledEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2_5,
	randomWordsFulfilledEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled,
	isNativeBilling bool,
) {
	l.Debug().
		Bool("Native Billing", isNativeBilling).
		Str("Coordinator", coordinator.Address()).
		Str("Total Payment", randomWordsFulfilledEvent.Payment.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Str("Subscription ID", randomWordsFulfilledEvent.SubId.String()).
		Str("Request ID", randomWordsFulfilledEvent.RequestId.String()).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Msg("RandomWordsFulfilled Event (TX metadata)")
}

func LogMigrationCompletedEvent(l zerolog.Logger, migrationCompletedEvent *vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted, vrfv2PlusContracts *VRFV2_5Contracts) {
	l.Debug().
		Str("Subscription ID", migrationCompletedEvent.SubId.String()).
		Str("Migrated From Coordinator", vrfv2PlusContracts.Coordinator.Address()).
		Str("Migrated To Coordinator", migrationCompletedEvent.NewCoordinator.String()).
		Msg("MigrationCompleted Event")
}

func LogSubDetailsAfterMigration(l zerolog.Logger, newCoordinator contracts.VRFCoordinatorV2PlusUpgradedVersion, subID *big.Int, migratedSubscription vrf_v2plus_upgraded_version.GetSubscription) {
	l.Debug().
		Str("New Coordinator", newCoordinator.Address()).
		Str("Subscription ID", subID.String()).
		Str("Juels Balance", migratedSubscription.Balance.String()).
		Str("Native Token Balance", migratedSubscription.NativeBalance.String()).
		Str("Subscription Owner", migratedSubscription.Owner.String()).
		Interface("Subscription Consumers", migratedSubscription.Consumers).
		Msg("Subscription Data After Migration to New Coordinator")
}

func LogFulfillmentDetailsLinkBilling(
	l zerolog.Logger,
	wrapperConsumerJuelsBalanceBeforeRequest *big.Int,
	wrapperConsumerJuelsBalanceAfterRequest *big.Int,
	consumerStatus vrfv2plus_wrapper_load_test_consumer.GetRequestStatus,
	randomWordsFulfilledEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled,
) {
	l.Debug().
		Str("Consumer Balance Before Request (Link)", (*commonassets.Link)(wrapperConsumerJuelsBalanceBeforeRequest).Link()).
		Str("Consumer Balance After Request (Link)", (*commonassets.Link)(wrapperConsumerJuelsBalanceAfterRequest).Link()).
		Bool("Fulfilment Status", consumerStatus.Fulfilled).
		Str("Paid by Consumer Contract (Link)", (*commonassets.Link)(consumerStatus.Paid).Link()).
		Str("Paid by Coordinator Sub (Link)", (*commonassets.Link)(randomWordsFulfilledEvent.Payment).Link()).
		Str("RequestTimestamp", consumerStatus.RequestTimestamp.String()).
		Str("FulfilmentTimestamp", consumerStatus.FulfilmentTimestamp.String()).
		Str("RequestBlockNumber", consumerStatus.RequestBlockNumber.String()).
		Str("FulfilmentBlockNumber", consumerStatus.FulfilmentBlockNumber.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Msg("Random Words Fulfilment Details For Link Billing")
}

func LogFulfillmentDetailsNativeBilling(
	l zerolog.Logger,
	wrapperConsumerBalanceBeforeRequestWei *big.Int,
	wrapperConsumerBalanceAfterRequestWei *big.Int,
	consumerStatus vrfv2plus_wrapper_load_test_consumer.GetRequestStatus,
	randomWordsFulfilledEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled,
) {
	l.Debug().
		Str("Consumer Balance Before Request", assets.FormatWei(wrapperConsumerBalanceBeforeRequestWei)).
		Str("Consumer Balance After Request", assets.FormatWei(wrapperConsumerBalanceAfterRequestWei)).
		Bool("Fulfilment Status", consumerStatus.Fulfilled).
		Str("Paid by Consumer Contract", assets.FormatWei(consumerStatus.Paid)).
		Str("Paid by Coordinator Sub", assets.FormatWei(randomWordsFulfilledEvent.Payment)).
		Str("RequestTimestamp", consumerStatus.RequestTimestamp.String()).
		Str("FulfilmentTimestamp", consumerStatus.FulfilmentTimestamp.String()).
		Str("RequestBlockNumber", consumerStatus.RequestBlockNumber.String()).
		Str("FulfilmentBlockNumber", consumerStatus.FulfilmentBlockNumber.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Msg("Random Words Request Fulfilment Details For Native Billing")
}

func logRandRequest(
	l zerolog.Logger,
	consumer string,
	coordinator string,
	subID *big.Int,
	isNativeBilling bool,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16) {
	l.Debug().
		Str("Consumer", consumer).
		Str("Coordinator", coordinator).
		Str("SubID", subID.String()).
		Bool("IsNativePayment", isNativeBilling).
		Uint16("MinimumConfirmations", minimumConfirmations).
		Uint32("CallbackGasLimit", callbackGasLimit).
		Uint32("NumberOfWords", numberOfWords).
		Uint16("RandomnessRequestCountPerRequest", randomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", randomnessRequestCountPerRequestDeviation).
		Msg("Requesting randomness")
}
