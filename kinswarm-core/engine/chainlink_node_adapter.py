import eth_abi
import hashlib
from engine.sovereign_metadata import get_sovereign_prefix

class KinSwarmNode:
    def __init__(self, expected_root):
        self.expected_root = expected_root
        self.identity = get_sovereign_prefix()

    def prepare_onchain_payload(self, total_amount):
        """
        ABI-encodes the audit results for the Solidity fulfillSettlement function.
        Format: (bytes32 root, uint256 amount, string identity)
        """
        # Convert hex string to bytes32
        root_bytes = bytes.fromhex(self.expected_root.replace("0x", ""))
        
        # ABI Encode for EVM
        encoded_payload = eth_abi.encode(
            ['bytes32', 'uint256', 'string'],
            [root_bytes, total_amount, self.identity]
        )
        
        print(f"\n--- Chainlink Node Payload Generated ---")
        print(f"Encoded Hex: 0x{encoded_payload.hex()[:64]}...")
        print(f"Attestation: {self.identity}")
        return encoded_payload

if __name__ == "__main__":
    # Integration test with your verified root
    ROOT = "6a99f4d0755e0ce9dba8afb3f3bde5c0c23a364ad47e886ebbaeca8ba75914b2"
    node = KinSwarmNode(ROOT)
    node.prepare_onchain_payload(350000000)
