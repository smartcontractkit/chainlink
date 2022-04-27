package soltxm

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promSolTxmSuccessfulTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "soltxm_num_tx_successful",
		Help: "Number of transactions that are included and successfully executed on chain",
	}, []string{"chainID"})
	promSolTxmRevertedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "soltxm_num_tx_reverted",
		Help: "Number of transactions that are included but reverted on chain",
	}, []string{"chainID"})
	promSolTxmFailedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "soltxm_num_tx_failed",
		Help: "Number of transactions that are failed sending to chain or failed simulation. Note that txs that failed simulation could still be included onchain",
	}, []string{"chainID"})
	promSolTxmTimedOutTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "soltxm_num_tx_timedout",
		Help: "Number of transactions that timed out during tx retry (never included or never moved beyond 'processed'). Note processed txs may still be included onchain",
	}, []string{"chainID"})
	promSolTxmInflightTxs = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "soltxm_num_tx_inflight",
		Help: "Number of transactions that are currently being retried and confirmed",
	}, []string{"chainID"})
)
