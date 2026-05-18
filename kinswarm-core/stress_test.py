import time
import os
from engine.worker_profile import WorkerProfile
from engine.network_adapter import EVMAdapter, CosmosAdapter, SolanaAdapter
from engine.swarm_manager import SwarmManager
from engine.ccip_manager import CCIPManager
from engine.sovereign_metadata import get_sovereign_prefix

def run_chainlink_endgame():
    prefix = get_sovereign_prefix()
    print(f"{prefix}")
    print(f"--- Chainlink Pivot: 10,000,000 Worker Epoch ---")
    
    # 1. Clean Slate
    if os.path.exists("shared"):
        import shutil
        shutil.rmtree("shared")
    os.makedirs("shared")
    
    start_time = time.time()
    
    # 2. Worker Generation (10M at $35/hr)
    workers = [WorkerProfile(f"w_{i}", 35) for i in range(10000000)]
    
    # 3. Off-Chain Aggregation (OCR Layer)
    adapters = [EVMAdapter(), CosmosAdapter(), SolanaAdapter()]
    swarm = SwarmManager(workers, adapters, shard_count=12)
    g_root = swarm.run_epoch(1)
    
    total_volume = 10000000 * 35 # $350M
    
    # 4. Chainlink CCIP Bridge
    ccip = CCIPManager()
    ccip.broadcast_to_asset_hub(g_root, total_volume)
    
    elapsed = round(time.time() - start_time, 2)
    print(f"\n--- Protocol State: FINALIZED ---")
    print(f"Time: {elapsed}s | Throughput: {int(10000000/elapsed)} tx/s")
    print(f"Merkle Anchor: 0x{g_root}")

if __name__ == "__main__":
    run_chainlink_endgame()
