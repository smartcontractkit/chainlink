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
	promSolTxmDropTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_error_drop",
		Help: "Number of transactions that timed out during confirmation. Note: tx is likely dropped from the chain, but may still be included.",
	}, []string{"chainID"})
	promSolTxmSimRevertTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_error_sim_revert",
		Help: "Number of transactions that reverted during simulation. Note: tx may still be included onchain",
	}, []string{"chainID"})
	promSolTxmSimOtherTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solana_txm_tx_error_sim_other",
		Help: "Number of transactions that failed simulation with an unrecognized error. Note: tx may still be included onchain",
	}, []string{"chainID"})
)
