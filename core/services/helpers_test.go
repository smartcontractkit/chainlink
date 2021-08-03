package services_test

import (
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/stretchr/testify/mock"
)

func NewEthClientAndSubMock(t mock.TestingT) (*mocks.Client, *mocks.Subscription) {
	mockSub := new(mocks.Subscription)
	mockSub.Test(t)
	mockEth := new(mocks.Client)
	mockEth.Test(t)
	return mockEth, mockSub
}

func NewEthClientMock(t mock.TestingT) *mocks.Client {
	mockEth := new(mocks.Client)
	mockEth.Test(t)
	return mockEth
}
