package mercury_test

import (
	"context"
	"errors"
	"os/exec"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	v2 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	v4 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v4"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

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
	v4FeedId = [32]uint8{00, 04, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}

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

	v4jsonCfg = job.JSONConfig{
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
			Relay:      relay.NetworkEVM,
			ChainID:    "1",
		},
		PipelineSpec:   &pipeline.Spec{},
		PipelineSpecID: int32(1),
	}

	// this is kind of gross, but it's the best way to test return values of the services
	expectedEmbeddedServiceCnt = 3
	expectedLoopServiceCnt     = expectedEmbeddedServiceCnt + 2 // factory server and loop unregisterer
)

func TestNewServices(t *testing.T) {
	type args struct {
		pluginConfig job.JSONConfig
		feedID       utils.FeedID
		cfg          mercuryocr2.Config
	}
	testCases := []struct {
		name            string
		args            args
		loopMode        bool
		wantLoopFactory any
		wantServiceCnt  int
		wantErr         bool
		wantErrStr      string
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
			name: "v4 legacy",
			args: args{
				pluginConfig: v4jsonCfg,
				feedID:       v4FeedId,
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
		{
			name:     "v3 loop err",
			loopMode: true,
			args: args{
				pluginConfig: v3jsonCfg,
				feedID:       v3FeedId,
				cfg:          mercuryocr2.NewMercuryConfig(1, 1, &testRegistrarConfig{failRegister: true}),
			},
			wantServiceCnt:  expectedLoopServiceCnt,
			wantErr:         true,
			wantLoopFactory: &loop.MercuryV3Service{},
			wantErrStr:      "failed to init loop for feed",
		},
		{
			name:     "v4 loop",
			loopMode: true,
			args: args{
				pluginConfig: v4jsonCfg,
				feedID:       v4FeedId,
			},
			wantServiceCnt:  expectedLoopServiceCnt,
			wantErr:         false,
			wantLoopFactory: &loop.MercuryV4Service{},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.loopMode {
				t.Setenv(string(env.MercuryPlugin.Cmd), "fake_cmd")
				assert.NotEmpty(t, env.MercuryPlugin.Cmd.Get())
			}
			// use default config if not provided
			if tt.args.cfg == nil {
				tt.args.cfg = testCfg
			}
			got, err := newServicesTestWrapper(t, tt.args.pluginConfig, tt.args.feedID, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.wantErrStr != "" {
					assert.Contains(t, err.Error(), tt.wantErrStr)
				}
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

	t.Run("restartable loop", func(t *testing.T) {
		// setup a real loop registry to test restartability
		registry := plugins.NewTestLoopRegistry(logger.TestLogger(t))
		loopRegistrarConfig := plugins.NewRegistrarConfig(loop.GRPCOpts{}, registry.Register, registry.Unregister)
		prodCfg := mercuryocr2.NewMercuryConfig(1, 1, loopRegistrarConfig)
		type args struct {
			pluginConfig job.JSONConfig
			feedID       utils.FeedID
			cfg          mercuryocr2.Config
		}
		testCases := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{
				name: "v1 loop",
				args: args{
					pluginConfig: v1jsonCfg,
					feedID:       v1FeedId,
					cfg:          prodCfg,
				},
				wantErr: false,
			},
			{
				name: "v2 loop",
				args: args{
					pluginConfig: v2jsonCfg,
					feedID:       v2FeedId,
					cfg:          prodCfg,
				},
				wantErr: false,
			},
			{
				name: "v3 loop",
				args: args{
					pluginConfig: v3jsonCfg,
					feedID:       v3FeedId,
					cfg:          prodCfg,
				},
				wantErr: false,
			},
			{
				name: "v4 loop",
				args: args{
					pluginConfig: v4jsonCfg,
					feedID:       v4FeedId,
					cfg:          prodCfg,
				},
				wantErr: false,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				t.Setenv(string(env.MercuryPlugin.Cmd), "fake_cmd")
				assert.NotEmpty(t, env.MercuryPlugin.Cmd.Get())

				got, err := newServicesTestWrapper(t, tt.args.pluginConfig, tt.args.feedID, tt.args.cfg)
				if (err != nil) != tt.wantErr {
					t.Errorf("NewServices() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// hack to simulate a restart. we don't have enough boilerplate to start the oracle service
				// only care about the subservices so we start all except the oracle, which happens to be the last one
				for i := 0; i < len(got)-1; i++ {
					require.NoError(t, got[i].Start(tests.Context(t)))
				}
				// if we don't close the services, we get conflicts with the loop registry
				_, err = newServicesTestWrapper(t, tt.args.pluginConfig, tt.args.feedID, tt.args.cfg)
				require.ErrorContains(t, err, "plugin already registered")

				// close all services and try again
				for i := len(got) - 2; i >= 0; i-- {
					require.NoError(t, got[i].Close())
				}
				_, err = newServicesTestWrapper(t, tt.args.pluginConfig, tt.args.feedID, tt.args.cfg)
				require.NoError(t, err)
			})
		}
	})
}

// we are only varying the version via feedID (and the plugin config)
// this wrapper supplies dummy values for the rest of the arguments
func newServicesTestWrapper(t *testing.T, pluginConfig job.JSONConfig, feedID utils.FeedID, cfg mercuryocr2.Config) ([]job.ServiceCtx, error) {
	t.Helper()
	jb := testJob
	jb.OCR2OracleSpec.PluginConfig = pluginConfig
	return mercuryocr2.NewServices(jb, &testProvider{}, nil, logger.TestLogger(t), testArgsNoPlugin, cfg, nil, &testDataSourceORM{}, feedID, false)
}

type testProvider struct{}

// ChainReader implements types.MercuryProvider.
func (*testProvider) ContractReader() commontypes.ContractReader { panic("unimplemented") }

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

// ReportCodecV4 implements types.MercuryProvider.
func (*testProvider) ReportCodecV4() v4.ReportCodec { return nil }

// Start implements types.MercuryProvider.
func (*testProvider) Start(context.Context) error { return nil }

var _ commontypes.MercuryProvider = (*testProvider)(nil)

type testRegistrarConfig struct {
	failRegister bool
}

func (c *testRegistrarConfig) UnregisterLOOP(ID string) {}

// RegisterLOOP implements plugins.RegistrarConfig.
func (c *testRegistrarConfig) RegisterLOOP(config plugins.CmdConfig) (func() *exec.Cmd, loop.GRPCOpts, error) {
	if c.failRegister {
		return nil, loop.GRPCOpts{}, errors.New("failed to register")
	}
	return nil, loop.GRPCOpts{}, nil
}

var _ plugins.RegistrarConfig = (*testRegistrarConfig)(nil)

type testDataSourceORM struct{}

// LatestReport implements types.DataSourceORM.
func (*testDataSourceORM) LatestReport(ctx context.Context, feedID [32]byte) (report []byte, err error) {
	return []byte{1, 2, 3}, nil
}

var _ types.DataSourceORM = (*testDataSourceORM)(nil)
