package v1_0

import (
	"errors"
	"fmt"

	commoncldchangesets "github.com/smartcontractkit/cld-changesets/pkg/cldfutil"

	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
)

type CacheView struct {
	commoncldchangesets.ContractMetaData
}

// GenerateDataFeedsCacheView generates a CacheView from a DataFeedsCache contract.
func GenerateDataFeedsCacheView(cache *cache.DataFeedsCache) (CacheView, error) {
	if cache == nil {
		return CacheView{}, errors.New("cannot generate view for nil DataFeedsCache")
	}
	meta, err := commoncldchangesets.NewContractMetaData(cache, cache.Address())
	if err != nil {
		return CacheView{}, fmt.Errorf("failed to generate contract metadata for DataFeedsCache: %w", err)
	}

	return CacheView{
		ContractMetaData: meta,
	}, nil
}
