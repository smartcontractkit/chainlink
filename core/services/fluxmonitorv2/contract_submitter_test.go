package fluxmonitorv2_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	fmmocks "github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFluxAggregatorContractSubmitter_Submit(t *testing.T) {
	var (
		fluxAggregator    = mocks.NewFluxAggregator(t)
		orm               = fmmocks.NewORM(t)
		keyStore          = fmmocks.NewKeyStoreInterface(t)
		gasLimit          = uint32(2100)
		forwardingAllowed = false
		submitter         = fluxmonitorv2.NewFluxAggregatorContractSubmitter(fluxAggregator, orm, keyStore, gasLimit, forwardingAllowed, testutils.FixtureChainID)

		toAddress   = testutils.NewAddress()
		fromAddress = testutils.NewAddress()
		roundID     = big.NewInt(1)
		submission  = big.NewInt(2)
	)

	payload, err := fluxmonitorv2.FluxAggregatorABI.Pack("submit", roundID, submission)
	assert.NoError(t, err)

	keyStore.On("GetRoundRobinAddress", testutils.FixtureChainID).Return(fromAddress, nil)
	fluxAggregator.On("Address").Return(toAddress)
	orm.On("CreateEthTransaction", fromAddress, toAddress, payload, gasLimit).Return(nil)

	err = submitter.Submit(roundID, submission)
	assert.NoError(t, err)
}
