package keystore

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func mustNewEthKey(t *testing.T) *ethkey.KeyV2 {
	key, err := ethkey.NewV2()
	require.NoError(t, err)
	return &key
}

type ExportedEncryptedKeyRing = encryptedKeyRing

func ExposedNewMaster(db *gorm.DB) *master {
	return newMaster(db, utils.FastScryptParams)
}

func (m *master) ExportedSave() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.save()
}

func (m *master) ResetXXXTestOnly() {
	m.keyRing = newKeyRing()
	m.keyStates = newKeyStates()
	m.password = ""
}
