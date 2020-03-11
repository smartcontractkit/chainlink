package store

import (
	"chainlink/core/assets"
	"chainlink/core/eth"
	"chainlink/core/logger"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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

func updatePrometheusEthBalance(balance *assets.Eth, from common.Address) {
	balanceFloat, err := approximateFloat64(balance)

	if err != nil {
		logger.Error(fmt.Errorf("updatePrometheusEthBalance: %v", err))
		return
	}

	promETHBalance.WithLabelValues(from.Hex()).Set(balanceFloat)
}

func approximateFloat64(e *assets.Eth) (float64, error) {
	ef := new(big.Float).SetInt(e.ToInt())
	bf := new(big.Float).Quo(ef, eth.WeiPerEth)
	f64, _ := bf.Float64()
	if f64 == math.Inf(1) || f64 == math.Inf(-1) {
		return math.Inf(1), errors.New("assets.Eth.Float64: Could not approximate Eth value into float")
	}
	return f64, nil
}
