-- +goose Up

-- evm.heads.id is redundant since a head is uniquely identified by (evm_chain_id, hash).
-- Using an id column to resolve (evm_chain_id, number) tie breaks on initial load is an overkill.
--
-- The migration soft drops the id column by renaming it to deprecated_id and dropping the primary key constraint that includes it.
-- This allows us to double-check that nothing is using the id column before we hard drop it in a future migration.
ALTER TABLE evm.heads DROP CONSTRAINT heads_pkey1;
ALTER TABLE evm.heads RENAME COLUMN id TO deprecated_id;

-- +goose Down

-- Revert the column rename and restore the primary key constraint
ALTER TABLE evm.heads RENAME COLUMN deprecated_id TO id;
ALTER TABLE evm.heads ADD CONSTRAINT heads_pkey1 PRIMARY KEY (id);

