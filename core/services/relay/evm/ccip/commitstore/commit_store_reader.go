package commitstore

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// Common to all versions
type CommitOnchainConfig commit_store.CommitStoreDynamicConfig

func (d CommitOnchainConfig) AbiString() string {
	return `
	[
		{
			"components": [
				{"name": "priceRegistry", "type": "address"}
			],
			"type": "tuple"
		}
	]`
}

func (d CommitOnchainConfig) Validate() error {
	if d.PriceRegistry == (common.Address{}) {
		return errors.New("must set Price Registry address")
	}
	return nil
}

func NewCommitOffchainConfig(
	gasPriceDeviationPPB uint32,
	gasPriceHeartBeat time.Duration,
	tokenPriceDeviationPPB uint32,
	tokenPriceHeartBeat time.Duration,
	inflightCacheExpiry time.Duration,
	priceReportingDisabled bool,
) cciptypes.CommitOffchainConfig {
	return cciptypes.CommitOffchainConfig{
		GasPriceDeviationPPB:   gasPriceDeviationPPB,
		GasPriceHeartBeat:      gasPriceHeartBeat,
		TokenPriceDeviationPPB: tokenPriceDeviationPPB,
		TokenPriceHeartBeat:    tokenPriceHeartBeat,
		InflightCacheExpiry:    inflightCacheExpiry,
		PriceReportingDisabled: priceReportingDisabled,
	}
}

// Do not change the JSON format of this struct without consulting with the RDD people first.
type JSONCommitOffchainConfig struct {
	SourceFinalityDepth      uint32
	DestFinalityDepth        uint32
	GasPriceHeartBeat        config.Duration
	DAGasPriceDeviationPPB   uint32
	ExecGasPriceDeviationPPB uint32
	TokenPriceHeartBeat      config.Duration
	TokenPriceDeviationPPB   uint32
	InflightCacheExpiry      config.Duration
	PriceReportingDisabled   bool
}

func (c JSONCommitOffchainConfig) Validate() error {
	if c.GasPriceHeartBeat.Duration() == 0 {
		return errors.New("must set GasPriceHeartBeat")
	}
	if c.ExecGasPriceDeviationPPB == 0 {
		return errors.New("must set ExecGasPriceDeviationPPB")
	}
	if c.TokenPriceHeartBeat.Duration() == 0 {
		return errors.New("must set TokenPriceHeartBeat")
	}
	if c.TokenPriceDeviationPPB == 0 {
		return errors.New("must set TokenPriceDeviationPPB")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}
	// DAGasPriceDeviationPPB is not validated because it can be 0 on non-rollups

	return nil
}
