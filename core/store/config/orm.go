package config

import (
	"context"
	"encoding"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
)

type ORM struct {
	db *gorm.DB
}

func NewORM(db *gorm.DB) *ORM {
	return &ORM{db}
}

// GetConfigValue returns the value for a named configuration entry
func (orm *ORM) GetConfigValue(field string, value encoding.TextUnmarshaler) error {
	name := EnvVarName(field)
	config := models.Configuration{}
	if err := orm.db.First(&config, "name = ?", name).Error; err != nil {
		return err
	}
	return value.UnmarshalText([]byte(config.Value))
}

// GetConfigBoolValue returns a boolean value for a named configuration entry
func (orm *ORM) GetConfigBoolValue(field string) (*bool, error) {
	name := EnvVarName(field)
	config := models.Configuration{}
	if err := orm.db.First(&config, "name = ?", name).Error; err != nil {
		return nil, err
	}
	value, err := strconv.ParseBool(config.Value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// SetConfigValue returns the value for a named configuration entry
func (orm *ORM) SetConfigValue(field string, value encoding.TextMarshaler) error {
	name := EnvVarName(field)
	textValue, err := value.MarshalText()
	if err != nil {
		return err
	}
	return orm.db.Where(models.Configuration{Name: name}).
		Assign(models.Configuration{Name: name, Value: string(textValue)}).
		FirstOrCreate(&models.Configuration{}).Error
}

// SetConfigValue returns the value for a named configuration entry
func (orm *ORM) SetConfigStrValue(ctx context.Context, field string, value string) error {
	name := EnvVarName(field)
	return orm.db.WithContext(ctx).Where(models.Configuration{Name: name}).
		Assign(models.Configuration{Name: name, Value: value}).
		FirstOrCreate(&models.Configuration{}).Error
}
