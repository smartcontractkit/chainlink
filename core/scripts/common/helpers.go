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
	"github.com/umbracle/fastrlp"

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

	case RoninChainID:
		prefix = "https://app.roninchain.com"
	case RoninSaigonChainID:
		prefix = "https://saigon-app.roninchain.com"
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
	parts := strings.SplitSeq(arg, ",")
	for part := range parts {
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
	parts := strings.SplitSeq(arg, ",")
	for part := range parts {
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
	headers = make([][]byte, len(blockNumbers))
	hashes = make([]string, len(blockNumbers))

	offset := big.NewInt(0)
	if getParentBlocks {
		offset = big.NewInt(1)
	}

	batchElems := make([]rpc.BatchElem, len(blockNumbers))
	switch {
	case IsAvaxNetwork(env.ChainID):
		return getRlpHeaders[*AvaHeader](env, blockNumbers, offset)
	case IsAvaxSubnet(env.ChainID) &&
		// For some reason, Nexon chain does not work with AvaSubnetHeader
		!IsNexonChain(env.ChainID):
		return getRlpHeaders[*AvaSubnetHeader](env, blockNumbers, offset)
	case IsPolygonEdgeNetwork(env.ChainID):
		hs := make([]*PolygonEdgeHeader, len(blockNumbers))
		for i, blockNum := range blockNumbers {
			// Get child block since it's the one that has the parent hash in its header.
			nextBlockNum := new(big.Int).Set(blockNum).Add(blockNum, offset)
			batchElems[i] = rpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []any{"0x" + nextBlockNum.Text(16), false},
				Result: &hs[i],
			}
		}
		err := batchCallContext(context.Background(), env.Jc, batchElems)
		if err != nil {
			return nil, hashes, fmt.Errorf("failed to get header: %w", err)
		}
		for i, h := range hs {
			ar := &fastrlp.Arena{}
			val, err := MarshalRLPWith(ar, h)
			if err != nil {
				return nil, hashes, fmt.Errorf("failed to encode rlp: %w", err)
			}

			rlpHeader := make([]byte, 0)
			rlpHeader = val.MarshalTo(rlpHeader)

			hashes[i] = h.Hash.String()
			headers[i] = rlpHeader
		}
		return headers, hashes, nil
	case IsRoninChain(env.ChainID):
		return getRlpHeaders[*RoninHeader](env, blockNumbers, offset)
	default:
		return getRlpHeaders[*types.Header](env, blockNumbers, offset)
	}
}

type Hashable interface {
	Hash() common.Hash
}

func getRlpHeaders[HEADER Hashable](env Environment, blockNumbers []*big.Int, offset *big.Int) (headers [][]byte, hashes []string, err error) {
	headers = make([][]byte, len(blockNumbers))
	hashes = make([]string, len(blockNumbers))

	batchElems := make([]rpc.BatchElem, len(blockNumbers))
	hs := make([]HEADER, len(blockNumbers))
	for i, blockNum := range blockNumbers {
		// Get child block since it's the one that has the parent hash in its header.
		nextBlockNum := new(big.Int).Set(blockNum).Add(blockNum, offset)
		batchElems[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []any{hexutil.EncodeBig(nextBlockNum), false},
			Result: &hs[i],
		}
	}
	err = batchCallContext(context.Background(), env.Jc, batchElems)
	if err != nil {
		return nil, hashes, fmt.Errorf("failed to get header: %w", err)
	}
	for i, h := range hs {
		rlpHeader, err := rlp.EncodeToBytes(h)
		if err != nil {
			return nil, hashes, fmt.Errorf("failed to encode rlp: %w", err)
		}
		hashes[i] = h.Hash().String()
		headers[i] = rlpHeader
	}
	return
}

// batchCallContext is a wrapper around rpc.Client.BatchCallContext that deals with RPC node batch size limitations.
func batchCallContext(ctx context.Context, client *rpc.Client, batchElems []rpc.BatchElem) error {
	err := client.BatchCallContext(ctx, batchElems)
	if err != nil {
		// Try again with a batch size of 1
		err = client.BatchCallContext(ctx, batchElems[0:1])
		if err != nil {
			// The error is unlikely to be due to a batch size issue, so return it.
			return err
		}
		// If we got here, the batch size of 1 worked, so we can try to find the maximum batch size that works.
		loBatchSize := 1
		hiBatchSize := len(batchElems)
		for start := 1; start < len(batchElems); {
			batchSize := (hiBatchSize + loBatchSize) / 2
			end := min(start+batchSize, len(batchElems))
			err = client.BatchCallContext(ctx, batchElems[start:end])
			if err != nil {
				hiBatchSize = batchSize
			} else {
				loBatchSize = batchSize
				start += batchSize
			}
		}
	}
	return err
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

func UpkeepLink(chainID int64, upkeepID *big.Int) string {
	return fmt.Sprintf("https://automation.chain.link/%s/%s", automationExplorerNetworkName(chainID), upkeepID.String())
}
