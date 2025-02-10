package ocrcommon

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	evmclient "github.com/smartcontractkit/chainlink-integrations/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// ArbitrumBlockTranslator uses Arbitrum's special L1BlockNumber to optimise log lookups
// Performance matters here hence aggressive use of the cache
// We want to minimise fetches because calling eth_getBlockByNumber is
// relatively expensive
type ArbitrumBlockTranslator struct {
	ethClient evmclient.Client
	lggr      logger.Logger
	// l2->l1 cache
	cache   map[int64]int64
	cacheMu sync.RWMutex
	l2Locks utils.KeyedMutex
}

// NewArbitrumBlockTranslator returns a concrete ArbitrumBlockTranslator
func NewArbitrumBlockTranslator(ethClient evmclient.Client, lggr logger.Logger) *ArbitrumBlockTranslator {
	return &ArbitrumBlockTranslator{
		ethClient,
		logger.Named(lggr, "ArbitrumBlockTranslator"),
		make(map[int64]int64),
		sync.RWMutex{},
		utils.KeyedMutex{},
	}
}

// NumberToQueryRange implements BlockTranslator interface
func (a *ArbitrumBlockTranslator) NumberToQueryRange(ctx context.Context, changedInL1Block uint64) (fromBlock *big.Int, toBlock *big.Int) {
	var err error
	fromBlock, toBlock, err = a.BinarySearch(ctx, int64(changedInL1Block))
	if err != nil {
		a.lggr.Warnw("Failed to binary search L2->L1, falling back to slow scan over entire chain", "err", err)
		return big.NewInt(0), nil
	}

	return
}

// BinarySearch uses both cache and RPC calls to find the smallest possible range of L2 block numbers that encompasses the given L1 block number
//
// Imagine as a virtual array of L1 block numbers indexed by L2 block numbers
// L1 values are likely duplicated so it looks something like
// [42, 42, 42, 42, 42, 155, 155, 155, 430, 430, 430, 430, 430, ...]
// Theoretical max difference between L1 values is typically about 5, "worst case" is 6545 but can be arbitrarily high if sequencer is broken
// The returned range of L2s from leftmost thru rightmost represent all possible L2s that correspond to the L1 value we are looking for
// nil can be returned as a rightmost value if the range has no upper bound
func (a *ArbitrumBlockTranslator) BinarySearch(ctx context.Context, targetL1 int64) (l2lowerBound *big.Int, l2upperBound *big.Int, err error) {
	mark := time.Now()
	var n int
	defer func() {
		duration := time.Since(mark)
		a.lggr.Debugw(fmt.Sprintf("BinarySearch completed in %s with %d total lookups", duration, n), "finishedIn", duration, "err", err, "nLookups", n)
	}()
	var h *evmtypes.Head

	// l2lower..l2upper is the inclusive range of L2 block numbers in which
	// transactions that called block.number will return the given L1 block
	// number
	var l2lower int64
	var l2upper int64

	var skipUpperBound bool

	{
		var maybeL2Upper *int64
		l2lower, maybeL2Upper = a.reverseLookup(targetL1)
		if maybeL2Upper != nil {
			l2upper = *maybeL2Upper
		} else {
			// Initial query to get highest L1 and L2 numbers
			h, err = a.ethClient.HeadByNumber(ctx, nil)
			n++
			if err != nil {
				return nil, nil, err
			}
			if h == nil {
				return nil, nil, errors.New("got nil head")
			}
			if !h.L1BlockNumber.Valid {
				return nil, nil, errors.New("head was missing L1 block number")
			}
			currentL1 := h.L1BlockNumber.Int64
			currentL2 := h.Number

			a.cachePut(currentL2, currentL1)

			// NOTE: This case shouldn't ever happen but we ought to handle it in the least broken way possible
			if targetL1 > currentL1 {
				// real upper must always be nil, we can skip the upper limit part of the binary search
				a.lggr.Debugf("BinarySearch target of %d is above current L1 block number of %d, using nil for upper bound", targetL1, currentL1)
				return big.NewInt(currentL2), nil, nil
			} else if targetL1 == currentL1 {
				// NOTE: If the latest seen L2 block corresponds to the target L1
				// block, we have to leave the top end of the range open because future
				// L2 blocks can be produced that would also match
				skipUpperBound = true
			}
			l2upper = currentL2
		}
	}

	a.lggr.Debugf("TRACE: BinarySearch starting search for L2 range wrapping L1 block number %d between bounds [%d, %d]", targetL1, l2lower, l2upper)

	var exactMatch bool

	// LEFT EDGE
	// First, use binary search to find the smallest L2 block number for which L1 >= changedInBlock
	// This L2 block number represents the lower bound on a range of L2s corresponding to this L1
	{
		l2lower, err = search(l2lower, l2upper+1, func(l2 int64) (bool, error) {
			l1, miss, err2 := a.arbL2ToL1(ctx, l2)
			if miss {
				n++
			}
			if err2 != nil {
				return false, err2
			}
			if targetL1 == l1 {
				exactMatch = true
			}
			return l1 >= targetL1, nil
		})
		if err != nil {
			return nil, nil, err
		}
	}

	// RIGHT EDGE
	// Second, use binary search again to find the smallest L2 block number for which L1 > changedInBlock
	// Now we can subtract one to get the largest L2 that corresponds to this L1
	// This can be skipped if we know we are already at the top of the range, and the upper limit will be returned as nil
	if !skipUpperBound {
		var r int64
		r, err = search(l2lower, l2upper+1, func(l2 int64) (bool, error) {
			l1, miss, err2 := a.arbL2ToL1(ctx, l2)
			if miss {
				n++
			}
			if err2 != nil {
				return false, err2
			}
			if targetL1 == l1 {
				exactMatch = true
			}
			return l1 > targetL1, nil
		})
		if err != nil {
			return nil, nil, err
		}
		l2upper = r - 1
		l2upperBound = big.NewInt(l2upper)
	}

	// NOTE: We expect either left or right search to make an exact match, if they don't something has gone badly wrong
	if !exactMatch {
		return nil, nil, errors.Errorf("target L1 block number %d is not represented by any L2 block", targetL1)
	}
	return big.NewInt(l2lower), l2upperBound, nil
}

