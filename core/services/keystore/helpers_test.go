package keystore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ethkey"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/models"
	"github.com/smartcontractkit/chainlink-common/keystore/scrypt"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

func mustNewEthKey(t *testing.T) *ethkey.KeyV2 {
	key, err := ethkey.NewV2()
	require.NoError(t, err)
	return &key
}

func ExposedNewMaster(t *testing.T, ds sqlutil.DataSource) *master {
	return newMaster(ds, scrypt.FastScryptParams, logger.Test(t).Infof)
}

func (m *master) ExportedSave(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.save(ctx)
}

func (m *master) ResetXXXTestOnly() {
	m.keyRing = models.NewKeyRing()
	m.keyStates = models.NewKeyStates()
	m.password = ""
}

func (m *master) SetPassword(pw string) {
	m.password = pw
}
