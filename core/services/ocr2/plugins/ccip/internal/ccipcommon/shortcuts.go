package ccipcommon

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

const (
	offRampBatchSizeLimit = 30
)

func GetMessageIDsAsHexString(messages []cciptypes.EVM2EVMMessage) []string {
	messageIDs := make([]string, 0, len(messages))
	for _, m := range messages {
		messageIDs = append(messageIDs, "0x"+hex.EncodeToString(m.MessageID[:]))
	}
	return messageIDs
}

type BackfillArgs struct {
	SourceLP, DestLP                 logpoller.LogPoller
	SourceStartBlock, DestStartBlock uint64
}

func GetSortedChainTokens(ctx context.Context, offRamps []ccipdata.OffRampReader, priceRegistry cciptypes.PriceRegistryReader) (chainTokens []cciptypes.Address, err error) {
	return getSortedChainTokensWithBatchLimit(ctx, offRamps, priceRegistry, offRampBatchSizeLimit)
}

// GetChainTokens returns union of all tokens supported on the destination chain, including fee tokens from the provided price registry
// and the bridgeable tokens from all the offRamps living on the chain.
func getSortedChainTokensWithBatchLimit(ctx context.Context, offRamps []ccipdata.OffRampReader, priceRegistry cciptypes.PriceRegistryReader, batchSize int) (chainTokens []cciptypes.Address, err error) {
	if batchSize == 0 {
		return nil, fmt.Errorf("batch size must be greater than 0")
	}

	eg := new(errgroup.Group)
	eg.SetLimit(batchSize)

	var destFeeTokens []cciptypes.Address
	var destBridgeableTokens []cciptypes.Address
	mu := &sync.RWMutex{}

	eg.Go(func() error {
		tokens, err := priceRegistry.GetFeeTokens(ctx)
		if err != nil {
			return fmt.Errorf("get dest fee tokens: %w", err)
		}
		destFeeTokens = tokens
		return nil
	})

	for _, o := range offRamps {
		offRamp := o
		eg.Go(func() error {
			tokens, err := offRamp.GetTokens(ctx)
			if err != nil {
				return fmt.Errorf("get dest bridgeable tokens: %w", err)
			}
			mu.Lock()
			destBridgeableTokens = append(destBridgeableTokens, tokens.DestinationTokens...)
			mu.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// same token can be returned by multiple offRamps, and fee token can overlap with bridgeable tokens,
	// we need to dedup them to arrive at chain token set
	chainTokens = FlattenUniqueSlice(destFeeTokens, destBridgeableTokens)

	// return the tokens in deterministic order to aid with testing and debugging
	sort.Slice(chainTokens, func(i, j int) bool {
		return chainTokens[i] < chainTokens[j]
	})

	return chainTokens, nil
}

// GetDestinationTokens returns the destination chain fee tokens from the provided price registry
// and the bridgeable tokens from the offramp.
func GetDestinationTokens(ctx context.Context, offRamp ccipdata.OffRampReader, priceRegistry cciptypes.PriceRegistryReader) (fee, bridged []cciptypes.Address, err error) {
	eg := new(errgroup.Group)

	var destFeeTokens []cciptypes.Address
	var destBridgeableTokens []cciptypes.Address

	eg.Go(func() error {
		tokens, err := priceRegistry.GetFeeTokens(ctx)
		if err != nil {
			return fmt.Errorf("get dest fee tokens: %w", err)
		}
		destFeeTokens = tokens
		return nil
	})

	eg.Go(func() error {
		tokens, err := offRamp.GetTokens(ctx)
		if err != nil {
			return fmt.Errorf("get dest bridgeable tokens: %w", err)
		}
		destBridgeableTokens = tokens.DestinationTokens
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}

	return destFeeTokens, destBridgeableTokens, nil
}

// FlattenUniqueSlice returns a flattened slice that contains unique elements by preserving their order.
func FlattenUniqueSlice[T comparable](slices ...[]T) []T {
	seen := make(map[T]struct{})
	flattened := make([]T, 0)

	for _, sl := range slices {
		for _, el := range sl {
			if _, exists := seen[el]; !exists {
				flattened = append(flattened, el)
				seen[el] = struct{}{}
			}
		}
	}
	return flattened
}

func IsTxRevertError(err error) bool {
	if err == nil {
		return false
	}

	// Geth eth_call reverts with "execution reverted"
	// Nethermind, Parity, OpenEthereum eth_call reverts with "VM execution error"
	// See: https://github.com/ethereum/go-ethereum/issues/21886
	return strings.Contains(err.Error(), "execution reverted") || strings.Contains(err.Error(), "VM execution error")
}
