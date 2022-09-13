package chainlink

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	soldb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	stkdb "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"
	tercfg "github.com/smartcontractkit/chainlink-terra/pkg/terra/config"
	terdb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	evmtyp "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	starknet "github.com/smartcontractkit/chainlink/core/chains/starknet/types"
	tertyp "github.com/smartcontractkit/chainlink/core/chains/terra/types"
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

	EVM EVMConfigs `toml:",omitempty"`

	Solana SolanaConfigs `toml:",omitempty"`

	Starknet StarknetConfigs `toml:",omitempty"`

	Terra TerraConfigs `toml:",omitempty"`
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
	//TODO core defaults c.Core.SetDefaults() https://app.shortcut.com/chainlinklabs/story/33615/create-new-implementation-of-chainscopedconfig-generalconfig-interfaces-that-sources-config-from-a-config-toml-file
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

// Use ENV vars or flags to override secrets
func (s *Secrets) SetOverrides(context *cli.Context) error {
	// Override DB and Explorer secrets from ENV vars, if present
	v := viper.New()
	v.AutomaticEnv()
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
	if context != nil && context.IsSet("password") {
		keystorePasswordFileName := context.String("password")
		keystorePwd, err := utils.PasswordFromFile(keystorePasswordFileName)
		if err != nil {
			return err
		}
		s.KeystorePassword = &keystorePwd
	}
	if context != nil && context.IsSet("vrfpassword") {
		vrfPasswordFileName := context.String("vrfpassword")
		vrfPwd, err := utils.PasswordFromFile(vrfPasswordFileName)
		if err != nil {
			return err
		}
		s.VRFPassword = &vrfPwd
	}
	return nil
}

type EVMConfigs []*EVMConfig

func (cs EVMConfigs) ValidateConfig() (err error) {
	chainIDs := map[string]struct{}{}
	for i, c := range cs {
		if c.ChainID == nil {
			continue
		}
		chainID := c.ChainID.String()
		if chainID == "" {
			continue
		}
		if _, ok := chainIDs[chainID]; ok {
			err = multierr.Append(err, config.ErrInvalid{Name: fmt.Sprintf("%d: ChainID", i), Msg: "duplicate - must be unique", Value: chainID})
		} else {
			chainIDs[chainID] = struct{}{}
		}
	}
	return
}

type EVMNodes []*evmcfg.Node

func (ns EVMNodes) ValidateConfig() (err error) {
	names := map[string]struct{}{}
	for i, n := range ns {
		if n.Name == nil || *n.Name == "" {
			continue
		}
		if _, ok := names[*n.Name]; ok {
			err = multierr.Append(err, config.ErrInvalid{Name: fmt.Sprintf("%d: Name", i), Msg: "duplicate - must be unique", Value: *n.Name})
		}
		names[*n.Name] = struct{}{}
	}
	return
}

type EVMConfig struct {
	ChainID *utils.Big
	Enabled *bool
	evmcfg.Chain
	Nodes EVMNodes
}

// Ensure that the embedded struct will be validated (w/o requiring a pointer receiver).
var _ config.Validated = evmcfg.Chain{}

func (c *EVMConfig) setFromDB(ch evmtyp.DBChain, nodes []evmtyp.Node) error {
	c.ChainID = &ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n evmcfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, &n)
	}
	return nil
}

func (c *EVMConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = multierr.Append(err, config.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if c.ChainID.String() == "" {
		err = multierr.Append(err, config.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}

	return
}

type SolanaConfigs []*SolanaConfig

func (cs SolanaConfigs) ValidateConfig() (err error) {
	chainIDs := map[string]struct{}{}
	for i, c := range cs {
		if c.ChainID == nil {
			continue
		}
		chainID := *c.ChainID
		if chainID == "" {
			continue
		}
		if _, ok := chainIDs[chainID]; ok {
			err = multierr.Append(err, config.ErrInvalid{Name: fmt.Sprintf("%d: ChainID", i), Msg: "duplicate - must be unique", Value: chainID})
		} else {
			chainIDs[chainID] = struct{}{}
		}
	}
	return
}

type SolanaNodes []*solcfg.Node

func (ns SolanaNodes) ValidateConfig() (err error) {
	names := map[string]struct{}{}
	for i, n := range ns {
		if n.Name == nil || *n.Name == "" {
			continue
		}
		if _, ok := names[*n.Name]; ok {
			err = multierr.Append(err, config.ErrInvalid{Name: fmt.Sprintf("%d: Name", i), Msg: "duplicate - must be unique", Value: *n.Name})
		}
		names[*n.Name] = struct{}{}
	}
	return
}

type SolanaConfig struct {
	ChainID *string
	Enabled *bool
	solcfg.Chain
	Nodes SolanaNodes
}

func (c *SolanaConfig) setFromDB(ch solana.DBChain, nodes []soldb.Node) error {
	c.ChainID = &ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n solcfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, &n)
	}
	return nil
}

