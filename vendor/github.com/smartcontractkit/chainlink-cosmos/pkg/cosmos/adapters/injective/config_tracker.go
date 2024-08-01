package injective

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tmtypes "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	chaintypes "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters/injective/types"
)

var _ types.ContractConfigTracker = &CosmosModuleConfigTracker{}

type CosmosModuleConfigTracker struct {
	feedID                  string
	injectiveClient         chaintypes.QueryClient
	tendermintServiceClient tmtypes.ServiceClient
}

func NewCosmosModuleConfigTracker(feedID string, queryClient chaintypes.QueryClient, serviceClient tmtypes.ServiceClient) *CosmosModuleConfigTracker {
	return &CosmosModuleConfigTracker{
		feedID:                  feedID,
		injectiveClient:         queryClient,
		tendermintServiceClient: serviceClient,
	}
}

// Notify may optionally emit notification events when the contract's
// configuration changes. This is purely used as an optimization reducing
// the delay between a configuration change and its enactment. Implementors
// who don't care about this may simply return a nil channel.
//
// The returned channel should never be closed.
func (c *CosmosModuleConfigTracker) Notify() <-chan struct{} {
	// TODO: track events from Tendermint WS
	return nil
}

// LatestConfigDetails returns information about the latest configuration,
// but not the configuration itself.
func (c *CosmosModuleConfigTracker) LatestConfigDetails(
	ctx context.Context,
) (
	changedInBlock uint64,
	configDigest types.ConfigDigest,
	err error,
) {
	resp, err := c.injectiveClient.FeedConfigInfo(ctx, &chaintypes.QueryFeedConfigInfoRequest{
		FeedId: c.feedID,
	})
	if err != nil {
		return 0, types.ConfigDigest{}, err
	}

	if resp.FeedConfigInfo == nil {
		err = fmt.Errorf("feed config not found: %s", c.feedID)
		return 0, types.ConfigDigest{}, err
	}

	changedInBlock = uint64(resp.FeedConfigInfo.LatestConfigBlockNumber)
	configDigest = configDigestFromBytes(resp.FeedConfigInfo.LatestConfigDigest)
	return changedInBlock, configDigest, nil
}

// LatestConfig returns the latest configuration.
func (c *CosmosModuleConfigTracker) LatestConfig(
	ctx context.Context,
	changedInBlock uint64,
) (types.ContractConfig, error) {
	resp, err := c.injectiveClient.FeedConfig(ctx, &chaintypes.QueryFeedConfigRequest{
		FeedId: c.feedID,
	})
	if err != nil {
		return types.ContractConfig{}, err
	}

	signers := make([]types.OnchainPublicKey, 0, len(resp.FeedConfig.Signers))
	for _, addr := range resp.FeedConfig.Signers {
		acc := sdk.MustAccAddressFromBech32(addr)
		signers = append(signers, types.OnchainPublicKey(acc.Bytes()))
	}

	transmitters := make([]types.Account, 0, len(resp.FeedConfig.Transmitters))
	for _, addr := range resp.FeedConfig.Transmitters {
		acc := sdk.MustAccAddressFromBech32(addr)
		transmitters = append(transmitters, types.Account(acc.String()))
	}

	config := types.ContractConfig{
		ConfigDigest:          configDigestFromBytes(resp.FeedConfigInfo.LatestConfigDigest),
		ConfigCount:           resp.FeedConfigInfo.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     uint8(resp.FeedConfig.F),
		OnchainConfig:         resp.FeedConfig.OnchainConfig,
		OffchainConfigVersion: resp.FeedConfig.OffchainConfigVersion,
		OffchainConfig:        resp.FeedConfig.OffchainConfig,
	}

	return config, nil
}

// TODO: duplicated from wasm adapter
// LatestBlockHeight returns the height of the most recent block in the chain.
func (c *CosmosModuleConfigTracker) LatestBlockHeight(
	ctx context.Context,
) (
	blockHeight uint64,
	err error,
) {
	b, err := c.tendermintServiceClient.GetLatestBlock(context.Background(), &tmtypes.GetLatestBlockRequest{})
	if err != nil {
		return 0, err
	}
	return uint64(b.SdkBlock.Header.Height), nil
}
