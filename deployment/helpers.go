package deployment

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// OCRSecrets are used to disseminate a shared secret to OCR nodes
// through the blockchain where OCR configuration is stored. Its a low value secret used
// to derive transmission order etc. They are extracted here such that they can common
// across signers when multiple signers are signing the same OCR config.
type OCRSecrets struct {
	SharedSecret [16]byte
	EphemeralSk  [32]byte
}

func (s OCRSecrets) IsEmpty() bool {
	return s.SharedSecret == [16]byte{} || s.EphemeralSk == [32]byte{}
}

func XXXGenerateTestOCRSecrets() OCRSecrets {
	var s OCRSecrets
	copy(s.SharedSecret[:], crypto.Keccak256([]byte("shared"))[:16])
	copy(s.EphemeralSk[:], crypto.Keccak256([]byte("ephemeral")))
	return s
}

// SimTransactOpts is useful to generate just the calldata for a given gethwrapper method.
func SimTransactOpts() *bind.TransactOpts {
	return &bind.TransactOpts{Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
		return transaction, nil
	}, From: common.HexToAddress("0x0"), NoSend: true, GasLimit: 1_000_000}
}

func GetErrorReasonFromTx(client bind.ContractBackend, from common.Address, tx types.Transaction, receipt *types.Receipt) (string, error) {
	call := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Data:     tx.Data(),
		Value:    tx.Value(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
	}
	_, err := client.CallContract(context.Background(), call, receipt.BlockNumber)
	if err != nil {
		errorReason, err := parseError(err)
		if err == nil {
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
		return "", errors.Errorf("please use an archive node")
	}

	return callErr.Data, nil
}

func ParseErrorFromABI(errorString string, contractABI string) (string, error) {
	parsedAbi, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return "", errors.Wrap(err, "error loading ABI")
	}
	errorString = strings.TrimPrefix(errorString, "Reverted ")
	errorString = strings.TrimPrefix(errorString, "0x")

	data, err := hex.DecodeString(errorString)
	if err != nil {
		return "", errors.Wrap(err, "error decoding error string")
	}
	for errorName, abiError := range parsedAbi.Errors {
		if bytes.Equal(data[:4], abiError.ID.Bytes()[:4]) {
			// Found a matching error
			v, err3 := abiError.Unpack(data)
			if err3 != nil {
				return "", errors.Wrap(err3, "error unpacking data")
			}
			return fmt.Sprintf("error is \"%v\" args %v\n", errorName, v), nil
		}
	}
	return "", errors.New("error not found in ABI")
}
