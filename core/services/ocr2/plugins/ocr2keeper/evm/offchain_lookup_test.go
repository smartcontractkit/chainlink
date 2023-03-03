package evm

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
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

	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/upkeep_apifetch_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2keeper/evm/mocks"
)

// setups up an evm registry for tests.
func setupEVMRegistry(t *testing.T) *EvmRegistry {
	lggr := logger.TestLogger(t)
	addr := common.Address{}
	keeperRegistryABI, err := abi.JSON(strings.NewReader(keeper_registry_wrapper2_0.KeeperRegistryABI))
	require.Nil(t, err, "need registry abi")
	apiFetchABI, err := abi.JSON(strings.NewReader(upkeep_apifetch_wrapper.UpkeepAPIFetchABI))
	require.Nil(t, err, "need API Fetch abi")
	upkeepInfoCache, cooldownCache, apiErrCache := setupCaches(DefaultUpkeepExpiration, DefaultCooldownExpiration, DefaultApiErrExpiration, CleanupInterval)
	var headTracker httypes.HeadTracker
	var headBroadcaster httypes.HeadBroadcaster
	var logPoller logpoller.LogPoller
	mockRegistry := mocks.KeeperRegistryInterface{}
	client := new(evmmocks.Client)

	r := &EvmRegistry{
		HeadProvider: HeadProvider{
			ht:     headTracker,
			hb:     headBroadcaster,
			chHead: make(chan types.BlockKey, 1),
		},
		lggr:          lggr,
		poller:        logPoller,
		addr:          addr,
		client:        client,
		txHashes:      make(map[string]bool),
		registry:      &mockRegistry,
		abi:           keeperRegistryABI,
		apiFetchABI:   apiFetchABI,
		active:        make(map[string]activeUpkeep),
		packer:        &evmRegistryPackerV2_0{abi: keeperRegistryABI},
		headFunc:      func(types.BlockKey) {},
		chLog:         make(chan logpoller.Log, 1000),
		upkeepCache:   upkeepInfoCache,
		cooldownCache: cooldownCache,
		apiErrCache:   apiErrCache,
	}
	return r
}

