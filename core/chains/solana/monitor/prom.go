package monitor

import (
	"github.com/gagliardetto/solana-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var promTerraBalance = promauto.NewGaugeVec(
	prometheus.GaugeOpts{Name: "terra_balance", Help: "Terra account balances"},
	[]string{"account", "terraChainID", "denomination"},
)

func (b *balanceMonitor) updateProm(acc solana.PublicKey, lamports uint64) {
	var v float64
	v = float64(lamports) / 1_000_000_000 // convert from lamports to SOL
	promTerraBalance.WithLabelValues(acc.String(), b.chainID, v)
}