func (c *SolanaConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = multierr.Append(err, config.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if *c.ChainID == "" {
		err = multierr.Append(err, config.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}

	return
}

type StarknetConfigs []*StarknetConfig

func (cs StarknetConfigs) ValidateConfig() (err error) {
	chainIDs := map[string]struct{}{}
	for _, c := range cs {
		if c.ChainID == nil {
			continue
		}
		chainID := *c.ChainID
		if chainID == "" {
			continue
		}
		if _, ok := chainIDs[chainID]; ok {
			err = multierr.Append(err, config.ErrInvalid{Name: "ChainID", Msg: "duplicate - must be unique", Value: chainID})
		} else {
			chainIDs[chainID] = struct{}{}
		}
	}
	return
}

type StarknetConfig struct {
	ChainID *string
	Enabled *bool
	stkcfg.Chain
	Nodes StarknetNodes
}

func (c *StarknetConfig) setFromDB(ch starknet.DBChain, nodes []stkdb.Node) error {
	c.ChainID = &ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n stkcfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, &n)
	}

	return nil
}

func (c *StarknetConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = multierr.Append(err, config.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if *c.ChainID == "" {
		err = multierr.Append(err, config.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}

	return
}

type StarknetNodes []*stkcfg.Node

func (ns StarknetNodes) ValidateConfig() (err error) {
	names := map[string]struct{}{}
	for _, n := range ns {
		if n.Name == nil || *n.Name == "" {
			continue
		}
		if _, ok := names[*n.Name]; ok {
			err = multierr.Append(err, config.ErrInvalid{Name: "Name", Msg: "duplicate - must be unique", Value: *n.Name})
		}
		names[*n.Name] = struct{}{}
	}
	return
}

type TerraConfigs []*TerraConfig

func (cs TerraConfigs) ValidateConfig() (err error) {
	chainIDs := map[string]struct{}{}
	for i, c := range cs {
		if c.ChainID == nil {
			continue
		}
		chainID := *c.ChainID
		if chainID == "" {
			continue
		}
		if _, ok := chainIDs[chainID]; ok {
			err = multierr.Append(err, config.ErrInvalid{Name: fmt.Sprintf("%d: ChainID", i), Msg: "duplicate - must be unique", Value: chainID})
		} else {
			chainIDs[chainID] = struct{}{}
		}
	}
	return
}

type TerraNodes []*tercfg.Node

func (ns TerraNodes) ValidateConfig() (err error) {
	names := map[string]struct{}{}
	for i, n := range ns {
		if n.Name == nil || *n.Name == "" {
			continue
		}
		if _, ok := names[*n.Name]; ok {
			err = multierr.Append(err, config.ErrInvalid{Name: fmt.Sprintf("%d: Name", i), Msg: "duplicate - must be unique", Value: *n.Name})
		}
		names[*n.Name] = struct{}{}
	}
	return
}

type TerraConfig struct {
	ChainID *string
	Enabled *bool
	tercfg.Chain
	Nodes TerraNodes
}

func (c *TerraConfig) setFromDB(ch tertyp.DBChain, nodes []terdb.Node) error {
	c.ChainID = &ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n tercfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, &n)
	}
	return nil
}

func (c *TerraConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = multierr.Append(err, config.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if *c.ChainID == "" {
		err = multierr.Append(err, config.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}

	return
}
