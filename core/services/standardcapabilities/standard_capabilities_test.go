package standardcapabilities

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

// TestStandardCapabilities_ForwardsPluginEnvFile tests that the standard capabilities
// LOOPP launcher reads the env file referenced by CL_CAPABILITIES_ENV and forwards
// its contents to the plugin. This is how operator-supplied env vars (e.g. chain-selector's
// EXTRA_SELECTORS_FILE) make it into the LOOPP subprocess despite plugins.NewCmdFactory
// not inheriting os.Environ().
func TestStandardCapabilities_ForwardsPluginEnvFile(t *testing.T) {
	writeEnvFile := func(t *testing.T, contents string) string {
		t.Helper()
		path := filepath.Join(t.TempDir(), "capabilities.env")
		require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))
		return path
	}

	startAndCapture := func(t *testing.T) (plugins.CmdConfig, error) {
		t.Helper()
		var captured plugins.CmdConfig
		registrar := &capturingRegistrar{
			onRegister: func(cfg plugins.CmdConfig) { captured = cfg },
		}
		std := NewStandardCapabilities(
			logger.TestLogger(t),
			"not/found/path/to/binary",
			"{}",
			registrar,
			core.StandardCapabilitiesDependencies{},
		)
		// Start runs RegisterLOOP synchronously; the capturingRegistrar
		// returns an error so Start unwinds before any LOOPP lifecycle code
		// runs. The returned error from Start is propagated to the caller so
		// each sub-test can assert on it as appropriate.
		return captured, std.Start(t.Context())
	}

	t.Run("env file is parsed and forwarded via CmdConfig.Env", func(t *testing.T) {
		envFile := writeEnvFile(t, "EXTRA_SELECTORS_FILE=/extraConfig/extra-selectors.yaml\nFOO=bar\n")
		t.Setenv(string(env.CapabilitiesPlugin.Env), envFile)

		cfg, err := startAndCapture(t)
		require.Error(t, err, "expected synthetic Start failure from capturingRegistrar")
		require.Contains(t, err.Error(), capturingRegistrarErr)

		require.Contains(t, cfg.Env, "EXTRA_SELECTORS_FILE=/extraConfig/extra-selectors.yaml",
			"standard-capability LOOPP launcher should forward EXTRA_SELECTORS_FILE from the operator-supplied env file so chain-selectors' init() in the LOOPP can read operator-provided selectors")
		require.Contains(t, cfg.Env, "FOO=bar",
			"all entries in the operator-supplied env file should be forwarded, not just well-known ones")
	})

	t.Run("env file unset results in empty CmdConfig.Env", func(t *testing.T) {
		t.Setenv(string(env.CapabilitiesPlugin.Env), "")

		cfg, err := startAndCapture(t)
		require.Error(t, err, "expected synthetic Start failure from capturingRegistrar")
		require.Contains(t, err.Error(), capturingRegistrarErr)

		require.Empty(t, cfg.Env,
			"no operator-provided env vars should be forwarded when CL_CAPABILITIES_ENV is unset")
	})

	t.Run("missing env file fails Start before RegisterLOOP", func(t *testing.T) {
		missingPath := filepath.Join(t.TempDir(), "does-not-exist.env")
		t.Setenv(string(env.CapabilitiesPlugin.Env), missingPath)

		var registerCalled bool
		registrar := &capturingRegistrar{
			onRegister: func(plugins.CmdConfig) { registerCalled = true },
		}
		std := NewStandardCapabilities(
			logger.TestLogger(t),
			"not/found/path/to/binary",
			"{}",
			registrar,
			core.StandardCapabilitiesDependencies{},
		)

		err := std.Start(t.Context())
		require.Error(t, err, "Start should fail when the env file is missing")
		require.Contains(t, err.Error(), "failed to parse capabilities env file")
		require.False(t, registerCalled, "RegisterLOOP must not be called when env-file parsing fails")
	})
}

const capturingRegistrarErr = "capturingRegistrar: stop after capture"

// capturingRegistrar records the CmdConfig handed to RegisterLOOP and then
// returns a sentinel error so the caller's Start flow exits before any
// LOOPP lifecycle code runs.
type capturingRegistrar struct {
	onRegister func(plugins.CmdConfig)
}

func (c *capturingRegistrar) RegisterLOOP(cfg plugins.CmdConfig) (func() *exec.Cmd, loop.GRPCOpts, error) {
	if c.onRegister != nil {
		c.onRegister(cfg)
	}
	return nil, loop.GRPCOpts{}, errors.New(capturingRegistrarErr)
}

func (c *capturingRegistrar) UnregisterLOOP(string) {}

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

		dependencies := core.StandardCapabilitiesDependencies{
			Config:             spec.Config,
			TelemetryService:   &telemetryServiceMock{},
			Store:              &kvstoreMock{},
			CapabilityRegistry: registry,
			ErrorLog:           &errorLogMock{},
			PipelineRunner:     &pipelineRunnerServiceMock{},
			RelayerSet:         &relayerSetMock{},
			OracleFactory:      &oracleFactoryMock{},
			GatewayConnector:   &gatewayConnectorMock{},
			P2PKeystore:        &keystoreMock{},
		}
		standardCapability := NewStandardCapabilities(lggr, spec.Command, spec.Config, pluginRegistrar, dependencies)
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

type gatewayConnectorMock struct {
	core.UnimplementedGatewayConnector
}

func (g *gatewayConnectorMock) Start(context.Context) error {
	return nil
}

func (g *gatewayConnectorMock) Close() error {
	return nil
}

func (g *gatewayConnectorMock) AddHandler(ctx context.Context, methods []string, handler core.GatewayConnectorHandler) error {
	return nil
}

func (g *gatewayConnectorMock) RemoveHandler(ctx context.Context, methods []string) error {
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
