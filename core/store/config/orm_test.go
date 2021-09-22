package config_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestORM_SetConfigStrValue(t *testing.T) {
	t.Parallel()
	db := pgtest.NewGormDB(t)
	orm := config.NewORM(db)

	fieldName := "LogSQLStatements"
	name := config.EnvVarName(fieldName)
	isSqlStatementEnabled := true
	res := models.Configuration{}

	// Store db config entry as true
	err := orm.SetConfigStrValue(context.TODO(), fieldName, strconv.FormatBool(isSqlStatementEnabled))
	require.NoError(t, err)

	err = db.First(&res, "name = ?", name).Error
	require.NoError(t, err)
	require.Equal(t, strconv.FormatBool(isSqlStatementEnabled), res.Value)

	// Update db config entry as false
	isSqlStatementEnabled = false
	err = orm.SetConfigStrValue(context.TODO(), fieldName, strconv.FormatBool(isSqlStatementEnabled))
	require.NoError(t, err)

	err = db.First(&res, "name = ?", name).Error
	require.NoError(t, err)
	require.Equal(t, strconv.FormatBool(isSqlStatementEnabled), res.Value)
}

func TestORM_GetConfigBoolValue(t *testing.T) {
	t.Parallel()
	db := pgtest.NewGormDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	cfg.SetDB(db)

	isSqlStatementEnabled := true
	err := cfg.SetLogSQLStatements(context.TODO(), isSqlStatementEnabled)
	require.NoError(t, err)
	assert.Equal(t, isSqlStatementEnabled, cfg.LogSQLStatements())
}
