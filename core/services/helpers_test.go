package services_test

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/stretchr/testify/mock"
)

func NewEthClientMock(t mock.TestingT) *mocks.Client {
	mockEth := new(mocks.Client)
	mockEth.Test(t)
	mockEth.On("ChainID").Maybe().Return(big.NewInt(0))
	return mockEth
}
