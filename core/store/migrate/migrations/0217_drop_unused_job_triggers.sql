-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS notify_job_created ON PUBLIC.jobs;
DROP FUNCTION IF EXISTS PUBLIC.notifyjobcreated();

DROP TRIGGER IF EXISTS notify_job_deleted ON PUBLIC.jobs;
DROP FUNCTION IF EXISTS PUBLIC.notifyjobdeleted();

DROP TRIGGER IF EXISTS notify_pipeline_run_started ON PUBLIC.pipeline_runs;
DROP FUNCTION IF EXISTS PUBLIC.notifypipelinerunstarted();
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
CREATE FUNCTION PUBLIC.notifyjobcreated() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
        BEGIN
            PERFORM pg_notify('insert_on_jobs', NEW.id::text);
            RETURN NEW;
        END
        $$;
CREATE TRIGGER notify_job_created AFTER INSERT ON PUBLIC.jobs FOR EACH ROW EXECUTE PROCEDURE PUBLIC.notifyjobcreated();

CREATE FUNCTION PUBLIC.notifyjobdeleted() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
	BEGIN
		PERFORM pg_notify('delete_from_jobs', OLD.id::text);
		RETURN OLD;
	END
	$$;
CREATE TRIGGER notify_job_deleted AFTER DELETE ON PUBLIC.jobs FOR EACH ROW EXECUTE PROCEDURE PUBLIC.notifyjobdeleted();

CREATE FUNCTION PUBLIC.notifypipelinerunstarted() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
	BEGIN
		IF NEW.finished_at IS NULL THEN
			PERFORM pg_notify('pipeline_run_started', NEW.id::text);
		END IF;
		RETURN NEW;
	END
	$$;
CREATE TRIGGER notify_pipeline_run_started AFTER INSERT ON PUBLIC.pipeline_runs FOR EACH ROW EXECUTE PROCEDURE PUBLIC.notifypipelinerunstarted();

-- +goose StatementEnd
