package evm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/mocks"

	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/feed_lookup_compatible_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
)

// setups up an evm registry for tests.
func setupEVMRegistry(t *testing.T) *EvmRegistry {
	lggr := logger.TestLogger(t)
	addr := common.HexToAddress("0x6cA639822c6C241Fa9A7A6b5032F6F7F1C513CAD")
	keeperRegistryABI, err := abi.JSON(strings.NewReader(i_keeper_registry_master_wrapper_2_1.IKeeperRegistryMasterABI))
	require.Nil(t, err, "need registry abi")
	feedLookupCompatibleABI, err := abi.JSON(strings.NewReader(feed_lookup_compatible_interface.FeedLookupCompatibleInterfaceABI))
	require.Nil(t, err, "need mercury abi")
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
			chHead: make(chan ocr2keepers.BlockKey, 1),
		},
		lggr:     lggr,
		poller:   logPoller,
		addr:     addr,
		client:   client,
		txHashes: make(map[string]bool),
		registry: mockRegistry,
		abi:      keeperRegistryABI,
		active:   make(map[string]activeUpkeep),
		packer:   &evmRegistryPackerV2_1{abi: keeperRegistryABI},
		headFunc: func(ocr2keepers.BlockKey) {},
		chLog:    make(chan logpoller.Log, 1000),
		mercury: &MercuryConfig{
			cred: &models.MercuryCredentials{
				URL:      "https://google.com",
				Username: "FakeClientID",
				Password: "FakeClientKey",
			},
			abi:            feedLookupCompatibleABI,
			allowListCache: cache.New(DefaultAllowListExpiration, CleanupInterval),
		},
		hc: mockHttpClient,
	}
	return r
}

// helper for mocking the http requests
//func (r *EvmRegistry) buildRevertBytesHelper() []byte {
//	mercuryErr := r.mercury.abi.Errors["FeedLookup"]
//	mercuryLookupSelector := [4]byte{0x62, 0xe8, 0xa5, 0x0d}
//	ml := FeedLookup{
//		feedParamKey:  "feedIDHex",
//		feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
//		timeParamKey: "blockNumber",
//		time:      big.NewInt(8586948),
//		extraData:  []byte{},
//	}
//	// check if Pack does not add selector for me
//	pack, err := mercuryErr.Inputs.Pack(ml.feedParamKey, ml.feeds, ml.timeParamKey, ml.time, ml.extraData)
//	if err != nil {
//		log.Fatal("failed to build revert")
//	}
//	var payload []byte
//	payload = append(payload, mercuryLookupSelector[:]...)
//	payload = append(payload, pack...)
//	return payload
//}

