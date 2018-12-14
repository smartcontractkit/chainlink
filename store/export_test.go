package store

func ExportedSetTxManagerDev(txm TxManager, dev bool) {
	typed := txm.(*EthTxManager)
	typed.config.Dev = dev
}
