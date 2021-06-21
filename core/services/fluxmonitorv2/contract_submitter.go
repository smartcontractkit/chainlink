package fluxmonitorv2

import (
	"math/big"

	"gorm.io/gorm"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

//go:generate mockery --name ContractSubmitter --output ./mocks/ --case=underscore

// FluxAggregatorABI initializes the Flux Aggregator ABI
var FluxAggregatorABI = eth.MustGetABI(flux_aggregator_wrapper.FluxAggregatorABI)

// ContractSubmitter defines an interface to submit an eth tx.
type ContractSubmitter interface {
	Submit(db *gorm.DB, roundID *big.Int, submission *big.Int) error
}

// FluxAggregatorContractSubmitter submits the polled answer in an eth tx.
type FluxAggregatorContractSubmitter struct {
	flux_aggregator_wrapper.FluxAggregatorInterface
	orm      ORM
	keyStore KeyStoreInterface
	gasLimit uint64
}

// NewFluxAggregatorContractSubmitter constructs a new NewFluxAggregatorContractSubmitter
func NewFluxAggregatorContractSubmitter(
	contract flux_aggregator_wrapper.FluxAggregatorInterface,
	orm ORM,
	keyStore KeyStoreInterface,
	gasLimit uint64,
) *FluxAggregatorContractSubmitter {
	return &FluxAggregatorContractSubmitter{
		FluxAggregatorInterface: contract,
		orm:                     orm,
		keyStore:                keyStore,
		gasLimit:                gasLimit,
	}
}

// Submit submits the answer by writing a EthTx for the bulletprooftxmanager to
// pick up
func (c *FluxAggregatorContractSubmitter) Submit(db *gorm.DB, roundID *big.Int, submission *big.Int) error {
	fromAddress, err := c.keyStore.GetRoundRobinAddress()
	if err != nil {
		return err
	}

	payload, err := FluxAggregatorABI.Pack("submit", roundID, submission)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	return errors.Wrap(
		c.orm.CreateEthTransaction(db, fromAddress, c.Address(), payload, c.gasLimit),
		"failed to send Eth transaction",
	)
}
