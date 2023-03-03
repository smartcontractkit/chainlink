package testsetups

import (
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
)

type FromVersionSpec struct {
	Image string
	Tag   string
}

type ToVersionSpec struct {
	Image string
	Tag   string
}

type DBMigrationSpec struct {
	FromSpec          FromVersionSpec
	ToSpec            ToVersionSpec
	KeepConnection    bool
	RemoveOnInterrupt bool
}

// DBMigration returns an environment with DB migrated from FromVersionSpec to ToVersionSpec
func DBMigration(spec *DBMigrationSpec) (*environment.Environment, error) {
	e := environment.New(nil).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"chainlink": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   spec.FromSpec.Image,
					"version": spec.FromSpec.Tag,
				},
			},
			"db": map[string]interface{}{
				"stateful": true,
				"capacity": "1Gi",
			},
		}))
	err := e.Run()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to setup initial deployment for version: %s:%s", spec.FromSpec.Image, spec.FromSpec.Tag)
	}
	e.Cfg.KeepConnection = spec.KeepConnection
	e.Cfg.RemoveOnInterrupt = spec.RemoveOnInterrupt
	e.Cfg.UpdateWaitInterval = 10 * time.Second
	err = e.
		ModifyHelm("chainlink-0", chainlink.New(0, map[string]interface{}{
			"chainlink": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   spec.ToSpec.Image,
					"version": spec.ToSpec.Tag,
				},
			},
			"db": map[string]interface{}{
				"stateful": true,
				"capacity": "1Gi",
			},
		})).Run()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to migrate to version: %s:%s", spec.ToSpec.Image, spec.ToSpec.Tag)
	}
	return e, nil
}
