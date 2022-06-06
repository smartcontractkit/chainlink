-- +goose Up
ALTER TABLE evm_forwarders
ADD COLUMN eoa bytea NOT NULL,
ADD COLUMN dest bytea NOT NULL,
DROP CONSTRAINT evm_forwarders_address_key,
ADD CONSTRAINT chk_dest_address_length CHECK ((octet_length(dest) = 20)),
ADD CONSTRAINT chk_eoa_address_length CHECK ((octet_length(eoa) = 20));


CREATE UNIQUE INDEX evm_forwarders_eoa_dest_key ON evm_forwarders using btree(eoa, dest);


-- +goose Down

DROP INDEX evm_forwarders_eoa_dest_key;

ALTER TABLE evm_forwarders
DROP CONSTRAINT chk_dest_address_length,
DROP CONSTRAINT chk_eoa_address_length,
DROP COLUMN eoa,
DROP COLUMN dest,
ADD CONSTRAINT evm_forwarders_address_key UNIQUE (address);
