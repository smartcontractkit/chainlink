-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS notify_insert_on_logs_topics ON EVM.logs;
DROP FUNCTION IF EXISTS evm.notifysavedlogtopics();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

CREATE FUNCTION evm.notifysavedlogtopics() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    PERFORM pg_notify(
        'evm.insert_on_logs'::text,
        -- hex encoded address plus comma separated list of hex encoded topic values
        -- e.g. "<address>:<topicVal1>,<topicVal2>"
        encode(NEW.address, 'hex') || ':' || array_to_string(array(SELECT encode(unnest(NEW.topics), 'hex')), ',')
    );
    RETURN NULL;
END
$$;

CREATE TRIGGER notify_insert_on_logs_topics AFTER INSERT ON evm.logs FOR EACH ROW EXECUTE PROCEDURE evm.notifysavedlogtopics();
-- +goose StatementEnd
