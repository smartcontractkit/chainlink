package coordinator

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	"github.com/smartcontractkit/chainlink/core/logger"
)

//go:generate mockery --quiet --name VRFBeaconCoordinator --output ./mocks/ --case=underscore

// vrfRouter routes requests to VRFBeacon and VRFCoordinator go wrappers
type vrfRouter struct {
	lggr        logger.Logger
	beacon      vrf_beacon.VRFBeaconInterface
	coordinator vrf_coordinator.VRFCoordinatorInterface
}

// vrfRouter implements VRFBeaconCoordinator interface
var _ vrfBeaconCoordinator = &vrfRouter{}

func newRouter(
	lggr logger.Logger,
	beaconAddress common.Address,
	coordinatorAddress common.Address,
	client evmclient.Client,
) (vrfBeaconCoordinator, error) {
	beacon, err := vrf_beacon.NewVRFBeacon(beaconAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "beacon wrapper creation")
	}
	coordinator, err := vrf_coordinator.NewVRFCoordinator(coordinatorAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "coordinator wrapper creation")
	}
	return &vrfRouter{
		lggr:        lggr,
		beacon:      beacon,
		coordinator: coordinator,
	}, nil
}

// sProvingKeyHash retrieves the proving key hash from the on-chain contract.
// Calls VRF beacon wrapper to retrieve proving key hash
func (v *vrfRouter) sProvingKeyHash(opts *bind.CallOpts) ([32]byte, error) {
	return v.beacon.SProvingKeyHash(opts)
}

// sKeyID retrieves the keyID from the on-chain contract.
// Calls VRF beacon wrapper to retrieve key ID
func (v *vrfRouter) sKeyID(opts *bind.CallOpts) ([32]byte, error) {
	return v.beacon.SKeyID(opts)
}

// iBeaconPeriodBlocks retrieves the beacon period in blocks from the on-chain contract.
// Calls VRF coordinator wrapper to beacon period blocks
func (v *vrfRouter) iBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	return v.coordinator.IBeaconPeriodBlocks(opts)
}

// parseLog parses the raw log data and topics into a go object.
// The returned object must be casted to the expected type.
// Calls either VRF beacon wrapper or VRF coordinator wrapper depending on the addresses of the log
func (v *vrfRouter) parseLog(log types.Log) (generated.AbigenLog, error) {
	if log.Address == v.beacon.Address() {
		return v.beacon.ParseLog(log)
	} else if log.Address == v.coordinator.Address() {
		return v.coordinator.ParseLog(log)
	} else {
		return nil, errors.Errorf("failed to parse log. contractAddress: %x logs: %x", log.Address, log.Topics)
	}
}

// getConfirmationDelays retrieves confirmation delays from the on-chain contract.
// Calls VRF coordinator to retrieve confirmation delays
func (v *vrfRouter) getConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error) {
	return v.coordinator.GetConfirmationDelays(opts)
}
