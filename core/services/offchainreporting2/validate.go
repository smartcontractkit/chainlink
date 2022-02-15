package offchainreporting2

import (
	"github.com/lib/pq"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"
)

// ValidatedOracleSpecToml validates an oracle spec that came from TOML
func ValidatedOracleSpecToml(config Config, tomlString string) (job.Job, error) {
	var jb = job.Job{}
	var spec job.OffchainReporting2OracleSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load")
	}
	// Note this validates all the fields which implement an UnmarshalText
	// i.e. TransmitterAddress, PeerID...
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on spec")
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}
	jb.Offchainreporting2OracleSpec = &spec
	if jb.Offchainreporting2OracleSpec.P2PBootstrapPeers == nil {
		// Empty but non-null, field is non-nullable.
		jb.Offchainreporting2OracleSpec.P2PBootstrapPeers = pq.StringArray{}
	}

	if jb.Type != job.OffchainReporting2 {
		return jb, errors.Errorf("the only supported type is currently 'offchainreporting2', got %s", jb.Type)
	}
	if !tree.Has("isBootstrapPeer") {
		return jb, errors.New("isBootstrapPeer is not defined")
	}
	if _, ok := relay.SupportedRelayers[spec.Relay]; !ok {
		return jb, errors.Errorf("no such relay %v supported", spec.Relay)
	}
	if len(spec.P2PBootstrapPeers) > 0 {
		_, err := ocrcommon.ParseBootstrapPeers(spec.P2PBootstrapPeers)
		if err != nil {
			return jb, err
		}
	}

	if spec.IsBootstrapPeer {
		if err := validateBootstrapSpec(tree); err != nil {
			return jb, err
		}
	} else if err := validateNonBootstrapSpec(tree, jb); err != nil {
		return jb, err
	}
	if err := validateTimingParameters(config, spec); err != nil {
		return jb, err
	}
	return jb, nil
}

// Parameters that must be explicitly set by the operator.
var (
	// Common to both bootstrap and non-boostrap
	params = map[string]struct{}{
		"type":            {},
		"schemaVersion":   {},
		"contractID":      {},
		"isBootstrapPeer": {},
		"relay":           {},
		"relayConfig":     {},
	}
	// Boostrap and non-bootstrap parameters
	// are mutually exclusive.
	bootstrapParams    = map[string]struct{}{}
	nonBootstrapParams = map[string]struct{}{
		"observationSource":     {},
		"juelsPerFeeCoinSource": {},
	}
)

func validateTimingParameters(config Config, spec job.OffchainReporting2OracleSpec) error {
	lc := ToLocalConfig(config, spec)
	return libocr2.SanityCheckLocalConfig(lc)
}

func validateBootstrapSpec(tree *toml.Tree) error {
	expected, notExpected := ocrcommon.CloneSet(params), ocrcommon.CloneSet(nonBootstrapParams)
	for k := range bootstrapParams {
		expected[k] = struct{}{}
	}
	return ocrcommon.ValidateExplicitlySetKeys(tree, expected, notExpected, "bootstrap")
}

func validateNonBootstrapSpec(tree *toml.Tree, spec job.Job) error {
	expected, notExpected := ocrcommon.CloneSet(params), ocrcommon.CloneSet(bootstrapParams)
	for k := range nonBootstrapParams {
		expected[k] = struct{}{}
	}
	if err := ocrcommon.ValidateExplicitlySetKeys(tree, expected, notExpected, "non-bootstrap"); err != nil {
		return err
	}
	if spec.Pipeline.Source == "" {
		return errors.New("no pipeline specified")
	}
	// validate that the JuelsPerFeeCoinPipeline is valid (not checked later because it's not a normal pipeline)
	if _, err := pipeline.Parse(spec.Offchainreporting2OracleSpec.JuelsPerFeeCoinPipeline); err != nil {
		return errors.Wrap(err, "invalid juelsPerFeeCoinSource pipeline")
	}

	return nil
}
