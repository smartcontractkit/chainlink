package pg

// Postgres channel to listen for new evm.txes
const (
	ChannelInsertOnTx        = "evm.insert_on_txes"
	ChannelInsertOnCosmosMsg = "insert_on_cosmos_msg"
	ChannelInsertOnEVMLogs   = "insert_on_evm_logs"
)
