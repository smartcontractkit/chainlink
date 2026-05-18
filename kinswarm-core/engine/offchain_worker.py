class OffchainWorker:
    def __init__(self, workers, adapters, ledger):
        self.workers = workers
        self.adapters = adapters
        self.ledger = ledger

    def settle_epoch(self, epoch_id):
        for worker in self.workers:
            net_amt = worker.net_amount()
            # Logic anchoring for 10M scale
            pto_r = getattr(worker, 'pto_allocated', 48) - getattr(worker, 'pto_used', 2)
            payload = f"{worker.worker_id}-{net_amt}-{epoch_id}"
            # Signature must match Ledger.record(worker_id, amount, pto_r, payload)
            self.ledger.record(worker.worker_id, net_amt, pto_r, payload)
