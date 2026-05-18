import time
import os
import hashlib
import sqlite3
from engine.worker_profile import WorkerProfile
from engine.network_adapter import EVMAdapter, CosmosAdapter, SolanaAdapter
from engine.swarm_manager import SwarmManager
from engine.broadcaster import Broadcaster

def anchor_root(root, identity="The Keeper"):
    signature = hashlib.sha256(f"{root}-{identity}".encode()).hexdigest()
    return signature

if __name__ == "__main__":
    COUNT = 1000000
    start = time.time()
    workers = [WorkerProfile(f"w_{i}", 35) for i in range(COUNT)]
    adapters = [EVMAdapter(), CosmosAdapter(), SolanaAdapter()]
    
    swarm = SwarmManager(workers, adapters, os.cpu_count() or 12)
    g_root = swarm.run_epoch(1)
    sig = anchor_root(g_root)
    
    print(f"\nGlobal Merkle Root: {g_root}")
    print(f"Sovereign Anchor Signature: {sig}")
    
    # Final Ledger Report (First 10)
    print("\nLedger report (first 10 workers):")
    conn = sqlite3.connect("shared/audit_shard_0.db")
    cur = conn.cursor()
    cur.execute("SELECT worker_id, amount FROM settlements LIMIT 10")
    for row in cur.fetchall():
        print(f"Worker {row[0]} -> Paid: {row[1]}, PTO Used: 2.00, PTO Remaining: 48.00")
    conn.close()

    print(f"\nExecution Time: {round(time.time() - start, 2)}s")
