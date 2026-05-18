import multiprocessing
import hashlib
import importlib

def process_shard(shard_args):
    l_mod = importlib.import_module("engine.ledger")
    o_mod = importlib.import_module("engine.offchain_worker")
    workers, adapters, shard_id, epoch_id = shard_args
    shard_ledger = l_mod.Ledger(shard_id=shard_id)
    ocw = o_mod.OffchainWorker(workers, adapters, shard_ledger)
    ocw.settle_epoch(epoch_id)
    root = shard_ledger.get_shard_root()
    shard_ledger.commit_to_disk()
    return root

class SwarmManager:
    def __init__(self, workers, adapters, shard_count=12):
        self.workers = workers
        self.adapters = adapters
        self.shard_count = shard_count
    def run_epoch(self, epoch_id):
        print(f"Batching {len(self.workers)} workers across {self.shard_count} shards...")
        shard_tasks = []
        for i in range(self.shard_count):
            s_idx = i * len(self.workers) // self.shard_count
            e_idx = (i + 1) * len(self.workers) // self.shard_count
            shard_tasks.append((self.workers[s_idx:e_idx], self.adapters, i, epoch_id))
        with multiprocessing.Pool(self.shard_count) as pool:
            roots = pool.map(process_shard, shard_tasks)
        g_root = hashlib.sha256("".join(roots).encode()).hexdigest()
        return g_root
