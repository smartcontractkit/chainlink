package llo_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
)

type mockHTTPClient struct {
	resp *http.Response
	err  error
	mu   sync.Mutex
}

func (h *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.resp, h.err
}

func (h *mockHTTPClient) SetResponse(resp *http.Response, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.resp = resp
	h.err = err
}

type MockReadCloser struct {
	data   []byte
	mu     sync.Mutex
	reader *bytes.Reader
}

func NewMockReadCloser(data []byte) *MockReadCloser {
	return &MockReadCloser{
		data:   data,
		reader: bytes.NewReader(data),
	}
}

// Read reads from the underlying data
func (m *MockReadCloser) Read(p []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.reader.Read(p)
}

// Close resets the reader to the beginning of the data
func (m *MockReadCloser) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, err := m.reader.Seek(0, io.SeekStart)
	return err
}

func Test_ChannelDefinitionCache_Integration(t *testing.T) {
	var (
		invalidDefinitions    = []byte(`{{{`)
		invalidDefinitionsSHA = sha3.Sum256(invalidDefinitions)

		sampleDefinitions = llotypes.ChannelDefinitions{
			1: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{
						StreamID:   1,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   2,
						Aggregator: llotypes.AggregatorMode,
					},
				},
			},
			2: {
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams: []llotypes.Stream{
					{
						StreamID:   1,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   2,
						Aggregator: llotypes.AggregatorMedian,
					},
					{
						StreamID:   3,
						Aggregator: llotypes.AggregatorQuote,
					},
				},
				Opts: llotypes.ChannelOpts([]byte(`{"baseUSDFee":"0.1","expirationWindow":86400,"feedId":"0x0003aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","multiplier":"1000000000000000000"}`)),
			},
		}
	)

	sampleDefinitionsJSON, err := json.MarshalIndent(sampleDefinitions, "", "  ")
	require.NoError(t, err)
	sampleDefinitionsSHA := sha3.Sum256(sampleDefinitionsJSON)

	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	db := pgtest.NewSqlxDB(t)
	const ETHMainnetChainSelector uint64 = 5009297550715157269
	orm := llo.NewORM(db, ETHMainnetChainSelector)

	steve := testutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := core.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	backend.Commit() // ensure starting block number at least 1

	ethClient := client.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID)

	configStoreAddress, _, configStoreContract, err := channel_config_store.DeployChannelConfigStore(steve, backend)
	require.NoError(t, err)

	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            1,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	ht := headtracker.NewSimulatedHeadTracker(ethClient, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(testutils.SimulatedChainID, db, lggr), ethClient, lggr, ht, lpOpts)
	servicetest.Run(t, lp)

	client := &mockHTTPClient{}
	donID := rand.Uint32()

	cdc := llo.NewChannelDefinitionCache(lggr, orm, client, lp, configStoreAddress, donID, 0, llo.WithLogPollInterval(100*time.Millisecond))
	servicetest.Run(t, cdc)

	t.Run("before any logs, returns empty Definitions", func(t *testing.T) {
		assert.Empty(t, cdc.Definitions())
	})

	{
		rc := NewMockReadCloser(invalidDefinitions)
		client.SetResponse(&http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)

		url := "http://example.com/foo"
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, sampleDefinitionsSHA)))

		backend.Commit()
	}

	t.Run("with sha mismatch, should not update", func(t *testing.T) {
		// clear the log messages
		t.Cleanup(func() { observedLogs.TakeAll() })

		testutils.WaitForLogMessage(t, observedLogs, "Got new channel definitions from chain")
		le := testutils.WaitForLogMessage(t, observedLogs, "Error while fetching channel definitions")
		fields := le.ContextMap()
		assert.Contains(t, fields, "err")
		assert.Equal(t, fmt.Sprintf("SHA3 mismatch: expected %x, got %x", sampleDefinitionsSHA, invalidDefinitionsSHA), fields["err"])

		assert.Empty(t, cdc.Definitions())
	})

	{
		rc := NewMockReadCloser(invalidDefinitions)
		client.SetResponse(&http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)

		url := "http://example.com/foo"
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, invalidDefinitionsSHA)))
		backend.Commit()
	}

	t.Run("after correcting sha with new channel definitions set on-chain, but with invalid JSON at url, should not update", func(t *testing.T) {
		// clear the log messages
		t.Cleanup(func() { observedLogs.TakeAll() })

		testutils.WaitForLogMessage(t, observedLogs, "Got new channel definitions from chain")
		testutils.WaitForLogMessageWithField(t, observedLogs, "Error while fetching channel definitions", "err", "failed to decode JSON: invalid character '{' looking for beginning of object key string")

		assert.Empty(t, cdc.Definitions())
	})

	{
		rc := NewMockReadCloser([]byte("not found"))
		client.SetResponse(&http.Response{
			StatusCode: 404,
			Body:       rc,
		}, nil)

		url := "http://example.com/foo3"
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, sampleDefinitionsSHA)))
		backend.Commit()
	}

	t.Run("if server returns 404, should not update", func(t *testing.T) {
		// clear the log messages
		t.Cleanup(func() { observedLogs.TakeAll() })

		testutils.WaitForLogMessageWithField(t, observedLogs, "Error while fetching channel definitions", "err", "got error from http://example.com/foo3: (status code: 404, response body: not found)")
	})

	{
		rc := NewMockReadCloser([]byte{})
		client.SetResponse(&http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
	}

	t.Run("if server starts returning empty body, still does not update", func(t *testing.T) {
		// clear the log messages
		t.Cleanup(func() { observedLogs.TakeAll() })

		testutils.WaitForLogMessageWithField(t, observedLogs, "Error while fetching channel definitions", "err", fmt.Sprintf("SHA3 mismatch: expected %x, got %x", sampleDefinitionsSHA, sha3.Sum256([]byte{})))
	})

	{
		rc := NewMockReadCloser(sampleDefinitionsJSON)
		client.SetResponse(&http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
	}

	t.Run("when URL starts returning valid JSON, updates even without needing new logs", func(t *testing.T) {
		// clear the log messages
		t.Cleanup(func() { observedLogs.TakeAll() })

		le := testutils.WaitForLogMessage(t, observedLogs, "Set new channel definitions")
		fields := le.ContextMap()
		assert.Contains(t, fields, "version")
		assert.Contains(t, fields, "url")
		assert.Contains(t, fields, "sha")
		assert.Contains(t, fields, "donID")
		assert.NotContains(t, fields, "err")

		assert.Equal(t, uint32(3), fields["version"])
		assert.Equal(t, "http://example.com/foo3", fields["url"])
		assert.Equal(t, fmt.Sprintf("%x", sampleDefinitionsSHA), fields["sha"])
		assert.Equal(t, donID, fields["donID"])

		assert.Equal(t, sampleDefinitions, cdc.Definitions())

		t.Run("latest channel definitions are persisted", func(t *testing.T) {
			pd, err := orm.LoadChannelDefinitions(testutils.Context(t), configStoreAddress, donID)
			require.NoError(t, err)
			assert.Equal(t, ETHMainnetChainSelector, pd.ChainSelector)
			assert.Equal(t, configStoreAddress, pd.Address)
			assert.Equal(t, sampleDefinitions, pd.Definitions)
			assert.Equal(t, donID, pd.DonID)
			assert.Equal(t, uint32(3), pd.Version)
		})

		t.Run("new cdc with same config should load from DB", func(t *testing.T) {
			// fromBlock far in the future to ensure logs are not used
			cdc2 := llo.NewChannelDefinitionCache(lggr, orm, client, lp, configStoreAddress, donID, 1000)
			servicetest.Run(t, cdc2)

			assert.Equal(t, sampleDefinitions, cdc.Definitions())
		})
	})

	{
		url := "not a real URL"
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, sampleDefinitionsSHA)))

		backend.Commit()

		client.SetResponse(nil, errors.New("failed; not a real URL"))
	}

	t.Run("new log with invalid channel definitions URL does not affect old channel definitions", func(t *testing.T) {
		// clear the log messages
		t.Cleanup(func() { observedLogs.TakeAll() })

		le := testutils.WaitForLogMessage(t, observedLogs, "Error while fetching channel definitions")
		fields := le.ContextMap()
		assert.Contains(t, fields, "err")
		assert.Equal(t, "error making http request: failed; not a real URL", fields["err"])
	})

	{
		// add a new definition, it should get loaded
		sampleDefinitions[3] = llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatJSON,
			Streams: []llotypes.Stream{
				{
					StreamID:   6,
					Aggregator: llotypes.AggregatorMedian,
				},
			},
		}
		var err error
		sampleDefinitionsJSON, err = json.MarshalIndent(sampleDefinitions, "", "  ")
		require.NoError(t, err)
		sampleDefinitionsSHA = sha3.Sum256(sampleDefinitionsJSON)
		rc := NewMockReadCloser(sampleDefinitionsJSON)
		client.SetResponse(&http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)

		url := "http://example.com/foo5"
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, sampleDefinitionsSHA)))

		backend.Commit()
	}

	t.Run("successfully updates to new channel definitions with new log", func(t *testing.T) {
		t.Cleanup(func() { observedLogs.TakeAll() })

		le := testutils.WaitForLogMessage(t, observedLogs, "Set new channel definitions")
		fields := le.ContextMap()
		assert.Contains(t, fields, "version")
		assert.Contains(t, fields, "url")
		assert.Contains(t, fields, "sha")
		assert.Contains(t, fields, "donID")
		assert.NotContains(t, fields, "err")

		assert.Equal(t, uint32(5), fields["version"])
		assert.Equal(t, "http://example.com/foo5", fields["url"])
		assert.Equal(t, fmt.Sprintf("%x", sampleDefinitionsSHA), fields["sha"])
		assert.Equal(t, donID, fields["donID"])

		assert.Equal(t, sampleDefinitions, cdc.Definitions())
	})

	t.Run("latest channel definitions are persisted and overwrite previous value", func(t *testing.T) {
		pd, err := orm.LoadChannelDefinitions(testutils.Context(t), configStoreAddress, donID)
		require.NoError(t, err)
		assert.Equal(t, ETHMainnetChainSelector, pd.ChainSelector)
		assert.Equal(t, configStoreAddress, pd.Address)
		assert.Equal(t, sampleDefinitions, pd.Definitions)
		assert.Equal(t, donID, pd.DonID)
		assert.Equal(t, uint32(5), pd.Version)
	})
}
