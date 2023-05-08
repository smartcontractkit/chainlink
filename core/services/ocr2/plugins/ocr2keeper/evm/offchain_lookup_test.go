package evm

import (
	"bytes"
	"context"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2keepers/pkg/chain"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_lookup_compatible_interface"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm/mocks"
)

// setups up an evm registry for tests.
func setupEVMRegistry(t *testing.T) *EvmRegistry {
	lggr := logger.TestLogger(t)
	addr := common.Address{}
	keeperRegistryABI, err := abi.JSON(strings.NewReader(keeper_registry_wrapper2_0.KeeperRegistryABI))
	require.Nil(t, err, "need registry abi")
	mercuryCompatibleABI, err := abi.JSON(strings.NewReader(mercury_lookup_compatible_interface.MercuryLookupCompatibleInterfaceABI))
	require.Nil(t, err, "need mercury abi")
	upkeepInfoCache, cooldownCache, apiErrCache := setupCaches(DefaultUpkeepExpiration, DefaultCooldownExpiration, DefaultApiErrExpiration, CleanupInterval)
	var headTracker httypes.HeadTracker
	var headBroadcaster httypes.HeadBroadcaster
	var logPoller logpoller.LogPoller
	mockRegistry := mocks.NewRegistry(t)
	mockHttpClient := mocks.NewHttpClient(t)
	client := evmClientMocks.NewClient(t)

	r := &EvmRegistry{
		HeadProvider: HeadProvider{
			ht:     headTracker,
			hb:     headBroadcaster,
			chHead: make(chan types.BlockKey, 1),
		},
		lggr:     lggr,
		poller:   logPoller,
		addr:     addr,
		client:   client,
		txHashes: make(map[string]bool),
		registry: mockRegistry,
		abi:      keeperRegistryABI,
		active:   make(map[string]activeUpkeep),
		packer:   &evmRegistryPackerV2_0{abi: keeperRegistryABI},
		headFunc: func(types.BlockKey) {},
		chLog:    make(chan logpoller.Log, 1000),
		mercury: MercuryConfig{
			cred: &models.MercuryCredentials{
				URL:      "https://google.com",
				Username: "FakeClientID",
				Password: "FakeClientKey",
			},
			abi:           mercuryCompatibleABI,
			upkeepCache:   upkeepInfoCache,
			cooldownCache: cooldownCache,
			apiErrCache:   apiErrCache,
		},
		hc: mockHttpClient,
	}
	return r
}

// helper for mocking the http requests
func (r *EvmRegistry) buildRevertBytesHelper() []byte {
	mercuryErr := r.mercury.abi.Errors["MercuryLookup"]
	mercuryLookupSelector := [4]byte{0x62, 0xe8, 0xa5, 0x0d}
	ml := MercuryLookup{
		feedLabel:  "feedIDStr",
		feeds:      []string{"ETH-USD-ARBITRUM-TESTNET", "BTC-USD-ARBITRUM-TESTNET"},
		queryLabel: "blockNumber",
		query:      big.NewInt(8586948),
		extraData:  []byte{},
	}
	// check if Pack does not add selector for me
	pack, err := mercuryErr.Inputs.Pack(ml.feedLabel, ml.feeds, ml.queryLabel, ml.query, ml.extraData)
	if err != nil {
		log.Fatal("failed to build revert")
	}
	var payload []byte
	payload = append(payload, mercuryLookupSelector[:]...)
	payload = append(payload, pack...)
	return payload
}

