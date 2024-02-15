package ccipdata

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
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
) cciptypes.CommitOffchainConfig {
	return cciptypes.CommitOffchainConfig{
		GasPriceDeviationPPB:   gasPriceDeviationPPB,
		GasPriceHeartBeat:      gasPriceHeartBeat,
		TokenPriceDeviationPPB: tokenPriceDeviationPPB,
		TokenPriceHeartBeat:    tokenPriceHeartBeat,
		InflightCacheExpiry:    inflightCacheExpiry,
	}
}

//go:generate mockery --quiet --name CommitStoreReader --filename commit_store_reader_mock.go --case=underscore
type CommitStoreReader interface {
	cciptypes.CommitStoreReader
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
