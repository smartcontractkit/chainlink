package headtracker

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

func NewHeadTracker(
	lggr logger.Logger,
	ethClient httypes.Client,
	config commontypes.Config,
	htConfig commontypes.HeadTrackerConfig,
	headBroadcaster httypes.HeadBroadcaster,
	headSaver httypes.HeadSaver,
	mailMon *mailbox.Monitor,
) httypes.HeadTracker {
	return headtracker.NewHeadTracker[*evmtypes.Head, ethereum.Subscription](
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
func (*nullTracker) Backfill(ctx context.Context, headWithChain *evmtypes.Head) (err error) {
	return nil
}
func (*nullTracker) LatestChain() *evmtypes.Head { return nil }
func (*nullTracker) LatestAndFinalizedBlock(ctx context.Context) (latest, finalized *evmtypes.Head, err error) {
	return nil, nil, nil
}
