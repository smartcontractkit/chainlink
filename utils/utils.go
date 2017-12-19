package utils

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func SenderFromTxHex(value string, chainID int64) (common.Address, error) {
	tx, err := DecodeTxFromHex(value, chainID)
	if err != nil {
		return common.Address{}, err
	}
	signer := types.NewEIP155Signer(big.NewInt(chainID))
	return types.Sender(signer, &tx)
}

func DecodeTxFromHex(value string, chainID int64) (types.Transaction, error) {
	buffer := bytes.NewBuffer(common.FromHex(value))
	rlpStream := rlp.NewStream(buffer, 0)
	tx := types.Transaction{}
	err := tx.DecodeRLP(rlpStream)
	return tx, err
}
