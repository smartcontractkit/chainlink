package ccipdata

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-evm/pkg/gas"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
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

type CommitStoreReader interface {
	cciptypes.CommitStoreReader
	SetGasEstimator(ctx context.Context, gpe gas.EvmFeeEstimator) error
	SetSourceMaxGasPrice(ctx context.Context, sourceMaxGasPrice *big.Int) error
}
