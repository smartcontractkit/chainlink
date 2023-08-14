package handler

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/router"
)

// RevertReasonFromErrorCodeString attempts to decode an error code string
func (h *BaseHandler) RevertReasonFromErrorCodeString(errorCodeString string) string {
	errorCodeString = strings.TrimPrefix(errorCodeString, "0x")
	return decodeErrorStringFromABI(errorCodeString, getAllABIs())
}

// RevertReasonFromTx attempts to fetch more info on failed TX
func (h *BaseHandler) RevertReasonFromTx(txHash string) string {
	// Need a node URL
	// NOTE: this node needs to run in archive mode
	ethUrl := h.cfg.NodeURL
	if ethUrl == "" {
		panicErr(errors.New("You must define ETH_NODE env variable"))
	}
	requester := h.cfg.FromAddress

	ec, ethErr := ethclient.Dial(ethUrl)
	panicErr(ethErr)
	errorString, _ := getErrorForTx(ec, txHash, requester)
	// Some nodes prepend "Reverted " and we also remove the 0x
	trimmed := strings.TrimPrefix(errorString, "Reverted ")[2:]

	contractABIs := getAllABIs()

	return decodeErrorStringFromABI(trimmed, contractABIs)
}

func decodeErrorStringFromABI(errorString string, contractABIs []string) string {
	builder := strings.Builder{}

	data, err := hex.DecodeString(errorString)
	panicErr(err)

	for _, contractABI := range contractABIs {
		parsedAbi, err2 := abi.JSON(strings.NewReader(contractABI))
		panicErr(err2)

		for k, abiError := range parsedAbi.Errors {
			if bytes.Equal(data[:4], abiError.ID.Bytes()[:4]) {
				// Found a matching error
				v, err3 := abiError.Unpack(data)
				panicErr(err3)
				builder.WriteString(fmt.Sprintf("Error is \"%v\" args %v\n", k, v))
				return builder.String()
			}
		}
	}

	if len(errorString) > 8 && errorString[:8] == "4e487b71" {
		builder.WriteString("Decoded error: Assertion failure\n")
		indicator := errorString[len(errorString)-2:]
		switch indicator {
		case "01":
			builder.WriteString("If you call assert with an argument that evaluates to false.\n")
		case "11":
			builder.WriteString("If an arithmetic operation results in underflow or overflow outside of an unchecked { ... } block.\n")
		case "12":
			builder.WriteString("If you divide or modulo by zero (e.g. 5 / 0 or 23 modulo 0).\n")
		case "21":
			builder.WriteString("If you convert a value that is too big or negative into an enum type.\n")
		case "31":
			builder.WriteString("If you call .pop() on an empty array.\n")
		case "32":
			builder.WriteString("If you access an array, bytesN or an array slice at an out-of-bounds or negative index (i.e. x[i] where i >= x.length or i < 0).\n")
		case "41":
			builder.WriteString("If you allocate too much memory or create an array that is too large.\n")
		case "51":
			builder.WriteString("If you call a zero-initialized variable of internal function type.\n")
		default:
			builder.WriteString(fmt.Sprintf("This is a revert produced by an assertion failure. Exact code not found \"%s\"\n", indicator))
		}
		return builder.String()
	}

	stringErr, err := abi.UnpackRevert(data)
	if err == nil {
		builder.WriteString("String error thrown")
		builder.WriteString(fmt.Sprintf("error: %s", stringErr))
		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("Cannot match error with contract ABI. Error code \"%v\"\n", "trimmed"))
	return builder.String()
}

func getAllABIs() []string {
	return []string{
		arm_contract.ARMContractABI,
		lock_release_token_pool.LockReleaseTokenPoolABI,
		burn_mint_token_pool.BurnMintTokenPoolABI,
		commit_store.CommitStoreABI,
		price_registry.PriceRegistryABI,
		evm_2_evm_onramp.EVM2EVMOnRampABI,
		evm_2_evm_offramp.EVM2EVMOffRampABI,
		router.RouterABI,
	}
}

func getErrorForTx(client *ethclient.Client, txHash string, requester string) (string, common.Address) {
	tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	panicErr(err)
	re, err := client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	panicErr(err)

	call := ethereum.CallMsg{
		From:     common.HexToAddress(requester),
		To:       tx.To(),
		Data:     tx.Data(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
	}
	_, err = client.CallContract(context.Background(), call, re.BlockNumber)
	if err == nil {
		panic("no error calling contract")
	}

	return parseError(err), *tx.To()
}

func parseError(txError error) string {
	b, err := json.Marshal(txError)
	panicErr(err)
	var callErr struct {
		Code    int
		Data    string `json:"data"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(b, &callErr)
	panicErr(err)

	if callErr.Data == "" && strings.Contains(callErr.Message, "missing trie node") {
		panic("Use an archive node")
	}
	return callErr.Data
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
