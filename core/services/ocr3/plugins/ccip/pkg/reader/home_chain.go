package reader

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	reader_internal "github.com/smartcontractkit/ccipocr3/internal/reader"
)

type HomeChain = reader_internal.HomeChain

type ChainConfig = reader_internal.ChainConfig

type ChainConfigInfo = reader_internal.ChainConfigInfo

func NewHomeChainReader(
	homeChainReader types.ContractReader,
	lggr logger.Logger,
	pollingInterval time.Duration,
) HomeChain {
	return reader_internal.NewHomeChainConfigPoller(homeChainReader, lggr, pollingInterval)
}
