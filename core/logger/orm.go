package logger

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
)

// LogConfig stores key value pairs for configuring package specific logging
type LogConfig struct {
	ID          int64
	ServiceName string
	LogLevel    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ORM interface {
	GetServiceLogLevel(serviceName string) (level string, ok bool)
	SetServiceLogLevel(ctx context.Context, serviceName string, level string) error
}

type orm struct {
	db   *sqlx.DB
	lggr Logger
}

// NewORM initializes a new ORM
func NewORM(db *sqlx.DB, lggr Logger) *orm {
	return &orm{db, lggr.Named("LoggerORM")}
}

// GetServiceLogLevel returns the log level for a configured service
func (orm *orm) GetServiceLogLevel(serviceName string) (string, bool) {
	config := LogConfig{}
	if err := orm.db.Get(&config, "SELECT * FROM log_configs WHERE service_name = $1", serviceName); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			orm.lggr.Errorf("Error while trying to fetch %s service log level: %v", serviceName, err)
		}
		return "", false
	}
	return config.LogLevel, true
}

func (orm *orm) SetServiceLogLevel(ctx context.Context, serviceName string, level string) error {
	_, err := orm.db.ExecContext(ctx, `
INSERT INTO log_configs (
	service_name, log_level, created_at, updated_at
) VALUES (
	$1, $2, NOW(), NOW()
) ON CONFLICT (service_name)
DO UPDATE SET log_level = EXCLUDED.log_level
    `, serviceName, level)
	return errors.Wrap(err, "LogOrm#SetServiceLogLevel failed")
}