// TODO remove later just using this to test setups
func TestRevertData(t *testing.T) {
	r := setupEVMRegistry(t)
	o := OffchainLookup{
		url:              "https://pokeapi.co/api/v2/pokemon/1",
		extraData:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		fields:           []string{"id", "name"},
		callbackFunction: [4]byte{183, 114, 215, 10},
	}
	revertPerformData := []byte{218, 139, 50, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 96, 183, 114, 215, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 35, 104, 116, 116, 112, 115, 58, 47, 47, 112, 111, 107, 101, 97, 112, 105, 46, 99, 111, 47, 97, 112, 105, 47, 118, 50, 47, 112, 111, 107, 101, 109, 111, 110, 47, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 105, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 110, 97, 109, 101, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	e := r.apiFetchABI.Errors["ChainlinkAPIFetch"]
	// error bytes4 selector
	apiFetch4Byte := [4]byte{0xda, 0x8b, 0x32, 0x14}

	pack, err := e.Inputs.Pack(o.url, o.extraData, o.fields, o.callbackFunction)
	assert.Nil(t, err, t.Name())
	var payload []byte
	payload = append(payload, apiFetch4Byte[:]...)
	payload = append(payload, pack...)

	assert.Equal(t, revertPerformData, payload, t.Name())

	callbackResp := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	wantPerformData := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	got, err := r.apiFetchABI.Methods["callback"].Outputs.Pack(true, wantPerformData)
	assert.Nil(t, err, t.Name())

	assert.Equal(t, callbackResp, got, t.Name())

}

// helper for mocking the http requests
func buildRevertBytesHelper(r *EvmRegistry, baseURL string) []byte {
	apiFetchError := r.apiFetchABI.Errors["ChainlinkAPIFetch"]
	// "ChainlinkAPIFetch" error bytes4 selector
	apiFetch4Byte := [4]byte{0xda, 0x8b, 0x32, 0x14}
	url := fmt.Sprintf("%s/api/v2/pokemon/1", baseURL)
	offchainLookup := OffchainLookup{
		url:              url,
		extraData:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		fields:           []string{"id", "name"},
		callbackFunction: [4]byte{183, 114, 215, 10},
	}
	pack, err := apiFetchError.Inputs.Pack(offchainLookup.url, offchainLookup.extraData, offchainLookup.fields, offchainLookup.callbackFunction)
	if err != nil {
		log.Fatal("failed to build revert")
	}
	var payload []byte
	payload = append(payload, apiFetch4Byte[:]...)
	payload = append(payload, pack...)
	return payload
}

func TestEvmRegistry_offchainLookup(t *testing.T) {
	setupRegistry := setupEVMRegistry(t)
	// load json response for testing
	content, e := os.ReadFile("poke_api.json")
	assert.Nil(t, e)
	// mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write(content)
	}))
	defer server.Close()
	upkeepKey := chain.UpkeepKey("8586948|104693970376404490964326661740762530070325727241549215715800663219260698550627")
	_, upkeepId, err := blockAndIdFromKey(upkeepKey)
	assert.Nil(t, err, t.Name())
	// builds revert data with mock server url
	revertPerformData := buildRevertBytesHelper(setupRegistry, server.URL)
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
	target := common.HexToAddress("0x79D8aDb571212b922089A48956c54A453D889dBe")
	callbackResp := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	upkeepNeededFalseResp, err := setupRegistry.apiFetchABI.Methods["callback"].Outputs.Pack(false, []byte{})
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
		upkeepInfo    keeper_registry_wrapper2_0.UpkeepInfo
		upkeepInfoErr error

		want    []types.UpkeepResult
		wantErr error
	}{
		{
			name:         "success - cached upkeep",
			input:        []types.UpkeepResult{upkeepResult},
			callbackResp: callbackResp,
			upkeepCache:  true,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:     target,
				ExecuteGas: 5000000,
			},

			want: []types.UpkeepResult{wantUpkeepResult},
		},
		{
			name:         "success - no cached upkeep",
			input:        []types.UpkeepResult{upkeepResult},
			callbackResp: callbackResp,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:     target,
				ExecuteGas: 5000000,
			},

			want: []types.UpkeepResult{wantUpkeepResult},
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
					FailureReason: UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
					PerformData:   []byte{},
				},
			},
		},
		{
			name:          "skip - error - no upkeep",
			input:         []types.UpkeepResult{upkeepResult},
			callbackResp:  callbackResp,
			upkeepInfoErr: errors.New("ouch"),

			want: []types.UpkeepResult{upkeepResult},
		},
		{
			name:         "skip - upkeep not needed",
			input:        []types.UpkeepResult{upkeepResult},
			callbackResp: upkeepNeededFalseResp,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:     target,
				ExecuteGas: 5000000,
			},

			want: []types.UpkeepResult{upkeepResult},
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
			client := new(evmmocks.Client)
			r.client = client
			client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(tt.callbackResp, tt.callbackErr)

			if tt.inCooldown {
				r.cooldownCache.Set(upkeepId.String(), nil, DefaultCooldownExpiration)
			}

			// either set cache or mock registry return
			if tt.upkeepCache {
				r.upkeepCache.Set(upkeepId.String(), tt.upkeepInfo, cache.DefaultExpiration)
			} else {
				mockReg := &mocks.KeeperRegistryInterface{}
				r.registry = mockReg
				mockReg.On("GetUpkeep", mock.Anything, mock.Anything).Return(tt.upkeepInfo, tt.upkeepInfoErr)
			}

			got, err := r.offchainLookup(context.Background(), tt.input)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), tt.name)
				assert.NotNil(t, err, tt.name)
			}
			assert.Equal(t, tt.want, got, tt.name)
		})
	}
}

