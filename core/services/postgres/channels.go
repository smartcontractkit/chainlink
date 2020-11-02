package postgres

const (
	ChannelJobCreated   = "insert_on_jobs"
	ChannelRunStarted   = "pipeline_run_started"
	ChannelRunCompleted = "pipeline_run_completed"

	// Postgres channel to listen for new eth_txes
	ChannelInsertOnEthTx = "insert_on_eth_txes"
)
