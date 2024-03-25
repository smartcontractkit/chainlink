package mercury_common_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	mercuryv1test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v1/test"
	mercuryv2test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v2/test"
	mercuryv3test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v3/test"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	mercuryv1types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	mercuryv2types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	mercuryv3types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

func PluginMercury(t *testing.T, p types.PluginMercury) {
	PluginMercuryTest{MercuryProvider}.TestPluginMercury(t, p)
}

type PluginMercuryTest struct {
	types.MercuryProvider
}

func (m PluginMercuryTest) TestPluginMercury(t *testing.T, p types.PluginMercury) {
	t.Run("PluginMercuryV3", func(t *testing.T) {
		ctx := tests.Context(t)
		factory, err := p.NewMercuryV3Factory(ctx, m.MercuryProvider, mercuryv3test.DataSource)
		require.NoError(t, err)
		require.NotNil(t, factory)

		MercuryPluginFactory(t, factory)
	})

	t.Run("PluginMercuryV2", func(t *testing.T) {
		ctx := tests.Context(t)
		factory, err := p.NewMercuryV2Factory(ctx, m.MercuryProvider, mercuryv2test.DataSource)
		require.NoError(t, err)
		require.NotNil(t, factory)

		MercuryPluginFactory(t, factory)
	})

	t.Run("PluginMercuryV1", func(t *testing.T) {
		ctx := tests.Context(t)
		factory, err := p.NewMercuryV1Factory(ctx, m.MercuryProvider, mercuryv1test.DataSource)
		require.NoError(t, err)
		require.NotNil(t, factory)

		MercuryPluginFactory(t, factory)
	})
}

var FactoryServer = staticMercuryServer{
	provider:     MercuryProvider,
	dataSourceV1: mercuryv1test.DataSource,
	dataSourceV2: mercuryv2test.DataSource,
	dataSourceV3: mercuryv3test.DataSource,
}

var _ types.PluginMercury = staticMercuryServer{}

type staticMercuryServer struct {
	provider     staticMercuryProvider
	dataSourceV1 mercuryv1test.DataSourceEvaluator
	dataSourceV2 mercuryv2test.DataSourceEvaluator
	dataSourceV3 mercuryv3test.DataSourceEvaluator
}

var _ types.PluginMercury = staticMercuryServer{}

func (s staticMercuryServer) commonValidation(ctx context.Context, provider types.MercuryProvider) error {
	ocd := provider.OffchainConfigDigester()
	err := s.provider.offchainDigester.Evaluate(ctx, ocd)
	if err != nil {
		return fmt.Errorf("failed to evaluate offchainDigester: %w", err)
	}

	cct := provider.ContractConfigTracker()
	err = s.provider.contractTracker.Evaluate(ctx, cct)
	if err != nil {
		return fmt.Errorf("failed to evaluate contractTracker: %w", err)
	}

	ct := provider.ContractTransmitter()
	err = s.provider.contractTransmitter.Evaluate(ctx, ct)
	if err != nil {
		return fmt.Errorf("failed to evaluate contractTransmitter: %w", err)
	}

	occ := provider.OnchainConfigCodec()
	err = s.provider.onchainConfigCodec.Evaluate(ctx, occ)
	if err != nil {
		return fmt.Errorf("failed to evaluate onchainConfigCodec: %w", err)
	}
	return nil
}

func (s staticMercuryServer) NewMercuryV3Factory(ctx context.Context, provider types.MercuryProvider, dataSource mercuryv3types.DataSource) (types.MercuryPluginFactory, error) {
	var err error
	defer func() {
		if err != nil {
			panic(fmt.Sprintf("provider %v, %T: %s", provider, provider, err))
		}
	}()
	err = s.commonValidation(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed commonValidation: %w", err)
	}

	rc := provider.ReportCodecV3()
	err = s.provider.reportCodecV3.Evaluate(ctx, rc)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate reportCodecV3: %w", err)
	}

	err = s.dataSourceV3.Evaluate(ctx, dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate dataSource: %w", err)
	}

	return staticMercuryPluginFactory{}, nil
}

