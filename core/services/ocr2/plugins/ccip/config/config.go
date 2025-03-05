package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/bytes"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
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
	AggregatorPrices map[common.Address]AggregatorPriceConfig `json:"aggregatorPrices"`
	StaticPrices     map[common.Address]StaticPriceConfig     `json:"staticPrices"`
}

// AggregatorPriceConfig specifies a price retrieved from an aggregator contract.
type AggregatorPriceConfig struct {
	ChainID                   uint64         `json:"chainID,string"`
	AggregatorContractAddress common.Address `json:"contractAddress"`
}

// StaticPriceConfig specifies a price defined statically.
type StaticPriceConfig struct {
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

func (c *DynamicPriceGetterConfig) Validate() error {
	for addr, v := range c.AggregatorPrices {
		if addr == utils.ZeroAddress {
			return fmt.Errorf("token address is zero")
		}
		if v.AggregatorContractAddress == utils.ZeroAddress {
			return fmt.Errorf("aggregator contract address is zero")
		}
		if v.ChainID == 0 {
			return fmt.Errorf("chain id is zero")
		}
	}

	for addr, v := range c.StaticPrices {
		if addr == utils.ZeroAddress {
			return fmt.Errorf("token address is zero")
		}
		if v.ChainID == 0 {
			return fmt.Errorf("chain id is zero")
		}
	}

	// Ensure no duplication in token price resolution rules.
	if c.AggregatorPrices != nil && c.StaticPrices != nil {
		for tk := range c.AggregatorPrices {
			if _, exists := c.StaticPrices[tk]; exists {
				return fmt.Errorf("token %s defined in both aggregator and static price rules", tk)
			}
		}
	}
	return nil
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
