package streams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury"
	v02 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury/v02"
	v03 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury/v03"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	autov2common "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_v21_plus_common"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type MockMercuryConfigProvider struct {
	mock.Mock
}

func (m *MockMercuryConfigProvider) Credentials() *types.MercuryCredentials {
	mc := &types.MercuryCredentials{
		LegacyURL: "https://google.old.com",
		URL:       "https://google.com",
		Username:  "FakeClientID",
		Password:  "FakeClientKey",
	}
	return mc
}

func (m *MockMercuryConfigProvider) IsUpkeepAllowed(s string) (interface{}, bool) {
	args := m.Called(s)
	return args.Get(0), args.Bool(1)
}

func (m *MockMercuryConfigProvider) SetUpkeepAllowed(s string, i interface{}, d time.Duration) {
	m.Called(s, i, d)
}

func (m *MockMercuryConfigProvider) GetPluginRetry(s string) (interface{}, bool) {
	args := m.Called(s)
	return args.Get(0), args.Bool(1)
}

func (m *MockMercuryConfigProvider) SetPluginRetry(s string, i interface{}, d time.Duration) {
	m.Called(s, i, d)
}

type MockBlockSubscriber struct {
	mock.Mock
}

func (b *MockBlockSubscriber) LatestBlock() *ocr2keepers.BlockKey {
	return nil
}

type MockHttpClient struct {
	mock.Mock
}

func (mock *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := mock.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

type mockRegistry struct {
	GetUpkeepPrivilegeConfigFn func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)
	CheckCallbackFn            func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error)
}

func (r *mockRegistry) Address() common.Address {
	return common.HexToAddress("0x6cA639822c6C241Fa9A7A6b5032F6F7F1C513CAD")
}

func (r *mockRegistry) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	return r.GetUpkeepPrivilegeConfigFn(opts, upkeepId)
}

func (r *mockRegistry) CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
	return r.CheckCallbackFn(opts, id, values, extraData)
}

// setups up a streams object for tests.
func setupStreams(t *testing.T) *streams {
	lggr := logger.TestLogger(t)
	mercuryConfig := new(MockMercuryConfigProvider)
	blockSubscriber := new(MockBlockSubscriber)
	registry := &mockRegistry{}
	client := clienttest.NewClient(t)

	streams := NewStreamsLookup(
		mercuryConfig,
		blockSubscriber,
		client,
		registry,
		lggr,
	)
	return streams
}

