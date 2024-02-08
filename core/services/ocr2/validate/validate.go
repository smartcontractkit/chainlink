package validate

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/pelletier/go-toml"
	pkgerrors "github.com/pkg/errors"

	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	dkgconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/dkg/config"
	mercuryconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	ocr2vrfconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2vrf/config"
	rebalancermodels "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// ValidatedOracleSpecToml validates an oracle spec that came from TOML
func ValidatedOracleSpecToml(config OCR2Config, insConf InsecureConfig, tomlString string) (job.Job, error) {
	var jb = job.Job{}
	var spec job.OCR2OracleSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, pkgerrors.Wrapf(err, "toml error on load %v", tomlString)
	}
	// Note this validates all the fields which implement an UnmarshalText
	// i.e. TransmitterAddress, PeerID...
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, pkgerrors.Wrap(err, "toml unmarshal error on spec")
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, pkgerrors.Wrap(err, "toml unmarshal error on job")
	}
	jb.OCR2OracleSpec = &spec
	if jb.OCR2OracleSpec.P2PV2Bootstrappers == nil {
		// Empty but non-null, field is non-nullable.
		jb.OCR2OracleSpec.P2PV2Bootstrappers = pq.StringArray{}
	}

	if jb.Type != job.OffchainReporting2 {
		return jb, pkgerrors.Errorf("the only supported type is currently 'offchainreporting2', got %s", jb.Type)
	}
	if _, ok := relay.SupportedRelays[spec.Relay]; !ok {
		return jb, pkgerrors.Errorf("no such relay %v supported", spec.Relay)
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
	if err = validateTimingParameters(config, insConf, spec); err != nil {
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

func validateTimingParameters(ocr2Conf OCR2Config, insConf InsecureConfig, spec job.OCR2OracleSpec) error {
	lc, err := ToLocalConfig(ocr2Conf, insConf, spec)
	if err != nil {
		return err
	}
	return libocr2.SanityCheckLocalConfig(lc)
}

func validateSpec(tree *toml.Tree, spec job.Job) error {
	expected, notExpected := ocrcommon.CloneSet(params), ocrcommon.CloneSet(notExpectedParams)
	if err := ocrcommon.ValidateExplicitlySetKeys(tree, expected, notExpected, "ocr2"); err != nil {
		return err
	}

	switch spec.OCR2OracleSpec.PluginType {
	case types.Median:
		if spec.Pipeline.Source == "" {
			return errors.New("no pipeline specified")
		}
	case types.DKG:
		return validateDKGSpec(spec.OCR2OracleSpec.PluginConfig)
	case types.OCR2VRF:
		return validateOCR2VRFSpec(spec.OCR2OracleSpec.PluginConfig)
	case types.OCR2Keeper:
		return validateOCR2KeeperSpec(spec.OCR2OracleSpec.PluginConfig)
	case types.Functions:
		// TODO validator for DR-OCR spec: https://app.shortcut.com/chainlinklabs/story/54054/ocr-plugin-for-directrequest-ocr
		return nil
	case types.Mercury:
		return validateOCR2MercurySpec(spec.OCR2OracleSpec.PluginConfig, *spec.OCR2OracleSpec.FeedID)
	case types.CCIPExecution:
		return validateOCR2CCIPExecutionSpec(spec.OCR2OracleSpec.PluginConfig)
	case types.CCIPCommit:
		return validateOCR2CCIPCommitSpec(spec.OCR2OracleSpec.PluginConfig)
	case types.GenericPlugin:
		return validateOCR2GenericPluginSpec(spec.OCR2OracleSpec.PluginConfig)
	case "rebalancer":
		return validateRebalancerSpec(spec.OCR2OracleSpec.PluginConfig)
	case "":
		return errors.New("no plugin specified")
	default:
		return pkgerrors.Errorf("invalid pluginType %s", spec.OCR2OracleSpec.PluginType)
	}

	return nil
}

type PipelineSpec struct {
	Name string `json:"name"`
	Spec string `json:"spec"`
}

type Config struct {
	Pipelines    []PipelineSpec `json:"pipelines"`
	PluginConfig map[string]any `json:"pluginConfig"`
}

type innerConfig struct {
	Command       string            `json:"command"`
	EnvVars       map[string]string `json:"envVars"`
	ProviderType  string            `json:"providerType"`
	PluginName    string            `json:"pluginName"`
	TelemetryType string            `json:"telemetryType"`
	Config
}

type OCR2GenericPluginConfig struct {
	innerConfig
}

func (o *OCR2GenericPluginConfig) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &o.innerConfig)
	if err != nil {
		return nil
	}

	m := map[string]any{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	o.PluginConfig = m
	return nil
}

func validateOCR2GenericPluginSpec(jsonConfig job.JSONConfig) error {
	p := OCR2GenericPluginConfig{}
	err := json.Unmarshal(jsonConfig.Bytes(), &p)
	if err != nil {
		return err
	}

	if p.PluginName == "" {
		return errors.New("generic config invalid: must provide plugin name")
	}

	if p.TelemetryType == "" {
		return errors.New("generic config invalid: must provide telemetry type")
	}

	return nil
}

func validateRebalancerSpec(jsonConfig job.JSONConfig) error {
	if jsonConfig == nil {
		return errors.New("pluginConfig is empty")
	}
	var pluginConfig rebalancermodels.PluginConfig
	err := json.Unmarshal(jsonConfig.Bytes(), &pluginConfig)
	if err != nil {
		return pkgerrors.Wrap(err, "error while unmarshalling plugin config")
	}
	if pluginConfig.LiquidityManagerNetwork == 0 {
		return errors.New("liquidityManagerNetwork must be provided")
	}
	if pluginConfig.ClosePluginTimeoutSec <= 0 {
		return errors.New("closePluginTimeoutSec must be positive")
	}
	if err := rebalancermodels.ValidateRebalancerConfig(pluginConfig.RebalancerConfig); err != nil {
		return fmt.Errorf("rebalancer config invalid: %w", err)
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
		return pkgerrors.Wrap(err, "error while unmarshalling plugin config")
	}
	err = validateHexString(pluginConfig.EncryptionPublicKey, 32)
	if err != nil {
		return pkgerrors.Wrap(err, "validation error for encryptedPublicKey")
	}
	err = validateHexString(pluginConfig.SigningPublicKey, 32)
	if err != nil {
		return pkgerrors.Wrap(err, "validation error for signingPublicKey")
	}
	err = validateHexString(pluginConfig.KeyID, 32)
	if err != nil {
		return pkgerrors.Wrap(err, "validation error for keyID")
	}

	return nil
}

func validateHexString(val string, expectedLengthInBytes uint) error {
	decoded, err := hex.DecodeString(val)
	if err != nil {
		return pkgerrors.Wrapf(err, "expected hex string but received %s", val)
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
		return pkgerrors.Wrap(err, "json unmarshal plugin config")
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

func validateOCR2MercurySpec(jsonConfig job.JSONConfig, feedId [32]byte) error {
	var pluginConfig mercuryconfig.PluginConfig
	err := json.Unmarshal(jsonConfig.Bytes(), &pluginConfig)
	if err != nil {
		return pkgerrors.Wrap(err, "error while unmarshalling plugin config")
	}
	return pkgerrors.Wrap(mercuryconfig.ValidatePluginConfig(pluginConfig, feedId), "Mercury PluginConfig is invalid")
}

func validateOCR2CCIPExecutionSpec(jsonConfig job.JSONConfig) error {
	if jsonConfig == nil {
		return errors.New("pluginConfig is empty")
	}
	var cfg config.ExecutionPluginJobSpecConfig
	err := json.Unmarshal(jsonConfig.Bytes(), &cfg)
	if err != nil {
		return pkgerrors.Wrap(err, "error while unmarshalling plugin config")
	}
	if cfg.USDCConfig != (config.USDCConfig{}) {
		return cfg.USDCConfig.ValidateUSDCConfig()
	}
	return nil
}

func validateOCR2CCIPCommitSpec(jsonConfig job.JSONConfig) error {
	if jsonConfig == nil {
		return errors.New("pluginConfig is empty")
	}
	var cfg config.CommitPluginJobSpecConfig
	err := json.Unmarshal(jsonConfig.Bytes(), &cfg)
	if err != nil {
		return pkgerrors.Wrap(err, "error while unmarshalling plugin config")
	}

	// Ensure that either the tokenPricesUSDPipeline or the priceGetterConfig is set, but not both.
	emptyPipeline := strings.Trim(cfg.TokenPricesUSDPipeline, "\n\t ") == ""
	emptyPriceGetter := cfg.PriceGetterConfig == nil
	if emptyPipeline && emptyPriceGetter {
		return fmt.Errorf("either tokenPricesUSDPipeline or priceGetterConfig must be set")
	}
	if !emptyPipeline && !emptyPriceGetter {
		return fmt.Errorf("only one of tokenPricesUSDPipeline or priceGetterConfig must be set: %s and %v", cfg.TokenPricesUSDPipeline, cfg.PriceGetterConfig)
	}

	if !emptyPipeline {
		_, err = pipeline.Parse(cfg.TokenPricesUSDPipeline)
		if err != nil {
			return pkgerrors.Wrap(err, "invalid token prices pipeline")
		}
	} else {
		// Validate prices config (like it was done for the pipeline).
		if emptyPriceGetter {
			return pkgerrors.New("priceGetterConfig is empty")
		}
	}

	return nil
}
