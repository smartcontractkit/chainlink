package validate

import (
	"github.com/lib/pq"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/relay"
)

// ValidatedOracleSpecToml validates an oracle spec that came from TOML
func ValidatedOracleSpecToml(config Config, tomlString string) (job.Job, error) {
	var jb = job.Job{}
	var spec job.OCR2OracleSpec
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
	jb.OCR2OracleSpec = &spec
	if jb.OCR2OracleSpec.P2PBootstrapPeers == nil {
		// Empty but non-null, field is non-nullable.
		jb.OCR2OracleSpec.P2PBootstrapPeers = pq.StringArray{}
	}

	if jb.Type != job.OffchainReporting2 {
		return jb, errors.Errorf("the only supported type is currently 'offchainreporting2', got %s", jb.Type)
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

	if err := validateSpec(tree, jb); err != nil {
		return jb, err
	}
	if err := validateTimingParameters(config, spec); err != nil {
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
		"pluginType":    {},
		"pluginConfig":  {},
	}
	notExpectedParams = map[string]struct{}{
		"isBootstrapPeer":       {},
		"juelsPerFeeCoinSource": {},
	}
)

func validateTimingParameters(config Config, spec job.OCR2OracleSpec) error {
	lc := ToLocalConfig(config, spec)
	return libocr2.SanityCheckLocalConfig(lc)
}

func validateSpec(tree *toml.Tree, spec job.Job) error {
	expected, notExpected := ocrcommon.CloneSet(params), ocrcommon.CloneSet(notExpectedParams)
	if err := ocrcommon.ValidateExplicitlySetKeys(tree, expected, notExpected, "ocr2"); err != nil {
		return err
	}

	switch spec.OCR2OracleSpec.PluginType {
	case job.Median:
		if spec.Pipeline.Source == "" {
			return errors.New("no pipeline specified")
		}
	case "":
		return errors.New("no plugin specified")
	default:
		return errors.Errorf("invalid pluginType %s", spec.OCR2OracleSpec.PluginType)
	}

	return nil
}
