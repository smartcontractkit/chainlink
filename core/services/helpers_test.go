package services_test

import (
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/stretchr/testify/mock"
)

func NewEthClientMock(t mock.TestingT) *mocks.Client {
	mockEth := new(mocks.Client)
	mockEth.Test(t)
	return mockEth
}
