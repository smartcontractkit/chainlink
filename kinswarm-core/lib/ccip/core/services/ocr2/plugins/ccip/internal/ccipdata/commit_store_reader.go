package ccipdata

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
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

// FetchCommitStoreStaticConfig provides access to a commitStore's static config, which is required to access the source chain ID.
func FetchCommitStoreStaticConfig(address common.Address, ec client.Client) (commit_store.CommitStoreStaticConfig, error) {
	commitStore, err := loadCommitStore(address, ec)
	if err != nil {
		return commit_store.CommitStoreStaticConfig{}, err
	}
	return commitStore.GetStaticConfig(&bind.CallOpts{})
}

func loadCommitStore(commitStoreAddress common.Address, client client.Client) (commit_store.CommitStoreInterface, error) {
	_, err := ccipconfig.VerifyTypeAndVersion(commitStoreAddress, client, ccipconfig.CommitStore)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid commitStore contract")
	}
	return commit_store.NewCommitStore(commitStoreAddress, client)
}
