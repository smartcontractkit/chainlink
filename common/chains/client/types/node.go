package types

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type NodeClientAPI[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
] interface {
	Close()
	ChainID(context.Context) (CHAIN_ID, error)
	Dial(callerCtx context.Context) error
	DialHTTP() error
	DisconnectAll()
	Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (SUB, error)
	ClientVersion(context.Context) (string, error)
}

type NodeConfig interface {
	PollFailureThreshold() uint32
	PollInterval() time.Duration
	SelectionMode() string
	SyncThreshold() uint32
}

type SendOnlyClientAPI[
	CHAIN_ID types.ID,
] interface {
	Close()
	ChainID(context.Context) (CHAIN_ID, error)
	DialHTTP() error
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
