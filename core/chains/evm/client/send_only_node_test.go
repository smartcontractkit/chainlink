package client_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client/mocks"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestNewSendOnlyNode(t *testing.T) {
	t.Parallel()

	urlFormat := "http://user:%s@testurl.com"
	password := "pass"
	url := testutils.MustParseURL(t, fmt.Sprintf(urlFormat, password))
	redacted := fmt.Sprintf(urlFormat, "xxxxx")
	lggr := logger.TestLogger(t)
	name := "TestNewSendOnlyNode"
	chainID := testutils.NewRandomEVMChainID()

	node := evmclient.NewSendOnlyNode(lggr, *url, name, chainID)
	assert.NotNil(t, node)

	// Must contain name & url with redacted password
	assert.Contains(t, node.String(), fmt.Sprintf("%s:%s", name, redacted))
	assert.Equal(t, node.ChainID(), chainID)
}

func TestStartSendOnlyNode(t *testing.T) {
	t.Parallel()

	t.Run("Start with Random ChainID", func(t *testing.T) {
		t.Parallel()
		chainID := testutils.NewRandomEVMChainID()
		r := chainIDResp{chainID.Int64(), nil}
		url := r.newHTTPServer(t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.WarnLevel)
		s := evmclient.NewSendOnlyNode(lggr, *url, t.Name(), chainID)
		defer s.Close()
		err := s.Start(testutils.Context(t))
		assert.NoError(t, err)                 // No errors expected
		assert.Equal(t, 0, observedLogs.Len()) // No warnings expected
	})

	t.Run("Start with ChainID=0", func(t *testing.T) {
		t.Parallel()
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.WarnLevel)
		chainID := testutils.FixtureChainID
		r := chainIDResp{chainID.Int64(), nil}
		url := r.newHTTPServer(t)
		s := evmclient.NewSendOnlyNode(lggr, *url, t.Name(), testutils.FixtureChainID)

		defer s.Close()
		err := s.Start(testutils.Context(t))
		assert.NoError(t, err)
		// getChainID() should return Error if ChainID = 0
		// This should get converted into a warning from Start()
		testutils.WaitForLogMessage(t, observedLogs, "ChainID verification skipped")
	})
}

func createSignedTx(t *testing.T, chainID *big.Int, nonce uint64, data []byte) *types.Transaction {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	sender, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	require.NoError(t, err)
	tx := types.NewTransaction(
		nonce, sender.From,
		assets.Ether(100),
		21000, big.NewInt(1000000000), data,
	)
	signedTx, err := sender.Signer(sender.From, tx)
	require.NoError(t, err)
	return signedTx
}

func TestSendTransaction(t *testing.T) {
	t.Parallel()

	chainID := testutils.FixtureChainID
	lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
	url := testutils.MustParseURL(t, "http://place.holder")
	s := evmclient.NewSendOnlyNode(lggr,
		*url,
		t.Name(),
		testutils.FixtureChainID).(evmclient.TestableSendOnlyNode)
	require.NotNil(t, s)

	signedTx := createSignedTx(t, chainID, 1, []byte{1, 2, 3})

	mockTxSender := new(mocks.TxSender)
	mockTxSender.Test(t)

	mockTxSender.On("SendTransaction", mock.Anything, mock.MatchedBy(
		func(tx *types.Transaction) bool {
			if tx.Nonce() != uint64(1) {
				return false
			}
			return true
		},
	)).Once().Return(nil)
	s.SetEthClient(nil, mockTxSender)

	err := s.SendTransaction(testutils.TestCtx(t), signedTx)
	assert.NoError(t, err)
	testutils.WaitForLogMessage(t, observedLogs, "SendOnly RPC call")
	mockTxSender.AssertExpectations(t)
}

func TestBatchCallContext(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	chainID := testutils.FixtureChainID
	url := testutils.MustParseURL(t, "http://place.holder")
	s := evmclient.NewSendOnlyNode(
		lggr,
		*url, "TestBatchCallContext",
		chainID).(evmclient.TestableSendOnlyNode)

	blockNum := hexutil.EncodeBig(big.NewInt(42))
	req := []rpc.BatchElem{
		{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{blockNum, true},
			Result: &types.Block{},
		},
		{
			Method: "method",
			Args:   []interface{}{1, false}},
	}

	mockBatchSender := new(mocks.BatchSender)
	mockBatchSender.Test(t)
	mockBatchSender.On("BatchCallContext", mock.Anything,
		mock.MatchedBy(
			func(b []rpc.BatchElem) bool {
				return len(b) == 2 &&
					b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == blockNum && b[0].Args[1] == true
			})).Return(nil).Once().Return(nil)

	s.SetEthClient(mockBatchSender, nil)

	err := s.BatchCallContext(context.Background(), req)
	assert.NoError(t, err)
	mockBatchSender.AssertExpectations(t)
}
