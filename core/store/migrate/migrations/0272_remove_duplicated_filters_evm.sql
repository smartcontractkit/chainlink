-- +goose Up
-- +goose StatementBegin
ALTER TABLE evm.log_poller_filters ADD COLUMN IF NOT EXISTS is_legacy_name BOOLEAN DEFAULT FALSE;
UPDATE evm.log_poller_filters SET is_legacy_name = true;
with uniques as (
	select id from (
		select id,
		ROW_NUMBER() OVER (
			PARTITION BY 
			topic2, topic3, topic4, address, event, evm_chain_id, split_part(name, '.', 1)
			ORDER BY
			CASE WHEN retention = 0 then 0 else 1 end,
				retention desc
		) as rn
		from evm.log_poller_filters lp
	) ids where rn = 1
)
DELETE from evm.log_poller_filters where id not in (select id from uniques);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE evm.log_poller_filters DROP COLUMN is_legacy_name;
-- +goose StatementEnd
