package common

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_v3_aggregator_contract"
)

type Environment struct {
	Owner *bind.TransactOpts
	Ec    *ethclient.Client

	Jc *rpc.Client

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

	insecureSkipVerify := os.Getenv("INSECURE_SKIP_VERIFY") == "true"
	tr := &http.Transport{
		// User enables this at their own risk!
		// #nosec G402
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	httpClient := &http.Client{Transport: tr}
	rpcConfig := rpc.WithHTTPClient(httpClient)
	jsonRPCClient, err := rpc.DialOptions(context.Background(), ethURL, rpcConfig)
	PanicErr(err)
	ec := ethclient.NewClient(jsonRPCClient)

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
	fmt.Println("Suggested Gas Price:", gp, "wei")
	owner.GasPrice = gp
	gasLimit, set := os.LookupEnv("GAS_LIMIT")
	if set {
		parsedGasLimit, err := strconv.ParseUint(gasLimit, 10, 64)
		if err != nil {
			panic("Failure while parsing GAS_LIMIT: " + gasLimit)
		}
		owner.GasLimit = parsedGasLimit
	}

	if overrideNonce {
		block, err := ec.BlockNumber(context.Background())
		PanicErr(err)

		nonce, err := ec.NonceAt(context.Background(), owner.From, big.NewInt(int64(block)))
		PanicErr(err)

		owner.Nonce = big.NewInt(int64(nonce))
	}
	owner.GasPrice = gp.Mul(gp, big.NewInt(2))
	fmt.Println("Modified Gas Price that will be set:", owner.GasPrice, "wei")
	// the execution environment for the scripts
	return Environment{
		Owner:   owner,
		Ec:      ec,
		Jc:      jsonRPCClient,
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

func explorerLinkPrefix(chainID int64) (prefix string) {
	switch chainID {
	case 1: // ETH mainnet
		prefix = "https://etherscan.io"
	case 4: // Rinkeby
		prefix = "https://rinkeby.etherscan.io"
	case 5: // Goerli
		prefix = "https://goerli.etherscan.io"
	case 42: // Kovan
		prefix = "https://kovan.etherscan.io"
	case 11155111: // Sepolia
		prefix = "https://sepolia.etherscan.io"

	case 420: // Optimism Goerli
		prefix = "https://goerli-optimism.etherscan.io"

	case ArbitrumGoerliChainID: // Arbitrum Goerli
		prefix = "https://goerli.arbiscan.io"
	case ArbitrumOneChainID: // Arbitrum mainnet
		prefix = "https://arbiscan.io"
	case ArbitrumSepoliaChainID: // Arbitrum Sepolia
		prefix = "https://sepolia.arbiscan.io"

	case 56: // BSC mainnet
		prefix = "https://bscscan.com"
	case 97: // BSC testnet
		prefix = "https://testnet.bscscan.com"

	case 137: // Polygon mainnet
		prefix = "https://polygonscan.com"
	case 80001: // Polygon Mumbai testnet
		prefix = "https://mumbai.polygonscan.com"

	case 250: // Fantom mainnet
		prefix = "https://ftmscan.com"
	case 4002: // Fantom testnet
		prefix = "https://testnet.ftmscan.com"

	case 43114: // Avalanche mainnet
		prefix = "https://snowtrace.io"
	case 43113: // Avalanche testnet
		prefix = "https://testnet.snowtrace.io"
	case 335: // Defi Kingdoms testnet
		prefix = "https://subnets-test.avax.network/defi-kingdoms"
	case 53935: // Defi Kingdoms mainnet
		prefix = "https://subnets.avax.network/defi-kingdoms"

	case 1666600000, 1666600001, 1666600002, 1666600003: // Harmony mainnet
		prefix = "https://explorer.harmony.one"
	case 1666700000, 1666700001, 1666700002, 1666700003: // Harmony testnet
		prefix = "https://explorer.testnet.harmony.one"

	case 84531:
		prefix = "https://goerli.basescan.org"
	case 8453:
		prefix = "https://basescan.org"

	case 280: // zkSync Goerli testnet
		prefix = "https://goerli.explorer.zksync.io"
	case 324: // zkSync mainnet
		prefix = "https://explorer.zksync.io"

	default: // Unknown chain, return prefix as-is
		prefix = ""
	}
	return
}

func automationExplorerNetworkName(chainID int64) (prefix string) {
	switch chainID {
	case 1: // ETH mainnet
		prefix = "mainnet"
	case 5: // Goerli
		prefix = "goerli"
	case 11155111: // Sepolia
		prefix = "sepolia"

	case 420: // Optimism Goerli
		prefix = "optimism-goerli"

	case ArbitrumGoerliChainID: // Arbitrum Goerli
		prefix = "arbitrum-goerli"
	case ArbitrumOneChainID: // Arbitrum mainnet
		prefix = "arbitrum"
	case ArbitrumSepoliaChainID: // Arbitrum Sepolia
		prefix = "arbitrum-sepolia"

	case 56: // BSC mainnet
		prefix = "bsc"
	case 97: // BSC testnet
		prefix = "bnb-chain-testnet"

	case 137: // Polygon mainnet
		prefix = "polygon"
	case 80001: // Polygon Mumbai testnet
		prefix = "mumbai"

	case 250: // Fantom mainnet
		prefix = "fantom"
	case 4002: // Fantom testnet
		prefix = "fantom-testnet"

	case 43114: // Avalanche mainnet
		prefix = "avalanche"
	case 43113: // Avalanche testnet
		prefix = "fuji"

	default: // Unknown chain, return prefix as-is
		prefix = "<NOT IMPLEMENTED>"
	}
	return
}

// ExplorerLink creates a block explorer link for the given transaction hash. If the chain ID is
// unrecognized, the hash is returned as-is.
func ExplorerLink(chainID int64, txHash common.Hash) string {
	prefix := explorerLinkPrefix(chainID)
	if prefix != "" {
		return fmt.Sprintf("%s/tx/%s", prefix, txHash.String())
	}
	return txHash.String()
}

// ContractExplorerLink creates a block explorer link for the given contract address.
// If the chain ID is unrecognized the address is returned as-is.
func ContractExplorerLink(chainID int64, contractAddress common.Address) string {
	prefix := explorerLinkPrefix(chainID)
	if prefix != "" {
		return fmt.Sprintf("%s/address/%s", prefix, contractAddress.Hex())
	}
	return contractAddress.Hex()
}

func TenderlySimLink(simID string) string {
	return "https://dashboard.tenderly.co/simulator/" + simID
}

// ConfirmTXMined confirms that the given transaction is mined and prints useful execution information.
func ConfirmTXMined(context context.Context, client *ethclient.Client, transaction *types.Transaction, chainID int64, txInfo ...string) (receipt *types.Receipt) {
	if transaction == nil {
		fmt.Println("No transaction to confirm")
		return
	}

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
	fmt.Println("Contract explorer link:", ContractExplorerLink(chainID, contractAddress))
	return contractAddress
}

func ConfirmCodeAt(ctx context.Context, client *ethclient.Client, addr common.Address, chainID int64) {
	fmt.Println("Confirming contract deployment:", addr)
	timeout := time.After(time.Minute)
	for {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("getting code at", addr)
			code, err := client.CodeAt(ctx, addr, nil)
			PanicErr(err)
			if len(code) > 0 {
				fmt.Println("contract deployment confirmed:", ContractExplorerLink(chainID, addr))
				return
			}
		case <-timeout:
			fmt.Println("Could not confirm contract deployment:", addr)
			return
		}
	}
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

func ParseHexSlice(arg string) (ret [][]byte) {
	parts := strings.Split(arg, ",")
	for _, part := range parts {
		ret = append(ret, hexutil.MustDecode(part))
	}
	return
}

func FundNodes(e Environment, transmitters []string, fundingAmount *big.Int) {
	for _, transmitter := range transmitters {
		FundNode(e, transmitter, fundingAmount)
	}
}

func FundNode(e Environment, address string, fundingAmount *big.Int) {
	block, err := e.Ec.BlockNumber(context.Background())
	PanicErr(err)

	nonce, err := e.Ec.NonceAt(context.Background(), e.Owner.From, big.NewInt(int64(block)))
	PanicErr(err)
	// Special case for Arbitrum since gas estimation there is different.

	var gasLimit uint64
	if IsArbitrumChainID(e.ChainID) {
		to := common.HexToAddress(address)
		estimated, err2 := e.Ec.EstimateGas(context.Background(), ethereum.CallMsg{
			From:  e.Owner.From,
			To:    &to,
			Value: fundingAmount,
		})
		PanicErr(err2)
		gasLimit = estimated
	} else {
		gasLimit = uint64(21_000)
	}
	toAddress := common.HexToAddress(address)

	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: e.Owner.GasPrice,
			Gas:      gasLimit,
			To:       &toAddress,
			Value:    fundingAmount,
			Data:     nil,
		})

	signedTx, err := e.Owner.Signer(e.Owner.From, tx)
	PanicErr(err)
	err = e.Ec.SendTransaction(context.Background(), signedTx)
	PanicErr(err)
	fmt.Printf("Sending to %s: %s\n", address, ExplorerLink(e.ChainID, signedTx.Hash()))
	PanicErr(err)
	_, err = bind.WaitMined(context.Background(), e.Ec, signedTx)
	PanicErr(err)
}

// binarySearch finds the highest value within the range bottom-top at which the test function is
// true.
func BinarySearch(top, bottom *big.Int, test func(amount *big.Int) bool) *big.Int {
	var runs int
	// While the difference between top and bottom is > 1
	for new(big.Int).Sub(top, bottom).Cmp(big.NewInt(1)) > 0 {
		// Calculate midpoint between top and bottom
		midpoint := new(big.Int).Sub(top, bottom)
		midpoint.Div(midpoint, big.NewInt(2))
		midpoint.Add(midpoint, bottom)

		// Check if the midpoint amount is withdrawable
		if test(midpoint) {
			bottom = midpoint
		} else {
			top = midpoint
		}

		runs++
		if runs%10 == 0 {
			fmt.Printf("Searching... current range %s-%s\n", bottom.String(), top.String())
		}
	}

	return bottom
}

// GetRlpHeaders gets RLP encoded headers of a list of block numbers
// Makes RPC network call eth_getBlockByNumber to blockchain RPC node
// to fetch header info
func GetRlpHeaders(env Environment, blockNumbers []*big.Int, getParentBlocks bool) (headers [][]byte, hashes []string, err error) {
	hashes = make([]string, 0)

	offset := big.NewInt(0)
	if getParentBlocks {
		offset = big.NewInt(1)
	}

	headers = [][]byte{}
	var rlpHeader []byte
	for _, blockNum := range blockNumbers {
		// Avalanche block headers are special, handle them by using the avalanche rpc client
		// rather than the regular go-ethereum ethclient.
		if IsAvaxNetwork(env.ChainID) {
			var h AvaHeader
			// Get child block since it's the one that has the parent hash in its header.
			nextBlockNum := new(big.Int).Set(blockNum).Add(blockNum, offset)
			err2 := env.Jc.CallContext(context.Background(), &h, "eth_getBlockByNumber", hexutil.EncodeBig(nextBlockNum), false)
			if err2 != nil {
				return nil, hashes, fmt.Errorf("failed to get header: %+w", err2)
			}
			// We can still use vanilla go-ethereum rlp.EncodeToBytes, see e.g
			// https://github.com/ava-labs/coreth/blob/e3ca41bf5295a9a7ca1aeaf29d541fcbb94f79b1/core/types/hashing.go#L49-L57.
			rlpHeader, err2 = rlp.EncodeToBytes(h)
			if err2 != nil {
				return nil, hashes, fmt.Errorf("failed to encode rlp: %+w", err2)
			}

			hashes = append(hashes, h.Hash().String())

			// Sanity check - can be un-commented if storeVerifyHeader is failing due to unexpected
			// blockhash.
			// bh := crypto.Keccak256Hash(rlpHeader)
			// fmt.Println("Calculated BH:", bh.String(),
			//	"fetched BH:", h.Hash(),
			//	"block number:", new(big.Int).Set(blockNum).Add(blockNum, offset).String())
		} else if IsAvaxSubnet(env.ChainID) {
			var h AvaSubnetHeader
			// Get child block since it's the one that has the parent hash in its header.
			nextBlockNum := new(big.Int).Set(blockNum).Add(blockNum, offset)
			err2 := env.Jc.CallContext(context.Background(), &h, "eth_getBlockByNumber", hexutil.EncodeBig(nextBlockNum), false)
			if err2 != nil {
				return nil, hashes, fmt.Errorf("failed to get header: %w", err2)
			}
			rlpHeader, err2 = rlp.EncodeToBytes(h)
			if err2 != nil {
				return nil, hashes, fmt.Errorf("failed to encode rlp: %w", err2)
			}

			hashes = append(hashes, h.Hash().String())
		} else if IsPolygonEdgeNetwork(env.ChainID) {
			// Get child block since it's the one that has the parent hash in its header.
			nextBlockNum := new(big.Int).Set(blockNum).Add(blockNum, offset)
			var hash string
			rlpHeader, hash, err = GetPolygonEdgeRLPHeader(env.Jc, nextBlockNum)
			if err != nil {
				return nil, hashes, fmt.Errorf("failed to encode rlp: %w", err)
			}

			hashes = append(hashes, hash)
		} else if IsRoninChain(env.ChainID) {
			var h RoninHeader
			// Get child block since it's the one that has the parent hash in its header.
			nextBlockNum := new(big.Int).Set(blockNum).Add(blockNum, offset)
			err2 := env.Jc.CallContext(context.Background(), &h, "eth_getBlockByNumber", hexutil.EncodeBig(nextBlockNum), false)
			if err2 != nil {
				return nil, hashes, fmt.Errorf("failed to get header: %w", err2)
			}
			rlpHeader, err2 = rlp.EncodeToBytes(h)
			if err2 != nil {
				return nil, hashes, fmt.Errorf("failed to encode rlp: %w", err2)
			}

			hashes = append(hashes, h.Hash().String())
		} else {
			// Get child block since it's the one that has the parent hash in its header.
			h, err2 := env.Ec.HeaderByNumber(
				context.Background(),
				new(big.Int).Set(blockNum).Add(blockNum, offset),
			)
			if err2 != nil {
				return nil, hashes, fmt.Errorf("failed to get header: %w", err2)
			}
			rlpHeader, err2 = rlp.EncodeToBytes(h)
			if err2 != nil {
				return nil, hashes, fmt.Errorf("failed to encode rlp: %w", err2)
			}

			hashes = append(hashes, h.Hash().String())
		}

		headers = append(headers, rlpHeader)
	}
	return
}

// IsPolygonEdgeNetwork returns true if the given chain ID corresponds to an Pologyon Edge network.
func IsPolygonEdgeNetwork(chainID int64) bool {
	return chainID == 100 || // Nexon test supernet
		chainID == 500 // Nexon test supernet
}

func CalculateLatestBlockHeader(env Environment, blockNumberInput int) (err error) {
	blockNumber := uint64(blockNumberInput)
	if blockNumberInput == -1 {
		blockNumber, err = env.Ec.BlockNumber(context.Background())
		if err != nil {
			return fmt.Errorf("failed to fetch latest block: %w", err)
		}
	}

	// GetRLPHeaders method increments the blockNum sent by 1 and then fetches
	// block headers for the child block.
	blockNumber = blockNumber - 1

	blockNumberBigInts := []*big.Int{big.NewInt(int64(blockNumber))}
	headers, hashes, err := GetRlpHeaders(env, blockNumberBigInts, true)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rlpHeader := headers[0]
	bh := crypto.Keccak256Hash(rlpHeader)
	fmt.Println("Calculated BH:", bh.String(),
		"\nfetched BH:", hashes[0],
		"\nRLP encoding of header: ", hex.EncodeToString(rlpHeader), ", len: ", len(rlpHeader),
		"\nblock number:", new(big.Int).Set(blockNumberBigInts[0]).Add(blockNumberBigInts[0], big.NewInt(1)).String(),
		fmt.Sprintf("\nblock number hex: 0x%x\n", blockNumber+1))

	return err
}

// IsAvaxNetwork returns true if the given chain ID corresponds to an avalanche network.
func IsAvaxNetwork(chainID int64) bool {
	return chainID == 43114 || // C-chain mainnet
		chainID == 43113 // Fuji testnet
}

// IsAvaxSubnet returns true if the given chain ID corresponds to an avalanche subnet.
func IsAvaxSubnet(chainID int64) bool {
	return chainID == 335 || // DFK testnet
		chainID == 53935 || // DFK mainnet
		chainID == 5668 || // Nexon Dev
		chainID == 595581 || // Nexon Test
		chainID == 807424 || // Nexon QA
		chainID == 847799 || // Nexon Stage
		chainID == 60118 || // Nexon Mainnet (Actually a testnet)
		chainID == 68414 // Nexon Henesys Mainnet
}

func IsRoninChain(chainID int64) bool {
	return chainID == 2020 || // Ronin Mainnet
		chainID == 2021 // Ronin Saigon testnet
}

func UpkeepLink(chainID int64, upkeepID *big.Int) string {
	return fmt.Sprintf("https://automation.chain.link/%s/%s", automationExplorerNetworkName(chainID), upkeepID.String())
}
