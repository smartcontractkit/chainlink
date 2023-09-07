package vrfv2plus

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_constants"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/utils"
	"math/big"
	"time"
)

var (
	ErrNodePrimaryKey         = "error getting node's primary ETH key"
	ErrCreatingProvingKeyHash = "error creating a keyHash from the proving key"
	ErrRegisteringProvingKey  = "error registering a proving key on Coordinator contract"
	ErrRegisterProvingKey     = "error registering proving keys"
	ErrEncodingProvingKey     = "error encoding proving key"
	ErrCreatingVRFv2PlusKey   = "error creating VRFv2Plus key"
	ErrDeployBlockHashStore   = "error deploying blockhash store"
	ErrDeployCoordinator      = "error deploying VRF CoordinatorV2Plus"
	ErrAdvancedConsumer       = "error deploying VRFv2Plus Advanced Consumer"
	ErrABIEncodingFunding     = "error Abi encoding subscriptionID"
	ErrSendingLinkToken       = "error sending Link token"
	ErrCreatingVRFv2PlusJob   = "error creating VRFv2Plus job"
	ErrParseJob               = "error parsing job definition"

	ErrDeployVRFV2PlusContracts = "error deploying VRFV2Plus contracts"
	ErrSetVRFCoordinatorConfig  = "error setting config for VRF Coordinator contract"
	ErrCreateVRFSubscription    = "error creating VRF Subscription"
	ErrFindSubID                = "error finding created subscription ID"
	ErrAddConsumerToSub         = "error adding consumer to VRF Subscription"
	ErrFundSubWithNativeToken   = "error funding subscription with native token"
	ErrSetLinkETHLinkFeed       = "error setting Link and ETH/LINK feed for VRF Coordinator contract"
	ErrFundSubWithLinkToken     = "error funding subscription with Link tokens"
	ErrCreateVRFV2PlusJobs      = "error creating VRF V2 Plus Jobs"
	ErrGetPrimaryKey            = "error getting primary ETH key address"
	ErrRestartCLNode            = "error restarting CL node"
	ErrWaitTXsComplete          = "error waiting for TXs to complete"

	ErrRequestRandomness             = "error requesting randomness"
	ErrWaitRandomWordsRequestedEvent = "error waiting for RandomWordsRequested event"

	ErrWaitRandomWordsFulfilledEvent = "error waiting for RandomWordsFulfilled event"
)

func DeployVRFV2PlusContracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
) (*VRFV2PlusContracts, error) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployBlockHashStore)
	}
	//err = chainClient.WaitForEvents()
	if err != nil {
		return nil, errors.Wrap(err, ErrWaitTXsComplete)
	}
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2Plus(bhs.Address())
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployCoordinator)
	}
	//err = chainClient.WaitForEvents()
	if err != nil {
		return nil, errors.Wrap(err, ErrWaitTXsComplete)
	}
	loadTestConsumer, err := contractDeployer.DeployVRFv2PlusLoadTestConsumer(coordinator.Address())
	if err != nil {
		return nil, errors.Wrap(err, ErrAdvancedConsumer)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, errors.Wrap(err, ErrWaitTXsComplete)
	}
	return &VRFV2PlusContracts{coordinator, bhs, loadTestConsumer}, nil
}

func CreateVRFV2PlusJob(
	chainlinkNode *client.ChainlinkClient,
	coordinatorAddress string,
	nativeTokenPrimaryKeyAddress string,
	pubKeyCompressed string,
	chainID string,
	minIncomingConfirmations uint16,
) (*client.Job, error) {
	jobUUID := uuid.New()
	os := &client.VRFV2PlusTxPipelineSpec{
		Address: coordinatorAddress,
	}
	ost, err := os.String()
	if err != nil {
		return nil, errors.Wrap(err, ErrParseJob)
	}

	job, err := chainlinkNode.MustCreateJob(&client.VRFV2PlusJobSpec{
		Name:                     fmt.Sprintf("vrf-v2-plus-%s", jobUUID),
		CoordinatorAddress:       coordinatorAddress,
		FromAddresses:            []string{nativeTokenPrimaryKeyAddress},
		EVMChainID:               chainID,
		MinIncomingConfirmations: int(minIncomingConfirmations),
		PublicKey:                pubKeyCompressed,
		ExternalJobID:            jobUUID.String(),
		ObservationSource:        ost,
		BatchFulfillmentEnabled:  false,
	})
	if err != nil {
		return nil, errors.Wrap(err, ErrCreatingVRFv2PlusJob)
	}

	return job, nil
}

