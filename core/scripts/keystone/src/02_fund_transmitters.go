package src

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type fundTransmitters struct{}

func NewFundTransmittersCommand() *fundTransmitters {
	return &fundTransmitters{}
}

func (g *fundTransmitters) Name() string {
	return "fund-transmitters"
}

func (g *fundTransmitters) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ExitOnError)
	// create flags for all of the env vars then set the env vars to normalize the interface
	// this is a bit of a hack but it's the easiest way to make this work
	chainID := fs.Int64("chainid", 1337, "chain ID of the Ethereum network to deploy to")
	ethUrl := fs.String("ethurl", "", "URL of the Ethereum node")
	accountKey := fs.String("accountkey", "", "private key of the account to deploy from")
	nodeSetsPath := fs.String("nodesets", defaultNodeSetsPath, "Custom node sets location")
	nodeSetSize := fs.Int("nodesetsize", 5, "number of nodes in a nodeset")

	err := fs.Parse(args)

	if err != nil ||
		*ethUrl == "" || ethUrl == nil ||
		*accountKey == "" || accountKey == nil {
		fs.Usage()
		os.Exit(1)
	}

	os.Setenv("ETH_URL", *ethUrl)
	os.Setenv("ETH_CHAIN_ID", fmt.Sprintf("%d", *chainID))
	os.Setenv("ACCOUNT_KEY", *accountKey)
	os.Setenv("INSECURE_SKIP_VERIFY", "true")
	env := helpers.SetupEnv(false)

	nodeSet := downloadNodeSets(*chainID, *nodeSetsPath, *nodeSetSize).Workflow
	distributeFunds(nodeSet, env)
}

func distributeFunds(nodeSet NodeSet, env helpers.Environment) {
	fmt.Println("Funding transmitters...")
	transmittersStr := []string{}
	fundingAmount := big.NewInt(500000000000000000) // 0.5 ETH
	minThreshold := big.NewInt(50000000000000000)   // 0.05 ETH

	for _, n := range nodeSet.NodeKeys {
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
