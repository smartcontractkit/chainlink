package channeldefinitions

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/types"
)

type mockLogPoller struct {
	latestBlock     logpoller.Block
	latestBlockErr  error
	filteredLogs    []logpoller.Log
	adderLogs       []logpoller.Log
	ownerLogs       []logpoller.Log
	filteredLogsErr error

	unregisteredFilterNames []string
}

func (m *mockLogPoller) RegisterFilter(ctx context.Context, filter logpoller.Filter) error {
	return nil
}
func (m *mockLogPoller) LatestBlock(ctx context.Context) (logpoller.Block, error) {
	return m.latestBlock, m.latestBlockErr
}
func (m *mockLogPoller) FilteredLogs(ctx context.Context, filter []query.Expression, limitAndSort query.LimitAndSort, queryName string) ([]logpoller.Log, error) {
	// Return different logs based on query name to simulate separate adder/owner queries
	if queryName == "ChannelDefinitionCachePoller - NewAdderChannelDefinition" {
		if len(m.adderLogs) > 0 {
			return m.adderLogs, m.filteredLogsErr
		}
		return m.filteredLogs, m.filteredLogsErr
	}
	if queryName == "ChannelDefinitionCachePoller - NewOwnerChannelDefinition" {
		if len(m.ownerLogs) > 0 {
			return m.ownerLogs, m.filteredLogsErr
		}
		return m.filteredLogs, m.filteredLogsErr
	}
	return m.filteredLogs, m.filteredLogsErr
}
func (m *mockLogPoller) UnregisterFilter(ctx context.Context, name string) error {
	m.unregisteredFilterNames = append(m.unregisteredFilterNames, name)
	return nil
}

var _ HTTPClient = &mockHTTPClient{}

type mockHTTPClient struct {
	resp *http.Response
	err  error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

var _ ChannelDefinitionCacheORM = &mockCDCORM{}

type mockCDCORM struct {
	err error

	lastPersistedAddr     common.Address
	lastPersistedDonID    uint32
	lastPersistedVersion  uint32
	lastPersistedDfns     map[uint32]types.SourceDefinition
	lastPersistedBlockNum int64
	lastPersistedFormat   uint32
}

func (m *mockCDCORM) LoadChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) (pd *types.PersistedDefinitions, err error) {
	panic("not implemented")
}
func (m *mockCDCORM) StoreChannelDefinitions(ctx context.Context, addr common.Address, donID, version uint32, dfns json.RawMessage, blockNum int64, format uint32) (err error) {
	m.lastPersistedAddr = addr
	m.lastPersistedDonID = donID
	m.lastPersistedVersion = version
	m.lastPersistedBlockNum = blockNum
	m.lastPersistedFormat = format
	// Unmarshal the json.RawMessage to store in lastPersistedDfns for test assertions
	if err := json.Unmarshal(dfns, &m.lastPersistedDfns); err != nil {
		return err
	}
	return m.err
}

func (m *mockCDCORM) CleanupChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) (err error) {
	panic("not implemented")
}

func makeLog(t *testing.T, donID, version uint32, url string, sha [32]byte) logpoller.Log {
	data := makeLogData(t, donID, version, url, sha)
	return logpoller.Log{EventSig: NewChannelDefinition, Topics: [][]byte{NewChannelDefinition[:], makeDonIDTopic(donID)}, Data: data, BlockNumber: int64(version) + 1000}
}

func makeLogData(t *testing.T, donID, version uint32, url string, sha [32]byte) []byte {
	event := channelConfigStoreABI.Events[newChannelDefinitionEventName]
	// donID is indexed
	// version, url, sha
	data, err := event.Inputs.NonIndexed().Pack(version, url, sha)
	require.NoError(t, err)
	return data
}

func makeAdderLog(t *testing.T, donID, adderID uint32, url string, sha [32]byte, blockNumber int64) logpoller.Log {
	data := makeAdderLogData(t, donID, adderID, url, sha)
	return logpoller.Log{EventSig: ChannelDefinitionAdded, Topics: [][]byte{ChannelDefinitionAdded[:], makeDonIDTopic(donID), makeDonIDTopic(adderID)}, Data: data, BlockNumber: blockNumber}
}

func makeAdderLogData(t *testing.T, donID, adderID uint32, url string, sha [32]byte) []byte {
	event := channelConfigStoreABI.Events[channelDefinitionAddedEventName]
	// donID and adderID are indexed (in Topics)
	// url, sha are non-indexed (in Data)
	data, err := event.Inputs.NonIndexed().Pack(url, sha)
	require.NoError(t, err)
	return data
}

func makeDonIDTopic(donID uint32) []byte {
	return common.BigToHash(big.NewInt(int64(donID))).Bytes()
}

// drainChannel drains all values from a channel
func drainChannel[T any](ch chan T) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}

// collectTriggers collects all available triggers from a channel up to maxCount
func collectTriggers(ch chan types.Trigger, maxCount int) []types.Trigger {
	triggers := make([]types.Trigger, 0, maxCount)
	for i := 0; i < maxCount; i++ {
		select {
		case trigger := <-ch:
			triggers = append(triggers, trigger)
		default:
			return triggers
		}
	}
	return triggers
}

// makeChannelDefinition creates a simple channel definition for testing
func makeChannelDefinition(channelID uint32, source uint32) llotypes.ChannelDefinition {
	return llotypes.ChannelDefinition{
		ReportFormat: llotypes.ReportFormatJSON,
		Streams:      []llotypes.Stream{{StreamID: channelID, Aggregator: llotypes.AggregatorMedian}},
		Source:       source,
		Tombstone:    false,
	}
}

