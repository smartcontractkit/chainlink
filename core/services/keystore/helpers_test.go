package keystore

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/models"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

func ExposedNewMaster(t *testing.T, ds sqlutil.DataSource) *master {
	return newMaster(ds, keystore.FastScryptParams, logger.Test(t).Infof)
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
