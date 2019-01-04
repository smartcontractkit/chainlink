package store

func ExportedSetTxManagerDev(txm TxManager, dev bool) {
	typed := txm.(*EthTxManager)
	typed.config.Set("CHAINLINK_DEV", dev)
}
