package view

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment/common/view"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/view/v1_0"
)

type ChainView struct {
	// v1.0
	DataFeedsCache  map[string]v1_0.CacheView `json:"dataFeedsCache,omitempty"`
	AggregatorProxy map[string]v1_0.ProxyView `json:"aggregatorProxy,omitempty"`
	FeedConfig      v1_0.FeedState            `json:"feedConfig,omitempty"`
}

func NewChain() ChainView {
	return ChainView{
		// v1.0
		DataFeedsCache:  make(map[string]v1_0.CacheView),
		AggregatorProxy: make(map[string]v1_0.ProxyView),
	}
}

type DataFeedsView struct {
	Chains map[string]ChainView    `json:"chains,omitempty"`
	Nops   map[string]view.NopView `json:"nops,omitempty"`
}

func (v DataFeedsView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias DataFeedsView
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
