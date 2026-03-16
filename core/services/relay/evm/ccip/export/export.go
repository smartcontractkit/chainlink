package export

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/ccipdata/factory"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/commitstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/versionfinder"
)

type JSONCommitOffchainConfigV1_2_0 = commitstore.JSONCommitOffchainConfig
type CommitOnchainConfig = commitstore.CommitOnchainConfig

func NewCommitOffchainConfig(
	gasPriceDeviationPPB uint32,
	gasPriceHeartBeat time.Duration,
	tokenPriceDeviationPPB uint32,
	tokenPriceHeartBeat time.Duration,
	inflightCacheExpiry time.Duration,
	priceReportingDisabled bool,
) ccip.CommitOffchainConfig {
	return commitstore.NewCommitOffchainConfig(gasPriceDeviationPPB, gasPriceHeartBeat, tokenPriceDeviationPPB, tokenPriceHeartBeat, inflightCacheExpiry, priceReportingDisabled)
}

func NewCommitStoreReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, address ccip.Address, ec client.Client, lp logpoller.LogPoller, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) (types.CommitStoreReader, error) {
	return factory.NewCommitStoreReader(ctx, lggr, versionFinder, address, ec, lp, feeEstimatorConfig)
}

func CloseCommitStoreReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, address ccip.Address, ec client.Client, lp logpoller.LogPoller, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) error {
	return factory.CloseCommitStoreReader(ctx, lggr, versionFinder, address, ec, lp, feeEstimatorConfig)
}
