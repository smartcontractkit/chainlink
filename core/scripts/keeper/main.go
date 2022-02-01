package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
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
	NodeURL        string   `mapstructure:"NODE_URL"`
	ChainID        int64    `mapstructure:"CHAIN_ID"`
	PrivateKey     string   `mapstructure:"PRIVATE_KEY"`
	LinkTokenAddr  string   `mapstructure:"LINK_TOKEN_ADDR"`
	Keepers        []string `mapstructure:"KEEPERS"`
	ApproveAmount  string   `mapstructure:"APPROVE_AMOUNT"`
	AddFundsAmount string   `mapstructure:"ADD_FUNDS_AMOUNT"`
	GasLimit       uint64   `mapstructure:"GAS_LIMIT"`

	// Keeper config
	LinkETHFeedAddr      string `mapstructure:"LINK_ETH_FEED"`
	FastGasFeedAddr      string `mapstructure:"FAST_GAS_FEED"`
	PaymentPremiumPBB    uint32 `mapstructure:"PAYMENT_PREMIUM_PBB"`
	FlatFeeMicroLink     uint32 `mapstructure:"FLAT_FEE_MICRO_LINK"`
	BlockCountPerTurn    int64  `mapstructure:"BLOCK_COUNT_PER_TURN"`
	CheckGasLimit        uint32 `mapstructure:"CHECK_GAS_LIMIT"`
	StalenessSeconds     int64  `mapstructure:"STALENESS_SECONDS"`
	GasCeilingMultiplier uint16 `mapstructure:"GAS_CEILING_MULTIPLIER"`
	FallbackGasPrice     int64  `mapstructure:"FALLBACK_GAS_PRICE"`
	FallbackLinkPrice    int64  `mapstructure:"FALLBACK_LINK_PRICE"`

	// Upkeep Config
	UpkeepTestRange                 int64  `mapstructure:"UPKEEP_TEST_RANGE"`
	UpkeepAverageEligibilityCadence int64  `mapstructure:"UPKEEP_AVERAGE_ELIGIBILITY_CADENCE"`
	UpkeepCheckData                 string `mapstructure:"UPKEEP_CHECK_DATA"`
	UpkeepGasLimit                  uint32 `mapstructure:"UPKEEP_GAS_LIMIT"`
	UpkeepCount                     int64  `mapstructure:"UPKEEP_COUNT"`
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
	viper.SetDefault("GAS_LIMIT", 8000000)
	viper.SetDefault("PAYMENT_PREMIUM_PBB", 200000000)
	viper.SetDefault("FLAT_FEE_MICRO_LINK", 0)
	viper.SetDefault("BLOCK_COUNT_PER_TURN", 1)
	viper.SetDefault("CHECK_GAS_LIMIT", 650000000)
	viper.SetDefault("STALENESS_SECONDS", 90000)
	viper.SetDefault("GAS_CEILING_MULTIPLIER", 3)
	viper.SetDefault("FALLBACK_GAS_PRICE", 10000000000)
	viper.SetDefault("FALLBACK_LINK_PRICE", 200000000000000000)
	viper.SetDefault("UPKEEP_TEST_RANGE", 1)
	viper.SetDefault("UPKEEP_AVERAGE_ELIGIBILITY_CADENCE", 1)
	viper.SetDefault("UPKEEP_CHECK_DATA", "0x00")
	viper.SetDefault("UPKEEP_GAS_LIMIT", 500000)
	viper.SetDefault("UPKEEP_COUNT", 5)

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
	registryAddr, deployKeeperRegistryTx, registryInstance, err := keeper.DeployKeeperRegistry(buildTxOpts(), client,
		common.HexToAddress(config.LinkTokenAddr),
		common.HexToAddress(config.LinkETHFeedAddr),
		common.HexToAddress(config.FastGasFeedAddr),
		config.PaymentPremiumPBB,
		config.FlatFeeMicroLink,
		big.NewInt(config.BlockCountPerTurn),
		config.CheckGasLimit,
		big.NewInt(config.StalenessSeconds),
		config.GasCeilingMultiplier,
		big.NewInt(config.FallbackGasPrice),
		big.NewInt(config.FallbackLinkPrice),
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	waitForTx(deployKeeperRegistryTx.Hash())
	log.Println("KeeperRegistry deployed:", registryAddr.Hex(), "-", deployKeeperRegistryTx.Hash().Hex())

	// Approve keeper registry
	approveRegistryTx, err := linkToken.Approve(buildTxOpts(), registryAddr, approveAmount)
	if err != nil {
		log.Fatal("Approve failed: ", err)
	}
	waitForTx(approveRegistryTx.Hash())
	log.Println("KeeperRegistry approved:", registryAddr.Hex(), "-", approveRegistryTx.Hash().Hex())

	// Deploy Upkeeps
	deployUpkeeps(registryInstance)

	// Set Keepers
	keepers, owners := config.keepers()
	setKeepersTx, err := registryInstance.SetKeepers(buildTxOpts(), keepers, owners)
	if err != nil {
		log.Fatal("SetKeepers failed: ", err)
	}
	waitForTx(setKeepersTx.Hash())
	log.Println("Keepers registered:", setKeepersTx.Hash().Hex())
}

func deployUpkeeps(registryInstance *keeper.KeeperRegistry) {
	fmt.Println()
	for i := int64(0); i < config.UpkeepCount; i++ {
		// Deploy
		upkeepAddr, deployUpkeepTx, _, err := upkeep.DeployUpkeepPerformCounterRestrictive(buildTxOpts(), client,
			big.NewInt(config.UpkeepTestRange), big.NewInt(config.UpkeepAverageEligibilityCadence),
		)
		if err != nil {
			log.Fatal(i, "- DeployAbi failed: ", err)
		}
		waitForTx(deployUpkeepTx.Hash())
		log.Println(i, "- Upkeep deployed:", upkeepAddr.Hex(), "-", deployUpkeepTx.Hash().Hex())

		// Approve
		approveUpkeepTx, err := linkToken.Approve(buildTxOpts(), upkeepAddr, approveAmount)
		if err != nil {
			log.Fatal(i, "- Approve failed: ", err)
		}
		waitForTx(approveUpkeepTx.Hash())
		log.Println(i, "- Upkeep approved:", upkeepAddr.Hex(), "-", approveUpkeepTx.Hash().Hex())

		// Register
		registerUpkeepTx, err := registryInstance.RegisterUpkeep(buildTxOpts(),
			upkeepAddr, config.UpkeepGasLimit, fromAddr, []byte(config.UpkeepCheckData),
		)
		if err != nil {
			log.Fatal(i, "- RegisterUpkeep failed: ", err)
		}
		waitForTx(registerUpkeepTx.Hash())
		log.Println(i, "- Upkeep registered:", upkeepAddr.Hex(), "-", registerUpkeepTx.Hash().Hex())

		// Fund
		addFundsTx, err := registryInstance.AddFunds(buildTxOpts(), big.NewInt(int64(i)), addFundsAmount)
		if err != nil {
			log.Fatal(i, "- AddFunds failed: ", err)
		}
		waitForTx(addFundsTx.Hash())
		log.Println(i, "- Upkeep funded:", upkeepAddr.Hex(), "-", addFundsTx.Hash().Hex())
	}
	fmt.Println()
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
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = config.GasLimit // in units
	auth.GasPrice = gasPrice

	return auth
}
