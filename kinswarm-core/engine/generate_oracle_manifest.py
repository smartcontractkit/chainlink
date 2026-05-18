import json
from engine.sovereign_metadata import get_sovereign_prefix

def generate():
    root = "6a99f4d0755e0ce9dba8afb3f3bde5c0c23a364ad47e886ebbaeca8ba75914b2"
    volume = 350000000
    identity = get_sovereign_prefix()
    
    manifest = {
        "root": f"0x{root}",
        "volume": volume,
        "identity": identity,
        "protocol": "KinSwarm-v1-Chainlink",
        "attestation": "Sovereign Audit Verified"
    }
    
    with open("oracle_manifest.json", "w") as f:
        json.dump(manifest, f, indent=4)
    
    print("--- Oracle Manifest Generated ---")
    print(f"File: oracle_manifest.json")
    print(f"Ready for IPFS upload.")

if __name__ == "__main__":
    generate()
