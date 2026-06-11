package resolve

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/manifest"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/modulemap"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/ref"
)

// Resolver reads versions from manifests and formats install targets.
type Resolver struct {
	store *manifest.Store
}

func New(store *manifest.Store) *Resolver {
	return &Resolver{store: store}
}

func (r *Resolver) Get(key string) (string, error) {
	return r.store.Lookup(key)
}

func (r *Resolver) Ref(key string) (string, error) {
	v, err := r.store.Lookup(key)
	if err != nil {
		return "", err
	}
	return ref.ForConsumer(v), nil
}

func (r *Resolver) Target(key string) (string, error) {
	var module string
	if strings.Contains(key, "/") {
		module = key
	} else {
		var err error
		module, err = modulemap.ModulePath(key)
		if err != nil {
			return "", err
		}
	}
	version, err := r.store.Lookup(key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s@%s", module, ref.ForInstall(version)), nil
}

func (r *Resolver) List() ([]manifest.Entry, error) {
	return r.store.List()
}

func (r *Resolver) ManagedModules() []string {
	mods := modulemap.Modules()
	mods = append(mods, r.store.GoToolModules()...)
	return mods
}
