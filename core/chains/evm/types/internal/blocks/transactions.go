package blocks

import (
	"bytes"
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

type TxType uint8

// NOTE: Need to roll our own unmarshaller since geth's hexutil.Uint64 does not
// handle double zeroes e.g. 0x00
func (txt *TxType) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte(`"0x00"`)) {
		data = []byte(`"0x0"`)
	}
	var hx hexutil.Uint64
	if err := (&hx).UnmarshalJSON(data); err != nil {
		return err
	}
	if hx > math.MaxUint8 {
		return errors.Errorf("expected 'type' to fit into a single byte, got: '%s'", data)
	}
	*txt = TxType(hx)
	return nil
}

func (txt *TxType) MarshalText() ([]byte, error) {
	hx := (hexutil.Uint64)(*txt)
	return hx.MarshalText()
}

// TransactionInternal is JSON-serialization optimized intermediate representation between EVM blocks
// and our public representation
type TransactionInternal struct {
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Type                 *TxType         `json:"type"`
	Hash                 common.Hash     `json:"hash"`
}
