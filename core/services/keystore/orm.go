package keystore

import (
	"database/sql"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"

	"github.com/ethereum/go-ethereum/common"
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
	if err := orm.q.Select(&ethkeystates, `SELECT id, address, evm_chain_id, next_nonce, disabled, created_at, updated_at FROM evm_key_states`); err != nil {
		return ks, errors.Wrap(err, "error loading evm_key_states from DB")
	}
	for _, state := range ethkeystates {
		ks.add(state)
	}
	return ks, nil
}

// getNextNonce returns evm_key_states.next_nonce for the given address
func (orm ksORM) getNextNonce(address common.Address, chainID *big.Int, qopts ...pg.QOpt) (nonce int64, err error) {
	q := orm.q.WithOpts(qopts...)
	err = q.Get(&nonce, "SELECT next_nonce FROM evm_key_states WHERE address = $1 AND evm_chain_id = $2 AND disabled = false", address, chainID.String())
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.Wrapf(sql.ErrNoRows, "key with address %s is not enabled for chain %s", address.Hex(), chainID.String())
	}
	return nonce, errors.Wrap(err, "failed to load next nonce")
}

// incrementNextNonce increments evm_key_states.next_nonce by 1
func (orm ksORM) incrementNextNonce(address common.Address, chainID *big.Int, currentNonce int64, qopts ...pg.QOpt) (incrementedNonce int64, err error) {
	q := orm.q.WithOpts(qopts...)
	err = q.Get(&incrementedNonce, "UPDATE evm_key_states SET next_nonce = next_nonce + 1, updated_at = NOW() WHERE address = $1 AND next_nonce = $2 AND evm_chain_id = $3 AND disabled = false RETURNING next_nonce", address, currentNonce, chainID.String())
	return incrementedNonce, errors.Wrap(err, "IncrementNextNonce failed to update keys")
}

// ~~~~~~~~~~~~~~~~~~~~ LEGACY FUNCTIONS FOR V1 MIGRATION ~~~~~~~~~~~~~~~~~~~~

func (orm ksORM) GetEncryptedV1CSAKeys() (retrieved []csakey.Key, err error) {
	return retrieved, orm.q.Select(&retrieved, `SELECT * FROM csa_keys`)
}

func (orm ksORM) GetEncryptedV1EthKeys() (retrieved []ethkey.Key, err error) {
	return retrieved, orm.q.Select(&retrieved, `SELECT * FROM keys WHERE deleted_at IS NULL`)
}

func (orm ksORM) GetEncryptedV1OCRKeys() (retrieved []ocrkey.EncryptedKeyBundle, err error) {
	return retrieved, orm.q.Select(&retrieved, `SELECT * FROM encrypted_ocr_key_bundles WHERE deleted_at IS NULL`)
}

func (orm ksORM) GetEncryptedV1P2PKeys() (retrieved []p2pkey.EncryptedP2PKey, err error) {
	return retrieved, orm.q.Select(&retrieved, `SELECT * FROM encrypted_p2p_keys WHERE deleted_at IS NULL`)
}

func (orm ksORM) GetEncryptedV1VRFKeys() (retrieved []vrfkey.EncryptedVRFKey, err error) {
	return retrieved, orm.q.Select(&retrieved, `SELECT * FROM encrypted_vrf_keys WHERE deleted_at IS NULL`)
}
