-- +goose Up
ALTER TABLE mercury_transmit_requests DROP CONSTRAINT mercury_transmit_requests_pkey;
DELETE FROM mercury_transmit_requests;
ALTER TABLE mercury_transmit_requests ADD COLUMN server_url TEXT NOT NULL;
ALTER TABLE mercury_transmit_requests ADD PRIMARY KEY (server_url, payload_hash);

-- +goose Down
ALTER TABLE mercury_transmit_requests DROP COLUMN server_url;
ALTER TABLE mercury_transmit_requests ADD PRIMARY KEY (payload_hash);
