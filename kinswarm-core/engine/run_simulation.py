import time
import os
from worker_profile import WorkerProfile
from network_adapter import EVMAdapter, CosmosAdapter, SolanaAdapter
from ledger import Ledger
from swarm_manager import SwarmManager

if __name__ == "__main__":
    COUNT = 1000000 
    print(f"🚀 KinSwarm Core: Executing 1M Settlement with Audit Persistence...")
    
    start_time = time.time()
    
    # Generate 1M worker identities
    workers = [WorkerProfile(f"w_{i}", 35, 40, 2) for i in range(COUNT)]
    
    # Initializing global ledger and adapters
    adapters = [EVMAdapter(), CosmosAdapter(), SolanaAdapter()]
    cores = os.cpu_count() or 4
    
    # Start Swarm Manager
    from ledger import Ledger as GlobalLedger
    swarm = SwarmManager(workers, adapters, GlobalLedger(), shard_count=cores)
    swarm.run_epoch(epoch_id=1)
    
    end_time = time.time()
    print(f"\n✅ 1M-Worker Settlement Complete.")
    print(f"📦 Audit Trail saved to shared/audit_shard_*.db")
    print(f"⏱️  Total Execution & Persistence Time: {round(end_time - start_time, 2)}s")
