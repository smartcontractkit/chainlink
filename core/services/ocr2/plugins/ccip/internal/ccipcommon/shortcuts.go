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

// GetFilteredSortedChainTokens returns union of all tokens supported on the destination chain, including fee tokens from the provided price registry
// and the bridgeable tokens from all the offRamps living on the chain. Bridgeable tokens are only included if they are configured on the pricegetter
// Fee tokens are not filtered as they must always be priced
func GetFilteredSortedChainTokens(ctx context.Context, offRamps []ccipdata.OffRampReader, priceRegistry cciptypes.PriceRegistryReader, priceGetter cciptypes.PriceGetter) (chainTokens []cciptypes.Address, excludedTokens []cciptypes.Address, err error) {
	destFeeTokens, destBridgeableTokens, err := getTokensWithBatchLimit(ctx, offRamps, priceRegistry, offRampBatchSizeLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("get tokens with batch limit: %w", err)
	}

	destTokensWithPrice, destTokensWithoutPrice, err := priceGetter.FilterConfiguredTokens(ctx, destBridgeableTokens)
	if err != nil {
		return nil, nil, fmt.Errorf("filter for priced tokens: %w", err)
	}

	return flattenedAndSortedChainTokens(destFeeTokens, destTokensWithPrice), destTokensWithoutPrice, nil
}

func flattenedAndSortedChainTokens(slices ...[]cciptypes.Address) (chainTokens []cciptypes.Address) {
	// same token can be returned by multiple offRamps, and fee token can overlap with bridgeable tokens,
	// we need to dedup them to arrive at chain token set
	chainTokens = FlattenUniqueSlice(slices...)

	// return the tokens in deterministic order to aid with testing and debugging
	sort.Slice(chainTokens, func(i, j int) bool {
		return chainTokens[i] < chainTokens[j]
	})

	return chainTokens
}

func getTokensWithBatchLimit(ctx context.Context, offRamps []ccipdata.OffRampReader, priceRegistry cciptypes.PriceRegistryReader, batchSize int) (destFeeTokens []cciptypes.Address, destBridgeableTokens []cciptypes.Address, err error) {
	if batchSize == 0 {
		return nil, nil, fmt.Errorf("batch size must be greater than 0")
	}

	eg := new(errgroup.Group)
	eg.SetLimit(batchSize)

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
		return nil, nil, err
	}

	// same token can be returned by multiple offRamps
	return destFeeTokens, flattenedAndSortedChainTokens(destBridgeableTokens), nil
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
