-- +goose Up
-- Schemas for the CRE processes the node launches beside itself: the p2p proxy (crecore) and the EVM
-- capability. Each keeps its own tables, and a schema per process is what keeps them from landing in
-- public beside the node's - or, in the EVM capability's case, on top of the node's own evm schema,
-- where its log poller and transaction manager would run over the node's tables.
--
-- They are created here rather than by the processes themselves because the database is the node's:
-- this is the migration set that owns it, and a process launched by the node should not need the
-- right to create a schema in it.
CREATE SCHEMA IF NOT EXISTS crecore;
CREATE SCHEMA IF NOT EXISTS evm_capability;

-- +goose Down
-- Dropped with what is in them: the tables belong to those processes, and a schema kept without them
-- is not the node's to hold either.
DROP SCHEMA IF EXISTS crecore CASCADE;
DROP SCHEMA IF EXISTS evm_capability CASCADE;
