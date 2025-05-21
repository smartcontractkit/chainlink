package mercury

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	automationTypes "github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
)

func TestGenerateHMACFn(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		path     string
		body     []byte
		clientId string
		secret   string
		ts       int64
		expected string
	}{
		{
			name:     "generate hmac function",
			method:   "GET",
			path:     "/example",
			body:     []byte(""),
			clientId: "yourClientId",
			secret:   "yourSecret",
			ts:       1234567890,
			expected: "17b0bb6b14f7b48ef9d24f941ff8f33ad2d5e94ac343380be02c2f1ca32fdbd8",
		},
		{
			name:     "generate hmac function with non-empty body",
			method:   "POST",
			path:     "/api",
			body:     []byte("request body"),
			clientId: "anotherClientId",
			secret:   "anotherSecret",
			ts:       1597534567,
			expected: "d326c168c50c996e271d6b3b4c97944db01163994090f73fcf4fd42f23f06bbb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateHMACFn(tc.method, tc.path, tc.body, tc.clientId, tc.secret, tc.ts)

			if result != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, result)
			}
		})
	}
}

func TestPacker_DecodeStreamsLookupRequest(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected *StreamsLookupError
		state    uint8
		err      error
	}{
		{
			name: "success - decode to streams lookup",
			data: hexutil.MustDecode("0xf055e4a200000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002400000000000000000000000000000000000000000000000000000000002435eb50000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000000966656564496448657800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000423078343535343438326435353533343432643431353234323439353435323535346432643534343535333534346534353534303030303030303030303030303030300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000042307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b626c6f636b4e756d62657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
			expected: &StreamsLookupError{
				FeedParamKey: "feedIdHex",
				Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				TimeParamKey: "blockNumber",
				Time:         big.NewInt(37969589),
				ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			},
		},
		{
			name: "failure - unpack error",
			data: []byte{1, 2, 3, 4},
			err:  errors.New("unpack error: invalid identifier, have 0x01020304 want 0xf055e4a2"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packer := NewAbiPacker()
			fl, err := packer.DecodeStreamsLookupRequest(tt.data)
			assert.Equal(t, tt.expected, fl)
			if tt.err != nil {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestPacker_UnpackGetUpkeepPrivilegeConfig(t *testing.T) {
	tests := []struct {
		name    string
		raw     []byte
		errored bool
	}{
		{
			name: "happy path",
			raw: func() []byte {
				b, _ := hexutil.Decode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000177b226d657263757279456e61626c6564223a747275657d000000000000000000")

				return b
			}(),
			errored: false,
		},
		{
			name: "error empty config",
			raw: func() []byte {
				b, _ := hexutil.Decode("0x")

				return b
			}(),
			errored: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			packer := NewAbiPacker()

			b, err := packer.UnpackGetUpkeepPrivilegeConfig(test.raw)

			if !test.errored {
				require.NoError(t, err, "should unpack bytes from abi encoded value")

				// the actual struct to unmarshal into is not available to this
				// package so basic json encoding is the limit of the following test
				var data map[string]interface{}
				err = json.Unmarshal(b, &data)

				assert.NoError(t, err, "packed data should unmarshal using json encoding")
				assert.Equal(t, []byte(`{"mercuryEnabled":true}`), b)
			} else {
				assert.Error(t, err, "error expected from unpack function")
			}
		})
	}
}

func TestPacker_PackGetUpkeepPrivilegeConfig(t *testing.T) {
	tests := []struct {
		name     string
		upkeepId *big.Int
		raw      []byte
		errored  bool
	}{
		{
			name: "happy path",
			upkeepId: func() *big.Int {
				id, _ := new(big.Int).SetString("52236098515066839510538748191966098678939830769967377496848891145101407612976", 10)

				return id
			}(),
			raw: func() []byte {
				b, _ := hexutil.Decode("0x19d97a94737c9583000000000000000000000001ea8ed6d0617dd5b3b87374020efaf030")

				return b
			}(),
			errored: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			packer := NewAbiPacker()

			b, err := packer.PackGetUpkeepPrivilegeConfig(test.upkeepId)

			if !test.errored {
				require.NoError(t, err, "no error expected from packing")

				assert.Equal(t, test.raw, b, "raw bytes for output should match expected")
			} else {
				assert.Error(t, err, "error expected from packing function")
			}
		})
	}
}

func TestPacker_UnpackCheckCallbackResult(t *testing.T) {
	tests := []struct {
		Name          string
		CallbackResp  []byte
		UpkeepNeeded  bool
		PerformData   []byte
		FailureReason encoding.UpkeepFailureReason
		GasUsed       *big.Int
		ErrorString   string
		State         encoding.PipelineExecutionState
	}{
		{
			Name:          "unpack upkeep needed",
			CallbackResp:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 46, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded:  true,
			PerformData:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			FailureReason: encoding.UpkeepFailureReasonNone,
			GasUsed:       big.NewInt(11796),
		},
		{
			Name:          "unpack upkeep not needed",
			CallbackResp:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 50, 208, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded:  false,
			PerformData:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			FailureReason: encoding.UpkeepFailureReasonUpkeepNotNeeded,
			GasUsed:       big.NewInt(13008),
		},
		{
			Name:         "unpack malformed data",
			CallbackResp: []byte{0, 0, 0, 23, 4, 163, 66, 91, 228, 102, 200, 84, 144, 233, 218, 44, 168, 192, 191, 253, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded: false,
			PerformData:  nil,
			ErrorString:  "abi: improperly encoded boolean value: unpack checkUpkeep return: ",
			State:        encoding.PackUnpackDecodeFailed,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer := NewAbiPacker()

			state, needed, pd, failureReason, gasUsed, err := packer.UnpackCheckCallbackResult(test.CallbackResp)

			if test.ErrorString != "" {
				assert.EqualError(t, err, test.ErrorString+hexutil.Encode(test.CallbackResp))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.UpkeepNeeded, needed)
			assert.Equal(t, test.PerformData, pd)
			assert.Equal(t, test.FailureReason, failureReason)
			assert.Equal(t, test.GasUsed, gasUsed)
			assert.Equal(t, test.State, state)
		})
	}
}

func TestPacker_PackUserCheckErrorHandler(t *testing.T) {
	tests := []struct {
		name      string
		errCode   encoding.ErrCode
		extraData []byte
		rawOutput []byte
		errored   bool
	}{
		{
			name:    "happy path",
			errCode: encoding.ErrCodeStreamsBadRequest,
			extraData: func() []byte {
				b, _ := hexutil.Decode("0x19d97a94737c9583000000000000000000000001ea8ed6d0617dd5b3b87374020efaf030")

				return b
			}(),
			rawOutput: []byte{0xf, 0xb1, 0x72, 0xfb, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc, 0x55, 0xd0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x24, 0x19, 0xd9, 0x7a, 0x94, 0x73, 0x7c, 0x95, 0x83, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0xea, 0x8e, 0xd6, 0xd0, 0x61, 0x7d, 0xd5, 0xb3, 0xb8, 0x73, 0x74, 0x2, 0xe, 0xfa, 0xf0, 0x30, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			errored:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			packer := NewAbiPacker()

			b, err := packer.PackUserCheckErrorHandler(test.errCode, test.extraData)

			if !test.errored {
				require.NoError(t, err, "no error expected from packing")

				assert.Equal(t, test.rawOutput, b, "raw bytes for output should match expected")
			} else {
				assert.Error(t, err, "error expected from packing function")
			}
		})
	}
}

func Test_CalculateRetryConfigFn(t *testing.T) {
	tests := []struct {
		name       string
		times      int
		upkeepType automationTypes.UpkeepType
		expected   time.Duration
	}{
		{
			name:       "first retry",
			times:      1,
			upkeepType: automationTypes.LogTrigger,
			expected:   1 * time.Second,
		},
		{
			name:       "second retry",
			times:      2,
			upkeepType: automationTypes.LogTrigger,
			expected:   1 * time.Second,
		},
		{
			name:       "fifth retry",
			times:      5,
			upkeepType: automationTypes.LogTrigger,
			expected:   1 * time.Second,
		},
		{
			name:       "sixth retry",
			times:      6,
			upkeepType: automationTypes.LogTrigger,
			expected:   5 * time.Second,
		},
		{
			name:       "timeout",
			times:      totalMediumPluginRetries + 1,
			upkeepType: automationTypes.LogTrigger,
			expected:   RetryIntervalTimeout,
		},
		{
			name:       "conditional first timeout",
			times:      1,
			upkeepType: automationTypes.ConditionTrigger,
			expected:   RetryIntervalTimeout,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := newMercuryConfigMock()
			var result time.Duration
			for i := 0; i < tc.times; i++ {
				result = CalculateStreamsRetryConfigFn(tc.upkeepType, "prk", cfg)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

type mercuryConfigMock struct {
	pluginRetryCache *cache.Cache
}

func newMercuryConfigMock() *mercuryConfigMock {
	return &mercuryConfigMock{
		pluginRetryCache: cache.New(10*time.Second, time.Minute),
	}
}

func (c *mercuryConfigMock) Credentials() *types.MercuryCredentials {
	return nil
}

func (c *mercuryConfigMock) IsUpkeepAllowed(k string) (interface{}, bool) {
	return nil, false
}

func (c *mercuryConfigMock) SetUpkeepAllowed(k string, v interface{}, d time.Duration) {
}

func (c *mercuryConfigMock) GetPluginRetry(k string) (interface{}, bool) {
	return c.pluginRetryCache.Get(k)
}

func (c *mercuryConfigMock) SetPluginRetry(k string, v interface{}, d time.Duration) {
	c.pluginRetryCache.Set(k, v, d)
}
