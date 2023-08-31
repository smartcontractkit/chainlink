package vrfv2plus

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
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
		return nil, err
	}
	return &VRFV2PlusContracts{coordinator, bhs, loadTestConsumer}, nil
}

func CreateVRFV2PlusJobs(
	chainlinkNodes []*client.ChainlinkClient,
	coordinator contracts.VRFCoordinatorV2Plus,
	c blockchain.EVMClient,
	minIncomingConfirmations uint16,
) ([]VRFV2PlusJobInfo, error) {
	jobInfo := make([]VRFV2PlusJobInfo, 0)
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
		ji := VRFV2PlusJobInfo{
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
