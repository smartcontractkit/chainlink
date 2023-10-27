package vrfv2_actions

import (
	"context"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	vrfConst "github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_constants"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
)

var (
	ErrNodePrimaryKey         = "error getting node's primary ETH key"
	ErrCreatingProvingKeyHash = "error creating a keyHash from the proving key"
	ErrCreatingProvingKey     = "error creating a keyHash from the proving key"
	ErrRegisterProvingKey     = "error registering proving keys"
	ErrEncodingProvingKey     = "error encoding proving key"
	ErrCreatingVRFv2Key       = "error creating VRFv2 key"
	ErrDeployBlockHashStore   = "error deploying blockhash store"
	ErrDeployCoordinator      = "error deploying VRFv2 CoordinatorV2"
	ErrAdvancedConsumer       = "error deploying VRFv2 Advanced Consumer"
	ErrABIEncodingFunding     = "error Abi encoding subscriptionID"
	ErrSendingLinkToken       = "error sending Link token"
	ErrCreatingVRFv2Job       = "error creating VRFv2 job"
	ErrParseJob               = "error parsing job definition"
)

func DeployVRFV2Contracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	linkTokenContract contracts.LinkToken,
	linkEthFeedContract contracts.MockETHLINKFeed,
) (*VRFV2Contracts, error) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployBlockHashStore)
	}
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkTokenContract.Address(), bhs.Address(), linkEthFeedContract.Address())
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployCoordinator)
	}
	loadTestConsumer, err := contractDeployer.DeployVRFv2LoadTestConsumer(coordinator.Address())
	if err != nil {
		return nil, errors.Wrap(err, ErrAdvancedConsumer)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, err
	}
	return &VRFV2Contracts{coordinator, bhs, loadTestConsumer}, nil
}

func CreateVRFV2Jobs(
	chainlinkNodes []*client.ChainlinkClient,
	coordinator contracts.VRFCoordinatorV2,
	c blockchain.EVMClient,
	minIncomingConfirmations uint16,
) ([]VRFV2JobInfo, error) {
	jobInfo := make([]VRFV2JobInfo, 0)
	for _, chainlinkNode := range chainlinkNodes {
		vrfKey, err := chainlinkNode.MustCreateVRFKey()
		if err != nil {
			return nil, errors.Wrap(err, ErrCreatingVRFv2Key)
		}
		pubKeyCompressed := vrfKey.Data.ID
		jobUUID := uuid.New()
		os := &client.VRFV2TxPipelineSpec{
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
		job, err := chainlinkNode.MustCreateJob(&client.VRFV2JobSpec{
			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
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
			return nil, errors.Wrap(err, ErrCreatingVRFv2Job)
		}
		provingKey, err := VRFV2RegisterProvingKey(vrfKey, nativeTokenPrimaryKeyAddress, coordinator)
		if err != nil {
			return nil, errors.Wrap(err, ErrCreatingProvingKey)
		}
		keyHash, err := coordinator.HashOfKey(context.Background(), provingKey)
		if err != nil {
			return nil, errors.Wrap(err, ErrCreatingProvingKeyHash)
		}
		ji := VRFV2JobInfo{
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
	coordinator contracts.VRFCoordinatorV2,
) (VRFV2EncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return VRFV2EncodedProvingKey{}, errors.Wrap(err, ErrEncodingProvingKey)
	}
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	if err != nil {
		return VRFV2EncodedProvingKey{}, errors.Wrap(err, ErrRegisterProvingKey)
	}
	return provingKey, nil
}

func FundVRFCoordinatorV2Subscription(linkToken contracts.LinkToken, coordinator contracts.VRFCoordinatorV2, chainClient blockchain.EVMClient, subscriptionID uint64, linkFundingAmount *big.Int) error {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subscriptionID)
	if err != nil {
		return errors.Wrap(err, ErrABIEncodingFunding)
	}
	_, err = linkToken.TransferAndCall(coordinator.Address(), big.NewInt(0).Mul(linkFundingAmount, big.NewInt(1e18)), encodedSubId)
	if err != nil {
		return errors.Wrap(err, ErrSendingLinkToken)
	}
	return chainClient.WaitForEvents()
}

/* setup for load tests */

func SetupLocalLoadTestEnv(nodesFunding *big.Float, subFundingLINK *big.Int) (*test_env.CLClusterTestEnv, *VRFV2Contracts, [32]byte, error) {
	env, err := test_env.NewCLTestEnvBuilder().
		WithGeth().
		WithLogWatcher().
		WithMockAdapter().
		WithCLNodes(1).
		WithFunding(nodesFunding).
		WithLogWatcher().
		Build()
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	env.ParallelTransactions(true)

	mockFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, vrfConst.LinkEthFeedResponse)
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	lt, err := actions.DeployLINKToken(env.ContractDeployer)
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	vrfv2Contracts, err := DeployVRFV2Contracts(env.ContractDeployer, env.EVMClient, lt, mockFeed)
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	err = vrfv2Contracts.Coordinator.SetConfig(
		vrfConst.MinimumConfirmations,
		vrfConst.MaxGasLimitVRFCoordinatorConfig,
		vrfConst.StalenessSeconds,
		vrfConst.GasAfterPaymentCalculation,
		vrfConst.LinkEthFeedResponse,
		vrfConst.VRFCoordinatorV2FeeConfig,
	)
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	err = vrfv2Contracts.Coordinator.CreateSubscription()
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	err = vrfv2Contracts.Coordinator.AddConsumer(vrfConst.SubID, vrfv2Contracts.LoadTestConsumer.Address())
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	err = FundVRFCoordinatorV2Subscription(lt, vrfv2Contracts.Coordinator, env.EVMClient, vrfConst.SubID, subFundingLINK)
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	jobs, err := CreateVRFV2Jobs(env.ClCluster.NodeAPIs(), vrfv2Contracts.Coordinator, env.EVMClient, vrfConst.MinimumConfirmations)
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	addr, err := env.ClCluster.Nodes[0].API.PrimaryEthAddress()
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	nodeConfig := node.NewConfig(env.ClCluster.Nodes[0].NodeConfig,
		node.WithVRFv2EVMEstimator(addr),
	)
	err = env.ClCluster.Nodes[0].Restart(nodeConfig)
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	return env, vrfv2Contracts, jobs[0].KeyHash, nil
}
