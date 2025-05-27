package chainlink

import (
	"errors"
	"fmt"
	"slices"

	"github.com/imdario/mergo"
	"go.uber.org/multierr"

	gotoml "github.com/pelletier/go-toml/v2"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	configtoml "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config/docs"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"
)

// Config is the root type used for TOML configuration.
//
// See docs at /docs/CONFIG.md generated via config.GenerateDocs from /internal/config/docs.toml
//
// When adding a new field:
//   - consider including a unit suffix with the field name
//   - TOML is limited to int64/float64, so fields requiring greater range/precision must use non-standard types
//     implementing encoding.TextMarshaler/TextUnmarshaler, like big.Big and decimal.Decimal
//   - std lib types that don't implement encoding.TextMarshaler/TextUnmarshaler (time.Duration, url.URL, big.Int) won't
//     work as expected, and require wrapper types. See commonconfig.Duration, commonconfig.URL, big.Big.
type Config struct {
	toml.Core

	EVM configtoml.EVMConfigs `toml:",omitempty"`

	Cosmos RawConfigs `toml:",omitempty"`

	Solana solcfg.TOMLConfigs `toml:",omitempty"`

	Starknet RawConfigs `toml:",omitempty"`

	Aptos RawConfigs `toml:",omitempty"`

	Tron RawConfigs `toml:",omitempty"`

	TON RawConfigs `toml:",omitempty"`
}

// RawConfigs is a list of RawConfig.
type RawConfigs []RawConfig

func (rs RawConfigs) SetDefaults() {
	for _, r := range rs {
		r.SetDefaults()
	}
}

func (rs *RawConfigs) SetFrom(configs RawConfigs) error {
	if err := configs.validateKeys(); err != nil {
		return err
	}

	for _, config := range configs {
		chainID := config.ChainID()
		i := slices.IndexFunc(*rs, func(r RawConfig) bool {
			otherChainID := r.ChainID()
			return otherChainID != "" && chainID == otherChainID
		})
		if i != -1 {
			if err := (*rs)[i].SetFrom(config); err != nil {
				return err
			}
		} else {
			*rs = append(*rs, config)
		}
	}

	return nil
}

