package evm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
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
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_lookup_compatible_interface"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm/mocks"
)

// setups up an evm registry for tests.
func setupEVMRegistry(t *testing.T) *EvmRegistry {
	lggr := logger.TestLogger(t)
	addr := common.Address{}
	keeperRegistryABI, err := abi.JSON(strings.NewReader(i_keeper_registry_master_wrapper_2_1.IKeeperRegistryMasterABI))
	require.Nil(t, err, "need registry abi")
	mercuryCompatibleABI, err := abi.JSON(strings.NewReader(mercury_lookup_compatible_interface.MercuryLookupCompatibleInterfaceABI))
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
		headFunc: func(types.BlockKey) {},
		chLog:    make(chan logpoller.Log, 1000),
		mercury: &MercuryConfig{
			cred: &models.MercuryCredentials{
				URL:      "https://google.com",
				Username: "FakeClientID",
				Password: "FakeClientKey",
			},
			abi:            mercuryCompatibleABI,
			allowListCache: cache.New(DefaultAllowListExpiration, CleanupInterval),
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
		feedLabel:  "feedIDHex",
		feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
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

func TestEvmRegistry_mercuryLookup(t *testing.T) {
	setupRegistry := setupEVMRegistry(t)
	// load json response for testdata
	btcBlob, e := os.ReadFile("./testdata/btc-usd.json")
	assert.Nil(t, e)
	ethBlob, e := os.ReadFile("./testdata/eth-usd.json")
	assert.Nil(t, e)

	var block uint32 = 8586948
	upkeepId, ok := new(big.Int).SetString("520376062160720574742736856650455852347405918082346589375578334001045521721", 10)
	assert.True(t, ok, t.Name())

	// builds revert data with mock server query
	revertPerformData := setupRegistry.buildRevertBytesHelper()
	// prepare input upkeepResult
	upkeepResult := EVMAutomationUpkeepResult21{
		Block:            block,
		ID:               upkeepId,
		Eligible:         false,
		FailureReason:    UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
		GasUsed:          big.NewInt(27071),
		PerformData:      revertPerformData,
		FastGasWei:       big.NewInt(2000000000),
		LinkNative:       big.NewInt(4391095484380865),
		CheckBlockNumber: 8586947,
		CheckBlockHash:   [32]byte{230, 67, 97, 54, 73, 238, 133, 239, 200, 124, 171, 132, 40, 18, 124, 96, 102, 97, 232, 17, 96, 237, 173, 166, 112, 42, 146, 204, 46, 17, 67, 34},
		ExecuteGas:       5000000,
	}
	upkeepResultReasonMercury := EVMAutomationUpkeepResult21{
		Block:            block,
		ID:               upkeepId,
		Eligible:         false,
		FailureReason:    UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
		GasUsed:          big.NewInt(27071),
		PerformData:      revertPerformData,
		FastGasWei:       big.NewInt(2000000000),
		LinkNative:       big.NewInt(4391095484380865),
		CheckBlockNumber: 8586947,
		CheckBlockHash:   [32]byte{230, 67, 97, 54, 73, 238, 133, 239, 200, 124, 171, 132, 40, 18, 124, 96, 102, 97, 232, 17, 96, 237, 173, 166, 112, 42, 146, 204, 46, 17, 67, 34},
		ExecuteGas:       5000000,
	}
	//target := common.HexToAddress("0x79D8aDb571212b922089A48956c54A453D889dBe")
	callbackResp := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	upkeepNeededFalseResp, err := setupRegistry.abi.Methods["mercuryCallback"].Outputs.Pack(false, []byte{})
	assert.Nil(t, err, t.Name())

	// desired outcomes
	wantPerformData := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	wantUpkeepResult := EVMAutomationUpkeepResult21{
		Block:            8586948,
		ID:               upkeepId,
		Eligible:         true,
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
		{
			name:         "success - cached upkeep",
			input:        []EVMAutomationUpkeepResult21{upkeepResult},
			callbackResp: callbackResp,

			want:           []EVMAutomationUpkeepResult21{wantUpkeepResult},
			hasHttpCalls:   true,
			callbackNeeded: true,
		},
		{
			name:         "success - no cached upkeep",
			input:        []EVMAutomationUpkeepResult21{upkeepResult},
			callbackResp: callbackResp,

			want:           []EVMAutomationUpkeepResult21{wantUpkeepResult},
			hasHttpCalls:   true,
			callbackNeeded: true,
		},
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
		{
			name: "skip - revert data does not decode to mercury lookup, not surfacing errors",
			input: []EVMAutomationUpkeepResult21{
				{
					Block:         block,
					ID:            upkeepId,
					Eligible:      false,
					FailureReason: UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
					PerformData:   []byte{},
				},
			},

			want: []EVMAutomationUpkeepResult21{
				{
					Block:         block,
					ID:            upkeepId,
					Eligible:      false,
					FailureReason: UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
					PerformData:   []byte{},
				},
			},
		},
		{
			name:         "skip - error - no upkeep",
			input:        []EVMAutomationUpkeepResult21{upkeepResult},
			callbackResp: callbackResp,

			want:    []EVMAutomationUpkeepResult21{upkeepResultReasonMercury},
			wantErr: errors.New("ouch"),
		},
		{
			name:         "skip - upkeep not needed",
			input:        []EVMAutomationUpkeepResult21{upkeepResult},
			callbackResp: upkeepNeededFalseResp,

			want: []EVMAutomationUpkeepResult21{{
				Block:            block,
				ID:               upkeepId,
				Eligible:         false,
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
			name:  "skip - cooldown cache",
			input: []EVMAutomationUpkeepResult21{upkeepResult},
			want:  []EVMAutomationUpkeepResult21{upkeepResult},
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

func TestEvmRegistry_decodeMercuryLookup(t *testing.T) {
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
				feedLabel:  "feedIDHex",
				feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
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
				feedLabel:  "feedIDHex",
				feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				queryLabel: "blockNumber",
				query:      big.NewInt(18952433),
				// this is the address of precompile contract ArbSys(0x0000000000000000000000000000000000000064)
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

func TestEvmRegistry_SingleFeedRequest(t *testing.T) {
	upkeepId := big.NewInt(123456789)
	tests := []struct {
		name         string
		index        int
		ml           *MercuryLookup
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
			ml: &MercuryLookup{
				feedLabel:  "feedIDHex",
				feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				queryLabel: "blockNumber",
				query:      big.NewInt(123456),
				extraData:  nil,
			},
			mv:   job.MercuryV02,
			blob: "0xab2123dc00000012",
		},
		{
			name:  "success - retry for 404",
			index: 0,
			ml: &MercuryLookup{
				feedLabel:  "feedIDHex",
				feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				queryLabel: "blockNumber",
				query:      big.NewInt(123456),
				extraData:  nil,
			},
			mv:          job.MercuryV02,
			blob:        "0xab2123dcbabbad",
			retryNumber: 1,
			statusCode:  http.StatusNotFound,
		},
		{
			name:  "success - retry for 500",
			index: 0,
			ml: &MercuryLookup{
				feedLabel:  "feedIDHex",
				feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				queryLabel: "blockNumber",
				query:      big.NewInt(123456),
				extraData:  nil,
			},
			mv:          job.MercuryV02,
			blob:        "0xab2123dcbbabad",
			retryNumber: 2,
			statusCode:  http.StatusInternalServerError,
		},
		{
			name:  "failure - returns retryable",
			index: 0,
			ml: &MercuryLookup{
				feedLabel:  "feedIDHex",
				feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				queryLabel: "blockNumber",
				query:      big.NewInt(123456),
				extraData:  nil,
			},
			mv:           job.MercuryV02,
			blob:         "0xab2123dc",
			retryNumber:  TotalAttempt,
			statusCode:   http.StatusNotFound,
			retryable:    true,
			errorMessage: "All attempts fail:\n#1: retry\n#2: retry\n#3: retry",
		},
		{
			name:  "failure - returns not retryable",
			index: 0,
			ml: &MercuryLookup{
				feedLabel:  "feedIDHex",
				feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				queryLabel: "blockNumber",
				query:      big.NewInt(123456),
				extraData:  nil,
			},
			mv:           job.MercuryV02,
			blob:         "0xab2123dc",
			statusCode:   http.StatusBadGateway,
			retryable:    false,
			errorMessage: "All attempts fail:\n#1: MercuryLookup upkeep 123456789 block 123456 received status code 502 for feed 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
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

func TestEvmRegistry_MercuryCallback(t *testing.T) {
	//executeGas := uint32(100)
	from := common.HexToAddress("0x6cA639822c6C241Fa9A7A6b5032F6F7F1C513CAD")
	to := common.HexToAddress("0x79D8aDb571212b922089A48956c54A453D889dBe")
	bs := []byte{183, 114, 215, 10, 0, 0, 0, 0, 0, 0}
	values := [][]byte{bs}
	tests := []struct {
		name          string
		mercuryLookup *MercuryLookup
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
		{
			name: "success - empty extra data",
			mercuryLookup: &MercuryLookup{
				feedLabel:  "feedIDHex",
				feeds:      []string{"ETD-USD", "BTC-ETH"},
				queryLabel: "blockNumber",
				query:      big.NewInt(100),
				extraData:  []byte{48, 120, 48, 48},
			},
			values:       values,
			statusCode:   http.StatusOK,
			upkeepId:     big.NewInt(123456789),
			blockNumber:  999,
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			upkeepNeeded: true,
			performData:  []byte{48, 120, 48, 48},
		},
		{
			name: "success - with extra data",
			mercuryLookup: &MercuryLookup{
				feedLabel:  "feedIDHex",
				feeds:      []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				queryLabel: "blockNumber",
				query:      big.NewInt(18952430),
				// this is the address of precompile contract ArbSys(0x0000000000000000000000000000000000000064)
				extraData: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			},
			values:       values,
			statusCode:   http.StatusOK,
			upkeepId:     big.NewInt(123456789),
			blockNumber:  999,
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			upkeepNeeded: true,
			performData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
		},
		{
			name: "failure - bad response",
			mercuryLookup: &MercuryLookup{
				feedLabel:  "feedIDHex",
				feeds:      []string{"ETD-USD", "BTC-ETH"},
				queryLabel: "blockNumber",
				query:      big.NewInt(100),
				extraData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			values:       values,
			statusCode:   http.StatusOK,
			upkeepId:     big.NewInt(123456789),
			blockNumber:  999,
			callbackResp: []byte{},
			wantErr:      errors.New("callback output unpack error: abi: attempting to unmarshall an empty string while arguments are expected"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(evmClientMocks.Client)
			r := setupEVMRegistry(t)
			r.addr = from
			r.client = client
			payload, err := r.abi.Pack("mercuryCallback", tt.upkeepId, values, tt.mercuryLookup.extraData)
			require.Nil(t, err)
			callbackMsg := ethereum.CallMsg{
				To:   &to,
				Data: payload,
			}
			// if args don't match, just use mock.Anything
			client.On("CallContract", mock.Anything, callbackMsg, tt.blockNumber).Return(tt.callbackResp, tt.callbackErr)

			upkeepNeeded, performData, _, _, err := r.mercuryCallback(context.Background(), tt.upkeepId, tt.values, tt.mercuryLookup.extraData, tt.blockNumber)
			assert.Equal(t, tt.upkeepNeeded, upkeepNeeded, tt.name)
			assert.Equal(t, tt.performData, performData, tt.name)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), tt.name)
				assert.NotNil(t, err, tt.name)
			}
		})
	}
}
