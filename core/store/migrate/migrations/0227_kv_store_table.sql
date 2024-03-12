-- +goose Up

CREATE TABLE job_kv_store (
      id SERIAL PRIMARY KEY,
      key VARCHAR,
      val JSONB,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      FOREIGN KEY (id) REFERENCES jobs(id),
      CONSTRAINT uk_keyperjob UNIQUE (id, key)
);

-- +goose Down
DROP TABLE job_kv_store;
