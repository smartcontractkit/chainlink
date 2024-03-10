-- +goose Up

CREATE TABLE bridge_last_value (
    spec_id int NOT NULL REFERENCES public.pipeline_specs(id) ON DELETE CASCADE DEFERRABLE,
    dot_id text NOT NULL,
    value bytea NOT NULL,
    finished_at timestamp NOT NULL,
    CONSTRAINT bridge_last_value_pkey PRIMARY KEY (spec_id, dot_id)
);

CREATE INDEX idx_bridge_last_value_optimise_finding_last_value ON bridge_last_value USING btree (finished_at);


-- +goose Down
DROP INDEX idx_bridge_last_value_optimise_finding_last_value;
DROP TABLE bridge_last_value;
