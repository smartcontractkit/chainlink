package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/bytes"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

// CommitPluginJobSpecConfig contains the plugin specific variables for the ccip.CCIPCommit plugin.
type CommitPluginJobSpecConfig struct {
	SourceStartBlock, DestStartBlock uint64            // Only for first time job add.
	OffRamp                          cciptypes.Address `json:"offRamp"`
	// TokenPricesUSDPipeline should contain a token price pipeline for the following tokens:
	//		The SOURCE chain wrapped native
	// 		The DESTINATION supported tokens (including fee tokens) as defined in destination OffRamp and PriceRegistry.
	TokenPricesUSDPipeline string `json:"tokenPricesUSDPipeline,omitempty"`
	// PriceGetterConfig defines where to get the token prices from (i.e. static or aggregator source).
	PriceGetterConfig *DynamicPriceGetterConfig `json:"priceGetterConfig,omitempty"`
}

type CommitPluginConfig struct {
	IsSourceProvider                 bool
	SourceStartBlock, DestStartBlock uint64
}

func (c CommitPluginConfig) Encode() ([]byte, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// DynamicPriceGetterConfig specifies which configuration to use for getting the price of tokens (map keys).
type DynamicPriceGetterConfig struct {
	// use TokenPrices instead this is Deprecated
	AggregatorPrices map[common.Address]AggregatorPriceConfig `json:"aggregatorPrices"`
	// use TokenPrices instead this is Deprecated
	StaticPrices map[common.Address]StaticPriceConfig `json:"staticPrices"`

	TokenPrices []TokenPriceConfig `json:"tokenPrices"`
}

// IsDeprecated returns true if the config uses the deprecated fields.
// If the config uses both old fields and the new field it is not deprecated.
func (c *DynamicPriceGetterConfig) IsDeprecated() bool {
	return (len(c.AggregatorPrices) > 0 || len(c.StaticPrices) > 0) && len(c.TokenPrices) == 0
}

// TokenPriceConfig specifies the configuration for a token price.
type TokenPriceConfig struct {
	// TokenAddress is the address of the token on the chain that is deployed on.
	TokenAddress common.Address `json:"tokenAddress"`
	// ChainSelector is the chain selector of the chain that the token is deployed on (source or dest).
	ChainSelector uint64 `json:"chainSelector,string"`
	// Exactly one of AggregatorConfig or StaticConfig must be set. It defines the source of the price.
	AggregatorConfig *AggregatorPriceConfig `json:"aggregatorConfig,omitempty"`
	StaticConfig     *StaticPriceConfig     `json:"staticConfig,omitempty"`
}

// MoveDeprecatedFields moves the deprecated fields to the new TokenPrices field.
func (c *DynamicPriceGetterConfig) MoveDeprecatedFields(
	sourceChainSelector, destChainSelector uint64,
	sourceChainNativeToken common.Address,
) error {
	if !c.IsDeprecated() {
		return nil
	}
	if len(c.TokenPrices) > 0 {
		return errors.New("config is deprecated but contains the new 'TokenPrices' field - remove deprecated fields")
	}

	tokenPricesCfg := make([]TokenPriceConfig, 0, len(c.AggregatorPrices)+len(c.StaticPrices))

	for tokenAddr, aggrCfg := range c.AggregatorPrices {
		chainSel := destChainSelector
		if tokenAddr == sourceChainNativeToken {
			chainSel = sourceChainSelector
		}

		tokenPricesCfg = append(tokenPricesCfg, TokenPriceConfig{
			TokenAddress:     tokenAddr,
			ChainSelector:    chainSel,
			AggregatorConfig: &aggrCfg,
		})
	}

	for tokenAddr, staticCfg := range c.StaticPrices {
		chainSel := destChainSelector
		if tokenAddr == sourceChainNativeToken {
			chainSel = sourceChainSelector
		}

		tokenPricesCfg = append(tokenPricesCfg, TokenPriceConfig{
			TokenAddress:  tokenAddr,
			ChainSelector: chainSel,
			StaticConfig:  &staticCfg,
		})
	}

	sort.Slice(tokenPricesCfg, func(i, j int) bool {
		if tokenPricesCfg[i].TokenAddress == tokenPricesCfg[j].TokenAddress {
			return tokenPricesCfg[i].ChainSelector < tokenPricesCfg[j].ChainSelector
		}
		return tokenPricesCfg[i].TokenAddress.String() < tokenPricesCfg[j].TokenAddress.String()
	})

	c.TokenPrices = tokenPricesCfg
	c.AggregatorPrices = nil
	c.StaticPrices = nil
	return nil
}

// Validate checks the configuration for errors.
func (c *DynamicPriceGetterConfig) Validate() error {
	type tokenKey struct {
		ChainSelector uint64
		TokenAddress  common.Address
	}

	seenTokens := make(map[tokenKey]struct{})

	for _, cfg := range c.TokenPrices {
		if cfg.TokenAddress == utils.ZeroAddress {
			return fmt.Errorf("token address is zero: %v", cfg)
		}

		if cfg.ChainSelector == 0 {
			return fmt.Errorf("chain selector is zero: %v", cfg)
		}

		if cfg.AggregatorConfig != nil && cfg.StaticConfig != nil {
			return fmt.Errorf("both aggregator and static price configuration is defined: %v", cfg)
		}

		k := tokenKey{ChainSelector: cfg.ChainSelector, TokenAddress: cfg.TokenAddress}
		if _, seen := seenTokens[k]; seen {
			return fmt.Errorf("duplicate token price configuration, (token, chain) pair appears twice: %v", cfg)
		}
		seenTokens[k] = struct{}{}

		if cfg.AggregatorConfig != nil {
			if cfg.AggregatorConfig.AggregatorContractAddress == utils.ZeroAddress {
				return fmt.Errorf("aggregator contract address is zero: %v", cfg)
			}
			if cfg.AggregatorConfig.ChainID == 0 {
				return fmt.Errorf("aggregator chain id is zero: %v", cfg)
			}
		}

		if cfg.StaticConfig != nil {
			if cfg.StaticConfig.Price == nil {
				return fmt.Errorf("static price is nil: %v", cfg)
			}
			if cfg.StaticConfig.ChainID == 0 {
				return fmt.Errorf("static chain id is zero: %v", cfg)
			}
		}

		if cfg.AggregatorConfig == nil && cfg.StaticConfig == nil {
			return fmt.Errorf("no price configuration defined: %v", cfg)
		}
	}
	return nil
}

// AggregatorPriceConfig specifies a price retrieved from an aggregator contract.
type AggregatorPriceConfig struct {
	ChainID                   uint64         `json:"chainID,string"`
	AggregatorContractAddress common.Address `json:"contractAddress"`
}

// StaticPriceConfig specifies a price defined statically.
type StaticPriceConfig struct {
	// Deprecated: ChainID field is not used.
	ChainID uint64   `json:"chainID,string"`
	Price   *big.Int `json:"price"`
}

// UnmarshalJSON provides a custom un-marshaller to handle JSON embedded in Toml content.
func (c *DynamicPriceGetterConfig) UnmarshalJSON(data []byte) error {
	type Alias DynamicPriceGetterConfig
	if bytes.HasQuotes(data) {
		trimmed := string(bytes.TrimQuotes(data))
		trimmed = strings.ReplaceAll(trimmed, "\\n", "")
		trimmed = strings.ReplaceAll(trimmed, "\\t", "")
		trimmed = strings.ReplaceAll(trimmed, "\\", "")
		return json.Unmarshal([]byte(trimmed), (*Alias)(c))
	}
	return json.Unmarshal(data, (*Alias)(c))
}

// ExecPluginJobSpecConfig contains the plugin specific variables for the ccip.CCIPExecution plugin.
type ExecPluginJobSpecConfig struct {
	SourceStartBlock, DestStartBlock uint64 // Only for first time job add.
	USDCConfig                       USDCConfig
	LBTCConfig                       LBTCConfig
}

type USDCConfig struct {
	SourceTokenAddress              common.Address
	SourceMessageTransmitterAddress common.Address
	AttestationAPI                  string
	AttestationAPITimeoutSeconds    uint
	// AttestationAPIIntervalMilliseconds can be set to -1 to disable or 0 to use a default interval.
	AttestationAPIIntervalMilliseconds int
}

type LBTCConfig struct {
	SourceTokenAddress           common.Address
	AttestationAPI               string
	AttestationAPITimeoutSeconds uint
	// AttestationAPIIntervalMilliseconds can be set to -1 to disable or 0 to use a default interval.
	AttestationAPIIntervalMilliseconds int
}

type ExecPluginConfig struct {
	SourceStartBlock, DestStartBlock uint64 // Only for first time job add.
	IsSourceProvider                 bool
	USDCConfig                       USDCConfig
	LBTCConfig                       LBTCConfig
	JobID                            string
}

func (e ExecPluginConfig) Encode() ([]byte, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (uc *USDCConfig) ValidateUSDCConfig() error {
	if uc.AttestationAPI == "" {
		return errors.New("USDCConfig: AttestationAPI is required")
	}
	if uc.AttestationAPIIntervalMilliseconds < -1 {
		return errors.New("USDCConfig: AttestationAPIIntervalMilliseconds must be -1 to disable, 0 for default or greater to define the exact interval")
	}
	if uc.SourceTokenAddress == utils.ZeroAddress {
		return errors.New("USDCConfig: SourceTokenAddress is required")
	}
	if uc.SourceMessageTransmitterAddress == utils.ZeroAddress {
		return errors.New("USDCConfig: SourceMessageTransmitterAddress is required")
	}

	return nil
}

func (lc *LBTCConfig) ValidateLBTCConfig() error {
	if lc.AttestationAPI == "" {
		return errors.New("LBTCConfig: AttestationAPI is required")
	}
	if lc.AttestationAPIIntervalMilliseconds < -1 {
		return errors.New("LBTCConfig: AttestationAPIIntervalMilliseconds must be -1 to disable, 0 for default or greater to define the exact interval")
	}
	if lc.SourceTokenAddress == utils.ZeroAddress {
		return errors.New("LBTCConfig: SourceTokenAddress is required")
	}
	return nil
}
