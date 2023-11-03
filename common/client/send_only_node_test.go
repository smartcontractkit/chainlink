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

	"github.com/smartcontractkit/chainlink-relay/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
		err := s.Start(tests.Context(t))
		require.NoError(t, err)

		assert.Equal(t, nodeStateUnusable, s.State())
		tests.RequireLogMessage(t, observedLogs, "Dial failed: SendOnly Node is unusable")
	})
	t.Run("Default ChainID produces warn and skips checks", func(t *testing.T) {
		t.Parallel()
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.WarnLevel)
		client := newMockSendOnlyClient[types.ID](t)
		client.On("Close").Once()
		client.On("DialHTTP").Return(nil).Once()
		s := NewSendOnlyNode(lggr, url.URL{}, t.Name(), types.NewIDFromInt(0), client)

		defer func() { assert.NoError(t, s.Close()) }()
		err := s.Start(tests.Context(t))
		require.NoError(t, err)

		assert.Equal(t, nodeStateAlive, s.State())
		tests.RequireLogMessage(t, observedLogs, "sendonly rpc ChainID verification skipped")
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
		err := s.Start(tests.Context(t))
		require.NoError(t, err)

		assert.Equal(t, nodeStateUnreachable, s.State())
		tests.WaitForLogMessageCount(t, observedLogs, fmt.Sprintf("Verify failed: %v", expectedError), failuresCount)
		tests.AssertEventually(t, func() bool {
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
		err := s.Start(tests.Context(t))
		require.NoError(t, err)

		assert.Equal(t, nodeStateInvalidChainID, s.State())
		tests.WaitForLogMessageCount(t, observedLogs, "sendonly rpc ChainID doesn't match local chain ID", failuresCount)
		tests.AssertEventually(t, func() bool {
			return s.State() == nodeStateAlive
		})
	})
}
