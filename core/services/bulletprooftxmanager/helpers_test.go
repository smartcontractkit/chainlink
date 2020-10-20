package bulletprooftxmanager

func ExportedTriggerChan(eb EthBroadcaster) <-chan struct{} {
	return eb.(*ethBroadcaster).trigger
}

func ExportedMustStartEthTxInsertListener(eb EthBroadcaster) {
	go eb.(*ethBroadcaster).ethTxInsertTriggerer()
}
