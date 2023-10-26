package ccipdata

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type LeafHasherInterface[H hashlib.Hash] interface {
	HashLeaf(log types.Log) (H, error)
}

const (
	COMMIT_CCIP_SENDS = "Commit ccip sends"
)

//go:generate mockery --quiet --name OnRampReader --output . --filename onramp_reader_mock.go --inpackage --case=underscore
type OnRampReader interface {
	Closer
	// GetSendRequestsGteSeqNum returns all the message send requests with sequence number greater than or equal to the provided.
	// If checkFinalityTags is set to true then confs param is ignored, the latest finalized block is used in the query.
	GetSendRequestsGteSeqNum(ctx context.Context, seqNum uint64, confs int) ([]Event[internal.EVM2EVMMessage], error)
	// GetSendRequestsBetweenSeqNums returns all the message send requests in the provided sequence numbers range (inclusive).
	GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, confs int) ([]Event[internal.EVM2EVMMessage], error)
	// Get router configured in the onRamp
	RouterAddress() (common.Address, error)
}

// NewOnRampReader determines the appropriate version of the onramp and returns a reader for it
func NewOnRampReader(lggr logger.Logger, sourceSelector, destSelector uint64, onRampAddress common.Address, sourceLP logpoller.LogPoller, source client.Client, finalityTags bool, qopts ...pg.QOpt) (OnRampReader, error) {
	contractType, version, err := ccipconfig.TypeAndVersion(onRampAddress, source)
	if err != nil {
		return nil, errors.Errorf("expected '%v' got '%v' (%v)", ccipconfig.EVM2EVMOnRamp, contractType, err)
	}
	switch version.String() {
	case V1_0_0:
		return NewOnRampV1_0_0(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source, finalityTags)
	case V1_1_0:
		return NewOnRampV1_1_0(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source, finalityTags)
	case V1_2_0:
		return NewOnRampV1_2_0(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source, finalityTags)
	default:
		return nil, errors.Errorf("got unexpected version %v", version.String())
	}
}

func latestFinalizedBlockHash(ctx context.Context, client client.Client) (common.Hash, error) {
	// If the chain is based on explicit finality we only examine logs less than or equal to the latest finalized block number.
	// NOTE: there appears to be a bug in ethclient whereby BlockByNumber fails with "unsupported txtype" when trying to parse the block
	// when querying L2s, headers however work.
	// TODO (CCIP-778): Migrate to core finalized tags, below doesn't work for some chains e.g. Celo.
	latestFinalizedHeader, err := client.HeaderByNumber(
		ctx,
		big.NewInt(rpc.FinalizedBlockNumber.Int64()),
	)
	if err != nil {
		return common.Hash{}, err
	}

	if latestFinalizedHeader == nil {
		return common.Hash{}, errors.New("latest finalized header is nil")
	}
	if latestFinalizedHeader.Number == nil {
		return common.Hash{}, errors.New("latest finalized number is nil")
	}
	return latestFinalizedHeader.Hash(), nil
}
