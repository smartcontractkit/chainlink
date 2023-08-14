package rhea

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/secrets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// DefaultGasTipFee is the default gas tip fee of 1 gwei.
var DefaultGasTipFee = big.NewInt(1e9)

// EVMGasSettings specifies the gas configuration for an EVM chain.
type EVMGasSettings struct {
	EIP1559   bool
	GasPrice  *big.Int
	GasTipCap *big.Int
}

type ChainDeploySettings struct {
	DeployARM           bool
	DeployTokenPools    bool
	DeployRouter        bool
	DeployUpgradeRouter bool
	DeployPriceRegistry bool
	DeployedAtBlock     uint64
}

type LaneDeploySettings struct {
	DeployLane         bool
	DeployPingPongDapp bool
	DeployedAtBlock    uint64
}

type Chain string

const (
	// Testnets
	Sepolia        Chain = "ethereum-testnet-sepolia"
	AvaxFuji       Chain = "avalanche-testnet-fuji"
	Goerli         Chain = "ethereum-testnet-goerli"
	PolygonMumbai  Chain = "polygon-testnet-mumbai"
	Quorum         Chain = "quorum-testnet-swift"
	BSCTestnet     Chain = "binance_smart_chain-testnet"
	ArbitrumGoerli Chain = "ethereum-testnet-goerli-arbitrum-1"
	OptimismGoerli Chain = "ethereum-testnet-goerli-optimism-1"
	BASEGoerli     Chain = "ethereum-testnet-goerli-base-1"
	// Mainnets
	Ethereum Chain = "ethereum-mainnet"
	Optimism Chain = "optimism-mainnet"
	Avax     Chain = "avax-mainnet"
	Arbitrum Chain = "arbitrum-mainnet"
	Polygon  Chain = "polygon-mainnet"
	Base     Chain = "base-mainnet"
	BSC      Chain = "binance_smart_chain-mainnet"
)

var evmChainIdToChainSelector = map[uint64]uint64{
	// Testnets
	97:       13264668187771770619, // BSC Testnet
	420:      2664363617261496610,  // Optimism Goerli
	1337:     3379446385462418246,  // Quorem
	43113:    14767482510784806043, // Avax Fuji
	84531:    5790810961207155433,  // BASE Goerli
	80001:    12532609583862916517, // Polygon Mumbai
	421613:   6101244977088475029,  // Arbitrum Goerli
	11155111: 16015286601757825753, // Sepolia
	// Mainnets
	1:     5009297550715157269,  // Ethereum
	10:    3734403246176062136,  // Optimism
	56:    11344663589394136015, // BSC
	137:   4051577828743386545,  // Polygon
	8453:  15971525489660198786, // BASE
	42161: 4949039107694359620,  // Arbitrum
	43114: 6433500567565415381,  // Avalanche

}

func GetCCIPChainSelector(EVMChainId uint64) uint64 {
	selector, ok := evmChainIdToChainSelector[EVMChainId]
	if !ok {
		panic(fmt.Sprintf("no chain selector for %d", EVMChainId))
	}
	return selector
}

type Token string

const (
	LINK       Token = "Link"
	WETH       Token = "WETH"
	WAVAX      Token = "WAVAX"
	WMATIC     Token = "WMATIC"
	WBNB       Token = "WBNB"
	CACHEGOLD  Token = "CACHE.gold"
	InsurAce   Token = "InsurAce"
	ZUSD       Token = "zUSD"
	STEADY     Token = "STEADY"
	SUPER      Token = "SUPER"
	BondToken  Token = "BondToken"
	BankToken  Token = "BankToken"
	SNXUSD     Token = "snxUSD"
	FUGAZIUSDC Token = "FugaziUSDCToken"
	Alongside  Token = "Alongside"
	CCIP_BnM   Token = "CCIP-BnM"
	CCIP_LnM   Token = "clCCIP-LnM"
	A_DC       Token = "A$DC"
	NZ_DC      Token = "NZ$DC"
	SG_DC      Token = "SG$DC"
	BetSwirl   Token = "BETS"
)

func GetAllTokens() []Token {
	return []Token{
		LINK, WETH, WAVAX, WBNB,
		WMATIC, CACHEGOLD,
		InsurAce, ZUSD, STEADY,
		SUPER, BondToken, BankToken,
		SNXUSD, FUGAZIUSDC, Alongside,
		CCIP_BnM, CCIP_LnM, A_DC, NZ_DC, SG_DC, BetSwirl,
	}
}

