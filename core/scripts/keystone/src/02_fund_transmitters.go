package src

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func distributeFunds(nodeKeys []NodeKeys, env helpers.Environment) {
	fmt.Println("Funding transmitters...")
	transmittersStr := []string{}
	fundingAmount := big.NewInt(500000000000000000) // 0.5 ETH
	minThreshold := big.NewInt(50000000000000000)   // 0.05 ETH

	for _, n := range nodeKeys {
		balance, err := getBalance(n.EthAddress, env)
		if err != nil {
			fmt.Printf("Error fetching balance for %s: %v\n", n.EthAddress, err)
			continue
		}
		if balance.Cmp(minThreshold) < 0 {
			fmt.Printf(
				"Transmitter %s has insufficient funds, funding with %s ETH. Current balance: %s, threshold: %s\n",
				n.EthAddress,
				conversions.WeiToEther(fundingAmount).String(),
				conversions.WeiToEther(balance).String(),
				conversions.WeiToEther(minThreshold).String(),
			)
			transmittersStr = append(transmittersStr, n.EthAddress)
		}
	}

	if len(transmittersStr) > 0 {
		helpers.FundNodes(env, transmittersStr, fundingAmount)
	} else {
		fmt.Println("All transmitters have sufficient funds.")
	}
}

func getBalance(address string, env helpers.Environment) (*big.Int, error) {
	balance, err := env.Ec.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