func TestEvmRegistry_mercuryLookup(t *testing.T) {
	//setupRegistry := setupEVMRegistry(t)
	// load json response for testdata
	btcBlob, e := os.ReadFile("./testdata/btc-usd.json")
	assert.Nil(t, e)
	ethBlob, e := os.ReadFile("./testdata/eth-usd.json")
	assert.Nil(t, e)

	var block uint32 = 8586948
	upkeepId, ok := new(big.Int).SetString("520376062160720574742736856650455852347405918082346589375578334001045521721", 10)
	assert.True(t, ok, t.Name())

	// builds revert data with mock server time
	//revertPerformData := setupRegistry.buildRevertBytesHelper()
	// prepare input upkeepResult
	//upkeepResult := EVMAutomationUpkeepResult21{
	//	Block:            block,
	//	ID:               upkeepId,
	//	Eligible:         false,
	//	FailureReason:    UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
	//	GasUsed:          big.NewInt(27071),
	//	PerformData:      revertPerformData,
	//	FastGasWei:       big.NewInt(2000000000),
	//	LinkNative:       big.NewInt(4391095484380865),
	//	CheckBlockNumber: 8586947,
	//	CheckBlockHash:   [32]byte{230, 67, 97, 54, 73, 238, 133, 239, 200, 124, 171, 132, 40, 18, 124, 96, 102, 97, 232, 17, 96, 237, 173, 166, 112, 42, 146, 204, 46, 17, 67, 34},
	//	ExecuteGas:       5000000,
	//}
	//upkeepResultReasonMercury := EVMAutomationUpkeepResult21{
	//	Block:            block,
	//	ID:               upkeepId,
	//	Eligible:         false,
	//	FailureReason:    UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
	//	GasUsed:          big.NewInt(27071),
	//	PerformData:      revertPerformData,
	//	FastGasWei:       big.NewInt(2000000000),
	//	LinkNative:       big.NewInt(4391095484380865),
	//	CheckBlockNumber: 8586947,
	//	CheckBlockHash:   [32]byte{230, 67, 97, 54, 73, 238, 133, 239, 200, 124, 171, 132, 40, 18, 124, 96, 102, 97, 232, 17, 96, 237, 173, 166, 112, 42, 146, 204, 46, 17, 67, 34},
	//	ExecuteGas:       5000000,
	//}
	//target := common.HexToAddress("0x79D8aDb571212b922089A48956c54A453D889dBe")
	//callbackResp := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	//upkeepNeededFalseResp, err := setupRegistry.abi.Methods["checkCallback"].Outputs.Pack(false, []byte{})
	//assert.Nil(t, err, t.Name())

	// desired outcomes
	//wantPerformData := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	//wantUpkeepResult := EVMAutomationUpkeepResult21{
	//	Block:            8586948,
	//	ID:               upkeepId,
	//	Eligible:         true,
	//	FailureReason:    UPKEEP_FAILURE_REASON_NONE,
	//	GasUsed:          big.NewInt(27071),
	//	PerformData:      wantPerformData,
	//	FastGasWei:       big.NewInt(2000000000),
	//	LinkNative:       big.NewInt(4391095484380865),
	//	CheckBlockNumber: 8586947,
	//	CheckBlockHash:   [32]byte{230, 67, 97, 54, 73, 238, 133, 239, 200, 124, 171, 132, 40, 18, 124, 96, 102, 97, 232, 17, 96, 237, 173, 166, 112, 42, 146, 204, 46, 17, 67, 34},
	//	ExecuteGas:       5000000,
	//}
	tests := []struct {
		name           string
		input          []EVMAutomationUpkeepResult21
		callbackResp   []byte
		callbackErr    error
		want           []EVMAutomationUpkeepResult21
		wantErr        error
		hasHttpCalls   bool
		callbackNeeded bool
		cachedAdminCfg bool
		upkeepId       *big.Int
	}{
		//{
		//	name:         "success - cached upkeep",
		//	input:        []EVMAutomationUpkeepResult21{upkeepResult},
		//	callbackResp: callbackResp,
		//
		//	want:           []EVMAutomationUpkeepResult21{wantUpkeepResult},
		//	hasHttpCalls:   true,
		//	callbackNeeded: true,
		//},
		//{
		//	name:         "success - no cached upkeep",
		//	input:        []EVMAutomationUpkeepResult21{upkeepResult},
		//	callbackResp: callbackResp,
		//
		//	want:           []EVMAutomationUpkeepResult21{wantUpkeepResult},
		//	hasHttpCalls:   true,
		//	callbackNeeded: true,
		//},
		{
			name: "skip - failure reason",
			input: []EVMAutomationUpkeepResult21{
				{
					Block:         block,
					ID:            upkeepId,
					Eligible:      false,
					FailureReason: UPKEEP_FAILURE_REASON_INSUFFICIENT_BALANCE,
					PerformData:   []byte{},
				},
			},

			want: []EVMAutomationUpkeepResult21{
				{
					Block:         block,
					ID:            upkeepId,
					Eligible:      false,
					FailureReason: UPKEEP_FAILURE_REASON_INSUFFICIENT_BALANCE,
					PerformData:   []byte{},
				},
			},
		},
		//{
		//	name: "skip - revert data does not decode to mercury lookup, not surfacing errors",
		//	input: []EVMAutomationUpkeepResult21{
		//		{
		//			Block:         block,
		//			ID:            upkeepId,
		//			Eligible:      false,
		//			FailureReason: UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
		//			PerformData:   []byte{},
		//		},
		//	},
		//
		//	want: []EVMAutomationUpkeepResult21{
		//		{
		//			Block:         block,
		//			ID:            upkeepId,
		//			Eligible:      false,
		//			FailureReason: UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
		//			PerformData:   []byte{},
		//		},
		//	},
		//},
		//{
		//	name:         "skip - error - no upkeep",
		//	input:        []EVMAutomationUpkeepResult21{upkeepResult},
		//	callbackResp: callbackResp,
		//
		//	want:    []EVMAutomationUpkeepResult21{upkeepResultReasonMercury},
		//	wantErr: errors.New("ouch"),
		//},
		//{
		//	name:         "skip - upkeep not needed",
		//	input:        []EVMAutomationUpkeepResult21{upkeepResult},
		//	callbackResp: upkeepNeededFalseResp,
		//
		//	want: []EVMAutomationUpkeepResult21{{
		//		Block:            block,
		//		ID:               upkeepId,
		//		Eligible:         false,
		//		FailureReason:    UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED,
		//		GasUsed:          big.NewInt(27071),
		//		PerformData:      revertPerformData,
		//		FastGasWei:       big.NewInt(2000000000),
		//		LinkNative:       big.NewInt(4391095484380865),
		//		CheckBlockNumber: 8586947,
		//		CheckBlockHash:   [32]byte{230, 67, 97, 54, 73, 238, 133, 239, 200, 124, 171, 132, 40, 18, 124, 96, 102, 97, 232, 17, 96, 237, 173, 166, 112, 42, 146, 204, 46, 17, 67, 34},
		//		ExecuteGas:       5000000,
		//	}},
		//	hasHttpCalls:   true,
		//	callbackNeeded: true,
		//},
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

			got, err := r.feedLookup(context.Background(), tt.input)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), tt.name)
				assert.NotNil(t, err, tt.name)
			} else {
				assert.Equal(t, tt.want, got, tt.name)
			}
		})
	}
}

