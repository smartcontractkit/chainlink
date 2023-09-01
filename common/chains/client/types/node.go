package types

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Head interface {
	BlockNumber() int64
	BlockDifficulty() *utils.Big
}

type NodeClient[
	CHAIN_ID types.ID,
	HEAD Head,
] interface {
	Close()
	ChainID(context.Context) (CHAIN_ID, error)
	Dial(callerCtx context.Context) error
	DialHTTP() error
	DisconnectAll()
	Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (types.Subscription, error)
	ClientVersion(context.Context) (string, error)
}

type NodeConfig interface {
	PollFailureThreshold() uint32
	PollInterval() time.Duration
	SelectionMode() string
	SyncThreshold() uint32
}

type NodeTier int

const (
	Primary = NodeTier(iota)
	Secondary
)

func (n NodeTier) String() string {
	switch n {
	case Primary:
		return "primary"
	case Secondary:
		return "secondary"
	default:
		return fmt.Sprintf("NodeTier(%d)", n)
	}
}
