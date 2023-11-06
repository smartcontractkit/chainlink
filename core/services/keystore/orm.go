package keystore

import (
	"database/sql"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) ksORM {
	namedLogger := lggr.Named("KeystoreORM")
	return ksORM{
		q:    pg.NewQ(db, namedLogger, cfg),
		lggr: namedLogger,
	}
}

type ksORM struct {
	q    pg.Q
	lggr logger.Logger
}

func (orm ksORM) isEmpty() (bool, error) {
	var count int64
	err := orm.q.QueryRow("SELECT count(*) FROM encrypted_key_rings").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (orm ksORM) saveEncryptedKeyRing(kr *encryptedKeyRing, callbacks ...func(pg.Queryer) error) error {
	return orm.q.Transaction(func(tx pg.Queryer) error {
		_, err := tx.Exec(`
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

func (orm ksORM) getEncryptedKeyRing() (kr encryptedKeyRing, err error) {
	err = orm.q.Get(&kr, `SELECT * FROM encrypted_key_rings LIMIT 1`)
	if errors.Is(err, sql.ErrNoRows) {
		sql := `INSERT INTO encrypted_key_rings (encrypted_keys, updated_at) VALUES (NULL, NOW()) RETURNING *;`
		err2 := orm.q.Get(&kr, sql)

		if err2 != nil {
			return kr, err2
		}
	} else if err != nil {
		return kr, err
	}
	return kr, nil
}

func (orm ksORM) loadKeyStates() (*keyStates, error) {
	ks := newKeyStates()
	var ethkeystates []*ethkey.State
	if err := orm.q.Select(&ethkeystates, `SELECT id, address, evm_chain_id, disabled, created_at, updated_at FROM evm.key_states`); err != nil {
		return ks, errors.Wrap(err, "error loading evm.key_states from DB")
	}
	for _, state := range ethkeystates {
		ks.add(state)
	}
	return ks, nil
}
