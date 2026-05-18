class Broadcaster:
    def __init__(self, adapters):
        self.adapters = adapters

    def publish_manifest(self, g_root, signature):
        # We broadcast the root and signature as the final epoch proof
        for adapter in self.adapters:
            adapter.broadcast(f"{g_root}{signature}")

    def stream_batch_sample(self, sample_payloads):
        # Simulates the mid-run "Sent" bursts from the reference
        for p in sample_payloads:
            for adapter in self.adapters:
                adapter.broadcast(p)
