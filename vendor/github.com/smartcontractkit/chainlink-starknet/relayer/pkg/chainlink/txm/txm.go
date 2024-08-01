package txm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	starknetaccount "github.com/NethermindEth/starknet.go/account"
	starknetrpc "github.com/NethermindEth/starknet.go/rpc"
	starknetutils "github.com/NethermindEth/starknet.go/utils"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
)

const (
	MaxQueueLen = 1000
)

type TxManager interface {
	Enqueue(accountAddress *felt.Felt, publicKey *felt.Felt, txFn starknetrpc.FunctionCall) error
	InflightCount() (int, int)
}

type Tx struct {
	publicKey      *felt.Felt
	accountAddress *felt.Felt
	call           starknetrpc.FunctionCall
}

type StarkTXM interface {
	services.Service
	TxManager
}

type starktxm struct {
	starter utils.StartStopOnce
	lggr    logger.Logger
	done    sync.WaitGroup
	stop    chan struct{}
	queue   chan Tx
	ks      KeystoreAdapter
	cfg     Config

	client       *utils.LazyLoad[*starknet.Client]
	feederClient *utils.LazyLoad[*starknet.FeederClient]
	accountStore *AccountStore
}

func New(lggr logger.Logger, keystore loop.Keystore, cfg Config, getClient func() (*starknet.Client, error),
	getFeederClient func() (*starknet.FeederClient, error)) (StarkTXM, error) {
	txm := &starktxm{
		lggr:         logger.Named(lggr, "StarknetTxm"),
		queue:        make(chan Tx, MaxQueueLen),
		stop:         make(chan struct{}),
		client:       utils.NewLazyLoad(getClient),
		feederClient: utils.NewLazyLoad(getFeederClient),
		ks:           NewKeystoreAdapter(keystore),
		cfg:          cfg,
		accountStore: NewAccountStore(),
	}

	return txm, nil
}

func (txm *starktxm) Name() string {
	return txm.lggr.Name()
}

func (txm *starktxm) Start(ctx context.Context) error {
	return txm.starter.StartOnce("starktxm", func() error {
		txm.done.Add(2) // waitgroup: broadcast loop and confirm loop
		go txm.broadcastLoop()
		go txm.confirmLoop()

		return nil
	})
}

func (txm *starktxm) broadcastLoop() {
	defer txm.done.Done()

	ctx, cancel := utils.ContextFromChan(txm.stop)
	defer cancel()

	txm.lggr.Debugw("broadcastLoop: started")
	for {
		select {
		case <-txm.stop:
			txm.lggr.Debugw("broadcastLoop: stopped")
			return
		case tx := <-txm.queue:
			if _, err := txm.client.Get(); err != nil {
				txm.lggr.Errorw("failed to fetch client: skipping processing tx", "error", err)
				continue
			}

			// broadcast tx serially - wait until accepted by mempool before processing next
			hash, err := txm.broadcast(ctx, tx.publicKey, tx.accountAddress, tx.call)
			if err != nil {
				txm.lggr.Errorw("transaction failed to broadcast", "error", err, "tx", tx.call)
			} else {
				txm.lggr.Infow("transaction broadcast", "txhash", hash)
			}
		}
	}
}

const FeeMargin uint32 = 115
const RPCNonceErrMsg = "Invalid transaction nonce"

