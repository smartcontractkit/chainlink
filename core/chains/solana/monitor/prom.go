package monitor

import (
	"github.com/gagliardetto/solana-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	solanaRelay "github.com/smartcontractkit/chainlink-solana/pkg/solana"
)

var promSolanaBalance = promauto.NewGaugeVec(
	prometheus.GaugeOpts{Name: "solana_balance", Help: "Solana account balances"},
	[]string{"account", "chainID", "chainSet", "denomination"},
)

func (b *balanceMonitor) updateProm(acc solana.PublicKey, lamports uint64) {
	v := solanaRelay.LamportsToSol(lamports) // convert from lamports to SOL
	promSolanaBalance.WithLabelValues(acc.String(), b.chainID, "solana", "SOL").Set(v)
}