func TestEvmRegistry_DecodeFeedLookup(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected *FeedLookup
		err      error
	}{
		{
			name: "success - decode to feed lookup",
			data: []byte{125, 221, 147, 62, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 138, 215, 253, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 102, 101, 101, 100, 73, 68, 72, 101, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 66, 48, 120, 52, 53, 53, 52, 52, 56, 50, 100, 53, 53, 53, 51, 52, 52, 50, 100, 52, 49, 53, 50, 52, 50, 52, 57, 53, 52, 53, 50, 53, 53, 52, 100, 50, 100, 53, 52, 52, 53, 53, 51, 53, 52, 52, 101, 52, 53, 53, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 66, 48, 120, 52, 50, 53, 52, 52, 51, 50, 100, 53, 53, 53, 51, 52, 52, 50, 100, 52, 49, 53, 50, 52, 50, 52, 57, 53, 52, 53, 50, 53, 53, 52, 100, 50, 100, 53, 52, 52, 53, 53, 51, 53, 52, 52, 101, 52, 53, 53, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 98, 108, 111, 99, 107, 78, 117, 109, 98, 101, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: &FeedLookup{
				feedParamKey: "feedIDHex",
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: "blockNumber",
				time:         big.NewInt(25876477),
				extraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			},
		},
		{
			name: "failure - unpack error",
			data: []byte{1, 2, 3, 4},
			err:  errors.New("unpack error: invalid data for unpacking"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)
			fl, err := r.decodeFeedLookup(tt.data)
			assert.Equal(t, tt.expected, fl)
			if tt.err != nil {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestEvmRegistry_AllowedToUseMercury(t *testing.T) {
	upkeepId := big.NewInt(123456789)
	tests := []struct {
		name         string
		cached       bool
		allowed      bool
		errorMessage string
	}{
		{
			name:    "success - allowed via cache",
			cached:  true,
			allowed: true,
		},
		{
			name:    "success - allowed via fetching admin offchain config",
			cached:  false,
			allowed: true,
		},
		{
			name:    "success - not allowed via cache",
			cached:  true,
			allowed: false,
		},
		{
			name:    "success - not allowed via fetching admin offchain config",
			cached:  false,
			allowed: false,
		},
		{
			name:         "failure - cannot unmarshal admin offchain config",
			cached:       false,
			errorMessage: "failed to unmarshal admin offchain config for upkeep ID 123456789: invalid character '\\x00' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)

			if tt.errorMessage != "" {
				mockRegistry := mocks.NewRegistry(t)
				mockRegistry.On("GetUpkeepAdminOffchainConfig", mock.Anything, upkeepId).Return([]byte{0, 1}, nil)
				r.registry = mockRegistry
			} else {
				if tt.cached {
					r.mercury.allowListCache.Set(upkeepId.String(), tt.allowed, cache.DefaultExpiration)
				} else {
					mockRegistry := mocks.NewRegistry(t)
					cfg := AdminOffchainConfig{MercuryEnabled: tt.allowed}
					b, err := json.Marshal(cfg)
					assert.Nil(t, err)
					mockRegistry.On("GetUpkeepAdminOffchainConfig", mock.Anything, upkeepId).Return(b, nil)
					r.registry = mockRegistry
				}
			}

			allowed, err := r.allowedToUseMercury(nil, upkeepId)
			if tt.errorMessage != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.allowed, allowed)
			}
		})
	}
}

func TestEvmRegistry_DoMercuryRequest(t *testing.T) {
	upkeepId := big.NewInt(0)
	upkeepId.SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)

	tests := []struct {
		name               string
		ml                 *FeedLookup
		mockHttpStatusCode int
		mockChainlinkBlobs []string
		expectedValues     [][]byte
		expectedRetryable  bool
		expectedError      error
	}{
		{
			name: "success",
			ml: &FeedLookup{
				feedParamKey: "feedIDHex",
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: "blockNumber",
				time:         big.NewInt(25880526),
				extraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			},
			mockHttpStatusCode: http.StatusOK,
			mockChainlinkBlobs: []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:     [][]byte{{0, 6, 109, 252, 209, 237, 45, 149, 177, 140, 148, 141, 188, 91, 214, 76, 104, 122, 254, 147, 228, 202, 125, 102, 61, 222, 193, 76, 32, 9, 10, 216, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 20, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 128, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 137, 28, 152, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 40, 154, 216, 211, 103, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 40, 154, 207, 11, 56, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 40, 155, 61, 164, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 138, 231, 206, 116, 217, 250, 37, 42, 137, 131, 151, 110, 171, 96, 13, 199, 89, 12, 119, 141, 4, 129, 52, 48, 132, 27, 198, 231, 101, 195, 76, 216, 26, 22, 141, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 138, 231, 203, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 137, 28, 152, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 96, 65, 43, 148, 229, 37, 202, 108, 237, 201, 245, 68, 253, 134, 247, 118, 6, 213, 47, 231, 49, 165, 208, 105, 219, 232, 54, 168, 191, 192, 251, 140, 145, 25, 99, 176, 174, 122, 20, 151, 31, 59, 70, 33, 191, 251, 128, 46, 240, 96, 83, 146, 185, 166, 200, 156, 127, 171, 29, 248, 99, 58, 90, 222, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 69, 0, 194, 245, 33, 248, 63, 186, 94, 252, 43, 243, 239, 250, 174, 221, 228, 61, 10, 74, 223, 247, 133, 193, 33, 59, 113, 42, 58, 237, 13, 129, 87, 100, 42, 132, 50, 77, 176, 207, 150, 149, 235, 210, 119, 8, 212, 96, 142, 176, 51, 126, 13, 216, 123, 14, 67, 240, 250, 112, 199, 0, 217, 17}},
			expectedRetryable:  false,
			expectedError:      nil,
		},
		{
			name: "failure - retryable",
			ml: &FeedLookup{
				feedParamKey: "feedIDHex",
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: "blockNumber",
				time:         big.NewInt(25880526),
				extraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			},
			mockHttpStatusCode: http.StatusInternalServerError,
			mockChainlinkBlobs: []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:     [][]byte{nil},
			expectedRetryable:  true,
			expectedError:      errors.New("All attempts fail:\n#1: 500\n#2: 500\n#3: 500"),
		},
		{
			name: "failure - not retryable",
			ml: &FeedLookup{
				feedParamKey: "feedIDHex",
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: "blockNumber",
				time:         big.NewInt(25880526),
				extraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			},
			mockHttpStatusCode: http.StatusBadGateway,
			mockChainlinkBlobs: []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:     [][]byte{nil},
			expectedRetryable:  false,
			expectedError:      errors.New("All attempts fail:\n#1: FeedLookup upkeep 88786950015966611018675766524283132478093844178961698330929478019253453382042 block 25880526 received status code 502 for feed 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)
			hc := mocks.NewHttpClient(t)

			for _, blob := range tt.mockChainlinkBlobs {
				mr := MercuryResponse{ChainlinkBlob: blob}
				b, err := json.Marshal(mr)
				assert.Nil(t, err)

				resp := &http.Response{
					StatusCode: tt.mockHttpStatusCode,
					Body:       io.NopCloser(bytes.NewReader(b)),
				}
				if tt.expectedError != nil && tt.expectedRetryable {
					hc.On("Do", mock.Anything).Return(resp, nil).Times(TotalAttempt)
				} else {
					hc.On("Do", mock.Anything).Return(resp, nil).Once()
				}
			}
			r.hc = hc

			values, retryable, reqErr := r.doMercuryRequest(context.Background(), tt.ml, upkeepId)
			assert.Equal(t, tt.expectedValues, values)
			assert.Equal(t, tt.expectedRetryable, retryable)
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError.Error(), reqErr.Error())
			}
		})
	}
}

