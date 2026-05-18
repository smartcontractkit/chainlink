-- +goose Up
CREATE TABLE flux_monitor_round_stats_v2 (
	id BIGSERIAL PRIMARY KEY,
	aggregator bytea NOT NULL,
	round_id integer NOT NULL,
	num_new_round_logs integer NOT NULL DEFAULT 0,
	num_submissions integer NOT NULL DEFAULT 0,
	pipeline_run_id bigint REFERENCES pipeline_runs(id) ON DELETE CASCADE,
	CONSTRAINT flux_monitor_round_stats_v2_aggregator_round_id_key UNIQUE (aggregator, round_id)
);

-- +goose Down
DROP TABLE flux_monitor_round_stats_v2;
