import sqlite3
import hashlib
import os

def calculate_shard_hash(db_path):
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    # Force a hash of the actual data rows: worker_id + amount
    cursor.execute("SELECT worker_id, amount FROM settlements ORDER BY worker_id")
    
    sha = hashlib.sha256()
    for row in cursor.fetchall():
        sha.update(f"{row[0]}{row[1]}".encode())
    
    conn.close()
    return sha.hexdigest()

def run_audit():
    expected_root = "6a99f4d0755e0ce9dba8afb3f3bde5c0c23a364ad47e886ebbaeca8ba75914b2"
    shard_hashes = []
    
    print("--- Initiating RAW DATA Forensic Audit ---")
    for i in range(12):
        db_path = f"./shared/audit_shard_{i}.db"
        h = calculate_shard_hash(db_path)
        shard_hashes.append(h)
        print(f"  [Shard {i}] Hash: {h[:16]}...")

    # Final Merkle Root of the Shard Hashes
    final_sha = hashlib.sha256()
    for sh in sorted(shard_hashes):
        final_sha.update(sh.encode())
    
    reconstructed = final_sha.hexdigest()
    
    print(f"\nAudit Summary:")
    print(f"Reconstructed: {reconstructed}")
    print(f"Expected:      {expected_root}")
    
    if reconstructed == expected_root:
        print("✅ INTEGRITY VERIFIED")
    else:
        print("❌ ALERT: HASH MISMATCH DETECTED. DATA TAMPERED.")

if __name__ == "__main__":
    run_audit()
