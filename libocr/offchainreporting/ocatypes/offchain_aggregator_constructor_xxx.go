package ocatypes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type BillingParameters struct {
	MaximumGasPrice uint32

	ReasonableGasPrice uint32

	MicroLinkPerEth uint32

	LinkGweiPerObservation uint32

	LinkGweiPerTransmission uint32
}

type OffchainAggregatorDependencyContracts struct {
	LinkToken, Validator, WriteAccessController common.Address
	ValidatorMinAnswer, ValidatorMaxAnswer      *big.Int
}
