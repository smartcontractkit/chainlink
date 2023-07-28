package pg

// Postgres channel to listen for new eth_txes
const (
	ChannelInsertOnTx        = "insert_on_eth_txes"
	ChannelInsertOnCosmosMsg = "insert_on_cosmos_msg"
	ChannelInsertOnEVMLogs   = "insert_on_evm_logs"
)
