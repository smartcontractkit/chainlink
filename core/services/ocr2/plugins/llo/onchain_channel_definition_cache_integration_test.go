package llo_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/channeldefinitions"
	llotypes2 "github.com/smartcontractkit/chainlink/v2/core/services/llo/types"
)

type mockHTTPClient struct {
	responses map[string]*http.Response
	errors    map[string]error
	mu        sync.Mutex
}

func (h *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	url := req.URL.String()
	// Check for URL-specific response first
	if err, ok := h.errors[url]; ok {
		return nil, err
	}
	if resp, ok := h.responses[url]; ok {
		return resp, nil
	}
	// Fall back to default response (for backward compatibility with old tests)
	if err, ok := h.errors[""]; ok {
		return nil, err
	}
	if resp, ok := h.responses[""]; ok {
		return resp, nil
	}
	return nil, fmt.Errorf("no response configured for URL: %s", url)
}

func (h *mockHTTPClient) SetResponseForURL(url string, resp *http.Response, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.responses == nil {
		h.responses = make(map[string]*http.Response)
		h.errors = make(map[string]error)
	}
	if err != nil {
		h.errors[url] = err
		delete(h.responses, url)
	} else {
		h.responses[url] = resp
		delete(h.errors, url)
	}
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

// extractChannelDefinitions unmarshals json.RawMessage and merges all channel definitions from source definitions into a single map
func extractChannelDefinitions(defsJSON json.RawMessage) llotypes.ChannelDefinitions {
	var sourceDefs map[uint32]llotypes2.SourceDefinition
	if err := json.Unmarshal(defsJSON, &sourceDefs); err != nil {
		return make(llotypes.ChannelDefinitions)
	}
	result := make(llotypes.ChannelDefinitions)
	for _, sourceDef := range sourceDefs {
		for channelID, def := range sourceDef.Definitions {
			result[channelID] = def
		}
	}
	return result
}

// countChannels unmarshals json.RawMessage and counts the total number of channels across all source definitions
func countChannels(defsJSON json.RawMessage) int {
	var sourceDefs map[uint32]llotypes2.SourceDefinition
	if err := json.Unmarshal(defsJSON, &sourceDefs); err != nil {
		return 0
	}
	count := 0
	for _, sourceDef := range sourceDefs {
		count += len(sourceDef.Definitions)
	}
	return count
}

func Test_ChannelDefinitionCache_Integration(t *testing.T) {
	t.Parallel()
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
				Tombstone: false,
				Source:    channeldefinitions.SourceOwner,
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
				Opts:      llotypes.ChannelOpts([]byte(`{"baseUSDFee":"0.1","expirationWindow":86400,"feedId":"0x0003aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","multiplier":"1000000000000000000"}`)),
				Tombstone: false,
				Source:    channeldefinitions.SourceOwner,
			},
		}
	)

	sampleDefinitionsJSON, err := json.MarshalIndent(sampleDefinitions, "", "  ")
	require.NoError(t, err)
	sampleDefinitionsSHA := sha3.Sum256(sampleDefinitionsJSON)

	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	db := pgtest.NewSqlxDB(t)
	const ETHMainnetChainSelector uint64 = 5009297550715157269
	orm := llo.NewChainScopedORM(db, ETHMainnetChainSelector)

	steve := evmtestutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := types.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	backend.Commit() // ensure starting block number at least 1

	ethClient := client.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID)

	configStoreAddress, _, configStoreContract, err := channel_config_store.DeployChannelConfigStore(steve, backend.Client())
	require.NoError(t, err)

	backend.Commit()

	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            1,
		BackfillBatchSize:        3,
		RPCBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	ht := headstest.NewSimulatedHeadTracker(ethClient, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(testutils.SimulatedChainID, db, lggr), ethClient, lggr, ht, lpOpts)
	servicetest.Run(t, lp)

	client := &mockHTTPClient{}
	donID := rand.Uint32()

	cdc := channeldefinitions.NewChannelDefinitionCache(lggr, orm, client, lp, configStoreAddress, donID, 0, channeldefinitions.WithLogPollInterval(100*time.Millisecond))
	servicetest.Run(t, cdc)

	t.Run("before any logs, returns empty Definitions", func(t *testing.T) {
		assert.Empty(t, cdc.Definitions(llotypes.ChannelDefinitions{}))
	})

	t.Run("with sha mismatch, should not update", func(t *testing.T) {
		// clear the log messages
		t.Cleanup(func() { observedLogs.TakeAll() })

		{
			url := "http://example.com/foo"
			rc := NewMockReadCloser(invalidDefinitions)
			client.SetResponseForURL(url, &http.Response{
				StatusCode: 200,
				Body:       rc,
			}, nil)

			require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, sampleDefinitionsSHA)))

			backend.Commit()
		}

		testutils.WaitForLogMessageWithField(t, observedLogs,
			"Error while fetching channel definitions",
			"err", "SHA3 mismatch for channel definitions")

		assert.Empty(t, cdc.Definitions(llotypes.ChannelDefinitions{}))
	})

	t.Run("after correcting sha with new channel definitions set on-chain, but with invalid JSON at url, should not update", func(t *testing.T) {
		// clear the log messages before waiting for new ones
		observedLogs.TakeAll()

		{
			url := "http://example.com/foo"
			rc := NewMockReadCloser(invalidDefinitions)
			client.SetResponseForURL(url, &http.Response{
				StatusCode: 200,
				Body:       rc,
			}, nil)

			require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, invalidDefinitionsSHA)))
			backend.Commit()
		}

		testutils.WaitForLogMessageWithField(t, observedLogs,
			"Error while fetching channel definitions",
			"err", "failed to fetch channel definitions: failed to decode channel definitions JSON from http://example.com/foo: invalid character '{' looking for beginning of object key string")
		assert.Empty(t, cdc.Definitions(llotypes.ChannelDefinitions{}))
	})

	t.Run("if server returns 404, should not update", func(t *testing.T) {
		// clear the log messages before waiting for new ones
		observedLogs.TakeAll()

		{
			rc := NewMockReadCloser([]byte("not found"))
			url := "http://example.com/foo3"
			client.SetResponseForURL(url, &http.Response{
				StatusCode: 404,
				Body:       rc,
			}, nil)

			require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, sampleDefinitionsSHA)))
			backend.Commit()
		}

		testutils.WaitForLogMessageWithField(t, observedLogs,
			"Error while fetching channel definitions",
			"err", "(status 404): not found")
	})

	t.Run("if server starts returning empty body, still does not update", func(t *testing.T) {
		// clear the log messages before waiting for new ones
		observedLogs.TakeAll()

		{
			rc := NewMockReadCloser([]byte{})
			url := "http://example.com/foo3"
			client.SetResponseForURL(url, &http.Response{
				StatusCode: 200,
				Body:       rc,
			}, nil)
		}

		testutils.WaitForLogMessageWithField(t, observedLogs,
			"Error while fetching channel definitions", "err", "failed to fetch channel definitions: SHA3 mismatch for channel definitions")
	})

	t.Run("when URL starts returning valid JSON, updates even without needing new logs", func(t *testing.T) {
		// clear the log messages before waiting for new ones
		observedLogs.TakeAll()

		{
			rc := NewMockReadCloser(sampleDefinitionsJSON)
			url := "http://example.com/foo3"
			client.SetResponseForURL(url, &http.Response{
				StatusCode: 200,
				Body:       rc,
			}, nil)
		}

		// Wait for the log trigger to be processed
		le := testutils.WaitForLogMessage(t, observedLogs, "Set channel definitions for source")
		fields := le.ContextMap()
		assert.Contains(t, fields, "source")
		assert.Contains(t, fields, "url")
		assert.Contains(t, fields, "sha")
		assert.Contains(t, fields, "blockNum")
		assert.NotContains(t, fields, "err")

		assert.Equal(t, channeldefinitions.SourceOwner, fields["source"])
		assert.Equal(t, "http://example.com/foo3", fields["url"])
		assert.Equal(t, hex.EncodeToString(sampleDefinitionsSHA[:]), fields["sha"])

		// Wait for definitions to be fetched and merged
		require.Eventually(t, func() bool {
			defs := cdc.Definitions(llotypes.ChannelDefinitions{})
			return len(defs) > 0
		}, 5*time.Second, 100*time.Millisecond, "definitions should be available")

		assert.Equal(t, sampleDefinitions, cdc.Definitions(llotypes.ChannelDefinitions{}))

		t.Run("latest channel definitions are persisted", func(t *testing.T) {
			// Wait for initial persistence to complete (persistLoop periodically persists source definitions)
			var prevOutcome *llotypes2.PersistedDefinitions
			require.Eventually(t, func() bool {
				loaded, err := orm.LoadChannelDefinitions(testutils.Context(t), configStoreAddress, donID)
				if err != nil || loaded == nil {
					return false
				}
				// Check if we have the expected number of channels across all sources
				if countChannels(loaded.Definitions) != len(sampleDefinitions) {
					return false
				}
				prevOutcome = loaded
				return true
			}, 5*time.Second, 100*time.Millisecond, "channel definitions should be persisted")
			require.NotNil(t, prevOutcome, "previous outcome should be loaded from database")

			// Simulate plugin behavior: call Definitions() with merged channel definitions from previous outcome
			// Definitions() merges source definitions with prev and returns the result
			// Persistence happens separately via persistLoop, which stores c.definitions.Sources
			_ = cdc.Definitions(extractChannelDefinitions(prevOutcome.Definitions))

			// Wait for persistence to complete after calling Definitions() with previous outcome
			var pd *llotypes2.PersistedDefinitions
			require.Eventually(t, func() bool {
				loaded, err := orm.LoadChannelDefinitions(testutils.Context(t), configStoreAddress, donID)
				if err != nil || loaded == nil {
					return false
				}
				// Check if we have the expected number of channels across all sources
				if countChannels(loaded.Definitions) != len(sampleDefinitions) {
					return false
				}
				pd = loaded
				return true
			}, 5*time.Second, 100*time.Millisecond, "channel definitions should be persisted after calling Definitions() with previous outcome")
			require.NotNil(t, pd)
			assert.Equal(t, ETHMainnetChainSelector, pd.ChainSelector)
			assert.Equal(t, configStoreAddress, pd.Address)
			// Verify the structure matches - extract and compare channel definitions
			extractedDefs := extractChannelDefinitions(pd.Definitions)
			assert.Len(t, extractedDefs, len(sampleDefinitions))
			for channelID, expectedDef := range sampleDefinitions {
				actualDef, exists := extractedDefs[channelID]
				assert.True(t, exists, "channel %d should exist", channelID)
				assert.Equal(t, expectedDef.ReportFormat, actualDef.ReportFormat)
				assert.Equal(t, expectedDef.Streams, actualDef.Streams)
			}
			assert.Equal(t, donID, pd.DonID)
			// persist() stores c.definitions.Sources (source definitions) to the database.
			// The version comes from c.definitions.Version which is set from the latest owner trigger in the source definitions.
			assert.GreaterOrEqual(t, pd.Version, prevOutcome.Version, "version should be >= previous outcome version")
		})

		t.Run("new cdc with same config should load from DB", func(t *testing.T) {
			// fromBlock far in the future to ensure logs are not used
			cdc2 := channeldefinitions.NewChannelDefinitionCache(logger.NullLogger, orm, client, lp, configStoreAddress, donID, 1000)
			servicetest.Run(t, cdc2)
			// Load the persisted source definitions from DB
			// The cache loads source definitions (map[uint32]types.SourceDefinition) from the database
			// Definitions(prev) merges the loaded source definitions from c.definitions.Sources with prev
			// Since source definitions are loaded from DB for a new cache, it should merge them with prev
			loaded, err := orm.LoadChannelDefinitions(testutils.Context(t), configStoreAddress, donID)
			require.NoError(t, err)
			require.NotNil(t, loaded)
			require.Equal(t, sampleDefinitions, extractChannelDefinitions(loaded.Definitions))
		})
	})

	t.Run("new log with invalid channel definitions URL does not affect old channel definitions", func(t *testing.T) {
		// clear the log messages
		observedLogs.TakeAll()
		{
			url := "not a real URL"
			require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, sampleDefinitionsSHA)))
			client.SetResponseForURL(url, nil, errors.New("failed; not a real URL"))
			backend.Commit()
		}

		testutils.WaitForLogMessageWithField(t, observedLogs, "Error while fetching channel definitions", "err", "invalid URI for request")
	})

	t.Run("new valid definitions set on-chain, should update", func(t *testing.T) {
		// clear the log messages before waiting for new ones
		observedLogs.TakeAll()

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
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			}
			var err error
			sampleDefinitionsJSON, err = json.MarshalIndent(sampleDefinitions, "", "  ")
			require.NoError(t, err)
			sampleDefinitionsSHA = sha3.Sum256(sampleDefinitionsJSON)
			rc := NewMockReadCloser(sampleDefinitionsJSON)
			url := "http://example.com/foo5"
			client.SetResponseForURL(url, &http.Response{
				StatusCode: 200,
				Body:       rc,
			}, nil)

			require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, sampleDefinitionsSHA)))

			backend.Commit()
		}

		// Wait for the log trigger to be processed
		le := testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", "http://example.com/foo5")
		fields := le.ContextMap()
		assert.Contains(t, fields, "source")
		assert.Contains(t, fields, "url")
		assert.Contains(t, fields, "sha")
		assert.Contains(t, fields, "blockNum")
		assert.NotContains(t, fields, "err")

		assert.Equal(t, channeldefinitions.SourceOwner, fields["source"])
		assert.Equal(t, "http://example.com/foo5", fields["url"])
		assert.Equal(t, hex.EncodeToString(sampleDefinitionsSHA[:]), fields["sha"])

		// Wait for definitions to be fetched and merged
		require.Eventually(t, func() bool {
			defs := cdc.Definitions(llotypes.ChannelDefinitions{})
			return len(defs) == len(sampleDefinitions)
		}, 5*time.Second, 100*time.Millisecond, "definitions should be updated")

		assert.Equal(t, sampleDefinitions, cdc.Definitions(llotypes.ChannelDefinitions{}))
	})

	t.Run("latest channel definitions are persisted and overwrite previous value", func(t *testing.T) {
		// Wait for initial persistence to complete (persistLoop periodically persists source definitions)
		var prev *llotypes2.PersistedDefinitions
		require.Eventually(t, func() bool {
			loaded, err := orm.LoadChannelDefinitions(testutils.Context(t), configStoreAddress, donID)
			if err != nil || loaded == nil {
				return false
			}
			// Check if we have the expected number of channels across all sources
			// Definitions is a map[uint32]types.SourceDefinition, so we need to count channels across all sources
			if countChannels(loaded.Definitions) != len(sampleDefinitions) {
				return false
			}
			prev = loaded
			return true
		}, 5*time.Second, 100*time.Millisecond, "latest channel definitions should be loaded from database")
		require.NotNil(t, prev, "latest channel definitions should be loaded from database")

		// Simulate plugin behavior: call Definitions() with merged channel definitions from previous outcome
		// Definitions() merges source definitions with prev and returns the result
		// Persistence happens separately via persistLoop, which stores c.definitions.Sources
		_ = cdc.Definitions(extractChannelDefinitions(prev.Definitions))

		// Wait for persistence to complete after calling Definitions() with previous outcome
		var pd *llotypes2.PersistedDefinitions
		require.Eventually(t, func() bool {
			loaded, err := orm.LoadChannelDefinitions(testutils.Context(t), configStoreAddress, donID)
			if err != nil || loaded == nil {
				return false
			}
			// Check if we have the expected number of channels across all sources
			if countChannels(loaded.Definitions) != len(sampleDefinitions) {
				return false
			}
			pd = loaded
			return true
		}, 5*time.Second, 100*time.Millisecond, "channel definitions should be persisted after calling Definitions() with previous outcome")
		require.NotNil(t, pd)
		assert.Equal(t, ETHMainnetChainSelector, pd.ChainSelector)
		assert.Equal(t, configStoreAddress, pd.Address)
		// Verify the structure matches - extract channel definitions from persisted source definitions
		extractedDefs := extractChannelDefinitions(pd.Definitions)
		assert.Len(t, extractedDefs, len(sampleDefinitions))
		for channelID, expectedDef := range sampleDefinitions {
			actualDef, exists := extractedDefs[channelID]
			assert.True(t, exists, "channel %d should exist", channelID)
			assert.Equal(t, expectedDef.ReportFormat, actualDef.ReportFormat)
			assert.Equal(t, expectedDef.Streams, actualDef.Streams)
		}
		assert.Equal(t, donID, pd.DonID)
		// persist() stores c.definitions.Sources (source definitions) to the database.
		// The version comes from c.definitions.Version which is set from the latest owner trigger in the source definitions.
		assert.GreaterOrEqual(t, pd.Version, prev.Version, "version should be >= previous outcome version")
	})

	t.Run("migration from SingleChannelDefinitionsFormat to MultiChannelDefinitionsFormat preserves metadata", func(t *testing.T) {
		migrationDonID := rand.Uint32()
		migrationVersion := uint32(1)
		migrationBlockNum := int64(1)
		migrationChainSelector := ETHMainnetChainSelector

		// Create old format definitions (just ChannelDefinitions, no source wrapper)
		oldFormatDefs := llotypes.ChannelDefinitions{
			1: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
					{StreamID: 2, Aggregator: llotypes.AggregatorMode},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
			2: {
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams: []llotypes.Stream{
					{StreamID: 3, Aggregator: llotypes.AggregatorQuote},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
		}

		oldFormatJSON, err := json.Marshal(oldFormatDefs)
		require.NoError(t, err)

		pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, chain_selector, don_id, definitions, block_num, version, updated_at, format)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7)
		`, configStoreAddress, migrationChainSelector, migrationDonID, oldFormatJSON, migrationBlockNum, migrationVersion, channeldefinitions.SingleChannelDefinitionsFormat)

		// Verify old format data in database
		oldPD, err := orm.LoadChannelDefinitions(testutils.Context(t), configStoreAddress, migrationDonID)
		require.NoError(t, err)
		require.NotNil(t, oldPD)
		assert.Equal(t, migrationChainSelector, oldPD.ChainSelector)
		assert.Equal(t, configStoreAddress, oldPD.Address)
		assert.Equal(t, migrationDonID, oldPD.DonID)
		assert.Equal(t, migrationVersion, oldPD.Version)
		assert.Equal(t, migrationBlockNum, oldPD.BlockNum)
		assert.Equal(t, channeldefinitions.SingleChannelDefinitionsFormat, oldPD.Format)

		// Create a new cache - it should load the metadata but not the definitions
		cdcMigration := channeldefinitions.NewChannelDefinitionCache(logger.NullLogger, orm, client, lp, configStoreAddress, migrationDonID, 0, channeldefinitions.WithLogPollInterval(100*time.Millisecond))
		servicetest.Run(t, cdcMigration)

		// Verify that metadata was loaded but definitions are empty
		// The cache should have loaded Version and BlockNum from the old format data
		defs := cdcMigration.Definitions(llotypes.ChannelDefinitions{})
		assert.Empty(t, defs, "definitions should be empty when format is SingleChannelDefinitionsFormat")

		// Now trigger new definitions to be persisted (this will migrate to new format)
		// Set up new definitions that will be fetched
		newDefinitions := llotypes.ChannelDefinitions{
			1: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
					{StreamID: 2, Aggregator: llotypes.AggregatorMode},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
			2: {
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams: []llotypes.Stream{
					{StreamID: 3, Aggregator: llotypes.AggregatorQuote},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
			3: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 4, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
		}

		newDefinitionsJSON, err := json.MarshalIndent(newDefinitions, "", "  ")
		require.NoError(t, err)
		newDefinitionsSHA := sha3.Sum256(newDefinitionsJSON)

		// Set up HTTP client to return new definitions
		rc := NewMockReadCloser(newDefinitionsJSON)
		url := "http://example.com/migration-test.json"
		client.SetResponseForURL(url, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)

		// Trigger new channel definitions on-chain
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, migrationDonID, url, newDefinitionsSHA)))
		backend.Commit()

		// Wait for definitions to be fetched and persisted
		require.Eventually(t, func() bool {
			defs := cdcMigration.Definitions(llotypes.ChannelDefinitions{})
			return len(defs) == len(newDefinitions)
		}, 5*time.Second, 100*time.Millisecond, "new definitions should be available")

		// Wait for persistence to complete
		var migratedPD *llotypes2.PersistedDefinitions
		require.Eventually(t, func() bool {
			var loaded *llotypes2.PersistedDefinitions
			if loaded, err = orm.LoadChannelDefinitions(testutils.Context(t), configStoreAddress, migrationDonID); err != nil || loaded == nil {
				return false
			}
			// Check that format has been migrated
			if loaded.Format != channeldefinitions.MultiChannelDefinitionsFormat {
				return false
			}
			migratedPD = loaded
			return true
		}, 5*time.Second, 100*time.Millisecond, "definitions should be migrated to MultiChannelDefinitionsFormat")

		require.NotNil(t, migratedPD)

		// Verify that all metadata is preserved
		assert.Equal(t, migrationChainSelector, migratedPD.ChainSelector, "ChainSelector should be preserved")
		assert.Equal(t, configStoreAddress, migratedPD.Address, "Address should be preserved")
		assert.Equal(t, migrationDonID, migratedPD.DonID, "DonID should be preserved")
		// Version should be preserved or updated (not reset to 0)
		assert.GreaterOrEqual(t, migratedPD.Version, migrationVersion, "Version should be preserved or updated, not reset")
		// BlockNum should be preserved or updated (not reset to 0)
		assert.GreaterOrEqual(t, migratedPD.BlockNum, migrationBlockNum, "BlockNum should be preserved or updated, not reset")

		// Verify format has been migrated
		assert.Equal(t, channeldefinitions.MultiChannelDefinitionsFormat, migratedPD.Format, "Format should be migrated to MultiChannelDefinitionsFormat")

		// Verify definitions are in new format (map[uint32]SourceDefinition)
		var newFormatDefs map[uint32]llotypes2.SourceDefinition
		err = json.Unmarshal(migratedPD.Definitions, &newFormatDefs)
		require.NoError(t, err)
		require.NotEmpty(t, newFormatDefs, "definitions should be in new format")

		// Verify the definitions contain the expected channels
		extractedDefs := extractChannelDefinitions(migratedPD.Definitions)
		assert.Len(t, extractedDefs, len(newDefinitions))
		for channelID, expectedDef := range newDefinitions {
			actualDef, exists := extractedDefs[channelID]
			assert.True(t, exists, "channel %d should exist", channelID)
			assert.Equal(t, expectedDef.ReportFormat, actualDef.ReportFormat)
			assert.Equal(t, expectedDef.Streams, actualDef.Streams)
		}
	})
}

func Test_ChannelDefinitionCache_OwnerAndAdderMerging(t *testing.T) {
	t.Parallel()
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	db := pgtest.NewSqlxDB(t)
	const ETHMainnetChainSelector uint64 = 5009297550715157269
	orm := llo.NewChainScopedORM(db, ETHMainnetChainSelector)

	steve := evmtestutils.MustNewSimTransactor(t) // config contract deployer and owner
	// Create adder accounts before creating backend
	adder1 := evmtestutils.MustNewSimTransactor(t)
	adder2 := evmtestutils.MustNewSimTransactor(t)
	genesisData := types.GenesisAlloc{
		steve.From:  {Balance: assets.Ether(1000).ToInt()},
		adder1.From: {Balance: assets.Ether(1000).ToInt()},
		adder2.From: {Balance: assets.Ether(1000).ToInt()},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	backend.Commit() // ensure starting block number at least 1

	ethClient := client.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID)

	configStoreAddress, _, configStoreContract, err := channel_config_store.DeployChannelConfigStore(steve, backend.Client())
	require.NoError(t, err)

	backend.Commit()

	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            1,
		BackfillBatchSize:        3,
		RPCBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	ht := headstest.NewSimulatedHeadTracker(ethClient, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(testutils.SimulatedChainID, db, lggr), ethClient, lggr, ht, lpOpts)
	servicetest.Run(t, lp)

	client := &mockHTTPClient{}
	donID := rand.Uint32()

	cdc := channeldefinitions.NewChannelDefinitionCache(lggr, orm, client, lp, configStoreAddress, donID, 0, channeldefinitions.WithLogPollInterval(100*time.Millisecond))
	servicetest.Run(t, cdc)

	// Configure adders on the contract
	// Adder IDs start from 1000
	adder1ID := uint32(1001)
	adder2ID := uint32(1002)

	require.NoError(t, utils.JustError(configStoreContract.SetChannelAdderAddress(steve, adder1ID, adder1.From)))
	backend.Commit()
	require.NoError(t, utils.JustError(configStoreContract.SetChannelAdderAddress(steve, adder2ID, adder2.From)))
	backend.Commit()

	// Enable adders
	require.NoError(t, utils.JustError(configStoreContract.SetChannelAdder(steve, donID, adder1ID, true)))
	backend.Commit()
	require.NoError(t, utils.JustError(configStoreContract.SetChannelAdder(steve, donID, adder2ID, true)))
	backend.Commit()

	t.Run("adder can add new channels", func(t *testing.T) {
		observedLogs.TakeAll()

		adder1Definitions := llotypes.ChannelDefinitions{
			100: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    adder1ID,
				Tombstone: false,
			},
			101: {
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams: []llotypes.Stream{
					{StreamID: 2, Aggregator: llotypes.AggregatorMode},
				},
				Source:    adder1ID,
				Tombstone: false,
			},
		}

		adder1DefinitionsJSON, err := json.MarshalIndent(adder1Definitions, "", "  ")
		require.NoError(t, err)
		adder1DefinitionsSHA := sha3.Sum256(adder1DefinitionsJSON)

		url := "http://example.com/adder1-defs.json"
		rc := NewMockReadCloser(adder1DefinitionsJSON)
		client.SetResponseForURL(url, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder1, donID, adder1ID, url, adder1DefinitionsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url)

		require.Eventually(t, func() bool {
			defs := cdc.Definitions(llotypes.ChannelDefinitions{})
			return len(defs) >= 2
		}, 5*time.Second, 100*time.Millisecond, "adder definitions should be available")

		defs := cdc.Definitions(llotypes.ChannelDefinitions{})
		assert.Equal(t, adder1Definitions[100], defs[100])
		assert.Equal(t, adder1Definitions[101], defs[101])
	})

	t.Run("adder cannot overwrite existing owner channels", func(t *testing.T) {
		observedLogs.TakeAll()

		// Owner sets channel definitions first
		ownerDefinitions := llotypes.ChannelDefinitions{
			200: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
		}

		ownerDefinitionsJSON, err := json.MarshalIndent(ownerDefinitions, "", "  ")
		require.NoError(t, err)
		ownerDefinitionsSHA := sha3.Sum256(ownerDefinitionsJSON)

		url := "http://example.com/owner-defs.json"
		rc := NewMockReadCloser(ownerDefinitionsJSON)
		client.SetResponseForURL(url, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, ownerDefinitionsSHA)))
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url)

		require.Eventually(t, func() bool {
			defs := cdc.Definitions(llotypes.ChannelDefinitions{})
			return len(defs) >= 1 && defs[200].Source == channeldefinitions.SourceOwner
		}, 5*time.Second, 100*time.Millisecond, "owner definitions should be available")

		// Now adder tries to add the same channel ID
		observedLogs.TakeAll()

		adderAttemptDefinitions := llotypes.ChannelDefinitions{
			200: {
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams: []llotypes.Stream{
					{StreamID: 999, Aggregator: llotypes.AggregatorQuote},
				},
				Source:    adder1ID,
				Tombstone: false,
			},
		}

		adderAttemptDefinitionsJSON, err := json.MarshalIndent(adderAttemptDefinitions, "", "  ")
		require.NoError(t, err)
		adderAttemptDefinitionsSHA := sha3.Sum256(adderAttemptDefinitionsJSON)

		url2 := "http://example.com/adder-attempt.json"
		rc = NewMockReadCloser(adderAttemptDefinitionsJSON)
		client.SetResponseForURL(url2, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder1, donID, adder1ID, url2, adderAttemptDefinitionsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url2)

		// Wait a bit for processing
		time.Sleep(500 * time.Millisecond)

		// Verify adder's definition was skipped - owner's definition should still be there
		defs := cdc.Definitions(llotypes.ChannelDefinitions{})
		assert.Equal(t, channeldefinitions.SourceOwner, defs[200].Source, "channel 200 should still be from owner")
		assert.Equal(t, ownerDefinitions[200].Streams, defs[200].Streams, "channel 200 should have owner's streams")

		// Check for conflict warning log
		testutils.WaitForLogMessageWithField(t, observedLogs, "channel adder conflict",
			"channelID", "200")
	})

	t.Run("adder cannot overwrite existing adder channels", func(t *testing.T) {
		observedLogs.TakeAll()

		// First adder adds a channel
		adder1Defs := llotypes.ChannelDefinitions{
			300: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    adder1ID,
				Tombstone: false,
			},
		}

		adder1DefsJSON, err := json.MarshalIndent(adder1Defs, "", "  ")
		require.NoError(t, err)
		adder1DefsSHA := sha3.Sum256(adder1DefsJSON)

		url := "http://example.com/adder1-channel300.json"
		rc := NewMockReadCloser(adder1DefsJSON)
		client.SetResponseForURL(url, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder1, donID, adder1ID, url, adder1DefsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url)

		require.Eventually(t, func() bool {
			defs := cdc.Definitions(llotypes.ChannelDefinitions{})
			return defs[300].Source == adder1ID
		}, 5*time.Second, 100*time.Millisecond, "adder1 channel 300 should be available")

		// Second adder tries to add the same channel ID
		observedLogs.TakeAll()

		adder2Defs := llotypes.ChannelDefinitions{
			300: {
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams: []llotypes.Stream{
					{StreamID: 999, Aggregator: llotypes.AggregatorQuote},
				},
				Source:    adder2ID,
				Tombstone: false,
			},
		}

		adder2DefsJSON, err := json.MarshalIndent(adder2Defs, "", "  ")
		require.NoError(t, err)
		adder2DefsSHA := sha3.Sum256(adder2DefsJSON)

		url2 := "http://example.com/adder2-channel300.json"
		rc = NewMockReadCloser(adder2DefsJSON)
		client.SetResponseForURL(url2, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder2, donID, adder2ID, url2, adder2DefsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url2)

		// Wait a bit for processing
		time.Sleep(500 * time.Millisecond)

		// Verify second adder's definition was skipped
		defs := cdc.Definitions(llotypes.ChannelDefinitions{})
		assert.Equal(t, adder1ID, defs[300].Source, "channel 300 should still be from adder1")
		assert.Equal(t, adder1Defs[300].Streams, defs[300].Streams, "channel 300 should have adder1's streams")

		// Check for conflict warning log
		testutils.WaitForLogMessageWithField(t, observedLogs, "channel adder conflict",
			"channelID", "300")
	})

	t.Run("adder cannot tombstone channels", func(t *testing.T) {
		observedLogs.TakeAll()

		// Adder tries to add a channel with Tombstone: true
		adderTombstoneDefs := llotypes.ChannelDefinitions{
			400: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    adder1ID,
				Tombstone: true, // Adders cannot tombstone
			},
		}

		adderTombstoneDefsJSON, err := json.MarshalIndent(adderTombstoneDefs, "", "  ")
		require.NoError(t, err)
		adderTombstoneDefsSHA := sha3.Sum256(adderTombstoneDefsJSON)

		url := "http://example.com/adder-tombstone.json"
		rc := NewMockReadCloser(adderTombstoneDefsJSON)
		client.SetResponseForURL(url, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder1, donID, adder1ID, url, adderTombstoneDefsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url)

		// Wait a bit for processing
		time.Sleep(500 * time.Millisecond)

		// Verify tombstone channel was skipped
		defs := cdc.Definitions(llotypes.ChannelDefinitions{})
		_, exists := defs[400]
		assert.False(t, exists, "channel 400 should not exist (tombstone skipped)")

		// Check for tombstone warning log
		testutils.WaitForLogMessageWithField(t, observedLogs, "invalid channel tombstone",
			"channelID", "400")
	})

	t.Run("owner can overwrite adder channels", func(t *testing.T) {
		observedLogs.TakeAll()

		// Adder adds a channel first
		adderDefs := llotypes.ChannelDefinitions{
			500: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    adder1ID,
				Tombstone: false,
			},
		}

		adderDefsJSON, err := json.MarshalIndent(adderDefs, "", "  ")
		require.NoError(t, err)
		adderDefsSHA := sha3.Sum256(adderDefsJSON)

		url := "http://example.com/adder-channel500.json"
		rc := NewMockReadCloser(adderDefsJSON)
		client.SetResponseForURL(url, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder1, donID, adder1ID, url, adderDefsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url)

		require.Eventually(t, func() bool {
			defs := cdc.Definitions(llotypes.ChannelDefinitions{})
			return defs[500].Source == adder1ID
		}, 5*time.Second, 100*time.Millisecond, "adder channel 500 should be available")

		// Owner sets new definitions that include the same channel ID with different values
		observedLogs.TakeAll()

		ownerDefs := llotypes.ChannelDefinitions{
			500: {
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams: []llotypes.Stream{
					{StreamID: 999, Aggregator: llotypes.AggregatorQuote},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
		}

		ownerDefsJSON, err := json.MarshalIndent(ownerDefs, "", "  ")
		require.NoError(t, err)
		ownerDefsSHA := sha3.Sum256(ownerDefsJSON)

		url2 := "http://example.com/owner-overwrite.json"
		rc = NewMockReadCloser(ownerDefsJSON)
		client.SetResponseForURL(url2, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url2, ownerDefsSHA)))
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url2)

		require.Eventually(t, func() bool {
			defs := cdc.Definitions(llotypes.ChannelDefinitions{})
			return defs[500].Source == channeldefinitions.SourceOwner
		}, 5*time.Second, 100*time.Millisecond, "owner should have overwritten channel 500")

		// Verify owner's definition overwrote the adder's
		defs := cdc.Definitions(llotypes.ChannelDefinitions{})
		assert.Equal(t, channeldefinitions.SourceOwner, defs[500].Source, "channel 500 should be from owner")
		assert.Equal(t, ownerDefs[500].Streams, defs[500].Streams, "channel 500 should have owner's streams")
	})

	t.Run("owner cannot implicitly remove channels", func(t *testing.T) {
		observedLogs.TakeAll()

		// Start with channels from owner and adders
		ownerDefs := llotypes.ChannelDefinitions{
			600: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
			601: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 2, Aggregator: llotypes.AggregatorMode},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
		}

		ownerDefsJSON, err := json.MarshalIndent(ownerDefs, "", "  ")
		require.NoError(t, err)
		ownerDefsSHA := sha3.Sum256(ownerDefsJSON)

		url := "http://example.com/owner-channels600-601.json"
		rc := NewMockReadCloser(ownerDefsJSON)
		client.SetResponseForURL(url, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, ownerDefsSHA)))
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url)

		// Adder adds a channel
		observedLogs.TakeAll()

		adderDefs := llotypes.ChannelDefinitions{
			602: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 3, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    adder1ID,
				Tombstone: false,
			},
		}

		adderDefsJSON, err := json.MarshalIndent(adderDefs, "", "  ")
		require.NoError(t, err)
		adderDefsSHA := sha3.Sum256(adderDefsJSON)

		url2 := "http://example.com/adder-channel602.json"
		rc = NewMockReadCloser(adderDefsJSON)
		client.SetResponseForURL(url2, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder1, donID, adder1ID, url2, adderDefsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url2)

		require.Eventually(t, func() bool {
			defs := cdc.Definitions(llotypes.ChannelDefinitions{})
			return len(defs) >= 3
		}, 5*time.Second, 100*time.Millisecond, "all channels should be available")

		// Owner sets new definitions that exclude channel 600
		observedLogs.TakeAll()

		ownerDefsUpdated := llotypes.ChannelDefinitions{
			601: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 2, Aggregator: llotypes.AggregatorMode},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
			// Channel 600 is excluded - should be removed
		}

		ownerDefsUpdatedJSON, err := json.MarshalIndent(ownerDefsUpdated, "", "  ")
		require.NoError(t, err)
		ownerDefsUpdatedSHA := sha3.Sum256(ownerDefsUpdatedJSON)

		url3 := "http://example.com/owner-removed-600.json"
		rc = NewMockReadCloser(ownerDefsUpdatedJSON)
		client.SetResponseForURL(url3, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url3, ownerDefsUpdatedSHA)))
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url3)

		require.Eventually(t, func() bool {
			defs := cdc.Definitions(ownerDefs)
			_, has600 := defs[600]
			_, has601 := defs[601]
			_, has602 := defs[602]
			return has600 && has601 && has602
		}, 5*time.Second, 100*time.Millisecond, "channel 600, 601 and 602 should still be present")
	})

	t.Run("owner can remove channels explicitly", func(t *testing.T) {
		observedLogs.TakeAll()

		// Start with channels from owner and adders
		ownerDefs := llotypes.ChannelDefinitions{
			600: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
			601: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 2, Aggregator: llotypes.AggregatorMode},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
		}

		// Owner sets new definitions that exclude channel 600
		ownerDefsUpdated := llotypes.ChannelDefinitions{
			600: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: true,
			},
		}

		ownerDefsUpdatedJSON, err := json.MarshalIndent(ownerDefsUpdated, "", "  ")
		require.NoError(t, err)
		ownerDefsUpdatedSHA := sha3.Sum256(ownerDefsUpdatedJSON)

		url3 := "http://example.com/owner-removed-600.json"
		rc := NewMockReadCloser(ownerDefsUpdatedJSON)
		client.SetResponseForURL(url3, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url3, ownerDefsUpdatedSHA)))
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url3)

		require.Eventually(t, func() bool {
			defs := cdc.Definitions(ownerDefs)
			def600 := defs[600]
			_, has601 := defs[601]
			_, has602 := defs[602]
			return def600.Tombstone && has601 && has602
		}, 5*time.Second, 100*time.Millisecond, "channel 600 should be removed, 601 and 602 should still be present")
	})

	t.Run("multiple adders can add different channels", func(t *testing.T) {
		observedLogs.TakeAll()

		// Adder1 adds channels
		adder1Defs := llotypes.ChannelDefinitions{
			700: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    adder1ID,
				Tombstone: false,
			},
		}

		adder1DefsJSON, err := json.MarshalIndent(adder1Defs, "", "  ")
		require.NoError(t, err)
		adder1DefsSHA := sha3.Sum256(adder1DefsJSON)

		url := "http://example.com/adder1-channel700.json"
		rc := NewMockReadCloser(adder1DefsJSON)
		client.SetResponseForURL(url, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder1, donID, adder1ID, url, adder1DefsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url)

		// Adder2 adds different channels
		observedLogs.TakeAll()

		adder2Defs := llotypes.ChannelDefinitions{
			701: {
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams: []llotypes.Stream{
					{StreamID: 2, Aggregator: llotypes.AggregatorMode},
				},
				Source:    adder2ID,
				Tombstone: false,
			},
		}

		adder2DefsJSON, err := json.MarshalIndent(adder2Defs, "", "  ")
		require.NoError(t, err)
		adder2DefsSHA := sha3.Sum256(adder2DefsJSON)

		url2 := "http://example.com/adder2-channel701.json"
		rc = NewMockReadCloser(adder2DefsJSON)
		client.SetResponseForURL(url2, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder2, donID, adder2ID, url2, adder2DefsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url2)

		require.Eventually(t, func() bool {
			defs := cdc.Definitions(llotypes.ChannelDefinitions{})
			_, has700 := defs[700]
			_, has701 := defs[701]
			return has700 && has701
		}, 5*time.Second, 100*time.Millisecond, "both adder channels should be available")

		// Verify all channels from both adders are present
		defs := cdc.Definitions(llotypes.ChannelDefinitions{})
		assert.Equal(t, adder1ID, defs[700].Source, "channel 700 should be from adder1")
		assert.Equal(t, adder2ID, defs[701].Source, "channel 701 should be from adder2")
		assert.Equal(t, adder1Defs[700].Streams, defs[700].Streams, "channel 700 should have adder1's streams")
		assert.Equal(t, adder2Defs[701].Streams, defs[701].Streams, "channel 701 should have adder2's streams")
	})

	t.Run("adder limit enforcement", func(t *testing.T) {
		observedLogs.TakeAll()

		// Create definitions with more than MaxChannelsPerAdder (100) channels
		tooManyDefs := make(llotypes.ChannelDefinitions)
		for i := uint32(800); i < 800+channeldefinitions.MaxChannelsPerAdder+1; i++ {
			tooManyDefs[i] = llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: i, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    adder1ID,
				Tombstone: false,
			}
		}

		tooManyDefsJSON, err := json.MarshalIndent(tooManyDefs, "", "  ")
		require.NoError(t, err)
		tooManyDefsSHA := sha3.Sum256(tooManyDefsJSON)

		url := "http://example.com/too-many-channels.json"
		rc := NewMockReadCloser(tooManyDefsJSON)
		client.SetResponseForURL(url, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder1, donID, adder1ID, url, tooManyDefsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url)

		// Wait a bit for processing
		time.Sleep(500 * time.Millisecond)

		// Call Definitions() to trigger the merge and error logging
		_ = cdc.Definitions(llotypes.ChannelDefinitions{})

		// Verify error is logged and channels are not merged
		testutils.WaitForLogMessageWithField(t, observedLogs, "adder limit exceeded, skipping remaining definitions for source",
			"source", strconv.FormatUint(uint64(adder1ID), 10))

		// Verify no channels above the limit were added
		defs := cdc.Definitions(llotypes.ChannelDefinitions{})
		var addedDefinitionsCount int
		for _, def := range defs {
			if def.Source == adder1ID {
				addedDefinitionsCount++
			}
		}
		require.Equal(t, channeldefinitions.MaxChannelsPerAdder, addedDefinitionsCount)
	})

	t.Run("deterministic processing order", func(t *testing.T) {
		observedLogs.TakeAll()

		// Add definitions from owner and adders at different block numbers
		// We'll add them in a specific order and verify the final result respects block/log ordering

		// First, adder1 adds channel 900
		adder1Defs := llotypes.ChannelDefinitions{
			900: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams: []llotypes.Stream{
					{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
				},
				Source:    adder1ID,
				Tombstone: false,
			},
		}

		adder1DefsJSON, err := json.MarshalIndent(adder1Defs, "", "  ")
		require.NoError(t, err)
		adder1DefsSHA := sha3.Sum256(adder1DefsJSON)

		url := "http://example.com/adder1-channel900.json"
		rc := NewMockReadCloser(adder1DefsJSON)
		client.SetResponseForURL(url, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		_, err = configStoreContract.AddChannelDefinitions(adder1, donID, adder1ID, url, adder1DefsSHA)
		require.NoError(t, err)
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url)

		// Then owner adds channel 900 (should overwrite)
		observedLogs.TakeAll()

		ownerDefs := llotypes.ChannelDefinitions{
			900: {
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams: []llotypes.Stream{
					{StreamID: 999, Aggregator: llotypes.AggregatorQuote},
				},
				Source:    channeldefinitions.SourceOwner,
				Tombstone: false,
			},
		}

		ownerDefsJSON, err := json.MarshalIndent(ownerDefs, "", "  ")
		require.NoError(t, err)
		ownerDefsSHA := sha3.Sum256(ownerDefsJSON)

		url2 := "http://example.com/owner-channel900.json"
		rc = NewMockReadCloser(ownerDefsJSON)
		client.SetResponseForURL(url2, &http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)
		require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url2, ownerDefsSHA)))
		backend.Commit()

		testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", url2)

		require.Eventually(t, func() bool {
			defs := cdc.Definitions(llotypes.ChannelDefinitions{})
			return defs[900].Source == channeldefinitions.SourceOwner
		}, 5*time.Second, 100*time.Millisecond, "owner should have overwritten channel 900")

		// Verify final result respects ordering (owner's definition should win)
		defs := cdc.Definitions(llotypes.ChannelDefinitions{})
		assert.Equal(t, channeldefinitions.SourceOwner, defs[900].Source, "channel 900 should be from owner (processed later)")
		assert.Equal(t, ownerDefs[900].Streams, defs[900].Streams, "channel 900 should have owner's streams")
	})
}
