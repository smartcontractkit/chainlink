package utils

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func SenderFromTxHex(value string, chainID uint64) (common.Address, error) {
	tx, err := DecodeTxFromHex(value, chainID)
	if err != nil {
		return common.Address{}, err
	}
	signer := types.NewEIP155Signer(big.NewInt(int64(chainID)))
	return types.Sender(signer, &tx)
}

func DecodeTxFromHex(value string, chainID uint64) (types.Transaction, error) {
	buffer := bytes.NewBuffer(common.FromHex(value))
	rlpStream := rlp.NewStream(buffer, 0)
	tx := types.Transaction{}
	err := tx.DecodeRLP(rlpStream)
	return tx, err
}

func HexToUint64(hex string) (uint64, error) {
	if strings.ToLower(hex[0:2]) == "0x" {
		hex = hex[2:]
	}
	return strconv.ParseUint(hex, 16, 64)
}

func Uint64ToHex(i uint64) string {
	return fmt.Sprintf("0x%x", i)
}

func EncodeTxToHex(tx *types.Transaction) (string, error) {
	rlp := new(bytes.Buffer)
	if err := tx.EncodeRLP(rlp); err != nil {
		return "", err
	}
	return common.ToHex(rlp.Bytes()), nil
}

func StringToHash(str string) (common.Hash, error) {
	b, err := hexutil.Decode(str)
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(b), nil
}

func StringToAddress(str string) (common.Address, error) {
	b, err := hexutil.Decode(str)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(b), nil
}
