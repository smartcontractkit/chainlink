package llo

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
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

type mockLogPoller struct {
	latestBlock     logpoller.LogPollerBlock
	latestBlockErr  error
	logsWithSigs    []logpoller.Log
	logsWithSigsErr error

	unregisteredFilterNames []string
}

func (m *mockLogPoller) RegisterFilter(ctx context.Context, filter logpoller.Filter) error {
	return nil
}
func (m *mockLogPoller) LatestBlock(ctx context.Context) (logpoller.LogPollerBlock, error) {
	return m.latestBlock, m.latestBlockErr
}
func (m *mockLogPoller) LogsWithSigs(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
	return m.logsWithSigs, m.logsWithSigsErr
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

var _ ChannelDefinitionCacheORM = &mockORM{}

type mockORM struct {
	err error

	lastPersistedAddr     common.Address
	lastPersistedDonID    uint32
	lastPersistedVersion  uint32
	lastPersistedDfns     llotypes.ChannelDefinitions
	lastPersistedBlockNum int64
}

func (m *mockORM) LoadChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) (pd *PersistedDefinitions, err error) {
	panic("not implemented")
}
func (m *mockORM) StoreChannelDefinitions(ctx context.Context, addr common.Address, donID, version uint32, dfns llotypes.ChannelDefinitions, blockNum int64) (err error) {
	m.lastPersistedAddr = addr
	m.lastPersistedDonID = donID
	m.lastPersistedVersion = version
	m.lastPersistedDfns = dfns
	m.lastPersistedBlockNum = blockNum
	return m.err
}

func (m *mockORM) CleanupChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) (err error) {
	panic("not implemented")
}

func makeLog(t *testing.T, donID, version uint32, url string, sha [32]byte) logpoller.Log {
	data := makeLogData(t, donID, version, url, sha)
	return logpoller.Log{EventSig: topicNewChannelDefinition, Topics: [][]byte{topicNewChannelDefinition[:], makeDonIDTopic(donID)}, Data: data}
}

func makeLogData(t *testing.T, donID, version uint32, url string, sha [32]byte) []byte {
	event := channelConfigStoreABI.Events[newChannelDefinitionEventName]
	// donID is indexed
	// version, url, sha
	data, err := event.Inputs.NonIndexed().Pack(version, url, sha)
	require.NoError(t, err)
	return data
}

