package keystore_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_ORM(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := pgtest.NewPGCfg(false)

	orm := keystore.NewORM(db, logger.TestLogger(t), cfg)

	csa := csakey.MustNewV2XXXTestingOnly(big.NewInt(42))
	p1 := cltest.MustRandomP2PPeerID(t)
	p2 := cltest.MustRandomP2PPeerID(t)
	v1 := vrfkey.MustNewV2XXXTestingOnly(big.NewInt(1))
	v2 := vrfkey.MustNewV2XXXTestingOnly(big.NewInt(2))
	err := multierr.Combine(
		// csakeys has no deleted_at column
		utils.JustError(db.Exec(`INSERT INTO csa_keys (public_key, encrypted_private_key, created_at, updated_at) VALUES ($1, '{}', NOW(), NOW())`, csa.PublicKey)),

		// two per key-type, one deleted and one not deleted
		utils.JustError(db.Exec(`INSERT INTO keys (address, json, created_at, updated_at, next_nonce, is_funding, deleted_at) VALUES ($1, '{}', NOW(), NOW(), 0, false, NULL)`, cltest.NewAddress())),
		utils.JustError(db.Exec(`INSERT INTO keys (address, json, created_at, updated_at, next_nonce, is_funding, deleted_at) VALUES ($1, '{}', NOW(), NOW(), 0, false, NOW())`, cltest.NewAddress())),
		utils.JustError(db.Exec(`INSERT INTO encrypted_ocr_key_bundles (id, on_chain_signing_address, off_chain_public_key, encrypted_private_keys, created_at, updated_at, config_public_key, deleted_at) VALUES ($1, $2, $3, '{}', NOW(), NOW(), $4, NULL)`, utils.NewHash(), cltest.NewAddress(), utils.NewHash(), utils.NewHash())),
		utils.JustError(db.Exec(`INSERT INTO encrypted_ocr_key_bundles (id, on_chain_signing_address, off_chain_public_key, encrypted_private_keys, created_at, updated_at, config_public_key, deleted_at) VALUES ($1, $2, $3, '{}', NOW(), NOW(), $4, NOW())`, utils.NewHash(), cltest.NewAddress(), utils.NewHash(), utils.NewHash())),
		utils.JustError(db.Exec(`INSERT INTO encrypted_p2p_keys (peer_id, pub_key, encrypted_priv_key, created_at, updated_at, deleted_at) VALUES ($1, $2, '{}', NOW(), NOW(), NULL)`, p1.Pretty(), utils.NewHash())),
		utils.JustError(db.Exec(`INSERT INTO encrypted_p2p_keys (peer_id, pub_key, encrypted_priv_key, created_at, updated_at, deleted_at) VALUES ($1, $2, '{}', NOW(), NOW(), NOW())`, p2.Pretty(), utils.NewHash())),
		utils.JustError(db.Exec(`INSERT INTO encrypted_vrf_keys (public_key, vrf_key, created_at, updated_at, deleted_at) VALUES ($1, '{}',  NOW(), NOW(), NULL)`, v1.PublicKey)),
		utils.JustError(db.Exec(`INSERT INTO encrypted_vrf_keys (public_key, vrf_key, created_at, updated_at, deleted_at) VALUES ($1, '{}',  NOW(), NOW(), NOW())`, v2.PublicKey)),
	)
	require.NoError(t, err)

	t.Run("legacy functions for V1 migration", func(t *testing.T) {
		t.Run("GetEncryptedV1CSAKeys", func(t *testing.T) {
			ks, err := orm.GetEncryptedV1CSAKeys()
			require.NoError(t, err)
			assert.Len(t, ks, 1)
		})
		t.Run("GetEncryptedV1EthKeys", func(t *testing.T) {
			ks, err := orm.GetEncryptedV1EthKeys()
			require.NoError(t, err)
			assert.Len(t, ks, 1)
		})
		t.Run("GetEncryptedV1OCRKeys", func(t *testing.T) {
			ks, err := orm.GetEncryptedV1OCRKeys()
			require.NoError(t, err)
			assert.Len(t, ks, 1)
		})
		t.Run("GetEncryptedV1P2PKeys", func(t *testing.T) {
			ks, err := orm.GetEncryptedV1P2PKeys()
			require.NoError(t, err)
			assert.Len(t, ks, 1)
		})
		t.Run("GetEncryptedV1VRFKeys", func(t *testing.T) {
			ks, err := orm.GetEncryptedV1VRFKeys()
			require.NoError(t, err)
			assert.Len(t, ks, 1)
		})
	})

}
