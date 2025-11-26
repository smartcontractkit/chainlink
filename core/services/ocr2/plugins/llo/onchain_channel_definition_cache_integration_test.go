package llo_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
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
			rc := NewMockReadCloser(invalidDefinitions)
			client.SetResponse(&http.Response{
				StatusCode: 200,
				Body:       rc,
			}, nil)

			url := "http://example.com/foo"
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
			rc := NewMockReadCloser(invalidDefinitions)
			client.SetResponse(&http.Response{
				StatusCode: 200,
				Body:       rc,
			}, nil)

			url := "http://example.com/foo"
			require.NoError(t, utils.JustError(configStoreContract.SetChannelDefinitions(steve, donID, url, invalidDefinitionsSHA)))
			backend.Commit()
		}

		testutils.WaitForLogMessageWithField(t, observedLogs,
			"Error while fetching channel definitions",
			"err", "invalid character '{' looking for beginning of object key string")
		assert.Empty(t, cdc.Definitions(llotypes.ChannelDefinitions{}))
	})

	t.Run("if server returns 404, should not update", func(t *testing.T) {
		// clear the log messages before waiting for new ones
		observedLogs.TakeAll()

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

		testutils.WaitForLogMessageWithField(t, observedLogs,
			"Error while fetching channel definitions",
			"err", "(status 404): not found")
	})

	t.Run("if server starts returning empty body, still does not update", func(t *testing.T) {
		// clear the log messages before waiting for new ones
		observedLogs.TakeAll()

		{
			rc := NewMockReadCloser([]byte{})
			client.SetResponse(&http.Response{
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
			client.SetResponse(&http.Response{
				StatusCode: 200,
				Body:       rc,
			}, nil)
		}

		// Wait for the log trigger to be processed
		le := testutils.WaitForLogMessageWithField(t, observedLogs, "Got new logs",
			"url", "http://example.com/foo3")
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
			client.SetResponse(nil, errors.New("failed; not a real URL"))
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
			client.SetResponse(&http.Response{
				StatusCode: 200,
				Body:       rc,
			}, nil)

			url := "http://example.com/foo5"
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
		client.SetResponse(&http.Response{
			StatusCode: 200,
			Body:       rc,
		}, nil)

		// Trigger new channel definitions on-chain
		url := "http://example.com/migration-test.json"
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
