package logger

import (
	"context"

	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type ORM interface {
	GetServiceLogLevel(serviceName string) (string, error)
	SetServiceLogLevel(ctx context.Context, serviceName string, level zapcore.Level) error
}

type orm struct {
	DB *gorm.DB
}

// NewORM initializes a new ORM
func NewORM(db *gorm.DB) *orm {
	return &orm{db}
}

// GetServiceLogLevel returns the log level for a configured service
func (orm *orm) GetServiceLogLevel(serviceName string) (string, error) {
	config := LogConfig{}
	if err := orm.DB.First(&config, "service_name = ?", serviceName).Error; err != nil {
		return "", err
	}
	return config.LogLevel, nil
}

func (orm *orm) SetServiceLogLevel(ctx context.Context, serviceName string, level zapcore.Level) error {
	return orm.DB.WithContext(ctx).Exec(`
        INSERT INTO log_configs (
            service_name, log_level
        ) VALUES (
            ?, ?
        ) ON CONFLICT (service_name) 
		DO UPDATE SET log_level = EXCLUDED.log_level
    `, serviceName, level.String()).Error
}
