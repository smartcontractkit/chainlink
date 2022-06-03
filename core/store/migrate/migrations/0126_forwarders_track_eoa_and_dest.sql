-- +goose Up
ALTER TABLE evm_forwarders
ADD COLUMN eoa bytea NOT NULL,
ADD COLUMN dest bytea NOT NULL,
DROP CONSTRAINT evm_forwarders_address_key,
ADD CONSTRAINT chk_dest_address_length CHECK ((octet_length(dest) = 20)),
ADD CONSTRAINT chk_eoa_address_length CHECK ((octet_length(eoa) = 20)),
ADD CONSTRAINT evm_forwarders_addr_eoa_dest_key UNIQUE (address, eoa, dest);

-- +goose Down
ALTER TABLE evm_forwarders
DROP CONSTRAINT chk_dest_address_length,
DROP CONSTRAINT chk_eoa_address_length,
DROP CONSTRAINT evm_forwarders_addr_eoa_dest_key,
DROP COLUMN eoa,
DROP COLUMN dest,
ADD CONSTRAINT evm_forwarders_address_key UNIQUE (address);
