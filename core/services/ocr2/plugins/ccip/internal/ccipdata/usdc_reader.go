package ccipdata

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

var (
	// shortLivedInMemLogsCacheExpiration is used for the short-lived in meme logs cache.
	// Value should usually be set to just a few seconds, a larger duration will not increase performance and might
	// cause performance issues on re-orged logs.
	shortLivedInMemLogsCacheExpiration = 20 * time.Second
)

const (
	MESSAGE_SENT_FILTER_NAME = "USDC message sent"
)

var _ USDCReader = &USDCReaderImpl{}

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

	// shortLivedInMemLogs is a short-lived cache (items expire every few seconds)
	// used to prevent frequent log fetching from the log poller
	shortLivedInMemLogs *cache.Cache
}

func (u *USDCReaderImpl) Close() error {
	// FIXME Dim pgOpts removed from LogPoller
	return u.lp.UnregisterFilter(context.Background(), u.filter.Name)
}

func (u *USDCReaderImpl) RegisterFilters() error {
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
	var lpLogs []logpoller.Log

	// fetch all the usdc logs for the provided tx hash
	k := fmt.Sprintf("usdc-%s", txHash) // custom prefix to avoid key collision if someone re-uses the cache
	if rawLogs, foundInMem := u.shortLivedInMemLogs.Get(k); foundInMem {
		inMemLogs, ok := rawLogs.([]logpoller.Log)
		if !ok {
			return nil, errors.Errorf("unexpected in-mem logs type %T", rawLogs)
		}
		u.lggr.Debugw("found logs in memory", "k", k, "len", len(inMemLogs))
		lpLogs = inMemLogs
	}

	if len(lpLogs) == 0 {
		u.lggr.Debugw("fetching logs from lp", "k", k)
		logs, err := u.lp.IndexedLogsByTxHash(
			ctx,
			u.usdcMessageSent,
			u.transmitterAddress,
			common.HexToHash(txHash),
		)
		if err != nil {
			return nil, err
		}
		lpLogs = logs
		u.shortLivedInMemLogs.Set(k, logs, cache.DefaultExpiration)
		u.lggr.Debugw("fetched logs from lp", "logs", len(lpLogs))
	}

	// collect the logs with log index less than the provided log index
	allUsdcTokensData := make([][]byte, 0)
	for _, current := range lpLogs {
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

func NewUSDCReader(lggr logger.Logger, jobID string, transmitter common.Address, lp logpoller.LogPoller, registerFilters bool) (*USDCReaderImpl, error) {
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
		transmitterAddress:  transmitter,
		shortLivedInMemLogs: cache.New(shortLivedInMemLogsCacheExpiration, 2*shortLivedInMemLogsCacheExpiration),
	}

	if registerFilters {
		if err := r.RegisterFilters(); err != nil {
			return nil, fmt.Errorf("register filters: %w", err)
		}
	}
	return r, nil
}

func CloseUSDCReader(lggr logger.Logger, jobID string, transmitter common.Address, lp logpoller.LogPoller) error {
	r, err := NewUSDCReader(lggr, jobID, transmitter, lp, false)
	if err != nil {
		return err
	}
	return r.Close()
}
