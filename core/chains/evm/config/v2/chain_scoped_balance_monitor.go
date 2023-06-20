package v2

type balanceMonitorConfig struct {
	c BalanceMonitor
}

func (b *balanceMonitorConfig) Enabled() bool {
	return *b.c.Enabled
}
