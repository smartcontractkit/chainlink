package vault

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/scylladb/go-reflectx"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/vault"
)

func newVaultORM(nodeIndex, externalPort int) (vault.ORM, *sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", externalPort, postgres.User, postgres.Password, fmt.Sprintf("db_%d", nodeIndex))
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, db, err
	}

	db.MapperFunc(reflectx.CamelToSnakeASCII)
	return vault.NewVaultORM(db), db, nil
}

func GetResultPackageCount(ctx context.Context, nodeIndex, externalPort int) (int64, error) {
	orm, db, err := newVaultORM(nodeIndex, externalPort)
	if err != nil {
		return 0, err
	}

	defer db.Close()
	return orm.GetResultPackageCount(ctx)
}