func TestStreams_CheckErrorHandler(t *testing.T) {
	upkeepId := big.NewInt(123456789)
	blockNumber := uint64(999)
	tests := []struct {
		name         string
		lookup       *mercury.StreamsLookup
		checkResults []ocr2keepers.CheckResult
		errCode      encoding.ErrCode

		callbackResp []byte
		callbackErr  error

		upkeepNeeded bool
		performData  []byte
		wantErr      assert.ErrorAssertionFunc

		state     encoding.PipelineExecutionState
		retryable bool
	}{
		{
			name: "success - empty extra data",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"ETD-USD", "BTC-ETH"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(100),
					ExtraData:    []byte{48, 120, 48, 48},
				},
				UpkeepId: upkeepId,
				Block:    blockNumber,
			},
			checkResults: []ocr2keepers.CheckResult{
				{},
			},
			errCode:      encoding.ErrCodeNil,
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			upkeepNeeded: true,
			performData:  []byte{48, 120, 48, 48},
			wantErr:      assert.NoError,
		},
		{
			name: "success - with extra data",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(18952430),
					// this is the address of precompile contract ArbSys(0x0000000000000000000000000000000000000064)
					ExtraData: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
				Block:    blockNumber,
			},
			checkResults: []ocr2keepers.CheckResult{
				{},
			},
			errCode:      encoding.ErrCodeNil,
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			upkeepNeeded: true,
			performData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			wantErr:      assert.NoError,
		},
		{
			name: "failure - checkCallback eth_call failure",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"ETD-USD", "BTC-ETH"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(100),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				UpkeepId: upkeepId,
				Block:    blockNumber,
			},
			checkResults: []ocr2keepers.CheckResult{
				{},
			},
			errCode:      encoding.ErrCodeNil,
			callbackResp: []byte{},
			callbackErr:  errors.New("bad response"),
			wantErr:      assert.Error,
			state:        encoding.RpcFlakyFailure,
			retryable:    true,
			upkeepNeeded: false,
		},
		{
			name: "failure - unpack error because of invalid response bytes from checkCallback eth_call",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"ETD-USD", "BTC-ETH"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(100),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				UpkeepId: upkeepId,
				Block:    blockNumber,
			},
			checkResults: []ocr2keepers.CheckResult{
				{},
			},
			errCode:      encoding.ErrCodeNil,
			callbackResp: []byte{0x01, 0xAB, 0x45, 0xFF, 0x32}, // invalid response bytes
			wantErr:      assert.Error,
			state:        encoding.PackUnpackDecodeFailed,
			retryable:    false,
			upkeepNeeded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(clienttest.Client)
			s := setupStreams(t)
			defer s.Close()

			userPayload, err := s.packer.PackUserCheckErrorHandler(tt.errCode, tt.lookup.ExtraData)
			require.NoError(t, err)
			payload, err := s.abi.Pack("executeCallback", tt.lookup.UpkeepId, userPayload)
			require.NoError(t, err)

			args := map[string]interface{}{
				"from": zeroAddress,
				"to":   s.registry.Address().Hex(),
				"data": hexutil.Bytes(payload),
			}
			client.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", args, hexutil.EncodeUint64(tt.lookup.Block)).Return(tt.callbackErr).
				Run(func(args mock.Arguments) {
					by := args.Get(1).(*hexutil.Bytes)
					*by = tt.callbackResp
				}).Once()
			s.client = client

			err = s.CheckErrorHandler(testutils.Context(t), tt.errCode, tt.lookup, tt.checkResults, 0)
			tt.wantErr(t, err, fmt.Sprintf("Error assertion failed: %v", tt.name))
			assert.Equal(t, uint8(tt.state), tt.checkResults[0].PipelineExecutionState)
			assert.Equal(t, tt.retryable, tt.checkResults[0].Retryable)
			assert.Equal(t, tt.upkeepNeeded, tt.checkResults[0].Eligible)
			assert.Equal(t, tt.performData, tt.checkResults[0].PerformData)
		})
	}
}

