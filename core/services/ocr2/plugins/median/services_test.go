package median

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	libocr_median "github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	clnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

// mockRelayer captures the PluginArgs passed to NewPluginProvider
type mockRelayer struct {
	loop.Relayer
	capturedPluginArgs commontypes.PluginArgs
}

func (m *mockRelayer) NewPluginProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.PluginProvider, error) {
	m.capturedPluginArgs = pargs
	return &mockPluginProvider{}, nil
}

func TestNewMedianServices_GasLimitOverride(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	lggr := logger.TestLogger(t)

	testCases := []struct {
		name                 string
		jobGasLimit          clnull.Uint32
		pluginConfigGasLimit *uint32
		expectedGasLimit     *uint32
		description          string
	}{
		{
			name:                 "gas limit from job spec overrides nil plugin config",
			jobGasLimit:          clnull.Uint32From(50000),
			pluginConfigGasLimit: nil,
			expectedGasLimit:     uint32Ptr(50000),
			description:          "When job spec has gas limit and plugin config doesn't, job spec gas limit should be used",
		},
		{
			name:                 "plugin config gas limit takes precedence when both are set",
			jobGasLimit:          clnull.Uint32From(50000),
			pluginConfigGasLimit: uint32Ptr(60000),
			expectedGasLimit:     uint32Ptr(60000),
			description:          "When both job spec and plugin config have gas limit, plugin config should take precedence",
		},
		{
			name:                 "no gas limit when job spec gas limit is invalid",
			jobGasLimit:          clnull.Uint32{},
			pluginConfigGasLimit: nil,
			expectedGasLimit:     nil,
			description:          "When job spec gas limit is invalid and plugin config doesn't have one, no gas limit should be set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pluginCfg := config.PluginConfig{
				JuelsPerFeeCoinPipeline: `ds1 [type=bridge name=voter_turnout];`,
				GasLimit:                tc.pluginConfigGasLimit,
				JuelsPerFeeCoinCache: &config.JuelsPerFeeCoinCache{
					Disable: true, // disable cache to avoid goroutine cleanup issues in tests
				},
			}
			pluginConfigBytes, err := json.Marshal(pluginCfg)
			require.NoError(t, err)

			// Unmarshal plugin config bytes into JSONConfig
			var pluginConfigMap map[string]any
			err = json.Unmarshal(pluginConfigBytes, &pluginConfigMap)
			require.NoError(t, err)

			// Create job with gas limit
			jb := job.Job{
				ID:               1,
				ExternalJobID:    uuid.New(),
				OCR2OracleSpecID: int32Ptr(7),
				GasLimit:         tc.jobGasLimit,
				OCR2OracleSpec: &job.OCR2OracleSpec{
					ID:            7,
					ContractID:    "0x1234567890123456789012345678901234567890",
					Relay:         "evm",
					ChainID:       "1",
					PluginConfig:  job.JSONConfig(pluginConfigMap),
					RelayConfig:   job.JSONConfig(map[string]any{}),
					TransmitterID: null.StringFrom("transmitter1"),
				},
				PipelineSpec: &pipeline.Spec{
					ID: 1,
				},
			}

			// Create mock relayer
			mockRelayer := &mockRelayer{}

			// Create OCR2OracleArgs with development mode to skip validation
			argsNoPlugin := libocr.OCR2OracleArgs{
				LocalConfig: ocrtypes.LocalConfig{
					DevelopmentMode: ocrtypes.EnableDangerousDevelopmentMode,
				},
			}

			// Call NewMedianServices
			_, err = NewMedianServices(
				ctx,
				jb,
				false, // isNewlyCreatedJob
				mockRelayer,
				nil, // kv store
				nil, // pipeline runner
				lggr,
				argsNoPlugin,
				&medianConfig{},
				nil, // enhanced telemetry channel
				nil, // error log
			)

			// Verify no error occurred
			require.NoError(t, err, tc.description)

			// Unmarshal the captured plugin config
			var capturedPluginConfig config.PluginConfig
			err = json.Unmarshal(mockRelayer.capturedPluginArgs.PluginConfig, &capturedPluginConfig)
			require.NoError(t, err, "Failed to unmarshal captured plugin config")

			// Verify gas limit
			if tc.expectedGasLimit == nil {
				assert.Nil(t, capturedPluginConfig.GasLimit, "Expected gas limit to be nil: %s", tc.description)
			} else {
				require.NotNil(t, capturedPluginConfig.GasLimit, "Expected gas limit to be set: %s", tc.description)
				assert.Equal(t, *tc.expectedGasLimit, *capturedPluginConfig.GasLimit, "Gas limit mismatch: %s", tc.description)
			}
		})
	}
}

// mockPluginProvider implements commontypes.PluginProvider and commontypes.MedianProvider
type mockPluginProvider struct {
	commontypes.PluginProvider
}

func (m *mockPluginProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return nil
}

func (m *mockPluginProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return nil
}

func (m *mockPluginProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return nil
}

func (m *mockPluginProvider) Start(context.Context) error { return nil }
func (m *mockPluginProvider) Close() error                { return nil }
func (m *mockPluginProvider) Ready() error                { return nil }
func (m *mockPluginProvider) HealthReport() map[string]error {
	return map[string]error{"mock": nil}
}
func (m *mockPluginProvider) Name() string { return "mock" }

// Type assertions
var _ commontypes.MedianProvider = (*mockPluginProvider)(nil)

func (m *mockPluginProvider) ReportCodec() libocr_median.ReportCodec { return nil }
func (m *mockPluginProvider) MedianContract() libocr_median.MedianContract {
	return nil
}
func (m *mockPluginProvider) OnchainConfigCodec() libocr_median.OnchainConfigCodec {
	return nil
}
func (m *mockPluginProvider) ContractReader() commontypes.ContractReader { return nil }
func (m *mockPluginProvider) Codec() commontypes.Codec                   { return nil }

func uint32Ptr(v uint32) *uint32 {
	return &v
}

func int32Ptr(v int32) *int32 {
	return &v
}
