package evm

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

var ErrBlockKeyNotParsable = fmt.Errorf("block identifier not parsable")

type BlockKey string

func (k BlockKey) After(kk types.BlockKey) (bool, error) {
	a, ok := big.NewInt(0).SetString(k.String(), 10)
	if !ok {
		return false, ErrBlockKeyNotParsable
	}

	b, ok := big.NewInt(0).SetString(kk.String(), 10)
	if !ok {
		return false, ErrBlockKeyNotParsable
	}

	gt := a.Cmp(b)
	return gt > 0, nil
}

func (k BlockKey) String() string {
	return string(k)
}
