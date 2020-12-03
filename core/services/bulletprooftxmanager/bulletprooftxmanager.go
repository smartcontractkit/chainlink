package bulletprooftxmanager

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"math/big"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethAccounts "github.com/ethereum/go-ethereum/accounts"
	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// For more information about the BulletproofTxManager architecture, see the design doc:
// https://www.notion.so/chainlink/BulletproofTxManager-Architecture-Overview-9dc62450cd7a443ba9e7dceffa1a8d6b

const (
	// maxEthNodeRequestTime is the worst case time we will wait for a response
	// from the eth node before we consider it to be an error
	maxEthNodeRequestTime = 15 * time.Second
)

var (
	promNumGasBumps = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_num_gas_bumps",
		Help: "Number of gas bumps",
	})

	promGasBumpExceedsLimit = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tx_manager_gas_bump_exceeds_limit",
		Help: "Number of times gas bumping failed from exceeding the configured limit. Any counts of this type indicate a serious problem.",
	})
)

// SendEther creates a transaction that transfers the given value of ether
func SendEther(s *strpkg.Store, from, to gethCommon.Address, value assets.Eth) (etx models.EthTx, err error) {
	if to == utils.ZeroAddress {
		return etx, errors.New("cannot send ether to zero address")
	}
	etx = models.EthTx{
		FromAddress:    from,
		ToAddress:      to,
		EncodedPayload: []byte{},
		Value:          value,
		GasLimit:       s.Config.EthGasLimitDefault(),
		State:          models.EthTxUnstarted,
	}
	err = s.DB.Create(&etx).Error
	return etx, err
}

func newAttempt(s *strpkg.Store, etx models.EthTx, gasPrice *big.Int) (models.EthTxAttempt, error) {
	attempt := models.EthTxAttempt{}
	account, err := s.KeyStore.GetAccountByAddress(etx.FromAddress)
	if err != nil {
		return attempt, errors.Wrapf(err, "error getting account %s for transaction %v", etx.FromAddress.String(), etx.ID)
	}

	transaction := gethTypes.NewTransaction(uint64(*etx.Nonce), etx.ToAddress, etx.Value.ToInt(), etx.GasLimit, gasPrice, etx.EncodedPayload)
	hash, signedTxBytes, err := signTx(s.KeyStore, account, transaction, s.Config.ChainID())
	if err != nil {
		return attempt, errors.Wrapf(err, "error using account %s to sign transaction %v", etx.FromAddress.String(), etx.ID)
	}

	attempt.State = models.EthTxAttemptInProgress
	attempt.SignedRawTx = signedTxBytes
	attempt.EthTxID = etx.ID
	attempt.GasPrice = *utils.NewBig(gasPrice)
	attempt.Hash = hash

	return attempt, nil
}

func signTx(keyStore strpkg.KeyStoreInterface, account gethAccounts.Account, tx *gethTypes.Transaction, chainID *big.Int) (gethCommon.Hash, []byte, error) {
	signedTx, err := keyStore.SignTx(account, tx, chainID)
	if err != nil {
		return gethCommon.Hash{}, nil, errors.Wrap(err, "signTx failed")
	}
	rlp := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(rlp); err != nil {
		return gethCommon.Hash{}, nil, errors.Wrap(err, "signTx failed")
	}
	return signedTx.Hash(), rlp.Bytes(), nil

}

// send broadcasts the transaction to the ethereum network, writes any relevant
// data onto the attempt and returns an error (or nil) depending on the status
func sendTransaction(ctx context.Context, ethClient eth.Client, a models.EthTxAttempt) *eth.SendError {
	signedTx, err := a.GetSignedTx()
	if err != nil {
		return eth.NewFatalSendError(err)
	}

	ctx, cancel := context.WithTimeout(ctx, maxEthNodeRequestTime)
	defer cancel()
	//m, err := signedTx.AsMessage(gethTypes.HomesteadSigner{})
	//if err != nil {
	//	return eth.NewFatalSendError(err)
	//}
	g, err := ethClient.EstimateGas(ctx, ethereum.CallMsg{
		From:     a.EthTx.FromAddress,
		To:       signedTx.To(),
		Gas:      signedTx.Gas(),
		GasPrice: signedTx.GasPrice(),
		Value:    signedTx.Value(),
		Data:     signedTx.Data(),
	})
	fmt.Println("gas", g, "err", err)
	if err != nil {
		return eth.NewFatalSendError(err)
	}

	err = ethClient.SendTransaction(ctx, signedTx)
	err = errors.WithStack(err)

	logger.Debugw("BulletproofTxManager: Broadcasting transaction", "ethTxAttemptID", a.ID, "txHash", signedTx.Hash(), "gasPriceWei", a.GasPrice.ToInt().Int64())
	sendErr := eth.NewSendError(err)
	if sendErr.IsTransactionAlreadyInMempool() {
		logger.Debugw("transaction already in mempool", "txHash", signedTx.Hash(), "nodeErr", sendErr.Error())
		fmt.Println("already in mem")
		return nil
	}
	return eth.NewSendError(err)
}

