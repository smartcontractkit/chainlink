package keystore

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func mustNewEthKey(t *testing.T) *ethkey.KeyV2 {
	key, err := ethkey.NewV2()
	require.NoError(t, err)
	return &key
}

func ExposedNewMaster(t *testing.T, db *sqlx.DB, cfg pg.QConfig) *master {
	return newMaster(db, utils.FastScryptParams, logger.TestLogger(t), cfg)
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

func (m *master) SetPassword(pw string) {
	m.password = pw
}