var tokenSymbols = map[Token]string{
	LINK:       "LINK",
	WETH:       "wETH",
	WAVAX:      "wAVAX",
	WMATIC:     "wMATIC",
	WBNB:       "wBNB",
	CACHEGOLD:  "CGT",
	InsurAce:   "INSUR",
	ZUSD:       "zUSD",
	STEADY:     "Steadefi",
	SUPER:      "SuperDuper",
	BondToken:  "BondToken",
	BankToken:  "BankToken",
	SNXUSD:     "snxUSD",
	FUGAZIUSDC: "FUGAZIUSDC",
	Alongside:  "AMKT",
	CCIP_BnM:   "CCIP-BnM",
	CCIP_LnM:   "clCCIP-LnM",
	A_DC:       "A$DC",
	NZ_DC:      "NZ$DC",
	SG_DC:      "SG$DC",
	BetSwirl:   "BETS",
}

func (token Token) Symbol() string {
	return tokenSymbols[token]
}

var tokenDecimalMultiplier = map[Token]uint8{
	LINK:       18,
	WETH:       18,
	WAVAX:      18,
	WMATIC:     18,
	WBNB:       18,
	CACHEGOLD:  8,
	InsurAce:   18,
	ZUSD:       18,
	STEADY:     18,
	SUPER:      18,
	BondToken:  18,
	BankToken:  18,
	SNXUSD:     18,
	FUGAZIUSDC: 6,
	Alongside:  18,
	CCIP_BnM:   18,
	CCIP_LnM:   18,
	A_DC:       6,
	NZ_DC:      6,
	SG_DC:      6,
	BetSwirl:   18,
}

func (token Token) Decimals() uint8 {
	return tokenDecimalMultiplier[token]
}

// Price is a mapping from a whole Token to dollar with 18 decimals precision
// This means a coin that costs $1 will have a price of 1e18 per whole token
func (token Token) Price() *big.Int {
	// Token prices in $ per whole coin
	var TokenPrices = map[Token]*big.Float{
		LINK:       big.NewFloat(6.5),
		WETH:       big.NewFloat(1800),
		WAVAX:      big.NewFloat(15),
		WMATIC:     big.NewFloat(0.85),
		WBNB:       big.NewFloat(200),
		CACHEGOLD:  big.NewFloat(60),
		InsurAce:   big.NewFloat(0.08),
		ZUSD:       big.NewFloat(1),
		STEADY:     big.NewFloat(1),
		SUPER:      big.NewFloat(1),
		BondToken:  big.NewFloat(1),
		BankToken:  big.NewFloat(1),
		SNXUSD:     big.NewFloat(1),
		FUGAZIUSDC: big.NewFloat(1),
		Alongside:  big.NewFloat(1),
		CCIP_BnM:   big.NewFloat(0.0000000001),
		CCIP_LnM:   big.NewFloat(0.0000000001),
		A_DC:       big.NewFloat(1),
		NZ_DC:      big.NewFloat(1),
		SG_DC:      big.NewFloat(1),
		BetSwirl:   big.NewFloat(1),
	}

	tokenValue := big.NewInt(0)
	new(big.Float).Mul(TokenPrices[token], big.NewFloat(1e18)).Int(tokenValue)

	return tokenValue
}

func (token Token) Multiplier() *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(token.Decimals())), nil)
}

type TokenPoolType string

const (
	LockRelease  TokenPoolType = "lockRelease"
	BurnMint     TokenPoolType = "burnMint"
	Wrapped      TokenPoolType = "wrapped"
	FeeTokenOnly TokenPoolType = "feeTokenOnly"
)

type EVMChainConfig struct {
	EvmChainId  uint64
	GasSettings EVMGasSettings

	SupportedTokens    map[Token]EVMBridgedToken
	FeeTokens          []Token
	WrappedNative      Token
	Router             gethcommon.Address
	UpgradeRouter      gethcommon.Address
	ARM                gethcommon.Address
	ARMProxy           gethcommon.Address
	PriceRegistry      gethcommon.Address
	AllowList          []gethcommon.Address
	ARMConfig          *arm_contract.ARMConfig // Setting ARMConfig to nil will deploy MockARM
	DeploySettings     ChainDeploySettings
	TunableChainValues TunableChainValues
}

type TunableChainValues struct {
	FinalityDepth            uint32
	OptimisticConfirmations  uint32
	BatchGasLimit            uint32
	RelativeBoostPerWaitHour float64
	FeeUpdateHeartBeat       models.Duration
	FeeUpdateDeviationPPB    uint32
	MaxGasPrice              uint64
	InflightCacheExpiry      models.Duration
	RootSnoozeTime           models.Duration
}

