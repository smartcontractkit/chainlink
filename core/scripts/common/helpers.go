package common

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

type Environment struct {
	Owner   *bind.TransactOpts
	Ec      *ethclient.Client
	ChainID int64
}

func SetupEnv() Environment {
	ethURL, set := os.LookupEnv("ETH_URL")
	if !set {
		panic("need eth url")
	}

	chainIDEnv, set := os.LookupEnv("ETH_CHAIN_ID")
	if !set {
		panic("need chain ID")
	}

	accountKey, set := os.LookupEnv("ACCOUNT_KEY")
	if !set {
		panic("need account key")
	}

	ec, err := ethclient.Dial(ethURL)
	PanicErr(err)

	chainID, err := strconv.ParseInt(chainIDEnv, 10, 64)
	PanicErr(err)

	// Owner key. Make sure it has eth
	b, err := hex.DecodeString(accountKey)
	PanicErr(err)
	d := new(big.Int).SetBytes(b)

	pkX, pkY := crypto.S256().ScalarBaseMult(d.Bytes())
	privateKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: crypto.S256(),
			X:     pkX,
			Y:     pkY,
		},
		D: d,
	}
	owner, err := bind.NewKeyedTransactorWithChainID(&privateKey, big.NewInt(chainID))
	PanicErr(err)
	// Explicitly set gas price to ensure non-eip 1559
	gp, err := ec.SuggestGasPrice(context.Background())
	PanicErr(err)
	owner.GasPrice = gp
	gasLimit, set := os.LookupEnv("GAS_LIMIT")
	if set {
		parsedGasLimit, err := strconv.ParseUint(gasLimit, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Failure while parsing GAS_LIMIT: %s", gasLimit))
		}
		owner.GasLimit = parsedGasLimit
	}

	// the execution environment for the scripts
	return Environment{owner, ec, chainID}

	// Uncomment the block below if transactions are not getting picked up due to nonce issues:
	//
	//block, err := ec.BlockNumber(context.Background())
	//helpers.PanicErr(err)
	//
	//nonce, err := ec.NonceAt(context.Background(), owner.From, big.NewInt(int64(block)))
	//helpers.PanicErr(err)
	//
	//owner.Nonce = big.NewInt(int64(nonce))
	//owner.GasPrice = gp.Mul(gp, big.NewInt(2))
}

// PanicErr panic if error detected
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ParseArgs parses arguments and ensures required args are set
func ParseArgs(flagSet *flag.FlagSet, args []string, requiredArgs ...string) {
	PanicErr(flagSet.Parse(args))
	seen := map[string]bool{}
	argValues := map[string]string{}
	flagSet.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
		argValues[f.Name] = f.Value.String()
	})
	for _, req := range requiredArgs {
		if !seen[req] {
			panic(fmt.Errorf("missing required -%s argument/flag", req))
		}
	}
}

// ExplorerLink creates a block explorer link for the given transaction hash. If the chain ID is
// unrecognized, the hash is returned as-is.
func ExplorerLink(chainID int64, txHash common.Hash) string {
	var fmtURL string
	switch chainID {
	case 1: // ETH mainnet
		fmtURL = "https://etherscan.io/tx/%s"
	case 4: // Rinkeby
		fmtURL = "https://rinkeby.etherscan.io/tx/%s"
	case 42: // Kovan
		fmtURL = "https://kovan.etherscan.io/tx/%s"

	case 56: // BSC mainnet
		fmtURL = "https://bscscan.com/tx/%s"
	case 97: // BSC testnet
		fmtURL = "https://testnet.bscscan.com/tx/%s"

	case 137: // Polygon mainnet
		fmtURL = "https://polygonscan.com/tx/%s"
	case 80001: // Polygon Mumbai testnet
		fmtURL = "https://mumbai.polygonscan.com/tx/%s"

	case 250: // Fantom mainnet
		fmtURL = "https://ftmscan.com/tx/%s"
	case 4002: // Fantom testnet
		fmtURL = "https://testnet.ftmscan.com/tx/%s"

	case 43114: // Avalanche mainnet
		fmtURL = "https://snowtrace.io/tx/%s"
	case 43113: // Avalanche testnet
		fmtURL = "https://testnet.snowtrace.io/tx/%s"

	case 1666600000, 1666600001, 1666600002, 1666600003: // Harmony mainnet
		fmtURL = "https://explorer.harmony.one/tx/%s"
	case 1666700000, 1666700001, 1666700002, 1666700003: // Harmony testnet
		fmtURL = "https://explorer.testnet.harmony.one/tx/%s"

	default: // Unknown chain, return TX as-is
		fmtURL = "%s"
	}

	return fmt.Sprintf(fmtURL, txHash.String())
}

func ConfirmTXMined(context context.Context, client *ethclient.Client, transaction *types.Transaction, chainID int64, txInfo ...string) {
	fmt.Println("Executing TX", ExplorerLink(chainID, transaction.Hash()), txInfo)
	receipt, err := bind.WaitMined(context, client, transaction)
	PanicErr(err)
	fmt.Println("TX", receipt.TxHash, "mined. \nBlock Number:", receipt.BlockNumber, "\nGas Used: ", receipt.GasUsed)
}

func ConfirmContractDeployed(context context.Context, client *ethclient.Client, transaction *types.Transaction, chainID int64) (address common.Address) {
	fmt.Println("Executing contract deployment, TX:", ExplorerLink(chainID, transaction.Hash()))
	contractAddress, err := bind.WaitDeployed(context, client, transaction)
	PanicErr(err)
	fmt.Println("Contract Address:", contractAddress.String())
	return contractAddress
}

func ParseBigIntSlice(arg string) (ret []*big.Int) {
	parts := strings.Split(arg, ",")
	ret = []*big.Int{}
	for _, part := range parts {
		ret = append(ret, decimal.RequireFromString(part).BigInt())
	}
	return ret
}

func ParseIntSlice(arg string) (ret []int) {
	parts := strings.Split(arg, ",")
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		PanicErr(err)
		ret = append(ret, num)
	}
	return ret
}

func ParseAddressSlice(arg string) (ret []common.Address) {
	parts := strings.Split(arg, ",")
	ret = []common.Address{}
	for _, part := range parts {
		ret = append(ret, common.HexToAddress(part))
	}
	return
}

func ParseHashSlice(arg string) (ret []common.Hash) {
	parts := strings.Split(arg, ",")
	ret = []common.Hash{}
	for _, part := range parts {
		ret = append(ret, common.HexToHash(part))
	}
	return
}
