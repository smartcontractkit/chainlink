package chainlink

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"

	gotoml "github.com/pelletier/go-toml/v2"

	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
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

	EVM evmcfg.EVMConfigs `toml:",omitempty"`

	Cosmos coscfg.TOMLConfigs `toml:",omitempty"`

	Solana solcfg.TOMLConfigs `toml:",omitempty"`

	Starknet stkcfg.TOMLConfigs `toml:",omitempty"`

	Aptos RawConfigs `toml:",omitempty"`
}

// RawConfigs is a list of RawConfig.
type RawConfigs []RawConfig

// RawConfig is the config used for chains that are not embedded.
type RawConfig map[string]any

// ValidateConfig returns an error if the Config is not valid for use, as-is.
func (c *RawConfig) ValidateConfig() (err error) {
	if v, ok := (*c)["Enabled"]; ok {
		if _, ok := v.(bool); !ok {
			err = multierr.Append(err, commonconfig.ErrInvalid{Name: "Enabled", Value: v, Msg: "expected bool"})
		}
	}
	if v, ok := (*c)["ChainID"]; ok {
		if _, ok := v.(string); !ok {
			err = multierr.Append(err, commonconfig.ErrInvalid{Name: "ChainID", Value: v, Msg: "expected string"})
		}
	}
	return err
}

func (c *RawConfig) IsEnabled() bool {
	if c == nil {
		return false
	}

	enabled, ok := (*c)["Enabled"].(bool)
	return ok && enabled
}

func (c *RawConfig) ChainID() string {
	if c == nil {
		return ""
	}

	chainID, _ := (*c)["ChainID"].(string)
	return chainID
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

	for i := range c.EVM {
		if input := c.EVM[i]; input == nil {
			c.EVM[i] = &evmcfg.EVMConfig{Chain: evmcfg.Defaults(nil)}
		} else {
			input.Chain = evmcfg.Defaults(input.ChainID, &input.Chain)
		}
	}

	for i := range c.Cosmos {
		if c.Cosmos[i] == nil {
			c.Cosmos[i] = new(coscfg.TOMLConfig)
		}
		c.Cosmos[i].Chain.SetDefaults()
	}

	for i := range c.Solana {
		if c.Solana[i] == nil {
			c.Solana[i] = new(solcfg.TOMLConfig)
		}
		c.Solana[i].Chain.SetDefaults()
	}

	for i := range c.Starknet {
		if c.Starknet[i] == nil {
			c.Starknet[i] = new(stkcfg.TOMLConfig)
		}
		c.Starknet[i].Chain.SetDefaults()
	}
}

func (c *Config) SetFrom(f *Config) (err error) {
	c.Core.SetFrom(&f.Core)

	if err1 := c.EVM.SetFrom(&f.EVM); err1 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err1, "EVM"))
	}

	if err2 := c.Cosmos.SetFrom(&f.Cosmos); err2 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err2, "Cosmos"))
	}

	if err3 := c.Solana.SetFrom(&f.Solana); err3 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err3, "Solana"))
	}

	if err4 := c.Starknet.SetFrom(&f.Starknet); err4 != nil {
		err = multierr.Append(err, commonconfig.NamedMultiErrorList(err4, "Starknet"))
	}

	// the plugin should handle it's own defaults and merging
	c.Aptos = f.Aptos

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

var ErrInvalidSecrets = errors.New("invalid secrets")

// Validate validates every consitutent secret and return an accumulated error
func (s *Secrets) Validate() error {
	if err := commonconfig.Validate(s); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidSecrets, err)
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
		return fmt.Errorf("%w: %s", ErrInvalidSecrets, err)
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
