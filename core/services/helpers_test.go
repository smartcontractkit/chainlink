package services_test

import (
	"math/big"

	ethmocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/stretchr/testify/mock"
)

func NewEthClientMock(t mock.TestingT) *ethmocks.Client {
	mockEth := new(ethmocks.Client)
	mockEth.Test(t)
	mockEth.On("ChainID").Maybe().Return(big.NewInt(0))
	return mockEth
}