func (txm *starktxm) estimateFriFee(ctx context.Context, client *starknet.Client, accountAddress *felt.Felt, tx starknetrpc.InvokeTxnV3) (*starknetrpc.FeeEstimate, *felt.Felt, error) {
	// skip prevalidation, which is known to overestimate amount of gas needed and error with L1GasBoundsExceedsBalance
	simFlags := []starknetrpc.SimulationFlag{starknetrpc.SKIP_VALIDATE}

	var largestEstimateNonce *felt.Felt

	for i := 1; i <= 5; i++ {
		txm.lggr.Infow("attempt to estimate fee", "attempt", i)

		estimateNonce, err := client.AccountNonce(ctx, accountAddress)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to check account nonce: %+w", err)
		}
		tx.Nonce = estimateNonce

		if largestEstimateNonce == nil || estimateNonce.Cmp(largestEstimateNonce) > 0 {
			largestEstimateNonce = estimateNonce
		}

		feeEstimate, err := client.Provider.EstimateFee(ctx, []starknetrpc.BroadcastTxn{tx}, simFlags, starknetrpc.BlockID{Tag: "pending"})
		if err != nil {
			var dataErr *starknetrpc.RPCError
			if !errors.As(err, &dataErr) {
				return nil, nil, fmt.Errorf("failed to read EstimateFee error: %T %+v", err, err)
			}
			data := dataErr.Data
			dataStr := fmt.Sprintf("%+v", data)

			txm.lggr.Errorw("failed to estimate fee", "attempt", i, "error", err, "data", dataStr)

			if strings.Contains(dataStr, RPCNonceErrMsg) {
				continue
			}

			return nil, nil, fmt.Errorf("failed to estimate fee: %T %+v", err, err)
		}

		// track the FRI estimate, but keep looping so we print out all estimates
		var friEstimate *starknetrpc.FeeEstimate
		for j, f := range feeEstimate {
			txm.lggr.Infow("Estimated fee", "attempt", i, "index", j, "EstimateNonce", estimateNonce, "GasConsumed", f.GasConsumed, "GasPrice", f.GasPrice, "DataGasConsumed", f.DataGasConsumed, "DataGasPrice", f.DataGasPrice, "OverallFee", f.OverallFee, "FeeUnit", string(f.FeeUnit))
			if f.FeeUnit == "FRI" {
				friEstimate = &feeEstimate[j]
			}
		}
		if friEstimate != nil {
			return friEstimate, largestEstimateNonce, nil
		}

		txm.lggr.Errorw("No FRI estimate was returned", "attempt", i)
	}

	txm.lggr.Errorw("all attempts to estimate fee failed")
	return nil, nil, fmt.Errorf("all attempts to estimate fee failed")
}

