package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"

	keeper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	link "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	upkeep "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/upkeep_perform_counter_restrictive_wrapper"
)

const (
	defaultGasLimit        = 8000000
	defaultUpkeepGasLimit  = 500000
	defaultCheckUpkeepData = "0x00"
)

var (
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	linkToken  *link.LinkToken
	fromAddr   common.Address

	approveAmount  *big.Int
	addFundsAmount *big.Int

	config Config
)

// Config represents configuration fields
type Config struct {
	NodeURL         string   `mapstructure:"NODE_URL"`
	ChainID         int64    `mapstructure:"CHAIN_ID"`
	PrivateKey      string   `mapstructure:"PRIVATE_KEY"`
	LinkTokenAddr   string   `mapstructure:"LINK_TOKEN_ADDR"`
	LinkETHFeedAddr string   `mapstructure:"LINK_ETH_FEED"`
	FastGasFeedAddr string   `mapstructure:"FAST_GAS_FEED"`
	Keepers         []string `mapstructure:"KEEPERS"`
	ApproveAmount   string   `mapstructure:"APPROVE_AMOUNT"`
	AddFundsAmount  string   `mapstructure:"ADD_FUNDS_AMOUNT"`
}

func (c *Config) keepers() ([]common.Address, []common.Address) {
	var addrs []common.Address
	var fromAddrs []common.Address
	for _, addr := range c.Keepers {
		addrs = append(addrs, common.HexToAddress(addr))
		fromAddrs = append(fromAddrs, fromAddr)
	}
	return addrs, fromAddrs
}

func init() {
	viper.SetDefault("APPROVE_AMOUNT", "1000000000000000000000")
	viper.SetDefault("ADD_FUNDS_AMOUNT", "1000000000000000000")

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("failed to read config: ", err)
	}

	if err = viper.Unmarshal(&config); err != nil {
		log.Fatal("failed to unmarshal config: ", err)
	}

	// Deal ETH client by the given addr
	client, err = ethclient.Dial(config.NodeURL)
	if err != nil {
		log.Fatal("failed to deal with ETH node", err)
	}

	// Parse private key
	d := new(big.Int).SetBytes(common.FromHex(config.PrivateKey))
	pkX, pkY := crypto.S256().ScalarBaseMult(d.Bytes())
	privateKey = &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: crypto.S256(),
			X:     pkX,
			Y:     pkY,
		},
		D: d,
	}

	// Init from address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddr = crypto.PubkeyToAddress(*publicKeyECDSA)

	// Create link token wrapper
	linkToken, err = link.NewLinkToken(common.HexToAddress(config.LinkTokenAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	approveAmount = big.NewInt(0)
	approveAmount.SetString(config.ApproveAmount, 10)

	addFundsAmount = big.NewInt(0)
	addFundsAmount.SetString(config.AddFundsAmount, 10)
}

func main() {
	// Deploy keeper registry
	registryAddr, tx, registryInstance, err := keeper.DeployKeeperRegistry(buildTxOpts(), client,
		common.HexToAddress(config.LinkTokenAddr),
		common.HexToAddress(config.LinkETHFeedAddr),
		common.HexToAddress(config.FastGasFeedAddr),
		200000000,
		0,
		big.NewInt(1),
		650000000,
		big.NewInt(90000),
		3,
		big.NewInt(10000000000),
		big.NewInt(200000000000000000),
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	waitForTx(tx.Hash())
	log.Println("KeeperRegistry deployed:", registryAddr.Hex(), "-", tx.Hash().Hex())

	// Approve keeper registry
	tx, err = linkToken.Approve(buildTxOpts(), registryAddr, approveAmount)
	if err != nil {
		log.Fatal("Approve failed: ", err)
	}
	waitForTx(tx.Hash())
	log.Println("KeeperRegistry approved:", registryAddr.Hex(), "-", tx.Hash().Hex())

	// Deploy Upkeeps
	for i := 0; i < 5; i++ {
		// Deploy
		upkeepAddr, tx, _, err := upkeep.DeployUpkeepPerformCounterRestrictive(buildTxOpts(), client, big.NewInt(1), big.NewInt(1))
		if err != nil {
			log.Fatal(i, " - DeployAbi failed: ", err)
		}
		waitForTx(tx.Hash())
		log.Println(i, " - Upkeep deployed:", upkeepAddr.Hex(), "-", tx.Hash().Hex())

		// Approve
		if tx, err = linkToken.Approve(buildTxOpts(), upkeepAddr, approveAmount); err != nil {
			log.Fatal(i, " - Approve failed: ", err)
		}
		waitForTx(tx.Hash())
		log.Println(i, " - Upkeep approved:", upkeepAddr.Hex(), "-", tx.Hash().Hex())

		// Register
		if tx, err = registryInstance.RegisterUpkeep(buildTxOpts(), upkeepAddr, defaultUpkeepGasLimit, fromAddr, []byte(defaultCheckUpkeepData)); err != nil {
			log.Fatal(i, " - RegisterUpkeep failed: ", err)
		}
		waitForTx(tx.Hash())
		log.Println(i, " - Upkeep registered:", upkeepAddr.Hex(), "-", tx.Hash().Hex())

		// Fund
		if tx, err = registryInstance.AddFunds(buildTxOpts(), big.NewInt(int64(i)), addFundsAmount); err != nil {
			log.Fatal(i, " - AddFunds failed: ", err)
		}
		waitForTx(tx.Hash())
		log.Println(i, " - Upkeep funded:", upkeepAddr.Hex(), "-", tx.Hash().Hex())
	}

	// Set Keepers
	keepers, owners := config.keepers()
	if tx, err = registryInstance.SetKeepers(buildTxOpts(), keepers, owners); err != nil {
		log.Fatal("SetKeepers failed: ", err)
	}
	waitForTx(tx.Hash())
	log.Println("Keepers registered:", tx.Hash().Hex())
}

func waitForTx(tx common.Hash) int {
	soc := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), soc)
	if err != nil {
		log.Fatal("SubscribeNewHead failed: ", err)
	}

	timeout := time.NewTimer(time.Minute).C
	for {
		select {
		case err = <-sub.Err():
			log.Fatal("SubscribeNewHead error: ", err)
		case <-timeout:
			log.Fatal("there is no receipt")
		case <-soc:
			transactionStatus := checkTransactionReceipt(tx)
			if transactionStatus == 0 {
				sub.Unsubscribe()
				return 0
			} else if transactionStatus == 1 {
				sub.Unsubscribe()
				return 1
			}
		}
	}
}

func checkTransactionReceipt(txHash common.Hash) int {
	tx, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return -1
	}
	return int(tx.Status)
}

func buildTxOpts() *bind.TransactOpts {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		log.Fatal("PendingNonceAt failed: ", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("SuggestGasPrice failed: ", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(config.ChainID))
	if err != nil {
		log.Fatal("NewKeyedTransactorWithChainID failed: ", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)              // in wei
	auth.GasLimit = uint64(defaultGasLimit) // in units
	auth.GasPrice = gasPrice

	return auth
}
