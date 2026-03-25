package chainlink

import (
	"context"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

const (
	evmSkipReasonChainConstruct     = "chain_construct"
	evmSkipReasonRelayerInit        = "relayer_init"
	evmSkipReasonUnknownSelector    = "unknown_chain_selector"
	evmSkipReasonLoopPluginRegister = "loop_plugin_register"
)

var (
	evmChainConfigSkipped = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evm_chain_config_skipped_total",
			Help: "Enabled EVM chains skipped during startup (e.g. missing chain-selectors entry or relayer init failure).",
		},
		[]string{"reason"},
	)
	evmChainConfigDegraded = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "evm_chain_config_degraded",
			Help: "1 if at least one enabled EVM chain was skipped at startup; 0 otherwise.",
		},
	)

	beholderEVMSkippedOnce sync.Once
	beholderEVMSkipped     func(context.Context, int64)
)

func incEVMChainConfigSkipped(reason string) {
	evmChainConfigSkipped.WithLabelValues(reason).Inc()
	evmChainConfigDegraded.Set(1)

	beholderEVMSkippedOnce.Do(func() {
		c, err := beholder.GetMeter().Int64Counter("platform_node_evm_chain_config_skipped_total")
		if err != nil {
			return
		}
		beholderEVMSkipped = func(ctx context.Context, n int64) {
			c.Add(ctx, n)
		}
	})
	if beholderEVMSkipped != nil {
		beholderEVMSkipped(context.Background(), 1)
	}
}
