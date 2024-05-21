-- +goose Up

-- +goose StatementBegin
DROP TRIGGER IF EXISTS insert_on_terra_msg ON PUBLIC.cosmos_msgs;
DROP FUNCTION IF EXISTS PUBLIC.notify_terra_msg_insert;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
CREATE FUNCTION notify_terra_msg_insert() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    PERFORM pg_notify('insert_on_terra_msg'::text, NOW()::text);
    RETURN NULL;
END
$$;
CREATE TRIGGER notify_terra_msg_insertion AFTER INSERT ON cosmos_msgs FOR EACH STATEMENT EXECUTE PROCEDURE notify_terra_msg_insert();
-- +goose StatementEnd