func (txm *starktxm) broadcast(ctx context.Context, publicKey *felt.Felt, accountAddress *felt.Felt, call starknetrpc.FunctionCall) (txhash string, err error) {
	client, err := txm.client.Get()
	if err != nil {
		txm.client.Reset()
		return txhash, fmt.Errorf("broadcast: failed to fetch client: %+w", err)
	}

	txStore := txm.accountStore.GetTxStore(accountAddress)
	if txStore == nil {
		initialNonce, accountNonceErr := client.AccountNonce(ctx, accountAddress)
		if accountNonceErr != nil {
			return txhash, fmt.Errorf("failed to check account nonce during TxStore creation: %+w", accountNonceErr)
		}
		newTxStore, createErr := txm.accountStore.CreateTxStore(accountAddress, initialNonce)
		if createErr != nil {
			return txhash, fmt.Errorf("failed to create TxStore: %+w", createErr)
		}
		txStore = newTxStore
	}

	// create new account
	cairoVersion := 2
	account, err := starknetaccount.NewAccount(client.Provider, accountAddress, publicKey.String(), txm.ks, cairoVersion)
	if err != nil {
		return txhash, fmt.Errorf("failed to create new account: %+w", err)
	}

	tx := starknetrpc.InvokeTxnV3{
		Type:          starknetrpc.TransactionType_Invoke,
		SenderAddress: account.AccountAddress,
		Version:       starknetrpc.TransactionV3,
		Signature:     []*felt.Felt{},
		Nonce:         &felt.Zero, // filled in below
		ResourceBounds: starknetrpc.ResourceBoundsMapping{
			L1Gas: starknetrpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
			L2Gas: starknetrpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
		},
		Tip:                   "0x0",
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         starknetrpc.DAModeL1,
		FeeMode:               starknetrpc.DAModeL1,
	}

	// Building the Calldata with the help of FmtCalldata where we pass in the FnCall struct along with the Cairo version
	tx.Calldata, err = account.FmtCalldata([]starknetrpc.FunctionCall{call})
	if err != nil {
		return txhash, err
	}

	friEstimate, largestEstimateNonce, err := txm.estimateFriFee(ctx, client, accountAddress, tx)
	if err != nil {
		return txhash, fmt.Errorf("failed to get FRI estimate: %+w", err)
	}

	nonce := txStore.GetNextNonce()
	if largestEstimateNonce.Cmp(nonce) > 0 {
		// The nonce value returned from the node during estimation is greater than our expected next nonce
		// - which means that we are behind, due to a resync. Fast forward our locally tracked nonce value.
		// See resyncNonce for a more detailed explanation.
		staleTxs := txStore.SetNextNonce(largestEstimateNonce)
		txm.lggr.Infow("fast-forwarding nonce after resync", "previousNonce", nonce, "updatedNonce", largestEstimateNonce, "staleTxs", len(staleTxs))
		if len(staleTxs) > 0 {
			txm.lggr.Errorw("unexpected stale transactions after nonce fast-forward", "accountAddress", accountAddress)
		}
		nonce = largestEstimateNonce
	}

	// TODO: consider making this configurable
	// pad estimate to 250% (add extra because estimate did not include validation)
	gasConsumed := friEstimate.GasConsumed.BigInt(new(big.Int))
	expandedGas := new(big.Int).Mul(gasConsumed, big.NewInt(250))
	maxGas := new(big.Int).Div(expandedGas, big.NewInt(100))
	tx.ResourceBounds.L1Gas.MaxAmount = starknetrpc.U64(starknetutils.BigIntToFelt(maxGas).String())

	// pad by 150%
	gasPrice := friEstimate.GasPrice.BigInt(new(big.Int))
	overallFee := friEstimate.OverallFee.BigInt(new(big.Int)) // overallFee = gas_used*gas_price + data_gas_used*data_gas_price

	// TODO: consider making this configurable
	// pad estimate to 150% (add extra because estimate did not include validation)
	gasUnits := new(big.Int).Div(overallFee, gasPrice)
	expandedGasUnits := new(big.Int).Mul(gasUnits, big.NewInt(150))
	maxGasUnits := new(big.Int).Div(expandedGasUnits, big.NewInt(100))
	tx.ResourceBounds.L1Gas.MaxAmount = starknetrpc.U64(starknetutils.BigIntToFelt(maxGasUnits).String())

	// pad by 150%
	expandedGasPrice := new(big.Int).Mul(gasPrice, big.NewInt(150))
	maxGasPrice := new(big.Int).Div(expandedGasPrice, big.NewInt(100))
	tx.ResourceBounds.L1Gas.MaxPricePerUnit = starknetrpc.U128(starknetutils.BigIntToFelt(maxGasPrice).String())

	txm.lggr.Infow("Set resource bounds", "L1MaxAmount", tx.ResourceBounds.L1Gas.MaxAmount, "L1MaxPricePerUnit", tx.ResourceBounds.L1Gas.MaxPricePerUnit)

	tx.Nonce = nonce
	// Re-sign transaction now that we've determined MaxFee
	// TODO: SignInvokeTransaction for V3 is missing so we do it by hand
	hash, err := account.TransactionHashInvoke(tx)
	if err != nil {
		return txhash, err
	}
	signature, err := account.Sign(ctx, hash)
	if err != nil {
		return txhash, err
	}
	tx.Signature = signature

	execCtx, execCancel := context.WithTimeout(ctx, txm.cfg.TxTimeout())
	defer execCancel()

	// finally, transmit the invoke
	res, err := account.AddInvokeTransaction(execCtx, tx)
	if err != nil {
		// TODO: handle initial broadcast errors - what kind of errors occur?
		var dataErr *starknetrpc.RPCError
		var dataStr string
		if !errors.As(err, &dataErr) {
			return txhash, fmt.Errorf("failed to read EstimateFee error: %T %+v", err, err)
		}
		data := dataErr.Data
		dataStr = fmt.Sprintf("%+v", data)
		txm.lggr.Errorw("failed to invoke tx", "accountAddress", accountAddress, "error", err, "data", dataStr)

		if strings.Contains(dataStr, RPCNonceErrMsg) {
			// if we see an invalid nonce error at the broadcast stage, that means that we are out of sync.
			// see the comment at resyncNonce for more details.
			if resyncErr := txm.resyncNonce(ctx, client, accountAddress); resyncErr != nil {
				txm.lggr.Errorw("failed to resync nonce after unsuccessful invoke", "error", err, "resyncError", resyncErr)
				return txhash, fmt.Errorf("failed to resync after bad invoke: %+w", err)
			}
		}
		return txhash, fmt.Errorf("failed to invoke tx: %+w", err)
	}
	// handle nil pointer
	if res == nil {
		return txhash, errors.New("execute response and error are nil")
	}

	// update nonce if transaction is successful
	txhash = res.TransactionHash.String()
	err = txStore.AddUnconfirmed(nonce, txhash, call, publicKey)
	if err != nil {
		return txhash, fmt.Errorf("failed to add unconfirmed tx: %+w", err)
	}
	return txhash, nil
}

