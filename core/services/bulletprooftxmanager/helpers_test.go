package bulletprooftxmanager

func ExportedTriggerChan(eb EthBroadcaster) <-chan struct{} {
	return eb.(*ethBroadcaster).trigger
}

func ExportedMustStartEthTxInsertListener(eb EthBroadcaster) {
	if err := eb.(*ethBroadcaster).ethTxInsertListener.Start(); err != nil {
		panic(err)
	}
	go eb.(*ethBroadcaster).ethTxInsertTriggerer()
}
