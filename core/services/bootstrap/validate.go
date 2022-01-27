package bootstrap

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"go.uber.org/multierr"
)

// ValidatedBootstrapSpecToml validates a bootstrap spec that came from TOML
func ValidatedBootstrapSpecToml(tomlString string) (job.Job, error) {
	var jb = job.Job{}
	var spec job.BootstrapSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load")
	}
	// Note this validates all the fields which implement an UnmarshalText
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on spec")
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}
	jb.BootstrapSpec = &spec

	if jb.Type != job.Bootstrap {
		return jb, errors.Errorf("the only supported type is currently 'bootstrap', got %s", jb.Type)
	}

	if err := validateExplicitlySetKeys(tree, cloneSet(params), cloneSet(nonBootstrapParams)); err != nil {
		return jb, err
	}

	return jb, nil
}

// Parameters that must be explicitly set by the operator.
var (
	params = map[string]struct{}{
		"type":          {},
		"schemaVersion": {},
		"contractID":    {},
		"relay":         {},
		"relayConfig":   {},
	}
	// Parameters that should not be set
	nonBootstrapParams = map[string]struct{}{
		"isBootstrapPeer":       {},
		"juelsPerFeeCoinSource": {},
	}
)

func cloneSet(in map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{})
	for k, v := range in {
		out[k] = v
	}
	return out
}

func validateExplicitlySetKeys(tree *toml.Tree, expected map[string]struct{}, notExpected map[string]struct{}) error {
	var err error
	// top level keys only
	for _, k := range tree.Keys() {
		if _, ok := notExpected[k]; ok {
			err = multierr.Append(err, errors.Errorf("unrecognised key %s", k))
		}
		delete(expected, k)
	}
	for missing := range expected {
		err = multierr.Append(err, errors.Errorf("missing required key %s", missing))
	}
	return err
}
