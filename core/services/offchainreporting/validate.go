package offchainreporting

import (
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/libocr/offchainreporting"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"go.uber.org/multierr"
)

type ValidationConfig interface {
	Dev() bool
	OCRBlockchainTimeout(override time.Duration) time.Duration
	OCRContractConfirmations(override uint16) uint16
	OCRContractPollInterval(override time.Duration) time.Duration
	OCRContractSubscribeInterval(override time.Duration) time.Duration
	OCRContractTransmitterTransmitTimeout() time.Duration
	OCRDatabaseTimeout() time.Duration
	OCRObservationTimeout(override time.Duration) time.Duration
	OCRObservationGracePeriod() time.Duration
}

// ValidatedOracleSpecToml validates an oracle spec that came from TOML
func ValidatedOracleSpecToml(config ValidationConfig, tomlString string) (job.Job, error) {
	var jb = job.Job{}
	var spec job.OffchainReportingOracleSpec
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
	jb.OffchainreportingOracleSpec = &spec

	// TODO(#175801038): upstream support for time.Duration defaults in go-toml
	if jb.Type != job.OffchainReporting {
		return jb, errors.Errorf("the only supported type is currently 'offchainreporting', got %s", jb.Type)
	}
	if !tree.Has("isBootstrapPeer") {
		return jb, errors.New("isBootstrapPeer is not defined")
	}
	for i := range spec.P2PBootstrapPeers {
		if _, err := multiaddr.NewMultiaddr(spec.P2PBootstrapPeers[i]); err != nil {
			return jb, errors.Wrapf(err, "p2p bootstrap peer %v is invalid", spec.P2PBootstrapPeers[i])
		}
	}
	if spec.IsBootstrapPeer {
		if err := validateBootstrapSpec(tree, jb); err != nil {
			return jb, err
		}
	} else if err := validateNonBootstrapSpec(tree, config, jb); err != nil {
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
		"contractAddress": {},
		"isBootstrapPeer": {},
	}
	// Boostrap and non-bootstrap parameters
	// are mutually exclusive.
	bootstrapParams    = map[string]struct{}{}
	nonBootstrapParams = map[string]struct{}{
		"observationSource": {},
	}
)

func cloneSet(in map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{})
	for k, v := range in {
		out[k] = v
	}
	return out
}

func validateTimingParameters(config ValidationConfig, spec job.OffchainReportingOracleSpec) error {
	lc := types.LocalConfig{
		BlockchainTimeout:                      config.OCRBlockchainTimeout(time.Duration(spec.BlockchainTimeout)),
		ContractConfigConfirmations:            config.OCRContractConfirmations(spec.ContractConfigConfirmations),
		ContractConfigTrackerPollInterval:      config.OCRContractPollInterval(time.Duration(spec.ContractConfigTrackerPollInterval)),
		ContractConfigTrackerSubscribeInterval: config.OCRContractSubscribeInterval(time.Duration(spec.ContractConfigTrackerSubscribeInterval)),
		ContractTransmitterTransmitTimeout:     config.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        config.OCRDatabaseTimeout(),
		DataSourceTimeout:                      config.OCRObservationTimeout(time.Duration(spec.ObservationTimeout)),
		DataSourceGracePeriod:                  config.OCRObservationGracePeriod(),
	}
	if config.Dev() {
		lc.DevelopmentMode = types.EnableDangerousDevelopmentMode
	}

	return offchainreporting.SanityCheckLocalConfig(lc)
}

func validateBootstrapSpec(tree *toml.Tree, spec job.Job) error {
	expected, notExpected := cloneSet(params), cloneSet(nonBootstrapParams)
	for k := range bootstrapParams {
		expected[k] = struct{}{}
	}
	return validateExplicitlySetKeys(tree, expected, notExpected, "bootstrap")
}

func validateNonBootstrapSpec(tree *toml.Tree, config ValidationConfig, spec job.Job) error {
	expected, notExpected := cloneSet(params), cloneSet(bootstrapParams)
	for k := range nonBootstrapParams {
		expected[k] = struct{}{}
	}
	if err := validateExplicitlySetKeys(tree, expected, notExpected, "non-bootstrap"); err != nil {
		return err
	}
	if spec.Pipeline.Source == "" {
		return errors.New("no pipeline specified")
	}
	observationTimeout := config.OCRObservationTimeout(time.Duration(spec.OffchainreportingOracleSpec.ObservationTimeout))
	if time.Duration(spec.MaxTaskDuration) > observationTimeout {
		return errors.Errorf("max task duration must be < observation timeout")
	}
	for _, task := range spec.Pipeline.Tasks {
		timeout, set := task.TaskTimeout()
		if set && timeout > observationTimeout {
			return errors.Errorf("individual max task duration must be < observation timeout")
		}
	}
	return nil
}

func validateExplicitlySetKeys(tree *toml.Tree, expected map[string]struct{}, notExpected map[string]struct{}, peerType string) error {
	var err error
	// top level keys only
	for _, k := range tree.Keys() {
		// TODO(#175801577): upstream a way to check for children in go-toml
		if _, ok := notExpected[k]; ok {
			err = multierr.Append(err, errors.Errorf("unrecognised key for %s peer: %s", peerType, k))
		}
		delete(expected, k)
	}
	for missing := range expected {
		err = multierr.Append(err, errors.Errorf("missing required key %s", missing))
	}
	return err
}
