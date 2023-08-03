package test_env

import (
	"github.com/pkg/errors"
)

var (
	ErrDeployBHSV1            = "failed to deploy BlockHashStoreV1 contract"
	ErrDeployVRFCootrinatorV1 = "failed to deploy VRFv1 Coordinator contract"
	ErrDeployVRFConsumerV1    = "failed to deploy VRFv1 Consumer contract"
)

func (m *CLClusterTestEnv) DeployVRFContracts() error {
	bhs, err := m.Geth.ContractDeployer.DeployBlockhashStore()
	if err != nil {
		return errors.Wrap(err, ErrDeployBHSV1)
	}
	m.BHSV1 = bhs
	coordinator, err := m.Geth.ContractDeployer.DeployVRFCoordinator(m.LinkToken.Address(), bhs.Address())
	if err != nil {
		return errors.Wrap(err, ErrDeployVRFCootrinatorV1)
	}
	m.CoordinatorV1 = coordinator
	consumer, err := m.Geth.ContractDeployer.DeployVRFConsumer(m.LinkToken.Address(), coordinator.Address())
	if err != nil {
		return errors.Wrap(err, ErrDeployVRFConsumerV1)
	}
	m.ConsumerV1 = consumer
	return m.Geth.EthClient.WaitForEvents()
}
