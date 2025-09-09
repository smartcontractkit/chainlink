package fakes

import (
	"context"
	"encoding/json"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type TelemetryServiceMock struct{}

func (t *TelemetryServiceMock) Send(ctx context.Context, network string, chainID string, contractID string, telemetryType string, payload []byte) error {
	return nil
}

type KVStoreMock struct {
	core.UnimplementedKeystore
}

func (k *KVStoreMock) Store(ctx context.Context, key string, val []byte) error {
	return nil
}
func (k *KVStoreMock) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}
func (k *KVStoreMock) PruneExpiredEntries(ctx context.Context, maxAge time.Duration) (int64, error) {
	return 0, nil
}

type KeystoreMock struct {
	core.UnimplementedKeystore
}

func (k *KeystoreMock) Accounts(ctx context.Context) (accounts []string, err error) {
	return nil, nil
}
func (k *KeystoreMock) Sign(ctx context.Context, account string, data []byte) (signed []byte, err error) {
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

type GatewayConnectorMock struct {
	core.UnimplementedGatewayConnector
}

func (g *GatewayConnectorMock) Start(context.Context) error {
	return nil
}

func (g *GatewayConnectorMock) Close() error {
	return nil
}

func (g *GatewayConnectorMock) AddHandler(context.Context, []string, core.GatewayConnectorHandler) error {
	return nil
}

func (g *GatewayConnectorMock) SendToGateway(context.Context, string, *jsonrpc.Response[json.RawMessage]) error {
	return nil
}

func (g *GatewayConnectorMock) SignMessage(context.Context, []byte) ([]byte, error) {
	return nil, nil
}

func (g *GatewayConnectorMock) GatewayIDs(context.Context) ([]string, error) {
	return nil, nil
}

func (g *GatewayConnectorMock) DonID(context.Context) (string, error) {
	return "", nil
}

func (g *GatewayConnectorMock) AwaitConnection(context.Context, string) error {
	return nil
}
