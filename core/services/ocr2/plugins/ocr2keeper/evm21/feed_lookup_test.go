package evm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/mocks"

	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/feed_lookup_compatible_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
)

// setups up an evm registry for tests.
func setupEVMRegistry(t *testing.T) *EvmRegistry {
	lggr := logger.TestLogger(t)
	addr := common.HexToAddress("0x6cA639822c6C241Fa9A7A6b5032F6F7F1C513CAD")
	keeperRegistryABI, err := abi.JSON(strings.NewReader(i_keeper_registry_master_wrapper_2_1.IKeeperRegistryMasterABI))
	require.Nil(t, err, "need registry abi")
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
	require.Nil(t, err, "need utils abi")
	feedLookupCompatibleABI, err := abi.JSON(strings.NewReader(feed_lookup_compatible_interface.FeedLookupCompatibleInterfaceABI))
	require.Nil(t, err, "need mercury abi")
	var logPoller logpoller.LogPoller
	mockRegistry := mocks.NewRegistry(t)
	mockHttpClient := mocks.NewHttpClient(t)
	client := evmClientMocks.NewClient(t)

	r := &EvmRegistry{
		lggr:     lggr,
		poller:   logPoller,
		addr:     addr,
		client:   client,
		txHashes: make(map[string]bool),
		registry: mockRegistry,
		abi:      keeperRegistryABI,
		active:   make(map[string]activeUpkeep),
		packer:   NewEvmRegistryPackerV2_1(keeperRegistryABI, utilsABI),
		headFunc: func(ocr2keepers.BlockKey) {},
		chLog:    make(chan logpoller.Log, 1000),
		mercury: &MercuryConfig{
			cred: &models.MercuryCredentials{
				URL:      "https://google.com",
				Username: "FakeClientID",
				Password: "FakeClientKey",
			},
			abi:            feedLookupCompatibleABI,
			allowListCache: cache.New(defaultAllowListExpiration, cleanupInterval),
		},
		hc: mockHttpClient,
	}
	return r
}

