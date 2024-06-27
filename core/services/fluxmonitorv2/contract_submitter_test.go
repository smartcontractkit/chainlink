package fluxmonitorv2_test

import (
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
	fmmocks "github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2/mocks"
)

func TestFluxAggregatorContractSubmitter_Submit(t *testing.T) {
	t.Parallel()
	var (
		fluxAggregator    = mocks.NewFluxAggregator(t)
		orm               = fmmocks.NewORM(t)
		keyStore          = fmmocks.NewKeyStoreInterface(t)
		gasLimit          = uint64(2100)
		forwardingAllowed = false
		submitter         = fluxmonitorv2.NewFluxAggregatorContractSubmitter(fluxAggregator, orm, keyStore, gasLimit, forwardingAllowed, testutils.FixtureChainID)

		toAddress   = testutils.NewAddress()
		fromAddress = testutils.NewAddress()
		roundID     = big.NewInt(1)
		submission  = big.NewInt(2)
	)

	payload, err := fluxmonitorv2.FluxAggregatorABI.Pack("submit", roundID, submission)
	assert.NoError(t, err)

	keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID).Return(fromAddress, nil)
	fluxAggregator.On("Address").Return(toAddress)

	idempotencyKey := uuid.New().String()
	orm.On("CreateEthTransaction", mock.Anything, fromAddress, toAddress, payload, gasLimit, &idempotencyKey).Return(nil)

	err = submitter.Submit(testutils.Context(t), roundID, submission, &idempotencyKey)
	assert.NoError(t, err)
}