func TestStreams_CheckCallback(t *testing.T) {
	upkeepId := big.NewInt(123456789)
	bn := uint64(999)
	bs := []byte{183, 114, 215, 10, 0, 0, 0, 0, 0, 0}
	values := [][]byte{bs}
	tests := []struct {
		name       string
		lookup     *mercury.StreamsLookup
		input      []ocr2keepers.CheckResult
		values     [][]byte
		statusCode int

		callbackResp []byte
		callbackErr  error

		upkeepNeeded bool
		performData  []byte
		wantErr      assert.ErrorAssertionFunc

		state     encoding.PipelineExecutionState
		retryable bool
		registry  streamRegistry
	}{
		{
			name: "success - empty extra data",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"ETD-USD", "BTC-ETH"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(100),
					ExtraData:    []byte{48, 120, 48, 48},
				},
				UpkeepId: upkeepId,
				Block:    bn,
			},
			input: []ocr2keepers.CheckResult{
				{},
			},
			values:       values,
			statusCode:   http.StatusOK,
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			upkeepNeeded: true,
			performData:  []byte{48, 120, 48, 48},
			wantErr:      assert.NoError,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(`{"mercuryEnabled":true}`), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{
						UpkeepNeeded: true,
						PerformData:  []byte{48, 120, 48, 48},
					}, nil
				},
			},
		},
		{
			name: "success - with extra data",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(18952430),
					// this is the address of precompile contract ArbSys(0x0000000000000000000000000000000000000064)
					ExtraData: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
				Block:    bn,
			},
			input: []ocr2keepers.CheckResult{
				{},
			},
			values:       values,
			statusCode:   http.StatusOK,
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			upkeepNeeded: true,
			performData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			wantErr:      assert.NoError,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(`{"mercuryEnabled":true}`), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{
						UpkeepNeeded: true,
						PerformData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
					}, nil
				},
			},
		},
		{
			name: "failure - bad response",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"ETD-USD", "BTC-ETH"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(100),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				UpkeepId: upkeepId,
				Block:    bn,
			},
			input: []ocr2keepers.CheckResult{
				{},
			},
			values:       values,
			statusCode:   http.StatusOK,
			callbackResp: []byte{},
			callbackErr:  errors.New("bad response"),
			wantErr:      assert.Error,
			state:        encoding.RpcFlakyFailure,
			retryable:    true,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(`{"mercuryEnabled":true}`), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{}, errors.New("bad response")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(clienttest.Client)
			s := setupStreams(t)
			defer s.Close()
			payload, err := s.abi.Pack("checkCallback", tt.lookup.UpkeepId, values, tt.lookup.ExtraData)
			require.NoError(t, err)
			args := map[string]interface{}{
				"from": zeroAddress,
				"to":   s.registry.Address().Hex(),
				"data": hexutil.Bytes(payload),
			}
			client.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", args, hexutil.EncodeUint64(tt.lookup.Block)).Return(tt.callbackErr).
				Run(func(args mock.Arguments) {
					by := args.Get(1).(*hexutil.Bytes)
					*by = tt.callbackResp
				}).Once()
			s.client = client

			err = s.CheckCallback(testutils.Context(t), tt.values, tt.lookup, tt.input, 0)
			tt.wantErr(t, err, fmt.Sprintf("Error assertion failed: %v", tt.name))
			assert.Equal(t, uint8(tt.state), tt.input[0].PipelineExecutionState)
			assert.Equal(t, tt.retryable, tt.input[0].Retryable)
			assert.Equal(t, tt.upkeepNeeded, tt.input[0].Eligible)
		})
	}
}

