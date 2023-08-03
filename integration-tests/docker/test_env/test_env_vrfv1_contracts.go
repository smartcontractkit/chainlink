package test_env

import (
	"github.com/pkg/errors"
)

var (
	ErrDeployBHSV1            = "failed to deploy BlockHashStoreV1 contract"
	ErrDeployVRFCootrinatorV1 = "failed to deploy VRFv1 Coordinator contract"
	ErrDeployVRFConsumerV1    = "failed to deploy VRFv1 Consumer contract"
)

func (te *CLClusterTestEnv) DeployVRFContracts() error {
	bhs, err := te.Geth.ContractDeployer.DeployBlockhashStore()
	if err != nil {
		return errors.Wrap(err, ErrDeployBHSV1)
	}
	te.BHSV1 = bhs
	coordinator, err := te.Geth.ContractDeployer.DeployVRFCoordinator(te.LinkToken.Address(), bhs.Address())
	if err != nil {
		return errors.Wrap(err, ErrDeployVRFCootrinatorV1)
	}
	te.CoordinatorV1 = coordinator
	consumer, err := te.Geth.ContractDeployer.DeployVRFConsumer(te.LinkToken.Address(), coordinator.Address())
	if err != nil {
		return errors.Wrap(err, ErrDeployVRFConsumerV1)
	}
	te.ConsumerV1 = consumer
	return te.Geth.EthClient.WaitForEvents()
}
