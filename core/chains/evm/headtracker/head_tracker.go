package headtracker

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type headTracker = headtracker.HeadTracker[*evmtypes.Head, ethereum.Subscription, *big.Int, common.Hash]

var _ commontypes.HeadTracker[*evmtypes.Head, common.Hash] = (*headTracker)(nil)

func NewHeadTracker(
	lggr logger.Logger,
	ethClient evmclient.Client,
	config Config,
	htConfig HeadTrackerConfig,
	headBroadcaster httypes.HeadBroadcaster,
	headSaver httypes.HeadSaver,
	mailMon *utils.MailboxMonitor,
) httypes.HeadTracker {
	return headtracker.NewHeadTracker[*evmtypes.Head, ethereum.Subscription, *big.Int, common.Hash](
		lggr,
		ethClient,
		config,
		htConfig,
		headBroadcaster,
		headSaver,
		mailMon,
		func() *evmtypes.Head { return nil },
	)
}

var NullTracker httypes.HeadTracker = &nullTracker{}

type nullTracker struct{}

func (*nullTracker) Start(context.Context) error    { return nil }
func (*nullTracker) Close() error                   { return nil }
func (*nullTracker) Ready() error                   { return nil }
func (*nullTracker) HealthReport() map[string]error { return map[string]error{} }
func (*nullTracker) Name() string                   { return "" }
func (*nullTracker) SetLogLevel(zapcore.Level)      {}
func (*nullTracker) Backfill(ctx context.Context, headWithChain *evmtypes.Head, depth uint) (err error) {
	return nil
}
func (*nullTracker) LatestChain() *evmtypes.Head { return nil }
