package ring

import (
	"context"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

func TestFactory_NewFactory(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	shardOrchestratorStore := shardorchestrator.NewStore(lggr)
	arbiter := &mockArbiter{}

	tests := []struct {
		name      string
		arbiter   ringpb.ArbiterScalerClient
		config    *ConsensusConfig
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "with_nil_config",
			arbiter: arbiter,
			config:  nil,
			wantErr: false,
		},
		{
			name:    "with_custom_config",
			arbiter: arbiter,
			config:  &ConsensusConfig{BatchSize: 50},
			wantErr: false,
		},
		{
			name:      "nil_arbiter_returns_error",
			arbiter:   nil,
			config:    nil,
			wantErr:   true,
			errSubstr: "arbiterScaler is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFactory(store, shardOrchestratorStore, tt.arbiter, lggr, tt.config)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, f)
			}
		})
	}
}

func TestFactory_NewReportingPlugin(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	f, err := NewFactory(store, nil, &mockArbiter{}, lggr, nil)
	require.NoError(t, err)

	config := ocr3types.ReportingPluginConfig{N: 4, F: 1}
	plugin, info, err := f.NewReportingPlugin(context.Background(), config)
	require.NoError(t, err)
	require.NotNil(t, plugin)
	require.NotEmpty(t, info.Name)
	require.Equal(t, "RingPlugin", info.Name)
	require.Equal(t, defaultMaxReportCount, info.Limits.MaxReportCount)
}

func TestFactory_Lifecycle(t *testing.T) {
	lggr := logger.Test(t)
	store := NewStore()
	f, err := NewFactory(store, nil, &mockArbiter{}, lggr, nil)
	require.NoError(t, err)

	err = f.Start(context.Background())
	require.NoError(t, err)

	name := f.Name()
	require.NotEmpty(t, name)

	report := f.HealthReport()
	require.NotNil(t, report)
	require.Contains(t, report, name)

	err = f.Close()
	require.NoError(t, err)
}