func TestEvmRegistry_offchainLookup(t *testing.T) {
	setupRegistry := setupEVMRegistry(t)
	// load json response for testdata
	btcBlob, e := os.ReadFile("./testdata/btc-usd.json")
	assert.Nil(t, e)
	ethBlob, e := os.ReadFile("./testdata/eth-usd.json")
	assert.Nil(t, e)

	upkeepKey := chain.UpkeepKey("8586948|520376062160720574742736856650455852347405918082346589375578334001045521721")
	_, upkeepId, err := blockAndIdFromKey(upkeepKey)
	assert.Nil(t, err, t.Name())
	// builds revert data with mock server query
	revertPerformData := setupRegistry.buildRevertBytesHelper()
	// prepare input upkeepResult
	upkeepResult := types.UpkeepResult{
		Key:              upkeepKey,
		State:            types.NotEligible,
		FailureReason:    UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
		GasUsed:          big.NewInt(27071),
		PerformData:      revertPerformData,
		FastGasWei:       big.NewInt(2000000000),
		LinkNative:       big.NewInt(4391095484380865),
		CheckBlockNumber: 8586947,
		CheckBlockHash:   [32]byte{230, 67, 97, 54, 73, 238, 133, 239, 200, 124, 171, 132, 40, 18, 124, 96, 102, 97, 232, 17, 96, 237, 173, 166, 112, 42, 146, 204, 46, 17, 67, 34},
		ExecuteGas:       5000000,
	}
	upkeepResultReasonOffchain := types.UpkeepResult{
		Key:              upkeepKey,
		State:            types.NotEligible,
		FailureReason:    UPKEEP_FAILURE_REASON_MERCURY_LOOKUP_ERROR,
		GasUsed:          big.NewInt(27071),
		PerformData:      revertPerformData,
		FastGasWei:       big.NewInt(2000000000),
		LinkNative:       big.NewInt(4391095484380865),
		CheckBlockNumber: 8586947,
		CheckBlockHash:   [32]byte{230, 67, 97, 54, 73, 238, 133, 239, 200, 124, 171, 132, 40, 18, 124, 96, 102, 97, 232, 17, 96, 237, 173, 166, 112, 42, 146, 204, 46, 17, 67, 34},
		ExecuteGas:       5000000,
	}
	target := common.HexToAddress("0x79D8aDb571212b922089A48956c54A453D889dBe")
	callbackResp := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	upkeepNeededFalseResp, err := setupRegistry.mercury.abi.Methods["mercuryCallback"].Outputs.Pack(false, []byte{})
	assert.Nil(t, err, t.Name())

	// desired outcomes
	wantPerformData := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	wantUpkeepResult := types.UpkeepResult{
		Key:              upkeepKey,
		State:            types.Eligible,
		FailureReason:    UPKEEP_FAILURE_REASON_NONE,
		GasUsed:          big.NewInt(27071),
		PerformData:      wantPerformData,
		FastGasWei:       big.NewInt(2000000000),
		LinkNative:       big.NewInt(4391095484380865),
		CheckBlockNumber: 8586947,
		CheckBlockHash:   [32]byte{230, 67, 97, 54, 73, 238, 133, 239, 200, 124, 171, 132, 40, 18, 124, 96, 102, 97, 232, 17, 96, 237, 173, 166, 112, 42, 146, 204, 46, 17, 67, 34},
		ExecuteGas:       5000000,
	}
	tests := []struct {
		name  string
		input []types.UpkeepResult

		inCooldown bool

		callbackResp []byte
		callbackErr  error

		upkeepCache   bool
		mockGetUpkeep bool
		upkeepInfo    keeper_registry_wrapper2_0.UpkeepInfo
		upkeepInfoErr error

		want           []types.UpkeepResult
		wantErr        error
		hasHttpCalls   bool
		callbackNeeded bool
	}{
		{
			name:         "success - cached upkeep",
			input:        []types.UpkeepResult{upkeepResult},
			callbackResp: callbackResp,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:     target,
				ExecuteGas: 5000000,
			},

			want:           []types.UpkeepResult{wantUpkeepResult},
			hasHttpCalls:   true,
			callbackNeeded: true,
			upkeepCache:    true,
		},
		{
			name:          "success - no cached upkeep",
			input:         []types.UpkeepResult{upkeepResult},
			callbackResp:  callbackResp,
			mockGetUpkeep: true,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:     target,
				ExecuteGas: 5000000,
			},

			want:           []types.UpkeepResult{wantUpkeepResult},
			hasHttpCalls:   true,
			callbackNeeded: true,
		},
		{
			name: "skip - failure reason",
			input: []types.UpkeepResult{
				{
					Key:           upkeepKey,
					State:         types.NotEligible,
					FailureReason: UPKEEP_FAILURE_REASON_INSUFFICIENT_BALANCE,
					PerformData:   []byte{},
				},
			},

			want: []types.UpkeepResult{
				{
					Key:           upkeepKey,
					State:         types.NotEligible,
					FailureReason: UPKEEP_FAILURE_REASON_INSUFFICIENT_BALANCE,
					PerformData:   []byte{},
				},
			},
		},
		{
			name: "skip - revert data does not decode to offchain lookup",
			input: []types.UpkeepResult{
				{
					Key:           upkeepKey,
					State:         types.NotEligible,
					FailureReason: UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
					PerformData:   []byte{},
				},
			},

			want: []types.UpkeepResult{
				{
					Key:           upkeepKey,
					State:         types.NotEligible,
					FailureReason: UPKEEP_FAILURE_REASON_MERCURY_LOOKUP_ERROR,
					PerformData:   []byte{},
				},
			},
			wantErr: errors.New("unpack error: invalid data for unpacking"),
		},
		{
			name:          "skip - error - no upkeep",
			input:         []types.UpkeepResult{upkeepResult},
			callbackResp:  callbackResp,
			upkeepInfoErr: errors.New("ouch"),

			want:          []types.UpkeepResult{upkeepResultReasonOffchain},
			mockGetUpkeep: true,
			wantErr:       errors.New("ouch"),
		},
		{
			name:          "skip - upkeep not needed",
			input:         []types.UpkeepResult{upkeepResult},
			mockGetUpkeep: true,
			callbackResp:  upkeepNeededFalseResp,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:     target,
				ExecuteGas: 5000000,
			},

			want: []types.UpkeepResult{{
				Key:              upkeepKey,
				State:            types.NotEligible,
				FailureReason:    UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED,
				GasUsed:          big.NewInt(27071),
				PerformData:      revertPerformData,
				FastGasWei:       big.NewInt(2000000000),
				LinkNative:       big.NewInt(4391095484380865),
				CheckBlockNumber: 8586947,
				CheckBlockHash:   [32]byte{230, 67, 97, 54, 73, 238, 133, 239, 200, 124, 171, 132, 40, 18, 124, 96, 102, 97, 232, 17, 96, 237, 173, 166, 112, 42, 146, 204, 46, 17, 67, 34},
				ExecuteGas:       5000000,
			}},
			hasHttpCalls:   true,
			callbackNeeded: true,
		},
		{
			name:       "skip - cooldown cache",
			input:      []types.UpkeepResult{upkeepResult},
			inCooldown: true,
			want:       []types.UpkeepResult{upkeepResult},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)
			client := new(evmClientMocks.Client)
			if tt.callbackNeeded {
				client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(tt.callbackResp, tt.callbackErr)
			}
			r.client = client

			mockHttpClient := mocks.NewHttpClient(t)
			if tt.hasHttpCalls {
				// mock the http client with Once() so the first call returns ETH-USD blob and the second call returns BTC-USD blob
				mockHttpClient.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(ethBlob)),
					},
					nil).Once()
				mockHttpClient.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(btcBlob)),
					},
					nil).Once()
			}
			r.hc = mockHttpClient

			if tt.inCooldown {
				r.mercury.cooldownCache.Set(upkeepId.String(), nil, DefaultCooldownExpiration)
			}

			// either set cache or mock getUpkeep
			if tt.mockGetUpkeep {
				mockReg := mocks.NewRegistry(t)
				r.registry = mockReg
				mockReg.On("GetUpkeep", mock.Anything, mock.Anything).Return(tt.upkeepInfo, tt.upkeepInfoErr)
			}
			if tt.upkeepCache {
				r.mercury.upkeepCache.Set(upkeepId.String(), tt.upkeepInfo, cache.DefaultExpiration)
			}

			got, err := r.mercuryLookup(context.Background(), tt.input)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), tt.name)
				assert.NotNil(t, err, tt.name)
			} else {
				assert.Equal(t, tt.want, got, tt.name)
			}
		})
	}
}

