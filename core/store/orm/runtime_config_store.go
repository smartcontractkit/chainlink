package orm

import (
	"encoding"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// RuntimeConfigStore represents configuration values specified while chainlink is running
type RuntimeConfigStore struct {
	store ConfigStore
	orm   *ORM
}

// NewRuntimeConfigStore returns a config store which prefers saved values in the orm
func NewRuntimeConfigStore(store ConfigStore, orm *ORM) *RuntimeConfigStore {
	return &RuntimeConfigStore{
		store: store,
		orm:   orm,
	}
}

// Get returns the setting from the orm store, and if not found from the config store
func (r RuntimeConfigStore) Get(name string, value encoding.TextUnmarshaler) error {
	if err := r.orm.GetConfigValue(name, value); err == nil {
		return nil
	} else if errors.Cause(err) == gorm.ErrRecordNotFound {
		r.store.Get(name, value)
		return nil
	} else {
		return err
	}
}

// SetMarshaler saves the value in the config store (using the ORM)
func (r RuntimeConfigStore) SetMarshaler(name string, value encoding.TextMarshaler) error {
	return r.orm.SetConfigValue(name, value)
}

// SetString saves the string in the config store
func (r RuntimeConfigStore) SetString(name, value string) error {
	return r.SetMarshaler(name, StringMarshaler(value))
}

// SetStringer saves the string erin the config store
func (r RuntimeConfigStore) SetStringer(name string, value fmt.Stringer) error {
	return r.SetString(name, value.String())
}
