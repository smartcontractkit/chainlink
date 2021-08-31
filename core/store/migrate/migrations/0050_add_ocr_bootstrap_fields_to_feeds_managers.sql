-- +goose Up
ALTER TABLE feeds_managers
DROP COLUMN network,
ADD COLUMN ocr_bootstrap_peer_multiaddr VARCHAR,
ADD CONSTRAINT chk_ocr_bootstrap_peer_multiaddr CHECK ( NOT (
	is_ocr_bootstrap_peer AND
	(
		ocr_bootstrap_peer_multiaddr IS NULL OR
		ocr_bootstrap_peer_multiaddr = ''
	)
));

-- +goose Down
ALTER TABLE feeds_managers
ADD COLUMN network VARCHAR (100),
DROP CONSTRAINT chk_ocr_bootstrap_peer_multiaddr,
DROP COLUMN ocr_bootstrap_peer_multiaddr;
