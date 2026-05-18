package keystore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func mustNewEthKey(t *testing.T) *ethkey.KeyV2 {
	key, err := ethkey.NewV2()
	require.NoError(t, err)
	return &key
}

func ExposedNewMaster(t *testing.T, ds sqlutil.DataSource) *master {
	return newMaster(ds, utils.FastScryptParams, logger.TestLogger(t))
}

func (m *master) ExportedSave(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.save(ctx)
}

func (m *master) ResetXXXTestOnly() {
	m.keyRing = newKeyRing()
	m.keyStates = newKeyStates()
	m.password = ""
}

func (m *master) SetPassword(pw string) {
	m.password = pw
}
