import sqlite3
import os

def inject_corruption():
    db_path = "./shared/audit_shard_5.db"
    
    if not os.path.exists(db_path):
        print(f"Error: Shard not found.")
        return

    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    table_name = "settlements"
    target_col = "amount"
    id_col = "worker_id"
    
    print(f"--- Poisoning {table_name} | Target: {target_col} | ID: {id_col} ---")
    
    # Injecting the 1.337 discrepancy
    cursor.execute(f"UPDATE {table_name} SET {target_col} = {target_col} + 1.337 WHERE {id_col} = (SELECT {id_col} FROM {table_name} LIMIT 1)")
    
    conn.commit()
    conn.close()
    print("SUCCESS: Cryptographic seal broken. Shard 5 is now adversarial.")

if __name__ == "__main__":
    inject_corruption()
