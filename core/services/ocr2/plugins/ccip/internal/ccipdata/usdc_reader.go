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

var _ USDCReader = &USDCReaderImpl{}

//go:generate mockery --quiet --name USDCReader --filename usdc_reader_mock.go --case=underscore
type USDCReader interface {
	// GetUSDCMessagePriorToLogIndexInTx returns the specified USDC message data.
	// e.g. if msg contains 3 tokens: [usdc1, wETH, usdc2] ignoring non-usdc tokens
	// if usdcTokenIndexOffset is 0 we select usdc2
	// if usdcTokenIndexOffset is 1 we select usdc1
	// The message logs are found using the provided transaction hash.
	GetUSDCMessagePriorToLogIndexInTx(ctx context.Context, logIndex int64, usdcTokenIndexOffset int, txHash string) ([]byte, error)
}

type USDCReaderImpl struct {
	usdcMessageSent    common.Hash
	lp                 logpoller.LogPoller
	filter             logpoller.Filter
	lggr               logger.Logger
	transmitterAddress common.Address
}

func (u *USDCReaderImpl) Close(qopts ...pg.QOpt) error {
	// FIXME Dim pgOpts removed from LogPoller
	return u.lp.UnregisterFilter(context.Background(), u.filter.Name)
}

func (u *USDCReaderImpl) RegisterFilters(qopts ...pg.QOpt) error {
	// FIXME Dim pgOpts removed from LogPoller
	return u.lp.RegisterFilter(context.Background(), u.filter)
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

func (u *USDCReaderImpl) GetUSDCMessagePriorToLogIndexInTx(ctx context.Context, logIndex int64, usdcTokenIndexOffset int, txHash string) ([]byte, error) {
	// fetch all the usdc logs for the provided tx hash
	logs, err := u.lp.IndexedLogsByTxHash(
		ctx,
		u.usdcMessageSent,
		u.transmitterAddress,
		common.HexToHash(txHash),
	)
	if err != nil {
		return nil, err
	}

	// collect the logs with log index less than the provided log index
	allUsdcTokensData := make([][]byte, 0)
	for _, current := range logs {
		if current.LogIndex < logIndex {
			u.lggr.Infow("Found USDC message", "logIndex", current.LogIndex, "txHash", current.TxHash.Hex(), "data", hexutil.Encode(current.Data))
			allUsdcTokensData = append(allUsdcTokensData, current.Data)
		}
	}

	usdcTokenIndex := (len(allUsdcTokensData) - 1) - usdcTokenIndexOffset

	if usdcTokenIndex < 0 || usdcTokenIndex >= len(allUsdcTokensData) {
		u.lggr.Errorw("usdc message not found",
			"logIndex", logIndex,
			"allUsdcTokenData", len(allUsdcTokensData),
			"txHash", txHash,
			"usdcTokenIndex", usdcTokenIndex,
		)
		return nil, errors.Errorf("usdc token index %d is not valid", usdcTokenIndex)
	}
	return parseUSDCMessageSent(allUsdcTokensData[usdcTokenIndex])
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
			Retention: CommitExecLogsRetention,
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
