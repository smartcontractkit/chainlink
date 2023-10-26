package client

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/testutils"
)

func TestNewSendOnlyNode(t *testing.T) {
	t.Parallel()

	urlFormat := "http://user:%s@testurl.com"
	password := "pass"
	u, err := url.Parse(fmt.Sprintf(urlFormat, password))
	require.NoError(t, err)
	redacted := fmt.Sprintf(urlFormat, "xxxxx")
	lggr := logger.TestLogger(t)
	name := "TestNewSendOnlyNode"
	chainID := types.RandomID()
	client := newMockSendOnlyClient[types.ID](t)

	node := NewSendOnlyNode(lggr, *u, name, chainID, client)
	assert.NotNil(t, node)

	// Must contain name & url with redacted password
	assert.Contains(t, node.String(), fmt.Sprintf("%s:%s", name, redacted))
	assert.Equal(t, node.ConfiguredChainID(), chainID)
}

func TestStartSendOnlyNode(t *testing.T) {
	t.Parallel()
	t.Run("becomes unusable if initial dial fails", func(t *testing.T) {
		t.Parallel()
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.WarnLevel)
		client := newMockSendOnlyClient[types.ID](t)
		client.On("Close").Once()
		expectedError := errors.New("some http error")
		client.On("DialHTTP").Return(expectedError).Once()
		s := NewSendOnlyNode(lggr, url.URL{}, t.Name(), types.RandomID(), client)

		defer func() { assert.NoError(t, s.Close()) }()
		err := s.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, nodeStateUnusable, s.State())
		testutils.RequireLogMessage(t, observedLogs, "Dial failed: SendOnly Node is unusable")
	})
	t.Run("Default ChainID produces warn and skips checks", func(t *testing.T) {
		t.Parallel()
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.WarnLevel)
		client := newMockSendOnlyClient[types.ID](t)
		client.On("Close").Once()
		client.On("DialHTTP").Return(nil).Once()
		s := NewSendOnlyNode(lggr, url.URL{}, t.Name(), types.NewIDFromInt(0), client)

		defer func() { assert.NoError(t, s.Close()) }()
		err := s.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, nodeStateAlive, s.State())
		testutils.RequireLogMessage(t, observedLogs, "sendonly rpc ChainID verification skipped")
	})
	t.Run("Can recover from chainID verification failure", func(t *testing.T) {
		t.Parallel()
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.WarnLevel)
		client := newMockSendOnlyClient[types.ID](t)
		client.On("Close").Once()
		client.On("DialHTTP").Return(nil)
		expectedError := errors.New("failed to get chain ID")
		chainID := types.RandomID()
		const failuresCount = 2
		client.On("ChainID", mock.Anything).Return(types.RandomID(), expectedError).Times(failuresCount)
		client.On("ChainID", mock.Anything).Return(chainID, nil)

		s := NewSendOnlyNode(lggr, url.URL{}, t.Name(), chainID, client)

		defer func() { assert.NoError(t, s.Close()) }()
		err := s.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, nodeStateUnreachable, s.State())
		testutils.WaitForLogMessageCount(t, observedLogs, fmt.Sprintf("Verify failed: %v", expectedError), failuresCount)
		testutils.AssertEventually(t, func() bool {
			return s.State() == nodeStateAlive
		})
	})
	t.Run("Can remover from chainID mismatch", func(t *testing.T) {
		t.Parallel()
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.WarnLevel)
		client := newMockSendOnlyClient[types.ID](t)
		client.On("Close").Once()
		client.On("DialHTTP").Return(nil).Once()
		configuredChainID := types.NewIDFromInt(11)
		rpcChainID := types.NewIDFromInt(20)
		const failuresCount = 2
		client.On("ChainID", mock.Anything).Return(rpcChainID, nil).Times(failuresCount)
		client.On("ChainID", mock.Anything).Return(configuredChainID, nil)
		s := NewSendOnlyNode(lggr, url.URL{}, t.Name(), configuredChainID, client)

		defer func() { assert.NoError(t, s.Close()) }()
		err := s.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, nodeStateInvalidChainID, s.State())
		testutils.WaitForLogMessageCount(t, observedLogs, "sendonly rpc ChainID doesn't match local chain ID", failuresCount)
		testutils.AssertEventually(t, func() bool {
			return s.State() == nodeStateAlive
		})
	})

	/*t.Run("Start with Random ChainID", func(t *testing.T) {
		chainID := types.RandomID()
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.WarnLevel)
		client := newMockSendOnlyClient[types.ID](t)
		s := NewSendOnlyNode(lggr, url.URL{}, t.Name(), chainID, client)
		defer func() { assert.NoError(t, s.Close()) }()
		err := s.Start(testutils.Context(t))
		assert.NoError(t, err)                 // No errors expected
		assert.Equal(t, 0, observedLogs.Len()) // No warnings expected
	})*/

	/**/
}

/*
func createSignedTx(t *testing.T, chainID *big.Int, nonce uint64, data []byte) *types.Transaction {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	sender, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	require.NoError(t, err)
	tx := types.NewTransaction(
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
	lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
	url := testutils.MustParseURL(t, "http://place.holder")
	s := evmclient.NewSendOnlyNode(lggr,
		*url,
		t.Name(),
		testutils.FixtureChainID).(evmclient.TestableSendOnlyNode)
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
}*/
