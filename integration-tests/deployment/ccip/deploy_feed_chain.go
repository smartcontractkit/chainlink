package ccipdeployment

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
)

const (
	LINK     = "LINK"
	WETH     = "WETH"
	DECIMALS = 18
)

var (
	LINK_PRICE = big.NewInt(5e18)
)

func DeployFeeds(lggr logger.Logger, chain deployment.Chain) (deployment.AddressBook, map[string]common.Address, error) {
	ab := deployment.NewMemoryAddressBook()
	//TODO: Maybe append LINK to the contract name
	linkTV := deployment.NewTypeAndVersion(PriceFeed, deployment.Version1_0_0)
	mockLinkFeed, err := deployContract(lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface] {
			linkFeed, tx, _, err1 := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
				chain.DeployerKey,
				chain.Client,
				DECIMALS,   // decimals
				LINK_PRICE, // initialAnswer
			)
			aggregatorCr, err2 := aggregator_v3_interface.NewAggregatorV3Interface(linkFeed, chain.Client)

			var err error
			if err1 != nil || err2 != nil {
				err = fmt.Errorf("linkFeedError: %v, AggregatorInterfaceError: %v", err1, err2)
			}
			return ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface]{
				Address: linkFeed, Contract: aggregatorCr, Tv: linkTV, Tx: tx, Err: err,
			}
		})

	if err != nil {
		lggr.Errorw("Failed to deploy link feed", "err", err)
		return ab, nil, err
	}

	lggr.Infow("deployed mockLinkFeed", "addr", mockLinkFeed.Address)

	tvToAddress := map[string]common.Address{
		linkTV.String(): mockLinkFeed.Address,
	}
	return ab, tvToAddress, nil
}