func TestEvmRegistry_SingleFeedRequest(t *testing.T) {
	upkeepId := big.NewInt(123456789)
	tests := []struct {
		name         string
		index        int
		ml           *FeedLookup
		mv           job.MercuryVersion
		blob         string
		statusCode   int
		retryNumber  int
		retryable    bool
		errorMessage string
	}{
		{
			name:  "success - mercury responds in the first try",
			index: 0,
			ml: &FeedLookup{
				feedParamKey: "feedIDHex",
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: "blockNumber",
				time:         big.NewInt(123456),
				extraData:    nil,
			},
			mv:   job.MercuryV02,
			blob: "0xab2123dc00000012",
		},
		{
			name:  "success - retry for 404",
			index: 0,
			ml: &FeedLookup{
				feedParamKey: "feedIDHex",
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: "blockNumber",
				time:         big.NewInt(123456),
				extraData:    nil,
			},
			mv:          job.MercuryV02,
			blob:        "0xab2123dcbabbad",
			retryNumber: 1,
			statusCode:  http.StatusNotFound,
		},
		{
			name:  "success - retry for 500",
			index: 0,
			ml: &FeedLookup{
				feedParamKey: "feedIDHex",
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: "blockNumber",
				time:         big.NewInt(123456),
				extraData:    nil,
			},
			mv:          job.MercuryV02,
			blob:        "0xab2123dcbbabad",
			retryNumber: 2,
			statusCode:  http.StatusInternalServerError,
		},
		{
			name:  "failure - returns retryable",
			index: 0,
			ml: &FeedLookup{
				feedParamKey: "feedIDHex",
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: "blockNumber",
				time:         big.NewInt(123456),
				extraData:    nil,
			},
			mv:           job.MercuryV02,
			blob:         "0xab2123dc",
			retryNumber:  TotalAttempt,
			statusCode:   http.StatusNotFound,
			retryable:    true,
			errorMessage: "All attempts fail:\n#1: 404\n#2: 404\n#3: 404",
		},
		{
			name:  "failure - returns not retryable",
			index: 0,
			ml: &FeedLookup{
				feedParamKey: "feedIDHex",
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: "blockNumber",
				time:         big.NewInt(123456),
				extraData:    nil,
			},
			mv:           job.MercuryV02,
			blob:         "0xab2123dc",
			statusCode:   http.StatusBadGateway,
			retryable:    false,
			errorMessage: "All attempts fail:\n#1: FeedLookup upkeep 123456789 block 123456 received status code 502 for feed 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)
			hc := mocks.NewHttpClient(t)

			mr := MercuryResponse{ChainlinkBlob: tt.blob}
			b, err := json.Marshal(mr)
			assert.Nil(t, err)

			if tt.retryNumber == 0 {
				if tt.errorMessage != "" {
					resp := &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(bytes.NewReader(b)),
					}
					hc.On("Do", mock.Anything).Return(resp, nil).Once()
				} else {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(b)),
					}
					hc.On("Do", mock.Anything).Return(resp, nil).Once()
				}
			} else if tt.retryNumber > 0 && tt.retryNumber < TotalAttempt {
				retryResp := &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(bytes.NewReader(b)),
				}
				hc.On("Do", mock.Anything).Return(retryResp, nil).Times(tt.retryNumber)

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(b)),
				}
				hc.On("Do", mock.Anything).Return(resp, nil).Once()
			} else {
				resp := &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(bytes.NewReader(b)),
				}
				hc.On("Do", mock.Anything).Return(resp, nil).Times(tt.retryNumber)
			}
			r.hc = hc

			ch := make(chan MercuryBytes, 1)
			r.singleFeedRequest(context.Background(), ch, upkeepId, tt.index, tt.ml, tt.mv)

			m := <-ch
			assert.Equal(t, tt.index, m.Index)
			assert.Equal(t, tt.retryable, m.Retryable)
			if tt.retryNumber >= TotalAttempt || tt.errorMessage != "" {
				assert.Equal(t, tt.errorMessage, m.Error.Error())
				assert.Nil(t, m.Bytes)
			} else {
				blobBytes, err := hexutil.Decode(tt.blob)
				assert.Nil(t, err)
				assert.Nil(t, m.Error)
				assert.Equal(t, blobBytes, m.Bytes)
			}
		})
	}
}

