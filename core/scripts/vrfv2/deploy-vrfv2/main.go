package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/services/vrf"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

var weiPerUnitLink = decimal.RequireFromString("10000000000000000")
var waitForMine = 10 * time.Second

// Send eth from prefunded account.
// Amount is number of ETH not wei.
func sendEth(chainID int64, ec ethclient.Client, to common.Address, amount int) {
	key, err := crypto.HexToECDSA("34d2ee6c703f755f9a205e322c68b8ff3425d915072ca7483190ac69684e548c")
	panicErr(err)
	nonce, err := ec.PendingNonceAt(context.Background(), common.HexToAddress("9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f"))
	panicErr(err)
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	panicErr(err)
	tx := types.NewTransaction(nonce, to, big.NewInt(0).Mul(big.NewInt(int64(amount)), big.NewInt(1000000000000000000)), uint64(21000), gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), key)
	panicErr(err)
	err = ec.SendTransaction(context.Background(), signedTx)
	panicErr(err)
	time.Sleep(waitForMine)
}

func main() {
	ec, err := ethclient.Dial("http://127.0.0.1:8545")
	panicErr(err)

	chainID := int64(34055)
	// Preloaded devnet key https://github.com/smartcontractkit/devnet/blob/master/passwords.json
	key, err := crypto.HexToECDSA("34d2ee6c703f755f9a205e322c68b8ff3425d915072ca7483190ac69684e548c")
	panicErr(err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(chainID))
	panicErr(err)
	fmt.Println(user)

	// --- First time deploy ----
	// Fund oracle address
	oracleAddress := common.HexToAddress("0x8EC241824833726911a6EE5dD41C7304C4d4897c")
	sendEth(chainID, *ec, oracleAddress, 100)
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
		uint16(1),                              // minRequestConfirmations
		uint32(1000000),                        // max gas limit
		uint32(60*60*24),                       // stalenessSeconds
		uint32(vrf.GasAfterPaymentCalculation), // gasAfterPaymentCalculation
		big.NewInt(10000000000000000),          // 0.01 eth per link fallbackLinkPrice
		vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
			FulfillmentFlatFeeLinkPPMTier1: uint32(10000),
			FulfillmentFlatFeeLinkPPMTier2: uint32(1000),
			FulfillmentFlatFeeLinkPPMTier3: uint32(100),
			FulfillmentFlatFeeLinkPPMTier4: uint32(10),
			FulfillmentFlatFeeLinkPPMTier5: uint32(1),
			ReqsForTier2:                   big.NewInt(10),
			ReqsForTier3:                   big.NewInt(20),
			ReqsForTier4:                   big.NewInt(30),
			ReqsForTier5:                   big.NewInt(40),
		},
	)
	panicErr(err)
	time.Sleep(waitForMine)
	minreq, maxgas, kh, err := coordinatorContract.GetRequestConfig(nil)
	panicErr(err)
	fmt.Printf("Coordinator config %v %v %v %v %v\n", coordinatorAddress, linkAddress, minreq, maxgas, kh)
	// Deploy consumer
	consumerContractAddress, _, consumerContract, err :=
		vrf_consumer_v2.DeployVRFConsumerV2(
			user, ec, coordinatorAddress, linkAddress)
	panicErr(err)
	time.Sleep(waitForMine)
	// Transfer it 1000 link
	link := int64(1000)
	_, err = linkContract.Transfer(user, consumerContractAddress, big.NewInt(0).Mul(big.NewInt(link), big.NewInt(1000000000000000000))) // Actually, LINK
	panicErr(err)
	time.Sleep(waitForMine)
	// Create an fund subscription with a link
	_, err = consumerContract.TestCreateSubscriptionAndFund(user, big.NewInt(0).Mul(big.NewInt(link), big.NewInt(1000000000000000000)))
	panicErr(err)
	time.Sleep(waitForMine)
	subID, err := consumerContract.SSubId(nil)
	panicErr(err)
	fmt.Println("Sub ID", subID)
	// Register the proving key
	// Note the 04 is a version byte, it means uncompressed pubkey.
	pubBytes, err := hex.DecodeString("041d460efc55a1cbce820409f29986574b9e99358f52c46073d9014c579aaad770e85024a1b8aaf115606f2014ff41db40af47763e0a64df060f797d5f2092ed2e")
	panicErr(err)
	pk, err := crypto.UnmarshalPubkey(pubBytes)
	panicErr(err)
	_, err = coordinatorContract.RegisterProvingKey(user, oracleAddress, [2]*big.Int{pk.X, pk.Y})
	panicErr(err)
	sub, err := coordinatorContract.GetSubscription(nil, subID)
	panicErr(err)
	fmt.Printf("Sub %+v\n", sub)
	fmt.Printf("Coordinator: %v, Link %v, Consumer %v, SubID %v", coordinatorAddress, linkAddress, consumerContractAddress, subID)
}
