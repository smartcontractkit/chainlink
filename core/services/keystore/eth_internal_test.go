package keystore

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EthKeyStore(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	keyStore := ExposedNewMaster(t, db, pgtest.NewQConfig(true))
	err := keyStore.Unlock(testutils.Password)
	require.NoError(t, err)
	ks := keyStore.Eth()

	t.Run("returns V1 keys as V2", func(t *testing.T) {
		ethAddress := testutils.NewAddress()
		err = utils.JustError(db.Exec(`INSERT INTO keys (address, json, created_at, updated_at, next_nonce, is_funding, deleted_at) VALUES ($1, '{"address":"6fdac88ddfd811d130095373986889ed90e0d622","crypto":{"cipher":"aes-128-ctr","ciphertext":"557f5324e770c3d203751c1f0f7fb5076386c49f5b05e3f20b3abb59758fd3c3","cipherparams":{"iv":"bd9472543fab7cc63027cbcd039daff0"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":2,"p":1,"r":8,"salt":"647b54770a3fda830b4440ae57c44cf7506297295fe4d72b1ff943e3a8ddb94a"},"mac":"0c654ee29ee06b3816fc0040d84ebd648c557144a77ccc55b9568355f53397b3"},"id":"6fdac88d-dfd8-11d1-3009-5373986889ed","version":3}', NOW(), NOW(), 0, false, NULL)`, ethAddress))
		require.NoError(t, err)

		keys, nonces, fundings, err := ks.(*eth).getV1KeysAsV2()
		require.NoError(t, err)

		assert.Len(t, keys, 1)
		assert.Equal(t, fmt.Sprintf("EthKeyV2{PrivateKey: <redacted>, Address: %s}", keys[0].Address), keys[0].GoString())
		assert.Equal(t, int64(0), nonces[0])
		assert.Equal(t, false, fundings[0])
	})
}
