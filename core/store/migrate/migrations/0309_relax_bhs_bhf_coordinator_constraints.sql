-- +goose Up
-- Relax the at_least_one_coordinator constraints so that a VRF V2Plus
-- coordinator address alone satisfies them. The constraints previously only
-- accepted V1/V2 coordinator addresses, which blocked persisting V2Plus-only
-- blockhash store / block header feeder jobs (the only supported kind since
-- the V1/V2 coordinator support was removed).
ALTER TABLE blockhash_store_specs DROP CONSTRAINT at_least_one_coordinator_chk;
ALTER TABLE blockhash_store_specs
    ADD CONSTRAINT at_least_one_coordinator_chk CHECK (
            coordinator_v1_address IS NOT NULL
            OR coordinator_v2_address IS NOT NULL
            OR coordinator_v2_plus_address IS NOT NULL
        );

ALTER TABLE block_header_feeder_specs DROP CONSTRAINT at_least_one_coordinator_chk;
ALTER TABLE block_header_feeder_specs
    ADD CONSTRAINT at_least_one_coordinator_chk CHECK (
            coordinator_v1_address IS NOT NULL
            OR coordinator_v2_address IS NOT NULL
            OR coordinator_v2_plus_address IS NOT NULL
        );

-- +goose Down
ALTER TABLE block_header_feeder_specs DROP CONSTRAINT at_least_one_coordinator_chk;
ALTER TABLE block_header_feeder_specs
    ADD CONSTRAINT at_least_one_coordinator_chk CHECK (
            coordinator_v1_address IS NOT NULL OR coordinator_v2_address IS NOT NULL
        );

ALTER TABLE blockhash_store_specs DROP CONSTRAINT at_least_one_coordinator_chk;
ALTER TABLE blockhash_store_specs
    ADD CONSTRAINT at_least_one_coordinator_chk CHECK (
            coordinator_v1_address IS NOT NULL OR coordinator_v2_address IS NOT NULL
        );
