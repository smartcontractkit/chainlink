-- +goose Up

ALTER TABLE "s4".functions RENAME TO shared;

ALTER TABLE "s4".shared ADD COLUMN IF NOT EXISTS namespace TEXT NOT NULL DEFAULT '';

DROP INDEX IF EXISTS "s4".functions_address_slot_id_idx;
DROP INDEX IF EXISTS "s4".functions_expiration_idx;
DROP INDEX IF EXISTS "s4".functions_confirmed_idx;

CREATE UNIQUE INDEX shared_namespace_address_slot_id_idx ON "s4".shared(namespace, address, slot_id);
CREATE INDEX shared_namespace_expiration_idx ON "s4".shared(namespace, expiration);
CREATE INDEX shared_namespace_confirmed_idx ON "s4".shared(namespace, confirmed);

-- +goose Down

DROP INDEX IF EXISTS "s4".shared_namespace_address_slot_id_idx;
DROP INDEX IF EXISTS "s4".shared_namespace_expiration_idx;
DROP INDEX IF EXISTS "s4".shared_namespace_confirmed_idx;

ALTER TABLE "s4".shared DROP COLUMN IF EXISTS namespace;

ALTER TABLE "s4".shared RENAME TO functions;

CREATE UNIQUE INDEX functions_address_slot_id_idx ON "s4".functions(address, slot_id);
CREATE INDEX functions_expiration_idx ON "s4".functions(expiration);
CREATE INDEX functions_confirmed_idx ON "s4".functions(confirmed);

