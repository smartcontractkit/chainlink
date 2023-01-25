package validate

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"

	"github.com/smartcontractkit/chainlink/core/services/job"
	dkgconfig "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/dkg/config"
	ocr2vrfconfig "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/config"
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
	if jb.OCR2OracleSpec.P2PV2Bootstrappers == nil {
		// Empty but non-null, field is non-nullable.
		jb.OCR2OracleSpec.P2PV2Bootstrappers = pq.StringArray{}
	}

	if jb.Type != job.OffchainReporting2 {
		return jb, errors.Errorf("the only supported type is currently 'offchainreporting2', got %s", jb.Type)
	}
	if _, ok := relay.SupportedRelays[spec.Relay]; !ok {
		return jb, errors.Errorf("no such relay %v supported", spec.Relay)
	}
	if len(spec.P2PV2Bootstrappers) > 0 {
		_, err = ocrcommon.ParseBootstrapPeers(spec.P2PV2Bootstrappers)
		if err != nil {
			return jb, err
		}
	}

	if err = validateSpec(tree, jb); err != nil {
		return jb, err
	}
	if err = validateTimingParameters(config, spec); err != nil {
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
	case job.DKG:
		return validateDKGSpec(spec.OCR2OracleSpec.PluginConfig)
	case job.OCR2VRF:
		return validateOCR2VRFSpec(spec.OCR2OracleSpec.PluginConfig)
	case job.OCR2Keeper:
		return validateOCR2KeeperSpec(spec.OCR2OracleSpec.PluginConfig)
	case job.OCR2Functions:
		// TODO validator for DR-OCR spec: https://app.shortcut.com/chainlinklabs/story/54054/ocr-plugin-for-directrequest-ocr
		return nil
	case "":
		return errors.New("no plugin specified")
	default:
		return errors.Errorf("invalid pluginType %s", spec.OCR2OracleSpec.PluginType)
	}

	return nil
}

func validateDKGSpec(jsonConfig job.JSONConfig) error {
	if jsonConfig == nil {
		return errors.New("pluginConfig is empty")
	}
	var pluginConfig dkgconfig.PluginConfig
	err := json.Unmarshal(jsonConfig.Bytes(), &pluginConfig)
	if err != nil {
		return errors.Wrap(err, "error while unmarshaling plugin config")
	}
	err = validateHexString(pluginConfig.EncryptionPublicKey, 32)
	if err != nil {
		return errors.Wrap(err, "validation error for encryptedPublicKey")
	}
	err = validateHexString(pluginConfig.SigningPublicKey, 32)
	if err != nil {
		return errors.Wrap(err, "validation error for signingPublicKey")
	}
	err = validateHexString(pluginConfig.KeyID, 32)
	if err != nil {
		return errors.Wrap(err, "validation error for keyID")
	}

	return nil
}

func validateHexString(val string, expectedLengthInBytes uint) error {
	decoded, err := hex.DecodeString(val)
	if err != nil {
		return errors.Wrapf(err, "expected hex string but received %s", val)
	}
	if len(decoded) != int(expectedLengthInBytes) {
		return fmt.Errorf("value: %s has unexpected length. Expected %d bytes", val, expectedLengthInBytes)
	}
	return nil
}

func validateOCR2VRFSpec(jsonConfig job.JSONConfig) error {
	if jsonConfig == nil {
		return errors.New("pluginConfig is empty")
	}
	var cfg ocr2vrfconfig.PluginConfig
	err := json.Unmarshal(jsonConfig.Bytes(), &cfg)
	if err != nil {
		return errors.Wrap(err, "json unmarshal plugin config")
	}
	err = validateDKGSpec(job.JSONConfig{
		"encryptionPublicKey": cfg.DKGEncryptionPublicKey,
		"signingPublicKey":    cfg.DKGSigningPublicKey,
		"keyID":               cfg.DKGKeyID,
	})
	if err != nil {
		return err
	}
	if cfg.LinkEthFeedAddress == "" {
		return errors.New("linkEthFeedAddress must be provided")
	}
	if cfg.DKGContractAddress == "" {
		return errors.New("dkgContractAddress must be provided")
	}
	return nil
}

func validateOCR2KeeperSpec(jsonConfig job.JSONConfig) error {
	return nil
}
