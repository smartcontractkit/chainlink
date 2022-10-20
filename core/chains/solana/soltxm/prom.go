package soltxm

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// successful transactions
	promSolTxmSuccessTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_success",
		Help: "Number of transactions that are included and successfully executed on chain",
	}, []string{"chainID"})

	// inflight transactions
	promSolTxmPendingTxs = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "solana_txm_tx_pending",
		Help: "Number of transactions that are pending confirmation",
	}, []string{"chainID"})

	// error cases
	promSolTxmErrorTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_error",
		Help: "Number of transactions that have errored across all cases",
	}, []string{"chainID"})
	promSolTxmRevertTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_error_revert",
		Help: "Number of transactions that are included and failed onchain",
	}, []string{"chainID"})
	promSolTxmRejectTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_error_reject",
		Help: "Number of transactions that the RPC immediately rejected",
	}, []string{"chainID"})
	promSolTxmInvalidBlockhash = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_error_invalidBlockhash",
		Help: "Number of transactions that included an invalid blockhash",
	}, []string{"chainID"})
)
