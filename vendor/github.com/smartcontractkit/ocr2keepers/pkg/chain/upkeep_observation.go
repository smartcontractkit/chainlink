package chain

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

var (
	ErrInvalidBlockKey         = fmt.Errorf("invalid block key")
	ErrInvalidUpkeepIdentifier = fmt.Errorf("invalid upkeep identifier")
)

// TODO (AUTO-2014), find a better place for these concrete types than chain package
type UpkeepObservation struct {
	BlockKey          BlockKey                 `json:"1"`
	UpkeepIdentifiers []types.UpkeepIdentifier `json:"2"`
}

type upkeepObservation UpkeepObservation

func (u *UpkeepObservation) UnmarshalJSON(b []byte) error {
	var upkeep upkeepObservation
	if err := json.Unmarshal(b, &upkeep); err != nil {
		return err
	}

	if err := u.validate(upkeep); err != nil {
		return err
	}

	*u = UpkeepObservation(upkeep)

	return nil
}

func (u *UpkeepObservation) validate(uo upkeepObservation) error {
	maxBlockNumber := new(big.Int)
	maxBlockNumber, _ = maxBlockNumber.SetString("18446744073709551615", 10) // 2 ** 64 -1
	maxUpkeepIdentifer := new(big.Int)
	maxUpkeepIdentifer, _ = maxUpkeepIdentifer.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10) // 2 ** 256 -1

	bl, ok := big.NewInt(0).SetString(uo.BlockKey.String(), 10)
	if !ok {
		return ErrBlockKeyNotParsable
	}
	if bl.String() != uo.BlockKey.String() {
		return ErrInvalidBlockKey
	}

	if bl.Cmp(big.NewInt(0)) <= 0 || bl.Cmp(maxBlockNumber) > 0 {
		// Block number should be positive and not greater than 2**64 -1
		return ErrInvalidBlockKey
	}

	for _, ui := range uo.UpkeepIdentifiers {
		uiInt, ok := ui.BigInt()
		if !ok {
			return ErrUpkeepKeyNotParsable
		}
		if uiInt.String() != string(ui) {
			return ErrInvalidUpkeepIdentifier
		}
		if uiInt.Cmp(big.NewInt(0)) == -1 || uiInt.Cmp(maxUpkeepIdentifer) > 0 {
			// UpkeepId should be non negative and not greater than 2**256 - 1
			return ErrInvalidUpkeepIdentifier
		}
	}

	return nil
}
