package test_env

import (
	"math/big"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	ErrCreatingProvingKey   = "error creating a keyHash from the proving key"
	ErrDeployBlockHashStore = "error deploying blockhash store"
	ErrDeployCoordinator    = "error deploying VRFv2 CoordinatorV2"
	ErrAdvancedConsumer     = "error deploying VRFv2 Advanced Consumer"
	ErrABIEncodingFunding   = "error Abi encoding subscriptionID"
	ErrSendingLinkToken     = "error sending Link token"
)

type VRFV2EncodedProvingKey [2]*big.Int

type VRFV2Contracts struct {
	Coordinator      contracts.VRFCoordinatorV2
	BHS              contracts.BlockHashStore
	LoadTestConsumer contracts.VRFv2LoadTestConsumer
}

func VRFV2RegisterProvingKey(
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2,
) (VRFV2EncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return VRFV2EncodedProvingKey{}, err
	}
	if err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	); err != nil {
		return VRFV2EncodedProvingKey{}, err
	}
	return provingKey, nil
}

func (m *CLClusterTestEnv) DeployVRFV2Contracts() error {
	bhs, err := m.Geth.ContractDeployer.DeployBlockhashStore()
	if err != nil {
		return errors.Wrap(err, ErrDeployBlockHashStore)
	}
	m.BHSV2 = bhs
	coordinator, err := m.Geth.ContractDeployer.DeployVRFCoordinatorV2(m.LinkToken.Address(), bhs.Address(), m.MockETHLinkFeed.Address())
	if err != nil {
		return errors.Wrap(err, ErrDeployCoordinator)
	}
	m.CoordinatorV2 = coordinator
	loadTestConsumer, err := m.Geth.ContractDeployer.DeployVRFv2LoadTestConsumer(coordinator.Address())
	if err != nil {
		return errors.Wrap(err, ErrAdvancedConsumer)
	}
	m.LoadTestConsumer = loadTestConsumer
	return m.WaitForEvents()
}

func (m *CLClusterTestEnv) FundVRFCoordinatorV2Subscription(subscriptionID uint64, linkFundingAmount *big.Int) error {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subscriptionID)
	if err != nil {
		return errors.Wrap(err, ErrABIEncodingFunding)
	}
	_, err = m.LinkToken.TransferAndCall(m.CoordinatorV2.Address(), big.NewInt(0).Mul(linkFundingAmount, big.NewInt(1e18)), encodedSubId)
	if err != nil {
		return errors.Wrap(err, ErrSendingLinkToken)
	}
	return m.WaitForEvents()
}
