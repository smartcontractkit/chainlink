package keystore

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

func NewORM(ds sqlutil.DataSource, lggr logger.Logger) ksORM {
	namedLogger := lggr.Named("KeystoreORM")
	return ksORM{
		ds:   ds,
		lggr: namedLogger,
	}
}

type ksORM struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

func (orm ksORM) isEmpty(ctx context.Context) (bool, error) {
	var count int64
	err := orm.ds.QueryRowxContext(ctx, "SELECT count(*) FROM encrypted_key_rings").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (orm ksORM) saveEncryptedKeyRing(ctx context.Context, kr *encryptedKeyRing, callbacks ...func(sqlutil.DataSource) error) error {
	return sqlutil.TransactDataSource(ctx, orm.ds, nil, func(tx sqlutil.DataSource) error {
		_, err := tx.ExecContext(ctx, `
		UPDATE encrypted_key_rings
		SET encrypted_keys = $1
	`, kr.EncryptedKeys)
		if err != nil {
			return errors.Wrap(err, "while saving keyring")
		}
		for _, callback := range callbacks {
			err = callback(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (orm ksORM) getEncryptedKeyRing(ctx context.Context) (kr encryptedKeyRing, err error) {
	err = orm.ds.GetContext(ctx, &kr, `SELECT * FROM encrypted_key_rings LIMIT 1`)
	if errors.Is(err, sql.ErrNoRows) {
		sql := `INSERT INTO encrypted_key_rings (encrypted_keys, updated_at) VALUES (NULL, NOW()) RETURNING *;`
		err2 := orm.ds.GetContext(ctx, &kr, sql)

		if err2 != nil {
			return kr, err2
		}
	} else if err != nil {
		return kr, err
	}
	return kr, nil
}

func (orm ksORM) loadKeyStates(ctx context.Context) (*keyStates, error) {
	ks := newKeyStates()
	var ethkeystates []*ethkey.State
	if err := orm.ds.SelectContext(ctx, &ethkeystates, `SELECT id, address, evm_chain_id, disabled, created_at, updated_at FROM evm.key_states`); err != nil {
		return ks, errors.Wrap(err, "error loading evm.key_states from DB")
	}
	for _, state := range ethkeystates {
		ks.add(state)
	}
	return ks, nil
}
