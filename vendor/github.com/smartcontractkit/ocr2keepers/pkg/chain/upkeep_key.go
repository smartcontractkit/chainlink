package chain

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

// TODO (AUTO-2014), find a better place for these concrete types than chain package
type UpkeepKey []byte

// NewUpkeepKey is the constructor of UpkeepKey
func NewUpkeepKey(block, id *big.Int) UpkeepKey {
	return UpkeepKey(fmt.Sprintf("%s%s%s", block, separator, id))
}

func NewUpkeepKeyFromBlockAndID(block types.BlockKey, id types.UpkeepIdentifier) UpkeepKey {
	return UpkeepKey(fmt.Sprintf("%s%s%s", block, separator, string(id)))
}

func (u UpkeepKey) BlockKeyAndUpkeepID() (types.BlockKey, types.UpkeepIdentifier, error) {
	components := strings.Split(u.String(), "|")
	if len(components) != 2 {
		return nil, nil, fmt.Errorf("%w: missing data in upkeep key", ErrUpkeepKeyNotParsable)
	}

	return BlockKey(components[0]), types.UpkeepIdentifier(components[1]), nil
}

func (u UpkeepKey) String() string {
	return string(u)
}
