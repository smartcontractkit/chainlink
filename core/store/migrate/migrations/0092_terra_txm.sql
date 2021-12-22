-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION notify_terra_msg_insert() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    PERFORM pg_notify('insert_on_terra_msg'::text, NOW()::text);
    RETURN NULL;
END
$$;
CREATE TABLE terra_msgs (
	id BIGSERIAL PRIMARY KEY,
	contract_id text NOT NULL,
    msg bytea NOT NULL,
	state text NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL
);
CREATE TRIGGER notify_terra_msg_insertion AFTER INSERT ON terra_msgs FOR EACH STATEMENT EXECUTE PROCEDURE notify_terra_msg_insert();

-- +goose StatementEnd

-- +goose Down
DROP TABLE terra_msgs;
DROP FUNCTION notify_terra_msg_insert;
-- +goose StatementBegin
-- +goose StatementEnd
