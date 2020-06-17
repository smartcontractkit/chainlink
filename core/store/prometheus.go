package store

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"

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

func promUpdateEthBalance(balance *assets.Eth, from common.Address) {
	balanceFloat, err := approximateFloat64(balance)

	if err != nil {
		logger.Error(fmt.Errorf("updatePrometheusEthBalance: %v", err))
		return
	}

	promETHBalance.WithLabelValues(from.Hex()).Set(balanceFloat)
}

func approximateFloat64(e *assets.Eth) (float64, error) {
	ef := new(big.Float).SetInt(e.ToInt())
	weif := new(big.Float).SetInt(models.WeiPerEth)
	bf := new(big.Float).Quo(ef, weif)
	f64, _ := bf.Float64()
	if f64 == math.Inf(1) || f64 == math.Inf(-1) {
		return math.Inf(1), errors.New("assets.Eth.Float64: Could not approximate Eth value into float")
	}
	return f64, nil
}
