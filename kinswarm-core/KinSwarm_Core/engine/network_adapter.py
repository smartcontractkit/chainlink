import hashlib

class NetworkAdapter:
    def __init__(self, name):
        self.name = name

    def broadcast(self, payload):
        # Generate a raw byte signature of the payload
        binary_payload = hashlib.sha256(str(payload).encode()).digest()
        # Truncate for terminal display like the reference
        display_bytes = binary_payload[:12]
        print(f"[{self.name}] Sent: b'{str(display_bytes)[3:-1]}...'")

class EVMAdapter(NetworkAdapter):
    def __init__(self): super().__init__("EVMAdapter")

class CosmosAdapter(NetworkAdapter):
    def __init__(self): super().__init__("CosmosAdapter")

class SolanaAdapter(NetworkAdapter):
    def __init__(self): super().__init__("SolanaAdapter")
