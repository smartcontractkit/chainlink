package migrations

import (
	"gorm.io/gorm"
)

const up38 = `
CREATE TABLE csa_keys(
    id BIGSERIAL PRIMARY KEY,
    public_key bytea NOT NULL CHECK (octet_length(public_key) = 32) UNIQUE,
    encrypted_private_key jsonb NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

`
const down38 = `
DROP TABLE csa_keys;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0038_create_csa_keys",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up38).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down38).Error
		},
	})
}
