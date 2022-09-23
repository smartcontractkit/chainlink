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

var _ VRFBeaconCoordinator = &VRFProxy{}

//go:generate mockery --name VRFBeaconCoordinator --output ./mocks/ --case=underscore

// VRFProxy routes requests to VRFBeacon and VRFCoordinator go wrappers and implements VRFBeaconCoordinator interface
type VRFProxy struct {
	lggr        logger.Logger
	beacon      vrf_beacon.VRFBeacon
	coordinator vrf_coordinator.VRFCoordinator
	evmClient   evmclient.Client
}

func NewProxy(
	lggr logger.Logger,
	beaconAddress common.Address,
	coordinatorAddress common.Address,
	client evmclient.Client,
) (VRFBeaconCoordinator, error) {
	beacon, err := vrf_beacon.NewVRFBeacon(beaconAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "beacon wrapper creation")
	}
	coordinator, err := vrf_coordinator.NewVRFCoordinator(coordinatorAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "coordinator wrapper creation")
	}
	return &VRFProxy{
		lggr:        lggr,
		beacon:      *beacon,
		coordinator: *coordinator,
		evmClient:   client,
	}, nil
}

// SProvingKeyHash retrieves the proving key hash from the on-chain contract.
func (v *VRFProxy) SProvingKeyHash(opts *bind.CallOpts) ([32]byte, error) {
	return v.beacon.SProvingKeyHash(opts)
}

// SKeyID retrieves the keyID from the on-chain contract.
func (v *VRFProxy) SKeyID(opts *bind.CallOpts) ([32]byte, error) {
	return v.beacon.SKeyID(opts)
}

// IBeaconPeriodBlocks retrieves the beacon period in blocks from the on-chain contract.
func (v *VRFProxy) IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	return v.coordinator.IBeaconPeriodBlocks(opts)
}

// ParseLog parses the raw log data and topics into a go object.
// The returned object must be casted to the expected type.
func (v *VRFProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	if log.Address == v.beacon.Address() {
		return v.beacon.ParseLog(log)
	} else if log.Address == v.coordinator.Address() {
		return v.coordinator.ParseLog(log)
	} else {
		return nil, errors.Errorf("failed to parse log. contractAddress: %x logs: %x", log.Address, log.Topics)
	}
}

// GetConfirmationDelays retrieves confirmation delays from the on-chain contract.
func (v *VRFProxy) GetConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error) {
	return v.coordinator.GetConfirmationDelays(opts)
}
