package s4_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/stretchr/testify/assert"
)

func setupORM(t *testing.T) s4.ORM {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	return s4.NewPostgresORM(db, lggr, pgtest.NewQConfig(true), "test")
}

func TestNewPostgresOrm(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	assert.NotNil(t, orm)
}
