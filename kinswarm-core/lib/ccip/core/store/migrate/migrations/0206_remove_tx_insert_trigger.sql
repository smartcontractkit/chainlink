-- +goose Up
DROP TRIGGER IF EXISTS notify_tx_insertion on evm.txes; 
DROP FUNCTION IF EXISTS evm.notifyethtxinsertion();


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

CREATE TRIGGER notify_tx_insertion AFTER INSERT ON evm.txes FOR EACH ROW EXECUTE PROCEDURE evm.notifytxinsertion();
-- +goose StatementEnd