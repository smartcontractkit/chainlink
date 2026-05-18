package bridges_test

import (
	"encoding/json"
	"math/big"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/math"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBridgeTypeRequest(t *testing.T) {
	u, err := url.Parse("http://example.com/test")
	require.NoError(t, err)
	r := bridges.BridgeTypeRequest{
		Name:                   bridges.MustParseBridgeName("test-bridge-name"),
		URL:                    models.WebURL(*u),
		Confirmations:          math.MaxUint32,
		MinimumContractPayment: (*assets.Link)(big.NewInt(1000)),
	}
	assert.Equal(t, "bridges", r.GetName())
	assert.Equal(t, "test-bridge-name", r.GetID())
	const validID = "abc123foo_bar-test"
	assert.NoError(t, r.SetID(validID))
	assert.Equal(t, validID, r.GetID())
	assert.Error(t, r.SetID("abc123.,<>/.foobar"))
}

func TestBridgeType_Authenticate(t *testing.T) {
	t.Parallel()

	bta, bt := cltest.NewBridgeType(t, cltest.BridgeOpts{})
	tests := []struct {
		name, token string
		wantError   bool
	}{
		{"correct", bta.IncomingToken, false},
		{"incorrect", "gibberish", true},
		{"empty incorrect", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok, err := bridges.AuthenticateBridgeType(bt, test.token)
			require.NoError(t, err)

			if test.wantError {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
		})
	}
}

func BenchmarkParseBridgeName(b *testing.B) {
	const valid = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_`
	for _, l := range []int{1, 10, 20, 50, 100, 1000, 10000} {
		b.Run(strconv.Itoa(l), func(b *testing.B) {
			var sb strings.Builder
			for i := 0; i < l; i++ {
				sb.WriteByte(valid[rand.Intn(len(valid))])
			}
			name := sb.String()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := bridges.ParseBridgeName(name)
				if err != nil {
					b.Fatalf("failed to parse %q: %v\n", name, err)
				}
			}
		})
	}
}

func TestBridgeName_UnmarshalJSON(t *testing.T) {
	var b bridges.BridgeName
	require.NoError(t, json.Unmarshal([]byte(`"asdf123test"`), &b))
	require.Equal(t, "asdf123test", b.String())

	got, err := json.Marshal(b)
	require.NoError(t, err)
	require.Equal(t, []byte(`"asdf123test"`), got)

	require.Error(t, json.Unmarshal([]byte(`"invalid,.<>/asdf?"`), &b))
}

func TestMarshalBridgeMetaData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		latestAnswer *big.Int
		updatedAt    *big.Int
		want         map[string]any
	}{
		{"nil", nil, nil,
			map[string]any{"latestAnswer": nil, "updatedAt": nil}},
		{"zero", big.NewInt(0), big.NewInt(0),
			map[string]any{"latestAnswer": float64(0), "updatedAt": float64(0)}},
		{"one", big.NewInt(1), big.NewInt(1),
			map[string]any{"latestAnswer": float64(1), "updatedAt": float64(1)}},
		{"negative", big.NewInt(-100), big.NewInt(-10),
			map[string]any{"latestAnswer": float64(-100), "updatedAt": float64(-10)}},
		// 9223372036854775807000
		{"large", new(big.Int).Mul(big.NewInt(math.MaxInt64), big.NewInt(1000)), big.NewInt(1),
			map[string]any{"latestAnswer": float64(9.223372036854776e+21), "updatedAt": float64(1)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bridges.MarshalBridgeMetaData(tt.latestAnswer, tt.updatedAt)
			require.NoError(t, err)
			assert.Equalf(t, tt.want, got, "MarshalBridgeMetaData(%v, %v)", tt.latestAnswer, tt.updatedAt)
		})
	}
}
