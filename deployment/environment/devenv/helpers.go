package devenv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func getErrorReasonFromTx(client bind.ContractBackend, from common.Address, tx *types.Transaction, receipt *types.Receipt) (string, error) {
	call := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Data:     tx.Data(),
		Value:    tx.Value(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
	}
	_, callContractErr := client.CallContract(context.Background(), call, receipt.BlockNumber)
	if callContractErr != nil {
		errorReason, parsingErr := parseError(callContractErr)
		// If we get no information from parsing the error, we return the original error from CallContract
		if errorReason == "" {
			return callContractErr.Error(), nil
		}
		// If the errorReason exists and we had no issues parsing it, we return it
		if parsingErr == nil {
			return errorReason, nil
		}
	}
	return "", fmt.Errorf("tx %s reverted with no reason", tx.Hash().Hex())
}

func parseError(txError error) (string, error) {
	b, err := json.Marshal(txError)
	if err != nil {
		return "", err
	}
	var callErr struct {
		Code    int
		Data    string `json:"data"`
		Message string `json:"message"`
	}
	if json.Unmarshal(b, &callErr) != nil {
		return "", err
	}

	if callErr.Data == "" && strings.Contains(callErr.Message, "missing trie node") {
		return "", errors.New("please use an archive node")
	}

	return callErr.Data, nil
}
