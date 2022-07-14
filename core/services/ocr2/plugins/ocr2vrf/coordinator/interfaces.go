package coordinator

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

//go:generate mockery --name VRFBeaconCoordinator --output ./mocks/ --case=underscore

// VRFBeaconCoordinator is a narrow interface implemented by the contract go wrappers.
type VRFBeaconCoordinator interface {
	// SProvingKeyHash retrieves the proving key hash from the on-chain contract.
	SProvingKeyHash(opts *bind.CallOpts) ([32]byte, error)

	// IBeaconPeriodBlocks retrieves the beacon period in blocks from the on-chain contract.
	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)
}

type DKG interface {
}
