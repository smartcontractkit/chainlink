# engine/chainlink_adapter.py
import json
from engine.sovereign_metadata import get_sovereign_prefix

class ChainlinkAdapter:
    def __init__(self, oracle_address):
        self.oracle = oracle_address

    def fulfill_request(self, root, total_amount):
        """
        Formats the 10M-worker settlement for a Chainlink consumer contract.
        """
        prefix = get_sovereign_prefix()
        payload = {
            "merkle_root": f"0x{root}",
            "total_settlement": total_amount,
            "attestation": prefix,
            "protocol": "KinSwarm-Chainlink-v1"
        }
        # This JSON is what the 'Functions' consumer receives
        return json.dumps(payload)