// makeChannelDefinitionWithFeedID creates a channel definition with a FeedID in options for testing
func makeChannelDefinitionWithFeedID(channelID uint32, source uint32, feedID common.Hash) llotypes.ChannelDefinition {
	optsJSON := fmt.Sprintf(`{"feedId":"%s"}`, feedID.Hex())
	return llotypes.ChannelDefinition{
		ReportFormat: llotypes.ReportFormatJSON,
		Streams:      []llotypes.Stream{{StreamID: channelID, Aggregator: llotypes.AggregatorMedian}},
		Source:       source,
		Tombstone:    false,
		Opts:         llotypes.ChannelOpts(optsJSON),
	}
}

// addChannelDefinitions adds channel definitions to the given map for a range of channel IDs
func addChannelDefinitions(defs llotypes.ChannelDefinitions, startID, endID uint32, source uint32) {
	for i := startID; i <= endID; i++ {
		defs[i] = makeChannelDefinition(i, source)
	}
}

func Test_ChannelDefinitionCache(t *testing.T) {
	donID := rand.Uint32()

	t.Run("Definitions", func(t *testing.T) {
		// NOTE: this is covered more thoroughly in the integration tests
		prev := llotypes.ChannelDefinitions(map[llotypes.ChannelID]llotypes.ChannelDefinition{
			1: {
				ReportFormat: llotypes.ReportFormat(43),
				Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 2, Aggregator: llotypes.AggregatorMode}, {StreamID: 3, Aggregator: llotypes.AggregatorQuote}},
				Opts:         llotypes.ChannelOpts{1, 2, 3},
				Source:       SourceOwner,
			},
		})

		// Test that Definitions() returns prev when sourceDefinitions is empty
		cdc := &channelDefinitionCache{
			lggr: logger.TestSugared(t),
			definitions: Definitions{
				Sources: make(map[uint32]types.SourceDefinition),
			},
			orm: &mockCDCORM{}, // Required for persist() call in Definitions()
		}

		result := cdc.Definitions(prev)
		require.Equal(t, prev, result)

		// Test merging from sourceDefinitions
		adderID := uint32(100)
		sourceDefs := llotypes.ChannelDefinitions{
			2: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams:      []llotypes.Stream{{StreamID: 2, Aggregator: llotypes.AggregatorMedian}},
				Source:       adderID,
			},
		}
		cdc.definitions.Sources[adderID] = types.SourceDefinition{
			Trigger: types.Trigger{
				Source:   adderID,
				BlockNum: 1000,
			},
			Definitions: sourceDefs,
		}

		result = cdc.Definitions(prev)
		// Should contain both prev channel 1 and adder channel 2
		require.Contains(t, result, llotypes.ChannelID(1))
		require.Contains(t, result, llotypes.ChannelID(2))
		require.Equal(t, SourceOwner, result[1].Source)
		require.Equal(t, adderID, result[2].Source)

		// Test tombstone removal
		tombstoneDefs := llotypes.ChannelDefinitions{
			1: {
				ReportFormat: llotypes.ReportFormat(43),
				Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 2, Aggregator: llotypes.AggregatorMode}, {StreamID: 3, Aggregator: llotypes.AggregatorQuote}},
				Opts:         llotypes.ChannelOpts{1, 2, 3},
				Source:       SourceOwner,
				Tombstone:    false,
			},
			3: {
				ReportFormat: llotypes.ReportFormatJSON,
				Streams:      []llotypes.Stream{{StreamID: 3, Aggregator: llotypes.AggregatorMedian}},
				Source:       SourceOwner,
				Tombstone:    true,
			},
		}
		cdc.definitions.Sources[SourceOwner] = types.SourceDefinition{
			Trigger: types.Trigger{
				Source:   SourceOwner,
				BlockNum: 2000,
			},
			Definitions: tombstoneDefs,
		}

		result = cdc.Definitions(prev)
		// Tombstoned channel should be kept in definitions with Tombstone: true
		require.Contains(t, result, llotypes.ChannelID(3))
		require.True(t, result[3].Tombstone, "channel 3 should be tombstoned")
		// Channels 1 and 2 should still be present
		require.Contains(t, result, llotypes.ChannelID(1))
		require.Contains(t, result, llotypes.ChannelID(2))
	})

	t.Run("readLogs", func(t *testing.T) {
		lp := &mockLogPoller{latestBlockErr: sql.ErrNoRows}
		fetchTriggerCh := make(chan types.Trigger, 100)
		cdc := &channelDefinitionCache{
			donID:          donID,
			lp:             lp,
			lggr:           logger.TestSugared(t),
			fetchTriggerCh: fetchTriggerCh,
			definitions: Definitions{
				Sources: make(map[uint32]types.SourceDefinition),
			},
		}

		t.Run("skips if logpoller has no blocks", func(t *testing.T) {
			ctx := t.Context()
			err := cdc.readLogs(ctx)
			require.NoError(t, err)
		})
		t.Run("returns error on LatestBlock failure", func(t *testing.T) {
			ctx := t.Context()
			lp.latestBlockErr = errors.New("test error")

			err := cdc.readLogs(ctx)
			require.EqualError(t, err, "test error")
		})
		t.Run("does nothing if LatestBlock older or the same as current channel definitions block", func(t *testing.T) {
			ctx := t.Context()
			lp.latestBlockErr = nil
			lp.latestBlock = logpoller.Block{BlockNumber: 42}
			cdc.definitions.LastBlockNum = 43

			err := cdc.readLogs(ctx)
			require.NoError(t, err)
		})
		t.Run("returns error if FilteredLogs fails", func(t *testing.T) {
			ctx := t.Context()
			cdc.definitions.LastBlockNum = 0
			lp.filteredLogsErr = errors.New("test error 2")

			err := cdc.readLogs(ctx)
			require.EqualError(t, err, "test error 2")
		})
		t.Run("ignores logs with different topic", func(t *testing.T) {
			ctx := t.Context()
			lp.filteredLogsErr = nil
			// Set logs with different event signature (not NewChannelDefinition or ChannelDefinitionAdded)
			lp.ownerLogs = []logpoller.Log{{EventSig: common.Hash{1, 2, 3, 4}}}
			lp.adderLogs = []logpoller.Log{}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)
		})
		t.Run("logs warning and continues if log is malformed", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			cdc.definitions.LastBlockNum = 0
			cdc.initialBlockNum = 0
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			lp.latestBlockErr = nil
			lp.filteredLogsErr = nil
			// Set malformed owner log (has correct event sig but missing data)
			lp.ownerLogs = []logpoller.Log{{EventSig: NewChannelDefinition}}
			lp.adderLogs = []logpoller.Log{}

			err := cdc.readLogs(ctx)
			require.NoError(t, err, "should not return error for malformed log, should log warning and continue")
			// Should not send trigger for malformed log
			select {
			case <-fetchTriggerCh:
				t.Fatal("should not send trigger for malformed log")
			default:
				// Expected - no trigger
			}
		})
		t.Run("sends trigger on channel if FilteredLogs returns new event with a later version", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			cdc.definitions.LastBlockNum = 0
			cdc.initialBlockNum = 0
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			lp.latestBlockErr = nil
			lp.filteredLogsErr = nil
			// Set owner logs
			lp.ownerLogs = []logpoller.Log{makeLog(t, donID, uint32(43), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4})}
			// Set empty adder logs
			lp.adderLogs = []logpoller.Log{}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)

			// Check that fetch trigger was sent
			select {
			case trigger := <-fetchTriggerCh:
				require.Equal(t, SourceOwner, trigger.Source)
				require.Equal(t, uint32(43), trigger.Version)
				require.Equal(t, "http://example.com/xxx.json", trigger.URL)
				require.Equal(t, [32]byte{1, 2, 3, 4}, trigger.SHA)
			default:
				t.Fatal("expected fetch trigger signal in channel")
			}
		})
		t.Run("sends triggers for all logs", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			cdc.definitions.LastBlockNum = 0
			cdc.initialBlockNum = 0
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			lp.filteredLogsErr = nil
			// Set owner logs (readLogs calls FilteredLogs for owner logs)
			lp.ownerLogs = []logpoller.Log{
				makeLog(t, donID, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
				makeLog(t, donID, uint32(43), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
			}
			// Set empty adder logs
			lp.adderLogs = []logpoller.Log{}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)
			// Should receive triggers for both owner logs
			triggers := collectTriggers(fetchTriggerCh, 4)
			require.Len(t, triggers, 2, "expected 2 triggers")
			// Find the trigger with version 43 (latest)
			var found43 bool
			for _, trigger := range triggers {
				if trigger.Version == 43 {
					found43 = true
					break
				}
			}
			require.True(t, found43, "expected trigger with version 43")
		})
		t.Run("in case of multiple logs, sends triggers for all", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			cdc.definitions.LastBlockNum = 0
			cdc.initialBlockNum = 0
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			lp.filteredLogsErr = nil
			// Set owner logs (readLogs calls FilteredLogs for owner logs)
			lp.ownerLogs = []logpoller.Log{
				makeLog(t, donID, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
				makeLog(t, donID, uint32(45), "http://example.com/xxx2.json", [32]byte{2, 2, 3, 4}),
				makeLog(t, donID, uint32(44), "http://example.com/xxx3.json", [32]byte{3, 2, 3, 4}),
				makeLog(t, donID, uint32(43), "http://example.com/xxx4.json", [32]byte{4, 2, 3, 4}),
			}
			// Set empty adder logs
			lp.adderLogs = []logpoller.Log{}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)

			// Check that fetch triggers were sent for all owner logs
			triggers := collectTriggers(fetchTriggerCh, 8)
			require.Len(t, triggers, 4, "expected 4 triggers")
			// Find the trigger with version 45 (latest)
			var latestTrigger *types.Trigger
			for i := range triggers {
				if triggers[i].Version == 45 {
					latestTrigger = &triggers[i]
					break
				}
			}
			require.NotNil(t, latestTrigger, "expected trigger with version 45")
			require.Equal(t, "http://example.com/xxx2.json", latestTrigger.URL)
			require.Equal(t, [32]byte{2, 2, 3, 4}, latestTrigger.SHA)
		})
		t.Run("ignores logs with incorrect don ID", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			lp.filteredLogsErr = nil
			// Set owner logs with wrong donID
			lp.ownerLogs = []logpoller.Log{
				makeLog(t, donID+1, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
			}
			// Set empty adder logs
			lp.adderLogs = []logpoller.Log{}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)

			// Check that no fetch trigger was sent
			select {
			case trigger := <-fetchTriggerCh:
				t.Fatalf("did not expect fetch trigger signal for log with wrong donID, got: %+v", trigger)
			default:
				// No signal, as expected
			}
		})
		t.Run("ignores logs with wrong number of topics", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			lp.filteredLogsErr = nil
			lg := makeLog(t, donID, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4})
			lg.Topics = lg.Topics[:1]
			// Set owner log with wrong number of topics
			lp.ownerLogs = []logpoller.Log{lg}
			// Set empty adder logs
			lp.adderLogs = []logpoller.Log{}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)

			// Check that no fetch trigger was sent
			select {
			case trigger := <-fetchTriggerCh:
				t.Fatalf("did not expect fetch trigger signal for log with missing topics, got: %+v", trigger)
			default:
				// No signal, as expected
			}
		})
		t.Run("reads adder logs and sends triggers", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			lp.filteredLogsErr = nil
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			cdc.definitions.LastBlockNum = 0
			cdc.initialBlockNum = 0
			adderID1 := uint32(100)
			adderID2 := uint32(200)
			// Set adder logs (readLogs calls FilteredLogs for adder logs first)
			lp.adderLogs = []logpoller.Log{
				makeAdderLog(t, donID, adderID1, "http://example.com/adder1.json", [32]byte{1, 1, 1, 1}, 1500),
				makeAdderLog(t, donID, adderID2, "http://example.com/adder2.json", [32]byte{2, 2, 2, 2}, 1600),
			}
			// Set empty owner logs
			lp.ownerLogs = []logpoller.Log{}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)

			// Check that fetch triggers were sent for both adders
			triggers := collectTriggers(fetchTriggerCh, 2)
			require.Len(t, triggers, 2, "expected 2 triggers")
			// Verify adder triggers
			for _, trigger := range triggers {
				require.NotEqual(t, SourceOwner, trigger.Source, "should not be owner")
				require.True(t, trigger.Source == adderID1 || trigger.Source == adderID2, "should be one of the adder IDs")
				if trigger.Source == adderID1 {
					require.Equal(t, "http://example.com/adder1.json", trigger.URL)
					require.Equal(t, [32]byte{1, 1, 1, 1}, trigger.SHA)
				} else {
					require.Equal(t, "http://example.com/adder2.json", trigger.URL)
					require.Equal(t, [32]byte{2, 2, 2, 2}, trigger.SHA)
				}
			}
		})
		t.Run("reads both owner and adder logs in one call", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			lp.filteredLogsErr = nil
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			cdc.definitions.LastBlockNum = 0
			cdc.initialBlockNum = 0
			adderID := uint32(100)
			// Set both adder and owner logs
			lp.adderLogs = []logpoller.Log{
				makeAdderLog(t, donID, adderID, "http://example.com/adder.json", [32]byte{6, 6, 6, 6}, 1500),
			}
			lp.ownerLogs = []logpoller.Log{
				makeLog(t, donID, uint32(50), "http://example.com/owner.json", [32]byte{5, 5, 5, 5}),
			}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)

			// Should have triggers for both adder and owner logs
			triggers := collectTriggers(fetchTriggerCh, 2)
			require.Len(t, triggers, 2, "expected 2 triggers (one adder, one owner)")
			// Verify we have both types
			var foundOwner, foundAdder bool
			for _, trigger := range triggers {
				switch trigger.Source {
				case SourceOwner:
					foundOwner = true
					require.Equal(t, uint32(50), trigger.Version)
					require.Equal(t, "http://example.com/owner.json", trigger.URL)
				case adderID:
					foundAdder = true
					require.Equal(t, "http://example.com/adder.json", trigger.URL)
					require.Equal(t, [32]byte{6, 6, 6, 6}, trigger.SHA)
				}
			}
			require.True(t, foundOwner, "expected owner trigger")
			require.True(t, foundAdder, "expected adder trigger")
		})
		t.Run("ignores adder logs with incorrect don ID", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			lp.filteredLogsErr = nil
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			cdc.definitions.LastBlockNum = 0
			cdc.initialBlockNum = 0
			adderID := uint32(100)
			// Set adder logs with wrong donID
			lp.adderLogs = []logpoller.Log{
				makeAdderLog(t, donID+1, adderID, "http://example.com/adder.json", [32]byte{1, 1, 1, 1}, 1500),
			}
			// Set empty owner logs
			lp.ownerLogs = []logpoller.Log{}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)
			// Should not send trigger for wrong donID
			select {
			case trigger := <-fetchTriggerCh:
				t.Fatalf("did not expect fetch trigger signal for log with wrong donID, got: %+v", trigger)
			default:
				// No signal, as expected
			}
		})
	})

	t.Run("fetchChannelDefinitions", func(t *testing.T) {
		c := &mockHTTPClient{}
		cdc := &channelDefinitionCache{
			lggr:      logger.TestSugared(t),
			client:    c,
			httpLimit: 2048,
		}

		t.Run("invalid URL returns error", func(t *testing.T) {
			ctx := t.Context()
			// Set up mock to return error for invalid URL scheme
			c.err = errors.New("unsupported protocol scheme")
			c.resp = nil

			// Use a URL with invalid scheme that will fail at HTTP client level
			// This avoids panic from URL parsing in the HTTP library
			trigger := types.Trigger{
				Source:   SourceOwner,
				URL:      "http://[::1",
				SHA:      [32]byte{},
				BlockNum: 0,
				Version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			// The error could be from URL parsing or HTTP client - both are acceptable
			require.Error(t, err)
		})

		t.Run("networking error while making request returns error", func(t *testing.T) {
			ctx := t.Context()
			c.resp = nil
			c.err = errors.New("http request failed")

			trigger := types.Trigger{
				Source:   SourceOwner,
				URL:      "http://example.com/definitions.json",
				SHA:      [32]byte{},
				BlockNum: 0,
				Version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			require.Contains(t, err.Error(), "failed to make HTTP request to channel definitions URL")
			require.Contains(t, err.Error(), "http request failed")
		})

		t.Run("server returns 500 returns error", func(t *testing.T) {
			ctx := t.Context()
			c.err = nil
			c.resp = &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte{1, 2, 3}))}

			trigger := types.Trigger{
				Source:   SourceOwner,
				URL:      "http://example.com/definitions.json",
				SHA:      [32]byte{},
				BlockNum: 0,
				Version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			require.Contains(t, err.Error(), "HTTP error from channel definitions URL http://example.com/definitions.json (status 500)")
			require.Contains(t, err.Error(), "\x01\x02\x03")
		})

		var largeBody = make([]byte, 2048)
		for i := range largeBody {
			largeBody[i] = 'a'
		}

		t.Run("server returns 404 returns error (and does not log entirety of huge response body)", func(t *testing.T) {
			ctx := t.Context()
			c.err = nil
			c.resp = &http.Response{StatusCode: 404, Body: io.NopCloser(bytes.NewReader(largeBody))}

			trigger := types.Trigger{
				Source:   SourceOwner,
				URL:      "http://example.com/definitions.json",
				SHA:      [32]byte{},
				BlockNum: 0,
				Version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			require.Contains(t, err.Error(), "HTTP error from channel definitions URL http://example.com/definitions.json (status 404)")
			require.Contains(t, err.Error(), "failed to read response body")
			require.Contains(t, err.Error(), "http: request body too large")
		})

		var hugeBody = make([]byte, 8096)
		c.resp = &http.Response{Body: io.NopCloser(bytes.NewReader(hugeBody))}

		t.Run("server returns body that is too large", func(t *testing.T) {
			ctx := t.Context()
			c.err = nil
			c.resp = &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(hugeBody))}

			trigger := types.Trigger{
				Source:   SourceOwner,
				URL:      "http://example.com/definitions.json",
				SHA:      [32]byte{},
				BlockNum: 0,
				Version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			require.Contains(t, err.Error(), "failed to read channel definitions response body from")
			require.Contains(t, err.Error(), "http: request body too large")
		})

		t.Run("server returns invalid JSON returns error", func(t *testing.T) {
			ctx := t.Context()
			c.err = nil
			c.resp = &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader([]byte{1, 2, 3}))}

			expectedSha := common.HexToHash("0xfd1780a6fc9ee0dab26ceb4b3941ab03e66ccd970d1db91612c66df4515b0a0a")
			trigger := types.Trigger{
				Source:   SourceOwner,
				URL:      "http://example.com/definitions.json",
				SHA:      [32]byte(expectedSha),
				BlockNum: 0,
				Version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			require.Contains(t, err.Error(), "failed to decode channel definitions JSON from")
			require.Contains(t, err.Error(), "invalid character '\\x01' looking for beginning of value")
		})

		t.Run("SHA mismatch returns error", func(t *testing.T) {
			ctx := t.Context()
			c.err = nil
			c.resp = &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader([]byte(`{"foo":"bar"}`)))}

			trigger := types.Trigger{
				Source:   SourceOwner,
				URL:      "http://example.com/definitions.json",
				SHA:      [32]byte{},
				BlockNum: 0,
				Version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			require.Contains(t, err.Error(), "SHA3 mismatch for channel definitions from")
			require.Contains(t, err.Error(), "expected 0000000000000000000000000000000000000000000000000000000000000000")
			require.Contains(t, err.Error(), "got 4d3304d0d87c27a031cbb6bdf95da79b7b4552c3d0bef2e5a94f50810121e1e0")
		})

		t.Run("valid JSON matching SHA returns channel definitions", func(t *testing.T) {
			ctx := t.Context()
			chainSelector := 4949039107694359620 // arbitrum mainnet
			feedID := [32]byte{00, 03, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
			expirationWindow := 3600
			multiplier := big.NewInt(1e18)
			baseUSDFee := 10
			valid := fmt.Sprintf(`
{
	"42": {
		"reportFormat": %d,
		"chainSelector": %d,
		"streams": [{"streamId": 52, "aggregator": %d}, {"streamId": 53, "aggregator": %d}, {"streamId": 55, "aggregator": %d}],
		"opts": {
			"feedId": "0x%x",
			"expirationWindow": %d,
			"multiplier": "%s",
			"baseUSDFee": "%d"
		}
	}
}`, llotypes.ReportFormatEVMPremiumLegacy, chainSelector, llotypes.AggregatorMedian, llotypes.AggregatorMedian, llotypes.AggregatorQuote, feedID, expirationWindow, multiplier.String(), baseUSDFee)

			c.err = nil
			c.resp = &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader([]byte(valid)))}

			expectedSha := common.HexToHash("0x367bbc75f7b6c9fc66a98ea99f837ea7ac4a3c2d6a9ee284de018bd02c41b52d")
			trigger := types.Trigger{
				Source:   SourceOwner,
				URL:      "http://example.com/definitions.json",
				SHA:      [32]byte(expectedSha),
				BlockNum: 0,
				Version:  0,
			}
			cd, err := cdc.fetchChannelDefinitions(ctx, trigger)
			require.NoError(t, err)
			expectedDef := llotypes.ChannelDefinition{
				ReportFormat: 0x1,
				Streams:      []llotypes.Stream{{StreamID: 0x34, Aggregator: 0x1}, {StreamID: 0x35, Aggregator: 0x1}, {StreamID: 0x37, Aggregator: 0x3}},
				Opts:         llotypes.ChannelOpts{0x7b, 0x22, 0x62, 0x61, 0x73, 0x65, 0x55, 0x53, 0x44, 0x46, 0x65, 0x65, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x22, 0x2c, 0x22, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x22, 0x3a, 0x33, 0x36, 0x30, 0x30, 0x2c, 0x22, 0x66, 0x65, 0x65, 0x64, 0x49, 0x64, 0x22, 0x3a, 0x22, 0x30, 0x78, 0x30, 0x30, 0x30, 0x33, 0x36, 0x62, 0x34, 0x61, 0x61, 0x37, 0x65, 0x35, 0x37, 0x63, 0x61, 0x37, 0x62, 0x36, 0x38, 0x61, 0x65, 0x31, 0x62, 0x66, 0x34, 0x35, 0x36, 0x35, 0x33, 0x66, 0x35, 0x36, 0x62, 0x36, 0x35, 0x36, 0x66, 0x64, 0x33, 0x61, 0x61, 0x33, 0x33, 0x35, 0x65, 0x66, 0x37, 0x66, 0x61, 0x65, 0x36, 0x39, 0x36, 0x62, 0x36, 0x36, 0x33, 0x66, 0x31, 0x62, 0x38, 0x34, 0x37, 0x32, 0x22, 0x2c, 0x22, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x22, 0x7d},
				Source:       SourceOwner,
			}
			require.Equal(t, llotypes.ChannelDefinitions{0x2a: expectedDef}, cd)
		})
	})

	t.Run("persist", func(t *testing.T) {
		definitions := llotypes.ChannelDefinitions{
			1: {
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 2, Aggregator: llotypes.AggregatorMode}, {StreamID: 3, Aggregator: llotypes.AggregatorQuote}},
				Opts:         llotypes.ChannelOpts(`{"foo":"bar"}`),
			},
		}
		cdc := &channelDefinitionCache{
			lggr:  logger.TestSugared(t),
			orm:   nil,
			addr:  testutils.NewAddress(),
			donID: donID,
			definitions: Definitions{
				LastBlockNum: 142,
			},
		}

		t.Run("persists current definitions", func(t *testing.T) {
			ctx := t.Context()
			orm := &mockCDCORM{}
			cdc.orm = orm
			cdc.definitions.Version = 42
			cdc.persistedBlockNum = 141
			cdc.definitions.LastBlockNum = 142
			cdc.definitions.Sources = map[uint32]types.SourceDefinition{
				SourceOwner: {
					Trigger: types.Trigger{
						Source:   SourceOwner,
						BlockNum: 142,
						Version:  42,
					},
					Definitions: definitions,
				},
			}

			// persist() always persists c.definitions (no comparison logic)
			memoryBlockNum, persistedBlockNum, err := cdc.persist(ctx)
			require.NoError(t, err)
			require.Equal(t, int64(142), memoryBlockNum)
			require.Equal(t, int64(142), persistedBlockNum)
			require.Equal(t, int64(142), cdc.persistedBlockNum)
			require.Equal(t, cdc.definitions.Sources, orm.lastPersistedDfns)
		})

		orm := &mockCDCORM{}
		cdc.orm = orm

		t.Run("returns error on db failure and does not update persisted block number", func(t *testing.T) {
			ctx := t.Context()
			cdc.persistedBlockNum = 141
			cdc.definitions.Version = 43
			cdc.definitions.LastBlockNum = 143
			cdc.definitions.Sources = map[uint32]types.SourceDefinition{
				SourceOwner: {
					Trigger: types.Trigger{
						Source:   SourceOwner,
						BlockNum: 143,
						Version:  43,
					},
					Definitions: definitions,
				},
			}
			orm.err = errors.New("test error")

			// persist() always persists c.definitions
			memoryBlockNum, persistedBlockNum, err := cdc.persist(ctx)
			require.Contains(t, err.Error(), "test error")
			require.Equal(t, int64(143), memoryBlockNum)
			require.Equal(t, int64(141), persistedBlockNum)
			require.Equal(t, int64(141), cdc.persistedBlockNum)
		})

		t.Run("updates persisted block number on success", func(t *testing.T) {
			ctx := t.Context()
			cdc.definitions.Version = 43
			cdc.definitions.LastBlockNum = 143
			cdc.definitions.Sources = map[uint32]types.SourceDefinition{
				SourceOwner: {
					Trigger: types.Trigger{
						Source:   SourceOwner,
						BlockNum: 143,
						Version:  43,
					},
					Definitions: definitions,
				},
			}
			cdc.persistedBlockNum = 141
			orm.err = nil

			// persist() always persists c.definitions
			memoryBlockNum, persistedBlockNum, err := cdc.persist(ctx)
			require.NoError(t, err)
			require.Equal(t, int64(143), memoryBlockNum)
			require.Equal(t, int64(143), persistedBlockNum)
			require.Equal(t, int64(143), cdc.persistedBlockNum)

			require.Equal(t, cdc.addr, orm.lastPersistedAddr)
			require.Equal(t, cdc.donID, orm.lastPersistedDonID)
			require.Equal(t, cdc.definitions.Version, orm.lastPersistedVersion)
			require.Equal(t, cdc.definitions.Sources, orm.lastPersistedDfns)
			require.Equal(t, cdc.definitions.LastBlockNum, orm.lastPersistedBlockNum)
		})
	})

	t.Run("adder limits", func(t *testing.T) {
		cdc := &channelDefinitionCache{
			lggr: logger.TestSugared(t),
		}

		adderID := uint32(100)

		t.Run("rejects adder definition file with more than MaxChannelsPerAdder channels", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Create a new definition file with MaxChannelsPerAdder + 1 channels
			// The limit is enforced based on existing channels plus new channels being added
			// When trying to add the (MaxChannelsPerAdder+1)th channel, numberOfChannels will be MaxChannelsPerAdder
			addChannelDefinitions(newDefinitions, 1, uint32(MaxChannelsPerAdder+1), adderID)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(adderID, currentDefinitions, newDefinitions, feedIDToChannelID)
			// The implementation stops processing at MaxChannelsPerAdder and doesn't return an error
			// Verify that only MaxChannelsPerAdder channels were added
			require.LessOrEqual(t, len(currentDefinitions), MaxChannelsPerAdder, "should not exceed MaxChannelsPerAdder")
			// Count channels from this adder source
			adderChannelCount := uint32(0)
			for _, def := range currentDefinitions {
				if def.Source == adderID {
					adderChannelCount++
				}
			}
			require.Equal(t, MaxChannelsPerAdder, int(adderChannelCount), "should have exactly MaxChannelsPerAdder channels from this adder")
		})

		t.Run("allows adder definition file with channels up to MaxChannelsPerAdder when most are existing", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Pre-populate with 90 existing channels (MaxChannelsPerAdder - 10)
			// This tests that existing channels + new channels can total up to MaxChannelsPerAdder
			existingEnd := uint32(90)
			addChannelDefinitions(currentDefinitions, 1, existingEnd, adderID)
			// Include these existing channels in the new definition file (they'll be skipped)
			addChannelDefinitions(newDefinitions, 1, existingEnd, adderID)

			// Add 9 new channels to reach exactly MaxChannelsPerAdder total in the file
			addChannelDefinitions(newDefinitions, existingEnd+1, existingEnd+9, adderID)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(adderID, currentDefinitions, newDefinitions, feedIDToChannelID)
			// Should have 90 existing + 9 new = 99 (below MaxChannelsPerAdder)
			require.Len(t, currentDefinitions, 99)
		})

		t.Run("owner definitions are not subject to adder limits", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Owner can add any number of channels (not subject to MaxChannelsPerAdder limit)
			addChannelDefinitions(newDefinitions, 1, 20, SourceOwner)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)
			require.Len(t, currentDefinitions, 20)
		})

		t.Run("owner can have more than MaxChannelsPerAdder channels", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Owner can have more than MaxChannelsPerAdder channels
			addChannelDefinitions(newDefinitions, 1, uint32(MaxChannelsPerAdder+10), SourceOwner)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)
			require.Len(t, currentDefinitions, MaxChannelsPerAdder+10)
		})
	})

	t.Run("owner removal", func(t *testing.T) {
		cdc := &channelDefinitionCache{
			lggr: logger.TestSugared(t),
		}

		t.Run("does not remove owner-defined channels missing from new definitions", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Set up current definitions with owner-defined channels 1, 2, 3, 4, 5
			addChannelDefinitions(currentDefinitions, 1, 5, SourceOwner)

			// New definitions only include channels 1, 3, 5 (missing 2 and 4)
			newDefinitions[1] = makeChannelDefinition(1, SourceOwner)
			newDefinitions[3] = makeChannelDefinition(3, SourceOwner)
			newDefinitions[5] = makeChannelDefinition(5, SourceOwner)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)

			// Channels 1, 3, 5 should be present (updated from newDefinitions)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			require.Contains(t, currentDefinitions, llotypes.ChannelID(3))
			require.Contains(t, currentDefinitions, llotypes.ChannelID(5))

			// Channels 2 and 4 should remain (not removed, just not in newDefinitions)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(2))
			require.Contains(t, currentDefinitions, llotypes.ChannelID(4))

			// Result should contain all 5 channels (2 and 4 remain from currentDefinitions)
			require.Len(t, currentDefinitions, 5)
		})

		t.Run("preserves non-owner channels when owner updates definitions", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			adderID := uint32(100)

			// Set up current definitions with owner channels 1, 2 and adder channel 10
			addChannelDefinitions(currentDefinitions, 1, 2, SourceOwner)
			currentDefinitions[10] = makeChannelDefinition(10, adderID)

			// New owner definitions only include channel 1 (missing channel 2)
			newDefinitions[1] = makeChannelDefinition(1, SourceOwner)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)

			// Owner channel 1 should be present (updated from newDefinitions)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			require.Equal(t, SourceOwner, currentDefinitions[1].Source)

			// Owner channel 2 should remain (not removed, just not in newDefinitions)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(2))
			require.Equal(t, SourceOwner, currentDefinitions[2].Source)

			// Adder channel 10 should be preserved
			require.Contains(t, currentDefinitions, llotypes.ChannelID(10))
			require.Equal(t, adderID, currentDefinitions[10].Source)

			// Result should contain channel 1 (owner, updated), channel 2 (owner, preserved), and channel 10 (adder)
			require.Len(t, currentDefinitions, 3)
		})

		t.Run("owner removal only happens when source is SourceOwner", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			adderID := uint32(200)

			// Set up current definitions with owner-defined channels 1, 2
			addChannelDefinitions(currentDefinitions, 1, 2, SourceOwner)

			// New definitions from adder only includes channel 1
			newDefinitions[1] = makeChannelDefinition(1, adderID)

			// When source is an adder (not SourceOwner), owner channels should NOT be removed
			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(adderID, currentDefinitions, newDefinitions, feedIDToChannelID)

			// Owner channel 1 should still be present (adder can't overwrite it)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			require.Equal(t, SourceOwner, currentDefinitions[1].Source, "channel 1 should still have owner source")

			// Owner channel 2 should still be present (not removed because source is not SourceOwner)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(2))
			require.Equal(t, SourceOwner, currentDefinitions[2].Source)

			// Result should contain both owner channels (adder's attempt to add channel 1 is ignored)
			require.Len(t, currentDefinitions, 2)
		})

		t.Run("owner can tombstone owner channels", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			adderID := uint32(300)

			// Set up current definitions with owner channel 1 and adder channel 2
			currentDefinitions[1] = makeChannelDefinition(1, SourceOwner)
			currentDefinitions[2] = makeChannelDefinition(2, adderID)

			// Owner tries to tombstone owner channel 1 (should succeed)
			newDefinitions[1] = llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatJSON,
				Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}},
				Source:       SourceOwner,
				Tombstone:    true,
			}

			// Owner tries to tombstone adder channel 2 (should succeed)
			newDefinitions[2] = llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatJSON,
				Source:       SourceOwner,
				Tombstone:    true,
			}

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)

			// Result should contain both channels
			require.Len(t, currentDefinitions, 2)

			// Owner channel 1 should be present (tombstone succeeded)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			require.Equal(t, SourceOwner, currentDefinitions[1].Source)
			require.True(t, currentDefinitions[1].Tombstone, "channel 1 should be tombstoned")

			// Adder channel 2 should be kept in definitions with Tombstone: true (tombstone succeeded)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(2))
			require.True(t, currentDefinitions[2].Tombstone, "channel 2 should be tombstoned")
		})
	})

	t.Run("feedID uniqueness", func(t *testing.T) {
		cdc := &channelDefinitionCache{
			lggr: logger.TestSugared(t),
		}

		adderID := uint32(100)
		feedID1 := common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
		feedID2 := common.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222")

		t.Run("skips new channel with colliding FeedID", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Existing channel 1 with feedID1
			currentDefinitions[1] = makeChannelDefinitionWithFeedID(1, SourceOwner, feedID1)

			// New channel 2 with same feedID1 (collision)
			newDefinitions[2] = makeChannelDefinitionWithFeedID(2, SourceOwner, feedID1)

			// New channel 3 with unique feedID2 (should be added)
			newDefinitions[3] = makeChannelDefinitionWithFeedID(3, SourceOwner, feedID2)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)

			// Channel 1 should be present (existing)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			// Channel 2 should NOT be present (collision, skipped)
			require.NotContains(t, currentDefinitions, llotypes.ChannelID(2))
			// Channel 3 should be present (unique FeedID)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(3))
			require.Len(t, currentDefinitions, 2)
		})

		t.Run("allows owner to update same channel with same FeedID", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Existing channel 1 with feedID1
			currentDefinitions[1] = makeChannelDefinitionWithFeedID(1, SourceOwner, feedID1)

			// Owner updates channel 1 with same feedID1 (should be allowed)
			updatedDef := makeChannelDefinitionWithFeedID(1, SourceOwner, feedID1)
			updatedDef.Streams = []llotypes.Stream{{StreamID: 999, Aggregator: llotypes.AggregatorMedian}}
			newDefinitions[1] = updatedDef

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)

			// Channel 1 should be present and updated
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			require.Equal(t, uint32(999), currentDefinitions[1].Streams[0].StreamID)
			require.Len(t, currentDefinitions, 1)
		})

		t.Run("skips owner update to same channel with colliding FeedID", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Existing channel 1 with feedID1
			currentDefinitions[1] = makeChannelDefinitionWithFeedID(1, SourceOwner, feedID1)
			// Existing channel 2 with feedID2
			currentDefinitions[2] = makeChannelDefinitionWithFeedID(2, SourceOwner, feedID2)

			// Owner tries to update channel 1 with feedID2 (collides with channel 2)
			newDefinitions[1] = makeChannelDefinitionWithFeedID(1, SourceOwner, feedID2)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)

			// Channel 1 should still have feedID1 (update was skipped)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			require.Equal(t, feedID1, extractFeedID(currentDefinitions[1].Opts))
			// Channel 2 should still be present
			require.Contains(t, currentDefinitions, llotypes.ChannelID(2))
			require.Len(t, currentDefinitions, 2)
		})

		t.Run("skips adder channel with colliding FeedID", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Existing owner channel 1 with feedID1
			currentDefinitions[1] = makeChannelDefinitionWithFeedID(1, SourceOwner, feedID1)

			// Adder tries to add channel 2 with same feedID1 (collision)
			newDefinitions[2] = makeChannelDefinitionWithFeedID(2, adderID, feedID1)

			// Adder tries to add channel 3 with unique feedID2 (should be added)
			newDefinitions[3] = makeChannelDefinitionWithFeedID(3, adderID, feedID2)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(adderID, currentDefinitions, newDefinitions, feedIDToChannelID)

			// Channel 1 should be present (existing)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			// Channel 2 should NOT be present (collision, skipped)
			require.NotContains(t, currentDefinitions, llotypes.ChannelID(2))
			// Channel 3 should be present (unique FeedID)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(3))
			require.Len(t, currentDefinitions, 2)
		})

		t.Run("allows owner to update channel with new unique FeedID", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Existing channel 1 with feedID1
			currentDefinitions[1] = makeChannelDefinitionWithFeedID(1, SourceOwner, feedID1)

			// Owner updates channel 1 with new unique feedID2 (should be allowed)
			newDefinitions[1] = makeChannelDefinitionWithFeedID(1, SourceOwner, feedID2)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)

			// Channel 1 should be present with new feedID2
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			require.Equal(t, feedID2, extractFeedID(currentDefinitions[1].Opts))
			require.Len(t, currentDefinitions, 1)
		})

		t.Run("skips multiple channels with same colliding FeedID", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Existing channel 1 with feedID1
			currentDefinitions[1] = makeChannelDefinitionWithFeedID(1, SourceOwner, feedID1)

			// Multiple new channels with same colliding feedID1
			newDefinitions[2] = makeChannelDefinitionWithFeedID(2, SourceOwner, feedID1)
			newDefinitions[3] = makeChannelDefinitionWithFeedID(3, SourceOwner, feedID1)
			newDefinitions[4] = makeChannelDefinitionWithFeedID(4, SourceOwner, feedID1)

			// One channel with unique feedID2
			newDefinitions[5] = makeChannelDefinitionWithFeedID(5, SourceOwner, feedID2)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)

			// Channel 1 should be present (existing)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			// Channels 2, 3, 4 should NOT be present (all collided)
			require.NotContains(t, currentDefinitions, llotypes.ChannelID(2))
			require.NotContains(t, currentDefinitions, llotypes.ChannelID(3))
			require.NotContains(t, currentDefinitions, llotypes.ChannelID(4))
			// Channel 5 should be present (unique FeedID)
			require.Contains(t, currentDefinitions, llotypes.ChannelID(5))
			require.Len(t, currentDefinitions, 2)
		})

		t.Run("allows channels without FeedID", func(t *testing.T) {
			currentDefinitions := make(llotypes.ChannelDefinitions)
			newDefinitions := make(llotypes.ChannelDefinitions)

			// Existing channel 1 with feedID1
			currentDefinitions[1] = makeChannelDefinitionWithFeedID(1, SourceOwner, feedID1)

			// New channel 2 without FeedID (should be allowed)
			newDefinitions[2] = makeChannelDefinition(2, SourceOwner)

			// New channel 3 with unique feedID2 (should be allowed)
			newDefinitions[3] = makeChannelDefinitionWithFeedID(3, SourceOwner, feedID2)

			feedIDToChannelID := buildFeedIDMap(currentDefinitions)
			cdc.mergeDefinitions(SourceOwner, currentDefinitions, newDefinitions, feedIDToChannelID)

			// All channels should be present
			require.Contains(t, currentDefinitions, llotypes.ChannelID(1))
			require.Contains(t, currentDefinitions, llotypes.ChannelID(2))
			require.Contains(t, currentDefinitions, llotypes.ChannelID(3))
			require.Len(t, currentDefinitions, 3)
		})
	})
}

func Test_filterName(t *testing.T) {
	s := types.ChannelDefinitionCacheFilterName(common.Address{1, 2, 3}, 654)
	require.Equal(t, "OCR3 LLO ChannelDefinitionCachePoller - 0x0102030000000000000000000000000000000000:654", s)
}