func TestStreams_AllowedToUseMercury(t *testing.T) {
	upkeepId, ok := new(big.Int).SetString("71022726777042968814359024671382968091267501884371696415772139504780367423725", 10)
	assert.True(t, ok, t.Name())
	tests := []struct {
		name       string
		cached     bool
		allowed    bool
		ethCallErr error
		err        error
		state      encoding.PipelineExecutionState
		reason     encoding.UpkeepFailureReason
		registry   streamRegistry
		retryable  bool
		config     []byte
	}{
		{
			name:    "success - allowed via cache",
			cached:  true,
			allowed: true,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return nil, nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{}, nil
				},
			},
		},
		{
			name:    "success - allowed via fetching privilege config",
			allowed: true,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(`{"mercuryEnabled":true}`), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{}, nil
				},
			},
		},
		{
			name:    "success - not allowed via cache",
			cached:  true,
			allowed: false,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return nil, nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{}, nil
				},
			},
		},
		{
			name:    "success - not allowed via fetching privilege config",
			allowed: false,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(`{"mercuryEnabled":false}`), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{}, nil
				},
			},
		},
		{
			name:   "failure - cannot unmarshal privilege config",
			err:    errors.New("failed to unmarshal privilege config: invalid character '\\x00' looking for beginning of value"),
			state:  encoding.PrivilegeConfigUnmarshalError,
			config: []byte{0, 1},
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(`{`), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{}, nil
				},
			},
		},
		{
			name:       "failure - flaky RPC",
			retryable:  true,
			err:        errors.New("failed to get upkeep privilege config: flaky RPC"),
			state:      encoding.RpcFlakyFailure,
			ethCallErr: errors.New("flaky RPC"),
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return nil, errors.New("flaky RPC")
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{}, nil
				},
			},
		},
		{
			name:   "failure - empty upkeep privilege config",
			err:    errors.New("upkeep privilege config is empty"),
			reason: encoding.UpkeepFailureReasonMercuryAccessNotAllowed,
			config: []byte{},
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(``), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{}, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupStreams(t)
			defer s.Close()
			s.registry = tt.registry

			client := new(clienttest.Client)
			s.client = client

			mc := new(MockMercuryConfigProvider)
			mc.On("IsUpkeepAllowed", mock.Anything).Return(tt.allowed, tt.cached).Once()
			mc.On("SetUpkeepAllowed", mock.Anything, mock.Anything, mock.Anything).Return().Once()
			s.mercuryConfig = mc

			if !tt.cached {
				if tt.err != nil {
					bContractCfg, err := s.abi.Methods["getUpkeepPrivilegeConfig"].Outputs.PackValues([]interface{}{tt.config})
					require.NoError(t, err)

					payload, err := s.abi.Pack("getUpkeepPrivilegeConfig", upkeepId)
					require.NoError(t, err)

					args := map[string]interface{}{
						"to":   s.registry.Address().Hex(),
						"data": hexutil.Bytes(payload),
					}

					client.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", args, mock.AnythingOfType("string")).
						Return(tt.ethCallErr).
						Run(func(args mock.Arguments) {
							b := args.Get(1).(*hexutil.Bytes)
							*b = bContractCfg
						}).Once()
				} else {
					cfg := UpkeepPrivilegeConfig{MercuryEnabled: tt.allowed}
					bCfg, err := json.Marshal(cfg)
					require.NoError(t, err)

					bContractCfg, err := s.abi.Methods["getUpkeepPrivilegeConfig"].Outputs.PackValues([]interface{}{bCfg})
					require.NoError(t, err)

					payload, err := s.abi.Pack("getUpkeepPrivilegeConfig", upkeepId)
					require.NoError(t, err)

					args := map[string]interface{}{
						"to":   s.registry.Address().Hex(),
						"data": hexutil.Bytes(payload),
					}

					client.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", args, mock.AnythingOfType("string")).Return(nil).
						Run(func(args mock.Arguments) {
							b := args.Get(1).(*hexutil.Bytes)
							*b = bContractCfg
						}).Once()
				}
			}

			opts := &bind.CallOpts{
				BlockNumber: big.NewInt(10),
			}

			state, reason, retryable, allowed, err := s.AllowedToUseMercury(opts, upkeepId)
			if tt.err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err.Error())
			}
			assert.Equal(t, tt.allowed, allowed)
			assert.Equal(t, tt.state, state)
			assert.Equal(t, tt.reason, reason)
			assert.Equal(t, tt.retryable, retryable)
		})
	}
}

