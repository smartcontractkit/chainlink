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

	avaxclient "github.com/ava-labs/coreth/ethclient"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/mock_v3_aggregator_contract"
)

type Environment struct {
	Owner *bind.TransactOpts
	Ec    *ethclient.Client

	// AvaxEc is appropriately set if the environment is configured to interact with an avalanche
	// chain. It should be used instead of the regular Ec field because avalanche calculates
	// blockhashes differently and the regular Ec will give consistently incorrect results on basic
	// queries (like e.g eth_blockByNumber).
	AvaxEc  avaxclient.Client
	ChainID int64
}

func DeployLinkToken(e Environment) common.Address {
	_, tx, _, err := link_token_interface.DeployLinkToken(e.Owner, e.Ec)
	PanicErr(err)
	return ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func DeployLinkEthFeed(e Environment, linkAddress string, weiPerUnitLink *big.Int) common.Address {
	_, tx, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			e.Owner, e.Ec, 18, weiPerUnitLink)
	PanicErr(err)
	return ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

// SetupEnv returns an Environment object populated from environment variables.
// If overrideNonce is set to true, the nonce will be set to what is returned
// by NonceAt (rather than the typical PendingNonceAt).
func SetupEnv(overrideNonce bool) Environment {
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

	var avaxClient avaxclient.Client
	if IsAvaxNetwork(chainID) {
		avaxClient, err = avaxclient.Dial(ethURL)
		PanicErr(err)
	}

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

	if overrideNonce {
		block, err := ec.BlockNumber(context.Background())
		PanicErr(err)

		nonce, err := ec.NonceAt(context.Background(), owner.From, big.NewInt(int64(block)))
		PanicErr(err)

		owner.Nonce = big.NewInt(int64(nonce))
		owner.GasPrice = gp.Mul(gp, big.NewInt(2))
	}

	// the execution environment for the scripts
	return Environment{
		Owner:   owner,
		Ec:      ec,
		AvaxEc:  avaxClient,
		ChainID: chainID,
	}
}

// PanicErr panics if error the given error is non-nil.
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ParseArgs parses arguments and ensures required args are set.
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
	case 5: // Goerli
		fmtURL = "https://goerli.etherscan.io/tx/%s"
	case 42: // Kovan
		fmtURL = "https://kovan.etherscan.io/tx/%s"
	case 11155111: // Sepolia
		fmtURL = "https://sepolia.etherscan.io/tx/%s"

	case 420: // Optimism Goerli
		fmtURL = "https://goerli-optimism.etherscan.io/tx/%s"

	case ArbitrumGoerliChainID: // Arbitrum Goerli
		fmtURL = "https://goerli-rollup-explorer.arbitrum.io/tx/%s"

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

// ConfirmTXMined confirms that the given transaction is mined and prints useful execution information.
func ConfirmTXMined(context context.Context, client *ethclient.Client, transaction *types.Transaction, chainID int64, txInfo ...string) (receipt *types.Receipt) {
	fmt.Println("Executing TX", ExplorerLink(chainID, transaction.Hash()), txInfo)
	receipt, err := bind.WaitMined(context, client, transaction)
	PanicErr(err)
	fmt.Println("TX", receipt.TxHash, "mined. \nBlock Number:", receipt.BlockNumber,
		"\nGas Used: ", receipt.GasUsed,
		"\nBlock hash: ", receipt.BlockHash.String())
	return
}

// ConfirmContractDeployed confirms that the given contract deployment transaction completed and prints useful execution information.
func ConfirmContractDeployed(context context.Context, client *ethclient.Client, transaction *types.Transaction, chainID int64) (address common.Address) {
	fmt.Println("Executing contract deployment, TX:", ExplorerLink(chainID, transaction.Hash()))
	contractAddress, err := bind.WaitDeployed(context, client, transaction)
	PanicErr(err)
	fmt.Println("Contract Address:", contractAddress.String())
	return contractAddress
}

// ParseBigIntSlice parses the given comma-separated string of integers into a slice
// of *big.Int objects.
func ParseBigIntSlice(arg string) (ret []*big.Int) {
	parts := strings.Split(arg, ",")
	ret = []*big.Int{}
	for _, part := range parts {
		ret = append(ret, decimal.RequireFromString(part).BigInt())
	}
	return ret
}

// ParseIntSlice parses the given comma-separated string of integers into a slice
// of int.
func ParseIntSlice(arg string) (ret []int) {
	parts := strings.Split(arg, ",")
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		PanicErr(err)
		ret = append(ret, num)
	}
	return ret
}

// ParseAddressSlice parses the given comma-separated string of addresses into a slice
// of common.Address objects.
func ParseAddressSlice(arg string) (ret []common.Address) {
	parts := strings.Split(arg, ",")
	ret = []common.Address{}
	for _, part := range parts {
		ret = append(ret, common.HexToAddress(part))
	}
	return
}

// ParseHashSlice parses the given comma-separated string of hashes into a slice of
// common.Hash objects.
func ParseHashSlice(arg string) (ret []common.Hash) {
	parts := strings.Split(arg, ",")
	ret = []common.Hash{}
	for _, part := range parts {
		ret = append(ret, common.HexToHash(part))
	}
	return
}

func FundNodes(e Environment, transmitters []string, fundingAmount *big.Int) {
	block, err := e.Ec.BlockNumber(context.Background())
	PanicErr(err)

	nonce, err := e.Ec.NonceAt(context.Background(), e.Owner.From, big.NewInt(int64(block)))
	PanicErr(err)

	for i := 0; i < len(transmitters); i++ {
		// Special case for Arbitrum since gas estimation there is different.
		var gasLimit uint64
		if IsArbitrumChainID(e.ChainID) {
			to := common.HexToAddress(transmitters[i])
			estimated, err := e.Ec.EstimateGas(context.Background(), ethereum.CallMsg{
				From:  e.Owner.From,
				To:    &to,
				Value: fundingAmount,
			})
			PanicErr(err)
			gasLimit = estimated
		} else {
			gasLimit = uint64(21_000)
		}

		tx := types.NewTransaction(
			nonce+uint64(i),
			common.HexToAddress(transmitters[i]),
			fundingAmount,
			gasLimit,
			e.Owner.GasPrice,
			nil,
		)
		signedTx, err := e.Owner.Signer(e.Owner.From, tx)
		PanicErr(err)
		err = e.Ec.SendTransaction(context.Background(), signedTx)
		PanicErr(err)

		fmt.Printf("Sending to %s: %s\n", transmitters[i], ExplorerLink(e.ChainID, signedTx.Hash()))
	}
}

// IsAvaxNetwork returns true if the given chain ID corresponds to an avalanche network or subnet.
func IsAvaxNetwork(chainID int64) bool {
	return chainID == 43114 || // C-chain mainnet
		chainID == 43113 // Fuji testnet
}
