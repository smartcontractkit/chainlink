package monitor

import (
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/internal"
)

var (
	promSolanaBalance = promauto.NewGaugeVec(
		prometheus.GaugeOpts{Name: "solana_balance", Help: "Solana account balances"},
		[]string{"account", "chainID", "chainSet", "denomination"},
	)
	promCacheTimestamp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{Name: "solana_cache_last_update_unix", Help: "Solana relayer cache last update timestamp"},
		[]string{"type", "chainID", "account"},
	)
	promClientReq = promauto.NewGaugeVec(
		prometheus.GaugeOpts{Name: "solana_client_latency_ms", Help: "Solana client request latency"},
		[]string{"request", "url"},
	)
)

func (b *balanceMonitor) updateProm(acc solana.PublicKey, lamports uint64) {
	v := internal.LamportsToSol(lamports) // convert from lamports to SOL
	promSolanaBalance.WithLabelValues(acc.String(), b.chainID, "solana", "SOL").Set(v)
}

func SetCacheTimestamp(t time.Time, cacheType, chainID, account string) {
	promCacheTimestamp.With(prometheus.Labels{
		"type":    cacheType,
		"chainID": chainID,
		"account": account,
	}).Set(float64(t.Unix()))
}

func SetClientLatency(d time.Duration, request, url string) {
	promClientReq.With(prometheus.Labels{
		"request": request,
		"url":     url,
	}).Set(float64(d.Milliseconds()))
}

func GetClientLatency(request, url string) (prometheus.Gauge, error) {
	return promClientReq.GetMetricWith(prometheus.Labels{
		"request": request,
		"url":     url,
	})
}