func (rs RawConfigs) validateKeys() (err error) {
	chainIDs := commonconfig.UniqueStrings{}
	for i, config := range rs {
		chainID := config.ChainID()
		if chainIDs.IsDupe(&chainID) {
			err = errors.Join(err, commonconfig.NewErrDuplicate(fmt.Sprintf("%d.ChainID", i), chainID))
		}
	}

	nodeNames := commonconfig.UniqueStrings{}
	for i, config := range rs {
		configNodeNames := config.NodeNames()
		for j, nodeName := range configNodeNames {
			if nodeNames.IsDupe(&nodeName) {
				err = errors.Join(err, commonconfig.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.Name", i, j), nodeName))
			}
		}
	}
	return
}

func (rs RawConfigs) ValidateConfig() (err error) {
	return rs.validateKeys()
}

// RawConfig is the config used for chains that are not embedded.
type RawConfig map[string]any

type parsedRawConfig struct {
	chainID    string
	nodesExist bool
	nodes      []map[string]any
	nodeNames  []string
}

func (c RawConfig) parse() (*parsedRawConfig, error) {
	var err error
	if v, ok := c["Enabled"]; ok {
		if _, ok := v.(bool); !ok {
			err = multierr.Append(err, commonconfig.ErrInvalid{Name: "Enabled", Value: v, Msg: "expected bool"})
		}
	}

	parsedRawConfig := &parsedRawConfig{}
	chainID, exists := c["ChainID"]
	if !exists {
		err = multierr.Append(err, commonconfig.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else {
		chainIDStr, ok := chainID.(string)
		switch {
		case !ok:
			err = multierr.Append(err, commonconfig.ErrInvalid{Name: "ChainID", Value: chainID, Msg: "expected string"})
		case chainIDStr == "":
			err = multierr.Append(err, commonconfig.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
		default:
			parsedRawConfig.chainID = chainIDStr
		}
	}
	nodes, nodesExist := c["Nodes"]
	parsedRawConfig.nodesExist = nodesExist
	if nodesExist {
		nodeMaps, ok := nodes.([]any)
		switch {
		case !ok:
			err = multierr.Append(err, commonconfig.ErrInvalid{Name: "Nodes", Value: nodes, Msg: "expected array of node configs"})
		default:
			for i, node := range nodeMaps {
				nodeConfig, ok := node.(map[string]any)
				if !ok {
					err = multierr.Append(err, commonconfig.ErrInvalid{Name: fmt.Sprintf("Nodes.%d", i), Value: nodeConfig, Msg: "expected node config map"})
				} else {
					parsedRawConfig.nodes = append(parsedRawConfig.nodes, nodeConfig)
					nodeName, exists := nodeConfig["Name"]
					if !exists {
						err = multierr.Append(err, commonconfig.ErrMissing{Name: fmt.Sprintf("Nodes.%d.Name", i), Msg: "required for all nodes"})
					} else {
						nodeNameStr, ok := nodeName.(string)
						switch {
						case !ok:
							err = multierr.Append(err, commonconfig.ErrInvalid{Name: fmt.Sprintf("Nodes.%d.Name", i), Value: nodeName, Msg: "expected string"})
						case nodeNameStr == "":
							err = multierr.Append(err, commonconfig.ErrEmpty{Name: fmt.Sprintf("Nodes.%d.Name", i), Msg: "required for all nodes"})
						default:
							parsedRawConfig.nodeNames = append(parsedRawConfig.nodeNames, nodeNameStr)
						}
					}
				}
			}
		}
	}

	return parsedRawConfig, err
}

// ValidateConfig returns an error if the Config is not valid for use, as-is.
func (c RawConfig) ValidateConfig() error {
	parsedRawConfig, err := c.parse()
	if !parsedRawConfig.nodesExist {
		err = multierr.Append(err, commonconfig.ErrMissing{Name: "Nodes", Msg: "expected at least one node"})
	} else if len(parsedRawConfig.nodes) == 0 {
		err = multierr.Append(err, commonconfig.ErrEmpty{Name: "Nodes", Msg: "expected at least one node"})
	}
	return err
}

func (c RawConfig) IsEnabled() bool {
	s := c["Enabled"]
	if s == nil {
		return true // default to true if omitted
	}
	b, ok := s.(bool)
	return ok && b
}

func (c RawConfig) ChainID() string {
	chainID, _ := c["ChainID"].(string)
	return chainID
}

func (c *RawConfig) SetFrom(config RawConfig) error {
	if e := config["Enabled"]; e != nil {
		(*c)["Enabled"] = e
	}
	parsedRawConfig, err := c.parse()
	if err != nil {
		return err
	}

	incomingParsedRawConfig, err := config.parse()
	if err != nil {
		return err
	}

	// Create a copy of config without nodes to merge other fields
	configWithoutNodes := make(RawConfig)
	for k, v := range config {
		if k != "Nodes" {
			configWithoutNodes[k] = v
		}
	}

	// Merge all non-node fields
	if err := mergo.Merge(c, configWithoutNodes, mergo.WithOverride); err != nil {
		return err
	}

	// Handle node merging
	for i, nodeConfig := range incomingParsedRawConfig.nodes {
		nodeName := incomingParsedRawConfig.nodeNames[i]
		i := slices.Index(parsedRawConfig.nodeNames, nodeName)
		if i != -1 {
			if err := mergo.Merge(&parsedRawConfig.nodes[i], nodeConfig, mergo.WithOverride); err != nil {
				return err
			}
		} else {
			parsedRawConfig.nodes = append(parsedRawConfig.nodes, nodeConfig)
		}
	}

	// Subsequence SetFrom invocations will call parse(), and expect to be able to cast c["Nodes"] to []any,
	// so we can't directly assign parsedRawConfig.nodes back to c["Nodes"].
	anyConfigs := []any{}
	for _, nodeConfig := range parsedRawConfig.nodes {
		anyConfigs = append(anyConfigs, nodeConfig)
	}

	(*c)["Nodes"] = anyConfigs
	return nil
}

func (c RawConfig) NodeNames() []string {
	nodes, _ := c["Nodes"].([]any)
	nodeNames := []string{}
	for _, node := range nodes {
		config, _ := node.(map[string]any)
		nodeName, _ := config["Name"].(string)
		nodeNames = append(nodeNames, nodeName)
	}
	return nodeNames
}

func (c RawConfig) SetDefaults() {
	if e, ok := c["Enabled"].(bool); ok && e {
		// already enabled by default so drop it
		delete(c, "Enabled")
	}
}

// TOMLString returns a TOML encoded string.
func (c *Config) TOMLString() (string, error) {
	b, err := gotoml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// warnings aggregates warnings from valueWarnings and deprecationWarnings
func (c *Config) warnings() (err error) {
	deprecationErr := c.deprecationWarnings()

	warningErr := c.valueWarnings()
	err = multierr.Append(deprecationErr, warningErr)
	_, list := commonconfig.MultiErrorList(err)
	return list
}

// valueWarnings returns an error if the Config contains values that hint at misconfiguration before defaults are applied.
func (c *Config) valueWarnings() (err error) {
	if c.Tracing.Enabled != nil && *c.Tracing.Enabled {
		if c.Tracing.Mode != nil && *c.Tracing.Mode == "unencrypted" {
			if c.Tracing.TLSCertPath != nil {
				err = multierr.Append(err, config.ErrInvalid{Name: "Tracing.TLSCertPath", Value: *c.Tracing.TLSCertPath, Msg: "must be empty when Tracing.Mode is 'unencrypted'"})
			}
		}
	}
	return
}

// deprecationWarnings returns an error if the Config contains deprecated fields.
// This is typically used before defaults have been applied, with input from the user.
func (c *Config) deprecationWarnings() (err error) {
	return nil
}

// Validate returns an error if the Config is not valid for use, as-is.
// This is typically used after defaults have been applied.
func (c *Config) Validate() error {
	if err := commonconfig.Validate(c); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	return nil
}

// setDefaults initializes unset fields with default values.
func (c *Config) setDefaults() {
	core := docs.CoreDefaults()
	core.SetFrom(&c.Core)
	c.Core = core

	c.Aptos.SetDefaults()

	for i := range c.EVM {
		if input := c.EVM[i]; input == nil {
			c.EVM[i] = &configtoml.EVMConfig{Chain: configtoml.Defaults(nil)}
		} else {
			input.Chain = configtoml.Defaults(input.ChainID, &input.Chain)
		}
	}

	c.Cosmos.SetDefaults()

	for i := range c.Solana {
		if c.Solana[i] == nil {
			c.Solana[i] = new(solcfg.TOMLConfig)
		}
		c.Solana[i].Chain.SetDefaults()
	}

	c.Starknet.SetDefaults()

	c.Tron.SetDefaults()

	c.TON.SetDefaults()
}

func (c *Config) SetFrom(f *Config) (err error) {
	c.Core.SetFrom(&f.Core)

	appendErr := func(e error, name string) {
		if e != nil {
			err = multierr.Append(err, commonconfig.NamedMultiErrorList(e, name))
		}
	}

	appendErr(c.EVM.SetFrom(&f.EVM), "EVM")
	appendErr(c.Cosmos.SetFrom(f.Cosmos), "Cosmos")
	appendErr(c.Solana.SetFrom(&f.Solana), "Solana")
	appendErr(c.Starknet.SetFrom(f.Starknet), "Starknet")
	appendErr(c.Aptos.SetFrom(f.Aptos), "Aptos")
	appendErr(c.Tron.SetFrom(f.Tron), "Tron")
	appendErr(c.TON.SetFrom(f.TON), "TON")

	_, err = commonconfig.MultiErrorList(err)
	return err
}

type Secrets struct {
	toml.Secrets
}

func (s *Secrets) SetFrom(f *Secrets) (err error) {
	if err2 := s.Database.SetFrom(&f.Database); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "Database"))
	}

	if err2 := s.Password.SetFrom(&f.Password); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "Password"))
	}

	if err2 := s.WebServer.SetFrom(&f.WebServer); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "WebServer"))
	}

	if err2 := s.Pyroscope.SetFrom(&f.Pyroscope); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "Pyroscope"))
	}

	if err2 := s.Prometheus.SetFrom(&f.Prometheus); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "Prometheus"))
	}

	if err2 := s.Mercury.SetFrom(&f.Mercury); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "Mercury"))
	}

	if err2 := s.Threshold.SetFrom(&f.Threshold); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "Threshold"))
	}

	if err2 := s.EVM.SetFrom(&f.EVM); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "EthKeys"))
	}

	if err2 := s.P2PKey.SetFrom(&f.P2PKey); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "P2PKey"))
	}

	if err2 := s.CRE.SetFrom(&f.CRE); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "CRE"))
	}

	_, err = commonconfig.MultiErrorList(err)

	return err
}

