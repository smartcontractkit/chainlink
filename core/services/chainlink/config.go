package chainlink

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink/core/chains/starknet"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	config "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/store/models"
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
	config.Core

	EVM evmcfg.EVMConfigs `toml:",omitempty"`

	Solana solana.SolanaConfigs `toml:",omitempty"`

	Starknet starknet.StarknetConfigs `toml:",omitempty"`
}

// TOMLString returns a TOML encoded string.
func (c *Config) TOMLString() (string, error) {
	b, err := toml.Marshal(c)
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
	core := config.CoreDefaults()
	core.SetFrom(&c.Core)
	c.Core = core

	for i := range c.EVM {
		if input := c.EVM[i]; input == nil {
			c.EVM[i] = &evmcfg.EVMConfig{Chain: evmcfg.Defaults(nil)}
		} else {
			input.Chain = evmcfg.Defaults(input.ChainID, &input.Chain)
		}
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

func (c *Config) SetFrom(f *Config) {
	c.Core.SetFrom(&f.Core)
	c.EVM.SetFrom(&f.EVM)
	c.Solana.SetFrom(&f.Solana)
	c.Starknet.SetFrom(&f.Starknet)
}

type Secrets struct {
	config.Secrets
}

// TOMLString returns a TOML encoded string with secret values redacted.
func (s *Secrets) TOMLString() (string, error) {
	b, err := toml.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *Secrets) Validate() error {
	if err := config.Validate(s); err != nil {
		return fmt.Errorf("invalid secrets: %w", err)
	}
	return nil
}

// setEnv overrides fields from ENV vars, if present.
func (s *Secrets) setEnv() error {
	if dbURL := config.EnvDatabaseURL.Get(); dbURL != "" {
		s.Database.URL = new(models.SecretURL)
		if err := s.Database.URL.UnmarshalText([]byte(dbURL)); err != nil {
			return err
		}
	}
	if dbBackupUrl := config.EnvDatabaseBackupURL.Get(); dbBackupUrl != "" {
		s.Database.BackupURL = new(models.SecretURL)
		if err := s.Database.BackupURL.UnmarshalText([]byte(dbBackupUrl)); err != nil {
			return err
		}
	}
	if config.EnvDatabaseAllowSimplePasswords.IsTrue() {
		s.Database.AllowSimplePasswords = true
	}
	if explorerKey := config.EnvExplorerAccessKey.Get(); explorerKey != "" {
		s.Explorer.AccessKey = &explorerKey
	}
	if explorerSecret := config.EnvExplorerSecret.Get(); explorerSecret != "" {
		s.Explorer.Secret = &explorerSecret
	}
	if keystorePassword := config.EnvPasswordKeystore.Get(); keystorePassword != "" {
		s.Password.Keystore = &keystorePassword
	}
	if vrfPassword := config.EnvPasswordVRF.Get(); vrfPassword != "" {
		s.Password.VRF = &vrfPassword
	}
	if pyroscopeAuthToken := config.EnvPyroscopeAuthToken.Get(); pyroscopeAuthToken != "" {
		s.Pyroscope.AuthToken = &pyroscopeAuthToken
	}
	return nil
}
