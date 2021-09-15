package main

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_consumer_v2"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Make a request to an already deployed setup
	chainID := int64(34055)
	key, err := crypto.HexToECDSA("34d2ee6c703f755f9a205e322c68b8ff3425d915072ca7483190ac69684e548c")
	panicErr(err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(chainID))
	panicErr(err)
	ec, err := ethclient.Dial("http://127.0.0.1:8545")
	panicErr(err)
	consumerAddress := "0x9E79d9A7F68D136ec4c1C0187B97c271CEa6008B"
	consumer, err := vrf_consumer_v2.NewVRFConsumerV2(common.HexToAddress(consumerAddress), ec)
	panicErr(err)
	// keyhash of offchain VRF proving key
	provingKey := "0x801b8899c3169bec9413ce6003a0117a2683b5d50db25c209b19a6bb9375890f"
	_, err = consumer.TestRequestRandomness(user, common.HexToHash(provingKey), uint64(1), uint16(2), uint32(300000), uint32(3))
	panicErr(err)
	time.Sleep(2 * time.Second)
}