func TestStreams_StreamsLookup(t *testing.T) {
	upkeepId, ok := new(big.Int).SetString("71022726777042968814359024671382968091267501884371696415772139504780367423725", 10)
	var upkeepIdentifier [32]byte
	copy(upkeepIdentifier[:], upkeepId.Bytes())
	assert.True(t, ok, t.Name())
	blockNum := ocr2keepers.BlockNumber(37974374)
	tests := []struct {
		name              string
		input             []ocr2keepers.CheckResult
		blobs             map[string]string
		callbackResp      []byte
		expectedResults   []ocr2keepers.CheckResult
		callbackNeeded    bool
		extraData         []byte
		checkCallbackResp []byte
		values            [][]byte
		cachedAdminCfg    bool
		hasError          bool
		hasPermission     bool
		v3                bool
		registry          streamRegistry
	}{
		{
			name: "success - happy path no cache",
			input: []ocr2keepers.CheckResult{
				{
					PerformData: hexutil.MustDecode("0xf055e4a200000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000000966656564496448657800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000423078343535343438326435353533343432643431353234323439353435323535346432643534343535333534346534353534303030303030303030303030303030300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000042307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b626c6f636b4e756d62657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: blockNum,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonTargetCheckReverted),
				},
			},
			blobs: map[string]string{
				"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000": "0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa3",
				"0x4254432d5553442d415242495452554d2d544553544e45540000000000000000": "0x0006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d",
			},
			cachedAdminCfg:    false,
			extraData:         hexutil.MustDecode("0x0000000000000000000000000000000000000064"),
			callbackNeeded:    true,
			checkCallbackResp: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000063a400000000000000000000000000000000000000000000000000000000000006e0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000002e000066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa300000000000000000000000000000000000000000000000000000000000002e00006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
			values:            [][]byte{hexutil.MustDecode("0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa3"), hexutil.MustDecode("0x0006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d")},
			expectedResults: []ocr2keepers.CheckResult{
				{
					Eligible:    true,
					PerformData: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000002e000066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa300000000000000000000000000000000000000000000000000000000000002e00006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: blockNum,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonNone),
				},
			},
			hasPermission: true,
			v3:            false,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(`{"mercuryEnabled":true}`), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{
						UpkeepNeeded: true,
						PerformData:  hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000002e000066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa300000000000000000000000000000000000000000000000000000000000002e00006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
					}, nil
				},
			},
		},
		{
			name: "two CheckResults: Mercury success E2E and No Mercury because of insufficient balance",
			input: []ocr2keepers.CheckResult{
				{
					PerformData: []byte{},
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: 26046145,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonInsufficientBalance),
				},
				{
					PerformData: hexutil.MustDecode("0xf055e4a200000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000000966656564496448657800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000423078343535343438326435353533343432643431353234323439353435323535346432643534343535333534346534353534303030303030303030303030303030300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000042307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b626c6f636b4e756d62657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: blockNum,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonTargetCheckReverted),
				},
			},
			blobs: map[string]string{
				"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000": "0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa3",
				"0x4254432d5553442d415242495452554d2d544553544e45540000000000000000": "0x0006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d",
			},
			cachedAdminCfg:    false,
			extraData:         hexutil.MustDecode("0x0000000000000000000000000000000000000064"),
			callbackNeeded:    true,
			checkCallbackResp: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000063a400000000000000000000000000000000000000000000000000000000000006e0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000002e000066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa300000000000000000000000000000000000000000000000000000000000002e00006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
			values:            [][]byte{hexutil.MustDecode("0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa3"), hexutil.MustDecode("0x0006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d")},
			expectedResults: []ocr2keepers.CheckResult{
				{
					Eligible:    false,
					PerformData: []byte{},
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: 26046145,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonInsufficientBalance),
				},
				{
					Eligible:    true,
					PerformData: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000002e000066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa300000000000000000000000000000000000000000000000000000000000002e00006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: blockNum,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonNone),
				},
			},
			hasPermission: true,
			v3:            false,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(`{"mercuryEnabled":true}`), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{
						UpkeepNeeded: true,
						PerformData:  hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000002e000066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa300000000000000000000000000000000000000000000000000000000000002e00006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
					}, nil
				},
			},
		},
		{
			name: "success - happy path no cache - v0.3",
			input: []ocr2keepers.CheckResult{
				{
					PerformData: hexutil.MustDecode("0xf055e4a200000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000000766656564494473000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000423078343535343438326435353533343432643431353234323439353435323535346432643534343535333534346534353534303030303030303030303030303030300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000042307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b626c6f636b4e756d62657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: blockNum,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonTargetCheckReverted),
				},
			},
			blobs: map[string]string{
				"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000": "0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa3",
				"0x4254432d5553442d415242495452554d2d544553544e45540000000000000000": "0x0006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d",
			},
			cachedAdminCfg:    false,
			extraData:         hexutil.MustDecode("0x0000000000000000000000000000000000000064"),
			callbackNeeded:    true,
			checkCallbackResp: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000063a400000000000000000000000000000000000000000000000000000000000006e0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000002e000066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa300000000000000000000000000000000000000000000000000000000000002e00006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
			values:            [][]byte{hexutil.MustDecode("0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa3"), hexutil.MustDecode("0x0006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d")},
			expectedResults: []ocr2keepers.CheckResult{
				{
					Eligible:    true,
					PerformData: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000002e000066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa300000000000000000000000000000000000000000000000000000000000002e00006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: blockNum,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonNone),
				},
			},
			hasPermission: true,
			v3:            true,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(`{"mercuryEnabled":true}`), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{
						UpkeepNeeded: true,
						PerformData:  hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000006a000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000002e000066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000004555638000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280010100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000000269ecbb83b000000000000000000000000000000000000000000000000000000269e4a4e14000000000000000000000000000000000000000000000000000000269f4d0edb000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002381e91cffa9502c20de1ddcee350db3f715a5ab449448e3184a5b03c682356c6e2115f20663b3731e373cf33465a96da26f2876debb548f281e62e48f64c374200000000000000000000000000000000000000000000000000000000000000027db99e34135098d4e0bb9ae143ec9cd72fd63150c6d0cc5b38f4aa1aa42408377e8fe8e5ac489c9b7f62ff5aa7b05d2e892e7dee4cac631097247969b3b03fa300000000000000000000000000000000000000000000000000000000000002e00006da4a86c4933dd4a87b21dd2871aea29f706bcde43c70039355ac5b664fb5000000000000000000000000000000000000000000000000000000000454d118000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204254432d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064f0d4a0000000000000000000000000000000000000000000000000000002645f00877a000000000000000000000000000000000000000000000000000002645e1e1010000000000000000000000000000000000000000000000000000002645fe2fee4000000000000000000000000000000000000000000000000000000000243716664b42d20423a47fb13ad3098b49b37f667548e6745fff958b663afe25a845f6100000000000000000000000000000000000000000000000000000000024371660000000000000000000000000000000000000000000000000000000064f0d4a00000000000000000000000000000000000000000000000000000000000000002a0373c0bce7393673f819eb9681cac2773c2d718ce933eb858252195b17a9c832d7b0bee173c02c3c25fb65912b8b13b9302ede8423bab3544cb7a8928d5eb3600000000000000000000000000000000000000000000000000000000000000027d7b79d7646383a5dbf51edf14d53bd3ad0a9f3ca8affab3165e89d3ddce9cb17b58e892fafe4ecb24d2fde07c6a756029e752a5114c33c173df4e7d309adb4d00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
					}, nil
				},
			},
		},
		{
			name: "skip - failure reason is insufficient balance",
			input: []ocr2keepers.CheckResult{
				{
					PerformData: []byte{},
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: 26046145,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonInsufficientBalance),
				},
			},
			expectedResults: []ocr2keepers.CheckResult{
				{
					Eligible:    false,
					PerformData: []byte{},
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: 26046145,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonInsufficientBalance),
				},
			},
			hasError: true,
			registry: &mockRegistry{
				GetUpkeepPrivilegeConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return []byte(`{"mercuryEnabled":true}`), nil
				},
				CheckCallbackFn: func(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error) {
					return autov2common.CheckCallback{
						UpkeepNeeded: true,
						PerformData:  []byte{},
					}, nil
				},
			},
		},
		{
			name: "skip - invalid revert data",
			input: []ocr2keepers.CheckResult{
				{
					PerformData: []byte{},
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: 26046145,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonTargetCheckReverted),
				},
			},
			expectedResults: []ocr2keepers.CheckResult{
				{
					Eligible:    false,
					PerformData: []byte{},
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: 26046145,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonTargetCheckReverted),
				},
			},
			hasError: true,
		},
		{
			name: "failure - invalid mercury version",
			input: []ocr2keepers.CheckResult{
				{
					// This Perform data contains invalid FeedParamKey: {feedIdHex:RandomString [ETD-USD BTC-ETH] blockNumber 100 [48 120 48 48]}
					PerformData: hexutil.MustDecode("0xf055e4a200000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000166665656449644865783a52616e646f6d537472696e670000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000074554442d5553440000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000074254432d45544800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b626c6f636b4e756d62657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000043078303000000000000000000000000000000000000000000000000000000000"),
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: blockNum,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonTargetCheckReverted),
				},
			},
			expectedResults: []ocr2keepers.CheckResult{
				{
					PerformData: hexutil.MustDecode("0xf055e4a200000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000166665656449644865783a52616e646f6d537472696e670000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000074554442d5553440000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000074254432d45544800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b626c6f636b4e756d62657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000043078303000000000000000000000000000000000000000000000000000000000"),
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: blockNum,
					},
					IneligibilityReason: uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput),
				},
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupStreams(t)
			defer s.Close()
			s.registry = tt.registry

			client := new(clienttest.Client)
			s.client = client

			mc := new(MockMercuryConfigProvider)
			mc.On("IsUpkeepAllowed", mock.Anything).Return(tt.hasPermission, tt.cachedAdminCfg).Once()
			mc.On("SetUpkeepAllowed", mock.Anything, mock.Anything, mock.Anything).Return().Once()
			s.mercuryConfig = mc

			if !tt.cachedAdminCfg && !tt.hasError {
				cfg := UpkeepPrivilegeConfig{MercuryEnabled: tt.hasPermission}
				bCfg, err := json.Marshal(cfg)
				require.NoError(t, err)

				bContractCfg, err := s.abi.Methods["getUpkeepPrivilegeConfig"].Outputs.PackValues([]interface{}{bCfg})
				require.NoError(t, err)

				payload, err := s.abi.Pack("getUpkeepPrivilegeConfig", upkeepId)
				require.NoError(t, err)

				args := map[string]interface{}{
					"to":   s.registry.Address().Hex(),
					"data": hexutil.Bytes(payload),
				}

				client.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", args, mock.AnythingOfType("string")).Return(nil).
					Run(func(args mock.Arguments) {
						b := args.Get(1).(*hexutil.Bytes)
						*b = bContractCfg
					}).Once()
			}

			if len(tt.blobs) > 0 {
				if tt.v3 {
					hc03 := new(MockHttpClient)
					v03HttpClient := v03.NewClient(s.mercuryConfig, hc03, s.threadCtrl, s.lggr)
					s.v03Client = v03HttpClient

					mr1 := v03.MercuryV03Response{
						Reports: []v03.MercuryV03Report{{FullReport: tt.blobs["0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"]}, {FullReport: tt.blobs["0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"]}}}
					b1, err := json.Marshal(mr1)
					assert.NoError(t, err)
					resp1 := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(b1)),
					}
					hc03.On("Do", mock.Anything).Return(resp1, nil).Once()
				} else {
					hc02 := new(MockHttpClient)
					v02HttpClient := v02.NewClient(s.mercuryConfig, hc02, s.threadCtrl, s.lggr)
					s.v02Client = v02HttpClient

					mr1 := v02.MercuryV02Response{ChainlinkBlob: tt.blobs["0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"]}
					b1, err := json.Marshal(mr1)
					assert.NoError(t, err)
					resp1 := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(b1)),
					}
					mr2 := v02.MercuryV02Response{ChainlinkBlob: tt.blobs["0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"]}
					b2, err := json.Marshal(mr2)
					assert.NoError(t, err)
					resp2 := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(b2)),
					}

					hc02.On("Do", mock.MatchedBy(func(req *http.Request) bool {
						return strings.Contains(req.URL.String(), "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000")
					})).Return(resp2, nil).Once()

					hc02.On("Do", mock.MatchedBy(func(req *http.Request) bool {
						return strings.Contains(req.URL.String(), "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000")
					})).Return(resp1, nil).Once()
				}
			}

			if tt.callbackNeeded {
				payload, err := s.abi.Pack("checkCallback", upkeepId, tt.values, tt.extraData)
				require.NoError(t, err)
				args := map[string]interface{}{
					"from": zeroAddress,
					"to":   s.registry.Address().Hex(),
					"data": hexutil.Bytes(payload),
				}
				client.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", args, hexutil.EncodeUint64(uint64(blockNum))).Return(nil).
					Run(func(args mock.Arguments) {
						b := args.Get(1).(*hexutil.Bytes)
						*b = tt.checkCallbackResp
					}).Once()
			}

			got := s.Lookup(testutils.Context(t), tt.input)
			assert.Equal(t, tt.expectedResults, got, tt.name)
		})
	}
}