// reverseLookup takes an l1 and returns lower and upper bounds for an L2 based on cache data
func (a *ArbitrumBlockTranslator) reverseLookup(targetL1 int64) (from int64, to *int64) {
	type val struct {
		l1 int64
		l2 int64
	}
	vals := make([]val, 0)

	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()

	for l2, l1 := range a.cache {
		vals = append(vals, val{l1, l2})
	}

	sort.Slice(vals, func(i, j int) bool { return vals[i].l1 < vals[j].l1 })

	for _, val := range vals {
		if val.l1 < targetL1 {
			from = val.l2
		} else if val.l1 > targetL1 && to == nil {
			// workaround golang footgun; can't take a pointer to val
			l2 := val.l2
			to = &l2
		}
	}
	return
}

func (a *ArbitrumBlockTranslator) arbL2ToL1(ctx context.Context, l2 int64) (l1 int64, cacheMiss bool, err error) {
	// This locking block synchronises access specifically around one l2 number so we never fetch the same data concurrently
	// One thread will wait while the other fetches
	unlock := a.l2Locks.LockInt64(l2)
	defer unlock()

	var exists bool
	if l1, exists = a.cacheGet(l2); exists {
		return l1, false, err
	}

	h, err := a.ethClient.HeadByNumber(ctx, big.NewInt(l2))
	if err != nil {
		return 0, true, err
	}
	if h == nil {
		return 0, true, errors.New("got nil head")
	}
	if !h.L1BlockNumber.Valid {
		return 0, true, errors.New("head was missing L1 block number")
	}
	l1 = h.L1BlockNumber.Int64

	a.cachePut(l2, l1)

	return l1, true, nil
}

func (a *ArbitrumBlockTranslator) cacheGet(l2 int64) (l1 int64, exists bool) {
	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()
	l1, exists = a.cache[l2]
	return
}

func (a *ArbitrumBlockTranslator) cachePut(l2, l1 int64) {
	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()
	a.cache[l2] = l1
}

// stolen from golang standard library and modified for 64-bit ints,
// customisable range and erroring function
// see: https://golang.org/src/sort/search.go
func search(i, j int64, f func(int64) (bool, error)) (int64, error) {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	for i < j {
		h := int64(uint64(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		is, err := f(h)
		if err != nil {
			return 0, err
		}
		if !is {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i, nil
}
