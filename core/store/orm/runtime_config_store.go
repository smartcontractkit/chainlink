package orm

import (
	"encoding"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// RuntimeConfigStore represents configuration values specified while chainlink is running
type RuntimeConfigStore struct {
	store *BootstrapConfigStore
	orm   *ORM
}

// NewRuntimeConfigStore returns a config store which prefers saved values in the orm
func NewRuntimeConfigStore(store *BootstrapConfigStore, orm *ORM) *RuntimeConfigStore {
	return &RuntimeConfigStore{
		store: store,
		orm:   orm,
	}
}

// Get returns the setting from the orm store, and if not found from the config store
func (r *RuntimeConfigStore) Get(name string, value encoding.TextUnmarshaler) error {
	if err := r.orm.GetConfigValue(name, value); err == nil {
		return nil
	} else if errors.Cause(err) == gorm.ErrRecordNotFound {
		r.store.Get(name, value)
		return nil
	} else {
		return err
	}
}

// Set saves the value in the orm store
func (r *RuntimeConfigStore) Set(name string, value encoding.TextMarshaler) error {
	return r.orm.SetConfigValue(name, value)
}
