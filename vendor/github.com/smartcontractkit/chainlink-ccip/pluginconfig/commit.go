package pluginconfig

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type CommitPluginConfig struct {
	// DestChain is the ccip destination chain configured for the commit plugin DON.
	DestChain cciptypes.ChainSelector `json:"destChain"`

	// NewMsgScanBatchSize is the number of max new messages to scan, typically set to 256.
	NewMsgScanBatchSize int `json:"newMsgScanBatchSize"`

	// SyncTimeout is the timeout for syncing the commit plugin reader.
	SyncTimeout time.Duration `json:"syncTimeout"`

	// SyncFrequency is the frequency at which the commit plugin reader should sync.
	SyncFrequency time.Duration `json:"syncFrequency"`

	// OffchainConfig is the offchain config set for the commit DON.
	OffchainConfig CommitOffchainConfig `json:"offchainConfig"`
}

func (c CommitPluginConfig) Validate() error {
	if c.DestChain == cciptypes.ChainSelector(0) {
		return fmt.Errorf("destChain not set")
	}

	if c.NewMsgScanBatchSize == 0 {
		return fmt.Errorf("newMsgScanBatchSize not set")
	}

	return c.OffchainConfig.Validate()
}

// ArbitrumPriceSource is the source of the TOKEN/USD price data of a particular token
// on Arbitrum.
// The commit plugin will use this to fetch prices for a particular token.
// See the PriceSources mapping in the CommitOffchainConfig struct.
type ArbitrumPriceSource struct {
	// AggregatorAddress is the address of the price feed TOKEN/USD aggregator on arbitrum.
	AggregatorAddress string `json:"aggregatorAddress"`

	// DeviationPPB is the deviation in parts per billion that the price feed is allowed to deviate
	// from the last written price on-chain before we write a new price.
	DeviationPPB cciptypes.BigInt `json:"deviationPPB"`
}

func (a ArbitrumPriceSource) Validate() error {
	if a.AggregatorAddress == "" {
		return errors.New("aggregatorAddress not set")
	}

	// aggregator must be an ethereum address
	decoded, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(a.AggregatorAddress, "0x")))
	if err != nil {
		return fmt.Errorf("aggregatorAddress must be a valid ethereum address (i.e hex encoded 20 bytes): %w", err)
	}
	if len(decoded) != 20 {
		return fmt.Errorf("aggregatorAddress must be a valid ethereum address, got %d bytes expected 20", len(decoded))
	}

	if a.DeviationPPB.Int.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("deviationPPB not set or negative, must be positive")
	}

	return nil
}

// CommitOffchainConfig is the OCR offchainConfig for the commit plugin.
// This is posted onchain as part of the OCR configuration process of the commit plugin.
// Every plugin is provided this configuration in its encoded form in the NewReportingPlugin
// method on the ReportingPluginFactory interface.
type CommitOffchainConfig struct {
	// RemoteGasPriceBatchWriteFrequency is the frequency at which the commit plugin
	// should write gas prices to the remote chain.
	RemoteGasPriceBatchWriteFrequency commonconfig.Duration `json:"remoteGasPriceBatchWriteFrequency"`

	// TokenPriceBatchWriteFrequency is the frequency at which the commit plugin should
	// write token prices to the remote chain.
	// If set to zero, no prices will be written (i.e keystone feeds would be active).
	TokenPriceBatchWriteFrequency commonconfig.Duration `json:"tokenPriceBatchWriteFrequency"`

	// PriceSources is a map of Arbitrum price sources for each token.
	// Note that the token address is that on the remote chain.
	PriceSources map[types.Account]ArbitrumPriceSource `json:"priceSources"`

	// TokenPriceChainSelector is the chain selector for the chain on which
	// the token prices are read from.
	// This will typically be an arbitrum testnet/mainnet chain depending on
	// the deployment.
	TokenPriceChainSelector uint64 `json:"tokenPriceChainSelector"`
}

func (c CommitOffchainConfig) Validate() error {
	if c.RemoteGasPriceBatchWriteFrequency.Duration() == 0 {
		return errors.New("remoteGasPriceBatchWriteFrequency not set")
	}

	// Note that commit may not have to submit prices if keystone feeds
	// are enabled for the chain.
	// If price sources are provided the batch write frequency and token price chain selector
	// config fields MUST be provided.
	if len(c.PriceSources) > 0 &&
		(c.TokenPriceBatchWriteFrequency.Duration() == 0 || c.TokenPriceChainSelector == 0) {
		return fmt.Errorf("tokenPriceBatchWriteFrequency (%s) or tokenPriceChainSelector (%d) not set",
			c.TokenPriceBatchWriteFrequency, c.TokenPriceChainSelector)
	}

	// if len(c.PriceSources) == 0 the other fields are ignored.

	return nil
}

// EncodeCommitOffchainConfig encodes a CommitOffchainConfig into bytes using JSON.
func EncodeCommitOffchainConfig(c CommitOffchainConfig) ([]byte, error) {
	return json.Marshal(c)
}

// DecodeCommitOffchainConfig JSON decodes a CommitOffchainConfig from bytes.
func DecodeCommitOffchainConfig(encodedCommitOffchainConfig []byte) (CommitOffchainConfig, error) {
	var c CommitOffchainConfig
	if err := json.Unmarshal(encodedCommitOffchainConfig, &c); err != nil {
		return c, err
	}
	return c, nil
}
