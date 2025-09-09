package standardcapabilities

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestStandardCapabilityStart(t *testing.T) {
	t.Run("NOK-not_found_binary_does_not_block", func(t *testing.T) {
		ctx := t.Context()
		lggr := logger.TestLogger(t)

		pluginRegistrar := plugins.NewRegistrarConfig(loop.GRPCOpts{}, func(name string) (*plugins.RegisteredLoop, error) { return &plugins.RegisteredLoop{}, nil }, func(loopId string) {})
		registry := mocks.NewCapabilitiesRegistry(t)

		spec := &job.StandardCapabilitiesSpec{
			Command: "not/found/path/to/binary",
			OracleFactory: job.OracleFactoryConfig{
				Enabled: true,
				BootstrapPeers: []string{
					"12D3KooWEBVwbfdhKnicois7FTYVsBFGFcoMhMCKXQC57BQyZMhz@localhost:6690",
				},
				OCRContractAddress: "0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6",
				ChainID:            "31337",
			}}

		standardCapability := NewStandardCapabilities(lggr, spec, pluginRegistrar, &telemetryServiceMock{}, &kvstoreMock{}, registry, &errorLogMock{}, &pipelineRunnerServiceMock{}, &relayerSetMock{}, &oracleFactoryMock{}, &gatewayConnectorMock{}, &keystoreMock{})
		standardCapability.startTimeout = 1 * time.Second
		err := standardCapability.Start(ctx)
		require.NoError(t, err)

		standardCapability.wg.Wait()
	})
}

type telemetryServiceMock struct{}

func (t *telemetryServiceMock) Send(ctx context.Context, network string, chainID string, contractID string, telemetryType string, payload []byte) error {
	return nil
}

type kvstoreMock struct{}

func (k *kvstoreMock) Store(ctx context.Context, key string, val []byte) error {
	return nil
}
func (k *kvstoreMock) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}
func (k *kvstoreMock) PruneExpiredEntries(ctx context.Context, maxAge time.Duration) (int64, error) {
	return 0, nil
}

type keystoreMock struct{ core.UnimplementedKeystore }

func (k *keystoreMock) Accounts(ctx context.Context) (accounts []string, err error) {
	return nil, nil
}
func (k *keystoreMock) Sign(ctx context.Context, account string, data []byte) (signed []byte, err error) {
	return nil, nil
}

type errorLogMock struct{}

func (e *errorLogMock) SaveError(ctx context.Context, msg string) error {
	return nil
}

type relayerSetMock struct{}

func (r *relayerSetMock) Get(ctx context.Context, relayID types.RelayID) (core.Relayer, error) {
	return nil, nil
}
func (r *relayerSetMock) List(ctx context.Context, relayIDs ...types.RelayID) (map[types.RelayID]core.Relayer, error) {
	return nil, nil
}

type pipelineRunnerServiceMock struct{}

func (p *pipelineRunnerServiceMock) ExecuteRun(ctx context.Context, spec string, vars core.Vars, options core.Options) (core.TaskResults, error) {
	return nil, nil
}

type oracleFactoryMock struct{}

func (o *oracleFactoryMock) NewOracle(ctx context.Context, args core.OracleArgs) (core.Oracle, error) {
	return &oracleMock{}, nil
}

type oracleMock struct{}

func (o *oracleMock) Start(ctx context.Context) error {
	return nil
}
func (o *oracleMock) Close(ctx context.Context) error {
	return nil
}

type gatewayConnectorMock struct{}

func (g *gatewayConnectorMock) Start(context.Context) error {
	return nil
}

func (g *gatewayConnectorMock) Close() error {
	return nil
}

func (g *gatewayConnectorMock) AddHandler(ctx context.Context, methods []string, handler core.GatewayConnectorHandler) error {
	return nil
}

func (g *gatewayConnectorMock) SendToGateway(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) error {
	return nil
}

func (g *gatewayConnectorMock) SignMessage(ctx context.Context, msg []byte) ([]byte, error) {
	return nil, nil
}

func (g *gatewayConnectorMock) GatewayIDs(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (g *gatewayConnectorMock) DonID(ctx context.Context) (string, error) {
	return "", nil
}

func (g *gatewayConnectorMock) AwaitConnection(ctx context.Context, gatewayID string) error {
	return nil
}
