package monitor

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var promTerraBalance = promauto.NewGaugeVec(
	prometheus.GaugeOpts{Name: "solana_balance", Help: "Solana account balances"},
	[]string{"account", "solanaChainID", "balanceSOL"},
)

func (b *balanceMonitor) updateProm(acc solana.PublicKey, lamports uint64) {
	var v float64
	v = float64(lamports) / 1_000_000_000 // convert from lamports to SOL
	promTerraBalance.WithLabelValues(acc.String(), b.chainID, fmt.Sprintf("%.9f", v))
}
