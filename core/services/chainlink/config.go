package chainlink

import (
	"fmt"

	"errors"

	"go.uber.org/multierr"

	gotoml "github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"

	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/config/docs"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// Config is the root type used for TOML configuration.
//
// See docs at /docs/CONFIG.md generated via config.GenerateDocs from /internal/config/docs.toml
//
// When adding a new field:
//   - consider including a unit suffix with the field name
//   - TOML is limited to int64/float64, so fields requiring greater range/precision must use non-standard types
//     implementing encoding.TextMarshaler/TextUnmarshaler, like utils.Big and decimal.Decimal
//   - std lib types that don't implement encoding.TextMarshaler/TextUnmarshaler (time.Duration, url.URL, big.Int) won't
//     work as expected, and require wrapper types. See models.Duration, models.URL, utils.Big.
type Config struct {
	toml.Core

	EVM evmcfg.EVMConfigs `toml:",omitempty"`

	Cosmos cosmos.CosmosConfigs `toml:",omitempty"`

	Solana solana.SolanaConfigs `toml:",omitempty"`

	Starknet starknet.StarknetConfigs `toml:",omitempty"`
}

// TOMLString returns a TOML encoded string.
func (c *Config) TOMLString() (string, error) {
	b, err := gotoml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Config) Validate() error {
	if err := config.Validate(c); err != nil {
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
			c.Cosmos[i] = new(cosmos.CosmosConfig)
		}
		c.Cosmos[i].Chain.SetDefaults()
	}

	for i := range c.Solana {
		if c.Solana[i] == nil {
			c.Solana[i] = new(solana.SolanaConfig)
		}
		c.Solana[i].Chain.SetDefaults()
	}

	for i := range c.Starknet {
		if c.Starknet[i] == nil {
			c.Starknet[i] = new(starknet.StarknetConfig)
		}
		c.Starknet[i].Chain.SetDefaults()
	}
}

func (c *Config) SetFrom(f *Config) (err error) {
	c.Core.SetFrom(&f.Core)

	if err1 := c.EVM.SetFrom(&f.EVM); err1 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err1, "EVM"))
	}

	if err2 := c.Cosmos.SetFrom(&f.Cosmos); err2 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err2, "Cosmos"))
	}

	if err3 := c.Solana.SetFrom(&f.Solana); err3 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err3, "Solana"))
	}

	if err4 := c.Starknet.SetFrom(&f.Starknet); err4 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err4, "Starknet"))
	}

	_, err = utils.MultiErrorList(err)

	return err
}

type Secrets struct {
	toml.Secrets
}

func (s *Secrets) SetFrom(f *Secrets) (err error) {
	if err1 := s.Database.SetFrom(&f.Database); err1 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err1, "Database"))
	}

	if err2 := s.Explorer.SetFrom(&f.Explorer); err2 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err2, "Explorer"))
	}

	if err3 := s.Password.SetFrom(&f.Password); err3 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err3, "Password"))
	}

	if err4 := s.Pyroscope.SetFrom(&f.Pyroscope); err4 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err4, "Pyroscope"))
	}

	if err5 := s.Prometheus.SetFrom(&f.Prometheus); err5 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err5, "Prometheus"))
	}

	if err6 := s.Mercury.SetFrom(&f.Mercury); err6 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err6, "Mercury"))
	}

	if err7 := s.Threshold.SetFrom(&f.Threshold); err7 != nil {
		err = multierr.Append(err, config.NamedMultiErrorList(err7, "Threshold"))
	}

	_, err = utils.MultiErrorList(err)

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
	if err := config.Validate(s); err != nil {
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
	if err := config.Validate(v); err != nil {
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
	if explorerKey := env.ExplorerAccessKey.Get(); explorerKey != "" {
		s.Explorer.AccessKey = &explorerKey
	}
	if explorerSecret := env.ExplorerSecret.Get(); explorerSecret != "" {
		s.Explorer.Secret = &explorerSecret
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
