package client_test

import (
	"fmt"
	"math/big"
	"net/url"
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

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestNewSendOnlyNode(t *testing.T) {
	t.Parallel()

	urlFormat := "http://user:%s@testurl.com"
	password := "pass"
	url := testutils.MustParseURL(t, fmt.Sprintf(urlFormat, password))
	redacted := fmt.Sprintf(urlFormat, "xxxxx")
	lggr := logger.Test(t)
	name := "TestNewSendOnlyNode"
	chainID := testutils.NewRandomEVMChainID()

	node := client.NewSendOnlyNode(lggr, *url, name, chainID)
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
		lggr, observedLogs := logger.TestObserved(t, zap.WarnLevel)
		s := client.NewSendOnlyNode(lggr, *url, t.Name(), chainID)
		defer func() { assert.NoError(t, s.Close()) }()
		err := s.Start(testutils.Context(t))
		assert.NoError(t, err)                 // No errors expected
		assert.Equal(t, 0, observedLogs.Len()) // No warnings expected
	})

	t.Run("Start with ChainID=0", func(t *testing.T) {
		t.Parallel()
		lggr, observedLogs := logger.TestObserved(t, zap.WarnLevel)
		chainID := testutils.FixtureChainID
		r := chainIDResp{chainID.Int64(), nil}
		url := r.newHTTPServer(t)
		s := client.NewSendOnlyNode(lggr, *url, t.Name(), testutils.FixtureChainID)

		defer func() { assert.NoError(t, s.Close()) }()
		err := s.Start(testutils.Context(t))
		assert.NoError(t, err)
		// If ChainID = 0, this should get converted into a warning from Start()
		testutils.WaitForLogMessage(t, observedLogs, "ChainID verification skipped")
	})

	t.Run("becomes unusable (and remains undialed) if initial dial fails", func(t *testing.T) {
		t.Parallel()
		lggr, observedLogs := logger.TestObserved(t, zap.WarnLevel)
		invalidURL := url.URL{Scheme: "some rubbish", Host: "not a valid host"}
		s := client.NewSendOnlyNode(lggr, invalidURL, t.Name(), testutils.FixtureChainID)

		defer func() { assert.NoError(t, s.Close()) }()
		err := s.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.False(t, client.IsDialed(s))
		testutils.RequireLogMessage(t, observedLogs, "Dial failed: EVM SendOnly Node is unusable")
	})
}

func createSignedTx(t *testing.T, chainID *big.Int, nonce uint64, data []byte) *types.Transaction {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	sender, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	require.NoError(t, err)
	tx := cltest.NewLegacyTransaction(
		nonce, sender.From,
		assets.Ether(100).ToInt(),
		21000, big.NewInt(1000000000), data,
	)
	signedTx, err := sender.Signer(sender.From, tx)
	require.NoError(t, err)
	return signedTx
}

func TestSendTransaction(t *testing.T) {
	t.Parallel()

	chainID := testutils.FixtureChainID
	lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
	url := testutils.MustParseURL(t, "http://place.holder")
	s := client.NewSendOnlyNode(lggr,
		*url,
		t.Name(),
		testutils.FixtureChainID).(client.TestableSendOnlyNode)
	require.NotNil(t, s)

	signedTx := createSignedTx(t, chainID, 1, []byte{1, 2, 3})

	mockTxSender := mocks.NewTxSender(t)
	mockTxSender.On("SendTransaction", mock.Anything, mock.MatchedBy(
		func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(1)
		},
	)).Once().Return(nil)
	s.SetEthClient(nil, mockTxSender)

	err := s.SendTransaction(testutils.Context(t), signedTx)
	assert.NoError(t, err)
	testutils.WaitForLogMessage(t, observedLogs, "SendOnly RPC call")
}

func TestBatchCallContext(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	chainID := testutils.FixtureChainID
	url := testutils.MustParseURL(t, "http://place.holder")
	s := client.NewSendOnlyNode(
		lggr,
		*url, "TestBatchCallContext",
		chainID).(client.TestableSendOnlyNode)

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

	mockBatchSender := mocks.NewBatchSender(t)
	mockBatchSender.On("BatchCallContext", mock.Anything,
		mock.MatchedBy(
			func(b []rpc.BatchElem) bool {
				return len(b) == 2 &&
					b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == blockNum && b[0].Args[1].(bool)
			})).Return(nil).Once().Return(nil)

	s.SetEthClient(mockBatchSender, nil)

	err := s.BatchCallContext(testutils.Context(t), req)
	assert.NoError(t, err)
}
