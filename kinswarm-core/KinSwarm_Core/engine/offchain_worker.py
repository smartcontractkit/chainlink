class OffchainWorker:
    def __init__(self, workers, adapters, ledger):
        self.workers = workers
        self.adapters = adapters
        self.ledger = ledger
    def settle_epoch(self, epoch_id):
        for worker in self.workers:
            net_amt = worker.net_amount()
            payload = f"{worker.worker_id}-{net_amt}-{epoch_id}"
            self.ledger.record(worker.worker_id, net_amt, payload)
