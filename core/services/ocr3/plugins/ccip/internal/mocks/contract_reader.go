package mocks

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/stretchr/testify/mock"
)

type ContractReaderMock struct {
	*mock.Mock
}

func NewContractReaderMock() *ContractReaderMock {
	return &ContractReaderMock{
		Mock: &mock.Mock{},
	}
}

// GetLatestValue Returns given configs at initialization
func (hr *ContractReaderMock) GetLatestValue(ctx context.Context, contractName, method string, params, returnVal any) error {
	args := hr.Called(ctx, contractName, method, params, returnVal)
	return args.Error(0)
}

func (hr *ContractReaderMock) Bind(ctx context.Context, bindings []types.BoundContract) error {
	args := hr.Called(ctx, bindings)
	return args.Error(0)
}

func (hr *ContractReaderMock) QueryKey(ctx context.Context, contractName string, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	args := hr.Called(ctx, contractName, filter, limitAndSort, sequenceDataType)
	return args.Get(0).([]types.Sequence), args.Error(1)
}

func (hr *ContractReaderMock) Start(ctx context.Context) error {
	args := hr.Called(ctx)
	return args.Error(0)
}

func (hr *ContractReaderMock) Close() error {
	args := hr.Called()
	return args.Error(0)
}

func (hr *ContractReaderMock) Ready() error {
	args := hr.Called()
	return args.Error(0)
}

func (hr *ContractReaderMock) HealthReport() map[string]error {
	args := hr.Called()
	return args.Get(0).(map[string]error)
}

func (hr *ContractReaderMock) Name() string {
	return "ContractReaderMock"
}
