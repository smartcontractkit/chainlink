package vrfv1

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var (
	ErrDeployBHSV1            = "error deploying BlockHashStoreV1 contract"
	ErrDeployVRFCootrinatorV1 = "error deploying VRFv1 Coordinator contract"
	ErrDeployVRFConsumerV1    = "error deploying VRFv1 Consumer contract"
)

type Contracts struct {
	BHS         contracts.BlockHashStore
	Coordinator contracts.VRFCoordinator
	Consumer    contracts.VRFConsumer
}

func DeployVRFContracts(cd contracts.ContractDeployer, bc blockchain.EVMClient, lt contracts.LinkToken) (*Contracts, error) {
	bhs, err := cd.DeployBlockhashStore()
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployBHSV1)
	}
	coordinator, err := cd.DeployVRFCoordinator(lt.Address(), bhs.Address())
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployVRFCootrinatorV1)
	}
	consumer, err := cd.DeployVRFConsumer(lt.Address(), coordinator.Address())
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployVRFConsumerV1)
	}
	if err := bc.WaitForEvents(); err != nil {
		return nil, err
	}
	return &Contracts{bhs, coordinator, consumer}, nil
}
