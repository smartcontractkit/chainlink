package liquiditymanager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge"
	liquiditymanager "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/chain/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/discoverer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/rebalalgo"
)

const (
	PluginName = "LiquidityManager"

	// OCR limits

	// maxQueryLength should be 0 as no queries are performed
	maxQueryLength = 0
	// maxObservationLength should be 1M bytes
	maxObservationLength = 1_000_000
	// maxOutcomeLength should be 1M bytes
	maxOutcomeLength = 1_000_000
	// maxReportLength should be 1M bytes
	maxReportLength = 1_000_000
	// maxReportCount should be 100
	maxReportCount = 100
)

type PluginFactory struct {
	lggr              logger.Logger
	config            models.PluginConfig
	lmFactory         liquiditymanager.Factory
	discovererFactory discoverer.Factory
	bridgeFactory     bridge.Factory
}

func NewPluginFactory(
	lggr logger.Logger,
	pluginConfigBytes []byte,
	lmFactory liquiditymanager.Factory,
	discovererFactory discoverer.Factory,
	bridgeFactory bridge.Factory,
) (*PluginFactory, error) {
	var pluginConfig models.PluginConfig
	if err := json.Unmarshal(pluginConfigBytes, &pluginConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin config: %w", err)
	}
	return &PluginFactory{
		lggr:              lggr.Named(PluginName),
		config:            pluginConfig,
		lmFactory:         lmFactory,
		discovererFactory: discovererFactory,
		bridgeFactory:     bridgeFactory,
	}, nil
}

func (p PluginFactory) buildRebalancer() (rebalalgo.RebalancingAlgo, error) {
	switch p.config.RebalancerConfig.Type {
	case models.RebalancerTypePingPong:
		return rebalalgo.NewPingPong(), nil
	case models.RebalancerTypeMinLiquidity:
		return rebalalgo.NewMinLiquidityRebalancer(p.lggr), nil
	case models.RebalancerTypeTargetAndMin:
		return rebalalgo.NewTargetMinBalancer(p.lggr, p.config), nil
	default:
		return nil, fmt.Errorf("invalid rebalancer type %s", p.config.RebalancerConfig.Type)
	}
}

func (p PluginFactory) NewReportingPlugin(config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[models.Report], ocr3types.ReportingPluginInfo, error) {
	liquidityRebalancer, err := p.buildRebalancer()
	if err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, fmt.Errorf("failed to build rebalancer: %w", err)
	}

	closePluginTimeout := 30 * time.Second
	if p.config.ClosePluginTimeoutSec > 0 {
		closePluginTimeout = time.Duration(p.config.ClosePluginTimeoutSec) * time.Second
	}

	discoverer, err := p.discovererFactory.NewDiscoverer(p.config.LiquidityManagerNetwork, p.config.LiquidityManagerAddress)
	if err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, fmt.Errorf("init discoverer: %w", err)
	}

	return NewPlugin(
			config.F,
			closePluginTimeout,
			p.config.LiquidityManagerNetwork,
			p.config.LiquidityManagerAddress,
			p.lmFactory,
			discoverer,
			p.bridgeFactory,
			liquidityRebalancer,
			liquiditymanager.NewEvmReportCodec(),
			p.lggr,
		),
		ocr3types.ReportingPluginInfo{
			Name: models.PluginName,
			Limits: ocr3types.ReportingPluginLimits{
				MaxQueryLength:       maxQueryLength,
				MaxObservationLength: maxObservationLength,
				MaxOutcomeLength:     maxOutcomeLength,
				MaxReportLength:      maxReportLength,
				MaxReportCount:       maxReportCount,
			},
		},
		nil
}
