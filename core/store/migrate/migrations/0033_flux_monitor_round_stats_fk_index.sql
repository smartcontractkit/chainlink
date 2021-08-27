-- +goose Up
CREATE INDEX flux_monitor_round_stats_job_run_id_idx ON flux_monitor_round_stats (job_run_id);
CREATE INDEX flux_monitor_round_stats_v2_pipeline_run_id_idx ON flux_monitor_round_stats_v2 (pipeline_run_id);

-- + goose Down
DROP INDEX flux_monitor_round_stats_job_run_id_idx;
DROP INDEX flux_monitor_round_stats_v2_pipeline_run_id_idx;
