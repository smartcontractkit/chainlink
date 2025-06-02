package v1_0

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/aggregator_proxy"
)

type ProxyView struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Address        common.Address `json:"address,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
	Description    string         `json:"description,omitempty"`
	Aggregator     common.Address `json:"aggregator,omitempty"`
}

// GenerateAggregatorProxyView generates a ProxyView from a AggregatorProxy contract.
func GenerateAggregatorProxyView(proxy *proxy.AggregatorProxy) (ProxyView, error) {
	if proxy == nil {
		return ProxyView{}, errors.New("cannot generate view for nil AggregatorProxy")
	}

	description, err := proxy.Description(nil)
	if err != nil {
		fmt.Printf("failed to get description for AggregatorProxy: %v\n", err)
		description = ""
	}

	owner, err := proxy.Owner(nil)
	if err != nil {
		return ProxyView{}, fmt.Errorf("failed to get owner for AggregatorProxy: %w", err)
	}

	aggregator, err := proxy.Aggregator(nil)
	if err != nil {
		return ProxyView{}, fmt.Errorf("failed to get aggregator for AggregatorProxy: %w", err)
	}

	return ProxyView{
		Address:        proxy.Address(),
		Owner:          owner,
		Description:    description,
		TypeAndVersion: "AggregatorProxy 1.0.0",
		Aggregator:     aggregator,
	}, nil
}
