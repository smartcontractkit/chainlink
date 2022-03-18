-- +goose Up
CREATE UNIQUE INDEX idx_unique_ws_url ON evm_nodes (ws_url);
CREATE UNIQUE INDEX idx_unique_http_url ON evm_nodes (http_url);

-- +goose Down
DROP INDEX idx_unique_ws_url;
DROP INDEX idx_unique_http_url;