func TestEvmRegistry_FeedLookup(t *testing.T) {
	upkeepId, ok := new(big.Int).SetString("71022726777042968814359024671382968091267501884371696415772139504780367423725", 10)
	var upkeepIdentifier [32]byte
	copy(upkeepIdentifier[:], upkeepId.Bytes())
	assert.True(t, ok, t.Name())
	tests := []struct {
		name              string
		input             []ocr2keepers.CheckResult
		blob              string
		callbackResp      []byte
		expectedResults   []ocr2keepers.CheckResult
		callbackNeeded    bool
		extraData         []byte
		checkCallbackResp []byte
		values            [][]byte
		cachedAdminCfg    bool
		hasError          bool
		hasPermission     bool
	}{
		{
			name: "success - happy path no cache",
			input: []ocr2keepers.CheckResult{
				{
					PerformData: []byte{125, 221, 147, 62, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 141, 110, 193, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 102, 101, 101, 100, 73, 100, 72, 101, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 66, 48, 120, 52, 53, 53, 52, 52, 56, 50, 100, 53, 53, 53, 51, 52, 52, 50, 100, 52, 49, 53, 50, 52, 50, 52, 57, 53, 52, 53, 50, 53, 53, 52, 100, 50, 100, 53, 52, 52, 53, 53, 51, 53, 52, 52, 101, 52, 53, 53, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 98, 108, 111, 99, 107, 78, 117, 109, 98, 101, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: 26046145,
					},
					IneligibilityReason: uint8(UpkeepFailureReasonTargetCheckReverted),
				},
			},
			blob:              "0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000159761000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e4554000000000000000000000000000000000000000000000000000000000000000000000000648a1fbb000000000000000000000000000000000000000000000000000000274421041500000000000000000000000000000000000000000000000000000027437c6ecd0000000000000000000000000000000000000000000000000000002744c5995d00000000000000000000000000000000000000000000000000000000018d6ec108936dfe39c48715572a51ac868129958f937fb95ef5abdf73a239cf86a4fee700000000000000000000000000000000000000000000000000000000018d6ec100000000000000000000000000000000000000000000000000000000648a1fbb00000000000000000000000000000000000000000000000000000000000000028a26e557ee2feb91ccb116f3ab4eb1469afe5c3b012538cb151dbe3fbceaf6f117b24ac2a82cff25b286ae0a9b903dc6badaa16f6e67bf0983461b008574e30a00000000000000000000000000000000000000000000000000000000000000020db5c5924481061b98df59caefd9c4c1e72657c4976bf7c7568730fbdaf828080bff6b1edea2c8fed5e8bbac5574aa94cf809d898f5055cb1db14a16f1493726",
			cachedAdminCfg:    false,
			extraData:         []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			callbackNeeded:    true,
			checkCallbackResp: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000428200000000000000000000000000000000000000000000000000000000000003c0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000003800000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002e000066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000159761000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e4554000000000000000000000000000000000000000000000000000000000000000000000000648a1fbb000000000000000000000000000000000000000000000000000000274421041500000000000000000000000000000000000000000000000000000027437c6ecd0000000000000000000000000000000000000000000000000000002744c5995d00000000000000000000000000000000000000000000000000000000018d6ec108936dfe39c48715572a51ac868129958f937fb95ef5abdf73a239cf86a4fee700000000000000000000000000000000000000000000000000000000018d6ec100000000000000000000000000000000000000000000000000000000648a1fbb00000000000000000000000000000000000000000000000000000000000000028a26e557ee2feb91ccb116f3ab4eb1469afe5c3b012538cb151dbe3fbceaf6f117b24ac2a82cff25b286ae0a9b903dc6badaa16f6e67bf0983461b008574e30a00000000000000000000000000000000000000000000000000000000000000020db5c5924481061b98df59caefd9c4c1e72657c4976bf7c7568730fbdaf828080bff6b1edea2c8fed5e8bbac5574aa94cf809d898f5055cb1db14a16f14937260000000000000000000000000000000000000000000000000000000000000008786f657a5a362c01000000000000000000000000000000000000000000000000"),
			values:            [][]byte{{0, 6, 109, 252, 209, 237, 45, 149, 177, 140, 148, 141, 188, 91, 214, 76, 104, 122, 254, 147, 228, 202, 125, 102, 61, 222, 193, 76, 32, 9, 10, 216, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 21, 151, 97, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 128, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 138, 31, 187, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 39, 68, 33, 4, 21, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 39, 67, 124, 110, 205, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 39, 68, 197, 153, 93, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 141, 110, 193, 8, 147, 109, 254, 57, 196, 135, 21, 87, 42, 81, 172, 134, 129, 41, 149, 143, 147, 127, 185, 94, 245, 171, 223, 115, 162, 57, 207, 134, 164, 254, 231, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 141, 110, 193, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 138, 31, 187, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 138, 38, 229, 87, 238, 47, 235, 145, 204, 177, 22, 243, 171, 78, 177, 70, 154, 254, 92, 59, 1, 37, 56, 203, 21, 29, 190, 63, 188, 234, 246, 241, 23, 178, 74, 194, 168, 44, 255, 37, 178, 134, 174, 10, 155, 144, 61, 198, 186, 218, 161, 111, 110, 103, 191, 9, 131, 70, 27, 0, 133, 116, 227, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 13, 181, 197, 146, 68, 129, 6, 27, 152, 223, 89, 202, 239, 217, 196, 193, 231, 38, 87, 196, 151, 107, 247, 199, 86, 135, 48, 251, 218, 248, 40, 8, 11, 255, 107, 30, 222, 162, 200, 254, 213, 232, 187, 172, 85, 116, 170, 148, 207, 128, 157, 137, 143, 80, 85, 203, 29, 177, 74, 22, 241, 73, 55, 38}},
			expectedResults: []ocr2keepers.CheckResult{
				{
					Eligible:    true,
					PerformData: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 224, 0, 6, 109, 252, 209, 237, 45, 149, 177, 140, 148, 141, 188, 91, 214, 76, 104, 122, 254, 147, 228, 202, 125, 102, 61, 222, 193, 76, 32, 9, 10, 216, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 21, 151, 97, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 128, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 138, 31, 187, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 39, 68, 33, 4, 21, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 39, 67, 124, 110, 205, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 39, 68, 197, 153, 93, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 141, 110, 193, 8, 147, 109, 254, 57, 196, 135, 21, 87, 42, 81, 172, 134, 129, 41, 149, 143, 147, 127, 185, 94, 245, 171, 223, 115, 162, 57, 207, 134, 164, 254, 231, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 141, 110, 193, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 138, 31, 187, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 138, 38, 229, 87, 238, 47, 235, 145, 204, 177, 22, 243, 171, 78, 177, 70, 154, 254, 92, 59, 1, 37, 56, 203, 21, 29, 190, 63, 188, 234, 246, 241, 23, 178, 74, 194, 168, 44, 255, 37, 178, 134, 174, 10, 155, 144, 61, 198, 186, 218, 161, 111, 110, 103, 191, 9, 131, 70, 27, 0, 133, 116, 227, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 13, 181, 197, 146, 68, 129, 6, 27, 152, 223, 89, 202, 239, 217, 196, 193, 231, 38, 87, 196, 151, 107, 247, 199, 86, 135, 48, 251, 218, 248, 40, 8, 11, 255, 107, 30, 222, 162, 200, 254, 213, 232, 187, 172, 85, 116, 170, 148, 207, 128, 157, 137, 143, 80, 85, 203, 29, 177, 74, 22, 241, 73, 55, 38, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: 26046145,
					},
					IneligibilityReason: uint8(UpkeepFailureReasonNone),
				},
			},
			hasError:      false,
			hasPermission: true,
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
					IneligibilityReason: uint8(UpkeepFailureReasonInsufficientBalance),
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
					IneligibilityReason: uint8(UpkeepFailureReasonInsufficientBalance),
				},
			},
			hasError: true,
		},
		{
			name: "skip - no mercury permission",
			input: []ocr2keepers.CheckResult{
				{
					PerformData: []byte{},
					UpkeepID:    upkeepIdentifier,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: 26046145,
					},
					IneligibilityReason: uint8(UpkeepFailureReasonTargetCheckReverted),
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
					IneligibilityReason: uint8(UpkeepFailureReasonMercuryAccessNotAllowed),
				},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)

			if !tt.cachedAdminCfg && !tt.hasError {
				mockRegistry := mocks.NewRegistry(t)
				cfg := AdminOffchainConfig{MercuryEnabled: tt.hasPermission}
				b, err := json.Marshal(cfg)
				assert.Nil(t, err)
				mockRegistry.On("GetUpkeepPrivilegeConfig", mock.Anything, upkeepId).Return(b, nil)
				r.registry = mockRegistry
			}

			if tt.blob != "" {
				hc := mocks.NewHttpClient(t)
				mr := MercuryResponse{ChainlinkBlob: tt.blob}
				b, err := json.Marshal(mr)
				assert.Nil(t, err)
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(b)),
				}
				hc.On("Do", mock.Anything).Return(resp, nil).Once()
				r.hc = hc
			}

			if tt.callbackNeeded {
				payload, err := r.abi.Pack("checkCallback", upkeepId, tt.values, tt.extraData)
				require.Nil(t, err)
				args := map[string]interface{}{
					"to":   r.addr.Hex(),
					"data": hexutil.Bytes(payload),
				}
				client := new(evmClientMocks.Client)
				client.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", args, hexutil.EncodeUint64(uint64(26046145))).Return(nil).
					Run(func(args mock.Arguments) {
						b := args.Get(1).(*hexutil.Bytes)
						*b = tt.checkCallbackResp
					}).Once()
				r.client = client
			}

			got := r.feedLookup(context.Background(), tt.input)
			assert.Equal(t, tt.expectedResults, got, tt.name)
		})
	}
}

