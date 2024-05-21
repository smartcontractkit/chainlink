package mercury_test

import (
	"context"
	"os/exec"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	v2 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"

	mercuryocr2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury"

	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	libocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

var (
	v1FeedId = [32]uint8{00, 01, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
	v2FeedId = [32]uint8{00, 02, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
	v3FeedId = [32]uint8{00, 03, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}

	testArgsNoPlugin = libocr2.MercuryOracleArgs{
		LocalConfig: libocr2types.LocalConfig{
			DevelopmentMode: libocr2types.EnableDangerousDevelopmentMode,
		},
	}

	testCfg = mercuryocr2.NewMercuryConfig(1, 1, &testRegistrarConfig{})

	v1jsonCfg = job.JSONConfig{
		"serverURL":          "example.com:80",
		"serverPubKey":       "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93",
		"initialBlockNumber": 1234,
	}

	v2jsonCfg = job.JSONConfig{
		"serverURL":    "example.com:80",
		"serverPubKey": "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93",
		"linkFeedID":   "0x00026b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472",
		"nativeFeedID": "0x00036b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472",
	}

	v3jsonCfg = job.JSONConfig{
		"serverURL":    "example.com:80",
		"serverPubKey": "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93",
		"linkFeedID":   "0x00026b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472",
		"nativeFeedID": "0x00036b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472",
	}

	testJob = job.Job{
		ID:               1,
		ExternalJobID:    uuid.Must(uuid.NewRandom()),
		OCR2OracleSpecID: ptr(int32(7)),
		OCR2OracleSpec: &job.OCR2OracleSpec{
			ID:         7,
			ContractID: "phony",
			FeedID:     ptr(common.BytesToHash([]byte{1, 2, 3})),
			Relay:      commontypes.NetworkEVM,
			ChainID:    "1",
		},
		PipelineSpec:   &pipeline.Spec{},
		PipelineSpecID: int32(1),
	}

	// this is kind of gross, but it's the best way to test return values of the services
	expectedEmbeddedServiceCnt = 3
	expectedLoopServiceCnt     = expectedEmbeddedServiceCnt + 1
)

func TestNewServices(t *testing.T) {
	type args struct {
		pluginConfig job.JSONConfig
		feedID       utils.FeedID
	}
	tests := []struct {
		name            string
		args            args
		loopMode        bool
		wantLoopFactory any
		wantServiceCnt  int
		wantErr         bool
	}{
		{
			name: "no plugin config error ",
			args: args{
				feedID: v1FeedId,
			},
			wantServiceCnt: 0,
			wantErr:        true,
		},

		{
			name: "v1 legacy",
			args: args{
				pluginConfig: v1jsonCfg,
				feedID:       v1FeedId,
			},
			wantServiceCnt: expectedEmbeddedServiceCnt,
			wantErr:        false,
		},
		{
			name: "v2 legacy",
			args: args{
				pluginConfig: v2jsonCfg,
				feedID:       v2FeedId,
			},
			wantServiceCnt: expectedEmbeddedServiceCnt,
			wantErr:        false,
		},
		{
			name: "v3 legacy",
			args: args{
				pluginConfig: v3jsonCfg,
				feedID:       v3FeedId,
			},
			wantServiceCnt: expectedEmbeddedServiceCnt,
			wantErr:        false,
		},
		{
			name:     "v1 loop",
			loopMode: true,
			args: args{
				pluginConfig: v1jsonCfg,
				feedID:       v1FeedId,
			},
			wantServiceCnt:  expectedLoopServiceCnt,
			wantErr:         false,
			wantLoopFactory: &loop.MercuryV1Service{},
		},
		{
			name:     "v2 loop",
			loopMode: true,
			args: args{
				pluginConfig: v2jsonCfg,
				feedID:       v2FeedId,
			},
			wantServiceCnt:  expectedLoopServiceCnt,
			wantErr:         false,
			wantLoopFactory: &loop.MercuryV2Service{},
		},
		{
			name:     "v3 loop",
			loopMode: true,
			args: args{
				pluginConfig: v3jsonCfg,
				feedID:       v3FeedId,
			},
			wantServiceCnt:  expectedLoopServiceCnt,
			wantErr:         false,
			wantLoopFactory: &loop.MercuryV3Service{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.loopMode {
				t.Setenv(string(env.MercuryPlugin.Cmd), "fake_cmd")
				assert.NotEmpty(t, env.MercuryPlugin.Cmd.Get())
			}
			got, err := newServicesTestWrapper(t, tt.args.pluginConfig, tt.args.feedID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, got, tt.wantServiceCnt)
			if tt.loopMode {
				foundLoopFactory := false
				for i := 0; i < len(got); i++ {
					if reflect.TypeOf(got[i]) == reflect.TypeOf(tt.wantLoopFactory) {
						foundLoopFactory = true
						break
					}
				}
				assert.True(t, foundLoopFactory)
			}
		})
	}
}

// we are only varying the version via feedID (and the plugin config)
// this wrapper supplies dummy values for the rest of the arguments
func newServicesTestWrapper(t *testing.T, pluginConfig job.JSONConfig, feedID utils.FeedID) ([]job.ServiceCtx, error) {
	t.Helper()
	jb := testJob
	jb.OCR2OracleSpec.PluginConfig = pluginConfig
	return mercuryocr2.NewServices(jb, &testProvider{}, nil, logger.TestLogger(t), testArgsNoPlugin, testCfg, nil, &testDataSourceORM{}, feedID)
}

type testProvider struct{}

// ChainReader implements types.MercuryProvider.
func (*testProvider) ChainReader() commontypes.ChainReader { panic("unimplemented") }

// Close implements types.MercuryProvider.
func (*testProvider) Close() error { return nil }

// Codec implements types.MercuryProvider.
func (*testProvider) Codec() commontypes.Codec { panic("unimplemented") }

// ContractConfigTracker implements types.MercuryProvider.
func (*testProvider) ContractConfigTracker() libocr2types.ContractConfigTracker {
	panic("unimplemented")
}

// ContractTransmitter implements types.MercuryProvider.
func (*testProvider) ContractTransmitter() libocr2types.ContractTransmitter {
	panic("unimplemented")
}

// HealthReport implements types.MercuryProvider.
func (*testProvider) HealthReport() map[string]error { panic("unimplemented") }

// MercuryChainReader implements types.MercuryProvider.
func (*testProvider) MercuryChainReader() mercury.ChainReader { return nil }

// MercuryServerFetcher implements types.MercuryProvider.
func (*testProvider) MercuryServerFetcher() mercury.ServerFetcher { return nil }

// Name implements types.MercuryProvider.
func (*testProvider) Name() string { panic("unimplemented") }

// OffchainConfigDigester implements types.MercuryProvider.
func (*testProvider) OffchainConfigDigester() libocr2types.OffchainConfigDigester {
	panic("unimplemented")
}

// OnchainConfigCodec implements types.MercuryProvider.
func (*testProvider) OnchainConfigCodec() mercury.OnchainConfigCodec {
	return nil
}

// Ready implements types.MercuryProvider.
func (*testProvider) Ready() error { panic("unimplemented") }

// ReportCodecV1 implements types.MercuryProvider.
func (*testProvider) ReportCodecV1() v1.ReportCodec { return nil }

// ReportCodecV2 implements types.MercuryProvider.
func (*testProvider) ReportCodecV2() v2.ReportCodec { return nil }

// ReportCodecV3 implements types.MercuryProvider.
func (*testProvider) ReportCodecV3() v3.ReportCodec { return nil }

// Start implements types.MercuryProvider.
func (*testProvider) Start(context.Context) error { panic("unimplemented") }

var _ commontypes.MercuryProvider = (*testProvider)(nil)

type testRegistrarConfig struct{}

func (c *testRegistrarConfig) UnregisterLOOP(ID string) {}

// RegisterLOOP implements plugins.RegistrarConfig.
func (*testRegistrarConfig) RegisterLOOP(config plugins.CmdConfig) (func() *exec.Cmd, loop.GRPCOpts, error) {
	return nil, loop.GRPCOpts{}, nil
}

var _ plugins.RegistrarConfig = (*testRegistrarConfig)(nil)

type testDataSourceORM struct{}

// LatestReport implements types.DataSourceORM.
func (*testDataSourceORM) LatestReport(ctx context.Context, feedID [32]byte) (report []byte, err error) {
	return []byte{1, 2, 3}, nil
}

var _ types.DataSourceORM = (*testDataSourceORM)(nil)
