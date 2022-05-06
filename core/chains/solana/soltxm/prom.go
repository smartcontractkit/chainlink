package soltxm

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promSolTxmSuccessfulTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_successful",
		Help: "Number of transactions that are included and successfully executed on chain",
	}, []string{"chainID"})
	promSolTxmRevertedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_reverted",
		Help: "Number of transactions that are included but reverted on chain",
	}, []string{"chainID"})
	promSolTxmFailedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_failed",
		Help: "Number of transactions that failed sending to chain or simulated with an unrecognized error.",
	}, []string{"chainID"})
	promSolTxmTimedOutTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_timedout",
		Help: "Number of transactions that timed out during tx retry and were not confirmed within the timeout. Note: processed txs may still be included onchain",
	}, []string{"chainID"})
	promSolTxmInflightTxs = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "solana_txm_tx_inflight",
		Help: "Number of transactions that are currently being retried and confirmed",
	}, []string{"chainID"})
)
