package migrations

import (
	"gorm.io/gorm"
)

const up39 = `

CREATE TABLE public.blocks (
    id bigserial PRIMARY KEY,
    hash bytea NOT NULL,
    number bigint NOT NULL,
    parent_hash bytea NOT NULL,
    transactions jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    CONSTRAINT chk_hash_size CHECK ((octet_length(hash) = 32)),
    CONSTRAINT chk_parent_hash_size CHECK ((octet_length(parent_hash) = 32))
);

CREATE UNIQUE INDEX idx_blocks_hash ON public.blocks USING btree (hash);

`
const down39 = `
DROP TABLE public.blocks;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0039_blocks_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up39).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down39).Error
		},
	})
}