// sendEmptyTransaction sends a transaction with 0 Eth and an empty payload to the burn address
// May be useful for clearing stuck nonces
func sendEmptyTransaction(
	ethClient eth.Client,
	keyStore strpkg.KeyStoreInterface,
	nonce uint64,
	gasLimit uint64,
	gasPriceWei *big.Int,
	account gethAccounts.Account,
	chainID *big.Int,
) (_ *gethTypes.Transaction, err error) {
	defer utils.WrapIfError(&err, "sendEmptyTransaction failed")

	to := utils.ZeroAddress
	value := big.NewInt(0)
	payload := []byte{}
	tx := gethTypes.NewTransaction(nonce, to, value, gasLimit, gasPriceWei, payload)
	signedTx, err := keyStore.SignTx(account, tx, chainID)
	if err != nil {
		return signedTx, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), maxEthNodeRequestTime)
	defer cancel()
	err = ethClient.SendTransaction(ctx, signedTx)
	return signedTx, err
}

func saveReplacementInProgressAttempt(store *strpkg.Store, oldAttempt models.EthTxAttempt, replacementAttempt *models.EthTxAttempt) error {
	if oldAttempt.State != models.EthTxAttemptInProgress || replacementAttempt.State != models.EthTxAttemptInProgress {
		return errors.New("expected attempts to be in_progress")
	}
	if oldAttempt.ID == 0 {
		return errors.New("expected oldAttempt to have an ID")
	}
	return store.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE id = ? `, oldAttempt.ID).Error; err != nil {
			return errors.Wrap(err, "saveReplacementInProgressAttempt failed")
		}
		return errors.Wrap(tx.Create(replacementAttempt).Error, "saveReplacementInProgressAttempt failed")
	})
}

// BumpGas computes the next gas price to attempt as the largest of:
// - A configured percentage bump (ETH_GAS_BUMP_PERCENT) on top of the baseline price.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline price.
// The baseline price is the maximum of the previous gas price attempt and the node's current gas price.
func BumpGas(config orm.ConfigReader, originalGasPrice *big.Int) (*big.Int, error) {
	baselinePrice := max(originalGasPrice, config.EthGasPriceDefault())

	var priceByPercentage = new(big.Int)
	priceByPercentage.Mul(baselinePrice, big.NewInt(int64(100+config.EthGasBumpPercent())))
	priceByPercentage.Div(priceByPercentage, big.NewInt(100))

	var priceByIncrement = new(big.Int)
	priceByIncrement.Add(baselinePrice, config.EthGasBumpWei())

	bumpedGasPrice := max(priceByPercentage, priceByIncrement)
	if bumpedGasPrice.Cmp(config.EthMaxGasPriceWei()) > 0 {
		promGasBumpExceedsLimit.Inc()
		return config.EthMaxGasPriceWei(), errors.Errorf("bumped gas price of %s would exceed configured max gas price of %s (original price was %s)",
			bumpedGasPrice.String(), config.EthMaxGasPriceWei(), originalGasPrice.String())
	} else if bumpedGasPrice.Cmp(originalGasPrice) == 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// ETH_GAS_BUMP_PERCENT and ETH_GAS_BUMP_WEI in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedGasPrice, errors.Errorf("bumped gas price of %s is equal to original gas price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI", bumpedGasPrice.String(), originalGasPrice.String())
	}
	promNumGasBumps.Inc()
	return bumpedGasPrice, nil
}

func max(a, b *big.Int) *big.Int {
	if a.Cmp(b) >= 0 {
		return a
	}
	return b
}
