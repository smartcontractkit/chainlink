import sqlite3
import hashlib
class Ledger:
    def __init__(self, shard_id=0):
        self.shard_id = shard_id
        self.records = []
        self.db_path = f"shared/audit_shard_{shard_id}.db"
        self._conn = None
    def _get_conn(self):
        if self._conn is None:
            self._conn = sqlite3.connect(self.db_path)
            self._conn.execute("PRAGMA journal_mode=WAL")
            self._conn.execute("PRAGMA synchronous=OFF")
            self._conn.execute("CREATE TABLE IF NOT EXISTS settlements (worker_id TEXT, amount INTEGER, payload TEXT)")
        return self._conn
    def record(self, worker_id, amount, payload):
        self.records.append((worker_id, amount, str(payload)))
    def commit_to_disk(self):
        conn = self._get_conn()
        conn.executemany("INSERT INTO settlements VALUES (?, ?, ?)", self.records)
        conn.commit()
        conn.close()
        self._conn = None
    def get_shard_root(self):
        combined = "".join([str(r[2]) for r in self.records]).encode()
        return hashlib.sha256(combined).hexdigest()
