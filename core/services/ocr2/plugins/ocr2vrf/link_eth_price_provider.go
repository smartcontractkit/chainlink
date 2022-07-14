package ocr2vrf

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2vrf/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/aggregator_v3_interface"
)

// linkEthPriceProvider provides conversation rate between Link and native token using price feeds
type linkEthPriceProvider struct {
	aggregator aggregator_v3_interface.AggregatorV3InterfaceInterface
}

var _ types.JuelsPerFeeCoin = (*linkEthPriceProvider)(nil)

func NewLinkEthPriceProvider(linkEthFeedAddress common.Address, client evmclient.Client) (types.JuelsPerFeeCoin, error) {
	aggregator, err := aggregator_v3_interface.NewAggregatorV3Interface(linkEthFeedAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "new aggregator v3 interface")
	}
	return &linkEthPriceProvider{aggregator: aggregator}, nil
}

func (p *linkEthPriceProvider) JuelsPerFeeCoin() (*big.Int, error) {
	roundData, err := p.aggregator.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		return nil, errors.Wrap(err, "get aggregator latest answer")
	}
	return roundData.Answer, nil
}