func TestEvmRegistry_CheckCallback(t *testing.T) {
	bs := []byte{183, 114, 215, 10, 0, 0, 0, 0, 0, 0}
	values := [][]byte{bs}
	tests := []struct {
		name          string
		mercuryLookup *FeedLookup
		values        [][]byte
		statusCode    int
		upkeepId      *big.Int
		blockNumber   uint32

		callbackMsg  ethereum.CallMsg
		callbackResp []byte
		callbackErr  error

		upkeepNeeded bool
		performData  []byte
		wantErr      error
	}{
		//{
		//	name: "success - empty extra data",
		//	mercuryLookup: &FeedLookup{
		//		feedParamKey: "feedIDHex",
		//		feeds:        []string{"ETD-USD", "BTC-ETH"},
		//		timeParamKey: "blockNumber",
		//		time:         big.NewInt(100),
		//		extraData:    []byte{48, 120, 48, 48},
		//	},
		//	values:       values,
		//	statusCode:   http.StatusOK,
		//	upkeepId:     big.NewInt(123456789),
		//	blockNumber:  999,
		//	callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		//	upkeepNeeded: true,
		//	performData:  []byte{48, 120, 48, 48},
		//},
		//{
		//	name: "success - with extra data",
		//	mercuryLookup: &FeedLookup{
		//		feedParamKey:  "feedIDHex",
		//		feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
		//		timeParamKey: "blockNumber",
		//		time:      big.NewInt(18952430),
		//		// this is the address of precompile contract ArbSys(0x0000000000000000000000000000000000000064)
		//		extraData: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
		//	},
		//	values:       values,
		//	statusCode:   http.StatusOK,
		//	upkeepId:     big.NewInt(123456789),
		//	blockNumber:  999,
		//	callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		//	upkeepNeeded: true,
		//	performData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
		//},
		//{
		//	name: "failure - bad response",
		//	mercuryLookup: &FeedLookup{
		//		feedParamKey:  "feedIDHex",
		//		feeds:      []string{"ETD-USD", "BTC-ETH"},
		//		timeParamKey: "blockNumber",
		//		time:      big.NewInt(100),
		//		extraData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		//	},
		//	values:       values,
		//	statusCode:   http.StatusOK,
		//	upkeepId:     big.NewInt(123456789),
		//	blockNumber:  999,
		//	callbackResp: []byte{},
		//	wantErr:      errors.New("callback output unpack error: abi: attempting to unmarshall an empty string while arguments are expected"),
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(evmClientMocks.Client)
			r := setupEVMRegistry(t)
			payload, err := r.abi.Pack("checkCallback", tt.upkeepId, values, tt.mercuryLookup.extraData)
			require.Nil(t, err)
			args := map[string]interface{}{
				"to":   r.addr.Hex(),
				"data": hexutil.Bytes(payload),
			}
			// if args don't match, just use mock.Anything
			client.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", args, hexutil.EncodeUint64(uint64(tt.blockNumber))).
				Run(func(args mock.Arguments) {
					//b := args.Get(1).(*hexutil.Bytes)
					//b = &hexutil.Bytes{0, 1}
				}).Once()
			r.client = client

			_, err = r.checkCallback(context.Background(), tt.upkeepId, tt.values, tt.mercuryLookup.extraData, tt.blockNumber)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), tt.name)
				assert.NotNil(t, err, tt.name)
			}
		})
	}
}
