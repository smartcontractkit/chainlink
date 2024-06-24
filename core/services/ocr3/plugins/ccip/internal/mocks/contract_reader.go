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
func (cr *ContractReaderMock) GetLatestValue(ctx context.Context, contractName, method string, params, returnVal any) error {
	args := cr.Called(ctx, contractName, method, params, returnVal)
	return args.Error(0)
}

func (cr *ContractReaderMock) Bind(ctx context.Context, bindings []types.BoundContract) error {
	args := cr.Called(ctx, bindings)
	return args.Error(0)
}

func (cr *ContractReaderMock) QueryKey(ctx context.Context, contractName string, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	args := cr.Called(ctx, contractName, filter, limitAndSort, sequenceDataType)
	return args.Get(0).([]types.Sequence), args.Error(1)
}

func (cr *ContractReaderMock) Start(ctx context.Context) error {
	args := cr.Called(ctx)
	return args.Error(0)
}

func (cr *ContractReaderMock) Close() error {
	args := cr.Called()
	return args.Error(0)
}

func (cr *ContractReaderMock) Ready() error {
	args := cr.Called()
	return args.Error(0)
}

func (cr *ContractReaderMock) HealthReport() map[string]error {
	args := cr.Called()
	return args.Get(0).(map[string]error)
}

func (cr *ContractReaderMock) Name() string {
	return "ContractReaderMock"
}