func (txm *starktxm) confirmLoop() {
	defer txm.done.Done()

	ctx, cancel := utils.ContextFromChan(txm.stop)
	defer cancel()

	tick := time.After(txm.cfg.ConfirmationPoll())

	txm.lggr.Debugw("confirmLoop: started")

	for {
		var start time.Time
		select {
		case <-tick:
			start = time.Now()
			client, err := txm.client.Get()
			if err != nil {
				txm.lggr.Errorw("failed to load client", "error", err)
				break
			}

			allUnconfirmedTxs := txm.accountStore.GetAllUnconfirmed()
			for accountAddressStr, unconfirmedTxs := range allUnconfirmedTxs {
				accountAddress, err := new(felt.Felt).SetString(accountAddressStr)
				// this should never occur because the acccount address string key was created from the account address felt.
				if err != nil {
					txm.lggr.Errorw("could not recreate account address felt", "accountAddress", accountAddressStr)
					continue
				}
				for _, unconfirmedTx := range unconfirmedTxs {
					hash := unconfirmedTx.Hash
					f, err := starknetutils.HexToFelt(hash)
					if err != nil {
						txm.lggr.Errorw("invalid felt value", "hash", hash)
						continue
					}
					response, err := client.Provider.GetTransactionStatus(ctx, f)

					// tx can be rejected due to a nonce error. but we cannot know from the Starknet RPC directly  so we have to wait for
					// a broadcasted tx to fail in order to fix the nonce errors

					if err != nil {
						txm.lggr.Errorw("failed to fetch transaction status", "hash", hash, "nonce", unconfirmedTx.Nonce, "error", err)
						continue
					}

					finalityStatus := response.FinalityStatus
					executionStatus := response.ExecutionStatus

					// any finalityStatus other than received
					if finalityStatus == starknetrpc.TxnStatus_Accepted_On_L1 || finalityStatus == starknetrpc.TxnStatus_Accepted_On_L2 || finalityStatus == starknetrpc.TxnStatus_Rejected {
						txm.lggr.Debugw(fmt.Sprintf("tx confirmed: %s", finalityStatus), "hash", hash, "nonce", unconfirmedTx.Nonce, "finalityStatus", finalityStatus)
						if err := txm.accountStore.GetTxStore(accountAddress).Confirm(unconfirmedTx.Nonce, hash); err != nil {
							txm.lggr.Errorw("failed to confirm tx in TxStore", "hash", hash, "accountAddress", accountAddress, "error", err)
						}
					}

					// currently, feeder client is only way to get rejected reason
					if finalityStatus == starknetrpc.TxnStatus_Rejected {
						// we assume that all rejected transactions results in a unused rejected nonce, so
						// resync. see the comment at resyncNonce for more details.
						if resyncErr := txm.resyncNonce(ctx, client, accountAddress); resyncErr != nil {
							txm.lggr.Errorw("resync failed for rejected tx", "error", resyncErr)
						}

						go txm.logFeederError(ctx, hash, f)
					}

					if executionStatus == starknetrpc.TxnExecutionStatusREVERTED {
						// TODO: get revert reason?
						txm.lggr.Errorw("transaction reverted", "hash", hash)
					}
				}
			}
		case <-txm.stop:
			txm.lggr.Debugw("confirmLoop: stopped")
			return
		}
		t := txm.cfg.ConfirmationPoll() - time.Since(start)
		tick = time.After(utils.WithJitter(t.Abs()))
	}
}

