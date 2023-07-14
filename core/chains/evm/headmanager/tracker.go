package headmanager

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/common/headmanager"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	hmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headmanager/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type tracker = headmanager.Tracker[*evmtypes.Head, ethereum.Subscription, *big.Int, common.Hash]

var _ commontypes.Tracker[*evmtypes.Head, common.Hash] = (*tracker)(nil)

func NewTracker(
	lggr logger.Logger,
	ethClient evmclient.Client,
	config Config,
	htConfig HeadTrackerConfig,
	headBroadcaster hmtypes.Broadcaster,
	headSaver hmtypes.Saver,
	mailMon *utils.MailboxMonitor,
) *tracker {
	return headmanager.NewTracker[*evmtypes.Head, ethereum.Subscription, *big.Int, common.Hash](
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

var NullTracker hmtypes.Tracker = &nullTracker{}

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
