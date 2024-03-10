-- +goose Up
ALTER TABLE ocr2_oracle_specs
    -- NOTE: The cleanest way to do this would be to allow NULL feed_id and use
    -- postgres 15's NULLS NOT DISTINCT feature on the index.
    -- However, it isn't reasonable to expect all users to upgrade to pg 15 at
    -- this time, so we require all specs to have a feed ID and use the zero
    -- value to indicate a missing feed ID.
    ADD COLUMN feed_id bytea CHECK (octet_length(feed_id) = 32) NOT NULL DEFAULT '\x0000000000000000000000000000000000000000000000000000000000000000', 
    DROP CONSTRAINT offchainreporting2_oracle_specs_unique_contract_addr;
;
CREATE UNIQUE INDEX offchainreporting2_oracle_specs_unique_contract_addr ON ocr2_oracle_specs (contract_id, feed_id);

-- NOTE: bootstrap_specs did not originally have a unique index, so we do not add one here
ALTER TABLE bootstrap_specs ADD COLUMN feed_id bytea CHECK (feed_id IS NULL OR octet_length(feed_id) = 32);

-- +goose Down
ALTER TABLE ocr2_oracle_specs DROP COLUMN feed_id, ADD CONSTRAINT offchainreporting2_oracle_specs_unique_contract_addr UNIQUE (contract_id);
ALTER TABLE bootstrap_specs DROP COLUMN feed_id;
