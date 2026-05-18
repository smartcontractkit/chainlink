-- +goose Up

CREATE SCHEMA "s4";

CREATE TABLE "s4".functions(
    id BIGSERIAL PRIMARY KEY,
    address NUMERIC(78,0) NOT NULL,
    slot_id INT NOT NULL,
    version INT NOT NULL,
    expiration BIGINT NOT NULL,
    confirmed BOOLEAN NOT NULL,
    payload BYTEA NOT NULL,
    signature BYTEA NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX functions_address_slot_id_idx ON "s4".functions(address, slot_id);
CREATE INDEX functions_expiration_idx ON "s4".functions(expiration);
CREATE INDEX functions_confirmed_idx ON "s4".functions(confirmed);

-- +goose Down

DROP INDEX IF EXISTS functions_address_slot_id_idx;
DROP INDEX IF EXISTS functions_expiration_idx;
DROP INDEX IF EXISTS functions_confirmed_idx;

DROP TABLE "s4".functions;

DROP SCHEMA "s4";