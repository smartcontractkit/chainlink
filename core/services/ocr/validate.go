package ocr

import (
	"math/big"
	"time"

	"github.com/lib/pq"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting"

	evmconfig "github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/chaintype"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

type GeneralConfig interface {
	OCR() coreconfig.OCR
	Insecure() coreconfig.Insecure
}

type ValidationConfig interface {
	ChainType() chaintype.ChainType
}

type OCRValidationConfig interface {
	BlockchainTimeout() time.Duration
	CaptureEATelemetry() bool
	ContractPollInterval() time.Duration
	ContractSubscribeInterval() time.Duration
	KeyBundleID() (string, error)
	ObservationTimeout() time.Duration
	TransmitterAddress() (types.EIP55Address, error)
}

type insecureConfig interface {
	OCRDevelopmentMode() bool
}

// ValidatedOracleSpecToml validates an oracle spec that came from TOML
func ValidatedOracleSpecToml(gcfg GeneralConfig, legacyChains legacyevm.LegacyChainContainer, tomlString string) (job.Job, error) {
	return ValidatedOracleSpecTomlCfg(gcfg, func(id *big.Int) (evmconfig.ChainScopedConfig, error) {
		c, err := legacyChains.Get(id.String())
		if err != nil {
			return nil, err
		}
		return c.Config(), nil
	}, tomlString)
}

func ValidatedOracleSpecTomlCfg(gcfg GeneralConfig, configFn func(id *big.Int) (evmconfig.ChainScopedConfig, error), tomlString string) (job.Job, error) {
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
		if err := validateBootstrapSpec(tree); err != nil {
			return jb, err
		}
	} else if err := validateNonBootstrapSpec(tree, jb, gcfg.OCR().ObservationTimeout()); err != nil {
		return jb, err
	}
	if err := validateTimingParameters(cfg.EVM(), cfg.EVM().OCR(), gcfg.Insecure(), spec, gcfg.OCR()); err != nil {
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

func validateTimingParameters(cfg ValidationConfig, evmOcrCfg evmconfig.OCR, insecureCfg insecureConfig, spec job.OCROracleSpec, ocrCfg job.OCRConfig) error {
	lc := toLocalConfig(cfg, evmOcrCfg, insecureCfg, spec, ocrCfg)
	return errors.Wrap(offchainreporting.SanityCheckLocalConfig(lc), "offchainreporting.SanityCheckLocalConfig failed")
}

func validateBootstrapSpec(tree *toml.Tree) error {
	expected, notExpected := ocrcommon.CloneSet(params), ocrcommon.CloneSet(nonBootstrapParams)
	for k := range bootstrapParams {
		expected[k] = struct{}{}
	}
	return ocrcommon.ValidateExplicitlySetKeys(tree, expected, notExpected, "bootstrap")
}

func validateNonBootstrapSpec(tree *toml.Tree, spec job.Job, ocrObservationTimeout time.Duration) error {
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
		observationTimeout = ocrObservationTimeout
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
