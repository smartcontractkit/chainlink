package ccipdata

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	MESSAGE_SENT_FILTER_NAME = "USDC message sent"
)

//go:generate mockery --quiet --name USDCReader --filename usdc_reader_mock.go --case=underscore
type USDCReader interface {
	// GetLastUSDCMessagePriorToLogIndexInTx returns the last USDC message that was sent
	// before the provided log index in the given transaction.
	GetLastUSDCMessagePriorToLogIndexInTx(ctx context.Context, logIndex int64, txHash string) ([]byte, error)
}

type USDCReaderImpl struct {
	usdcMessageSent    common.Hash
	lp                 logpoller.LogPoller
	filter             logpoller.Filter
	lggr               logger.Logger
	transmitterAddress common.Address
}

func (u *USDCReaderImpl) Close(qopts ...pg.QOpt) error {
	return u.lp.UnregisterFilter(u.filter.Name, qopts...)
}

func (u *USDCReaderImpl) RegisterFilters(qopts ...pg.QOpt) error {
	return u.lp.RegisterFilter(u.filter, qopts...)
}

// usdcPayload has to match the onchain event emitted by the USDC message transmitter
type usdcPayload []byte

func (d usdcPayload) AbiString() string {
	return `[{"type": "bytes"}]`
}

func (d usdcPayload) Validate() error {
	if len(d) == 0 {
		return errors.New("must be non-empty")
	}
	return nil
}

func parseUSDCMessageSent(logData []byte) ([]byte, error) {
	decodeAbiStruct, err := abihelpers.DecodeAbiStruct[usdcPayload](logData)
	if err != nil {
		return nil, err
	}
	return decodeAbiStruct, nil
}

func (u *USDCReaderImpl) GetLastUSDCMessagePriorToLogIndexInTx(ctx context.Context, logIndex int64, txHash string) ([]byte, error) {
	logs, err := u.lp.IndexedLogsByTxHash(
		u.usdcMessageSent,
		u.transmitterAddress,
		common.HexToHash(txHash),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}

	for i := range logs {
		current := logs[len(logs)-i-1]
		if current.LogIndex < logIndex {
			u.lggr.Infow("Found USDC message", "logIndex", current.LogIndex, "txHash", current.TxHash.Hex(), "data", hexutil.Encode(current.Data))
			return parseUSDCMessageSent(current.Data)
		}
	}
	return nil, errors.Errorf("no USDC message found prior to log index %d in tx %s", logIndex, txHash)
}

func NewUSDCReader(lggr logger.Logger, jobID string, transmitter common.Address, lp logpoller.LogPoller, registerFilters bool, qopts ...pg.QOpt) (*USDCReaderImpl, error) {
	eventSig := utils.Keccak256Fixed([]byte("MessageSent(bytes)"))

	r := &USDCReaderImpl{
		lggr:            lggr,
		lp:              lp,
		usdcMessageSent: eventSig,
		filter: logpoller.Filter{
			Name:      logpoller.FilterName(MESSAGE_SENT_FILTER_NAME, jobID, transmitter.Hex()),
			EventSigs: []common.Hash{eventSig},
			Addresses: []common.Address{transmitter},
		},
		transmitterAddress: transmitter,
	}

	if registerFilters {
		if err := r.RegisterFilters(qopts...); err != nil {
			return nil, fmt.Errorf("register filters: %w", err)
		}
	}
	return r, nil
}

func CloseUSDCReader(lggr logger.Logger, jobID string, transmitter common.Address, lp logpoller.LogPoller, qopts ...pg.QOpt) error {
	r, err := NewUSDCReader(lggr, jobID, transmitter, lp, false)
	if err != nil {
		return err
	}
	return r.Close(qopts...)
}
