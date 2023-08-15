package types

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type FeedIDPrefix uint16

const (
	_         FeedIDPrefix = 0 // reserved to prevent errors where a zero-default creeps through somewhere
	REPORT_V1 FeedIDPrefix = 1
	REPORT_V2 FeedIDPrefix = 2
	REPORT_V3 FeedIDPrefix = 3
	_         FeedIDPrefix = 0xFFFF // reserved for future use
)

type FeedID [32]byte

func (f FeedID) String() string {
	return fmt.Sprintf("%x", f[:])
}

func (f FeedID) Hex() string {
	return f.String()
}

// UnmarshalText parses a hash in hex syntax.
func (f *FeedID) UnmarshalText(input []byte) error {
	return (*common.Hash)(f).UnmarshalText(input)
}

// UnmarshalJSON parses a hash in hex syntax.
func (f *FeedID) UnmarshalJSON(input []byte) error {
	return (*common.Hash)(f).UnmarshalJSON(input)
}

func (f FeedID) Version() FeedIDPrefix {
	return FeedIDPrefix(binary.BigEndian.Uint16(f[:2]))
}

func (f FeedID) IsV1() bool { return f.Version() == REPORT_V1 }
func (f FeedID) IsV2() bool { return f.Version() == REPORT_V2 }
func (f FeedID) IsV3() bool { return f.Version() == REPORT_V3 }

//go:generate mockery --quiet --name ChainHeadTracker --output ../mocks/ --case=underscore
type ChainHeadTracker interface {
	Client() evmclient.Client
	HeadTracker() httypes.HeadTracker
}

type DataSourceORM interface {
	LatestReport(ctx context.Context, feedID [32]byte, qopts ...pg.QOpt) (report []byte, err error)
}