func makeDonIDTopic(donID uint32) []byte {
	return common.BigToHash(big.NewInt(int64(donID))).Bytes()
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
		newLogCh := make(chan *channel_config_store.ChannelConfigStoreNewChannelDefinition, 100)
		cdc := &channelDefinitionCache{donID: donID, lp: lp, lggr: logger.TestSugared(t), newLogCh: newLogCh}

		t.Run("skips if logpoller has no blocks", func(t *testing.T) {
			err := cdc.readLogs()
			assert.NoError(t, err)
			assert.Nil(t, cdc.newLog)
		})
		t.Run("returns error on LatestBlock failure", func(t *testing.T) {
			lp.latestBlockErr = errors.New("test error")

			err := cdc.readLogs()
			assert.EqualError(t, err, "test error")
			assert.Nil(t, cdc.newLog)
		})
		t.Run("does nothing if LatestBlock older or the same as current channel definitions block", func(t *testing.T) {
			lp.latestBlockErr = nil
			lp.latestBlock = logpoller.LogPollerBlock{BlockNumber: 42}
			cdc.definitionsBlockNum = 43

			err := cdc.readLogs()
			assert.NoError(t, err)
			assert.Nil(t, cdc.newLog)
		})
		t.Run("returns error if LogsWithSigs fails", func(t *testing.T) {
			cdc.definitionsBlockNum = 0
			lp.logsWithSigsErr = errors.New("test error 2")

			err := cdc.readLogs()
			assert.EqualError(t, err, "test error 2")
			assert.Nil(t, cdc.newLog)
		})
		t.Run("ignores logs with different topic", func(t *testing.T) {
			lp.logsWithSigsErr = nil
			lp.logsWithSigs = []logpoller.Log{{EventSig: common.Hash{1, 2, 3, 4}}}

			err := cdc.readLogs()
			assert.NoError(t, err)
			assert.Nil(t, cdc.newLog)
		})
		t.Run("returns error if log is malformed", func(t *testing.T) {
			lp.logsWithSigsErr = nil
			lp.logsWithSigs = []logpoller.Log{{EventSig: topicNewChannelDefinition}}

			err := cdc.readLogs()
			assert.EqualError(t, err, "failed to unpack log data: abi: attempting to unmarshal an empty string while arguments are expected")
			assert.Nil(t, cdc.newLog)
		})
		t.Run("sets definitions and sends on channel if LogsWithSigs returns new event with a later version", func(t *testing.T) {
			lp.logsWithSigsErr = nil
			lp.logsWithSigs = []logpoller.Log{makeLog(t, donID, uint32(43), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4})}

			err := cdc.readLogs()
			require.NoError(t, err)
			require.NotNil(t, cdc.newLog)
			assert.Equal(t, uint32(43), cdc.newLog.Version)
			assert.Equal(t, "http://example.com/xxx.json", cdc.newLog.Url)
			assert.Equal(t, [32]byte{1, 2, 3, 4}, cdc.newLog.Sha)
			assert.Equal(t, int64(donID), cdc.newLog.DonId.Int64())

			func() {
				for {
					select {
					case log := <-newLogCh:
						assert.Equal(t, cdc.newLog, log)
					default:
						return
					}
				}
			}()
		})
		t.Run("does nothing if version older or the same as the one currently set", func(t *testing.T) {
			lp.logsWithSigsErr = nil
			lp.logsWithSigs = []logpoller.Log{
				makeLog(t, donID, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
				makeLog(t, donID, uint32(43), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
			}

			err := cdc.readLogs()
			require.NoError(t, err)
			assert.Equal(t, uint32(43), cdc.newLog.Version)
		})
		t.Run("in case of multiple logs, takes the latest", func(t *testing.T) {
			lp.logsWithSigsErr = nil
			lp.logsWithSigs = []logpoller.Log{
				makeLog(t, donID, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
				makeLog(t, donID, uint32(45), "http://example.com/xxx2.json", [32]byte{2, 2, 3, 4}),
				makeLog(t, donID, uint32(44), "http://example.com/xxx3.json", [32]byte{3, 2, 3, 4}),
				makeLog(t, donID, uint32(43), "http://example.com/xxx4.json", [32]byte{4, 2, 3, 4}),
			}

			err := cdc.readLogs()
			require.NoError(t, err)
			assert.Equal(t, uint32(45), cdc.newLog.Version)
			assert.Equal(t, "http://example.com/xxx2.json", cdc.newLog.Url)
			assert.Equal(t, [32]byte{2, 2, 3, 4}, cdc.newLog.Sha)
			assert.Equal(t, int64(donID), cdc.newLog.DonId.Int64())

			func() {
				for {
					select {
					case log := <-newLogCh:
						assert.Equal(t, cdc.newLog, log)
					default:
						return
					}
				}
			}()
		})
		t.Run("ignores logs with incorrect don ID", func(t *testing.T) {
			lp.logsWithSigsErr = nil
			lp.logsWithSigs = []logpoller.Log{
				makeLog(t, donID+1, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4}),
			}

			err := cdc.readLogs()
			require.NoError(t, err)
			assert.Equal(t, uint32(45), cdc.newLog.Version)

			func() {
				for {
					select {
					case log := <-newLogCh:
						t.Fatal("did not expect log with wrong donID, got: ", log)
					default:
						return
					}
				}
			}()
		})
		t.Run("ignores logs with wrong number of topics", func(t *testing.T) {
			lp.logsWithSigsErr = nil
			lg := makeLog(t, donID, uint32(42), "http://example.com/xxx.json", [32]byte{1, 2, 3, 4})
			lg.Topics = lg.Topics[:1]
			lp.logsWithSigs = []logpoller.Log{lg}

			err := cdc.readLogs()
			require.NoError(t, err)
			assert.Equal(t, uint32(45), cdc.newLog.Version)

			func() {
				for {
					select {
					case log := <-newLogCh:
						t.Fatal("did not expect log with missing topics, got: ", log)
					default:
						return
					}
				}
			}()
		})
	})

	t.Run("fetchChannelDefinitions", func(t *testing.T) {
		c := &mockHTTPClient{}
		cdc := &channelDefinitionCache{
			lggr:      logger.TestSugared(t),
			client:    c,
			httpLimit: 2048,
		}

		t.Run("nil ctx returns error", func(t *testing.T) {
			_, err := cdc.fetchChannelDefinitions(nil, "notvalid://foos", [32]byte{}) //nolint
			assert.EqualError(t, err, "failed to create http.Request; net/http: nil Context")
		})

		t.Run("networking error while making request returns error", func(t *testing.T) {
			ctx := tests.Context(t)
			c.resp = nil
			c.err = errors.New("http request failed")

			_, err := cdc.fetchChannelDefinitions(ctx, "http://example.com/definitions.json", [32]byte{})
			assert.EqualError(t, err, "error making http request: http request failed")
		})

		t.Run("server returns 500 returns error", func(t *testing.T) {
			ctx := tests.Context(t)
			c.err = nil
			c.resp = &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte{1, 2, 3}))}

			_, err := cdc.fetchChannelDefinitions(ctx, "http://example.com/definitions.json", [32]byte{})
			assert.EqualError(t, err, "got error from http://example.com/definitions.json: (status code: 500, response body: \x01\x02\x03)")
		})

		var largeBody = make([]byte, 2048)
		for i := range largeBody {
			largeBody[i] = 'a'
		}

		t.Run("server returns 404 returns error (and does not log entirety of huge response body)", func(t *testing.T) {
			ctx := tests.Context(t)
			c.err = nil
			c.resp = &http.Response{StatusCode: 404, Body: io.NopCloser(bytes.NewReader(largeBody))}

			_, err := cdc.fetchChannelDefinitions(ctx, "http://example.com/definitions.json", [32]byte{})
			assert.EqualError(t, err, "got error from http://example.com/definitions.json: (status code: 404, error reading response body: http: request body too large, response body: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa)")
		})

		var hugeBody = make([]byte, 8096)
		c.resp.Body = io.NopCloser(bytes.NewReader(hugeBody))

		t.Run("server returns body that is too large", func(t *testing.T) {
			ctx := tests.Context(t)
			c.err = nil
			c.resp = &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(hugeBody))}

			_, err := cdc.fetchChannelDefinitions(ctx, "http://example.com/definitions.json", [32]byte{})
			assert.EqualError(t, err, "failed to read from body: http: request body too large")
		})

		t.Run("server returns invalid JSON returns error", func(t *testing.T) {
			ctx := tests.Context(t)
			c.err = nil
			c.resp = &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader([]byte{1, 2, 3}))}

			_, err := cdc.fetchChannelDefinitions(ctx, "http://example.com/definitions.json", common.HexToHash("0xfd1780a6fc9ee0dab26ceb4b3941ab03e66ccd970d1db91612c66df4515b0a0a"))
			assert.EqualError(t, err, "failed to decode JSON: invalid character '\\x01' looking for beginning of value")
		})

		t.Run("SHA mismatch returns error", func(t *testing.T) {
			ctx := tests.Context(t)
			c.err = nil
			c.resp = &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader([]byte(`{"foo":"bar"}`)))}

			_, err := cdc.fetchChannelDefinitions(ctx, "http://example.com/definitions.json", [32]byte{})
			assert.EqualError(t, err, "SHA3 mismatch: expected 0000000000000000000000000000000000000000000000000000000000000000, got 4d3304d0d87c27a031cbb6bdf95da79b7b4552c3d0bef2e5a94f50810121e1e0")
		})

		t.Run("valid JSON matching SHA returns channel definitions", func(t *testing.T) {
			ctx := tests.Context(t)
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

			cd, err := cdc.fetchChannelDefinitions(ctx, "http://example.com/definitions.json", common.HexToHash("0x367bbc75f7b6c9fc66a98ea99f837ea7ac4a3c2d6a9ee284de018bd02c41b52d"))
			assert.NoError(t, err)
			assert.Equal(t, llotypes.ChannelDefinitions{0x2a: llotypes.ChannelDefinition{ReportFormat: 0x1, Streams: []llotypes.Stream{llotypes.Stream{StreamID: 0x34, Aggregator: 0x1}, llotypes.Stream{StreamID: 0x35, Aggregator: 0x1}, llotypes.Stream{StreamID: 0x37, Aggregator: 0x3}}, Opts: llotypes.ChannelOpts{0x7b, 0x22, 0x62, 0x61, 0x73, 0x65, 0x55, 0x53, 0x44, 0x46, 0x65, 0x65, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x22, 0x2c, 0x22, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x22, 0x3a, 0x33, 0x36, 0x30, 0x30, 0x2c, 0x22, 0x66, 0x65, 0x65, 0x64, 0x49, 0x64, 0x22, 0x3a, 0x22, 0x30, 0x78, 0x30, 0x30, 0x30, 0x33, 0x36, 0x62, 0x34, 0x61, 0x61, 0x37, 0x65, 0x35, 0x37, 0x63, 0x61, 0x37, 0x62, 0x36, 0x38, 0x61, 0x65, 0x31, 0x62, 0x66, 0x34, 0x35, 0x36, 0x35, 0x33, 0x66, 0x35, 0x36, 0x62, 0x36, 0x35, 0x36, 0x66, 0x64, 0x33, 0x61, 0x61, 0x33, 0x33, 0x35, 0x65, 0x66, 0x37, 0x66, 0x61, 0x65, 0x36, 0x39, 0x36, 0x62, 0x36, 0x36, 0x33, 0x66, 0x31, 0x62, 0x38, 0x34, 0x37, 0x32, 0x22, 0x2c, 0x22, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x22, 0x7d}}}, cd)
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

		t.Run("does nothing if persisted version is up-to-date", func(t *testing.T) {
			ctx := tests.Context(t)
			cdc.definitionsVersion = 42
			cdc.persistedVersion = 42

			memoryVersion, persistedVersion, err := cdc.persist(ctx)
			assert.NoError(t, err)
			assert.Equal(t, uint32(42), memoryVersion)
			assert.Equal(t, uint32(42), persistedVersion)
			assert.Equal(t, uint32(42), cdc.persistedVersion)
		})

		orm := &mockORM{}
		cdc.orm = orm

		t.Run("returns error on db failure and does not update persisted version", func(t *testing.T) {
			ctx := tests.Context(t)
			cdc.persistedVersion = 42
			cdc.definitionsVersion = 43
			orm.err = errors.New("test error")

			memoryVersion, persistedVersion, err := cdc.persist(ctx)
			assert.EqualError(t, err, "test error")
			assert.Equal(t, uint32(43), memoryVersion)
			assert.Equal(t, uint32(42), persistedVersion)
			assert.Equal(t, uint32(42), cdc.persistedVersion)
		})

		t.Run("updates persisted version on success", func(t *testing.T) {
			ctx := tests.Context(t)
			cdc.definitionsVersion = 43
			orm.err = nil

			memoryVersion, persistedVersion, err := cdc.persist(ctx)
			assert.NoError(t, err)
			assert.Equal(t, uint32(43), memoryVersion)
			assert.Equal(t, uint32(43), persistedVersion)
			assert.Equal(t, uint32(43), cdc.persistedVersion)

			assert.Equal(t, cdc.addr, orm.lastPersistedAddr)
			assert.Equal(t, cdc.donID, orm.lastPersistedDonID)
			assert.Equal(t, cdc.persistedVersion, orm.lastPersistedVersion)
			assert.Equal(t, cdc.definitions, orm.lastPersistedDfns)
			assert.Equal(t, cdc.definitionsBlockNum, orm.lastPersistedBlockNum)
		})
	})
}

func Test_filterName(t *testing.T) {
	s := filterName(common.Address{1, 2, 3}, 654)
	assert.Equal(t, "OCR3 LLO ChannelDefinitionCachePoller - 0x0102030000000000000000000000000000000000:654", s)
}
