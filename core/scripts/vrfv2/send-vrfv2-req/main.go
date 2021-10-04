package main

import (
	"context"
	"fmt"
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
	consumerAddress := "0xf682e3491BeD0D71ef8B7144AC98c17A87fc301F"
	consumer, err := vrf_consumer_v2.NewVRFConsumerV2(common.HexToAddress(consumerAddress), ec)
	panicErr(err)

	// keyhash of offchain VRF proving key
	provingKey := "0x6ae4cabd964d1c04ad06518dfa47d8f0d17dcc0365d8e4de6ddabfdb1fdedf6e"
	numReqs := 1024
	var hashes []common.Hash
	for i := 0; i < numReqs; i++ {
		opts := user
		// Note cannot rely on estimate gas here, costs will change as request start going through.
		opts.GasLimit = 500000
		tx, err := consumer.TestRequestRandomness(user, common.HexToHash(provingKey), uint64(1), uint16(2), uint32(300000), uint32(1))
		panicErr(err)
		fmt.Println(i, tx.Hash())
		hashes = append(hashes, tx.Hash())
	}
	time.Sleep(20 * time.Second)
	failed := 0
	succeeded := 0
	notfound := 0
	for i := 0; i < numReqs; i++ {
		re, err := ec.TransactionReceipt(context.Background(), hashes[i])
		if err != nil {
			notfound++
			continue
		}
		fmt.Println(hashes[i], re.Status)
		if re.Status == 0 {
			failed++
		} else {
			succeeded++
		}
	}
	fmt.Printf("not found %v failed %v succeeded %v\n", notfound, failed, succeeded)
}
