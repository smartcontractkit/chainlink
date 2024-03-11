-- +goose Up

CREATE TABLE job_kv_store (
      id SERIAL PRIMARY KEY,
      key varchar,
      val JSONB,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      FOREIGN KEY (id) REFERENCES jobs(id),
      CONSTRAINT unique_id_key UNIQUE (id, key)
);

-- +goose Down
Drop table job_kv_store;
