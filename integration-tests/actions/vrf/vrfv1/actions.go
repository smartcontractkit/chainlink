package vrfv1

import (
	"fmt"

	"github.com/smartcontractkit/seth"

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

func DeployVRFContracts(client *seth.Client, linkTokenAddress string) (*Contracts, error) {
	bhs, err := contracts.DeployBlockhashStore(client)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployBHSV1, err)
	}
	coordinator, err := contracts.DeployVRFCoordinator(client, linkTokenAddress, bhs.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployVRFCootrinatorV1, err)
	}
	consumer, err := contracts.DeployVRFConsumer(client, linkTokenAddress, coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployVRFConsumerV1, err)
	}
	return &Contracts{bhs, coordinator, consumer}, nil
}
