package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// pollInterval is the maximum amount of time in addition to
// EthTxResendAfterThreshold that we will wait before resending an attempt
const DefaultResenderPollInterval = 5 * time.Second

// EthResender periodically picks up transactions that have been languishing
// unconfirmed for a configured amount of time without being sent, and sends
// their highest priced attempt again. This helps to defend against geth/parity
// silently dropping txes, or txes being ejected from the mempool.
//
// Previously we relied on the bumper to do this for us implicitly but there
// can occasionally be problems with this (e.g. abnormally long block times, or
// if gas bumping is disabled)
type EthResender[ADDR types.Hashable, TX_HASH types.Hashable, BLOCK_HASH types.Hashable] struct {
	txStorageService txmgrtypes.TxStorageService[ADDR, big.Int, TX_HASH, BLOCK_HASH, NewTx[ADDR], *evmtypes.Receipt, EthTx[ADDR, TX_HASH], EthTxAttempt[ADDR, TX_HASH], int64, int64]
	ethClient        evmclient.Client
	ks               txmgrtypes.KeyStore[ADDR, *big.Int, int64]
	chainID          big.Int
	interval         time.Duration
	config           Config
	logger           logger.Logger

	ctx    context.Context
	cancel context.CancelFunc
	chDone chan struct{}
}

// NewEthResender creates a new concrete EthResender
func NewEthResender(
	lggr logger.Logger,
	txStorageService EvmTxStorageService,
	ethClient evmclient.Client, ks EvmKeyStore,
	pollInterval time.Duration,
	config Config,
) *EvmEthResender {
	if config.EthTxResendAfterThreshold() == 0 {
		panic("EthResender requires a non-zero threshold")
	}
	// todo: add context to evmTxStorageService
	ctx, cancel := context.WithCancel(context.Background())
	return &EvmEthResender{
		txStorageService,
		ethClient,
		ks,
		*ethClient.ChainID(),
		pollInterval,
		config,
		lggr.Named("EthResender"),
		ctx,
		cancel,
		make(chan struct{}),
	}
}

// Start is a comment which satisfies the linter
func (er *EthResender[ADDR, TX_HASH, BLOCK_HASH]) Start() {
	er.logger.Debugf("Enabled with poll interval of %s and age threshold of %s", er.interval, er.config.EthTxResendAfterThreshold())
	go er.runLoop()
}

// Stop is a comment which satisfies the linter
func (er *EthResender[ADDR, TX_HASH, BLOCK_HASH]) Stop() {
	er.cancel()
	<-er.chDone
}

func (er *EthResender[ADDR, TX_HASH, BLOCK_HASH]) runLoop() {
	defer close(er.chDone)

	if err := er.resendUnconfirmed(); err != nil {
		er.logger.Warnw("Failed to resend unconfirmed transactions", "err", err)
	}

	ticker := time.NewTicker(utils.WithJitter(er.interval))
	defer ticker.Stop()
	for {
		select {
		case <-er.ctx.Done():
			return
		case <-ticker.C:
			if err := er.resendUnconfirmed(); err != nil {
				er.logger.Warnw("Failed to resend unconfirmed transactions", "err", err)
			}
		}
	}
}

func (er *EthResender[ADDR, TX_HASH, BLOCK_HASH]) resendUnconfirmed() error {
	keys, err := er.ks.EnabledAddressesForChain(&er.chainID)
	if err != nil {
		return errors.Wrapf(err, "EthResender failed getting enabled keys for chain %s", er.chainID.String())
	}
	ageThreshold := er.config.EthTxResendAfterThreshold()
	maxInFlightTransactions := er.config.EvmMaxInFlightTransactions()
	olderThan := time.Now().Add(-ageThreshold)
	var allAttempts []EthTxAttempt[ADDR, TX_HASH]
	for _, k := range keys {
		var attempts []EthTxAttempt[ADDR, TX_HASH]
		attempts, err = er.txStorageService.FindEthTxAttemptsRequiringResend(olderThan, maxInFlightTransactions, er.chainID, k)
		if err != nil {
			return errors.Wrap(err, "failed to FindEthTxAttemptsRequiringResend")
		}

		allAttempts = append(allAttempts, attempts...)
	}

	if len(allAttempts) == 0 {
		return nil
	}
	er.logger.Infow(fmt.Sprintf("Re-sending %d unconfirmed transactions that were last sent over %s ago. These transactions are taking longer than usual to be mined. %s", len(allAttempts), ageThreshold, label.NodeConnectivityProblemWarning), "n", len(allAttempts))

	batchSize := int(er.config.EvmRPCDefaultBatchSize())
	ctx, cancel := context.WithTimeout(er.ctx, batchSendTransactionTimeout)
	defer cancel()
	reqs, err := batchSendTransactions(ctx, er.txStorageService, allAttempts, batchSize, er.logger, er.ethClient)
	if err != nil {
		return errors.Wrap(err, "failed to re-send transactions")
	}
	logResendResult(er.logger, reqs)

	return nil
}

func logResendResult(lggr logger.Logger, reqs []rpc.BatchElem) {
	var nNew int
	var nFatal int
	for _, req := range reqs {
		serr := evmclient.NewSendError(req.Error)
		if serr == nil {
			nNew++
		} else if serr.Fatal() {
			nFatal++
		}
	}
	lggr.Debugw("Completed", "n", len(reqs), "nNew", nNew, "nFatal", nFatal)
}
