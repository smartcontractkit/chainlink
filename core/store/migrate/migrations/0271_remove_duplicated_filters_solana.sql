-- +goose Up
-- +goose StatementBegin
ALTER TABLE solana.log_poller_filters ADD COLUMN IF NOT EXISTS is_legacy_name BOOLEAN;
UPDATE solana.log_poller_filters SET is_legacy_name = true;
with uniques as (
	select id from (
		select id,
		ROW_NUMBER() OVER (
			PARTITION BY 
			event_sig, subkey_paths::text, address, event_idl, event_name, chain_id, split_part(name, '.', 1), split_part(name, '.', 2)
			ORDER BY
			CASE WHEN retention = 0 then 0 else 1 end,
				retention desc
		) as rn
		from solana.log_poller_filters lp
	) ids where rn = 1
)
UPDATE solana.log_poller_filters set is_deleted = true where id not in (select id from uniques);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE solana.log_poller_filters DROP COLUMN is_legacy_name;
-- +goose StatementEnd
