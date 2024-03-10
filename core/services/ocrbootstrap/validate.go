package ocrbootstrap

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

// ValidatedBootstrapSpecToml validates a bootstrap spec that came from TOML
func ValidatedBootstrapSpecToml(tomlString string) (jb job.Job, err error) {
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
	expected, notExpected := ocrcommon.CloneSet(params), ocrcommon.CloneSet(nonBootstrapParams)
	if err := ocrcommon.ValidateExplicitlySetKeys(tree, expected, notExpected, "bootstrap"); err != nil {
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
		"observationSource":     {},
	}
)
