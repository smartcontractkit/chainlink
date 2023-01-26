package coordinator

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
)

//go:generate mockery --quiet --name vrfBeaconCoordinator --output ./mocks/ --case=underscore

// vrfBeaconCoordinator is an interface that defines methods needed by the off-chain coordinator
type vrfBeaconCoordinator interface {
	// sProvingKeyHash retrieves the proving key hash from the on-chain contract.
	sProvingKeyHash(opts *bind.CallOpts) ([32]byte, error)

	// sKeyID retrieves the keyID from the on-chain contract.
	sKeyID(opts *bind.CallOpts) ([32]byte, error)

	// iBeaconPeriodBlocks retrieves the beacon period in blocks from the on-chain contract.
	iBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	// parseLog parses the raw log data and topics into a go object.
	// The returned object must be casted to the expected type.
	parseLog(log types.Log) (generated.AbigenLog, error)

	// getConfirmationDelays retrieves confirmation delays from the on-chain contract.
	getConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error)
}