func TestEvmRegistry_DecodeFeedLookup(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected *FeedLookup
		state    PipelineExecutionState
		err      error
	}{
		{
			name:  "success - decode to feed lookup",
			data:  []byte{125, 221, 147, 62, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 138, 215, 253, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 102, 101, 101, 100, 73, 100, 72, 101, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 66, 48, 120, 52, 53, 53, 52, 52, 56, 50, 100, 53, 53, 53, 51, 52, 52, 50, 100, 52, 49, 53, 50, 52, 50, 52, 57, 53, 52, 53, 50, 53, 53, 52, 100, 50, 100, 53, 52, 52, 53, 53, 51, 53, 52, 52, 101, 52, 53, 53, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 66, 48, 120, 52, 50, 53, 52, 52, 51, 50, 100, 53, 53, 53, 51, 52, 52, 50, 100, 52, 49, 53, 50, 52, 50, 52, 57, 53, 52, 53, 50, 53, 53, 52, 100, 50, 100, 53, 52, 52, 53, 53, 51, 53, 52, 52, 101, 52, 53, 53, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 98, 108, 111, 99, 107, 78, 117, 109, 98, 101, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			state: NoPipelineError,
			expected: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(25876477),
				extraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			},
		},
		{
			name:  "failure - unpack error",
			data:  []byte{1, 2, 3, 4},
			err:   errors.New("unpack error: invalid data for unpacking"),
			state: PackUnpackDecodeFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)
			state, fl, err := r.decodeFeedLookup(tt.data)
			assert.Equal(t, tt.expected, fl)
			assert.Equal(t, tt.state, state)
			if tt.err != nil {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestEvmRegistry_AllowedToUseMercury(t *testing.T) {
	upkeepId, ok := new(big.Int).SetString("71022726777042968814359024671382968091267501884371696415772139504780367423725", 10)
	assert.True(t, ok, t.Name())
	tests := []struct {
		name         string
		cached       bool
		allowed      bool
		errorMessage string
		state        PipelineExecutionState
		retryable    bool
	}{
		{
			name:    "success - allowed via cache",
			cached:  true,
			allowed: true,
		},
		{
			name:    "success - allowed via fetching privilege config",
			cached:  false,
			allowed: true,
		},
		{
			name:    "success - not allowed via cache",
			cached:  true,
			allowed: false,
		},
		{
			name:    "success - not allowed via fetching privilege config",
			cached:  false,
			allowed: false,
		},
		{
			name:         "failure - cannot unmarshal privilege config",
			cached:       false,
			errorMessage: "failed to unmarshal privilege config for upkeep ID 71022726777042968814359024671382968091267501884371696415772139504780367423725: invalid character '\\x00' looking for beginning of value",
			state:        MercuryUnmarshalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupEVMRegistry(t)

			if tt.errorMessage != "" {
				mockRegistry := mocks.NewRegistry(t)
				mockRegistry.On("GetUpkeepPrivilegeConfig", mock.Anything, upkeepId).Return([]byte{0, 1}, nil)
				r.registry = mockRegistry
			} else {
				if tt.cached {
					r.mercury.allowListCache.Set(upkeepId.String(), tt.allowed, cache.DefaultExpiration)
				} else {
					mockRegistry := mocks.NewRegistry(t)
					cfg := AdminOffchainConfig{MercuryEnabled: tt.allowed}
					b, err := json.Marshal(cfg)
					assert.Nil(t, err)
					mockRegistry.On("GetUpkeepPrivilegeConfig", mock.Anything, upkeepId).Return(b, nil)
					r.registry = mockRegistry
				}
			}

			state, retryable, allowed, err := r.allowedToUseMercury(nil, upkeepId)
			if tt.errorMessage != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.allowed, allowed)
			}
			assert.Equal(t, tt.state, state)
			assert.Equal(t, tt.retryable, retryable)
		})
	}
}

