package chainlink

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"

	"github.com/smartcontractkit/chainlink/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/core/chains/terra"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	config "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
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

	Terra terra.TerraConfigs `toml:",omitempty"`
}

// TODO these method have to be on the generalConfig, not just the Config
func (c *Config) EVMConfigs() evmcfg.EVMConfigs {
	return c.EVM
}

func (c *Config) SolanaConfigs() solana.SolanaConfigs {
	return c.Solana
}

func (c *Config) StarknetConfigs() starknet.StarknetConfigs {
	return c.Starknet
}

func (c *Config) TerraConfigs() terra.TerraConfigs {
	return c.Terra
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
	return config.Validate(c)
}

// SetDefaults initializes unset fields with default values.
func (c *Config) SetDefaults() {
	core := config.CoreDefaults()
	core.SetFrom(&c.Core)
	c.Core = core

	for _, input := range c.EVM {
		ch, _ := evmcfg.Defaults(input.ChainID)
		ch.SetFrom(&input.Chain)
		input.Chain = ch
	}

	//TODO terra and solana defaults https://app.shortcut.com/chainlinklabs/story/37975/chains-nodes-should-be-read-from-the-config-interface
}

type Secrets struct {
	config.Secrets
}

func (s *Secrets) Validate() error {
	return config.Validate(s)
}

// SetOverrides overrides fields with values from ENV vars and password files.
func (s *Secrets) SetOverrides(keystorePasswordFileName, vrfPasswordFileName *string) error {
	// Override DB and Explorer secrets from ENV vars, if present
	v := viper.New()
	v.AutomaticEnv()
	//TODO CL_ prefix: https://app.shortcut.com/chainlinklabs/story/23679/prefix-all-env-vars-with-cl
	if dbURL := v.GetString("DATABASE_URL"); dbURL != "" {
		parsedURL, err := models.ParseURL(dbURL)
		if err != nil {
			return err
		}
		s.DatabaseURL = parsedURL
	}
	if dbBackupUrl := v.GetString("DATABASE_BACKUP_URL"); dbBackupUrl != "" {
		parsedURL, err := models.ParseURL(dbBackupUrl)
		if err != nil {
			return err
		}
		s.DatabaseBackupURL = parsedURL
	}
	if explorerKey := v.GetString("EXPLORER_ACCESS_KEY"); explorerKey != "" {
		s.ExplorerAccessKey = &explorerKey
	}
	if explorerSecret := v.GetString("EXPLORER_SECRET"); explorerSecret != "" {
		s.ExplorerSecret = &explorerSecret
	}

	// Override Keystore and VRF passwords from corresponding files, if present
	if keystorePasswordFileName != nil {
		keystorePwd, err := utils.PasswordFromFile(*keystorePasswordFileName)
		if err != nil {
			return err
		}
		s.KeystorePassword = &keystorePwd
	}
	if vrfPasswordFileName != nil {
		vrfPwd, err := utils.PasswordFromFile(*vrfPasswordFileName)
		if err != nil {
			return err
		}
		s.VRFPassword = &vrfPwd
	}
	return nil
}
