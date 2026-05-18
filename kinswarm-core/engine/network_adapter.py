import asyncio

class NetworkAdapter:
    async def broadcast_batch(self, batch_payloads):
        # Simulates non-blocking network I/O to multichain endpoints
        return True

class EVMAdapter(NetworkAdapter): pass
class CosmosAdapter(NetworkAdapter): pass
class SolanaAdapter(NetworkAdapter): pass

class Broadcaster:
    def __init__(self):
        self.adapters = [EVMAdapter(), CosmosAdapter(), SolanaAdapter()]
    
    async def push(self, payloads):
        tasks = [a.broadcast_batch(payloads) for a in self.adapters]
        await asyncio.gather(*tasks)
