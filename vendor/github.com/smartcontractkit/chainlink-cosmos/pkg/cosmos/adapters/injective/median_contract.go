package injective

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	chaintypes "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters/injective/types"
)

var _ median.MedianContract = &CosmosMedianReporter{}

type CosmosMedianReporter struct {
	feedID      string
	queryClient chaintypes.QueryClient
}

func NewCosmosMedianReporter(feedID string, queryClient chaintypes.QueryClient) *CosmosMedianReporter {
	return &CosmosMedianReporter{
		feedID,
		queryClient,
	}
}

func (c *CosmosMedianReporter) LatestTransmissionDetails(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	latestAnswer *big.Int,
	latestTimestamp time.Time,
	err error,
) {
	var resp *chaintypes.QueryLatestTransmissionDetailsResponse
	if resp, err = c.queryClient.LatestTransmissionDetails(ctx, &chaintypes.QueryLatestTransmissionDetailsRequest{
		FeedId: c.feedID,
	}); err != nil {
		return
	}

	if resp.ConfigDigest == nil {
		err = fmt.Errorf("unable to receive config digest for for feedID=%s", c.feedID)
		return
	}

	configDigest = configDigestFromBytes(resp.ConfigDigest)

	if resp.EpochAndRound != nil {
		epoch = uint32(resp.EpochAndRound.Epoch)
		round = uint8(resp.EpochAndRound.Round)
	}

	if resp.Data != nil {
		latestAnswer = resp.Data.Answer.BigInt()
		latestTimestamp = time.Unix(resp.Data.TransmissionTimestamp, 0)
	} else {
		latestAnswer = big.NewInt(0)
	}

	err = nil

	return
}

func (c *CosmosMedianReporter) LatestRoundRequested(
	ctx context.Context,
	lookback time.Duration,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	err error,
) {
	// TODO:
	return
}
