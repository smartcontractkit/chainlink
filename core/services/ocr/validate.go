package ocr

import (
	"math/big"
	"time"

	"github.com/lib/pq"
	"github.com/multiformats/go-multiaddr"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	config2 "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
)

type ValidationConfig interface {
	ChainType() config.ChainType
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
	return ValidatedOracleSpecTomlCfg(func(id *big.Int) (config2.ChainScopedConfig, error) {
		c, err := chainSet.Get(id)
		if err != nil {
			return nil, err
		}
		return c.Config(), nil
	}, tomlString)
}

func ValidatedOracleSpecTomlCfg(configFn func(id *big.Int) (config2.ChainScopedConfig, error), tomlString string) (job.Job, error) {
	var jb = job.Job{}
	var spec job.OCROracleSpec
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
	jb.OCROracleSpec = &spec

	if jb.OCROracleSpec.P2PV2Bootstrappers == nil {
		// Empty but non-null, field is non-nullable.
		jb.OCROracleSpec.P2PV2Bootstrappers = pq.StringArray{}
	}

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

	if len(spec.P2PV2Bootstrappers) > 0 {
		_, err = ocrcommon.ParseBootstrapPeers(spec.P2PV2Bootstrappers)
		if err != nil {
			return jb, err
		}
	}

	cfg, err := configFn(jb.OCROracleSpec.EVMChainID.ToInt())
	if err != nil {
		return jb, err
	}

	if spec.IsBootstrapPeer {
		if err := validateBootstrapSpec(tree, jb); err != nil {
			return jb, err
		}
	} else if err := validateNonBootstrapSpec(tree, cfg, jb); err != nil {
		return jb, err
	}
	if err := validateTimingParameters(cfg, spec); err != nil {
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

func validateTimingParameters(cfg ValidationConfig, spec job.OCROracleSpec) error {
	lc := toLocalConfig(cfg, spec)
	return errors.Wrap(offchainreporting.SanityCheckLocalConfig(lc), "offchainreporting.SanityCheckLocalConfig failed")
}

func validateBootstrapSpec(tree *toml.Tree, spec job.Job) error {
	expected, notExpected := ocrcommon.CloneSet(params), ocrcommon.CloneSet(nonBootstrapParams)
	for k := range bootstrapParams {
		expected[k] = struct{}{}
	}
	return ocrcommon.ValidateExplicitlySetKeys(tree, expected, notExpected, "bootstrap")
}

func validateNonBootstrapSpec(tree *toml.Tree, config ValidationConfig, spec job.Job) error {
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
	var observationTimeout time.Duration
	if spec.OCROracleSpec.ObservationTimeout != 0 {
		observationTimeout = spec.OCROracleSpec.ObservationTimeout.Duration()
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
