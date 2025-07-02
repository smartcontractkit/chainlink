package chainsupport

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func InfuraUrl(network string) string {
	return fmt.Sprintf("wss://%s.infura.io/ws/v3/de4f73b9679f41219d9a0c386367be1b", network)
}

var NetworkToChainID = map[string]int{
	"sepolia":        11155111,
	"avalanche-fuji": 43113,
}

func EthClient(rpcURL string) (client *ethclient.Client) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
