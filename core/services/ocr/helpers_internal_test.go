package ocr

import (
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func (c *ConfigOverriderImpl) ExportedUpdateFlagsStatus() error {
	return c.updateFlagsStatus()
}

func NewTestDB(t *testing.T, sqldb *sqlx.DB, oracleSpecID int32) *db {
	return NewDB(sqldb, oracleSpecID, logger.TestLogger(t), pgtest.NewQConfig(true))
}
