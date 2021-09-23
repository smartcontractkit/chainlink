package logger

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ORM interface {
	GetServiceLogLevel(serviceName string) (level string, ok bool)
	SetServiceLogLevel(ctx context.Context, serviceName string, level string) error
}

type orm struct {
	DB *gorm.DB
}

// NewORM initializes a new ORM
func NewORM(db *gorm.DB) *orm {
	return &orm{db}
}

// GetServiceLogLevel returns the log level for a configured service
func (orm *orm) GetServiceLogLevel(serviceName string) (string, bool) {
	config := LogConfig{}
	if err := orm.DB.First(&config, "service_name = ?", serviceName).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			Warnf("Error while trying to fetch %s service log level: %v", serviceName, err)
		}
		return "", false
	}
	return config.LogLevel, true
}

func (orm *orm) SetServiceLogLevel(ctx context.Context, serviceName string, level string) error {
	return orm.DB.WithContext(ctx).Exec(`
        INSERT INTO log_configs (
            service_name, log_level
        ) VALUES (
            ?, ?
        ) ON CONFLICT (service_name) 
		DO UPDATE SET log_level = EXCLUDED.log_level
    `, serviceName, level).Error
}