func (s staticMercuryServer) NewMercuryV2Factory(ctx context.Context, provider types.MercuryProvider, dataSource mercuryv2types.DataSource) (types.MercuryPluginFactory, error) {
	var err error
	defer func() {
		if err != nil {
			panic(fmt.Sprintf("provider %v, %T: %s", provider, provider, err))
		}
	}()
	err = s.commonValidation(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed commonValidation: %w", err)
	}

	rc := provider.ReportCodecV2()
	err = s.provider.reportCodecV2.Evaluate(ctx, rc)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate reportCodecV2: %w", err)
	}

	err = s.dataSourceV2.Evaluate(ctx, dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate dataSource: %w", err)
	}
	return staticMercuryPluginFactory{}, nil
}

func (s staticMercuryServer) NewMercuryV1Factory(ctx context.Context, provider types.MercuryProvider, dataSource mercuryv1types.DataSource) (types.MercuryPluginFactory, error) {
	var err error
	defer func() {
		if err != nil {
			panic(fmt.Sprintf("provider %v, %T: %s", provider, provider, err))
		}
	}()
	err = s.commonValidation(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed commonValidation: %w", err)
	}

	rc := provider.ReportCodecV1()
	err = s.provider.reportCodecV1.Evaluate(ctx, rc)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate reportCodecV1: %w", err)
	}

	err = s.dataSourceV1.Evaluate(ctx, dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate dataSource: %w", err)
	}

	return staticMercuryPluginFactory{}, nil
}

type staticMercuryPluginFactory struct{}

func (s staticMercuryPluginFactory) Name() string { panic("implement me") }

func (s staticMercuryPluginFactory) Start(ctx context.Context) error { return nil }

func (s staticMercuryPluginFactory) Close() error { return nil }

func (s staticMercuryPluginFactory) Ready() error { panic("implement me") }

func (s staticMercuryPluginFactory) HealthReport() map[string]error { panic("implement me") }

func (s staticMercuryPluginFactory) NewMercuryPlugin(ctx context.Context, config ocr3types.MercuryPluginConfig) (ocr3types.MercuryPlugin, ocr3types.MercuryPluginInfo, error) {
	if config.ConfigDigest != mercuryPluginConfig.ConfigDigest {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected ConfigDigest %x but got %x", mercuryPluginConfig.ConfigDigest, config.ConfigDigest)
	}
	if config.OracleID != mercuryPluginConfig.OracleID {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected OracleID %d but got %d", mercuryPluginConfig.OracleID, config.OracleID)
	}
	if config.F != mercuryPluginConfig.F {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected F %d but got %d", mercuryPluginConfig.F, config.F)
	}
	if config.N != mercuryPluginConfig.N {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected N %d but got %d", mercuryPluginConfig.N, config.N)
	}
	if !bytes.Equal(config.OnchainConfig, mercuryPluginConfig.OnchainConfig) {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected OnchainConfig %x but got %x", mercuryPluginConfig.OnchainConfig, config.OnchainConfig)
	}
	if !bytes.Equal(config.OffchainConfig, mercuryPluginConfig.OffchainConfig) {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected OffchainConfig %x but got %x", mercuryPluginConfig.OffchainConfig, config.OffchainConfig)
	}
	if config.EstimatedRoundInterval != mercuryPluginConfig.EstimatedRoundInterval {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected EstimatedRoundInterval %d but got %d", mercuryPluginConfig.EstimatedRoundInterval, config.EstimatedRoundInterval)
	}

	if config.MaxDurationObservation != mercuryPluginConfig.MaxDurationObservation {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected MaxDurationObservation %d but got %d", mercuryPluginConfig.MaxDurationObservation, config.MaxDurationObservation)
	}

	return OCR3Plugin, mercuryPluginInfo, nil
}

func MercuryPluginFactory(t *testing.T, factory types.MercuryPluginFactory) {
	expectedMercuryPlugin := OCR3Plugin
	t.Run("ReportingPluginFactory", func(t *testing.T) {
		ctx := tests.Context(t)
		rp, gotRPI, err := factory.NewMercuryPlugin(ctx, mercuryPluginConfig)
		require.NoError(t, err)
		assert.Equal(t, mercuryPluginInfo, gotRPI)
		t.Cleanup(func() { assert.NoError(t, rp.Close()) })
		t.Run("ReportingPlugin", func(t *testing.T) {
			expectedMercuryPlugin.AssertEqual(ctx, t, rp)
		})
	})
}
