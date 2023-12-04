package ccipdata

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

type OnRampDynamicConfig struct {
	Router                            common.Address
	MaxNumberOfTokensPerMsg           uint16
	DestGasOverhead                   uint32
	DestGasPerPayloadByte             uint16
	DestDataAvailabilityOverheadGas   uint32
	DestGasPerDataAvailabilityByte    uint16
	DestDataAvailabilityMultiplierBps uint16
	PriceRegistry                     common.Address
	MaxDataBytes                      uint32
	MaxPerMsgGasLimit                 uint32
}

//go:generate mockery --quiet --name OnRampReader --filename onramp_reader_mock.go --case=underscore
type OnRampReader interface {
	Closer
	// GetSendRequestsBetweenSeqNums returns all the finalized message send requests in the provided sequence numbers range (inclusive).
	GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, finalized bool) ([]Event[internal.EVM2EVMMessage], error)
	// Get router configured in the onRamp
	RouterAddress() (common.Address, error)
	Address() (common.Address, error)
	GetDynamicConfig() (OnRampDynamicConfig, error)
	RegisterFilters(qopts ...pg.QOpt) error
}

// NewOnRampReader determines the appropriate version of the onramp and returns a reader for it
func NewOnRampReader(lggr logger.Logger, sourceSelector, destSelector uint64, onRampAddress common.Address, sourceLP logpoller.LogPoller, source client.Client) (OnRampReader, error) {
	contractType, version, err := ccipconfig.TypeAndVersion(onRampAddress, source)
	if err != nil {
		return nil, errors.Errorf("expected '%v' got '%v' (%v)", ccipconfig.EVM2EVMOnRamp, contractType, err)
	}
	switch version.String() {
	case V1_0_0:
		return NewOnRampV1_0_0(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source)
	case V1_1_0:
		return NewOnRampV1_1_0(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source)
	case V1_2_0:
		return NewOnRampV1_2_0(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source)
	default:
		return nil, errors.Errorf("got unexpected version %v", version.String())
	}
}

func logsConfirmations(finalized bool) logpoller.Confirmations {
	if finalized {
		return logpoller.Finalized
	}
	return logpoller.Confirmations(0)
}
