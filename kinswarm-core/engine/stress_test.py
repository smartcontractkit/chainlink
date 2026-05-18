import time
import os
import hashlib
import sqlite3
import multiprocessing
from engine.worker_profile import WorkerProfile
from engine.network_adapter import EVMAdapter, CosmosAdapter, SolanaAdapter
from engine.swarm_manager import SwarmManager

def run_stress_test(total_workers=10000000):
    start_time = time.time()
    print(f"--- Starting 10M Worker Stress Test ---")
    
    # Initialize Adapters
    adapters = [EVMAdapter(), CosmosAdapter(), SolanaAdapter()]
    
    # Generate Workers in chunks to save RAM
    print("Generating 10M Worker Profiles...")
    workers = [WorkerProfile(f"w_{i}", 35) for i in range(total_workers)]
    
    # Initialize Swarm with max physical cores
    shards = os.cpu_count() or 12
    swarm = SwarmManager(workers, adapters, shard_count=shards)
    
    # Execute Epoch
    print(f"Executing Settlement across {shards} shards...")
    g_root = swarm.run_epoch(1)
    
    end_time = time.time()
    elapsed = round(end_time - start_time, 2)
    
    print(f"\n--- Stress Test Results ---")
    print(f"Total Workers: {total_workers}")
    print(f"Global Merkle Root: {g_root}")
    print(f"Total Execution Time: {elapsed}s")
    print(f"Throughput: {int(total_workers / elapsed)} workers/sec")

if __name__ == "__main__":
    run_stress_test()
