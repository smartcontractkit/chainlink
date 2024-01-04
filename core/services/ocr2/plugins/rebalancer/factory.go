package rebalancer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquidityrebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type PluginFactory struct{}

func NewPluginFactory() *PluginFactory {
	return &PluginFactory{}
}

func (p PluginFactory) NewReportingPlugin(config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[models.ReportMetadata], ocr3types.ReportingPluginInfo, error) {
	var offchainConfig models.PluginConfig
	if err := json.Unmarshal(config.OffchainConfig, &offchainConfig); err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, fmt.Errorf("invalid config: %w", err)
	}

	liquidityRebalancer := liquidityrebalancer.NewDummyRebalancer()
	liquidityGraph := liquiditygraph.NewGraph()
	liquidityManagerFactory := liquiditymanager.NewBaseLiquidityManagerFactory()

	closePluginTimeout := 30 * time.Second
	if offchainConfig.ClosePluginTimeoutSec > 0 {
		closePluginTimeout = time.Duration(offchainConfig.ClosePluginTimeoutSec) * time.Second
	}

	return NewPlugin(
			config.F,
			closePluginTimeout,
			offchainConfig.LiquidityManagerNetwork,
			offchainConfig.LiquidityManagerAddress,
			liquidityManagerFactory,
			liquidityGraph,
			liquidityRebalancer,
		),
		ocr3types.ReportingPluginInfo{
			Name: models.PluginName,
			Limits: ocr3types.ReportingPluginLimits{
				MaxQueryLength:       0,
				MaxObservationLength: 0,
				MaxOutcomeLength:     0,
				MaxReportLength:      0,
				MaxReportCount:       0,
			},
		},
		nil
}