func (s *Secrets) setDefaults() {
	if nil == s.Database.AllowSimplePasswords {
		s.Database.AllowSimplePasswords = new(bool)
	}
}

// TOMLString returns a TOML encoded string with secret values redacted.
func (s *Secrets) TOMLString() (string, error) {
	b, err := gotoml.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type InvalidSecretsError struct {
	err error
}

func (e InvalidSecretsError) Error() string {
	return fmt.Sprintf("invalid secrets: %v", e.err)
}

func (e InvalidSecretsError) Unwrap() error {
	return e.err
}

func (e InvalidSecretsError) Is(err error) bool {
	_, ok := err.(InvalidSecretsError) //nolint:errcheck // implementing errors.Is
	return ok
}

// Validate validates every consitutent secret and return an accumulated error
func (s *Secrets) Validate() error {
	if err := commonconfig.Validate(s); err != nil {
		return InvalidSecretsError{err: err}
	}
	return nil
}

// ValidateDB only validates the encompassed DatabaseSecret
func (s *Secrets) ValidateDB() error {
	// This implementation was chosen so that error reporting is uniform
	// when validating all the secret or only the db secrets,
	// and so we could reuse config.Validate, which contains fearsome reflection logic.
	// This meets the current needs, but if we ever wanted to compose secret
	// validation we may need to rethink this approach and instead find a way to
	// toggle on/off the validation of the embedded secrets.

	type dbValidationType struct {
		// choose field name to match that of Secrets.Database so we have
		// consistent error messages.
		Database toml.DatabaseSecrets
	}
	s.setDefaults()
	v := &dbValidationType{s.Database}
	if err := commonconfig.Validate(v); err != nil {
		return InvalidSecretsError{err: err}
	}
	return nil
}

// setEnv overrides fields from ENV vars, if present.
func (s *Secrets) setEnv() error {
	if dbURL := env.DatabaseURL.Get(); dbURL != "" {
		s.Database.URL = new(models.SecretURL)
		if err := s.Database.URL.UnmarshalText([]byte(dbURL)); err != nil {
			return err
		}
	}
	if dbBackupUrl := env.DatabaseBackupURL.Get(); dbBackupUrl != "" {
		s.Database.BackupURL = new(models.SecretURL)
		if err := s.Database.BackupURL.UnmarshalText([]byte(dbBackupUrl)); err != nil {
			return err
		}
	}
	if env.DatabaseAllowSimplePasswords.IsTrue() {
		s.Database.AllowSimplePasswords = new(bool)
		*s.Database.AllowSimplePasswords = true
	}
	if keystorePassword := env.PasswordKeystore.Get(); keystorePassword != "" {
		s.Password.Keystore = &keystorePassword
	}
	if vrfPassword := env.PasswordVRF.Get(); vrfPassword != "" {
		s.Password.VRF = &vrfPassword
	}
	if pyroscopeAuthToken := env.PyroscopeAuthToken.Get(); pyroscopeAuthToken != "" {
		s.Pyroscope.AuthToken = &pyroscopeAuthToken
	}
	if prometheusAuthToken := env.PrometheusAuthToken.Get(); prometheusAuthToken != "" {
		s.Prometheus.AuthToken = &prometheusAuthToken
	}
	if thresholdKeyShare := env.ThresholdKeyShare.Get(); thresholdKeyShare != "" {
		s.Threshold.ThresholdKeyShare = &thresholdKeyShare
	}
	return nil
}
