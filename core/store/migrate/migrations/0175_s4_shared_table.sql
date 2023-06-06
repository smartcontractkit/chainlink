-- +goose Up

ALTER TABLE "s4".functions RENAME TO shared;

ALTER TABLE "s4".shared ADD COLUMN IF NOT EXISTS namespace TEXT NOT NULL DEFAULT '';

CREATE INDEX shared_namespace_idx ON "s4".shared(namespace);

ALTER INDEX "s4".functions_address_slot_id_idx RENAME TO shared_address_slot_id_idx;
ALTER INDEX "s4".functions_expiration_idx RENAME TO shared_expiration_idx;
ALTER INDEX "s4".functions_confirmed_idx RENAME TO shared_confirmed_idx;

-- +goose Down

ALTER INDEX "s4".shared_address_slot_id_idx RENAME TO functions_address_slot_id_idx;
ALTER INDEX "s4".shared_expiration_idx RENAME TO functions_expiration_idx;
ALTER INDEX "s4".shared_confirmed_idx RENAME TO functions_confirmed_idx;

DROP INDEX IF EXISTS "s4".shared_namespace_idx;

ALTER TABLE "s4".shared DROP COLUMN IF EXISTS namespace;

ALTER TABLE "s4".shared RENAME TO functions;
