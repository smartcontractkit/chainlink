package fakes

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type TelemetryServiceMock struct{}

func (t *TelemetryServiceMock) Send(ctx context.Context, network string, chainID string, contractID string, telemetryType string, payload []byte) error {
	return nil
}

type KVStoreMock struct{}

func (k *KVStoreMock) Store(ctx context.Context, key string, val []byte) error {
	return nil
}
func (k *KVStoreMock) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

type ErrorLogMock struct{}

func (e *ErrorLogMock) SaveError(ctx context.Context, msg string) error {
	return nil
}

type RelayerSetMock struct{}

func (r *RelayerSetMock) Get(ctx context.Context, relayID types.RelayID) (core.Relayer, error) {
	return nil, nil
}
func (r *RelayerSetMock) List(ctx context.Context, relayIDs ...types.RelayID) (map[types.RelayID]core.Relayer, error) {
	return nil, nil
}

type PipelineRunnerServiceMock struct{}

func (p *PipelineRunnerServiceMock) ExecuteRun(ctx context.Context, spec string, vars core.Vars, options core.Options) (core.TaskResults, error) {
	return nil, nil
}

type OracleFactoryMock struct{}

func (o *OracleFactoryMock) NewOracle(ctx context.Context, args core.OracleArgs) (core.Oracle, error) {
	return &OracleMock{}, nil
}

type OracleMock struct{}

func (o *OracleMock) Start(ctx context.Context) error { return nil }
func (o *OracleMock) Close(ctx context.Context) error { return nil }
