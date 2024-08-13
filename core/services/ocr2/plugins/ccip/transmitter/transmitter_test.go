package transmitter

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	FixtureChainID = *testutils.FixtureChainID
	Password       = testutils.Password
)

func newMockTxStrategy(t *testing.T) *commontxmmocks.TxStrategy {
	return commontxmmocks.NewTxStrategy(t)
}

func Test_DefaultTransmitter_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := NewKeyStore(t, db).Eth()

	_, fromAddress := MustInsertRandomKey(t, ethKeyStore)

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := fromAddress
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	transmitter, err := ocrcommon.NewTransmitter(
		txm,
		[]common.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
	)
	require.NoError(t, err)

	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         gasLimit,
		ForwarderAddress: common.Address{},
		Meta:             nil,
		Strategy:         strategy,
	}).Return(txmgr.Tx{}, nil).Once()
	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := NewKeyStore(t, db).Eth()

	_, fromAddress := MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress2 := MustInsertRandomKey(t, ethKeyStore)

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := common.Address{}
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	transmitter, err := ocrcommon.NewTransmitter(
		txm,
		[]common.Address{fromAddress, fromAddress2},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
	)
	require.NoError(t, err)

	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         gasLimit,
		ForwarderAddress: common.Address{},
		Meta:             nil,
		Strategy:         strategy,
	}).Return(txmgr.Tx{}, nil).Once()
	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		FromAddress:      fromAddress2,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         gasLimit,
		ForwarderAddress: common.Address{},
		Meta:             nil,
		Strategy:         strategy,
	}).Return(txmgr.Tx{}, nil).Once()
	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction_Round_Robin_Error(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := NewKeyStore(t, db).Eth()

	fromAddress := common.Address{}

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := common.Address{}
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	transmitter, err := ocrcommon.NewTransmitter(
		txm,
		[]common.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
	)
	require.NoError(t, err)
	require.Error(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction_No_Keystore_Error(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := NewKeyStore(t, db).Eth()

	_, fromAddress := MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress2 := MustInsertRandomKey(t, ethKeyStore)

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := common.Address{}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	_, err := ocrcommon.NewTransmitter(
		txm,
		[]common.Address{fromAddress, fromAddress2},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		nil,
	)
	require.Error(t, err)
}

func Test_Transmitter_With_StatusChecker_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := NewKeyStore(t, db).Eth()

	_, fromAddress := MustInsertRandomKey(t, ethKeyStore)

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := fromAddress
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	idempotencyKey := "1-0"
	txMeta := &txmgr.TxMeta{MessageIDs: []string{"1"}}

	transmitter, err := NewTransmitterWithStatusChecker(
		txm,
		[]common.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
	)
	require.NoError(t, err)

	// This case is for when the message ID was not found in the status checker
	txm.On("GetTransactionStatus", mock.Anything, idempotencyKey).Return(types.Unknown, errors.New("dummy")).Once()

	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		IdempotencyKey:   &idempotencyKey,
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         gasLimit,
		ForwarderAddress: common.Address{},
		Meta:             txMeta,
		Strategy:         strategy,
	}).Return(txmgr.Tx{}, nil).Once()

	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, txMeta))
	txm.AssertExpectations(t)
}

func NewKeyStore(t testing.TB, ds sqlutil.DataSource) keystore.Master {
	ctx := testutils.Context(t)
	keystore := keystore.NewInMemory(ds, utils.FastScryptParams, logger.TestLogger(t))
	require.NoError(t, keystore.Unlock(ctx, Password))
	return keystore
}

type RandomKey struct {
	Nonce    int64
	Disabled bool

	chainIDs []ubig.Big // nil: Fixture, set empty for none
}

func (r RandomKey) MustInsert(t testing.TB, keystore keystore.Eth) (ethkey.KeyV2, common.Address) {
	ctx := testutils.Context(t)
	chainIDs := r.chainIDs
	if chainIDs == nil {
		chainIDs = []ubig.Big{*ubig.New(&FixtureChainID)}
	}

	key := MustGenerateRandomKey(t)
	keystore.XXXTestingOnlyAdd(ctx, key)

	for _, cid := range chainIDs {
		require.NoError(t, keystore.Add(ctx, key.Address, cid.ToInt()))
		require.NoError(t, keystore.Enable(ctx, key.Address, cid.ToInt()))
		if r.Disabled {
			require.NoError(t, keystore.Disable(ctx, key.Address, cid.ToInt()))
		}
	}

	return key, key.Address
}

func MustInsertRandomKey(t testing.TB, keystore keystore.Eth, chainIDs ...ubig.Big) (ethkey.KeyV2, common.Address) {
	r := RandomKey{}
	if len(chainIDs) > 0 {
		r.chainIDs = chainIDs
	}
	return r.MustInsert(t, keystore)
}

func MustGenerateRandomKey(t testing.TB) ethkey.KeyV2 {
	key, err := ethkey.NewV2()
	require.NoError(t, err)
	return key
}
