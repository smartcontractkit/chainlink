package juelsfeecoin

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2vrf/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/aggregator_v3_interface"
)

// linkEthPriceProvider provides conversation rate between Link and native token using price feeds
type linkEthPriceProvider struct {
	aggregator aggregator_v3_interface.AggregatorV3InterfaceInterface
	timeout    time.Duration
	stubbed    bool
}

var _ types.JuelsPerFeeCoin = (*linkEthPriceProvider)(nil)

func NewLinkEthPriceProvider(linkEthFeedAddress common.Address, client evmclient.Client, timeout time.Duration) (types.JuelsPerFeeCoin, error) {
	aggregator, err := aggregator_v3_interface.NewAggregatorV3Interface(linkEthFeedAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "new aggregator v3 interface")
	}
	// Return the stubbed implementation, as we are not currently using juelsPerFeeCoin.
	return &linkEthPriceProvider{aggregator: aggregator, timeout: timeout, stubbed: true}, nil
}

func (p *linkEthPriceProvider) JuelsPerFeeCoin() (*big.Int, error) {
	if p.stubbed {
		return p.juelsPerFeeCoinStubbed()
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	roundData, err := p.aggregator.LatestRoundData(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, errors.Wrap(err, "get aggregator latest answer")
	}
	return roundData.Answer, nil
}

// Stubbed implementation, returns 0 and does not make an RPC call.
func (p *linkEthPriceProvider) juelsPerFeeCoinStubbed() (*big.Int, error) {
	return big.NewInt(0), nil
}