type EVMBridgedToken struct {
	ChainId uint64
	Token   gethcommon.Address
	Pool    gethcommon.Address
	TokenPriceType
	Price    *big.Int
	Decimals uint8
	PriceFeed
	TokenPoolType
	PoolAllowList []gethcommon.Address // empty slice indicates allowList is not enabled
}

type TokenPriceType string

const (
	TokenPrices TokenPriceType = "TokenPrices"
	PriceFeeds  TokenPriceType = "PriceFeeds"
)

type PriceFeed struct {
	Aggregator gethcommon.Address
	Multiplier *big.Int
}

type EVMLaneConfig struct {
	OnRamp         gethcommon.Address
	OffRamp        gethcommon.Address
	CommitStore    gethcommon.Address
	PingPongDapp   gethcommon.Address
	DeploySettings LaneDeploySettings
}

type EvmDeploymentConfig struct {
	Owner  *bind.TransactOpts
	Client *ethclient.Client
	Logger logger.Logger

	ChainConfig       EVMChainConfig
	LaneConfig        EVMLaneConfig
	UpgradeLaneConfig EVMLaneConfig
}

type EvmConfig struct {
	Owner       *bind.TransactOpts
	Client      *ethclient.Client
	Logger      logger.Logger
	ChainConfig *EVMChainConfig
}

func (chain *EvmDeploymentConfig) OnlyEvmConfig() EvmConfig {
	return EvmConfig{
		Owner:       chain.Owner,
		Client:      chain.Client,
		Logger:      chain.Logger,
		ChainConfig: &chain.ChainConfig,
	}
}

func (chain *EvmDeploymentConfig) SetupChain(t *testing.T, ownerPrivateKey string) {
	chain.Owner = GetOwner(t, ownerPrivateKey, chain.ChainConfig.EvmChainId, chain.ChainConfig.GasSettings)
	chain.Client = GetClient(t, secrets.GetRPC(chain.ChainConfig.EvmChainId))
	chain.Logger = logger.TestLogger(t).Named(ccip.ChainName(int64(chain.ChainConfig.EvmChainId)))
	chain.Logger.Info("Completed chain setup")
}

func (chain *EvmDeploymentConfig) SetupReadOnlyChain(lggr logger.Logger) error {
	client, err := ethclient.Dial(secrets.GetRPC(chain.ChainConfig.EvmChainId))
	if err != nil {
		return err
	}
	chain.Logger = lggr
	chain.Client = client

	return nil
}

// GetOwner sets the owner user credentials and ensures a GasTipCap is set for the resulting user.
func GetOwner(t *testing.T, ownerPrivateKey string, chainId uint64, gasSettings EVMGasSettings) *bind.TransactOpts {
	ownerKey, err := crypto.HexToECDSA(ownerPrivateKey)
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(ownerKey, big.NewInt(int64(chainId)))
	require.NoError(t, err)
	fmt.Println("--- Owner address ")
	fmt.Println(user.From.Hex())
	SetGasFees(user, gasSettings)

	return user
}

// GetClient dials a given EVM client url and returns the resulting client.
func GetClient(t *testing.T, ethUrl string) *ethclient.Client {
	client, err := ethclient.Dial(ethUrl)
	require.NoError(t, err)
	return client
}

// SetGasFees configures the chain client with the given EVMGasSettings. This method is needed for EIP txs
// to function because of the geth-only tip fee method.
func SetGasFees(owner *bind.TransactOpts, config EVMGasSettings) {
	if config.EIP1559 {
		// to not use geth-only tip fee method when EIP1559 is enabled
		// https://github.com/ethereum/go-ethereum/pull/23484
		owner.GasTipCap = config.GasTipCap
	} else {
		owner.GasPrice = config.GasPrice
	}
}

// GetPricePer1e18Units returns the price, in USD with 18 decimals, per 1e18 of the smallest token denomination.
// For example,
//
//	1 USDC = 1.00 USD per full token, each full token is 1e6 units -> 1 * 1e18 * 1e18 / 1e6 = 1e30
//	1 ETH = 2,000 USD per full token, each full token is 1e18 units -> 2000 * 1e18 * 1e18 / 1e18 = 2_000e18
//	1 LINK = 5.00 USD per full token, each full token is 1e18 units -> 5 * 1e18 * 1e18 / 1e18 = 5e18
func GetPricePer1e18Units(price *big.Int, decimals uint8) *big.Int {
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	return new(big.Int).Quo(new(big.Int).Mul(price, big.NewInt(1e18)), multiplier)
}
