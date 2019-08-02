package orm

import (
	"encoding"
	"fmt"
	"os"
	"reflect"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/spf13/viper"
)

// BootstrapConfigStore is the initial configuration implementation that is used by clients and to bootstrap the runtime configuration store
type BootstrapConfigStore struct {
	viper *viper.Viper
}

var configFileNotFoundError = reflect.TypeOf(viper.ConfigFileNotFoundError{})

// NewBootstrapConfigStore returns a config store that primarily uses Viper for serving values
func NewBootstrapConfigStore() BootstrapConfigStore {
	v := viper.New()
	return newConfigWithViper(v)
}

func newConfigWithViper(v *viper.Viper) BootstrapConfigStore {
	schemaT := reflect.TypeOf(Schema{})
	for index := 0; index < schemaT.NumField(); index++ {
		item := schemaT.FieldByIndex([]int{index})
		name := item.Tag.Get("env")
		v.SetDefault(name, item.Tag.Get("default"))
		v.BindEnv(name, name)
	}

	return BootstrapConfigStore{viper: v}
}

// LoadConfigFile loads the root dir's config file if it's present
func (c *BootstrapConfigStore) LoadConfigFile(rootDir string) error {
	if err := os.MkdirAll(rootDir, os.FileMode(0700)); err != nil {
		return errors.Wrap(err, "Error creating rooot directory")
	}

	c.viper.SetConfigName("chainlink")
	c.viper.AddConfigPath(rootDir)
	err := c.viper.ReadInConfig()
	if err != nil && reflect.TypeOf(err) != configFileNotFoundError {
		return errors.Wrap(err, "Unable to load config file")
	}

	return nil
}

// SetMarshaler a specific configuration variable, takes precedence over all other values
func (c BootstrapConfigStore) SetMarshaler(name string, value encoding.TextMarshaler) error {
	schemaT := reflect.TypeOf(Schema{})
	for index := 0; index < schemaT.NumField(); index++ {
		item := schemaT.FieldByIndex([]int{index})
		envName := item.Tag.Get("env")
		if envName == name {
			c.viper.Set(name, value)
			return nil
		}
	}
	logger.Panicf("No configuration parameter for %s", name)
	return nil // unreachable
}

// SetString saves the string in the config store
func (c BootstrapConfigStore) SetString(name, value string) error {
	return c.SetMarshaler(name, StringMarshaler(value))
}

// SetStringer saves the string erin the config store
func (c BootstrapConfigStore) SetStringer(name string, value fmt.Stringer) error {
	return c.SetString(name, value.String())
}

// Get a value by name
func (c BootstrapConfigStore) Get(name string, value encoding.TextUnmarshaler) error {
	source := c.viper.GetString(name)
	fmt.Println("name", name, "source", source)
	return value.UnmarshalText([]byte(source))
}