func (txm *starktxm) logFeederError(ctx context.Context, hash string, f *felt.Felt) {
	feederClient, err := txm.feederClient.Get()
	if err != nil {
		txm.lggr.Errorw("failed to load feeder client", "error", err)
		return
	}

	rejectedTx, err := feederClient.TransactionFailure(ctx, f)
	if err != nil {
		txm.lggr.Errorw("failed to fetch reason for transaction failure", "hash", hash, "error", err)
		return
	}

	txm.lggr.Errorw("feeder rejected reason", "hash", hash, "errorMessage", rejectedTx.ErrorMessage)
}

func (txm *starktxm) resyncNonce(ctx context.Context, client *starknet.Client, accountAddress *felt.Felt) error {
	/*
	   the follow errors indicate that there could be a problem with our locally tracked nonce value:
	       1. a EstimateFee was successful, but broadcasting using the locally tracked nonce results in a nonce error,
	       2. a transaction was rejected after a successful broadcast.

	   for these cases, we call starknet_getNonce from the RPC node and resync the locally tracked next nonce
	   with the RPC node's value.

	   however, while the value returned by starknet_getNonce is eventually consistent, it can be lower than the actual
	   next nonce value when pending transactions haven't yet been processed - resulting in more category 1
	   invalid nonce broadcast errors.

	   in order to recover from these cases, each time we do starknet_getNonce during estimation (see estimateFriFee),
	   we compare it with our locally tracked nonce - if it is greater, than that means our locally tracked value is
	   behind, and we fast forward. this ensures our locally tracked value will also eventually be correct.
	*/

	rpcNonce, err := client.AccountNonce(ctx, accountAddress)
	if err != nil {
		return fmt.Errorf("failed to check nonce during resync: %+w", err)
	}

	txStore := txm.accountStore.GetTxStore(accountAddress)
	currentNonce := txStore.GetNextNonce()

	if rpcNonce.Cmp(currentNonce) == 0 {
		txm.lggr.Infow("resync nonce skipped, nonce value is the same", "accountAddress", accountAddress, "nonce", currentNonce)
		return nil
	}

	staleTxs := txStore.SetNextNonce(rpcNonce)

	txm.lggr.Infow("resynced nonce", "accountAddress", "accountAddress", "previousNonce", currentNonce, "updatedNonce", rpcNonce, "staleTxCount", len(staleTxs))

	return nil
}

func (txm *starktxm) Close() error {
	return txm.starter.StopOnce("starktxm", func() error {
		close(txm.stop)
		txm.done.Wait()
		return nil
	})
}

func (txm *starktxm) Healthy() error {
	return txm.starter.Healthy()
}

func (txm *starktxm) Ready() error {
	return txm.starter.Ready()
}

func (txm *starktxm) HealthReport() map[string]error {
	return map[string]error{txm.Name(): txm.Healthy()}
}

func (txm *starktxm) Enqueue(accountAddress, publicKey *felt.Felt, tx starknetrpc.FunctionCall) error {
	// validate key exists for sender
	// use the embedded Loopp Keystore to do this; the spec and design
	// encourage passing nil data to the loop.Keystore.Sign as way to test
	// existence of a key
	if _, err := txm.ks.Loopp().Sign(context.Background(), publicKey.String(), nil); err != nil {
		return fmt.Errorf("enqueue: failed to sign: %+w", err)
	}

	select {
	case txm.queue <- Tx{publicKey: publicKey, accountAddress: accountAddress, call: tx}: // TODO fix naming here
	default:
		return fmt.Errorf("failed to enqueue transaction: %+v", tx)
	}

	return nil
}

func (txm *starktxm) InflightCount() (queue int, unconfirmed int) {
	return len(txm.queue), txm.accountStore.GetTotalInflightCount()
}
