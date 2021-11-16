package offchainreporting

import (
	"database/sql"
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"
)

func (c *ConfigOverriderImpl) ExportedUpdateFlagsStatus() error {
	return c.updateFlagsStatus()
}

func NewTestDB(t *testing.T, sqldb *sql.DB, oracleSpecID int32) *db {
	return NewDB(sqldb, oracleSpecID, logger.TestLogger(t))
}