func TestEvmRegistry_DoMercuryRequest(t *testing.T) {
	upkeepId := big.NewInt(0)
	upkeepId.SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)

	tests := []struct {
		name               string
		lookup             *FeedLookup
		mockHttpStatusCode int
		mockChainlinkBlobs []string
		expectedValues     [][]byte
		expectedRetryable  bool
		expectedError      error
		state              PipelineExecutionState
	}{
		{
			name: "success",
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(25880526),
				extraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				upkeepId:     upkeepId,
			},
			mockHttpStatusCode: http.StatusOK,
			mockChainlinkBlobs: []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:     [][]byte{{0, 6, 109, 252, 209, 237, 45, 149, 177, 140, 148, 141, 188, 91, 214, 76, 104, 122, 254, 147, 228, 202, 125, 102, 61, 222, 193, 76, 32, 9, 10, 216, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 20, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 128, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 137, 28, 152, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 40, 154, 216, 211, 103, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 40, 154, 207, 11, 56, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 40, 155, 61, 164, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 138, 231, 206, 116, 217, 250, 37, 42, 137, 131, 151, 110, 171, 96, 13, 199, 89, 12, 119, 141, 4, 129, 52, 48, 132, 27, 198, 231, 101, 195, 76, 216, 26, 22, 141, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 138, 231, 203, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 137, 28, 152, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 96, 65, 43, 148, 229, 37, 202, 108, 237, 201, 245, 68, 253, 134, 247, 118, 6, 213, 47, 231, 49, 165, 208, 105, 219, 232, 54, 168, 191, 192, 251, 140, 145, 25, 99, 176, 174, 122, 20, 151, 31, 59, 70, 33, 191, 251, 128, 46, 240, 96, 83, 146, 185, 166, 200, 156, 127, 171, 29, 248, 99, 58, 90, 222, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 69, 0, 194, 245, 33, 248, 63, 186, 94, 252, 43, 243, 239, 250, 174, 221, 228, 61, 10, 74, 223, 247, 133, 193, 33, 59, 113, 42, 58, 237, 13, 129, 87, 100, 42, 132, 50, 77, 176, 207, 150, 149, 235, 210, 119, 8, 212, 96, 142, 176, 51, 126, 13, 216, 123, 14, 67, 240, 250, 112, 199, 0, 217, 17}},
			expectedRetryable:  false,
			expectedError:      nil,
		},
		{
			name: "failure - retryable",
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(25880526),
				extraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				upkeepId:     upkeepId,
			},
			mockHttpStatusCode: http.StatusInternalServerError,
			mockChainlinkBlobs: []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:     [][]byte{nil},
			expectedRetryable:  true,
			expectedError:      errors.New("All attempts fail:\n#1: 500\n#2: 500\n#3: 500"),
			state:              MercuryFlakyFailure,
		},
		{
			name: "failure - not retryable",
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(25880526),
				extraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				upkeepId:     upkeepId,
			},
			mockHttpStatusCode: http.StatusBadGateway,
			mockChainlinkBlobs: []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:     [][]byte{nil},
			expectedRetryable:  false,
			expectedError:      errors.New("All attempts fail:\n#1: FeedLookup upkeep 88786950015966611018675766524283132478093844178961698330929478019253453382042 block 25880526 received status code 502 for feed 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"),
			state:              InvalidMercuryRequest,
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

			state, values, retryable, reqErr := r.doMercuryRequest(context.Background(), tt.lookup)
			assert.Equal(t, tt.expectedValues, values)
			assert.Equal(t, tt.expectedRetryable, retryable)
			assert.Equal(t, tt.state, state)
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError.Error(), reqErr.Error())
			}
		})
	}
}

