-- +goose Up

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION PUBLIC.notifysavedlogtopics() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    PERFORM pg_notify(
        'insert_on_evm_logs'::text,
        -- hex encoded address plus comma separated list of hex encoded topic values
        -- e.g. "<address>:<topicVal1>,<topicVal2>"
        encode(NEW.address, 'hex') || ':' || array_to_string(array(SELECT encode(unnest(NEW.topics), 'hex')), ',')
    );
    RETURN NULL;
END
$$;

DROP TRIGGER IF EXISTS notify_insert_on_evm_logs_topics ON PUBLIC.evm_logs;
CREATE TRIGGER notify_insert_on_evm_logs_topics AFTER INSERT ON PUBLIC.evm_logs FOR EACH ROW EXECUTE PROCEDURE PUBLIC.notifysavedlogtopics();
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TRIGGER IF EXISTS notify_insert_on_evm_logs_topics ON PUBLIC.evm_logs;
DROP FUNCTION IF EXISTS PUBLIC.notifysavedlogtopics;
-- +goose StatementEnd
