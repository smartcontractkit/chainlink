package postgres

import (
	"database/sql"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func IsSerializationAnomaly(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(errors.Cause(err).Error(), "could not serialize access due to concurrent update")
}

func IsRecordNotFound(err error) bool {
	if err == nil {
		return false
	}
	return gorm.IsRecordNotFoundError(errors.Cause(err))
}

var (
	DefaultSqlTxOptions = sql.TxOptions{
		// NOTE: This is the default level in Postgres anyway, we just make it
		// explicit here
		Isolation: sql.LevelReadCommitted,
	}
)