func TestEvmRegistry_decodeOffchainLookup(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *MercuryLookup
		wantErr error
	}{
		{
			name: "success",
			data: []byte{98, 232, 165, 13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 33, 20, 213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 102, 101, 101, 100, 73, 68, 83, 116, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 66, 84, 67, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 98, 108, 111, 99, 107, 78, 117, 109, 98, 101, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want: &MercuryLookup{
				feedLabel:  "feedIDStr",
				feeds:      []string{"ETH-USD-ARBITRUM-TESTNET", "BTC-USD-ARBITRUM-TESTNET"},
				queryLabel: "blockNumber",
				query:      big.NewInt(18945237),
				extraData:  []byte{48, 120, 48, 48},
			},
			wantErr: nil,
		},
		{
			name: "success - with extra data",
			data: []byte{98, 232, 165, 13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 33, 48, 241, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 102, 101, 101, 100, 73, 68, 83, 116, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 66, 84, 67, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 98, 108, 111, 99, 107, 78, 117, 109, 98, 101, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want: &MercuryLookup{
				feedLabel:  "feedIDStr",
				feeds:      []string{"ETH-USD-ARBITRUM-TESTNET", "BTC-USD-ARBITRUM-TESTNET"},
				queryLabel: "blockNumber",
				query:      big.NewInt(18952433),
				// this is the address of precompile contradt ArbSys(0x0000000000000000000000000000000000000064)
				extraData: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			},
			wantErr: nil,
		},
		{
			name:    "fail",
			data:    []byte{},
			want:    nil,
			wantErr: errors.New("unpack error: invalid data for unpacking"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)
			got, err := r.decodeMercuryLookup(tt.data)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), tt.name)
				assert.NotNil(t, err, tt.name)
			}
			assert.Equal(t, tt.want, got, tt.name)
		})
	}
}

