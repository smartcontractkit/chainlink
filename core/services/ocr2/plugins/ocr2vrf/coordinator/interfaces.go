package coordinator

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
)

//go:generate mockery --quiet --name VRFBeaconCoordinator --output ./mocks/ --case=underscore

// VRFBeaconCoordinator is an interface that defines methods needed by the off-chain coordinator
type VRFBeaconCoordinator interface {
	// SProvingKeyHash retrieves the proving key hash from the on-chain contract.
	SProvingKeyHash(opts *bind.CallOpts) ([32]byte, error)

	// SKeyID retrieves the keyID from the on-chain contract.
	SKeyID(opts *bind.CallOpts) ([32]byte, error)

	// IBeaconPeriodBlocks retrieves the beacon period in blocks from the on-chain contract.
	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	// ParseLog parses the raw log data and topics into a go object.
	// The returned object must be casted to the expected type.
	ParseLog(log types.Log) (generated.AbigenLog, error)

	// GetConfirmationDelays retrieves confirmation delays from the on-chain contract.
	GetConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error)
}
