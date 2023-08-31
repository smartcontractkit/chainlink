package vrfv2plus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_constants"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/utils"
	"math/big"
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
)

func DeployVRFV2PlusContracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
) (*VRFV2PlusContracts, error) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployBlockHashStore)
	}
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2Plus(bhs.Address())
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployCoordinator)
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

func CreateVRFV2PlusJobs(
	chainlinkNodes []*client.ChainlinkClient,
	coordinator contracts.VRFCoordinatorV2Plus,
	c blockchain.EVMClient,
	minIncomingConfirmations uint16,
) ([]*VRFV2PlusJobInfo, error) {
	jobInfo := make([]*VRFV2PlusJobInfo, 0)
	for _, chainlinkNode := range chainlinkNodes {
		vrfKey, err := chainlinkNode.MustCreateVRFKey()
		if err != nil {
			return nil, errors.Wrap(err, ErrCreatingVRFv2PlusKey)
		}
		pubKeyCompressed := vrfKey.Data.ID
		jobUUID := uuid.New()
		os := &client.VRFV2PlusTxPipelineSpec{
			Address: coordinator.Address(),
		}
		ost, err := os.String()
		if err != nil {
			return nil, errors.Wrap(err, ErrParseJob)
		}
		nativeTokenPrimaryKeyAddress, err := chainlinkNode.PrimaryEthAddress()
		if err != nil {
			return nil, errors.Wrap(err, ErrNodePrimaryKey)
		}
		job, err := chainlinkNode.MustCreateJob(&client.VRFV2PlusJobSpec{
			Name:                     fmt.Sprintf("vrf-v2-plus-%s", jobUUID),
			CoordinatorAddress:       coordinator.Address(),
			FromAddresses:            []string{nativeTokenPrimaryKeyAddress},
			EVMChainID:               c.GetChainID().String(),
			MinIncomingConfirmations: int(minIncomingConfirmations),
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
			BatchFulfillmentEnabled:  false,
		})
		if err != nil {
			return nil, errors.Wrap(err, ErrCreatingVRFv2PlusJob)
		}
		provingKey, err := VRFV2RegisterProvingKey(vrfKey, nativeTokenPrimaryKeyAddress, coordinator)
		if err != nil {
			return nil, errors.Wrap(err, ErrRegisteringProvingKey)
		}
		keyHash, err := coordinator.HashOfKey(context.Background(), provingKey)
		if err != nil {
			return nil, errors.Wrap(err, ErrCreatingProvingKeyHash)
		}
		ji := &VRFV2PlusJobInfo{
			Job:               job,
			VRFKey:            vrfKey,
			EncodedProvingKey: provingKey,
			KeyHash:           keyHash,
		}
		jobInfo = append(jobInfo, ji)
	}
	return jobInfo, nil
}

func VRFV2RegisterProvingKey(
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

func SetupVRFV2PlusEnvironment(
	env *test_env.CLClusterTestEnv,
	linkAddress contracts.LinkToken,
	mockETHLinkFeedAddress contracts.MockETHLINKFeed,
) (*test_env.CLClusterTestEnv, *VRFV2PlusContracts, *big.Int, *VRFV2PlusJobInfo, error) {

	vrfv2PlusContracts, err := DeployVRFV2PlusContracts(env.ContractDeployer, env.EVMClient)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrDeployVRFV2PlusContracts)
	}

	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrWaitTXsComplete)
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
		return nil, nil, nil, nil, errors.Wrap(err, ErrSetVRFCoordinatorConfig)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrWaitTXsComplete)
	}

	err = vrfv2PlusContracts.Coordinator.CreateSubscription()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrCreateVRFSubscription)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrWaitTXsComplete)
	}

	subID, err := vrfv2PlusContracts.Coordinator.FindSubscriptionID()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrFindSubID)
	}

	err = vrfv2PlusContracts.Coordinator.AddConsumer(subID, vrfv2PlusContracts.LoadTestConsumer.Address())
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrAddConsumerToSub)
	}

	//Native Billing
	err = vrfv2PlusContracts.Coordinator.FundSubscriptionWithEth(subID, big.NewInt(0).Mul(vrfv2plus_constants.VRFSubscriptionFundingAmountNativeToken, big.NewInt(1e18)))
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrFundSubWithNativeToken)
	}

	//Link Billing
	err = vrfv2PlusContracts.Coordinator.SetLINKAndLINKETHFeed(linkAddress.Address(), mockETHLinkFeedAddress.Address())
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrSetLinkETHLinkFeed)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrWaitTXsComplete)
	}
	err = FundVRFCoordinatorV2PlusSubscription(linkAddress, vrfv2PlusContracts.Coordinator, env.EVMClient, subID, vrfv2plus_constants.VRFSubscriptionFundingAmountLink)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrFundSubWithLinkToken)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrWaitTXsComplete)
	}

	vrfV2PlusJobs, err := CreateVRFV2PlusJobs(env.GetAPIs(), vrfv2PlusContracts.Coordinator, env.EVMClient, vrfv2plus_constants.MinimumConfirmations)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrCreateVRFV2PlusJobs)
	}

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	addr, err := env.CLNodes[0].API.PrimaryEthAddress()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrGetPrimaryKey)
	}
	nodeConfig := node.NewConfig(env.CLNodes[0].NodeConfig,
		node.WithVRFv2EVMEstimator(addr),
	)
	err = env.CLNodes[0].Restart(nodeConfig)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, ErrRestartCLNode)
	}
	return env, vrfv2PlusContracts, subID, vrfV2PlusJobs[0], nil
}