func TestEvmRegistry_offchainLookupCallback(t *testing.T) {
	executeGas := uint32(100)
	gas := uint32(200000) + uint32(6500000) + uint32(300000) + executeGas
	from := common.HexToAddress("0x6cA639822c6C241Fa9A7A6b5032F6F7F1C513CAD")
	to := common.HexToAddress("0x79D8aDb571212b922089A48956c54A453D889dBe")
	bs := []byte{183, 114, 215, 10, 0, 0, 0, 0, 0, 0}
	values := [][]byte{bs}
	tests := []struct {
		name          string
		mercuryLookup *MercuryLookup
		values        [][]byte
		statusCode    int
		upkeepInfo    keeper_registry_wrapper2_0.UpkeepInfo
		opts          *bind.CallOpts

		callbackMsg  ethereum.CallMsg
		callbackResp []byte
		callbackErr  error

		upkeepNeeded bool
		performData  []byte
		wantErr      error
	}{
		{
			name: "success - empty extra data",
			mercuryLookup: &MercuryLookup{
				feedLabel:  "feedIDStr",
				feeds:      []string{"ETD-USD", "BTC-ETH"},
				queryLabel: "blockNumber",
				query:      big.NewInt(100),
				extraData:  []byte{48, 120, 48, 48},
			},
			values:     values,
			statusCode: http.StatusOK,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:         to,
				ExecuteGas:     executeGas,
				OffchainConfig: nil,
			},
			opts: &bind.CallOpts{
				BlockNumber: big.NewInt(999),
			},
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			upkeepNeeded: true,
			performData:  []byte{48, 120, 48, 48},
		},
		{
			name: "success - with extra data",
			mercuryLookup: &MercuryLookup{
				feedLabel:  "feedIDStr",
				feeds:      []string{"ETH-USD-ARBITRUM-TESTNET", "BTC-USD-ARBITRUM-TESTNET"},
				queryLabel: "blockNumber",
				query:      big.NewInt(18952430),
				// this is the address of precompile contradt ArbSys(0x0000000000000000000000000000000000000064)
				extraData: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			},
			values:     values,
			statusCode: http.StatusOK,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:         to,
				ExecuteGas:     executeGas,
				OffchainConfig: nil,
			},
			opts: &bind.CallOpts{
				BlockNumber: big.NewInt(18952430),
			},
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			upkeepNeeded: true,
			performData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
		},
		{
			name: "failure - bad response",
			mercuryLookup: &MercuryLookup{
				feedLabel:  "feedIDStr",
				feeds:      []string{"ETD-USD", "BTC-ETH"},
				queryLabel: "blockNumber",
				query:      big.NewInt(100),
				extraData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			values:     values,
			statusCode: http.StatusOK,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:         to,
				ExecuteGas:     executeGas,
				OffchainConfig: nil,
			},
			opts: &bind.CallOpts{
				BlockNumber: big.NewInt(999),
			},
			callbackResp: []byte{},

			wantErr: errors.New("callback output unpack error: abi: attempting to unmarshall an empty string while arguments are expected"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(evmClientMocks.Client)
			r := setupEVMRegistry(t)
			r.addr = from
			r.client = client
			payload, err := r.mercury.abi.Pack("mercuryCallback", values, tt.mercuryLookup.extraData)
			require.Nil(t, err)
			callbackMsg := ethereum.CallMsg{
				From: from,
				To:   &to,
				Gas:  uint64(gas),
				Data: payload,
			}
			client.On("CallContract", mock.Anything, callbackMsg, tt.opts.BlockNumber).Return(tt.callbackResp, tt.callbackErr)

			upkeepNeeded, performData, err := r.mercuryLookupCallback(context.Background(), tt.mercuryLookup, tt.values, tt.upkeepInfo, tt.opts)
			assert.Equal(t, tt.upkeepNeeded, upkeepNeeded, tt.name)
			assert.Equal(t, tt.performData, performData, tt.name)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), tt.name)
				assert.NotNil(t, err, tt.name)
			}
		})
	}
}