func VRFV2PlusRegisterProvingKey(
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2Plus,
) (VRFV2PlusEncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return VRFV2PlusEncodedProvingKey{}, errors.Wrap(err, ErrEncodingProvingKey)
	}
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	if err != nil {
		return VRFV2PlusEncodedProvingKey{}, errors.Wrap(err, ErrRegisterProvingKey)
	}
	return provingKey, nil
}

func VRFV2PlusUpgradedVersionRegisterProvingKey(
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
) (VRFV2PlusEncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return VRFV2PlusEncodedProvingKey{}, errors.Wrap(err, ErrEncodingProvingKey)
	}
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	if err != nil {
		return VRFV2PlusEncodedProvingKey{}, errors.Wrap(err, ErrRegisterProvingKey)
	}
	return provingKey, nil
}

func FundVRFCoordinatorV2PlusSubscription(linkToken contracts.LinkToken, coordinator contracts.VRFCoordinatorV2Plus, chainClient blockchain.EVMClient, subscriptionID *big.Int, linkFundingAmount *big.Int) error {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint256"}]`, subscriptionID)
	if err != nil {
		return errors.Wrap(err, ErrABIEncodingFunding)
	}
	_, err = linkToken.TransferAndCall(coordinator.Address(), big.NewInt(0).Mul(linkFundingAmount, big.NewInt(1e18)), encodedSubId)
	if err != nil {
		return errors.Wrap(err, ErrSendingLinkToken)
	}
	return chainClient.WaitForEvents()
}

func SetupVRFV2PlusEnvironment(env *test_env.CLClusterTestEnv, linkAddress contracts.LinkToken, mockETHLinkFeedAddress contracts.MockETHLINKFeed) (*VRFV2PlusContracts, *big.Int, *VRFV2PlusData, error) {

	vrfv2PlusContracts, err := DeployVRFV2PlusContracts(env.ContractDeployer, env.EVMClient)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrDeployVRFV2PlusContracts)
	}

	err = vrfv2PlusContracts.Coordinator.SetConfig(
		vrfv2plus_constants.MinimumConfirmations,
		vrfv2plus_constants.MaxGasLimitVRFCoordinatorConfig,
		vrfv2plus_constants.StalenessSeconds,
		vrfv2plus_constants.GasAfterPaymentCalculation,
		vrfv2plus_constants.LinkEthFeedResponse,
		vrfv2plus_constants.VRFCoordinatorV2PlusFeeConfig,
	)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrSetVRFCoordinatorConfig)
	}

	err = vrfv2PlusContracts.Coordinator.CreateSubscription()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrCreateVRFSubscription)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrWaitTXsComplete)
	}

	subID, err := vrfv2PlusContracts.Coordinator.FindSubscriptionID()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrFindSubID)
	}

	err = vrfv2PlusContracts.Coordinator.AddConsumer(subID, vrfv2PlusContracts.LoadTestConsumer.Address())
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrAddConsumerToSub)
	}

	err = SetupBilling(env, linkAddress, mockETHLinkFeedAddress, vrfv2PlusContracts, subID)
	if err != nil {
		return nil, nil, nil, err
	}

	vrfKey, err := env.GetAPIs()[0].MustCreateVRFKey()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrCreatingVRFv2PlusKey)
	}
	pubKeyCompressed := vrfKey.Data.ID

	nativeTokenPrimaryKeyAddress, err := env.GetAPIs()[0].PrimaryEthAddress()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrNodePrimaryKey)
	}
	provingKey, err := VRFV2PlusRegisterProvingKey(vrfKey, nativeTokenPrimaryKeyAddress, vrfv2PlusContracts.Coordinator)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrRegisteringProvingKey)
	}
	keyHash, err := vrfv2PlusContracts.Coordinator.HashOfKey(context.Background(), provingKey)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrCreatingProvingKeyHash)
	}

	chainID := env.EVMClient.GetChainID()

	job, err := CreateVRFV2PlusJob(
		env.GetAPIs()[0],
		vrfv2PlusContracts.Coordinator.Address(),
		nativeTokenPrimaryKeyAddress,
		pubKeyCompressed,
		chainID.String(),
		vrfv2plus_constants.MinimumConfirmations,
	)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrCreateVRFV2PlusJobs)
	}

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	addr, err := env.CLNodes[0].API.PrimaryEthAddress()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrGetPrimaryKey)
	}
	nodeConfig := node.NewConfig(env.CLNodes[0].NodeConfig,
		node.WithVRFv2EVMEstimator(addr),
	)
	err = env.CLNodes[0].Restart(nodeConfig)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, ErrRestartCLNode)
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

	return vrfv2PlusContracts, subID, &data, nil
}