func TestEvmRegistry_SingleFeedRequest(t *testing.T) {
	upkeepId := big.NewInt(123456789)
	tests := []struct {
		name           string
		index          int
		lookup         *FeedLookup
		mv             MercuryVersion
		blob           string
		statusCode     int
		lastStatusCode int
		retryNumber    int
		retryable      bool
		errorMessage   string
	}{
		{
			name:  "success - mercury responds in the first try",
			index: 0,
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			mv:   MercuryV02,
			blob: "0xab2123dc00000012",
		},
		{
			name:  "success - retry for 404",
			index: 0,
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			mv:             MercuryV02,
			blob:           "0xab2123dcbabbad",
			retryNumber:    1,
			statusCode:     http.StatusNotFound,
			lastStatusCode: http.StatusOK,
		},
		{
			name:  "success - retry for 500",
			index: 0,
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			mv:             MercuryV02,
			blob:           "0xab2123dcbbabad",
			retryNumber:    2,
			statusCode:     http.StatusInternalServerError,
			lastStatusCode: http.StatusOK,
		},
		{
			name:  "failure - returns retryable",
			index: 0,
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			mv:           MercuryV02,
			blob:         "0xab2123dc",
			retryNumber:  TotalAttempt,
			statusCode:   http.StatusNotFound,
			retryable:    true,
			errorMessage: "All attempts fail:\n#1: 404\n#2: 404\n#3: 404",
		},
		{
			name:  "failure - returns retryable and then non-retryable",
			index: 0,
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			mv:             MercuryV02,
			blob:           "0xab2123dc",
			retryNumber:    1,
			statusCode:     http.StatusNotFound,
			lastStatusCode: http.StatusBadGateway,
			errorMessage:   "All attempts fail:\n#1: 404\n#2: FeedLookup upkeep 123456789 block 123456 received status code 502 for feed 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
		},
		{
			name:  "failure - returns not retryable",
			index: 0,
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			mv:           MercuryV02,
			blob:         "0xab2123dc",
			statusCode:   http.StatusBadGateway,
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
					StatusCode: tt.lastStatusCode,
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

			ch := make(chan MercuryData, 1)
			r.singleFeedRequest(context.Background(), ch, tt.index, tt.lookup, tt.mv)

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

func TestEvmRegistry_MultiFeedRequest(t *testing.T) {
	upkeepId := big.NewInt(123456789)
	tests := []struct {
		name           string
		lookup         *FeedLookup
		blob           string
		statusCode     int
		lastStatusCode int
		retryNumber    int
		retryable      bool
		errorMessage   string
	}{
		{
			name: "success - mercury responds in the first try",
			lookup: &FeedLookup{
				feedParamKey: FeedId,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: Timestamp,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			blob: "0xab2123dc00000012",
		},
		{
			name: "success - retry for 404",
			lookup: &FeedLookup{
				feedParamKey: FeedId,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: Timestamp,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			blob:           "0xab2123dcbabbad",
			retryNumber:    1,
			statusCode:     http.StatusNotFound,
			lastStatusCode: http.StatusOK,
		},
		{
			name: "success - retry for 500",
			lookup: &FeedLookup{
				feedParamKey: FeedId,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: Timestamp,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			blob:           "0xab2123dcbbabad",
			retryNumber:    2,
			statusCode:     http.StatusInternalServerError,
			lastStatusCode: http.StatusOK,
		},
		{
			name: "failure - returns retryable",
			lookup: &FeedLookup{
				feedParamKey: FeedId,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: Timestamp,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			blob:         "0xab2123dc",
			retryNumber:  TotalAttempt,
			statusCode:   http.StatusNotFound,
			retryable:    true,
			errorMessage: "All attempts fail:\n#1: 404\n#2: 404\n#3: 404",
		},
		{
			name: "failure - returns retryable and then non-retryable",
			lookup: &FeedLookup{
				feedParamKey: FeedId,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: Timestamp,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			blob:           "0xab2123dc",
			retryNumber:    1,
			statusCode:     http.StatusNotFound,
			lastStatusCode: http.StatusBadGateway,
			errorMessage:   "All attempts fail:\n#1: 404\n#2: FeedLookup upkeep 123456789 block 123456 received status code 502 for multi feed",
		},
		{
			name: "failure - returns not retryable",
			lookup: &FeedLookup{
				feedParamKey: FeedId,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: Timestamp,
				time:         big.NewInt(123456),
				upkeepId:     upkeepId,
			},
			blob:         "0xab2123dc",
			statusCode:   http.StatusBadGateway,
			errorMessage: "All attempts fail:\n#1: FeedLookup upkeep 123456789 block 123456 received status code 502 for multi feed",
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
					StatusCode: tt.lastStatusCode,
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

			ch := make(chan MercuryData, 1)
			r.multiFeedsRequest(context.Background(), ch, tt.lookup)

			m := <-ch
			assert.Equal(t, 0, m.Index)
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
	upkeepId := big.NewInt(123456789)
	blockNumber := uint64(999)
	bs := []byte{183, 114, 215, 10, 0, 0, 0, 0, 0, 0}
	values := [][]byte{bs}
	tests := []struct {
		name       string
		lookup     *FeedLookup
		values     [][]byte
		statusCode int

		callbackResp []byte
		callbackErr  error

		upkeepNeeded bool
		performData  []byte
		wantErr      assert.ErrorAssertionFunc

		state     PipelineExecutionState
		retryable bool
	}{
		{
			name: "success - empty extra data",
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"ETD-USD", "BTC-ETH"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(100),
				extraData:    []byte{48, 120, 48, 48},
				upkeepId:     upkeepId,
				block:        blockNumber,
			},
			values:       values,
			statusCode:   http.StatusOK,
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			upkeepNeeded: true,
			performData:  []byte{48, 120, 48, 48},
			wantErr:      assert.NoError,
		},
		{
			name: "success - with extra data",
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(18952430),
				// this is the address of precompile contract ArbSys(0x0000000000000000000000000000000000000064)
				extraData: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				upkeepId:  upkeepId,
				block:     blockNumber,
			},
			values:       values,
			statusCode:   http.StatusOK,
			callbackResp: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			upkeepNeeded: true,
			performData:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			wantErr:      assert.NoError,
		},
		{
			name: "failure - bad response",
			lookup: &FeedLookup{
				feedParamKey: FeedIdHex,
				feeds:        []string{"ETD-USD", "BTC-ETH"},
				timeParamKey: BlockNumber,
				time:         big.NewInt(100),
				extraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 48, 120, 48, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				upkeepId:     upkeepId,
				block:        blockNumber,
			},
			values:       values,
			statusCode:   http.StatusOK,
			callbackResp: []byte{},
			callbackErr:  errors.New("bad response"),
			wantErr:      assert.Error,
			state:        RpcFlakyFailure,
			retryable:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(evmClientMocks.Client)
			r := setupEVMRegistry(t)
			payload, err := r.abi.Pack("checkCallback", tt.lookup.upkeepId, values, tt.lookup.extraData)
			require.Nil(t, err)
			args := map[string]interface{}{
				"to":   r.addr.Hex(),
				"data": hexutil.Bytes(payload),
			}
			client.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", args, hexutil.EncodeUint64(tt.lookup.block)).Return(tt.callbackErr).
				Run(func(args mock.Arguments) {
					by := args.Get(1).(*hexutil.Bytes)
					*by = tt.callbackResp
				}).Once()
			r.client = client

			state, retryable, _, err := r.checkCallback(context.Background(), tt.values, tt.lookup)
			tt.wantErr(t, err, fmt.Sprintf("Error asserion failed: %v", tt.name))
			assert.Equal(t, tt.state, state)
			assert.Equal(t, tt.retryable, retryable)
		})
	}
}
