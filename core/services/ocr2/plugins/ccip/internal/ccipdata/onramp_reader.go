package ccipdata

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

type LeafHasherInterface[H hashutil.Hash] interface {
	HashLeaf(log types.Log) (H, error)
}

const (
	COMMIT_CCIP_SENDS = "Commit ccip sends"
	CONFIG_CHANGED    = "Dynamic config changed"
)

//go:generate mockery --quiet --name OnRampReader --filename onramp_reader_mock.go --case=underscore
type OnRampReader interface {
	cciptypes.OnRampReader
}
