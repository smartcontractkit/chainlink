package chains

import (
	"errors"
	"fmt"
	"sync"
)

type ChainsKV[T ChainService] struct {
	mu     sync.Mutex
	chains map[string]T
}

var ErrNoSuchChainID = errors.New("chain id does not exist")

func NewChainsKV[T ChainService]() *ChainsKV[T] {
	return &ChainsKV[T]{
		chains: map[string]T{},
	}
}
func (c *ChainsKV[T]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.chains)
}
func (c *ChainsKV[T]) lazyInit() {
	if c.chains == nil {
		c.chains = make(map[string]T)
	}
}
func (c *ChainsKV[T]) Get(id string) (T, error) {
	var dflt T
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lazyInit()
	chn, exist := c.chains[id]
	if !exist {
		return dflt, fmt.Errorf("%w: %s", ErrNoSuchChainID, id)
	}
	return chn, nil
}

func (c *ChainsKV[T]) Put(id string, chn T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lazyInit()
	c.chains[id] = chn
}

func (c *ChainsKV[T]) List(ids ...string) ([]T, error) {
	if len(ids) == 0 {
		return c.Slice(), nil
	}

	var (
		result []T
		err    error
	)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lazyInit()

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
	var result []T
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lazyInit()
	for _, chn := range c.chains {
		result = append(result, chn)
	}
	return result
}
