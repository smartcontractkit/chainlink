package ccipdata

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
)

type LeafHasherInterface[H hashlib.Hash] interface {
	HashLeaf(log types.Log) (H, error)
}

const (
	COMMIT_CCIP_SENDS = "Commit ccip sends"
)

//go:generate mockery --quiet --name OnRampReader --filename onramp_reader_mock.go --case=underscore
type OnRampReader interface {
	cciptypes.OnRampReader
}
