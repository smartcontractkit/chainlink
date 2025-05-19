-- +goose Up

CREATE TABLE job_kv_store (
      job_id INTEGER NOT NULL REFERENCES jobs(id),
      key TEXT NOT NULL,
      val JSONB NOT NULL,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      PRIMARY KEY (job_id, key)
);

-- +goose Down
DROP TABLE job_kv_store;