func TestEvmRegistry_decodeOffchainLookup(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    OffchainLookup
		wantErr error
	}{
		{
			name: "success",
			data: []byte{218, 139, 50, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 96, 183, 114, 215, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 35, 104, 116, 116, 112, 115, 58, 47, 47, 112, 111, 107, 101, 97, 112, 105, 46, 99, 111, 47, 97, 112, 105, 47, 118, 50, 47, 112, 111, 107, 101, 109, 111, 110, 47, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 105, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 110, 97, 109, 101, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want: OffchainLookup{
				url:              "https://pokeapi.co/api/v2/pokemon/1",
				extraData:        []uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x30, 0x78, 0x30, 0x30, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				fields:           []string{"id", "name"},
				callbackFunction: [4]uint8{0xb7, 0x72, 0xd7, 0xa},
			},
			wantErr: nil,
		},
		{
			name:    "fail",
			data:    []byte{},
			want:    OffchainLookup{},
			wantErr: errors.New("unpack error: invalid data for unpacking"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)
			got, err := r.decodeOffchainLookup(tt.data)
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
	callbackPayload := []byte{183, 114, 215, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 200, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		name           string
		offchainLookup OffchainLookup
		values         []string
		statusCode     int
		upkeepInfo     keeper_registry_wrapper2_0.UpkeepInfo
		opts           *bind.CallOpts

		callbackMsg  ethereum.CallMsg
		callbackResp []byte
		callbackErr  error

		upkeepNeeded bool
		performData  []byte
		wantErr      error
	}{
		{
			name: "success",
			offchainLookup: OffchainLookup{
				url:              "https://pokeapi.co/api/v2/pokemon/1",
				extraData:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				fields:           []string{"id", "name"},
				callbackFunction: [4]byte{183, 114, 215, 10},
			},
			values:     []string{"1", "bulbasaur"},
			statusCode: http.StatusOK,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:         to,
				ExecuteGas:     executeGas,
				OffchainConfig: nil,
			},
			opts: &bind.CallOpts{
				BlockNumber: big.NewInt(999),
			},

			callbackMsg: ethereum.CallMsg{
				From: from,
				To:   &to,
				Gas:  uint64(gas),
				Data: callbackPayload,
			},
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},

			upkeepNeeded: true,
			performData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 98, 117, 108, 98, 97, 115, 97, 117, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "failure - bad response",
			offchainLookup: OffchainLookup{
				url:              "https://pokeapi.co/api/v2/pokemon/1",
				extraData:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				fields:           []string{"id", "name"},
				callbackFunction: [4]byte{183, 114, 215, 10},
			},
			values:     []string{"1", "bulbasaur"},
			statusCode: http.StatusOK,
			upkeepInfo: keeper_registry_wrapper2_0.UpkeepInfo{
				Target:         to,
				ExecuteGas:     executeGas,
				OffchainConfig: nil,
			},
			opts: &bind.CallOpts{
				BlockNumber: big.NewInt(999),
			},

			callbackMsg: ethereum.CallMsg{
				From: from,
				To:   &to,
				Gas:  uint64(gas),
				Data: callbackPayload,
			},
			callbackResp: []byte{},

			wantErr: errors.New("callback ouput unpack error:: abi: attempting to unmarshall an empty string while arguments are expected"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(evmmocks.Client)
			r := setupEVMRegistry(t)
			r.addr = from
			r.client = client
			client.On("CallContract", mock.Anything, tt.callbackMsg, tt.opts.BlockNumber).Return(tt.callbackResp, tt.callbackErr)

			upkeepNeeded, performData, err := r.offchainLookupCallback(context.Background(), tt.offchainLookup, tt.values, tt.statusCode, tt.upkeepInfo, tt.opts)
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

			val, exp, b := r.apiErrCache.GetWithExpiration(cacheKey)
			assert.True(t, b, "cache key found in apiErrCache")
			assert.NotNil(t, exp, "expiration found in apiErrCache")
			assert.GreaterOrEqual(t, exp, now.Add(DefaultApiErrExpiration-1*time.Minute), "expiration found in apiErrCache >= Default-1Minute")
			assert.Equal(t, tt.rounds, val, "err count correct")
			errCount := val.(int)

			val, exp, b = r.cooldownCache.GetWithExpiration(cacheKey)
			assert.True(t, b, "cache key found in cooldownCache")
			assert.NotNil(t, exp, "expiration found in cooldownCache")
			cooldown := time.Second * time.Duration(2^errCount)
			assert.GreaterOrEqual(t, exp, now.Add(cooldown/2), "expiration found in cooldownCache >= cooldown/2")
			assert.Equal(t, nil, val, "err count correct")
		})
	}
}

// TODO really test parsing
func TestOffchainLookup_parseJson(t *testing.T) {
	content, e := os.ReadFile("poke_api.json")
	assert.Nil(t, e)
	tests := []struct {
		name           string
		offchainLookup OffchainLookup
		body           []byte

		want    []string
		wantErr error
	}{
		{
			name: "success",
			offchainLookup: OffchainLookup{
				fields: []string{"id", "name"},
			},
			body: content,
			want: []string{"1", "bulbasaur"},
		},
		{
			name: "fail to unmarshal",
			offchainLookup: OffchainLookup{
				fields: []string{"", ""},
			},
			body:    []byte{},
			wantErr: errors.New("unexpected end of JSON input"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.offchainLookup.parseJson(tt.body)
			assert.Equal(t, tt.want, got, tt.name)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), tt.name)
				assert.NotNil(t, err, tt.name)
			}
		})
	}
}
