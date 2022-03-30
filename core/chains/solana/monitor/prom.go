package monitor

import (
	"github.com/gagliardetto/solana-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var promSolanaBalance = promauto.NewGaugeVec(
	prometheus.GaugeOpts{Name: "solana_balance", Help: "Solana account balances"},
	[]string{"account", "chainID", "chainSet", "denomination"},
)

func (b *balanceMonitor) updateProm(acc solana.PublicKey, lamports uint64) {
	// TODO: These kinds of converting utilities should be exposed by the `chainlink-solana` package.
	v := float64(lamports) / 1_000_000_000 // convert from lamports to SOL
	promSolanaBalance.WithLabelValues(acc.String(), b.chainID, "solana", "SOL").Set(v)
}
