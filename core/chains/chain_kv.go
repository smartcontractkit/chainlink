package chains

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type ChainsKV[T types.ChainService] struct {
	// note: this is read only after construction so no need for mutex
	chains map[string]T
}

var ErrNoSuchChainID = errors.New("chain id does not exist")

func NewChainsKV[T types.ChainService](cs map[string]T) *ChainsKV[T] {
	return &ChainsKV[T]{
		chains: cs,
	}
}
func (c *ChainsKV[T]) Len() int {
	return len(c.chains)
}

// Get return [ErrNoSuchChainID] if [id] is not found
func (c *ChainsKV[T]) Get(id string) (T, error) {
	var dflt T
	chn, exist := c.chains[id]
	if !exist {
		return dflt, fmt.Errorf("%w: %s", ErrNoSuchChainID, id)
	}
	return chn, nil
}

func (c *ChainsKV[T]) List(ids ...string) ([]T, error) {
	if len(ids) == 0 {
		return c.Slice(), nil
	}

	var (
		result []T
		err    error
	)

	for _, id := range ids {
		chn, exists := c.chains[id]
		if !exists {
			err2 := fmt.Errorf("%w: %s", ErrNoSuchChainID, id)
			err = errors.Join(err, err2)
			continue
		}
		result = append(result, chn)
	}

	return result, err
}

func (c *ChainsKV[T]) Slice() []T {
	return maps.Values(c.chains)
}
