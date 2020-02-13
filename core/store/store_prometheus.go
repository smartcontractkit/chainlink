package store

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var promETHBalance = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "eth_balance",
		Help: "Each Ethereum account's balance",
	},
	[]string{"account"},
)