func TestEvmRegistry_setCachesOnAPIErr(t *testing.T) {
	tests := []struct {
		name     string
		upkeepId *big.Int
		rounds   int
	}{
		{
			name:     "success - 1,round",
			upkeepId: big.NewInt(100),
			rounds:   1,
		},
		{
			name:     "success - 2,rounds",
			upkeepId: big.NewInt(100),
			rounds:   2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)
			cacheKey := tt.upkeepId.String()
			for i := 0; i < tt.rounds; i++ {
				r.setCachesOnAPIErr(tt.upkeepId)
			}
			now := time.Now()

			val, exp, b := r.mercury.apiErrCache.GetWithExpiration(cacheKey)
			assert.True(t, b, "cache key found in apiErrCache")
			assert.NotNil(t, exp, "expiration found in apiErrCache")
			assert.GreaterOrEqual(t, exp, now.Add(DefaultApiErrExpiration-1*time.Minute), "expiration found in apiErrCache >= Default-1Minute")
			assert.Equal(t, tt.rounds, val, "err count correct")
			errCount := val.(int)

			val, exp, b = r.mercury.cooldownCache.GetWithExpiration(cacheKey)
			assert.True(t, b, "cache key found in cooldownCache")
			assert.NotNil(t, exp, "expiration found in cooldownCache")
			cooldown := time.Second * time.Duration(2^errCount)
			assert.GreaterOrEqual(t, exp, now.Add(cooldown/2), "expiration found in cooldownCache >= cooldown/2")
			assert.Equal(t, nil, val, "err count correct")
		})
	}
}
