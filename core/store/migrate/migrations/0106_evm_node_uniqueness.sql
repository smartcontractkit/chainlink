-- +goose Up
-- Delete sendonlys if they redundantly duplicate a primary
DELETE FROM
    evm_nodes a
    USING evm_nodes b
WHERE
    a.http_url = b.http_url
    AND a.id != b.id
    AND a.send_only;

CREATE UNIQUE INDEX idx_unique_ws_url ON evm_nodes (ws_url);
CREATE UNIQUE INDEX idx_unique_http_url ON evm_nodes (http_url);

-- +goose Down
DROP INDEX idx_unique_ws_url;
DROP INDEX idx_unique_http_url;
