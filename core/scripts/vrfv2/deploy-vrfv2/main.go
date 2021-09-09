package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

var weiPerUnitLink = decimal.RequireFromString("10000000000000000")
var waitForMine = 8 * time.Second

func main() {
	ec, err := ethclient.Dial("http://127.0.0.1:8545")
	panicErr(err)

	chainID := int64(34055)
	// Preloaded devnet key https://github.com/smartcontractkit/devnet/blob/master/passwords.json
	key, err := crypto.HexToECDSA("34d2ee6c703f755f9a205e322c68b8ff3425d915072ca7483190ac69684e548c")
	panicErr(err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(chainID))
	panicErr(err)

	// --- First time deploy ----
	// Fund oracle address
	oracleAddress := "0x5f1bbb70AEeb5754BD68EdF856a7234B232e6858"
	nonce, err := ec.PendingNonceAt(context.Background(), common.HexToAddress("9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f"))
	panicErr(err)
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	panicErr(err)
	tx := types.NewTransaction(nonce, common.HexToAddress(oracleAddress), big.NewInt(500000000000000000), uint64(21000), gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), key)
	panicErr(err)
	err = ec.SendTransaction(context.Background(), signedTx)
	panicErr(err)
	time.Sleep(waitForMine)
	// Deploy link
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		user, ec)
	panicErr(err)
	time.Sleep(waitForMine)
	// Deploy feed
	linkEthFeed, _, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			user, ec, 18, weiPerUnitLink.BigInt()) // 0.01 eth per link
	panicErr(err)
	time.Sleep(waitForMine)
	// Deploy coordinator
	coordinatorAddress, _, coordinatorContract, err :=
		vrf_coordinator_v2.DeployVRFCoordinatorV2(
			user, ec, linkAddress, common.Address{}, linkEthFeed)
	panicErr(err)
	time.Sleep(waitForMine)
	// Set coordinators config
	_, err = coordinatorContract.SetConfig(user,
		uint16(1),    // minRequestConfirmations
		uint32(1000), // 0.0001 link flat fee
		uint32(1000000),
		uint32(60*60*24),                       // stalenessSeconds
		uint32(vrf.GasAfterPaymentCalculation), // gasAfterPaymentCalculation
		big.NewInt(10000000000000000),          // 0.01 eth per link fallbackLinkPrice
		big.NewInt(1000000000000000000),        // Minimum subscription balance 0.01 link
	)
	panicErr(err)
	time.Sleep(waitForMine)
	c, err := coordinatorContract.GetConfig(nil)
	panicErr(err)
	fmt.Printf("Coordinator config %+v\n", c)
	// Deploy consumer
	consumerContractAddress, _, consumerContract, err :=
		vrf_consumer_v2.DeployVRFConsumerV2(
			user, ec, coordinatorAddress, linkAddress)
	panicErr(err)
	time.Sleep(waitForMine)
	// Transfer it a link
	_, err = linkContract.Transfer(user, consumerContractAddress, big.NewInt(1000000000000000000)) // Actually, LINK
	panicErr(err)
	time.Sleep(waitForMine)
	// Create an fund subscription with a link
	_, err = consumerContract.TestCreateSubscriptionAndFund(user, big.NewInt(1000000000000000000))
	panicErr(err)
	time.Sleep(waitForMine)
	subID, err := consumerContract.SSubId(nil)
	panicErr(err)
	fmt.Println("Sub ID", subID)
	// Register the proving key
	// Note the 04 is a version byte, it means uncompressed pubkey.
	pubBytes, err := hex.DecodeString("046bef3bd2b3043e6c8a5d482ed22c23c9c31c05e854676e8856869a54fe781922578f695613f819b8fe71d49401a13470e495ebc62dda8d268a8615fb30e72b95")
	panicErr(err)
	pk, err := crypto.UnmarshalPubkey(pubBytes)
	panicErr(err)
	_, err = coordinatorContract.RegisterProvingKey(user, common.HexToAddress(oracleAddress), [2]*big.Int{pk.X, pk.Y})
	panicErr(err)
	sub, err := coordinatorContract.GetSubscription(nil, subID)
	panicErr(err)
	fmt.Printf("Sub %+v\n", sub)
	fmt.Printf("Coordinator: %v, Link %v, Consumer %v, SubID %v", coordinatorAddress, linkAddress, consumerContractAddress, subID)
}
