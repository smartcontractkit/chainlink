package export

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/commitstore"
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
