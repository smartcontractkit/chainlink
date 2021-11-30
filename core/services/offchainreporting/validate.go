package offchainreporting

import (
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/libocr/offchainreporting"
)

type ValidationConfig interface {
	ChainType() chains.ChainType
	Dev() bool
	OCRBlockchainTimeout() time.Duration
	OCRContractConfirmations() uint16
	OCRContractPollInterval() time.Duration
	OCRContractSubscribeInterval() time.Duration
	OCRContractTransmitterTransmitTimeout() time.Duration
	OCRDatabaseTimeout() time.Duration
	OCRKeyBundleID() (string, error)
	OCRObservationGracePeriod() time.Duration
	OCRObservationTimeout() time.Duration
	OCRTransmitterAddress() (ethkey.EIP55Address, error)
	P2PPeerID() p2pkey.PeerID
}

// ValidatedOracleSpecToml validates an oracle spec that came from TOML
func ValidatedOracleSpecToml(chainSet evm.ChainSet, tomlString string) (job.Job, error) {
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

	if jb.Type != job.OffchainReporting {
		return jb, errors.Errorf("the only supported type is currently 'offchainreporting', got %s", jb.Type)
	}
	if !tree.Has("isBootstrapPeer") {
		return jb, errors.New("isBootstrapPeer is not defined")
	}
	for i := range spec.P2PBootstrapPeers {
		if _, err = multiaddr.NewMultiaddr(spec.P2PBootstrapPeers[i]); err != nil {
			return jb, errors.Wrapf(err, "p2p bootstrap peer %v is invalid", spec.P2PBootstrapPeers[i])
		}
	}

	chain, err := chainSet.Get(jb.OffchainreportingOracleSpec.EVMChainID.ToInt())
	if err != nil {
		return jb, err
	}

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

func validateTimingParameters(cfg ValidationConfig, spec job.OffchainReportingOracleSpec) error {
	lc := toLocalConfig(cfg, spec)
	return errors.Wrap(offchainreporting.SanityCheckLocalConfig(lc), "offchainreporting.SanityCheckLocalConfig failed")
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
	if spec.OffchainreportingOracleSpec.ObservationTimeout != 0 {
		observationTimeout = spec.OffchainreportingOracleSpec.ObservationTimeout.Duration()
	} else {
		observationTimeout = config.OCRObservationTimeout()
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
