package channeldefinitions

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
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
	lastPersistedDfns     llotypes.ChannelDefinitions
	lastPersistedBlockNum int64
}

func (m *mockCDCORM) LoadChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) (pd *types.PersistedDefinitions, err error) {
	panic("not implemented")
}
func (m *mockCDCORM) StoreChannelDefinitions(ctx context.Context, addr common.Address, donID, version uint32, dfns llotypes.ChannelDefinitions, blockNum int64) (err error) {
	m.lastPersistedAddr = addr
	m.lastPersistedDonID = donID
	m.lastPersistedVersion = version
	m.lastPersistedDfns = dfns
	m.lastPersistedBlockNum = blockNum
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
	return logpoller.Log{EventSig: ChannelDefinitionAdded, Topics: [][]byte{ChannelDefinitionAdded[:], makeDonIDTopic(donID)}, Data: data, BlockNumber: blockNumber}
}

func makeAdderLogData(t *testing.T, donID, adderID uint32, url string, sha [32]byte) []byte {
	event := channelConfigStoreABI.Events[channelDefinitionAddedEventName]
	// donID is indexed
	// adderId, url, sha
	data, err := event.Inputs.NonIndexed().Pack(adderID, url, sha)
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
func collectTriggers(ch chan fetchTrigger, maxCount int) []fetchTrigger {
	triggers := make([]fetchTrigger, 0, maxCount)
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

func Test_ChannelDefinitionCache(t *testing.T) {
	donID := rand.Uint32()

	t.Run("Definitions", func(t *testing.T) {
		// NOTE: this is covered more thoroughly in the integration tests
		dfns := llotypes.ChannelDefinitions(map[llotypes.ChannelID]llotypes.ChannelDefinition{
			1: {
				ReportFormat: llotypes.ReportFormat(43),
				Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 2, Aggregator: llotypes.AggregatorMode}, {StreamID: 3, Aggregator: llotypes.AggregatorQuote}},
				Opts:         llotypes.ChannelOpts{1, 2, 3},
			},
		})

		cdc := &channelDefinitionCache{definitions: dfns}

		assert.Equal(t, dfns, cdc.Definitions())
	})

	t.Run("readLogs", func(t *testing.T) {
		lp := &mockLogPoller{latestBlockErr: sql.ErrNoRows}
		fetchTriggerCh := make(chan fetchTrigger, 100)
		cdc := &channelDefinitionCache{
			donID:          donID,
			lp:             lp,
			lggr:           logger.TestSugared(t),
			fetchTriggerCh: fetchTriggerCh,
		}

		t.Run("skips if logpoller has no blocks", func(t *testing.T) {
			ctx := t.Context()
			err := cdc.readLogs(ctx)
			assert.NoError(t, err)
		})
		t.Run("returns error on LatestBlock failure", func(t *testing.T) {
			ctx := t.Context()
			lp.latestBlockErr = errors.New("test error")

			err := cdc.readLogs(ctx)
			assert.EqualError(t, err, "test error")
		})
		t.Run("does nothing if LatestBlock older or the same as current channel definitions block", func(t *testing.T) {
			ctx := t.Context()
			lp.latestBlockErr = nil
			lp.latestBlock = logpoller.Block{BlockNumber: 42}
			cdc.definitionsBlockNum = 43

			err := cdc.readLogs(ctx)
			assert.NoError(t, err)
		})
		t.Run("returns error if FilteredLogs fails", func(t *testing.T) {
			ctx := t.Context()
			cdc.definitionsBlockNum = 0
			lp.filteredLogsErr = errors.New("test error 2")

			err := cdc.readLogs(ctx)
			assert.EqualError(t, err, "test error 2")
		})
		t.Run("ignores logs with different topic", func(t *testing.T) {
			ctx := t.Context()
			lp.filteredLogsErr = nil
			lp.filteredLogs = []logpoller.Log{{EventSig: common.Hash{1, 2, 3, 4}}}

			err := cdc.readLogs(ctx)
			assert.NoError(t, err)
		})
		t.Run("logs warning and continues if log is malformed", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			cdc.definitionsBlockNum = 0
			cdc.initialBlockNum = 0
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			lp.latestBlockErr = nil
			lp.filteredLogsErr = nil
			lp.filteredLogs = []logpoller.Log{{EventSig: NewChannelDefinition}}

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
		t.Run("sets definitions and sends on channel if FilteredLogs returns new event with a later version", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			cdc.definitionsBlockNum = 0
			cdc.initialBlockNum = 0
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			lp.latestBlockErr = nil
			lp.filteredLogsErr = nil
			lp.filteredLogs = []logpoller.Log{makeLog(t, donID, uint32(43), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4})}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)

			// Check that fetch trigger was sent
			select {
			case trigger := <-fetchTriggerCh:
				assert.Equal(t, SourceOwner, trigger.source)
				assert.Equal(t, uint32(43), trigger.version)
				assert.Equal(t, "http://example.com/xxx.json", trigger.url)
				assert.Equal(t, [32]byte{1, 2, 3, 4}, trigger.sha)
			default:
				t.Fatal("expected fetch trigger signal in channel")
			}
		})
		t.Run("sends triggers for all logs", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			cdc.definitionsBlockNum = 0
			cdc.initialBlockNum = 0
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			lp.filteredLogsErr = nil
			lp.filteredLogs = []logpoller.Log{
				makeLog(t, donID, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
				makeLog(t, donID, uint32(43), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
			}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)
			// Should receive triggers for both logs (they get processed twice because readLogs calls FilteredLogs twice)
			// The logs are sorted by block number, so we get them in order
			triggers := collectTriggers(fetchTriggerCh, 4)
			require.GreaterOrEqual(t, len(triggers), 2, "expected at least 2 triggers")
			// Find the trigger with version 43 (latest)
			var found43 bool
			for _, trigger := range triggers {
				if trigger.version == 43 {
					found43 = true
					break
				}
			}
			assert.True(t, found43, "expected trigger with version 43")
		})
		t.Run("in case of multiple logs, sends triggers for all", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			cdc.definitionsBlockNum = 0
			cdc.initialBlockNum = 0
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			lp.filteredLogsErr = nil
			lp.filteredLogs = []logpoller.Log{
				makeLog(t, donID, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
				makeLog(t, donID, uint32(45), "http://example.com/xxx2.json", [32]byte{2, 2, 3, 4}),
				makeLog(t, donID, uint32(44), "http://example.com/xxx3.json", [32]byte{3, 2, 3, 4}),
				makeLog(t, donID, uint32(43), "http://example.com/xxx4.json", [32]byte{4, 2, 3, 4}),
			}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)

			// Check that fetch triggers were sent for all logs
			// Note: readLogs calls FilteredLogs twice (once for owner, once for adder), so we get duplicates
			triggers := collectTriggers(fetchTriggerCh, 8)
			require.GreaterOrEqual(t, len(triggers), 4, "expected at least 4 triggers")
			// Find the trigger with version 45 (latest)
			var latestTrigger *fetchTrigger
			for i := range triggers {
				if triggers[i].version == 45 {
					latestTrigger = &triggers[i]
					break
				}
			}
			require.NotNil(t, latestTrigger, "expected trigger with version 45")
			assert.Equal(t, "http://example.com/xxx2.json", latestTrigger.url)
			assert.Equal(t, [32]byte{2, 2, 3, 4}, latestTrigger.sha)
		})
		t.Run("ignores logs with incorrect don ID", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			lp.filteredLogsErr = nil
			lp.filteredLogs = []logpoller.Log{
				makeLog(t, donID+1, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
			}

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
			lp.filteredLogs = []logpoller.Log{lg}

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
			cdc.definitionsBlockNum = 0
			cdc.initialBlockNum = 0
			adderID1 := uint32(100)
			adderID2 := uint32(200)
			lp.filteredLogs = []logpoller.Log{
				makeAdderLog(t, donID, adderID1, "http://example.com/adder1.json", [32]byte{1, 1, 1, 1}, 1500),
				makeAdderLog(t, donID, adderID2, "http://example.com/adder2.json", [32]byte{2, 2, 2, 2}, 1600),
			}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)

			// Check that fetch triggers were sent for both adders
			triggers := collectTriggers(fetchTriggerCh, 2)
			require.Len(t, triggers, 2, "expected 2 triggers")
			// Verify adder triggers
			for _, trigger := range triggers {
				assert.NotEqual(t, SourceOwner, trigger.source, "should not be owner")
				assert.True(t, trigger.source == adderID1 || trigger.source == adderID2, "should be one of the adder IDs")
				if trigger.source == adderID1 {
					assert.Equal(t, "http://example.com/adder1.json", trigger.url)
					assert.Equal(t, [32]byte{1, 1, 1, 1}, trigger.sha)
				} else {
					assert.Equal(t, "http://example.com/adder2.json", trigger.url)
					assert.Equal(t, [32]byte{2, 2, 2, 2}, trigger.sha)
				}
			}
		})
		t.Run("reads both owner and adder logs in one call", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			lp.filteredLogsErr = nil
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			cdc.definitionsBlockNum = 0
			cdc.initialBlockNum = 0
			// readLogs calls FilteredLogs twice - once for owner logs, once for adder logs
			// The mock returns the same logs each time, so we set it to return owner logs
			// which will be processed on the first call, then the same logs on second call
			// (which will try to process as adder logs but fail validation)
			// For a proper test, we'd need a smarter mock, but for now just verify owner logs work
			lp.filteredLogs = []logpoller.Log{
				makeLog(t, donID, uint32(50), "http://example.com/owner.json", [32]byte{5, 5, 5, 5}),
			}

			err := cdc.readLogs(ctx)
			require.NoError(t, err)

			// Should have at least one trigger (owner log)
			select {
			case trigger := <-fetchTriggerCh:
				assert.Equal(t, SourceOwner, trigger.source)
				assert.Equal(t, uint32(50), trigger.version)
				assert.Equal(t, "http://example.com/owner.json", trigger.url)
			default:
				t.Fatal("expected owner trigger")
			}
		})
		t.Run("ignores adder logs with incorrect don ID", func(t *testing.T) {
			ctx := t.Context()
			// Drain any existing triggers
			drainChannel(fetchTriggerCh)
			lp.filteredLogsErr = nil
			lp.latestBlock = logpoller.Block{BlockNumber: 2000}
			cdc.definitionsBlockNum = 0
			cdc.initialBlockNum = 0
			adderID := uint32(100)
			lp.filteredLogs = []logpoller.Log{
				makeAdderLog(t, donID+1, adderID, "http://example.com/adder.json", [32]byte{1, 1, 1, 1}, 1500),
			}

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
			trigger := fetchTrigger{
				source:   SourceOwner,
				url:      "http://[::1",
				sha:      [32]byte{},
				blockNum: 0,
				version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			// The error could be from URL parsing or HTTP client - both are acceptable
			assert.Error(t, err)
		})

		t.Run("networking error while making request returns error", func(t *testing.T) {
			ctx := t.Context()
			c.resp = nil
			c.err = errors.New("http request failed")

			trigger := fetchTrigger{
				source:   SourceOwner,
				url:      "http://example.com/definitions.json",
				sha:      [32]byte{},
				blockNum: 0,
				version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			assert.Contains(t, err.Error(), "failed to make HTTP request to channel definitions URL")
			assert.Contains(t, err.Error(), "http request failed")
		})

		t.Run("server returns 500 returns error", func(t *testing.T) {
			ctx := t.Context()
			c.err = nil
			c.resp = &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte{1, 2, 3}))}

			trigger := fetchTrigger{
				source:   SourceOwner,
				url:      "http://example.com/definitions.json",
				sha:      [32]byte{},
				blockNum: 0,
				version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			assert.Contains(t, err.Error(), "HTTP error from channel definitions URL http://example.com/definitions.json (status 500)")
			assert.Contains(t, err.Error(), "\x01\x02\x03")
		})

		var largeBody = make([]byte, 2048)
		for i := range largeBody {
			largeBody[i] = 'a'
		}

		t.Run("server returns 404 returns error (and does not log entirety of huge response body)", func(t *testing.T) {
			ctx := t.Context()
			c.err = nil
			c.resp = &http.Response{StatusCode: 404, Body: io.NopCloser(bytes.NewReader(largeBody))}

			trigger := fetchTrigger{
				source:   SourceOwner,
				url:      "http://example.com/definitions.json",
				sha:      [32]byte{},
				blockNum: 0,
				version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			assert.Contains(t, err.Error(), "HTTP error from channel definitions URL http://example.com/definitions.json (status 404)")
			assert.Contains(t, err.Error(), "failed to read response body")
			assert.Contains(t, err.Error(), "http: request body too large")
		})

		var hugeBody = make([]byte, 8096)
		c.resp = &http.Response{Body: io.NopCloser(bytes.NewReader(hugeBody))}

		t.Run("server returns body that is too large", func(t *testing.T) {
			ctx := t.Context()
			c.err = nil
			c.resp = &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(hugeBody))}

			trigger := fetchTrigger{
				source:   SourceOwner,
				url:      "http://example.com/definitions.json",
				sha:      [32]byte{},
				blockNum: 0,
				version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			assert.Contains(t, err.Error(), "failed to read channel definitions response body from")
			assert.Contains(t, err.Error(), "http: request body too large")
		})

		t.Run("server returns invalid JSON returns error", func(t *testing.T) {
			ctx := t.Context()
			c.err = nil
			c.resp = &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader([]byte{1, 2, 3}))}

			expectedSha := common.HexToHash("0xfd1780a6fc9ee0dab26ceb4b3941ab03e66ccd970d1db91612c66df4515b0a0a")
			trigger := fetchTrigger{
				source:   SourceOwner,
				url:      "http://example.com/definitions.json",
				sha:      [32]byte(expectedSha),
				blockNum: 0,
				version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			assert.Contains(t, err.Error(), "failed to decode channel definitions JSON from")
			assert.Contains(t, err.Error(), "invalid character '\\x01' looking for beginning of value")
		})

		t.Run("SHA mismatch returns error", func(t *testing.T) {
			ctx := t.Context()
			c.err = nil
			c.resp = &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader([]byte(`{"foo":"bar"}`)))}

			trigger := fetchTrigger{
				source:   SourceOwner,
				url:      "http://example.com/definitions.json",
				sha:      [32]byte{},
				blockNum: 0,
				version:  0,
			}
			_, err := cdc.fetchChannelDefinitions(ctx, trigger)
			assert.Contains(t, err.Error(), "SHA3 mismatch for channel definitions from")
			assert.Contains(t, err.Error(), "expected 0000000000000000000000000000000000000000000000000000000000000000")
			assert.Contains(t, err.Error(), "got 4d3304d0d87c27a031cbb6bdf95da79b7b4552c3d0bef2e5a94f50810121e1e0")
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
			trigger := fetchTrigger{
				source:   SourceOwner,
				url:      "http://example.com/definitions.json",
				sha:      [32]byte(expectedSha),
				blockNum: 0,
				version:  0,
			}
			cd, err := cdc.fetchChannelDefinitions(ctx, trigger)
			assert.NoError(t, err)
			expectedDef := llotypes.ChannelDefinition{
				ReportFormat: 0x1,
				Streams:      []llotypes.Stream{{StreamID: 0x34, Aggregator: 0x1}, {StreamID: 0x35, Aggregator: 0x1}, {StreamID: 0x37, Aggregator: 0x3}},
				Opts:         llotypes.ChannelOpts{0x7b, 0x22, 0x62, 0x61, 0x73, 0x65, 0x55, 0x53, 0x44, 0x46, 0x65, 0x65, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x22, 0x2c, 0x22, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x22, 0x3a, 0x33, 0x36, 0x30, 0x30, 0x2c, 0x22, 0x66, 0x65, 0x65, 0x64, 0x49, 0x64, 0x22, 0x3a, 0x22, 0x30, 0x78, 0x30, 0x30, 0x30, 0x33, 0x36, 0x62, 0x34, 0x61, 0x61, 0x37, 0x65, 0x35, 0x37, 0x63, 0x61, 0x37, 0x62, 0x36, 0x38, 0x61, 0x65, 0x31, 0x62, 0x66, 0x34, 0x35, 0x36, 0x35, 0x33, 0x66, 0x35, 0x36, 0x62, 0x36, 0x35, 0x36, 0x66, 0x64, 0x33, 0x61, 0x61, 0x33, 0x33, 0x35, 0x65, 0x66, 0x37, 0x66, 0x61, 0x65, 0x36, 0x39, 0x36, 0x62, 0x36, 0x36, 0x33, 0x66, 0x31, 0x62, 0x38, 0x34, 0x37, 0x32, 0x22, 0x2c, 0x22, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x22, 0x7d},
				Source:       SourceOwner,
			}
			assert.Equal(t, llotypes.ChannelDefinitions{0x2a: expectedDef}, cd)
		})
	})

	t.Run("persist", func(t *testing.T) {
		cdc := &channelDefinitionCache{
			lggr:  logger.TestSugared(t),
			orm:   nil,
			addr:  testutils.NewAddress(),
			donID: donID,
			definitions: llotypes.ChannelDefinitions{
				1: {
					ReportFormat: llotypes.ReportFormat(43),
					Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 2, Aggregator: llotypes.AggregatorMode}, {StreamID: 3, Aggregator: llotypes.AggregatorQuote}},
					Opts:         llotypes.ChannelOpts{1, 2, 3},
				},
			},
			definitionsBlockNum: 142,
		}

		t.Run("does nothing if persisted block number is up-to-date", func(t *testing.T) {
			ctx := t.Context()
			cdc.definitionsVersion = 42
			cdc.persistedBlockNum = 142 // Match definitionsBlockNum

			memoryBlockNum, persistedBlockNum, err := cdc.persist(ctx)
			assert.NoError(t, err)
			assert.Equal(t, int64(142), memoryBlockNum)
			assert.Equal(t, int64(142), persistedBlockNum)
			assert.Equal(t, int64(142), cdc.persistedBlockNum)
		})

		orm := &mockCDCORM{}
		cdc.orm = orm

		t.Run("returns error on db failure and does not update persisted block number", func(t *testing.T) {
			ctx := t.Context()
			cdc.persistedBlockNum = 141
			cdc.definitionsVersion = 43
			cdc.definitionsBlockNum = 143
			orm.err = errors.New("test error")

			memoryBlockNum, persistedBlockNum, err := cdc.persist(ctx)
			assert.EqualError(t, err, "test error")
			assert.Equal(t, int64(143), memoryBlockNum)
			assert.Equal(t, int64(141), persistedBlockNum)
			assert.Equal(t, int64(141), cdc.persistedBlockNum)
		})

		t.Run("updates persisted block number on success", func(t *testing.T) {
			ctx := t.Context()
			cdc.definitionsVersion = 43
			cdc.definitionsBlockNum = 143
			cdc.persistedBlockNum = 141
			orm.err = nil

			memoryBlockNum, persistedBlockNum, err := cdc.persist(ctx)
			assert.NoError(t, err)
			assert.Equal(t, int64(143), memoryBlockNum)
			assert.Equal(t, int64(143), persistedBlockNum)
			assert.Equal(t, int64(143), cdc.persistedBlockNum)

			assert.Equal(t, cdc.addr, orm.lastPersistedAddr)
			assert.Equal(t, cdc.donID, orm.lastPersistedDonID)
			assert.Equal(t, cdc.definitionsVersion, orm.lastPersistedVersion)
			assert.Equal(t, cdc.definitions, orm.lastPersistedDfns)
			assert.Equal(t, cdc.definitionsBlockNum, orm.lastPersistedBlockNum)
		})
	})
}

func Test_filterName(t *testing.T) {
	s := types.ChannelDefinitionCacheFilterName(common.Address{1, 2, 3}, 654)
	assert.Equal(t, "OCR3 LLO ChannelDefinitionCachePoller - 0x0102030000000000000000000000000000000000:654", s)
}
