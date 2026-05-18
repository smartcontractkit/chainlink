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
    
    # Extract sample for terminal output
    sample = [r[2] for r in shard_ledger.records[:2]] 
    
    root = shard_ledger.get_shard_root()
    shard_ledger.commit_to_disk()
    return root, sample

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
            results = pool.map(process_shard, shard_tasks)
        
        roots = [r[0] for r in results]
        samples = [r[1] for r in results]
        
        # Stream the binary output to terminal (flattened)
        from engine.broadcaster import Broadcaster
        b = Broadcaster(self.adapters)
        for s_list in samples:
            b.stream_batch_sample(s_list)

        g_root = hashlib.sha256("".join(roots).encode()).hexdigest()
        return g_root
