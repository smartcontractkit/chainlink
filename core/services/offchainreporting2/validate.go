package offchainreporting2

import (
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/services/job"
	offchainreporting "github.com/smartcontractkit/libocr/offchainreporting2"
	"go.uber.org/multierr"
)

type ValidationConfig interface {
	Dev() bool
	OCR2BlockchainTimeout() time.Duration
	OCR2ContractConfirmations() uint16
	OCR2ContractPollInterval() time.Duration
	OCR2ContractTransmitterTransmitTimeout() time.Duration
	OCR2DatabaseTimeout() time.Duration
	OCR2ObservationTimeout() time.Duration
}

// ValidatedOracleSpecToml validates an oracle spec that came from TOML
func ValidatedOracleSpecToml(chainSet evm.ChainSet, tomlString string) (job.Job, error) {
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

	// TODO(#175801038): upstream support for time.Duration defaults in go-toml
	if jb.Type != job.OffchainReporting2 {
		return jb, errors.Errorf("the only supported type is currently 'offchainreporting2', got %s", jb.Type)
	}
	if !tree.Has("isBootstrapPeer") {
		return jb, errors.New("isBootstrapPeer is not defined")
	}
	for i := range spec.P2PBootstrapPeers {
		if _, err := multiaddr.NewMultiaddr(spec.P2PBootstrapPeers[i]); err != nil {
			return jb, errors.Wrapf(err, "p2p bootstrap peer %v is invalid", spec.P2PBootstrapPeers[i])
		}
	}

	chain, err := chainSet.Get(jb.OffchainreportingOracleSpec.EVMChainID.ToInt())
	if err != nil {
		return jb, err
	}
	panic(chain)

	if spec.IsBootstrapPeer {
		if err := validateBootstrapSpec(tree, jb); err != nil {
			return jb, err
		}
	} else if err := validateNonBootstrapSpec(tree, chain.Config(), jb); err != nil {
		return jb, err
	}
	if err := validateTimingParameters(chain.Config(), spec); err != nil {
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

func validateTimingParameters(config ValidationConfig, spec job.OffchainReporting2OracleSpec) error {
	lc := computeLocalConfig(config, spec)
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

	var observationTimeout time.Duration
	if spec.Offchainreporting2OracleSpec.ObservationTimeout != 0 {
		observationTimeout = time.Duration(spec.Offchainreporting2OracleSpec.ObservationTimeout)
	} else {
		observationTimeout = config.OCR2ObservationTimeout()
	}
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
