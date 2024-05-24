-- +goose Up
UPDATE offchainreporting_oracle_specs SET contract_config_confirmations = 0 where contract_config_confirmations is NULL;
ALTER TABLE offchainreporting_oracle_specs
	ALTER COLUMN contract_config_confirmations SET NOT NULL;
ALTER TABLE external_initiators ADD CONSTRAINT "access_key_unique" UNIQUE ("access_key");

-- +goose Down
ALTER TABLE offchainreporting_oracle_specs 
	ALTER COLUMN contract_config_confirmations DROP NOT NULL;
ALTER TABLE external_initiators DROP CONSTRAINT "access_key_unique";
