package vault

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/scylladb/go-reflectx"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/vault"
)

const defaultDBHost = "127.0.0.1"

func newVaultORM(nodeIndex int, host string, externalPort int) (vault.ORM, *sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, externalPort, postgres.User, postgres.Password, fmt.Sprintf("db_%d", nodeIndex))
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, db, err
	}

	db.MapperFunc(reflectx.CamelToSnakeASCII)
	return vault.NewVaultORM(db), db, nil
}

func GetResultPackageCount(ctx context.Context, nodeIndex, externalPort int) (int64, error) {
	orm, db, err := newVaultORM(nodeIndex, defaultDBHost, externalPort)
	if err != nil {
		return 0, err
	}

	defer db.Close()
	return orm.GetResultPackageCount(ctx)
}

func GetResultPackageCountRemoteAware(ctx context.Context, nodeIndex, externalPort int, isRemoteNodeSet bool) (int64, error) {
	host, err := resolveDBHostForNodeSet(isRemoteNodeSet)
	if err != nil {
		return 0, err
	}

	orm, db, err := newVaultORM(nodeIndex, host, externalPort)
	if err != nil {
		return 0, err
	}

	defer db.Close()
	return orm.GetResultPackageCount(ctx)
}

func resolveDBHostForNodeSet(isRemoteNodeSet bool) (string, error) {
	if !isRemoteNodeSet {
		return defaultDBHost, nil
	}
	if !runtimecfg.IsDirectMode() {
		return defaultDBHost, nil
	}
	return runtimecfg.DirectHostIP()
}
