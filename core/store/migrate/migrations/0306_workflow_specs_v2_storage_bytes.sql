-- +goose Up
-- Durable storage footprint of each registered workflow spec, in bytes
-- (raw decoded binary + config), for storage-based metering. Set at
-- registration and intentionally NOT cleared by pause tombstoning: bytes are
-- billed for the registration lifetime (register -> delete), so the value must
-- survive payload clearing to keep deltas self-balancing and snapshots
-- consistent. Pre-migration rows are backfilled from the persisted payload
-- where it still exists; already-paused tombstones backfill to 0 (the original
-- bytes are unrecoverable).
ALTER TABLE workflow_specs_v2 ADD COLUMN storage_bytes bigint NOT NULL DEFAULT 0;
UPDATE workflow_specs_v2 SET storage_bytes = OCTET_LENGTH(workflow)/2 + OCTET_LENGTH(config) WHERE workflow <> '';

-- +goose Down
ALTER TABLE workflow_specs_v2 DROP COLUMN storage_bytes;
