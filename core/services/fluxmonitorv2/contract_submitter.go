package fluxmonitorv2

import (
	"context"
	"math/big"

	"github.com/pkg/errors"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flux_aggregator_wrapper"
)

//go:generate mockery --quiet --name ContractSubmitter --output ./mocks/ --case=underscore

// FluxAggregatorABI initializes the Flux Aggregator ABI
var FluxAggregatorABI = evmtypes.MustGetABI(flux_aggregator_wrapper.FluxAggregatorABI)

// ContractSubmitter defines an interface to submit an eth tx.
type ContractSubmitter interface {
	Submit(ctx context.Context, roundID *big.Int, submission *big.Int, idempotencyKey *string) error
}

// FluxAggregatorContractSubmitter submits the polled answer in an eth tx.
type FluxAggregatorContractSubmitter struct {
	flux_aggregator_wrapper.FluxAggregatorInterface
	orm               ORM
	keyStore          KeyStoreInterface
	gasLimit          uint64
	forwardingAllowed bool
	chainID           *big.Int
}

// NewFluxAggregatorContractSubmitter constructs a new NewFluxAggregatorContractSubmitter
func NewFluxAggregatorContractSubmitter(
	contract flux_aggregator_wrapper.FluxAggregatorInterface,
	orm ORM,
	keyStore KeyStoreInterface,
	gasLimit uint64,
	forwardingAllowed bool,
	chainID *big.Int,
) *FluxAggregatorContractSubmitter {
	return &FluxAggregatorContractSubmitter{
		FluxAggregatorInterface: contract,
		orm:                     orm,
		keyStore:                keyStore,
		gasLimit:                gasLimit,
		forwardingAllowed:       forwardingAllowed,
		chainID:                 chainID,
	}
}

// Submit submits the answer by writing a EthTx for the txmgr to
// pick up
func (c *FluxAggregatorContractSubmitter) Submit(ctx context.Context, roundID *big.Int, submission *big.Int, idempotencyKey *string) error {
	fromAddress, err := c.keyStore.GetRoundRobinAddress(ctx, c.chainID)
	if err != nil {
		return err
	}

	payload, err := FluxAggregatorABI.Pack("submit", roundID, submission)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	return errors.Wrap(
		c.orm.CreateEthTransaction(ctx, fromAddress, c.Address(), payload, c.gasLimit, idempotencyKey),
		"failed to send Eth transaction",
	)
}
