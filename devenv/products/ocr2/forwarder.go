package ocr2

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/operator_factory"
	"github.com/smartcontractkit/chainlink/devenv/products"
)

// DeployForwarders deploys a forwarder factory and creates a number of operators and forwarders using it.
func DeployForwarders(ctx context.Context, c *ethclient.Client, auth *bind.TransactOpts, linkAddr common.Address, nodes int) ([]common.Address, []common.Address, error) { //nolint:revive // not less confusing than using named imports at all
	factoryAddr, tx, operatorFactory, err := operator_factory.DeployOperatorFactory(auth, c, linkAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy forwarder factory: %w", err)
	}
	_, err = products.WaitMinedFast(ctx, c, tx.Hash())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to wait for forwarder deployment: %w", err)
	}

	var (
		operators  []common.Address
		forwarders []common.Address
	)

	abi, err := operator_factory.OperatorFactoryMetaData.GetAbi()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get operator ABI: %w", err)
	}

	for range nodes {
		tx, err := operatorFactory.DeployNewOperatorAndForwarder(auth)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to deploy forwarder: %w", err)
		}

		r, err := products.WaitMinedFast(ctx, c, tx.Hash())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to wait for forwarder deployment: %w", err)
		}

		for _, log := range r.Logs {
			if log.Address != factoryAddr {
				continue
			}
			switch log.Topics[0] {
			case abi.Events["OperatorCreated"].ID:
				operator := common.BytesToAddress(log.Topics[1].Bytes())
				operators = append(operators, operator)
			case abi.Events["AuthorizedForwarderCreated"].ID:
				forwarder := common.BytesToAddress(log.Topics[1].Bytes())
				forwarders = append(forwarders, forwarder)
			}
		}
	}

	return operators, forwarders, nil
}