func SetupBilling(env *test_env.CLClusterTestEnv, linkAddress contracts.LinkToken, mockETHLinkFeedAddress contracts.MockETHLINKFeed, vrfv2PlusContracts *VRFV2PlusContracts, subID *big.Int) error {
	//Native Billing
	err := vrfv2PlusContracts.Coordinator.FundSubscriptionWithEth(subID, big.NewInt(0).Mul(vrfv2plus_constants.VRFSubscriptionFundingAmountNativeToken, big.NewInt(1e18)))
	if err != nil {
		return errors.Wrap(err, ErrFundSubWithNativeToken)
	}

	//Link Billing
	err = vrfv2PlusContracts.Coordinator.SetLINKAndLINKETHFeed(linkAddress.Address(), mockETHLinkFeedAddress.Address())
	if err != nil {
		return errors.Wrap(err, ErrSetLinkETHLinkFeed)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return errors.Wrap(err, ErrWaitTXsComplete)
	}

	err = FundVRFCoordinatorV2PlusSubscription(linkAddress, vrfv2PlusContracts.Coordinator, env.EVMClient, subID, vrfv2plus_constants.VRFSubscriptionFundingAmountLink)
	if err != nil {
		return errors.Wrap(err, ErrFundSubWithLinkToken)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return errors.Wrap(err, ErrWaitTXsComplete)
	}

	return nil
}

func RequestRandomnessAndWaitForFulfillment(
	consumer contracts.VRFv2PlusLoadTestConsumer,
	//coordinator contracts.VRFCoordinatorV2Plus,
	coordinator contracts.VRFCoordinatorV2Plus,
	vrfv2PlusData *VRFV2PlusData,
	subID *big.Int,
	isNativeBilling bool,
	l zerolog.Logger,
) (*vrf_coordinator_v2plus.VRFCoordinatorV2PlusRandomWordsFulfilled, error) {
	_, err := consumer.RequestRandomness(
		vrfv2PlusData.KeyHash,
		subID,
		vrfv2plus_constants.MinimumConfirmations,
		vrfv2plus_constants.CallbackGasLimit,
		isNativeBilling,
		vrfv2plus_constants.NumberOfWords,
		vrfv2plus_constants.RandomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrRequestRandomness)
	}

	randomWordsRequestedEvent, err := coordinator.WaitForRandomWordsRequestedEvent(
		[][32]byte{vrfv2PlusData.KeyHash},
		[]*big.Int{subID},
		[]common.Address{common.HexToAddress(consumer.Address())},
		time.Minute*1,
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrWaitRandomWordsRequestedEvent)
	}

	l.Debug().
		Interface("Request ID", randomWordsRequestedEvent.RequestId).
		Interface("Subscription ID", randomWordsRequestedEvent.SubId).
		Interface("Sender Address", randomWordsRequestedEvent.Sender.String()).
		Interface("Keyhash", randomWordsRequestedEvent.KeyHash).
		Interface("Callback Gas Limit", randomWordsRequestedEvent.CallbackGasLimit).
		Interface("Number of Words", randomWordsRequestedEvent.NumWords).
		Interface("Minimum Request Confirmations", randomWordsRequestedEvent.MinimumRequestConfirmations).
		Msg("RandomnessRequested Event")

	randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
		[]*big.Int{subID},
		[]*big.Int{randomWordsRequestedEvent.RequestId},
		time.Minute*2,
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrWaitRandomWordsFulfilledEvent)
	}

	l.Debug().
		Interface("Total Payment in Juels", randomWordsFulfilledEvent.Payment).
		Interface("TX Hash", randomWordsFulfilledEvent.Raw.TxHash).
		Interface("Subscription ID", randomWordsFulfilledEvent.SubID).
		Interface("Request ID", randomWordsFulfilledEvent.RequestId).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Msg("RandomWordsFulfilled Event (TX metadata)")
	return randomWordsFulfilledEvent, err
}
