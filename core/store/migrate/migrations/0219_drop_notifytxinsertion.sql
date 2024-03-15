-- +goose Up
-- +goose StatementBegin
DROP FUNCTION IF EXISTS evm.notifytxinsertion();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION evm.notifytxinsertion() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
        BEGIN
		PERFORM pg_notify('evm.insert_on_txes'::text, encode(NEW.from_address, 'hex'));
		RETURN NULL;
        END
        $$;
-- +goose StatementEnd
