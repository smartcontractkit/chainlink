package monitor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var promTerraBalance = promauto.NewGaugeVec(
	prometheus.GaugeOpts{Name: "terra_balance", Help: "Terra account balances"},
	[]string{"account", "terraChainID", "denomination"},
)

func (b *balanceMonitor) updateProm(acc sdk.AccAddress, bal *sdk.DecCoin) {
	balF, err := bal.Amount.Float64()
	if err != nil {
		b.lggr.Errorw("Failed to convert balance to float", "err", err)
		return
	}
	promTerraBalance.WithLabelValues(acc.String(), b.chainID, bal.GetDenom()).Set(balF)
}
