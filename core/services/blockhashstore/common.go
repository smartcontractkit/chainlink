package blockhashstore

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

// Coordinator defines an interface for fetching request and fulfillment metadata from a VRF
// coordinator.
type Coordinator interface {
	// Requests fetches VRF requests that occurred within the specified blocks.
	Requests(ctx context.Context, fromBlock uint64, toBlock uint64) ([]Event, error)

	// Fulfillments fetches VRF fulfillments that occurred since the specified block.
	Fulfillments(ctx context.Context, fromBlock uint64) ([]Event, error)
}

// Event contains metadata about a VRF randomness request or fulfillment.
type Event struct {
	// ID of the relevant VRF request. For a VRF V1 request, this will an encoded 32 byte array.
	// For VRF V2, it will be an integer in string form.
	ID string

	// Block that the request or fulfillment was included in.
	Block uint64
}

// BHS defines an interface for interacting with a BlockhashStore contract.
type BHS interface {
	// Store the hash associated with blockNum.
	Store(ctx context.Context, blockNum uint64) error

	// IsStored checks whether the hash associated with blockNum is already stored.
	IsStored(ctx context.Context, blockNum uint64) (bool, error)

	// StoreEarliest stores the earliest possible blockhash (i.e. block.number - 256)
	StoreEarliest(ctx context.Context) error
}

func GetUnfulfilledBlocksAndRequests(
	ctx context.Context,
	lggr logger.Logger,
	coordinator Coordinator,
	fromBlock, toBlock uint64,
) (map[uint64]map[string]struct{}, error) {
	blockToRequests := make(map[uint64]map[string]struct{})
	requestIDToBlock := make(map[string]uint64)

	reqs, err := coordinator.Requests(ctx, uint64(fromBlock), uint64(toBlock))
	if err != nil {
		lggr.Errorw("Failed to fetch VRF requests",
			"error", err)
		return nil, errors.Wrap(err, "fetching VRF requests")
	}
	for _, req := range reqs {
		if _, ok := blockToRequests[req.Block]; !ok {
			blockToRequests[req.Block] = make(map[string]struct{})
		}
		blockToRequests[req.Block][req.ID] = struct{}{}
		requestIDToBlock[req.ID] = req.Block
	}

	fuls, err := coordinator.Fulfillments(ctx, uint64(fromBlock))
	if err != nil {
		lggr.Errorw("Failed to fetch VRF fulfillments",
			"error", err)
		return nil, errors.Wrap(err, "fetching VRF fulfillments")
	}
	for _, ful := range fuls {
		requestBlock, ok := requestIDToBlock[ful.ID]
		if !ok {
			continue
		}
		delete(blockToRequests[requestBlock], ful.ID)
	}

	return blockToRequests, nil
}

// LimitReqIDs converts a set of request IDs to a slice limited to maxLength.
func LimitReqIDs(reqs map[string]struct{}, maxLength int) []string {
	var reqIDs []string
	for id := range reqs {
		reqIDs = append(reqIDs, id)
		if len(reqIDs) >= maxLength {
			break
		}
	}
	return reqIDs
}

// DecreasingBlockRange creates a contiguous block range starting with
// block `start` (inclusive) and ending at block `end` (inclusive).
func DecreasingBlockRange(start, end *big.Int) (ret []*big.Int, err error) {
	if start.Cmp(end) == -1 {
		return nil, fmt.Errorf("start (%s) must be greater than end (%s)", start.String(), end.String())
	}
	ret = []*big.Int{}
	for i := new(big.Int).Set(start); i.Cmp(end) >= 0; i.Sub(i, big.NewInt(1)) {
		ret = append(ret, new(big.Int).Set(i))
	}
	return
}

// GetSearchWindow returns the search window (fromBlock, toBlock) given the latest block number, wait blocks and lookback blocks
func GetSearchWindow(latestBlock, waitBlocks, lookbackBlocks int) (uint64, uint64) {
	var (
		fromBlock = latestBlock - lookbackBlocks
		toBlock   = latestBlock - waitBlocks
	)

	if fromBlock < 0 {
		fromBlock = 0
	}
	if toBlock < 0 {
		toBlock = 0
	}

	return uint64(fromBlock), uint64(toBlock)
}

// SendingKeys returns a list of sending keys (common.Address) given EIP55 addresses
func SendingKeys(fromAddresses []ethkey.EIP55Address) []common.Address {
	var keys []common.Address
	for _, a := range fromAddresses {
		keys = append(keys, a.Address())
	}
	return keys
}
